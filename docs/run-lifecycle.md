# Run Lifecycle

本文定义 Run、Turn、ToolCall 和 PendingInteraction 的产品状态机。Eino checkpoint/interrupt 只提供内部暂停恢复能力；产品状态必须由 Product Facts 管理。

## 状态层级

```text
Run
  -> Turn
  -> ToolCall / ToolResult
  -> PendingInteraction
```

状态枚举和事实字段见 `facts-contract.md`。本文只定义跨状态行为。

## Run 状态

| 状态 | 含义 | 可继续执行 |
| --- | --- | --- |
| `created` | 已创建未开始 | 是 |
| `running` | 模型或工具正在运行 | 是 |
| `waiting` | 等待 approval 或 clarification | 只能 resume/cancel/timeout |
| `succeeded` | 正常完成 | 否 |
| `failed` | 不可自动恢复或安全失败 | 否 |
| `cancelled` | 用户取消 | 否 |
| `stopped` | 系统或用户停止 | 否 |

终态必须幂等。重复 stop/cancel/resume 查询只能返回当前事实投影，不得重新执行工具或 mutation。

## Cancel / Stop / Timeout

| 操作 | 触发者 | 行为 |
| --- | --- | --- |
| cancel | 用户 | 当前 run 进入 `cancelled`；pending 进入 `cancelled`；不得继续工具或 mutation。 |
| stop | 用户或系统 | 当前 run 进入 `stopped`；保留已完成事实；不得新增业务工具结果。 |
| provider timeout | 系统 | 当前 tool 进入 `failed`；run 根据是否可继续转 `failed` 或 partial safe result。 |
| pending timeout | 系统 | pending 进入 `expired`；run 进入 `failed`，用户需重新发起。 |
| budget exceeded | 系统 | 预算判断层已确认超限后触发；`running`/`waiting` run 进入 `failed`，写入 `safe_error=budget_exceeded`，active tool 进入 `cancelled`，waiting pending 进入 `expired`，并写入 `event_type=budget` audit。 |
| checkpoint missing | 系统 | resume 失败，run 进入安全失败或保持 waiting 并返回安全错误，按实现阶段门禁固定。 |

终态 run 对重复 lifecycle 请求保持幂等 no-op，返回当前事实投影，不再重新写 audit、tool 或 pending。`budget_exceeded` 的“只允许 running/waiting”指会改变事实的有效状态范围；对已终态 run 发送该 action 只能得到当前终态投影。

## Retry

Retry 只能基于安全 retry policy：

- schema invalid：用户修正后可重新发起。
- provider timeout：只可重试只读工具；mutation 不得自动重试。
- credential missing：配置凭据后重新发起。
- approval/clarification 终态后不得通过 retry 复用旧 resume_ref。

retry 请求必须有新的 `client_request_id` 或明确 retry action，并写入 audit。

## Mutation 幂等

写域 mutation 必须满足：

- 未 approval 不执行。
- approve 后同一 idempotency key 最多执行一次。
- reject/cancel/expired 后不得执行。
- duplicate approve 返回当前终态，不重复 mutation。
- mutation 执行结果、业务安全引用和 policy decision 必须写入 Product Facts。

## Replay 规则

Replay 从 Product Facts 重建：

- 已完成的 Turn、ToolCall、ToolResult 保留。
- failed/cancelled/stopped 显示 run notice。
- pending 终态显示只读卡片。
- replay 不依赖前端缓存或 Eino raw event。

## 验收

- cancel waiting run 不继续执行工具。
- stop running run 进入 `stopped` 并可 replay。
- provider timeout 产生安全错误，不泄漏 raw provider body。
- pending timeout 后不得 submit/approve。
- budget exceeded 后 active tool 取消、waiting pending 过期、audit/replay 同源且不泄漏 raw 或 credential。
- duplicate resume/approve/cancel 幂等。
- mutation approval 后最多执行一次。

参考资料：

- Eino HITL：`https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/`
- Eino Interrupt & CheckPoint：`https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/`
