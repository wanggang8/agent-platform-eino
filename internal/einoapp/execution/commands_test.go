package execution_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
)

func TestStaticCommandsStartMessageUsesProvidedRunID(t *testing.T) {
	commands := execution.NewStaticCommands()

	accepted, err := commands.StartMessage(context.Background(), execution.MessageCommand{
		WorkspaceID: "ws_123",
		RunID:       "run_existing",
	})
	if err != nil {
		t.Fatal(err)
	}

	if accepted.RunID != "run_existing" {
		t.Fatalf("RunID = %q, want provided run id", accepted.RunID)
	}
	if accepted.Status != "accepted" {
		t.Fatalf("Status = %q, want accepted", accepted.Status)
	}
}

func TestStaticCommandsStartMessageCreatesOpaqueRunID(t *testing.T) {
	commands := execution.NewStaticCommands()

	accepted, err := commands.StartMessage(context.Background(), execution.MessageCommand{
		WorkspaceID: "ws_123",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(accepted.RunID, "run_") || strings.Contains(accepted.RunID, "ws_123") {
		t.Fatalf("RunID = %q", accepted.RunID)
	}
}

func TestStaticCommandsStartActionCreatesOpaqueRunID(t *testing.T) {
	commands := execution.NewStaticCommands()

	accepted, err := commands.StartAction(context.Background(), execution.ActionCommand{
		ActionID: "action-demo",
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(accepted.RunID, "run_") || strings.Contains(accepted.RunID, "action-demo") {
		t.Fatalf("RunID = %q", accepted.RunID)
	}
}

func TestFactCommandsWriteAcceptedRunFact(t *testing.T) {
	// 命令层接受请求时必须先写 Product Facts，避免 Workbench 和 Action API 分裂。
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository)

	accepted, err := commands.StartMessage(context.Background(), execution.MessageCommand{
		WorkspaceID: "ws_123",
	})
	if err != nil {
		t.Fatal(err)
	}

	run, err := repository.GetRun(context.Background(), accepted.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if run.WorkspaceID != "ws_123" || run.Status != facts.RunStatusCreated {
		t.Fatalf("run fact not populated: %+v", run)
	}
}

func TestFactCommandsAppendInputTurnForMessageAndAction(t *testing.T) {
	// message/action 入口的用户输入必须先成为 Product Facts，SSE 和 projection 才能同源读取。
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository)
	ctx := context.Background()

	messageRun, err := commands.StartMessage(ctx, execution.MessageCommand{
		WorkspaceID:     "ws_123",
		Message:         "查询资产",
		ClientRequestID: "client-message",
	})
	if err != nil {
		t.Fatal(err)
	}
	actionRun, err := commands.StartAction(ctx, execution.ActionCommand{
		WorkspaceID:     "ws_123",
		ActionID:        "action-demo",
		ClientRequestID: "client-action",
		InputText:       "查询漏洞",
	})
	if err != nil {
		t.Fatal(err)
	}

	messageSnapshot, err := repository.GetSnapshot(ctx, messageRun.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if len(messageSnapshot.Turns) != 1 || messageSnapshot.Turns[0].Content != "查询资产" {
		t.Fatalf("message turns = %+v", messageSnapshot.Turns)
	}
	actionSnapshot, err := repository.GetSnapshot(ctx, actionRun.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if len(actionSnapshot.Turns) != 1 || actionSnapshot.Turns[0].Content != "查询漏洞" {
		t.Fatalf("action turns = %+v", actionSnapshot.Turns)
	}
}

func TestRunnerCommandsExecuteChatModelRunnerAfterRunCreated(t *testing.T) {
	// Phase 3 服务路径必须在创建 run 后执行 ChatModelRunner，确保 assistant/context facts 落库。
	repository := facts.NewMemoryRepository()
	runner := &recordingRunner{}
	commands := execution.NewRunnerCommands(repository, runner)

	accepted, err := commands.StartMessage(context.Background(), execution.MessageCommand{
		WorkspaceID:     "ws_123",
		Message:         "hello",
		ClientRequestID: "client-runner",
	})
	if err != nil {
		t.Fatal(err)
	}
	if runner.runID != accepted.RunID {
		t.Fatalf("runner runID = %q, want %q", runner.runID, accepted.RunID)
	}
}

func TestFactCommandsResumeRequiresExistingRun(t *testing.T) {
	// resume 不能凭前端传参创建隐式 run，必须绑定已存在的 run 生命周期。
	repository := facts.NewMemoryRepository()
	commands := execution.NewFactCommands(repository)

	_, err := commands.Resume(context.Background(), execution.ResumeCommand{
		WorkspaceID: "ws_123",
		RunID:       "run_missing",
		ResumeRef:   "resume_ref_1",
	})
	if !errors.Is(err, execution.ErrRunNotFound) {
		t.Fatalf("err = %v, want ErrRunNotFound", err)
	}
}

type recordingRunner struct {
	runID string
}

func (runner *recordingRunner) Run(_ context.Context, runID string) error {
	runner.runID = runID
	return nil
}
