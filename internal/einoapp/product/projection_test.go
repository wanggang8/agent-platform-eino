package product_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

func TestEmptyProjectionReturnsWorkbenchViewSchemaDirectly(t *testing.T) {
	projection := product.NewEmptyProjection()

	view, err := projection.WorkbenchView(context.Background(), "ws-demo")
	if err != nil {
		t.Fatal(err)
	}

	if view.SchemaVersion != "eino_workbench_view.v1" {
		t.Fatalf("schema_version = %q", view.SchemaVersion)
	}
	if view.WorkspaceID != "ws-demo" {
		t.Fatalf("workspace_id = %q", view.WorkspaceID)
	}
	if view.Timeline == nil {
		t.Fatal("timeline must be an empty array, not nil")
	}
}

func TestProductErrorConvertsToSafeAPIError(t *testing.T) {
	err := product.NewSafeError("projection_unavailable", "投影暂不可用", false)

	if err.Code != "projection_unavailable" {
		t.Fatalf("code = %q", err.Code)
	}
	if err.SafeDetail == "" {
		t.Fatal("safe detail is empty")
	}
}

func TestFactsProjectionReadsLatestRunFromProductFacts(t *testing.T) {
	repository := facts.NewMemoryRepository()
	if err := repository.CreateRun(context.Background(), facts.Run{
		RunID:       "run-facts",
		WorkspaceID: "ws-demo",
		Status:      facts.RunStatusCreated,
	}); err != nil {
		t.Fatal(err)
	}
	projection := product.NewFactsProjection(repository)

	view, err := projection.WorkbenchView(context.Background(), "ws-demo")
	if err != nil {
		t.Fatal(err)
	}

	if view.RunID != "run-facts" {
		t.Fatalf("run_id = %q", view.RunID)
	}
}

func TestFactsProjectionBuildsWorkbenchActionReplayAndStreamFromSnapshot(t *testing.T) {
	// 产品投影测试同一组 facts 的多个出口，防止 Workbench、Action API、Replay 和 SSE 分裂。
	ctx := context.Background()
	repository := facts.NewMemoryRepository()
	now := time.Unix(100, 0).UTC()
	if err := repository.CreateRun(ctx, facts.Run{
		RunID:       "run-rich",
		WorkspaceID: "ws-demo",
		Status:      facts.RunStatusWaiting,
		CreatedAt:   now,
		UpdatedAt:   now.Add(time.Second),
		ModelLabel:  "Mock Chat",
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendTurn(ctx, facts.Turn{TurnID: "turn-user", RunID: "run-rich", Role: facts.TurnRoleUser, Content: "查询资产", Sequence: 1, CreatedAt: now}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendTurn(ctx, facts.Turn{TurnID: "turn-assistant", RunID: "run-rich", Role: facts.TurnRoleAssistant, Content: "需要审批", Sequence: 4, CreatedAt: now}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendToolCall(ctx, facts.ToolCall{ToolCallID: "call-1", RunID: "run-rich", ToolID: "cap.asset.read", DisplayName: "资产查询", Status: facts.ToolCallSucceeded, ArgsPreview: "scope=own", CreatedAt: now}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendToolResult(ctx, facts.ToolResult{
		ResultID:   "result-1",
		ToolCallID: "call-1",
		Status:     facts.ToolResultSucceeded,
		StructuredResult: facts.StructuredResultRef{
			SchemaVersion: "tool.structured_result.v1",
			ResultRef:     "result:call-1",
			SafeSummary:   "发现 1 个资产",
		},
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendPendingInteraction(ctx, facts.PendingInteraction{
		PendingID:   "pending-1",
		RunID:       "run-rich",
		Kind:        facts.PendingKindApproval,
		Status:      facts.PendingStatusWaiting,
		ResumeRef:   "resume-safe-1",
		RiskSummary: "写域操作需要审批",
	}); err != nil {
		t.Fatal(err)
	}
	if err := repository.AppendAuditEvent(ctx, facts.AuditEvent{AuditID: "audit-1", RunID: "run-rich", EventType: "tool", SafeSummary: "tool completed", Actor: "tool", CreatedAt: now}); err != nil {
		t.Fatal(err)
	}

	projection := product.NewFactsProjection(repository)
	view, err := projection.RunSnapshot(ctx, "ws-demo", "run-rich")
	if err != nil {
		t.Fatal(err)
	}
	if view.Status != "waiting" {
		t.Fatalf("view status = %q", view.Status)
	}
	if len(view.Timeline) != 4 {
		t.Fatalf("timeline = %+v", view.Timeline)
	}
	if view.Timeline[1].Kind != "tool_card" || view.Timeline[1].SafeSummary != "发现 1 个资产" {
		t.Fatalf("tool card not projected from StructuredResult: %+v", view.Timeline[1])
	}
	if view.Timeline[2].Kind != "approval_card" || view.Timeline[2].PendingID != "pending-1" {
		t.Fatalf("pending card not projected: %+v", view.Timeline[2])
	}
	if view.Inspector.Runtime.Status != "waiting" || view.Inspector.Runtime.PendingID != "pending-1" {
		t.Fatalf("runtime inspector mismatch: %+v", view.Inspector.Runtime)
	}
	if len(view.Inspector.Audit) != 1 || view.Inspector.Audit[0].AuditID != "audit-1" {
		t.Fatalf("audit not projected: %+v", view.Inspector.Audit)
	}

	action, err := projection.ActionResult(ctx, "ws-demo", "action-demo", "run-rich")
	if err != nil {
		t.Fatal(err)
	}
	if action.Status != "waiting" || len(action.ResultCards) != 1 || action.ResultCards[0].SafeSummary != "发现 1 个资产" {
		t.Fatalf("action result not facts-backed: %+v", action)
	}
	if len(action.ApprovalRefs) != 1 || action.ApprovalRefs[0] != "resume-safe-1" {
		t.Fatalf("approval refs = %+v", action.ApprovalRefs)
	}

	replay, err := projection.ReplayView(ctx, "ws-demo", "run-rich")
	if err != nil {
		t.Fatal(err)
	}
	if len(replay.Events) == 0 || replay.View.RunID != "run-rich" {
		t.Fatalf("replay not built from facts: %+v", replay)
	}

	events, err := projection.StreamEvents(ctx, "ws-demo", "run-rich")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 || events[0].EventID != "run-rich:000001" {
		t.Fatalf("stream events not facts cursor based: %+v", events)
	}
}

func TestFactsProjectionDoesNotSynthesizeSpecificMissingRun(t *testing.T) {
	// current view 可以有空态；指定 run 的 snapshot/action/stream 不能伪造不存在的 Product Facts。
	projection := product.NewFactsProjection(facts.NewMemoryRepository())
	ctx := context.Background()

	if _, err := projection.RunSnapshot(ctx, "ws-demo", "run-missing"); !errors.Is(err, facts.ErrNotFound) {
		t.Fatalf("RunSnapshot err = %v, want facts.ErrNotFound", err)
	}
	if _, err := projection.ActionResult(ctx, "ws-demo", "action-demo", "run-missing"); !errors.Is(err, facts.ErrNotFound) {
		t.Fatalf("ActionResult err = %v, want facts.ErrNotFound", err)
	}
	if _, err := projection.StreamEvents(ctx, "ws-demo", "run-missing"); !errors.Is(err, facts.ErrNotFound) {
		t.Fatalf("StreamEvents err = %v, want facts.ErrNotFound", err)
	}
}
