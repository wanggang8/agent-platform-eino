# Observability 与 Budgets

本文定义运行观测、审计和预算门禁。目标是能定位问题、控制成本和防止失控执行，同时不把运维 telemetry 误当成产品事实。

## 两类输出

| 类型 | 用途 | 可给用户看 | 事实来源 |
| --- | --- | --- | --- |
| Product audit | Workbench、Replay、Inspector、合规追踪 | 是，脱敏 | Product Facts / AuditEvent |
| Operational telemetry | tracing、metrics、latency、token、错误定位 | 否，默认内部 | Eino callback / middleware |

Operational telemetry 不得驱动主 SSE，不得作为 Workbench 或 ActionResult 的事实来源。

## Callback 使用边界

Eino callback 用于：

- trace span。
- model latency。
- tool latency。
- token usage。
- error category。
- internal diagnostics。

Callback 不用于：

- 直接拼接 Workbench view。
- 直接发送产品 SSE。
- 保存 raw provider payload 到 Product Facts。
- 绕过 Safety Gate 输出 assistant 或 tool result。

## Trace 字段

最小 telemetry 字段：

| 字段 | 要求 |
| --- | --- |
| `trace_id` | 内部追踪 id。 |
| `run_id` | 产品 run id。 |
| `workspace_id` | 工作区。 |
| `provider` / `model` | 安全 label，不含 token/base url secret。 |
| `gen_ai.operation.name` | chat/tool/replay/resume 等。 |
| `latency_ms` | 总耗时或阶段耗时。 |
| `input_tokens` / `output_tokens` | 可用时记录。 |
| `tool_count` | 本轮工具调用数。 |
| `failure_category` | 安全错误类别。 |

默认不记录完整 prompt、completion、tool args 或 provider body。若未来开启采样，必须有显式环境开关、脱敏器和验收记录。

## Budgets

P0/P1 最小预算：

| Budget | 默认行为 |
| --- | --- |
| max model calls per run | 超限后停止 tool loop，返回 safe partial result。 |
| max tool calls per run | 超限后不再调用工具，写 audit。 |
| max run duration | 超时进入 stopped 或 failed，按场景安全投影。 |
| max input tokens | 截断低优先级 context，并记录 context snapshot。 |
| max output tokens | 由 provider 配置控制，超限返回安全摘要。 |

P2 可增加 cost estimate、workspace quota、rate limit 和 per-connector budgets。

## Product Audit

必须写入 audit 的事件：

- run created / completed / failed / cancelled / stopped。
- model call started / completed / failed 的安全摘要。
- tool selected / started / completed / failed。
- policy decision。
- credential binding 状态变化。
- approval / clarification requested and resolved。
- budget exceeded。
- safety redaction。

AuditEvent 不得包含 raw prompt、Authorization、API key、credential ref、checkpoint id、interrupt id、raw provider body 或 reusable resume token。

## 报告

real model、live read、live write 和 skip report 必须包含：

- command / scenario。
- selected tool。
- sanitized args summary。
- latency 或 safe duration summary。
- failure category。
- screenshot/report path。
- blocks_claims（skip 时）。

## 验收

- callback 事件不会直接进入 Workbench SSE。
- telemetry 不含 secret、raw prompt、raw provider body。
- budget exceeded 可安全停止并 replay。
- token/latency/tool count 至少进入内部 telemetry 或安全报告。
- audit 与 replay 同源 Product Facts。

参考资料：

- Eino Callback Manual：`https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/`
- OpenTelemetry GenAI semantic conventions：`https://opentelemetry.io/docs/specs/semconv/registry/attributes/gen-ai/`
