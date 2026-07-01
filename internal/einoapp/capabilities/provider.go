package capabilities

import "time"

// RiskLevel 表示能力调用的风险域，写域必须进入审批策略。
type RiskLevel string

const (
	RiskReadOnly RiskLevel = "read_only"
	RiskWrite    RiskLevel = "write"
)

// Capability 是 provider 注册到平台的能力元数据。
// execution 只能通过该元数据和 policy 选择工具，不允许硬编码工具名分支。
type Capability struct {
	ID               string
	ProviderID       string
	ToolName         string
	DisplayName      string
	Description      string
	InputSchema      JSONSchema
	ResultSchema     string
	RiskLevel        RiskLevel
	ApprovalRequired bool
	Timeout          time.Duration
}

// JSONSchema 是 Phase 1/3 的轻量 schema 描述，后续会映射到 Eino tool schema。
type JSONSchema struct {
	SchemaVersion string
	Properties    map[string]string
}

// Provider 定义能力提供方的注册入口，具体业务 provider 不能反向依赖 execution/httpapi。
type Provider interface {
	ID() string
	ListCapabilities() ([]Capability, error)
}
