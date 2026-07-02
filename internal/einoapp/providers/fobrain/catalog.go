package fobrain

import (
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
)

// ProviderID 是 Fobrain provider 在新 capability registry 中的稳定标识。
const ProviderID = "fobrain"

// ConnectorID 是 Fobrain connector 的安全展示标识，不包含任何真实连接参数。
const ConnectorID = "fobrain"

// CapabilityCurrentUserContext 是 Batch A 当前用户读取能力。
const CapabilityCurrentUserContext = "tool.fobrain.current_user_context"

// CapabilityMyPermissions 是 Batch A 权限范围读取能力。
const CapabilityMyPermissions = "tool.fobrain.my_permissions"

// CapabilityConnectorSecurity 是 Batch A connector 状态能力，不替代业务读取工具。
const CapabilityConnectorSecurity = "connector.fobrain.security"

const (
	// CapabilityListAssetsByOwner 按负责人查询资产列表。
	CapabilityListAssetsByOwner = "tool.fobrain.list_assets_by_owner"
	// CapabilityListVulnerabilitiesByOwner 按负责人查询漏洞列表。
	CapabilityListVulnerabilitiesByOwner = "tool.fobrain.list_vulnerabilities_by_owner"
	// CapabilityListAssetsByDepartment 按部门查询资产列表。
	CapabilityListAssetsByDepartment = "tool.fobrain.list_assets_by_department"
	// CapabilityListVulnerabilitiesByDepartment 按部门查询漏洞列表。
	CapabilityListVulnerabilitiesByDepartment = "tool.fobrain.list_vulnerabilities_by_department"
	// CapabilityListAssetsByIP 按 IP 查询资产列表。
	CapabilityListAssetsByIP = "tool.fobrain.list_assets_by_ip"
	// CapabilityListVulnerabilitiesByIP 按 IP 查询漏洞列表。
	CapabilityListVulnerabilitiesByIP = "tool.fobrain.list_vulnerabilities_by_ip"
)

const (
	// CapabilityGetAssetDetail 查询资产详情。
	CapabilityGetAssetDetail = "tool.fobrain.get_asset_detail"
	// CapabilityGetVulnerabilityDetail 查询漏洞详情。
	CapabilityGetVulnerabilityDetail = "tool.fobrain.get_vulnerability_detail"
	// CapabilityBusinessRiskSummary 查询业务系统风险聚合。
	CapabilityBusinessRiskSummary = "tool.fobrain.business_risk_summary"
	// CapabilityThreatRelevanceList 查询漏洞/威胁关联资产列表。
	CapabilityThreatRelevanceList = "tool.fobrain.threat_relevance_list"
)

// PolicyRead 是 Fobrain 只读能力的稳定策略引用。
const PolicyRead = "policy:fobrain:read:v1"

// PolicyConnectorRead 是 Fobrain connector 状态读取的稳定策略引用。
const PolicyConnectorRead = "policy:fobrain:connector-read:v1"

// BusinessResultSchemaVersion 是 Fobrain 业务展示 payload schema，不是 Product Facts 事实 schema。
const BusinessResultSchemaVersion = "fobrain.tool_result.v2"

const currentUserContextTimeout = 10 * time.Second

// providerCatalog 返回当前阶段允许注册的 Fobrain 能力目录。
// 后续恢复只读工具时只能扩展目录数据和 provider mapper，不能在 execution/httpapi 写工具名分支。
func providerCatalog(includeBatchD bool, includeBatchE bool) []capabilities.Capability {
	catalog := []capabilities.Capability{
		{
			ID:          CapabilityCurrentUserContext,
			ProviderID:  ProviderID,
			ToolName:    CapabilityCurrentUserContext,
			DisplayName: "读取当前 Fobrain 用户信息",
			Description: "读取当前 Fobrain 用户身份、部门和角色的安全摘要。",
			InputSchema: capabilities.JSONSchema{
				SchemaVersion: "json_schema.v1",
				Properties:    map[string]string{},
				Required:      []string{},
			},
			ResultSchema:            BusinessResultSchemaVersion,
			RiskLevel:               capabilities.RiskLow,
			SideEffect:              capabilities.SideEffectReadExternal,
			PolicyRef:               PolicyRead,
			PermissionScope:         capabilities.PermissionScopeWorkspace,
			CredentialBindingPolicy: capabilities.CredentialBindingRequired,
			ConnectorID:             ConnectorID,
			ApprovalRequired:        false,
			IdempotencyRequired:     false,
			Timeout:                 currentUserContextTimeout,
		},
		{
			ID:          CapabilityMyPermissions,
			ProviderID:  ProviderID,
			ToolName:    CapabilityMyPermissions,
			DisplayName: "读取我的 Fobrain 权限范围",
			Description: "读取当前 Fobrain 用户权限、菜单和数据范围的安全摘要。",
			InputSchema: capabilities.JSONSchema{
				SchemaVersion: "json_schema.v1",
				Properties:    map[string]string{},
				Required:      []string{},
			},
			ResultSchema:            BusinessResultSchemaVersion,
			RiskLevel:               capabilities.RiskLow,
			SideEffect:              capabilities.SideEffectReadExternal,
			PolicyRef:               PolicyRead,
			PermissionScope:         capabilities.PermissionScopeWorkspace,
			CredentialBindingPolicy: capabilities.CredentialBindingRequired,
			ConnectorID:             ConnectorID,
			ApprovalRequired:        false,
			IdempotencyRequired:     false,
			Timeout:                 currentUserContextTimeout,
		},
		{
			ID:          CapabilityConnectorSecurity,
			ProviderID:  ProviderID,
			ToolName:    CapabilityConnectorSecurity,
			DisplayName: "读取 Fobrain 连接器状态",
			Description: "读取 Fobrain connector 可用性和凭据绑定安全摘要，不读取业务数据。",
			InputSchema: capabilities.JSONSchema{
				SchemaVersion: "json_schema.v1",
				Properties:    map[string]string{},
				Required:      []string{},
			},
			ResultSchema:            BusinessResultSchemaVersion,
			RiskLevel:               capabilities.RiskLow,
			SideEffect:              capabilities.SideEffectReadExternal,
			PolicyRef:               PolicyConnectorRead,
			PermissionScope:         capabilities.PermissionScopeWorkspace,
			CredentialBindingPolicy: capabilities.CredentialBindingOptional,
			ApprovalRequired:        false,
			IdempotencyRequired:     false,
			Timeout:                 currentUserContextTimeout,
		},
	}
	if includeBatchD {
		catalog = append(catalog, batchDParameterizedCatalog()...)
	}
	if includeBatchE {
		catalog = append(catalog, batchEDetailRiskCatalog()...)
	}
	return catalog
}

func batchDParameterizedCatalog() []capabilities.Capability {
	return []capabilities.Capability{
		batchDParameterizedCapability(CapabilityListAssetsByOwner, "查询负责人资产", "按负责人查询资产列表。", map[string]string{"person_name": "string", "person_staff_id": "string", "page": "integer", "page_size": "integer", "severity": "string"}, []string{"person_name"}),
		batchDParameterizedCapability(CapabilityListVulnerabilitiesByOwner, "查询负责人漏洞", "按负责人查询漏洞列表。", map[string]string{"person_name": "string", "person_staff_id": "string", "page": "integer", "page_size": "integer", "severity": "string"}, []string{"person_name"}),
		batchDParameterizedCapability(CapabilityListAssetsByDepartment, "查询部门资产", "按部门查询资产列表。", map[string]string{"department_name": "string", "page": "integer", "page_size": "integer", "severity": "string"}, []string{"department_name"}),
		batchDParameterizedCapability(CapabilityListVulnerabilitiesByDepartment, "查询部门漏洞", "按部门查询漏洞列表。", map[string]string{"department_name": "string", "page": "integer", "page_size": "integer", "severity": "string"}, []string{"department_name"}),
		batchDParameterizedCapability(CapabilityListAssetsByIP, "查询 IP 资产", "按 IP 查询资产列表。", map[string]string{"ip": "string", "page": "integer", "page_size": "integer", "severity": "string", "status": "string"}, []string{"ip"}),
		batchDParameterizedCapability(CapabilityListVulnerabilitiesByIP, "查询 IP 漏洞", "按 IP 查询漏洞列表。", map[string]string{"ip": "string", "page": "integer", "page_size": "integer", "severity": "string", "status": "string"}, []string{"ip"}),
	}
}

func batchDParameterizedCapability(id string, displayName string, description string, properties map[string]string, required []string) capabilities.Capability {
	return capabilities.Capability{
		ID:          id,
		ProviderID:  ProviderID,
		ToolName:    id,
		DisplayName: displayName,
		Description: description,
		InputSchema: capabilities.JSONSchema{
			SchemaVersion: "json_schema.v1",
			Properties:    properties,
			Required:      required,
		},
		ResultSchema:            BusinessResultSchemaVersion,
		RiskLevel:               capabilities.RiskLow,
		SideEffect:              capabilities.SideEffectReadExternal,
		PolicyRef:               PolicyRead,
		PermissionScope:         capabilities.PermissionScopeWorkspace,
		CredentialBindingPolicy: capabilities.CredentialBindingRequired,
		ConnectorID:             ConnectorID,
		ApprovalRequired:        false,
		IdempotencyRequired:     false,
		Timeout:                 currentUserContextTimeout,
	}
}

func batchEDetailRiskCatalog() []capabilities.Capability {
	return []capabilities.Capability{
		fobrainReadCapability(CapabilityGetAssetDetail, "查询资产详情", "按资产 ID 查询资产安全详情。", map[string]string{"asset_id": "string", "network_type": "string"}, []string{"asset_id"}),
		fobrainReadCapability(CapabilityGetVulnerabilityDetail, "查询漏洞详情", "按漏洞 ID 查询漏洞安全详情。", map[string]string{"vulnerability_id": "string"}, []string{"vulnerability_id"}),
		fobrainReadCapability(CapabilityBusinessRiskSummary, "查询业务风险摘要", "按业务系统名称聚合关联漏洞风险。", map[string]string{"business_name": "string"}, []string{"business_name"}),
		fobrainReadCapability(CapabilityThreatRelevanceList, "查询威胁关联列表", "按漏洞或威胁名称查询关联资产。", map[string]string{"vulnerability_name": "string", "ip": "string", "business_name": "string", "page": "integer", "page_size": "integer"}, []string{"vulnerability_name"}),
	}
}

func fobrainReadCapability(id string, displayName string, description string, properties map[string]string, required []string) capabilities.Capability {
	return capabilities.Capability{
		ID:          id,
		ProviderID:  ProviderID,
		ToolName:    id,
		DisplayName: displayName,
		Description: description,
		InputSchema: capabilities.JSONSchema{
			SchemaVersion: "json_schema.v1",
			Properties:    properties,
			Required:      required,
		},
		ResultSchema:            BusinessResultSchemaVersion,
		RiskLevel:               capabilities.RiskLow,
		SideEffect:              capabilities.SideEffectReadExternal,
		PolicyRef:               PolicyRead,
		PermissionScope:         capabilities.PermissionScopeWorkspace,
		CredentialBindingPolicy: capabilities.CredentialBindingRequired,
		ConnectorID:             ConnectorID,
		ApprovalRequired:        false,
		IdempotencyRequired:     false,
		Timeout:                 currentUserContextTimeout,
	}
}
