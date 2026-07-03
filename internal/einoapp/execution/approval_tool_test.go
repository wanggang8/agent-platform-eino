package execution_test

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/store/sqlite"
)

func TestApprovalInterruptCreatesWaitingPendingAndDoesNotExecuteMutation(t *testing.T) {
	// 写域 capability 必须先转成 approval pending 和内部 checkpoint，审批前绝不触发 provider。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	toolRunner := &recordingCapabilityRunner{}
	checkpoints := newMemoryApprovalCheckpointStore()
	commands := approvalCommands(repository, toolRunner).WithApprovalCheckpointStore(checkpoints).WithRunIDGenerator(fixedRunID("run-approval-interrupt"))

	accepted, err := commands.StartAction(ctx, approvalActionCommand("client-start-approval"))
	if err != nil {
		t.Fatal(err)
	}
	if toolRunner.runID != "" {
		t.Fatalf("approval-required capability executed before approval: %+v", toolRunner)
	}
	snapshot, err := repository.GetSnapshot(ctx, accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusWaiting {
		t.Fatalf("run status = %q, want waiting", snapshot.Run.Status)
	}
	if len(snapshot.PendingInteractions) != 1 {
		t.Fatalf("pending interactions = %+v", snapshot.PendingInteractions)
	}
	pending := snapshot.PendingInteractions[0]
	if pending.Kind != facts.PendingKindApproval || pending.Status != facts.PendingStatusWaiting || !strings.HasPrefix(pending.ResumeRef, "resume_ref:") {
		t.Fatalf("approval pending mismatch: %+v", pending)
	}
	if !facts.SafeCheckpointRef(pending.CheckpointRef) || pending.OperationName != "更新业务数据" || pending.RiskSummary == "" {
		t.Fatalf("approval pending safe fields mismatch: %+v", pending)
	}
	checkpointID, ok, err := checkpoints.ResolveCheckpointID(ctx, pending.CheckpointRef, pending.RunID, pending.PendingID)
	if err != nil || !ok {
		t.Fatalf("checkpoint binding missing: checkpointID=%q ok=%v err=%v", checkpointID, ok, err)
	}
	payload, ok, err := checkpoints.Get(ctx, checkpointID)
	if err != nil || !ok || !bytes.Contains(payload, []byte("danger.write")) {
		t.Fatalf("approval continuation payload missing capability: payload=%s ok=%v err=%v", payload, ok, err)
	}
	if len(snapshot.AuditEvents) < 2 || snapshot.AuditEvents[len(snapshot.AuditEvents)-1].EventType != "approval" {
		t.Fatalf("approval request audit missing: %+v", snapshot.AuditEvents)
	}
}

func TestApprovalInterruptResumeDuplicateApproveDoesNotRerunMutation(t *testing.T) {
	// 同一个 client_request_id 重放 approve 时只能命中幂等记录，不能重复执行写域能力。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	toolRunner := &recordingCapabilityRunner{}
	commands := approvalCommands(repository, toolRunner).
		WithApprovalCheckpointStore(newMemoryApprovalCheckpointStore()).
		WithRunIDGenerator(fixedRunID("run-approval-approve"))

	accepted, pending := startApprovalRun(t, ctx, commands, repository)
	for i := 0; i < 2; i++ {
		got, err := commands.Resume(ctx, execution.ResumeCommand{
			WorkspaceID:     "ws-hitl",
			RunID:           accepted.RunID,
			ResumeRef:       pending.ResumeRef,
			ClientRequestID: "client-approve-once",
			Decision:        "approve",
			Comment:         "确认执行",
		})
		if err != nil {
			t.Fatal(err)
		}
		if got.RunID != accepted.RunID || got.Status != "accepted" {
			t.Fatalf("accepted = %+v", got)
		}
	}
	if toolRunner.calls != 1 || toolRunner.capabilityID != "danger.write" || toolRunner.inputText != "执行写入" {
		t.Fatalf("tool runner should execute exactly once after approve: %+v", toolRunner)
	}
	snapshot, err := repository.GetSnapshot(ctx, accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.PendingInteractions[0].Status != facts.PendingStatusApproved {
		t.Fatalf("pending status = %q, want approved", snapshot.PendingInteractions[0].Status)
	}
}

func TestApprovalInterruptResumeRejectCannotApprove(t *testing.T) {
	// reject 会终止 run 并关闭 pending；后续 approve 即使换幂等键也不能再次打开写域执行。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	toolRunner := &recordingCapabilityRunner{}
	commands := approvalCommands(repository, toolRunner).
		WithApprovalCheckpointStore(newMemoryApprovalCheckpointStore()).
		WithRunIDGenerator(fixedRunID("run-approval-reject"))

	accepted, pending := startApprovalRun(t, ctx, commands, repository)
	_, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		ResumeRef:       pending.ResumeRef,
		ClientRequestID: "client-reject",
		Decision:        "reject",
		Comment:         "拒绝",
	})
	if err != nil {
		t.Fatal(err)
	}
	if toolRunner.calls != 0 {
		t.Fatalf("rejected approval must not execute mutation: %+v", toolRunner)
	}
	snapshot, err := repository.GetSnapshot(ctx, accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusFailed || snapshot.Run.SafeError != "approval_rejected" {
		t.Fatalf("run after reject mismatch: %+v", snapshot.Run)
	}
	if snapshot.PendingInteractions[0].Status != facts.PendingStatusRejected {
		t.Fatalf("pending status after reject = %q", snapshot.PendingInteractions[0].Status)
	}
	_, err = commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		ResumeRef:       pending.ResumeRef,
		ClientRequestID: "client-reject",
		Decision:        "approve",
	})
	if !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("approve with rejected idempotency key err = %v, want ErrResumeNotAllowed", err)
	}
	_, err = commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		ResumeRef:       pending.ResumeRef,
		ClientRequestID: "client-approve-after-reject",
		Decision:        "approve",
	})
	if !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("approve after reject err = %v, want ErrResumeNotAllowed", err)
	}
	if toolRunner.calls != 0 {
		t.Fatalf("approve after reject must not execute mutation: %+v", toolRunner)
	}
}

func TestApprovalInterruptResumeRejectDoesNotRequireCheckpoint(t *testing.T) {
	// reject 不需要恢复 mutation continuation；即使 checkpoint 丢失，也必须能关闭审批。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	toolRunner := &recordingCapabilityRunner{}
	checkpoints := newMemoryApprovalCheckpointStore()
	commands := approvalCommands(repository, toolRunner).
		WithApprovalCheckpointStore(checkpoints).
		WithRunIDGenerator(fixedRunID("run-approval-reject-missing-checkpoint"))

	accepted, pending := startApprovalRun(t, ctx, commands, repository)
	checkpoints.values = map[string][]byte{}
	_, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		ResumeRef:       pending.ResumeRef,
		ClientRequestID: "client-reject-missing-checkpoint",
		Decision:        "reject",
	})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := repository.GetSnapshot(ctx, accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.PendingInteractions[0].Status != facts.PendingStatusRejected || snapshot.Run.SafeError != "approval_rejected" {
		t.Fatalf("reject without checkpoint should close pending safely: %+v", snapshot)
	}
	if toolRunner.calls != 0 {
		t.Fatalf("reject without checkpoint must not execute mutation: %+v", toolRunner)
	}
}

func TestApprovalInterruptResumeCancelDoesNotRequireCheckpoint(t *testing.T) {
	// cancel 表示用户取消整个审批和 run，不需要恢复 checkpoint，也不能触发写域能力。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	toolRunner := &recordingCapabilityRunner{}
	checkpoints := newMemoryApprovalCheckpointStore()
	commands := approvalCommands(repository, toolRunner).
		WithApprovalCheckpointStore(checkpoints).
		WithRunIDGenerator(fixedRunID("run-approval-cancel"))

	accepted, pending := startApprovalRun(t, ctx, commands, repository)
	checkpoints.values = map[string][]byte{}
	_, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		ResumeRef:       pending.ResumeRef,
		ClientRequestID: "client-cancel-approval",
		Decision:        "cancel",
	})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := repository.GetSnapshot(ctx, accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.PendingInteractions[0].Status != facts.PendingStatusCancelled {
		t.Fatalf("pending status after cancel = %q, want cancelled", snapshot.PendingInteractions[0].Status)
	}
	if snapshot.Run.Status != facts.RunStatusCancelled || snapshot.Run.SafeError != "approval_cancelled" {
		t.Fatalf("run after cancel mismatch: %+v", snapshot.Run)
	}
	if toolRunner.calls != 0 {
		t.Fatalf("cancelled approval must not execute mutation: %+v", toolRunner)
	}
}

func TestApprovalInterruptResumeCancelIsIdempotentAndCannotApprove(t *testing.T) {
	// cancel 后旧 resume_ref 进入终态；同一幂等键可重放，不同键 approve 不能重新打开写域能力。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	toolRunner := &recordingCapabilityRunner{}
	commands := approvalCommands(repository, toolRunner).
		WithApprovalCheckpointStore(newMemoryApprovalCheckpointStore()).
		WithRunIDGenerator(fixedRunID("run-approval-cancel-idempotent"))

	accepted, pending := startApprovalRun(t, ctx, commands, repository)
	for i := 0; i < 2; i++ {
		got, err := commands.Resume(ctx, execution.ResumeCommand{
			WorkspaceID:     "ws-hitl",
			RunID:           accepted.RunID,
			ResumeRef:       pending.ResumeRef,
			ClientRequestID: "client-cancel-approval-once",
			Decision:        "cancel",
		})
		if err != nil {
			t.Fatal(err)
		}
		if got.RunID != accepted.RunID || got.Status != "accepted" {
			t.Fatalf("accepted = %+v", got)
		}
	}
	_, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		ResumeRef:       pending.ResumeRef,
		ClientRequestID: "client-approve-after-cancel",
		Decision:        "approve",
	})
	if !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("approve after cancel err = %v, want ErrResumeNotAllowed", err)
	}
	if toolRunner.calls != 0 {
		t.Fatalf("approve after cancel must not execute mutation: %+v", toolRunner)
	}
}

func TestApprovalInterruptResumeExpiredCannotApprove(t *testing.T) {
	// pending timeout 会把 approval 转为 expired；过期后不能再通过旧 resume_ref 执行 mutation。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	toolRunner := &recordingCapabilityRunner{}
	commands := approvalCommands(repository, toolRunner).
		WithApprovalCheckpointStore(newMemoryApprovalCheckpointStore()).
		WithRunIDGenerator(fixedRunID("run-approval-expired"))

	accepted, pending := startApprovalRun(t, ctx, commands, repository)
	if _, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		Action:          execution.LifecycleActionPendingTimeout,
		ClientRequestID: "client-timeout",
	}); err != nil {
		t.Fatal(err)
	}
	_, err := commands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		ResumeRef:       pending.ResumeRef,
		ClientRequestID: "client-approve-after-expired",
		Decision:        "approve",
	})
	if !errors.Is(err, execution.ErrResumeNotAllowed) {
		t.Fatalf("approve after expired err = %v, want ErrResumeNotAllowed", err)
	}
	if toolRunner.calls != 0 {
		t.Fatalf("expired approval must not execute mutation: %+v", toolRunner)
	}
	snapshot, err := repository.GetSnapshot(ctx, accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.PendingInteractions[0].Status != facts.PendingStatusExpired {
		t.Fatalf("pending status after timeout = %q", snapshot.PendingInteractions[0].Status)
	}
}

func TestApprovalInterruptResumeAfterRestartContinuesMutation(t *testing.T) {
	// Product Facts 与 checkpoint store 重启后，approval resume 仍只能通过安全 checkpoint_ref 恢复内部 continuation。
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "approval.db")
	repository, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	checkpoints, err := sqlite.OpenCheckpointStore(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	startCommands := approvalCommands(repository, &recordingCapabilityRunner{}).
		WithApprovalCheckpointStore(checkpoints).
		WithRunIDGenerator(fixedRunID("run-approval-restart"))
	accepted, pending := startApprovalRun(t, ctx, startCommands, repository)
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
	toolRunner := &recordingCapabilityRunner{}
	resumeCommands := approvalCommands(reopenedRepository, toolRunner).WithApprovalCheckpointStore(reopenedCheckpoints)
	_, err = resumeCommands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		ResumeRef:       pending.ResumeRef,
		ClientRequestID: "client-restart-approve",
		Decision:        "approve",
	})
	if err != nil {
		t.Fatal(err)
	}
	if toolRunner.calls != 1 || toolRunner.runID != accepted.RunID || toolRunner.capabilityID != "danger.write" {
		t.Fatalf("restart resume did not continue stored mutation: %+v", toolRunner)
	}
}

func TestApprovalInterruptResumeCancelAfterRestartClosesRun(t *testing.T) {
	// SQLite Product Facts 重启后，cancel 仍必须只关闭审批和 run，不恢复 mutation continuation。
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "approval-cancel.db")
	repository, err := sqlite.Open(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	checkpoints, err := sqlite.OpenCheckpointStore(ctx, dbPath)
	if err != nil {
		t.Fatal(err)
	}
	startCommands := approvalCommands(repository, &recordingCapabilityRunner{}).
		WithApprovalCheckpointStore(checkpoints).
		WithRunIDGenerator(fixedRunID("run-approval-cancel-restart"))
	accepted, pending := startApprovalRun(t, ctx, startCommands, repository)
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
	toolRunner := &recordingCapabilityRunner{}
	resumeCommands := approvalCommands(reopenedRepository, toolRunner).WithApprovalCheckpointStore(reopenedCheckpoints)
	_, err = resumeCommands.Resume(ctx, execution.ResumeCommand{
		WorkspaceID:     "ws-hitl",
		RunID:           accepted.RunID,
		ResumeRef:       pending.ResumeRef,
		ClientRequestID: "client-restart-cancel",
		Decision:        "cancel",
	})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := reopenedRepository.GetSnapshot(ctx, accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.PendingInteractions[0].Status != facts.PendingStatusCancelled || snapshot.Run.Status != facts.RunStatusCancelled || snapshot.Run.SafeError != "approval_cancelled" {
		t.Fatalf("restart cancel state mismatch: %+v", snapshot)
	}
	if toolRunner.calls != 0 {
		t.Fatalf("restart cancel must not execute mutation: %+v", toolRunner)
	}
}

func approvalCommands(repository facts.Repository, runner *recordingCapabilityRunner) execution.StaticCommands {
	registry := capabilities.NewRegistry()
	if err := registry.Register(approvalCapability()); err != nil {
		panic(err)
	}
	return execution.NewToolRunnerCommandsWithRegistry(repository, &recordingRunner{}, runner, registry)
}

func approvalCapability() capabilities.Capability {
	return capabilities.Capability{
		ID:                  "danger.write",
		ProviderID:          "demo",
		ToolName:            "danger_write",
		DisplayName:         "更新业务数据",
		Description:         "更新业务数据",
		RiskLevel:           capabilities.RiskWrite,
		SideEffect:          capabilities.SideEffectWriteExternal,
		ApprovalRequired:    true,
		IdempotencyRequired: true,
	}
}

func approvalActionCommand(clientRequestID string) execution.ActionCommand {
	return execution.ActionCommand{
		WorkspaceID:     "ws-hitl",
		ActionID:        "action-demo",
		ClientRequestID: clientRequestID,
		CapabilityHint:  "danger.write",
		InputText:       "执行写入",
	}
}

func startApprovalRun(t *testing.T, ctx context.Context, commands execution.StaticCommands, repository facts.Repository) (execution.AcceptedRun, facts.PendingInteraction) {
	t.Helper()
	accepted, err := commands.StartAction(ctx, approvalActionCommand("client-start-"+time.Now().Format("150405.000000000")))
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := repository.GetSnapshot(ctx, accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.PendingInteractions) != 1 {
		t.Fatalf("pending interactions = %+v", snapshot.PendingInteractions)
	}
	return accepted, snapshot.PendingInteractions[0]
}

func fixedRunID(runID string) func() (string, error) {
	return func() (string, error) { return runID, nil }
}

type memoryApprovalCheckpointStore struct {
	values   map[string][]byte
	bindings map[string]memoryCheckpointBinding
}

type memoryCheckpointBinding struct {
	checkpointID string
	runID        string
	pendingID    string
}

func newMemoryApprovalCheckpointStore() *memoryApprovalCheckpointStore {
	return &memoryApprovalCheckpointStore{
		values:   map[string][]byte{},
		bindings: map[string]memoryCheckpointBinding{},
	}
}

func (store *memoryApprovalCheckpointStore) Set(_ context.Context, checkpointID string, payload []byte) error {
	store.values[checkpointID] = append([]byte(nil), payload...)
	return nil
}

func (store *memoryApprovalCheckpointStore) Get(_ context.Context, checkpointID string) ([]byte, bool, error) {
	payload, ok := store.values[checkpointID]
	return append([]byte(nil), payload...), ok, nil
}

func (store *memoryApprovalCheckpointStore) BindCheckpointRef(_ context.Context, checkpointRef string, checkpointID string, runID string, pendingID string, _ time.Time) error {
	store.bindings[checkpointRef] = memoryCheckpointBinding{checkpointID: checkpointID, runID: runID, pendingID: pendingID}
	return nil
}

func (store *memoryApprovalCheckpointStore) ResolveCheckpointID(_ context.Context, checkpointRef string, runID string, pendingID string) (string, bool, error) {
	binding, ok := store.bindings[checkpointRef]
	if !ok || binding.runID != runID || binding.pendingID != pendingID {
		return "", false, nil
	}
	if _, ok := store.values[binding.checkpointID]; !ok {
		return "", false, nil
	}
	return binding.checkpointID, true, nil
}
