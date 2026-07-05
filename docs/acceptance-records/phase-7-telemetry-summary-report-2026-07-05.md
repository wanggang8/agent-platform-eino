# Phase 7 Telemetry Summary Report Acceptance Record

## Scope

本记录覆盖 Phase 7.2 的内部 telemetry summary JSON report 切片：从安全 `UsageSummary` 构建 `eino.telemetry_usage_summary_report.v1` 文件报告。不声明 HTTP exporter、Workbench 展示、Product Facts 写入、workspace quota、rate limit 或限制型 cost budget 已完成。

## Change Summary

- 新增 `observability.UsageSummaryReport` 和 `UsageSummaryReportSchemaVersion`。
- 新增 `NewUsageSummaryReport`，从 `UsageSummary` 构建稳定 JSON report。
- 新增 `WriteUsageSummaryReport`，创建父目录并以 0600 权限写入 JSON 文件；已有宽权限文件会被替换/重置为 0600。
- report 输出前会再次归一化 summary label 和数值，不写 raw telemetry events。
- report 固定 `blocks_claims=[]`，schema 禁止非空阻断声明。
- 新增 `docs/schemas/telemetry_usage_summary_report.v1.schema.json`、fixture 和 manifest 入口。

## Verification

```bash
gofmt -l internal/einoapp/observability/report.go internal/einoapp/observability/report_test.go internal/einoapp/observability/summary.go internal/einoapp/observability/summary_test.go
go test ./internal/einoapp/observability -run 'Report|Summary|Telemetry|Cost' -count=1
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget|RunLifecycle|ImportBoundary|Usage|Cost|Summary|Report' -count=1
npm run eino-workbench:contract-test
go test ./...
bash scripts/eino_workbench_server_smoke.sh --scenario budget
git diff --check
```

结果：命令通过；`gofmt -l` 无输出。首次 `npm run eino-workbench:contract-test` 检测到新增 schema 后 `web/eino-workbench/src/contracts/generated.ts` 过期，已运行 `npm run eino-workbench:contract-generate` 更新契约索引；最终复跑显示 contract index is current for 37 schemas。

## Remaining Risks

- 当前 report writer 是内部库能力，尚未提供 CLI 或 OTel exporter。
- report 只覆盖单 run summary，不做 workspace 汇总或 quota。
