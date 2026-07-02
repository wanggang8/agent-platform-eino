package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"agent-platform-eino/internal/einoapp/bootstrap"
	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

const (
	batchAReportSchemaVersion = "eino.fobrain_batch_a_live_report.v1"
	batchAScenario            = "fobrain-batch-a"
)

// batchAReport 是 Batch A live smoke 的脱敏验收报告。
// 报告只记录 StructuredResult 安全摘要和稳定原因码，不记录 raw 响应或凭据。
type batchAReport struct {
	SchemaVersion      string                   `json:"schema_version"`
	Scenario           string                   `json:"scenario"`
	Status             string                   `json:"status"`
	ProviderMode       string                   `json:"provider_mode"`
	WorkspaceID        string                   `json:"workspace_id"`
	TLSPrivateCertMode bool                     `json:"tls_private_cert_mode"`
	CapabilityResults  []batchACapabilityReport `json:"capability_results"`
	RedactionChecks    []string                 `json:"redaction_checks"`
	FailureCategory    string                   `json:"failure_category"`
	BlocksClaims       []string                 `json:"blocks_claims,omitempty"`
	ReportCreatedAt    string                   `json:"report_created_at"`
}

// batchACapabilityReport 记录单个 capability 的事实出口，不包含 provider raw payload。
type batchACapabilityReport struct {
	CapabilityID           string `json:"capability_id"`
	Status                 string `json:"status"`
	StructuredResultSchema string `json:"structured_result_schema"`
	ResultRef              string `json:"result_ref,omitempty"`
	SafeSummary            string `json:"safe_summary,omitempty"`
	PolicyDecision         string `json:"policy_decision"`
	FailureCategory        string `json:"failure_category"`
}

func main() {
	configPath := flag.String("config", "configs/eino-workbench.local.yaml", "local config file")
	outputPath := flag.String("output", "test-results/eino-workbench-fobrain-batch-a-live-report.json", "report output path")
	flag.Parse()

	if err := runBatchALiveSmoke(*configPath, *outputPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runBatchALiveSmoke 从本地配置执行 Batch A provider live 验收并写入报告。
func runBatchALiveSmoke(configPath string, outputPath string) error {
	cfg, err := bootstrap.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	report, err := executeBatchALiveSmoke(context.Background(), cfg)
	if writeErr := writeBatchAReport(outputPath, report); writeErr != nil {
		return writeErr
	}
	if err != nil {
		return err
	}
	return nil
}

// executeBatchALiveSmoke 通过 provider catalog 调用当前阶段已注册的 Batch A 能力。
func executeBatchALiveSmoke(ctx context.Context, cfg bootstrap.Config) (batchAReport, error) {
	report := newBatchAReport(cfg)
	if !cfg.Fobrain.Enabled || cfg.Fobrain.ConnectorStatus.Mode != "live" {
		report.Status = "failed"
		report.FailureCategory = "live_config_unavailable"
		report.BlocksClaims = []string{"fobrain-batch-a live pass"}
		return report, errors.New("fobrain live config unavailable")
	}
	provider, err := fobrainProviderFromConfig(cfg.Fobrain)
	if err != nil {
		report.Status = "failed"
		report.FailureCategory = stableFailureCategory(err)
		report.BlocksClaims = []string{"fobrain-batch-a live pass"}
		return report, err
	}
	catalog, err := provider.ListCapabilities()
	if err != nil {
		report.Status = "failed"
		report.FailureCategory = "catalog_unavailable"
		report.BlocksClaims = []string{"fobrain-batch-a live pass"}
		return report, err
	}
	batchACatalog := batchACapabilities(catalog)
	if len(batchACatalog) != 3 {
		report.Status = "failed"
		report.FailureCategory = "batch_a_catalog_mismatch"
		report.BlocksClaims = []string{"fobrain-batch-a live pass"}
		return report, fmt.Errorf("batch a catalog length = %d", len(batchACatalog))
	}

	var firstErr error
	policyContext := fobrainPolicyContextFromConfig(cfg.Fobrain)
	gate := product.NewStructuredResultSafetyGate()
	for _, capability := range batchACatalog {
		item := batchACapabilityReport{
			CapabilityID:           capability.ID,
			Status:                 "failed",
			StructuredResultSchema: facts.StructuredResultSchemaVersion,
			PolicyDecision:         "allowed",
			FailureCategory:        "none",
		}
		decision := capabilities.EvaluatePolicy(capability, policyContext)
		if !decision.Allowed {
			item.PolicyDecision = string(decision.ReasonCode)
			item.FailureCategory = string(decision.ReasonCode)
			if firstErr == nil {
				firstErr = fmt.Errorf("capability %s blocked: %s", capability.ID, decision.ReasonCode)
			}
			report.CapabilityResults = append(report.CapabilityResults, item)
			continue
		}
		candidate, err := provider.Invoke(ctx, capabilities.InvocationRequest{
			CapabilityID:  capability.ID,
			PolicyContext: policyContext,
		})
		if err != nil {
			item.FailureCategory = stableFailureCategory(err)
			if firstErr == nil {
				firstErr = fmt.Errorf("capability %s failed: %w", capability.ID, err)
			}
			report.CapabilityResults = append(report.CapabilityResults, item)
			continue
		}
		structuredResult, err := gate.Approve(candidate)
		if err != nil {
			item.FailureCategory = "structured_result_rejected"
			if firstErr == nil {
				firstErr = fmt.Errorf("capability %s structured result rejected: %w", capability.ID, err)
			}
			report.CapabilityResults = append(report.CapabilityResults, item)
			continue
		}
		item.Status = "passed"
		item.ResultRef = structuredResult.ResultRef
		item.SafeSummary = structuredResult.SafeSummary
		report.CapabilityResults = append(report.CapabilityResults, item)
	}
	if firstErr != nil {
		report.Status = "failed"
		report.FailureCategory = "capability_failed"
		report.BlocksClaims = []string{"fobrain-batch-a live pass"}
		return report, firstErr
	}
	report.Status = "passed"
	report.FailureCategory = "none"
	return report, nil
}

// batchACapabilities 只保留 Batch A 验收范围，避免后续 catalog 扩展误伤 live smoke。
func batchACapabilities(catalog []capabilities.Capability) []capabilities.Capability {
	out := []capabilities.Capability{}
	for _, capability := range catalog {
		switch capability.ID {
		case fobrain.CapabilityConnectorSecurity, fobrain.CapabilityCurrentUserContext, fobrain.CapabilityMyPermissions:
			out = append(out, capability)
		}
	}
	return out
}

// newBatchAReport 初始化报告公共字段，避免失败分支缺少门禁字段。
func newBatchAReport(cfg bootstrap.Config) batchAReport {
	return batchAReport{
		SchemaVersion:      batchAReportSchemaVersion,
		Scenario:           batchAScenario,
		Status:             "failed",
		ProviderMode:       cfg.Fobrain.ConnectorStatus.Mode,
		WorkspaceID:        cfg.Fobrain.WorkspaceID,
		TLSPrivateCertMode: cfg.Fobrain.TLS.InsecureSkipVerify,
		CapabilityResults:  []batchACapabilityReport{},
		RedactionChecks: []string{
			"secret_absent",
			"auth_header_absent",
			"raw_payload_absent",
			"binding_reference_absent",
			"checkpoint_absent",
			"interrupt_absent",
		},
		FailureCategory: "not_run",
		ReportCreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// fobrainProviderFromConfig 复用服务启动路径的 provider 边界规则，但只服务 smoke 验收。
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

// fobrainPolicyContextFromConfig 只传递可展示的 policy 材料，不传真实凭据。
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

// fobrainCredentialBindingFromConfig 转换配置派生的安全凭据摘要，字段不包含真实 secret。
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

// writeBatchAReport 以稳定 JSON 格式落盘，供 acceptance record 引用。
func writeBatchAReport(outputPath string, report batchAReport) error {
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

// stableFailureCategory 提取 provider 稳定原因码；未知错误统一折叠。
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
