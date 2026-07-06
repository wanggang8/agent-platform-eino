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

func TestRunBatchCLiveSmokeWritesSafePassedReport(t *testing.T) {
	// Batch C live report 只能记录 StructuredResult 引用和样本存在性，不能泄漏 token、筛选值或 raw payload。
	const secret = "batch-c-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/business":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{"id": "biz-1", "name": "核心业务", "level": "重要", "owner": "张三"}}}})
		case "/api/external_ip_asset":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{"id": "asset-1", "ip": "203.0.113.10", "risk_level": "high", "business": "核心业务"}}}})
		case "/api/threat_center/count":
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want POST", r.Method)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": []map[string]any{{"key": "待修复", "count": 3}}})
		case "/api/ticket/pending":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{"id": "ticket-1", "title": "漏洞修复", "status": "pending", "assignee": "张三"}}}})
		case "/api/threat_center/relevance/ip_stats":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": []map[string]any{{"key": "外网", "count": 5}}})
		case "/api/threat_center/relevance/vul_stats":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": []map[string]any{{"key": "高危", "count": 2}}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "batch-c-report.json")
	writeBatchCTestConfig(t, configPath, server.URL, secret)

	err := runBatchCLiveSmoke(configPath, reportPath, batchCSamples{
		keyword:      "核心业务",
		businessName: "核心业务",
		severity:     "high",
		status:       "pending",
		person:       "张三",
		field:        "network",
		timeRange:    "7d",
	})
	if err != nil {
		t.Fatal(err)
	}
	report := readBatchCReport(t, reportPath)
	if report.Status != "passed" || report.FailureCategory != "none" || len(report.BlocksClaims) != 0 {
		t.Fatalf("report status mismatch: %+v", report)
	}
	if !report.SampleInputs.KeywordPresent || !report.SampleInputs.BusinessNamePresent ||
		!report.SampleInputs.SeverityPresent || !report.SampleInputs.StatusPresent || !report.SampleInputs.PersonPresent ||
		!report.SampleInputs.FieldPresent || !report.SampleInputs.TimeRangePresent {
		t.Fatalf("sample input mismatch: %+v", report.SampleInputs)
	}
	if len(report.CapabilityResults) != 6 {
		t.Fatalf("capability result count = %d", len(report.CapabilityResults))
	}
	for _, result := range report.CapabilityResults {
		if result.Status != "passed" || result.ResultState != "resolved" || result.ResultRef == "" || result.PolicyDecision != "allowed" || result.FailureCategory != "none" || result.ItemCount == 0 {
			t.Fatalf("capability result mismatch: %+v", result)
		}
	}
	assertBatchCReportNoLeak(t, reportPath, secret, "核心业务", "张三", "203.0.113.10", "ticket-1", "asset-1")
}

func TestRunBatchCLiveSmokeBlocksWhenAnyCapabilityReturnsEmpty(t *testing.T) {
	// Batch C 真实验收不能用空结果冒充通过；空 StructuredResult 必须阻断 live pass 声明。
	const secret = "batch-c-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/business", "/api/external_ip_asset", "/api/ticket/pending":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{}}})
		case "/api/threat_center/count", "/api/threat_center/relevance/ip_stats", "/api/threat_center/relevance/vul_stats":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": []map[string]any{}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "batch-c-report.json")
	writeBatchCTestConfig(t, configPath, server.URL, secret)

	err := runBatchCLiveSmoke(configPath, reportPath, batchCSamples{})
	if err == nil || !strings.Contains(err.Error(), "empty_result") {
		t.Fatalf("runBatchCLiveSmoke err = %v, want empty_result", err)
	}
	report := readBatchCReport(t, reportPath)
	if report.Status != "blocked" || report.FailureCategory != "empty_result" || len(report.BlocksClaims) == 0 {
		t.Fatalf("blocked report mismatch: %+v", report)
	}
	for _, result := range report.CapabilityResults {
		if result.Status != "blocked" || result.ResultState != "empty" || result.FailureCategory != "empty_result" {
			t.Fatalf("empty capability mismatch: %+v", result)
		}
	}
	assertBatchCReportNoLeak(t, reportPath, secret)
}

func readBatchCReport(t *testing.T, outputPath string) batchCReport {
	t.Helper()
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	var report batchCReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	return report
}

func assertBatchCReportNoLeak(t *testing.T, outputPath string, forbidden ...string) {
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

func writeBatchCTestConfig(t *testing.T, path string, baseURL string, secret string) {
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
