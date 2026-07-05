# Phase 7 Telemetry Run Summary Acceptance Record

## Scope

本记录覆盖 Phase 7.2 的内部 run summary 切片：从安全 telemetry events 聚合 run 级 token/cost/latency/failure 统计。不声明 HTTP exporter、Workbench 展示、Product Facts 写入、workspace quota、rate limit 或限制型 cost budget 已完成。

## Change Summary

- 新增 `observability.UsageSummary`，只包含安全 run/workspace label 和数值统计。
- 新增 `MemorySink.RunSummary(run_id)`，按 run 聚合已脱敏 telemetry events。
- 新增 `observability.SummarizeRunEvents`，对传入事件再次执行安全归一化，避免绕过 MemorySink 时泄漏不安全 label 或负数；不安全 run_id 查询只返回空摘要，避免多个 redacted run 串账。
- summary 覆盖 event/model/failure count、input/output/total token、estimated cost、total latency 和首末事件时间。
- summary 不读取或写入 Product Facts，不进入 Workbench SSE、Action API、Replay 或 audit。

## Verification

```bash
gofmt -l internal/einoapp/observability/summary.go internal/einoapp/observability/summary_test.go
go test ./internal/einoapp/observability -run 'Summary|Telemetry|Cost' -count=1
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget|RunLifecycle|ImportBoundary|Usage|Cost|Summary' -count=1
go test ./...
npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario budget
git diff --check
```

结果：命令通过；`gofmt -l` 无输出。

## Remaining Risks

- summary 目前只在内部内存 sink 中可用，尚未实现文件报告或 OTel exporter。
- 跨进程/持久化 telemetry 聚合仍待后续 exporter/reporting 切片实现。
