package fobrain_test

import (
	"context"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestProviderInvokesMyPermissionsThroughStructuredResult(t *testing.T) {
	// my_permissions 复用 provider 边界的安全权限材料，不能返回 raw policy payload。
	client := &recordingBatchAClient{
		permissions: fobrain.MyPermissionsResult{
			DisplayName:         "王五",
			PermissionNames:     []string{"资产只读", "漏洞只读"},
			DataPermissionNames: []string{"本部门数据"},
		},
	}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	candidate, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityMyPermissions,
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.permissionCalls != 1 || client.currentUserCalls != 0 {
		t.Fatalf("client calls current=%d permissions=%d", client.currentUserCalls, client.permissionCalls)
	}
	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion ||
		!strings.Contains(candidate.SafeSummary, "资产只读") ||
		strings.Contains(candidate.SafeSummary, "raw") {
		t.Fatalf("permissions candidate mismatch: %+v", candidate)
	}
	if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
		t.Fatalf("permissions candidate rejected: %v", err)
	}
}

func TestProviderInvokesConnectorSecurityWithoutBusinessRead(t *testing.T) {
	// connector 状态只展示配置和凭据摘要，不调用业务读取接口替代真实工具。
	client := &recordingBatchAClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:       "ws_fobrain",
		CredentialBinding: boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:   capabilities.ConnectorStatusAvailable,
		Client:            client,
	})

	candidate, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID:  fobrain.CapabilityConnectorSecurity,
		PolicyContext: policyContextForFobrain("ws_fobrain"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.currentUserCalls != 0 || client.permissionCalls != 0 {
		t.Fatalf("connector status should not call business client: current=%d permissions=%d", client.currentUserCalls, client.permissionCalls)
	}
	if candidate.SchemaVersion != facts.StructuredResultSchemaVersion ||
		!strings.Contains(candidate.SafeSummary, "Fobrain connector") ||
		!strings.Contains(candidate.SafeSummary, "bound") {
		t.Fatalf("connector candidate mismatch: %+v", candidate)
	}
	if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
		t.Fatalf("connector candidate rejected: %v", err)
	}
}

func TestConnectorStatusCredentialBindingSummaryIsSafe(t *testing.T) {
	// connector 状态摘要只能暴露绑定状态和 workspace，不得泄漏凭据引用、认证参数或 token。
	policyContext := capabilities.PolicyContext{
		WorkspaceID: "ws_fobrain",
		CredentialBinding: capabilities.CredentialBinding{
			WorkspaceID: "ws_fobrain",
			System:      "fobrain",
			Status:      capabilities.CredentialStatusBound,
			DisplayRef:  "bound:fobrain:local",
			OwnerScope:  capabilities.PermissionScopeWorkspace,
			AuditRef:    "audit:fobrain:credential",
		},
		ConnectorStatus: capabilities.ConnectorStatusAvailable,
	}

	candidate := fobrain.BuildConnectorSecurityStructuredResult(policyContext, capabilities.ConnectorStatusAvailable)
	if candidate.ResultRef != "result:fobrain:connector-security" ||
		!strings.Contains(candidate.SafeSummary, "workspace：ws_fobrain") ||
		!strings.Contains(candidate.SafeSummary, "绑定状态：bound") {
		t.Fatalf("connector summary mismatch: %+v", candidate)
	}
	for _, forbidden := range []string{"authorization", "api_token", "token", "credential_ref", "bound:fobrain:local", "audit:fobrain:credential"} {
		if strings.Contains(strings.ToLower(candidate.SafeSummary), strings.ToLower(forbidden)) {
			t.Fatalf("connector summary leaked %q: %s", forbidden, candidate.SafeSummary)
		}
	}
}

func policyContextForFobrain(workspaceID string) capabilities.PolicyContext {
	return capabilities.PolicyContext{
		WorkspaceID:       workspaceID,
		CredentialBinding: boundFobrainCredential(workspaceID),
		ConnectorStatus:   capabilities.ConnectorStatusAvailable,
	}
}

type recordingBatchAClient struct {
	currentUserCalls int
	permissionCalls  int
	currentUser      fobrain.CurrentUserContextResult
	permissions      fobrain.MyPermissionsResult
	err              error
}

func (client *recordingBatchAClient) CurrentUserContext(_ context.Context, _ fobrain.ResolvedCredential) (fobrain.CurrentUserContextResult, error) {
	client.currentUserCalls++
	return client.currentUser, client.err
}

func (client *recordingBatchAClient) MyPermissions(_ context.Context, _ fobrain.ResolvedCredential) (fobrain.MyPermissionsResult, error) {
	client.permissionCalls++
	return client.permissions, client.err
}
