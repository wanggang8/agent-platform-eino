# Phase 7 Telemetry Summary Report Smoke

阶段：Phase 7.2 Audit callbacks, observability, and budgets  
日期：2026-07-06  
方案：新增内部 telemetry summary report 可复跑 smoke 生成器  
执行人：Codex

## 环境

- Go：go1.23.12 darwin/arm64
- Node：v25.8.1
- npm：11.11.0
- OS：macOS 26.5.2 25F84
- 模型 provider：不使用
- Fobrain 环境：不使用

## 命令

```bash
go test ./scripts/telemetry_summary_report ./internal/einoapp/observability -run 'TelemetrySummary|Report|Summary|Cost' -count=1
go run ./scripts/telemetry_summary_report --output test-results/eino-workbench-telemetry-usage-summary-report.json
node scripts/eino_workbench_report_validate.mjs --schema docs/schemas/telemetry_usage_summary_report.v1.schema.json --report test-results/eino-workbench-telemetry-usage-summary-report.json
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget|RunLifecycle|ImportBoundary|Usage|Cost|Summary|Report' -count=1
go test ./...
npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario budget
git diff --check
```

## 报告

- telemetry summary：`test-results/eino-workbench-telemetry-usage-summary-report.json`
- schema：`docs/schemas/telemetry_usage_summary_report.v1.schema.json`

## 结论

- 通过：本记录中的命令已在本地执行通过。
- 该报告是内部诊断材料，不是 Product Facts、Workbench、Action API、Replay 或 audit 的事实来源。
