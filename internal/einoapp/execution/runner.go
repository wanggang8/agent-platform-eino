package execution

import (
	"context"
	"errors"
	"fmt"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/llm"
	"github.com/cloudwego/eino/adk"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// ChatModelRunnerConfig 提供可替换时间源，保证 Runner 写入事实的测试稳定。
type ChatModelRunnerConfig struct {
	Now func() time.Time
}

// ChatModelRunner 使用 Eino ChatModelAgent 执行一次对话，并把事件落入 Product Facts。
type ChatModelRunner struct {
	repository facts.Repository
	provider   llm.Provider
	config     llm.Config
	now        func() time.Time
}

// NewChatModelRunner 创建 Phase 3 的 Eino ChatModel runner。
func NewChatModelRunner(repository facts.Repository, provider llm.Provider, config llm.Config, runnerConfig ChatModelRunnerConfig) ChatModelRunner {
	now := runnerConfig.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return ChatModelRunner{
		repository: repository,
		provider:   provider,
		config:     config,
		now:        now,
	}
}

// Run 执行指定 run 的 ChatModelAgent，并只通过 EventMapper 写入产品事实。
func (runner ChatModelRunner) Run(ctx context.Context, runID string) error {
	if err := runner.repository.UpdateRunStatus(ctx, runID, facts.RunStatusRunning, "", runner.now()); err != nil {
		return err
	}
	projector := NewContextProjector(runner.repository, ContextProjectorConfig{Now: runner.now})
	safeContext, err := projector.Project(ctx, runID)
	if err != nil {
		return err
	}

	model, err := runner.provider.NewChatModel(ctx, runner.config)
	if err != nil {
		return err
	}
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "eino-workbench-chat",
		Description: "Agent Workbench chat runner",
		Model:       einoModelAdapter{model: model},
		GenModelInput: func(context.Context, string, *adk.AgentInput) ([]*schema.Message, error) {
			return toEinoMessages(safeContext), nil
		},
		MaxIterations: 1,
	})
	if err != nil {
		return err
	}

	einoRunner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent, EnableStreaming: false})
	iterator := einoRunner.Query(ctx, lastUserContent(safeContext))
	mapper := NewEventMapper(runner.repository, EventMapperConfig{Now: runner.now})
	sequence, err := nextSequence(ctx, runner.repository, runID)
	if err != nil {
		return err
	}
	assistantWritten := false
	for {
		event, ok := iterator.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			_ = runner.repository.UpdateRunStatus(ctx, runID, facts.RunStatusFailed, "模型执行失败", runner.now())
			return event.Err
		}
		for _, runnerEvent := range EinoAgentEventToRunnerEvents(runID, sequence, event) {
			if runnerEvent.Kind == RunnerEventAssistantMessage && assistantWritten {
				continue
			}
			if err := mapper.Map(ctx, runnerEvent); err != nil {
				return err
			}
			if runnerEvent.Kind == RunnerEventAssistantMessage {
				assistantWritten = true
			}
		}
	}
	if !assistantWritten {
		_ = runner.repository.UpdateRunStatus(ctx, runID, facts.RunStatusFailed, "模型未返回可展示内容", runner.now())
		return errors.New("chat model returned no assistant message")
	}
	if err := runner.repository.AppendAuditEvent(ctx, facts.AuditEvent{
		AuditID:     fmt.Sprintf("%s:audit:message:%d", runID, runner.now().UnixNano()),
		RunID:       runID,
		EventType:   "message",
		SafeSummary: "assistant message completed",
		Actor:       "model",
		CreatedAt:   runner.now(),
	}); err != nil {
		return err
	}
	return runner.repository.UpdateRunStatus(ctx, runID, facts.RunStatusSucceeded, "", runner.now())
}

// einoModelAdapter 把项目 llm.ChatModel 适配为 Eino BaseChatModel。
type einoModelAdapter struct {
	model llm.ChatModel
}

// Generate 将 Eino messages 转成项目安全消息，再返回 Eino assistant message。
func (adapter einoModelAdapter) Generate(ctx context.Context, input []*schema.Message, _ ...einomodel.Option) (*schema.Message, error) {
	response, err := adapter.model.Generate(ctx, llm.ChatRequest{Messages: fromEinoMessages(input)})
	if err != nil {
		return nil, err
	}
	return schema.AssistantMessage(response.Content, nil), nil
}

// Stream 当前 Phase 3 不启用 Eino streaming；SSE 只从 Product Facts 投影生成。
func (adapter einoModelAdapter) Stream(context.Context, []*schema.Message, ...einomodel.Option) (*schema.StreamReader[*schema.Message], error) {
	return nil, errors.New("streaming model is not enabled in phase 3")
}

func nextSequence(ctx context.Context, repository facts.Repository, runID string) (int64, error) {
	snapshot, err := repository.GetSnapshot(ctx, runID)
	if err != nil {
		return 0, err
	}
	maxSequence := int64(0)
	for _, turn := range snapshot.Turns {
		if turn.Sequence > maxSequence {
			maxSequence = turn.Sequence
		}
	}
	return maxSequence + 1, nil
}

func toEinoMessages(messages []llm.Message) []*schema.Message {
	result := make([]*schema.Message, 0, len(messages))
	for _, message := range messages {
		switch message.Role {
		case "system":
			result = append(result, schema.SystemMessage(message.Content))
		case "assistant":
			result = append(result, schema.AssistantMessage(message.Content, nil))
		default:
			result = append(result, schema.UserMessage(message.Content))
		}
	}
	return result
}

func fromEinoMessages(messages []*schema.Message) []llm.Message {
	result := make([]llm.Message, 0, len(messages))
	for _, message := range messages {
		result = append(result, llm.Message{
			Role:    string(message.Role),
			Content: message.Content,
		})
	}
	return result
}

func lastUserContent(messages []llm.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			return messages[i].Content
		}
	}
	return ""
}
