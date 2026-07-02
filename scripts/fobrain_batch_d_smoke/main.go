package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/bootstrap"
	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

const (
	batchDReportSchemaVersion = "eino.fobrain_batch_d_live_report.v1"
	batchDScenario            = "fobrain-batch-d"
)

// batchDReport 是 Batch D live smoke 的脱敏机器报告，不记录样本值、token 或 raw payload。
type batchDReport struct {
	SchemaVersion      string                   `json:"schema_version"`
	Scenario           string                   `json:"scenario"`
	Status             string                   `json:"status"`
	ProviderMode       string                   `json:"provider_mode"`
	WorkspaceID        string                   `json:"workspace_id"`
	TLSPrivateCertMode bool                     `json:"tls_private_cert_mode"`
	SampleInputs       batchDSampleInputs       `json:"sample_inputs"`
	CapabilityResults  []batchDCapabilityResult `json:"capability_results"`
	RedactionChecks    []string                 `json:"redaction_checks"`
	FailureCategory    string                   `json:"failure_category"`
	BlocksClaims       []string                 `json:"blocks_claims"`
	ReportCreatedAt    string                   `json:"report_created_at"`
}

// batchDSampleInputs 只记录样本是否存在，不能记录负责人、部门、IP 的真实值。
type batchDSampleInputs struct {
	OwnerPresent      bool `json:"owner_present"`
	DepartmentPresent bool `json:"department_present"`
	IPPresent         bool `json:"ip_present"`
}

// batchDCapabilityResult 记录单个工具的 StructuredResult 出口和稳定失败分类。
type batchDCapabilityResult struct {
	CapabilityID           string `json:"capability_id"`
	Status                 string `json:"status"`
	StructuredResultSchema string `json:"structured_result_schema"`
	ResultRef              string `json:"result_ref,omitempty"`
	PolicyDecision         string `json:"policy_decision"`
	FailureCategory        string `json:"failure_category"`
	SampleSource           string `json:"sample_source"`
}

type batchDSamples struct {
	owner            string
	department       string
	ip               string
	ownerSource      string
	departmentSource string
	ipSource         string
}

func main() {
	configPath := flag.String("config", "configs/eino-workbench.local.yaml", "local config file")
	outputPath := flag.String("output", "test-results/eino-workbench-fobrain-batch-d-live-report.json", "report output path")
	owner := flag.String("owner", "", "owner sample value; redacted from report")
	department := flag.String("department", "", "department sample value; redacted from report")
	ip := flag.String("ip", "", "ip sample value; redacted from report")
	flag.Parse()

	if err := runBatchDLiveSmoke(*configPath, *outputPath, batchDSamples{owner: *owner, department: *department, ip: *ip}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runBatchDLiveSmoke(configPath string, outputPath string, samples batchDSamples) error {
	cfg, err := bootstrap.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	report, err := executeBatchDLiveSmoke(context.Background(), cfg, samples)
	if writeErr := writeBatchDReport(outputPath, report); writeErr != nil {
		return writeErr
	}
	if err != nil {
		return err
	}
	return nil
}

func executeBatchDLiveSmoke(ctx context.Context, cfg bootstrap.Config, samples batchDSamples) (batchDReport, error) {
	report := newBatchDReport(cfg)
	if !cfg.Fobrain.Enabled || cfg.Fobrain.ConnectorStatus.Mode != "live" {
		report.Status = "failed"
		report.FailureCategory = "live_config_unavailable"
		return report, errors.New("fobrain live config unavailable")
	}
	provider, err := fobrainProviderFromConfig(cfg.Fobrain)
	if err != nil {
		report.Status = "failed"
		report.FailureCategory = stableFailureCategory(err)
		return report, err
	}
	samples = enrichSamplesFromCurrentUser(ctx, provider, cfg.Fobrain, samples)
	report.SampleInputs = batchDSampleInputs{
		OwnerPresent:      strings.TrimSpace(samples.owner) != "",
		DepartmentPresent: strings.TrimSpace(samples.department) != "",
		IPPresent:         strings.TrimSpace(samples.ip) != "",
	}

	policyContext := fobrainPolicyContextFromConfig(cfg.Fobrain)
	gate := product.NewStructuredResultSafetyGate()
	var firstErr error
	for _, request := range batchDRequests(samples) {
		item := batchDCapabilityResult{
			CapabilityID:           request.capabilityID,
			Status:                 "failed",
			StructuredResultSchema: facts.StructuredResultSchemaVersion,
			PolicyDecision:         "allowed",
			FailureCategory:        "none",
			SampleSource:           request.sampleSource,
		}
		if request.missingSample {
			item.Status = "blocked"
			item.FailureCategory = "missing_sample_input"
			item.PolicyDecision = "permission_denied"
			report.CapabilityResults = append(report.CapabilityResults, item)
			continue
		}
		candidate, err := provider.Invoke(ctx, capabilities.InvocationRequest{
			CapabilityID:  request.capabilityID,
			Arguments:     request.arguments,
			PolicyContext: policyContext,
		})
		if err != nil {
			item.FailureCategory = stableFailureCategory(err)
			if firstErr == nil {
				firstErr = fmt.Errorf("capability %s failed: %w", request.capabilityID, err)
			}
			report.CapabilityResults = append(report.CapabilityResults, item)
			continue
		}
		structuredResult, err := gate.Approve(candidate)
		if err != nil {
			item.FailureCategory = "structured_result_rejected"
			if firstErr == nil {
				firstErr = fmt.Errorf("capability %s structured result rejected: %w", request.capabilityID, err)
			}
			report.CapabilityResults = append(report.CapabilityResults, item)
			continue
		}
		item.Status = "passed"
		item.ResultRef = structuredResult.ResultRef
		report.CapabilityResults = append(report.CapabilityResults, item)
	}
	report.Status = batchDOverallStatus(report.CapabilityResults)
	report.FailureCategory = batchDFailureCategory(report.Status, firstErr)
	if report.Status == "passed" {
		report.BlocksClaims = []string{}
	}
	if firstErr != nil {
		return report, firstErr
	}
	return report, nil
}

func enrichSamplesFromCurrentUser(ctx context.Context, provider *fobrain.Provider, config bootstrap.FobrainConfig, samples batchDSamples) batchDSamples {
	if strings.TrimSpace(samples.owner) != "" && strings.TrimSpace(samples.department) != "" {
		return samples
	}
	candidate, err := provider.Invoke(ctx, capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityCurrentUserContext,
		PolicyContext: fobrainPolicyContextFromConfig(config),
	})
	if err != nil {
		return samples
	}
	// current_user StructuredResult 只用于判断 owner/department 是否可从当前用户安全摘要补齐；不写入 Batch D 报告。
	summary := candidate.SafeSummary
	if strings.TrimSpace(samples.owner) == "" {
		samples.owner = valueAfterLabel(summary, "Fobrain 当前用户：")
		samples.ownerSource = "current_user"
	}
	if strings.TrimSpace(samples.department) == "" {
		samples.department = valueAfterLabel(summary, "部门：")
		samples.departmentSource = "current_user"
	}
	return samples
}

func valueAfterLabel(summary string, label string) string {
	index := strings.Index(summary, label)
	if index < 0 {
		return ""
	}
	value := summary[index+len(label):]
	if comma := strings.Index(value, "，"); comma >= 0 {
		value = value[:comma]
	}
	return strings.TrimSpace(value)
}

type batchDInvocationRequest struct {
	capabilityID  string
	arguments     map[string]any
	sampleSource  string
	missingSample bool
}

func batchDRequests(samples batchDSamples) []batchDInvocationRequest {
	ownerSource := sourceOrProvided(samples.ownerSource)
	departmentSource := sourceOrProvided(samples.departmentSource)
	ipSource := sourceOrProvided(samples.ipSource)
	return []batchDInvocationRequest{
		batchDRequest(fobrain.CapabilityListAssetsByOwner, "person_name", samples.owner, ownerSource),
		batchDRequest(fobrain.CapabilityListVulnerabilitiesByOwner, "person_name", samples.owner, ownerSource),
		batchDRequest(fobrain.CapabilityListAssetsByDepartment, "department_name", samples.department, departmentSource),
		batchDRequest(fobrain.CapabilityListVulnerabilitiesByDepartment, "department_name", samples.department, departmentSource),
		batchDRequest(fobrain.CapabilityListAssetsByIP, "ip", samples.ip, ipSource),
		batchDRequest(fobrain.CapabilityListVulnerabilitiesByIP, "ip", samples.ip, ipSource),
	}
}

func batchDRequest(capabilityID string, key string, value string, source string) batchDInvocationRequest {
	value = strings.TrimSpace(value)
	if value == "" {
		return batchDInvocationRequest{capabilityID: capabilityID, arguments: map[string]any{}, sampleSource: "missing", missingSample: true}
	}
	return batchDInvocationRequest{
		capabilityID: capabilityID,
		arguments: map[string]any{
			key:         value,
			"page":      1,
			"page_size": 20,
		},
		sampleSource: source,
	}
}

func sourceOrProvided(source string) string {
	if strings.TrimSpace(source) == "" {
		return "provided"
	}
	return source
}

func batchDOverallStatus(results []batchDCapabilityResult) string {
	hasFailed := false
	hasBlocked := false
	for _, result := range results {
		switch result.Status {
		case "failed":
			hasFailed = true
		case "blocked":
			hasBlocked = true
		}
	}
	if hasFailed {
		return "failed"
	}
	if hasBlocked {
		return "blocked"
	}
	return "passed"
}

func batchDFailureCategory(status string, firstErr error) string {
	if status == "passed" {
		return "none"
	}
	if status == "blocked" {
		return "missing_sample_input"
	}
	return stableFailureCategory(firstErr)
}

func newBatchDReport(cfg bootstrap.Config) batchDReport {
	return batchDReport{
		SchemaVersion:      batchDReportSchemaVersion,
		Scenario:           batchDScenario,
		Status:             "failed",
		ProviderMode:       cfg.Fobrain.ConnectorStatus.Mode,
		WorkspaceID:        cfg.Fobrain.WorkspaceID,
		TLSPrivateCertMode: cfg.Fobrain.TLS.InsecureSkipVerify,
		CapabilityResults:  []batchDCapabilityResult{},
		RedactionChecks: []string{
			"secret_absent",
			"auth_header_absent",
			"raw_payload_absent",
			"binding_reference_absent",
			"checkpoint_absent",
			"interrupt_absent",
		},
		FailureCategory: "not_run",
		BlocksClaims: []string{
			"fobrain-batch-d live pass",
			"Fobrain 24 readonly final acceptance",
		},
		ReportCreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func writeBatchDReport(outputPath string, report batchDReport) error {
	if outputPath == "" {
		return errors.New("output path is required")
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(outputPath, data, 0o600)
}

// fobrainProviderFromConfig 复用服务启动路径的 provider 边界规则，但只服务 Batch D smoke。
func fobrainProviderFromConfig(config bootstrap.FobrainConfig) (*fobrain.Provider, error) {
	status := capabilities.ConnectorStatusUnavailable
	if config.ConnectorStatus.Available {
		status = capabilities.ConnectorStatusAvailable
	}
	client := fobrain.FobrainClient(fobrain.MockClient{})
	if config.ConnectorStatus.Mode == "live" {
		httpClient, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
			BaseURL:            config.BaseURL,
			Timeout:            config.Timeout,
			InsecureSkipVerify: config.TLS.InsecureSkipVerify,
		})
		if err != nil {
			return nil, err
		}
		client = httpClient
	}
	return fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        config.WorkspaceID,
		CredentialBinding:  fobrainCredentialBindingFromConfig(config.CredentialBinding),
		ConnectorStatus:    status,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: config.WorkspaceID, AuthParam: config.Credential.AuthParam, APIToken: config.Credential.APIToken},
		Client:             client,
	}), nil
}

// fobrainPolicyContextFromConfig 只传递可展示的 policy 材料，不传真实 token。
func fobrainPolicyContextFromConfig(config bootstrap.FobrainConfig) capabilities.PolicyContext {
	status := capabilities.ConnectorStatusUnavailable
	if config.ConnectorStatus.Available {
		status = capabilities.ConnectorStatusAvailable
	}
	return capabilities.PolicyContext{
		WorkspaceID:       config.WorkspaceID,
		CredentialBinding: fobrainCredentialBindingFromConfig(config.CredentialBinding),
		ConnectorStatus:   status,
	}
}

func fobrainCredentialBindingFromConfig(config bootstrap.CredentialBinding) capabilities.CredentialBinding {
	return capabilities.CredentialBinding{
		SchemaVersion: config.SchemaVersion,
		WorkspaceID:   config.WorkspaceID,
		System:        config.System,
		Status:        capabilities.CredentialStatus(config.Status),
		DisplayRef:    config.DisplayRef,
		OwnerScope:    capabilities.PermissionScope(config.OwnerScope),
		UpdatedAt:     config.UpdatedAt,
		AuditRef:      config.AuditRef,
	}
}

func stableFailureCategory(err error) string {
	if err == nil {
		return "none"
	}
	var stable interface{ StableReasonCode() string }
	if errors.As(err, &stable) {
		return stable.StableReasonCode()
	}
	return "unexpected_error"
}
