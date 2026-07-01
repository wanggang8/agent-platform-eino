package execution

import "agent-platform-eino/internal/einoapp/capabilities"

// SelectionMode 表示入口选择结果；普通自然语言默认进入 chat，不直接绑定工具。
type SelectionMode string

const (
	SelectionModeChat       SelectionMode = "chat"
	SelectionModeCapability SelectionMode = "capability"
	SelectionModeRejected   SelectionMode = "rejected"
)

// SelectionRequest 是执行入口的选择输入，只允许显式 hint 影响候选能力。
type SelectionRequest struct {
	InputText      string
	CapabilityHint string
	IntentHint     string
	ProductAction  string
}

// SelectionResult 记录 registry/policy 之后的选择结果，供后续写入 Product Facts。
type SelectionResult struct {
	Mode             SelectionMode
	CapabilityID     string
	RequiresApproval bool
	PolicyReason     string
}

// SelectCapability 只通过 Capability Registry 和 policy 处理 hint。
// 它不会根据自然语言关键词或 Fobrain 工具名硬编码选择业务工具。
func SelectCapability(registry *capabilities.Registry, request SelectionRequest) SelectionResult {
	if request.CapabilityHint == "" {
		return SelectionResult{
			Mode:         SelectionModeChat,
			PolicyReason: "chat_default",
		}
	}

	capability, ok := registry.Get(request.CapabilityHint)
	if !ok {
		return SelectionResult{
			Mode:         SelectionModeRejected,
			PolicyReason: "capability_not_registered",
		}
	}
	decision := capabilities.EvaluatePolicy(capability)
	return SelectionResult{
		Mode:             SelectionModeCapability,
		CapabilityID:     capability.ID,
		RequiresApproval: decision.RequiresApproval,
		PolicyReason:     decision.ReasonCode,
	}
}
