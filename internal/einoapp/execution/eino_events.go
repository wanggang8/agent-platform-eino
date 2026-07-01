package execution

import (
	"agent-platform-eino/internal/einoapp/facts"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// EinoAgentEventToRunnerEvents 是 Eino AgentEvent 进入 Product Facts 前的唯一转换点。
// 它只提取安全产品事件；raw Eino event 不允许进入 HTTP、product、audit 或 replay。
func EinoAgentEventToRunnerEvents(runID string, sequence int64, event *adk.AgentEvent) []RunnerEvent {
	if event == nil {
		return nil
	}
	if event.Err != nil {
		return []RunnerEvent{{
			RunID:       runID,
			Kind:        RunnerEventLifecycle,
			Sequence:    sequence,
			Status:      string(facts.RunStatusFailed),
			SafeSummary: "模型执行失败",
		}}
	}
	if event.Output != nil && event.Output.MessageOutput != nil && event.Output.MessageOutput.Message != nil {
		message := event.Output.MessageOutput.Message
		if message.Role == schema.Assistant && message.Content != "" {
			return []RunnerEvent{{
				RunID:    runID,
				Kind:     RunnerEventAssistantMessage,
				Sequence: sequence,
				Content:  message.Content,
			}}
		}
	}
	return nil
}
