# 开发前复核记录：2026-07-01

## 范围

- 目标 Phase：Phase 3，Eino chat、Product Facts、SSE、Action API。
- 新项目路径：`/Users/vick/Desktop/project/agent-platform-eino`。
- 旧项目路径：`/Users/vick/Desktop/project/ai-agent`，只读参考，未修改。
- 结论：允许继续 Phase 3；不允许声明 P0/P1/P2 完成，不允许声明 Phase 6 HITL/checkpoint/resume 或 Phase 8 Fobrain 能力完成。

## 外部资料

| 主题 | 访问日期 | 结论 | 对本阶段影响 |
| --- | --- | --- | --- |
| Eino ChatModelAgent / Runner | 2026-07-01 | 官方文档入口为 `https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/`；本阶段以本地锁定源码 `github.com/cloudwego/eino@v0.9.12` 为实现真相，确认 `adk.NewChatModelAgent`、`adk.NewRunner`、`Runner.Query/Run` 可用于 mock ChatModelAgent 路径。 | 允许接入 ChatModelAgent Runner，但 Runner event 必须先落 Product Facts。 |
| Eino HITL / Checkpoint / Callback | 2026-07-01 | 官方入口为 `agent_hitl`、`checkpoint_interrupt`、`callback_manual`；本阶段只记录 resume/checkpoint 边界，不实现 durable CheckPointStore、approval/clarification interrupt 或 restart resume。 | Phase 6 前必须重新复核并补专项验收。 |
| MCP lifecycle / tools | 2026-07-01 | 参考 `https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle` 与 `/server/tools`；当前 Phase 不实现 MCP adapter。 | Phase 4 前重新复核。 |
| OpenAPI / JSON Schema | 2026-07-01 | 继续使用 OpenAPI 3.1 与 JSON Schema 2020-12；schema/fixture 仍由 `scripts/eino_workbench_schema_validate.mjs` 校验。 | Contract 门禁保持不变。 |

## 版本命令摘要

执行：

```bash
GOTOOLCHAIN=local go list -m github.com/cloudwego/eino
GOTOOLCHAIN=local go list -m -versions github.com/cloudwego/eino
GOTOOLCHAIN=local go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
GOTOOLCHAIN=local go list -m -versions github.com/eino-contrib/jsonschema
```

结果摘要：

- 当前项目：`github.com/cloudwego/eino v0.9.12`。
- 可见 Eino 稳定版本到 `v0.9.12`，另有 `v0.10.0-alpha.*`；不采用 alpha。
- `github.com/cloudwego/eino-ext/components/model/openai` 可见到 `v0.1.13`，Phase 4 真实 provider 再固定。
- `github.com/eino-contrib/jsonschema` 可见到 `v1.0.3`，当前因 Eino schema 传递依赖进入 `go.sum`，不作为本项目 schema 生成入口。

## 当前架构复核

已只读参考旧项目：

- `docs/ARCHITECTURE.md`：旧项目 Workbench 与 Integration API 共用 Runtime 事实，工具结果统一 `StructuredResult`。
- `docs/STATUS.md`：旧项目已覆盖 Runtime P0、Workbench、安全投影、Fobrain 24 只读、connector、实体消歧、写域审批。
- `docs/design/workbench-chat-contract.md`：旧项目强调 Workbench 只从 shared contract 和 `structured_result` 渲染，不返回 CSS/DOM/provider raw 字段。

对新项目 Phase 3 的约束：

- 不复制旧 runtime、旧 step runner、旧 tool pair、旧 Workbench DOM 或旧接口兼容层。
- Eino Runner 是执行输入，不是产品契约。
- Workbench、Action API、Replay、SSE 必须从 Product Facts / product projection 读取。
- 工具结果仍只能以 StructuredResult 安全投影进入产品层。

## 新设计确认

- Workbench 和 Action API 共用 Product Facts：是，`cmd/eino-workbench` 注入同一个 SQLite repository 给 execution commands 和 product projection。
- 意图识别不按关键词路由：是，Phase 3 仅保留 registry/policy selection 边界；Fobrain 真实选择留到 Phase 4/5/8。
- 模型上下文只来自 safe Product Facts / StructuredResult：是，`ContextProjector` 从 snapshot 构造 messages 并保存 context snapshot。
- Eino event 只作为 Product Facts 输入：是，HTTP/SSE 不 import Eino event；Runner 输出经 mapper 写 assistant turn/lifecycle。
- HITL/checkpoint/resume：Phase 3 仅保留边界；Phase 6 前不得声明生产恢复能力。
- Fobrain 24 只读与写域审批：设计和 fixture 已覆盖范围；Phase 8 前不得声明能力等价。

## 风险与处理

| 风险 | 处理 |
| --- | --- |
| Phase 3 Runner 目前只覆盖 mock ChatModelAgent assistant 输出 | 允许作为 Phase 3 P0 执行链路；工具 loop、real provider、StructuredResult conversion 进入 Phase 4。 |
| `chat-stream` / `action-basic` smoke 可能假阳性 | 已要求断言 assistant message、completed ActionResult 和 SQLite context snapshot。 |
| SSE event id 当前由 projection 重算 | Phase 3 可作为静态 snapshot cursor；Phase 7 replay/reconnect 前必须持久化 facts cursor 或扩展 sequence 事实。 |
| Missing run 合成空视图可能掩盖事实缺失 | 当前仅用于 current view/未创建 run 空态；具体 run 的严格 404 可在后续 API hardening 中收紧。 |

## 结论

- 是否允许继续 Phase 3：是。
- 是否允许进入 Phase 4：否，需完成 Phase 3 收尾、记录验收并重新确认 tool/provider 方案。
- 是否允许进入 Phase 6：否，需 HITL/checkpoint/resume 专项复核。
- 是否允许进入 Phase 8：否，需 Fobrain 字段级恢复和 live 门禁。
- 是否允许声明 P0/P1/P2 完成：否。
- 关联 ADR：`docs/adr/2026-06-30-pin-eino-version.md`、`docs/adr/2026-07-01-phase-3-facts-sqlite-decisions.md`。
