package execution_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func TestRunnerEventMapperWritesAssistantDeltaToFacts(t *testing.T) {
	repository := facts.NewMemoryRepository()
	mapper := execution.NewEventMapper(repository, execution.EventMapperConfig{
		Now: func() time.Time { return time.Unix(100, 0).UTC() },
	})
	ctx := context.Background()
	run := facts.Run{RunID: "run-1", WorkspaceID: "ws-1", Status: facts.RunStatusRunning}
	if err := repository.CreateRun(ctx, run); err != nil {
		t.Fatal(err)
	}

	// assistant 事件必须先进入 Product Facts，不能直接拼 Workbench 消息。
	if err := mapper.Map(ctx, execution.RunnerEvent{
		RunID:    "run-1",
		Kind:     execution.RunnerEventAssistantMessage,
		Content:  "完成",
		Sequence: 3,
	}); err != nil {
		t.Fatal(err)
	}

	snapshot, err := repository.GetSnapshot(ctx, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Turns) != 1 || snapshot.Turns[0].Role != facts.TurnRoleAssistant || snapshot.Turns[0].Content != "完成" || snapshot.Turns[0].Sequence != 3 {
		t.Fatalf("assistant message not mapped to facts: %+v", snapshot.Turns)
	}
}

func TestRunnerEventMapperWritesToolResultAsStructuredResultFact(t *testing.T) {
	repository := facts.NewMemoryRepository()
	mapper := execution.NewEventMapper(repository, execution.EventMapperConfig{
		Now: func() time.Time { return time.Unix(200, 0).UTC() },
	})
	ctx := context.Background()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-1", WorkspaceID: "ws-1", Status: facts.RunStatusRunning}); err != nil {
		t.Fatal(err)
	}

	// 工具调用和结果分两步写入，结果只保留 StructuredResult 安全引用。
	if err := mapper.Map(ctx, execution.RunnerEvent{
		RunID:       "run-1",
		Kind:        execution.RunnerEventToolCall,
		ToolCallID:  "call-1",
		ToolID:      "tool.read",
		DisplayName: "只读查询",
		Status:      "running",
		ArgsPreview: "target=example",
	}); err != nil {
		t.Fatal(err)
	}
	if err := mapper.Map(ctx, execution.RunnerEvent{
		RunID:       "run-1",
		Kind:        execution.RunnerEventToolResult,
		ToolCallID:  "call-1",
		ResultID:    "result-1",
		Status:      "succeeded",
		ResultRef:   "result:call-1",
		SafeSummary: "查询完成",
	}); err != nil {
		t.Fatal(err)
	}

	snapshot, err := repository.GetSnapshot(ctx, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.ToolCalls) != 1 || snapshot.ToolCalls[0].ToolID != "tool.read" {
		t.Fatalf("tool call not mapped: %+v", snapshot.ToolCalls)
	}
	if len(snapshot.ToolResults) != 1 || snapshot.ToolResults[0].StructuredResult.SchemaVersion != "tool.structured_result.v1" || snapshot.ToolResults[0].StructuredResult.SafeSummary != "查询完成" {
		t.Fatalf("tool result not mapped as structured result fact: %+v", snapshot.ToolResults)
	}
}

func TestRunnerEventMapperRejectsUnsafeAssistantAndToolResult(t *testing.T) {
	repository := facts.NewMemoryRepository()
	mapper := execution.NewEventMapper(repository, execution.EventMapperConfig{
		Now: func() time.Time { return time.Unix(250, 0).UTC() },
	})
	ctx := context.Background()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-1", WorkspaceID: "ws-1", Status: facts.RunStatusRunning}); err != nil {
		t.Fatal(err)
	}
	if err := mapper.Map(ctx, execution.RunnerEvent{
		RunID:    "run-1",
		Kind:     execution.RunnerEventAssistantMessage,
		Content:  "Authorization: Bearer local-secret",
		Sequence: 1,
	}); !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe assistant err = %v", err)
	}

	if err := mapper.Map(ctx, execution.RunnerEvent{
		RunID:       "run-1",
		Kind:        execution.RunnerEventToolCall,
		ToolCallID:  "call-1",
		ToolID:      "tool.read",
		DisplayName: "只读查询",
		Status:      "running",
		ArgsPreview: "target=example",
	}); err != nil {
		t.Fatal(err)
	}
	if err := mapper.Map(ctx, execution.RunnerEvent{
		RunID:       "run-1",
		Kind:        execution.RunnerEventToolResult,
		ToolCallID:  "call-1",
		ResultID:    "result-1",
		Status:      "succeeded",
		ResultRef:   "result:call-1",
		SafeSummary: "raw provider body: {...}",
	}); !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe structured result err = %v", err)
	}
	if err := mapper.Map(ctx, execution.RunnerEvent{
		RunID:         "run-1",
		Kind:          execution.RunnerEventPending,
		PendingID:     "pending-1",
		PendingKind:   "approval",
		PendingStatus: "waiting",
		ResumeRef:     "resume-safe-1",
		CheckpointRef: "checkpoint-raw-eino-id",
	}); !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe pending err = %v", err)
	}

	snapshot, err := repository.GetSnapshot(ctx, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Turns) != 0 || len(snapshot.ToolResults) != 0 || len(snapshot.PendingInteractions) != 0 {
		t.Fatalf("unsafe material must not be written: turns=%+v results=%+v pending=%+v", snapshot.Turns, snapshot.ToolResults, snapshot.PendingInteractions)
	}
}

func TestRunnerEventMapperWritesPendingAndLifecycleFacts(t *testing.T) {
	repository := facts.NewMemoryRepository()
	mapper := execution.NewEventMapper(repository, execution.EventMapperConfig{
		Now: func() time.Time { return time.Unix(300, 0).UTC() },
	})
	ctx := context.Background()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-1", WorkspaceID: "ws-1", Status: facts.RunStatusRunning}); err != nil {
		t.Fatal(err)
	}

	// pending 和 run lifecycle 都必须落 Product Facts，供 resume/audit/replay 同源读取。
	if err := mapper.Map(ctx, execution.RunnerEvent{
		RunID:         "run-1",
		Kind:          execution.RunnerEventPending,
		PendingID:     "pending-1",
		PendingKind:   "approval",
		PendingStatus: "waiting",
		ResumeRef:     "resume-safe-1",
		CheckpointRef: "checkpoint-safe-1",
		RiskSummary:   "需要审批",
	}); err != nil {
		t.Fatal(err)
	}
	if err := mapper.Map(ctx, execution.RunnerEvent{
		RunID:  "run-1",
		Kind:   execution.RunnerEventLifecycle,
		Status: string(facts.RunStatusWaiting),
	}); err != nil {
		t.Fatal(err)
	}

	snapshot, err := repository.GetSnapshot(ctx, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusWaiting {
		t.Fatalf("run status = %q", snapshot.Run.Status)
	}
	if len(snapshot.PendingInteractions) != 1 || snapshot.PendingInteractions[0].ResumeRef != "resume-safe-1" {
		t.Fatalf("pending not mapped: %+v", snapshot.PendingInteractions)
	}
}

func TestEinoAgentEventConvertsAssistantOutputToRunnerEvent(t *testing.T) {
	// Eino AgentEvent shape 只在 execution 内转换，不进入 HTTP、product 或前端契约。
	event := adk.EventFromMessage(schema.AssistantMessage("完成", nil), nil, schema.Assistant, "")

	events := execution.EinoAgentEventToRunnerEvents("run-1", 7, event)
	if len(events) != 1 {
		t.Fatalf("events = %+v", events)
	}
	if events[0].Kind != execution.RunnerEventAssistantMessage || events[0].Content != "完成" || events[0].Sequence != 7 {
		t.Fatalf("runner event mismatch: %+v", events[0])
	}
}
