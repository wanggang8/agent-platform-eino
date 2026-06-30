package capabilities

type PolicyDecision struct {
	Allowed          bool
	RequiresApproval bool
	ReasonCode       string
}

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
