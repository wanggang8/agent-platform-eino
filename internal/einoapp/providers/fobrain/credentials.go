package fobrain

import (
	"context"

	"agent-platform-eino/internal/einoapp/capabilities"
)

// ResolvedCredential 是 provider client 内部使用的凭据材料。
// APIToken 不得进入 Product Facts、日志、报告或 Workbench 投影。
type ResolvedCredential struct {
	WorkspaceID string
	APIToken    string
}

// CredentialResolver 按 workspace 解析 Fobrain 凭据，防止跨工作区复用 token。
type CredentialResolver interface {
	ResolveFobrainCredential(context.Context, string, capabilities.CredentialBinding) (ResolvedCredential, error)
}

// StaticCredentialResolver 是本地配置和测试使用的最小凭据解析器。
type StaticCredentialResolver struct {
	WorkspaceID string
	APIToken    string
}

// ResolveFobrainCredential 只在 workspace 匹配时返回 provider client 私有凭据。
func (resolver StaticCredentialResolver) ResolveFobrainCredential(_ context.Context, workspaceID string, binding capabilities.CredentialBinding) (ResolvedCredential, error) {
	if resolver.WorkspaceID != "" && resolver.WorkspaceID != workspaceID {
		return ResolvedCredential{}, NewSafeError(capabilities.PolicyReasonCredentialScopeDenied, "Fobrain 凭据不属于当前工作区")
	}
	if binding.WorkspaceID != "" && binding.WorkspaceID != workspaceID {
		return ResolvedCredential{}, NewSafeError(capabilities.PolicyReasonCredentialScopeDenied, "Fobrain 凭据不属于当前工作区")
	}
	if resolver.APIToken == "" {
		return ResolvedCredential{}, NewSafeError(capabilities.PolicyReasonCredentialMissing, "Fobrain 凭据未配置")
	}
	return ResolvedCredential{WorkspaceID: workspaceID, APIToken: resolver.APIToken}, nil
}
