package fobrain_test

import (
	"context"
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestClientReceivesResolvedCredentialOnlyAfterPolicyAllowed(t *testing.T) {
	client := &recordingCurrentUserClient{result: fobrain.CurrentUserContextResult{
		DisplayName: "张三",
		Department:  "安全部",
		Role:        "安全运营",
	}}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:        "ws_fobrain",
		CredentialBinding:  boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:    capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{WorkspaceID: "ws_fobrain", AuthParam: "authorization", APIToken: "local-secret"},
		Client:             client,
	})

	_, err := provider.Invoke(context.Background(), invocationForWorkspace("ws_fobrain"))
	if err != nil {
		t.Fatal(err)
	}
	if client.calls != 1 {
		t.Fatalf("client calls = %d, want 1", client.calls)
	}
	if client.lastCredential.APIToken != "local-secret" ||
		client.lastCredential.AuthParam != "authorization" ||
		client.lastCredential.WorkspaceID != "ws_fobrain" {
		t.Fatalf("client did not receive resolved credential: %+v", client.lastCredential)
	}
}

func boundFobrainCredential(workspaceID string) capabilities.CredentialBinding {
	return capabilities.CredentialBinding{
		WorkspaceID: workspaceID,
		System:      "fobrain",
		Status:      capabilities.CredentialStatusBound,
		DisplayRef:  "bound:fobrain:local",
		OwnerScope:  capabilities.PermissionScopeWorkspace,
	}
}

func invocationForWorkspace(workspaceID string) capabilities.InvocationRequest {
	return capabilities.InvocationRequest{
		CapabilityID: fobrain.CapabilityCurrentUserContext,
		PolicyContext: capabilities.PolicyContext{
			WorkspaceID: workspaceID,
		},
	}
}

type recordingCurrentUserClient struct {
	calls          int
	lastCredential fobrain.ResolvedCredential
	result         fobrain.CurrentUserContextResult
	err            error
}

func (client *recordingCurrentUserClient) CurrentUserContext(_ context.Context, credential fobrain.ResolvedCredential) (fobrain.CurrentUserContextResult, error) {
	client.calls++
	client.lastCredential = credential
	return client.result, client.err
}
