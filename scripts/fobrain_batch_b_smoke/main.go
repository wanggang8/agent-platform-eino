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
	batchBReportSchemaVersion = "eino.fobrain_batch_b_live_report.v1"
	batchBScenario            = "fobrain-batch-b"
)

// batchBReport 是 Batch B live smoke 的脱敏机器报告，只记录 Product Facts 出口和稳定失败分类。
type batchBReport struct {
	SchemaVersion      string                   `json:"schema_version"`
	Scenario           string                   `json:"scenario"`
	Status             string                   `json:"status"`
	ProviderMode       string                   `json:"provider_mode"`
	WorkspaceID        string                   `json:"workspace_id"`
	TLSPrivateCertMode bool                     `json:"tls_private_cert_mode"`
	SampleInputs       batchBSampleInputs       `json:"sample_inputs"`
	CapabilityResults  []batchBCapabilityResult `json:"capability_results"`
	RedactionChecks    []string                 `json:"redaction_checks"`
	FailureCategory    string                   `json:"failure_category"`
	BlocksClaims       []string                 `json:"blocks_claims"`
	ReportCreatedAt    string                   `json:"report_created_at"`
}

// batchBSampleInputs 只记录外部筛选是否存在；当前用户/部门由 provider 内部安全上下文派生。
type batchBSampleInputs struct {
	KeywordPresent bool `json:"keyword_present"`
}

// batchBCapabilityResult 记录单个 Batch B 工具的 StructuredResult 出口；不记录当前用户、部门、token 或 raw body。
type batchBCapabilityResult struct {
	CapabilityID           string `json:"capability_id"`
	Status                 string `json:"status"`
	ResultState            string `json:"result_state"`
	StructuredResultSchema string `json:"structured_result_schema"`
	ResultRef              string `json:"result_ref,omitempty"`
	ItemCount              int    `json:"item_count"`
	PolicyDecision         string `json:"policy_decision"`
	FailureCategory        string `json:"failure_category"`
	SampleSource           string `json:"sample_source"`
}

type batchBSamples struct {
	keyword string
}

func main() {
	configPath := flag.String("config", "configs/eino-workbench.local.yaml", "local config file")
	outputPath := flag.String("output", "test-results/eino-workbench-fobrain-batch-b-live-report.json", "report output path")
	keyword := flag.String("keyword", "", "optional keyword sample value; redacted from report")
	flag.Parse()

	if err := runBatchBLiveSmoke(*configPath, *outputPath, batchBSamples{keyword: *keyword}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runBatchBLiveSmoke(configPath string, outputPath string, samples batchBSamples) error {
	cfg, err := bootstrap.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	report, err := executeBatchBLiveSmoke(context.Background(), cfg, samples)
	if writeErr := writeBatchBReport(outputPath, report); writeErr != nil {
		return writeErr
	}
	if err != nil {
		return err
	}
	return nil
}

func executeBatchBLiveSmoke(ctx context.Context, cfg bootstrap.Config, samples batchBSamples) (batchBReport, error) {
	report := newBatchBReport(cfg)
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
	report.SampleInputs = batchBSampleInputs{KeywordPresent: strings.TrimSpace(samples.keyword) != ""}

	policyContext := fobrainPolicyContextFromConfig(cfg.Fobrain)
	gate := product.NewStructuredResultSafetyGate()
	var firstErr error
	for _, request := range batchBRequests(samples) {
		item := batchBCapabilityResult{
			CapabilityID:           request.capabilityID,
			Status:                 "failed",
			ResultState:            "not_run",
			StructuredResultSchema: facts.StructuredResultSchemaVersion,
			PolicyDecision:         "allowed",
			FailureCategory:        "none",
			SampleSource:           request.sampleSource,
		}
		candidate, err := provider.Invoke(ctx, capabilities.InvocationRequest{
			CapabilityID:  request.capabilityID,
			Arguments:     request.arguments,
			PolicyContext: policyContext,
		})
		if err != nil {
			item.FailureCategory = batchBStableFailureCategory(err)
			if item.FailureCategory == string(fobrain.PolicyReasonMissingCurrentUserScope) {
				item.Status = "blocked"
				report.CapabilityResults = append(report.CapabilityResults, item)
				continue
			}
			if firstErr == nil {
				firstErr = fmt.Errorf("capability %s failed: %w", request.capabilityID, err)
			}
			report.CapabilityResults = append(report.CapabilityResults, item)
			continue
		}
		item.ItemCount = candidate.ItemCount
		item.ResultState = "resolved"
		if candidate.ItemCount == 0 {
			item.Status = "blocked"
			item.ResultState = "empty"
			item.FailureCategory = "empty_result"
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
	report.Status = batchBOverallStatus(report.CapabilityResults)
	report.FailureCategory = batchBFailureCategory(report.CapabilityResults, firstErr)
	if report.Status == "passed" {
		report.BlocksClaims = []string{}
	}
	if firstErr != nil {
		return report, firstErr
	}
	if report.Status != "passed" {
		return report, errors.New(report.FailureCategory)
	}
	return report, nil
}

type batchBInvocationRequest struct {
	capabilityID string
	arguments    map[string]any
	sampleSource string
}

func batchBRequests(samples batchBSamples) []batchBInvocationRequest {
	return []batchBInvocationRequest{
		batchBRequest(fobrain.CapabilityMyAssets, samples.keyword),
		batchBRequest(fobrain.CapabilityMyDepartmentAssets, samples.keyword),
		batchBRequest(fobrain.CapabilityMyVulnerabilities, samples.keyword),
		batchBRequest(fobrain.CapabilityMyDepartmentVulnerabilities, samples.keyword),
		batchBRequest(fobrain.CapabilityMyBusinessSystems, samples.keyword),
		batchBRequest(fobrain.CapabilityMyImportantBusinessSystems, samples.keyword),
	}
}

func batchBRequest(capabilityID string, keyword string) batchBInvocationRequest {
	keyword = strings.TrimSpace(keyword)
	source := "current_user"
	arguments := map[string]any{"page": 1, "page_size": 20}
	if keyword != "" {
		arguments["keyword"] = keyword
		source = "provided_filter"
	}
	return batchBInvocationRequest{capabilityID: capabilityID, arguments: arguments, sampleSource: source}
}

func batchBOverallStatus(results []batchBCapabilityResult) string {
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

func batchBFailureCategory(results []batchBCapabilityResult, firstErr error) string {
	status := batchBOverallStatus(results)
	if status == "passed" {
		return "none"
	}
	if status == "blocked" {
		for _, result := range results {
			if result.FailureCategory == string(fobrain.PolicyReasonMissingCurrentUserScope) {
				return string(fobrain.PolicyReasonMissingCurrentUserScope)
			}
		}
		return "empty_result"
	}
	return stableFailureCategory(firstErr)
}

func batchBStableFailureCategory(err error) string {
	if err == nil {
		return "none"
	}
	return stableFailureCategory(err)
}

func newBatchBReport(cfg bootstrap.Config) batchBReport {
	return batchBReport{
		SchemaVersion:      batchBReportSchemaVersion,
		Scenario:           batchBScenario,
		Status:             "failed",
		ProviderMode:       cfg.Fobrain.ConnectorStatus.Mode,
		WorkspaceID:        cfg.Fobrain.WorkspaceID,
		TLSPrivateCertMode: cfg.Fobrain.TLS.InsecureSkipVerify,
		CapabilityResults:  []batchBCapabilityResult{},
		RedactionChecks: []string{
			"secret_absent",
			"auth_header_absent",
			"raw_payload_absent",
			"sample_values_absent",
			"binding_reference_absent",
			"checkpoint_absent",
			"interrupt_absent",
		},
		FailureCategory: "not_run",
		BlocksClaims: []string{
			"fobrain-batch-b live pass",
			"Fobrain 24 readonly final acceptance",
		},
		ReportCreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func writeBatchBReport(outputPath string, report batchBReport) error {
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

// fobrainProviderFromConfig 复用服务启动路径的 provider 边界规则，但只服务 Batch B smoke。
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
	return "connector_execution_failed"
}
