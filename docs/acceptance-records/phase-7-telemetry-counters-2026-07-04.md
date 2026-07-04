# Phase 7 Telemetry Counters Acceptance Record

## Scope

本记录覆盖 Phase 7.2 的内部 telemetry/counter 基础切片，不声明 OpenTelemetry exporter、真实 provider token usage、cost estimate 或 workspace quota 已完成。

## Change Summary

- 新增 `internal/einoapp/observability`，提供安全 `Event`、`Sink`、`NoopSink` 和测试用 `MemorySink`。
- ChatModelRunner 记录 `chat.model.generate` 内部 telemetry 事件，包含 run/workspace、operation、provider/model label、latency 和 failure category，不写 Product Facts。
- 配置文件新增 `max_model_calls_per_run`、`max_tool_calls_per_run`、`max_input_tokens_per_run`，旧配置会派生安全默认值。
- execution 新增预算评估器，只返回安全 `BudgetDecision`，不直接修改 run 状态。
- telemetry label 采用保守脱敏：secret/token/key/cookie/raw payload/URL/超长或非安全字符 label 均写为 `redacted`。
- architecture import boundary 禁止 `httpapi` / `product` 反向导入 `observability`，防止 telemetry 进入产品出口层。

## Verification

```bash
go test ./internal/einoapp/observability ./internal/einoapp/execution ./internal/einoapp/bootstrap -run 'Telemetry|Budget|Config' -count=1
go test ./internal/einoapp/observability ./internal/einoapp/execution ./internal/einoapp/bootstrap -run 'Telemetry|Budget|Config|FobrainRedactedSummary' -count=1
go test ./internal/einoapp/observability ./internal/einoapp/architecture ./internal/einoapp/execution ./internal/einoapp/bootstrap -run 'Telemetry|ImportBoundary|Budget|Config|FobrainRedactedSummary' -count=1
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget|RunLifecycle' -count=1
go test ./...
npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario budget
git diff --check
```

结果：命令通过。

## Remaining Risks

- Eino ChatModel callback handler 已接入安全 telemetry sink；真实 provider usage 映射仍需接入具体 provider 返回。
- OTel exporter、cost estimate、workspace quota 和 rate limit 仍待后续任务实现。
