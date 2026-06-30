package facts_test

import (
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

func TestToolResultFactRequiresStructuredResultOnly(t *testing.T) {
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
