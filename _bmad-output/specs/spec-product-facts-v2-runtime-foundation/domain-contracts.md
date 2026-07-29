# Domain Contracts

本文件固定 v2 领域所有权和跨层可编译接缝。字段名是契约语义；正式 JSON Schema/Go 类型可按命名规范落地，但不得改变所有权、基数和引用方向。

## 所有权

| 聚合/值对象 | 定义所有者 | 唯一命令所有者 | 只读消费者 |
| --- | --- | --- | --- |
| Run、Turn、ToolCall、ToolResult | `facts` | `execution` | `product`、context assembler |
| QueryResultSnapshot、SnapshotItem | `facts` | `execution` | resolver、product、ActionDraft builder |
| ActionDraft、ActionItem、ActionAttempt | `facts` | `execution` | policy、product、audit/replay |
| EntitySubject、EntityIdentityFingerprint、ProviderLocator、SemanticMutationClaim | `facts` | `execution`（claim）／resolver（identity/locator read） | capabilities、provider adapter、retention scanner |
| PendingInteraction、ContinuationRef | `facts` | `execution` | product、httpapi |
| FactEvent、VerificationEvidence、AuditEvent | `facts` | `execution` | product、observability |
| Capability/Catalog v2 policy refs | `capabilities` | registry/bootstrap | execution、Eino tool adapter |
| OperationListSnapshot | `facts` | `product` 通过 repository port 创建只读物化快照 | httpapi、Workbench、retention scanner |
| OperationRecordView、OperationRecordDetail | `product` | 无，只读投影 | httpapi、Workbench |

`product` 不推进状态，provider 不写 facts，Workbench/Action API 不维护平行事实。

## 版本边界

| Schema epoch | 新建 Run | 可读取 | 可创建 QuerySnapshot/ActionDraft |
| --- | --- | --- | --- |
| 非目标／legacy epoch | 否 | 否；`unsupported_schema_epoch` | 否 |
| Product Facts v2 目标 epoch | 是 | 是 | 是，但受有效期、归属、权限与 readiness gate 限制 |

Product Facts v2 是唯一运行时模型。bootstrap 只创建／接受目标 epoch；不提供旧投影 adapter、历史回填或 dual write。

## StructuredResultValue

Safety Gate 通过后，ToolResult 必须按值持久化：

| 字段 | 规则 |
| --- | --- |
| `schema_version` | 必须来自登记 schema。 |
| `result_ref` | workspace 内稳定安全引用。 |
| `safe_payload` | 完整安全 JSON；不得包含 raw provider payload。 |
| `content_digest` | 对规范化安全 JSON 计算 SHA-256。 |
| `byte_size` | Safety Gate 限额依据。 |
| `created_at` | RFC 3339。 |

ToolResult 的 presentation 只能从该值重新派生，不能成为事实来源。

## QueryResultSnapshot

必填语义：

- `query_snapshot_id`、`snapshot_version`、`snapshot_digest`
- `workspace_id`、`conversation_id`、`actor_subject_id`
- `source_run_id`、`source_tool_call_id`、`source_result_ref`
- `query_descriptor`：安全、结构化的查询语义与参数摘要
- `coverage`：`complete_set | partial_page | incomplete`
- `order_spec`：排序字段、方向和稳定 tie-breaker
- `total_count_known`、可选 `total_count`
- `captured_at`、`expires_at`
- 按确认显示顺序排列的 `items[]`

`SnapshotItem` 必填：`item_ref`（产品 opaque `entity_ref`）、`entity_subject_key`（仅内部）、`position`、`entity_type`、`facts_schema_ref`、冻结安全事实、`facts_digest`。跨 workspace、conversation 或 actor 禁止复用；任一分页失败不能发布 `complete_set`。

## Opaque EntityRef 与 ProviderLocator

产品可见 `entity_ref` 必须是不可逆、不可拆解的随机 opaque 安全引用。`facts` 首次识别 provider 业务对象时用 CSPRNG 生成随机、不可变的内部 `entity_subject_key=esk.v1:base64url(32_bytes)`；它不由 secret 推导。provider business identity 只在 provider／locator 边界出现。`EntityIdentityFingerprint(workspace, derivation_key_version, fingerprint, entity_subject_key)` 使用受控版本化 workspace subject-derivation key 生成 `esf.v1.<key_version>:base64url(HMAC-SHA-256(key,JCS(provider_instance_id,entity_type,canonical_provider_business_identity)))`，作为查找既有 subject 的内部别名；`(workspace,key_version,fingerprint)` 必须唯一映射一个 subject。`EntityRefMapping(workspace, entity_ref, entity_subject_key, source_capability, expires_at)` 连接一次查询引用与内部主体。

每次解析必须对 active 与所有 retained derivation key version 计算 fingerprint；命中任一 alias 时复用同一 entity_subject_key，并在同一事务为 active version 补 alias；多个 alias 指向不同 subject 时返回稳定 `entity_identity_conflict` 并 fail closed，全部未命中才允许新建 subject。旧 key version 只有在全部保留 subject 已迁移 alias，且 snapshot/draft/attempt/semantic claim/locator/open attention/audit retention 均结束后才能退役；locator encryption key 独立轮换且不得改变 subject。轮换前后的同一 provider 对象必须得到同一 entity_subject_key、semantic mutation key、claim namespace 与 locator lineage。

Provider 的真实对象 ID、资产类型、路由选择和 locator version 组成独立 `ProviderLocator`，以 `(workspace_id, entity_subject_key, locator_version)` 为唯一键，由 `facts.LocatorRepository` 加密／安全持久化；`capabilities.LocatorResolver` 是 provider adapter 唯一只读端口。解析链只能是 `entity_ref → EntityRefMapping → entity_subject_key → ProviderLocator`。Workbench、Action API、LLM、日志和安全 StructuredResult 不得看到 entity_subject_key/locator，也不得通过字符串拆解恢复。

列表进入详情时，resolver 必须重新校验 workspace、actor、source capability、locator version、对象权限和有效期；任一项 unknown/denied/expired 时 fail closed。跨查询去重、ActionItem identity 与 semantic mutation key 只使用 entity_subject_key，永不使用随机 entity_ref、fingerprint 或 locator。subject/fingerprint alias/mapping/locator 在任何 snapshot、draft、attempt、semantic claim、open attention 或 audit retention 引用存在期间禁止清理；locator version 迁移必须 additive，旧引用可读直到保留期结束。FR-21～FR-24 的列表与详情 fixture 必须证明产品 ref 与 provider ID 不相等、不同对象不碰撞、错误类型不会默认为 internal asset、来源时间和采集时间均保留、subject secret 轮换不改变 identity/claim/locator，缺失事实使用契约定义的稳定缺失态。

## FOBrain 只读事实基线

精确“新增漏洞”查询必须使用 provider 可验证的新增状态语义，按 `discovered_at` 倒序，并以稳定业务身份作为同时间戳 tie-breaker。结果必须区分：

- 查询成功且集合为空：产品明确投影“没有待派发漏洞”；
- 当前 actor 无权或查询范围无效：拒绝或权限错误；
- provider 失败或分页不完整：失败/不完整，不能伪装为空集合。

每条漏洞 `SnapshotItem` 的统一安全事实至少包含：POC 名称、IP、在线状态、业务系统、当前修复负责人。派发/转发判断还应在 provider 可用时包含业务系统负责人和运维负责人；缺失值必须保持空或使用契约定义的稳定缺失态，禁止由模型猜测。Workbench、确认摘要、ActionResult、replay、audit 和模型上下文均从这一行级事实投影，不得分别拼装。

人员候选查询返回当前 FOBrain 可见的全部人员，并为每个人提供稳定业务身份和显示名。模型不得根据姓名、职位或历史行为自动推荐或选择接收人；人工选择后，程序按稳定业务身份写入草案。

## Reference Resolution

Resolver 输入为当前 scope、用户指代表达和候选快照元数据；输出只能是：

- `resolved(query_snapshot_id)`：唯一有效候选；
- `clarification_required(candidate_refs[])`：多个合理候选；
- `invalid(reason)`：不存在、过期、越权、跨 scope 或不可验证。

模型只提取序号、时间、标签或显式 ref，不拥有最终选择权。上下文摘要只保存 safe refs；具体事实每次从 repository 重载并重验 scope。

## ActionDraft

必填语义：

- `action_draft_id`、`draft_version`、`draft_digest`
- `workspace_id`、`conversation_id`、`actor_subject_id`
- `capability_id/version`、`policy_ref/version`、`verification_policy_ref/version`
- `provider_route_version`、`credential_binding_version`
- `history_access_policy_ref/version`
- `query_snapshot_id`、`snapshot_digest`
- `action_type`、规范化人工参数、`items[]`
- `status`、`created_at`、`expires_at`
- 可选 `retry_of_draft_id`；retry draft 的每个 item 还必须绑定 `retry_of_action_item_id`

`draft_digest` 固定为 `v1:base64url(SHA-256(JCS(confirmable_payload)))`。confirmable payload 包含 actor、全部版本引用、snapshot id/digest、按确认显示顺序排列的 item identity/expected facts、动作参数和 expires_at。服务端生成并常量时间比较；任何字段变化产生新 version/digest。

一个 ActionDraft 对应一个 approval PendingInteraction 和一个一次性 approval ref；一次确认覆盖整个草案，不能为每个 item 单独创建竞争审批。

## ActionItem 与 ActionAttempt

`ActionItem` 必填：

- `action_item_id`、`action_draft_id`、`source_item_ref`
- `entity_subject_key`、`expected_facts`、规范化动作参数、`group_key`
- `mutation_key`、`semantic_mutation_key`
- 可选 `retry_of_action_item_id`
- `status`、`attention_status=not_applicable|open`、attempt refs、最终 outcome 或安全错误；首版 writer 不接受 `closed`

`mutation_key = v1:sha256(workspace|draft_id|draft_version|entity_subject_key|action_type)`。`semantic_mutation_key = smk.v1:base64url(SHA-256(JCS(canonical_payload)))`，payload 字段固定为 `version/workspace_id/provider_instance_id/entity_subject_key/action_type/normalized_target`。target 由 capability target schema 规范化：字符串 Unicode NFC；时间转 UTC RFC3339 纳秒；missing 与 null 分离；对象键由 JCS 排序；声明为 set 的数组按元素 digest 排序去重，其他数组保留顺序；数字必须满足 schema 且使用 JCS 表示。任何规则变化必须升级 `smk` version。

`SemanticMutationClaim` 以 `(workspace_id, semantic_mutation_key)` 为数据库唯一主键，必填 owner item、lineage root、claim status、`claim_origin=initial|failed_retry`、可空 predecessor owner/status、created/updated sequence。Confirm 必须在同一 SQLite 事务中先插入／转移 claim，再预留 `reserved_not_sent` attempt；并发普通草案只有一个能成功。既有 `reserved_not_sent/sent/verifying/reconciling/succeeded/manual_attention/failed` claim 全部阻断普通草案；`released_zero_write` 不阻断后续普通新决策，但新 owner 必须在同一事务取得该唯一行。只有新人工确认的 retry draft 逐条引用同 lineage 的 definitive failed item 时，才能在同一事务把 failed claim 转移到新 item并保存 predecessor；缺 lineage、夹带其他状态或普通草案命中 failed 均返回 `semantic_mutation_conflict` 且 mutation count=0。首次发送前 cancel/stop 必须以有效 fencing token 在一个事务内 CAS `reserved_not_sent→cancelled_before_send`、写入 mutation count=0 证据并按来源处置 claim：`initial` 转为 `released_zero_write`；`failed_retry` 原子恢复 predecessor owner/status=`failed`，取消的 child 只保留审计 tombstone。若 `sent` 转换先提交，则 cancel/stop 返回 `action_already_sent`，claim 不变且 Run 继续核验。audit tombstone 不能代替权威 claim 状态；其余 claim 不自动释放，至少保留至 audit/draft retention、reconcile window、locator retention 与 open attention 的最晚者。

`ActionAttempt` 必须实现为 `phase` 判别的封闭 union：

| phase | 必填 | 允许 | 禁止 |
| --- | --- | --- | --- |
| `reserved_not_sent` | attempt/item/lineage/claim ref、created_at | 无 sent-only 字段 | `sent_at`、deadline、policy、epoch、预算、provider 摘要、next poll、evidence |
| `cancelled_before_send` | 全部 reserved 字段、cancel/stop decision、decided_at、mutation_count=0、`claim_disposition=released_zero_write|restored_failed_predecessor`、claim disposition audit ref | 安全取消原因 | sent-only 字段、provider 摘要、evidence、任何 mutation/verifier 调用 |
| `sent` | 前项 + immutable `sent_at`、`settle_deadline`、verification policy/verifier version、lease epoch、`verification_calls_consumed=0` | provider 接受安全摘要（收到后） | evidence/next poll 在核验前出现 |
| `verifying` / `reconciling` | 全部 sent 字段、已预留的调用计数 | `next_verification_at`、零到多个 evidence refs、provider 安全摘要 | 修改 sent_at/deadline/policy/epoch 或减少计数 |
| `terminal` | 全部 sent 字段、terminal outcome、最终 evidence ref | 安全错误／attention ref | 任何后续 mutation 或 verifier 调用 |

`reserved_not_sent → sent` 与 draft 级 cancel/stop 的 `reserved_not_sent → cancelled_before_send` 必须用 phase CAS 互斥。取消事务全有或全无：仅当全部 attempts 均未 sent 时才逐项取消、证明零 mutation，并依据每项 claim 来源分别写 `released_zero_write` 或恢复 predecessor failed claim；随后写 Run/ActionResult=`cancelled/cancelled` 或 `stopped/stopped`、outcome=`none_succeeded`。任一 sent+ 存在则整体返回 `action_already_sent` 且不修改任何 item/claim。`sent` 后禁止 mutation 重发或伪报取消。每次回读先事务预留并扣减一次调用预算，崩溃不返还。不得用空字符串、Unix epoch 或零值伪造尚不存在的 sent-only 字段，也不得保存 raw request/response。

## VerificationEvidence 与唯一 Reducer

`VerificationEvidence` 必须是以 `evidence_kind` 判别的封闭 union。共同必填：evidence id、attempt/item id、call index、verification policy ref/version、verifier id/version、expected digest、provider acceptance=`accepted|rejected|unknown`、`control_certainty=held|lost`、conclusiveness、proof reason code、decided_at 与 decision。分支规则：

- `readback_observation`：必填真实只读 observed StructuredResult ref/digest、observed_at 与 relation=`match|mismatch|unknown`；不得使用 synthetic/空 observed。
- `conclusive_non_application`：必填 relation=`conclusive_not_applied`、provider acceptance=`rejected`、control certainty=`held` 和 policy 允许的证明；只有确有 readback 时才允许 observed ref/digest，禁止为满足 schema 伪造。
- `control_loss`：必填 control certainty=`lost`、lease epoch/loss reason/lost_at；observed ref/digest 仅在失控前已有真实 readback 时允许，decision 只能是 `manual_attention`。

唯一 reducer 顺序为：① control certainty=`lost` → manual_attention；② held + accepted + conclusive readback match → succeeded；③ held + rejected + policy 允许的 conclusive_non_application → definitive failed；④ 其他组合在 deadline 与预算均未耗尽 → continue_reconciling；⑤ 否则 → manual_attention。`rejected|unknown + match` 不得成功。批次含 manual_attention 时 Run=`failed`、ActionResult=`blocked`，outcome 仍按成功项数量固定为 `partial` 或 `none_succeeded`，不得新增或留空枚举。

## ActorSubject

当前 `actor_subject_id = provider_instance_id + stable_user_business_id`。显示名不参与身份比较，credential fingerprint 只作审计上下文。prepare 与 confirm 都必须重新解析 current user 并校验一致；HTTP、前端和模型不得提交或覆盖 actor。

## Catalog v2

Catalog v2 直接定义以下稳定引用：

- `freshness_policy_ref`
- `verification_policy_ref`
- `max_items`
- `group_key`
- `recipient_source_capability`
- `actor_role_predicate`
- `object_permission_predicate`
- `confirmation_projection`

`VerificationPolicy` 固定 verifier id/version、回读 capability、expected/observed schema、settle timeout、poll interval、最大调用数和 outcome mapping。Provider 只返回 observed StructuredResult candidate；项目 `ActionVerifier` 生成 VerificationEvidence。

Outcome mapping 必须完整消费 control certainty、provider acceptance、evidence kind、relation、conclusiveness 与 deadline／预算：control lost 优先 manual_attention；held+accepted+conclusive real-readback match 才 succeeded；held+rejected+policy-approved conclusive non-application 才 definitive failed；其余在 deadline／预算内 reconciling、耗尽后 manual_attention。不得由 provider adapter 自行选择 outcome。settle deadline、已用调用次数和下一次回读时间按上述 ActionAttempt 事务持久化，重启不得重置；核验不能再次 mutation。

## Operation Record Projection

`OperationRecordView` 与详情不是第二套事实聚合，只能从 Product Facts、ActionResult、`attention_status` 和安全 audit refs 投影。列表至少包含业务动作、状态、`operation_at`、可选 `decision_at`、来源查询安全标识、对象数量、目标人员安全显示，以及完成／待核验／待排查计数；详情复用逐条安全 ActionResult，不保存或展示内部 ref、raw payload、凭据和 checkpoint。`operation_at` 在记录首次发布时赋值且永不变化：有 Pending 的流程一律取 Pending `published_at`，无 Pending 的确定性操作取服务端 decision time；后续审批只写 decision_at。

ActionDraft 冻结 `history_access_policy_ref/version`，Pending、ActionResult 与 OperationRecord index 原样携带。查询要求当前 `(workspace_id, actor_subject_id)` 与记录精确相等，并按该 frozen ref 通过当前授权 capability 证明整批对象均仍可见；ref 缺失／损坏返回稳定 `history_policy_unavailable`，当前 unknown/denied 返回 `history_access_denied`，均整批 fail closed，不过滤 item、不重算原 ActionResult。provider instance 改变不做身份别名。

首个列表请求不得提交 `as_of`；服务端生成 RFC3339 纳秒 as_of。转为 Asia/Shanghai 后，cutoff 取六个日历月前：目标年月向前 6，日=`min(as_of 日, 目标月最后一日)`，保留本地时分秒与纳秒，包含 `operation_at == cutoff`。在一个 SQLite 事务中，后端从当前 OperationRecord index、Run lifecycle 与 attention 求窗口内记录和 `lifecycle_open=true` 的去重 union，按 `(operation_at DESC, run_id DESC)` 排序，并把全部安全列表行、记录版本、total、scope、filters digest、as_of/cutoff 与排序键物化为 immutable `OperationListSnapshot`。open 仅指 Run waiting/running 或 attention open；首版 definitive failed/manual_attention 创建 open attention，`attention_status` 不存在 closed 值。

`OperationListSnapshot` 使用随机 opaque `operation_list_snapshot_ref`，明确不等于 Run-local FactEvent `snapshot_sequence`，也不能成为 ActionDraft、audit 或 replay 的事实来源。cursor 是服务端签发的 opaque、tamper-evident token，绑定 schema version、workspace/actor、filters digest、as_of、cutoff、operation_list_snapshot_ref 与 last keys；后续页只读取该快照的有序安全行，并对快照内全部 frozen history policy 重新验权，任一 denied/unknown 整批 fail closed。其他 Run 的状态／attention 变化、新记录和原记录更新都不改变当前快照。篡改、scope/filter 不匹配、快照缺失或配置化 TTL 过期返回 `invalid_cursor`；前端不二次合并。快照只为分页一致性保留到配置化 TTL，历史 Product Facts 仍按独立 retention 保存。

## FOBrain 动作约束

| 动作 | 基数/分组 | 候选和权限 | 核验目标 |
| --- | --- | --- | --- |
| 派发 | 同一接收人可合并；不同接收人拆 draft | 接收人来自全员列表；对象在当前可见范围 | 修复负责人等于目标人员 |
| 转发 | 同一接收人可合并；不同接收人拆 draft | 接收人来自全员列表；对象在当前可见范围 | 修复负责人等于目标人员 |
| 延时 | `max_items=1` | actor 为管理员；期限由人决定，原因可空 | 状态为延时且期限等于目标值 |
| 误报 | 按已确认对象执行 | 当前 actor 负责对象；判断由人作出 | 状态为误报 |

多条确认默认展示前 5 条，并提供完整冻结集合；仅数量摘要不能替代可核对明细。

派发和转发允许把已经逐条判断为同一接收人的对象合并到一个草案；不得仅凭查询集合相邻、字段相似或模型推断自动批量归组。不同接收人必须拆分草案。

## Transport Mapping

以下路径由目标 OpenAPI 统一定义；保留某条路径必须因为它符合目标 API，而不是为了旧客户端兼容：

| Command | HTTP 映射 | v2 必需材料 |
| --- | --- | --- |
| PrepareAction | `POST /api/workspaces/:workspace_id/agent/actions` | result ref、人工参数、client request id、scope；返回 draft/version/digest 与 approval waiting 投影 |
| ConfirmAction | `POST /api/workspaces/:workspace_id/runs/:run_id/resume` | resume ref、decision、client request id、draft id/version/digest |
| Query Run | 既有 run snapshot/action result endpoint | facts version、snapshot sequence、item outcomes |
| Subscribe | 既有 SSE endpoint | Last-Event-ID、持久化 sequence、view replacement |
| Query operation records | 新增 product query endpoint，具体路径由 OpenAPI v2 固定 | actor/workspace、时间/动作/状态筛选、稳定分页、默认六个月与未关闭补集 |
| Query operation record detail | 新增 product query endpoint，具体路径由 OpenAPI v2 固定 | 原 ActionResult/Product Facts 的逐条安全只读投影 |

同一 `client_request_id` 在 `(workspace, actor, operation)` 作用域内：相同 request digest 返回既有投影，不同 digest 返回 conflict。错误与状态契约分别给出 `read_retryable`、`verification_resumable` 和 `mutation_retry_allowed`；post-sent unknown 的最后一项恒为 false。传输只发布目标 schema version，不生成 v1/v2 兼容联合；写域只接受 Product Facts v2。
