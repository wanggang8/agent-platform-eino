package bootstrap_test

import (
	"os"
	"path/filepath"
	"testing"

	"agent-platform-eino/internal/einoapp/bootstrap"
)

func TestLoadM1ConfigAcceptsOnlyWalkingSkeletonComposition(t *testing.T) {
	path := filepath.Join(t.TempDir(), "m1.yaml")
	content := []byte("server:\n  addr: 127.0.0.1:8080\n  read_timeout: 10s\n  write_timeout: 30s\ndatabase:\n  driver: sqlite\n  dsn: facts.db\n  busy_timeout: 3s\n  pool_size: 2\nfixture:\n  path: docs/fixtures/new-vulnerability-walking-skeleton.v1.json\n  case: resolved\nfreshness:\n  policy_version: m1.v1\n  ttl: 30m\ncapability:\n  timeout: 5s\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := bootstrap.LoadM1Config(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.Database.PoolSize != 2 || config.Fixture.Case != "resolved" {
		t.Fatalf("config = %+v", config)
	}
}

func TestLoadM1ConfigRejectsLegacyProviderFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "m1.yaml")
	if err := os.WriteFile(path, []byte("llm:\n  provider: openai_compatible\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := bootstrap.LoadM1Config(path); err == nil {
		t.Fatal("legacy provider config must not be silently accepted")
	}
}
