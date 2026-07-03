package product

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
)

// Projection 是 Workbench、Action API 和 Replay 的唯一产品投影入口。
type Projection interface {
	WorkbenchView(ctx context.Context, workspaceID string) (WorkbenchView, error)
	RunSnapshot(ctx context.Context, workspaceID string, runID string) (WorkbenchView, error)
	ActionResult(ctx context.Context, workspaceID string, actionID string, runID string) (ActionResult, error)
	ResumeResult(ctx context.Context, workspaceID string, runID string) (ActionResult, error)
	ReplayView(ctx context.Context, workspaceID string, runID string) (ReplayView, error)
	StreamEvents(ctx context.Context, workspaceID string, runID string) ([]StreamEvent, error)
}

// WorkbenchView 是 Workbench 前端消费的产品视图。
type WorkbenchView struct {
	SchemaVersion string         `json:"schema_version"`
	WorkspaceID   string         `json:"workspace_id"`
	RunID         string         `json:"run_id"`
	Status        string         `json:"status"`
	Timeline      []TimelineItem `json:"timeline"`
	Inspector     Inspector      `json:"inspector"`
}

// TimelineItem 是主时间线条目，必须由 Product Facts 派生。
type TimelineItem struct {
	ItemID           string                   `json:"item_id"`
	Kind             string                   `json:"kind"`
	ToolCallID       string                   `json:"tool_call_id,omitempty"`
	PendingID        string                   `json:"pending_id,omitempty"`
	Content          string                   `json:"content,omitempty"`
	Status           string                   `json:"status,omitempty"`
	SafeSummary      string                   `json:"safe_summary,omitempty"`
	InputMode        string                   `json:"input_mode,omitempty"`
	Candidates       []facts.PendingCandidate `json:"candidates,omitempty"`
	StructuredResult map[string]any           `json:"structured_result,omitempty"`
}

// Inspector 是右侧证据/结构化/运行/审计面板摘要。
type Inspector struct {
	Tabs       []string         `json:"tabs"`
	Evidence   []EvidenceItem   `json:"evidence,omitempty"`
	Structured *StructuredPanel `json:"structured,omitempty"`
	Runtime    RuntimePanel     `json:"runtime,omitempty"`
	Audit      []AuditEvent     `json:"audit,omitempty"`
}

// EvidenceItem 是 Inspector 中可展示的安全证据摘要。
type EvidenceItem struct {
	Label      string    `json:"label"`
	Value      string    `json:"value"`
	TargetRef  string    `json:"target_ref,omitempty"`
	ObservedAt time.Time `json:"observed_at,omitempty"`
}

// StructuredPanel 指向最近一次 StructuredResult 的安全投影。
type StructuredPanel struct {
	ToolCallID       string         `json:"tool_call_id,omitempty"`
	ResultRef        string         `json:"result_ref,omitempty"`
	StructuredResult map[string]any `json:"structured_result,omitempty"`
}

// RuntimePanel 展示 run 生命周期和等待态摘要，不包含 checkpoint/resume 原始材料。
type RuntimePanel struct {
	Status     string `json:"status,omitempty"`
	ModelLabel string `json:"model_label,omitempty"`
	SafeError  string `json:"safe_error,omitempty"`
	PendingID  string `json:"pending_id,omitempty"`
}

// ActionResult 是外部 Action API 返回体，必须和 Workbench 同源。
type ActionResult struct {
	SchemaVersion string        `json:"schema_version"`
	WorkspaceID   string        `json:"workspace_id"`
	ActionID      string        `json:"action_id"`
	RunID         string        `json:"run_id"`
	Status        string        `json:"status"`
	SnapshotURL   string        `json:"snapshot_url,omitempty"`
	EventsURL     string        `json:"events_url,omitempty"`
	StreamRef     string        `json:"stream_ref,omitempty"`
	FinalAnswer   string        `json:"final_answer,omitempty"`
	ResultCards   []ResultCard  `json:"result_cards"`
	Waiting       *WaitingState `json:"waiting,omitempty"`
	ApprovalRefs  []string      `json:"approval_refs,omitempty"`
	ResumeRefs    []string      `json:"resume_refs,omitempty"`
	WaitingReason string        `json:"waiting_reason,omitempty"`
	AuditRefs     []string      `json:"audit_refs"`
}

// ResultCard 是 ActionResult 中的紧凑结果卡片。
type ResultCard struct {
	CardID           string         `json:"card_id"`
	ToolCallID       string         `json:"tool_call_id,omitempty"`
	Title            string         `json:"title"`
	Status           string         `json:"status"`
	SafeSummary      string         `json:"safe_summary,omitempty"`
	StructuredResult map[string]any `json:"structured_result,omitempty"`
}

// WaitingState 是 ActionResult 暴露给非 Workbench 客户端的等待态摘要。
type WaitingState struct {
	Kind          string                   `json:"kind"`
	Question      string                   `json:"question"`
	ApprovalRefs  []string                 `json:"approval_refs,omitempty"`
	ResumeRefs    []string                 `json:"resume_refs,omitempty"`
	RiskSummary   string                   `json:"risk_summary,omitempty"`
	TargetSummary string                   `json:"target_summary,omitempty"`
	InputMode     string                   `json:"input_mode,omitempty"`
	Candidates    []facts.PendingCandidate `json:"candidates,omitempty"`
}

// ReplayView 是从 Product Facts 重建的回放视图。
type ReplayView struct {
	SchemaVersion string        `json:"schema_version"`
	WorkspaceID   string        `json:"workspace_id"`
	RunID         string        `json:"run_id"`
	Events        []any         `json:"events"`
	View          WorkbenchView `json:"view"`
}

// AuditEvent 是产品出口的审计事件，字段必须已经脱敏。
type AuditEvent struct {
	SchemaVersion string    `json:"schema_version"`
	AuditID       string    `json:"audit_id"`
	RunID         string    `json:"run_id"`
	EventType     string    `json:"event_type"`
	SafeSummary   string    `json:"safe_summary"`
	Actor         string    `json:"actor"`
	CreatedAt     time.Time `json:"created_at"`
}

// StreamEvent 是 Workbench SSE 消费的产品事件，事件 id 来自 facts 游标。
type StreamEvent struct {
	SchemaVersion string         `json:"schema_version"`
	EventID       string         `json:"event_id"`
	RunID         string         `json:"run_id"`
	Type          string         `json:"type"`
	Sequence      int64          `json:"sequence"`
	CreatedAt     time.Time      `json:"created_at"`
	Message       *MessagePatch  `json:"message,omitempty"`
	Tool          *ToolPatch     `json:"tool,omitempty"`
	Pending       map[string]any `json:"pending,omitempty"`
	Run           *RunPatch      `json:"run,omitempty"`
	View          *WorkbenchView `json:"view,omitempty"`
	Audit         *AuditEvent    `json:"audit,omitempty"`
}

// MessagePatch 是 stream 中的消息变更，不携带模型内部事件。
type MessagePatch struct {
	MessageID string `json:"message_id"`
	Role      string `json:"role"`
	Content   string `json:"content,omitempty"`
	Status    string `json:"status,omitempty"`
}

// ToolPatch 是 stream 中的工具卡变更，只携带 StructuredResult 安全投影。
type ToolPatch struct {
	ToolCallID       string         `json:"tool_call_id"`
	ToolID           string         `json:"tool_id,omitempty"`
	DisplayName      string         `json:"display_name,omitempty"`
	Status           string         `json:"status"`
	SafeSummary      string         `json:"safe_summary,omitempty"`
	StructuredResult map[string]any `json:"structured_result,omitempty"`
}

// RunPatch 是 stream 中的 run 生命周期变更。
type RunPatch struct {
	Status    string    `json:"status"`
	SafeError string    `json:"safe_error,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// EmptyProjection 是早期 Phase 的空投影替身，不作为产品事实来源。
type EmptyProjection struct{}

// NewEmptyProjection 创建空投影替身。
func NewEmptyProjection() EmptyProjection {
	return EmptyProjection{}
}

// WorkbenchView 返回空 Workbench 视图。
func (projection EmptyProjection) WorkbenchView(_ context.Context, workspaceID string) (WorkbenchView, error) {
	return newWorkbenchView(workspaceID, "run_initial"), nil
}

// RunSnapshot 返回空 run 快照。
func (projection EmptyProjection) RunSnapshot(_ context.Context, workspaceID string, runID string) (WorkbenchView, error) {
	return newWorkbenchView(workspaceID, runID), nil
}

// ActionResult 返回空 Action API 结果。
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

// ResumeResult 返回空 resume 结果。
func (projection EmptyProjection) ResumeResult(ctx context.Context, workspaceID string, runID string) (ActionResult, error) {
	return projection.ActionResult(ctx, workspaceID, "resume", runID)
}

// ReplayView 返回空 replay 视图。
func (projection EmptyProjection) ReplayView(_ context.Context, workspaceID string, runID string) (ReplayView, error) {
	return ReplayView{
		SchemaVersion: "eino_replay_view.v1",
		WorkspaceID:   workspaceID,
		RunID:         runID,
		Events:        []any{},
		View:          newWorkbenchView(workspaceID, runID),
	}, nil
}

// StreamEvents 返回空投影替身的最小 run 事件。
func (projection EmptyProjection) StreamEvents(_ context.Context, _ string, runID string) ([]StreamEvent, error) {
	return []StreamEvent{newRunStreamEvent(runID, facts.Run{RunID: runID, Status: facts.RunStatusCreated}, 1)}, nil
}

// newWorkbenchView 构建基础 Workbench 视图，后续由 facts projection 填充时间线。
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

// FactsProjection 从 Product Facts repository 读取事实并投影到产品 API。
type FactsProjection struct {
	repository facts.Repository
}

// NewFactsProjection 创建基于 Product Facts 的产品投影。
func NewFactsProjection(repository facts.Repository) FactsProjection {
	return FactsProjection{repository: repository}
}

// WorkbenchView 读取 workspace 最新 run 并返回 Workbench 视图。
func (projection FactsProjection) WorkbenchView(ctx context.Context, workspaceID string) (WorkbenchView, error) {
	run, err := projection.repository.LatestRun(ctx, workspaceID)
	if errors.Is(err, facts.ErrNotFound) {
		return newWorkbenchView(workspaceID, ""), nil
	}
	if err != nil {
		return WorkbenchView{}, err
	}
	return projection.RunSnapshot(ctx, run.WorkspaceID, run.RunID)
}

// RunSnapshot 返回指定 run 的 Workbench 快照。
func (projection FactsProjection) RunSnapshot(ctx context.Context, workspaceID string, runID string) (WorkbenchView, error) {
	snapshot, err := projection.snapshotForWorkspace(ctx, workspaceID, runID)
	if err != nil {
		return WorkbenchView{}, err
	}
	return workbenchFromSnapshot(snapshot), nil
}

// ActionResult 从同一 run facts 投影 Action API 结果。
func (projection FactsProjection) ActionResult(ctx context.Context, workspaceID string, actionID string, runID string) (ActionResult, error) {
	snapshot, err := projection.snapshotForWorkspace(ctx, workspaceID, runID)
	if err != nil {
		return ActionResult{}, err
	}
	return actionResultFromSnapshot(workspaceID, actionID, snapshot), nil
}

// ResumeResult 复用 ActionResult 路径，避免 resume 产生第二套事实。
func (projection FactsProjection) ResumeResult(ctx context.Context, workspaceID string, runID string) (ActionResult, error) {
	return projection.ActionResult(ctx, workspaceID, "resume", runID)
}

// ReplayView 基于 run snapshot 构建回放视图。
func (projection FactsProjection) ReplayView(ctx context.Context, workspaceID string, runID string) (ReplayView, error) {
	view, err := projection.RunSnapshot(ctx, workspaceID, runID)
	if err != nil {
		return ReplayView{}, err
	}
	events, err := projection.StreamEvents(ctx, workspaceID, runID)
	if err != nil {
		return ReplayView{}, err
	}
	replayEvents := make([]any, 0, len(events)+len(view.Inspector.Audit))
	for _, event := range events {
		replayEvents = append(replayEvents, event)
	}
	for _, event := range view.Inspector.Audit {
		replayEvents = append(replayEvents, event)
	}
	return ReplayView{
		SchemaVersion: "eino_replay_view.v1",
		WorkspaceID:   view.WorkspaceID,
		RunID:         view.RunID,
		Events:        replayEvents,
		View:          view,
	}, nil
}

// StreamEvents 从 Product Facts 构造 SSE 事件，event_id 使用 run_id + 单调序号。
func (projection FactsProjection) StreamEvents(ctx context.Context, workspaceID string, runID string) ([]StreamEvent, error) {
	snapshot, err := projection.snapshotForWorkspace(ctx, workspaceID, runID)
	if err != nil {
		return nil, err
	}
	return streamEventsFromSnapshot(snapshot), nil
}

// snapshotForWorkspace 是所有按 run_id 查询产品出口的 workspace 隔离边界。
func (projection FactsProjection) snapshotForWorkspace(ctx context.Context, workspaceID string, runID string) (facts.Snapshot, error) {
	snapshot, err := projection.repository.GetSnapshot(ctx, runID)
	if err != nil {
		return facts.Snapshot{}, err
	}
	if snapshot.Run.WorkspaceID != workspaceID {
		return facts.Snapshot{}, facts.ErrNotFound
	}
	return snapshot, nil
}

func actionResultFromSnapshot(workspaceID string, actionID string, snapshot facts.Snapshot) ActionResult {
	run := snapshot.Run
	if run.WorkspaceID == "" {
		run.WorkspaceID = workspaceID
	}
	result := ActionResult{
		SchemaVersion: "eino_action_result.v1",
		WorkspaceID:   run.WorkspaceID,
		ActionID:      actionID,
		RunID:         run.RunID,
		Status:        actionStatus(run.Status),
		SnapshotURL:   fmt.Sprintf("/api/workspaces/%s/runs/%s", run.WorkspaceID, run.RunID),
		EventsURL:     fmt.Sprintf("/api/workspaces/%s/runs/%s/stream", run.WorkspaceID, run.RunID),
		StreamRef:     fmt.Sprintf("sse:/api/workspaces/%s/runs/%s/stream", run.WorkspaceID, run.RunID),
		ResultCards:   resultCards(snapshot),
		AuditRefs:     auditRefs(run.WorkspaceID, run.RunID, snapshot.AuditEvents),
	}
	result.FinalAnswer = finalAssistantAnswer(snapshot.Turns)
	if pending, ok := latestWaitingPending(snapshot.PendingInteractions); ok && run.Status == facts.RunStatusWaiting {
		result.WaitingReason = waitingReason(pending)
		result.Waiting = waitingState(pending)
		if pending.Kind == facts.PendingKindApproval {
			result.ApprovalRefs = []string{pending.ResumeRef}
			result.ResumeRefs = []string{}
		} else {
			result.ApprovalRefs = []string{}
			result.ResumeRefs = []string{pending.ResumeRef}
		}
	}
	return result
}

func workbenchFromSnapshot(snapshot facts.Snapshot) WorkbenchView {
	view := newWorkbenchView(snapshot.Run.WorkspaceID, snapshot.Run.RunID)
	view.Status = string(snapshot.Run.Status)
	view.Timeline = timelineItems(snapshot)
	view.Inspector = inspectorFromSnapshot(snapshot)
	return view
}

func timelineItems(snapshot facts.Snapshot) []TimelineItem {
	items := []TimelineItem{}
	turns := append([]facts.Turn(nil), snapshot.Turns...)
	sort.SliceStable(turns, func(i, j int) bool { return turns[i].Sequence < turns[j].Sequence })
	for _, turn := range turns {
		if turn.Role == facts.TurnRoleUser {
			items = append(items, turnTimelineItem(turn))
		}
	}
	for _, call := range snapshot.ToolCalls {
		items = append(items, toolTimelineItem(call, resultForCall(snapshot.ToolResults, call.ToolCallID)))
	}
	for _, pending := range snapshot.PendingInteractions {
		items = append(items, pendingTimelineItem(pending))
	}
	for _, turn := range turns {
		if turn.Role != facts.TurnRoleUser {
			items = append(items, turnTimelineItem(turn))
		}
	}
	if snapshot.Run.SafeError != "" {
		items = append(items, TimelineItem{
			ItemID:  snapshot.Run.RunID + ":notice",
			Kind:    "run_notice",
			Content: snapshot.Run.SafeError,
			Status:  string(snapshot.Run.Status),
		})
	}
	return items
}

func turnTimelineItem(turn facts.Turn) TimelineItem {
	kind := "assistant_message"
	if turn.Role == facts.TurnRoleUser {
		kind = "user_message"
	}
	if turn.Role == facts.TurnRoleSystemNotice {
		kind = "run_notice"
	}
	return TimelineItem{
		ItemID:  turn.TurnID,
		Kind:    kind,
		Content: turn.Content,
	}
}

func toolTimelineItem(call facts.ToolCall, result *facts.ToolResult) TimelineItem {
	item := TimelineItem{
		ItemID:     call.ToolCallID,
		Kind:       "tool_card",
		ToolCallID: call.ToolCallID,
		Status:     string(call.Status),
	}
	if result != nil {
		item.Status = string(result.Status)
		item.SafeSummary = result.StructuredResult.SafeSummary
		item.StructuredResult = structuredResult(result.StructuredResult)
	}
	if item.SafeSummary == "" {
		item.SafeSummary = call.ArgsPreview
	}
	return item
}

func pendingTimelineItem(pending facts.PendingInteraction) TimelineItem {
	kind := "clarification_card"
	content := pending.Question
	if pending.Kind == facts.PendingKindApproval {
		kind = "approval_card"
		content = pending.RiskSummary
	}
	return TimelineItem{
		ItemID:     pending.PendingID,
		Kind:       kind,
		PendingID:  pending.PendingID,
		Content:    content,
		Status:     string(pending.Status),
		InputMode:  string(pending.InputMode),
		Candidates: pending.Candidates,
	}
}

func inspectorFromSnapshot(snapshot facts.Snapshot) Inspector {
	inspector := Inspector{
		Tabs:     []string{"evidence", "structured", "runtime", "audit"},
		Evidence: evidenceItems(snapshot),
		Runtime: RuntimePanel{
			Status:     string(snapshot.Run.Status),
			ModelLabel: snapshot.Run.ModelLabel,
			SafeError:  snapshot.Run.SafeError,
		},
		Audit: auditEvents(snapshot.AuditEvents),
	}
	if pending, ok := latestWaitingPending(snapshot.PendingInteractions); ok {
		inspector.Runtime.PendingID = pending.PendingID
	}
	if call, result, ok := latestStructuredResult(snapshot); ok {
		inspector.Structured = &StructuredPanel{
			ToolCallID:       call.ToolCallID,
			ResultRef:        result.StructuredResult.ResultRef,
			StructuredResult: structuredResult(result.StructuredResult),
		}
	}
	return inspector
}

func evidenceItems(snapshot facts.Snapshot) []EvidenceItem {
	items := []EvidenceItem{}
	for _, call := range snapshot.ToolCalls {
		label := call.DisplayName
		if label == "" {
			label = call.ToolID
		}
		items = append(items, EvidenceItem{
			Label:      "工具",
			Value:      label,
			TargetRef:  call.ToolCallID,
			ObservedAt: call.CreatedAt,
		})
	}
	return items
}

func resultCards(snapshot facts.Snapshot) []ResultCard {
	cards := []ResultCard{}
	for _, call := range snapshot.ToolCalls {
		result := resultForCall(snapshot.ToolResults, call.ToolCallID)
		card := ResultCard{
			CardID:     call.ToolCallID,
			ToolCallID: call.ToolCallID,
			Title:      call.DisplayName,
			Status:     string(call.Status),
		}
		if card.Title == "" {
			card.Title = call.ToolID
		}
		if result != nil {
			card.Status = string(result.Status)
			card.SafeSummary = result.StructuredResult.SafeSummary
			card.StructuredResult = structuredResult(result.StructuredResult)
		}
		cards = append(cards, card)
	}
	return cards
}

func streamEventsFromSnapshot(snapshot facts.Snapshot) []StreamEvent {
	events := []StreamEvent{newRunStreamEvent(snapshot.Run.RunID, snapshot.Run, 1)}
	sequence := int64(2)
	turns := append([]facts.Turn(nil), snapshot.Turns...)
	sort.SliceStable(turns, func(i, j int) bool { return turns[i].Sequence < turns[j].Sequence })
	for _, turn := range turns {
		events = append(events, StreamEvent{
			SchemaVersion: "eino_workbench_stream_event.v1",
			EventID:       eventID(snapshot.Run.RunID, sequence),
			RunID:         snapshot.Run.RunID,
			Type:          "message.updated",
			Sequence:      sequence,
			CreatedAt:     nonZeroTime(turn.CreatedAt, snapshot.Run.UpdatedAt),
			Message: &MessagePatch{
				MessageID: turn.TurnID,
				Role:      string(turn.Role),
				Content:   turn.Content,
				Status:    "completed",
			},
		})
		sequence++
	}
	for _, call := range snapshot.ToolCalls {
		result := resultForCall(snapshot.ToolResults, call.ToolCallID)
		patch := &ToolPatch{
			ToolCallID:  call.ToolCallID,
			ToolID:      call.ToolID,
			DisplayName: call.DisplayName,
			Status:      string(call.Status),
		}
		if result != nil {
			patch.Status = string(result.Status)
			patch.SafeSummary = result.StructuredResult.SafeSummary
			patch.StructuredResult = structuredResult(result.StructuredResult)
		}
		events = append(events, StreamEvent{
			SchemaVersion: "eino_workbench_stream_event.v1",
			EventID:       eventID(snapshot.Run.RunID, sequence),
			RunID:         snapshot.Run.RunID,
			Type:          "tool.updated",
			Sequence:      sequence,
			CreatedAt:     nonZeroTime(call.EndedAt, nonZeroTime(call.CreatedAt, snapshot.Run.UpdatedAt)),
			Tool:          patch,
		})
		sequence++
	}
	for _, pending := range snapshot.PendingInteractions {
		events = append(events, StreamEvent{
			SchemaVersion: "eino_workbench_stream_event.v1",
			EventID:       eventID(snapshot.Run.RunID, sequence),
			RunID:         snapshot.Run.RunID,
			Type:          "pending.updated",
			Sequence:      sequence,
			CreatedAt:     snapshot.Run.UpdatedAt,
			Pending:       pendingPatch(pending),
		})
		sequence++
	}
	return events
}

func newRunStreamEvent(runID string, run facts.Run, sequence int64) StreamEvent {
	return StreamEvent{
		SchemaVersion: "eino_workbench_stream_event.v1",
		EventID:       eventID(runID, sequence),
		RunID:         runID,
		Type:          "run.updated",
		Sequence:      sequence,
		CreatedAt:     nonZeroTime(run.UpdatedAt, time.Now().UTC()),
		Run: &RunPatch{
			Status:    string(run.Status),
			SafeError: run.SafeError,
			UpdatedAt: run.UpdatedAt,
		},
	}
}

func eventID(runID string, sequence int64) string {
	return fmt.Sprintf("%s:%06d", runID, sequence)
}

func resultForCall(results []facts.ToolResult, toolCallID string) *facts.ToolResult {
	for i := range results {
		if results[i].ToolCallID == toolCallID {
			return &results[i]
		}
	}
	return nil
}

func structuredResult(result facts.StructuredResultRef) map[string]any {
	status := "resolved"
	if result.SafeSummary == "" {
		status = "empty"
	}
	return map[string]any{
		"schema_version": result.SchemaVersion,
		"status":         status,
		"data": map[string]any{
			"summary": result.SafeSummary,
			"facts":   []any{},
		},
		"metadata": map[string]any{
			"safe":       true,
			"result_ref": result.ResultRef,
			"source":     "product_facts",
		},
	}
}

func latestWaitingPending(pending []facts.PendingInteraction) (facts.PendingInteraction, bool) {
	for i := len(pending) - 1; i >= 0; i-- {
		if pending[i].Status == facts.PendingStatusWaiting {
			return pending[i], true
		}
	}
	return facts.PendingInteraction{}, false
}

func latestStructuredResult(snapshot facts.Snapshot) (facts.ToolCall, facts.ToolResult, bool) {
	for i := len(snapshot.ToolCalls) - 1; i >= 0; i-- {
		call := snapshot.ToolCalls[i]
		if result := resultForCall(snapshot.ToolResults, call.ToolCallID); result != nil {
			return call, *result, true
		}
	}
	return facts.ToolCall{}, facts.ToolResult{}, false
}

func auditEvents(events []facts.AuditEvent) []AuditEvent {
	result := make([]AuditEvent, 0, len(events))
	for _, event := range events {
		result = append(result, AuditEvent{
			SchemaVersion: "eino_audit_event.v1",
			AuditID:       event.AuditID,
			RunID:         event.RunID,
			EventType:     event.EventType,
			SafeSummary:   event.SafeSummary,
			Actor:         event.Actor,
			CreatedAt:     event.CreatedAt,
		})
	}
	return result
}

func auditRefs(workspaceID string, runID string, events []facts.AuditEvent) []string {
	if len(events) == 0 {
		return []string{}
	}
	return []string{fmt.Sprintf("/api/workspaces/%s/runs/%s/replay", workspaceID, runID)}
}

func finalAssistantAnswer(turns []facts.Turn) string {
	for i := len(turns) - 1; i >= 0; i-- {
		if turns[i].Role == facts.TurnRoleAssistant {
			return turns[i].Content
		}
	}
	return ""
}

func actionStatus(status facts.RunStatus) string {
	switch status {
	case facts.RunStatusCreated:
		return "accepted"
	case facts.RunStatusRunning:
		return "running"
	case facts.RunStatusWaiting:
		return "waiting"
	case facts.RunStatusSucceeded:
		return "completed"
	case facts.RunStatusFailed:
		return "failed"
	case facts.RunStatusStopped, facts.RunStatusCancelled:
		return "stopped"
	default:
		return "accepted"
	}
}

func waitingState(pending facts.PendingInteraction) *WaitingState {
	state := &WaitingState{
		Kind:          string(pending.Kind),
		Question:      pending.Question,
		RiskSummary:   pending.RiskSummary,
		TargetSummary: pending.TargetSummary,
		InputMode:     string(pending.InputMode),
		Candidates:    pending.Candidates,
	}
	if state.Question == "" && pending.Kind == facts.PendingKindApproval {
		state.Question = "是否批准继续执行？"
	}
	if state.Question == "" {
		state.Question = "请补充信息后继续。"
	}
	if pending.Kind == facts.PendingKindApproval {
		state.ApprovalRefs = []string{pending.ResumeRef}
		return state
	}
	state.ResumeRefs = []string{pending.ResumeRef}
	return state
}

func waitingReason(pending facts.PendingInteraction) string {
	if pending.Kind == facts.PendingKindApproval {
		return "approval_required"
	}
	return "clarification_required"
}

func pendingPatch(pending facts.PendingInteraction) map[string]any {
	patch := map[string]any{
		"schema_version": "eino_workbench_pending_interaction.v1",
		"pending_id":     pending.PendingID,
		"run_id":         pending.RunID,
		"kind":           string(pending.Kind),
		"status":         string(pending.Status),
		"question":       pending.Question,
		"operation_name": pending.OperationName,
		"risk_summary":   pending.RiskSummary,
		"target_summary": pending.TargetSummary,
		"resume_ref":     pending.ResumeRef,
	}
	if pending.InputMode != "" {
		patch["input_mode"] = string(pending.InputMode)
	}
	if len(pending.Candidates) > 0 {
		patch["candidates"] = pending.Candidates
	}
	return patch
}

func nonZeroTime(value time.Time, fallback time.Time) time.Time {
	if !value.IsZero() {
		return value
	}
	return fallback
}
