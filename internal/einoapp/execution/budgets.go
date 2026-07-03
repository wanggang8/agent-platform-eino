package execution

// BudgetLimits 是 execution 层的预算阈值快照，来源应由配置或 policy 显式注入，不能在业务逻辑中硬编码。
type BudgetLimits struct {
	MaxModelCallsPerRun  int
	MaxToolCallsPerRun   int
	MaxInputTokensPerRun int
}

// BudgetUsage 是一次 run 当前已知的安全计数，不包含 prompt、provider payload 或 tool args。
type BudgetUsage struct {
	ModelCalls  int
	ToolCalls   int
	InputTokens int
}

// BudgetDecision 是预算评估结果；Exceeded=true 时调用方必须通过 lifecycle 写入 Product Facts。
type BudgetDecision struct {
	Exceeded    bool
	Reason      string
	SafeSummary string
}

// Evaluate 根据安全计数判断是否超预算，不直接修改 run、tool 或 pending 状态。
func (limits BudgetLimits) Evaluate(usage BudgetUsage) BudgetDecision {
	if limits.MaxModelCallsPerRun > 0 && usage.ModelCalls > limits.MaxModelCallsPerRun {
		return budgetExceeded("model_calls_exceeded")
	}
	if limits.MaxToolCallsPerRun > 0 && usage.ToolCalls > limits.MaxToolCallsPerRun {
		return budgetExceeded("tool_calls_exceeded")
	}
	if limits.MaxInputTokensPerRun > 0 && usage.InputTokens > limits.MaxInputTokensPerRun {
		return budgetExceeded("input_tokens_exceeded")
	}
	return BudgetDecision{}
}

func budgetExceeded(reason string) BudgetDecision {
	return BudgetDecision{
		Exceeded:    true,
		Reason:      reason,
		SafeSummary: "budget exceeded",
	}
}
