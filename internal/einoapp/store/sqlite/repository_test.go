package sqlite_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/store/sqlite"
)

func TestRunTurnEventPersistsWithSequence(t *testing.T) {
	repository := openTestRepository(t)
	ctx := context.Background()
	now := time.Unix(100, 0).UTC()

	run := facts.Run{
		RunID:       "run-1",
		WorkspaceID: "ws-1",
		Status:      facts.RunStatusCreated,
		CreatedAt:   now,
		UpdatedAt:   now,
		ModelLabel:  "mock-model",
	}
	if err := repository.CreateRun(ctx, run); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendTurn(ctx, facts.Turn{
		TurnID:    "turn-1",
		RunID:     run.RunID,
		Role:      facts.TurnRoleUser,
		Content:   "hello",
		Sequence:  1,
		CreatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}

	loaded, err := repository.GetRun(ctx, run.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.RunID != run.RunID || loaded.Status != facts.RunStatusCreated || loaded.ModelLabel != "mock-model" {
		t.Fatalf("loaded run mismatch: %+v", loaded)
	}

	latest, err := repository.LatestRun(ctx, run.WorkspaceID)
	if err != nil {
		t.Fatal(err)
	}
	if latest.RunID != run.RunID {
		t.Fatalf("latest run = %q", latest.RunID)
	}
}

func TestToolCallResultContextSnapshotAndAuditPersist(t *testing.T) {
	repository := openTestRepository(t)
	ctx := context.Background()
	now := time.Unix(200, 0).UTC()

	if err := repository.CreateRun(ctx, facts.Run{
		RunID:       "run-1",
		WorkspaceID: "ws-1",
		Status:      facts.RunStatusRunning,
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendToolCall(ctx, facts.ToolCall{
		ToolCallID:  "call-1",
		RunID:       "run-1",
		ToolID:      "tool.read",
		DisplayName: "只读查询",
		Status:      facts.ToolCallRunning,
		ArgsPreview: "target=example",
		CreatedAt:   now,
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendToolResult(ctx, facts.ToolResult{
		ResultID:   "result-1",
		ToolCallID: "call-1",
		Status:     facts.ToolResultSucceeded,
		StructuredResult: facts.StructuredResultRef{
			SchemaVersion: "tool.structured_result.v1",
			ResultRef:     "result:call-1",
			SafeSummary:   "查询完成",
		},
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.SaveContextSnapshot(ctx, facts.ContextSnapshot{
		SnapshotID:  "snapshot-1",
		RunID:       "run-1",
		SafeSummary: "safe context",
		CreatedAt:   now,
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendAuditEvent(ctx, facts.AuditEvent{
		AuditID:     "audit-1",
		RunID:       "run-1",
		EventType:   "tool",
		SafeSummary: "tool completed",
		Actor:       "system",
		CreatedAt:   now,
	}); err != nil {
		t.Fatal(err)
	}

	snapshot, err := repository.GetSnapshot(ctx, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.ToolCalls) != 1 || snapshot.ToolCalls[0].ToolID != "tool.read" || snapshot.ToolCalls[0].ArgsPreview != "target=example" {
		t.Fatalf("tool calls not projected from facts: %+v", snapshot.ToolCalls)
	}
	if len(snapshot.ToolResults) != 1 || snapshot.ToolResults[0].StructuredResult.SafeSummary != "查询完成" {
		t.Fatalf("tool results not projected from facts: %+v", snapshot.ToolResults)
	}
	if len(snapshot.ContextSnapshots) != 1 || snapshot.ContextSnapshots[0].SafeSummary != "safe context" {
		t.Fatalf("context snapshots not projected from facts: %+v", snapshot.ContextSnapshots)
	}
	if len(snapshot.AuditEvents) != 1 || snapshot.AuditEvents[0].SafeSummary != "tool completed" {
		t.Fatalf("audit events not projected from facts: %+v", snapshot.AuditEvents)
	}
}

func TestPendingResumeRefIsConsumedOnce(t *testing.T) {
	repository := openTestRepository(t)
	ctx := context.Background()
	now := time.Unix(300, 0).UTC()

	if err := repository.CreateRun(ctx, facts.Run{
		RunID:       "run-1",
		WorkspaceID: "ws-1",
		Status:      facts.RunStatusWaiting,
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendPendingInteraction(ctx, facts.PendingInteraction{
		PendingID:     "pending-1",
		RunID:         "run-1",
		Kind:          facts.PendingKindApproval,
		Status:        facts.PendingStatusWaiting,
		ResumeRef:     "resume-safe-1",
		CheckpointRef: "checkpoint-safe-1",
		RiskSummary:   "需要审批",
		ExpiresAt:     now.Add(time.Hour),
	}); err != nil {
		t.Fatal(err)
	}

	pending, err := repository.ConsumeResumeRef(ctx, "resume-safe-1", facts.PendingStatusApproved)
	if err != nil {
		t.Fatal(err)
	}
	if pending.Status != facts.PendingStatusConsumed || pending.CheckpointRef != "checkpoint-safe-1" {
		t.Fatalf("consumed pending mismatch: %+v", pending)
	}

	again, err := repository.ConsumeResumeRef(ctx, "resume-safe-1", facts.PendingStatusApproved)
	if !errors.Is(err, facts.ErrResumeAlreadyConsumed) {
		t.Fatalf("second consume err = %v", err)
	}
	if again.PendingID != "pending-1" || again.Status != facts.PendingStatusConsumed {
		t.Fatalf("duplicate consume must return terminal pending: %+v", again)
	}
}

func TestConsumeResumeRefWithIdempotencyRecordsSameTransaction(t *testing.T) {
	repository := openTestRepository(t)
	ctx := context.Background()
	now := time.Unix(325, 0).UTC()

	if err := repository.CreateRun(ctx, facts.Run{
		RunID:       "run-1",
		WorkspaceID: "ws-1",
		Status:      facts.RunStatusWaiting,
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendPendingInteraction(ctx, facts.PendingInteraction{
		PendingID:     "pending-1",
		RunID:         "run-1",
		Kind:          facts.PendingKindApproval,
		Status:        facts.PendingStatusWaiting,
		ResumeRef:     "resume-safe-1",
		CheckpointRef: "checkpoint-safe-1",
	}); err != nil {
		t.Fatal(err)
	}

	pending, record, existed, err := repository.ConsumeResumeRefWithIdempotency(ctx, "resume-safe-1", facts.PendingStatusApproved, facts.IdempotencyRecord{
		Scope:       facts.IdempotencyScopeResume,
		Key:         "resume-request-1",
		RunID:       "run-1",
		ResourceRef: "resume-safe-1",
		Status:      "consumed",
		CreatedAt:   now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if existed || pending.Status != facts.PendingStatusConsumed || record.Key != "resume-request-1" {
		t.Fatalf("resume transaction mismatch: pending=%+v record=%+v existed=%v", pending, record, existed)
	}

	duplicate, existed, err := repository.RecordIdempotency(ctx, facts.IdempotencyRecord{
		Scope:       facts.IdempotencyScopeResume,
		Key:         "resume-request-1",
		RunID:       "run-other",
		ResourceRef: "resume-safe-1",
		Status:      "duplicate",
		CreatedAt:   now.Add(time.Second),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !existed || duplicate.RunID != "run-1" || duplicate.Status != "consumed" {
		t.Fatalf("resume idempotency record was not committed with consume: %+v existed=%v", duplicate, existed)
	}

	_, _, _, err = repository.ConsumeResumeRefWithIdempotency(ctx, "resume-other", facts.PendingStatusApproved, facts.IdempotencyRecord{
		Scope:       facts.IdempotencyScopeResume,
		Key:         "resume-request-1",
		RunID:       "run-1",
		ResourceRef: "resume-other",
		Status:      "duplicate",
		CreatedAt:   now.Add(2 * time.Second),
	})
	if !errors.Is(err, facts.ErrIdempotencyConflict) {
		t.Fatalf("conflicting idempotency replay err = %v", err)
	}

	if err := repository.AppendPendingInteraction(ctx, facts.PendingInteraction{
		PendingID:     "pending-2",
		RunID:         "run-1",
		Kind:          facts.PendingKindApproval,
		Status:        facts.PendingStatusWaiting,
		ResumeRef:     "resume-other-existing",
		CheckpointRef: "checkpoint-safe-2",
	}); err != nil {
		t.Fatal(err)
	}
	_, _, _, err = repository.ConsumeResumeRefWithIdempotency(ctx, "resume-other-existing", facts.PendingStatusApproved, facts.IdempotencyRecord{
		Scope:       facts.IdempotencyScopeResume,
		Key:         "resume-request-1",
		RunID:       "run-1",
		ResourceRef: "resume-safe-1",
		Status:      "duplicate",
		CreatedAt:   now.Add(3 * time.Second),
	})
	if !errors.Is(err, facts.ErrIdempotencyConflict) {
		t.Fatalf("mismatched resumeRef/resourceRef err = %v", err)
	}
}

func TestIdempotencyRecordReturnsExistingRecord(t *testing.T) {
	repository := openTestRepository(t)
	ctx := context.Background()
	now := time.Unix(350, 0).UTC()

	record := facts.IdempotencyRecord{
		Scope:       facts.IdempotencyScopeMutation,
		Key:         "mutation-key-1",
		RunID:       "run-1",
		ResourceRef: "mutation-safe-1",
		Status:      "started",
		CreatedAt:   now,
	}
	first, existed, err := repository.RecordIdempotency(ctx, record)
	if err != nil {
		t.Fatal(err)
	}
	if existed || first.Key != record.Key {
		t.Fatalf("first idempotency record mismatch: existed=%v record=%+v", existed, first)
	}

	duplicate := record
	duplicate.RunID = "run-other"
	duplicate.Status = "replayed"
	second, existed, err := repository.RecordIdempotency(ctx, duplicate)
	if err != nil {
		t.Fatal(err)
	}
	if !existed {
		t.Fatal("duplicate idempotency key must return existing record")
	}
	if second.RunID != "run-1" || second.Status != "started" {
		t.Fatalf("duplicate must not overwrite original record: %+v", second)
	}
}

func TestUnsafeFactMaterialIsRejected(t *testing.T) {
	repository := openTestRepository(t)
	ctx := context.Background()
	now := time.Unix(375, 0).UTC()

	if err := repository.CreateRun(ctx, facts.Run{
		RunID:       "run-1",
		WorkspaceID: "ws-1",
		Status:      facts.RunStatusWaiting,
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		t.Fatal(err)
	}
	err := repository.AppendPendingInteraction(ctx, facts.PendingInteraction{
		PendingID:     "pending-1",
		RunID:         "run-1",
		Kind:          facts.PendingKindApproval,
		Status:        facts.PendingStatusWaiting,
		ResumeRef:     "resume-safe-1",
		CheckpointRef: "checkpoint-raw-eino-id",
	})
	if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
		t.Fatalf("raw checkpoint ref err = %v", err)
	}

	for _, testCase := range []struct {
		name    string
		preview string
	}{
		{name: "authorization", preview: "Authorization: Bearer token"},
		{name: "credential", preview: "credential_ref=abc"},
		{name: "raw provider", preview: "raw provider body"},
		{name: "secret", preview: "secret=value"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			err := repository.AppendToolCall(ctx, facts.ToolCall{
				ToolCallID:  "call-" + testCase.name,
				RunID:       "run-1",
				ToolID:      "tool.read",
				DisplayName: "只读查询",
				Status:      facts.ToolCallRunning,
				ArgsPreview: testCase.preview,
				CreatedAt:   now,
			})
			if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
				t.Fatalf("unsafe args preview err = %v", err)
			}
		})
	}
}

func TestLifecycleTransitionUpdatesRunStatus(t *testing.T) {
	repository := openTestRepository(t)
	ctx := context.Background()
	createdAt := time.Unix(400, 0).UTC()
	updatedAt := time.Unix(401, 0).UTC()

	if err := repository.CreateRun(ctx, facts.Run{
		RunID:       "run-1",
		WorkspaceID: "ws-1",
		Status:      facts.RunStatusRunning,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.UpdateRunStatus(ctx, "run-1", facts.RunStatusFailed, "safe failure", updatedAt); err != nil {
		t.Fatal(err)
	}

	loaded, err := repository.GetRun(ctx, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Status != facts.RunStatusFailed || loaded.SafeError != "safe failure" || !loaded.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("run lifecycle transition mismatch: %+v", loaded)
	}
}

func openTestRepository(t *testing.T) *sqlite.Repository {
	t.Helper()

	path := filepath.Join(t.TempDir(), "facts.db")
	repository, err := sqlite.Open(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := repository.Close(); err != nil {
			t.Fatal(err)
		}
	})
	return repository
}
