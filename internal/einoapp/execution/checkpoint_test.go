package execution_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/store/sqlite"
)

func TestResumeCheckpointMissingReturnsSafeError(t *testing.T) {
	// resume 必须先确认 checkpoint 存在；缺失时不能消费 pending 或泄漏内部 checkpoint id。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository).WithCheckpointResolver(missingCheckpointResolver{})
	createRun(t, repository, "run-missing-checkpoint", "ws-hitl", facts.RunStatusWaiting, "")
	appendPendingWithStatus(t, repository, "run-missing-checkpoint", "pending-missing", "resume-safe-missing", "checkpoint_ref:missing", facts.PendingStatusWaiting)

	_, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-missing-checkpoint",
		ResumeRef:       "resume-safe-missing",
		ClientRequestID: "client-resume-missing",
		Decision:        "approved",
	})
	if !errors.Is(err, execution.ErrCheckpointMissing) {
		t.Fatalf("resume err = %v, want ErrCheckpointMissing", err)
	}
	snapshot, err := repository.GetSnapshot(ctx, "run-missing-checkpoint")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusWaiting || snapshot.PendingInteractions[0].Status != facts.PendingStatusWaiting {
		t.Fatalf("missing checkpoint should not consume pending: %+v", snapshot)
	}
}

func TestResumeAfterRestartResolvesCheckpointBeforeAccepting(t *testing.T) {
	// checkpoint store 和 Product Facts 都重启后，resume 仍能通过安全 checkpoint_ref 找到内部 checkpoint。
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "hitl.db")
	now := time.Unix(710, 0).UTC()

	repository, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	checkpoints, err := sqlite.OpenCheckpointStore(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	createRun(t, repository, "run-restart", "ws-hitl", facts.RunStatusWaiting, "")
	appendPendingWithStatus(t, repository, "run-restart", "pending-restart", "resume-safe-restart", "checkpoint_ref:restart", facts.PendingStatusWaiting)
	if err := checkpoints.Set(ctx, "eino-internal-checkpoint-restart", []byte("serialized")); err != nil {
		t.Fatal(err)
	}
	if err := checkpoints.BindCheckpointRef(ctx, sqlite.CheckpointBinding{
		CheckpointRef: "checkpoint_ref:restart",
		CheckpointID:  "eino-internal-checkpoint-restart",
		RunID:         "run-restart",
		PendingID:     "pending-restart",
		CreatedAt:     now,
	}); err != nil {
		t.Fatal(err)
	}
	if err := checkpoints.Close(); err != nil {
		t.Fatal(err)
	}
	if err := repository.Close(); err != nil {
		t.Fatal(err)
	}

	reopenedRepository, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer reopenedRepository.Close()
	reopenedCheckpoints, err := sqlite.OpenCheckpointStore(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer reopenedCheckpoints.Close()

	accepted, err := execution.NewFactCommands(reopenedRepository).
		WithCheckpointResolver(reopenedCheckpoints).
		Resume(ctx, execution.ResumeCommand{
			WorkspaceID:     "ws-hitl",
			RunID:           "run-restart",
			ResumeRef:       "resume-safe-restart",
			ClientRequestID: "client-resume-restart",
			Decision:        "approved",
		})
	if err != nil {
		t.Fatal(err)
	}
	if accepted.RunID != "run-restart" || accepted.Status != "accepted" {
		t.Fatalf("accepted = %+v", accepted)
	}
}

func TestResumeRejectsNonWaitingRunOrTerminalPending(t *testing.T) {
	// checkpoint 存在也不能让非等待态 run 或终态 pending 被接受。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	resolver := fixedCheckpointResolver{checkpointID: "eino-internal-checkpoint-1"}
	commands := execution.NewFactCommands(repository).WithCheckpointResolver(resolver)

	createRun(t, repository, "run-created", "ws-hitl", facts.RunStatusCreated, "")
	appendPendingWithStatus(t, repository, "run-created", "pending-created", "resume-safe-created", "checkpoint_ref:created", facts.PendingStatusWaiting)
	if _, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-created",
		ResumeRef:       "resume-safe-created",
		ClientRequestID: "client-resume-created",
	}); !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("created run resume err = %v, want ErrResumeNotAllowed", err)
	}

	createRun(t, repository, "run-expired", "ws-hitl", facts.RunStatusWaiting, "")
	appendPendingWithStatus(t, repository, "run-expired", "pending-expired", "resume-safe-expired", "checkpoint_ref:expired", facts.PendingStatusExpired)
	if _, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-expired",
		ResumeRef:       "resume-safe-expired",
		ClientRequestID: "client-resume-expired",
	}); !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("expired pending resume err = %v, want ErrResumeNotAllowed", err)
	}
}

type missingCheckpointResolver struct{}

func (missingCheckpointResolver) ResolveCheckpointID(context.Context, string, string, string) (string, bool, error) {
	return "", false, nil
}

type fixedCheckpointResolver struct {
	checkpointID string
}

func (resolver fixedCheckpointResolver) ResolveCheckpointID(context.Context, string, string, string) (string, bool, error) {
	return resolver.checkpointID, true, nil
}

func appendWaitingPending(t *testing.T, repository facts.Repository, runID string, pendingID string, resumeRef string, checkpointRef string) {
	t.Helper()
	appendPendingWithStatus(t, repository, runID, pendingID, resumeRef, checkpointRef, facts.PendingStatusWaiting)
}

func appendPendingWithStatus(t *testing.T, repository facts.Repository, runID string, pendingID string, resumeRef string, checkpointRef string, status facts.PendingStatus) {
	t.Helper()
	if err := repository.AppendPendingInteraction(context.Background(), facts.PendingInteraction{
		PendingID:     pendingID,
		RunID:         runID,
		Kind:          facts.PendingKindApproval,
		Status:        status,
		ResumeRef:     resumeRef,
		CheckpointRef: checkpointRef,
		OperationName: "更新工单状态",
		RiskSummary:   "需要审批",
	}); err != nil {
		t.Fatal(err)
	}
}
