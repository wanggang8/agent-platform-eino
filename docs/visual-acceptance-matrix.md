# Workbench 视觉验收矩阵

本文把 Workbench 视觉验收固定为可执行矩阵。目标是保持当前产品的专业、密集、聊天优先三栏工作台体验，而不是复用旧 DOM 或旧前端代码。

## 验收层级

| 层级 | 使用时机 | 必须证据 |
| --- | --- | --- |
| Contract fixture | 每个前端任务 | 从 `docs/fixtures/workbench-view-success.json` 渲染，不能按 `tool_id` 写业务分支。 |
| Block screenshot | 每个可见 block | desktop/mobile 截图、selector、状态说明。 |
| Real service | SSE、工具、审批、澄清、回放相关任务 | 运行服务截图、SSE payload、ActionResult 或 replay evidence。 |

Phase 2 必须补充目标截图路径、block crop 坐标和 Playwright baseline。没有目标图的状态仍必须有 fixture screenshot 和人工 pass/fail 记录。

## 产品化视觉门禁

视觉验收必须按产品交付界面验收，不按调试工作台验收。所有进入 `web/eino-workbench/tests/__screenshots__/` 的截图必须满足：

- 可见文案使用中文产品语言；不得出现 `tool_id`、schema 名、内部枚举、provider 字段名、run id、checkpoint/resume 标识或英文调试文案。
- 工具结果只能展示处理后的安全标题、摘要、指标、事实表格和审计摘要；不得展示原始 JSON、`StructuredResult` dump、raw provider payload 或对象字段名。
- 连接器、凭据、运行状态和审计事件必须映射为中文产品状态，例如“可用”“已绑定”“当前工作区”“工具执行”，不能展示 `available`、`bound`、`workspace`、`tool` 等内部值。
- `StructuredResult`、`Product Facts` 等契约名可以出现在技术文档和源码注释中，但不能作为产品截图可见文本。
- Playwright 视觉用例必须在截图前断言产品可见文本不包含 ASCII 英文调试词、raw JSON 标记或敏感材料。

机器可读视觉矩阵固定在：

```text
docs/schemas/visual_evidence_matrix.v1.schema.json
docs/fixtures/visual-evidence-matrix.json
```

该矩阵必须登记在 `docs/fixtures/manifest.json`。`schema-test` 会检查 8 个必选 block 是否精确覆盖，以及每个视觉状态引用的 fixture 是否存在。

## 本地旧验收参考

旧项目最终视觉参考已复制到：

```text
docs/assets/legacy/workbench/
```

这些图片只作为目标参考，不是新项目通过结果。新项目必须用 Playwright 重新生成 `web/eino-workbench/tests/__screenshots__/` 和 `test-results/eino-workbench-*`。

| 用途 | 本地路径 |
| --- | --- |
| Canonical target | `docs/assets/legacy/workbench/canonical/workbench-target-ui-v1-2026-06-26.png` |
| Task15 region contact sheets | `docs/assets/legacy/workbench/task15-contact-sheets/` |
| Task15 target crops | `docs/assets/legacy/workbench/task15-target-crops/` |
| Task17 scenario contact sheets | `docs/assets/legacy/workbench/task17-contact-sheets/` |
| Final UI block references | `docs/assets/legacy/workbench/final-ui-blocks/` |
| Final acceptance reports | `docs/assets/legacy/workbench/reports/` |

## 产品场景覆盖

| 场景 | 必须行为 |
| --- | --- |
| 普通聊天 | 用户消息、assistant delta、最终回答按时间线展示，不出现空工具区域。 |
| 单工具 | 一个紧凑工具卡，折叠态可读，展开态展示安全结构化结果。 |
| 多工具 | 保持真实顺序，不伪造成固定流程。 |
| assistant 插入文本 | 工具前后中间文本保持原始顺序。 |
| 审批等待 | 审批卡出现在等待点，resume 前不得显示 mutation 成功。 |
| 澄清等待 | 澄清卡展示候选或输入，提交后变只读。 |
| 失败/空/部分结果 | 状态留在原工具位置，显示安全摘要和 evidence。 |
| stop/cancel/timeout | 显示 run notice，不重复终态卡。 |
| 多轮追问 | 只能引用已投影的安全证据和结构化结果，不把内部引用当正文。 |

## 必选可见 block

| Block id | Selector 建议 | 状态 |
| --- | --- | --- |
| `shell` | `[data-testid="workbench-shell"]` | desktop 三栏，mobile 单栏/tabs。 |
| `sidebar` | `[data-testid="workspace-sidebar"]` | workspace、session、search、new chat。 |
| `timeline` | `[data-testid="chat-timeline"]` | 空态、聊天态、工具态、等待态、失败态。 |
| `composer` | `[data-testid="chat-composer"]` | 输入、发送、disabled/running/waiting。 |
| `tool-card` | `[data-testid="tool-card"]` | pending/running/completed/failed/cancelled。 |
| `approval-card` | `[data-testid="approval-card"]` | waiting/approved/rejected/expired。 |
| `clarification-card` | `[data-testid="clarification-card"]` | waiting/submitted/cancelled/expired。 |
| `inspector` | `[data-testid="inspector"]` | evidence、structured、runtime、audit tabs。 |

Selectors 可以调整，但必须在本矩阵和 screenshot 脚本同一任务更新。缺 selector 视为失败。

## Display Type 覆盖

前端 presenter 只能按 `structured_result.schema_version` 和 `display_type` 选择布局：

| display_type | 用途 | 必须渲染 |
| --- | --- | --- |
| `entity_collection` | 资产、漏洞、工单、业务系统列表 | `columns[]`、`items[]`、pagination/empty。 |
| `entity_detail` | 单个对象详情 | `facts[]`、`sections[]`。 |
| `metrics_summary` | 统计或风险摘要 | `metrics[]`、可选 drill-down sections。 |
| `operation_result` | 写域结果或审批相关结果 | summary/facts/actions，并与审批卡一致。 |
| `connector_status` | 连接器健康和绑定状态 | safe facts，不显示 credential ref。 |
| `entity_resolution` | 候选消歧 | candidates 或 resolved entity，并驱动澄清卡。 |

禁止按 `tool_id`、旧工具别名、raw object keys 或旧 DOM 类名写 presenter 分支。

## 全局失败条件

- 第一屏不是可工作的聊天工作台。
- 桌面出现 hero、营销说明、手机预览或演示玩具布局。
- 主聊天出现 raw prompt、provider payload、credential ref、Authorization、resume token、internal reason code、英文调试词或 raw JSON dump。
- Workbench、replay、audit、ActionResult 对同一 run 展示不一致。
- 工具卡展示内容不是从 StructuredResult 派生。
- Inspector 直接展示 raw runtime event 或 raw provider payload。
