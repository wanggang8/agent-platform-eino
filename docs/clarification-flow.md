# 澄清流程契约

本文定义 `PendingInteraction(kind=clarification)` 的产品行为。该流程用于实体消歧、目标不明确、参数不足等需要用户补充信息的场景。

## 触发条件

- 同名实体存在多个候选项。
- 用户问题缺少必须参数，且不能安全默认。
- 工具或 provider 返回需要用户选择的业务候选项。

## 请求结构

| 字段 | 说明 |
| --- | --- |
| `pending_id` | 澄清卡 id |
| `run_id` | 关联运行 |
| `kind` | 固定为 `clarification` |
| `status` | 初始为 `waiting` |
| `question` | 展示给用户的问题 |
| `input_mode` | `single_choice` / `multi_choice` / `free_text` / `mixed` |
| `candidates` | 候选项列表 |
| `resume_ref` | 不可复用恢复引用 |
| `expires_at` | 可选超时时间 |

### candidate

| 字段 | 说明 |
| --- | --- |
| `candidate_ref` | 安全候选引用 |
| `label` | 用户可读名称 |
| `description` | 可选说明 |
| `entity_type` | person / asset / vulnerability / department / business |
| `safe_fields` | 可展示字段 |

## 提交结构

| 字段 | 说明 |
| --- | --- |
| `resume_ref` | 当前澄清卡引用 |
| `selected_candidate_refs` | 用户选择的候选项 |
| `free_text` | 用户补充文本 |
| `client_request_id` | 客户端幂等 id |

提交数据先过 schema 和 Safety Gate，再转换为 Eino resume data。模型恢复上下文只能注入安全候选、用户选择和安全摘要，不能注入 raw provider payload。

## 状态机

`PendingInteraction.status` 和 `Run.status` 必须分开理解。澄清卡自身不会进入 `running`；用户提交后，澄清卡进入终态，Run 才恢复为 `running`。

```text
PendingInteraction.status:

waiting
  -> submitted
  -> consumed

waiting
  -> cancelled

waiting
  -> expired
```

| 当前 PendingInteraction.status | 用户动作 | 新 PendingInteraction.status | Run.status | 说明 |
| --- | --- | --- | --- | --- |
| `waiting` | submit | `submitted` -> `consumed` | `running` | 恢复数据验证通过后继续执行 |
| `waiting` | cancel | `cancelled` | `cancelled` | 默认取消本次运行，不继续执行原工具 |
| `waiting` | timeout | `expired` | `failed` | 返回安全超时摘要；用户需重新发起任务 |
| `submitted` / `consumed` | duplicate submit | 不变 | 当前 Run 状态 | 返回当前状态，不重复执行 |
| `cancelled` / `expired` | submit | 不变 | 不变 | 返回安全错误，不恢复 |

## 行为规则

- `waiting` 时 Run.status 为 `waiting`。
- `submitted` 后立刻消费恢复数据；成功恢复后 PendingInteraction.status 变为 `consumed`，Run.status 变为 `running`。
- `submit` 必须先确认 `checkpoint_ref` 绑定的内部 checkpoint 存在，缺失时不得消费 pending。
- `cancelled` 后 Run.status 变为 `cancelled`，不得继续执行原工具。
- `cancel` 不恢复执行，不依赖 checkpoint 存在。
- `expired` 后不得继续 submit；用户需重新发起任务。
- 进程重启后，`waiting` 澄清必须能从 store 恢复。
- duplicate submit 不得重复调用工具或 mutation。

## SSE patch

```json
{
  "event_id": "evt_001",
  "type": "pending.updated",
  "run_id": "run_001",
  "pending": {
    "pending_id": "pending_001",
    "kind": "clarification",
    "status": "waiting",
    "question": "请选择要查询的人员",
    "input_mode": "single_choice",
    "resume_ref": "resume_ref:...",
    "candidates": [
      {"candidate_ref": "candidate:fobrain:person:1", "label": "张三", "entity_type": "person"}
    ]
  }
}
```

## ActionResult waiting 示例

```json
{
  "run_id": "run_001",
  "status": "waiting",
  "waiting": {
    "kind": "clarification",
    "question": "请选择要查询的人员",
    "resume_refs": ["resume_ref:..."],
    "input_mode": "single_choice",
    "candidates": [
      {"candidate_ref": "candidate:fobrain:person:1", "label": "张三", "entity_type": "person"}
    ]
  }
}
```

## 测试 fixture

必须覆盖：

- clarification requested。
- clarification submitted。
- clarification cancel。
- clarification timeout。
- duplicate submit。
- submit after cancel。
- submit after restart。
- candidate redaction。
