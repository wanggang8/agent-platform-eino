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

当前实现提供内部 `observability.Sink` 和内存测试 sink。ChatModelRunner 通过 Eino ChatModel callback handler 记录 `chat.model.generate` 的安全 telemetry 事件，包含 run/workspace、provider/model label、latency、token/tool count 字段和 failure category；该事件不进入 Product Facts、Workbench SSE、ActionResult 或 Replay。OpenTelemetry exporter 仍需后续接入。

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

当前 Phase 7.2 最小实现先固定预算终态入口：execution lifecycle 接收 `budget_exceeded`，只允许作用于 `running` 或 `waiting` run；系统将 run 标记为 `failed`、写入 `safe_error=budget_exceeded`、取消 active tool、过期 waiting pending，并追加安全 `event_type=budget` audit event。该入口代表已由预算判断层触发的安全结果，不在 HTTP、前端或 telemetry 中直接写 Product Facts。完整 token/cost 估算、workspace quota 和 rate limit 仍需后续任务实现。

当前配置文件支持 `max_model_calls_per_run`、`max_tool_calls_per_run` 和 `max_input_tokens_per_run`。execution 预算评估器只返回安全 `BudgetDecision`；调用方必须再通过 `budget_exceeded` lifecycle 写入 Product Facts，不能由 telemetry 或 HTTP 直接修改产品状态。

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

real model provider、real model tool、live read、live write 和 skip report 必须包含：

- command / scenario。
- provider-only 报告包含 provider kind、model label、脱敏配置摘要和 redaction checks。
- tool/live 报告包含 selected tool、sanitized args summary 和 tool card rendering status。
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
