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

func TestRunBatchALiveSmokeWritesSafeThreeCapabilityReport(t *testing.T) {
	// 使用自签名 HTTPS 服务证明 smoke 读取配置中的 Fobrain 私有证书开关。
	const secret = "local-live-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/user" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		_, _ = w.Write([]byte(`{"data":{"display_name":"张三","department_name":"安全部","role":{"name":"安全运营"},"permissions":[{"name":"资产只读"}],"data_permissions":["本部门"]}}`))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	outputPath := filepath.Join(tempDir, "report.json")
	writeBatchATestConfig(t, configPath, server.URL, secret)

	if err := runBatchALiveSmoke(configPath, outputPath); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	encoded := strings.ToLower(string(data))
	for _, forbidden := range []string{secret, "bearer ", "api_token", "raw provider", "raw body", "credential_ref"} {
		if strings.Contains(encoded, strings.ToLower(forbidden)) {
			t.Fatalf("report leaked forbidden marker %q: %s", forbidden, encoded)
		}
	}

	var report batchAReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != "eino.fobrain_batch_a_live_report.v1" ||
		report.Scenario != "fobrain-batch-a" ||
		report.Status != "passed" ||
		report.ProviderMode != "live" {
		t.Fatalf("report header mismatch: %+v", report)
	}
	if len(report.CapabilityResults) != 3 {
		t.Fatalf("capability count = %d, want 3: %+v", len(report.CapabilityResults), report.CapabilityResults)
	}
	byID := map[string]batchACapabilityReport{}
	for _, capability := range report.CapabilityResults {
		byID[capability.CapabilityID] = capability
		if capability.Status != "passed" ||
			capability.StructuredResultSchema != "tool.structured_result.v1" ||
			capability.PolicyDecision != "allowed" ||
			strings.TrimSpace(capability.SafeSummary) == "" {
			t.Fatalf("capability report mismatch: %+v", capability)
		}
	}
	for _, capabilityID := range []string{
		"connector.fobrain.security",
		"tool.fobrain.current_user_context",
		"tool.fobrain.my_permissions",
	} {
		if _, ok := byID[capabilityID]; !ok {
			t.Fatalf("missing capability %s in %+v", capabilityID, report.CapabilityResults)
		}
	}
}

func writeBatchATestConfig(t *testing.T, path string, baseURL string, secret string) {
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
  base_url: "` + baseURL + `"
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
