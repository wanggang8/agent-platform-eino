# 审批流程契约

本文定义 `PendingInteraction(kind=approval)` 的产品行为。审批用于所有外部写入、工单状态变更、触发扫描、创建或修改外部对象等高风险操作。

## 触发条件

- capability 声明 `approval_required=true`。
- `side_effect` 为 `write_external` 或风险策略判定为 `high`。
- provider 返回 `pending_approval` StructuredResult candidate。
- 权限策略返回 `approval_required`。

未生成审批事实前，不得调用真实 mutation。

## 请求结构

| 字段 | 说明 |
| --- | --- |
| `pending_id` | 审批卡 id |
| `run_id` | 关联运行 |
| `kind` | 固定为 `approval` |
| `status` | 初始为 `waiting` |
| `operation_name` | 中文操作名 |
| `risk_summary` | 安全风险摘要 |
| `target_summary` | 影响对象，如工单、资产、业务系统 |
| `proposed_args_summary` | allowlist 参数摘要 |
| `resume_ref` | 不可复用恢复引用 |
| `expires_at` | 可选超时时间 |

审批卡不得包含 raw tool args、credential、checkpoint id、interrupt id 或 provider payload。

## 提交结构

| 字段 | 说明 |
| --- | --- |
| `resume_ref` | 当前审批卡引用 |
| `decision` | `approve` / `reject` |
| `client_request_id` | 客户端幂等 id |
| `comment` | 可选安全备注 |

提交数据先过 schema 和 Safety Gate，再映射为 Eino resume data。`approve` 后才允许执行真实 mutation；`reject` 后本次 mutation 永不执行。

## 状态机

```text
waiting
  -> approved
  -> consumed

waiting
  -> rejected

waiting
  -> cancelled

waiting
  -> expired
```

| 当前状态 | 用户动作 | 新状态 | Run.status | 说明 |
| --- | --- | --- | --- | --- |
| `waiting` | approve | `approved` -> `consumed` | `running` | 恢复执行并允许一次 mutation |
| `waiting` | reject | `rejected` | `cancelled` | 不执行 mutation |
| `waiting` | cancel | `cancelled` | `cancelled` | 用户取消审批，不执行 mutation |
| `waiting` | timeout | `expired` | `failed` | 返回安全超时摘要 |
| `approved` / `consumed` | duplicate approve | 不变 | 当前 Run 状态 | 返回当前状态，不重复 mutation |
| `rejected` / `cancelled` / `expired` | approve | 不变 | 不变 | 返回安全错误，不恢复 |

## 行为规则

- `waiting` 时 Run.status 必须为 `waiting`。
- 审批前 mutation count 必须为 0。
- approve 后同一 idempotency key 最多执行一次 mutation。
- reject、cancel、expired 后不得 approve。
- 进程重启后 waiting 审批必须能恢复。
- 审批请求、审批结果、mutation 执行结果必须进入 audit 和 replay。
- `approval_refs` 可出现在 ActionResult；可复用 `resume_token` 不得出现在任何 JSON 输出。

## SSE patch

```json
{
  "event_id": "evt_approval_001",
  "type": "pending.updated",
  "run_id": "run_001",
  "pending": {
    "pending_id": "pending_approval_001",
    "kind": "approval",
    "status": "waiting",
    "operation_name": "更新工单状态",
    "risk_summary": "将修改外部 Fobrain 工单状态",
    "target_summary": "ticket:T-1001 -> fixed",
    "resume_ref": "resume_ref_..."
  }
}
```

## ActionResult waiting 示例

```json
{
  "schema_version": "eino_action_result.v1",
  "run_id": "run_001",
  "status": "waiting",
  "waiting": {
    "kind": "approval",
    "question": "是否批准更新工单状态？",
    "approval_refs": ["resume_ref_..."]
  }
}
```

## 测试 fixture

必须覆盖：

- approval requested。
- approval approved。
- approval rejected。
- approval timeout。
- approval cancelled。
- duplicate approve。
- approve after reject。
- approve after restart。
- mutation not executed before approval。
- mutation executed once after approval。
- approval card redaction。
