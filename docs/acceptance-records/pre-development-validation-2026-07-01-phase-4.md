# Phase 4 Pre-development Validation

日期：2026-07-01
目标 Phase：Phase 4，StructuredResult、Safety Gate、Eino tool loop、真实 LLM provider 和 provider policy。

## 外部资料复核

| 主题 | 访问日期 | 结论 | 对本项目影响 |
| --- | --- | --- | --- |
| Eino ChatModelAgent | 2026-07-01 | 官方说明配置 Tools 后 ChatModelAgent 进入 ReAct loop：模型选择工具、执行工具、把结果注入上下文继续迭代。 | `tool-card` 应先用 mock tool loop 验证事实链路，再接真实模型。 |
| Eino Tool interface | 2026-07-01 | 官方 Tool 指南定义 `BaseTool.Info(ctx)` 和 `InvokableTool.InvokableRun(ctx, argumentsInJSON)`，ToolInfo 持有参数 schema。 | Capability Registry 必须转成 Eino ToolInfo；InvokableRun 只能返回安全 tool message。 |
| Eino OpenAI AgenticModel | 2026-07-01 | OpenAI agentic model 是 provider 集成能力，服务于 Agent 能力。 | OpenAI-compatible provider 应作为模型边界替换，不应改变 Product Facts/tool adapter。 |
| Eino HITL | 2026-07-01 | HITL 依赖 interrupt/resume/checkpoint 语义。 | 写域工具继续阻断到 Phase 6，Phase 4 不执行 mutation。 |
| Eino Checkpoint / Interrupt | 2026-07-01 | 官方 checkpoint/interrupt 支持暂停并从断点恢复。 | Phase 4 不实现 durable resume；只保留写域阻断。 |
| Eino Callback | 2026-07-01 | callback 用于日志、tracing、metrics、屏幕展示等横切能力注入。 | callback 不作为 Workbench SSE 主事实来源，仍只进 observability。 |
| MCP Lifecycle | 2026-07-01 | MCP lifecycle 覆盖 initialize、operation、shutdown 和 capability negotiation。 | Phase 4 MCP 只做 contract/mock，不接生产 server。 |
| MCP Tools | 2026-07-01 | MCP tools 以 name 唯一标识，包含 metadata/schema，可返回 structuredContent/isError。 | MCP 输出仍只是 StructuredResult candidate，必须过 Safety Gate。 |

参考链接：

- `https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/`
- `https://www.cloudwego.io/docs/eino/core_modules/components/tools_node_guide/`
- `https://www.cloudwego.io/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/`
- `https://www.cloudwego.io/docs/eino/ecosystem_integration/chat_model/agentic_model_openai/`
- `https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/`
- `https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/`
- `https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/`
- `https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle`
- `https://modelcontextprotocol.io/specification/2025-06-18/server/tools`

## 版本复核

执行命令：

```bash
go list -m github.com/cloudwego/eino
go list -m -f '{{.Dir}}' github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go list -m -versions github.com/eino-contrib/jsonschema
```

| 模块 | 当前项目版本 | 当前可见版本结论 | 是否升级 |
| --- | --- | --- | --- |
| `github.com/cloudwego/eino` | `v0.9.12` | `go list -m -versions` 可见最新稳定为 `v0.9.12`，另有 `v0.10.0-alpha.1` 到 `v0.10.0-alpha.9`。 | 不升级；不采用 alpha。 |
| `github.com/cloudwego/eino-ext/components/model/openai` | 未引入 | 当前可见到 `v0.1.13`。 | Phase 4.3 引入真实 provider 时再固定。 |
| `github.com/eino-contrib/jsonschema` | 未引入 | 当前可见到 `v1.0.3`。 | 只有在 ToolInfo JSON Schema 转换需要独立依赖时再引入。 |

本地 module source 复核：

- `/Users/vick/go/pkg/mod/github.com/cloudwego/eino@v0.9.12/components/tool/interface.go`
- `/Users/vick/go/pkg/mod/github.com/cloudwego/eino@v0.9.12/schema/tool.go`
- `/Users/vick/go/pkg/mod/github.com/cloudwego/eino@v0.9.12/adk/chatmodel.go`

结论：本地 API 与 Phase 4 ADR 一致；升级 Eino 后必须重跑 StructuredResult、tool adapter、tool loop、HITL 和 provider smoke。

## 旧项目只读参考

读取范围：

- `/Users/vick/Desktop/project/ai-agent/docs/ARCHITECTURE.md`
- `/Users/vick/Desktop/project/ai-agent/docs/STATUS.md`
- `/Users/vick/Desktop/project/ai-agent/docs/ROADMAP.md`
- `/Users/vick/Desktop/project/ai-agent/docs/PRODUCTION_CONNECTOR_GUIDE.md`
- `/Users/vick/Desktop/project/ai-agent/docs/design/workbench-chat-contract.md`
- `/Users/vick/Desktop/project/ai-agent/docs/design/workbench-chat-frontend-contract.md`
- `/Users/vick/Desktop/project/ai-agent/docs/design/workbench-ui-acceptance-matrix.md`
- `/Users/vick/Desktop/project/ai-agent/internal/interfaces/presentation`
- `/Users/vick/Desktop/project/ai-agent/internal/registry/capability`
- `/Users/vick/Desktop/project/ai-agent/internal/registry/providerapi`
- `/Users/vick/Desktop/project/ai-agent/internal/business/fobrain`
- `/Users/vick/Desktop/project/ai-agent/internal/infrastructure/llm`
- `/Users/vick/Desktop/project/ai-agent/internal/infrastructure/eino`

只读结论：

- 旧项目要求 Workbench 与 Integration API 共用 Runtime 事实；新项目继续用 Product Facts 同源替代，不能复制旧 runtime 类型。
- 旧项目明确工具结果只有一份 `StructuredResult` / `structured_result`；Phase 4.1 必须先固定 Safety Gate。
- 旧 ActionResult 禁止 reusable `resume_token`，只暴露安全 refs；Phase 4 不改变该约束。
- 旧 Workbench 合同要求前端按 `structured_result.schema_version` 和 Fobrain `display_type` 呈现，不按 `tool_id` 分支；新项目 Phase 4 tool-card 仍需遵守。
- 旧 Fobrain 基线包含 24 只读工具 + 1 connector、实体消歧、凭据绑定和写域审批；Phase 4 只做通用 tool loop，不声明 Fobrain 等价。
- 旧 Eino/LLM 适配只能作为经验参考；新实现必须通过当前 Eino-first 分层、Capability Registry、Product Facts 和 Safety Gate 重建。

## 当前新项目架构复核

- `facts.ToolResult` 已以 StructuredResult 作为事实材料，不保存 provider raw payload。
- `product.Projection` 已能从 ToolCall/ToolResult 生成工具卡、SSE tool patch 和 ActionResult result card。
- `execution.EventMapper` 已有 tool call/result 事件到 Product Facts 的映射测试。
- `capabilities.Registry` 已支持 capability metadata 和 risk policy；Phase 4 需要补 Eino ToolInfo adapter。
- `cmd/eino-workbench` 已从配置构建 registry，不内置 smoke capability。
- 当前 import boundary 仍阻止 `httpapi`、`product`、`facts`、`providers/fobrain` 直接依赖 Eino。

## 设计确认问题

| 问题 | 结论 |
| --- | --- |
| Workbench 和 Action API 是否继续共用 Product Facts？ | 是，Phase 4 tool facts 仍写同一 repository，再由 product projection 输出。 |
| 意图识别是否以 Eino native tool call 和 capability metadata 为主？ | 是，tool name/description/schema 来自 Capability Registry，不按关键词路由。 |
| 模型上下文是否只来自 safe Product Facts / StructuredResult？ | 是，`ContextProjector` 已有 context snapshot；Phase 4 增加 tool result context 回归。 |
| Run lifecycle 是否覆盖 cancel、stop、timeout、retry 和终态幂等？ | Phase 4 不新增完整生命周期；Phase 7 前不得声明 lifecycle 完成。 |
| callback、telemetry 和 budget 是否不绕过 Product Facts 或 Safety Gate？ | 是，Phase 4 不把 callback 作为产品 SSE 主来源；budget/telemetry 在后续门禁补齐。 |
| 工具结果是否仍只有 StructuredResult 一份事实材料？ | 是，Phase 4.1 先加 Safety Gate。 |
| JSON 输出是否禁止可复用 `resume_token`？ | 是，只允许安全 `approval_refs` / `resume_refs`；Phase 4 不新增 resume token。 |
| Eino event 是否只作为 Product Facts 输入？ | 是，HTTP/product 层不得 import Eino event。 |
| HITL / checkpoint / resume 是否已覆盖？ | 否，Phase 6 前不得声明；写域 mutation 继续阻断。 |
| Capability Provider 是否禁止 raw provider payload 进入产品输出？ | 是，Safety Gate 是强制门禁。 |
| MCP provider 是否覆盖 lifecycle、tool list cache、listChanged、annotations、structuredContent/isError？ | Phase 4 只做 contract/mock；生产 MCP 后移。 |
| Fobrain 24 只读、connector、实体消歧、写域审批是否映射到新 schema / fixture / report？ | 已在 docs/fobrain-tool-matrix.md 和 Phase 8 计划中映射；Phase 4 不声明等价。 |
| 前端视觉是否有参考截图、fixture、block matrix 和 baseline 更新规则？ | 是，Phase 2 文档和 visual evidence matrix 已覆盖；Phase 4 tool-card 仍需生成新证据。 |
| 旧项目最终通过证据是否只作范围参考？ | 是，不能作为新项目通过结果。 |
| 新设计是否不依赖旧 runtime 类型、旧 step runner、旧 tool pair、旧 Workbench 接口兼容层？ | 是，Phase 4 通过 Eino-first tool adapter + Product Facts 重建。 |

## 风险与处理

| 风险 | 处理 |
| --- | --- |
| 真实 LLM provider 提前接入导致 tool loop 失败难定位 | Phase 4 顺序改为 StructuredResult/Safety -> mock tool loop -> real LLM。 |
| Eino `InvokableRun` 返回 string 被误当成事实材料 | ADR 明确返回值只作为模型 tool message，事实必须来自 Safety Gate 后的 StructuredResult。 |
| 写域能力在无 HITL 时被执行 | Phase 4 只允许 read-only mock tool loop；写域继续阻断到 Phase 6。 |
| MCP 复杂度拖慢 P0 | Phase 4 只保留 MCP adapter contract，生产 MCP 后移。 |
| `real-model-chat` 与 P0/P1 门禁混淆 | `tool-card` 属于 P0/Phase 4 必跑；`real-model-chat` 归属 P1 真实模型 smoke，无凭据可 skip 但必须记录。 |

## 结论

- 是否允许进入 Phase 4.1：是。
- 是否允许打开 `tool-card` smoke：否，需完成 Phase 4.1 和 Phase 4.2。
- 是否允许接真实 LLM provider：否，需先完成 mock tool loop。
- 是否允许把 `real-model-chat` 作为 P0 必过门禁：否，它属于 P1 真实模型 smoke；无凭据可 skip 但不能声明真实模型能力通过。
- 是否允许执行写域 mutation：否，需 Phase 6 HITL/checkpoint/resume。
- 是否允许声明 P0/P1/P2 完成：否。
- 关联 ADR：`docs/adr/2026-07-01-phase-4-tool-loop-before-real-llm.md`。
