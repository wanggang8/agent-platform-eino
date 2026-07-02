# Phase 8.2 开发前复核记录

阶段：Phase 8.2 Fobrain presentation evidence matrix  
日期：2026-07-02  
执行人：Codex

## 外部资料复核

| 主题 | 官方资料 | 访问日期 | 结论 | 影响 |
| --- | --- | --- | --- | --- |
| Eino ChatModelAgent / Runner | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/ | 2026-07-02 | 官方仍描述 ChatModelAgent 通过 ReAct loop 由模型决定 ToolCall，再注入 Tool result 继续迭代。 | Phase 8.2 只做 Workbench 视觉 fixture，不改变 Eino runner；后续真实工具链仍必须由 capability metadata 和 tool call 驱动。 |
| Eino HITL | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/ | 2026-07-02 | HITL 仍围绕 interrupt/resume 人工介入流程。 | 本切片不新增 approval/resume；不得声明 HITL 已完成。 |
| Eino Checkpoint / Interrupt | https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/ | 2026-07-02 | Checkpoint/interrupt 仍是 waiting/resume 可靠恢复相关能力。 | 本切片只验证已完成工具结果展示，不触碰 checkpoint 设计。 |
| Eino Callback | https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/ | 2026-07-02 | Callback 仍适合 tracing/metrics 观测，不应作为前端事实出口。 | 视觉 fixture 继续消费 Product Facts/StructuredResult 投影，不直接暴露 callback event。 |
| MCP lifecycle | https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle | 2026-07-02 | MCP 2025-06-18 lifecycle 仍包含 initialize、operation、shutdown 阶段。 | 本切片不新增 MCP provider 行为；connector 状态仍在 Task 8.3。 |
| MCP tools | https://modelcontextprotocol.io/specification/2025-06-18/server/tools | 2026-07-02 | Tools spec 仍包含 tools/list、tools/call、structuredContent、isError 等字段。 | 视觉证据仍以 StructuredResult 为唯一展示事实，不引入 raw MCP payload。 |
| OpenAPI 3.1 | https://spec.openapis.org/oas/v3.1.2.html | 2026-07-02 | OAS 3.1.2 继续基于 JSON Schema 2020-12 语义。 | 本切片未改 OpenAPI schema；contract-test 仍作为门禁。 |
| JSON Schema 2020-12 | https://json-schema.org/draft/2020-12/json-schema-core | 2026-07-02 | JSON Schema 2020-12 仍以 `$id`、schema resource 和 dialect 作为版本化契约基础。 | 本切片复用现有 schema-test，不新增 schema dialect。 |
| Playwright visual comparisons | https://playwright.dev/docs/test-snapshots | 2026-07-02 | 官方仍以 `expect(locator).toHaveScreenshot()` / snapshot 目录管理视觉基线，并要求同一运行环境保证截图稳定。 | 继续沿用当前 Playwright visual fixture 和 committed screenshot baseline。 |
| React app from scratch | https://react.dev/learn/start-a-new-react-project | 2026-07-02 | React 官方文档显示当前文档版本为 React 19.2；本项目已固定 React 19，未改变架构。 | 不新增前端框架或状态库。 |
| Vite guide | https://vite.dev/guide/ | 2026-07-02 | Vite SPA/dev server 入口仍适用当前 Vite + React 项目。 | 不调整 Vite 架构。 |

## 版本复核

| 模块/工具 | 当前项目版本 | 可见版本结论 | 决策 |
| --- | --- | --- | --- |
| `github.com/cloudwego/eino` | `v0.9.12` | `go list -m -versions` 可见最新稳定到 `v0.9.12`，另有 `v0.10.0-alpha.*`。 | Phase 8.2 不升级，不采用 alpha。 |
| `github.com/cloudwego/eino-ext/components/model/openai` | 未引入 | 可见到 `v0.1.13`。 | 本任务不涉及模型 provider 变更。 |
| `github.com/eino-contrib/jsonschema` | 未引入 | 可见到 `v1.0.3`。 | 本任务不新增 schema 转换依赖。 |
| Go | `go1.23.12 darwin/arm64` | 本机当前版本。 | 可继续执行现有 Go 门禁。 |
| Node/npm | `node v25.8.1` / `npm 11.11.0` | 本机当前版本。 | 可继续执行前端测试。 |

## 旧项目只读复核

- 旧项目 Fobrain tool acceptance matrix：`/Users/vick/Desktop/project/ai-agent/test-results/fobrain-tool-acceptance/fobrain-tool-acceptance-matrix.json`。
- 矩阵摘要：24 个只读业务工具 + 1 个 connector，`25 passed`，`0 failed`。
- Batch A 业务读参考：
  - `tool.fobrain.current_user_context`：六区域 `audit`、`evidence`、`fresh-main-chat`、`internal-details`、`main-chat`、`process` 均通过。
  - `tool.fobrain.my_permissions`：六区域 `audit`、`evidence`、`fresh-main-chat`、`internal-details`、`main-chat`、`process` 均通过。
- 旧项目截图和 DOM 只作为覆盖参考；本任务必须由新项目 Playwright 生成截图证据。

## 新设计确认

- Workbench 和 Action API 仍共用 Product Facts；本任务只增加前端 fixture 投影，不新增事实模型。
- 工具结果仍以 Fobrain `StructuredResult` / `fobrain.tool_result.v2` fixture 作为安全展示材料。
- Workbench 可见文案必须从 StructuredResult 安全字段生成，不允许在前端 fixture 中补写第二套产品事实。
- 前端继续消费 `web/eino-workbench/src/contracts/generated.ts` 类型，不手写平行 DTO。
- 不新增 Fobrain presenter 分支，不复制旧项目 DOM、CSS class 或接口字段。
- 不进入 Batch B-E，不声明 24 个只读工具完整恢复。
- `connector.fobrain.security` 视觉证据按 `docs/07-implementation-plan.md` Task 8.3 单独处理。

## 允许进入 Phase 8.2

允许进入 Phase 8.2 Batch A business-read visual evidence slice。

前置门禁：

```bash
npm --workspace @agent-platform-eino/eino-workbench run test -- --run
npm --workspace @agent-platform-eino/eino-workbench run visual-test -- --grep fobrain
npm --workspace @agent-platform-eino/eino-workbench run typecheck
npm run eino-workbench:schema-test
```

不得声明：

- 24 个 Fobrain 只读工具完整恢复。
- connector 视觉证据已完成。
- 写域审批、resume 或 live write 可用。
