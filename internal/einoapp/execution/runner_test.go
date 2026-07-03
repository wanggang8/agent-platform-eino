package execution_test

import (
	"context"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/llm"
	"agent-platform-eino/internal/einoapp/observability"
)

func TestChatModelRunnerWritesContextAssistantAndLifecycleFacts(t *testing.T) {
	// ChatModelRunner 通过 Eino Runner 执行，但对外只写 Product Facts，不暴露 Eino event。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	now := time.Unix(500, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-chat", WorkspaceID: "ws-chat", Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendTurn(ctx, facts.Turn{TurnID: "turn-user", RunID: "run-chat", Role: facts.TurnRoleUser, Content: "你好", Sequence: 1, CreatedAt: now}); err != nil {
		t.Fatal(err)
	}

	runner := execution.NewChatModelRunner(repository, llm.NewMockProvider("你好，已收到。"), llm.Config{
		Provider: "mock",
		Model:    "mock-chat",
	}, execution.ChatModelRunnerConfig{
		Now: func() time.Time { return now.Add(time.Second) },
	})

	if err := runner.Run(ctx, "run-chat"); err != nil {
		t.Fatal(err)
	}
	snapshot, err := repository.GetSnapshot(ctx, "run-chat")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusSucceeded {
		t.Fatalf("run status = %q", snapshot.Run.Status)
	}
	if len(snapshot.ContextSnapshots) != 1 {
		t.Fatalf("context snapshots = %+v", snapshot.ContextSnapshots)
	}
	if len(snapshot.Turns) != 2 || snapshot.Turns[1].Role != facts.TurnRoleAssistant || snapshot.Turns[1].Content != "你好，已收到。" {
		t.Fatalf("assistant turn not written from runner: %+v", snapshot.Turns)
	}
	if len(snapshot.AuditEvents) == 0 || snapshot.AuditEvents[0].EventType != "message" {
		t.Fatalf("runner audit not written: %+v", snapshot.AuditEvents)
	}
}

func TestChatModelRunnerRecordsSafeTelemetry(t *testing.T) {
	// runner telemetry 是内部诊断材料，必须记录模型调用摘要但不能进入 Product Facts 或泄漏输入内容。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	sink := observability.NewMemorySink()
	now := time.Unix(700, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{RunID: "run-telemetry", WorkspaceID: "ws-chat", Status: facts.RunStatusCreated, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendTurn(ctx, facts.Turn{TurnID: "turn-user", RunID: "run-telemetry", Role: facts.TurnRoleUser, Content: "不要进入 telemetry 的用户输入", Sequence: 1, CreatedAt: now}); err != nil {
		t.Fatal(err)
	}

	runner := execution.NewChatModelRunner(repository, llm.NewMockProvider("安全回复"), llm.Config{
		Provider:   "mock",
		Model:      "mock-chat",
		ModelLabel: "mock-chat",
	}, execution.ChatModelRunnerConfig{
		Now:       func() time.Time { return now.Add(2 * time.Second) },
		Telemetry: sink,
	})

	if err := runner.Run(ctx, "run-telemetry"); err != nil {
		t.Fatal(err)
	}
	events := sink.Events()
	if len(events) != 1 {
		t.Fatalf("telemetry events = %+v", events)
	}
	if events[0].RunID != "run-telemetry" ||
		events[0].WorkspaceID != "ws-chat" ||
		events[0].OperationName != "chat.model.generate" ||
		events[0].Provider != "mock" ||
		events[0].Model != "mock-chat" ||
		events[0].LatencyMS <= 0 {
		t.Fatalf("telemetry event = %+v", events[0])
	}
	snapshot, err := repository.GetSnapshot(ctx, "run-telemetry")
	if err != nil {
		t.Fatal(err)
	}
	for _, audit := range snapshot.AuditEvents {
		if audit.EventType == "telemetry" {
			t.Fatalf("telemetry must not be Product Facts audit: %+v", snapshot.AuditEvents)
		}
	}
}
