package execution

import (
	"testing"

	"agent-platform-eino/internal/einoapp/llm"
)

func TestEinoTokenUsageMapsTotalTokens(t *testing.T) {
	// Eino callback usage 是内部 telemetry 输入，必须完整保留 provider 安全计数。
	usage := einoTokenUsage(llm.TokenUsage{InputTokens: 3, OutputTokens: 5, TotalTokens: 8})
	if usage == nil {
		t.Fatal("usage = nil")
	}
	if usage.PromptTokens != 3 || usage.CompletionTokens != 5 || usage.TotalTokens != 8 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestEinoTokenUsageOmitsEmptyUsage(t *testing.T) {
	// provider 未返回 usage 时不制造虚假 token 计数。
	if usage := einoTokenUsage(llm.TokenUsage{}); usage != nil {
		t.Fatalf("usage = %+v", usage)
	}
}
