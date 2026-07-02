package fobrain_test

import (
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestProviderCatalogRegistersBatchACapabilities(t *testing.T) {
	// Phase 8 Batch A 只注册 connector、当前用户和权限三类能力，不提前声明 Batch B-E。
	provider := fobrain.NewProvider(fobrain.ProviderConfig{})

	catalog, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog) != 3 {
		t.Fatalf("catalog length = %d, want Batch A three capabilities: %+v", len(catalog), catalog)
	}
	byID := map[string]capabilities.Capability{}
	for _, capability := range catalog {
		byID[capability.ID] = capability
	}

	capability := byID[fobrain.CapabilityCurrentUserContext]
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

	permissions := byID[fobrain.CapabilityMyPermissions]
	if permissions.ID != fobrain.CapabilityMyPermissions ||
		permissions.ResultSchema != fobrain.BusinessResultSchemaVersion ||
		permissions.CredentialBindingPolicy != capabilities.CredentialBindingRequired ||
		permissions.ConnectorID != fobrain.ConnectorID ||
		permissions.SideEffect != capabilities.SideEffectReadExternal {
		t.Fatalf("my permissions capability mismatch: %+v", permissions)
	}

	connector := byID[fobrain.CapabilityConnectorSecurity]
	if connector.ID != fobrain.CapabilityConnectorSecurity ||
		connector.ResultSchema != fobrain.BusinessResultSchemaVersion ||
		connector.CredentialBindingPolicy != capabilities.CredentialBindingOptional ||
		connector.ConnectorID != "" ||
		connector.SideEffect != capabilities.SideEffectReadExternal {
		t.Fatalf("connector capability mismatch: %+v", connector)
	}
}
