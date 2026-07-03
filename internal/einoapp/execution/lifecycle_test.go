package execution_test

import (
	"context"
	"errors"
	"testing"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
)

func TestRunCancelStopTimeoutAndTerminalIdempotency(t *testing.T) {
	// lifecycle 命令只能迁移 Product Facts，不能生成第二套运行状态。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository)

	createRun(t, repository, "run-cancel", "ws-life", facts.RunStatusRunning, "")
	appendRunningTool(t, repository, "run-cancel", "call-cancel-1")
	cancelled, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-cancel",
		Action:          execution.LifecycleActionCancel,
		ClientRequestID: "client-cancel-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if cancelled.RunID != "run-cancel" || cancelled.Status != string(facts.RunStatusCancelled) {
		t.Fatalf("cancel result = %+v", cancelled)
	}
	assertRunStatus(t, repository, "run-cancel", facts.RunStatusCancelled, "user_cancelled")
	assertToolStatus(t, repository, "run-cancel", facts.ToolCallCancelled)

	again, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-cancel",
		Action:          execution.LifecycleActionCancel,
		ClientRequestID: "client-cancel-2",
	})
	if err != nil {
		t.Fatal(err)
	}
	if again.Status != string(facts.RunStatusCancelled) {
		t.Fatalf("terminal cancel should be idempotent: %+v", again)
	}

	createRun(t, repository, "run-stop", "ws-life", facts.RunStatusRunning, "")
	appendRunningTool(t, repository, "run-stop", "call-stop-1")
	stopped, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-stop",
		Action:          execution.LifecycleActionStop,
		ClientRequestID: "client-stop-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if stopped.Status != string(facts.RunStatusStopped) {
		t.Fatalf("stop result = %+v", stopped)
	}
	assertRunStatus(t, repository, "run-stop", facts.RunStatusStopped, "user_stopped")
	assertToolStatus(t, repository, "run-stop", facts.ToolCallCancelled)

	createRun(t, repository, "run-timeout", "ws-life", facts.RunStatusRunning, "")
	appendRunningTool(t, repository, "run-timeout", "call-timeout-1")
	timedOut, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-timeout",
		Action:          execution.LifecycleActionProviderTimeout,
		ClientRequestID: "client-timeout-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if timedOut.Status != string(facts.RunStatusFailed) {
		t.Fatalf("provider timeout result = %+v", timedOut)
	}
	timeoutSnapshot, err := repository.GetSnapshot(ctx, "run-timeout")
	if err != nil {
		t.Fatal(err)
	}
	if timeoutSnapshot.Run.Status != facts.RunStatusFailed || timeoutSnapshot.Run.SafeError != "provider_timeout" {
		t.Fatalf("run timeout state = %+v", timeoutSnapshot.Run)
	}
	if len(timeoutSnapshot.ToolCalls) != 1 || timeoutSnapshot.ToolCalls[0].Status != facts.ToolCallFailed {
		t.Fatalf("tool timeout state = %+v", timeoutSnapshot.ToolCalls)
	}

	createRun(t, repository, "run-timeout-invalid", "ws-life", facts.RunStatusWaiting, "")
	if _, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-timeout-invalid",
		Action:          execution.LifecycleActionProviderTimeout,
		ClientRequestID: "client-timeout-invalid",
	}); !errors.Is(err, execution.ErrLifecycleNotAllowed) {
		t.Fatalf("provider timeout waiting run err = %v, want ErrLifecycleNotAllowed", err)
	}
}

func TestRunPendingTimeoutExpiresWaitingPending(t *testing.T) {
	// pending timeout 必须同时迁移 pending 和 run，避免后续 resume 继续旧 checkpoint。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository)
	createRun(t, repository, "run-pending", "ws-life", facts.RunStatusWaiting, "")
	if err := repository.AppendPendingInteraction(ctx, facts.PendingInteraction{
		PendingID:     "pending-1",
		RunID:         "run-pending",
		Kind:          facts.PendingKindClarification,
		Status:        facts.PendingStatusWaiting,
		ResumeRef:     "resume-safe-1",
		CheckpointRef: "checkpoint_ref:1",
		Question:      "请选择实体",
	}); err != nil {
		t.Fatal(err)
	}

	accepted, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-pending",
		Action:          execution.LifecycleActionPendingTimeout,
		ClientRequestID: "client-pending-timeout-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if accepted.Status != string(facts.RunStatusFailed) {
		t.Fatalf("pending timeout result = %+v", accepted)
	}
	snapshot, err := repository.GetSnapshot(ctx, "run-pending")
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Run.Status != facts.RunStatusFailed || snapshot.Run.SafeError != "pending_timeout" {
		t.Fatalf("run after pending timeout = %+v", snapshot.Run)
	}
	if len(snapshot.PendingInteractions) != 1 || snapshot.PendingInteractions[0].Status != facts.PendingStatusExpired {
		t.Fatalf("pending after timeout = %+v", snapshot.PendingInteractions)
	}
}

func TestRunRetryPolicyCreatesNewRunForRetryableFailure(t *testing.T) {
	// retry 只能基于安全失败摘要创建新 run，不复用旧 resume_ref 或旧 run 状态。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	nextRunID := "run-retry-new"
	commands := execution.NewFactCommands(repository).WithRunIDGenerator(func() (string, error) {
		return nextRunID, nil
	})
	createRun(t, repository, "run-retry-old", "ws-life", facts.RunStatusFailed, "schema_invalid")

	retry, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-retry-old",
		Action:          execution.LifecycleActionRetry,
		ClientRequestID: "client-retry-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if retry.RunID != "run-retry-new" || retry.Status != "accepted" {
		t.Fatalf("retry result = %+v", retry)
	}
	assertRunStatus(t, repository, "run-retry-new", facts.RunStatusCreated, "")
	nextRunID = "run-retry-duplicate"
	duplicate, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-retry-old",
		Action:          execution.LifecycleActionRetry,
		ClientRequestID: "client-retry-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if duplicate.RunID != "run-retry-new" {
		t.Fatalf("duplicate retry result = %+v", duplicate)
	}

	createRun(t, repository, "run-succeeded", "ws-life", facts.RunStatusSucceeded, "")
	if _, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-succeeded",
		Action:          execution.LifecycleActionRetry,
		ClientRequestID: "client-retry-denied",
	}); !errors.Is(err, execution.ErrLifecycleNotAllowed) {
		t.Fatalf("retry succeeded run err = %v, want ErrLifecycleNotAllowed", err)
	}
	createRun(t, repository, "run-provider-timeout", "ws-life", facts.RunStatusFailed, "provider_timeout")
	if _, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-life",
		RunID:           "run-provider-timeout",
		Action:          execution.LifecycleActionRetry,
		ClientRequestID: "client-retry-provider-timeout",
	}); !errors.Is(err, execution.ErrLifecycleNotAllowed) {
		t.Fatalf("provider timeout retry err = %v, want ErrLifecycleNotAllowed", err)
	}
}

func TestRunLifecycleRejectsCrossWorkspaceAccess(t *testing.T) {
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository)
	createRun(t, repository, "run-private", "ws-owner", facts.RunStatusRunning, "")

	if _, err := commands.RunLifecycle(ctx, execution.LifecycleCommand{
		WorkspaceID:     "ws-other",
		RunID:           "run-private",
		Action:          execution.LifecycleActionCancel,
		ClientRequestID: "client-cross",
	}); !errors.Is(err, execution.ErrRunNotFound) {
		t.Fatalf("cross workspace err = %v, want ErrRunNotFound", err)
	}
}

func createRun(t *testing.T, repository facts.Repository, runID string, workspaceID string, status facts.RunStatus, safeError string) {
	t.Helper()
	if err := repository.CreateRun(context.Background(), facts.Run{
		RunID:       runID,
		WorkspaceID: workspaceID,
		Status:      status,
		SafeError:   safeError,
	}); err != nil {
		t.Fatal(err)
	}
}

func appendRunningTool(t *testing.T, repository facts.Repository, runID string, toolCallID string) {
	t.Helper()
	if err := repository.AppendToolCall(context.Background(), facts.ToolCall{
		ToolCallID:  toolCallID,
		RunID:       runID,
		ToolID:      "cap.read",
		DisplayName: "只读查询",
		Status:      facts.ToolCallRunning,
	}); err != nil {
		t.Fatal(err)
	}
}

func assertToolStatus(t *testing.T, repository facts.Repository, runID string, status facts.ToolCallStatus) {
	t.Helper()
	snapshot, err := repository.GetSnapshot(context.Background(), runID)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.ToolCalls) != 1 || snapshot.ToolCalls[0].Status != status {
		t.Fatalf("tool calls for %s = %+v, want status=%s", runID, snapshot.ToolCalls, status)
	}
}

func assertRunStatus(t *testing.T, repository facts.Repository, runID string, status facts.RunStatus, safeError string) {
	t.Helper()
	run, err := repository.GetRun(context.Background(), runID)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != status || run.SafeError != safeError {
		t.Fatalf("run %s = %+v, want status=%s safe_error=%q", runID, run, status, safeError)
	}
}
