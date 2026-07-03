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
	batchEReportSchemaVersion = "eino.fobrain_batch_e_live_report.v1"
	batchEScenario            = "fobrain-batch-e"
)

// batchEReport 是 Batch E live smoke 的脱敏机器报告，不记录真实样本值、token 或 raw payload。
type batchEReport struct {
	SchemaVersion      string                   `json:"schema_version"`
	Scenario           string                   `json:"scenario"`
	Status             string                   `json:"status"`
	ProviderMode       string                   `json:"provider_mode"`
	WorkspaceID        string                   `json:"workspace_id"`
	TLSPrivateCertMode bool                     `json:"tls_private_cert_mode"`
	SampleInputs       batchESampleInputs       `json:"sample_inputs"`
	CapabilityResults  []batchECapabilityResult `json:"capability_results"`
	RedactionChecks    []string                 `json:"redaction_checks"`
	FailureCategory    string                   `json:"failure_category"`
	BlocksClaims       []string                 `json:"blocks_claims"`
	ReportCreatedAt    string                   `json:"report_created_at"`
}

// batchESampleInputs 只记录 Batch E 样本是否存在，不记录真实 ID、业务名或漏洞名。
type batchESampleInputs struct {
	AssetDetailPresent         bool `json:"asset_detail_present"`
	VulnerabilityDetailPresent bool `json:"vulnerability_detail_present"`
	BusinessPresent            bool `json:"business_present"`
	ThreatNamePresent          bool `json:"threat_name_present"`
}

// batchECapabilityResult 记录单个 Batch E 工具的 StructuredResult 出口和稳定失败分类。
type batchECapabilityResult struct {
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

type batchESamples struct {
	assetID           string
	assetNetworkType  string
	vulnerabilityID   string
	businessName      string
	vulnerabilityName string
}

func main() {
	configPath := flag.String("config", "configs/eino-workbench.local.yaml", "local config file")
	outputPath := flag.String("output", "test-results/eino-workbench-fobrain-batch-e-live-report.json", "report output path")
	assetID := flag.String("asset-id", "", "asset id sample value; redacted from report")
	assetNetworkType := flag.String("asset-network-type", "internal", "asset network type sample value; redacted from report")
	vulnerabilityID := flag.String("vulnerability-id", "", "vulnerability id sample value; redacted from report")
	businessName := flag.String("business-name", "", "business name sample value; redacted from report")
	vulnerabilityName := flag.String("vulnerability-name", "", "vulnerability name sample value; redacted from report")
	flag.Parse()

	err := runBatchELiveSmoke(*configPath, *outputPath, batchESamples{
		assetID:           *assetID,
		assetNetworkType:  *assetNetworkType,
		vulnerabilityID:   *vulnerabilityID,
		businessName:      *businessName,
		vulnerabilityName: *vulnerabilityName,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runBatchELiveSmoke(configPath string, outputPath string, samples batchESamples) error {
	cfg, err := bootstrap.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	report, err := executeBatchELiveSmoke(context.Background(), cfg, samples)
	if writeErr := writeBatchEReport(outputPath, report); writeErr != nil {
		return writeErr
	}
	if err != nil {
		return err
	}
	return nil
}

func executeBatchELiveSmoke(ctx context.Context, cfg bootstrap.Config, samples batchESamples) (batchEReport, error) {
	report := newBatchEReport(cfg)
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
	report.SampleInputs = batchESampleInputs{
		AssetDetailPresent:         strings.TrimSpace(samples.assetID) != "",
		VulnerabilityDetailPresent: strings.TrimSpace(samples.vulnerabilityID) != "",
		BusinessPresent:            strings.TrimSpace(samples.businessName) != "",
		ThreatNamePresent:          strings.TrimSpace(samples.vulnerabilityName) != "",
	}
	policyContext := fobrainPolicyContextFromConfig(cfg.Fobrain)
	gate := product.NewStructuredResultSafetyGate()
	var firstErr error
	for _, request := range batchERequests(samples) {
		item := batchECapabilityResult{
			CapabilityID:           request.capabilityID,
			Status:                 "failed",
			ResultState:            "not_run",
			StructuredResultSchema: facts.StructuredResultSchemaVersion,
			PolicyDecision:         "allowed",
			FailureCategory:        "none",
			SampleSource:           request.sampleSource,
		}
		if request.missingSample {
			item.Status = "blocked"
			item.ResultState = "not_run"
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
		item.ItemCount = candidate.ItemCount
		item.ResultState = "resolved"
		if candidate.ItemCount == 0 {
			item.ResultState = "empty"
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
	report.Status = batchEOverallStatus(report.CapabilityResults)
	report.FailureCategory = batchEFailureCategory(report.CapabilityResults, firstErr)
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

type batchEInvocationRequest struct {
	capabilityID  string
	arguments     map[string]any
	sampleSource  string
	missingSample bool
}

func batchERequests(samples batchESamples) []batchEInvocationRequest {
	return []batchEInvocationRequest{
		batchERequest(fobrain.CapabilityGetAssetDetail, map[string]any{"asset_id": samples.assetID, "network_type": samples.assetNetworkType}, samples.assetID),
		batchERequest(fobrain.CapabilityGetVulnerabilityDetail, map[string]any{"vulnerability_id": samples.vulnerabilityID}, samples.vulnerabilityID),
		batchERequest(fobrain.CapabilityBusinessRiskSummary, map[string]any{"business_name": samples.businessName}, samples.businessName),
		batchERequest(fobrain.CapabilityThreatRelevanceList, map[string]any{"vulnerability_name": samples.vulnerabilityName, "page": 1, "page_size": 20}, samples.vulnerabilityName),
	}
}

func batchERequest(capabilityID string, arguments map[string]any, requiredValue string) batchEInvocationRequest {
	if strings.TrimSpace(requiredValue) == "" {
		return batchEInvocationRequest{capabilityID: capabilityID, arguments: map[string]any{}, sampleSource: "missing", missingSample: true}
	}
	if capabilityID == fobrain.CapabilityGetAssetDetail && strings.TrimSpace(fmt.Sprint(arguments["network_type"])) == "" {
		arguments["network_type"] = "internal"
	}
	return batchEInvocationRequest{capabilityID: capabilityID, arguments: arguments, sampleSource: "provided"}
}

func batchEOverallStatus(results []batchECapabilityResult) string {
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

func batchEFailureCategory(results []batchECapabilityResult, firstErr error) string {
	status := batchEOverallStatus(results)
	if status == "passed" {
		return "none"
	}
	if status == "blocked" {
		return "missing_sample_input"
	}
	return stableFailureCategory(firstErr)
}

func newBatchEReport(cfg bootstrap.Config) batchEReport {
	return batchEReport{
		SchemaVersion:      batchEReportSchemaVersion,
		Scenario:           batchEScenario,
		Status:             "failed",
		ProviderMode:       cfg.Fobrain.ConnectorStatus.Mode,
		WorkspaceID:        cfg.Fobrain.WorkspaceID,
		TLSPrivateCertMode: cfg.Fobrain.TLS.InsecureSkipVerify,
		CapabilityResults:  []batchECapabilityResult{},
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
		BlocksClaims:    []string{"fobrain-batch-e live pass", "Fobrain 24 readonly final acceptance"},
		ReportCreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func writeBatchEReport(outputPath string, report batchEReport) error {
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

// fobrainProviderFromConfig 复用服务启动路径的 provider 边界规则，但只服务 Batch E smoke。
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
		WorkspaceID: config.WorkspaceID,
		System:      config.System,
		Status:      capabilities.CredentialStatus(config.Status),
		DisplayRef:  config.DisplayRef,
		OwnerScope:  capabilities.PermissionScope(config.OwnerScope),
		AuditRef:    config.AuditRef,
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
