package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBatchELiveSmokeWritesSafePassedReport(t *testing.T) {
	// Batch E live report 只能记录样本是否存在和 StructuredResult 摘要引用，不能记录真实 ID、token 或 raw payload。
	const secret = "batch-e-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/internal_asset/asset-1":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"id": "asset-1", "ip": "10.10.11.12", "hostname": "prod-web-01", "status": "online", "poc_num": 2}})
		case "/api/threat_center/vuln-1":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"id": "vuln-1", "name": "高危组件漏洞", "level": 3, "status_code": 10, "risk_num": 4}})
		case "/api/threat_center/count":
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want POST", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": []map[string]any{{"name": "核心业务", "count": 7}}})
		case "/api/threat_center/relevance/list":
			if r.URL.Query().Get("keyword") == "" || r.URL.Query().Get("vul_name") == "" {
				t.Fatalf("missing relevance query: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{"id": "rel-1", "vul_name": "高危组件漏洞", "level": 4, "relevance_num": 2}}}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "batch-e-report.json")
	writeBatchETestConfig(t, configPath, server.URL, secret)

	err := runBatchELiveSmoke(configPath, reportPath, batchESamples{
		assetID:           "asset-1",
		assetNetworkType:  "internal",
		vulnerabilityID:   "vuln-1",
		businessName:      "核心业务",
		vulnerabilityName: "高危组件漏洞",
	})
	if err != nil {
		t.Fatal(err)
	}
	report := readBatchEReport(t, reportPath)
	if report.Status != "passed" || report.FailureCategory != "none" || len(report.BlocksClaims) != 0 {
		t.Fatalf("report status mismatch: %+v", report)
	}
	if !report.SampleInputs.AssetDetailPresent ||
		!report.SampleInputs.VulnerabilityDetailPresent ||
		!report.SampleInputs.BusinessPresent ||
		!report.SampleInputs.ThreatNamePresent {
		t.Fatalf("sample input mismatch: %+v", report.SampleInputs)
	}
	if len(report.CapabilityResults) != 4 {
		t.Fatalf("capability result count = %d", len(report.CapabilityResults))
	}
	for _, result := range report.CapabilityResults {
		if result.Status != "passed" || result.ResultRef == "" || result.PolicyDecision != "allowed" || result.FailureCategory != "none" {
			t.Fatalf("capability result mismatch: %+v", result)
		}
	}
	assertBatchEReportNoLeak(t, reportPath, secret, "asset-1", "vuln-1", "核心业务", "高危组件漏洞", "10.10.11.12", "prod-web-01")
}

func TestRunBatchELiveSmokeBlocksWhenSamplesMissing(t *testing.T) {
	const secret = "batch-e-secret"
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "batch-e-report.json")
	writeBatchETestConfig(t, configPath, "https://fobrain.example.test", secret)

	err := runBatchELiveSmoke(configPath, reportPath, batchESamples{})
	if err == nil || !strings.Contains(err.Error(), "missing_sample_input") {
		t.Fatalf("runBatchELiveSmoke err = %v, want missing_sample_input", err)
	}
	report := readBatchEReport(t, reportPath)
	if report.Status != "blocked" || report.FailureCategory != "missing_sample_input" || len(report.BlocksClaims) == 0 {
		t.Fatalf("blocked report mismatch: %+v", report)
	}
}

func readBatchEReport(t *testing.T, outputPath string) batchEReport {
	t.Helper()
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	var report batchEReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	return report
}

func assertBatchEReportNoLeak(t *testing.T, outputPath string, forbidden ...string) {
	t.Helper()
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	encoded := strings.ToLower(string(data))
	for _, marker := range append(forbidden, "api_token", "credential_ref", "raw provider", "raw body", "safe_summary", "authorization") {
		if strings.TrimSpace(marker) == "" {
			continue
		}
		if strings.Contains(encoded, strings.ToLower(marker)) {
			t.Fatalf("report leaked forbidden marker %q: %s", marker, encoded)
		}
	}
}

func writeBatchETestConfig(t *testing.T, path string, baseURL string, secret string) {
	t.Helper()
	body := `
server:
  addr: "127.0.0.1:0"
  read_timeout: "10s"
  write_timeout: "30s"
database:
  driver: "sqlite"
  dsn: "` + filepath.Join(t.TempDir(), "eino-workbench.db") + `"
llm:
  provider: "mock"
  base_url: "http://127.0.0.1/mock-llm"
  model: "mock-chat"
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "60s"
  max_tool_timeout: "120s"
fobrain:
  enabled: true
  connector_id: "fobrain"
  workspace_id: "ws_fobrain"
  base_url: "` + baseURL + `/api"
  timeout: "10s"
  credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "authorization"
    api_token: "` + secret + `"
  connector_status:
    mode: "live"
    available: true
  tls:
    insecure_skip_verify: true
`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}
