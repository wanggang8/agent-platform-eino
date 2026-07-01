package execution

import (
	"context"
	"fmt"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
)

// RunnerEventKind 是 execution 内部事件类型，不允许直接作为 Workbench 或 Action API 契约。
type RunnerEventKind string

const (
	RunnerEventAssistantMessage RunnerEventKind = "assistant_message"
	RunnerEventToolCall         RunnerEventKind = "tool_call"
	RunnerEventToolResult       RunnerEventKind = "tool_result"
	RunnerEventPending          RunnerEventKind = "pending"
	RunnerEventLifecycle        RunnerEventKind = "lifecycle"
)

// RunnerEvent 是 Eino/Runner 事件进入 Product Facts 前的项目内中间形态。
// 它只服务于 mapper，不作为外部接口或前端事件暴露。
type RunnerEvent struct {
	RunID    string
	Kind     RunnerEventKind
	Sequence int64

	Content string

	ToolCallID  string
	ToolID      string
	DisplayName string
	ArgsPreview string
	ResultID    string
	ResultRef   string
	SafeSummary string

	PendingID     string
	PendingKind   string
	PendingStatus string
	ResumeRef     string
	CheckpointRef string
	RiskSummary   string

	Status string
}

// EventMapperConfig 提供可替换时间源，保证事件映射测试不依赖真实时钟。
type EventMapperConfig struct {
	Now func() time.Time
}

// EventMapper 负责把执行事件写入 Product Facts，是 execution event 到产品事实的唯一入口。
type EventMapper struct {
	repository facts.Repository
	now        func() time.Time
}

// NewEventMapper 创建事件映射器；未提供时钟时使用 UTC 当前时间。
func NewEventMapper(repository facts.Repository, config EventMapperConfig) EventMapper {
	now := config.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return EventMapper{repository: repository, now: now}
}

// Map 将单个执行事件转换为 Product Facts 写入。
// 这里不能返回 Workbench/Action API payload，避免 execution 事件绕过投影层。
func (mapper EventMapper) Map(ctx context.Context, event RunnerEvent) error {
	switch event.Kind {
	case RunnerEventAssistantMessage:
		return mapper.repository.AppendTurn(ctx, facts.Turn{
			TurnID:    turnID(event),
			RunID:     event.RunID,
			Role:      facts.TurnRoleAssistant,
			Content:   event.Content,
			Sequence:  event.Sequence,
			CreatedAt: mapper.now(),
		})
	case RunnerEventToolCall:
		return mapper.repository.AppendToolCall(ctx, facts.ToolCall{
			ToolCallID:  event.ToolCallID,
			RunID:       event.RunID,
			ToolID:      event.ToolID,
			DisplayName: event.DisplayName,
			Status:      facts.ToolCallStatus(event.Status),
			ArgsPreview: event.ArgsPreview,
			CreatedAt:   mapper.now(),
		})
	case RunnerEventToolResult:
		return mapper.repository.AppendToolResult(ctx, facts.ToolResult{
			ResultID:   event.ResultID,
			ToolCallID: event.ToolCallID,
			Status:     facts.ToolResultStatus(event.Status),
			StructuredResult: facts.StructuredResultRef{
				SchemaVersion: "tool.structured_result.v1",
				ResultRef:     event.ResultRef,
				SafeSummary:   event.SafeSummary,
			},
		})
	case RunnerEventPending:
		return mapper.repository.AppendPendingInteraction(ctx, facts.PendingInteraction{
			PendingID:     event.PendingID,
			RunID:         event.RunID,
			Kind:          facts.PendingKind(event.PendingKind),
			Status:        facts.PendingStatus(event.PendingStatus),
			ResumeRef:     event.ResumeRef,
			CheckpointRef: event.CheckpointRef,
			RiskSummary:   event.RiskSummary,
		})
	case RunnerEventLifecycle:
		return mapper.repository.UpdateRunStatus(ctx, event.RunID, facts.RunStatus(event.Status), event.SafeSummary, mapper.now())
	default:
		return fmt.Errorf("unsupported runner event kind %q", event.Kind)
	}
}

// turnID 生成稳定的 assistant turn id，保证同一 run 内按 sequence 可重放。
func turnID(event RunnerEvent) string {
	if event.Sequence > 0 {
		return fmt.Sprintf("%s:assistant:%d", event.RunID, event.Sequence)
	}
	return fmt.Sprintf("%s:assistant", event.RunID)
}
