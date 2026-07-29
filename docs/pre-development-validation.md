# 开发前复核门禁

本文定义每个纵向 Story 开始实现前必须完成的复核。目标是确认当前 Story 的需求、外部合同、架构边界、数据来源和验收条件可执行；不是恢复旧项目全部能力的检查表。

## 当前开工裁决

- Story 1.1 已完成，`G-TOOLCHAIN=PASS`。
- Implementation Readiness 已于 2026-07-20 达到 `READY`；Story 1.2 已完成 Create Story，开发前复核于 2026-07-27 裁决 `ALLOW`，可以开始其限定范围实现。
- 当前正式范围仅为 `epics.md` 中 27 个 Story；旧项目 24 工具齐套、connector、工单、移动端和旧接口兼容不构成本轮开工条件。

## 必须执行的时间点

- 每个正式 Story 第一次开始开发前。
- M1、M2、M3、M4、M5、M6 每个里程碑的首个 Story 前。
- 外部依赖、provider 合同、schema epoch、架构决策或授权范围变化后。
- 距离同一 Story / 同一外部依赖的上次复核超过 7 天时。
- M5 每一次目标部署写域取证前；授权不可跨日期或跨动作复用。

## 第一步：锁定工作项

复核记录必须写明：

- Milestone、Epic、Story ID、标题和 WorkItemType。
- `epics.md` 中逐项列出的 Requirements、Architecture decisions 和 Prerequisites。
- Inputs / Outputs、Scope、Non-goals、Affected directories、Acceptance commands 和 Exit gate。
- 当前 Sprint 状态与所有适用 Gate / OQ 状态。
- 本次允许的 mutation 上限；只读和 mock Story 必须为 0。

若 Story 标题、语义、前置或状态与 `epics.md`、Sprint 状态、实施计划不一致，立即 `BLOCKED`。

## 第二步：核对产品与架构边界

每次必须逐项回答：

- Workbench 和 Action API 是否只消费同一套 Product Facts？
- 工具结果是否只以 StructuredResult 进入产品层？
- raw provider payload 是否停留在 `providers/*` 边界内？
- LLM 上下文是否只来自安全投影，不包含 raw payload、secret 或未授权字段？
- 工具选择是否基于 capability registry、metadata、intent 和 policy，而非工具名或关键词分支？
- `httpapi`、`execution`、`facts`、`product`、`capabilities`、`llm`、`observability`、`providers/*`、`store/sqlite` 的依赖方向是否保持清晰？
- Eino event 是否只作为 Product Facts 输入，不直接成为前端或外部 API 契约？
- checkpoint、resume、SSE、audit 和 operation history 是否引用同一事实与 run 状态？
- HTTP / provider / LLM / resume 等不可信输入是否在边界校验并脱敏？
- 当前实现是否不需要复制旧 runtime、旧执行链、旧 Workbench API、旧数据库或旧 UI DOM？

任一回答为“不确定”或“否”时，不得开工。

## 第三步：按需复核官方资料

只复核当前 Story 实际使用或本次发生版本变化的资料，并在记录中写入访问日期、版本结论和设计影响。技术问题只使用官方来源。

| 当前 Story 使用的能力 | 官方资料 | 必须确认 |
| --- | --- | --- |
| Eino ChatModelAgent / Runner | `https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/` | event、tool call、stream 与当前最薄执行链一致 |
| Eino HITL | `https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/` | interrupt/resume 能映射到 pending interaction；仅 M4/M6 需要 |
| Eino checkpoint | `https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/` | 进程重启恢复、waiting 与 checkpoint 行为 |
| Eino callback | `https://www.cloudwego.io/docs/eino/core_modules/eino_adk/adk_agent_callback/` | callback 只做 observability，不裁决事实或终态 |
| OpenAPI 3.1 | `https://spec.openapis.org/oas/v3.1.2.html` | 与 JSON Schema 2020-12 和生成器一致，不使用 OAS 3.0 `nullable` |
| JSON Schema 2020-12 | `https://json-schema.org/draft/2020-12/json-schema-core` | `$id`、`$schema`、版本化和 validator 行为 |
| React / Vite | `https://react.dev/learn/start-a-new-react-project`、`https://vite.dev/guide/` | 当前脚手架、构建和开发服务器配置 |
| Playwright | `https://playwright.dev/docs/test-snapshots` | 当前 Story 的桌面快照环境、mask 和 threshold |
| SQLite / modernc SQLite | `https://www.sqlite.org/wal.html`、`https://www.sqlite.org/foreignkeys.html`、`https://gitlab.com/cznic/sqlite/-/blob/master/CHANGELOG.md` | WAL 返回值、每连接 foreign keys、运行时 SQLite 安全下限与 driver 对应版本；M1 起必须复核 |
| MCP | `https://modelcontextprotocol.io/specification/` | 只有正式 Story 明确引入 MCP 时才复核 lifecycle/tools；M1～M5 不因未实现 MCP 被阻断 |

版本复核按当前 Story 需要运行：

```bash
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go list -m -versions github.com/eino-contrib/jsonschema
go list -m modernc.org/sqlite
```

记录当前固定版本、官方可见稳定版本、是否升级、决定理由和升级后必须重跑的命令。不得自动采用 alpha/beta。

## 第四步：核对数据、fixture 与 provider 合同

- M1 只允许一个安全 fixture、一个只读 capability、一个页面和一个结果；真实 FOBrain、真实 LLM、Action 和 history data 均禁止。
- M2 只复核当前 Story 涉及的精确新增查询、七项代表性读取和 `business_list`；不要求 24 工具齐套。
- M3 复核 snapshot、result_ref、locator、SSE cursor、SQLite schema epoch 和重启恢复材料。
- M4 复核全部人员安全字段、stable identity、mock mutation、确认/取消、幂等、lease、verifier、reconcile 和操作记录；真实 provider mutation 必须为 0。
- M5 对每个动作单独复核目标接口、actor/owner/recipient identity、允许状态、写入请求、错误语义、独立 readback 和恢复方法。
- M6 只复核已通过门禁且已转换为正式 Story 的动作，不批量解锁其他动作。

所有 fixture 必须脱敏、版本化、可验证来源，并覆盖当前 Story 的空、部分、失败或安全拒绝边界。接口存在不等于产品事实完整。

## 第五步：旧项目只读核对

旧项目路径固定为 `/Users/vick/Desktop/project/ai-agent`，只允许用于核对当前 Story 对应的产品行为、验收基线和安全边界：

- 不修改旧项目。
- 不复制旧 runtime 类型、执行 loop、tool pair、Workbench API、数据库或 DOM/CSS。
- 不把旧项目通过记录当成新项目通过证据。
- 不以旧项目的全量 24 工具、connector、工单或移动端范围扩大当前 Story。

只读取与当前 Story 直接相关的文件；无直接映射时记录“不需要旧项目参考”。

## 第六步：M5 写域专项授权

任何目标部署 mutation 前必须同时记录：

- 授权人、授权日期、环境、动作类型和最大写入条数。
- 可恢复、非关键的 stable sample identity。
- 写前读取、唯一 mutation、独立 readback、恢复原状态和恢复后 readback 步骤。
- 超时、响应未知、部分失败和恢复失败时的停止条件。
- 允许修改的目录仅为 `scripts/acceptance/`、`docs/fixtures/`、`docs/acceptance-records/`、`_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md`。

缺少任一项必须 `exit 2`、`mutation=0` 并保持 Gate `BLOCKED`。Gate Evidence 通过也不允许注册生产 capability。

## 输出产物

每次复核生成：

```text
docs/acceptance-records/pre-development-validation-<story-id>-YYYY-MM-DD.md
```

至少包含：

- Story 与里程碑信息。
- 前置 Story、Gate、OQ 和 readiness 状态。
- 官方资料访问日期与版本结论。
- 数据来源、fixture/live 授权和敏感字段边界。
- 架构逐项确认结果。
- 受影响 schema、fixture、ADR、目录和验收命令。
- mutation 上限。
- 最终结论：`ALLOW` 或 `BLOCKED`，以及客观理由。

## 阻断条件

出现以下任一情况必须阻断：

- Implementation Readiness 未允许进入当前 Story。
- 前置 Story、适用门禁或开放问题未满足。
- 需求、数据来源、输入输出、验收标准或旧项目参考含义不清。
- 官方 API / spec 与当前设计冲突，或固定依赖版本无法支持所需行为。
- StructuredResult、Product Facts、schema、fixture、OpenAPI 或 generated contract 之间存在第二套真相。
- 需要硬编码工具名、关键词、provider 字段、凭据或环境值才能继续。
- 需要复制、兼容或降级到旧 runtime、旧数据库、旧接口或旧 UI。
- 未确认/未授权路径可能发生 mutation。
- M5 缺少授权、可恢复样本、独立回读或恢复方案。
- 验收命令不存在、无法区分 `PASS` 与 `exit 2`，或证据可能泄漏敏感信息。
