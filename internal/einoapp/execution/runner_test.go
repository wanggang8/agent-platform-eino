package execution_test

import (
	"context"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/llm"
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
