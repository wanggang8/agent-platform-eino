package execution_test

import (
	"context"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
)

func TestContextProjectionUsesSafeProductFactsOnly(t *testing.T) {
	// ContextProjector 只能从 Product Facts 的安全摘要构造 LLM 输入，不能读取 provider raw payload。
	repository := facts.NewMemoryRepository()
	projector := execution.NewContextProjector(repository, execution.ContextProjectorConfig{
		Now: func() time.Time { return time.Unix(400, 0).UTC() },
	})
	ctx := context.Background()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-1", WorkspaceID: "ws-1", Status: facts.RunStatusRunning}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendTurn(ctx, facts.Turn{TurnID: "turn-1", RunID: "run-1", Role: facts.TurnRoleUser, Content: "查询资产", Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendToolCall(ctx, facts.ToolCall{ToolCallID: "call-1", RunID: "run-1", ToolID: "tool.read", DisplayName: "只读查询", Status: facts.ToolCallSucceeded}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendToolResult(ctx, facts.ToolResult{
		ResultID:   "result-1",
		ToolCallID: "call-1",
		Status:     facts.ToolResultSucceeded,
		StructuredResult: facts.StructuredResultRef{
			SchemaVersion: "tool.structured_result.v1",
			ResultRef:     "result:call-1",
			SafeSummary:   "发现 1 个资产",
		},
	}); err != nil {
		t.Fatal(err)
	}

	messages, err := projector.Project(ctx, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 2 {
		t.Fatalf("messages = %+v", messages)
	}
	if messages[0].Role != "user" || messages[0].Content != "查询资产" {
		t.Fatalf("turn context mismatch: %+v", messages[0])
	}
	if messages[1].Role != "system" || messages[1].Content != "StructuredResult: 发现 1 个资产" {
		t.Fatalf("structured result context mismatch: %+v", messages[1])
	}

	snapshot, err := repository.GetSnapshot(ctx, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.ContextSnapshots) != 1 || snapshot.ContextSnapshots[0].SafeSummary == "" {
		t.Fatalf("context snapshot not persisted: %+v", snapshot.ContextSnapshots)
	}
}
