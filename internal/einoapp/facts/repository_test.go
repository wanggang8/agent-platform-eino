package facts_test

import (
	"context"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
)

func TestRunFactHasStableProductIdentity(t *testing.T) {
	run := facts.Run{
		RunID:       "run-1",
		WorkspaceID: "ws-1",
		Status:      facts.RunStatusCreated,
		CreatedAt:   time.Unix(1, 0).UTC(),
		UpdatedAt:   time.Unix(2, 0).UTC(),
	}

	if run.RunID == "" || run.WorkspaceID == "" {
		t.Fatalf("run identity must be present: %+v", run)
	}
	if run.Status != facts.RunStatusCreated {
		t.Fatalf("status = %q", run.Status)
	}
}

func TestRunStatusMatchesProductFactsContract(t *testing.T) {
	statuses := []facts.RunStatus{
		facts.RunStatusCreated,
		facts.RunStatusRunning,
		facts.RunStatusWaiting,
		facts.RunStatusSucceeded,
		facts.RunStatusFailed,
		facts.RunStatusCancelled,
		facts.RunStatusStopped,
	}

	for _, status := range statuses {
		if !status.Valid() {
			t.Fatalf("status %q must be valid", status)
		}
	}
	if facts.RunStatus("done").Valid() {
		t.Fatal("done must not be a valid run status")
	}
	if !facts.RunStatusSucceeded.Terminal() || !facts.RunStatusCancelled.Terminal() || !facts.RunStatusStopped.Terminal() {
		t.Fatal("terminal run statuses must be terminal")
	}
	if facts.RunStatusRunning.Terminal() || facts.RunStatusWaiting.Terminal() {
		t.Fatal("active run statuses must not be terminal")
	}
}

func TestToolResultFactRequiresStructuredResultOnly(t *testing.T) {
	// 工具结果测试只断言 StructuredResult 引用，避免把 raw provider payload 固化进契约。
	result := facts.ToolResult{
		ResultID:   "result-1",
		ToolCallID: "call-1",
		Status:     facts.ToolResultSucceeded,
		StructuredResult: facts.StructuredResultRef{
			SchemaVersion: "tool.structured_result.v1",
			ResultRef:     "result:call-1",
			SafeSummary:   "查询完成",
		},
	}

	if result.StructuredResult.SchemaVersion != "tool.structured_result.v1" {
		t.Fatalf("structured result schema = %q", result.StructuredResult.SchemaVersion)
	}
}

func TestFactsCarrySequenceAndSafeResumeReferences(t *testing.T) {
	// sequence、resume_ref、checkpoint_ref 是 replay/resume 的公共事实边界。
	turn := facts.Turn{
		TurnID:   "turn-1",
		RunID:    "run-1",
		Role:     facts.TurnRoleAssistant,
		Content:  "完成",
		Sequence: 7,
	}
	if turn.Sequence != 7 {
		t.Fatalf("sequence = %d", turn.Sequence)
	}

	pending := facts.PendingInteraction{
		PendingID:     "pending-1",
		RunID:         "run-1",
		Kind:          facts.PendingKindApproval,
		Status:        facts.PendingStatusWaiting,
		ResumeRef:     "resume-safe-1",
		CheckpointRef: "checkpoint-safe-1",
	}
	if pending.ResumeRef == "" || pending.CheckpointRef == "" {
		t.Fatalf("pending safe refs must be present: %+v", pending)
	}
	if pending.CheckpointRef == "checkpoint-raw-eino-id" {
		t.Fatal("checkpoint ref must not expose a raw Eino checkpoint id")
	}
}

func TestMemoryRepositoryTracksLatestRunByWorkspace(t *testing.T) {
	// 内存仓库是 HTTP 默认依赖和早期测试替身，必须和 SQLite 仓库保持同源语义。
	repository := facts.NewMemoryRepository()
	ctx := context.Background()

	first := facts.Run{RunID: "run-1", WorkspaceID: "ws-1", Status: facts.RunStatusCreated}
	second := facts.Run{RunID: "run-2", WorkspaceID: "ws-1", Status: facts.RunStatusRunning}
	if err := repository.CreateRun(ctx, first); err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateRun(ctx, second); err != nil {
		t.Fatal(err)
	}

	latest, err := repository.LatestRun(ctx, "ws-1")
	if err != nil {
		t.Fatal(err)
	}
	if latest.RunID != "run-2" {
		t.Fatalf("latest run = %q", latest.RunID)
	}
}
