---
id: SPEC-product-facts-v2-runtime-foundation
companions:
  - domain-contracts.md
  - state-machines.md
  - migration-and-gates.md
  - acceptance-matrix.md
  - ../../planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14/ARCHITECTURE-SPINE.md
  - ../../planning-artifacts/product-blueprint/prd.md
  - ../../planning-artifacts/product-blueprint/implementation-readiness-gate.md
  - ../../planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/DESIGN.md
  - ../../planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/EXPERIENCE.md
  - ../../../docs/adr/2026-07-14-product-facts-v2-eino-runtime-migration.md
  - ../../../docs/adr/2026-07-15-greenfield-no-legacy-compatibility.md
  - ../../../docs/facts-contract.md
  - ../../../docs/conversation-context.md
  - ../../../docs/capability-provider-contract.md
sources: []
---

> **Canonical contract.** 本 SPEC 与 `companions:` 共同构成 Product Facts v2 Runtime Foundation 的完整机器契约；下游 Epic、Story、实现和验证必须全部读取。

# Product Facts v2 Runtime Foundation

## Why

当前过渡性 Product Facts 只持久化工具结果摘要，缺少可冻结查询集合、不可变操作草案、逐条核验、稳定事件游标和真实 Eino 工具/HITL 执行，无法安全完成“查询—人工判断—确认—写入—回读”闭环。受影响的是产品中心运营、管理员、FOBrain 系统负责人以及所有消费 Workbench、Action API、replay 和 audit 的实现单元；必须直接建立统一、可恢复、可证明的 v2 事实基础并替换过渡实现，再允许任何真实外部写入。

## Capabilities

- **CAP-1 — Product Facts v2**
  - **intent:** 系统能够用版本化单一事实模型保存运行、工具、查询、操作、核验与审计事实。
  - **success:** v2 schema、领域类型、仓储、持久化和全部产品投影通过同源映射测试，且不存在 v1 adapter、dual write 或旧 epoch 产品读取。

- **CAP-2 — 安全结果与查询快照**
  - **intent:** 系统能够持久化安全 StructuredResult，并把一次查询冻结成范围、顺序和逐行事实确定的快照。
  - **success:** 刷新或重启后同一 `result_ref` 恢复相同对象集合；只有完整抓取所有分页的快照可用于“全部”动作。

- **CAP-3 — 确定性多轮引用**
  - **intent:** 系统能够理解用户指代并由程序唯一锁定其引用的历史结果。
  - **success:** 唯一候选自动解析，多候选进入澄清，过期、越权、跨会话或不存在的引用被拒绝；上下文压缩不改变对象集合。

- **CAP-4 — 不可变操作草案与确认**
  - **intent:** 系统能够从冻结事实和人工参数创建不可变操作草案，并在用户确认后才允许执行。
  - **success:** 确认绑定同一 actor、草案版本、摘要和对象范围；任何内容变化产生新草案，未确认、拒绝、取消或过期均产生零外部写入。

- **CAP-5 — 能力策略与身份约束**
  - **intent:** 系统能够按版本化能力策略约束新鲜度、核验、基数、分组、候选来源、角色和对象权限。
  - **success:** prepare 与 confirm 均完成策略检查；未知身份、权限、策略版本、凭据上下文或对象归属一律阻止写入。

- **CAP-6 — 逐条执行与客观核验**
  - **intent:** 系统能够独立执行、回读和判定每个操作对象，并保留批量部分成功事实。
  - **success:** 只有 provider 接受且回读证据符合目标的条目成功；未知结果进入核对或人工排查，成功条目在确认重试和进程恢复中都不重复写入。

- **CAP-7 — Eino-first 双入口执行**
  - **intent:** 自然语言入口能够使用 Eino 完成模型工具循环和人工确认，确定性 Action API 能够跳过模型但复用相同产品控制面。
  - **success:** 两个入口对同一操作产生相同草案、策略、审批、幂等、执行、核验和产品事实；确定性入口不伪造 Eino checkpoint。

- **CAP-8 — 持久化 Run 与崩溃恢复**
  - **intent:** 系统能够让 Run 独立于浏览器连接推进，并在进程重启后安全恢复或进入明确待排查状态。
  - **success:** waiting continuation 可恢复，确认只消费一次，执行租约可防迟到提交；不确定外部写入先回读，不能证明时进入 `manual_attention`。

- **CAP-9 — 同源投影与稳定 SSE**
  - **intent:** Workbench、ActionResult、replay、audit、Inspector 和模型上下文能够消费同一事实及稳定增量。
  - **success:** 初始高水位与后续持久化事件无缝衔接，重复、乱序、断线和过旧游标不会丢失、重复或复活已删除/替换状态。

- **CAP-10 — 可验证 Greenfield bootstrap 与运行门禁**
  - **intent:** 系统能够从空库直接建立目标 v2 schema，只在工具链、存储、恢复和安全门禁有证据后启用 writer，并对旧 epoch fail closed。
  - **success:** `G-ARCH-V2`、`G-TOOLCHAIN` 通过，legacy epoch 返回 `unsupported_schema_epoch` 且 readiness=false；真实写入仍由适用的 `G-FACT/G-WRITE/G-SAFE` 独立授权。

- **CAP-11 — 操作记录安全投影**
  - **intent:** 系统能够让当前执行人在离开会话或重新登录后，从固定入口重新找到其获授权的历史操作和未关闭结果。
  - **success:** 记录只从 Product Facts、ActionResult 与安全审计引用投影；默认发现最近六个月，未关闭对象持续可发现，详情与当前结果同源且只读，不重新调用 FOBrain、不生成重试或新草案。

## Constraints

- `implementation-readiness-gate.md` 是实施授权的唯一裁决源；所有适用门禁 PASS 前禁止实现或启用真实 FOBrain mutation。
- Product Facts v2 是唯一运行时事实模型；不实现 v1 历史读取、projection adapter、dual write 或旧 waiting Run 恢复。
- raw provider payload、凭据、raw prompt、checkpoint/interrupt ID 不得进入 Product Facts、模型上下文或产品出口。
- 采用当前聚合状态加不可变 FactEvent，不通过全量事件回放重建整个数据库。
- 所有修改类能力默认二次确认；AI 不决定接收人、误报、延时或期限，前端、HTTP 和 provider 不授予权限。
- 自然语言工具编排进入 Eino Runner；`facts/product/httpapi/providers/*/web` 保持 Eino-free。
- 首版采用受控单操作人、模块化单体、单后端实例和 SQLite；不声明多人审计或高可用。
- API 使用 OpenAPI 3.1 与 JSON Schema 2020-12，前端只消费目标 generated contracts；被替代的过渡 schema 直接退出目标运行时，不生成兼容 union。
- 首版只交付浅色桌面三栏；移动端、暗色和触屏专用体验不构成门禁。
- 下一实现阶段使用 Go module language 1.26.0、Go toolchain 1.26.5 与 Node 24 LTS；依赖解析由受控 module/lockfile 和可复现安装证明。
- v2 状态迁移、幂等、continuation、SSE、verifier、Greenfield bootstrap 和故障恢复必须先有 schema、fixture 和自动化验收，不能只验证 happy path。
- `reconciling` 唯一投影为“待核验”且 Run 保持 `running`；明确 `failed` 或 `manual_attention` 唯一投影为“待排查”。verification deadline 与调用预算必须持久化，重启不得重置。
- 操作记录读取必须重新校验 actor/workspace 和 frozen history policy，只能使用 Product Facts、ActionResult 和安全 audit refs；operation_at 首次发布后不可变，服务端把跨 Run 六个月窗口与 open lifecycle 并集物化为独立 immutable OperationListSnapshot，分页 cursor 不得复用 Run-local SSE sequence。
- attempt 采用 phase 判别 union，在首次外部 mutation 前原子进入 durable `sent` 并固定 settle deadline；`sent` 之后只允许 Verify/Reconcile。Confirm 事务必须获取唯一 SemanticMutationClaim；普通草案命中任一 claim 均拒绝，只有 definitive failed item 可通过 item-level lineage 和新人工确认事务转移 claim。cancel/stop 只能在 sent 前以互斥 CAS 进入 cancelled_before_send、证明零 mutation，并按 claim 来源把首次 claim 置为 `released_zero_write`、把 failed-retry transferred claim 恢复为 predecessor failed；sent 后不得伪报取消。
- VerificationEvidence 必须以 evidence_kind、provider acceptance 和 control certainty 构成可编译判别 union；唯一 reducer 先处理失租约，再要求 held+accepted+conclusive real-readback match 才成功，provider adapter 不决定产品 outcome。
- 产品 `entity_ref` 必须 opaque 且不可拆解；随机不可变 entity_subject_key 连接跨查询身份、semantic key 与独立 ProviderLocator，版本化 HMAC 只作为可轮换 fingerprint alias。解析只通过 facts LocatorRepository/capabilities LocatorResolver，secret 轮换不得改变 subject、claim 或 locator lineage。

## Non-goals

- 本 SPEC 不实现或授权派发、转发、延时、误报等真实 FOBrain 写入。
- 不包含自动派发、自动选择接收人、自动误报/延时、通知或完整漏洞管理闭环。
- 不包含多用户 SSO、每用户凭据、PostgreSQL、多实例、微服务或生产 SLA。
- 不包含完整 Event Sourcing/CQRS、跨 session 自然语言长期记忆、向量检索或生产 MCP 市场；CAP-11 的结构化操作记录不属于自然语言长期记忆。
- 不包含移动端、暗色主题或触屏专用交互。

## Success signal

在 v2 受控演示中，用户可以查询并冻结事实、唯一引用结果、创建和确认不可变草案，并看到逐条核验或安全模拟结果；浏览器断线和进程重启后，Workbench、Action API、replay 与 audit 仍展示同一事实。任何没有获得 readiness gate 授权的真实写入始终保持关闭。
