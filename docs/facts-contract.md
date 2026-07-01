# Product Facts 契约

本文定义 Workbench、Action API、replay、audit、Inspector 的同源事实模型。实现时不得为不同出口各自拼装业务结果。

## 状态枚举

### Run.status

| 值 | 含义 | 终态 |
| --- | --- | --- |
| `created` | 已创建但未开始执行 | 否 |
| `running` | 模型或工具正在执行 | 否 |
| `waiting` | 等待用户审批或澄清 | 否 |
| `succeeded` | 正常完成 | 是 |
| `failed` | 安全失败或不可恢复错误 | 是 |
| `cancelled` | 用户取消 | 是 |
| `stopped` | 系统或用户停止 | 是 |

### ToolCall.status

| 值 | 含义 |
| --- | --- |
| `queued` | 已计划调用 |
| `running` | 正在执行 |
| `succeeded` | 已返回安全结构化结果 |
| `failed` | 执行失败并产生安全错误 |
| `cancelled` | 已取消 |

### PendingInteraction.kind

| 值 | 含义 |
| --- | --- |
| `approval` | 高风险操作审批 |
| `clarification` | 用户补充实体、目标或参数 |

### PendingInteraction.status

| 值 | 含义 |
| --- | --- |
| `waiting` | 等待用户输入 |
| `submitted` | 已提交恢复数据 |
| `approved` | 审批通过 |
| `rejected` | 审批拒绝 |
| `cancelled` | 用户取消 |
| `expired` | 超时 |
| `consumed` | 恢复引用已消费 |

### Run / PendingInteraction 状态关系

| 场景 | PendingInteraction.status | Run.status | 说明 |
| --- | --- | --- | --- |
| 等待审批或澄清 | `waiting` | `waiting` | 暂停执行 |
| 审批通过 | `approved` -> `consumed` | `running` | 恢复执行 |
| 澄清提交 | `submitted` -> `consumed` | `running` | 恢复执行 |
| 审批拒绝 | `rejected` | `cancelled` | 不执行 mutation |
| 澄清取消 | `cancelled` | `cancelled` | 不继续执行原工具 |
| 等待超时 | `expired` | `failed` | 返回安全超时摘要 |

## 字段模型

### Run

| 字段 | 来源 | 可展示 | 说明 |
| --- | --- | --- | --- |
| `run_id` | 项目生成 | 是 | 外部稳定运行标识 |
| `workspace_id` | 请求 | 是 | 工作区 |
| `status` | Product Facts | 是 | 使用 Run.status |
| `created_at` / `updated_at` | store | 是 | ISO 8601 |
| `model_label` | config 安全名 | 是 | 不包含密钥、base url token |
| `safe_error` | Safety Gate | 是 | 失败时的安全错误摘要 |

### Turn

| 字段 | 来源 | 可展示 | 说明 |
| --- | --- | --- | --- |
| `turn_id` | 项目生成 | 是 | 用户轮次 |
| `run_id` | Run | 是 | 关联运行 |
| `role` | 请求或模型事件 | 是 | `user` / `assistant` / `system_notice` |
| `content` | Safety Gate | 是 | 用户文本或安全 assistant 文本 |
| `sequence` | store | 是 | 单调递增 |

### ToolCall

| 字段 | 来源 | 可展示 | 说明 |
| --- | --- | --- | --- |
| `tool_call_id` | Eino/tool adapter | 是 | 产品安全 id |
| `tool_id` | registry | 详情可见 | 主聊天不直接展示 |
| `display_name` | registry | 是 | 产品化中文名 |
| `status` | adapter | 是 | 使用 ToolCall.status |
| `args_hash` | Safety Gate | 否 | 仅审计内部定位 |
| `args_preview` | Safety Gate | 是 | allowlist 摘要 |
| `started_at` / `ended_at` | adapter | 是 | 工具时间线 |

### ToolResult

| 字段 | 来源 | 可展示 | 说明 |
| --- | --- | --- | --- |
| `tool_call_id` | ToolCall | 是 | 关联工具调用 |
| `structured_result` | StructuredResult Safety Gate | 是 | 唯一结果材料，当前 schema 为 `tool.structured_result.v1` |
| `presentation` | Product Mapper | 是 | 只能从 structured_result 派生 |
| `safe_error` | Safety Gate | 是 | 失败摘要 |
| `result_ref` | 项目生成 | 是 | 安全引用，不可反查 raw payload |

工具 adapter/provider 只能产出 StructuredResult candidate。candidate 必须先通过 product 层 Safety Gate，转成 `facts.StructuredResultRef` 后才能写入 Product Facts。`execution.EventMapper` 是 assistant/tool runner event 写入路径的强制接入点，内存 repository 和 SQLite repository 都保留 schema/result_ref/summary 的最后一道 unsafe material 防线。

### PendingInteraction

| 字段 | 来源 | 可展示 | 说明 |
| --- | --- | --- | --- |
| `pending_id` | 项目生成 | 是 | 稳定 pending 标识 |
| `kind` | execution | 是 | approval / clarification |
| `status` | Product Facts | 是 | 使用 PendingInteraction.status |
| `question` / `risk_summary` | Safety Gate | 是 | 澄清问题或审批风险摘要 |
| `options` | Product Mapper | 是 | 候选项或审批动作 |
| `resume_ref` | secure ref store | 是 | 不可复用安全引用 |
| `checkpoint_ref` | checkpoint store | 否 | 不暴露给前端/API |
| `expires_at` | policy | 是 | 可选 |

### AuditEvent

| 字段 | 来源 | 可展示 | 说明 |
| --- | --- | --- | --- |
| `audit_id` | 项目生成 | 是 | 审计事件 id |
| `run_id` | Run | 是 | 关联运行 |
| `event_type` | Product Facts | 是 | message/tool/pending/resume/safety/error |
| `safe_summary` | Safety Gate | 是 | 安全摘要 |
| `actor` | auth/context | 是 | user/system/tool |
| `redaction` | Safety Gate | 是 | 脱敏策略说明 |
| `created_at` | store | 是 | ISO 8601 |

## 投影矩阵

| 事实 | Workbench View | SSE | ActionResult | Replay | Inspector | Audit |
| --- | --- | --- | --- | --- | --- | --- |
| Run | 状态栏、运行摘要 | run patch | status、safe_error | 运行摘要 | runtime tab | run events |
| Turn | 主聊天 | message delta/replace | final_answer | 消息回放 | evidence | message event |
| ToolCall | 工具卡 | tool patch | tool status | 工具回放 | evidence/runtime | tool started |
| ToolResult | 工具卡展开 | tool result patch | structured tool results | 工具结果回放 | structured tab | tool completed |
| PendingInteraction | 审批/澄清卡 | pending patch | waiting refs | pending 回放 | runtime tab | pending/resume |
| AuditEvent | 审计列表 | 可选 audit patch | 不直接返回全集 | 审计回放 | audit tab | 审计事实 |

## 幂等规则

- `event_id` 必须稳定；重复 SSE 事件不得改变最终状态。
- 同一个 `resume_ref` 只能成功消费一次。
- duplicate resume 返回当前 PendingInteraction 终态，不重复执行工具或 mutation。
- rejected / cancelled / expired 后不得 approve 或 submit。
- `ActionResult` 查询同一 run 时必须返回同源事实投影，不能重新执行。
- 工具 mutation 必须以 `idempotency_key` 或业务安全引用防止重复提交。

## 存储不变量

- `run_id`、`turn_id`、`tool_call_id`、`pending_id`、`audit_id` 在 workspace 内唯一。
- 同一 `run_id` 内的 `sequence` 单调递增；SSE `event_id` 必须由稳定事实 id 和 sequence 派生。
- Product Facts 写入必须是 append 或显式状态迁移；不得由 Workbench、Action API、audit 或 replay 反向写工具结果。
- `ToolResult.structured_result` 写入后不可被展示层修改；presentation 只能重新派生，不能成为事实来源。
- 恢复引用消费、pending 状态迁移和 mutation idempotency 必须在同一事务边界内完成。
- replay 按 facts sequence 重建，遇到重复 event_id 必须去重，遇到乱序事件必须按 sequence 排序。

## 错误语义

| 错误 | 产品表现 | 是否可重试 |
| --- | --- | --- |
| schema invalid | 安全错误摘要 | 修复请求后可重试 |
| safety blocked | 安全 run notice | 否 |
| provider timeout | 工具失败卡 | 可重试 |
| checkpoint missing | 恢复失败安全摘要 | 通常不可重试 |
| resume consumed | 返回已消费状态 | 否 |
| credential missing | connector 状态卡 | 配置后可重试 |
