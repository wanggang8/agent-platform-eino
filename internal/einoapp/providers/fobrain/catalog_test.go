package fobrain_test

import (
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestProviderCatalogRegistersReadonlyPoC(t *testing.T) {
	// Phase 5 只允许注册一个只读 PoC 能力，避免提前声明完整 Fobrain 恢复。
	provider := fobrain.NewProvider(fobrain.ProviderConfig{})

	catalog, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog) != 1 {
		t.Fatalf("catalog length = %d, want one Phase 5 PoC capability: %+v", len(catalog), catalog)
	}

	capability := catalog[0]
	if capability.ID != fobrain.CapabilityCurrentUserContext {
		t.Fatalf("capability id = %q", capability.ID)
	}
	if capability.ProviderID != fobrain.ProviderID {
		t.Fatalf("provider id = %q", capability.ProviderID)
	}
	if capability.ToolName != fobrain.CapabilityCurrentUserContext {
		t.Fatalf("tool name = %q", capability.ToolName)
	}
	if capability.ResultSchema != fobrain.BusinessResultSchemaVersion {
		t.Fatalf("result schema = %q", capability.ResultSchema)
	}
	if capability.RiskLevel != capabilities.RiskLow {
		t.Fatalf("risk level = %q", capability.RiskLevel)
	}
	if capability.SideEffect != capabilities.SideEffectReadExternal {
		t.Fatalf("side effect = %q", capability.SideEffect)
	}
	if capability.PolicyRef != fobrain.PolicyRead {
		t.Fatalf("policy ref = %q", capability.PolicyRef)
	}
	if capability.PermissionScope != capabilities.PermissionScopeWorkspace {
		t.Fatalf("permission scope = %q", capability.PermissionScope)
	}
	if capability.CredentialBindingPolicy != capabilities.CredentialBindingRequired {
		t.Fatalf("credential binding policy = %q", capability.CredentialBindingPolicy)
	}
	if capability.ConnectorID != fobrain.ConnectorID {
		t.Fatalf("connector id = %q", capability.ConnectorID)
	}
	if capability.ApprovalRequired {
		t.Fatal("readonly PoC must not require approval")
	}
	if capability.IdempotencyRequired {
		t.Fatal("readonly PoC must not require idempotency key")
	}
}
