package bootstrap_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/bootstrap"
)

func TestLoadConfigReadsYAMLFile(t *testing.T) {
	// 后端配置必须来自本地 YAML 文件，避免把密钥和部署参数散落在环境变量或代码中。
	path := writeConfig(t, `
server:
  addr: "127.0.0.1:19091"
  read_timeout: "2s"
  write_timeout: "3s"
database:
  driver: "sqlite"
  dsn: "data/test.db"
llm:
  provider: "mock"
  base_url: "https://llm.example.test/v1"
  model: "mock-chat"
  timeout_ms: 8000
  credential_binding:
    schema_version: "eino.provider_credential_binding.v1"
    workspace_id: "ws-demo"
    system: "llm"
    status: "bound"
    display_ref: "bound:llm:local"
    owner_scope: "workspace"
  network_safety:
    require_https: true
    allow_local_http: false
    block_private_networks: true
    allow_redirects: false
security:
  redact_secrets: true
  allow_private_network: false
observability:
  log_level: "debug"
  enable_request_log: true
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
capabilities:
  - id: "cap.smoke.read"
    provider_id: "phase3-smoke"
    tool_name: "phase3_smoke_read"
    display_name: "Phase 3 只读验证"
    description: "验证配置驱动的能力注册"
    result_schema: "structured_result.v1"
    risk_level: "read_only"
    approval_required: false
    timeout: "5s"
`)
	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Server.Addr != "127.0.0.1:19091" {
		t.Fatalf("server addr = %q, want config file value", cfg.Server.Addr)
	}
	if cfg.Server.ReadTimeout != 2*time.Second {
		t.Fatalf("read timeout = %s, want 2s", cfg.Server.ReadTimeout)
	}
	if cfg.LLM.CredentialBinding.DisplayRef != "bound:llm:local" {
		t.Fatalf("credential binding was not loaded from local config")
	}
	if len(cfg.Capabilities) != 1 || cfg.Capabilities[0].ID != "cap.smoke.read" {
		t.Fatalf("capabilities were not loaded from local config: %+v", cfg.Capabilities)
	}
}

func TestLoadConfigRejectsMissingConfigPath(t *testing.T) {
	_, err := bootstrap.LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("LoadConfig missing file error = nil")
	}
}

func TestConfigValidationRejectsMissingFileBackedSettings(t *testing.T) {
	path := writeConfig(t, `
server:
  addr: "127.0.0.1:19091"
llm:
  provider: "mock"
  model: "mock-chat"
`)

	_, err := bootstrap.LoadConfig(path)
	if err == nil {
		t.Fatal("LoadConfig with missing required sections error = nil")
	}
}

func TestConfigValidationRejectsInvalidCapability(t *testing.T) {
	// capability 元数据缺失时必须在启动前失败，不能让未知工具进入运行期选择。
	path := writeConfig(t, `
server:
  addr: "127.0.0.1:19091"
  read_timeout: "2s"
  write_timeout: "3s"
database:
  driver: "sqlite"
  dsn: "data/test.db"
llm:
  provider: "mock"
  base_url: "https://llm.example.test/v1"
  model: "mock-chat"
  timeout_ms: 8000
  credential_binding:
    schema_version: "eino.provider_credential_binding.v1"
    workspace_id: "ws-demo"
    system: "llm"
    status: "bound"
    display_ref: "bound:llm:local"
  network_safety:
    require_https: true
    allow_local_http: false
    block_private_networks: true
    allow_redirects: false
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
capabilities:
  - id: "cap.invalid"
    provider_id: "phase3-smoke"
    tool_name: "phase3_smoke_read"
    risk_level: "medium"
`)

	_, err := bootstrap.LoadConfig(path)
	if err == nil || !strings.Contains(err.Error(), "risk_level") {
		t.Fatalf("LoadConfig invalid capability err = %v", err)
	}
}

func TestRedactedSummaryDoesNotExposeSecrets(t *testing.T) {
	// 配置摘要只允许进入日志的安全字段，URL 中的鉴权信息和查询密钥必须被移除。
	path := writeConfig(t, `
server:
  addr: "127.0.0.1:19091"
  read_timeout: "2s"
  write_timeout: "3s"
database:
  driver: "sqlite"
  dsn: "file:data/secret.db?password=db-secret"
llm:
  provider: "mock"
  base_url: "https://user:llm-secret@llm.example.test/v1?api_key=query-secret"
  model: "mock-chat"
  timeout_ms: 8000
  credential_binding:
    schema_version: "eino.provider_credential_binding.v1"
    workspace_id: "ws-demo"
    system: "llm"
    status: "bound"
    display_ref: "bound:llm:local"
    owner_scope: "workspace"
  network_safety:
    require_https: true
    allow_local_http: false
    block_private_networks: true
    allow_redirects: false
security:
  redact_secrets: true
  allow_private_network: false
observability:
  log_level: "debug"
  enable_request_log: true
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	summary := cfg.RedactedSummary()
	encoded := summary.String()
	for _, secret := range []string{"db-secret", "llm-secret", "query-secret", "api_key"} {
		if strings.Contains(encoded, secret) {
			t.Fatalf("redacted summary leaked %q: %s", secret, encoded)
		}
	}
	if !strings.Contains(encoded, "https://llm.example.test/v1") {
		t.Fatalf("redacted summary should keep safe base url origin/path: %s", encoded)
	}
}

func TestRedactedSummaryIncludesBudgetWithoutSecrets(t *testing.T) {
	path := writeConfig(t, `
server:
  addr: "127.0.0.1:19091"
  read_timeout: "2s"
  write_timeout: "3s"
database:
  driver: "sqlite"
  dsn: "data/test.db"
llm:
  provider: "mock"
  base_url: "https://llm.example.test/v1"
  model: "mock-chat"
  timeout_ms: 8000
  credential_binding:
    schema_version: "eino.provider_credential_binding.v1"
    workspace_id: "ws-demo"
    system: "llm"
    status: "bound"
    display_ref: "bound:llm:local"
    owner_scope: "workspace"
  network_safety:
    require_https: true
    allow_local_http: false
    block_private_networks: true
    allow_redirects: false
security:
  redact_secrets: true
  allow_private_network: false
observability:
  log_level: "debug"
  enable_request_log: true
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	encoded := cfg.RedactedSummary().String()
	if !strings.Contains(encoded, "default_timeout: 11s") || !strings.Contains(encoded, "max_tool_timeout: 22s") {
		t.Fatalf("redacted summary missing budgets: %s", encoded)
	}
}

func TestRedactedSummaryIncludesCapabilityMetadataOnly(t *testing.T) {
	path := writeConfig(t, `
server:
  addr: "127.0.0.1:19091"
  read_timeout: "2s"
  write_timeout: "3s"
database:
  driver: "sqlite"
  dsn: "data/test.db"
llm:
  provider: "mock"
  base_url: "https://llm.example.test/v1"
  model: "mock-chat"
  timeout_ms: 8000
  credential_binding:
    schema_version: "eino.provider_credential_binding.v1"
    workspace_id: "ws-demo"
    system: "llm"
    status: "bound"
    display_ref: "bound:llm:local"
  network_safety:
    require_https: true
    allow_local_http: false
    block_private_networks: true
    allow_redirects: false
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
capabilities:
  - id: "cap.smoke.read"
    provider_id: "phase3-smoke"
    tool_name: "phase3_smoke_read"
    description: "contains no secret"
    result_schema: "structured_result.v1"
    risk_level: "read_only"
`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	encoded := cfg.RedactedSummary().String()
	if !strings.Contains(encoded, "cap.smoke.read") || !strings.Contains(encoded, "phase3-smoke") {
		t.Fatalf("redacted summary missing capability metadata: %s", encoded)
	}
	if strings.Contains(encoded, "phase3_smoke_read") || strings.Contains(encoded, "contains no secret") {
		t.Fatalf("redacted summary leaked detailed capability internals: %s", encoded)
	}
}

func writeConfig(t *testing.T, body string) string {
	t.Helper()

	// 测试配置使用 0600，模拟本地密钥配置文件的最小权限约束。
	path := filepath.Join(t.TempDir(), "eino-workbench.yaml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
