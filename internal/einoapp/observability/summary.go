package observability

import (
	"strings"
	"time"
)

// UsageSummary 是内部 telemetry 的 run 级安全统计摘要，不是 Product Facts。
type UsageSummary struct {
	RunID                   string
	WorkspaceID             string
	EventCount              int
	ModelCallCount          int
	FailureCount            int
	InputTokens             int
	OutputTokens            int
	TotalTokens             int
	EstimatedCostMicrounits int64
	TotalLatencyMS          int64
	FirstEventAt            time.Time
	LastEventAt             time.Time
}

// RunSummary 返回指定 run 的内部统计摘要；未知 run 返回只包含安全 run_id 的空摘要。
func (sink *MemorySink) RunSummary(runID string) UsageSummary {
	sink.mu.Lock()
	defer sink.mu.Unlock()
	return SummarizeRunEvents(sink.events, runID)
}

// SummarizeRunEvents 聚合指定 run 的安全 telemetry 事件，不读取或写入 Product Facts。
func SummarizeRunEvents(events []Event, runID string) UsageSummary {
	safeRunID := safeLabel(runID)
	summary := UsageSummary{RunID: safeRunID}
	// 不安全 run_id 会被折叠成 redacted；直接返回空摘要，避免多个不安全 id 串账。
	if strings.TrimSpace(runID) != safeRunID {
		return summary
	}
	for _, event := range events {
		event = sanitizeEvent(event)
		if event.RunID != safeRunID {
			continue
		}
		summary.EventCount++
		if event.OperationName == "chat.model.generate" {
			summary.ModelCallCount++
		}
		if event.FailureCategory != "" {
			summary.FailureCount++
		}
		if summary.WorkspaceID == "" {
			summary.WorkspaceID = event.WorkspaceID
		}
		summary.InputTokens += event.InputTokens
		summary.OutputTokens += event.OutputTokens
		summary.TotalTokens += event.TotalTokens
		summary.EstimatedCostMicrounits += event.EstimatedCostMicrounits
		summary.TotalLatencyMS += event.LatencyMS
		if !event.CreatedAt.IsZero() && (summary.FirstEventAt.IsZero() || event.CreatedAt.Before(summary.FirstEventAt)) {
			summary.FirstEventAt = event.CreatedAt
		}
		if event.CreatedAt.After(summary.LastEventAt) {
			summary.LastEventAt = event.CreatedAt
		}
	}
	return summary
}
