# 开发前复核门禁

本文定义开始实现 Eino-first 重构前必须完成的复核。目标是确认外部框架资料、当前代码基线和新设计三者一致，避免在错误假设上开工。

该门禁不是旧实现迁移说明；它只确认新设计是否覆盖既有产品能力、安全边界和外部框架约束。

## 必须执行的时间点

以下节点开始编码前必须执行本复核：

- Phase 1：创建契约和服务骨架前。
- Phase 3：引入或升级 Eino 依赖前。
- Phase 4：实现 Capability Provider / MCP / Tool Adapter 前。
- Phase 6：实现 HITL / checkpoint / resume 前。
- Phase 8：恢复 Fobrain 产品能力前。

如果距离上一次复核超过 7 天，或外部依赖版本发生变化，也必须重新执行。

## 外部资料复核

必须重新打开官方资料，并在验收记录中写入访问日期、结论和影响。

| 主题 | 官方资料 | 需要确认的问题 |
| --- | --- | --- |
| Eino ChatModelAgent / Runner | `https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/` | Runner event、tool call、stream 行为是否仍符合 Product Facts 设计 |
| Eino HITL | `https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/` | interrupt/resume 数据结构是否可映射到 `PendingInteraction` |
| Eino Checkpoint / Interrupt | `https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/` | CheckPointStore 是否支持进程重启恢复和 waiting 状态保留 |
| Eino Callback | `https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/` | callback 是否只用于 tracing/metrics，不作为产品 SSE 主来源 |
| MCP lifecycle | `https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle` | initialize/session/auth/close 行为是否覆盖 provider contract |
| MCP tools | `https://modelcontextprotocol.io/specification/2025-06-18/server/tools` | tools/list、listChanged、annotations、structuredContent/isError 是否覆盖 |
| OpenAPI 3.1 | `https://spec.openapis.org/oas/v3.1.2.html` | OpenAPI 3.1 与 JSON Schema dialect、`type: ["x", "null"]`、oneOf 等生成器支持；不使用 OAS 3.0 `nullable` |
| JSON Schema 2020-12 | `https://json-schema.org/draft/2020-12/json-schema-core` | `$id`、`$schema`、版本化 schema、validator 行为 |
| Playwright visual comparisons | `https://playwright.dev/docs/test-snapshots` | snapshot 更新、threshold、mask、运行环境稳定性 |
| React / Vite | `https://react.dev/learn/start-a-new-react-project`、`https://vite.dev/guide/` | 前端脚手架、构建、测试和开发服务器配置是否仍适用 |

## 版本复核命令

```bash
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go list -m -versions github.com/eino-contrib/jsonschema
```

记录要求：

- 当前项目使用版本。
- 官方可见最新稳定版本。
- 是否存在 alpha/beta 线。
- 是否升级。
- 不升级的原因。
- 升级后必须重跑的验证命令。

2026-06-30 调研快照：

| 模块 | 当前项目版本 | 当日可见版本结论 |
| --- | --- | --- |
| `github.com/cloudwego/eino` | `v0.9.12` | `go list -m -versions` 可见最新稳定到 `v0.9.12`，并有 `v0.10.0-alpha.*`；不采用 alpha |
| `github.com/cloudwego/eino-ext/components/model/openai` | 未引入 | 当前可见到 `v0.1.13`；引入模型 provider 时再固定 |
| `github.com/eino-contrib/jsonschema` | 未引入 | 当前可见到 `v1.0.3`；需要独立 schema 转换时再固定 |

该快照已固化到 `go.mod` 和 `adr/2026-06-30-pin-eino-version.md`。Phase 3 真正接入 Eino runtime 前必须重新执行版本复核，并补齐 ChatModelAgent、Runner event、tool loop、interrupt、checkpoint、callback 验收。

## 当前项目架构复核

开始开发前必须只读参考旧项目并记录结论。旧项目路径固定为 `/Users/vick/Desktop/project/ai-agent`：

- `/Users/vick/Desktop/project/ai-agent/docs/ARCHITECTURE.md`
- `/Users/vick/Desktop/project/ai-agent/docs/STATUS.md`
- `/Users/vick/Desktop/project/ai-agent/docs/ROADMAP.md`
- `/Users/vick/Desktop/project/ai-agent/docs/PRODUCTION_CONNECTOR_GUIDE.md`
- `/Users/vick/Desktop/project/ai-agent/docs/design/workbench-chat-contract.md`
- `/Users/vick/Desktop/project/ai-agent/docs/design/workbench-chat-frontend-contract.md`
- `/Users/vick/Desktop/project/ai-agent/docs/design/workbench-ui-acceptance-matrix.md`
- `/Users/vick/Desktop/project/ai-agent/test-results/workbench-acceptance-audit/workbench-acceptance-audit.json`

必须检查的代码边界：

| 边界 | 路径 | 复核目标 |
| --- | --- | --- |
| Agent 主链 | `internal/orchestration/agent` | 识别旧模型工具循环职责，但不复制旧 loop 结构 |
| Runtime 事实 | `internal/runtime` | 识别现有 Run/Step/Pending/Tool Pair 能力，确认新设计由 Eino + Product Facts 替代 |
| 产品契约 | `internal/interfaces/presentation` | 确认 ActionResult、Workbench SSE/ViewModel、安全引用没有丢失 |
| Capability Registry | `internal/registry/capability`、`internal/registry/providerapi` | 确认 provider metadata、StructuredResult、风险策略和凭据绑定需要保留为新契约 |
| Fobrain 样板 | `internal/business/fobrain` | 确认 24 只读、connector、消歧、写域审批、live smoke 能力被新矩阵覆盖 |
| Eino adapter 现状 | `internal/infrastructure/llm`、`internal/infrastructure/eino` | 确认旧适配经验，不复制旧接口包袱 |

## 新设计确认问题

复核记录必须逐项回答：

- 新设计是否仍保持 Workbench 和 Action API 共用同一套 Product Facts？
- 意图识别是否仍以 Eino native tool call 和 capability metadata 为主，而不是后端关键词路由？
- 模型上下文是否只来自 safe Product Facts / StructuredResult，并有 context snapshot？
- Run lifecycle 是否覆盖 cancel、stop、timeout、retry 和终态幂等？
- callback、telemetry 和 budget 是否不绕过 Product Facts 或 Safety Gate？
- 工具结果是否仍只有 StructuredResult 一份事实材料？
- JSON 输出是否仍禁止可复用 `resume_token`，只暴露安全 `approval_refs` / `resume_refs`？
- Eino event 是否只作为 Product Facts 输入，而不是直接暴露给前端或外部 API？
- HITL / checkpoint / resume 是否能覆盖 approval、clarification、duplicate、restart、missing checkpoint？
- Capability Provider 是否禁止 raw provider payload 进入产品输出？
- MCP provider 是否覆盖 lifecycle、tool list cache、listChanged、annotations、structuredContent/isError？
- Fobrain 24 个只读工具、connector、实体消歧、写域审批是否全部映射到新 schema / fixture / report？
- 前端视觉是否有当前产品参考截图、fixture、block matrix 和 baseline 更新规则？
- 旧项目最终通过证据是否已按 `legacy-acceptance-evidence.md` 对齐，且没有被当成新项目通过结果？
- 新设计是否没有依赖旧 runtime 类型、旧 step runner、旧 tool pair、旧 Workbench 接口兼容层？

任一问题回答为“不确定”或“否”，不得开始对应 Phase 的实现。

## 输出产物

每次复核必须生成：

```text
docs/acceptance-records/pre-development-validation-YYYY-MM-DD.md
```

记录模板见 `acceptance-records/PRE_DEVELOPMENT_VALIDATION_TEMPLATE.md`。

记录必须包含：

- 外部资料访问日期和链接。
- 版本命令输出摘要。
- 当前代码架构复核结论。
- 新设计风险清单。
- 需要更新的 ADR。
- 是否允许进入对应 Phase。

## 阻断条件

出现以下任一情况必须阻断开发：

- 外部官方资料与当前设计冲突。
- Eino 版本 API 与设计中的 ChatModelAgent、Runner、HITL、checkpoint、callback 不匹配。
- MCP spec 与 provider contract 不匹配。
- OpenAPI / JSON Schema 工具链不支持本项目约定。
- 现有产品能力没有映射到新设计。
- 复核发现新设计需要复制旧 runtime 类型或旧接口才能跑通。
- Fobrain P2 能力没有字段级矩阵、fixture 和验收断言。
