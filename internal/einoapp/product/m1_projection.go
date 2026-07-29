package product

import (
	"context"
	"errors"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
)

const (
	M1WorkbenchViewSchemaVersion  = "eino_workbench_view.v2"
	M1WorkbenchEventSchemaVersion = "eino_workbench_stream_event.v2"
)

// M1ProjectionPort 是 HTTP/Workbench 共用的唯一只读产品投影端口。
type M1ProjectionPort interface {
	WorkbenchView(ctx context.Context, workspaceID string) (M1WorkbenchView, error)
	RunSnapshot(ctx context.Context, workspaceID, runID string) (M1WorkbenchView, error)
	StreamEvents(ctx context.Context, workspaceID, runID string) ([]M1StreamEvent, error)
}

type M1Message struct {
	MessageID string `json:"message_id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
}

// M1ResultCard 只展示查询序号、摘要、数量和观察时间。
type M1ResultCard struct {
	QuerySequence int64     `json:"query_sequence"`
	Summary       string    `json:"summary"`
	Count         int       `json:"count"`
	ObservedAt    time.Time `json:"observed_at"`
}

type M1RightPanel struct {
	Tabs         [2]string `json:"tabs"`
	EmptyMessage string    `json:"empty_message"`
}

// M1WorkbenchView 是 facts 的安全产品投影，不包含任何内部引用或行级技术材料。
type M1WorkbenchView struct {
	SchemaVersion  string        `json:"schema_version"`
	WorkspaceID    string        `json:"workspace_id"`
	ConversationID string        `json:"conversation_id"`
	RunID          string        `json:"run_id"`
	Status         string        `json:"status"`
	Messages       []M1Message   `json:"messages"`
	Result         *M1ResultCard `json:"result"`
	SafeError      string        `json:"safe_error,omitempty"`
	RightPanel     M1RightPanel  `json:"right_panel"`
}

// M1StreamEvent 的 id/sequence 直接使用持久 FactEvent。
type M1StreamEvent struct {
	SchemaVersion string          `json:"schema_version"`
	EventID       string          `json:"event_id"`
	RunID         string          `json:"run_id"`
	Sequence      int64           `json:"sequence"`
	CreatedAt     time.Time       `json:"created_at"`
	Type          string          `json:"type"`
	View          M1WorkbenchView `json:"view"`
}

type M1Projection struct{ repository facts.QueryRepository }

func NewM1Projection(repository facts.QueryRepository) M1Projection {
	return M1Projection{repository: repository}
}

func (projection M1Projection) WorkbenchView(ctx context.Context, workspaceID string) (M1WorkbenchView, error) {
	aggregate, err := projection.repository.LatestQuery(ctx, workspaceID)
	if errors.Is(err, facts.ErrNotFound) {
		return newM1IdleView(workspaceID), nil
	}
	if err != nil {
		return M1WorkbenchView{}, err
	}
	return projectM1Aggregate(aggregate), nil
}

func (projection M1Projection) RunSnapshot(ctx context.Context, workspaceID, runID string) (M1WorkbenchView, error) {
	aggregate, err := projection.repository.GetQuery(ctx, workspaceID, runID)
	if err != nil {
		return M1WorkbenchView{}, err
	}
	return projectM1Aggregate(aggregate), nil
}

func (projection M1Projection) StreamEvents(ctx context.Context, workspaceID, runID string) ([]M1StreamEvent, error) {
	aggregate, err := projection.repository.GetQuery(ctx, workspaceID, runID)
	if err != nil {
		return nil, err
	}
	view := projectM1Aggregate(aggregate)
	events := make([]M1StreamEvent, 0, len(aggregate.Events()))
	for _, factEvent := range aggregate.Events() {
		events = append(events, M1StreamEvent{
			SchemaVersion: M1WorkbenchEventSchemaVersion, EventID: factEvent.EventID(), RunID: aggregate.RunID(),
			Sequence: factEvent.Sequence(), CreatedAt: factEvent.CreatedAt(), Type: "view.replaced", View: view,
		})
	}
	return events, nil
}

func projectM1Aggregate(aggregate facts.QueryFacts) M1WorkbenchView {
	view := M1WorkbenchView{
		SchemaVersion: M1WorkbenchViewSchemaVersion, WorkspaceID: aggregate.WorkspaceID(), ConversationID: aggregate.ConversationID(), RunID: aggregate.RunID(),
		Messages:   []M1Message{{MessageID: aggregate.RunID() + ":user", Role: "user", Content: "查询新增的漏洞"}},
		RightPanel: M1RightPanel{Tabs: [2]string{"事实", "执行记录"}, EmptyMessage: "选择查询结果后查看安全事实。"},
	}
	if aggregate.Status() == facts.QueryFactsFailed {
		view.Status = "failed"
		view.SafeError = aggregate.SafeError().Message
		view.Messages = append(view.Messages, M1Message{MessageID: aggregate.RunID() + ":assistant", Role: "assistant", Content: aggregate.SafeError().Message})
		return view
	}
	result := aggregate.Result()
	view.Status = string(result.Status())
	view.Result = &M1ResultCard{QuerySequence: aggregate.QuerySequence(), Summary: result.Summary(), Count: result.Count(), ObservedAt: result.ObservedAt()}
	view.Messages = append(view.Messages, M1Message{MessageID: aggregate.RunID() + ":assistant", Role: "assistant", Content: result.Summary()})
	return view
}

func newM1IdleView(workspaceID string) M1WorkbenchView {
	return M1WorkbenchView{
		SchemaVersion: M1WorkbenchViewSchemaVersion, WorkspaceID: workspaceID, ConversationID: "conversation_initial", RunID: "run_initial", Status: "idle",
		Messages: []M1Message{}, Result: nil,
		RightPanel: M1RightPanel{Tabs: [2]string{"事实", "执行记录"}, EmptyMessage: "完成一次查询后，这里将显示安全事实。"},
	}
}
