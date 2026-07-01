package capabilities

// PolicyDecision 是能力选择后的安全策略结果，后续会写入 Product Facts/audit。
type PolicyDecision struct {
	Allowed          bool
	RequiresApproval bool
	ReasonCode       string
}

// EvaluatePolicy 执行最小风险策略：写域能力必须审批，只读能力可直接进入候选。
func EvaluatePolicy(capability Capability) PolicyDecision {
	if capability.RiskLevel == RiskWrite {
		return PolicyDecision{
			Allowed:          false,
			RequiresApproval: true,
			ReasonCode:       "write_requires_approval",
		}
	}
	return PolicyDecision{
		Allowed:          true,
		RequiresApproval: capability.ApprovalRequired,
		ReasonCode:       "allowed",
	}
}
