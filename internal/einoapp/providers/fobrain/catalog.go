package fobrain

import (
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
)

// ProviderID 是 Fobrain provider 在新 capability registry 中的稳定标识。
const ProviderID = "fobrain"

// ConnectorID 是 Fobrain connector 的安全展示标识，不包含任何真实连接参数。
const ConnectorID = "fobrain"

// CapabilityCurrentUserContext 是 Phase 5 只读 PoC 能力。
const CapabilityCurrentUserContext = "tool.fobrain.current_user_context"

// PolicyRead 是 Fobrain 只读能力的稳定策略引用。
const PolicyRead = "policy:fobrain:read:v1"

// BusinessResultSchemaVersion 是 Fobrain 业务展示 payload schema，不是 Product Facts 事实 schema。
const BusinessResultSchemaVersion = "fobrain.tool_result.v2"

const currentUserContextTimeout = 10 * time.Second

// phase5Catalog 返回当前阶段允许注册的最小 Fobrain 能力目录。
// 后续恢复 24 个只读工具时只能扩展该目录数据，不能在 execution/httpapi 写工具名分支。
func phase5Catalog() []capabilities.Capability {
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
	}
}
