# Workbench 视觉验收矩阵

本文把 Workbench 视觉验收固定为可执行矩阵。目标是保持当前产品的专业、密集、聊天优先三栏工作台体验，而不是复用旧 DOM 或旧前端代码。

## 验收层级

| 层级 | 使用时机 | 必须证据 |
| --- | --- | --- |
| Contract fixture | 每个前端任务 | 从 `docs/fixtures/workbench-view-success.json` 渲染，不能按 `tool_id` 写业务分支。 |
| Block screenshot | 每个可见 block | desktop/mobile 截图、selector、状态说明。 |
| Real service | SSE、工具、审批、澄清、回放相关任务 | 运行服务截图、SSE payload、ActionResult 或 replay evidence。 |

Phase 2 必须补充目标截图路径、block crop 坐标和 Playwright baseline。没有目标图的状态仍必须有 fixture screenshot 和人工 pass/fail 记录。

机器可读视觉矩阵固定在：

```text
docs/schemas/visual_evidence_matrix.v1.schema.json
docs/fixtures/visual-evidence-matrix.json
```

该矩阵必须登记在 `docs/fixtures/manifest.json`。`schema-test` 会检查 8 个必选 block 是否精确覆盖，以及每个视觉状态引用的 fixture 是否存在。

## 产品场景覆盖

| 场景 | 必须行为 |
| --- | --- |
| 普通聊天 | 用户消息、assistant delta、最终回答按时间线展示，不出现空工具区域。 |
| 单工具 | 一个紧凑工具卡，折叠态可读，展开态展示 StructuredResult。 |
| 多工具 | 保持真实顺序，不伪造成固定 workflow。 |
| assistant 插入文本 | 工具前后中间文本保持原始顺序。 |
| 审批等待 | 审批卡出现在等待点，resume 前不得显示 mutation 成功。 |
| 澄清等待 | 澄清卡展示候选或输入，提交后变只读。 |
| 失败/空/部分结果 | 状态留在原工具位置，显示安全摘要和 evidence。 |
| stop/cancel/timeout | 显示 run notice，不重复终态卡。 |
| 多轮追问 | 只能引用 prior safe evidence 和 StructuredResult，不把内部 ref 当正文。 |

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
- 主聊天出现 raw prompt、provider payload、credential ref、Authorization、resume token、internal reason code 或 raw JSON dump。
- Workbench、replay、audit、ActionResult 对同一 run 展示不一致。
- 工具卡展示内容不是从 StructuredResult 派生。
- Inspector 直接展示 raw runtime event 或 raw provider payload。
