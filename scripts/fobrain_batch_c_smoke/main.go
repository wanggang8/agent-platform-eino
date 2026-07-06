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
	batchCReportSchemaVersion = "eino.fobrain_batch_c_live_report.v1"
	batchCScenario            = "fobrain-batch-c"
)

// batchCReport 是 Batch C live smoke 的脱敏机器报告，只记录 StructuredResult 出口和稳定失败分类。
type batchCReport struct {
	SchemaVersion      string                   `json:"schema_version"`
	Scenario           string                   `json:"scenario"`
	Status             string                   `json:"status"`
	ProviderMode       string                   `json:"provider_mode"`
	WorkspaceID        string                   `json:"workspace_id"`
	TLSPrivateCertMode bool                     `json:"tls_private_cert_mode"`
	SampleInputs       batchCSampleInputs       `json:"sample_inputs"`
	CapabilityResults  []batchCCapabilityResult `json:"capability_results"`
	RedactionChecks    []string                 `json:"redaction_checks"`
	FailureCategory    string                   `json:"failure_category"`
	BlocksClaims       []string                 `json:"blocks_claims"`
	ReportCreatedAt    string                   `json:"report_created_at"`
}

// batchCSampleInputs 只记录可选筛选是否存在，不记录真实业务名、关键词、人员或状态值。
type batchCSampleInputs struct {
	KeywordPresent      bool `json:"keyword_present"`
	BusinessNamePresent bool `json:"business_name_present"`
	SeverityPresent     bool `json:"severity_present"`
	StatusPresent       bool `json:"status_present"`
	PersonPresent       bool `json:"person_present"`
	FieldPresent        bool `json:"field_present"`
	TimeRangePresent    bool `json:"time_range_present"`
}

// batchCCapabilityResult 记录单个 Batch C 工具的事实出口；raw provider body 不进入报告。
type batchCCapabilityResult struct {
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

type batchCSamples struct {
	keyword      string
	businessName string
	severity     string
	status       string
	person       string
	field        string
	timeRange    string
}

func main() {
	configPath := flag.String("config", "configs/eino-workbench.local.yaml", "local config file")
	outputPath := flag.String("output", "test-results/eino-workbench-fobrain-batch-c-live-report.json", "report output path")
	keyword := flag.String("keyword", "", "optional keyword sample value; redacted from report")
	businessName := flag.String("business-name", "", "optional business name sample value; redacted from report")
	severity := flag.String("severity", "", "optional severity sample value; redacted from report")
	status := flag.String("status", "", "optional status sample value; redacted from report")
	person := flag.String("person", "", "optional person sample value; redacted from report")
	field := flag.String("field", "", "optional field sample value; redacted from report")
	timeRange := flag.String("time-range", "", "optional time range sample value; redacted from report")
	flag.Parse()

	err := runBatchCLiveSmoke(*configPath, *outputPath, batchCSamples{
		keyword:      *keyword,
		businessName: *businessName,
		severity:     *severity,
		status:       *status,
		person:       *person,
		field:        *field,
		timeRange:    *timeRange,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runBatchCLiveSmoke(configPath string, outputPath string, samples batchCSamples) error {
	cfg, err := bootstrap.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	report, err := executeBatchCLiveSmoke(context.Background(), cfg, samples)
	if writeErr := writeBatchCReport(outputPath, report); writeErr != nil {
		return writeErr
	}
	if err != nil {
		return err
	}
	return nil
}

func executeBatchCLiveSmoke(ctx context.Context, cfg bootstrap.Config, samples batchCSamples) (batchCReport, error) {
	report := newBatchCReport(cfg)
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
	report.SampleInputs = batchCSampleInputs{
		KeywordPresent:      strings.TrimSpace(samples.keyword) != "",
		BusinessNamePresent: strings.TrimSpace(samples.businessName) != "",
		SeverityPresent:     strings.TrimSpace(samples.severity) != "",
		StatusPresent:       strings.TrimSpace(samples.status) != "",
		PersonPresent:       strings.TrimSpace(samples.person) != "",
		FieldPresent:        strings.TrimSpace(samples.field) != "",
		TimeRangePresent:    strings.TrimSpace(samples.timeRange) != "",
	}

	policyContext := fobrainPolicyContextFromConfig(cfg.Fobrain)
	gate := product.NewStructuredResultSafetyGate()
	var firstErr error
	for _, request := range batchCRequests(samples) {
		item := batchCCapabilityResult{
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
			item.FailureCategory = stableFailureCategory(err)
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
	report.Status = batchCOverallStatus(report.CapabilityResults)
	report.FailureCategory = batchCFailureCategory(report.CapabilityResults, firstErr)
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

type batchCInvocationRequest struct {
	capabilityID string
	arguments    map[string]any
	sampleSource string
}

func batchCRequests(samples batchCSamples) []batchCInvocationRequest {
	return []batchCInvocationRequest{
		batchCRequest(fobrain.CapabilityBusinessList, map[string]any{
			"business_name": samples.businessName,
			"keyword":       samples.keyword,
			"page":          1,
			"page_size":     20,
		}, samples.businessName, samples.keyword),
		batchCRequest(fobrain.CapabilityExternalHighRiskAssets, map[string]any{
			"severity":  samples.severity,
			"keyword":   samples.keyword,
			"page":      1,
			"page_size": 20,
		}, samples.severity, samples.keyword),
		batchCRequest(fobrain.CapabilityVulnerabilityStatusSummary, map[string]any{
			"severity": samples.severity,
			"status":   samples.status,
			"keyword":  samples.keyword,
		}, samples.severity, samples.status, samples.keyword),
		batchCRequest(fobrain.CapabilityPendingTickets, map[string]any{
			"person":    samples.person,
			"status":    samples.status,
			"keyword":   samples.keyword,
			"page":      1,
			"page_size": 20,
		}, samples.person, samples.status, samples.keyword),
		batchCRequest(fobrain.CapabilityIPStats, map[string]any{
			"field":      samples.field,
			"time_range": samples.timeRange,
			"keyword":    samples.keyword,
		}, samples.field, samples.timeRange, samples.keyword),
		batchCRequest(fobrain.CapabilityVulStats, map[string]any{
			"field":      samples.field,
			"time_range": samples.timeRange,
			"keyword":    samples.keyword,
		}, samples.field, samples.timeRange, samples.keyword),
	}
}

func batchCRequest(capabilityID string, arguments map[string]any, sampleValues ...string) batchCInvocationRequest {
	source := "default_filter"
	for _, value := range sampleValues {
		if strings.TrimSpace(value) != "" {
			source = "provided"
			break
		}
	}
	return batchCInvocationRequest{capabilityID: capabilityID, arguments: arguments, sampleSource: source}
}

func batchCOverallStatus(results []batchCCapabilityResult) string {
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

func batchCFailureCategory(results []batchCCapabilityResult, firstErr error) string {
	status := batchCOverallStatus(results)
	if status == "passed" {
		return "none"
	}
	if status == "blocked" {
		return "empty_result"
	}
	return stableFailureCategory(firstErr)
}

func newBatchCReport(cfg bootstrap.Config) batchCReport {
	return batchCReport{
		SchemaVersion:      batchCReportSchemaVersion,
		Scenario:           batchCScenario,
		Status:             "failed",
		ProviderMode:       cfg.Fobrain.ConnectorStatus.Mode,
		WorkspaceID:        cfg.Fobrain.WorkspaceID,
		TLSPrivateCertMode: cfg.Fobrain.TLS.InsecureSkipVerify,
		CapabilityResults:  []batchCCapabilityResult{},
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
		BlocksClaims:    []string{"fobrain-batch-c live pass", "Fobrain 24 readonly final acceptance"},
		ReportCreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func writeBatchCReport(outputPath string, report batchCReport) error {
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

// fobrainProviderFromConfig 复用服务启动路径的 provider 边界规则，但只服务 Batch C smoke。
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
