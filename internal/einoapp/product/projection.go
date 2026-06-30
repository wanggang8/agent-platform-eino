package product

import (
	"context"
	"errors"

	"agent-platform-eino/internal/einoapp/facts"
)

type Projection interface {
	WorkbenchView(ctx context.Context, workspaceID string) (WorkbenchView, error)
	RunSnapshot(ctx context.Context, workspaceID string, runID string) (WorkbenchView, error)
	ActionResult(ctx context.Context, workspaceID string, actionID string, runID string) (ActionResult, error)
	ResumeResult(ctx context.Context, workspaceID string, runID string) (ActionResult, error)
	ReplayView(ctx context.Context, workspaceID string, runID string) (ReplayView, error)
}

type WorkbenchView struct {
	SchemaVersion string         `json:"schema_version"`
	WorkspaceID   string         `json:"workspace_id"`
	RunID         string         `json:"run_id"`
	Status        string         `json:"status"`
	Timeline      []TimelineItem `json:"timeline"`
	Inspector     Inspector      `json:"inspector"`
}

type TimelineItem struct {
	ItemID  string `json:"item_id"`
	Kind    string `json:"kind"`
	Content string `json:"content,omitempty"`
	Status  string `json:"status,omitempty"`
}

type Inspector struct {
	Tabs []string `json:"tabs"`
}

type ActionResult struct {
	SchemaVersion string       `json:"schema_version"`
	WorkspaceID   string       `json:"workspace_id"`
	ActionID      string       `json:"action_id"`
	RunID         string       `json:"run_id"`
	Status        string       `json:"status"`
	ResultCards   []ResultCard `json:"result_cards"`
	AuditRefs     []string     `json:"audit_refs"`
}

type ResultCard struct {
	CardID string `json:"card_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type ReplayView struct {
	SchemaVersion string        `json:"schema_version"`
	WorkspaceID   string        `json:"workspace_id"`
	RunID         string        `json:"run_id"`
	Events        []any         `json:"events"`
	View          WorkbenchView `json:"view"`
}

type EmptyProjection struct{}

func NewEmptyProjection() EmptyProjection {
	return EmptyProjection{}
}

func (projection EmptyProjection) WorkbenchView(_ context.Context, workspaceID string) (WorkbenchView, error) {
	return newWorkbenchView(workspaceID, "run_initial"), nil
}

func (projection EmptyProjection) RunSnapshot(_ context.Context, workspaceID string, runID string) (WorkbenchView, error) {
	return newWorkbenchView(workspaceID, runID), nil
}

func (projection EmptyProjection) ActionResult(_ context.Context, workspaceID string, actionID string, runID string) (ActionResult, error) {
	return ActionResult{
		SchemaVersion: "eino_action_result.v1",
		WorkspaceID:   workspaceID,
		ActionID:      actionID,
		RunID:         runID,
		Status:        "accepted",
		ResultCards:   []ResultCard{},
		AuditRefs:     []string{},
	}, nil
}

func (projection EmptyProjection) ResumeResult(ctx context.Context, workspaceID string, runID string) (ActionResult, error) {
	return projection.ActionResult(ctx, workspaceID, "resume", runID)
}

func (projection EmptyProjection) ReplayView(_ context.Context, workspaceID string, runID string) (ReplayView, error) {
	return ReplayView{
		SchemaVersion: "eino_replay_view.v1",
		WorkspaceID:   workspaceID,
		RunID:         runID,
		Events:        []any{},
		View:          newWorkbenchView(workspaceID, runID),
	}, nil
}

func newWorkbenchView(workspaceID string, runID string) WorkbenchView {
	return WorkbenchView{
		SchemaVersion: "eino_workbench_view.v1",
		WorkspaceID:   workspaceID,
		RunID:         runID,
		Status:        "created",
		Timeline:      []TimelineItem{},
		Inspector: Inspector{
			Tabs: []string{"evidence", "structured", "runtime", "audit"},
		},
	}
}

type FactsProjection struct {
	repository facts.Repository
}

func NewFactsProjection(repository facts.Repository) FactsProjection {
	return FactsProjection{repository: repository}
}

func (projection FactsProjection) WorkbenchView(ctx context.Context, workspaceID string) (WorkbenchView, error) {
	run, err := projection.repository.LatestRun(ctx, workspaceID)
	if errors.Is(err, facts.ErrNotFound) {
		return newWorkbenchView(workspaceID, ""), nil
	}
	if err != nil {
		return WorkbenchView{}, err
	}
	return newWorkbenchView(run.WorkspaceID, run.RunID), nil
}

func (projection FactsProjection) RunSnapshot(ctx context.Context, workspaceID string, runID string) (WorkbenchView, error) {
	run, err := projection.repository.GetRun(ctx, runID)
	if errors.Is(err, facts.ErrNotFound) {
		return newWorkbenchView(workspaceID, runID), nil
	}
	if err != nil {
		return WorkbenchView{}, err
	}
	return newWorkbenchView(run.WorkspaceID, run.RunID), nil
}

func (projection FactsProjection) ActionResult(ctx context.Context, workspaceID string, actionID string, runID string) (ActionResult, error) {
	run, err := projection.repository.GetRun(ctx, runID)
	if errors.Is(err, facts.ErrNotFound) {
		run = facts.Run{WorkspaceID: workspaceID, RunID: runID}
	} else if err != nil {
		return ActionResult{}, err
	}
	return ActionResult{
		SchemaVersion: "eino_action_result.v1",
		WorkspaceID:   run.WorkspaceID,
		ActionID:      actionID,
		RunID:         run.RunID,
		Status:        "accepted",
		ResultCards:   []ResultCard{},
		AuditRefs:     []string{},
	}, nil
}

func (projection FactsProjection) ResumeResult(ctx context.Context, workspaceID string, runID string) (ActionResult, error) {
	return projection.ActionResult(ctx, workspaceID, "resume", runID)
}

func (projection FactsProjection) ReplayView(ctx context.Context, workspaceID string, runID string) (ReplayView, error) {
	view, err := projection.RunSnapshot(ctx, workspaceID, runID)
	if err != nil {
		return ReplayView{}, err
	}
	return ReplayView{
		SchemaVersion: "eino_replay_view.v1",
		WorkspaceID:   view.WorkspaceID,
		RunID:         view.RunID,
		Events:        []any{},
		View:          view,
	}, nil
}
