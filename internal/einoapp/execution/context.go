package execution

import (
	"context"
	"fmt"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/llm"
)

// ContextProjectorConfig 提供可替换时间源，保证 context snapshot 测试稳定。
type ContextProjectorConfig struct {
	Now func() time.Time
}

// ContextProjector 从 Product Facts 构造模型安全上下文。
// 它不读取 provider raw payload、checkpoint id 或 credential ref。
type ContextProjector struct {
	repository facts.Repository
	now        func() time.Time
}

// NewContextProjector 创建安全上下文投影器。
func NewContextProjector(repository facts.Repository, config ContextProjectorConfig) ContextProjector {
	now := config.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return ContextProjector{repository: repository, now: now}
}

// Project 从同一 run 的 Product Facts 生成 LLM messages，并记录 context snapshot。
func (projector ContextProjector) Project(ctx context.Context, runID string) ([]llm.Message, error) {
	snapshot, err := projector.repository.GetSnapshot(ctx, runID)
	if err != nil {
		return nil, err
	}

	messages := make([]llm.Message, 0, len(snapshot.Turns)+len(snapshot.ToolResults))
	for _, turn := range snapshot.Turns {
		messages = append(messages, llm.Message{
			Role:    string(turn.Role),
			Content: turn.Content,
		})
	}
	for _, result := range snapshot.ToolResults {
		if strings.TrimSpace(result.StructuredResult.SafeSummary) == "" {
			continue
		}
		messages = append(messages, llm.Message{
			Role:    "system",
			Content: "StructuredResult: " + result.StructuredResult.SafeSummary,
		})
	}

	now := projector.now()
	if err := projector.repository.SaveContextSnapshot(ctx, facts.ContextSnapshot{
		SnapshotID:  fmt.Sprintf("%s:context:%d", runID, now.UnixNano()),
		RunID:       runID,
		SafeSummary: contextSummary(messages),
		CreatedAt:   now,
	}); err != nil {
		return nil, err
	}
	return messages, nil
}

// contextSummary 保存摘要而不是 raw prompt，供 audit/replay 证明上下文来源。
func contextSummary(messages []llm.Message) string {
	if len(messages) == 0 {
		return "empty safe context"
	}
	return fmt.Sprintf("safe context messages=%d", len(messages))
}
