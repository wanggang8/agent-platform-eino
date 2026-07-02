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

// PolicyRead 是 Fobrain 只读能力的稳定策略引用。
const PolicyRead = "policy:fobrain:read:v1"

// PolicyConnectorRead 是 Fobrain connector 状态读取的稳定策略引用。
const PolicyConnectorRead = "policy:fobrain:connector-read:v1"

// BusinessResultSchemaVersion 是 Fobrain 业务展示 payload schema，不是 Product Facts 事实 schema。
const BusinessResultSchemaVersion = "fobrain.tool_result.v2"

const currentUserContextTimeout = 10 * time.Second

// batchACatalog 返回当前阶段允许注册的 Fobrain Batch A 能力目录。
// 后续恢复 24 个只读工具时只能扩展该目录数据，不能在 execution/httpapi 写工具名分支。
func batchACatalog() []capabilities.Capability {
	return []capabilities.Capability{
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
}
