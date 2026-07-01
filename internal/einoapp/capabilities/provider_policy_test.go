package capabilities_test

import (
	"encoding/json"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
)

func TestProviderPolicyAllowsReadWithBoundCredential(t *testing.T) {
	// 只读外部能力在凭据绑定且 connector 可用时才允许执行。
	decision := capabilities.EvaluatePolicy(providerPolicyCapability(), capabilities.PolicyContext{
		WorkspaceID: "ws_123",
		CredentialBinding: capabilities.CredentialBinding{
			WorkspaceID: "ws_123",
			System:      "fobrain",
			Status:      capabilities.CredentialStatusBound,
			DisplayRef:  "bound:fobrain:main",
			OwnerScope:  capabilities.PermissionScopeWorkspace,
		},
		ConnectorStatus: capabilities.ConnectorStatusAvailable,
	})

	if !decision.Allowed || decision.ReasonCode != capabilities.PolicyReasonAllowed {
		t.Fatalf("decision = %+v", decision)
	}
	if decision.PolicyRef != "policy:fobrain:read:v1" || decision.AuditRef == "" {
		t.Fatalf("policy metadata missing: %+v", decision)
	}
	if decision.CredentialBinding == nil || decision.CredentialBinding.DisplayRef != "bound:fobrain:main" {
		t.Fatalf("credential binding mismatch: %+v", decision.CredentialBinding)
	}
}

func TestProviderPolicyBlocksMissingCredential(t *testing.T) {
	decision := capabilities.EvaluatePolicy(providerPolicyCapability(), capabilities.PolicyContext{
		WorkspaceID: "ws_123",
		CredentialBinding: capabilities.CredentialBinding{
			WorkspaceID: "ws_123",
			System:      "fobrain",
			Status:      capabilities.CredentialStatusMissing,
			DisplayRef:  "missing",
			OwnerScope:  capabilities.PermissionScopeWorkspace,
		},
		ConnectorStatus: capabilities.ConnectorStatusAvailable,
	})

	if decision.Allowed || decision.ReasonCode != capabilities.PolicyReasonCredentialMissing {
		t.Fatalf("decision = %+v", decision)
	}
}

func TestProviderPolicyBlocksWorkspaceScopeMismatch(t *testing.T) {
	decision := capabilities.EvaluatePolicy(providerPolicyCapability(), capabilities.PolicyContext{
		WorkspaceID: "ws_123",
		CredentialBinding: capabilities.CredentialBinding{
			WorkspaceID: "ws_other",
			System:      "fobrain",
			Status:      capabilities.CredentialStatusBound,
			DisplayRef:  "bound:fobrain:main",
			OwnerScope:  capabilities.PermissionScopeWorkspace,
		},
		ConnectorStatus: capabilities.ConnectorStatusAvailable,
	})

	if decision.Allowed || decision.ReasonCode != capabilities.PolicyReasonCredentialScopeDenied {
		t.Fatalf("decision = %+v", decision)
	}
}

func TestProviderPolicyBlocksUnavailableConnector(t *testing.T) {
	decision := capabilities.EvaluatePolicy(providerPolicyCapability(), capabilities.PolicyContext{
		WorkspaceID: "ws_123",
		CredentialBinding: capabilities.CredentialBinding{
			WorkspaceID: "ws_123",
			System:      "fobrain",
			Status:      capabilities.CredentialStatusBound,
			DisplayRef:  "bound:fobrain:main",
			OwnerScope:  capabilities.PermissionScopeWorkspace,
		},
		ConnectorStatus: capabilities.ConnectorStatusUnavailable,
	})

	if decision.Allowed || decision.ReasonCode != capabilities.PolicyReasonConnectorTransportUnavailable {
		t.Fatalf("decision = %+v", decision)
	}
}

func TestProviderPolicyRequiresApprovalForWriteExternal(t *testing.T) {
	capability := providerPolicyCapability()
	capability.ID = "tool.fobrain.update_ticket_status"
	capability.PolicyRef = "policy:fobrain:ticket-mutation:v1"
	capability.RiskLevel = capabilities.RiskWrite
	capability.SideEffect = capabilities.SideEffectWriteExternal
	capability.ApprovalRequired = true
	capability.IdempotencyRequired = true

	decision := capabilities.EvaluatePolicy(capability, capabilities.PolicyContext{
		WorkspaceID:     "ws_123",
		ConnectorStatus: capabilities.ConnectorStatusAvailable,
		CredentialBinding: capabilities.CredentialBinding{
			WorkspaceID: "ws_123",
			System:      "fobrain",
			Status:      capabilities.CredentialStatusBound,
			DisplayRef:  "bound:fobrain:main",
			OwnerScope:  capabilities.PermissionScopeWorkspace,
		},
	})

	if decision.Allowed || !decision.RequiresApproval || decision.ReasonCode != capabilities.PolicyReasonApprovalRequired {
		t.Fatalf("decision = %+v", decision)
	}
}

func TestProviderPolicyApprovalTakesPriorityForWriteExternal(t *testing.T) {
	capability := providerPolicyCapability()
	capability.ID = "tool.fobrain.update_ticket_status"
	capability.PolicyRef = "policy:fobrain:ticket-mutation:v1"
	capability.RiskLevel = capabilities.RiskWrite
	capability.SideEffect = capabilities.SideEffectWriteExternal
	capability.ApprovalRequired = true
	capability.IdempotencyRequired = true

	decision := capabilities.EvaluatePolicy(capability, capabilities.PolicyContext{
		WorkspaceID:     "ws_123",
		ConnectorStatus: capabilities.ConnectorStatusAvailable,
		CredentialBinding: capabilities.CredentialBinding{
			WorkspaceID: "ws_123",
			System:      "fobrain",
			Status:      capabilities.CredentialStatusMissing,
			DisplayRef:  "missing",
			OwnerScope:  capabilities.PermissionScopeWorkspace,
		},
	})

	if decision.ReasonCode != capabilities.PolicyReasonApprovalRequired || !decision.RequiresApproval {
		t.Fatalf("write external must stop at approval before credential lookup: %+v", decision)
	}
}

func TestProviderPolicyDecisionDoesNotLeakCredentialMaterial(t *testing.T) {
	decision := capabilities.EvaluatePolicy(providerPolicyCapability(), capabilities.PolicyContext{
		WorkspaceID: "ws_123",
		CredentialBinding: capabilities.CredentialBinding{
			SchemaVersion: "token=secret",
			WorkspaceID:   "raw_config",
			System:        "fobrain-secret",
			Status:        capabilities.CredentialStatus("Authorization: Bearer secret"),
			DisplayRef:    "bound:local:raw-config",
			OwnerScope:    capabilities.PermissionScopeWorkspace,
			UpdatedAt:     "password=db-local",
			AuditRef:      "audit:raw_config",
		},
		ConnectorStatus: capabilities.ConnectorStatusAvailable,
	})

	encoded, err := json.Marshal(decision)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"credential_ref", "sk-secret", "token=secret", "Authorization", "password", "apikey", "raw_payload", "provider_payload", "raw_config", "raw-config", "raw config", "secret"} {
		if strings.Contains(string(encoded), forbidden) {
			t.Fatalf("decision leaked %q: %s", forbidden, encoded)
		}
	}
}

func TestProviderPolicyDecisionFoldsUnsafePolicyRef(t *testing.T) {
	for _, policyRef := range []string{"https://policy.example.test/v1", "dsn=file:data/db"} {
		capability := providerPolicyCapability()
		capability.PolicyRef = policyRef

		decision := capabilities.EvaluatePolicy(capability, capabilities.PolicyContext{
			WorkspaceID: "ws_123",
			CredentialBinding: capabilities.CredentialBinding{
				WorkspaceID: "ws_123",
				System:      "fobrain",
				Status:      capabilities.CredentialStatusBound,
				DisplayRef:  "bound:fobrain:main",
				OwnerScope:  capabilities.PermissionScopeWorkspace,
			},
			ConnectorStatus: capabilities.ConnectorStatusAvailable,
		})

		encoded, err := json.Marshal(decision)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(encoded), policyRef) {
			t.Fatalf("decision leaked policy_ref material %q: %s", policyRef, encoded)
		}
		if decision.PolicyRef != "policy:fobrain:low:v1" {
			t.Fatalf("unsafe policy_ref should fall back to derived safe ref, got %q", decision.PolicyRef)
		}
	}
}

func providerPolicyCapability() capabilities.Capability {
	return capabilities.Capability{
		ID:                      "tool.fobrain.asset.read",
		ProviderID:              "fobrain",
		ToolName:                "fobrain_asset_read",
		DisplayName:             "资产查询",
		RiskLevel:               capabilities.RiskReadOnly,
		SideEffect:              capabilities.SideEffectReadExternal,
		PolicyRef:               "policy:fobrain:read:v1",
		PermissionScope:         capabilities.PermissionScopeWorkspace,
		CredentialBindingPolicy: capabilities.CredentialBindingRequired,
		ConnectorID:             "fobrain",
	}
}
