# Phase 8.3 开发前复核记录

阶段：Phase 8.3 Connector 状态与凭据绑定  
日期：2026-07-02  
执行人：Codex

## 外部资料复核

本任务与 Phase 8.2 在同一日执行，继续采用 `docs/acceptance-records/pre-development-validation-2026-07-02-phase-8-2.md` 已复核的官方资料版本；本记录补充 connector/credential 相关影响判断。

| 主题 | 官方资料 | 访问日期 | 结论 | 影响 |
| --- | --- | --- | --- | --- |
| Eino ChatModelAgent / Runner | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/ | 2026-07-02 | ToolCall 仍由 Eino runner 驱动。 | connector 状态只作为 capability 结果投影，不新增后端关键词路由。 |
| Eino HITL | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/ | 2026-07-02 | HITL 仍围绕 interrupt/resume。 | 本任务不声明 approval/resume 完成。 |
| MCP lifecycle | https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle | 2026-07-02 | connector lifecycle 仍需明确初始化和关闭边界。 | 本任务只展示 Fobrain connector 状态，不新增 MCP lifecycle 实现。 |
| MCP tools | https://modelcontextprotocol.io/specification/2025-06-18/server/tools | 2026-07-02 | tools/call 可携带 structuredContent/isError。 | Workbench 仍只展示 StructuredResult 安全投影，不展示 raw tool payload。 |
| OpenAPI 3.1 | https://spec.openapis.org/oas/v3.1.2.html | 2026-07-02 | OAS 3.1.2 继续基于 JSON Schema 2020-12。 | 本任务不改 OpenAPI；schema/contract 门禁继续执行。 |
| JSON Schema 2020-12 | https://json-schema.org/draft/2020-12/json-schema-core | 2026-07-02 | 版本化 schema 仍以 `$id` 和 dialect 为基础。 | 新增 connector fixture 必须通过 `fobrain.tool_result.v2` schema。 |
| Playwright visual comparisons | https://playwright.dev/docs/test-snapshots | 2026-07-02 | snapshot baseline 仍是视觉回归门禁。 | 为 connector 生成 desktop/mobile 六区域截图。 |
| React / Vite | https://react.dev/learn/start-a-new-react-project / https://vite.dev/guide/ | 2026-07-02 | 当前 Vite + React 架构仍适用。 | 不引入新前端框架或状态库。 |

## 版本复核

沿用同日 Phase 8.2 版本复核：`github.com/cloudwego/eino v0.9.12`，未采用 `v0.10.0-alpha.*`；本任务不新增 Go/TS 依赖。

## 新设计确认

- connector 状态和 credential binding 只进入 `fobrain.tool_result.v2` 的安全 facts。
- Workbench 可见文案只来自 StructuredResult 安全字段，不引入第二套事实。
- 不展示 token、auth header、`api_token`、`credential_ref`、raw connector config 或 raw provider payload。
- 不声明 24/24、Batch B-E、实体消歧、写域审批或 live write 完成。

## 允许进入 Phase 8.3

允许进入 Phase 8.3 connector visual evidence slice。
