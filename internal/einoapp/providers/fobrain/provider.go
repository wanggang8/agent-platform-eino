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

// ListCapabilities 返回当前 provider client 可执行的 Fobrain capability。
func (provider *Provider) ListCapabilities() ([]capabilities.Capability, error) {
	catalog := providerCatalog(provider.supportsParameterizedQuery(), provider.supportsDetailRisk())
	out := make([]capabilities.Capability, len(catalog))
	copy(out, catalog)
	return out, nil
}

// Invoke 执行 Fobrain 能力。provider 自身再次执行 policy/credential 防线。
func (provider *Provider) Invoke(ctx context.Context, request capabilities.InvocationRequest) (product.StructuredResultCandidate, error) {
	capability, ok := provider.capabilityByID(request.CapabilityID)
	if !ok {
		return product.StructuredResultCandidate{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain capability 未注册")
	}
	policyContext, err := provider.policyContextForInvocation(request)
	if err != nil {
		return product.StructuredResultCandidate{}, err
	}
	decision := capabilities.EvaluatePolicy(capability, policyContext)
	if !decision.Allowed {
		return product.StructuredResultCandidate{}, NewSafeError(decision.ReasonCode, decision.SafeSummary)
	}
	if request.CapabilityID == CapabilityConnectorSecurity {
		return BuildConnectorSecurityStructuredResult(policyContext, provider.connectorStatus()), nil
	}

	credential, err := provider.resolveCredential(ctx, policyContext)
	if err != nil {
		return product.StructuredResultCandidate{}, err
	}
	client := provider.config.Client
	if client == nil {
		return product.StructuredResultCandidate{}, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain client 未配置")
	}
	switch request.CapabilityID {
	case CapabilityCurrentUserContext:
		result, err := client.CurrentUserContext(ctx, credential)
		if err != nil {
			return product.StructuredResultCandidate{}, foldProviderError(err)
		}
		candidate, _ := BuildCurrentUserStructuredResult(result)
		return candidate, nil
	case CapabilityMyPermissions:
		result, err := client.MyPermissions(ctx, credential)
		if err != nil {
			return product.StructuredResultCandidate{}, foldProviderError(err)
		}
		candidate, _ := BuildMyPermissionsStructuredResult(result)
		return candidate, nil
	default:
		if batchDParameterizedCapabilityByID(request.CapabilityID) {
			parameterizedClient, ok := client.(ParameterizedQueryClient)
			if !ok {
				return product.StructuredResultCandidate{}, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 参数化查询 client 未配置")
			}
			query, err := parameterizedQueryFromArguments(request.CapabilityID, request.Arguments)
			if err != nil {
				return product.StructuredResultCandidate{}, err
			}
			result, err := parameterizedClient.ParameterizedQuery(ctx, credential, request.CapabilityID, query)
			if err != nil {
				return product.StructuredResultCandidate{}, foldProviderError(err)
			}
			// Product Facts 只能使用平台已校验的安全查询条件，不能信任 client 回填的 raw 查询。
			result.ToolID = request.CapabilityID
			result.Query = query
			candidate, _ := BuildParameterizedQueryStructuredResult(result)
			return candidate, nil
		}
		if batchEDetailRiskCapabilityByID(request.CapabilityID) {
			detailClient, ok := client.(DetailRiskClient)
			if !ok {
				return product.StructuredResultCandidate{}, NewSafeError(capabilities.PolicyReasonConnectorTransportUnavailable, "Fobrain 详情风险 client 未配置")
			}
			return provider.invokeBatchEDetailRisk(ctx, detailClient, credential, request)
		}
		return product.StructuredResultCandidate{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain capability 未注册")
	}
}

// invokeBatchEDetailRisk 统一执行 Batch E 参数校验、client 调用和 StructuredResult 构建。
func (provider *Provider) invokeBatchEDetailRisk(ctx context.Context, client DetailRiskClient, credential ResolvedCredential, request capabilities.InvocationRequest) (product.StructuredResultCandidate, error) {
	var result DetailRiskResult
	var err error
	switch request.CapabilityID {
	case CapabilityGetAssetDetail:
		query, parseErr := assetDetailQueryFromArguments(request.Arguments)
		if parseErr != nil {
			return product.StructuredResultCandidate{}, parseErr
		}
		result, err = client.AssetDetail(ctx, credential, query)
		result.Target = query.AssetID
	case CapabilityGetVulnerabilityDetail:
		query, parseErr := vulnerabilityDetailQueryFromArguments(request.Arguments)
		if parseErr != nil {
			return product.StructuredResultCandidate{}, parseErr
		}
		result, err = client.VulnerabilityDetail(ctx, credential, query)
		result.Target = query.VulnerabilityID
	case CapabilityBusinessRiskSummary:
		query, parseErr := businessRiskQueryFromArguments(request.Arguments)
		if parseErr != nil {
			return product.StructuredResultCandidate{}, parseErr
		}
		result, err = client.BusinessRiskSummary(ctx, credential, query)
		result.Target = query.BusinessName
	case CapabilityThreatRelevanceList:
		query, parseErr := threatRelevanceQueryFromArguments(request.Arguments)
		if parseErr != nil {
			return product.StructuredResultCandidate{}, parseErr
		}
		result, err = client.ThreatRelevanceList(ctx, credential, query)
		result.Target = query.VulnerabilityName
	default:
		return product.StructuredResultCandidate{}, NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain 详情风险能力未注册")
	}
	if err != nil {
		return product.StructuredResultCandidate{}, foldProviderError(err)
	}
	result.ToolID = request.CapabilityID
	candidate, _ := BuildDetailRiskStructuredResult(result)
	return candidate, nil
}

// capabilityByID 从 provider catalog 查找能力元数据，避免 execution/httpapi 写 provider 分支。
func (provider *Provider) capabilityByID(capabilityID string) (capabilities.Capability, bool) {
	for _, capability := range providerCatalog(provider.supportsParameterizedQuery(), provider.supportsDetailRisk()) {
		if capability.ID == capabilityID {
			return capability, true
		}
	}
	return capabilities.Capability{}, false
}

func (provider *Provider) supportsParameterizedQuery() bool {
	if provider == nil || provider.config.Client == nil {
		return false
	}
	_, ok := provider.config.Client.(ParameterizedQueryClient)
	return ok
}

func (provider *Provider) supportsDetailRisk() bool {
	if provider == nil || provider.config.Client == nil {
		return false
	}
	_, ok := provider.config.Client.(DetailRiskClient)
	return ok
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
