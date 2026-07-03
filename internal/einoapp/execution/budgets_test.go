package execution_test

import (
	"testing"

	"agent-platform-eino/internal/einoapp/execution"
)

func TestBudgetLimitsDetectCounterExceeded(t *testing.T) {
	// 预算评估器只产出安全决策，不直接写 Product Facts；终态落库仍走 budget_exceeded lifecycle。
	limits := execution.BudgetLimits{
		MaxModelCallsPerRun:  1,
		MaxToolCallsPerRun:   2,
		MaxInputTokensPerRun: 100,
	}

	decision := limits.Evaluate(execution.BudgetUsage{ModelCalls: 2, ToolCalls: 1, InputTokens: 20})
	if !decision.Exceeded || decision.Reason != "model_calls_exceeded" || decision.SafeSummary != "budget exceeded" {
		t.Fatalf("model call decision = %+v", decision)
	}

	decision = limits.Evaluate(execution.BudgetUsage{ModelCalls: 1, ToolCalls: 3, InputTokens: 20})
	if !decision.Exceeded || decision.Reason != "tool_calls_exceeded" {
		t.Fatalf("tool call decision = %+v", decision)
	}

	decision = limits.Evaluate(execution.BudgetUsage{ModelCalls: 1, ToolCalls: 2, InputTokens: 101})
	if !decision.Exceeded || decision.Reason != "input_tokens_exceeded" {
		t.Fatalf("input token decision = %+v", decision)
	}
}
