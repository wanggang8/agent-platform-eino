package fobrain

import (
	"context"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/product"
)

// ProviderConfig 保存 Fobrain provider 的安全运行配置。
// secret 只允许传给 CredentialResolver 和 FobrainClient，不能进入产品事实。
type ProviderConfig struct {
	WorkspaceID        string
	CredentialBinding  capabilities.CredentialBinding
	ConnectorStatus    capabilities.ConnectorStatus
	CredentialResolver CredentialResolver
	Client             FobrainClient
}

// Provider 将 Fobrain capability catalog 暴露给平台注册表。
type Provider struct {
	config ProviderConfig
}

var _ capabilities.Provider = (*Provider)(nil)
var _ capabilities.Invoker = (*Provider)(nil)

// NewProvider 创建 Fobrain provider。配置只在 provider 边界内使用，不进入产品事实。
func NewProvider(config ProviderConfig) *Provider {
	return &Provider{config: config}
}

// ID 返回 provider 稳定标识。
func (provider *Provider) ID() string {
	return ProviderID
}

// ListCapabilities 返回 Phase 5 允许注册的 Fobrain PoC capability。
func (provider *Provider) ListCapabilities() ([]capabilities.Capability, error) {
	_ = provider.config
	catalog := phase5Catalog()
	out := make([]capabilities.Capability, len(catalog))
	copy(out, catalog)
	return out, nil
}

// Invoke 执行 Fobrain PoC 能力。provider 自身再次执行 policy/credential 防线。
func (provider *Provider) Invoke(ctx context.Context, request capabilities.InvocationRequest) (product.StructuredResultCandidate, error) {
	if request.CapabilityID != CapabilityCurrentUserContext {
		return product.StructuredResultCandidate{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain capability 未注册")
	}
	capability := phase5Catalog()[0]
	policyContext, err := provider.policyContextForInvocation(request)
	if err != nil {
		return product.StructuredResultCandidate{}, err
	}
	decision := capabilities.EvaluatePolicy(capability, policyContext)
	if !decision.Allowed {
		return product.StructuredResultCandidate{}, NewSafeError(decision.ReasonCode, decision.SafeSummary)
	}

	credential, err := provider.resolveCredential(ctx, policyContext)
	if err != nil {
		return product.StructuredResultCandidate{}, err
	}
	client := provider.config.Client
	if client == nil {
		return product.StructuredResultCandidate{}, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain client 未配置")
	}
	result, err := client.CurrentUserContext(ctx, credential)
	if err != nil {
		return product.StructuredResultCandidate{}, foldProviderError(err)
	}
	candidate, _ := BuildCurrentUserStructuredResult(result)
	return candidate, nil
}

// policyContextForInvocation 以调用方的 run/request workspace 为准，配置只补凭据摘要和 connector 状态。
func (provider *Provider) policyContextForInvocation(request capabilities.InvocationRequest) (capabilities.PolicyContext, error) {
	policyContext := request.PolicyContext
	if policyContext.WorkspaceID == "" {
		return capabilities.PolicyContext{}, NewSafeError(capabilities.PolicyReasonPermissionDenied, "Fobrain 调用缺少工作区上下文")
	}
	if policyContext.CredentialBinding.Status == "" && policyContext.CredentialBinding.WorkspaceID == "" {
		policyContext.CredentialBinding = provider.config.CredentialBinding
	}
	if policyContext.ConnectorStatus == capabilities.ConnectorStatusUnknown {
		policyContext.ConnectorStatus = provider.connectorStatus()
	}
	return policyContext, nil
}

// connectorStatus 默认按可用处理，配置显式 unavailable 时才阻断。
func (provider *Provider) connectorStatus() capabilities.ConnectorStatus {
	if provider.config.ConnectorStatus == capabilities.ConnectorStatusUnavailable {
		return capabilities.ConnectorStatusUnavailable
	}
	return capabilities.ConnectorStatusAvailable
}

// resolveCredential 在 provider 边界内解析凭据；workspace 必须来自调用上下文。
func (provider *Provider) resolveCredential(ctx context.Context, policyContext capabilities.PolicyContext) (ResolvedCredential, error) {
	if provider.config.CredentialResolver == nil {
		return ResolvedCredential{WorkspaceID: policyContext.WorkspaceID}, nil
	}
	return provider.config.CredentialResolver.ResolveFobrainCredential(ctx, policyContext.WorkspaceID, policyContext.CredentialBinding)
}
