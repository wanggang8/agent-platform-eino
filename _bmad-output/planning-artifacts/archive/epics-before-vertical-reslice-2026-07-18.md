---
stepsCompleted:
  - step-01-validate-prerequisites
  - step-02-design-epics
  - step-03-create-stories
  - step-04-final-validation
inputDocuments:
  - _bmad-output/planning-artifacts/product-blueprint/prd.md
  - _bmad-output/planning-artifacts/product-blueprint/experiment-scope-proposal.md
  - _bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md
  - _bmad-output/planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14/ARCHITECTURE-SPINE.md
  - _bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/DESIGN.md
  - _bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/EXPERIENCE.md
  - docs/adr/2026-07-14-product-facts-v2-eino-runtime-migration.md
  - docs/adr/2026-07-15-greenfield-no-legacy-compatibility.md
  - docs/fobrain-tool-matrix.md
  - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/SPEC.md
  - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/domain-contracts.md
  - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/state-machines.md
  - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/migration-and-gates.md
  - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/acceptance-matrix.md
  - _bmad-output/planning-artifacts/sprint-change-proposal-2026-07-15.md
courseCorrection: 2026-07-15
lastValidated: 2026-07-15
---

# Agent Platform Eino - Epic Breakdown

## Overview

本文档把 Agent Platform Eino 的 FOBrain 漏洞处置实验需求、UX 契约、架构约束与 Product Facts v2 机器契约拆解为可实施 Epic 和 Story。

## Requirements Inventory

### Functional Requirements

FR-01：产品必须只使用当前用户已授权的 FOBrain 事实展示漏洞及关联上下文。

FR-02：产品在漏洞列表、二次确认和结果展示中必须呈现统一漏洞事实。单条直接展示完整对象；多条时确认卡默认展示前 5 条并提供“查看全部”，使执行人在确认前可以核对完整冻结集合；不得仅用数量替代待写对象明细。

FR-03：产品必须区分真实空结果、无权限、数据不足和外部失败；不得由 AI 虚构“没有待处理漏洞”等结果。

FR-04：运营必须能够查看全部新增漏洞，按发现时间倒序；无结果时展示“没有待派发漏洞”。

FR-05：运营必须先查看漏洞、IP 在线状态、业务系统、业务系统负责人和运维负责人，再从 FOBrain 全部人员列表选择接收人；AI 不得替代该决定。

FR-06：仅已逐条判断且最终选择同一接收人的漏洞可合并派发；不同接收人不得混合。

FR-07：派发仅在二次确认后写入；每条漏洞只有接口成功且回读修复负责人为指定人员时才完成。

FR-08：当前用户可主动对其内置权限范围内的漏洞发起转发；不要求系统识别错派或先发生其他动作。

FR-09：转发接收人必须从 FOBrain 全部人员列表选择；同一接收人可合并、不同接收人必须拆分。

FR-10：转发仅在二次确认后写入；逐条回读负责人为目标人员时才完成。

FR-11：运营与相关人员人工沟通后确定新修复期限；产品不得自动决定是否延时或期限。

FR-12：直接延时一次只处理一条漏洞，延时原因可不填；只有管理员完成二次确认后可执行。

FR-13：只有回读状态为“延时”且修复期限等于指定新期限时，结果才显示“延时成功”。

FR-14：操作人可对自己确认负责的漏洞发起误报标记；误报判断必须由人作出，不能扩展为自动规则。

FR-15：误报标记仅在二次确认后写入；只有回读状态为“误报”时才显示标记成功。

FR-16：任一写动作在未完成二次确认或取消确认时，不得向外部系统写入。

FR-17：任一写动作必须逐条区分完成、待核验和待排查。当前无法确认且仍可沿同一执行事实链继续核验时显示“待核验”；明确失败，或核验期限／调用预算结束后仍不一致、未知时显示“待排查”。成功项可独立完成；任何未完成项不得伪报成功，也不得通过重复外部写入完成核验。

FR-18：本产品不得发送消息、邮件、webhook、催办或通知中心通知；FOBrain 自身通知不构成本产品功能或成功条件。

FR-19：产品必须能够读取当前 FOBrain 用户接口提供的安全身份、部门和角色上下文，用于解释凭据主体和权限范围；部门等源字段未提供时必须显示稳定的“源字段不可用／未提供”，不得推断或伪报已解析，也不得展示 token、原始账号对象或凭据。

FR-20：产品必须能够读取并安全展示当前用户接口提供的 FOBrain 权限和数据范围，并明确区分 `resolved`、`confirmed_empty/no_permission` 与 `source_field_unavailable/unknown`；未知不得折叠为无权限或通过，不得展示原始 policy payload。

FR-21：产品必须能够按 IP 查询资产列表及其安全在线状态，作为统一漏洞事实和资产核对上下文。

FR-22：产品必须能够从资产列表进入单个资产的安全详情核对；safe ref 不得冒充 provider 真实 asset ID。

FR-23：产品必须能够按 IP 查询关联漏洞列表，展示安全漏洞事实并支持资产—漏洞联查。

FR-24：产品必须能够从漏洞列表进入单个漏洞的安全详情核对，为人工处置前判断提供事实。

FR-25：产品必须能够通过本轮已批准的 `business_list` 能力查询 FOBrain 全部业务系统安全列表，为统一漏洞卡片和业务上下文提供事实；可见范围只由 FOBrain 接口内置权限裁剪，不得替换为“我的业务系统”等已排除能力，也不得展示原始 provider payload。

### NonFunctional Requirements

NFR-01（安全事实）：外部原始返回不能直接作为产品展示、审计或模型上下文事实；产品只消费安全、结构化的业务事实。

NFR-02（权限边界）：产品不得扩大 FOBrain 权限、代替 FOBrain 管理权限，或将查询不到的对象视为可操作对象。

NFR-03（写入完整性）：所有数据修改均须绑定执行人二次确认、客观回读和逐条失败可见性；未确认、取消、拒绝或过期时外部写入数必须为零。

NFR-04（证据边界）：不把目标部署 API 已验证等同于产品已集成；不得对写入结果、写域状态值域或权限作静默假设。

### Additional Requirements

AR-01：下一实施阶段开始前，将 CI、版本文件和发布环境固定到 Go 1.26.x 与 Node 24 LTS，并在该环境重跑现有 schema、Go、TypeScript、stream、browser、checkpoint 与 SQLite 基线；当前 Go 1.23／Node 25 结果不构成通过证据。

AR-02：保持六边形分层和明确 import boundary：`execution` 推进 Run/Action，`facts` 拥有领域与 repository ports，`product` 只读投影，`capabilities` 管注册和策略，`providers/*` 只接外部边界，`bootstrap` 是唯一组合根；禁止万能 `utils` 和跨层直调。

AR-03：自然语言 Agent、模型工具选择、ReAct loop、Runner event 与聊天 interrupt/checkpoint 使用 Eino v0.9.12；工具来自 Capability Registry。Product Facts、policy、approval、幂等、安全投影、HTTP/SSE、audit/replay 由项目控制面拥有。

AR-04：只有 `execution`、`capabilities`、`observability` 和 `store/sqlite` 的指定集成接缝可以直接 import Eino；`facts`、`product`、`httpapi`、`providers/*` 与 `web` 必须 Eino-free。

AR-05：先交付 Product Facts v2、Catalog v2、Action/Resume/Result/SSE v2 的 JSON Schema 2020-12、fixtures、OpenAPI 3.1、Go domain、repository ports、generated TypeScript types 和映射 contract tests；目标运行时不生成 v1/v2 union 或旧契约 adapter。

AR-06：Product Facts v2 必须包含 QueryResultSnapshot、ActionDraft、ActionItem、ActionAttempt、FactEvent 与 VerificationEvidence；`execution` 是唯一命令和状态迁移所有者。

AR-07：Safety Gate 通过后的 StructuredResult 以完整不可变安全 JSON 值持久化，并保存 schema version、稳定 result_ref、SHA-256 内容摘要、字节数和创建时间；presentation 只能重新派生。

AR-08：QueryResultSnapshot 必须按 `(workspace_id, conversation_id, actor_subject_id)` 归属，冻结查询语义、来源、coverage、稳定排序、有效期和全部安全行事实；只有所有分页成功的 `complete_set` 可用于“全部”动作。

AR-09：Reference Resolver 只允许 `resolved`、`clarification_required` 或 `invalid` 三种结果；唯一候选可自动锁定，多候选必须澄清，跨 scope、过期、越权或不可验证引用必须拒绝。模型不拥有最终选择权。

AR-10：ActionDraft 必须冻结 actor、版本化 capability/policy/verifier/provider route/credential binding、源 snapshot、全部 items、人工参数和 expires_at；使用 RFC 8785 JCS + SHA-256 生成服务端摘要，任何内容变化必须创建新版本并重新确认。

AR-11：Workbench 与 Action API 共用 `PrepareAction → ConfirmAction → ExecuteItem → VerifyItem` 控制面；聊天使用真实 Eino continuation，确定性 Action API 使用项目 continuation record，不伪造 Eino checkpoint。

AR-12：一个 ActionDraft 只对应一个 PendingInteraction 和一次性 approval ref；Confirm 事务原子消费 ref、固定 Pending 终态、复核 actor/policy/digest/version、批准 draft、预留 attempts 并追加 FactEvent，提交后才 resume 或执行。

AR-13：Capability Catalog v2 必须版本化声明 freshness、verification、max_items、group_key、recipient source、actor role、object permission 和 confirmation projection；Eino schema、prepare、policy 与 routing 均由注册表派生，不按工具名、关键词或 provider 私有字段硬编码。

AR-14：Actor 只能由 workspace FOBrain 凭据解析 current user，规范化为 provider instance + stable business user ID；前端、HTTP 和模型不能提交或覆盖 actor，prepare 与 confirm 都必须重新校验。

AR-15：每个 ActionItem 使用跨草案稳定的语义 mutation key，并保存 item-level retry lineage；同一确认、client request、技术重试和进程恢复必须幂等。Attempt 必须持久区分 `reserved_not_sent` 与 `sent`，一旦进入 `sent` 就只能核验，不得重发；只有可证明未发送且由人创建的新草案才可能再次写入。

AR-16：项目 ActionVerifier 必须按版本化 VerificationPolicy 比较 expected/observed 安全事实并生成 VerificationEvidence；provider 只返回 StructuredResult candidate，不决定产品成功。

AR-17：Run 只使用 `created/running/waiting/succeeded/failed/cancelled/stopped` 七态；Pending、Draft、Item、Attempt 和 ActionResult 必须遵循唯一映射，不创建竞争 Run 状态。

AR-18：Run 与浏览器/SSE 连接解耦，由后台 executor 使用持久化状态、lease epoch 和 fencing token 推进；旧 owner 的迟到提交必须被拒绝，未知外部写入进入 reconcile。

AR-19：每次聚合迁移与 FactEvent 在同一 SQLite 事务提交并获得 Run 内单调 sequence；view 返回 snapshot_sequence，SSE 从其后续传，只使用完整 upsert 或显式 tombstone，旧/未知游标触发 `view.replaced`。

AR-20：SQLite 采用统一 database handle/connection policy，启动创建／校验 Greenfield schema epoch，并验证 WAL、每连接 foreign keys、配置化 busy timeout、pool/transaction 策略、integrity/foreign-key check、可写探针和恢复扫描；任一失败时 readiness=false 且不接业务流量。

AR-21：新数据库必须直接创建目标 schema epoch 并可重复启动；旧 epoch 必须返回 `unsupported_schema_epoch` 且 readiness=false，不自动迁移或删库。Greenfield 基线后的未来 schema 变化才使用前向、可重复迁移与可恢复备份。

AR-22：Product Facts v2 是唯一运行时事实模型；不实现历史 v1 Run 读取、projection adapter、旧 waiting Run 恢复或 dual write。当前过渡数据库由操作者显式备份／删除／重建。

AR-23：waiting/running/reconciling 引用的 continuation、snapshot 和 draft，以及仍有 open attention 的结果禁止 cleanup；首版不自动物理删除 terminal Product Facts 或 audit；幂等记录保留到 audit/draft retention 与 reconcile window 的较晚者。

AR-24：API 使用 OpenAPI 3.1、JSON Schema 2020-12、业务成功 schema 和统一错误 envelope；时间使用 RFC 3339，产品统一显示 UTC+8；稳定产品 ID 不暴露内部 checkpoint、credential 或 provider 标识。

AR-25：前端只消费 `contracts/generated.ts`；TanStack Query 只管 server state，Zustand 只管 UI state；前端不得重建 Product Facts、Run lifecycle、approval/resume、provider DTO 或幂等状态机。

AR-26：FOBrain provider adapter 负责认证、配置化 timeout、输入映射、raw response 归一化和错误脱敏；真实字段映射必须有 schema、fixture 和 live evidence，provider 不得写 facts、view 或直接授予权限。

AR-27：配置、provider、模型、端口、URL、凭据、预算、timeout、freshness 和审批策略只能来自受控配置或注册表；日志/callback 仅记录脱敏诊断、trace、latency、token 和预算，不记录 raw prompt/completion/provider payload。

AR-28：首版部署为模块化单体：一个 Go 后端实例、独立 React 静态构建和单实例 SQLite；多实例、高可用、共享 SQLite、PostgreSQL 或微服务需要新的迁移 ADR。

AR-29：实施必须严格按 M-0 工具链、M-1 契约、M-2 存储、M-3 投影、M-4 Runtime、M-5 只读事实链、M-6 写域准入推进；M-1 至 M-5 不得分散进四项 FOBrain 写动作 Story。

AR-30：必须先以 schema、fixture、mock 和自动化失败注入覆盖 continuation 各崩溃点、重复确认、provider 响应丢失、回读最终一致、lease 失效、checkpoint 损坏、SSE 竞态/旧游标、fresh database bootstrap、legacy epoch 拒绝和备份恢复；未到 Phase 的测试必须明确 blocked/skipped，不能算 PASS。

AR-31：`implementation-readiness-gate.md` 是真实写域实施授权的唯一裁决源。当前只允许实现 G-READ-01/02 与 G-FACT-01 所需的精确新增查询、全员列表、统一行级事实和安全投影；适用的 G-ARCH-V2、G-TOOLCHAIN、G-WRITE 与 G-SAFE 未有机器证据 PASS 前，不得实现或启用真实 mutation。

AR-32：本轮 FOBrain 只读恢复只包含 Owner 已批准的 7 个代表性能力，以及主流程必需但当前缺失的精确新增查询和全员列表；选定能力的 tool id、input schema、StructuredResult、fixture、mock/live 断言、截图和敏感字段必须遵守 `docs/fobrain-tool-matrix.md`，其余 17 个矩阵工具不得进入本轮 Epic。

AR-33：PI-002/PI-004 写域启用前必须关闭 OQ-02，明确 243 人全量列表的搜索、筛选、排序、同名消歧、键盘选取和不可选规则；不得用 AI 推荐或静默默认替代。

AR-34：PI-005 写域启用前必须关闭 OQ-04，验证直接延时允许的原状态集合、新期限范围、时间粒度、时区和与当前期限的关系，并形成 schema、policy、fixture 和目标部署证据。

AR-35：跨会话操作记录采用已关闭的 RD-03／AD-26：IA-01 左侧提供固定“操作记录”入口，仅展示当前 workspace、当前凭据主体且重新授权通过的 Product Facts／ActionResult 安全只读投影；按服务器 UTC+8 `as_of` 的前六个自然月窗口与仍未闭合的生命周期对象并集查询，稳定游标分页，不从聊天摘要或 provider raw payload 重建，不提供草案、确认或重试入口。

AR-36：技术信息实现前必须关闭 OQ-08，明确可见角色、授权规则和脱敏字段 allowlist；未关闭时前后端均不提供技术信息入口或 payload。

AR-37：查询成功且结果为 0 条时可以发布可重开的空 `complete_set` 快照并投影“没有待派发漏洞”，但空快照不得创建 ActionDraft 或触发任何 mutation。

AR-38：多页查询任一页失败时只能形成 `partial_page/incomplete` 事实并明确未核验范围；该结果可以只读查看，但不得用于“全部处理”或被模型解释为完整集合。

AR-39：多个查询快照形成一个草案时，程序必须按稳定业务身份确定性去重并保留来源查询摘要；同一对象被指定不同目标或存在重叠归属冲突时必须澄清，禁止 last-write-wins。

AR-40：对象数量超过 capability 或 provider 的 `max_items` 时禁止后台静默分批；系统必须创建用户可见、范围固定且分别确认的草案，或在写入前明确拒绝。

AR-41：确认前重新校验 actor、对象权限、freshness policy 和动作关键事实；当前负责人、处置状态或其他 policy 声明的关键事实变化时，原草案必须失效并要求重新查询、准备和确认。

AR-42：FOBrain 接口 2xx、过滤后实际零对象或请求已接收均不构成条目成功；只有每个对象的回读事实达到目标才可完成，被 provider 过滤或无可证明变化的对象必须进入明确未完成、待核验或待排查状态。

AR-43：G-READ-03 覆盖 `current_user_context`、`my_permissions`、`list_assets_by_ip`、`get_asset_detail`、`list_vulnerabilities_by_ip`、`get_vulnerability_detail` 与 `business_list` 七项能力；既有接口证据只允许继续补 capability、schema、fixture、provider locator resolver、v2 行级事实与只读产品 smoke，不等于产品完成。

AR-44：产品层 `entity_ref` 必须是不可逆、不可解释 provider 身份的 opaque ref；内部 `entity_subject_key` 必须是随机不可变主体，版本化 workspace HMAC 只作为 active/retained fingerprint alias。subject secret 或 locator encryption key 轮换不得改变主体、semantic claim 或 locator lineage；provider locator 只能在受控 resolver 内按 workspace、actor、provider instance、policy version 和有效期解析。详情查询不得把 safe ref 当 provider 真实 ID，也不得把 subject/fingerprint/locator 进入 Product Facts 产品出口、前端契约、模型上下文、audit 或 replay。

AR-45：首次外部 mutation 前，Attempt 必须在同一事务从 `reserved_not_sent` 转为 `sent`，一次性固定 `sent_at`、settle deadline、policy version、lease epoch 和核验调用预算；该 CAS 与 cancel/stop 的 `cancelled_before_send + mutation_count=0 + source-aware claim disposition` 事务互斥。取消首次 claim 写 released_zero_write；取消 failed-retry transferred claim 恢复 predecessor owner/status=failed；sent 后 cancel 返回 `action_already_sent` 并继续只读核验。核验调用先事务预留预算，崩溃不返还预算，重启不得重算 deadline 或次数。

AR-46：Verifier outcome 必须由判别式 evidence 唯一映射：control lost 优先待排查；只有 held+provider accepted+conclusive 真实 readback 目标一致为完成；held+provider rejected+policy-approved conclusive non-application 为明确失败／待排查；其余在 deadline／预算内为待核验、耗尽后为待排查。rejected/unknown+match 不得成功，非 readback 分支不得伪造 observed；后续核验只能读取。

AR-47：操作记录的 actor 是当前凭据解析出的 stable subject；权限未知、拒绝或无法证明时整批 fail closed，不做静默过滤或重新归属。列表使用 `(operation_at desc, run_id desc)`；有 Pending 的流程在记录首次发布时把 `operation_at` 固定为 Pending `published_at`，无 Pending 的流程固定为服务端 decision time，后续审批只写 `decision_at`。首个请求必须把跨 Run 窗口/open union 的完整有序安全列表行物化为独立 immutable OperationListSnapshot；游标绑定服务端生成的 `as_of`、cutoff、opaque operation_list_snapshot_ref 和排序键，绝不复用 Run-local SSE sequence。

### UX Design Requirements

UX-DR-01：实现 DESIGN.md 的颜色、字体、圆角、间距和组件 token，并由项目 CSS 直接消费；使用系统字体、无样式 Radix Primitives 和 lucide-react，不引入 Tailwind、shadcn、外部字体或第三方默认视觉层。

UX-DR-02：首版只实现浅色桌面三栏骨架：232px 会话导航、`minmax(600px, 1fr)` 主工作区、360px 事实与执行栏，最小宽度 1180px；低于宽度时横向滚动或提示，不折叠为移动布局或抽屉。

UX-DR-03：只实现 IA-01 工作台与 IA-02 漏洞明细工作区两张表面；IA-01 中栏包含聊天模式和只读操作记录模式，左、右栏职责稳定，只有中栏切换。返回聊天时恢复来源卡、会话上下文和右栏上下文。

UX-DR-04：聊天是全部查询和动作意图的唯一入口；只读明细和右栏不得准备、确认、取消或执行写入，行选择只更新右栏且不得改变冻结动作范围。

UX-DR-05：会话导航必须显示固定“操作记录”入口、当前会话和待确认状态；切换会话或进入操作记录不得合并、覆盖或重建 result_ref，标题和摘要不得展示内部 ID 或敏感事实。

UX-DR-06：聊天时间线按时序承载消息、查询结果卡、澄清卡、审批卡和执行结果卡；百条对象不得在时间线全量展开，进入明细后返回应定位到来源卡。

UX-DR-07：查询结果卡显示用户可理解的查询序号、摘要、数量和观察时间，0 条与失败使用不同状态；“查看全部”进入该不可变快照的 IA-02，技术 result_ref 不对用户展示。

UX-DR-08：IA-02 按冻结快照每页 100 条展示，显示固定总数、总页数、当前页和加载状态，并支持首页、上一页、下一页、末页及页码跳转；新到漏洞不得改变已打开快照。

UX-DR-09：漏洞事实表采用交互式 grid 和单一 Tab 停靠点，支持方向键、Home/End、PageUp/PageDown、Enter，提供全局 aria-rowcount/aria-rowindex；选择始终只读，不出现复选框式写入范围。

UX-DR-10：漏洞行至少显示行序号、POC、IP、在线状态、业务系统、当前修复负责人和逐条状态；派发核对增加业务系统负责人、运维负责人和发现时间。表格可截断 POC/多 IP，但 accessible name 和右栏必须提供完整安全值。

UX-DR-11：缺失修复负责人固定显示“未分配”，缺失业务系统、业务系统负责人或运维负责人固定显示“未提供”；不得留空、显示猜测值或用风险色替代文字。

UX-DR-12：发现时间统一显示 `YYYY-MM-DD HH:mm:ss（UTC+8）`，IP、时间、计数和稳定数据使用等宽数据字体并支持纵向核对。

UX-DR-13：右栏固定 360px，以“事实”和“执行记录”标签展示当前安全事实；完整 POC 和全部 IP 必须可读。OQ-08 未关闭前隐藏全部技术信息，关闭后也只能显示授权 allowlist 投影。

UX-DR-14：引用缺失、候选不唯一、人员不能唯一确定或必要参数缺失时显示澄清卡；只列安全摘要、数量和时间，明确选择完成前不得创建草案或写入。

UX-DR-15：门禁后的人员选择卡从 FOBrain 人员安全投影人工单选，支持搜索、加载、空、失败、同名消歧和不可选状态；不得 AI 推荐、自由文本姓名、推荐排序或展示内部人员 ID。OQ-02 关闭前不得形成可写入口。

UX-DR-16：门禁后的审批卡只绑定不可变草案，显示来源查询、动作、目标变化/负责人/状态期限、总数和前 5 条，并提供完整冻结集合入口；确认和取消只在卡内，草案变化后必须重新确认。

UX-DR-17：审批提交后立即锁定确认与取消并进入中性执行态；重复提交只投影同一权威终态，聊天中的重复文字不构成确认。

UX-DR-18：执行结果卡先显示完成、待排查、待核验三类计数，再保留逐条结果并可进入 IA-02；混合结果不得整卡染绿/红，不能提供通知、自动处理或未定义的再次执行入口。

UX-DR-19：颜色必须遵守状态语义：蓝色仅查询/导航，琥珀仅尚未写入的待确认，绿色仅回读核验成功，红色仅明确失败；待核验、缺失、执行中、离线和普通事实使用中性表达，所有状态同时有文字或图标。

UX-DR-20：微文案使用中文、冷静且确定：查询空、查询失败、取消、部分失败、待核验、写域未启用和最终待排查均采用 EXPERIENCE.md 的确定性语义，不显示英文内部状态或虚构置信度。

UX-DR-21：写域未过门禁时，用户表达派发、转发、误报或延时时只显示“当前仅支持查询与只读核对，此写操作尚未启用”，不得生成可点击假审批或调用写接口。

UX-DR-22：实现 ST-01～ST-29 的状态投影，重点区分恢复、查询中、有结果、空结果、字段缺失、澄清、无权限、数据不足、外部失败、待确认、取消、执行/回读、混合结算、回读不一致、过期、重复提交、待核验、待排查和历史只读。

UX-DR-23：页面必须具有可命名 `nav`、唯一 `main`、补充 `aside` 和“跳到主内容”；DOM 顺序即焦点顺序，所有交互可用 Tab/Enter/Space，禁止正 tabindex。Esc 不得代替审批取消。

UX-DR-24：动态状态使用节流后的单一 polite status 摘要；最终阻断/失败只播报一次。后台更新不抢焦点；进入 IA-02 聚焦标题，返回聚焦来源卡，错误聚焦安全错误摘要。

UX-DR-25：满足 WCAG 2.2 AA，焦点环与相邻表面非文本对比至少 3:1，文字缩放 200% 时关键查询、澄清、审批、错误和结果不被遮挡；只允许 IA-02 表格自身水平滚动。

UX-DR-26：分页、当前行和滚动位置在表面切换或右栏阅读后恢复；页加载失败保留最后已核验页面并明确“不是空结果”，“重新加载本页”只能读取同一冻结快照。

UX-DR-27：历史查询和执行结果只读恢复必须来自 Product Facts；会话重新进入时按 waiting、过期、已提交、待核验、待排查和终态分流，不从聊天摘要重建范围，也不自动生成重试能力。

UX-DR-28：首版不提供暗色主题、移动断点、触屏专用交互、装饰插画、玻璃拟态、营销页式大留白或未在两份 spine 中命名的通用覆盖层。

UX-DR-29：Epic 1 必须把当前 Inspector 的 evidence、structured、runtime、audit 技术标签替换为目标“事实／执行记录”安全投影；OQ-08 未关闭前技术信息整体隐藏，不得因当前代码已存在技术标签而保留。

UX-DR-30：操作记录模式按操作时间倒序显示动作、对象摘要、三类结果计数和状态；支持稳定游标分页与进入逐条只读结果。时间窗口边界、生命周期并集、权限失败和重新登录行为必须与 AD-26 一致，前端不得合并多页、过滤越权项或重算历史状态。

### FR Coverage Map

FR-01：Epic 1 — 只使用当前 actor 已授权的 FOBrain 事实。

FR-02：Epic 1 — 在列表、明细及未来确认/结果投影中建立统一漏洞事实。

FR-03：Epic 1 — 区分空结果、无权限、数据不足和外部失败。

FR-04：Epic 1 — 精确查询全部新增漏洞并按发现时间倒序展示。

FR-05：Epic 2 — 运营基于统一事实并从 FOBrain 全员列表人工选择接收人。

FR-06：Epic 2 — 只合并人工判定为同一接收人的派发对象。

FR-07：Epic 2 — 派发经确认、逐条写入并按负责人回读核验。

FR-08：Epic 2 — 当前用户可主动转发其可见范围内的漏洞。

FR-09：Epic 2 — 转发使用 FOBrain 全员候选并按接收人拆分草案。

FR-10：Epic 2 — 转发经确认并逐条回读目标负责人。

FR-11：Epic 3 — 延时决定和新期限始终由人工沟通确定。

FR-12：Epic 3 — 管理员一次确认并延时一条漏洞，原因可空。

FR-13：Epic 3 — 仅状态和期限均回读一致时显示延时成功。

FR-14：Epic 3 — 操作人可对自己负责的漏洞作人工误报判断。

FR-15：Epic 3 — 误报确认后写入且仅在状态回读一致时成功。

FR-16：Epic 2、Epic 3 — 两类写域都复用 Epic 1 固定的二次确认协议，并分别证明未确认、取消、拒绝和过期时零写入。

FR-17：Epic 2、Epic 3 — 两类写域都必须逐条投影完成、待核验和待排查；批量负责人变更额外验证部分失败隔离。

FR-18：Epic 2、Epic 3 — 两类写域都不得发送产品通知或把 FOBrain 自身通知作为成功条件。

FR-19：Epic 1 — 读取当前 FOBrain 用户的安全身份上下文。

FR-20：Epic 1 — 读取当前用户的 FOBrain 权限和数据范围。

FR-21：Epic 1 — 按 IP 查询资产及安全在线状态。

FR-22：Epic 1 — 查看单个资产的安全详情。

FR-23：Epic 1 — 按 IP 查询关联漏洞列表。

FR-24：Epic 1 — 查看单个漏洞的安全详情。

FR-25：Epic 1 — 查询授权业务系统列表。

## Epic List

### Epic 1：可信、可恢复的 FOBrain 事实工作台

运营可以在浅色桌面聊天工作台中安全查询 FOBrain 当前用户、权限、资产、漏洞和业务系统，精确查看全部新增漏洞及统一行级事实，并通过冻结快照、百条分页和确定性引用在刷新、断线或重启后继续核对同一结果；所有真实写入保持关闭。

**FRs covered:** FR-01、FR-02、FR-03、FR-04、FR-19、FR-20、FR-21、FR-22、FR-23、FR-24、FR-25。

**实施说明：** 该 Epic 完整承载 M-0～M-5，包括工具链固定、Product Facts/Catalog/Action/SSE v2 唯一目标契约、Greenfield SQLite bootstrap、旧 epoch fail-closed、同源投影、稳定 SSE、Eino Runner/HITL 基础、7 个代表性只读能力、精确新增查询、全员列表和统一漏洞卡片。通用 ActionDraft、Prepare/Confirm、continuation、幂等、verifier ports，以及审批卡、澄清卡、执行结果卡和禁用态通过 schema、fixture、mock 固定，但不连接真实 mutation。必须通过 G-TOOLCHAIN、G-ARCH-V2 以及适用的 G-READ/G-FACT 产品证据；真实写入保持关闭。

### Epic 2：运营人工派发与责任人转发

运营可以基于 Epic 1 冻结的授权漏洞事实，从 FOBrain 全员列表人工选择接收人，将同一接收人的对象合并派发或主动转发；系统在聊天中二次确认，逐条写入、回读负责人，并清楚展示完成、待核验和待排查对象。

**FRs covered:** FR-05、FR-06、FR-07、FR-08、FR-09、FR-10、FR-16、FR-17、FR-18。

**实施说明：** 该 Epic 直接依赖 Epic 1，复用其不可变 ActionDraft、二阶段确认、policy 双检、幂等、verifier、reconcile、安全结果组件和 Story 1.36 操作记录，只增加负责人变更的分组、provider mutation、回读和批量部分失败行为。只有 G-ARCH-V2、G-TOOLCHAIN、G-FACT-01、G-WRITE-01/02、G-SAFE-01 及 OQ-02 的客观证据通过后，真实派发/转发 Story 才可开始。否则只保留契约、fixture、mock 与门禁提示。

### Epic 3：管理员延时与运营误报处置

管理员可以对人工沟通确定的一条漏洞执行修复延时，运营可以对自己负责并已人工判断的漏洞标记误报；两类动作都复用可信确认与核验链路，只有目标状态事实回读一致才显示成功。

**FRs covered:** FR-11、FR-12、FR-13、FR-14、FR-15、FR-16、FR-17、FR-18。

**实施说明：** 该 Epic 直接依赖 Epic 1，不依赖 Epic 2；复用 Epic 1 的统一写控制面、Story 1.36 操作记录和 UX 组件，不创建第二套状态机或确认协议，只增加延时/误报的角色、基数、目标状态、provider mutation 和回读规则。直接延时受 G-WRITE-03 与 OQ-04 阻塞，误报受 G-WRITE-04 阻塞；两者同时继承 G-FACT-01、G-WRITE-01、G-SAFE-01 及 FR-16～FR-18 的共同安全边界。

## Epic 完成与停止规则

1. Epic 1 的 Story 必须按 M-0→M-5 的依赖顺序推进；允许每个切片形成独立机器证据，但只有 READ-01“新增漏洞查询—冻结快照—百条明细—写域禁用提示—零外部写入”端到端通过，且 G-TOOLCHAIN、G-ARCH-V2、适用 G-READ/G-FACT 有效时，Epic 1 才能标记完成。
2. G-FACT-01 所需关键事实无法从目标部署安全获得时，停止统一卡片和后续写域，更新 PRD、schema、fixture、readiness gate 或范围；不得用摘要、硬编码、provider raw 字段或模型推断替代。
3. Epic 2/3 当前只允许创建和实施决策关闭、契约、fixture、mock、门禁取证及安全失败验证的正式 Story。真实 mutation 可以作为不可执行的目标 Story Candidate 保留在完整产品蓝图中；Candidate 不得进入 Sprint、不得交给开发 Agent、不得作为实施授权。只有该 Candidate `ConversionGates` 中列出的全部门禁与决策、环境授权和可回滚样本均有效 PASS，并重新执行 readiness 校验后，才能由 Create Story 生成新的正式 Story 并实施。
4. Epic 2 只有派发与转发在批准环境中分别具备权限、确认后单次写入、逐条回读、部分失败和错误脱敏证据，并能通过 Story 1.36 的固定“操作记录”入口重新找到历史逐条结果与待排查对象时才能完成；Epic 3 只有延时与误报分别具备相应角色/对象权限、目标状态、确认后写入和回读证据，并满足同一历史结果可找回条件时才能完成。
5. `blocked`、`skipped`、mock PASS、接口 2xx、前端控件存在或代码路径存在均不是 Epic 完成信号；机器报告和 readiness gate 验收记录是唯一完成证据。
6. 通用 Action 状态机、确认协议、幂等、reconcile、verifier port 和产品状态投影只能在 Epic 1 建立。Epic 2/3 若需要修改这些公共不变量，必须先 correct-course 并更新 SPEC/ADR，不能在动作 Story 中复制或静默分叉。
7. 生产 Workbench 是否显示审批/确认能力只能来自后端 capability/policy 安全投影；真实写域未授权时，mock 审批控件只能存在于测试环境，产品界面必须显示明确禁用态。
8. Epic 1 必须用机器证据证明当前过渡外壳已被目标能力直接替换且没有兼容分支：ToolResult 持久化完整安全 StructuredResult 值；v2 聚合和 repository 存在；Eino Runner 使用 registry ToolsConfig 与配置化迭代；SSE 使用持久化事实游标；Inspector 遵守 UX-DR-29；CI/本地验收使用 Go 1.26.0/1.26.5 与 Node 24 LTS。旧 adapter、摘要存储、查询时 sequence、`MaxIterations=1` 外壳或现有技术标签均不得计为完成。

## Story 产物分类

| 类型 | 用途 | 当前允许范围 | 是否可进入 Sprint／开发 |
| --- | --- | --- | --- |
| 正式 Story | 当前已获授权且具有完整实现上下文的可执行工作单元 | Epic 1；Epic 2/3 的决策、契约、fixture、mock、门禁取证与安全失败验证 | 只有自身前置门禁满足时可以 |
| Gate/Evidence Story | 关闭 OQ、验证目标部署权限/字段/写入/回读或形成机器验收记录 | Epic 1～3 | 可以，但不得顺带实现被阻塞 mutation |
| Target Story Candidate | 保留完整产品蓝图中的用户价值、目标 AC、依赖门禁和范围边界 | Epic 2/3 的真实派发、转发、延时、误报 | 不可以；只用于追踪和后续转换 |

Candidate 转正式 Story 必须同时满足：适用 readiness gate 有有效 PASS 证据；环境授权和可回滚样本明确；PRD、SPEC、ADR、schema/fixture 与目标 AC 无冲突；重新运行 implementation-readiness 检查；转换动作记录日期、证据和批准人。不得通过删除 `blocked` 文案、修改状态字段或复制 Candidate 文件绕过转换门禁。

## Epic 1：可信、可恢复的 FOBrain 事实工作台

运营可以在浅色桌面聊天工作台中安全查询 FOBrain 当前用户、权限、资产、漏洞和业务系统，精确查看全部新增漏洞及统一行级事实，并通过冻结快照、百条分页和确定性引用在刷新、断线或重启后继续核对同一结果；所有真实写入保持关闭。

### Story 1.1：固定可复现工具链与开工门禁

As a 交付负责人，
I want 项目固定 Go 1.26.x、Node 24 LTS 和开工校验，
So that 后续证据来自一致且受支持的构建环境。

Requirements:
  FR: []
  NFR: [NFR-01]
  AR: [AR-01, AR-29, AR-30]
  Architecture: [M-0]
  UX: []
  Gates: [G-TOOLCHAIN]
  Milestone: [M-0]

**Acceptance Criteria:**

**Given** 新环境或 CI 开始构建
**When** 执行版本、schema、Go、TypeScript、stream、browser、checkpoint 与 SQLite 基线校验
**Then** 只有 Go 1.26.x 和 Node 24 LTS 的完整结果可发布为 G-TOOLCHAIN 证据
**And** 其他版本必须明确失败或标记为非通过证据。

### Story 1.2：固定 StructuredResult 与查询快照契约

As a 运营人员，
I want 每次查询结果都形成不可变安全事实快照，
So that AI、界面和后续操作引用同一份结果。

Requirements:
  FR: [FR-01, FR-03]
  NFR: [NFR-01, NFR-04]
  AR: [AR-05, AR-07, AR-08, AR-37, AR-38]
  Architecture: [AD-03, AD-05]
  UX: [UX-DR-07, UX-DR-20]
  Gates: [G-ARCH-V2]
  Milestone: [M-1]

**Acceptance Criteria:**

**Given** provider 返回成功、空、部分分页或失败结果
**When** Safety Gate 接受并持久化查询结果
**Then** 生成带 schema version、result_ref、内容摘要、coverage 和完整安全行事实的 StructuredResult／QueryResultSnapshot
**And** raw provider payload 不得进入产品契约。

### Story 1.3：固定 Catalog 与 Action 核心契约

As a 产品控制面，
I want 能力声明和动作草案使用版本化契约，
So that 查询能力与未来写动作不依赖工具名硬编码。

Requirements:
  FR: [FR-16, FR-17, FR-18]
  NFR: [NFR-03, NFR-04]
  AR: [AR-05, AR-06, AR-10, AR-13]
  Architecture: [AD-07, AD-09, AD-13]
  UX: [UX-DR-14, UX-DR-16]
  Gates: [G-ARCH-V2]
  Milestone: [M-1]

**Acceptance Criteria:**

**Given** schema 与 fixtures 被生成和校验
**When** 读取 Catalog、ActionDraft、ActionItem、ActionAttempt 与 VerificationEvidence
**Then** capability、policy、verifier、route、credential binding、冻结对象和人工参数均具有版本化字段
**And** 过渡契约被目标契约替换，不生成旧契约 adapter 或 v1/v2 runtime union。

### Story 1.4：固定 Run、Pending、Continuation 与 Verification 契约

As a 操作人，
I want 等待确认、恢复和核验状态具有唯一语义，
So that 刷新或重启后不会重复确认或写入。

Requirements:
  FR: [FR-16, FR-17]
  NFR: [NFR-01, NFR-03]
  AR: [AR-06, AR-11, AR-12, AR-17, AR-46]
  Architecture: [AD-08, AD-10, AD-12, AD-23]
  UX: [UX-DR-17, UX-DR-18, UX-DR-22]
  Gates: [G-ARCH-V2]
  Milestone: [M-1]

**Acceptance Criteria:**

**Given** Run 进入等待、执行、核验或终态
**When** schema、fixture 与状态映射测试运行
**Then** Run 七态、Pending 终态、Attempt phase 与 verifier outcome 只有一个合法映射
**And** 待核验与待排查不被混为同一状态。

### Story 1.5：发布 Action API、OpenAPI 与生成契约

As a Workbench 与自动化调用方，
I want 两个入口共享同一 Action 控制面契约，
So that 确认、取消和结果不会形成两套真相。

Requirements:
  FR: [FR-16, FR-17]
  NFR: [NFR-01, NFR-03]
  AR: [AR-05, AR-11, AR-24, AR-25]
  Architecture: [AD-03, AD-08, AD-16]
  UX: [UX-DR-16, UX-DR-17]
  Gates: [G-ARCH-V2]
  Milestone: [M-1]

**Acceptance Criteria:**

**Given** OpenAPI 3.1 与 JSON Schema 2020-12 输入
**When** 生成 Go／TypeScript 契约并运行映射测试
**Then** PrepareAction、ConfirmAction、取消和 ActionResult 使用同一 Product Facts 类型
**And** 错误统一为脱敏 envelope。

### Story 1.6：固定持久化 SSE 事件契约

As a Workbench 用户，
I want 断线重连后继续接收确定的事实变化，
So that 页面状态不会因连接重建而漂移。

Requirements:
  FR: [FR-02, FR-03]
  NFR: [NFR-01]
  AR: [AR-05, AR-19, AR-24]
  Architecture: [AD-17, AD-24]
  UX: [UX-DR-24, UX-DR-27]
  Gates: [G-ARCH-V2]
  Milestone: [M-1]

**Acceptance Criteria:**

**Given** 客户端提供有效、旧或未知游标
**When** SSE 合约测试重放事实事件
**Then** 有效游标只收到后续完整 upsert／tombstone
**And** 旧或未知游标收到 `view.replaced`，不得依赖进程内消息序号。

### Story 1.7：建立 Greenfield SQLite bootstrap 与启动 readiness

As a 运维负责人，
I want 数据库只在目标 schema epoch 和安全连接条件满足时接流量，
So that Product Facts v2 不被旧结构或静默重建污染。

Requirements:
  FR: []
  NFR: [NFR-01]
  AR: [AR-20, AR-21, AR-22]
  Architecture: [AD-20, AD-22, AD-25]
  UX: []
  Gates: [G-ARCH-V2]
  Milestone: [M-2]

**Acceptance Criteria:**

**Given** 全新数据库或旧 schema epoch 数据库
**When** 启动执行目标 schema bootstrap、epoch、integrity、foreign-key、WAL 和可写探针检查
**Then** 全新数据库直接建立目标结构，重复启动保持一致
**And** 旧 epoch 返回 `unsupported_schema_epoch`、readiness=false，不自动迁移或删库。

### Story 1.8：持久化唯一安全 StructuredResult 事实

As a 运营人员，
I want 系统保存完整安全查询事实而非摘要，
So that 之后能原样恢复同一结果。

Requirements:
  FR: [FR-01, FR-02, FR-03]
  NFR: [NFR-01]
  AR: [AR-07, AR-22]
  Architecture: [AD-03, AD-05, AD-21]
  UX: [UX-DR-07, UX-DR-27]
  Gates: [G-ARCH-V2]
  Milestone: [M-2]

**Acceptance Criteria:**

**Given** StructuredResult candidate 已通过 schema 与 Safety Gate
**When** repository 写入或读取结果
**Then** 保存并恢复完整安全 JSON 值、schema version、result_ref、摘要、大小和时间
**And** 不存在旧摘要 adapter、历史回填或第二事实源。

### Story 1.9：建立 QuerySnapshot 与 Action Repository

As a 执行控制面，
I want 快照和动作聚合具有明确 repository ports，
So that 恢复流程不跨层直接访问 SQLite。

Requirements:
  FR: [FR-02, FR-16]
  NFR: [NFR-01]
  AR: [AR-02, AR-06, AR-08, AR-10]
  Architecture: [AD-04, AD-05]
  UX: [UX-DR-27]
  Gates: [G-ARCH-V2]
  Milestone: [M-2]

**Acceptance Criteria:**

**Given** facts repository interface 与 SQLite adapter
**When** 保存、按 scope 读取或拒绝跨 scope 的 snapshot／draft
**Then** workspace、conversation、actor、版本、有效期和完整 items 被一致保存
**And** facts/product/httpapi 不反向依赖 SQLite 实现。

### Story 1.10：原子提交聚合与 FactEvent

As a 审计与恢复机制，
I want 状态聚合和事实事件在同一事务提交，
So that 任何可见状态都有对应的单调事实序列。

Requirements:
  FR: [FR-17]
  NFR: [NFR-01]
  AR: [AR-06, AR-19]
  Architecture: [AD-04, AD-17]
  UX: [UX-DR-27]
  Gates: [G-ARCH-V2]
  Milestone: [M-2]

**Acceptance Criteria:**

**Given** 任一 v2 聚合状态变化
**When** SQLite transaction 提交
**Then** 聚合与不可变 FactEvent 同时成功并获得 Run 内单调 sequence
**And** 任一写入失败时两者都不发布。

### Story 1.11：持久化 continuation、attempt 与 lease 恢复材料

As a 后台执行器，
I want 进程重启后从持久事实安全接管执行，
So that 旧 owner 不会提交迟到结果。

Requirements:
  FR: [FR-16, FR-17]
  NFR: [NFR-01, NFR-03]
  AR: [AR-12, AR-18, AR-23, AR-45]
  Architecture: [AD-08, AD-11, AD-12, AD-18, AD-23]
  UX: [UX-DR-17, UX-DR-22]
  Gates: [G-ARCH-V2]
  Milestone: [M-2]

**Acceptance Criteria:**

**Given** waiting、running 或 reconciling Run 在任意崩溃点重启
**When** 新 executor 通过 lease epoch 与 fencing token 接管
**Then** continuation、attempt phase、sent_at、deadline 和核验预算保持原值
**And** 旧 token 的提交被拒绝，sent attempt 不得重发。

### Story 1.12：建立后端 Product Facts v2 唯一安全投影

As a Workbench 与 Action API 用户，
I want 所有出口读取同一份安全产品事实，
So that 界面、API、审计和模型上下文保持一致。

Requirements:
  FR: [FR-01, FR-02, FR-03]
  NFR: [NFR-01, NFR-04]
  AR: [AR-22, AR-24]
  Architecture: [AD-03, AD-16, AD-21]
  UX: [UX-DR-11, UX-DR-13, UX-DR-20]
  Gates: [G-ARCH-V2]
  Milestone: [M-3]

**Acceptance Criteria:**

**Given** Product Facts v2 聚合与安全 StructuredResult
**When** product projection 生成 Workbench、Action API、audit 或 replay view
**Then** 所有出口使用同一安全字段和稳定缺失语义
**And** 不允许旧 projection adapter、raw provider、checkpoint、credential 或 locator 字段泄漏。

### Story 1.13：发布持久事实 SSE View 与重连

As a Workbench 用户，
I want 刷新或断线后恢复权威页面视图，
So that 我不会看到重复、回退或丢失的结果。

Requirements:
  FR: [FR-02, FR-03]
  NFR: [NFR-01]
  AR: [AR-19, AR-24]
  Architecture: [AD-17, AD-24]
  UX: [UX-DR-24, UX-DR-26, UX-DR-27]
  Gates: [G-ARCH-V2]
  Milestone: [M-3]

**Acceptance Criteria:**

**Given** 已持久化 snapshot_sequence 与后续 FactEvent
**When** 客户端首次加载或携带游标重连
**Then** 服务端先返回同源 view，再从 snapshot_sequence 后传增量
**And** 未知游标使用完整替换，断线不改变冻结范围。

### Story 1.14：生成模型可用的安全事实上下文

As a 对话用户，
I want AI 只根据已核验安全事实回答和理解指代，
So that 上下文压缩不会导致虚构对象或越权数据。

Requirements:
  FR: [FR-01, FR-02, FR-03]
  NFR: [NFR-01, NFR-03]
  AR: [AR-07, AR-09]
  Architecture: [AD-03, AD-05, AD-06]
  UX: [UX-DR-07, UX-DR-14]
  Gates: [G-ARCH-V2]
  Milestone: [M-3]

**Acceptance Criteria:**

**Given** 一次或多次安全查询结果及对话压缩
**When** LLM context builder 组织事实
**Then** 只注入 Product Facts／StructuredResult 安全投影和用户可理解摘要
**And** raw payload、locator、内部 ID 与未经 resolver 锁定的候选不得进入上下文。

### Story 1.15：实现前端 SSE Reducer 与确定恢复

As a Workbench 用户，
I want 页面只按权威事件更新 server state，
So that 重连和乱序事件不会在浏览器重建第二套事实。

Requirements:
  FR: [FR-02, FR-03]
  NFR: [NFR-01]
  AR: [AR-19, AR-25]
  Architecture: [AD-17, AD-24]
  UX: [UX-DR-05, UX-DR-24, UX-DR-26]
  Gates: [G-ARCH-V2]
  Milestone: [M-3]

**Acceptance Criteria:**

**Given** 完整 view、增量 upsert、tombstone、重复或乱序事件
**When** 前端 reducer 消费事件
**Then** 只接受合法后续 sequence，并可用 `view.replaced` 原子替换
**And** TanStack Query／Zustand 不复制 Product Facts 或 Run 状态机。

### Story 1.16：建立 Capability Registry 与 Eino ToolsConfig

As a Agent 运行时，
I want 从能力注册表发现受控 FOBrain 工具，
So that 工具选择不依赖名称或自然语言关键词硬编码。

Requirements:
  FR: [FR-01]
  NFR: [NFR-02, NFR-03]
  AR: [AR-02, AR-03, AR-13, AR-27]
  Architecture: [AD-01, AD-09, AD-13]
  UX: []
  Gates: [G-ARCH-V2]
  Milestone: [M-4]

**Acceptance Criteria:**

**Given** 版本化 capability metadata 与当前 actor policy
**When** 运行时生成 Eino ToolsConfig
**Then** 只暴露允许的 schema、freshness、max_items 和安全能力描述
**And** provider route、凭据和工具私有字段不暴露给模型。

### Story 1.17：接入 Eino Runner 与受控工具循环

As a 对话用户，
I want AI 能连续查询、解释和澄清，
So that 复杂的资产与漏洞核对不受单轮工具调用限制。

Requirements:
  FR: [FR-01, FR-03]
  NFR: [NFR-02, NFR-03]
  AR: [AR-03, AR-11, AR-13]
  Architecture: [AD-01, AD-02, AD-08]
  UX: [UX-DR-06, UX-DR-14]
  Gates: [G-ARCH-V2]
  Milestone: [M-4]

**Acceptance Criteria:**

**Given** 已注册的只读 capability 与配置化迭代预算
**When** Eino Runner 执行多步工具循环
**Then** 每次工具结果先通过 Safety Gate 形成 StructuredResult 才可继续推理
**And** 模型不能直接改变 Run、Draft、Pending 或 Product Facts 状态。

### Story 1.18：建立 LLM Provider、网络策略与错误脱敏

As a 平台维护者，
I want 模型访问通过独立 provider interface 和受控网络策略，
So that 模型故障不会泄露凭据或污染业务事实。

Requirements:
  FR: [FR-03]
  NFR: [NFR-01, NFR-04]
  AR: [AR-02, AR-27]
  Architecture: [AD-01, AD-15]
  UX: [UX-DR-20]
  Gates: [G-ARCH-V2]
  Milestone: [M-4]

**Acceptance Criteria:**

**Given** mock 或真实 LLM provider 的成功、超时、拒绝与异常响应
**When** llm 边界调用模型
**Then** timeout、URL、模型和预算来自配置，错误以稳定脱敏分类返回
**And** prompt、completion、token、credential 和 raw response 不进入 Product Facts 或产品错误。

### Story 1.19：建立可观测性、脱敏与组合根

As a 平台维护者，
I want 运行时具有可诊断但不泄密的启动与观测链，
So that 问题可定位且分层边界保持清晰。

Requirements:
  FR: []
  NFR: [NFR-01]
  AR: [AR-02, AR-04, AR-27, AR-28]
  Architecture: [AD-01, AD-15, AD-19, AD-22]
  UX: []
  Gates: [G-ARCH-V2]
  Milestone: [M-4]

**Acceptance Criteria:**

**Given** 服务启动、查询、恢复和错误路径
**When** bootstrap 组合依赖并记录 trace、metrics 与日志
**Then** 只有组合根连接实现，import boundary 检查通过
**And** 观测数据仅包含脱敏诊断、latency、token 与预算摘要。

### Story 1.20：建立 PrepareAction、Policy 与不可变草案

As a 操作人，
I want 写意图先形成可核对的冻结草案，
So that AI 不能替我决定对象、人员或关键参数。

Requirements:
  FR: [FR-05, FR-06, FR-08, FR-09, FR-11, FR-14, FR-16]
  NFR: [NFR-02, NFR-03]
  AR: [AR-08, AR-09, AR-10, AR-13, AR-39, AR-40, AR-41]
  Architecture: [AD-06, AD-07, AD-09]
  UX: [UX-DR-14, UX-DR-15, UX-DR-16]
  Gates: [G-ARCH-V2]
  Milestone: [M-4]

**Acceptance Criteria:**

**Given** 已锁定 complete_set 快照、人工参数和当前权限
**When** PrepareAction 校验 scope、freshness、grouping 与 max_items
**Then** 只生成带摘要、版本、有效期和完整 items 的不可变草案
**And** 候选不唯一、事实变化或权限未知时只返回澄清／拒绝且零 mutation。

### Story 1.21：建立 ConfirmAction、审批消费与双 Continuation

As a 操作人，
I want 二次确认被原子消费且可安全恢复，
So that 重复点击或刷新不会产生第二次执行。

Requirements:
  FR: [FR-16]
  NFR: [NFR-01, NFR-03]
  AR: [AR-11, AR-12]
  Architecture: [AD-07, AD-08, AD-12]
  UX: [UX-DR-16, UX-DR-17]
  Gates: [G-ARCH-V2]
  Milestone: [M-4]

**Acceptance Criteria:**

**Given** 未消费、已消费、过期或内容变化的 approval ref
**When** Workbench Eino continuation 或 Action API continuation 确认／取消
**Then** 单一事务消费 ref、复核 actor/policy/digest/version 并固定 Pending 终态
**And** 只有首次合法确认可预留 attempts，其他路径外部 mutation 为零。

### Story 1.22：建立 mutation 幂等与 Attempt 发送预留

As a 执行控制面，
I want 每个对象的写入意图跨恢复保持唯一，
So that 响应丢失或新草案不会导致重复外部写入。

Requirements:
  FR: [FR-16, FR-17]
  NFR: [NFR-01, NFR-03]
  AR: [AR-15, AR-45]
  Architecture: [AD-11, AD-23]
  UX: [UX-DR-17, UX-DR-18]
  Gates: [G-ARCH-V2]
  Milestone: [M-4]

**Acceptance Criteria:**

**Given** 同一语义对象、动作与目标出现在重复请求、恢复或新草案
**When** executor 预留并发送 attempt
**Then** semantic key 使用版本化 canonical payload，Confirm 事务通过唯一 SemanticMutationClaim 保证两个并发普通草案只有一个可预留
**And** prepared/reserved_not_sent/sent/verifying/reconciling/succeeded/manual_attention/failed 的既有 claim 均阻断普通草案，只有绑定 definitive failed item 的人工 retry 可事务转移 claim
**And** `reserved_not_sent` 禁止 sent-only 字段，`reserved_not_sent→sent` 事务才固定 sent_at、deadline、policy、epoch 和预算
**And** cancel/stop 与 sent 使用互斥 phase CAS；取消获胜时写 `cancelled_before_send`、mutation count=0，首次 claim 进入 released_zero_write，failed-retry claim 恢复 predecessor failed；sent 获胜时返回 `action_already_sent` 且只能核验
**And** failed → retry reserved → cancelled_before_send 后，普通 draft 仍被恢复的 failed claim 阻断，新的合法 lineage retry 仍可事务转移。

### Story 1.23：建立 Executor Lease、Fencing 与崩溃接管

As a 系统运行者，
I want 后台执行不依赖浏览器连接且只由当前 owner 推进，
So that 崩溃和并发接管不会产生竞争提交。

Requirements:
  FR: [FR-16, FR-17]
  NFR: [NFR-01, NFR-03]
  AR: [AR-18, AR-23]
  Architecture: [AD-12, AD-18]
  UX: [UX-DR-17, UX-DR-27]
  Gates: [G-ARCH-V2]
  Milestone: [M-4]

**Acceptance Criteria:**

**Given** executor 在确认后、发送前、发送后或核验中崩溃
**When** 新 owner 取得更高 lease epoch
**Then** 仅新 fencing token 可提交后续事实，旧 owner 迟到提交被拒绝
**And** sent item 只进入核验恢复，不重新 mutation
**And** sent 前 cancel/stop 与发送 CAS 的竞态只有一个合法提交，不能留下 reserved orphan、伪造 sent 字段或永久 active claim。

### Story 1.24：建立 Verifier、Reconcile 与 ActionResult

As a 操作人，
I want 每条结果由目标事实回读决定，
So that 接口 2xx 不会被误报为业务成功。

Requirements:
  FR: [FR-07, FR-10, FR-13, FR-15, FR-17]
  NFR: [NFR-03, NFR-04]
  AR: [AR-16, AR-42, AR-45, AR-46]
  Architecture: [AD-10, AD-11, AD-23]
  UX: [UX-DR-18, UX-DR-19, UX-DR-20]
  Gates: [G-ARCH-V2]
  Milestone: [M-4]

**Acceptance Criteria:**

**Given** mutation 响应与零个或多个核验 observed facts
**When** 版本化 verifier 在固定 deadline／预算内比较 expected 与 observed
**Then** 每项由带 evidence_kind、provider acceptance、control certainty、policy/verifier version、call index、relation、conclusiveness 和条件 observed 字段的 VerificationEvidence 唯一映射
**And** control lost 优先 manual_attention；仅 held+accepted+conclusive real-readback match 完成，仅 held+rejected+policy-approved non-application 明确失败
**And** rejected/unknown+match、预算内失租约与无 observed 的明确拒绝 fixture 均有唯一结果，manual_attention 批次 outcome 仍固定为 partial 或 none_succeeded。

### Story 1.25：关闭 G-READ 与 G-FACT 数据证据

As a 产品负责人，
I want 在开发产品查询前确认目标部署能提供所需安全事实，
So that 界面不会靠模型推断或硬编码补齐字段。

Requirements:
  FR: [FR-01, FR-02, FR-03, FR-19, FR-20, FR-21, FR-22, FR-23, FR-24, FR-25]
  NFR: [NFR-01, NFR-02]
  AR: [AR-26, AR-31, AR-32, AR-43, AR-44]
  Architecture: [AD-27]
  UX: [UX-DR-10, UX-DR-11, UX-DR-13]
  Gates: [G-READ-01, G-READ-02, G-READ-03, G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 已批准的九项读取需求和七项代表性能力
**When** 通过 schema、fixture、mock/live assertion 与敏感字段审计核对目标部署
**Then** 每个字段标记 resolved、confirmed empty/no permission 或 source unavailable/unknown，并形成机器证据
**And** 关键事实不可得时 G-FACT-01 保持 BLOCKED，后续统一卡片和写域停止。

### Story 1.26：读取当前用户与权限三态

As a 运营人员，
I want 查看当前 FOBrain 身份和权限范围，
So that 我知道查询事实属于谁以及覆盖什么范围。

Requirements:
  FR: [FR-01, FR-19, FR-20]
  NFR: [NFR-01, NFR-02, NFR-04]
  AR: [AR-14, AR-26, AR-43]
  Architecture: [AD-14, AD-15]
  UX: [UX-DR-11, UX-DR-20]
  Gates: [G-READ-03, G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 当前 workspace 的 FOBrain 凭据
**When** 调用 current_user_context 与 my_permissions
**Then** 安全投影区分 resolved、confirmed_empty/no_permission 与 source_field_unavailable/unknown
**And** token、raw account、policy payload 和推断部门不得展示或记录。

### Story 1.27：按 IP 查询资产与安全详情

As a 运营人员，
I want 按 IP 查看资产列表和单个资产安全详情，
So that 我能核对主机在线状态和资产上下文。

Requirements:
  FR: [FR-21, FR-22]
  NFR: [NFR-01, NFR-02, NFR-04]
  AR: [AR-26, AR-43, AR-44]
  Architecture: [AD-14, AD-27]
  UX: [UX-DR-10, UX-DR-11, UX-DR-13]
  Gates: [G-READ-03, G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 当前用户有权查询的 IP
**When** 查询资产列表并从 opaque entity_ref 打开详情
**Then** 展示安全资产事实和稳定在线状态
**And** resolver 按 `entity_ref→entity_subject_key→ProviderLocator` 唯一链解析；entity_subject_key 随机不可变且只供跨查询身份与 semantic guard 使用
**And** subject-derivation key 与 locator encryption key 轮换前后，同一资产仍命中同一 subject、semantic claim 与 locator lineage，内部材料不进入产品出口。

### Story 1.28：按 IP 查询漏洞与安全详情

As a 运营人员，
I want 按 IP 查看关联漏洞和单个漏洞详情，
So that 我能在处置前核对完整漏洞事实。

Requirements:
  FR: [FR-23, FR-24]
  NFR: [NFR-01, NFR-02, NFR-04]
  AR: [AR-26, AR-43, AR-44]
  Architecture: [AD-14, AD-27]
  UX: [UX-DR-10, UX-DR-11, UX-DR-13]
  Gates: [G-READ-03, G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 当前用户有权查询的 IP 或漏洞安全引用
**When** 查询漏洞列表并打开详情
**Then** 列表与详情使用同一安全漏洞事实语义
**And** resolver 按 `entity_ref→entity_subject_key→ProviderLocator` 唯一链解析，无权限、空结果、字段不可用和外部失败保持不同状态
**And** active/retained fingerprint alias 在 secret 轮换前后必须复用同一漏洞 subject 与 locator lineage。

### Story 1.29：查询全部业务系统安全列表

As a 运营人员，
I want 查看 FOBrain 返回的全部业务系统，
So that 漏洞卡片能够关联真实业务上下文。

Requirements:
  FR: [FR-25]
  NFR: [NFR-01, NFR-02]
  AR: [AR-26, AR-32, AR-43]
  Architecture: [AD-13, AD-14]
  UX: [UX-DR-10, UX-DR-11]
  Gates: [G-READ-03, G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 当前用户凭据
**When** 调用已批准的 `business_list` capability
**Then** 返回经接口内置权限裁剪的全部业务系统安全列表
**And** 业务对象使用随机不可变 entity_subject_key 关联安全事实，subject secret 轮换只新增 fingerprint alias，不得生成新主体、替换为 `my_business_systems` 或暴露 raw payload／locator。

### Story 1.30：精确查询新增漏洞

As a 运营人员，
I want 按发现时间倒序查询全部新增漏洞，
So that 每天能准确看到需要人工处理的对象。

Requirements:
  FR: [FR-03, FR-04]
  NFR: [NFR-02, NFR-04]
  AR: [AR-31, AR-37, AR-38]
  Architecture: [AD-05]
  UX: [UX-DR-07, UX-DR-20]
  Gates: [G-READ-01, G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 查询覆盖全部分页
**When** 用户请求新增漏洞
**Then** 结果按发现时间倒序形成 complete_set 快照，0 条显示“没有待派发漏洞”
**And** 任一页失败只形成可读 incomplete 事实，不得声称完整或用于“全部”动作。

### Story 1.31：读取 FOBrain 全部人员列表

As a 运营人员，
I want 从 FOBrain 全部人员安全列表中查找人员，
So that 未来能人工指定负责人而不是由 AI 猜测。

Requirements:
  FR: [FR-05, FR-09]
  NFR: [NFR-02, NFR-03]
  AR: [AR-31, AR-33]
  Architecture: [AD-14]
  UX: [UX-DR-15]
  Gates: [G-READ-02, OQ-02]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** FOBrain 返回全部人员安全投影
**When** 用户搜索、筛选、键盘选择或遇到同名人员
**Then** 按已固定规则展示、消歧并只允许人工单选
**And** AI 不推荐、静默默认或接受自由文本伪造人员 ID。

### Story 1.32：建立统一漏洞行事实与多源冻结快照

As a 运营人员，
I want 每条漏洞统一呈现 POC、IP、在线状态、业务系统及负责人，
So that 我能基于同一事实进行判断。

Requirements:
  FR: [FR-02, FR-05]
  NFR: [NFR-01, NFR-04]
  AR: [AR-08, AR-39, AR-41]
  Architecture: [AD-03, AD-05]
  UX: [UX-DR-10, UX-DR-11, UX-DR-12]
  Gates: [G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 漏洞、资产、业务系统与人员等多个安全来源
**When** 产品组装统一漏洞行事实
**Then** 使用稳定业务身份确定性去重并保留每字段来源和缺失语义
**And** 冲突不得 last-write-wins，必须显示数据不足或要求澄清。

### Story 1.33：建立结果索引、引用解析与恢复

As a 对话用户，
I want 用“刚才第 2 次结果”等自然语言准确引用已有查询，
So that 查询多次或上下文压缩后仍锁定正确对象集合。

Requirements:
  FR: [FR-02, FR-03]
  NFR: [NFR-01, NFR-02]
  AR: [AR-08, AR-09, AR-44]
  Architecture: [AD-05, AD-06, AD-27]
  UX: [UX-DR-07, UX-DR-14, UX-DR-27]
  Gates: [G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 同一 scope 内存在一个、多个、过期或越权的 result_ref 候选
**When** resolver 解释用户指代
**Then** 只返回 resolved、clarification_required 或 invalid
**And** 唯一候选可锁定，多候选必须人工澄清，跨 scope／过期引用必须拒绝。

### Story 1.34：交付浅色桌面聊天 Workbench

As a 运营人员，
I want 在统一桌面聊天工作台查询并查看安全结果摘要，
So that 日常核对从自然语言入口完成。

Requirements:
  FR: [FR-01, FR-02, FR-03, FR-04]
  NFR: [NFR-01, NFR-03]
  AR: [AR-25]
  Architecture: [AD-16, AD-17]
  UX: [UX-DR-01, UX-DR-02, UX-DR-03, UX-DR-04, UX-DR-05, UX-DR-06, UX-DR-21, UX-DR-23, UX-DR-25, UX-DR-28, UX-DR-29]
  Gates: [G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 浅色桌面视口与只读能力
**When** 用户聊天查询、切换会话或查看结果卡
**Then** 三栏布局、焦点语义、空／失败／澄清状态和“事实／执行记录”安全投影符合 UX 契约
**And** 真实写域未启用时只显示明确禁用态，外部 mutation 为零。

### Story 1.35：交付百条分页漏洞明细与右侧事实区

As a 运营人员，
I want 在专门工作区分页核对全部漏洞和单条事实，
So that 百条结果也能高效、无歧义地检查。

Requirements:
  FR: [FR-02, FR-04]
  NFR: [NFR-01, NFR-04]
  AR: [AR-25]
  Architecture: [AD-16, AD-24]
  UX: [UX-DR-08, UX-DR-09, UX-DR-10, UX-DR-11, UX-DR-12, UX-DR-13, UX-DR-24, UX-DR-25, UX-DR-26]
  Gates: [G-FACT-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 一个冻结的多页漏洞快照
**When** 用户分页、键盘浏览、选择行、打开右栏或返回聊天
**Then** 总数、页码、行序、选中行、滚动位置和来源卡保持稳定
**And** 页加载失败保留最后核验页面并明确不是空结果。

### Story 1.36：交付跨会话操作记录只读入口

As a 运营人员，
I want 重新登录后从固定入口找到近期操作和仍待处理结果，
So that 我不依赖聊天摘要恢复事实。

Requirements:
  FR: [FR-17]
  NFR: [NFR-01, NFR-02]
  AR: [AR-23, AR-35, AR-47]
  Architecture: [AD-21, AD-24, AD-25, AD-26]
  UX: [UX-DR-03, UX-DR-05, UX-DR-18, UX-DR-27, UX-DR-30]
  Gates: [G-ARCH-V2]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** 当前凭据主体具有已授权历史操作和 open lifecycle 对象
**When** 用户从 IA-01 固定“操作记录”入口按稳定游标浏览
**Then** operation_at 在首次发布时固定；服务端生成 as_of/cutoff，并在一个事务中把跨 Run 的精确六个 Asia/Shanghai 日历月窗口与 open lifecycle 并集物化为 immutable OperationListSnapshot
**And** 快照冻结全部有序安全列表行、记录版本与 total，opaque tamper-evident cursor 绑定 actor/workspace、filters、as_of、cutoff、operation_list_snapshot_ref 和 last keys，明确不使用 Run-local SSE sequence
**And** 两个不同 Run 在翻页期间变化不得造成重复／遗漏；快照缺失／TTL 过期／篡改、history policy ref 缺失／损坏、权限未知／拒绝均整批 fail closed
**And** 页面只读且不生成草案、重试或重新归属。

### Story 1.37：完成 READ-01 端到端验收

As a 产品负责人，
I want 用机器证据验收完整只读主流程，
So that 只有真正可用的事实工作台才能进入 Sprint 完成状态。

Requirements:
  FR: [FR-01, FR-02, FR-03, FR-04, FR-19, FR-20, FR-21, FR-22, FR-23, FR-24, FR-25]
  NFR: [NFR-01, NFR-02, NFR-04]
  AR: [AR-29, AR-30, AR-31, AR-32]
  Architecture: [M-0, M-1, M-2, M-3, M-4, M-5]
  UX: [UX-DR-01, UX-DR-30]
  Gates: [G-TOOLCHAIN, G-ARCH-V2, G-READ-01, G-READ-02, G-READ-03, G-FACT-01, READ-01]
  Milestone: [M-5]

**Acceptance Criteria:**

**Given** M-0～M-5 前置 Story 与适用门禁均有有效证据
**When** 执行“新增漏洞查询—冻结快照—百条明细—引用恢复—写域禁用提示—零外部写入”E2E
**Then** schema、Go、TypeScript、SSE、browser、checkpoint、SQLite 和安全断言全部通过并记录证据
**And** blocked、skipped、mock PASS、接口 2xx 或控件存在均不得计为 READ-01 PASS。

## Archived Source Bundles（非 Sprint 工作项）

以下 Source Bundle 保留原始详细 AC，供正式 Story 实施时追溯；它们不是可执行 Story，不得被 Sprint Planning 或开发 Agent 当作任务。

> 2026-07-15 Greenfield 决策覆盖以下历史 bundle 中的兼容表述：不得实现 v1 只读、v1/v2 union、additive v1 migration、旧数据库自动迁移或旧投影 adapter。若历史 AC 与 `docs/adr/2026-07-15-greenfield-no-legacy-compatibility.md` 冲突，以 Greenfield ADR 和上方正式 Story 为准。

### Source Bundle SB-1.1：固定可复现工具链与验收基线

As a 项目维护者，
I want 开发、CI 和发布环境统一使用受支持的 Go 1.26.x 与 Node 24 LTS，
So that 后续 Product Facts v2 开发和验收具有可复现的统一环境。

**Acceptance Criteria:**

**AC-1.1.1**

**Given** Story 尚未开始实施
**When** 执行 `docs/pre-development-validation.md`
**Then** 形成 M-0 开工记录，明确影响目录、目标工具链、验收命令和 `G-TOOLCHAIN` 当前状态
**And** 未完成验证前不得修改工具链配置。

**AC-1.1.2**

**Given** 仓库当前使用 Go 1.23.12、Node 25.8.1，且声明范围未精确固定
**When** 完成工具链配置
**Then** Go、Node、npm 的版本文件、CI 和项目声明保持一致
**And** Go 必须属于 `1.26.x`，Node 必须属于 `24.x LTS`
**And** 精确 patch 版本必须在机器可读配置和验收记录中一致。

**AC-1.1.3**

**Given** 使用全新环境检出仓库
**When** 按项目说明安装工具链并执行 `npm ci`、Go 依赖解析
**Then** 安装可复现且不会意外修改 lockfile
**And** 不顺带升级 Eino、SQLite、Vite、Vitest、Radix 或其他业务依赖。

**AC-1.1.4**

**Given** 使用 Go 1.23 或 Node 25 等非目标环境
**When** 执行 CI 或工具链门禁
**Then** 检查必须明确失败并显示安全、可操作的版本错误
**And** 该结果不得记录为 `G-TOOLCHAIN PASS`。

**AC-1.1.5**

**Given** 目标工具链安装完成
**When** 执行 schema、contract、OpenAPI、Go、TypeScript、stream 和桌面 browser 现有基线
**Then** 所有当前适用检查通过
**And** 未到阶段的检查明确标记 blocked/skipped，不计为通过。

**AC-1.1.6**

**Given** 全部基线检查完成
**When** 生成验收记录
**Then** 记录精确工具版本、命令、结果和未覆盖风险
**And** 同步 `G-TOOLCHAIN`、实施计划及相关开发说明
**And** 不记录 token、凭据或本地敏感路径。

### Source Bundle SB-1.2：建立安全结果与查询快照 v2 契约

As a 产品中心运营，
I want 每次安全查询都形成范围、顺序和逐行事实确定的版本化快照，
So that 刷新、分页或上下文压缩后仍能引用同一批对象。

**Acceptance Criteria:**

**AC-1.2.1**

**Given** Story 1.1 已完成且 M-1 尚未开始
**When** 执行开发前验证
**Then** 明确本 Story 只影响安全 StructuredResult、QueryResultSnapshot、SnapshotItem、引用解析结果及其 schema、fixture、Go/TypeScript 契约
**And** 不加入 ActionDraft、数据库 bootstrap 或真实 writer。

**AC-1.2.2**

**Given** 工具结果通过 Safety Gate
**When** 生成安全 StructuredResultValue
**Then** 契约包含 schema version、稳定 `result_ref`、完整 `safe_payload`、SHA-256 内容摘要、字节数和创建时间
**And** 内容摘要规范化算法具有确定性测试
**And** raw provider payload、凭据和内部 checkpoint 引用无法通过契约或安全测试。

**AC-1.2.3**

**Given** 一次查询包含零条、一页或多页对象
**When** 发布 QueryResultSnapshot
**Then** 契约固定 workspace、conversation、actor、来源 Run/ToolResult、查询语义、coverage、稳定排序、总数、采集时间、有效期及有序 items
**And** 每个 item 固定安全业务身份、位置、事实 schema、完整安全事实和摘要
**And** 成功空结果可表示为零项 `complete_set`。

**AC-1.2.4**

**Given** 查询分页失败或未完整加载
**When** 创建快照事实
**Then** coverage 只能为 `partial_page` 或 `incomplete`
**And** 契约能够区分完整集合、部分页面和不完整结果
**And** 不使用数量或摘要伪装完整对象集合。

**AC-1.2.5**

**Given** 系统需要解析历史查询引用
**When** 生成引用解析结果
**Then** 契约只允许 `resolved`、`clarification_required` 或 `invalid`
**And** resolved 只引用一个安全 snapshot ID
**And** clarification 只包含安全候选摘要
**And** invalid 使用稳定原因码，不暴露内部 ID 或 provider 数据。

**AC-1.2.6**

**Given** 新增 v2 schema 和 fixtures
**When** 执行 schema 与 contract 校验
**Then** 所有 schema 使用 JSON Schema 2020-12、稳定 `$id`、明确 `schema_version` 和 `additionalProperties: false`
**And** fixtures 至少覆盖完整多页、成功空集、部分分页、字段缺失、越权引用、过期引用和不安全材料拒绝
**And** 全部 fixture 登记到 manifest。

**AC-1.2.7**

**Given** Go 与前端需要消费相同契约
**When** 生成或映射领域类型和 TypeScript 类型
**Then** 前端只通过 `contracts/generated.ts` 消费契约
**And** Go 类型、schema 和 fixtures 具有映射测试
**And** 不创建手写平行 DTO 或按 provider 字段硬编码的类型。

**AC-1.2.8**

**Given** 当前仓库仍有过渡性 Product Facts 实现
**When** 完成本 Story
**Then** 目标 schema、fixtures 和 generated types 不生成旧实现兼容 union
**And** database writer、ActionDraft 和真实 mutation 仍未启用
**And** schema、contract、Go facts 测试和 TypeScript typecheck 通过并形成验收记录。

### Source Bundle SB-1.3：建立 Action、Catalog、Run 与 SSE v2 契约

As a 产品中心运营，
I want 所有未来写动作使用同一套草案、确认、执行和结果契约，
So that 聊天与 Action API 都不能绕过人工确认或形成不同事实。

**Acceptance Criteria:**

**AC-1.3.1**

**Given** Story 1.2 已完成
**When** 执行 M-1 本 Story 的开发前验证
**Then** 范围只包含 Product Facts v2 组合、ActionDraft/Item/Attempt、Catalog v2、Run/Pending/Continuation、VerificationEvidence、Action/Resume/Result/SSE v2 schema、fixtures、OpenAPI 和生成契约
**And** 不实现 SQLite bootstrap、Eino resume、provider mutation 或真实 writer。

**AC-1.3.2**

**Given** 系统准备一个修改动作
**When** 使用 ActionDraft v2 契约表达该动作
**Then** 契约固定 actor、capability/policy/verifier/provider route/credential binding 版本、源 snapshot、完整 items、人工参数、有效期和状态
**And** `draft_digest` 使用 RFC 8785 JCS 与 SHA-256 覆盖全部可确认材料
**And** 每个 ActionItem 固定来源 item、expected facts、目标参数、group key 和 mutation key
**And** ActionAttempt 只保存安全接受摘要、lease epoch 和 verification evidence 引用，不保存 raw request/response。

**AC-1.3.3**

**Given** 当前过渡 Catalog 不包含完整写域约束
**When** 定义 Catalog v2
**Then** capability metadata 可以版本化声明 freshness、verification、max_items、group key、recipient source、actor role、object permission 和 confirmation projection
**And** 未知或缺失策略能够被后续实现 fail closed
**And** schema 不包含按具体 FOBrain 工具名称硬编码的业务分支。

**AC-1.3.4**

**Given** Run、Pending、Draft、Item 和 Attempt 存在不同生命周期
**When** 定义 Product Facts v2 状态契约
**Then** Run 仍只允许 `created/running/waiting/succeeded/failed/cancelled/stopped` 七态
**And** waiting、executing、verifying、reconciling、manual_attention 等状态按 Architecture Spine 的唯一映射归属正确聚合
**And** 一个 Draft 只对应一个 PendingInteraction 和一个一次性 approval ref。

**AC-1.3.5**

**Given** Workbench 与 Action API 必须复用相同控制面
**When** 定义 Action/Resume/Result v2 传输契约和 OpenAPI
**Then** prepare 与 resume 只发布目标版本契约，不生成旧接口判别联合
**And** prepare 返回 draft ID/version/digest 与 approval waiting 投影
**And** confirm 提交 resume ref、decision、client request ID 和 draft ID/version/digest
**And** 相同 client request digest 可返回既有结果，不同 digest 可以表达 conflict
**And** 写域只接受 v2 facts。

**AC-1.3.6**

**Given** 客户端需要从当前 view 接续 SSE
**When** 定义 view 与 stream event v2 契约
**Then** view 包含 `snapshot_sequence`
**And** 事件包含稳定 event ID、持久化 sequence、完整实体 upsert 或显式 tombstone
**And** 数组默认 replace，只有 schema 声明的 keyed collection 才允许按 key 合并
**And** 契约能够表达旧或未知游标触发的 `view.replaced`。

**AC-1.3.7**

**Given** 新增联合契约和 fixtures
**When** 执行 schema、contract 和 OpenAPI 校验
**Then** fixtures 至少覆盖等待确认、取消零写入、摘要冲突、草案过期、部分成功、待核验、manual attention、重复请求、乱序事件和旧游标替换
**And** 目标 generated TypeScript types 与 Go 映射测试均通过
**And** 所有 fixture 登记到 manifest。

**AC-1.3.8**

**Given** 本 Story 完成
**When** 检查实施状态和验收记录
**Then** 过渡 Product Facts 不进入目标运行时，旧 adapter 与 dual write 不存在
**And** database writer、真实审批恢复和 provider mutation 仍保持关闭
**And** `G-ARCH-V2` 只记录 M-1 契约切片证据，不得提前标记整体 PASS。

### Source Bundle SB-1.4：建立 Greenfield SQLite bootstrap 与启动 readiness

As a 项目维护者，
I want SQLite schema、连接和启动检查具有统一且可验证的安全策略，
So that 数据库不满足目标 epoch 或恢复条件时服务不会接收业务流量。

**Acceptance Criteria:**

**AC-1.4.1**

**Given** Story 1.3 已完成且 M-2 尚未开始
**When** 执行开发前验证
**Then** 明确本 Story 只建立 Greenfield schema bootstrap、统一连接策略、备份恢复、健康/readiness 和启动检查
**And** StructuredResult 与 v2 聚合表按 Story 1.5、1.6 的实际需要再创建
**And** 不提前创建未被本 Story 使用的业务表。

**AC-1.4.2**

**Given** facts 与 checkpoint 需要访问同一个 SQLite 文件
**When** bootstrap 创建数据库依赖
**Then** 两者复用同一 database handle 和 connection policy
**And** `journal_mode=WAL` 的返回值必须验证为 `wal`
**And** 每个连接保证 `foreign_keys=ON` 和配置化 `busy_timeout`
**And** pool 上限、事务模式和超时来自受控配置，不在业务代码硬编码。

**AC-1.4.3**

**Given** 数据库是全新文件或包含非目标 schema epoch
**When** 执行目标 schema bootstrap 与 epoch 校验
**Then** 全新数据库按固定顺序和事务边界创建目标结构
**And** 重复启动不会重复修改或破坏目标结构
**And** 非目标 epoch 返回 `unsupported_schema_epoch`、readiness=false，不自动迁移或删库。

**AC-1.4.4**

**Given** 操作者准备显式重建非目标数据库
**When** 执行重建前备份
**Then** WAL 一致性材料包含在受控备份范围
**And** 备份必须恢复到临时数据库并通过 integrity、foreign-key 和 schema-version 检查后才记为有效
**And** 备份、日志和验收记录不包含凭据或未投影业务数据。

**AC-1.4.5**

**Given** 后端进程启动
**When** bootstrap 准备监听业务端口
**Then** 先完成目标 schema bootstrap/epoch、WAL/foreign-key/busy-timeout 验证、integrity check、可写探针和当前可用的恢复扫描
**And** 任一步失败时 `/ready` 为 false 且业务流量不被接收
**And** `/health` 只表达进程存活
**And** FOBrain 或 LLM 不可用只使相关 capability degraded/fail closed，不使安全事实读取服务错误退出 readiness。

**AC-1.4.6**

**Given** schema bootstrap/epoch、连接初始化、备份恢复或可写探针发生故障
**When** 执行启动失败注入测试
**Then** 服务产生脱敏且可诊断的错误
**And** 不创建 v2 Run、不调用外部 provider、不修改 readiness 为 true
**And** 修复环境后重新启动能够安全重试。

**AC-1.4.7**

**Given** 本 Story 实现完成
**When** 在全新数据库、非目标 epoch 数据库和恢复备份数据库上执行测试
**Then** bootstrap 重复性、legacy epoch 拒绝、WAL、每连接外键、busy timeout、锁竞争、integrity、readiness 和 backup restore 测试全部通过
**And** Go race 测试覆盖连接和启动边界
**And** 实施计划、配置说明、schema bootstrap 记录和 `G-ARCH-V2` 的 M-2 部分证据同步更新
**And** v2 writer 与真实 mutation 仍保持关闭。

### Source Bundle SB-1.5：持久化完整安全 StructuredResult

As a 产品中心运营，
I want 查询得到的完整安全事实能够在进程重启后恢复，
So that 历史查询不会退化成无法核对对象的数量摘要。

**Acceptance Criteria:**

**AC-1.5.1**

**Given** Story 1.4 已完成
**When** 执行本 Story 的开发前验证
**Then** 范围只包含 StructuredResultValue 目标表、repository port/实现、完整安全值和安全测试
**And** 不创建 QueryResultSnapshot、ActionDraft、FactEvent 或真实 writer。

**AC-1.5.2**

**Given** 数据库包含非目标 schema epoch 或摘要型过渡 ToolResult
**When** bootstrap 校验 schema epoch
**Then** 返回 `unsupported_schema_epoch` 且不创建 Product Facts v2
**And** 不猜测生成缺失的 `safe_payload`、digest 或完整列表
**And** 不启用兼容读取 adapter。

**AC-1.5.3**

**Given** 一个工具候选结果已通过 schema 和 Safety Gate
**When** repository 持久化 StructuredResultValue
**Then** schema version、result_ref、完整 safe payload、内容摘要、字节数和创建时间在同一事实事务内保存
**And** presentation、safe summary 和列表展示均不是独立事实源
**And** raw provider request/response 不进入任何表。

**AC-1.5.4**

**Given** safe payload 超过配置化大小限制、摘要不匹配、schema 未登记或包含禁止材料
**When** 尝试持久化
**Then** repository 在事务提交前拒绝该结果
**And** 不留下半条 ToolResult 或孤立 payload
**And** 错误只包含稳定 code、脱敏摘要和 correlation ref。

**AC-1.5.5**

**Given** 安全 StructuredResultValue 已成功持久化
**When** 服务重启后按 result_ref 读取
**Then** 返回的 schema、payload、digest、大小和创建时间与写入值一致
**And** repository 重新验证 digest 和安全边界
**And** 损坏、缺失或不一致的数据 fail closed，不返回伪造的完整事实。

**AC-1.5.6**

**Given** Workbench、Action API、replay、audit 或模型上下文请求事实
**When** 当前数据库不是目标 Product Facts v2 epoch
**Then** 服务保持 readiness=false 并拒绝产品读取
**And** 不创建 QuerySnapshot、ActionDraft 或兼容投影
**And** 操作者必须显式重建环境后重新查询。

**AC-1.5.7**

**Given** 本 Story 实现完成
**When** 执行 fresh database、legacy epoch rejection、重启 round-trip、payload 上限、unsafe material、digest corruption 和事务回滚测试
**Then** schema、fixture、Go facts/store 测试和 race 测试全部通过
**And** 日志、备份和验收记录中不存在 raw payload、凭据或个人敏感信息
**And** v2 Run writer 与真实 mutation 仍保持关闭
**And** 实施计划、facts contract、Greenfield schema 记录和 M-2 验收证据同步更新。

### Source Bundle SB-1.6：持久化 v2 聚合、FactEvent 与恢复材料

As a 产品中心运营，
I want 查询、草案、确认和执行事实具有统一且可恢复的持久化结构，
So that 刷新、重启和并发恢复不会产生第二套真相或重复状态。

**Acceptance Criteria:**

**AC-1.6.1**

**Given** Story 1.5 已完成
**When** 执行本 Story 的开发前验证
**Then** 范围只包含 v2 聚合目标 schema、repository ports/实现、事务不变量、FactEvent sequence、retention 和恢复扫描材料
**And** 不实现 Eino tool loop、HTTP 产品入口、provider mutation 或业务 verifier。

**AC-1.6.2**

**Given** v2 契约已通过 M-1 校验
**When** 在 fresh database 创建目标持久化结构
**Then** 只创建本 Story repository 测试实际使用的 QuerySnapshot/Item、ActionDraft/Item/Attempt、Pending/ContinuationRef、FactEvent、VerificationEvidence、idempotency 和 lease/fencing 持久化结构
**And** 外键、唯一约束、版本字段、状态约束和索引与 schema 一致
**And** bootstrap 可在全新数据库重复执行，非目标 epoch 稳定拒绝。

**AC-1.6.3**

**Given** execution 提交一个 v2 聚合状态迁移
**When** repository 接受命令结果
**Then** 当前聚合状态与对应不可变 FactEvent 在同一 SQLite 事务提交
**And** repository 可以直接读取当前状态，不需要全量事件回放
**And** FactEvent 只用于来源、审计、诊断、replay 和增量投影，不宣称完整 Event Sourcing。

**AC-1.6.4**

**Given** 系统保存 QuerySnapshot、Draft 或 Attempt
**When** 执行 repository 校验
**Then** workspace、conversation、actor、版本、digest、来源引用和状态满足领域约束
**And** 一个 Draft 最多存在一个有效 Pending 和一个一次性 approval ref
**And** mutation key、client request key 和 attempt lineage 的唯一性由数据库与 repository 双重保护
**And** 跨 scope、缺失来源或 digest 不一致的写入在事务前被拒绝。

**AC-1.6.5**

**Given** 同一 Run 发生并发状态迁移
**When** repository 追加 FactEvent
**Then** sequence 在 Run 内持久化且严格单调递增
**And** 聚合版本冲突、重复请求和旧 lease epoch 不会覆盖较新状态
**And** 迟到提交返回稳定 conflict/reconcile 结果，不产生重复事实。

**AC-1.6.6**

**Given** 服务启动或 executor 接管未终结 Run
**When** 执行恢复扫描
**Then** waiting、running、reconciling 及其引用的 snapshot、draft、continuation 和 idempotency 材料保持可恢复
**And** 孤儿 staged continuation、可见 waiting continuation 和未知外部写入被区分记录
**And** 当前 Story 只形成恢复分类和安全待处理事实，不执行外部重试。

**AC-1.6.7**

**Given** cleanup 或 retention 任务运行
**When** 检查 Product Facts 引用关系
**Then** waiting/running/reconciling 使用中的材料不得删除
**And** 首版不自动物理删除 terminal Product Facts 或 audit
**And** snapshot/draft 只能逻辑过期
**And** 幂等记录至少保留到 audit/draft retention 与 reconcile window 的较晚者。

**AC-1.6.8**

**Given** 本 Story 完成
**When** 执行 fresh bootstrap、legacy epoch rejection、事务回滚、并发 sequence、重复请求、lease fencing、active retention、恢复扫描和 backup restore 测试
**Then** facts/store/sqlite 测试与 race 测试通过
**And** 不存在历史 v1 Run adapter 或 v1/v2 dual write
**And** v2 repository 仅在测试或受控内部配置中启用，生产 HTTP/Agent 入口仍不创建 v2 Run
**And** 实施计划、schema bootstrap 文档、恢复规则和 M-2 验收证据同步更新。

### Source Bundle SB-1.7：建立同源产品投影与稳定 SSE

As a 产品中心运营，  
I want 工作台、执行结果、历史回放和审计始终显示同一组事实，  
So that 断线重连或页面刷新不会丢失、重复或改变业务结果。

**Acceptance Criteria:**

**AC-1.7.1**

**Given** Story 1.6 已完成且 M-3 尚未开始  
**When** 执行开发前验证  
**Then** 范围只包含 `product` 唯一 v2 投影、HTTP view/SSE 边界、模型安全上下文和前端 stream reducer  
**And** 不实现 Workbench 视觉布局、Eino tool loop、provider mutation 或业务操作按钮。

**AC-1.7.2**

**Given** 同一 v2 Run 包含工具、快照、草案、执行、核验和审计事实  
**When** Workbench、ActionResult、replay、audit 或模型上下文请求投影  
**Then** 所有出口从同一 repository snapshot 和 Product Facts 映射  
**And** 不分别维护状态、列表或成功判断  
**And** raw provider payload、内部 checkpoint、credential ref 和未授权技术字段不进入任何出口。

**AC-1.7.3**

**Given** 客户端首次读取当前 Run view  
**When** 服务生成安全投影  
**Then** aggregate view 与当时最大持久化 sequence 在一致读取边界内返回  
**And** view 携带 `snapshot_sequence=N`  
**And** N 之后产生的事实不会被首次读取和 SSE 订阅之间的竞态遗漏。

**AC-1.7.4**

**Given** 客户端以 `Last-Event-ID` 或 snapshot sequence 订阅  
**When** SSE 发送后续事实  
**Then** 只发送 sequence 大于客户端高水位的持久化事件  
**And** event ID 在重新查询或重启后保持稳定  
**And** 重复订阅不会重新生成不同 event ID  
**And** 旧或未知游标收到 `view.replaced` 与新高水位后再继续增量。

**AC-1.7.5**

**Given** stream 包含 upsert、tombstone、数组 replace、keyed collection、重复或乱序事件  
**When** 前端 reducer 应用事件  
**Then** `sequence <= applied_sequence` 的事件被忽略  
**And** 完整实体 upsert、显式 tombstone 和数组替换遵循 v2 schema  
**And** 只有 schema 声明的 keyed collection 按 key 合并  
**And** 最终 view 与服务端同源投影一致，不产生幽灵 item 或状态回退。

**AC-1.7.6**

**Given** 数据库不是目标 Product Facts v2 epoch  
**When** 请求 view、replay 或 audit  
**Then** 服务返回稳定不可用错误且 readiness=false  
**And** 不创建兼容投影、snapshot sequence、StructuredResult 或 ActionDraft  
**And** 操作者获得显式重建指引。

**AC-1.7.7**

**Given** 模型需要多轮上下文  
**When** Context Assembler 构造输入  
**Then** 只使用安全 Product Facts 摘要、result refs 和当前授权范围  
**And** 具体行事实按 ref 从 repository 重载并重新校验 scope  
**And** 不从聊天摘要、前端状态或 provider raw data 重建对象集合。

**AC-1.7.8**

**Given** 本 Story 实现完成  
**When** 执行 snapshot-to-stream 竞态、重复/乱序、断线重连、过旧游标、tombstone、数组替换、legacy epoch 拒绝和安全泄漏测试  
**Then** product/httpapi/stream/reducer 测试及 TypeScript typecheck 全部通过  
**And** 当前查询时重新编号的过渡 SSE 逻辑被删除而非保留 fallback  
**And** 实施计划、SSE 契约、上下文文档和 M-3 验收证据同步更新  
**And** 生产 v2 Run 入口和真实 mutation 仍保持关闭。

### Source Bundle SB-1.8：接入 Eino Registry、ToolsConfig 与受控工具循环

As a 产品中心运营，  
I want Agent 能根据自然语言调用已登记且获准的工具，  
So that 我可以通过对话完成查询，而不依赖硬编码工具分支。

**Acceptance Criteria:**

**AC-1.8.1**

**Given** Story 1.7 已完成  
**When** 执行开发前验证  
**Then** 范围只包含 `execution`、`capabilities`、`llm`、`observability` 和 bootstrap 接线  
**And** 不实现真实 FOBrain 写入、Action 确认、业务操作按钮或完整 Workbench 页面。

**AC-1.8.2**

**Given** Capability Registry 中登记了可用能力  
**When** 创建 Eino Agent 的 `ToolsConfig`  
**Then** 工具定义来自 registry 的 capability metadata、输入输出 schema 和 policy metadata  
**And** 不按工具名称、自然语言关键词或 provider 私有字段硬编码选择逻辑  
**And** 未登记、未启用或不属于当前 scope 的能力不会提供给模型。

**AC-1.8.3**

**Given** 用户请求可能需要零次、一次或多次工具调用  
**When** Eino Agent 执行请求  
**Then** 使用 Eino `ChatModelAgent`、Runner 和原生工具调用机制完成受控循环  
**And** 最大迭代次数、工具调用次数、token、时间和并发预算均来自配置  
**And** 不另外实现一套与 Eino 竞争的手写 Agent runtime。

**AC-1.8.4**

**Given** 模型生成工具调用  
**When** execution 准备执行工具  
**Then** 输入必须通过登记 schema、actor scope、capability policy 和预算校验  
**And** 未知工具、非法参数、越权调用或禁用能力被稳定拒绝  
**And** 首版生产环境只开放已通过门禁的查询能力和 mock 能力。

**AC-1.8.5**

**Given** 模型、工具或 Runner 产生执行事件  
**When** EventMapper 接收事件  
**Then** 事件被转换为 Product Facts v2 的运行状态、ToolCall 和 StructuredResult 引用  
**And** Eino 内部对象、provider raw payload 和未投影错误不进入产品层  
**And** Workbench、replay、audit 和模型上下文继续使用同一组事实。

**AC-1.8.6**

**Given** 客户端断开、请求超时、用户取消或进程恢复  
**When** Run 仍未终结  
**Then** Run 生命周期独立于 HTTP/SSE 连接  
**And** timeout、cancel、budget exhausted 和执行失败映射为契约规定的稳定状态与脱敏错误  
**And** 不因客户端重连重新执行已经完成的工具调用。

**AC-1.8.7**

**Given** 模型接收上下文或返回文本与工具调用  
**When** LLM adapter 处理输入输出  
**Then** 输入只包含授权后的安全 Product Facts 投影  
**And** 输出被视为不可信并经过 schema、安全和能力边界校验  
**And** prompt injection、伪造 result_ref、越权工具调用和敏感信息泄漏不能绕过 registry 与 policy。

**AC-1.8.8**

**Given** 本 Story 实现完成  
**When** 使用 mock model 和 mock tools 执行测试  
**Then** 覆盖零工具、单工具、多工具、未知工具、非法参数、预算耗尽、超时、取消、断线重连和 prompt injection 场景  
**And** Go tests、race tests、schema/fixture 校验及相关边界测试全部通过  
**And** 实施计划、Eino runtime 说明、配置说明和 M-4 验收证据同步更新  
**And** 真实 FOBrain mutation 与 Action continuation 仍保持关闭。

### Source Bundle SB-1.9：建立不可变 Prepare/Confirm 与双入口 Continuation

As a 产品中心运营，  
I want 聊天与 Action API 都通过同一套准备和确认流程，  
So that 任何写动作都不能绕过人工确认或替换已确认对象。

**Acceptance Criteria:**

**AC-1.9.1**

**Given** Story 1.8 已完成  
**When** 执行开发前验证  
**Then** 范围只包含 `PrepareAction`、`ConfirmAction`、双 continuation backend、mock action executor 和对应 HTTP/Eino 接线  
**And** 不接入真实 FOBrain mutation、业务动作 verifier 或生产写能力。

**AC-1.9.2**

**Given** 用户指定有效的冻结查询结果、动作和人工参数  
**When** 执行 `PrepareAction`  
**Then** 系统复核 actor、scope、capability policy、对象权限、完整集合、`max_items` 和动作参数  
**And** 创建不可变 ActionDraft、唯一 PendingInteraction 和一次性 approval ref  
**And** digest 覆盖 actor、版本引用、snapshot、全部对象事实、人工参数和过期时间。

**AC-1.9.3**

**Given** PrepareAction 分别来自聊天和确定性 Action API  
**When** 系统建立等待确认状态  
**Then** 聊天路径使用真实 Eino staged checkpoint/interrupt  
**And** Action API 使用项目 continuation record，不创建伪 Eino checkpoint  
**And** 两者通过统一 `ContinuationRef` port 投影相同的 Run、Draft、Pending 和 FactEvent。

**AC-1.9.4**

**Given** ActionDraft 已创建  
**When** 产品层生成确认摘要  
**Then** 展示来源查询、动作、人工指定参数、对象总数和前5条统一行级事实  
**And** 提供查看完整冻结集合的入口  
**And** AI 不推荐或决定接收人、误报、是否延时、新期限及确认决定。

**AC-1.9.5**

**Given** 用户提交确认  
**When** 执行 `ConfirmAction`  
**Then** 系统重新校验 actor、approval ref、draft ID、version、digest、policy、权限、freshness 和关键事实  
**And** 在同一事务中消费 approval ref、固定 Pending 终态、批准 Draft、预留 Attempts 并追加 FactEvent  
**And** 只有事务提交成功后才允许 resume 或调用 mock executor。

**AC-1.9.6**

**Given** 用户取消、拒绝、未确认、确认过期或确认材料不匹配  
**When** 系统处理该决定  
**Then** mock mutation 调用次数为零  
**And** Draft、Pending、Run 和 audit/replay 显示明确且同源的结果  
**And** 对象范围、人工参数或关键事实变化时必须创建新版本并重新确认。

**AC-1.9.7**

**Given** continuation staging、Pending 发布或 Confirm 事务期间发生崩溃、重启或重复提交  
**When** 系统恢复该 Run  
**Then** 可见 waiting continuation 保持可恢复且不得清理  
**And** 未发布的孤儿 staged continuation 可安全清理  
**And** 同一 approval ref 只能成功消费一次，重复确认投影同一权威结果。

**AC-1.9.8**

**Given** 本 Story 实现完成  
**When** 分别通过聊天和 Action API 对同一草案执行安全 mock  
**Then** 两条路径产生一致的 policy、Draft、Pending、Attempt、FactEvent 和产品投影  
**And** 测试覆盖确认、取消、拒绝、过期、摘要冲突、actor 冲突、事实变化、各崩溃点、重启恢复和重复确认  
**And** Go tests、race tests、schema/fixture/OpenAPI 与生成契约校验全部通过  
**And** 实施计划、continuation 说明和 M-4 验收证据同步更新  
**And** 所有真实 FOBrain mutation 仍保持关闭。

### Source Bundle SB-1.10：建立幂等执行、租约接管与回读核验

As a 产品中心运营，  
I want 已确认的操作在重试、断线和服务重启后仍能可靠执行并核验，  
So that 系统不会重复修改数据或把未知结果错误显示为成功。

**Acceptance Criteria:**

**AC-1.10.1**

**Given** Story 1.9 已完成  
**When** 执行开发前验证  
**Then** 范围只包含 mock mutation adapter、幂等控制、executor lease/fencing、reconcile、Verifier port 和结果投影  
**And** 不连接真实 FOBrain mutation 或声明任何真实写能力可用。

**AC-1.10.2**

**Given** ActionItem 已通过确认并预留 Attempt  
**When** executor 准备执行  
**Then** 系统根据草案、动作、对象和版本化能力生成稳定 mutation key  
**And** 相同确认、相同 client request、技术重试和进程恢复复用同一幂等记录  
**And** 相同请求键但不同 digest 必须返回 conflict，不得覆盖既有操作。

**AC-1.10.3**

**Given** 多个 executor 可能同时发现同一待执行 Attempt  
**When** executor 获取或续约执行租约  
**Then** 同一时刻只有持有有效 epoch 与 fencing token 的 executor 可以推进  
**And** 旧租约持有者不能覆盖较新的状态或核验结果  
**And** 外部调用后失去租约时进入未知结果核对，不得直接重新写入。

**AC-1.10.4**

**Given** mock provider 接收一个已确认操作  
**When** Attempt 执行和恢复  
**Then** 状态按契约在 `reserved_not_sent`、`sent`、`verifying`、`reconciling` 和终态之间迁移  
**And** 每次迁移与 FactEvent、幂等材料在受控事务中保存  
**And** Run 继续只使用既有七态，不新增竞争状态。

**AC-1.10.5**

**Given** mock provider 返回成功响应  
**When** 系统判断业务操作结果  
**Then** provider 接受或 HTTP 2xx 只表示可以进入核验  
**And** Verifier 通过独立只读回读获取 VerificationEvidence 并比较预期事实  
**And** 只有回读事实与目标一致时 ActionItem 才能显示成功  
**And** 空数据、过滤结果或仅有成功状态码不能作为成功证据。

**AC-1.10.6**

**Given** provider 超时、连接中断、响应丢失或返回结果不确定  
**When** 系统无法证明 mutation 是否发生  
**Then** 不盲目重复调用 mutation  
**And** 进入 `reconciling` 并按配置化策略执行只读回读  
**And** 超过核对预算仍无法证明时进入 `manual_attention`，保留脱敏诊断与证据引用。

**AC-1.10.7**

**Given** 批量 mock 操作包含成功、失败、待核验和人工排查条目  
**When** 生成 ActionResult  
**Then** 每条对象保持独立 Attempt、VerificationEvidence 和最终状态  
**And** 已成功条目不会被后续失败覆盖或重复执行  
**And** 产品结果明确区分 `all_succeeded`、`partial`、`none_succeeded` 和 `blocked`。

**AC-1.10.8**

**Given** 本 Story 实现完成  
**When** 执行并发接管、重复确认、重复请求、迟到提交、响应丢失、最终一致回读、核验不一致、预算耗尽和进程重启测试  
**Then** mock mutation 的实际调用次数符合幂等预期  
**And** Go tests、race tests、fixture、状态机和安全投影测试全部通过  
**And** 实施计划、恢复规则、Verifier 契约和 M-4 验收证据同步更新  
**And** 所有真实 FOBrain mutation 仍保持关闭。

### Source Bundle SB-1.11：关闭 G-READ/G-FACT 数据可得性证据

As a 产品中心研发负责人，  
I want 用可复现的只读验证证明 FOBrain 数据能够形成完整安全事实，  
So that 后续能力不会建立在抽样、推断或不完整数据之上。

**Acceptance Criteria:**

**AC-1.11.1**

**Given** Story 1.10 已完成，当前 G-READ 的目标部署 API 已初步通过而产品能力仍阻塞  
**When** 执行开发前验证  
**Then** 范围只包含授权目标部署的只读取证、fixture、数据映射验证和门禁记录  
**And** 不调用任何派发、转发、延时、误报或其他 mutation 接口。

**AC-1.11.2**

**Given** 使用本机忽略提交的凭据配置访问目标部署  
**When** 验证精确新增漏洞查询  
**Then** 状态为新增的全部授权漏洞被完整分页读取  
**And** 全局结果按发现时间倒序  
**And** 空结果、权限失败、网络失败和 provider 失败能够明确区分  
**And** 机器证据记录页数、数量、顺序断言和安全摘要，不保存真实实体或凭据。

**AC-1.11.3**

**Given** 使用同一 actor 和授权范围  
**When** 验证 FOBrain 人员列表  
**Then** 全部可见人员被完整分页读取  
**And** 人员具有稳定业务标识和显示名称  
**And** 空结果、重复人员、分页重叠和接口失败被明确识别  
**And** 不通过手工名单或模型推断补齐人员。

**AC-1.11.4**

**Given** 新增漏洞源记录的行级字段可能缺失  
**When** 执行受控资产和业务关联查询  
**Then** 每条事实保留可证明的 POC、IP、在线状态、业务系统和当前修复负责人  
**And** 可用时关联业务系统负责人和运维负责人  
**And** 缺失、冲突和无关联结果使用契约定义的明确状态，不被猜测或静默覆盖。

**AC-1.11.5**

**Given** 同一漏洞或 IP 从多个只读来源获得事实  
**When** evidence mapper 合并数据  
**Then** 合并规则基于登记的 source identity、schema 和确定性优先级  
**And** 分页重复可安全去重，事实冲突必须显式报告  
**And** raw provider 字段不进入 Product Facts、模型上下文或证据报告。

**AC-1.11.6**

**Given** G-FACT-01 要求的关键事实无法从目标部署安全获得或稳定关联  
**When** 生成门禁结论  
**Then** G-FACT-01 保持 `BLOCKED` 并列明缺失字段、影响对象和复现证据  
**And** 后续依赖统一漏洞事实的能力停止实施  
**And** 不得用摘要、硬编码、模型推断或抽样成功替代门禁。

**AC-1.11.7**

**Given** 全量读取和事实关联验证通过  
**When** 生成验收记录  
**Then** 记录目标部署标识摘要、代码版本、schema/fixture 版本、执行时间、数量、分页、字段覆盖、缺失率、冲突率和脱敏命令  
**And** G-READ-01、G-READ-02 与 G-FACT-01 的“数据可得性”结论按客观证据更新  
**And** 在产品 capability 尚未实现前不得将其产品状态标记为 PASS。

**AC-1.11.8**

**Given** 本 Story 完成  
**When** 执行 fixture、分页、排序、去重、字段缺失、关联冲突、安全泄漏和授权 live-read smoke  
**Then** 机器报告与人工可读证据结论一致  
**And** `implementation-readiness-gate.md`、FOBrain 工具矩阵和 M-5 验收记录同步更新  
**And** 所有 G-WRITE、G-SAFE 和真实 mutation 继续保持阻塞。

### Source Bundle SB-1.12：识别当前用户与 FOBrain 权限范围

As a 产品中心运营，  
I want Agent 准确识别我的 FOBrain 身份和权限范围，  
So that 后续查询与操作只能使用我实际获准访问的数据。

**Acceptance Criteria:**

**AC-1.12.1**

**Given** Story 1.11 的数据可得性结论允许继续  
**When** 执行开发前验证  
**Then** 范围只包含 `current_user_context`、`my_permissions`、actor 绑定和安全投影  
**And** 不实现资产、漏洞、业务系统查询或任何真实 mutation。

**AC-1.12.2**

**Given** workspace 配置了有效 FOBrain 凭据  
**When** 调用当前用户能力  
**Then** StructuredResult 返回稳定用户业务标识、显示名称、部门和角色安全摘要  
**And** `actor_subject_id` 由 provider instance 与稳定用户业务标识组成  
**And** 显示名称、账号文本和 credential fingerprint 不参与身份比较。

**AC-1.12.3**

**Given** 当前 actor 已解析  
**When** 调用我的权限范围能力  
**Then** 返回 FOBrain 提供的权限与数据范围安全摘要  
**And** 空权限作为合法空态展示，不等同于 provider 失败  
**And** 产品不扩大、推断或重写 FOBrain 接口内置的数据权限。

**AC-1.12.4**

**Given** Run、QuerySnapshot、ActionDraft、policy decision 或 audit 需要记录 actor  
**When** 保存 Product Facts  
**Then** 统一绑定同一 `actor_subject_id` 和 provider instance  
**And** credential fingerprint 只作为脱敏审计上下文  
**And** token、账号 raw payload、完整权限策略和 credential ref 不进入产品投影。

**AC-1.12.5**

**Given** Eino Agent 需要选择身份或权限能力  
**When** 构造 ToolsConfig 并执行工具  
**Then** 工具来自 Capability Registry 的 metadata 和 schema  
**And** 不根据用户名、自然语言关键词或 provider 私有字段硬编码路由  
**And** StructuredResult 通过 Safety Gate 后才进入 Product Facts 和模型上下文。

**AC-1.12.6**

**Given** 用户身份缺少稳定业务标识、身份冲突、凭据失效、无权限或 provider 调用失败  
**When** 系统解析 actor  
**Then** 未认证、无权限、合法空态、数据不足和外部失败具有不同稳定结果  
**And** 无法证明唯一 actor 时写域 fail closed  
**And** 模型不得根据显示名称、聊天内容或历史结果猜测当前用户。

**AC-1.12.7**

**Given** 用户改名、凭据轮换、同名人员、共享凭据或 workspace 切换  
**When** 系统重新解析身份  
**Then** 稳定业务标识相同的凭据轮换仍对应同一 actor  
**And** 同名或改名不会改变身份比较结果  
**And** workspace 或 provider instance 不同的 actor 不得复用 scope、result ref 或草案  
**And** 共享凭据不得被宣称为多人独立审计。

**AC-1.12.8**

**Given** 本 Story 实现完成  
**When** 执行 mock、fixture 和授权 live-read smoke  
**Then** 覆盖正常身份、空权限、无权限、稳定标识缺失、重名、改名、凭据轮换、跨 workspace 和敏感信息泄漏  
**And** schema、StructuredResult、Capability Registry、Go tests 和安全投影测试全部通过  
**And** FOBrain 工具矩阵、实施计划及 M-5 验收证据同步更新  
**And** 所有真实 mutation 继续保持关闭。

### Source Bundle SB-1.13：按 IP 查询资产并查看资产详情

As a 产品中心运营，  
I want 通过 IP 查询资产并核对具体资产详情，  
So that 我可以确认漏洞关联主机的在线状态和业务归属。

**Acceptance Criteria:**

**AC-1.13.1**

**Given** Story 1.12 已完成  
**When** 执行开发前验证  
**Then** 范围只包含 `list_assets_by_ip`、`get_asset_detail` 及其安全 StructuredResult  
**And** 不实现漏洞查询、业务系统列表、统一漏洞快照或任何 mutation。

**AC-1.13.2**

**Given** 用户提供一个合法 IP  
**When** 调用按 IP 查询资产能力  
**Then** 系统使用经过 schema 校验的 IP 精确查询  
**And** 不把 IP 转换为模糊关键词或其他查询条件  
**And** 非法 IP 在 provider 调用前返回稳定输入错误。

**AC-1.13.3**

**Given** FOBrain 返回一条或多条资产  
**When** provider adapter 生成 StructuredResult  
**Then** 完整分页结果包含稳定资产标识、IP、在线状态、业务系统及安全资产摘要  
**And** 可用时包含资产重要程度、负责人和运维负责人  
**And** 缺失字段保持明确缺失态，不由模型补造。

**AC-1.13.4**

**Given** 查询没有匹配资产、用户无权限或 provider 失败  
**When** 生成产品结果  
**Then** “未找到资产”“无权限”“数据不足”和“外部查询失败”分别展示  
**And** 零结果显示为 0，不被描述为调用失败  
**And** 接口 2xx 但返回空集合时不得伪造资产。

**AC-1.13.5**

**Given** 用户从资产查询结果请求查看详情  
**When** 调用资产详情能力  
**Then** 使用 Product Facts 中保存的真实 provider asset identity  
**And** `result_ref`、安全 item ref 或前端行号不得冒充真实 `asset_id`  
**And** asset identity 不属于当前 actor、workspace 或来源结果时拒绝查询。

**AC-1.13.6**

**Given** 前一查询只有一个资产或包含多个资产  
**When** 用户通过上下文请求“查看这个资产详情”  
**Then** 唯一对象可以通过持久化 result ref 和 item identity 确定  
**And** 存在多个候选时返回可选择的安全澄清项  
**And** 不依赖聊天摘要或模型记忆猜测资产。

**AC-1.13.7**

**Given** 资产列表或详情已通过 Safety Gate  
**When** 保存和投影结果  
**Then** 完整安全 StructuredResult 可在重启后按 result ref 恢复  
**And** Workbench、replay、audit 和模型上下文消费同一安全事实  
**And** raw provider payload、credential、内部账号和未授权技术字段不得离开 provider 边界。

**AC-1.13.8**

**Given** 本 Story 实现完成  
**When** 执行 fixture、mock 和授权 live-read smoke  
**Then** 覆盖单资产、多资产、空结果、非法 IP、分页、字段缺失、在线状态、详情身份混淆、跨 scope 和 provider 失败  
**And** schema、Capability Registry、StructuredResult、Go tests 和安全泄漏测试全部通过  
**And** FOBrain 工具矩阵、实施计划及 M-5 验收证据同步更新  
**And** 所有真实 mutation 继续保持关闭。

### Source Bundle SB-1.14：按 IP 查询漏洞并查看漏洞详情

As a 产品中心运营，  
I want 通过 IP 查询漏洞并查看具体漏洞详情，  
So that 我可以核对某台主机上的漏洞事实和当前处置状态。

**Acceptance Criteria:**

**AC-1.14.1**

**Given** Story 1.13 已完成  
**When** 执行开发前验证  
**Then** 范围只包含 `list_vulnerabilities_by_ip`、`get_vulnerability_detail` 及其安全 StructuredResult  
**And** 不实现精确新增队列、统一漏洞快照、写动作或 Workbench 页面。

**AC-1.14.2**

**Given** 用户提供合法 IP 以及可选严重度、状态和分页条件  
**When** 调用按 IP 查询漏洞能力  
**Then** 所有输入先通过登记 schema 校验  
**And** IP 使用精确查询，不转换为模糊关键词  
**And** 非法筛选值在 provider 调用前返回稳定输入错误。

**AC-1.14.3**

**Given** 用户刚查询过一个资产结果  
**When** 用户要求“查询这个 IP 的漏洞”  
**Then** 系统从持久化 result ref 和选定资产事实获取 IP  
**And** 多个 IP 候选时返回安全澄清项  
**And** 不从聊天摘要或模型记忆猜测 IP。

**AC-1.14.4**

**Given** FOBrain 返回漏洞列表  
**When** provider adapter 生成 StructuredResult  
**Then** 完整分页结果包含稳定漏洞标识、POC或漏洞名称、IP、严重度、状态、发现时间和当前修复负责人  
**And** 每条结果必须与精确查询的 IP 一致  
**And** 不一致记录被标记为数据冲突，不静默过滤后宣称查询完整  
**And** 缺失字段保持明确缺失态。

**AC-1.14.5**

**Given** 查询没有漏洞、用户无权限、数据不完整或 provider 失败  
**When** 生成产品结果  
**Then** “没有相关漏洞”“无权限”“数据不足”和“外部查询失败”分别展示  
**And** 零结果显示为 0  
**And** HTTP 2xx、过滤后为空或缺少对象列表不能作为存在漏洞的证据。

**AC-1.14.6**

**Given** 用户从漏洞结果请求查看详情  
**When** 调用漏洞详情能力  
**Then** 使用 Product Facts 中保存的真实 provider vulnerability identity  
**And** `result_ref`、安全 item ref、CVE、POC 名称或前端行号不得冒充真实 `vulnerability_id`  
**And** identity 不属于当前 actor、workspace 或来源结果时拒绝查询。

**AC-1.14.7**

**Given** 漏洞列表或详情通过 Safety Gate  
**When** 保存和投影结果  
**Then** 完整安全 StructuredResult 可在重启后恢复  
**And** Workbench、replay、audit 和模型上下文消费同一事实  
**And** raw provider payload、请求参数、credential 和未授权漏洞技术材料不得离开 provider 边界。

**AC-1.14.8**

**Given** 本 Story 实现完成  
**When** 执行 fixture、mock 和授权 live-read smoke  
**Then** 覆盖单条、多条、空结果、非法筛选、分页、字段缺失、IP 不一致、上下文 IP 复用、详情身份混淆、跨 scope 和 provider 失败  
**And** schema、Capability Registry、StructuredResult、Go tests 和安全泄漏测试全部通过  
**And** FOBrain 工具矩阵、实施计划及 M-5 验收证据同步更新  
**And** 所有真实 mutation 继续保持关闭。

### Source Bundle SB-1.15：查询全部业务系统

As a 产品中心运营，  
I want 查看权限范围内的业务系统列表，  
So that 我可以确认资产和漏洞所属的业务范围。

**Acceptance Criteria:**

**AC-1.15.1**

**Given** Story 1.14 已完成  
**When** 执行开发前验证  
**Then** 范围只包含 `business_list`、分页读取和安全 StructuredResult  
**And** 不实现业务风险汇总、业务资产统计、漏洞快照或任何 mutation。

**AC-1.15.2**

**Given** 用户要求查看业务系统且没有提供筛选条件  
**When** 调用业务系统列表能力  
**Then** 系统直接读取当前 actor 权限范围内的全部业务系统  
**And** 不强制用户补充负责人、关键词或重要程度  
**And** 可选关键词、负责人和分页条件必须通过登记 schema 校验。

**AC-1.15.3**

**Given** FOBrain 返回多页业务系统  
**When** provider adapter 生成 StructuredResult  
**Then** 完整分页结果包含稳定业务系统标识、名称和安全业务摘要  
**And** 可用时包含重要程度、业务负责人和运维负责人  
**And** 资产数、漏洞数等统计只有在来源明确提供时才展示  
**And** 缺失的重要程度保持为空，不被自动归类。

**AC-1.15.4**

**Given** 不同页面出现重复标识、同名系统或字段冲突  
**When** 系统合并分页结果  
**Then** 稳定业务标识相同的记录按确定性规则去重  
**And** 相同名称但不同标识的业务系统保持为不同对象  
**And** 负责人、重要程度或其他事实冲突必须显式报告，不静默覆盖。

**AC-1.15.5**

**Given** 查询结果为空、用户无权限、分页不完整或 provider 失败  
**When** 生成产品结果  
**Then** “没有业务系统”“无权限”“结果不完整”和“外部查询失败”分别展示  
**And** 空集合显示总数 0  
**And** 分页中途失败不得把已读取部分宣称为完整结果。

**AC-1.15.6**

**Given** 用户后续引用某个业务系统  
**When** 系统解析业务系统对象  
**Then** 使用持久化 result ref 和稳定业务系统 identity  
**And** 同名多候选时返回安全澄清项  
**And** 不根据聊天摘要、名称相似度或模型推断选择对象。

**AC-1.15.7**

**Given** 业务系统结果通过 Safety Gate  
**When** 保存和投影结果  
**Then** 完整安全 StructuredResult 可在重启后恢复  
**And** Workbench、replay、audit 和模型上下文消费同一事实  
**And** raw provider payload、credential、内部组织字段和未授权人员信息不得离开 provider 边界。

**AC-1.15.8**

**Given** 本 Story 实现完成  
**When** 执行 fixture、mock 和授权 live-read smoke  
**Then** 覆盖无筛选全量读取、筛选、空结果、分页、重复标识、同名系统、字段缺失、事实冲突、无权限和 provider 失败  
**And** schema、Capability Registry、StructuredResult、Go tests 和安全泄漏测试全部通过  
**And** FOBrain 工具矩阵、实施计划及 M-5 验收证据同步更新  
**And** 所有真实 mutation 继续保持关闭。

### Source Bundle SB-1.16：精确查询新增漏洞与 FOBrain 全员列表

As a 产品中心运营，  
I want 查看全部新增漏洞并取得完整人员列表，  
So that 我可以准确识别待处理对象并自行选择负责人。

**Acceptance Criteria:**

**AC-1.16.1**

**Given** Story 1.15 已完成  
**When** 执行开发前验证  
**Then** 范围只包含精确新增漏洞查询、FOBrain 全员列表及其 Catalog、schema、StructuredResult 和安全投影  
**And** 不实现统一漏洞快照、负责人推荐、动作草案或任何 mutation。

**AC-1.16.2**

**Given** 用户要求查询新增漏洞  
**When** Agent 调用登记的精确新增查询能力  
**Then** provider adapter 使用 FOBrain 已验证的“新增”业务状态精确查询  
**And** “新增”不被解释为今天、最近若干小时或模型推断的时间范围  
**And** 当前 actor 权限内的全部结果被完整分页读取。

**AC-1.16.3**

**Given** 新增漏洞跨越多个分页  
**When** 系统形成完整结果  
**Then** 全局按发现时间倒序排列  
**And** 相同发现时间使用稳定对象 identity 作为确定性次序  
**And** 分页重复被安全去重  
**And** 分页中断、总数不一致或游标循环时结果标记为不完整，不得形成完整集合声明。

**AC-1.16.4**

**Given** 当前没有新增漏洞、用户无权限或 provider 查询失败  
**When** 生成产品结果  
**Then** 空集合显示“没有待派发漏洞”且总数为 0  
**And** “无权限”“结果不完整”和“外部查询失败”分别展示  
**And** 空结果文案来自确定性产品投影，不依赖大模型自由生成。

**AC-1.16.5**

**Given** 用户需要选择人员  
**When** Agent 调用登记的全员列表能力  
**Then** 完整分页读取 FOBrain 对当前 actor 可见的全部人员  
**And** 不按部门、业务系统、角色或模型判断预过滤  
**And** 每个人员包含稳定业务标识、显示名称及契约允许的安全摘要  
**And** 不允许用手工姓名列表替代。

**AC-1.16.6**

**Given** 人员分页包含重复标识、同名人员、改名记录或缺少稳定标识  
**When** 系统形成全员列表  
**Then** 相同稳定标识按确定性规则去重  
**And** 同名不同标识保持为不同人员  
**And** 缺少稳定业务标识的记录不可作为未来动作接收人  
**And** 不通过姓名相似度或模型推断合并人员。

**AC-1.16.7**

**Given** 新增漏洞结果或人员列表通过 Safety Gate  
**When** 保存和投影 StructuredResult  
**Then** 完整安全集合可在重启后按 result ref 恢复  
**And** Workbench、replay、audit 和模型上下文消费同一事实  
**And** raw provider payload、token、账号 raw 字段和未授权人员信息不得离开 provider 边界。

**AC-1.16.8**

**Given** 本 Story 实现完成  
**When** 执行 fixture、mock 和授权 live-read smoke  
**Then** 覆盖多页新增漏洞、全局排序、空结果、分页中断、重复对象、多页人员、同名人员、缺少稳定标识、无权限和 provider 失败  
**And** schema、Catalog、StructuredResult、Go tests 和安全泄漏测试全部通过  
**And** 只有产品 smoke 与既有目标部署证据同时有效时，G-READ-01/02 的产品 capability 才更新为 PASS  
**And** G-FACT、所有 G-WRITE、G-SAFE 和真实 mutation 继续保持阻塞。

### Source Bundle SB-1.17：建立统一漏洞事实、冻结快照与结果指代

As a 产品中心运营，  
I want 每次漏洞查询都形成可核对、可引用且不会变化的事实集合，  
So that 后续查看或操作能够准确使用我指定的查询结果。

**Acceptance Criteria:**

**AC-1.17.1**

**Given** Story 1.16 已完成且 G-READ 产品能力已通过  
**When** 执行开发前验证  
**Then** 范围只包含统一漏洞事实映射、QueryResultSnapshot、Conversation Result Index 和确定性 resolver  
**And** 不实现 Workbench 视觉页面、真实 ActionDraft 或任何 mutation。

**AC-1.17.2**

**Given** 漏洞、资产、业务系统和人员只读结果可用  
**When** 系统形成统一漏洞行级事实  
**Then** 每条记录包含来源可证明的 POC、IP、在线状态、业务系统和当前修复负责人  
**And** 可用时包含业务系统负责人和运维负责人  
**And** 关联基于稳定 identity、登记 schema 和确定性规则  
**And** 缺失、无关联或冲突保持明确状态，不由模型推断。

**AC-1.17.3**

**Given** 一次漏洞查询已经完成  
**When** 创建 QueryResultSnapshot  
**Then** 在同一事实事务中冻结 actor、workspace、conversation、查询条件、source result refs、完整对象集合、顺序、事实版本、创建时间和 digest  
**And** 完整空结果可以形成对象数为 0 的有效快照  
**And** 分页不完整、来源冲突未解决或数据校验失败的结果必须标记为不完整，不得冒充可操作完整集合。

**AC-1.17.4**

**Given** 同一会话已经执行五次或更多查询  
**When** 用户引用“刚才那个结果”“第二次查询”或带条件描述的某次结果  
**Then** 模型只生成结构化引用意图，程序通过持久化 Conversation Result Index 解析候选  
**And** 唯一候选绑定对应 result ref 和 snapshot  
**And** 多个候选时展示查询时间、条件、数量和安全摘要供用户选择  
**And** 不让模型直接重建或替换对象集合。

**AC-1.17.5**

**Given** 会话被压缩、浏览器断线或服务重启  
**When** 用户继续引用历史查询结果  
**Then** resolver 从 repository 重载 result index、snapshot 和 scope  
**And** 不依赖聊天摘要或前端内存恢复事实  
**And** 跨 actor、workspace、conversation 或 provider instance 的 result ref 被拒绝  
**And** 过期或已失效快照只能只读查看，不能静默刷新后继续使用。

**AC-1.17.6**

**Given** 用户选择全部对象、部分对象或某个具体漏洞  
**When** resolver 锁定对象范围  
**Then** 选择结果引用 snapshot 内稳定 item identity  
**And** 前端行号、当前排序、数量摘要或自然语言列表不得成为对象事实源  
**And** snapshot 之后的负责人、在线状态、处置状态等关键事实变化不会改写历史快照  
**And** 未来 PrepareAction 必须重新校验并使不再满足策略的草案失效。

**AC-1.17.7**

**Given** 统一漏洞事实和快照已保存  
**When** Workbench、确认摘要、ActionResult、replay、audit 或模型上下文请求数据  
**Then** 所有出口从同一 Product Facts 投影生成  
**And** 字段缺失和冲突状态保持一致  
**And** raw provider payload、未授权关联数据和内部 identity 不进入产品出口。

**AC-1.17.8**

**Given** 本 Story 实现完成  
**When** 执行 fixture、mock 和授权 live-read smoke  
**Then** 覆盖完整集合、空集合、不完整分页、多源关联、字段缺失、事实冲突、五次查询引用、歧义澄清、上下文压缩、重启恢复、跨 scope 和 digest 损坏  
**And** facts/product/resolver、schema、Go tests、race tests 和安全泄漏测试全部通过  
**And** 只有统一行级事实的机器证据和目标部署 smoke 同时有效时，G-FACT-01 才更新为 PASS  
**And** 所有 G-WRITE、G-SAFE 和真实 mutation 继续保持阻塞。

### Source Bundle SB-1.18：交付只读桌面聊天 Workbench

As a 产品中心运营，  
I want 在桌面工作台中通过聊天查询并打开对应结果，  
So that 我可以在一个界面内完成查询、结果定位和事实核对。

**Acceptance Criteria:**

**AC-1.18.1**

**Given** Story 1.17 已完成  
**When** 执行开发前验证  
**Then** 范围只包含浅色桌面 Workbench 骨架、会话导航、聊天查询、查询结果卡、明细路由基础和右侧安全事实面板  
**And** 不实现真实审批卡、人员选择、确认按钮、完整百条明细表或任何 mutation。

**AC-1.18.2**

**Given** 前端开始实现  
**When** 建立 Workbench 应用边界  
**Then** 使用 Vite、React、TypeScript、React Router、TanStack Query、Zustand、Radix UI Primitives、lucide-react 和项目 CSS  
**And** 组件只消费 `contracts/generated.ts` 的生成类型  
**And** 不引入 Next.js、Redux、XState、Tailwind、shadcn 或平行 DTO。

**AC-1.18.3**

**Given** 用户在宽度不小于 1180px 的桌面浏览器打开 Workbench  
**When** 页面完成布局  
**Then** 使用左侧232px会话导航、中间主区域和右侧360px事实面板的三栏结构  
**And** 只提供浅色主题  
**And** 不实现移动断点、抽屉导航、触屏专用交互或主题切换。

**AC-1.18.4**

**Given** 用户在聊天输入自然语言查询  
**When** Run 通过 SSE 返回 Product Facts 投影  
**Then** 时间线按顺序展示用户消息、Agent 安全摘要和查询结果卡  
**And** 查询结果卡显示查询序号、条件摘要、对象数量和观察时间  
**And** 百条对象不在聊天时间线展开  
**And** 空结果、查询失败和结果不完整使用不同状态。

**AC-1.18.5**

**Given** 用户点击查询结果卡的“查看全部”  
**When** 路由进入对应漏洞明细工作区  
**Then** 保留当前会话、来源查询、冻结快照和右栏上下文  
**And** 明细表面至少显示来源查询与快照安全摘要  
**And** 返回聊天时定位到原查询结果卡  
**And** 路由参数不能代替后端 result ref 与 scope 校验。

**AC-1.18.6**

**Given** 用户切换会话、选择查询或进入明细表面  
**When** 左栏和右栏更新  
**Then** 会话切换不得合并、覆盖或重建历史 result ref  
**And** 右栏只显示当前安全事实和执行记录摘要  
**And** 技术信息入口在授权和字段 allowlist 未确定前整体隐藏  
**And** 右栏选择变化不得改变冻结对象范围。

**AC-1.18.7**

**Given** G-WRITE 或 G-SAFE 尚未通过  
**When** 用户输入派发、转发、延时或误报意图  
**Then** 后端 capability/policy 投影明确返回“写域尚未启用”状态  
**And** 前端不显示人员选择、审批确认、取消或执行控件  
**And** 不创建 ActionDraft、不调用 PrepareAction，也不靠前端关键词判断权限  
**And** 用户重复发送“确认”文字不会触发写入。

**AC-1.18.8**

**Given** 本 Story 实现完成  
**When** 执行组件测试、路由测试、SSE 重连测试、TypeScript typecheck 和桌面视觉验收  
**Then** 覆盖查询中、有结果、空结果、不完整、失败、会话切换、断线重连、结果卡跳转和写域禁用状态  
**And** 页面具备命名的 `nav`、唯一 `main`、补充 `aside`、跳到主内容入口和可见键盘焦点  
**And** TanStack Query 只管理 server state，Zustand 只管理 UI state  
**And** 实施计划、UX 验收记录和 M-5 证据同步更新。

### Source Bundle SB-1.19：交付百条分页漏洞明细与右侧事实面板

As a 产品中心运营，  
I want 在专门工作区逐条查看完整漏洞集合和选中对象事实，  
So that 我可以高效核对百条级查询结果而不会改变冻结范围。

**Acceptance Criteria:**

**AC-1.19.1**

**Given** Story 1.18 已完成  
**When** 执行开发前验证  
**Then** 范围只包含 IA-02 只读漏洞明细、每页100条分页、右侧事实与执行面板、键盘交互和视觉验收  
**And** 不实现选择后写入、批量编辑、人员选择、确认按钮或任何 mutation。

**AC-1.19.2**

**Given** 查询快照包含一条或多条漏洞  
**When** 用户进入漏洞明细工作区  
**Then** 中栏按冻结顺序每页显示100条  
**And** 始终显示当前页、总页数、总数和加载状态  
**And** 支持首页、上一页、下一页、末页和合法页码跳转  
**And** 翻页只读取同一 snapshot，不重新执行 FOBrain 查询。

**AC-1.19.3**

**Given** 一条统一漏洞事实被投影为明细行  
**When** 页面渲染该行  
**Then** 显示快照行序号、POC、IP、在线状态、业务系统、当前修复负责人和逐条状态  
**And** 可用时显示业务系统负责人、运维负责人和发现时间  
**And** 发现时间统一显示为 `YYYY-MM-DD HH:mm:ss（UTC+8）`  
**And** 当前负责人缺失显示“未分配”，业务系统及其他负责人缺失显示“未提供”。

**AC-1.19.4**

**Given** POC 或 IP 内容超过表格单元格可读空间  
**When** 渲染列表和右侧面板  
**Then** 表格中允许视觉省略，多个 IP 显示首个 IP 与剩余数量  
**And** 单元格可访问名称保留完整安全值  
**And** 选中行后右栏完整换行展示 POC 和全部 IP  
**And** 完整值不能只通过鼠标悬停查看。

**AC-1.19.5**

**Given** 用户通过鼠标或键盘选择一行  
**When** 当前行发生变化  
**Then** 右栏“事实”区域显示该漏洞的完整安全 Product Facts  
**And** “执行记录”区域只显示当前上下文的安全 Run、查询和核验摘要  
**And** 不显示内部 result ref、checkpoint、raw args 或 provider 字段  
**And** 行选择只改变查看上下文，不改变 snapshot 或未来动作范围。

**AC-1.19.6**

**Given** 明细页加载失败、SSE 重连或用户返回聊天后再次进入  
**When** 工作区恢复  
**Then** 当前 snapshot、页码和选中对象从权威 server state 与允许的 UI state 恢复  
**And** 局部加载失败不会把完整查询标记为空  
**And** 重试只重新读取安全投影，不重新执行外部查询或写入  
**And** 过期、无权限和损坏快照分别显示明确状态。

**AC-1.19.7**

**Given** 用户使用键盘、读屏或200%文本缩放  
**When** 操作明细工作区  
**Then** 行网格使用 roving focus 且只有一个 Tab 停靠点  
**And** 焦点顺序遵循 DOM，不使用正数 `tabindex`  
**And** 关键事实和状态不会因缩放被遮挡  
**And** 只有明细表自身可在必要时水平滚动  
**And** 状态同时使用文字、图标和符合语义的颜色。

**AC-1.19.8**

**Given** 本 Story 实现完成  
**When** 使用0、1、100、101和628条 fixture 执行组件、分页、键盘、可访问性、路由恢复和视觉测试  
**Then** 每条对象恰好出现一次且页码、总数和顺序正确  
**And** TypeScript typecheck、前端测试和桌面宽度及200%缩放视觉验收通过  
**And** 不存在可触发 PrepareAction 或 mutation 的控件  
**And** UX 验收记录、视觉证据和 M-5 实施记录同步更新。

### Source Bundle SB-1.20：完成 READ-01 端到端验收与门禁记录

As a 产品中心研发负责人，  
I want 用可复现的端到端测试验证完整只读旅程，  
So that Epic 1 只有在查询、事实、恢复和零写入均有证据时才被判定完成。

**Acceptance Criteria:**

**AC-1.20.1**

**Given** Story 1.19 已完成  
**When** 执行开发前验证  
**Then** 范围只包含 READ-01 自动化、授权 live-read smoke、浏览器验收、机器报告和门禁记录  
**And** 不新增业务能力、不修改写域策略、不执行任何真实 mutation。

**AC-1.20.2**

**Given** 准备执行 READ-01  
**When** 验证前置门禁  
**Then** G-TOOLCHAIN、G-ARCH-V2、G-READ-01/02 和 G-FACT-01 必须具有当前代码版本对应的有效机器证据  
**And** 任一门禁缺失、过期或阻塞时测试明确以 blocked/`exit 2` 结束  
**And** blocked、skipped、mock PASS 或历史证据不得计为端到端通过。

**AC-1.20.3**

**Given** 授权运营 actor 在聊天中要求“查询全部新增漏洞”  
**When** Eino Agent 执行查询  
**Then** 通过 Capability Registry 选择精确新增查询能力  
**And** 完整分页读取当前权限范围内的新增漏洞  
**And** 结果按发现时间全局倒序形成不可变 QueryResultSnapshot  
**And** 聊天显示查询序号、总数、观察时间和安全摘要。

**AC-1.20.4**

**Given** 查询结果卡已经生成  
**When** 用户点击“查看全部”并浏览明细  
**Then** 每页100条展示冻结集合且总数、顺序和全局行序号保持不变  
**And** 每行展示统一漏洞事实和确定性缺失状态  
**And** 选择一行只更新右侧事实与查询记录  
**And** 返回聊天后焦点回到来源卡  
**And** 查询后新到漏洞不会改变已打开快照。

**AC-1.20.5**

**Given** 查询返回0条、provider 失败、分页中断、浏览器断线或服务重启  
**When** 执行 READ-01 失败与恢复路径  
**Then** 0条显示“没有待派发漏洞”  
**And** 失败和不完整结果不得显示为空集合  
**And** SSE 重连及重启恢复后仍引用同一 snapshot  
**And** 未加载范围被明确标记，不以聊天摘要或旧缓存补造事实。

**AC-1.20.6**

**Given** 用户随后要求派发这些漏洞  
**When** 写域门禁仍未通过  
**Then** 系统显示“当前仅支持查询与只读核对，此写操作尚未启用”  
**And** 不生成 ActionDraft、PendingInteraction、审批卡或确认按钮  
**And** 不调用 PrepareAction 或任何 FOBrain 写接口  
**And** 本轮外部 mutation 计数可客观证明为0。

**AC-1.20.7**

**Given** 自动化与人工浏览器验收完成  
**When** 生成验收材料  
**Then** 机器报告记录代码版本、工具链、schema/fixture 版本、测试命令、门禁状态、数量、分页、排序、恢复和零写入断言  
**And** 视觉证据覆盖聊天结果卡、百条明细、右侧事实和写域禁用状态  
**And** 凭据、真实漏洞实体、人员信息、raw payload 和内部引用不进入报告或截图。

**AC-1.20.8**

**Given** READ-01 所有断言均通过  
**When** 更新 Epic 1 完成记录  
**Then** G-TOOLCHAIN、G-ARCH-V2、G-READ-01/02 和 G-FACT-01 的状态与证据链接同步到唯一 readiness gate  
**And** Epic 1 可以标记完成  
**And** 任一断言失败时 Epic 1 保持未完成并记录阻塞原因  
**And** 所有 G-WRITE、G-SAFE 和真实 mutation 继续保持阻塞。

## Epic 2：运营人工派发与责任人转发

运营可以基于 Epic 1 冻结的授权漏洞事实，从 FOBrain 全员列表人工选择接收人，将同一接收人的对象合并派发或主动转发；系统在聊天中二次确认，逐条写入、回读负责人，并清楚展示完成、待核验和待排查对象。当前只允许实施决策关闭、契约、fixture、mock、门禁取证与安全失败验证；真实派发和转发保持为 Target Story Candidate，必须通过适用门禁后才能转为正式 Story。

### Story 2.1：固定全员列表的安全选人规则

Requirements:
  FR: [FR-05, FR-09]
  NFR: [NFR-02, NFR-03]
  AR: [AR-33]
  Architecture: [AD-14]
  UX: [UX-DR-15]
  Gates: [G-READ-02, OQ-02]
  Milestone: [M-6-EVIDENCE]

As a 产品中心运营，  
I want 从 FOBrain 全员列表中明确选择一名接收人，  
So that 派发和转发不会因重名、无效人员或 AI 推荐而选错对象。

**Acceptance Criteria:**

**AC-2.1.1**

**Given** Epic 1 已具备完整人员列表安全投影  
**When** 执行本 Story 的开发前验证  
**Then** 范围只包含 OQ-02 决策、人员选择契约、fixture、测试环境原型和验收证据  
**And** 不提供生产可写入口、不创建 ActionDraft、不执行 mutation。

**AC-2.1.2**

**Given** 用户打开人员选择  
**When** 加载候选人  
**Then** 默认展示当前 actor 可见的完整 FOBrain 人员快照  
**And** 不按部门、业务系统、角色或 AI 判断预过滤  
**And** 不允许自由输入任意姓名或使用手工人员名单。

**AC-2.1.3**

**Given** 完整人员快照已加载  
**When** 用户搜索、筛选或查看排序  
**Then** 搜索只匹配规范化后的显示名称和已批准安全消歧字段  
**And** 默认保持冻结快照中的 FOBrain 来源顺序，不提供推荐排序  
**And** 只提供“全部／可选／不可选”状态筛选  
**And** 清空搜索后恢复完整列表及原顺序。

**AC-2.1.4**

**Given** 人员包含稳定业务标识和 provider 状态  
**When** 判断是否可选  
**Then** 只有稳定标识存在且状态属于经目标部署验证的允许集合时才可选  
**And** 缺少稳定标识、状态未知、停用或不在允许集合中的人员明确显示为不可选  
**And** 前端不得自行解释 provider 状态值。

**AC-2.1.5**

**Given** 人员列表存在同名人员  
**When** 用户查看同名候选  
**Then** 仅展示经批准 allowlist 中的部门、员工编号或其他安全消歧字段  
**And** 不展示内部人员 ID、账号 raw 字段或敏感组织信息  
**And** 如果没有足够的安全消歧字段，同名候选保持不可选并说明原因  
**And** 不通过列表顺序或模型推断代替消歧。

**AC-2.1.6**

**Given** 用户使用键盘操作人员选择卡  
**When** 搜索、移动和选择候选  
**Then** 输入框与候选列表使用可访问的 combobox/listbox 语义  
**And** 方向键移动候选、Enter 明确选择、Escape 关闭并恢复触发点焦点  
**And** 候选列表使用单一 Tab 停靠点和可见焦点  
**And** 加载、空结果、失败和不可选状态均有文字说明。

**AC-2.1.7**

**Given** 用户明确选择一名可选人员  
**When** 系统保存选择结果  
**Then** 绑定当前人员快照、稳定人员 identity 和安全显示摘要  
**And** 不自动预选、不因搜索结果唯一而自动确认  
**And** PrepareAction 和 ConfirmAction 必须重新校验人员仍在当前 actor 可见范围且仍可选  
**And** 人员事实变化时原草案失效。

**AC-2.1.8**

**Given** 本 Story 完成  
**When** 执行243人、同名、停用、状态未知、缺少标识、搜索、筛选、键盘和安全泄漏测试  
**Then** 人员状态允许集合、安全消歧 allowlist 和交互规则具有目标部署证据与产品负责人批准  
**And** OQ-02 才能标记为关闭，并同步 UX、schema、fixture、PRD 和 readiness gate  
**And** 任一证据不足时 OQ-02 保持阻塞，生产人员选择与真实写入继续隐藏。

### Story 2.2：用 Mock 固定手动派发业务规则

Requirements:
  FR: [FR-05, FR-06, FR-07, FR-16, FR-17, FR-18]
  NFR: [NFR-03, NFR-04]
  AR: [AR-10, AR-12, AR-15, AR-16, AR-40, AR-42, AR-45, AR-46]
  Architecture: [AD-07, AD-10, AD-11, AD-23]
  UX: [UX-DR-15, UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-19]
  Gates: [G-WRITE-01-BLOCKED, G-WRITE-02-BLOCKED, G-SAFE-01]
  Milestone: [M-6-MOCK]

As a 产品中心运营，  
I want 将已核对的新增漏洞派发给我明确选择的同一接收人，  
So that 系统能够冻结正确范围、二次确认并逐条判断结果。

**Acceptance Criteria:**

**AC-2.2.1**

**Given** Story 2.1 已关闭 OQ-02  
**When** 执行开发前验证  
**Then** 范围只包含派发 capability metadata、policy、schema、fixture、mock adapter、mock verifier 和测试环境审批组件  
**And** 复用 Epic 1 的 ActionDraft、确认、幂等和核验控制面  
**And** 不连接真实 FOBrain mutation，不在生产 Workbench 显示派发入口。

**AC-2.2.2**

**Given** 运营已经逐条核对一个完整新增漏洞快照  
**When** 运营选择接收人和对象范围  
**Then** 接收人必须来自 Story 2.1 定义的可选 FOBrain 人员快照  
**And** 对象必须来自当前 actor、workspace 和 conversation 下的完整 QueryResultSnapshot  
**And** AI 不推荐接收人、不自动选择对象，也不把数量摘要当作对象集合。

**AC-2.2.3**

**Given** 用户为多条漏洞指定接收人  
**When** 系统组织派发草案  
**Then** 只有最终指定给同一接收人的对象可以进入同一 ActionDraft  
**And** 不同接收人的对象必须形成不同草案并分别确认  
**And** 同一漏洞不能同时进入多个有效派发草案  
**And** 超过 capability `max_items` 时必须展示固定拆分范围并分别确认，或在写入前明确拒绝，不得后台静默分批。

**AC-2.2.4**

**Given** 用户请求准备派发  
**When** 执行 PrepareAction  
**Then** 重新校验 actor、对象权限、完整快照、人员可选状态和派发 policy  
**And** 冻结每条漏洞的 POC、IP、在线状态、业务系统、当前修复负责人、业务系统负责人、运维负责人和目标接收人  
**And** 必要事实缺失或冲突时对应对象不得进入可确认草案  
**And** 任何对象范围或目标人员变化都产生新草案和新 digest。

**AC-2.2.5**

**Given** 派发草案有效  
**When** 生成测试环境审批卡  
**Then** 显示来源查询、动作“派发”、目标接收人、对象总数和前5条统一漏洞事实  
**And** 提供查看完整冻结集合的入口  
**And** 确认和取消只存在于审批卡  
**And** 聊天中的“确认”文字、行选择或查看详情不构成批准。

**AC-2.2.6**

**Given** 用户未确认、取消、拒绝、确认过期、人员状态变化或关键漏洞事实变化  
**When** mock 控制面处理该草案  
**Then** mock mutation 调用次数为0  
**And** 草案进入明确的取消、拒绝、过期或失效状态  
**And** 需要继续时必须重新查询、选择人员、准备并确认  
**And** 不发送产品通知。

**AC-2.2.7**

**Given** 用户确认有效草案并执行 mock 派发  
**When** mock adapter 返回逐条结果并由 mock verifier 回读  
**Then** 只有接口被接受且回读修复负责人等于目标人员的条目显示完成  
**And** control lost 优先 `manual_attention/待排查`；held+accepted+conclusive real-readback match 才完成；held+rejected+policy-approved non-application 为 `failed/待排查`；其他组合按 deadline／预算分为 `reconciling/待核验` 或 `manual_attention/待排查`  
**And** 同批成功条目保持完成且不被失败项覆盖  
**And** 结果逐条展示 POC、IP、在线状态、业务系统、目标人员和实际修复负责人。

**AC-2.2.8**

**Given** 本 Story 完成  
**When** 执行同一接收人合并、不同接收人拆分、超出 `max_items`、重复对象、缺失事实、草案变化、取消零写入、部分失败、回读不一致、重复确认和恢复测试  
**Then** schema、fixture、mock provider、policy、Verifier、Go tests、前端测试和生成契约校验全部通过  
**And** FR-05～07、FR-16～18 的 mock 行为具有机器证据  
**And** 生产派发 capability、审批控件和真实 mutation 继续保持关闭。

### Story 2.3：用 Mock 固定主动转发业务规则

Requirements:
  FR: [FR-08, FR-09, FR-10, FR-16, FR-17, FR-18]
  NFR: [NFR-03, NFR-04]
  AR: [AR-10, AR-12, AR-15, AR-16, AR-40, AR-42, AR-45, AR-46]
  Architecture: [AD-07, AD-10, AD-11, AD-23]
  UX: [UX-DR-15, UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-19]
  Gates: [G-WRITE-01-BLOCKED, G-WRITE-02-BLOCKED, G-SAFE-01]
  Milestone: [M-6-MOCK]

As a 产品中心运营，  
I want 将权限范围内的漏洞主动转发给我选择的人员，  
So that 我可以独立纠正负责人，而不依赖系统判断是否错派。

**Acceptance Criteria:**

**AC-2.3.1**

**Given** Story 2.2 已完成  
**When** 执行开发前验证  
**Then** 范围只包含转发 capability metadata、policy、schema、fixture、mock adapter、mock verifier 和测试环境审批组件  
**And** 复用 Epic 1 的统一 Action 控制面  
**And** 不连接真实 FOBrain mutation，不在生产 Workbench 显示转发入口。

**AC-2.3.2**

**Given** 当前 actor 能通过 FOBrain 内置权限查询到一条或多条漏洞  
**When** 用户主动表达转发意图  
**Then** 不要求系统先判断该漏洞是否错派  
**And** 不要求此前发生过派发、失败或其他动作  
**And** 用户当前负责或当前可见的漏洞均按登记 policy 判断是否允许转发  
**And** AI 不自动识别错派、不推荐是否转发。

**AC-2.3.3**

**Given** 用户选择目标人员和待转发对象  
**When** 系统解析选择  
**Then** 接收人必须来自可选 FOBrain 人员快照  
**And** 对象必须来自当前 actor、workspace 和 conversation 下的完整 QueryResultSnapshot  
**And** 同一接收人的多条漏洞可以合并  
**And** 不同接收人的对象必须拆成不同草案分别确认  
**And** 超过 `max_items` 时不得后台静默分批。

**AC-2.3.4**

**Given** 用户请求准备转发  
**When** 执行 PrepareAction  
**Then** 重新校验 actor、对象权限、当前修复负责人、人员状态、完整快照和转发 policy  
**And** 冻结每条漏洞的 POC、IP、在线状态、业务系统、当前修复负责人和目标接收人  
**And** 当前负责人、对象权限或目标人员发生变化时原草案失效  
**And** 不从聊天摘要或前端选中行重建对象范围。

**AC-2.3.5**

**Given** 转发草案有效  
**When** 生成测试环境审批卡  
**Then** 显示来源查询、动作“转发”、当前负责人、目标接收人、对象总数和前5条统一漏洞事实  
**And** 提供查看完整冻结集合的入口  
**And** 确认和取消只存在于审批卡  
**And** 一次确认只覆盖当前不可变草案。

**AC-2.3.6**

**Given** 用户未确认、取消、拒绝、确认过期、权限变化或关键事实变化  
**When** mock 控制面处理该草案  
**Then** mock mutation 调用次数为0  
**And** 不把转发转化为派发或其他补偿动作  
**And** 需要继续时必须重新查询、选择人员、准备并确认  
**And** 不发送产品通知。

**AC-2.3.7**

**Given** 用户确认有效草案并执行 mock 转发  
**When** mock adapter 返回逐条结果并由 mock verifier 回读  
**Then** 只有接口被接受且回读修复负责人等于目标人员的条目显示完成  
**And** control lost 优先 `manual_attention/待排查`；held+accepted+conclusive real-readback match 才完成；held+rejected+policy-approved non-application 为 `failed/待排查`；其他组合按 deadline／预算分为 `reconciling/待核验` 或 `manual_attention/待排查`  
**And** 同批成功项不被失败项覆盖或重复执行  
**And** 结果逐条显示 POC、IP、在线状态、业务系统、原负责人、目标人员和实际负责人。

**AC-2.3.8**

**Given** 本 Story 完成  
**When** 执行无前置派发、当前用户主动转发、同一接收人合并、不同接收人拆分、权限变化、取消零写入、部分失败、回读不一致、重复确认和恢复测试  
**Then** schema、fixture、mock provider、policy、Verifier、Go tests、前端测试和生成契约校验全部通过  
**And** FR-08～10、FR-16～18 的 mock 行为具有机器证据  
**And** 生产转发 capability、审批控件和真实 mutation 继续保持关闭。

### Story 2.4：验证目标部署的派发写入与回读证据

Requirements:
  FR: [FR-05, FR-06, FR-07, FR-16, FR-17]
  NFR: [NFR-02, NFR-03, NFR-04]
  AR: [AR-26, AR-30, AR-31, AR-42]
  Architecture: [AD-10, AD-14, AD-23]
  UX: [UX-DR-18, UX-DR-20]
  Gates: [G-WRITE-01, G-WRITE-02, G-SAFE-01]
  Milestone: [M-6-EVIDENCE]

As a 产品中心研发负责人，  
I want 在获准环境中验证派发权限、写入和负责人回读，  
So that 真实派发只有在行为可证明且样本可恢复时才允许进入实现。

**Acceptance Criteria:**

**AC-2.4.1**

**Given** Story 2.2 的派发契约和 mock 验收已通过  
**When** 执行开发前验证  
**Then** 本 Story 被标记为 Gate/Evidence Story  
**And** 范围只包含目标部署取证、最小验证脚本、脱敏记录和门禁更新  
**And** 不接入生产 Workbench、不启用真实派发 capability。

**AC-2.4.2**

**Given** 准备验证 FOBrain 派发能力  
**When** 核对 provider 合同  
**Then** 明确实际权限、请求字段、人员 identity、对象 identity、响应语义、错误语义、超时、重复请求行为和负责人回读入口  
**And** provider 私有字段只在 adapter 边界映射  
**And** 接口 2xx 或返回成功文本不能单独作为派发成功证据。

**AC-2.4.3**

**Given** 验证可能改变目标部署数据  
**When** 准备首次 mutation  
**Then** 必须存在带日期的环境负责人授权、验证窗口、最小样本、最大写入次数、恢复方案和失败联系人  
**And** 样本必须允许恢复到验证前事实并可回读证明  
**And** 凭据只从被忽略的本地配置读取  
**And** 任一条件缺失时以 blocked/`exit 2` 结束且 mutation 次数为0。

**AC-2.4.4**

**Given** 一个已授权、可恢复的漏洞和目标人员  
**When** 在确认后执行一次最小派发验证  
**Then** 记录脱敏请求摘要和 provider 接受结果  
**And** 通过独立只读查询回读当前修复负责人  
**And** 只有负责人等于目标人员时该样本通过  
**And** 验证结束后恢复原事实并再次回读确认。

**AC-2.4.5**

**Given** 相同验证草案分别处于未确认、取消、拒绝、过期或 digest 不匹配状态  
**When** 运行安全失败验证  
**Then** provider mutation transport 的实际调用次数为0  
**And** audit/replay 保留明确原因  
**And** 不发送产品消息、邮件、webhook 或其他通知  
**And** FOBrain 自身可能产生的通知不作为成功证据。

**AC-2.4.6**

**Given** 已批准执行最小批量部分失败验证  
**When** 对相互独立的授权样本逐条派发  
**Then** 至少能够客观区分已回读成功、明确失败和回读不一致  
**And** 单条失败不覆盖已成功条目  
**And** 每个成功样本只发生一次 mutation  
**And** 若批准环境无法安全制造部分失败，G-WRITE-02 派发子项保持阻塞，不使用 mock 替代目标部署证据。

**AC-2.4.7**

**Given** mutation 请求超时、响应丢失或结果不确定  
**When** 系统无法证明是否写入  
**Then** 不重复发送派发请求  
**And** 先执行只读负责人回读并进入 reconcile  
**And** 在核验期限与调用预算内保持 `reconciling` 并投影为“待核验”，仅在预算耗尽后仍无法证明时记录为 `manual_attention` 并投影为“待排查”  
**And** 产品证据、日志和错误中不包含 raw payload、token、账号或真实人员信息。

**AC-2.4.8**

**Given** 所有验证和恢复步骤完成  
**When** 生成派发门禁记录  
**Then** 记录环境授权摘要、代码版本、契约版本、样本数量、mutation 次数、回读、恢复、部分失败和零写入断言  
**And** 仅按客观结果更新 G-WRITE-01/02 与 G-SAFE-01 的派发子项  
**And** 任一恢复失败、证据缺失或结果不一致时门禁保持阻塞并停止后续真实派发工作  
**And** 生产派发入口和 capability 继续关闭。

### Story 2.5：验证目标部署的主动转发与负责人回读证据

Requirements:
  FR: [FR-08, FR-09, FR-10, FR-16, FR-17]
  NFR: [NFR-02, NFR-03, NFR-04]
  AR: [AR-26, AR-30, AR-31, AR-42]
  Architecture: [AD-10, AD-14, AD-23]
  UX: [UX-DR-18, UX-DR-20]
  Gates: [G-WRITE-01, G-WRITE-02, G-SAFE-01]
  Milestone: [M-6-EVIDENCE]

As a 产品中心研发负责人，  
I want 在获准环境中独立验证主动转发权限、写入和负责人回读，  
So that 系统不会把派发接口存在误认为转发场景已经可用。

**Acceptance Criteria:**

**AC-2.5.1**

**Given** Story 2.3 的转发契约和 mock 验收已通过  
**When** 执行开发前验证  
**Then** 本 Story 被标记为 Gate/Evidence Story  
**And** 范围只包含转发目标部署取证、最小验证脚本、脱敏记录和门禁更新  
**And** 不接入生产 Workbench、不启用真实转发 capability。

**AC-2.5.2**

**Given** 准备验证 FOBrain 主动转发  
**When** 核对 provider 合同  
**Then** 分别证明当前 actor 的对象权限、目标人员约束、负责人变更入口和回读入口  
**And** 不假设派发和转发一定使用相同 endpoint、参数或权限  
**And** 即使底层接口相同，也必须记录转发场景的独立 policy 和证据  
**And** provider 接受或 HTTP 2xx 不能单独判定成功。

**AC-2.5.3**

**Given** 转发验证会改变已有负责人  
**When** 准备首次 mutation  
**Then** 必须有环境负责人授权、验证窗口、当前负责人已知的可恢复样本、目标人员、最大写入次数、恢复方案和失败联系人  
**And** 必须能够恢复原负责人并通过只读回读证明  
**And** 任一条件缺失时以 blocked/`exit 2` 结束且 mutation 次数为0。

**AC-2.5.4**

**Given** 当前 actor 可见且当前负责人事实明确的漏洞  
**When** 用户不经过任何前置派发而确认主动转发  
**Then** 系统向已选目标人员执行一次最小转发  
**And** 当前用户负责的漏洞也允许按实际 policy 验证转发给他人  
**And** 通过独立只读查询回读实际修复负责人  
**And** 只有实际负责人等于目标人员时该样本通过  
**And** 验证后恢复原负责人并再次回读确认。

**AC-2.5.5**

**Given** 相同转发草案处于未确认、取消、拒绝、过期、权限变化或 digest 不匹配状态  
**When** 运行安全失败验证  
**Then** provider mutation transport 的实际调用次数为0  
**And** 不触发派发或其他补偿动作  
**And** audit/replay 保存明确原因  
**And** 产品不发送通知。

**AC-2.5.6**

**Given** 已批准最小批量转发与部分失败验证  
**When** 对同一目标人员的多个授权样本逐条执行  
**Then** 能够分别识别负责人回读成功、明确失败和回读不一致  
**And** 单条失败不覆盖成功条目  
**And** 成功条目只发生一次 mutation  
**And** 若无法安全形成目标部署部分失败样本，G-WRITE-02 转发子项保持阻塞。

**AC-2.5.7**

**Given** 转发请求超时、响应丢失或结果不确定  
**When** 系统无法证明负责人是否改变  
**Then** 不重新发送转发 mutation  
**And** 先回读当前负责人并进入 reconcile  
**And** 在核验期限与调用预算内保持 `reconciling` 并投影为“待核验”，仅在预算耗尽后仍无法证明时记录为 `manual_attention` 并投影为“待排查”  
**And** 错误和证据只保存脱敏摘要，不包含真实人员、raw payload 或凭据。

**AC-2.5.8**

**Given** 转发验证、恢复和安全失败测试完成  
**When** 生成门禁记录  
**Then** 记录环境授权摘要、代码版本、契约版本、样本数量、mutation 次数、原负责人、目标结果、回读、恢复和零写入断言的脱敏结论  
**And** 仅按客观结果更新 G-WRITE-01/02 与 G-SAFE-01 的转发子项  
**And** 任一恢复失败、权限不明或回读不一致时门禁保持阻塞  
**And** 生产转发入口和 capability 继续关闭。

### Source Bundle SB-2.6：固定跨会话操作记录与待处理结果入口

> 非 Sprint 工作项：内容已迁移到 Story 1.36，保留本节仅用于追溯原始 AC。

As a 产品中心运营，  
I want 重新登录或离开原会话后，仍能找到最近六个月的操作结果，  
So that 我可以继续核验失败、结果未知和待排查的漏洞。

**Acceptance Criteria:**

**AC-2.6.1**

**Given** RD-03／AD-26 已关闭操作记录决策  
**When** 追溯本 Source Bundle  
**Then** 范围仅用于解释历史查询契约、Product Facts 投影、fixture 和 UX 决策如何迁移到 Story 1.36  
**And** 不得把本 Source Bundle 作为真实写域启用依据。

**AC-2.6.2**

**Given** 用户进入 Workbench  
**When** 查看左侧导航  
**Then** 存在固定入口“操作记录”  
**And** 该入口独立于聊天会话  
**And** 用户重新登录后仍可访问。

**AC-2.6.3**

**Given** 用户打开操作记录  
**When** 系统查询历史  
**Then** 默认展示最近六个月的数据并按操作时间倒序排列  
**And** 显示动作类型、状态、操作时间、来源查询、对象数量、目标人员以及完成、待核验、待排查数量  
**And** 不显示内部引用、raw payload、凭据或未投影字段。

**AC-2.6.4**

**Given** 存在多种操作记录  
**When** 用户筛选记录  
**Then** 支持按时间、动作类型和状态筛选  
**And** 提供“待核验/待排查”快捷筛选  
**And** 无结果、加载失败和无权限状态均有确定提示。

**AC-2.6.5**

**Given** 用户打开一条历史记录  
**When** 查看操作详情  
**Then** 进入既有漏洞操作工作区  
**And** 从原 ActionResult 和 Product Facts 展示逐条结果  
**And** 不重新调用 FOBrain 重建历史事实  
**And** 每条对象显示 POC、IP、在线状态、业务系统、原负责人、目标人员、实际负责人和处理状态。

**AC-2.6.6**

**Given** 用户访问历史数据  
**When** 后端执行查询  
**Then** 重新校验 actor 和 workspace 权限  
**And** 只返回当前 actor 被授权查看的数据  
**And** 无法证明身份或范围时拒绝访问。

**AC-2.6.7**

**Given** 产品提供六个月可查询范围  
**When** 数据超过六个月  
**Then** 六个月仅表示产品界面的可发现窗口，不代表自动物理删除 Product Facts  
**And** 仍处于执行、待确认、核验或 `manual_attention` 状态的引用不得因时间窗口被清理  
**And** 更长期保存由独立审计保留策略决定。

**AC-2.6.8**

**Given** 本 Story 完成  
**When** 执行重新登录、服务重启、六个月边界、分页、actor 隔离、待处理筛选、损坏引用和安全字段测试  
**Then** schema、fixture、后端测试、前端测试和生成契约校验全部通过  
**And** 固定入口、展示范围和权限失败语义必须与 Story 1.36／AD-26 一致  
**And** 本 Source Bundle 始终不得进入 Sprint 或开发。

### Target Candidate TC-2.1：接入真实手动派发能力

TargetRequirements:
  FR: [FR-05, FR-06, FR-07, FR-16, FR-17, FR-18]
  NFR: [NFR-02, NFR-03, NFR-04]
  Architecture: [AD-07, AD-10, AD-11, AD-23]
  UX: [UX-DR-15, UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-30]
  ConversionGates: [G-TOOLCHAIN, G-ARCH-V2, G-FACT-01, G-WRITE-01, G-WRITE-02, G-SAFE-01, OQ-02]

As a 产品中心运营，  
I want 查看新增漏洞并将选中的漏洞派发给指定人员，  
So that 漏洞能够由明确的修复负责人继续处理。

**Acceptance Criteria:**

**AC-2.7.1**

**Given** 本 Story 当前为 Target Story Candidate  
**When** 任一适用门禁尚未通过  
**Then** 不得进入 Sprint、不得交给开发 Agent、不得启用真实派发  
**And** mock 通过、HTTP 2xx 或人工口头确认不能解除限制。

**AC-2.7.2**

**Given** G-ARCH-V2、G-TOOLCHAIN、G-FACT-01、派发对应的 G-WRITE-01/02、G-SAFE-01 均有有效 PASS 证据  
**And** OQ-02 已关闭且 Story 1.36 的操作记录验收有效  
**And** 环境授权、可恢复样本和失败联系人有效  
**When** 重新执行 implementation-readiness 检查并获得批准  
**Then** 记录日期、证据和批准人  
**And** 将本 Candidate 转换为正式 Story 后才允许实施。

**AC-2.7.3**

**Given** 用户从完整新增漏洞结果中选择对象和接收人  
**When** 系统准备派发  
**Then** 接收人来自已重新校验的 FOBrain 可选人员列表  
**And** 对象来自完整且未过期的 QueryResultSnapshot  
**And** 同一接收人的多条漏洞合并为一个草案  
**And** 不同接收人分别形成草案  
**And** 超过 `max_items` 时拒绝继续，不静默拆批。

**AC-2.7.4**

**Given** 派发草案通过权限和事实复检  
**When** 后端生成审批投影  
**Then** 聊天区显示操作摘要  
**And** 操作工作区显示接收人、对象总数及逐条 POC、IP、在线状态、业务系统和当前负责人  
**And** 确认控件仅由后端 capability、policy 和审批状态决定  
**And** 聊天文本不能替代审批卡确认。

**AC-2.7.5**

**Given** 用户未确认、取消、拒绝、确认过期、权限变化或摘要不一致  
**When** 系统处理草案  
**Then** FOBrain mutation 调用次数为0  
**And** 需要继续时必须重新准备并确认  
**And** audit 记录脱敏原因。

**AC-2.7.6**

**Given** 用户确认有效草案  
**When** 系统逐条执行真实派发  
**Then** 每条对象最多执行一次 mutation  
**And** 不因超时或未知结果盲目重试  
**And** provider 数据只在 adapter 边界转换为 StructuredResult  
**And** 产品层不接收 raw provider payload。

**AC-2.7.7**

**Given** FOBrain 已接受派发请求  
**When** 系统独立回读漏洞负责人  
**Then** 只有实际负责人等于目标人员的条目显示“完成”  
**And** 明确失败显示“待排查”  
**And** 写入结果未知或回读不一致时，在核验期限与调用预算内进入 `reconciling` 并显示“待核验”  
**And** 预算耗尽后仍无法证明目标事实时进入 `manual_attention` 并显示“待排查”  
**And** 单条失败不覆盖或重复执行其他成功条目。

**AC-2.7.8**

**Given** 派发执行结束  
**When** 用户查看操作工作区或“操作记录”  
**Then** 展示每条漏洞的 POC、IP、在线状态、业务系统、目标人员、实际负责人和最终状态  
**And** 当前会话与跨会话入口引用相同 Product Facts 和 ActionResult  
**And** 产品不主动发送通知。

**AC-2.7.9**

**Given** 正式 Story 实施完成  
**When** 执行单条派发、同人合并、多人拆分、取消零写入、权限变化、重复确认、部分失败、超时、回读不一致、重启恢复和跨会话找回测试  
**Then** schema、fixture、Go tests、前端测试、生成契约和批准环境验收全部通过  
**And** Workbench 与 Action API 对同一操作投影一致  
**And** 验收记录不包含凭据、真实人员信息或 raw payload。

### Target Candidate TC-2.2：接入真实主动转发能力

TargetRequirements:
  FR: [FR-08, FR-09, FR-10, FR-16, FR-17, FR-18]
  NFR: [NFR-02, NFR-03, NFR-04]
  Architecture: [AD-07, AD-10, AD-11, AD-23]
  UX: [UX-DR-15, UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-30]
  ConversionGates: [G-TOOLCHAIN, G-ARCH-V2, G-FACT-01, G-WRITE-01, G-WRITE-02, G-SAFE-01, OQ-02]

As a 产品中心运营，  
I want 将权限范围内的漏洞主动转发给指定人员，  
So that 当前负责人不合适时可以直接调整修复责任人。

**Acceptance Criteria:**

**AC-2.8.1**

**Given** 本 Story 当前为 Target Story Candidate  
**When** 任一适用门禁尚未通过  
**Then** 不得进入 Sprint、不得交给开发 Agent、不得启用真实转发  
**And** 不得用派发能力存在或 mock 验收通过替代转发证据。

**AC-2.8.2**

**Given** G-ARCH-V2、G-TOOLCHAIN、G-FACT-01、转发对应的 G-WRITE-01/02、G-SAFE-01 均有有效 PASS 证据  
**And** OQ-02 已关闭且 Story 1.36 的操作记录验收有效  
**And** 环境授权、可恢复样本和失败联系人有效  
**When** 重新执行 implementation-readiness 检查并获得批准  
**Then** 记录日期、证据和批准人  
**And** 将本 Candidate 转换为正式 Story 后才允许实施。

**AC-2.8.3**

**Given** 用户选择权限范围内的漏洞和目标人员  
**When** 系统准备转发  
**Then** 转发是独立动作，不要求系统先识别错派，也不要求此前发生过派发  
**And** 当前用户负责的漏洞也允许转发给其他人员  
**And** 接收人来自重新校验后的 FOBrain 可选人员列表  
**And** 对象来自完整且未过期的 QueryResultSnapshot。

**AC-2.8.4**

**Given** 多个漏洞需要转发  
**When** 系统形成草案  
**Then** 同一接收人的多条漏洞可以合并  
**And** 不同接收人分别形成草案  
**And** 冻结每条漏洞的当前负责人和目标人员  
**And** 超过 `max_items` 时拒绝继续，不静默拆批。

**AC-2.8.5**

**Given** 转发草案通过权限和事实复检  
**When** 后端生成审批投影  
**Then** 聊天区显示转发摘要  
**And** 操作工作区显示当前负责人、目标人员及逐条 POC、IP、在线状态和业务系统  
**And** 确认控件仅由后端 capability、policy 和审批状态决定  
**And** 聊天文本不能替代审批卡确认。

**AC-2.8.6**

**Given** 用户未确认、取消、拒绝、确认过期、权限变化、当前负责人变化或摘要不一致  
**When** 系统处理草案  
**Then** FOBrain mutation 调用次数为0  
**And** 不自动改为派发或触发其他补偿动作  
**And** 继续操作前必须重新查询、准备并确认。

**AC-2.8.7**

**Given** 用户确认有效草案  
**When** 系统逐条执行真实转发并回读  
**Then** 每条对象最多执行一次 mutation  
**And** 只有实际负责人等于目标人员的条目显示“完成”  
**And** 明确失败显示“待排查”  
**And** 超时、未知结果或回读不一致时，在核验期限与调用预算内进入 `reconciling` 并显示“待核验”  
**And** 预算耗尽后仍无法证明目标事实时进入 `manual_attention` 并显示“待排查”  
**And** 不盲目重试或覆盖其他成功条目。

**AC-2.8.8**

**Given** 转发执行结束  
**When** 用户查看操作工作区或“操作记录”  
**Then** 展示每条漏洞的 POC、IP、在线状态、业务系统、原负责人、目标人员、实际负责人和最终状态  
**And** 当前会话与跨会话入口使用同一 Product Facts 和 ActionResult  
**And** 产品不主动发送通知。

**AC-2.8.9**

**Given** 正式 Story 实施完成  
**When** 执行无前置派发、当前用户对象转发、同人合并、多人拆分、取消零写入、负责人变化、重复确认、部分失败、超时、回读不一致、恢复和跨会话找回测试  
**Then** schema、fixture、Go tests、前端测试、生成契约和批准环境验收全部通过  
**And** Workbench 与 Action API 的事实和状态一致  
**And** 验收记录不包含凭据、真实人员信息或 raw payload。

### Target Candidate TC-2.3：完成派发与转发端到端验收

TargetRequirements:
  FR: [FR-05, FR-06, FR-07, FR-08, FR-09, FR-10, FR-16, FR-17, FR-18]
  NFR: [NFR-01, NFR-02, NFR-03, NFR-04]
  Architecture: [AD-07, AD-10, AD-11, AD-23, AD-26]
  UX: [UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-27, UX-DR-30]
  ConversionGates: [G-TOOLCHAIN, G-ARCH-V2, G-FACT-01, G-WRITE-01, G-WRITE-02, G-SAFE-01, OQ-02]

As a 产品中心研发负责人，  
I want 用可复现的端到端测试验证派发和主动转发完整流程，  
So that Epic 2 只有在真实写入、回读、恢复和安全边界均有证据时才能完成。

**Acceptance Criteria:**

**AC-2.9.1**

**Given** 本 Story 当前为 Target Story Candidate  
**When** TC-2.1 或 TC-2.2 尚未正式转换并完成  
**Then** 本 Story 不得进入 Sprint、不得交给开发 Agent  
**And** blocked、skipped、mock PASS 或历史证据不能作为通过信号。

**AC-2.9.2**

**Given** TC-2.1、TC-2.2 已正式转换并完成且所有适用门禁仍有效  
**When** 转换本 Candidate  
**Then** 记录代码版本、门禁证据、环境授权、转换日期和批准人  
**And** 范围只包含端到端测试、浏览器验收、恢复验证和完成裁决  
**And** 不新增业务能力或复制通用 Action 状态机。

**AC-2.9.3**

**Given** 授权运营用户查询到新增漏洞  
**When** 分别完成手动派发和主动转发旅程  
**Then** 两条旅程均经过完整快照、人员选择、草案冻结、审批确认、逐条写入和独立负责人回读  
**And** 同一接收人的对象合并处理  
**And** 不同接收人的对象分别确认  
**And** 转发不依赖此前发生派发。

**AC-2.9.4**

**Given** 草案处于未确认、取消、拒绝、过期、权限变化、人员变化、负责人变化或摘要不一致状态  
**When** 执行安全失败矩阵  
**Then** 对应 provider mutation 调用次数为0  
**And** 不通过派发、转发或其他动作进行补偿  
**And** audit/replay 保存脱敏原因。

**AC-2.9.5**

**Given** 批量执行出现成功、明确失败、超时或回读不一致  
**When** 系统生成逐条结果  
**Then** 成功项只执行一次并显示实际负责人  
**And** 失败项显示“待排查”  
**And** 未知或不一致项在核验期限与调用预算内进入 `reconciling` 并显示“待核验”  
**And** 预算耗尽后仍无法证明目标事实的条目进入 `manual_attention` 并显示“待排查”  
**And** 不盲目重试、不覆盖成功项、不把 2xx 当作最终成功。

**AC-2.9.6**

**Given** 执行期间发生浏览器断线、服务重启或结果响应丢失  
**When** 系统恢复运行  
**Then** 继续引用原 ActionDraft、ActionResult 和 Product Facts  
**And** 不重复 mutation  
**And** 通过只读负责人回读恢复能够证明的状态  
**And** 无法证明的对象先在核验期限与调用预算内保持 `reconciling`，仅在预算耗尽后进入 `manual_attention`。

**AC-2.9.7**

**Given** 操作完成或存在待处理对象  
**When** 用户离开会话、重新登录并打开“操作记录”  
**Then** 可以在六个月产品窗口内找到原操作  
**And** 操作工作区展示相同的逐条事实和状态  
**And** Workbench 与 Action API 投影一致  
**And** 不重新调用 FOBrain 重建历史结果。

**AC-2.9.8**

**Given** 自动化和人工浏览器验收完成  
**When** 生成验收材料  
**Then** 覆盖 FR-05～10、FR-16～18  
**And** 记录派发与转发的确认后单次写入、逐条回读、部分失败、恢复、零写入和错误脱敏证据  
**And** 报告和截图不包含凭据、真实人员、真实漏洞、raw payload 或内部引用。

**AC-2.9.9**

**Given** 所有端到端断言通过  
**When** 更新唯一 readiness gate  
**Then** 同步 G-WRITE-01/02、G-SAFE-01、OQ-02 与 Story 1.36 的最终状态和证据链接  
**And** 仅在派发和转发两条旅程均通过时允许 Epic 2 标记完成  
**And** 任一断言失败时 Epic 2 保持未完成并记录阻塞原因。

## Epic 3：管理员延时与运营误报处置

管理员可以对人工沟通确定的一条漏洞执行修复延时，运营可以对自己负责并已人工判断的漏洞标记误报；两类动作都复用 Epic 1 的可信确认、逐条执行、回读核验和恢复链路，只有目标事实回读一致才显示成功。当前只允许实施决策关闭、契约、fixture、mock、门禁取证与安全失败验证；真实延时和误报保持为 Target Story Candidate，必须通过适用门禁后才能转为正式 Story。

### Story 3.1：验证并固定直接修复延时规则

Requirements:
  FR: [FR-11, FR-12, FR-13]
  NFR: [NFR-02, NFR-03]
  AR: [AR-34]
  Architecture: [AD-07, AD-14]
  UX: [UX-DR-14, UX-DR-16]
  Gates: [OQ-04, G-WRITE-03]
  Milestone: [M-6-EVIDENCE]

As a 产品中心研发负责人，  
I want 用目标部署证据固定直接修复延时的状态、期限和权限规则，  
So that 产品不会猜测哪些漏洞可以延时或接受什么新期限。

**Acceptance Criteria:**

**AC-3.1.1**

**Given** OQ-04 尚未关闭  
**When** 实施本 Story  
**Then** 本 Story 被标记为 Decision/Gate Story  
**And** 范围只包含只读合同取证、规则表、schema、policy、fixture、测试和文档更新  
**And** 不执行真实 mutation、不启用生产延时 capability。

**AC-3.1.2**

**Given** FOBrain 源码同时存在直接延时和申请审批延时流程  
**When** 确定本轮使用的能力  
**Then** 必须用目标部署合同、权限或产品负责人证据证明采用“直接修复延时”  
**And** 不把申请、审批、拒绝或催办流程混入本轮  
**And** 历史源码只作为候选线索，不能单独关闭 OQ-04。

**AC-3.1.3**

**Given** 目标部署提供延时规则证据  
**When** 固定允许操作的对象范围  
**Then** 明确列出允许的原状态集合及对应稳定状态值  
**And** 状态缺失、未知、不在集合或映射冲突时 fail closed  
**And** provider 私有状态只在 adapter 边界映射。

**AC-3.1.4**

**Given** 用户需要输入新修复期限  
**When** 固定期限合同  
**Then** 明确最早值、最晚值、时间粒度、规范时区以及与当前时间和现有期限的关系  
**And** 对等于边界、越界、与现期限相同、早于现期限和跨日期边界分别给出客观结论  
**And** 产品不通过模型推断或散落硬编码补齐未知规则。

**AC-3.1.5**

**Given** 期限规则已形成  
**When** 定义产品输入与事实表示  
**Then** 用户必须显式输入或选择新期限，AI 不得推荐、决定或自动修正  
**And** schema 使用带时区的规范时间表示并保留展示时区  
**And** 延时原因允许为空；填写时作为用户输入保存安全文本，不作为成功条件。

**AC-3.1.6**

**Given** 直接延时只有管理员可以执行  
**When** 定义 actor policy  
**Then** 管理员资格来自当前 FOBrain 身份和稳定角色映射  
**And** 不按显示名称、聊天内容或前端传值推断管理员  
**And** 未知、缺失或冲突角色均拒绝准备和确认。

**AC-3.1.7**

**Given** 状态、期限和角色规则已记录  
**When** 生成契约测试材料  
**Then** fixture 覆盖允许与拒绝状态、期限上下边界、粒度、时区、相同期限、原因空值、管理员、非管理员和未知权限  
**And** policy 在 prepare 与 confirm 两个时点使用同一版本化规则  
**And** 每次延时的 `max_items` 固定为1。

**AC-3.1.8**

**Given** 规则表、证据和测试完成  
**When** 产品负责人评审 OQ-04  
**Then** 只有目标部署证据、schema、policy、fixture 和产品决定一致时关闭 OQ-04  
**And** 同步 PRD、UX、readiness gate 和证据链接  
**And** 任一规则仍未知时 OQ-04 保持阻塞，任何目标部署写入证据或真实延时 Candidate 都不得进入实施阶段。

### Story 3.2：用 Mock 固定单条修复延时流程

Requirements:
  FR: [FR-11, FR-12, FR-13, FR-16, FR-17, FR-18]
  NFR: [NFR-03, NFR-04]
  AR: [AR-10, AR-12, AR-15, AR-16, AR-34, AR-42, AR-45, AR-46]
  Architecture: [AD-07, AD-10, AD-11, AD-23]
  UX: [UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-19]
  Gates: [G-WRITE-01-BLOCKED, G-WRITE-03-BLOCKED, G-SAFE-01]
  Milestone: [M-6-MOCK]

As a FOBrain 管理员，  
I want 在测试环境中对一条漏洞确认新的修复期限，  
So that 延时动作的人工输入、权限、确认和双事实核验可以先被机器验证。

**Acceptance Criteria:**

**AC-3.2.1**

**Given** Story 3.1 已固定 OQ-04 的版本化规则  
**When** 注册 mock 延时 capability  
**Then** capability 声明 `max_items=1`、管理员角色谓词、对象状态谓词、期限约束和 verifier 版本  
**And** 复用 Epic 1 的 ActionDraft、确认、幂等、reconcile 和结果投影  
**And** 不创建延时专用状态机或生产 provider 分支。

**AC-3.2.2**

**Given** 管理员与相关人员已在线下沟通  
**When** 管理员选择一条漏洞并输入新期限  
**Then** 是否延时和新期限完全来自用户输入  
**And** 原因可填写也可留空  
**And** AI 不推荐期限、不判断是否应延时、不自动确认。

**AC-3.2.3**

**Given** 用户请求准备延时  
**When** 执行 PrepareAction  
**Then** 重新校验 actor 为管理员、对象属于完整且未过期的授权快照、当前状态允许延时且新期限满足规则  
**And** 冻结 POC、IP、在线状态、业务系统、当前修复负责人、当前状态、当前期限、新期限和原因  
**And** 多于一条对象、规则未知或任一校验失败时不生成草案。

**AC-3.2.4**

**Given** 延时草案有效  
**When** 生成测试环境审批卡  
**Then** 完整展示唯一漏洞的统一事实、当前状态、当前期限、新期限和可选原因  
**And** 明确动作是“修复延时”  
**And** 确认和取消只存在于审批卡  
**And** 一次确认只覆盖当前不可变草案。

**AC-3.2.5**

**Given** 用户未确认、取消、拒绝、确认过期、管理员权限变化、对象状态变化、期限规则变化或 digest 不匹配  
**When** mock 控制面处理草案  
**Then** mock mutation 调用次数为0  
**And** 原草案失效并保留脱敏原因  
**And** 继续操作必须重新查询、输入期限、准备并确认。

**AC-3.2.6**

**Given** 管理员确认有效草案  
**When** mock adapter 接受延时并由 mock verifier 回读  
**Then** 只有回读状态为“延时”且实际修复期限等于冻结的新期限时显示“延时成功”  
**And** 仅状态一致、仅期限一致或任一事实缺失均不得显示成功  
**And** provider 接受或 mock 2xx 不能单独作为成功证据。

**AC-3.2.7**

**Given** mock 写入明确未应用、响应未知、回读不一致、预算耗尽或失租约未知  
**When** 系统投影结果  
**Then** control lost 优先 `manual_attention/待排查`；held+accepted+conclusive real-readback match 才完成；held+rejected+policy-approved non-application 为 `failed/待排查`；其他组合按 deadline／预算分为 `reconciling/待核验` 或 `manual_attention/待排查`  
**And** 不盲目重试、不发送产品通知  
**And** 结果继续显示原事实、目标期限和实际回读事实。

**AC-3.2.8**

**Given** 本 Story 完成  
**When** 执行原因空值、边界期限、非管理员、一次多条、取消零写入、权限变化、状态变化、重复确认、双事实不一致和重启恢复测试  
**Then** schema、fixture、mock provider、policy、Verifier、Go tests、前端测试和生成契约校验全部通过  
**And** FR-11～13、FR-16～18 的 mock 行为具有机器证据  
**And** 生产延时 capability、审批控件和真实 mutation 继续关闭。

### Story 3.3：用 Mock 固定人工误报标记流程

Requirements:
  FR: [FR-14, FR-15, FR-16, FR-17, FR-18]
  NFR: [NFR-03, NFR-04]
  AR: [AR-10, AR-12, AR-15, AR-16, AR-42, AR-45, AR-46]
  Architecture: [AD-07, AD-10, AD-11, AD-23]
  UX: [UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-19]
  Gates: [G-WRITE-01-BLOCKED, G-WRITE-04-BLOCKED, G-SAFE-01]
  Milestone: [M-6-MOCK]

As a 产品中心运营，  
I want 在测试环境中把自己负责且已人工判断的漏洞标记为误报，  
So that 系统能够验证人工控制、对象权限、确认和状态回读规则。

**Acceptance Criteria:**

**AC-3.3.1**

**Given** Epic 1 的通用 Action 控制面和统一漏洞事实已固定  
**When** 注册 mock 误报 capability  
**Then** capability 声明当前 actor 负责对象的权限谓词、允许原状态、目标状态“误报”、`max_items` 和 verifier 版本  
**And** 复用同一 ActionDraft、确认、幂等、逐条执行、reconcile 和结果投影  
**And** 不创建误报专用状态机。

**AC-3.3.2**

**Given** 操作人查看自己权限范围内的漏洞  
**When** 操作人选择对象并发起误报标记  
**Then** 误报结论必须由操作人明确作出  
**And** AI 不推荐、不推断、不批量自动选择误报对象  
**And** 当前修复负责人必须通过稳定 identity 与当前 actor 一致。

**AC-3.3.3**

**Given** 用户选择一条或多条漏洞  
**When** 执行 PrepareAction  
**Then** 对象来自完整且未过期的授权 QueryResultSnapshot  
**And** 逐条重新校验当前负责人、对象权限和允许原状态  
**And** 数量不得超过 capability 的 `max_items`，超过时拒绝且不静默拆批  
**And** 任一不符合对象必须明确指出，不得静默删除后继续。

**AC-3.3.4**

**Given** 误报草案有效  
**When** 生成测试环境审批卡  
**Then** 冻结并展示每条漏洞的 POC、IP、在线状态、业务系统、当前修复负责人、当前状态和目标状态“误报”  
**And** 多条默认展示前5条并可查看完整冻结集合  
**And** 确认和取消只存在于审批卡  
**And** 聊天文本不能替代确认。

**AC-3.3.5**

**Given** 用户未确认、取消、拒绝、确认过期、负责人变化、权限变化、状态变化或 digest 不匹配  
**When** mock 控制面处理草案  
**Then** mock mutation 调用次数为0  
**And** 原草案失效并记录脱敏原因  
**And** 不自动转为派发、转发或其他动作。

**AC-3.3.6**

**Given** 用户确认有效草案并执行 mock 误报标记  
**When** mock adapter 返回逐条结果并由 mock verifier 回读  
**Then** 只有 provider 接受且回读状态为“误报”的条目显示完成  
**And** control lost 优先 `manual_attention/待排查`；held+accepted+conclusive real-readback match 才完成；held+rejected+policy-approved non-application 为 `failed/待排查`；其他组合按 deadline／预算分为 `reconciling/待核验` 或 `manual_attention/待排查`  
**And** 单条失败不覆盖或重复执行其他成功条目。

**AC-3.3.7**

**Given** mock 操作结束  
**When** 用户查看逐条结果  
**Then** 每条展示统一漏洞事实、原状态、目标状态、实际状态和最终处理状态  
**And** 失败说明只使用中文脱敏语义  
**And** 产品不发送消息、邮件、webhook、催办或通知中心通知。

**AC-3.3.8**

**Given** 本 Story 完成  
**When** 执行自己负责、非自己负责、多条、超量、取消零写入、负责人变化、部分失败、回读不一致、重复确认和恢复测试  
**Then** schema、fixture、mock provider、policy、Verifier、Go tests、前端测试和生成契约校验全部通过  
**And** FR-14～18 的 mock 行为具有机器证据  
**And** 生产误报 capability、审批控件和真实 mutation 继续关闭。

### Story 3.4：验证目标部署的修复延时与双事实回读证据

Requirements:
  FR: [FR-11, FR-12, FR-13, FR-16, FR-17]
  NFR: [NFR-02, NFR-03, NFR-04]
  AR: [AR-26, AR-30, AR-31, AR-34, AR-42]
  Architecture: [AD-10, AD-14, AD-23]
  UX: [UX-DR-18, UX-DR-20]
  Gates: [G-WRITE-01, G-WRITE-03, G-SAFE-01, OQ-04]
  Milestone: [M-6-EVIDENCE]

As a 产品中心研发负责人，  
I want 在获准环境中验证直接修复延时的权限、写入和双事实回读，  
So that 真实延时只有在状态与期限都能被证明时才允许实现。

**Acceptance Criteria:**

**AC-3.4.1**

**Given** Story 3.1、3.2 已完成且 OQ-04 已关闭  
**When** 执行开发前验证  
**Then** 本 Story 被标记为 Gate/Evidence Story  
**And** 范围只包含目标部署取证、最小验证脚本、脱敏记录和门禁更新  
**And** 不接入生产 Workbench、不启用真实延时 capability。

**AC-3.4.2**

**Given** 准备验证 FOBrain 直接修复延时  
**When** 核对 provider 合同  
**Then** 证明使用的是直接延时而非申请审批流程  
**And** 明确管理员权限、允许原状态、期限格式、原因字段、响应语义、错误语义和状态／期限回读入口  
**And** provider 接受或 HTTP 2xx 不能单独判定成功。

**AC-3.4.3**

**Given** 延时验证会改变漏洞状态和期限  
**When** 准备首次 mutation  
**Then** 必须存在带日期的环境授权、验证窗口、当前状态与期限已知的单条可恢复样本、新期限、最大写入次数、恢复方案和失败联系人  
**And** 样本必须能够恢复原状态与原期限并通过只读回读证明  
**And** 任一条件缺失时以 blocked/`exit 2` 结束且 mutation 次数为0。

**AC-3.4.4**

**Given** 管理员和非管理员测试身份均已获准  
**When** 验证角色边界  
**Then** 管理员可以进入最小已确认验证  
**And** 非管理员、未知角色或角色映射冲突在写入前被拒绝  
**And** 被拒绝路径的 provider mutation transport 实际调用次数为0。

**AC-3.4.5**

**Given** 一个符合 OQ-04 规则的授权漏洞和新期限  
**When** 管理员确认后执行一次直接延时  
**Then** 记录脱敏请求摘要和 provider 接受结果  
**And** 独立回读实际状态与实际修复期限  
**And** 只有状态为“延时”且期限等于冻结新期限时样本通过  
**And** 验证结束后恢复原状态与期限并再次回读确认。

**AC-3.4.6**

**Given** 相同草案处于未确认、取消、拒绝、过期、权限变化、状态变化、期限规则变化或 digest 不匹配状态  
**When** 运行安全失败验证  
**Then** provider mutation transport 的实际调用次数为0  
**And** audit/replay 保存明确脱敏原因  
**And** 产品不发送通知，FOBrain 自身通知也不作为成功证据。

**AC-3.4.7**

**Given** 延时请求超时、响应丢失、仅一个目标事实一致或结果不确定  
**When** 系统无法证明最终状态  
**Then** 不重新发送延时 mutation  
**And** 先回读状态和期限并进入 reconcile  
**And** 在核验期限与调用预算内保持 `reconciling` 并投影为“待核验”，仅在预算耗尽后仍无法证明时记录为 `manual_attention` 并投影为“待排查”  
**And** 错误与证据不包含 raw payload、凭据或真实漏洞信息。

**AC-3.4.8**

**Given** 延时验证、恢复和安全失败测试完成  
**When** 生成门禁记录  
**Then** 记录授权摘要、代码与契约版本、样本数量、mutation 次数、原事实、目标事实、回读、恢复和零写入断言  
**And** 仅按客观结果更新 G-WRITE-01、G-WRITE-03 与 G-SAFE-01 的延时子项  
**And** 任一恢复失败、权限不明或双事实不一致时门禁保持阻塞  
**And** 生产延时入口和 capability 继续关闭。

### Story 3.5：验证目标部署的误报权限、写入与状态回读证据

Requirements:
  FR: [FR-14, FR-15, FR-16, FR-17]
  NFR: [NFR-02, NFR-03, NFR-04]
  AR: [AR-26, AR-30, AR-31, AR-42]
  Architecture: [AD-10, AD-14, AD-23]
  UX: [UX-DR-18, UX-DR-20]
  Gates: [G-WRITE-01, G-WRITE-04, G-SAFE-01]
  Milestone: [M-6-EVIDENCE]

As a 产品中心研发负责人，  
I want 在获准环境中验证误报标记的对象权限、写入和状态回读，  
So that 产品不会允许用户修改不属于自己的漏洞或把接口接受误报为完成。

**Acceptance Criteria:**

**AC-3.5.1**

**Given** Story 3.3 的误报契约和 mock 验收已通过  
**When** 执行开发前验证  
**Then** 本 Story 被标记为 Gate/Evidence Story  
**And** 范围只包含目标部署取证、最小验证脚本、脱敏记录和门禁更新  
**And** 不接入生产 Workbench、不启用真实误报 capability。

**AC-3.5.2**

**Given** 准备验证 FOBrain 误报能力  
**When** 核对 provider 合同  
**Then** 明确“当前 actor 负责”的稳定 identity 规则、允许原状态、目标状态值、请求字段、响应与错误语义和状态回读入口  
**And** 不依赖显示名称、聊天声明或前端传入负责人证明对象权限  
**And** provider 接受或 HTTP 2xx 不能单独判定成功。

**AC-3.5.3**

**Given** 误报验证会改变漏洞状态  
**When** 准备首次 mutation  
**Then** 必须存在带日期的环境授权、验证窗口、当前 actor 确认负责且原状态已知的可恢复样本、最大写入次数、恢复方案和失败联系人  
**And** 样本必须能够恢复原状态并通过只读回读证明  
**And** 任一条件缺失时以 blocked/`exit 2` 结束且 mutation 次数为0。

**AC-3.5.4**

**Given** 已授权验证自己负责与非自己负责的对象边界  
**When** 执行权限测试  
**Then** 只有稳定负责人 identity 与当前 actor 一致的对象可进入已确认写入  
**And** 非自己负责、负责人缺失、重名冲突或权限未知对象在写入前被拒绝  
**And** 所有拒绝路径的 provider mutation transport 实际调用次数为0。

**AC-3.5.5**

**Given** 操作人已人工判断一个授权样本为误报  
**When** 二次确认后执行一次最小误报标记  
**Then** 记录脱敏请求摘要和 provider 接受结果  
**And** 通过独立只读查询回读实际状态  
**And** 只有实际状态为“误报”时样本通过  
**And** 验证结束后恢复原状态并再次回读确认。

**AC-3.5.6**

**Given** 已批准最小多对象与部分失败验证  
**When** 对独立授权样本逐条执行  
**Then** 能够分别识别状态回读成功、明确失败和回读不一致  
**And** 单条失败不覆盖成功条目  
**And** 成功条目只发生一次 mutation  
**And** 若无法安全形成目标部署部分失败证据，则真实多对象误报保持阻塞。

**AC-3.5.7**

**Given** 草案未确认、取消、拒绝、过期、负责人变化或请求结果不确定  
**When** 执行安全与恢复验证  
**Then** 未确认类路径的 mutation 次数为0  
**And** 不确定写入不盲目重试，先回读状态并进入 reconcile  
**And** 在核验期限与调用预算内保持 `reconciling` 并投影为“待核验”，仅在预算耗尽后仍无法证明时记录为 `manual_attention` 并投影为“待排查”  
**And** 产品不发送通知，证据不包含 raw payload、凭据或真实人员信息。

**AC-3.5.8**

**Given** 误报验证、恢复和安全失败测试完成  
**When** 生成门禁记录  
**Then** 记录授权摘要、代码与契约版本、样本数量、mutation 次数、负责人边界、原状态、目标状态、回读、恢复和零写入断言  
**And** 仅按客观结果更新 G-WRITE-01、G-WRITE-04 与 G-SAFE-01 的误报子项  
**And** 任一恢复失败、权限不明或状态不一致时门禁保持阻塞  
**And** 生产误报入口和 capability 继续关闭。

### Target Candidate TC-3.1：接入真实单条修复延时能力

TargetRequirements:
  FR: [FR-11, FR-12, FR-13, FR-16, FR-17, FR-18]
  NFR: [NFR-02, NFR-03, NFR-04]
  Architecture: [AD-07, AD-10, AD-11, AD-23]
  UX: [UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-30]
  ConversionGates: [G-TOOLCHAIN, G-ARCH-V2, G-FACT-01, G-WRITE-01, G-WRITE-03, G-SAFE-01, OQ-04]

As a FOBrain 管理员，  
I want 对人工沟通确定的一条漏洞执行修复延时，  
So that 新修复期限能够被受控写入并通过状态和期限回读确认。

**Acceptance Criteria:**

**AC-3.6.1**

**Given** 本 Story 当前为 Target Story Candidate  
**When** 任一适用门禁尚未通过  
**Then** 不得进入 Sprint、不得交给开发 Agent、不得启用真实延时  
**And** mock 通过、接口存在、HTTP 2xx 或人工口头确认不能解除限制。

**AC-3.6.2**

**Given** G-ARCH-V2、G-TOOLCHAIN、G-FACT-01、延时对应的 G-WRITE-01、G-WRITE-03、G-SAFE-01 均有有效 PASS 证据  
**And** OQ-04 已关闭且 Story 1.36 的操作记录验收有效  
**And** 环境授权、可恢复样本和失败联系人有效  
**When** 重新执行 implementation-readiness 检查并获得批准  
**Then** 记录日期、证据和批准人  
**And** 将本 Candidate 转换为正式 Story 后才允许实施。

**AC-3.6.3**

**Given** 管理员已通过人工沟通确定是否延时和新期限  
**When** 从授权漏洞事实中选择对象  
**Then** 一次只能选择一条漏洞  
**And** 新期限必须由管理员显式输入并满足 OQ-04 的版本化规则  
**And** 原因允许为空  
**And** AI 不决定、推荐或自动修正延时结论和期限。

**AC-3.6.4**

**Given** 用户请求准备延时  
**When** 后端执行 PrepareAction  
**Then** 重新校验管理员身份、对象权限、允许原状态、当前期限和新期限  
**And** 从完整且未过期的 QueryResultSnapshot 冻结统一漏洞事实、当前状态、当前期限、新期限和原因  
**And** 任一事实或规则变化时原草案失效。

**AC-3.6.5**

**Given** 延时草案有效  
**When** 后端生成审批投影  
**Then** 聊天区显示操作摘要  
**And** 操作工作区完整显示唯一漏洞、当前状态、当前期限、新期限和原因  
**And** 确认控件只由后端 capability、policy 和审批状态决定  
**And** 聊天文本不能替代审批卡确认。

**AC-3.6.6**

**Given** 用户未确认、取消、拒绝、确认过期、权限变化、状态变化、期限规则变化或 digest 不匹配  
**When** 系统处理草案  
**Then** FOBrain mutation 调用次数为0  
**And** 继续操作必须重新查询、输入期限、准备并确认  
**And** audit 记录脱敏原因。

**AC-3.6.7**

**Given** 管理员确认有效草案  
**When** 系统执行一次真实延时并独立回读  
**Then** 每个草案最多执行一次 mutation  
**And** 只有实际状态为“延时”且实际期限等于目标期限时显示“延时成功”  
**And** 明确失败显示“待排查”  
**And** 超时、未知结果或任一事实不一致时，在核验期限与调用预算内进入 `reconciling` 并显示“待核验”  
**And** 预算耗尽后仍无法同时证明目标状态与期限时进入 `manual_attention` 并显示“待排查”。

**AC-3.6.8**

**Given** 延时执行结束  
**When** 用户查看操作工作区或“操作记录”  
**Then** 展示 POC、IP、在线状态、业务系统、修复负责人、原状态、原期限、目标期限、实际状态、实际期限和最终状态  
**And** 当前会话与跨会话入口使用同一 Product Facts 和 ActionResult  
**And** 产品不主动发送通知。

**AC-3.6.9**

**Given** 正式 Story 实施完成  
**When** 执行原因空值、期限边界、非管理员、一次多条、取消零写入、权限或状态变化、重复确认、超时、双事实不一致、重启恢复和跨会话找回测试  
**Then** schema、fixture、Go tests、前端测试、生成契约和批准环境验收全部通过  
**And** Workbench 与 Action API 对同一延时事实和状态投影一致  
**And** 验收记录不包含凭据、真实漏洞、raw payload 或内部引用。

### Target Candidate TC-3.2：接入真实人工误报标记能力

TargetRequirements:
  FR: [FR-14, FR-15, FR-16, FR-17, FR-18]
  NFR: [NFR-02, NFR-03, NFR-04]
  Architecture: [AD-07, AD-10, AD-11, AD-23]
  UX: [UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-30]
  ConversionGates: [G-TOOLCHAIN, G-ARCH-V2, G-FACT-01, G-WRITE-01, G-WRITE-04, G-SAFE-01]

As a 产品中心运营，  
I want 将自己负责且已人工判断的漏洞标记为误报，  
So that 已确认的误报对象不再持续显示为未处理漏洞。

**Acceptance Criteria:**

**AC-3.7.1**

**Given** 本 Story 当前为 Target Story Candidate  
**When** 任一适用门禁尚未通过  
**Then** 不得进入 Sprint、不得交给开发 Agent、不得启用真实误报  
**And** mock 通过、接口存在或前端控件存在不能解除限制。

**AC-3.7.2**

**Given** G-ARCH-V2、G-TOOLCHAIN、G-FACT-01、误报对应的 G-WRITE-01、G-WRITE-04、G-SAFE-01 均有有效 PASS 证据  
**And** Story 1.36 的操作记录验收有效  
**And** 环境授权、可恢复样本和失败联系人有效  
**When** 重新执行 implementation-readiness 检查并获得批准  
**Then** 记录日期、证据和批准人  
**And** 将本 Candidate 转换为正式 Story 后才允许实施。

**AC-3.7.3**

**Given** 操作人查看权限范围内的漏洞  
**When** 明确选择对象并发起误报标记  
**Then** 每条对象的当前修复负责人 stable identity 必须与当前 actor 一致  
**And** 误报判断完全由操作人作出  
**And** AI 不推荐、不推断、不自动选择误报对象。

**AC-3.7.4**

**Given** 用户选择一条或多条自己负责的漏洞  
**When** 后端执行 PrepareAction  
**Then** 对象来自完整且未过期的 QueryResultSnapshot  
**And** 逐条重新校验 actor、负责人、对象权限和允许原状态  
**And** 数量遵守 capability 的 `max_items`，超过时拒绝而不静默拆批  
**And** 冻结统一漏洞事实、原状态和目标状态“误报”。

**AC-3.7.5**

**Given** 误报草案有效  
**When** 后端生成审批投影  
**Then** 聊天区显示操作摘要  
**And** 操作工作区逐条展示 POC、IP、在线状态、业务系统、当前修复负责人、原状态和目标状态  
**And** 多条默认展示前5条并可查看完整冻结集合  
**And** 聊天文本不能替代审批卡确认。

**AC-3.7.6**

**Given** 用户未确认、取消、拒绝、确认过期、负责人变化、权限变化、状态变化或 digest 不匹配  
**When** 系统处理草案  
**Then** FOBrain mutation 调用次数为0  
**And** 不自动触发派发、转发或其他补偿动作  
**And** 继续操作必须重新查询、准备并确认。

**AC-3.7.7**

**Given** 用户确认有效草案  
**When** 系统逐条执行真实误报标记并回读  
**Then** 每条对象最多执行一次 mutation  
**And** 只有实际状态为“误报”的条目显示完成  
**And** 明确失败显示“待排查”  
**And** 超时、未知结果或状态不一致时，在核验期限与调用预算内进入 `reconciling` 并显示“待核验”  
**And** 预算耗尽后仍无法证明目标状态时进入 `manual_attention` 并显示“待排查”  
**And** 单条失败不覆盖或重复执行其他成功条目。

**AC-3.7.8**

**Given** 误报操作结束  
**When** 用户查看操作工作区或“操作记录”  
**Then** 展示每条漏洞的统一事实、原负责人、原状态、实际状态和最终处理状态  
**And** 当前会话与跨会话入口使用同一 Product Facts 和 ActionResult  
**And** 产品不主动发送通知。

**AC-3.7.9**

**Given** 正式 Story 实施完成  
**When** 执行自己负责、非自己负责、多条、超量、取消零写入、负责人变化、重复确认、部分失败、超时、状态不一致、重启恢复和跨会话找回测试  
**Then** schema、fixture、Go tests、前端测试、生成契约和批准环境验收全部通过  
**And** Workbench 与 Action API 对同一误报事实和状态投影一致  
**And** 验收记录不包含凭据、真实人员、真实漏洞或 raw payload。

### Target Candidate TC-3.3：完成延时与误报端到端验收

TargetRequirements:
  FR: [FR-11, FR-12, FR-13, FR-14, FR-15, FR-16, FR-17, FR-18]
  NFR: [NFR-01, NFR-02, NFR-03, NFR-04]
  Architecture: [AD-07, AD-10, AD-11, AD-23, AD-26]
  UX: [UX-DR-16, UX-DR-17, UX-DR-18, UX-DR-27, UX-DR-30]
  ConversionGates: [G-TOOLCHAIN, G-ARCH-V2, G-FACT-01, G-WRITE-01, G-WRITE-03, G-WRITE-04, G-SAFE-01, OQ-04]

As a 产品中心研发负责人，  
I want 用可复现的端到端测试验证修复延时和误报标记完整流程，  
So that Epic 3 只有在真实权限、写入、回读和恢复均有证据时才能完成。

**Acceptance Criteria:**

**AC-3.8.1**

**Given** 本 Story 当前为 Target Story Candidate  
**When** TC-3.1 或 TC-3.2 尚未正式转换并完成  
**Then** 本 Story 不得进入 Sprint、不得交给开发 Agent  
**And** blocked、skipped、mock PASS 或历史证据不能作为通过信号。

**AC-3.8.2**

**Given** TC-3.1、TC-3.2 已正式转换并完成且所有适用门禁仍有效  
**When** 转换本 Candidate  
**Then** 记录代码版本、门禁证据、环境授权、转换日期和批准人  
**And** 范围只包含端到端测试、浏览器验收、恢复验证和完成裁决  
**And** 不新增业务能力或复制通用 Action 状态机。

**AC-3.8.3**

**Given** 授权管理员和运营 actor 具有获准样本  
**When** 分别完成单条修复延时和人工误报旅程  
**Then** 两条旅程均经过授权事实、人工输入或判断、草案冻结、二次确认、单次写入和独立回读  
**And** 延时一次只处理一条并同时核验状态与期限  
**And** 误报只处理当前 actor 负责的对象并逐条核验状态。

**AC-3.8.4**

**Given** 草案处于未确认、取消、拒绝、过期、权限变化、负责人变化、状态变化、期限规则变化或 digest 不匹配状态  
**When** 执行安全失败矩阵  
**Then** 对应 provider mutation 调用次数为0  
**And** 不通过其他动作补偿或绕过确认  
**And** audit/replay 保存脱敏原因。

**AC-3.8.5**

**Given** provider 返回明确失败、超时、未知结果或回读不一致  
**When** 系统生成结果  
**Then** 已证明成功的对象只执行一次  
**And** 明确失败显示“待排查”  
**And** 未知或不一致时，在核验期限与调用预算内进入 `reconciling` 并显示“待核验”  
**And** 预算耗尽后仍无法证明目标事实时进入 `manual_attention` 并显示“待排查”  
**And** 延时缺少任一目标事实都不得显示“延时成功”  
**And** 误报状态不一致不得显示完成。

**AC-3.8.6**

**Given** 执行期间发生浏览器断线、服务重启或结果响应丢失  
**When** 系统恢复运行  
**Then** 继续引用原 ActionDraft、ActionResult 和 Product Facts  
**And** 不重复 mutation  
**And** 通过只读回读恢复能够证明的状态  
**And** 无法证明的对象先在核验期限与调用预算内保持 `reconciling`，仅在预算耗尽后进入 `manual_attention`。

**AC-3.8.7**

**Given** 操作完成或存在待处理对象  
**When** 用户离开会话、重新登录并打开“操作记录”  
**Then** 可以在六个月产品窗口内找到原操作  
**And** 操作工作区展示相同逐条事实、目标值、实际值和处理状态  
**And** Workbench 与 Action API 投影一致  
**And** 不重新调用 FOBrain 重建历史结果。

**AC-3.8.8**

**Given** 自动化和人工浏览器验收完成  
**When** 生成验收材料  
**Then** 覆盖 FR-11～18  
**And** 记录管理员与自己负责权限、人工决策、单条限制、确认后单次写入、状态／期限回读、部分失败、恢复、零写入和错误脱敏证据  
**And** 产品未发送任何消息、邮件、webhook、催办或通知中心通知  
**And** 报告和截图不包含凭据、真实人员、真实漏洞、raw payload 或内部引用。

**AC-3.8.9**

**Given** 所有端到端断言通过  
**When** 更新唯一 readiness gate  
**Then** 同步 G-WRITE-01、G-WRITE-03、G-WRITE-04、G-SAFE-01、OQ-04 与 Story 1.36 的最终状态与证据链接  
**And** 仅在延时和误报两条旅程均通过时允许 Epic 3 标记完成  
**And** 任一断言失败时 Epic 3 保持未完成并记录阻塞原因。
