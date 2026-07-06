package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunTelemetrySummaryReportWritesSafeReport(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "telemetry-summary.json")
	generatedAt := time.Date(2026, 7, 6, 3, 5, 0, 0, time.UTC)

	if err := runTelemetrySummaryReport(outputPath, generatedAt); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	encoded := strings.ToLower(string(data))
	// 泄漏断言覆盖文档禁止进入 telemetry summary report 的 raw 材料、凭据和可复用 token。
	for _, forbidden := range []string{
		"raw events",
		"raw prompt",
		"completion",
		"tool args",
		"provider_payload",
		"provider raw payload",
		"authorization",
		"api key",
		"api_key",
		"credential ref",
		"credential_ref",
		"resume token",
		"resume_ref:",
		"run-not-in-report",
		"ws-other",
	} {
		if strings.Contains(encoded, strings.ToLower(forbidden)) {
			t.Fatalf("report leaked forbidden marker %q: %s", forbidden, encoded)
		}
	}

	var report struct {
		SchemaVersion string `json:"schema_version"`
		ReportKind    string `json:"report_kind"`
		GeneratedAt   string `json:"generated_at"`
		Summary       struct {
			RunID                   string `json:"run_id"`
			WorkspaceID             string `json:"workspace_id"`
			EventCount              int    `json:"event_count"`
			ModelCallCount          int    `json:"model_call_count"`
			FailureCount            int    `json:"failure_count"`
			InputTokens             int    `json:"input_tokens"`
			OutputTokens            int    `json:"output_tokens"`
			TotalTokens             int    `json:"total_tokens"`
			EstimatedCostMicrounits int64  `json:"estimated_cost_microunits"`
			TotalLatencyMS          int64  `json:"total_latency_ms"`
		} `json:"summary"`
		BlocksClaims []string `json:"blocks_claims"`
	}
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != "eino.telemetry_usage_summary_report.v1" ||
		report.ReportKind != "telemetry_usage_summary" ||
		report.GeneratedAt != generatedAt.Format(time.RFC3339) {
		t.Fatalf("report header mismatch: %+v", report)
	}
	if report.Summary.RunID != telemetrySummaryRunID ||
		report.Summary.WorkspaceID != "ws-telemetry-summary" ||
		report.Summary.EventCount != 2 ||
		report.Summary.ModelCallCount != 2 ||
		report.Summary.FailureCount != 1 ||
		report.Summary.InputTokens != 50 ||
		report.Summary.OutputTokens != 25 ||
		report.Summary.TotalTokens != 75 ||
		report.Summary.EstimatedCostMicrounits != 375 ||
		report.Summary.TotalLatencyMS != 200 {
		t.Fatalf("summary mismatch: %+v", report.Summary)
	}
	if len(report.BlocksClaims) != 0 {
		t.Fatalf("blocks_claims = %+v, want empty", report.BlocksClaims)
	}
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("report mode = %v, want 0600", info.Mode().Perm())
	}
}

func TestRunTelemetrySummaryReportRejectsEmptyOutputPath(t *testing.T) {
	// 空输出路径应在脚本边界提前失败，避免误写当前工作目录。
	if err := runTelemetrySummaryReport("", time.Date(2026, 7, 6, 3, 5, 0, 0, time.UTC)); err == nil {
		t.Fatal("expected error")
	}
}
