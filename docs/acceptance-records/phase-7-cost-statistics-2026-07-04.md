# Phase 7 Cost Statistics Acceptance Record

## Scope

本记录覆盖 Phase 7.2 的只统计型 cost estimate 切片：配置化 token 单价、callback telemetry 总 token 和估算成本字段、负数归一化，以及 runner 路径验证。不声明 budget 限制、workspace quota、rate limit 或 OTel exporter 已完成。

## Change Summary

- `observability.Event` 新增 `TotalTokens` 和 `EstimatedCostMicrounits`。
- `observability.TokenCostRates` 使用显式配置的 micro-units/token 单价估算成本；默认 0 不内置模型价格。
- Eino ChatModel callback telemetry 记录 `total_tokens` 和 `estimated_cost_microunits`，只进入内部 sink。
- `observability.cost_statistics.input_unit_microunits` / `output_unit_microunits` 从 YAML 配置读取，负数配置启动失败。
- 成本统计不写 Product Facts、Workbench SSE、Action API、Replay 或 lifecycle。

## Verification

```bash
gofmt -l cmd/eino-workbench/main.go internal/einoapp/bootstrap/config.go internal/einoapp/bootstrap/config_test.go internal/einoapp/execution/runner.go internal/einoapp/execution/runner_test.go internal/einoapp/observability/cost.go internal/einoapp/observability/eino_callback.go internal/einoapp/observability/eino_callback_test.go internal/einoapp/observability/telemetry.go internal/einoapp/observability/telemetry_test.go
go test ./internal/einoapp/observability ./internal/einoapp/execution ./internal/einoapp/bootstrap -run 'Telemetry|Cost|Config' -count=1
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget|RunLifecycle|ImportBoundary|Usage|Cost' -count=1
go test ./...
npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario budget
git diff --check
```

结果：命令通过；`gofmt -l` 无输出。

## Remaining Risks

- 当前成本只在内部 telemetry 中估算，不做持久化报表或导出。
- workspace quota、rate limit、限制型 cost budget 和 OTel exporter 仍待后续任务实现。
- streaming usage 尚未接入成本统计。
