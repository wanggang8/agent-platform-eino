package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"agent-platform-eino/internal/einoapp/observability"
)

const telemetrySummaryRunID = "run-telemetry-summary-smoke"

func main() {
	outputPath := flag.String("output", "test-results/eino-workbench-telemetry-usage-summary-report.json", "telemetry summary report output path")
	flag.Parse()

	if err := runTelemetrySummaryReport(*outputPath, time.Now().UTC()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runTelemetrySummaryReport 生成内部 telemetry summary 文件报告；该入口不连接 HTTP、Product Facts 或产品投影。
func runTelemetrySummaryReport(outputPath string, generatedAt time.Time) error {
	if outputPath == "" {
		return errors.New("telemetry summary report output path is required")
	}
	sink := observability.NewMemorySink()
	if err := recordTelemetrySummarySmokeEvents(context.Background(), sink); err != nil {
		return err
	}
	summary := sink.RunSummary(telemetrySummaryRunID)
	return observability.WriteUsageSummaryReport(outputPath, summary, generatedAt)
}

// recordTelemetrySummarySmokeEvents 写入固定安全样本，验证 report 只聚合目标 run 的脱敏统计。
func recordTelemetrySummarySmokeEvents(ctx context.Context, sink observability.Sink) error {
	events := []observability.Event{
		{
			TraceID:                 "trace-telemetry-summary-smoke",
			RunID:                   telemetrySummaryRunID,
			WorkspaceID:             "ws-telemetry-summary",
			OperationName:           "chat.model.generate",
			Provider:                "mock",
			Model:                   "mock-chat",
			LatencyMS:               120,
			InputTokens:             40,
			OutputTokens:            20,
			TotalTokens:             60,
			EstimatedCostMicrounits: 300,
			ToolCount:               1,
			CreatedAt:               time.Date(2026, 7, 6, 3, 0, 0, 0, time.UTC),
		},
		{
			TraceID:                 "trace-telemetry-summary-smoke",
			RunID:                   telemetrySummaryRunID,
			WorkspaceID:             "ws-telemetry-summary",
			OperationName:           "chat.model.generate",
			Provider:                "mock",
			Model:                   "mock-chat",
			LatencyMS:               80,
			InputTokens:             10,
			OutputTokens:            5,
			TotalTokens:             15,
			EstimatedCostMicrounits: 75,
			FailureCategory:         "model_error",
			CreatedAt:               time.Date(2026, 7, 6, 3, 0, 5, 0, time.UTC),
		},
		{
			TraceID:                 "trace-other-run",
			RunID:                   "run-not-in-report",
			WorkspaceID:             "ws-other",
			OperationName:           "chat.model.generate",
			Provider:                "mock",
			Model:                   "mock-chat",
			LatencyMS:               999,
			InputTokens:             999,
			OutputTokens:            999,
			TotalTokens:             1_998,
			EstimatedCostMicrounits: 9_999,
			CreatedAt:               time.Date(2026, 7, 6, 3, 0, 10, 0, time.UTC),
		},
	}
	for _, event := range events {
		if err := sink.Record(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
