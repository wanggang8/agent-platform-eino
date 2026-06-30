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

func TestFactCommandsResumeRequiresExistingRun(t *testing.T) {
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
