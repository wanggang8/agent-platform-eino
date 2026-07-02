package capabilities

import "time"

// RiskLevel 表示能力调用的风险域，取值与 capability catalog schema 对齐。
type RiskLevel string

const (
	RiskNone   RiskLevel = "none"
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"

	// RiskReadOnly/RiskWrite 是旧测试与本地 mock 的兼容别名，值仍使用新契约枚举。
	RiskReadOnly RiskLevel = RiskLow
	RiskWrite    RiskLevel = RiskHigh
)

// SideEffect 描述能力是否触达外部系统或本地运行时，用于审批和审计策略。
type SideEffect string

const (
	SideEffectNone          SideEffect = "none"
	SideEffectReadExternal  SideEffect = "read_external"
	SideEffectWriteExternal SideEffect = "write_external"
	SideEffectLocalRuntime  SideEffect = "local_runtime"
)

// PermissionScope 标记策略按 workspace、caller 还是 system 维度授权。
type PermissionScope string

const (
	PermissionScopeWorkspace PermissionScope = "workspace"
	PermissionScopeCaller    PermissionScope = "caller"
	PermissionScopeSystem    PermissionScope = "system"
)

// CredentialBindingPolicy 定义能力调用是否需要凭据绑定。
type CredentialBindingPolicy string

const (
	CredentialBindingNone     CredentialBindingPolicy = "none"
	CredentialBindingRequired CredentialBindingPolicy = "required"
	CredentialBindingOptional CredentialBindingPolicy = "optional"
)

// Capability 是 provider 注册到平台的能力元数据。
// execution 只能通过该元数据和 policy 选择工具，不允许硬编码工具名分支。
type Capability struct {
	ID                      string
	ProviderID              string
	ToolName                string
	DisplayName             string
	Description             string
	InputSchema             JSONSchema
	ResultSchema            string
	RiskLevel               RiskLevel
	SideEffect              SideEffect
	PolicyRef               string
	PermissionScope         PermissionScope
	CredentialBindingPolicy CredentialBindingPolicy
	ConnectorID             string
	ApprovalRequired        bool
	IdempotencyRequired     bool
	Timeout                 time.Duration
}

// JSONSchema 是 Phase 1/3 的轻量 schema 描述，后续会映射到 Eino tool schema。
type JSONSchema struct {
	SchemaVersion string
	Properties    map[string]string
	Required      []string
}

// Provider 定义能力提供方的注册入口，具体业务 provider 不能反向依赖 execution/httpapi。
type Provider interface {
	ID() string
	ListCapabilities() ([]Capability, error)
}
