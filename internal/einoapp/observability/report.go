package observability

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const UsageSummaryReportSchemaVersion = "eino.telemetry_usage_summary_report.v1"

// UsageSummaryReport 是内部 telemetry summary 的文件报告，不是 Workbench 或 Action API 契约。
type UsageSummaryReport struct {
	SchemaVersion string       `json:"schema_version"`
	ReportKind    string       `json:"report_kind"`
	GeneratedAt   time.Time    `json:"generated_at"`
	Summary       UsageSummary `json:"summary"`
	BlocksClaims  []string     `json:"blocks_claims"`
}

// NewUsageSummaryReport 构建可写入文件的安全 summary 报告。
func NewUsageSummaryReport(summary UsageSummary, generatedAt time.Time) UsageSummaryReport {
	return UsageSummaryReport{
		SchemaVersion: UsageSummaryReportSchemaVersion,
		ReportKind:    "telemetry_usage_summary",
		GeneratedAt:   generatedAt.UTC(),
		Summary:       sanitizeUsageSummary(summary),
		BlocksClaims:  []string{},
	}
}

// WriteUsageSummaryReport 将内部 telemetry summary 写入 JSON 文件，文件权限固定为 0600。
func WriteUsageSummaryReport(path string, summary UsageSummary, generatedAt time.Time) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("usage summary report path is required")
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(NewUsageSummaryReport(summary, generatedAt), "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	removeTmp := true
	defer func() {
		if removeTmp {
			_ = os.Remove(tmpPath)
		}
	}()
	if _, err := tmp.Write(append(data, '\n')); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return err
	}
	removeTmp = false
	return os.Chmod(path, 0o600)
}

func sanitizeUsageSummary(summary UsageSummary) UsageSummary {
	summary.RunID = safeLabel(summary.RunID)
	summary.WorkspaceID = safeLabel(summary.WorkspaceID)
	summary.EventCount = nonNegativeInt(summary.EventCount)
	summary.ModelCallCount = nonNegativeInt(summary.ModelCallCount)
	summary.FailureCount = nonNegativeInt(summary.FailureCount)
	summary.InputTokens = nonNegativeInt(summary.InputTokens)
	summary.OutputTokens = nonNegativeInt(summary.OutputTokens)
	summary.TotalTokens = nonNegativeInt(summary.TotalTokens)
	summary.EstimatedCostMicrounits = nonNegativeInt64(summary.EstimatedCostMicrounits)
	summary.TotalLatencyMS = nonNegativeInt64(summary.TotalLatencyMS)
	return summary
}
