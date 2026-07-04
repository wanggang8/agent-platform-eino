# Phase 7 Provider Usage Mapping Acceptance Record

## Scope

本记录覆盖 Phase 7.2 的 provider usage 映射切片：OpenAI-compatible usage 解析、项目 LLM 安全响应承载、Eino callback output 映射和 runner telemetry 验证。不声明 cost estimate、workspace quota、rate limit 或 OTel exporter 已完成。

## Change Summary

- `llm.ChatResponse` 新增 `Usage TokenUsage`，只包含 input/output/total token 计数。
- OpenAI-compatible provider 解析 `usage.prompt_tokens`、`usage.completion_tokens`、`usage.total_tokens`。
- OpenAI-compatible provider 在 LLM 边界将不可信负数 usage 归一化为 0，避免 provider 外部输入穿透到 Eino callback。
- mock provider 支持测试用 usage。
- `einoModelAdapter.Generate` 将 `llm.TokenUsage` 映射为 Eino `model.CallbackOutput.TokenUsage`。
- ChatModelRunner telemetry 测试覆盖真实 runner 路径中的非零 input/output token。

## Verification

```bash
gofmt -l internal/einoapp/llm/provider.go internal/einoapp/llm/mock.go internal/einoapp/llm/openai_compatible.go internal/einoapp/llm/openai_compatible_test.go internal/einoapp/execution/runner.go internal/einoapp/execution/runner_test.go internal/einoapp/execution/runner_internal_test.go
go test ./internal/einoapp/llm ./internal/einoapp/execution -run 'Usage|Telemetry|EinoTokenUsage' -count=1
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget|RunLifecycle|ImportBoundary|Usage' -count=1
go test ./...
npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario budget
git diff --check
```

结果：命令通过；`gofmt -l` 无输出。

## Remaining Risks

- usage 目前只作为内部 telemetry 计数，不做 cost estimate 或 workspace quota。
- streaming usage 尚未接入。
- 不同 OpenAI-compatible provider 的扩展 usage 字段暂不解析。
