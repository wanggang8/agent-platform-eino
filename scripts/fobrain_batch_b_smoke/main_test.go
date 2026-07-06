package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestRunBatchBLiveSmokeWritesSafePassedReport(t *testing.T) {
	// Batch B live report 只能记录 StructuredResult 引用和筛选存在性，不能泄漏当前用户、部门、token 或 raw payload。
	const secret = "batch-b-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/v1/user":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{
				"display_name":    "王五",
				"department_name": "安全部",
			}})
		case "/api/asset":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{
				"id": "asset-1", "hostname": "prod-web-01", "ip": "10.10.11.12", "oper_info": []map[string]any{{"name": "王五"}},
			}}}})
		case "/api/threat_center":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{
				"id": "vul-1", "name": "高危组件漏洞", "level": 3, "status": 10, "person_info": []map[string]any{{"name": "王五"}},
			}}}})
		case "/api/business":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{
				"id": "biz-1", "business_name": "核心业务", "person_info": []map[string]any{{"name": "王五"}}, "risk_num": 2,
			}}}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "batch-b-report.json")
	writeBatchBTestConfig(t, configPath, server.URL, secret)

	err := runBatchBLiveSmoke(configPath, reportPath, batchBSamples{keyword: "核心业务"})
	if err != nil {
		t.Fatal(err)
	}
	report := readBatchBReport(t, reportPath)
	if report.Status != "passed" || report.FailureCategory != "none" || len(report.BlocksClaims) != 0 {
		t.Fatalf("report status mismatch: %+v", report)
	}
	if !report.SampleInputs.KeywordPresent {
		t.Fatalf("sample inputs mismatch: %+v", report.SampleInputs)
	}
	if len(report.CapabilityResults) != 6 {
		t.Fatalf("capability result count = %d", len(report.CapabilityResults))
	}
	for _, result := range report.CapabilityResults {
		if result.Status != "passed" || result.ResultState != "resolved" || result.ResultRef == "" || result.PolicyDecision != "allowed" || result.FailureCategory != "none" || result.ItemCount == 0 {
			t.Fatalf("capability result mismatch: %+v", result)
		}
	}
	assertBatchBReportNoLeak(t, reportPath, secret, "王五", "安全部", "核心业务", "10.10.11.12", "asset-1")
}

func TestRunBatchBLiveSmokeBlocksWhenAnyCapabilityReturnsEmpty(t *testing.T) {
	// Batch B 真实验收不能用空结果冒充通过；空 StructuredResult 必须阻断 live pass 声明。
	const secret = "batch-b-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/v1/user":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{
				"display_name":    "王五",
				"department_name": "安全部",
			}})
		case "/api/asset", "/api/threat_center", "/api/business":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{}}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "batch-b-report.json")
	writeBatchBTestConfig(t, configPath, server.URL, secret)

	err := runBatchBLiveSmoke(configPath, reportPath, batchBSamples{})
	if err == nil || !strings.Contains(err.Error(), "empty_result") {
		t.Fatalf("runBatchBLiveSmoke err = %v, want empty_result", err)
	}
	report := readBatchBReport(t, reportPath)
	if report.Status != "blocked" || report.FailureCategory != "empty_result" || len(report.BlocksClaims) == 0 {
		t.Fatalf("blocked report mismatch: %+v", report)
	}
	for _, result := range report.CapabilityResults {
		if result.Status != "blocked" || result.ResultState != "empty" || result.FailureCategory != "empty_result" {
			t.Fatalf("empty capability mismatch: %+v", result)
		}
	}
	assertBatchBReportNoLeak(t, reportPath, secret, "王五", "安全部")
}

func TestRunBatchBLiveSmokeBlocksWhenCurrentUserDepartmentMissing(t *testing.T) {
	// 当前用户接口缺少部门时，本部门工具不能静默跳过或伪装通过，必须生成 blocking report。
	const secret = "batch-b-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/v1/user":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"display_name": "王五"}})
		case "/api/asset", "/api/threat_center", "/api/business":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{
				"id": "row-1", "name": "安全结果",
			}}}})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "batch-b-report.json")
	writeBatchBTestConfig(t, configPath, server.URL, secret)

	err := runBatchBLiveSmoke(configPath, reportPath, batchBSamples{})
	if err == nil || !strings.Contains(err.Error(), string(fobrain.PolicyReasonMissingCurrentUserScope)) {
		t.Fatalf("runBatchBLiveSmoke err = %v, want missing_current_user_scope", err)
	}
	report := readBatchBReport(t, reportPath)
	if report.Status != "blocked" || report.FailureCategory != string(fobrain.PolicyReasonMissingCurrentUserScope) || len(report.BlocksClaims) == 0 {
		t.Fatalf("missing scope report mismatch: %+v", report)
	}
	blockedDepartmentTools := 0
	for _, result := range report.CapabilityResults {
		if result.FailureCategory == string(fobrain.PolicyReasonMissingCurrentUserScope) {
			blockedDepartmentTools++
			if result.Status != "blocked" || result.ResultState != "not_run" {
				t.Fatalf("missing scope capability mismatch: %+v", result)
			}
		}
	}
	if blockedDepartmentTools != 2 {
		t.Fatalf("blocked department tools = %d, want 2", blockedDepartmentTools)
	}
	assertBatchBReportNoLeak(t, reportPath, secret, "王五")
}

func readBatchBReport(t *testing.T, outputPath string) batchBReport {
	t.Helper()
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	var report batchBReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	return report
}

func assertBatchBReportNoLeak(t *testing.T, outputPath string, forbidden ...string) {
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

func writeBatchBTestConfig(t *testing.T, path string, baseURL string, secret string) {
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
