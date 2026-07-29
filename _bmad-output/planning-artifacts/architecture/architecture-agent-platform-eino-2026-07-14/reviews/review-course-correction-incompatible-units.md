# 纠偏后架构不兼容单元攻击评审

- 评审对象：`../ARCHITECTURE-SPINE.md`
- 评审范围：AD-10、AD-21、AD-23、AD-26，以及它们在 `execution / facts / product / httpapi / web` 间形成的接缝
- 评审方法：构造两个互不沟通的实现单元，要求双方都能逐字引用 Spine 证明自身合规，再检查组合后是否仍只有一个状态、发现范围、actor 边界和重试语义
- 评审性质：只读对抗性评审；不修改 Architecture、PRD、UX、SPEC、Epics 或代码

## 初审 Verdict（已由文末 Recheck 取代）

**NOT IMPLEMENTATION-SAFE。** AD-10/21/23 已固定大部分状态名称，AD-26 也确认了“操作记录 + 最近六个月”的产品方向，但四条规则还没有组成一个可交换协议。两个独立实现者可以完全遵守这些 AD，仍分别得到不同的 ActionResult 终态、不同的“待排查”长期可发现集合、不同的访问时权限判断，以及相反的外部 mutation 重发行为。最危险的不是页面文案不一致，而是同一条记录会在一个单元中被视为“已终结可从默认入口消失”，在另一个单元中被视为“未关闭必须永久可发现”；或在一个单元中被允许用同一幂等键技术重发，在另一个单元中被 AD-10 的禁令永久阻断。

当前 Spine 可以作为目标方向，不能直接作为 Story 1.3～1.6、1.12～1.15、1.22～1.24、1.36 独立实施后拼装的充分契约。至少需要在 v2 schema/port/OpenAPI/fixture 中补齐：权威 outcome 表、操作记录发现谓词、关闭生命周期、访问时 actor/permission 谓词，以及跨 draft mutation lineage 与允许调用矩阵。

## 可构造的两个合规实现单元

- **单元 A：控制面事实单元（`execution + facts`）。** 它把 verifier 的每次观察写为 VerificationEvidence，由 `execution` 迁移 ActionItem；全部 item 为 `succeeded/failed` 时按 AD-21 将 Run 终结为 `succeeded`，遇到 `manual_attention` 时将 Run 终结为 `failed`。它以 `confirmed_at` 作为“操作时间”，以 Run 是否 terminal 判断“未关闭”，repository 先按持久化 `actor_subject_id + workspace_id` 过滤。技术恢复只要复用 AD-11 的 mutation key，就允许 provider 调用层按其 idempotency 声明重发。

- **单元 B：产品传输单元（`product + httpapi + web`）。** 它不重建 execution 状态，只消费生成契约；看到 `failed/manual_attention` 都按 AD-21 显示“待排查”，并把这种仍需人工处置的产品标签理解为“未关闭”。它以最近一次状态变化时间排序操作记录，每次打开时重新验证当前 actor/workspace 及当前对象可见性。页面不从历史详情创建 draft，但把统一错误 envelope 的 `retryable` 理解为可重新提交原 operation。

两个单元都能声称遵守 AD-10（逐条状态和禁止历史自动重试）、AD-21（不另建内部状态机）、AD-23（只消费版本化 verifier 结果）、AD-26（安全只读投影并每次复核 actor/workspace）。然而 A 的 facts 无法提供 B 所假设的“产品关闭”事实，B 也无法仅从 A 的 Run 终态推导六个月例外；A 的技术重发和 B 的客户端重试又没有共同的允许调用语义。以下路径不需要任何一方故意违规即可触发。

## 不兼容路径一：同一核验结果产生不同终态

场景：provider 已接受写请求；期限内多次回读都得到与 expected 不一致的、结构有效的 observed；到 settle deadline 时仍不一致，但没有 provider 的明确拒绝证据。

1. 单元 A 可把“稳定 observed 证明目标未达成”解释为 definitive `failed`。AD-10 允许“能够证明 mutation 未完成”进入 `failed`，AD-23 又要求 outcome mapping 单独支持“明确失败”。若批次没有 `manual_attention`，A 按 AD-21 写 `Run=succeeded + outcome=none_succeeded/partial`。
2. 另一个同样合规的 verifier 实现可把“provider 接受后，deadline 结束仍不一致”解释为结果未知，写 `manual_attention`。AD-23 明确要求表示“预算结束仍未知／不一致”，AD-10 也允许期限结束仍不能证明目标时进入 `manual_attention`。该路径按 AD-21 写 `Run=failed + ActionResult=blocked`。
3. 两者的 observed、expected、deadline 和调用次数完全相同，只因“证明 mutation 未完成”与“预算结束仍不一致”的判定边界未定义，就形成相反 Run 终态。版本化 policy 只能保存分歧，不能替代全平台统一的分类约束；Catalog schema 若允许两种 mapping，产品层仍无法保证同类事实使用同一语义。
4. 即使 execution 已写 `failed`，AD-21 只固定 Run=`succeeded` 和 outcome，却没有固定该组合的 `ActionResult.status`。`product` 可以投影为 `completed + none_succeeded`，也可以因所有 item 均“待排查”投影为 `failed`。两种实现均没有新建状态机，但 HTTP 轮询终止值、SSE 最后一帧和 Web 卡片的顶层状态不兼容。

## 不兼容路径二：六个月与 unresolved discoverability 无法同时实现

场景：一条操作在七个月前执行，包含一个 definitive `failed` item，或包含一个 `manual_attention` item；用户从固定“操作记录”入口重新查找。

1. AD-21 将 `failed` 和 `manual_attention` 都映射为产品“待排查”，但 AD-26 的超六个月例外只枚举 `waiting/running/reconciling/manual_attention`，遗漏 definitive `failed`。因此 A 可按枚举把七个月前的 `failed` 记录排除；B 可按产品“待排查且未关闭”的语义继续展示。两边都逐字遵守自己负责的规则。
2. `manual_attention` 同时被 AD-21 定义为 `Run=failed + ActionResult=blocked` 的终结条件，又被 AD-26 写成“在终结前始终可发现”。A 可因 Run 已 terminal 而在六个月后排除，B 可因 `manual_attention` 被点名而无限期保留。Spine 没有规则决定“终结前”修饰 Run、Item、人工排查事项还是操作记录。
3. Spine 没有 `resolution_status/closed_at/resolved_by/resolution_reason` 一类产品关闭事实。`failed/manual_attention` 是执行终态，不等于人工排查已经关闭；仅凭现有 Product Facts，任何实现都无法证明 approved proposal 所说的“待排查对象在关闭前不受六个月限制”。
4. “操作时间”没有权威字段。A 使用 `confirmed_at`，B 使用 `last_transition_at`，第三种合理实现可使用 `mutation_started_at`。一条五个月前创建、七个月前确认不可能，但一条七个月前确认、五个月前进入 manual attention 的记录，会因时间锚点不同而分别落在窗口内外；“按操作时间倒序”也会产生不同游标和分页结果。
5. AD-26 没有固定默认窗口与超窗未关闭集合的合并、去重、排序和分页代数。HTTP 可以先取六个月分页再追加 unresolved，Web 也可以要求后端返回统一排序全集；前者会在页间重复或漏掉超窗记录，后者的 total/cursor 又无法与前者交换。

## 不兼容路径三：actor scope 只有名词，没有访问谓词

场景：原操作由当前 actor 完成，但此后该 actor 对其中部分漏洞的对象权限被收回；或 workspace 更换 FOBrain provider instance/凭据后，同一业务用户重新进入“操作记录”。

1. AD-26 的“每次读取重新校验当前 actor/workspace”没有说明只校验 `(record.actor_subject_id == current_actor && record.workspace_id == current_workspace)`，还是还要执行当前 `object_permission_predicate`。A 可以返回操作时已获授权的安全冻结事实；B 可以按 UX 源中的“当前用户仍有权限”重新过滤或整条拒绝。前者可能暴露已撤权对象，后者会改变历史批次的 count/outcome，二者都自称 actor/workspace-scoped。
2. 若访问时重新做对象权限校验，Spine 没有固定批次的投影规则：隐藏无权 item 后是否重算 `success_count/issue_count/outcome`，还是保留原总数并显示“部分内容不可见”。A 提供不可变 ActionResult，B 若重算就形成第二套结果；B 若不重算，汇总又可能泄漏已撤权对象数量。
3. AD-15 把 actor identity 固定为 `provider_instance_id + stable_user_business_id`。provider instance 更换会让同一业务用户成为新 actor；严格相等实现会使六个月内记录全部不可发现，做 instance alias/rebinding 的实现又违反字面 identity 公式或需要尚不存在的受审计迁移。当前规则没有选择 fail-closed 不可见、管理员迁移，还是稳定主体映射。
4. 共享凭据下多个真实操作者会得到同一个 `actor_subject_id`。AD-15 承认不能声称多人审计，但 AD-26 仍将 actor scope 描述为防止跨 actor 泄漏；`web` 若把它呈现为“我的操作记录”或“实际确认人”，会把凭据主体误当自然人。至少需要把首版 scope 明确为“凭据主体记录”，否则 httpapi/product 与 UX 会对 actor 文案和隐私边界采用不同含义。

## 不兼容路径四：重试禁令无法跨层机械执行

场景：一个批量 draft 有 8 个 succeeded、2 个 definitive failed；用户稍后重新查询相同对象并再次表达同类动作，或网络在 provider 已接收请求后中断。

1. AD-10 只给出 draft 级 `retry_of_draft_id`，没有 item 级 lineage。新 retry draft 若引用原批次，execution 无法仅靠该字段证明它只包含两个 `failed` item；包含 succeeded/manual_attention 的混合 retry draft、重复 item 或改变 action type 都没有明确拒绝规则。
2. AD-11 的 mutation key 包含 `draft_id/version`。用户重新查询后创建普通新 draft，即使业务 `item_identity + action_type + target parameters` 与既有 succeeded/manual_attention 完全相同，也天然得到新 key。AD-10 的“未知结果和成功项不得重复 mutation”因而不能由所声明的幂等键执行；A 可以维护额外跨 draft operation index 并阻断，另一个实现可以把新人工确认视为新操作并放行。
3. “历史详情只读且不自动创建 retry”只约束 AD-26 的 UI 入口，不约束聊天输入或 Action API。Web 不显示按钮并不证明 `POST /agent/actions` 拒绝从历史安全 ref、相同业务对象或复制参数创建新 draft；httpapi 也没有稳定错误码区分 `retry_forbidden_unknown`、`retry_requires_new_decision` 和普通 conflict。
4. AD-11 允许“进程恢复与技术重试复用同一 key”，AD-10 又禁止未知结果重复 mutation。provider 调用在请求发出后超时，A 可以基于 provider 声明的幂等能力用同 key 重发；另一实现可以认为任何 post-send 重发都属于盲重试，只允许 readback。两者都能引用 AD 文本。必须固定 `reserved/not_sent/sent/accepted/unknown` 的 durable 边界及每一状态允许的外部调用，而不是只规定 key。
5. 全局错误约定仍暴露一个无命名空间的 `retryable`。httpapi 可把“重试同一查询/SSE 连接”标成 true，Web 可把它解释为重新提交原 action；即使历史卡没有重试按钮，通用错误组件或客户端 SDK 仍可能重发 mutation。传输契约需要区分 `transport_retryable`、`verification_resumable` 和 `mutation_retry_allowed=false`，写操作未知结果不得只靠中文文案约束。

## Findings

- `failed` 与 `manual_attention` 共用“待排查”产品标签，却只有后者进入 AD-26 的超六个月枚举，无法得到统一 discoverability。
- `manual_attention` 在 AD-21 中终结 Run，在 AD-26 中又只保证“终结前”可发现，终结对象不明确。
- 执行终态与人工排查关闭态被混为一体；缺少权威的 unresolved/closed 事实与关闭审计。
- “操作时间”未绑定一个持久化字段，窗口、排序、游标会因 confirmed/started/updated 时间选择不同而分叉。
- 六个月集合与超窗 unresolved 集合没有统一查询、排序、分页和去重代数。
- deadline 后的 observed mismatch 既可被解释为 definitive failed，也可被解释为 manual attention，AD-10 与 AD-23 没有可测试判据。
- `Run=succeeded + outcome=none_succeeded/partial` 对应的 ActionResult.status 未固定，product/httpapi/web 仍可产生竞争顶层终态。
- 同一“待排查”标签抹平了可重新人工决策的 definitive failed 与永久禁止重复 mutation 的 manual attention，前端不能从文案推导安全动作。
- actor/workspace scope 没有固定访问时对象权限谓词，冻结历史的完整性与撤权后的最小披露无法同时证明。
- 访问时隐藏 item 后是否重算批次汇总未定义，可能形成第二套 ActionResult 或通过数量泄漏不可见对象。
- provider instance 参与 actor 主键但没有凭据轮换/instance 替换策略，六个月发现承诺可能因身份技术标识变化失效。
- 共享凭据主体与自然人操作者的语义没有在操作记录 API/文案中区分，actor scope 容易被过度声明。
- `retry_of_draft_id` 粒度不足以证明只重试 definitive failed items，也不能表达 item 级 lineage 和禁止集合。
- draft-scoped mutation key 无法阻止跨 draft 对 succeeded/manual_attention 的语义重复写入。
- “技术重试复用同 key”与“未知结果不得重复 mutation”缺少 durable sent 边界和允许调用矩阵，恢复器可采取相反行为。
- 通用 `retryable` 没有区分读取重连、继续核验和 mutation 重发，HTTP SDK/Web 通用重试可能绕过产品禁令。

## 最小收口要求

1. 在 AD-21/23 的 v2 fixture 中给出穷举 outcome 表：输入至少包含 provider acceptance、observed relation、evidence conclusiveness、deadline/budget、lease certainty；输出唯一固定 Item、Run、ActionResult、product status、stream terminality 和 `mutation_retry_allowed`。
2. 将操作记录的“执行终态”和“人工事项关闭态”分离，增加可审计的 `attention_status/closed_at`（或等价事实），明确 `failed/manual_attention` 何时属于 unresolved，并以一个持久化 `operation_at` 定义六个月窗口。
3. 固定 ListOperations port/OpenAPI：actor equality、访问时对象权限、字段遮蔽、汇总是否重算、六个月与 unresolved 的 union、排序、cursor 和 total 必须由后端单一投影完成，Web 不二次推导。
4. 定义 item 级 retry lineage 与跨 draft semantic mutation guard；Confirm 必须拒绝包含 succeeded、reconciling、manual_attention 或 lineage 不完整 item 的 retry draft。普通新 draft 命中既有受禁操作时也必须得到稳定拒绝，而不是靠 UI 隐藏入口。
5. 用持久化 attempt 阶段固定重发矩阵，并在 OpenAPI 错误契约中拆分读取重试、核验续跑和 mutation 重发；任何 post-send unknown 默认只允许 Verify/Reconcile，除非架构明确批准并证明 provider 端幂等重放语义。

在上述接缝进入 schema、fixture、OpenAPI、generated contract 和跨层 contract tests 前，Story 1.24 与 Story 1.36 不应被判定为可独立验收完成，真实写域门禁也不应因“已有 AD-10/21/23/26”而放行。

---

## 中间 Recheck — 2026-07-15（已由文末 Final Targeted Recheck 取代）

### Recheck Verdict

**NOT IMPLEMENTATION-SAFE — CHANGES REQUIRED。** 初审指出的大部分产品级分歧已经关闭：AD-10 现有穷举 outcome 表；AD-21 固定了 Run/ActionResult/outcome；AD-23 固定了 `reserved_not_sent → sent`、持久 deadline 和预算先扣减；AD-11 增加 item-level retry lineage、跨 draft semantic key 和传输层三种 retry 标志；AD-26 固定 actor/workspace 精确相等、访问时整批权限证明、六个月与 open lifecycle 的后端 union/cursor；AD-27 固定 opaque ref 与 ProviderLocator 分离。这些规则也已进入迁移 ADR、SPEC companions、Epics/AC 和 `run-lifecycle.md`，不再只是 Architecture 的孤立声明。

但当前机器契约仍留下四组可导致独立实现不兼容或重复 mutation 的阻断接缝：semantic key 没有规范化字节契约和并发 claim；ActionAttempt 的状态依赖字段被写成全阶段必填；operation history 的 immutable `operation_at`、日历窗口与 cursor snapshot 仍可产生两种实现；opaque ref、内部稳定业务身份和 ProviderLocator 之间缺少唯一桥接/所有权契约。唯一 outcome 的表虽然存在，`VerificationEvidence` 本身仍没有可编译 discriminant，`manual_attention` 的 aggregate outcome 在 SPEC state table 中也没有落成枚举值。

因此，当前文档足以证明“纠偏方向已传播”，还不足以证明 Story 1.3/1.4、1.9/1.10、1.22/1.24、1.27～1.29 和 1.36 由不同实现者完成后能够无设计补充地拼装。真实 mutation 继续受既有 gate 阻断；上述契约缺口应在 Sprint Planning 前关闭，而不是留给 Dev Story 静默选择。

### 复核覆盖矩阵

| 复核点 | 当前结论 | 已传播证据 | 剩余接缝 |
| --- | --- | --- | --- |
| 唯一 outcome | 部分关闭 | Architecture AD-10/21、迁移 ADR §4/§7、SPEC `state-machines.md`、AC-30、Epics AR-46/Story 1.24、`run-lifecycle.md` | VerificationEvidence 没有字段/discriminant；SPEC 对 manual-attention aggregate outcome 仍非枚举 |
| sent 原子边界、deadline、预算 | 方向关闭、领域契约未闭 | Architecture AD-11/23、迁移 ADR §6/§7、SPEC 崩溃图/AC-29、Epics AR-45/Story 1.22～1.24、lifecycle 68～70 | pre-sent 状态却要求 sent-only 字段；状态图仍使用旧 `reserved` |
| 跨 draft semantic guard | 未关闭 | Architecture AD-11、迁移 ADR §6、SPEC domain/AC-31、Epics AR-15/Story 1.22、lifecycle 63 | 缺规范化字节算法、并发 claim、活动状态覆盖和 definitive-failed 普通草案拒绝规则 |
| actor/permission | 基本关闭 | Architecture AD-26、迁移 ADR §9、SPEC OperationRecord/AC-28/33、Epics AR-47/Story 1.36 | `history_access_policy_ref` 未成为任何聚合的必填冻结字段，也没有解析失败的稳定传输状态 |
| calendar union/cursor/open attention | 未完全关闭 | Architecture AD-21/26、ADR §9、SPEC operation projection/AC-33、Epics AR-35/47/Story 1.36 | waiting→decision 时 operation_at 是否变化、自然月算法、server as_of、cursor tamper/snapshot 一致性未固定 |
| opaque ref/locator | 未完全关闭 | Architecture AD-27/G-READ-03、ADR §10、SPEC Domain/AC-32、Epics AR-44/Story 1.25/1.27～1.29 | stable semantic identity 与 opaque ref/locator 的桥接缺失；ProviderLocator 不在所有权表/port/retention 契约中 |
| 进入 SPEC/Epics/ADR/lifecycle | 已进入但不完全同构 | 四类文档均可检索到相应规则和验收 | 下述名称、字段必填性和 AC 强度仍分叉 |

### 当前仍可构造的两个合规实现单元

- **单元 A（`execution + facts + store/sqlite`）** 使用 Go struct 字段顺序和普通 JSON marshal 生成 `normalized_target` 字节；Prepare 时创建带相同 semantic key 的多个 `prepared` item，Confirm 只检查既有 `succeeded/reconciling/manual_attention`，直到 attempt 进入 `sent` 才认为 guard 被占用。它为 `reserved_not_sent` attempt 写零值 `sent_at/deadline` 以满足“必填”，waiting 记录永久使用 `published_at` 作为 immutable `operation_at`，cursor 只保存文档点名的 `as_of/cutoff/last keys`。

- **单元 B（`capabilities + product + httpapi + web`）** 使用 JCS、业务字段排序和 UTC 时间规范化得到 target digest；它假设 semantic key 在 Confirm 事务中先取得唯一 claim，任何 `prepared/reserved_not_sent/sent/verifying/reconciling/succeeded/manual_attention` 均阻止第二个普通 draft。它把 approval decision 当最终 `operation_at`，假设 server 返回不可篡改、绑定查询高水位的 opaque cursor，并假设 `entity_ref` 可跨查询稳定代表同一业务对象。

两者都可以引用最新文档证明自身合规：A 严格采用已写出的 hash 公式、被点名的阻断终态和 cursor 字段；B 严格执行“跨 draft guard”“operation_at 取 approval decision”“稳定分页”和 opaque ref。组合时，A/B 对同一 target 得到不同 semantic key；两个同时确认的草案可在 A 中都进入首次发送；B 无法接受 A 的零值 sent-only 字段；历史记录会在审批后移动排序位置或不移动；B 也无法证明 A 的 per-snapshot opaque ref 与跨 draft business identity 是同一对象。

### 剩余阻断路径

#### 路径一：两个未发送草案绕过 semantic mutation guard

1. 两个请求基于同一业务对象、动作和目标并发 Prepare，分别生成 draft A/B；两个 ActionItem 均为 `prepared`，semantic key 相同。
2. AD-11 只明确普通新 draft 命中 `succeeded/reconciling/manual_attention` 时拒绝，未枚举 `prepared/reserved_not_sent/sent/verifying`；domain contract 又只说命中“未解除”的 key，不定义 claim 状态。
3. 两个 Confirm 事务可依次看到“没有被列为禁止的 terminal key”，各自预留 `reserved_not_sent` attempt。随后两者都合法执行一次 `reserved_not_sent→sent` 和首次调用。
4. 单实例 SQLite 与 draft-scoped mutation key 都不能自动阻止该路径；必须有事务内唯一的 semantic claim/reservation，而不是事后检查 ActionItem 状态。
5. 同一缺口还允许普通新 draft 命中 definitive `failed` 后绕过 `retry_of_action_item_id`：Architecture 的普通 draft 拒绝枚举遗漏 `failed`，而“只有 failed 可 retry”没有机械规定普通 draft 必须被识别并转换为 retry lineage。

#### 路径二：attempt schema 无法同时表达 sent 前后

1. SPEC `domain-contracts.md` 把 `sent_at`、`settle_deadline`、`next_verification_at`、provider 接受摘要和 evidence ref 全部列为 ActionAttempt 必填。
2. 同一段又要求 attempt 先持久化为 `reserved_not_sent`，直到首次外部调用前的后续事务才初始化 `sent_at/deadline`；provider 摘要和 evidence 更不可能在预留时存在。
3. 实现 A 可用空字符串、Unix epoch 或零摘要满足 schema；实现 B 可把字段做 nullable/状态判别联合。两者的恢复扫描、deadline 比较和 generated contracts 不兼容，零值还可能使 pre-sent attempt 被误判预算耗尽。
4. `state-machines.md` 的权威状态图仍写 `prepared_item→reserved→sent`，而其崩溃图、Architecture、ADR、Epics 与 lifecycle 使用 `reserved_not_sent`；schema 作者可以合理选择两个不同 enum。

#### 路径三：operation history 的游标不是查询快照

1. 文档同时规定 `operation_at` 为 approval decision 时间、waiting 时使用 Pending `published_at`，并称写入后不再变化。waiting 后获批时，到底保持 published_at，还是改成 decision_at，没有唯一答案；改值会让已出现在第一页的记录在下一页前移动。
2. `as_of` 只固定 cutoff，没有明确必须由 server 生成、是否允许客户端提供，也没有要求 cursor opaque/签名；实现可接受客户端改写 as_of 或 last keys，另一实现可拒绝，HTTP 契约不兼容。
3. cursor 仅固定 `as_of/cutoff/last keys`，没有 Product Facts sequence、数据库 snapshot 或“按 as_of 读取 lifecycle/attention 状态”。七个月前的 Run 在翻页期间由 running 变为 terminal+open attention，union membership 发生变化但 operation_at 不变，keyset 分页可能漏项或重复。
4. Architecture/SPEC 使用“Asia/Shanghai 日历时间减六个月”，Epic AC 使用“前六个 UTC+8 自然月”。例如月末和当月中旬可分别被实现为滚动 `AddDate(0,-6,0)` 或当前月加前五个完整自然月；只有后端计算并不能消除契约语义分歧。

#### 路径四：opaque ref 与 semantic identity 没有共同内部主键

1. AD-27 要 Product Facts 只保留 opaque `entity_ref`，ProviderLocator 独立保存；semantic key 却使用 `item_business_identity`，QuerySnapshot 又要求“稳定业务身份安全值”。三者的等价关系没有契约。
2. 实现 A 可每次查询生成随机 opaque ref，并用 locator 中真实 provider id 计算 semantic key；实现 B 可把稳定 HMAC ref 同时当 business identity。二者都不泄漏 provider id，但跨 draft guard、引用去重和 locator lookup 不能交换。
3. `ProviderLocator` 未进入 SPEC 所有权表，也没有命名 repository/resolver port、mapping key、唯一性、版本迁移或与 snapshot/open attention 的 retention 关系；`facts repository` 和 `capabilities resolver` 可各自认为另一个包拥有生命周期。
4. `history_access_policy_ref` 存在相同问题：OperationRecord 查询要求使用“冻结 ref”，但 ActionDraft、ActionResult、AuditEvent 的必填字段都未承载它，Epics Story 1.36 也只验整批 fail closed，未验证 ref 的来源、版本和损坏恢复。

### Recheck Findings

- `VerificationEvidence` 只有所有权和名称，没有可编译字段、evidence conclusiveness discriminant、expected/observed digest、policy/version、调用序号或判定来源，两个 verifier 仍可用不同证据证明 `failed`。
- SPEC aggregate table 对 `manual_attention` 的 Outcome 写“成功项仍保留”，没有采用 Architecture 已固定的 `partial/none_succeeded` 枚举，ActionResult schema 仍可分叉。
- ActionAttempt 把只可能在 sent/verification 后产生的字段列为全阶段必填，与 `reserved_not_sent` 的先持久化要求直接矛盾。
- SPEC item 状态图仍使用 `reserved`，其余目标文档使用 `reserved_not_sent`，封闭 enum 无法同时忠实实现。
- `semantic_mutation_key` 没有版本化 canonical payload、JCS/字段边界、Unicode、时间、空值、集合排序或 target normalization 规则，同一语义可得到不同 hash。
- semantic guard 没有事务内 claim 聚合、唯一约束和 owner/lifecycle；并发 prepared/confirmed drafts 可以在任一 item sent 前同时通过。
- 普通新 draft 对 definitive `failed` 的处理未固定为“必须拒绝并要求 item-level retry lineage”，可绕过 retry 专用路径。
- guard 的阻断状态遗漏 `prepared/reserved_not_sent/sent/verifying`，仅阻断成功、reconciling 和 manual attention 不能实现跨 draft at-most-once。
- waiting 记录的 immutable `operation_at=published_at` 与批准后 `operation_at=decision_at` 互相冲突，排序键可能在生命周期中改变。
- 六个月“日历减法”和“前六个自然月”未给出同一个边界算法，月末和月中 cutoff 可不同。
- `as_of` 的服务端所有权、时钟精度、cursor 编码/签名和非法 cursor 错误没有固定，httpapi 与 Web 客户端可采用不兼容协议。
- cursor 没有绑定事实高水位或 lifecycle-as-of snapshot，open attention union 在多页读取中变化时不保证无重复、无遗漏。
- `history_access_policy_ref` 未成为被命名聚合的必填冻结字段，访问时无法稳定选择原 policy/version 或区分损坏与当前 denied。
- `attention_status` 允许 `closed`，同时首版声明没有 closure command；AC-33 又要求“已关闭终态”，当前 writer 是否可产生 closed 没有唯一答案。
- opaque `entity_ref`、安全稳定 business identity、semantic mutation identity 和 ProviderLocator mapping 没有共同内部主键/等价规则。
- ProviderLocator 虽被指定“由 facts repository 安全持久化”，却未进入领域所有权表或命名 port，也没有引用计数、过期、open attention 保留和版本迁移约束。
- Story 1.22 的 AC 只断言“阻止已 sent 语义”，未覆盖并发 prepared/reserved/verifying claim 或普通草案绕过 failed lineage；SPEC AC-31 的更强要求没有完整下沉到 Story AC。
- Story 1.36 的 AC 没有固定 inclusive cutoff、server-issued as_of、cursor snapshot、高水位和权限策略 ref 损坏路径，开发 Story 可在弱 AC 下提前完成。

### 最小收口要求

1. 把 ActionAttempt 写成按 phase 判别的领域/schema union：`reserved_not_sent` 禁止 sent-only 字段；`sent+` 必填 immutable `sent_at/deadline/policy/epoch/budget`；verification 后才允许 evidence/next poll。统一删除旧 `reserved` 名称。
2. 定义 `VerificationEvidence` schema 和唯一 reducer；补齐 manual-attention aggregate outcome。evidence 必须携带 policy/verifier version、expected/observed digest、relation、conclusive-not-applied 判据、调用预算序号和判定时间。
3. 新增事务性 `SemanticMutationClaim`（或等价唯一索引/聚合）：规定 canonical payload/version、claim 获取时点、所有阻断阶段、item owner、failed retry 交接和保留/释放条件；AC 必须覆盖两个并发普通 draft 在 sent 前竞争。
4. 将 operation history 定义为 snapshot-consistent query：选择一个永不变化的 `operation_at`，给出精确自然月算法，由 server 生成 `as_of`，cursor opaque/tamper-evident 并绑定 Product Facts 高水位或等价 snapshot；明确 lifecycle/attention 按该 snapshot 求值。
5. 增加内部稳定 `entity_subject_key`（或等价值）把 opaque ref、跨查询 identity、semantic key 和 ProviderLocator 绑定；在领域所有权表中命名 LocatorRepository/Resolver port，并固定 `history_access_policy_ref` 的存放聚合、版本、损坏语义和 retention。
6. 将以上收口同步到 Epics Story 1.22、1.24、1.27～1.29、1.36 的 AC 和 SPEC acceptance matrix；只在 Architecture/ADR 中补词不足以让独立 Story 自动继承。

完成这些修订后，可再次进行 targeted recheck。当前已经关闭的初审问题无需回退：三态产品文案、Run/ActionResult 顶层映射、open attention、actor 精确范围、整批权限 fail closed、post-sent 只读核验、retry flag 拆分和 locator 不外露应继续保留。

---

## Final Targeted Recheck — 2026-07-15（已由文末 Latest Targeted Recheck 取代）

### Final Verdict

**NOT IMPLEMENTATION-SAFE — 仍有 3 组精确阻断。** 上轮六项最小收口要求已经逐项落地，且不是只补在 Architecture：phase-discriminated ActionAttempt、VerificationEvidence 字段、`smk.v1` canonical payload、事务唯一 SemanticMutationClaim、immutable operation history/cursor、entity subject/locator/history policy 均已进入 SPEC、AC-34～39、Epics Story AC、迁移 ADR 和 lifecycle。两个独立实现者现在会生成相同 attempt union、相同 semantic claim 冲突、相同 calendar cutoff、相同 actor/permission fail-closed 和相同 locator 解析链。

最终攻击仍发现三个会改变权威状态或安全去重结果的接缝：唯一 reducer 没有消费自己契约中的 `provider_acceptance` 与失租约事实；跨会话操作记录使用了只在单 Run 内定义的 `snapshot_sequence`，无法证明跨 Run union 的分页快照；`entity_subject_key` 依赖版本化 workspace secret，却没有固定 secret 轮换期间的稳定主体/claim 继承。另有一个 evidence 字段条件性问题与第一项绑定：明确拒绝可形成 conclusive not-applied，但所有 evidence 又无条件要求 observed StructuredResult ref/digest。

这些不是“待实现 schema 会自然补齐”的普通细节。它们分别允许同一 evidence 在 Architecture reducer 与 SPEC reducer 得到不同结果、允许操作记录翻页漏项，以及允许 secret 轮换后同一业务对象获得新 semantic key 从而绕过旧 claim。因此当前不能给出 `IMPLEMENTATION_SAFE`。

### 上轮收口逐项复核

- ActionAttempt 已在 `domain-contracts.md` 成为按 phase 判别的封闭 union；`reserved_not_sent` 禁止 sent-only 字段，旧 `reserved` 已从目标状态图删除，AC-34 明确拒绝零值伪造。
- VerificationEvidence 已具备 policy/verifier version、digest、call index、relation、conclusiveness、decision/time；manual_attention aggregate 固定为 `partial|none_succeeded`，AC-35 覆盖唯一 reducer。
- semantic key 已固定 `smk.v1 + RFC 8785 JCS`、NFC、UTC 纳秒、missing/null 与 set 排序；SemanticMutationClaim 以数据库唯一键在 Confirm 事务内先 claim 后预留，所有已有状态阻断普通草案，failed retry 只能事务转移 lineage，AC-31/36 覆盖并发和 canonicalization。
- operation history 已固定首次发布时 immutable `operation_at`、独立 `decision_at`、服务端 as_of、Asia/Shanghai 月末 clamp、inclusive cutoff、opaque tamper-evident cursor 与 filter/scope 绑定，AC-33/37 覆盖篡改和跨页变化。
- actor/workspace 精确相等、frozen history policy 原样传播、损坏与 denied 的不同稳定错误、整批 fail closed 且不重算均进入 Domain Contract 与 AC-39。
- `entity_ref → EntityRefMapping → entity_subject_key → ProviderLocator` 已成为唯一解析链；LocatorRepository/Resolver 所有权、引用保留、additive locator version 和跨查询去重进入 Domain Contract 与 AC-38。
- Story 1.22、1.24、1.27～1.29、1.36 已获得对应强 AC；此前“SPEC 强、Story 弱”的落差已关闭。
- post-sent 禁止 mutation、预算先事务扣减、deadline/epoch 不重置及 retry flag 拆分仍保持一致，没有因本轮补强回退。
- definitive failed 与 manual_attention 均产生 open attention，首版 writer 不接受 closed；六个月窗口与 open lifecycle union 的产品语义一致。
- opaque refs、entity subject、locator 与 history policy 均明确禁止进入 Workbench/Action API/LLM/raw audit 出口，安全投影边界一致。

### 剩余阻断一：唯一 Reducer 与 AD-10 输入不等价

Architecture AD-10 规定成功条件是 **provider accepted 且 expected=observed**；`domain-contracts.md` 的 VerificationEvidence 也持久化 `provider_acceptance=accepted|rejected|unknown`。但紧接着定义的唯一 reducer 第一条只是 `relation=match → succeeded`，完全没有检查 provider acceptance。

两个实现者可对同一 evidence 合法地产生相反结果：

- evidence=`provider_acceptance=rejected, relation=match`。实现 A 按 AD-10 不允许 succeeded；实现 B 按 SPEC 唯一 reducer 直接 succeeded。
- evidence=`provider_acceptance=unknown, relation=match`。A 继续按未知写入处理，B 标记完成；这会改变 Run/ActionResult、attention、semantic claim 和 mutation retry 结论。

失租约路径也无法由当前 reducer 表达。AD-10、ADR、lifecycle 与 AC-30 都要求“失租约仍未知”直接进入 manual_attention；VerificationEvidence 没有 `lease_certainty/control_certainty` 字段，reducer 又只检查 relation、deadline 和预算。若失租约时预算尚未耗尽，实现 A 按 AD-10 进入 manual_attention，实现 B 按 reducer 继续 reconciling。

此外，Evidence 把 `observed StructuredResult ref/digest` 列为所有分支必填，但 `provider rejected + conclusive_not_applied` 可能在没有一次成功 readback candidate 时已经具备明确未应用证明。一个实现会伪造 synthetic observed result，另一个会让字段 nullable，schema 不可交换。

**精确收口：** reducer 输入与顺序必须固定为：control certainty/lease lost 优先；`provider_acceptance=accepted && relation=match` 才 succeeded；policy 允许的 conclusive-not-applied 才 failed；其余按 deadline/budget reconciling/manual_attention。Evidence 增加 lease/control certainty；observed ref/digest 改为按 evidence kind 判别的条件必填，或明确所有 conclusive rejection 必须先取得并持久化真实只读 observed candidate。同步 AC-35 增加 rejected+match、unknown+match、lease-lost-before-budget-exhaustion 和 rejection-without-observed 四个 fixture。

### 剩余阻断二：Operation History 的高水位没有跨 Run 定义

AD-24、Epics AR-19 和既有 SSE 契约只定义 **Run 内** 单调 FactEvent sequence。操作记录却跨 conversation、跨 Run 查询，AD-26/Domain Contract/AC-37 使用一个 Product Facts `snapshot_sequence` 冻结整个 union，但没有定义它是 workspace-global sequence、operation-index sequence、SQLite snapshot id，还是一组 `(run_id, sequence)`。

只使用当前 Run-local sequence 无法冻结以下路径：第一页读取后，另一个旧 Run 从 running 变为 definitive failed 并创建 open attention。其 `operation_at` 落在 cutoff 外但现在因 open attention 加入 union；它的 Run-local sequence 与第一页 cursor 中的单个高水位没有可比较关系。实现 A 可以把数据库当前状态重新查询而插入/遗漏该项；实现 B 可以物化第一页 id 集合；两者都能声称“按 snapshot_sequence 重建”。

而本项目明确采用“当前状态 + FactEvent，不承诺完整事件溯源”。即使补一个 workspace-global counter，如果 OperationRecord index 没有保留 valid-from/valid-to 或物化 membership，也无法从已变化的当前行重建旧高水位的 lifecycle/attention 状态。

**精确收口：** 为跨 Run 历史投影单独定义一种可恢复快照机制，例如：在操作索引每次变化时同事务取得 workspace-scoped `operation_index_sequence` 并保留版本行，cursor 绑定该 sequence；或首个请求物化 actor/filter 对应的 result-id snapshot 并让后续页引用。必须明确 `snapshot_sequence` 不是 Run SSE sequence，定义 owner、事务写入点、retention/invalid_cursor 条件，并让 AC-37 同时变化两个不同 Run 验证无重复遗漏。

### 剩余阻断三：Entity Subject 在 workspace secret 轮换后不稳定

`entity_subject_key` 由 `HMAC-SHA-256(workspace_secret, canonical provider identity)` 生成，文档同时要求 workspace secret 受控且版本化，却没有规定 secret version 轮换时如何复用旧 subject。若实现 A 用当前 secret 为新查询重算 key，同一 provider 对象在轮换后得到新 `entity_subject_key`、新 semantic key 和新 claim namespace；既有 succeeded/manual_attention claim 不再阻断它。实现 B 可永久 pin 首个 secret 或遍历历史 secret 复用 mapping。两者都满足“secret versioned”，但安全行为相反。

这也影响 LocatorRepository：新 subject 下无法自然找到旧 locator/version，可能重复建立 mapping 或使 open attention 的详情引用与新查询去重失配。

**精确收口：** 固定主体 key 的密钥生命周期。最小方案是首版声明每个 workspace 的 subject derivation key version 在所有引用/claim/audit retention 结束前不可轮换；更完整方案是在 `EntitySubject` 持久化 immutable subject id 与 derivation key version，解析 canonical provider identity 时对仍保留的 key versions 做受控匹配，轮换只改变新 locator 的加密密钥而不改变 subject id。AC-38 增加 workspace subject secret 轮换前后同对象仍命中同一 subject、semantic claim 和 locator lineage 的 fixture。

### Final Findings

- AD-10 要求 accepted+match 才成功，SPEC reducer 只检查 match，rejected+match 可产生竞争 outcome。
- unknown+match 在当前 reducer 中也会成功，违反未知写入不得伪报完成的控制面规则。
- VerificationEvidence 未携带 lease/control certainty，失租约且预算未耗尽时无法唯一进入 manual_attention。
- reducer 没有声明失租约与 match/conclusive/deadline 的优先顺序，执行层可绕开所谓唯一 reducer。
- conclusive provider rejection 可能没有 observed StructuredResult，但 Evidence 把 observed ref/digest 无条件列为必填。
- AC-35 没有覆盖 rejected+match、unknown+match、预算内失租约和无 observed 的明确拒绝。
- 操作记录跨 Run，但被绑定的 `snapshot_sequence` 仅在现有契约中定义为 Run-local sequence。
- 跨 Run lifecycle/open-attention 变化无法与单个 Run sequence 比较，keyset 分页仍可能漏项或重复。
- 当前状态 + 非完整 event sourcing 无法仅凭一个新全局数字重建旧 operation union，必须版本化索引或物化查询成员。
- AC-37 只说“新增事实/attention 变化”，没有明确两个不同 Run 的 sequence 不可比较攻击。
- workspace secret 被声明版本化，但 entity subject key 没有 derivation key version/轮换继承规则。
- secret 轮换可改变 entity_subject_key 和 semantic key，绕过既有 succeeded/manual_attention claim。
- secret 轮换还可能切断旧 ProviderLocator 与新查询主体的映射，破坏跨查询去重。
- AC-38 没有覆盖 subject derivation secret 轮换前后身份稳定性。

### Final Gate

在以下三项进入 Architecture、SPEC Domain/State、迁移 ADR、lifecycle 和对应 Story/AC 前，维持 `NOT IMPLEMENTATION-SAFE`：

1. reducer 纳入 provider acceptance、lease/control certainty与 evidence-kind 条件字段；
2. 定义跨 Run operation-history snapshot owner/sequence/版本化重建机制；
3. 固定 entity subject derivation secret 的轮换与 claim/locator 继承。

除这三项外，上轮列出的六组接缝已经关闭，不应回滚或重新设计。

---

## Latest Targeted Recheck — 2026-07-15

### Latest Verdict

**NOT IMPLEMENTATION-SAFE — 仅剩 2 个精确接缝。** 上次列出的三组阻断已完整关闭：判别式 Evidence reducer 现在消费 provider acceptance 与 control certainty；Operation History 使用 facts-owned immutable `OperationListSnapshot`，不再借用 Run-local sequence；随机 immutable EntitySubject 配合 active/retained fingerprint alias 和独立 locator encryption key，密钥轮换不再改变 semantic claim namespace。`cancelled_before_send` 的 sent/cancel draft 级全有或全无 CAS 也已横向进入 Architecture、SPEC、ADR、lifecycle、Epics Story 1.22/1.23 与 AC-29/34。

本轮没有再发现 outcome、deadline/预算、跨 Run 历史分页、actor/permission、opaque locator 或 subject secret rotation 的竞争实现。剩余两个问题都比上一轮窄，但其中一个仍会改变跨 draft retry 安全语义：**从 definitive failed 转移来的 retry claim 若在 sent 前取消，当前统一“释放为 unclaimed”的规则会丢失原 failed claim 的 retry-only 门禁。** 另一个是 canonical acceptance matrix 的旧术语：AC-33 仍要求固定 `snapshot_sequence`，与 AC-37 和全部主规则的 `OperationListSnapshot` 明确分离相冲突。

只要这两点修正，即可给 `IMPLEMENTATION_SAFE`；无需再重开已经关闭的七组主设计。

### 已关闭证据

- `VerificationEvidence` 已按 `readback_observation|conclusive_non_application|control_loss` 形成封闭 union，observed ref/digest 只在真实 readback 分支必填。
- reducer 固定 control lost 最高优先级，且只有 held+accepted+conclusive real-readback match 成功；rejected/unknown+match 不再可能成功。
- AC-35 已包含 rejected+match、unknown+match、无 observed 的明确拒绝和预算未耗尽 control lost 攻击 fixture。
- `OperationListSnapshot` 已由 facts 拥有、product 通过 repository port 在单个 SQLite 事务中物化完整有序安全行和 total。
- `operation_list_snapshot_ref` 是独立随机 opaque ref；cursor 明确不复用 Run-local SSE sequence，并覆盖 TTL、篡改、scope/filter 与跨两个 Run 变化。
- `entity_subject_key` 已固定为 CSPRNG 256-bit 随机 immutable key，不再由可轮换 secret 直接推导。
- fingerprint 已固定 `esf.v1.<key_version> + HMAC-SHA-256`，并有 `(workspace,key_version,fingerprint)` 唯一约束。
- active/retained aliases 指向不同 subject 时稳定返回 `entity_identity_conflict` 并 fail closed；AC-38 已覆盖。
- subject derivation key 退役受 alias migration 与全部引用 retention 门禁；locator encryption key 可独立轮换而不改变 subject/claim/lineage。
- `reserved_not_sent→sent` 与 `reserved_not_sent→cancelled_before_send` 使用互斥 phase CAS，批次 cancel/stop 全有或全无。
- 任一 sent+ item 存在时 cancel/stop 整体 `action_already_sent` 且不改 item/claim；取消获胜时证明 mutation_count=0。
- Story 1.22、1.23、1.24、1.27～1.29、1.36 与 AC-29/34～39 均已承接相应强约束。

### 剩余阻断一：取消 retry 会丢失 definitive-failed 门禁

当前 claim 状态机只有：

```text
failed --人工 retry 事务转移--> reserved_not_sent
reserved_not_sent --cancel/stop 零写入事务--> unclaimed
```

普通新 draft 命中任何 active/retained claim 时被拒绝，只有逐条绑定 definitive failed lineage 的 retry 才可转移；但 cancel 规则又无条件“释放 active claim”，状态图明确回到 `unclaimed`。于是可构造：

1. 原 item A definitive failed，claim owner=A/status=failed。
2. 人工 retry item B 合法转移同一 claim，进入 `reserved_not_sent`。
3. 用户在 sent 前取消 B；CAS 证明零 mutation 并把 claim 释放为 unclaimed，仅保留 audit tombstone/lineage。
4. 普通新 draft C 命中相同 semantic key。实现 X 按 `unclaimed` 允许 C，无需 `retry_of_action_item_id`；实现 Y 读取 audit tombstone，恢复 A 的 failed 门禁并要求 retry lineage。

两者都能引用当前规则：X 遵守“取消后释放 claim”，Y 遵守“definitive failed 只能 retry 且 retained claim 阻断普通草案”。更糟的是，如果 tombstone 被解释为 retained claim但状态不是 `failed`，未来合法 retry 也无法按“只能从 failed 转移”继续，形成永久阻断。

**精确收口：** claim cancel 必须区分来源：

- 普通首次 claim 在 `cancelled_before_send` 后可进入 `released_zero_write`，允许未来普通新决策；
- 从 definitive failed 转移来的 retry claim 在取消事务中必须原子恢复 predecessor owner/status=`failed`，或保留 root failed claim、只取消其 child execution reservation；不得变成自由 unclaimed；
- audit tombstone 不能代替权威 claim 状态。

AC-29/31 增加“failed → retry reserved → cancelled_before_send → 普通 draft 拒绝 → 新 lineage retry 仍可转移”的 fixture；Story 1.22 AC 同步明确 claim 恢复语义。

### 剩余阻断二：AC-33 仍使用旧 snapshot_sequence

所有最新主规则和 AC-37 都明确 Operation History 使用 `operation_list_snapshot_ref`，且它“不是 Run-local SSE snapshot_sequence”。但同一 canonical `acceptance-matrix.md` 的 AC-33 仍写：

> 在固定 snapshot_sequence 上做窗口与 lifecycle-open union

实现 A 可据 AC-33 为操作记录增加或复用 sequence；实现 B 可据 AC-37 只实现物化列表 ref。即使产品行为最终相同，schema、cursor payload、fixture 断言和 generated contract 仍不兼容。

**精确收口：** 将 AC-33 的 `snapshot_sequence` 替换为“同一 immutable OperationListSnapshot / operation_list_snapshot_ref”，并明确 AC-33 验 calendar union，AC-37 验跨 Run/TTL/tamper 分页；不得同时要求 operation history sequence 与 snapshot ref。

### Latest Findings

- 三类 evidence union 已闭合，非 readback 分支不再需要伪造 observed 事实。
- reducer 已闭合 provider acceptance、control certainty、conclusiveness 与 deadline/预算的优先顺序。
- lost lease 在预算未耗尽时也唯一进入 manual_attention。
- Operation History 已从 Run-local SSE sequence 解耦为物化列表快照。
- 跨两个 Run 的 lifecycle/attention 变化不会改变已签发分页集合。
- entity subject、fingerprint alias、locator 与 semantic claim 在 secret/key rotation 后保持同一 lineage。
- fingerprint alias 冲突已具有稳定 fail-closed 结果。
- sent/cancel CAS 已证明批次全有或全无，不会出现部分取消后继续发送。
- cancel 获胜时零 mutation 证据、claim release 和 audit tombstone 同事务提交。
- sent 获胜时 cancel 不再伪报 Run cancelled/stopped。
- retry claim 取消后统一回到 unclaimed，会丢失 predecessor definitive-failed 的 retry-only 门禁。
- audit tombstone 若被当 retained claim，当前又没有从 cancelled tombstone 到 retry 的合法转移，可能永久阻断后续 retry。
- AC-33 的旧 `snapshot_sequence` 与 AC-37 的 `operation_list_snapshot_ref` 仍是同一 canonical matrix 内的竞争契约。

### Latest Gate

维持 `NOT IMPLEMENTATION-SAFE`，直到：

1. `cancelled_before_send` 对普通 claim 与 failed-retry transferred claim 定义不同的释放／恢复结果，并有 AC；
2. AC-33 删除 Operation History 的旧 `snapshot_sequence` 术语。

完成后无需再次扩展架构范围，只需 targeted recheck 这两个点即可。

## Final Targeted Recheck — 2026-07-15

### Verdict

**IMPLEMENTATION_SAFE。** 上轮仅存的两个精确接缝均已关闭，并且修正已从 Architecture Spine 横向传播至 Runtime SPEC、domain contract、state machine、迁移 ADR、run lifecycle、Epics Story AC、acceptance matrix、实施计划与验收计划。当前未发现会让两个独立实现者产生不同权威 claim 状态、retry 门禁或 Operation History 分页契约的竞争规则。

### Closure Evidence

- `SemanticMutationClaim` 已固定 `claim_origin=initial|failed_retry`、predecessor owner/status 与 lineage root；retry 转移不再丢失 predecessor。
- 首次 claim 的 sent 前取消唯一进入 `released_zero_write`，明确不阻断未来普通新决策。
- failed-retry transferred claim 的 sent 前取消在同一事务原子恢复 predecessor owner/status=`failed`，不会进入自由 `unclaimed`。
- 取消 child 只保留 audit tombstone；tombstone 明确不能替代权威 claim 状态。
- 恢复后的 failed claim 继续阻断普通 draft，消除了“retry 取消后绕过 lineage 再次写入”的攻击路径。
- 新的合法 lineage retry 仍可从恢复后的 definitive failed claim 事务转移，消除了 tombstone 导致永久阻断的另一种竞争实现。
- sent/cancel 仍使用 draft 级全有或全无互斥 CAS；sent 获胜时不修改 claim，cancel 获胜时同时写零 mutation 证据和来源感知 disposition。
- AC-29 验证 initial/retry 两种取消 disposition 与 sent/cancel 竞争；AC-31 验证 failed → retry reserved → cancelled_before_send → 普通 draft 拒绝 → 合法 lineage retry 可转移；AC-34 验证 union/schema/migration 映射。
- AC-33 已改为在同一 immutable `OperationListSnapshot` 上完成 calendar window 与 lifecycle-open union，并由 cursor 绑定 `operation_list_snapshot_ref`。
- AC-33 明确不要求或复用 Run-local SSE `snapshot_sequence`；AC-37 继续独立验证跨 Run 变化、TTL、篡改和 scope/filter 失配，二者不再定义竞争 cursor 模型。
- Architecture AD-11/23/26、SPEC、domain/state machine、ADR、lifecycle 与 Story 1.22/1.36 对上述术语、事务边界和失败结果一致。
- 实施计划和验收计划已把来源感知 claim 恢复与独立 OperationListSnapshot 纳入开发顺序和纠偏门禁，不会只停留在设计说明。

### Final Gate

本次 targeted recheck 不再保留架构兼容性阻断；可进入后续 Implementation Readiness 检查与按 Milestone/Story 顺序实施。`IMPLEMENTATION_SAFE` 只表示本报告追踪的跨层契约已达到可实现的一致性，不替代 schema/fixture 生成、迁移测试、并发 CAS 测试、contract test 或项目既定阶段验收。
