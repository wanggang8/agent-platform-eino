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

	if cfg.Server.Addr != "127.0.0.1:19091" {
		t.Fatalf("server addr = %q, want config file value", cfg.Server.Addr)
	}
	if cfg.Server.ReadTimeout != 2*time.Second {
		t.Fatalf("read timeout = %s, want 2s", cfg.Server.ReadTimeout)
	}
	if cfg.LLM.CredentialBinding.DisplayRef != "bound:llm:local" {
		t.Fatalf("credential binding was not loaded from local config")
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

func TestRedactedSummaryDoesNotExposeSecrets(t *testing.T) {
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

	summary := cfg.RedactedSummary()
	encoded := summary.String()
	for _, secret := range []string{"db-secret"} {
		if strings.Contains(encoded, secret) {
			t.Fatalf("redacted summary leaked %q: %s", secret, encoded)
		}
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

func writeConfig(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "eino-workbench.yaml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
