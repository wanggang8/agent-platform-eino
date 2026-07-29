# State Machines

## 产品状态层级

```mermaid
stateDiagram-v2
  [*] --> created
  created --> running
  running --> waiting: approval or clarification
  waiting --> running: valid resume
  running --> succeeded: orchestration safely completed
  running --> failed: unsafe or manual attention
  running --> cancelled
  running --> stopped
  waiting --> failed: expired or unsafe
  waiting --> cancelled
  waiting --> stopped
  succeeded --> [*]
  failed --> [*]
  cancelled --> [*]
  stopped --> [*]
```

Run 只使用上述七态。approval/clarification 由 PendingInteraction 表达；写入和核验阶段由 ActionItem/Attempt 表达。

## Prepare 与 Confirm

```mermaid
sequenceDiagram
  actor U as 用户
  participant I as Chat或Action API
  participant X as execution
  participant R as facts repository
  participant C as continuation backend
  participant P as provider

  U->>I: 基于result_ref提出动作和人工参数
  I->>X: PrepareAction
  X->>R: 解析快照、身份、策略和权限
  X->>R: 保存不可变ActionDraft
  X->>C: stage continuation
  X->>R: 发布waiting Pending与FactEvent
  R-->>I: draft/version/digest和确认摘要
  U->>I: approve或reject
  I->>X: ConfirmAction
  X->>R: 原子消费ref、复核并预留attempt
  alt reject/cancel/expired
    X-->>I: 零外部写入的终态
  else approve
    X->>C: resume或推进deterministic continuation
    loop 每个ActionItem
      X->>P: 使用稳定mutation key写入
      X->>P: 回读安全事实
      X->>R: 保存VerificationEvidence与item outcome
    end
    R-->>I: 同源结果投影
  end
```

聊天 continuation backend 是真实 Eino checkpoint/interrupt；确定性 Action API backend 是项目 continuation record。两者实现统一 `ContinuationRef`，不能互相冒充。

## ActionDraft 与 Item 状态

```mermaid
stateDiagram-v2
  state Draft {
    [*] --> prepared
    prepared --> waiting_confirmation
    waiting_confirmation --> approved
    waiting_confirmation --> rejected
    waiting_confirmation --> expired
    approved --> executing
    executing --> completed
  }

  state Item {
    [*] --> prepared_item
    prepared_item --> reserved_not_sent: semantic claim acquired
    reserved_not_sent --> sent
    reserved_not_sent --> cancelled_before_send: cancel/stop CAS + source-aware claim disposition
    sent --> verifying
    verifying --> succeeded
    verifying --> failed
    verifying --> reconciling
    reconciling --> succeeded
    reconciling --> failed
    reconciling --> manual_attention
  }
```

聚合映射：

| Item 集合 | Run | ActionResult | Outcome |
| --- | --- | --- | --- |
| 全部 `succeeded` | `succeeded` | `completed` | `all_succeeded` |
| 成功与确定 `failed` 混合 | `succeeded` | `completed` | `partial` |
| 全部确定 `failed` | `succeeded` | `completed` | `none_succeeded` |
| 任一 `reconciling` | `running` | `running` | 未终结 |
| 任一 `manual_attention` 且存在成功项 | `failed` | `blocked` | `partial` |
| 任一 `manual_attention` 且无成功项 | `failed` | `blocked` | `none_succeeded` |
| 全部 `cancelled_before_send`，用户 cancel | `cancelled` | `cancelled` | `none_succeeded` |
| 全部 `cancelled_before_send`，系统／用户 stop | `stopped` | `stopped` | `none_succeeded` |

产品状态只采用以下映射：

| Item 状态／事实条件 | 产品状态 | 收敛规则 |
| --- | --- | --- |
| `reconciling`；mutation 可能已发生且仍在 verification policy 的 settle deadline／调用预算内 | 待核验 | Run 保持 `running`；沿同一 attempt 只回读，不得再次 mutation。 |
| `succeeded`；control held、provider accepted 且 conclusive 真实 readback 与 expected 一致 | 完成 | 终态，不得重复执行。 |
| definitive `failed`；provider 明确拒绝或可证明 mutation 未完成 | 待排查 | 条目终态；若要重试，必须重新人工判断并创建新 draft。 |
| `manual_attention`；deadline／预算结束仍无法证明目标，或失租约导致结果未知 | 待排查 | Run=`failed`、ActionResult=`blocked`；保留安全核验证据。 |
| `cancelled_before_send`；cancel/stop CAS 证明 mutation count=0 | 已取消／已停止 | attempt 与 Run 终态；首次 claim=`released_zero_write`，failed-retry claim 恢复 predecessor failed；child 审计 tombstone 保留。 |

VerificationEvidence reducer 的优先级不可配置：`control_certainty=lost` 首先进入 manual_attention；held + provider accepted + conclusive real-readback match 才 succeeded；held + provider rejected + policy 允许的 conclusive non-application 才 definitive failed；其余在 deadline／预算内 reconciling，耗尽后 manual_attention。`rejected|unknown + match` 不得成功；conclusive rejection 与 control-loss 不得伪造 observed StructuredResult。

verification deadline、已用调用次数和下一次允许回读时间必须持久化，服务重启不得重新计算或重置预算。`manual_attention` 明确由 FOBrain 系统负责人排查；本产品只呈现事实和排查指引，不发送通知、消息或工单。

definitive `failed` 与 `manual_attention` 都追加 `attention_status=open`，但不改变原执行终态。首版没有 attention 关闭命令；未来只能追加版本化 closure fact 和 audit，不能改写原 ActionResult。

## Semantic Mutation Claim

```mermaid
stateDiagram-v2
  [*] --> unclaimed
  unclaimed --> reserved_not_sent: 首次Confirm事务成功claim
  reserved_not_sent --> released_zero_write: 首次claim取消且零写入
  released_zero_write --> reserved_not_sent: 后续普通新决策取得claim
  reserved_not_sent --> sent
  sent --> verifying
  verifying --> reconciling
  verifying --> succeeded
  verifying --> failed
  reconciling --> succeeded
  reconciling --> manual_attention
  failed --> retry_reserved_not_sent: 新人工retry事务转移claim并保存predecessor
  retry_reserved_not_sent --> failed: sent前取消时恢复predecessor claim
  retry_reserved_not_sent --> sent
```

`(workspace_id, semantic_mutation_key)` 是 claim 唯一键。普通草案命中任何 active/retained claim 均失败；只有逐条绑定 definitive failed item 且 lineage 完整的新人工 retry 可以事务转移 claim。claim 插入／转移与 attempt 预留同事务，确保并发确认只有一个草案进入 attempt phase=`reserved_not_sent`；图中的 `retry_reserved_not_sent` 是 claim provenance，不是新的 Attempt phase。cancel/stop CAS 必须写 mutation count=0 和权威 claim disposition：首次 claim 进入 `released_zero_write`；retry transferred claim 恢复 predecessor owner/status=`failed`。取消 child 的 audit tombstone 不得替代该权威状态。若 sent CAS 先提交，处置失败并返回 `action_already_sent`。

## Confirm 事务与崩溃恢复

```mermaid
flowchart TD
  A["校验resume ref与request digest"] --> B["事务：消费ref、固定Pending终态"]
  B --> C["事务：校验draft/actor/policy/version"]
  C --> D["事务：draft approved、预留attempt、追加FactEvent"]
  D --> E["提交事务"]
  E --> F["resume Eino或推进项目continuation"]
  F --> C{"cancel/stop是否先到达"}
  C -->|是| X["事务CAS：cancelled_before_send + 零写入证据 + 来源感知claim处置"]
  C -->|否| S["事务CAS：reserved_not_sent→sent，固定deadline/预算/epoch"]
  S --> G["首次且唯一外部mutation调用"]
  G --> H{"是否有确定响应"}
  H -->|是| I["回读并核验"]
  H -->|否| J["reconciling：只回读不盲重试"]
  J --> I
  I --> K{"可证明结果"}
  K -->|是| L["item终态"]
  K -->|否，预算未结束| N["reconciling：待核验"]
  N --> I
  K -->|否，预算已结束| M["manual_attention：待排查"]
```

进程在 Confirm 事务提交前崩溃：确认未消费，可按同 request 重试。Confirm 提交后、`sent` 事务前崩溃：恢复同一 `reserved_not_sent` attempt，可执行一次首次调用；若已收到 cancel/stop，则只能用互斥 CAS 进入 `cancelled_before_send`，并将首次 claim 置为 `released_zero_write` 或恢复 failed-retry 的 predecessor failed claim。`sent` CAS 先提交后，无论调用是否真正到达 provider、响应是否丢失、租约是否变化或随后收到 cancel/stop，都只能回读／reconcile，不能再次写入或伪报取消。每次回读先原子扣减调用预算；崩溃不返还。

## Lease 与 Fencing

后台 executor 取得 `(run_id, lease_epoch, expires_at)`。所有事实提交携带 epoch；repository 拒绝小于当前 epoch 的迟到提交。新 owner 接管未知外部调用时只能回读或进入人工排查，不能凭租约接管而盲目重发。

## SSE 高水位

```mermaid
sequenceDiagram
  participant W as Workbench
  participant H as HTTP View
  participant S as SSE
  participant F as FactEvent Store

  W->>H: GET current view
  H->>F: 原子读取aggregate与max sequence
  F-->>H: view + snapshot_sequence=N
  H-->>W: 安全投影和N
  W->>S: subscribe after N
  S->>F: 读取sequence>N
  F-->>S: upsert/tombstone events
  S-->>W: 严格递增事件
  alt Last-Event-ID未知或过旧
    S-->>W: view.replaced + 新高水位
  end
```

事件 payload 只允许完整实体 upsert 或显式 tombstone；数组默认 replace，只有 schema 明确声明 keyed collection 时可按 key 合并。客户端忽略 `sequence <= applied_sequence` 的事件。

## 操作记录投影

```mermaid
flowchart LR
  U["当前凭据主体 actor/workspace"] --> Q["操作记录查询 + 当前权限证明"]
  Q --> F["Product Facts + ActionResult + attention + safe audit refs"]
  F --> L["IA-01 历史列表"]
  L --> D["IA-02 只读逐条详情"]
  C["conversation summary"] -. 禁止重建 .-> Q
  P["重新调用 FOBrain"] -. 禁止重建 .-> Q
  D -. 不自动创建 .-> A["ActionDraft / retry"]
```

`operation_at` 在记录首次发布时固定且永不变化：有 Pending 的流程始终使用 Pending `published_at`，无 Pending 的确定性流程使用服务端 decision time，后续审批只写 `decision_at`。首个请求不接受客户端 `as_of`，由服务端生成 RFC3339 纳秒时间；转为 Asia/Shanghai 后，目标年月向前六个月，日取原日与目标月末较小值并保留时分秒／纳秒，边界包含。后端在一个 SQLite 事务中把窗口内记录与 Run waiting/running 或 attention open 的去重 union 物化为 immutable `OperationListSnapshot`，冻结完整有序安全列表行、记录版本和 total；它有独立 opaque `operation_list_snapshot_ref`，不是 Run-local SSE `snapshot_sequence`。后续页只读同一物化快照，其他 Run 变化不进入；cursor 绑定 scope、filters、as_of、cutoff、operation_list_snapshot_ref 与 last keys，篡改、丢失或 TTL 过期返回 `invalid_cursor`，前端不二次合并。查询要求 actor/workspace 精确相等，并按冻结访问策略证明快照内整批对象当前仍可见，unknown/denied 整批 fail closed，不过滤或重算。provider instance 改变不做身份别名。六个月不触发物理删除；首版 open attention 持续可发现，历史详情不推进状态，也不自动创建重试或新草案。
