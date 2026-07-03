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

func TestClarificationSubmitConsumesPendingAndResumesRun(t *testing.T) {
	// clarification submit 只能使用 pending 中的安全候选引用，成功后关闭 pending 并恢复 run。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository).WithCheckpointResolver(fixedCheckpointResolver{checkpointID: "eino-checkpoint-clarify"})
	createRun(t, repository, "run-clarify-submit", "ws-hitl", facts.RunStatusWaiting, "")
	appendClarificationPending(t, repository, "run-clarify-submit", "pending-clarify-submit", "resume-clarify-submit", "checkpoint_ref:clarify-submit", facts.PendingStatusWaiting)

	accepted, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-submit",
		ResumeRef:       "resume-clarify-submit",
		ClientRequestID: "client-clarify-submit",
		Decision:        "submit",
		SelectedRefs:    []string{"candidate:fobrain:person:1"},
		FreeText:        "补充说明",
	})
	if err != nil {
		t.Fatal(err)
	}
	if accepted.RunID != "run-clarify-submit" || accepted.Status != "accepted" {
		t.Fatalf("accepted = %+v", accepted)
	}
	snapshot, err := repository.GetSnapshot(ctx, "run-clarify-submit")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusRunning {
		t.Fatalf("run status = %q, want running", snapshot.Run.Status)
	}
	if snapshot.PendingInteractions[0].Status != facts.PendingStatusConsumed {
		t.Fatalf("pending status = %q, want consumed", snapshot.PendingInteractions[0].Status)
	}
	if len(snapshot.AuditEvents) == 0 || snapshot.AuditEvents[len(snapshot.AuditEvents)-1].EventType != "clarification" {
		t.Fatalf("clarification audit missing: %+v", snapshot.AuditEvents)
	}
}

func TestClarificationDuplicateSubmitIsIdempotent(t *testing.T) {
	// 相同 client_request_id 重复提交只能命中幂等记录，不能重复追加 resume audit。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository).WithCheckpointResolver(fixedCheckpointResolver{checkpointID: "eino-checkpoint-clarify"})
	createRun(t, repository, "run-clarify-duplicate", "ws-hitl", facts.RunStatusWaiting, "")
	appendClarificationPending(t, repository, "run-clarify-duplicate", "pending-clarify-duplicate", "resume-clarify-duplicate", "checkpoint_ref:clarify-duplicate", facts.PendingStatusWaiting)

	for i := 0; i < 2; i++ {
		if _, err := commands.Resume(ctx, execution.ResumeCommand{
			WorkspaceID:     "ws-hitl",
			RunID:           "run-clarify-duplicate",
			ResumeRef:       "resume-clarify-duplicate",
			ClientRequestID: "client-clarify-duplicate",
			Decision:        "submit",
			SelectedRefs:    []string{"candidate:fobrain:person:1"},
		}); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := repository.GetSnapshot(ctx, "run-clarify-duplicate")
	if err != nil {
		t.Fatal(err)
	}
	clarificationAudits := 0
	for _, event := range snapshot.AuditEvents {
		if event.EventType == "clarification" {
			clarificationAudits++
		}
	}
	if clarificationAudits != 1 {
		t.Fatalf("clarification audit count = %d, want 1: %+v", clarificationAudits, snapshot.AuditEvents)
	}
}

func TestClarificationDuplicateSubmitRejectsDifferentPayload(t *testing.T) {
	// 同一 client_request_id 不能绑定不同选择，否则客户端会误以为新选择已生效。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository).WithCheckpointResolver(fixedCheckpointResolver{checkpointID: "eino-checkpoint-clarify"})
	createRun(t, repository, "run-clarify-conflict", "ws-hitl", facts.RunStatusWaiting, "")
	appendClarificationPending(t, repository, "run-clarify-conflict", "pending-clarify-conflict", "resume-clarify-conflict", "checkpoint_ref:clarify-conflict", facts.PendingStatusWaiting)

	if _, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-conflict",
		ResumeRef:       "resume-clarify-conflict",
		ClientRequestID: "client-clarify-conflict",
		Decision:        "submit",
		SelectedRefs:    []string{"candidate:fobrain:person:1"},
	}); err != nil {
		t.Fatal(err)
	}
	_, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-conflict",
		ResumeRef:       "resume-clarify-conflict",
		ClientRequestID: "client-clarify-conflict",
		Decision:        "submit",
		FreeText:        "另一个选择",
	})
	if !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("different payload with same idempotency key err = %v, want ErrResumeNotAllowed", err)
	}
}

func TestClarificationCancelClosesPendingWithoutCheckpoint(t *testing.T) {
	// cancel 不恢复执行，不需要 checkpoint 存在，也不能继续让旧 resume_ref submit。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository).WithCheckpointResolver(missingCheckpointResolver{})
	createRun(t, repository, "run-clarify-cancel", "ws-hitl", facts.RunStatusWaiting, "")
	appendClarificationPending(t, repository, "run-clarify-cancel", "pending-clarify-cancel", "resume-clarify-cancel", "checkpoint_ref:clarify-cancel", facts.PendingStatusWaiting)

	if _, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-cancel",
		ResumeRef:       "resume-clarify-cancel",
		ClientRequestID: "client-clarify-cancel",
		Decision:        "cancel",
	}); err != nil {
		t.Fatal(err)
	}
	snapshot, err := repository.GetSnapshot(ctx, "run-clarify-cancel")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusCancelled || snapshot.PendingInteractions[0].Status != facts.PendingStatusCancelled {
		t.Fatalf("cancel state mismatch: %+v", snapshot)
	}
	_, err = commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-cancel",
		ResumeRef:       "resume-clarify-cancel",
		ClientRequestID: "client-clarify-submit-after-cancel",
		Decision:        "submit",
		SelectedRefs:    []string{"candidate:fobrain:person:1"},
	})
	if !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("submit after cancel err = %v, want ErrResumeNotAllowed", err)
	}
}

func TestClarificationSubmitCheckpointMissingDoesNotConsumePending(t *testing.T) {
	// submit 需要 checkpoint；缺失时不能消费 pending，否则用户无法重新恢复。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository).WithCheckpointResolver(missingCheckpointResolver{})
	createRun(t, repository, "run-clarify-missing-checkpoint", "ws-hitl", facts.RunStatusWaiting, "")
	appendClarificationPending(t, repository, "run-clarify-missing-checkpoint", "pending-clarify-missing-checkpoint", "resume-clarify-missing-checkpoint", "checkpoint_ref:clarify-missing", facts.PendingStatusWaiting)

	_, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-missing-checkpoint",
		ResumeRef:       "resume-clarify-missing-checkpoint",
		ClientRequestID: "client-clarify-missing-checkpoint",
		Decision:        "submit",
		SelectedRefs:    []string{"candidate:fobrain:person:1"},
	})
	if !errors.Is(err, execution.ErrCheckpointMissing) {
		t.Fatalf("missing checkpoint submit err = %v, want ErrCheckpointMissing", err)
	}
	snapshot, err := repository.GetSnapshot(ctx, "run-clarify-missing-checkpoint")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusWaiting || snapshot.PendingInteractions[0].Status != facts.PendingStatusWaiting {
		t.Fatalf("missing checkpoint should not consume pending: %+v", snapshot)
	}
}

func TestClarificationSubmitRejectsExpiredAndUnsafeResumeData(t *testing.T) {
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository).WithCheckpointResolver(fixedCheckpointResolver{checkpointID: "eino-checkpoint-clarify"})
	createRun(t, repository, "run-clarify-expired", "ws-hitl", facts.RunStatusWaiting, "")
	appendClarificationPending(t, repository, "run-clarify-expired", "pending-clarify-expired", "resume-clarify-expired", "checkpoint_ref:clarify-expired", facts.PendingStatusExpired)

	_, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-expired",
		ResumeRef:       "resume-clarify-expired",
		ClientRequestID: "client-clarify-expired",
		Decision:        "submit",
		SelectedRefs:    []string{"candidate:fobrain:person:1"},
	})
	if !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("expired submit err = %v, want ErrResumeNotAllowed", err)
	}

	createRun(t, repository, "run-clarify-unsafe", "ws-hitl", facts.RunStatusWaiting, "")
	appendClarificationPending(t, repository, "run-clarify-unsafe", "pending-clarify-unsafe", "resume-clarify-unsafe", "checkpoint_ref:clarify-unsafe", facts.PendingStatusWaiting)
	_, err = commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-unsafe",
		ResumeRef:       "resume-clarify-unsafe",
		ClientRequestID: "client-clarify-unsafe",
		Decision:        "submit",
		SelectedRefs:    []string{"candidate:fobrain:person:missing"},
	})
	if !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("unknown candidate err = %v, want ErrResumeNotAllowed", err)
	}

	_, err = commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-unsafe",
		ResumeRef:       "resume-clarify-unsafe",
		ClientRequestID: "client-clarify-unsafe-text",
		Decision:        "submit",
		FreeText:        "Authorization: Bearer secret",
	})
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("unsafe free text err = %v, want ErrUnsafeFactMaterial", err)
	}
}

func TestClarificationAfterRestartSubmitsWithCheckpoint(t *testing.T) {
	// SQLite Product Facts 和 checkpoint store 重启后，submit 仍能校验安全 checkpoint_ref 并消费 pending。
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "clarification.db")
	repository, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	checkpoints, err := sqlite.OpenCheckpointStore(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	createRun(t, repository, "run-clarify-restart", "ws-hitl", facts.RunStatusWaiting, "")
	appendClarificationPending(t, repository, "run-clarify-restart", "pending-clarify-restart", "resume-clarify-restart", "checkpoint_ref:clarify-restart", facts.PendingStatusWaiting)
	if err := checkpoints.Set(ctx, "eino-checkpoint-clarify-restart", []byte("serialized")); err != nil {
		t.Fatal(err)
	}
	if err := checkpoints.BindCheckpointRef(ctx, "checkpoint_ref:clarify-restart", "eino-checkpoint-clarify-restart", "run-clarify-restart", "pending-clarify-restart", time.Now().UTC()); err != nil {
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

	commands := execution.NewFactCommands(reopenedRepository).WithCheckpointResolver(reopenedCheckpoints)
	if _, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           "run-clarify-restart",
		ResumeRef:       "resume-clarify-restart",
		ClientRequestID: "client-clarify-restart",
		Decision:        "submit",
		SelectedRefs:    []string{"candidate:fobrain:person:1"},
	}); err != nil {
		t.Fatal(err)
	}
	snapshot, err := reopenedRepository.GetSnapshot(ctx, "run-clarify-restart")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusRunning || snapshot.PendingInteractions[0].Status != facts.PendingStatusConsumed {
		t.Fatalf("restart clarification state mismatch: %+v", snapshot)
	}
}

func appendClarificationPending(t *testing.T, repository facts.Repository, runID string, pendingID string, resumeRef string, checkpointRef string, status facts.PendingStatus) {
	t.Helper()
	if err := repository.AppendPendingInteraction(context.Background(), facts.PendingInteraction{
		PendingID:     pendingID,
		RunID:         runID,
		Kind:          facts.PendingKindClarification,
		Status:        status,
		ResumeRef:     resumeRef,
		CheckpointRef: checkpointRef,
		Question:      "请选择要查询的人员",
		InputMode:     facts.PendingInputModeSingleChoice,
		Candidates: []facts.PendingCandidate{
			{
				CandidateRef: "candidate:fobrain:person:1",
				Label:        "张三",
				Description:  "安全部",
				EntityType:   "person",
				SafeFields: []facts.PendingCandidateField{
					{Label: "部门", Value: "安全部"},
				},
			},
		},
	}); err != nil {
		t.Fatal(err)
	}
}
