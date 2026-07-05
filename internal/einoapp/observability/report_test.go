package observability_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/observability"
)

func TestNewUsageSummaryReportSanitizesSummary(t *testing.T) {
	// report 构建阶段再次归一化输入，避免调用方绕过 MemorySink 后写出不安全字段。
	report := observability.NewUsageSummaryReport(observability.UsageSummary{
		RunID:                   "run-token-secret",
		WorkspaceID:             "https://workspace.example.test",
		EventCount:              -1,
		ModelCallCount:          -2,
		FailureCount:            -3,
		InputTokens:             -4,
		OutputTokens:            -5,
		TotalTokens:             -6,
		EstimatedCostMicrounits: -7,
		TotalLatencyMS:          -8,
	}, time.Unix(100, 0))

	if report.SchemaVersion != observability.UsageSummaryReportSchemaVersion ||
		report.ReportKind != "telemetry_usage_summary" ||
		!report.GeneratedAt.Equal(time.Unix(100, 0).UTC()) ||
		report.Summary.RunID != "redacted" ||
		report.Summary.WorkspaceID != "redacted" ||
		report.Summary.EventCount != 0 ||
		report.Summary.InputTokens != 0 ||
		report.Summary.EstimatedCostMicrounits != 0 ||
		len(report.BlocksClaims) != 0 {
		t.Fatalf("report = %+v", report)
	}
}

func TestWriteUsageSummaryReportWritesSafeJSONFile(t *testing.T) {
	// 文件报告只包含 summary，不写 raw events、prompt、completion 或 provider payload。
	path := filepath.Join(t.TempDir(), "nested", "telemetry-summary.json")
	summary := observability.UsageSummary{
		RunID:                   "run-safe",
		WorkspaceID:             "ws-safe",
		EventCount:              2,
		ModelCallCount:          2,
		FailureCount:            1,
		InputTokens:             10,
		OutputTokens:            16,
		TotalTokens:             26,
		EstimatedCostMicrounits: 62,
		TotalLatencyMS:          24,
		FirstEventAt:            time.Unix(10, 0).UTC(),
		LastEventAt:             time.Unix(20, 0).UTC(),
	}

	if err := observability.WriteUsageSummaryReport(path, summary, time.Unix(30, 0).UTC()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"run_id"`) || strings.Contains(string(data), `"RunID"`) {
		t.Fatalf("report must use stable snake_case JSON fields: %s", data)
	}
	for _, forbidden := range []string{"raw prompt", "raw completion", "provider_payload", "authorization", "secret", "events"} {
		if strings.Contains(strings.ToLower(string(data)), forbidden) {
			t.Fatalf("report leaked %q: %s", forbidden, data)
		}
	}
	var decoded observability.UsageSummaryReport
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.SchemaVersion != observability.UsageSummaryReportSchemaVersion ||
		decoded.Summary.RunID != "run-safe" ||
		decoded.Summary.TotalTokens != 26 ||
		decoded.Summary.EstimatedCostMicrounits != 62 {
		t.Fatalf("decoded report = %+v", decoded)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("report mode = %o", info.Mode().Perm())
	}
}

func TestWriteUsageSummaryReportResetsExistingFilePermission(t *testing.T) {
	// 已存在的宽权限文件必须被替换为 0600，不能依赖 os.WriteFile 的创建权限。
	path := filepath.Join(t.TempDir(), "telemetry-summary.json")
	if err := os.WriteFile(path, []byte(`{"old":true}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := observability.WriteUsageSummaryReport(path, observability.UsageSummary{RunID: "run-safe"}, time.Unix(30, 0).UTC()); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("report mode = %o", info.Mode().Perm())
	}
}

func TestUsageSummaryReportBlocksClaimsAlwaysEmpty(t *testing.T) {
	// telemetry summary report 不能携带阻断声明，避免内部统计报告被误用为限制信号。
	report := observability.NewUsageSummaryReport(observability.UsageSummary{RunID: "run-safe"}, time.Unix(30, 0).UTC())
	if report.BlocksClaims == nil || len(report.BlocksClaims) != 0 {
		t.Fatalf("blocks_claims = %+v", report.BlocksClaims)
	}
}

func TestWriteUsageSummaryReportRejectsEmptyPath(t *testing.T) {
	// report 输出路径必须显式提供，避免意外写入当前目录或覆盖不相关文件。
	if err := observability.WriteUsageSummaryReport(" ", observability.UsageSummary{}, time.Unix(30, 0).UTC()); err == nil {
		t.Fatal("WriteUsageSummaryReport empty path error = nil")
	}
}
