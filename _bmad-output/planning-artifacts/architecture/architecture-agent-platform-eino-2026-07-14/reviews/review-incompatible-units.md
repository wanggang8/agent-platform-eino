# 对抗性不兼容单元评审

## Verdict

**NOT IMPLEMENTATION-SAFE。** Spine 已给出正确的安全不变量，但还不是足以让下一级实现单元独立开发后直接拼装的架构契约。下面两个单元可以不违反 AD-01 至 AD-19 的任何字面规则，却会在所有关键接缝上产生不可调和的协议分歧。当前 Spine 因而不能证明“相同 application command、相同 Product Facts、相同恢复语义、相同安全投影”在代码层真的是同一个东西。

## 可构造的两个实现单元

- 单元 A（控制面内核）：实现 `execution + facts + store/sqlite`。它采用规范化的 `QueryResultSnapshot / ActionDraft / ActionItem / ActionAttempt` 聚合、显式 prepare/confirm 命令、数据库租约、逐项 attempt 和由 verifier 返回观察事实再由 execution 判定成功。所有状态变化与不可变事实在 SQLite 事务内提交，外部 payload 不越界，写操作始终两阶段且 fail closed。

- 单元 B（产品与外部适配面）：实现 `httpapi + product + capabilities + providers/fobrain + web`。它采用当前 OpenAPI/JSON Schema 中的 `Run / ToolResult / PendingInteraction / ActionResult` 形状，以安全 ref 连接审批，通过现有 `/agent/actions` 与 `/resume` 暴露操作，provider 只返回 candidate/typed error，前端只消费 generated contracts，SSE 只消费 Product Facts sequence。

两者分别忠实采用六边形依赖、Eino/product 边界、StructuredResult、二阶段审批、逐条回读、幂等、安全投影、单实例 SQLite 和桌面 Workbench 等全部 AD；下面的冲突都来自 AD 没有固定的接缝，而不是任一单元违反了 AD。

## Findings

- Spine 在 `ARCHITECTURE-SPINE.md:95-107, 185, 247-251` 把 `QueryResultSnapshot`、`ActionDraft`、`ActionItem`、`ActionAttempt` 当作一等实体，却没有给出核心类型、repository port、JSON Schema、必填字段、枚举、版本兼容或引用方向。单元 A 可以把四者做成规范化聚合和专用仓储；单元 B 可以继续把查询结果嵌在 `ToolResult.structured_result`、把审批意图嵌在 `PendingInteraction`、把逐项结果投影成 `ActionResult.result_cards`。两种实现都只有一套 Product Facts 且都使用 StructuredResult，但没有共同可编译的数据形状；当前 `docs/schemas/eino_product_facts.v1.schema.json` 事实上也没有这四类实体。必须先固定领域契约、传输投影和 repository 命令三者的映射。

- Spine 在 `ARCHITECTURE-SPINE.md:83, 95, 107, 269-270` 同时使用 `StructuredResult candidate`、`ToolResult`、`result_ref`、snapshot 和 item refs，却没有规定 StructuredResult 是按值保存、按 ref 保存，还是按值加内容摘要保存，也没有规定 schema version、result ref、snapshot version 和 presentation version 的关系。单元 A 可以只在事实中保存不可变 `StructuredResultRef` 和摘要；单元 B 可以按当前 schema 把完整 `structured_result` 内嵌到 ToolResult/ActionResult。两者都从同一安全材料投影且不保存 raw payload，但 API、replay、context assembler 和 verifier 无法交换同一种结果。

- Spine 在 `ARCHITECTURE-SPINE.md:95, 101, 241-250` 只写 snapshot 保存“归属”，没有决定它究竟归 workspace、conversation、原始 run、tool result 还是这些实体的组合，也没有决定跨 run 的 Action API 是否可复用 snapshot。单元 A 可以把 snapshot 设为原 run 的子实体，任何新 run 只能复制出新 snapshot；单元 B 可以把 snapshot 设为 workspace 级事实，只要 resolver 重新校验 conversation、actor 和有效期就跨 run 复用。两者都逐项执行 AD-05/06 的校验，但一个单元交出的 `query_snapshot_id` 在另一个单元创建的 confirm run 中天然不可解析。

- Spine 在 `ARCHITECTURE-SPINE.md:107, 113, 248-251` 没有固定 `ActionDraft`、`PendingInteraction`、approval ref 与 run 的基数和所有权。单元 A 可以为一个批量 draft 创建一个 pending/approval ref，确认后展开全部 ActionItem；单元 B 可以为每个 ActionItem 创建独立 pending/approval ref，再由 UI 聚合成一张审批卡。两者都冻结完整 item refs、都逐条执行、都要求一次人工确认，但 duplicate confirm、部分过期、cancel 和 replay 的输入输出完全不同。

- Spine 在 `ARCHITECTURE-SPINE.md:125, 131, 137, 190` 混用了 action item 的 `verification_failed/reconciling/manual_attention`、执行阶段的 `queued/resuming/reconciling/terminal` 和产品 Run 终态，却没有一张总状态机或跨聚合映射。单元 A 可以把 `reconciling` 建成 Run.status；单元 B 可以保持 Run.status=`running`，只让 ActionItem.status=`reconciling`，再把 ActionResult 投影为 `blocked`。这两种状态都被持久化且都没有让前端重建状态机，但当前 `docs/schemas/eino_product_facts.v1.schema.json` 的 Run 枚举只有 `created/running/waiting/succeeded/failed/cancelled/stopped`，`docs/schemas/eino_action_result.v1.schema.json` 又使用 `accepted/running/waiting/completed/failed/cancelled/stopped/denied/blocked`；没有规范映射时后端、API 和前端会对同一事实得出不同终态。

- Spine 在 `ARCHITECTURE-SPINE.md:125, 131, 282` 要求显式部分成功和不确定结果，却没有定义批量聚合的完成规则。单元 A 可以在所有 item 都达到 `succeeded/failed/verification_failed/manual_attention` 后把 run 标为 terminal，并以 `partial_success` 投影；单元 B 可以让任何 `reconciling/manual_attention` item 阻止 run 终结，直到人工处理。两者都不把 2xx 当成功、都不补偿且都保留逐项事实，但轮询终止条件、SSE 关闭条件、ActionResult status 和“成功项不重试”的执行门禁不兼容。

- Spine 在 `ARCHITECTURE-SPINE.md:89, 113, 131` 要求若干数据库变化同事务，却没有规定 approval consumption、draft transition、attempt reservation、idempotency record、pending transition 与 immutable fact append 的精确提交顺序。单元 A 可以在 resume 前一次提交 `approved + consumed + attempt_reserved + fact`；单元 B 可以先成功加载 checkpoint，再一次提交 `consumed + attempt_reserved + fact`，同时让 pending 保持 approved 作为审计终态。两者都满足被点名的写入在同一事务，但在“提交后、Eino resume 前”崩溃时，一个恢复器看到已消费 ref 和已预留 attempt，另一个看到仍可提交的 approval；相同客户端重试会得到不同结果。

- Spine 在 `ARCHITECTURE-SPINE.md:77, 113, 131, 137, 248-253` 没有定义 Product Facts 事务与 Eino checkpoint store 的一致性协议。单元 A 可以先持久化 checkpoint mapping，再提交 waiting pending fact；单元 B 可以先提交 waiting fact，再保存 checkpoint。checkpoint 不是 provider raw payload，也不必成为产品事实，所以两种顺序都不直接违反 AD-02/04/12；但任一步之间崩溃会分别产生“孤儿 checkpoint”或“无 checkpoint 的可见审批卡”，且 Spine 没有补偿扫描、不可见 staging 状态或原子同库要求来收敛这两种状态。

- Spine 在 `ARCHITECTURE-SPINE.md:131` 只规定 item 幂等身份“至少绑定” draft、item、action type 和 attempt lineage，没有固定规范化算法、命名空间、provider 长度限制或 lineage 在 retry 时是否换 key。单元 A 可以使用 `hash(draft_version,item_ref,action_type,lineage_root)`，让所有恢复 attempt 共享 provider key；单元 B 可以使用 `hash(draft_version,item_ref,action_type,attempt_id)`，让每次新 attempt 有新 key但通过本地 reservation 防重。两者都包含所有要求字段，但 provider 去重行为相反：前者可能把显式新草案误判重复，后者可能在进程崩溃后的恢复 attempt 再次写外部系统。

- Spine 在 `ARCHITECTURE-SPINE.md:107, 113, 275-276` 要求 version/digest，却没有定义 digest 的算法、canonical serialization、字段集合、item 排序、Unicode/时间规范、算法版本和 constant-time compare 责任。单元 A 可以对按 item identity 排序后的 canonical JSON 求 SHA-256；单元 B 可以对数据库字段按存储顺序串联后求摘要。二者都保证任何意图变化产生新 version，且客户端都只回显安全 digest，但单元 B 发出的确认卡永远无法被单元 A 的 ConfirmAction 验证。

- Spine 在 `ARCHITECTURE-SPINE.md:113, 161, 272-276` 宣布共同的 `PrepareAction → ConfirmAction`，却没有把这两个动作投影成 OpenAPI operation、请求/响应 schema、状态码和错误码。单元 A 可以把 prepare/confirm 都做成 `/agent/actions` 上的判别联合；单元 B 可以按现有 API 用 `/agent/actions` prepare、`/runs/{run_id}/resume` confirm。两者都走同一 application command 和安全门禁，但当前 `docs/schemas/eino_action_request.v1.schema.json` 只有 `action_id/client_request_id/input.text`，没有 result ref、draft version 或 digest；现有 OpenAPI 也没有显式 prepare/confirm operation。生成客户端无法表达单元 A 的命令，单元 A 也无法验证单元 B 的 resume payload。

- Spine 在 `ARCHITECTURE-SPINE.md:125, 143, 149, 278-280` 要求 capability verifier，却没有定义 verifier port、输入、输出、expected/observed 比较责任、typed error、证据 schema 或 verifier 与 read capability 的关系。单元 A 可以注册 `Verify(target, expected) -> StructuredResult candidate` 并由 execution 比较；单元 B 可以在 registry 中引用另一个只读 capability，执行后按声明式字段路径比较两个安全结果。两者都让 provider 只返回 candidate、都由项目而非 provider 决定成功，但函数签名、证据事实和失败分类不兼容。更直接的证据是当前 `docs/schemas/capability_catalog.v1.schema.json` 没有 Spine 所要求的 `freshness` 或 `verifier` 字段。

- Spine 在 `ARCHITECTURE-SPINE.md:131, 137, 173` 要求有期限租约和未知写入先 reconcile，却没有定义 lease owner、epoch/fencing token、renewal、过期接管和 reserved attempt 的恢复矩阵。即使只有单实例，旧进程退出与新进程启动仍可短暂重叠。单元 A 可以在租约过期后由新 executor 取得更高 epoch，先回读再继续；单元 B 可以把所有失去租约的 reserved attempt 直接变为 `manual_attention`，绝不自动接管。两者都禁止盲重试，但一个恢复运行、一个永久停止，且没有共同 fencing 字段可防止旧 executor 的迟到提交覆盖新事实。

- Spine 在 `ARCHITECTURE-SPINE.md:161, 188, 232` 规定前端只消费 generated contracts，却没有定义 SSE patch 的操作语义：缺失字段表示“不变”还是“清空”、数组是 replace 还是 keyed merge、同 sequence 的 replace 与 item patch 谁优先、终态 patch 能否被迟到 running patch 覆盖。单元 A 的 product projection 可以发完整数组替换；单元 B 的 React reducer 可以按 id 合并数组。两者都使用同一生成类型、按 sequence 去重且不在 Zustand 复制 server state，但运行卡、审批卡和逐项结果会出现幽灵项或无法清空的旧状态。JSON Schema 只能约束形状，不能替代 patch 代数。

上述缺口至少需要以版本化领域 schema/port、完整状态机与投影映射、prepare/confirm OpenAPI、数据库事务与 checkpoint/lease 恢复协议、verifier SPI、SSE patch/cursor 语义和跨层 contract tests 固定。否则继续拆分实现只会把“同源”“同 command”“可恢复”留成命名一致、行为不一致。

## Recheck — 2026-07-14

### Recheck Verdict

**IMPLEMENTATION_SAFE。** 此结论取代上方初审 verdict，但只表示更新后的 Spine 与 Product Facts v2/Eino runtime migration ADR 已把原 top 5 接缝固定到足以安全启动“契约与迁移前置 epic”的程度；它不表示当前 v1 实现已经具备写域能力，也不解除 readiness gate 对真实外部 mutation 的阻断。

### Top 5 接缝复核

- 核心实体所有权与版本接缝已关闭。`ARCHITECTURE-SPINE.md:94, 118, 130, 208` 与迁移 ADR 第 1、2 节把 `facts` 固定为 v2 类型/仓储端口所有者、`execution` 固定为唯一命令与迁移所有者，并明确 StructuredResult 按不可变安全 JSON 值保存、SnapshotItem 按值冻结、snapshot 归属 `(workspace, conversation, actor)`、一个 draft 对应一个 pending/approval ref、v1 只读与 v2 writer 单向切换。更关键的是 AD-20 禁止在 v2 schema、repository、migration、OpenAPI、fixture、generated union 与映射 contract tests 同时交付前启动 writer，因此两个下一级单元不能再各自发明平行实体形状。

- 状态映射接缝已关闭。`ARCHITECTURE-SPINE.md:160, 212-214` 与迁移 ADR 第 4 节固定 Run 只保留七态，waiting 由 Pending/Draft 表达，执行与 reconcile 由 Item/Attempt 表达，并固定全部确定终结、部分成功和 manual attention 时的 Run/ActionResult/outcome 映射。原先把 `reconciling` 同时解释成 Run 状态或 Item 状态的两种实现已不再同时合法；AD-21 又要求 schema fixture 与投影测试锁定映射。

- 事务、checkpoint 与 lease 接缝已关闭。`ARCHITECTURE-SPINE.md:216-220, 238` 与迁移 ADR 第 6、8、9 节固定 staged checkpoint 先写、Pending 发布与 retained 标记同 facts 事务、Confirm 的 ref 消费/终态/draft approved/attempt reservation/FactEvent 单次提交、提交后才 resume/execute，并要求统一 SQLite connection policy、epoch/fencing token、迟到提交拒绝以及失租约外部调用进入 reconcile。原先“先发布 pending”与“先存 checkpoint”、或“提交前 resume”与“提交后 resume”的竞争协议均已被排除。

- Prepare/Confirm API 接缝已关闭。`ARCHITECTURE-SPINE.md:136` 固定 `/agent/actions` 执行 prepare、`/runs/:run_id/resume` 执行 confirm，confirm payload 必须提交 `action_draft_id/version/digest` 与决定；迁移 ADR 第 5 节固定聊天走真实 Eino interrupt/resume、确定性 Action API 走项目 continuation，但二者共享 draft、policy、approval、idempotency、execution 与 verification commands。AD-08/20 又明确现有 v1 schema 不足且 v2 OpenAPI/generated contract 是 writer 前置，因而不能把 v1 文本 action request 当成写域 prepare 契约。

- Verifier SPI 与 SSE patch 接缝已关闭。`ARCHITECTURE-SPINE.md:223-232` 与迁移 ADR 第 3、7 节固定 Catalog v2 的 verification policy 内容、provider 只产 observed candidate、项目 ActionVerifier 比较 expected/observed 并写 VerificationEvidence；同时固定 view 的 `snapshot_sequence`、只发送高水位后事件、完整实体 upsert/显式 tombstone、数组默认 replace、schema 声明后才允许 keyed merge，以及旧游标触发 `view.replaced`。provider、product projection 与 React reducer 不再能对“成功”或 patch 代数采用互斥解释。

### 剩余阻断项

- `G-ARCH-V2` 仍为 BLOCKED：实际 Product Facts/Catalog/Action/Resume/SSE v2 schema、Go domain/repository ports、additive SQLite migration、fixtures、OpenAPI、generated TypeScript union 与跨层 mapping contract tests 尚未交付。该阻断是更新后架构主动设置的落地门禁，不是新的接缝缺口。

- 状态映射在允许 writer 前仍须以 v2 enum/fixture/投影测试固化，尤其验证 `partial/none_succeeded`、`manual_attention -> Run failed + ActionResult blocked` 和成功 item 不被聚合终态覆盖。

- checkpoint/lease 在允许恢复执行前仍须完成真实 Eino interrupt address/checkpoint 集成、staged/retained 崩溃点测试、Confirm 事务故障注入、epoch/fencing 迟到提交测试和进程重启恢复验收。

- Prepare/Confirm 在允许写域入口前仍须新增兼容 schema 版本、错误码、duplicate/conflict fixture、OpenAPI operation 映射和 generated client contract；现有 v1 `eino_action_request` 不能替代这些产物。

- Verifier/SSE 在允许真实写域和新前端 reducer 前仍须交付 Catalog v2/VerificationPolicy/VerificationEvidence schema、verifier port contract tests，以及 upsert/tombstone、数组 replace/keyed merge、旧游标 `view.replaced` 和 snapshot-to-stream 无缝衔接测试。

- `G-FACT-01`、适用的 `G-WRITE-01..04` 与 `G-SAFE-01` 仍必须按 readiness gate 取得 PASS；在此之前只能实现契约、迁移、mock、只读事实链和安全验证，不能实现或启用真实外部 mutation。
