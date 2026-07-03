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
  api_key: "sk-local-test"
  model: "mock-chat"
  timeout: "8s"
security:
  redact_secrets: true
  allow_private_network: false
observability:
  log_level: "debug"
  enable_request_log: true
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
  max_model_calls_per_run: 3
  max_tool_calls_per_run: 5
  max_input_tokens_per_run: 4096
capabilities:
  - id: "cap.smoke.read"
    provider_id: "phase3-smoke"
    tool_name: "phase3_smoke_read"
    display_name: "Phase 3 只读验证"
    description: "验证配置驱动的能力注册"
    result_schema: "structured_result.v1"
    risk_level: "low"
    side_effect: "read_external"
    policy_ref: "policy:smoke:read:v1"
    permission_scope: "workspace"
    credential_binding_policy: "none"
    approval_required: false
    idempotency_required: false
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
	if cfg.LLM.TimeoutMillis != 8000 {
		t.Fatalf("timeout millis = %d, want derived 8000", cfg.LLM.TimeoutMillis)
	}
	if cfg.LLM.CredentialBinding.DisplayRef != "bound:llm:mock" {
		t.Fatalf("credential binding was not derived from local config: %+v", cfg.LLM.CredentialBinding)
	}
	if cfg.LLM.CredentialBinding.OwnerScope != "system" {
		t.Fatalf("credential owner scope must match schema: %+v", cfg.LLM.CredentialBinding)
	}
	if cfg.LLM.NetworkSafety.AllowedHosts[0] != "llm.example.test" {
		t.Fatalf("network safety allowed hosts not derived: %+v", cfg.LLM.NetworkSafety)
	}
	if len(cfg.Capabilities) != 1 || cfg.Capabilities[0].ID != "cap.smoke.read" {
		t.Fatalf("capabilities were not loaded from local config: %+v", cfg.Capabilities)
	}
	if cfg.Budgets.MaxModelCallsPerRun != 3 || cfg.Budgets.MaxToolCallsPerRun != 5 || cfg.Budgets.MaxInputTokensPerRun != 4096 {
		t.Fatalf("budget counters were not loaded: %+v", cfg.Budgets)
	}
	if cfg.Capabilities[0].SideEffect != "read_external" ||
		cfg.Capabilities[0].PolicyRef != "policy:smoke:read:v1" ||
		cfg.Capabilities[0].PermissionScope != "workspace" ||
		cfg.Capabilities[0].CredentialBindingPolicy != "none" {
		t.Fatalf("capability policy fields were not loaded: %+v", cfg.Capabilities[0])
	}
}

func TestConfigAppliesDefaultBudgetCounters(t *testing.T) {
	// 新增预算计数字段必须能从旧配置安全派生默认值，避免本地配置升级时直接启动失败。
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
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Budgets.MaxModelCallsPerRun <= 0 || cfg.Budgets.MaxToolCallsPerRun <= 0 || cfg.Budgets.MaxInputTokensPerRun <= 0 {
		t.Fatalf("default budget counters not applied: %+v", cfg.Budgets)
	}
	summary := cfg.RedactedSummary()
	budgets, ok := summary["budgets"].(map[string]any)
	if !ok {
		t.Fatalf("budget summary missing: %+v", summary)
	}
	if budgets["max_model_calls_per_run"] == nil || budgets["max_tool_calls_per_run"] == nil || budgets["max_input_units_per_run"] == nil {
		t.Fatalf("budget counter summary missing: %+v", budgets)
	}
	encoded := summary.String()
	for _, forbidden := range []string{"token", "api_key", "api_token"} {
		if strings.Contains(encoded, forbidden) {
			t.Fatalf("budget summary leaked %q: %s", forbidden, encoded)
		}
	}
}

func TestConfigValidationRejectsNegativeBudgetCounter(t *testing.T) {
	// 用户显式写错的负数预算不能被默认值静默覆盖，否则会掩盖配置错误。
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
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
  max_model_calls_per_run: -1
`)

	_, err := bootstrap.LoadConfig(path)
	if err == nil || !strings.Contains(err.Error(), "budgets.max_model_calls_per_run") {
		t.Fatalf("negative budget counter err = %v", err)
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
  timeout: "8s"
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
    risk_level: "invalid"
    side_effect: "read_external"
    policy_ref: "policy:smoke:read:v1"
    permission_scope: "workspace"
    credential_binding_policy: "none"
`)

	_, err := bootstrap.LoadConfig(path)
	if err == nil || !strings.Contains(err.Error(), "risk_level") {
		t.Fatalf("LoadConfig invalid capability err = %v", err)
	}
}

func TestConfigValidationRejectsInvalidLLMBaseURL(t *testing.T) {
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
  base_url: "not-a-url"
  model: "mock-chat"
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
`)

	_, err := bootstrap.LoadConfig(path)
	if err == nil || !strings.Contains(err.Error(), "llm.base_url") {
		t.Fatalf("LoadConfig invalid base_url err = %v", err)
	}
}

func TestConfigValidationRequiresAPIKeyForOpenAICompatible(t *testing.T) {
	path := writeConfig(t, `
server:
  addr: "127.0.0.1:19091"
  read_timeout: "2s"
  write_timeout: "3s"
database:
  driver: "sqlite"
  dsn: "data/test.db"
llm:
  provider: "openai_compatible"
  base_url: "https://llm.example.test/v1"
  model: "real-chat"
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
`)

	_, err := bootstrap.LoadConfig(path)
	if err == nil || !strings.Contains(err.Error(), "llm.api_key") {
		t.Fatalf("LoadConfig missing api_key err = %v", err)
	}
}

func TestConfigNormalizesLLMProvider(t *testing.T) {
	path := writeConfig(t, `
server:
  addr: "127.0.0.1:19091"
  read_timeout: "2s"
  write_timeout: "3s"
database:
  driver: "sqlite"
  dsn: "data/test.db"
llm:
  provider: " mock "
  base_url: "http://127.0.0.1/mock-llm"
  model: "mock-chat"
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.LLM.Provider != "mock" || cfg.LLM.CredentialBinding.DisplayRef != "bound:llm:mock" {
		t.Fatalf("llm provider not normalized: provider=%q binding=%+v", cfg.LLM.Provider, cfg.LLM.CredentialBinding)
	}
}

func TestLoadConfigReadsMockMCPServers(t *testing.T) {
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
  timeout: "8s"
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
mcp_mock_servers:
  - server_id: "mock"
    tools:
      - name: "asset_lookup"
        title: "MCP 资产查询"
        description: "通过 mock MCP 查询资产"
        input_schema:
          type: "object"
          properties:
            query: "string"
          required: ["query"]
        output_schema:
          type: "object"
        annotations:
          read_only_hint: true
    results:
      asset_lookup:
        structured_content:
          result_ref: "result:mcp:asset_lookup"
          safe_summary: "MCP 资产查询完成"
`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.MCPMockServers) != 1 || cfg.MCPMockServers[0].ServerID != "mock" {
		t.Fatalf("mcp mock servers not loaded: %+v", cfg.MCPMockServers)
	}
	tool := cfg.MCPMockServers[0].Tools[0]
	if tool.InputSchema.Properties["query"] != "string" || !tool.Annotations.ReadOnlyHint {
		t.Fatalf("mcp mock tool metadata mismatch: %+v", tool)
	}
	if cfg.MCPMockServers[0].Results["asset_lookup"].StructuredContent["safe_summary"] != "MCP 资产查询完成" {
		t.Fatalf("mcp mock result mismatch: %+v", cfg.MCPMockServers[0].Results)
	}
}

func TestConfigValidationRejectsInvalidMockMCPProjectPolicy(t *testing.T) {
	for _, testCase := range []struct {
		name    string
		policy  string
		wantErr string
	}{
		{
			name: "invalid risk",
			policy: `
          risk_level: "unsafe"
          side_effect: "read_external"
          permission_scope: "workspace"`,
			wantErr: "risk_level",
		},
		{
			name: "invalid side effect",
			policy: `
          risk_level: "low"
          side_effect: "network"
          permission_scope: "workspace"`,
			wantErr: "side_effect",
		},
		{
			name: "invalid scope",
			policy: `
          risk_level: "low"
          side_effect: "read_external"
          permission_scope: "tenant"`,
			wantErr: "permission_scope",
		},
		{
			name: "write missing idempotency",
			policy: `
          risk_level: "high"
          side_effect: "write_external"
          permission_scope: "workspace"`,
			wantErr: "write_external",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
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
  timeout: "8s"
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
mcp_mock_servers:
  - server_id: "mock"
    tools:
      - name: "write_tool"
        description: "mock write"
        input_schema:
          type: "object"
          properties:
            ticket: "string"
        project_policy:
`+testCase.policy+`
`)

			_, err := bootstrap.LoadConfig(path)
			if err == nil || !strings.Contains(err.Error(), testCase.wantErr) {
				t.Fatalf("LoadConfig invalid mcp project policy err = %v, want %q", err, testCase.wantErr)
			}
		})
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
  api_key: "sk-local-secret"
  model: "mock-chat"
  timeout: "8s"
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
	for _, secret := range []string{"db-secret", "llm-secret", "query-secret", "sk-local-secret", "api_key"} {
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
  timeout: "8s"
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
  timeout: "8s"
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
    risk_level: "low"
    side_effect: "read_external"
    policy_ref: "policy:smoke:read:v1"
    permission_scope: "workspace"
    credential_binding_policy: "none"
    connector_id: "https://connector.example.test?api_key=db-local"
	`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	encoded := cfg.RedactedSummary().String()
	for _, expected := range []string{
		"cap.smoke.read",
		"phase3-smoke",
		"side_effect: read_external",
		"policy_ref: policy:smoke:read:v1",
		"permission_scope: workspace",
		"credential_binding_policy: none",
		"idempotency_required: false",
	} {
		if !strings.Contains(encoded, expected) {
			t.Fatalf("redacted summary missing capability metadata %q: %s", expected, encoded)
		}
	}
	if !strings.Contains(encoded, "approval_required: false") {
		t.Fatalf("redacted summary missing capability metadata: %s", encoded)
	}
	for _, forbidden := range []string{"phase3_smoke_read", "contains no secret", "structured_result.v1", "https://connector.example.test", "api_key", "db-local"} {
		if strings.Contains(encoded, forbidden) {
			t.Fatalf("redacted summary leaked detailed capability internals %q: %s", forbidden, encoded)
		}
	}
}

func TestRedactedSummaryFoldsUnsafeCapabilityPolicyMetadata(t *testing.T) {
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
  timeout: "8s"
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
    risk_level: "low"
    side_effect: "read_external"
    policy_ref: "https://policy.example.test?password=db-local"
    permission_scope: "workspace"
    credential_binding_policy: "none"
    connector_id: "dsn=file:data/db?raw_config=1"
`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	encoded := cfg.RedactedSummary().String()
	for _, forbidden := range []string{"https://policy.example.test", "password", "db-local", "dsn=file", "raw_config"} {
		if strings.Contains(encoded, forbidden) {
			t.Fatalf("redacted summary leaked unsafe policy metadata %q: %s", forbidden, encoded)
		}
	}
	if !strings.Contains(encoded, "policy_ref: policy:redacted") || !strings.Contains(encoded, "connector_id: redacted") {
		t.Fatalf("redacted summary leaked detailed capability internals: %s", encoded)
	}
}

func TestFobrainConfigDerivesSafeConnectorFields(t *testing.T) {
	// Fobrain provider 配置必须从 YAML 派生安全摘要，真实 token 只留在 provider 边界。
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
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
fobrain:
  enabled: true
  connector_id: "fobrain"
  workspace_id: "ws_fobrain"
  base_url: "https://fobrain.example.test/api"
  timeout: "9s"
  credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "authorization"
    api_token: "fobrain-local-secret"
  connector_status:
    mode: "mock"
    available: true
`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	if !cfg.Fobrain.Enabled || cfg.Fobrain.ConnectorID != "fobrain" || cfg.Fobrain.WorkspaceID != "ws_fobrain" {
		t.Fatalf("fobrain config not loaded: %+v", cfg.Fobrain)
	}
	if cfg.Fobrain.Timeout != 9*time.Second {
		t.Fatalf("fobrain timeout = %s", cfg.Fobrain.Timeout)
	}
	if cfg.Fobrain.CredentialBinding.System != "fobrain" ||
		cfg.Fobrain.CredentialBinding.Status != "bound" ||
		cfg.Fobrain.CredentialBinding.DisplayRef != "bound:fobrain:local" ||
		cfg.Fobrain.CredentialBinding.OwnerScope != "workspace" {
		t.Fatalf("fobrain safe credential binding mismatch: %+v", cfg.Fobrain.CredentialBinding)
	}
	if cfg.Fobrain.ConnectorStatus.Mode != "mock" || !cfg.Fobrain.ConnectorStatus.Available {
		t.Fatalf("fobrain connector status mismatch: %+v", cfg.Fobrain.ConnectorStatus)
	}
	if cfg.Fobrain.Credential.AuthParam != "authorization" {
		t.Fatalf("fobrain auth param = %q, want authorization", cfg.Fobrain.Credential.AuthParam)
	}
}

func TestFobrainConfigValidationRejectsInvalidFields(t *testing.T) {
	for _, testCase := range []struct {
		name        string
		fobrainBody string
		wantErr     string
	}{
		{name: "missing base url", fobrainBody: defaultFobrainConfigBody(`base_url: ""`), wantErr: "fobrain.base_url"},
		{name: "invalid base url", fobrainBody: defaultFobrainConfigBody(`base_url: "not-a-url"`), wantErr: "fobrain.base_url"},
		{name: "missing workspace", fobrainBody: defaultFobrainConfigBody(`workspace_id: ""`), wantErr: "fobrain.workspace_id"},
		{name: "invalid credential status", fobrainBody: defaultFobrainConfigBody(`credential:
    status: "ready"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"`), wantErr: "credential.status"},
		{name: "invalid connector mode", fobrainBody: defaultFobrainConfigBody(`connector_status:
    mode: "auto"
    available: true`), wantErr: "connector_status.mode"},
		{name: "invalid auth param", fobrainBody: defaultFobrainConfigBody(`credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "Authorization: Bearer"
    api_token: "fobrain-local-secret"`), wantErr: "credential.auth_param"},
		{name: "missing auth param", fobrainBody: defaultFobrainConfigBody(`credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    api_token: "fobrain-local-secret"`), wantErr: "credential.auth_param"},
		{name: "live mode requires authorization auth param", fobrainBody: `
enabled: true
connector_id: "fobrain"
workspace_id: "ws_fobrain"
base_url: "https://fobrain.example.test/api"
timeout: "9s"
credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "x-fobrain-auth"
    api_token: "fobrain-local-secret"
connector_status:
    mode: "live"
    available: true`, wantErr: "credential.auth_param must be authorization in live mode"},
		{name: "live bound missing token", fobrainBody: `
enabled: true
connector_id: "fobrain"
workspace_id: "ws_fobrain"
base_url: "https://fobrain.example.test/api"
timeout: "9s"
credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "authorization"
connector_status:
    mode: "live"
    available: true`, wantErr: "credential.api_token"},
		{name: "live mode requires bound credential", fobrainBody: `
enabled: true
connector_id: "fobrain"
workspace_id: "ws_fobrain"
base_url: "https://fobrain.example.test/api"
timeout: "9s"
credential:
    status: "missing"
    display_ref: "missing"
    owner_scope: "workspace"
    auth_param: "authorization"
connector_status:
    mode: "live"
    available: true`, wantErr: "credential.status must be configured or bound in live mode"},
		{name: "live mode rejects http base url", fobrainBody: `
enabled: true
connector_id: "fobrain"
workspace_id: "ws_fobrain"
base_url: "http://fobrain.example.test/api"
timeout: "9s"
credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "authorization"
    api_token: "fobrain-local-secret"
connector_status:
    mode: "live"
    available: true`, wantErr: "fobrain.base_url"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			path := writeConfig(t, fobrainConfigYAML(testCase.fobrainBody))
			_, err := bootstrap.LoadConfig(path)
			if err == nil || !strings.Contains(err.Error(), testCase.wantErr) {
				t.Fatalf("LoadConfig err = %v, want containing %q", err, testCase.wantErr)
			}
		})
	}
}

func TestFobrainConfigAllowsLiveModeWithLocalHTTP(t *testing.T) {
	// 本地验收可以通过 loopback HTTP 访问 Fobrain mock/live 代理；非本地 live HTTP 仍被拒绝。
	path := writeConfig(t, fobrainConfigYAML(`
enabled: true
connector_id: "fobrain"
workspace_id: "ws_fobrain"
base_url: "http://127.0.0.1:3001/api/v1"
timeout: "9s"
credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "authorization"
    api_token: "fobrain-local-secret"
connector_status:
    mode: "live"
    available: true`))

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Fobrain.BaseURL != "http://127.0.0.1:3001/api/v1" || cfg.Fobrain.ConnectorStatus.Mode != "live" {
		t.Fatalf("local live fobrain config mismatch: %+v", cfg.Fobrain)
	}
}

func TestFobrainConfigAllowsLiveModeWithWorkspaceToken(t *testing.T) {
	// Phase 8 Batch A 允许 live mode，但仍必须使用 workspace scope 和本地 ignored token。
	path := writeConfig(t, fobrainConfigYAML(`
enabled: true
connector_id: "fobrain"
workspace_id: "ws_fobrain"
base_url: "https://fobrain.example.test/api"
timeout: "9s"
credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "authorization"
    api_token: "fobrain-local-secret"
connector_status:
    mode: "live"
    available: true
tls:
    insecure_skip_verify: true`))

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Fobrain.ConnectorStatus.Mode != "live" ||
		cfg.Fobrain.Credential.AuthParam != "authorization" ||
		cfg.Fobrain.CredentialBinding.OwnerScope != "workspace" ||
		!cfg.Fobrain.TLS.InsecureSkipVerify {
		t.Fatalf("live fobrain config mismatch: %+v", cfg.Fobrain)
	}
}

func TestFobrainConfigAllowsDisabledWithoutConnectorFields(t *testing.T) {
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
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
fobrain:
  enabled: false
`)

	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Fobrain.Enabled {
		t.Fatalf("fobrain should remain disabled: %+v", cfg.Fobrain)
	}
}

func TestFobrainRedactedSummaryDoesNotExposeCredentialLeak(t *testing.T) {
	path := writeConfig(t, fobrainConfigYAML(defaultFobrainConfigBody(`base_url: "https://user:fobrain-url-secret@fobrain.example.test/api?api_key=query-secret"`)))
	cfg, err := bootstrap.LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	encoded := cfg.RedactedSummary().String()
	for _, forbidden := range []string{"fobrain-local-secret", "fobrain-url-secret", "query-secret", "api_key", "api_token", "token"} {
		if strings.Contains(encoded, forbidden) {
			t.Fatalf("fobrain redacted summary leaked %q: %s", forbidden, encoded)
		}
	}
	for _, expected := range []string{"enabled: true", "connector_id: fobrain", "workspace_id: ws_fobrain", "base_url: https://fobrain.example.test/api", "credential_status: bound", "auth_param: authorization", "connector_mode: mock", "tls_insecure_skip_verify: false"} {
		if !strings.Contains(encoded, expected) {
			t.Fatalf("fobrain redacted summary missing %q: %s", expected, encoded)
		}
	}
}

func fobrainConfigYAML(override string) string {
	if strings.TrimSpace(override) == "" {
		override = defaultFobrainConfigBody("")
	}
	return `
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
security:
  redact_secrets: true
observability:
  log_level: "debug"
budgets:
  default_timeout: "11s"
  max_tool_timeout: "22s"
fobrain:
` + indentYAML(override)
}

func defaultFobrainConfigBody(overrides string) string {
	values := map[string]string{
		"enabled":      `enabled: true`,
		"connector_id": `connector_id: "fobrain"`,
		"workspace_id": `workspace_id: "ws_fobrain"`,
		"base_url":     `base_url: "https://fobrain.example.test/api"`,
		"timeout":      `timeout: "9s"`,
		"credential": `credential:
  status: "bound"
  display_ref: "bound:fobrain:local"
  owner_scope: "workspace"
  auth_param: "authorization"
  api_token: "fobrain-local-secret"`,
		"connector_status": `connector_status:
  mode: "mock"
  available: true`,
		"tls": `tls:
  insecure_skip_verify: false`,
	}
	for _, block := range strings.Split(strings.TrimSpace(overrides), "\n") {
		key := strings.TrimSpace(strings.SplitN(block, ":", 2)[0])
		if _, ok := values[key]; ok {
			values[key] = strings.TrimSpace(overrides)
			break
		}
	}
	return strings.Join([]string{
		values["enabled"],
		values["connector_id"],
		values["workspace_id"],
		values["base_url"],
		values["timeout"],
		values["credential"],
		values["connector_status"],
		values["tls"],
	}, "\n")
}

func indentYAML(body string) string {
	lines := strings.Split(strings.TrimSpace(body), "\n")
	for index, line := range lines {
		lines[index] = "  " + line
	}
	return strings.Join(lines, "\n")
}

func writeConfig(t *testing.T, body string) string {
	t.Helper()

	// 测试配置使用 0600，模拟本地密钥配置文件的最小权限约束。
	path := filepath.Join(t.TempDir(), "eino-workbench.yaml")
	if err := os.WriteFile(path, []byte(strings.TrimSpace(body)+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
