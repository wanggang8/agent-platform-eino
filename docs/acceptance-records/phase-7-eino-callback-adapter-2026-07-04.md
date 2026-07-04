# Phase 7 Eino Callback Adapter Acceptance Record

## Scope

本记录覆盖 Phase 7.2 的 ChatModel Eino callback adapter。它只将 Eino callback 转为内部 telemetry，不声明 OpenTelemetry exporter、workspace quota 或真实 provider cost estimate 已完成。

## Change Summary

- 新增 `observability.NewEinoCallbackHandler`，接入 Eino `OnStart` / `OnEnd` / `OnError`。
- handler 只接受 ChatModel callback，忽略非模型 callback。
- handler 能从 `model.CallbackOutput.TokenUsage` 提取 input/output token 计数，并记录 latency、run/workspace、provider/model label 和安全 failure category；真实 runner/provider 路径的非零 usage 映射已由后续 provider usage 切片补齐。
- ChatModelRunner 使用 `callbacks.InitCallbacks` 注入 handler；`einoModelAdapter.Generate` 自行触发 Eino callback，并声明 `IsCallbacksEnabled` 避免框架重复包装。
- callback telemetry 不写 Product Facts、Workbench SSE、ActionResult 或 Replay。

## Verification

```bash
go test ./internal/einoapp/observability ./internal/einoapp/execution -run 'EinoCallback|Telemetry' -count=1
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget|RunLifecycle|ImportBoundary' -count=1
go test ./...
npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario budget
git diff --check
gofmt -l internal/einoapp/observability/eino_callback.go internal/einoapp/observability/eino_callback_test.go internal/einoapp/execution/runner.go internal/einoapp/execution/runner_test.go
```

结果：命令通过。

## Remaining Risks

- OpenTelemetry exporter 尚未实现。
- 真实 runner/provider token usage 映射已由后续 provider usage 切片补齐。
- error callback 当前只落安全 `model_error`，更细的 provider timeout/config 分类待后续接入 provider error taxonomy。
- workspace quota、rate limit 和 cost estimate 仍待后续任务实现。
