package fobrain_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestCredentialMissingBlocksClientCall(t *testing.T) {
	// provider 边界必须再次校验凭据，不能依赖 HTTP 或 execution 已经做过正确 policy。
	client := &recordingCurrentUserClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID: "ws_fobrain",
		CredentialBinding: capabilities.CredentialBinding{
			WorkspaceID: "ws_fobrain",
			System:      "fobrain",
			Status:      capabilities.CredentialStatusMissing,
			DisplayRef:  "missing",
			OwnerScope:  capabilities.PermissionScopeWorkspace,
		},
		ConnectorStatus: capabilities.ConnectorStatusAvailable,
		Client:          client,
	})

	_, err := provider.Invoke(context.Background(), invocationForWorkspace("ws_fobrain"))
	if !fobrain.HasReason(err, capabilities.PolicyReasonCredentialMissing) {
		t.Fatalf("err = %v, want credential missing", err)
	}
	if client.calls != 0 {
		t.Fatalf("client calls = %d, want blocked before client", client.calls)
	}
}

func TestCredentialWorkspaceMismatchBlocksClientCall(t *testing.T) {
	client := &recordingCurrentUserClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID: "ws_fobrain",
		CredentialBinding: capabilities.CredentialBinding{
			WorkspaceID: "ws_other",
			System:      "fobrain",
			Status:      capabilities.CredentialStatusBound,
			DisplayRef:  "bound:fobrain:local",
			OwnerScope:  capabilities.PermissionScopeWorkspace,
		},
		ConnectorStatus: capabilities.ConnectorStatusAvailable,
		Client:          client,
	})

	_, err := provider.Invoke(context.Background(), invocationForWorkspace("ws_fobrain"))
	if !fobrain.HasReason(err, capabilities.PolicyReasonCredentialScopeDenied) {
		t.Fatalf("err = %v, want credential scope denied", err)
	}
	if client.calls != 0 {
		t.Fatalf("client calls = %d, want blocked before client", client.calls)
	}
}

func TestInvocationWorkspaceMismatchBlocksClientCall(t *testing.T) {
	// provider 二次校验必须使用 run/request workspace，不能退回配置 workspace 放行凭据。
	client := &recordingCurrentUserClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:       "ws_fobrain",
		CredentialBinding: boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:   capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{
			WorkspaceID: "ws_fobrain",
			APIToken:    "local-secret",
		},
		Client: client,
	})

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{
		CapabilityID: fobrain.CapabilityCurrentUserContext,
		PolicyContext: capabilities.PolicyContext{
			WorkspaceID:       "ws_other",
			CredentialBinding: boundFobrainCredential("ws_fobrain"),
			ConnectorStatus:   capabilities.ConnectorStatusAvailable,
		},
	})
	if !fobrain.HasReason(err, capabilities.PolicyReasonCredentialScopeDenied) {
		t.Fatalf("err = %v, want credential scope denied", err)
	}
	if client.calls != 0 {
		t.Fatalf("client calls = %d, want blocked before client", client.calls)
	}
}

func TestInvocationRequiresWorkspaceContext(t *testing.T) {
	// 直接调用 provider 也必须显式携带 workspace，配置 workspace 不能作为调用事实来源。
	client := &recordingCurrentUserClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:       "ws_fobrain",
		CredentialBinding: boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:   capabilities.ConnectorStatusAvailable,
		CredentialResolver: fobrain.StaticCredentialResolver{
			WorkspaceID: "ws_fobrain",
			APIToken:    "local-secret",
		},
		Client: client,
	})

	_, err := provider.Invoke(context.Background(), capabilities.InvocationRequest{CapabilityID: fobrain.CapabilityCurrentUserContext})
	if !fobrain.HasReason(err, capabilities.PolicyReasonPermissionDenied) {
		t.Fatalf("err = %v, want permission denied", err)
	}
	if client.calls != 0 {
		t.Fatalf("client calls = %d, want blocked before client", client.calls)
	}
}

func TestStaticCredentialResolverRequiresExplicitAuthParam(t *testing.T) {
	resolver := fobrain.StaticCredentialResolver{
		WorkspaceID: "ws_fobrain",
		APIToken:    "local-secret",
	}

	_, err := resolver.ResolveFobrainCredential(context.Background(), "ws_fobrain", boundFobrainCredential("ws_fobrain"))
	if !fobrain.HasReason(err, capabilities.PolicyReasonCredentialMissing) {
		t.Fatalf("err = %v, want credential missing for absent auth param", err)
	}
}

func TestConnectorUnavailableBlocksClientCall(t *testing.T) {
	client := &recordingCurrentUserClient{}
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:       "ws_fobrain",
		CredentialBinding: boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:   capabilities.ConnectorStatusUnavailable,
		Client:            client,
	})

	_, err := provider.Invoke(context.Background(), invocationForWorkspace("ws_fobrain"))
	if !fobrain.HasReason(err, capabilities.PolicyReasonConnectorTransportUnavailable) {
		t.Fatalf("err = %v, want connector unavailable", err)
	}
	if client.calls != 0 {
		t.Fatalf("client calls = %d, want blocked before client", client.calls)
	}
}

func TestProviderErrorIsFoldedToSafeReason(t *testing.T) {
	clientErr := fobrain.NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 读取失败")
	provider := fobrain.NewProvider(fobrain.ProviderConfig{
		WorkspaceID:       "ws_fobrain",
		CredentialBinding: boundFobrainCredential("ws_fobrain"),
		ConnectorStatus:   capabilities.ConnectorStatusAvailable,
		Client:            &recordingCurrentUserClient{err: clientErr},
	})

	_, err := provider.Invoke(context.Background(), invocationForWorkspace("ws_fobrain"))
	if !fobrain.HasReason(err, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("err = %v, want connector execution failed", err)
	}
	encoded := err.Error()
	for _, forbidden := range []string{"token", "Authorization", "credential_ref", "raw provider"} {
		if strings.Contains(encoded, forbidden) {
			t.Fatalf("safe provider error leaked %q: %s", forbidden, encoded)
		}
	}
}

func TestHasReasonSupportsWrappedErrors(t *testing.T) {
	err := errors.New("plain")
	if fobrain.HasReason(err, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("plain error should not match provider reason")
	}
	wrapped := fobrain.WrapSafeError(fobrain.NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 读取失败"))
	if !fobrain.HasReason(wrapped, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("wrapped provider error reason not detected: %v", wrapped)
	}
}
