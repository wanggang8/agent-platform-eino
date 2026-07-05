package observability_test

import (
	"context"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/observability"
)

func TestMemorySinkRunSummaryAggregatesSafeStatistics(t *testing.T) {
	// run summary 只聚合内部安全计数，不回读 Product Facts 或 provider payload。
	sink := observability.NewMemorySink()
	mustRecordTelemetry(t, sink, observability.Event{
		RunID:                   "run-a",
		WorkspaceID:             "ws-a",
		OperationName:           "chat.model.generate",
		InputTokens:             3,
		OutputTokens:            5,
		TotalTokens:             8,
		EstimatedCostMicrounits: 19,
		LatencyMS:               11,
		CreatedAt:               time.Unix(20, 0).UTC(),
	})
	mustRecordTelemetry(t, sink, observability.Event{
		RunID:                   "run-a",
		WorkspaceID:             "ws-a",
		OperationName:           "chat.model.generate",
		InputTokens:             7,
		OutputTokens:            11,
		TotalTokens:             18,
		EstimatedCostMicrounits: 43,
		LatencyMS:               13,
		FailureCategory:         "model_error",
		CreatedAt:               time.Unix(10, 0).UTC(),
	})
	mustRecordTelemetry(t, sink, observability.Event{
		RunID:                   "run-b",
		WorkspaceID:             "ws-b",
		OperationName:           "chat.model.generate",
		InputTokens:             100,
		OutputTokens:            200,
		TotalTokens:             300,
		EstimatedCostMicrounits: 400,
		LatencyMS:               500,
		CreatedAt:               time.Unix(30, 0).UTC(),
	})

	summary := sink.RunSummary("run-a")
	if summary.RunID != "run-a" ||
		summary.WorkspaceID != "ws-a" ||
		summary.EventCount != 2 ||
		summary.ModelCallCount != 2 ||
		summary.FailureCount != 1 ||
		summary.InputTokens != 10 ||
		summary.OutputTokens != 16 ||
		summary.TotalTokens != 26 ||
		summary.EstimatedCostMicrounits != 62 ||
		summary.TotalLatencyMS != 24 ||
		!summary.FirstEventAt.Equal(time.Unix(10, 0).UTC()) ||
		!summary.LastEventAt.Equal(time.Unix(20, 0).UTC()) {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestRunSummaryReturnsEmptySummaryForUnknownRun(t *testing.T) {
	// 未知 run 只能返回空统计，不能猜测 workspace 或复制其它 run 的计数。
	sink := observability.NewMemorySink()
	mustRecordTelemetry(t, sink, observability.Event{RunID: "run-a", InputTokens: 10})

	summary := sink.RunSummary("run-missing")
	if summary.RunID != "run-missing" ||
		summary.EventCount != 0 ||
		summary.InputTokens != 0 ||
		summary.EstimatedCostMicrounits != 0 ||
		summary.WorkspaceID != "" {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestSummarizeRunEventsSanitizesInputEvents(t *testing.T) {
	// 聚合函数也按不可信输入处理，unsafe run_id 只返回空摘要，避免 redacted 串账。
	summary := observability.SummarizeRunEvents([]observability.Event{{
		RunID:                   "run-token-secret",
		WorkspaceID:             "https://workspace.example.test",
		OperationName:           "chat.model.generate",
		InputTokens:             -1,
		OutputTokens:            -2,
		TotalTokens:             -3,
		EstimatedCostMicrounits: -4,
		LatencyMS:               -5,
	}}, "run-token-secret")

	if summary.RunID != "redacted" ||
		summary.WorkspaceID != "" ||
		summary.EventCount != 0 ||
		summary.InputTokens != 0 ||
		summary.OutputTokens != 0 ||
		summary.TotalTokens != 0 ||
		summary.EstimatedCostMicrounits != 0 ||
		summary.TotalLatencyMS != 0 {
		t.Fatalf("summary = %+v", summary)
	}
}

func TestSummarizeRunEventsDoesNotMergeUnsafeRunLabels(t *testing.T) {
	// 不同 unsafe run_id 都会脱敏为 redacted，因此不能按 redacted 聚合。
	events := []observability.Event{
		{RunID: "run-token-a", InputTokens: 10},
		{RunID: "run-secret-b", InputTokens: 20},
	}

	summary := observability.SummarizeRunEvents(events, "run-token-a")
	if summary.RunID != "redacted" || summary.EventCount != 0 || summary.InputTokens != 0 {
		t.Fatalf("summary = %+v", summary)
	}
}

func mustRecordTelemetry(t *testing.T, sink observability.Sink, event observability.Event) {
	t.Helper()
	if err := sink.Record(context.Background(), event); err != nil {
		t.Fatal(err)
	}
}
