package bootstrap_test

import (
	"testing"

	"agent-platform-eino/internal/einoapp/bootstrap"
)

func TestLoadConfigUsesDefaultAddress(t *testing.T) {
	t.Setenv("EINO_WORKBENCH_ADDR", "")

	cfg := bootstrap.LoadConfig()

	if cfg.Addr != "127.0.0.1:8081" {
		t.Fatalf("Addr = %q, want default 127.0.0.1:8081", cfg.Addr)
	}
}

func TestLoadConfigAllowsAddressOverride(t *testing.T) {
	t.Setenv("EINO_WORKBENCH_ADDR", "127.0.0.1:19091")

	cfg := bootstrap.LoadConfig()

	if cfg.Addr != "127.0.0.1:19091" {
		t.Fatalf("Addr = %q, want env override", cfg.Addr)
	}
}
