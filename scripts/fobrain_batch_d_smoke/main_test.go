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

func TestRunBatchDLiveSmokeWritesSafePassedReport(t *testing.T) {
	// Batch D live report 只记录 StructuredResult 出口，不记录样本值、token 或 raw payload。
	const secret = "batch-d-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/user":
			_, _ = w.Write([]byte(`{"data":{"username":"当前用户","department_name":"安全部","role":{"name":"只读"}}}`))
		case "/api/asset":
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[{"id":"asset-1","ip":"10.10.11.12","hostname":["prod-web-01"],"oper_info":[{"name":"张三"}],"status":"online"}],"page":1,"per_page":20,"total":1}}`))
		case "/api/threat_center":
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[{"id":"vuln-1","name":"高危组件漏洞","level":"3","statusCode":10,"person_info":{"name":"张三"},"risk_num":2}],"page":1,"per_page":20,"total":1}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	outputPath := filepath.Join(tempDir, "report.json")
	writeBatchDTestConfig(t, configPath, server.URL, secret)

	err := runBatchDLiveSmoke(configPath, outputPath, batchDSamples{owner: "张三", department: "安全部", ip: "10.10.11.12"})
	if err != nil {
		t.Fatal(err)
	}
	report := readBatchDReport(t, outputPath)
	if report.Status != "passed" || report.FailureCategory != "none" || len(report.BlocksClaims) != 0 {
		t.Fatalf("report status mismatch: %+v", report)
	}
	if len(report.CapabilityResults) != 6 {
		t.Fatalf("capability count = %d", len(report.CapabilityResults))
	}
	for _, result := range report.CapabilityResults {
		if result.Status != "passed" || result.StructuredResultSchema != "tool.structured_result.v1" || result.ResultRef == "" {
			t.Fatalf("capability result mismatch: %+v", result)
		}
	}
	assertBatchDReportNoLeak(t, outputPath, secret, "张三", "安全部", "10.10.11.12")
}

func TestRunBatchDLiveSmokeBlocksWhenSamplesMissing(t *testing.T) {
	// 缺少部门/IP 样本时只能生成 blocking report，不能声明 Batch D live pass。
	const secret = "batch-d-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user" {
			_, _ = w.Write([]byte(`{"data":{"username":"当前用户"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"code":0,"data":{"items":[]}}`))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	outputPath := filepath.Join(tempDir, "report.json")
	writeBatchDTestConfig(t, configPath, server.URL, secret)

	err := runBatchDLiveSmoke(configPath, outputPath, batchDSamples{})
	if err != nil {
		t.Fatal(err)
	}
	report := readBatchDReport(t, outputPath)
	if report.Status != "blocked" || report.FailureCategory != "missing_sample_input" {
		t.Fatalf("report status mismatch: %+v", report)
	}
	if len(report.BlocksClaims) == 0 {
		t.Fatalf("blocked report must block claims: %+v", report)
	}
	blocked := 0
	for _, result := range report.CapabilityResults {
		if result.Status == "blocked" {
			blocked++
		}
	}
	if blocked == 0 {
		t.Fatalf("missing samples should block at least one capability: %+v", report.CapabilityResults)
	}
	assertBatchDReportNoLeak(t, outputPath, secret)
}

func readBatchDReport(t *testing.T, outputPath string) batchDReport {
	t.Helper()
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	var report batchDReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	return report
}

func assertBatchDReportNoLeak(t *testing.T, outputPath string, forbidden ...string) {
	t.Helper()
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	encoded := strings.ToLower(string(data))
	for _, marker := range append(forbidden, "api_token", "credential_ref", "raw provider", "raw body") {
		if strings.TrimSpace(marker) == "" {
			continue
		}
		if strings.Contains(encoded, strings.ToLower(marker)) {
			t.Fatalf("report leaked forbidden marker %q: %s", marker, encoded)
		}
	}
}

func writeBatchDTestConfig(t *testing.T, path string, baseURL string, secret string) {
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
  timeout: "30s"
security:
  redact_secrets: true
  allow_private_network: false
observability:
  log_level: "debug"
  enable_request_log: true
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
