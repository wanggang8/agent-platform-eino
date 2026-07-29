---
stepsCompleted:
  - step-01-validate-prerequisites
  - step-02-design-epics
  - step-03-create-stories
  - step-04-final-validation
inputDocuments:
  - _bmad-output/planning-artifacts/product-blueprint/prd.md
  - _bmad-output/planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14/ARCHITECTURE-SPINE.md
  - _bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/DESIGN.md
  - _bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/EXPERIENCE.md
  - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/SPEC.md
  - _bmad-output/planning-artifacts/sprint-change-proposal-2026-07-18.md
  - _bmad-output/planning-artifacts/implementation-readiness-report-2026-07-16.md
courseCorrection: 2026-07-18
---

# Agent Platform Eino - Epic Breakdown

## Overview

本文档把已批准的 FOBrain 漏洞处置实验需求、UX 契约、Architecture 不变量和 Product Facts v2 约束拆解为纵向、可独立验收的 Epic 与 Story。产品范围保持不变；本轮只纠正实施切片、依赖和门禁顺序。

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

AR-01：G-TOOLCHAIN 已由权威验收记录裁决 PASS；所有后续本地、CI 与发布证据必须继续使用 Go module 1.26.0、Go toolchain 1.26.5、Node 24 LTS 和 canonical linux/amd64 image，非目标工具链结果不得作为通过证据。

AR-02：保持六边形分层和明确 import boundary：`execution` 推进 Run/Action，`facts` 拥有领域与 repository ports，`product` 只读投影，`capabilities` 管注册和策略，`providers/*` 只接外部边界，`bootstrap` 是唯一组合根；禁止万能 `utils` 和跨层直调。

AR-03：自然语言 Agent、模型工具选择、ReAct loop、Runner event 与聊天 interrupt/checkpoint 使用 Eino v0.9.12；工具来自 Capability Registry。Product Facts、policy、approval、幂等、安全投影、HTTP/SSE、audit/replay 由项目控制面拥有。

AR-04：只有 `execution`、`capabilities`、`observability` 和 `store/sqlite` 的指定集成接缝可以直接 import Eino；`facts`、`product`、`httpapi`、`providers/*` 与 `web` 必须 Eino-free。

AR-05：每个纵向 Story 在启用其对应 runtime writer 或产品入口前，必须先交付该切片所需的 JSON Schema 2020-12、fixtures、OpenAPI 3.1、Go domain、repository ports、generated TypeScript types 和映射 contract tests。M-1 只建立 Product Facts/StructuredResult/QueryResultSnapshot/SSE 的最薄只读子集；Action/Resume/Result 契约随 M-4 mock Action 控制面建立。目标运行时始终不生成 v1/v2 union 或旧契约 adapter。

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

AR-29：实施按 M-0 已完成工具链、M-1 只读 walking skeleton、M-2 代表性真实读取、M-3 快照／引用／恢复与 READ-01、M-4 人员选择与 mock Action 控制面、M-5 目标部署 Gate Evidence、M-6 门禁后转换的 Target Story 推进；每个阶段必须形成用户可验证的纵向切片，不再按 schema、数据库、SSE、Eino 或 Action 控制面水平铺设大底座。

AR-30：必须随首次涉及相应行为的纵向 Story，以 schema、fixture、mock 和自动化失败注入覆盖 continuation 各崩溃点、重复确认、provider 响应丢失、回读最终一致、lease 失效、checkpoint 损坏、SSE 竞态/旧游标、fresh database bootstrap、legacy epoch 拒绝和备份恢复；未到对应 Story／Milestone 的测试必须明确 blocked/skipped，不能算 PASS。

AR-31：`implementation-readiness-gate.md` 是真实写域实施授权的唯一裁决源。G-TOOLCHAIN 已 PASS；当前只允许实现代表性只读事实链、人员安全读取／选择、mock Action 控制面和获授权 Gate Evidence。适用的 G-ARCH-V2、G-FACT、G-WRITE 与 G-SAFE 未有机器证据 PASS 前，不得实现或启用真实 mutation。

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

FR-02：Epic 1 — 建立列表、明细及后续动作共用的统一漏洞事实。

FR-03：Epic 1 — 区分真实空结果、无权限、数据不足和外部失败。

FR-04：Epic 1 — 精确查询全部新增漏洞并按发现时间倒序展示。

FR-05：Epic 2 — 运营先核对统一事实，再从 FOBrain 全员列表人工选择接收人。

FR-06：Epic 2 — 只合并人工判定为同一接收人的派发对象。

FR-07：Epic 2 — 派发经二次确认、逐条执行并按负责人回读核验。

FR-08：Epic 2 — 当前用户可主动转发其权限范围内的漏洞。

FR-09：Epic 2 — 转发从 FOBrain 全员列表选择接收人，并按接收人拆分草案。

FR-10：Epic 2 — 转发经二次确认并逐条回读目标负责人。

FR-11：Epic 3 — 延时决定和新期限始终由人工沟通确定。

FR-12：Epic 3 — 管理员一次确认并延时一条漏洞，原因可空。

FR-13：Epic 3 — 仅状态和期限均回读一致时显示延时成功。

FR-14：Epic 3 — 操作人可对自己负责的漏洞作人工误报判断。

FR-15：Epic 3 — 误报确认后写入，且仅在状态回读一致时成功。

FR-16：Epic 2 首次建立，Epic 3 复用并重新验收 — 所有写动作确认前、取消、拒绝或过期时零外部写入。

FR-17：Epic 2 首次建立，Epic 3 复用并重新验收 — 所有写动作逐条投影完成、待核验和待排查，禁止重复写入核验。

FR-18：Epic 2 首次建立，Epic 3 复用并重新验收 — 产品不发送通知，也不把 FOBrain 自身通知作为成功条件。

FR-19：Epic 1 — 读取当前 FOBrain 用户的安全身份上下文。

FR-20：Epic 1 — 读取当前用户的 FOBrain 权限和数据范围三态。

FR-21：Epic 1 — 按 IP 查询资产及安全在线状态。

FR-22：Epic 1 — 通过安全 opaque ref 查看单个资产详情。

FR-23：Epic 1 — 按 IP 查询关联漏洞列表。

FR-24：Epic 1 — 通过安全 opaque ref 查看单个漏洞详情。

FR-25：Epic 1 — 查询 `business_list` 返回的授权业务系统列表。

## Epic List

### Epic 1：可信、可恢复的 FOBrain 事实工作台

运营可以在浅色桌面聊天工作台中安全查询新增漏洞、当前身份权限、资产、漏洞和业务系统，通过冻结快照、百条分页、确定性引用与断线／重启恢复持续核对同一结果；所有真实写入保持关闭。

**FRs covered:** FR-01、FR-02、FR-03、FR-04、FR-19、FR-20、FR-21、FR-22、FR-23、FR-24、FR-25。

**Implementation boundary:** 以 read-only walking skeleton 开始，按用户可见纵向切片扩展真实代表性读取、统一事实、引用和恢复；不提前建设 ActionDraft、mutation、verifier 或真实写入口。

**Epic completion rule:** Stories 1.1～1.12 全部完成且 READ-01=PASS 后，Epic 1 才可标记 done。

### Epic 2：运营人工派发与责任人转发

运营可以读取并人工选择 FOBrain 人员，基于 Epic 1 的冻结漏洞事实准备、确认、执行和核验派发／转发，逐条看到完成、待核验和待排查结果，并可跨会话重新找到同源操作记录。

**FRs covered:** FR-05、FR-06、FR-07、FR-08、FR-09、FR-10、FR-16、FR-17、FR-18。

**Implementation boundary:** 先关闭 OQ-02，再通过 mock 派发纵向建立共用 Action 控制面、零写入、幂等、lease、verifier、reconcile 和操作记录；真实派发／转发仅在适用门禁 PASS 后由 Create Story 生成新的 M-6 正式 Target Story。

**Epic completion rule:** Stories 2.1～2.10 只完成 M-4/M-5 控制面与取证段；Epic 2 的产品能力只有在派发、转发各自适用门禁 PASS 且相应 M-6 Target Stories 完成后才可标记 done。任一动作保持 BLOCKED 时不得用 mock 或 Gate Evidence 宣称真实能力完成。

### Epic 3：管理员延时与运营误报处置

管理员可以对人工沟通确定的一条漏洞执行修复延时，运营可以对自己负责且已人工判断的漏洞标记误报；两类动作复用 Epic 2 的可信确认、执行、核验和操作记录，只有目标事实回读一致才显示成功。

**FRs covered:** FR-11、FR-12、FR-13、FR-14、FR-15，并复用／重新验收 FR-16、FR-17、FR-18。

**Implementation boundary:** 依赖 Epic 1 的统一只读事实链，并仅以 Epic 2 Stories 2.3～2.8 已完成的共用 Action 控制面、状态投影和操作记录为前置；不等待 Stories 2.9／2.10 的派发／转发 Gate Evidence，也不等待任何 M-6 Target Story。Epic 3 只增加延时／误报专属角色、基数、参数、目标事实与 verifier 规则，不复制状态机；真实延时／误报仍仅在自身适用门禁 PASS 后由 Create Story 生成新的 M-6 正式 Target Story。

**Epic completion rule:** Stories 3.1～3.5 只完成规则、mock 旅程与取证段；Epic 3 的产品能力只有在延时、误报各自适用门禁 PASS 且相应 M-6 Target Stories 完成后才可标记 done。一个动作 BLOCKED 不阻止另一个动作独立取证，但也不能把 Epic 整体标记 done。

## Epic 1：可信、可恢复的 FOBrain 事实工作台

运营可以在浅色桌面聊天工作台中安全查询新增漏洞、当前身份权限、资产、漏洞和业务系统，通过冻结快照、百条分页、确定性引用与断线／重启恢复持续核对同一结果；所有真实写入保持关闭。

### Story 1.1：固定可复现工具链与开工门禁

As a 交付负责人，
I want 使用唯一受支持的 Go、Node 与 CI 工具链产生验收证据，
So that 后续所有实现和裁决来自一致、可复现的环境。

**WorkItemType:** enabler / done

**Requirements:** NFR-01；AR-01、AR-29、AR-30；G-TOOLCHAIN；M-0。

**Architecture decisions:** AD-18、AD-19。

**Prerequisites:** 无。

**Inputs / Outputs:** 输入为固定工具链版本、canonical image 与基线检查清单；输出为可复现的本地/CI 证据和 Story 1.1 done、G-TOOLCHAIN PASS 的权威记录。

**Exit gate:** 两条验收命令在目标工具链通过，权威记录与当前门禁状态一致；否则不得进入 Story 1.2。

**Scope:** 固定 Go module 1.26.0、Go toolchain 1.26.5、Node 24 LTS、canonical linux/amd64 image、版本检查与权威验收记录。

**Non-goals:** 不实现任何产品能力，不接受其他工具链结果作为通过证据。

**Affected directories:** go.mod、go.work、.github/workflows/、build/toolchain/、scripts/、docs/acceptance-records/。

**Acceptance commands:** bash scripts/verify_toolchain.sh；bash scripts/run_toolchain_baseline.sh。

**Acceptance Criteria:**

**Given** 新环境、本机或 CI 开始构建  
**When** 执行工具链基线校验  
**Then** 只有固定 Go、Node 与 canonical image 的完整结果可发布为 G-TOOLCHAIN PASS 证据  
**And** 其他版本必须明确失败或标记为非通过证据。

**Given** Story 1.1 已有权威验收记录  
**When** 当前计划读取门禁状态  
**Then** Story 1.1 必须标记 done，G-TOOLCHAIN 必须标记 PASS  
**And** 历史失败只保留为历史记录，不覆盖当前裁决。

### Story 1.2：用单一安全 Fixture 跑通只读 Walking Skeleton

As a 产品中心运营，
I want 在浅色桌面 Workbench 中发起一次固定的新增漏洞查询并看到摘要，
So that 我能尽早验证从聊天入口到安全事实展示的完整产品链路。

**WorkItemType:** user-value

**Requirements:** FR-01、FR-02、FR-03、FR-04；NFR-01、NFR-03；AR-02、AR-03、AR-04、AR-05、AR-06、AR-07、AR-08、AR-17、AR-19、AR-20、AR-21、AR-22、AR-23、AR-24、AR-25、AR-28、AR-29、AR-30、AR-31、AR-36、AR-37；UX-DR-01、UX-DR-02、UX-DR-03、UX-DR-04、UX-DR-05、UX-DR-06、UX-DR-07、UX-DR-19、UX-DR-20、UX-DR-21、UX-DR-22、UX-DR-23、UX-DR-24、UX-DR-25、UX-DR-28、UX-DR-29；G-ARCH-V2；M-1。

**Architecture decisions:** AD-01、AD-02、AD-03、AD-04、AD-05、AD-13、AD-16、AD-17、AD-18、AD-20、AD-24、AD-25。

**Prerequisites:** Story 1.1。

**Inputs / Outputs:** 输入为一个版本化新增漏洞安全 fixture、全新数据库和固定聊天查询；输出为一个持久化 StructuredResult/QueryResultSnapshot、同源 Product Facts/SSE 投影及一张浅色桌面结果卡。

**Exit gate:** 单 fixture 纵向 E2E、schema、Go、TypeScript、browser 与零 mutation 断言全部通过，且未接真实 FOBrain、真实 LLM、Action 或历史数据；否则不得进入 Story 1.3。

**Scope:** 一个版本化安全 fixture、一个只读新增漏洞 capability、最小 Product Facts v2/StructuredResult/QueryResultSnapshot、Greenfield SQLite、mock Eino loop、同源 product projection/SSE、一个浅色桌面页面、一个查询结果卡和右栏“事实／执行记录”空态。

**Non-goals:** 不接真实 FOBrain、真实 LLM、ActionDraft、确认／写入、跨会话操作记录数据、完整分页、完整引用恢复或旧契约兼容。

**Affected directories:** docs/schemas/、docs/fixtures/、internal/einoapp/{bootstrap,httpapi,execution,facts,product,capabilities,store/sqlite}/、web/eino-workbench/src/、scripts/。

**Acceptance commands:** node scripts/eino_workbench_schema_validate.mjs；go test ./...；npm --prefix web/eino-workbench run typecheck；npm --prefix web/eino-workbench test；npm --prefix web/eino-workbench run browser-test。

**Acceptance Criteria:**

**Given** 全新数据库和唯一批准的新增漏洞安全 fixture  
**When** 用户在桌面聊天入口提交固定查询  
**Then** mock Eino loop 只能选择一个注册的只读 capability，结果经 Safety Gate 成为唯一 StructuredResult 与 QueryResultSnapshot 后持久化  
**And** Workbench 从同一 Product Facts 投影显示查询序号、摘要、数量和观察时间。

**Given** fixture 分别表示有结果、0 条和失败  
**When** 查询完成  
**Then** 0 条显示“没有待派发漏洞”，失败显示确定性失败语义，两者不得混淆  
**And** raw payload、token、provider locator、内部 result_ref 和旧 runtime 类型不得进入产品出口。

**Given** 用户查看首版桌面页面  
**When** 浏览三栏结构与键盘焦点  
**Then** 页面仅使用浅色桌面布局，左栏固定提供“操作记录”入口，中栏显示聊天，右栏仅显示安全“事实／执行记录”空态  
**And** 操作记录入口不查询历史数据，不出现写入入口，真实 FOBrain mutation 调用数为 0。

**Given** 服务重启后重新打开该结果  
**When** repository 恢复最小 Product Facts  
**Then** 同一结果仍由持久事实安全投影  
**And** 不要求实现完整 SSE 恢复、上下文引用或历史操作列表。

### Story 1.3：验证代表性读取与统一事实的数据可得性

As a 产品负责人，
I want 用目标部署证据确认精确新增查询和七项代表性读取能提供哪些字段，
So that 后续真实页面不会靠 AI 推断或硬编码补齐业务事实。

**WorkItemType:** gate-evidence

**Requirements:** FR-01、FR-02、FR-03、FR-04、FR-19、FR-20、FR-21、FR-22、FR-23、FR-24、FR-25；NFR-01、NFR-02、NFR-04；AR-26、AR-30、AR-31、AR-32、AR-37、AR-38、AR-43、AR-44；UX-DR-10、UX-DR-11、UX-DR-12、UX-DR-13、UX-DR-20；G-READ-01、G-READ-03、G-FACT 数据可得性；M-2。

**Architecture decisions:** AD-03、AD-05、AD-13、AD-14、AD-15、AD-16、AD-27。

**Prerequisites:** Story 1.2。

**Inputs / Outputs:** 输入为精确新增查询、七项代表性读取的目标合同、获准凭据与脱敏样本；输出为逐字段 resolved/empty/unavailable 证据矩阵、fixture/live assertions 及 G-READ/G-FACT 数据可得性裁决。

**Exit gate:** Story 必须产出可复核的 PASS/BLOCKED 证据记录；只有 Story 1.4 所需身份与权限字段被证明可安全投影时才允许继续，否则后续对应能力保持阻塞。

**Scope:** 精确新增漏洞查询，以及 current_user_context、my_permissions、list_assets_by_ip、get_asset_detail、list_vulnerabilities_by_ip、get_vulnerability_detail、business_list 的 schema、fixture、mock/live assertions、字段来源与敏感字段审计。

**Non-goals:** 不实现人员列表，不执行 mutation，不把“目标 API 可访问”声明为产品已集成。

**Affected directories:** docs/fobrain-tool-matrix.md、docs/schemas/、docs/fixtures/、internal/einoapp/providers/fobrain/、scripts/acceptance/、docs/acceptance-records/。

**Acceptance commands:** node scripts/eino_workbench_schema_validate.mjs；go test ./internal/einoapp/providers/...；bash scripts/acceptance/story_1_3_read_evidence.sh（未获授权或证据不足时必须 exit 2）。

**Acceptance Criteria:**

**Given** 目标部署凭据与获准只读窗口  
**When** 逐项采集八项读取需求的脱敏样本  
**Then** 每个产品字段必须被裁决为 resolved、confirmed_empty/no_permission 或 source_field_unavailable/unknown  
**And** unknown 不得折叠为无权限、空或通过。

**Given** 查询包含多页、0 条、权限拒绝、字段缺失或外部失败  
**When** 生成 evidence record  
**Then** complete_set、partial_page/incomplete、空结果和失败具有可机器核验的不同证据  
**And** raw payload、token、真实账号和内部 locator 不得写入记录。

**Given** 必要统一漏洞字段无法从目标部署证明  
**When** 更新 G-READ/G-FACT 数据可得性状态  
**Then** 对应门禁保持 BLOCKED 并明确缺口  
**And** Story 1.4 之后不得通过模型推断填补缺失字段。

### Story 1.4：查看当前用户与权限三态

As a 产品中心运营，
I want 查看当前 FOBrain 身份、部门、角色和数据权限范围，
So that 我知道页面事实属于谁以及查询覆盖什么范围。

**WorkItemType:** user-value

**Requirements:** FR-01、FR-03、FR-19、FR-20；NFR-01、NFR-02、NFR-04；AR-14、AR-26、AR-27、AR-43；UX-DR-05、UX-DR-11、UX-DR-13、UX-DR-20、UX-DR-22、UX-DR-23；G-READ-03；M-2。

**Architecture decisions:** AD-09、AD-14、AD-15、AD-16。

**Prerequisites:** Story 1.3 的 current_user_context 与 my_permissions 数据证据允许继续。

**Inputs / Outputs:** 输入为当前 workspace FOBrain 凭据以及两项读取的安全映射；输出为不可覆盖的 actor、身份/部门/角色/权限 Product Facts 与三态桌面投影。

**Exit gate:** actor 来源、三态语义、错误脱敏和敏感字段泄漏测试全部通过；否则不得进入 Story 1.5。

**Scope:** 两项真实只读 capability、actor 解析、安全身份/权限 Product Facts、三态投影及桌面右栏展示。

**Non-goals:** 不展示 token、raw account/policy payload，不允许前端或模型覆盖 actor，不创建人员选择或写入口。

**Affected directories:** internal/einoapp/{capabilities,facts,product,providers/fobrain}/、web/eino-workbench/src/、docs/fixtures/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；npm --prefix web/eino-workbench run typecheck。

**Acceptance Criteria:**

**Given** 当前 workspace 的 FOBrain 凭据  
**When** 查询 current_user_context 与 my_permissions  
**Then** actor 只能由 provider instance 与 stable business user ID 解析，界面展示获准的安全身份、部门和角色摘要  
**And** HTTP、前端和模型提交的 actor 覆盖值必须被忽略或拒绝。

**Given** 部门、权限或数据范围分别为已解析、确认空／无权限或源字段不可用／未知  
**When** 投影到聊天结果和右栏事实  
**Then** 三种状态使用稳定中文语义并保持彼此不同  
**And** 未提供字段不得留空、推断或伪报为已解析。

**Given** provider 超时、拒绝或返回不可信字段  
**When** Safety Gate 处理结果  
**Then** 页面显示脱敏外部失败且不保留旧成功值冒充当前事实  
**And** token、raw account、policy payload 与凭据不得进入日志、模型上下文或前端契约。

### Story 1.5：查询全部新增漏洞并区分完整、空、部分与失败

As a 产品中心运营，
I want 按发现时间倒序查看当前权限内的全部新增漏洞，
So that 我每天能够准确识别需要人工核对的对象。

**WorkItemType:** user-value

**Requirements:** FR-01、FR-02、FR-03、FR-04；NFR-01、NFR-02、NFR-04；AR-08、AR-26、AR-31、AR-37、AR-38、AR-42；UX-DR-06、UX-DR-07、UX-DR-10、UX-DR-11、UX-DR-12、UX-DR-20、UX-DR-22；G-READ-01；M-2。

**Architecture decisions:** AD-03、AD-05、AD-13、AD-14、AD-16、AD-27。

**Prerequisites:** Story 1.4。

**Inputs / Outputs:** 输入为当前 actor、精确新增查询条件和 provider 全部分页；输出为稳定倒序的 complete_set/空快照，或明确的 incomplete、无权限、外部失败事实。

**Exit gate:** 有结果、0 条、分页失败、无权限和外部失败五类自动化断言通过，且 ActionDraft/mutation 均为 0；否则不得进入 Story 1.6。

**Scope:** 真实精确新增查询 capability、全部分页读取、稳定倒序、coverage、安全漏洞基础事实、空态和错误态。

**Non-goals:** 不补齐资产/业务系统扩展字段，不提供“全部处理”或任何写动作。

**Affected directories:** internal/einoapp/{capabilities,facts,product,providers/fobrain}/、web/eino-workbench/src/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_1_5_new_vulnerabilities_smoke.sh。

**Acceptance Criteria:**

**Given** 当前 actor 可见的新增漏洞跨多个 provider 页  
**When** 用户请求“查询新增漏洞”  
**Then** 只有全部分页成功才发布 complete_set 快照，按发现时间降序和稳定次级键排序  
**And** 结果卡显示查询序号、固定数量与观察时间。

**Given** 全部分页成功且结果为 0 条  
**When** 查询完成  
**Then** 发布可重开的空 complete_set 快照并显示“没有待派发漏洞”  
**And** 不创建动作草案或触发 mutation。

**Given** 任一页失败、权限未知或外部调用失败  
**When** 查询结算  
**Then** partial_page/incomplete、无权限和外部失败分别展示，并明确未核验范围  
**And** 不得由 AI 改写为完整集合或空结果。

### Story 1.6：按 IP 查询资产并打开安全详情

As a 产品中心运营，
I want 按 IP 查看资产列表、在线状态和单个资产安全详情，
So that 我能核对主机是否真实在线及其资产上下文。

**WorkItemType:** user-value

**Requirements:** FR-21、FR-22；NFR-01、NFR-02、NFR-04；AR-26、AR-43、AR-44；UX-DR-10、UX-DR-11、UX-DR-12、UX-DR-13、UX-DR-19、UX-DR-22；G-READ-03；M-2。

**Architecture decisions:** AD-03、AD-05、AD-14、AD-16、AD-27。

**Prerequisites:** Story 1.5。

**Inputs / Outputs:** 输入为获授权 IP 或资产 opaque entity_ref；输出为资产列表、稳定在线状态、安全详情及受控 locator 解析结果。

**Exit gate:** list/detail、空/权限/失败、越权/过期引用和 key rotation 测试全部通过，内部身份材料零泄漏；否则不得进入 Story 1.7。

**Scope:** list_assets_by_ip、get_asset_detail、在线状态安全映射、opaque entity_ref 与受控 ProviderLocator resolver。

**Non-goals:** 不把 safe ref 当 provider asset ID，不展示 locator、subject key 或 fingerprint，不修改资产。

**Affected directories:** internal/einoapp/{capabilities,facts,product,providers/fobrain}/、web/eino-workbench/src/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_1_6_asset_list_detail_smoke.sh。

**Acceptance Criteria:**

**Given** 用户输入获授权 IP  
**When** 查询资产列表  
**Then** 每条只展示安全资产事实和稳定文字化在线状态  
**And** 空、无权限、字段不可用和外部失败保持不同状态。

**Given** 用户从列表打开一条资产详情  
**When** 服务端解析 opaque entity_ref  
**Then** 只能按 workspace、actor、provider instance、policy version 与有效期解析到受控 locator  
**And** safe ref、entity_subject_key、fingerprint alias 和 locator 均不得出现在产品出口。

**Given** subject 或 locator key 轮换  
**When** 再次查询同一资产  
**Then** 主体身份与 locator lineage 保持一致  
**And** 旧别名只在受控保留窗口内用于解析。

### Story 1.7：按 IP 查询漏洞并打开安全详情

As a 产品中心运营，
I want 按 IP 查看关联漏洞并打开单个漏洞的安全详情，
So that 我能在处置前核对完整漏洞事实。

**WorkItemType:** user-value

**Requirements:** FR-23、FR-24；NFR-01、NFR-02、NFR-04；AR-26、AR-43、AR-44；UX-DR-09、UX-DR-10、UX-DR-11、UX-DR-12、UX-DR-13、UX-DR-19、UX-DR-22；G-READ-03；M-2。

**Architecture decisions:** AD-03、AD-05、AD-14、AD-16、AD-27。

**Prerequisites:** Story 1.6。

**Inputs / Outputs:** 输入为获授权 IP 或漏洞 opaque entity_ref；输出为关联漏洞列表、统一安全详情和稳定漏洞主体解析结果。

**Exit gate:** list/detail、字段缺失、权限、失败、多 IP/长 POC、越权引用和 alias rotation 测试全部通过；否则不得进入 Story 1.8。

**Scope:** list_vulnerabilities_by_ip、get_vulnerability_detail、列表/详情统一安全事实与 opaque 引用解析。

**Non-goals:** 不准备、确认或执行漏洞动作，不允许列表与详情使用两套事实模型。

**Affected directories:** internal/einoapp/{capabilities,facts,product,providers/fobrain}/、web/eino-workbench/src/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_1_7_vulnerability_list_detail_smoke.sh。

**Acceptance Criteria:**

**Given** 用户输入获授权 IP  
**When** 查询关联漏洞  
**Then** 返回 POC、全部安全 IP、在线状态、发现时间和当前负责人等已证明字段  
**And** 缺失字段使用固定语义，不留空、不猜测。

**Given** 用户从列表打开漏洞详情  
**When** 服务端解析 entity_ref 并读取详情  
**Then** 详情必须与列表引用同一 entity_subject 和字段语义  
**And** 越权、过期、解析失败和外部失败均 fail closed。

**Given** 列表或详情包含多 IP、长 POC 或敏感 provider 字段  
**When** 投影到表格和右栏  
**Then** accessible name 与右栏提供完整安全值  
**And** raw payload、内部身份材料与未授权字段被移除。

### Story 1.8：查询全部业务系统并补齐统一漏洞行事实

As a 产品中心运营，
I want 查看 FOBrain 全部业务系统，并让漏洞行统一显示业务和负责人上下文，
So that 我能基于同一组事实判断每条漏洞。

**WorkItemType:** user-value

**Requirements:** FR-02、FR-05、FR-25；NFR-01、NFR-02、NFR-04；AR-08、AR-26、AR-32、AR-39、AR-41、AR-43、AR-44；UX-DR-10、UX-DR-11、UX-DR-12、UX-DR-13、UX-DR-19、UX-DR-20；G-READ-03、G-FACT-01；M-2。

**Architecture decisions:** AD-03、AD-05、AD-13、AD-14、AD-16、AD-27。

**Prerequisites:** Story 1.7。

**Inputs / Outputs:** 输入为 business_list 以及漏洞、资产、业务和负责人安全事实；输出为业务系统安全列表、字段来源摘要和确定性统一漏洞行事实。

**Exit gate:** 业务列表权限裁剪、多源去重、缺失固定文案和冲突 fail-closed 测试通过，G-FACT-01 对必要统一字段无未记录缺口；否则不得进入 Story 1.9。

**Scope:** business_list、业务系统安全投影、多源字段来源、稳定对象身份、统一漏洞行事实及缺失/冲突语义。

**Non-goals:** 不实现“我的业务系统”，不关注无关联 IP 的业务系统统计，不用 last-write-wins 解决冲突。

**Affected directories:** internal/einoapp/{capabilities,facts,product,providers/fobrain}/、web/eino-workbench/src/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_1_8_business_unified_fact_smoke.sh。

**Acceptance Criteria:**

**Given** 当前凭据调用获批准的 business_list capability  
**When** FOBrain 返回当前 actor 可见的业务系统  
**Then** 产品展示接口内置权限裁剪后的全部业务系统安全列表  
**And** 不替换为已排除的 my_business_systems，不暴露 raw payload 或 locator。

**Given** 漏洞、资产、业务系统和负责人事实来自多个安全结果  
**When** 组装统一漏洞行  
**Then** 每行至少稳定呈现 POC、IP、在线状态、业务系统、当前修复负责人、业务系统负责人、运维负责人和发现时间  
**And** 当前负责人缺失显示“未分配”，业务及负责人缺失显示“未提供”。

**Given** 多个来源指向同一对象或字段冲突  
**When** 进行确定性去重和字段合并  
**Then** 使用稳定业务身份并保留字段来源摘要  
**And** 冲突显示数据不足或要求澄清，不得静默覆盖。

### Story 1.9：在冻结快照中每页 100 条核对完整事实

As a 产品中心运营，
I want 在专门的漏洞明细工作区分页查看冻结结果，
So that 百条以上对象也能稳定、高效地逐条核对。

**WorkItemType:** user-value

**Requirements:** FR-02、FR-04；NFR-01、NFR-04；AR-07、AR-08、AR-19、AR-24、AR-25、AR-36、AR-37、AR-38；UX-DR-03、UX-DR-04、UX-DR-06、UX-DR-07、UX-DR-08、UX-DR-09、UX-DR-10、UX-DR-11、UX-DR-12、UX-DR-13、UX-DR-23、UX-DR-24、UX-DR-25、UX-DR-26、UX-DR-29；G-FACT-01；M-3。

**Architecture decisions:** AD-03、AD-04、AD-05、AD-16、AD-17、AD-20、AD-24、AD-27。

**Prerequisites:** Story 1.8。

**Inputs / Outputs:** 输入为一个持久化 complete_set 统一漏洞快照与页面导航意图；输出为不可变 IA-02 百条分页、grid/右栏安全投影和可恢复页面状态。

**Exit gate:** 100+ 条分页、键盘/ARIA、页面失败、返回定位、快照不漂移和只读边界 browser tests 全部通过；否则不得进入 Story 1.10。

**Scope:** 不可变 QueryResultSnapshot、完整安全 JSON 持久化、IA-02 每页 100 条、grid 键盘交互、右栏事实、页面/行/滚动恢复。

**Non-goals:** 不以勾选改变动作范围，不从明细页发起写入，不让新到漏洞改变已打开快照。

**Affected directories:** internal/einoapp/{facts,product,httpapi,store/sqlite}/、web/eino-workbench/src/、docs/schemas/、docs/fixtures/。

**Acceptance commands:** node scripts/eino_workbench_schema_validate.mjs --contract；go test ./internal/einoapp/...；npm --prefix web/eino-workbench run typecheck；npm --prefix web/eino-workbench test；npm --prefix web/eino-workbench run browser-test。

**Acceptance Criteria:**

**Given** 一个包含超过 100 条统一漏洞事实的 complete_set 快照  
**When** 用户从聊天结果卡选择“查看全部”  
**Then** IA-02 按冻结顺序每页 100 条显示固定总数、总页数、当前页与加载状态  
**And** 首页、上下页、末页和页码跳转不会改变 result_ref 或冻结集合。

**Given** 用户使用键盘浏览 grid、选择行并阅读右栏  
**When** 方向键、Home/End、PageUp/PageDown 或 Enter 被触发  
**Then** grid 使用单一 Tab 停靠点、正确 aria-rowcount/aria-rowindex，行选择只更新右栏  
**And** 不出现复选框式动作范围或任何写入口。

**Given** 页面加载失败、返回聊天或再次进入同一快照  
**When** 恢复表面状态  
**Then** 保留最后核验页、当前行、滚动位置和来源卡焦点，页失败明确“不是空结果”  
**And** 重新加载本页只读取同一冻结快照。

### Story 1.10：确定性解析多次查询结果指代

As a 对话用户，
I want 用“这些”“刚才第 2 次结果”等表达准确引用已有查询，
So that 多次查询或上下文压缩后系统仍锁定我指的对象集合。

**WorkItemType:** user-value

**Requirements:** FR-01、FR-02、FR-03；NFR-01、NFR-02；AR-08、AR-09、AR-39、AR-44；UX-DR-06、UX-DR-07、UX-DR-14、UX-DR-20、UX-DR-22、UX-DR-27；G-FACT-01；M-3。

**Architecture decisions:** AD-03、AD-05、AD-06、AD-09、AD-16、AD-27。

**Prerequisites:** Story 1.9。

**Inputs / Outputs:** 输入为当前 scope 的 Result Index、候选 result_ref 和用户指代表达；输出只能是 resolved、clarification_required 或 invalid 及对应安全澄清投影。

**Exit gate:** 唯一、多候选、跨 scope、过期、篡改和上下文压缩测试全部通过，模型不能覆盖程序解析结果；否则不得进入 Story 1.11。

**Scope:** Result Index、稳定 result_ref、scope/expiry 校验、resolved/clarification_required/invalid 三态、澄清卡和安全模型上下文。

**Non-goals:** 不让模型拥有最终选择权，不从聊天摘要重建集合，不暴露 result_ref。

**Affected directories:** internal/einoapp/{execution,facts,product}/、web/eino-workbench/src/、docs/schemas/、docs/fixtures/。

**Acceptance commands:** go test ./internal/einoapp/execution/... ./internal/einoapp/facts/... ./internal/einoapp/product/...；npm --prefix web/eino-workbench test。

**Acceptance Criteria:**

**Given** 同一 scope 只有一个可验证候选  
**When** 用户说“这些漏洞”  
**Then** resolver 可自动返回 resolved 并由程序绑定对应 result_ref  
**And** 模型只收到安全摘要，不拥有或覆盖最终引用。

**Given** 同一 scope 存在多个合理候选  
**When** 用户使用模糊指代  
**Then** resolver 返回 clarification_required，澄清卡只展示查询序号、摘要、数量和时间  
**And** 用户选择前不得创建动作草案或写入。

**Given** 引用跨 workspace/actor/conversation、已过期、不可验证或被篡改  
**When** resolver 校验  
**Then** 返回 invalid 并要求重新查询  
**And** 不回退到最近一次结果或聊天摘要。

### Story 1.11：在刷新、断线和服务重启后恢复同一只读事实

As a 产品中心运营，
I want 页面刷新、SSE 断线或服务重启后继续查看同一查询，
So that 我不会因连接变化丢失或误读冻结事实。

**WorkItemType:** user-value

**Requirements:** FR-01、FR-02、FR-03；NFR-01；AR-07、AR-17、AR-18、AR-19、AR-20、AR-21、AR-22、AR-23、AR-24、AR-25；UX-DR-05、UX-DR-22、UX-DR-24、UX-DR-26、UX-DR-27；G-FACT-01；M-3。

**Architecture decisions:** AD-03、AD-04、AD-05、AD-12、AD-16、AD-18、AD-20、AD-24、AD-25。

**Prerequisites:** Story 1.10。

**Inputs / Outputs:** 输入为持久化 Product Facts、FactEvent、snapshot_sequence 与有效/旧/未知 SSE cursor；输出为同源恢复 view、确定 patch 或 view.replaced。

**Exit gate:** transaction 同成同败、断线、乱序、刷新和服务重启矩阵通过，同一 result_ref/范围不漂移；否则不得进入 Story 1.12。

**Scope:** 聚合与 FactEvent 原子提交、Run 内单调 sequence、持久 view/SSE cursor、view.replaced、read-only continuation 恢复和 Greenfield restart。

**Non-goals:** 不实现 Action attempt、lease、verifier 或操作记录数据；不自动清理 terminal Product Facts。

**Affected directories:** internal/einoapp/{execution,facts,product,httpapi,store/sqlite}/、web/eino-workbench/src/、scripts/acceptance/。

**Acceptance commands:** go test ./internal/einoapp/...；node scripts/eino_workbench_stream_test.mjs；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_1_11_read_restart_smoke.sh。

**Acceptance Criteria:**

**Given** 查询结果、聚合状态和 FactEvent 待提交  
**When** SQLite transaction 成功或失败  
**Then** 聚合与事件必须同成同败并获得 Run 内单调 sequence  
**And** 不存在仅更新 view 或仅追加事件的中间真相。

**Given** 客户端携带有效、旧或未知 SSE 游标连接  
**When** 服务端发布当前 view  
**Then** 有效游标只收到 snapshot_sequence 之后的完整 upsert/tombstone，旧或未知游标收到原子 view.replaced  
**And** 重复或乱序事件不能使前端状态回退。

**Given** 浏览器刷新或服务在查询完成后重启  
**When** 用户重新进入会话和快照  
**Then** 同一 Product Facts 与 StructuredResult 恢复原 result_ref、冻结范围和页面位置  
**And** 不从模型上下文、聊天摘要或 provider raw payload 重建事实。

### Story 1.12：完成 READ-01 只读端到端验收

As a 产品负责人，
I want 用机器证据验收真实只读主流程并确认写域仍关闭，
So that 只有可用、可恢复且安全的事实工作台才能进入 Action 阶段。

**WorkItemType:** acceptance

**Requirements:** FR-01、FR-02、FR-03、FR-04、FR-19、FR-20、FR-21、FR-22、FR-23、FR-24、FR-25；NFR-01、NFR-02、NFR-04；AR-29、AR-30、AR-31、AR-32、AR-37、AR-38、AR-43、AR-44；UX-DR-01、UX-DR-02、UX-DR-03、UX-DR-04、UX-DR-05、UX-DR-06、UX-DR-07、UX-DR-08、UX-DR-09、UX-DR-10、UX-DR-11、UX-DR-12、UX-DR-13、UX-DR-14、UX-DR-19、UX-DR-20、UX-DR-21、UX-DR-22、UX-DR-23、UX-DR-24、UX-DR-25、UX-DR-26、UX-DR-27、UX-DR-28、UX-DR-29；G-TOOLCHAIN、G-ARCH-V2、G-READ-01、G-READ-03、G-FACT-01、READ-01；M-3。

**Architecture decisions:** AD-01、AD-02、AD-03、AD-04、AD-05、AD-06、AD-13、AD-14、AD-15、AD-16、AD-17、AD-18、AD-19、AD-20、AD-24、AD-25、AD-27。

**Prerequisites:** Stories 1.1～1.11。

**Inputs / Outputs:** 输入为 Stories 1.1～1.11 的代码、契约和证据；输出为 READ-01 端到端报告、更新后的只读门禁和可复现失败记录。

**Exit gate:** 全部验收命令与 READ-01 产品旅程通过，G-TOOLCHAIN/G-ARCH-V2/适用 G-READ/G-FACT 状态一致且真实 mutation 为 0；任何失败都阻止 Story 2.1。

**Scope:** schema/fixture/contract、Go、TypeScript、SSE、browser、SQLite/restart、安全泄漏、代表性真实读取和零 mutation 的完整 READ-01 证据。

**Non-goals:** 不验收人员选择、动作草案或任何真实写能力；不以目标 API 单测代替产品 E2E。

**Affected directories:** docs/acceptance-records/、_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md、scripts/、scripts/acceptance/、全体 Epic 1 受影响目录。

**Acceptance commands:** bash scripts/verify_toolchain.sh；node scripts/eino_workbench_schema_validate.mjs --contract；go test ./...；npm --prefix web/eino-workbench run typecheck；npm --prefix web/eino-workbench test；npm --prefix web/eino-workbench run browser-test；bash scripts/acceptance/read_01_product_smoke.sh。

**Acceptance Criteria:**

**Given** Stories 1.1～1.11 和适用只读门禁具有有效证据  
**When** 执行“身份权限—新增漏洞—资产/漏洞详情—业务系统—统一行事实—百条快照—结果引用—刷新/断线/重启恢复”E2E  
**Then** 所有表面只能从同源 Product Facts/StructuredResult 投影，客观区分有结果、0 条、无权限、字段不可用、部分分页和外部失败  
**And** 浏览器证据满足浅色桌面、键盘、焦点、WCAG 2.2 AA 和安全字段约束。

**Given** 用户表达派发、转发、延时或误报意图  
**When** READ-01 环境处理该消息  
**Then** 只显示“当前仅支持查询与只读核对，此写操作尚未启用”  
**And** ActionDraft 数、approval 数和真实 FOBrain mutation 调用数均为 0。

**Given** 任一 schema、fixture、恢复、安全或产品 E2E 断言失败  
**When** 更新 READ-01 与 readiness 状态  
**Then** Epic 1 保持未完成并记录可复现缺口  
**And** 不允许进入 Story 2.1 或把失败项记录为 PASS。

## Epic 2：运营人工派发与责任人转发

运营可以读取并人工选择 FOBrain 人员，基于 Epic 1 的冻结漏洞事实准备、确认、执行和核验派发／转发，逐条看到完成、待核验和待排查结果，并可跨会话重新找到同源操作记录。

### Story 2.1：读取 FOBrain 全部人员安全列表

As a 产品中心运营，
I want 查看当前权限内 FOBrain 返回的全部人员安全列表，
So that 后续可以基于真实人员事实人工指定接收人。

**WorkItemType:** user-value

**Requirements:** FR-05、FR-09；NFR-01、NFR-02；AR-14、AR-26、AR-31、AR-32、AR-33、AR-44；UX-DR-15、UX-DR-20、UX-DR-21、UX-DR-22、UX-DR-23；G-READ-02；M-4。

**Architecture decisions:** AD-03、AD-05、AD-09、AD-13、AD-14、AD-15、AD-16、AD-27。

**Prerequisites:** Story 1.12。

**Inputs / Outputs:** 输入为当前 actor、FOBrain 人员接口全部分页和安全字段 allowlist；输出为完整人员 QueryResultSnapshot、稳定身份、安全消歧字段及可选状态候选事实。

**Exit gate:** 243 人级分页、空/无权限/失败、同名/停用/未知/缺标识和敏感字段测试通过，G-READ-02 形成产品证据且 mutation 为 0；否则不得进入 Story 2.2。

**Scope:** 全员列表只读 capability、243 人级 fixture、稳定人员身份、安全 allowlist、完整分页、加载／空／失败／不可选投影。

**Non-goals:** 不实现人员选择决定、不 AI 推荐、不创建 ActionDraft、不连接 mutation。

**Affected directories:** docs/fobrain-tool-matrix.md、docs/schemas/、docs/fixtures/、internal/einoapp/{capabilities,facts,product,providers/fobrain}/、web/eino-workbench/src/、scripts/acceptance/。

**Acceptance commands:** node scripts/eino_workbench_schema_validate.mjs --contract；go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_2_1_people_list_smoke.sh。

**Acceptance Criteria:**

**Given** 当前 workspace 凭据可访问 FOBrain 人员接口  
**When** 用户请求查看人员列表  
**Then** 系统读取当前 actor 可见的完整安全人员快照并保持 provider 来源顺序  
**And** 不按部门、业务、角色或 AI 判断预过滤。

**Given** 人员存在稳定标识、停用、未知状态、标识缺失或同名情况  
**When** Safety Gate 建立人员事实  
**Then** 只投影获批准的显示名称、状态和安全消歧字段，并明确标记可选性候选事实  
**And** 内部人员 ID、账号 raw 字段、token 和敏感组织字段不得进入产品出口。

**Given** 全员列表为空、分页不完整、无权限或外部失败  
**When** 查询结算  
**Then** 四类状态分别展示且不发布伪 complete_set  
**And** ActionDraft 与真实 FOBrain mutation 调用数均为 0。

### Story 2.2：人工搜索、消歧并选择唯一接收人

As a 产品中心运营，
I want 从冻结的 FOBrain 人员列表中搜索、消歧并明确选择一人，
So that 派发和转发不会因重名、失效人员或 AI 猜测而选错对象。

**WorkItemType:** decision / user-value

**Requirements:** FR-05、FR-09；NFR-02、NFR-03；AR-33、AR-41；UX-DR-14、UX-DR-15、UX-DR-20、UX-DR-23、UX-DR-24、UX-DR-25；G-READ-02、OQ-02；M-4。

**Architecture decisions:** AD-05、AD-09、AD-14、AD-15、AD-16、AD-17、AD-27。

**Prerequisites:** Story 2.1。

**Inputs / Outputs:** 输入为冻结人员快照、搜索/筛选/键盘意图和目标部署状态证据；输出为唯一 selected recipient identity、安全摘要、版本化选择规则及 OQ-02 裁决。

**Exit gate:** 搜索、来源排序、同名消歧、不可选、键盘和二次校验测试通过且 OQ-02=CLOSED；否则 Story 2.3 保持阻塞。

**Scope:** 搜索字段、状态筛选、来源排序、同名消歧 allowlist、不可选规则、combobox/listbox 键盘单选、人员快照绑定与 OQ-02 证据。

**Non-goals:** 不自动预选、不接受自由文本姓名或伪造 ID、不推荐排序、不创建写草案。

**Affected directories:** docs/schemas/、docs/fixtures/、web/eino-workbench/src/、internal/einoapp/{facts,product}/、docs/acceptance-records/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench run typecheck；npm --prefix web/eino-workbench test；npm --prefix web/eino-workbench run browser-test。

**Acceptance Criteria:**

**Given** 完整人员快照已加载  
**When** 用户搜索、筛选或清空条件  
**Then** 搜索只匹配规范化显示名称和批准的安全消歧字段，筛选仅含全部／可选／不可选  
**And** 默认及清空后均恢复冻结快照来源顺序，不提供推荐排序。

**Given** 人员同名、停用、状态未知或缺少稳定标识  
**When** 用户查看候选  
**Then** 只有证据充分且状态属于批准集合的人员可选；无法安全消歧的同名项保持不可选并说明原因  
**And** 不通过列表位置、模型推断或隐藏内部 ID 完成消歧。

**Given** 用户仅使用键盘操作人员选择卡  
**When** 聚焦、移动、选择或关闭  
**Then** combobox/listbox 语义、单一 Tab 停靠点、方向键、Enter、Escape 和焦点恢复均符合 UX 契约  
**And** 搜索结果唯一也不得自动选择或确认。

**Given** 用户明确选择一名可选人员  
**When** 保存选择并重新校验人员事实  
**Then** 选择绑定人员快照、stable identity 与安全摘要，人员变化会使选择失效  
**And** 243 人、同名、停用、未知、缺标识和键盘证据全部通过后才能关闭 OQ-02；否则后续草案保持阻塞。

### Story 2.3：基于冻结事实准备不可变派发草案

As a 产品中心运营，
I want 将已核对的漏洞和人工选择的同一接收人冻结为派发草案，
So that 确认前能够完整核对将要修改的对象和目标人员。

**WorkItemType:** user-value / mock

**Requirements:** FR-02、FR-05、FR-06、FR-16、FR-18；NFR-02、NFR-03；AR-08、AR-09、AR-10、AR-13、AR-14、AR-39、AR-40、AR-41；UX-DR-04、UX-DR-14、UX-DR-15、UX-DR-16、UX-DR-21、UX-DR-22；G-ARCH-V2、OQ-02；M-4。

**Architecture decisions:** AD-03、AD-05、AD-07、AD-08、AD-09、AD-13、AD-14、AD-15、AD-16。

**Prerequisites:** Story 2.2 且 OQ-02 closed。

**Inputs / Outputs:** 输入为未过期 complete_set 漏洞快照、人工目标人员和派发意图；输出为 policy 校验结果或不可变 ActionDraft、完整 items、JCS digest、有效期和审批安全投影。

**Exit gate:** 分组、冲突、超量、权限/freshness 变化、前 5 条＋查看全部和零真实 mutation 测试全部通过；否则不得进入 Story 2.4。

**Scope:** 派发 Catalog v2 metadata、PrepareAction、policy 双检、不可变 ActionDraft/digest、同接收人分组、max_items、前 5 条＋查看全部审批投影；只连接 mock provider/verifier。

**Non-goals:** 不确认、不执行、不连接真实 FOBrain mutation，不从 IA-02 行选择或聊天摘要重建范围。

**Affected directories:** docs/schemas/、docs/fixtures/、internal/einoapp/{capabilities,execution,facts,product}/、web/eino-workbench/src/。

**Acceptance commands:** node scripts/eino_workbench_schema_validate.mjs --contract；go test ./internal/einoapp/...；npm --prefix web/eino-workbench test。

**Acceptance Criteria:**

**Given** 用户锁定一个未过期 complete_set 漏洞快照并人工选择接收人  
**When** 请求准备派发  
**Then** PrepareAction 重新校验 actor、对象权限、freshness、人员可选状态、group_key 与 max_items  
**And** AI 不推荐接收人、不自动选择对象，不把数量摘要当对象集合。

**Given** 多条漏洞最终指定给相同或不同接收人  
**When** 系统组织草案  
**Then** 同一接收人可进入一个草案，不同接收人必须形成独立草案并分别确认  
**And** 重叠对象、不同目标冲突或超量必须澄清／固定拆分／拒绝，不得 last-write-wins 或后台静默分批。

**Given** 所有校验通过  
**When** 发布不可变派发草案  
**Then** 冻结来源快照、actor、capability/policy/verifier/route/credential binding 版本、全部 items、目标人员、有效期和 JCS SHA-256 digest  
**And** 审批投影显示动作、来源查询、总数、前 5 条 POC/IP/在线/业务/负责人事实，并可查看完整冻结集合。

**Given** 必要事实缺失、权限未知、快照不完整或人员／关键漏洞事实变化  
**When** PrepareAction 或确认前复核发现变化  
**Then** 不创建可确认草案或使旧草案失效，要求重新查询与准备  
**And** 真实 FOBrain mutation 调用数为 0。

### Story 2.4：确认或取消派发并证明未批准路径零写入

As a 产品中心运营，
I want 在审批卡中二次确认或取消不可变派发草案，
So that 只有我明确批准的同一版本草案可以进入执行。

**WorkItemType:** user-value / mock

**Requirements:** FR-07、FR-16、FR-18；NFR-01、NFR-03；AR-11、AR-12、AR-17、AR-41；UX-DR-16、UX-DR-17、UX-DR-19、UX-DR-20、UX-DR-21、UX-DR-22、UX-DR-23、UX-DR-24；G-ARCH-V2、G-SAFE-01；M-4。

**Architecture decisions:** AD-04、AD-07、AD-08、AD-09、AD-12、AD-21、AD-22、AD-24。

**Prerequisites:** Story 2.3。

**Inputs / Outputs:** 输入为不可变 draft、一次性 approval ref 与明确 confirm/cancel 决定；输出为原子 Pending 终态、批准后的 attempt reservations 或零写入拒绝/取消结果。

**Exit gate:** 首次确认、取消、拒绝、过期、digest/policy 变化、重复提交及双 continuation 测试通过，所有未批准路径 mutation=0；否则不得进入 Story 2.5。

**Scope:** PendingInteraction、一次性 approval ref、ConfirmAction、取消/拒绝/过期、Eino continuation 与确定性 Action API continuation、原子消费及零写入证据。

**Non-goals:** 不执行真实 mutation，不以聊天中的“确认”文字、详情选择或页面动作代替审批卡按钮。

**Affected directories:** internal/einoapp/{execution,facts,product,httpapi,store/sqlite}/、web/eino-workbench/src/、docs/schemas/、docs/fixtures/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；node scripts/eino_workbench_stream_test.mjs。

**Acceptance Criteria:**

**Given** 一个有效且未消费的 approval ref  
**When** 操作人在审批卡确认  
**Then** 单一事务消费 ref、固定 Pending 终态、复核 actor/policy/digest/version、批准 draft、预留 attempts 并追加 FactEvent  
**And** 事务提交后才允许 continuation 恢复或后台执行。

**Given** 用户取消、拒绝、不确认、草案过期或关键事实变化  
**When** 控制面结算 Pending  
**Then** Pending 进入唯一终态且 mock/真实 mutation transport 调用数均为 0  
**And** 继续操作必须创建新草案并重新确认。

**Given** 重复点击、重复请求、刷新或聊天中再次输入确认  
**When** approval ref 已消费  
**Then** 只投影首次权威终态，不创建第二次 attempt 或第二个 continuation  
**And** 前端确认与取消在首次提交后立即锁定。

**Given** Workbench 与 Action API 分别准备同语义操作  
**When** 执行确认协议  
**Then** 两个入口共用 Product Facts 与确认不变量，Workbench 使用真实 Eino continuation，Action API 使用项目 continuation record  
**And** 不伪造 Eino checkpoint。

### Story 2.5：单条 Mock 派发只执行一次并回读完成

As a 产品中心运营，
I want 对一条漏洞执行 mock 派发并依据负责人回读看到结果，
So that 单项写入、幂等和成功判定可以在零真实写入下被验证。

**WorkItemType:** user-value / mock

**Requirements:** FR-07、FR-16、FR-17、FR-18；NFR-03、NFR-04；AR-15、AR-16、AR-42、AR-45、AR-46；UX-DR-17、UX-DR-18、UX-DR-19、UX-DR-20、UX-DR-21、UX-DR-22；G-ARCH-V2、G-SAFE-01；M-4。

**Architecture decisions:** AD-07、AD-08、AD-09、AD-10、AD-11、AD-21、AD-22、AD-23。

**Prerequisites:** Story 2.4。

**Inputs / Outputs:** 输入为一个已确认 ActionItem、mock provider 响应和负责人 readback；输出为唯一 Attempt/claim、VerificationEvidence 与单条完成/待核验/待排查 ActionResult。

**Exit gate:** 单条成功、拒绝、未知、不一致、重复和恢复测试通过，mock mutation 最大 1、真实 mutation=0；否则不得进入 Story 2.6。

**Scope:** 单条 mock adapter、SemanticMutationClaim、Attempt reservation/sent CAS、幂等、版本化 verifier、负责人 readback 与逐条结果卡。

**Non-goals:** 不调用真实 FOBrain mutation，不以 mock 2xx 或 provider accepted 单独判成功。

**Affected directories:** internal/einoapp/{execution,facts,product,capabilities}/、internal/einoapp/providers/mock/、web/eino-workbench/src/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_2_5_mock_dispatch_smoke.sh。

**Acceptance Criteria:**

**Given** 一个已确认的单条派发草案  
**When** executor 预留并发送 mock attempt  
**Then** 跨草案稳定的 semantic mutation key 只允许一个有效 claim，reserved_not_sent→sent CAS 一次性固定 sent_at、deadline、policy、lease epoch 和核验预算  
**And** 重复确认、技术重试或恢复不会产生第二次发送。

**Given** mock provider 接受请求并由独立 mock readback 返回目标负责人  
**When** verifier 比较 expected 与 observed 安全事实  
**Then** 仅 held control＋accepted＋conclusive readback match 显示完成  
**And** 结果显示 POC、IP、在线状态、业务系统、目标人员和实际负责人。

**Given** provider 拒绝、响应未知、负责人不一致或 observed 缺失  
**When** verifier 结算  
**Then** 按版本化 evidence 唯一映射为待核验或待排查，不得伪报成功  
**And** 后续核验只能读取，不得重新发送 mutation。

**Given** 本 Story 的所有正负测试运行  
**When** 检查 provider spy  
**Then** 真实 FOBrain mutation 调用数必须为 0，单条 mock mutation 最大为 1  
**And** 产品不发送任何通知。

### Story 2.6：批量 Mock 派发逐条结算并可恢复

As a 产品中心运营，
I want 同一接收人的多条 mock 派发逐条显示完成、待核验和待排查，
So that 部分失败、崩溃和恢复不会隐藏结果或重复写入。

**WorkItemType:** user-value / mock

**Requirements:** FR-06、FR-07、FR-16、FR-17、FR-18；NFR-01、NFR-03、NFR-04；AR-15、AR-18、AR-23、AR-30、AR-40、AR-42、AR-45、AR-46；UX-DR-18、UX-DR-19、UX-DR-20、UX-DR-22、UX-DR-24、UX-DR-27；G-ARCH-V2、G-SAFE-01；M-4。

**Architecture decisions:** AD-04、AD-10、AD-11、AD-12、AD-18、AD-20、AD-21、AD-22、AD-23、AD-24。

**Prerequisites:** Story 2.5。

**Inputs / Outputs:** 输入为同一接收人的多个已确认 items、失败注入和持久 lease/attempt 事实；输出为逐条三态结果、批次计数、reconcile evidence 与可恢复执行状态。

**Exit gate:** partial、crash、lease/fencing、cancel-vs-send、预算耗尽和 restart 矩阵全部通过，sent item 无重发且真实 mutation=0；否则不得进入 Story 2.7。

**Scope:** 多条逐项 attempt、partial 结算、lease/fencing、crash injection、reconcile、预算耗尽、restart recovery、结果卡和 IA-02 逐条结果。

**Non-goals:** 不整批染成一个成功/失败状态，不后台静默拆批，不自动重试 sent attempt。

**Affected directories:** internal/einoapp/{execution,facts,product,store/sqlite}/、internal/einoapp/providers/mock/、web/eino-workbench/src/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_2_6_mock_partial_restart_smoke.sh。

**Acceptance Criteria:**

**Given** 同一接收人的多条已确认 items 返回成功、明确拒绝、超时和回读不一致组合  
**When** executor 逐条执行与核验  
**Then** 每项独立进入完成、待核验或待排查，成功项不会被失败项覆盖  
**And** 结果卡先显示三类计数，再保留全部逐条结果并可进入 IA-02。

**Given** executor 在确认后、发送前、发送后或核验中崩溃  
**When** 新 owner 以更高 lease epoch 接管  
**Then** 旧 fencing token 的迟到提交被拒绝，sent item 只继续只读核验  
**And** sent_at、deadline、核验预算和 claim lineage 不因重启重算。

**Given** cancel/stop 与首次发送发生竞态  
**When** 两个事务竞争 phase CAS  
**Then** 只有 cancelled_before_send＋零 mutation 或 sent＋继续核验之一可以提交  
**And** 不留下 reserved orphan、active claim 泄漏或伪造 sent 字段。

**Given** 核验期限或调用预算耗尽仍不一致／未知  
**When** reconcile 结算  
**Then** 项目进入待排查且保留 evidence；期限内则保持待核验  
**And** 重开页面只能读取同一事实链，不生成重试入口。

### Story 2.7：主动转发复用同一 Action 控制面

As a 产品中心运营，
I want 将权限范围内的漏洞主动转发给人工选择的人员，
So that 我可以独立纠正负责人，而不依赖系统判断是否错派。

**WorkItemType:** user-value / mock

**Requirements:** FR-08、FR-09、FR-10、FR-16、FR-17、FR-18；NFR-02、NFR-03、NFR-04；AR-10、AR-13、AR-15、AR-16、AR-39、AR-40、AR-41、AR-42、AR-45、AR-46；UX-DR-14、UX-DR-15、UX-DR-16、UX-DR-17、UX-DR-18、UX-DR-19、UX-DR-20、UX-DR-21、UX-DR-22；G-ARCH-V2、OQ-02、G-SAFE-01；M-4。

**Architecture decisions:** AD-07、AD-08、AD-09、AD-10、AD-11、AD-13、AD-15、AD-16、AD-21、AD-22、AD-23。

**Prerequisites:** Story 2.6。

**Inputs / Outputs:** 输入为当前 actor 可见漏洞快照、主动转发意图和人工目标人员；输出为转发专属 draft/policy/verifier 结果及同源逐条 ActionResult。

**Exit gate:** 无前置派发、自己的漏洞、同/异接收人、事实变化、取消、部分失败和恢复测试通过，并证明未复制状态机、真实 mutation=0；否则不得进入 Story 2.8。

**Scope:** transfer capability/policy/verifier、主动意图、同接收人分组、当前/目标负责人事实、mock 执行与逐条回读；复用 Stories 2.3～2.6 控制面。

**Non-goals:** 不识别“错派”，不要求先发生派发，不复制 ActionDraft/Pending/Attempt/lease/verifier/reconcile 状态机，不连接真实 mutation。

**Affected directories:** internal/einoapp/{capabilities,execution,facts,product}/、internal/einoapp/providers/mock/、web/eino-workbench/src/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_2_7_mock_transfer_partial_smoke.sh。

**Acceptance Criteria:**

**Given** 当前 actor 能查询到一条或多条漏洞  
**When** 用户主动表达转发并人工选择目标人员  
**Then** 不要求系统先判断错派或存在前置派发，自己的漏洞也可按登记 policy 转发  
**And** AI 不推荐是否转发或接收人。

**Given** 同一或不同接收人的多条转发对象  
**When** PrepareAction 建立草案  
**Then** 同一接收人可合并、不同接收人必须分草案，冻结原负责人和目标负责人  
**And** 权限、负责人或人员事实变化使旧草案失效。

**Given** 用户确认并执行 mock 转发  
**When** verifier 回读负责人  
**Then** 仅目标负责人一致的条目完成，其余按统一 evidence 映射待核验／待排查  
**And** 逐条结果展示 POC、IP、在线状态、业务系统、原负责人、目标人员和实际负责人。

**Given** 未确认、取消、过期、重复确认、崩溃或部分失败  
**When** 共用控制面处理  
**Then** 零写入、幂等、lease、reconcile 和恢复行为与派发一致  
**And** 真实 FOBrain mutation 调用数为 0，不触发派发或通知作为补偿。

### Story 2.8：从固定“操作记录”入口恢复同源逐条结果

As a 产品中心运营，
I want 跨会话从“操作记录”找到近期操作和仍待处理的逐条结果，
So that 我不依赖聊天摘要恢复完成、待核验和待排查事实。

**WorkItemType:** user-value

**Requirements:** FR-17、FR-18；NFR-01、NFR-02；AR-23、AR-35、AR-47；UX-DR-03、UX-DR-05、UX-DR-13、UX-DR-18、UX-DR-20、UX-DR-22、UX-DR-27、UX-DR-30；OQ-05（CLOSED by RD-03／AD-26）、G-ARCH-V2；M-4。

**Architecture decisions:** AD-03、AD-04、AD-16、AD-20、AD-21、AD-24、AD-25、AD-26。

**Prerequisites:** Story 2.7。

**Inputs / Outputs:** 输入为当前 actor/workspace、持久化 Product Facts/ActionResult/open lifecycle 和筛选/分页参数；输出为 immutable OperationListSnapshot、tamper-evident cursor 和 IA-01/IA-02 只读历史投影。

**Exit gate:** 六个月＋open union、稳定排序/分页、重登录、权限 fail-closed、cursor/TTL 和只读边界测试通过，RD-03/OQ-05/AD-26 被实现；完成后 Epic 3 可开始。

**Scope:** 激活 Epic 1 固定入口的数据、OperationListSnapshot、六个 UTC+8 自然月＋open lifecycle 并集、actor fail-closed、稳定游标分页、逐条只读结果和右栏执行记录。

**Non-goals:** 不从聊天摘要/provider raw payload 重建历史，不提供草案、确认、重试或重新归属入口，不做完整事件溯源。

**Affected directories:** internal/einoapp/{facts,product,httpapi,store/sqlite}/、web/eino-workbench/src/、docs/schemas/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；npm --prefix web/eino-workbench run browser-test；bash scripts/acceptance/story_2_8_operation_history_smoke.sh。

**Acceptance Criteria:**

**Given** 当前 workspace 与凭据主体拥有已授权操作或 open lifecycle 对象  
**When** 首次打开固定“操作记录”入口  
**Then** 服务端以当前稳定 actor fail closed 授权，在一个事务中物化“前六个 UTC+8 自然月＋仍未闭合对象”的完整有序 OperationListSnapshot  
**And** operation_at 在 Pending published_at 或无 Pending 的 decision time 首次固定，后续审批只写 decision_at  
**And** 本入口直接实施已关闭的 RD-03／OQ-05 与 AD-26，不再把 OQ-05 作为等待关闭的门禁。

**Given** 用户按操作时间倒序分页  
**When** 请求下一页或跨会话重开  
**Then** 排序固定为 operation_at desc、run_id desc，opaque cursor 绑定 actor/workspace、filters、as_of、cutoff、snapshot ref 与 last keys  
**And** 翻页期间其他 Run 变化不会造成重复或遗漏，且不得复用 Run-local SSE sequence。

**Given** 用户进入一条操作记录  
**When** 查看列表和逐条结果  
**Then** 显示动作、对象摘要、三类计数及每项 POC、IP、在线、业务、目标/实际负责人和最终状态  
**And** 页面只读，不生成草案、确认、重试或通知。

**Given** 权限未知／拒绝、cursor 篡改、快照缺失／过期或 history policy 损坏  
**When** 服务端处理请求  
**Then** 整批 fail closed 并显示确定性错误，不静默过滤、重新归属或从聊天摘要回退  
**And** Epic 3 可以复用该投影而无需等待 Stories 2.9／2.10。

### Story 2.9：验证目标部署的派发权限、写入与回读

As a 产品中心研发负责人，
I want 在获准环境中验证真实派发的权限、单次写入、负责人回读和恢复，
So that 真实派发只有在行为可证明时才可转换为 Target Story。

**WorkItemType:** gate-evidence

**Requirements:** FR-05、FR-06、FR-07、FR-16、FR-17、FR-18；NFR-02、NFR-03、NFR-04；AR-26、AR-30、AR-31、AR-42、AR-45、AR-46；UX-DR-18、UX-DR-20、UX-DR-21；G-WRITE-01、G-WRITE-02、G-SAFE-01；M-5。

**Architecture decisions:** AD-09、AD-10、AD-11、AD-14、AD-15、AD-23、AD-27。

**Prerequisites:** Story 2.8。

**Inputs / Outputs:** 输入为带日期授权、获准验证窗口、最小可恢复漏洞/人员样本、写入上限和恢复方案；输出为隔离取证脚本结果、脱敏 mutation/readback/recovery 证据及派发 Gate PASS/BLOCKED 裁决。

**Exit gate:** 无论 PASS 或 BLOCKED 都必须有可复核记录；只有 G-WRITE-01/G-WRITE-02/G-SAFE 派发子项全部 PASS 才允许由 Create Story 生成新的 M-6 正式 Target Story，且本 Story 永不启用生产 capability。

**Scope:** 目标部署合同取证、带日期授权、最小可恢复样本、零写入负向路径、单次 mutation、独立负责人 readback、部分失败、响应未知/reconcile、脱敏证据和门禁更新。

**Non-goals:** 不接入生产 Workbench，不启用真实派发 capability；Gate PASS 只允许后续创建新的 M-6 正式 Target Story，不等于产品集成完成。

**Affected directories:** scripts/acceptance/、docs/fixtures/、docs/acceptance-records/、_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md。

**Acceptance commands:** bash scripts/acceptance/story_2_9_dispatch_gate.sh（缺少环境授权、恢复方案或安全样本时必须 exit 2 且 mutation=0）；go test ./internal/einoapp/providers/...。

**Acceptance Criteria:**

**Given** 派发验证可能改变目标数据  
**When** 准备首次真实 mutation  
**Then** 必须具备带日期环境授权、验证窗口、最小样本、目标人员、最大写入次数、恢复步骤和失败联系人  
**And** 任一条件缺失时 blocked/exit 2，真实 mutation 调用数为 0。

**Given** 未确认、取消、拒绝、过期、digest 不匹配或权限变化的草案  
**When** 运行安全负向验证  
**Then** provider mutation transport 实际调用数为 0  
**And** audit evidence 保存脱敏原因且产品通知数为 0。

**Given** 一个已授权可恢复样本在确认后执行派发  
**When** provider 接受并通过独立查询回读负责人  
**Then** 只有实际负责人等于目标人员时样本通过，每个成功样本只发生一次 mutation  
**And** 验证后恢复原负责人并再次回读证明。

**Given** 部分失败、超时、响应丢失或回读不一致  
**When** 目标部署验证结算  
**Then** 成功、明确失败和待核验／待排查可逐条区分，不确定请求不重发而先回读 reconcile  
**And** 无法安全制造的证据子项保持 BLOCKED，不用 mock 替代。

**Given** 所有样本与恢复步骤完成  
**When** 更新 G-WRITE/G-SAFE  
**Then** 仅按脱敏客观证据裁决派发子项并保持生产 capability 关闭  
**And** 任一恢复失败或证据缺失时门禁保持 BLOCKED。

### Story 2.10：验证目标部署的转发权限、写入与回读

As a 产品中心研发负责人，
I want 在获准环境中独立验证主动转发的对象权限、写入和负责人回读，
So that 系统不会把派发接口存在误认为转发场景已经可用。

**WorkItemType:** gate-evidence

**Requirements:** FR-08、FR-09、FR-10、FR-16、FR-17、FR-18；NFR-02、NFR-03、NFR-04；AR-26、AR-30、AR-31、AR-42、AR-45、AR-46；UX-DR-18、UX-DR-20、UX-DR-21；G-WRITE-01、G-WRITE-02、G-SAFE-01；M-5。

**Architecture decisions:** AD-09、AD-10、AD-11、AD-14、AD-15、AD-23、AD-27。

**Prerequisites:** Story 2.8；不依赖 Story 2.9 PASS。

**Inputs / Outputs:** 输入为转发场景独立授权、当前/目标负责人已知的可恢复样本、写入上限和恢复方案；输出为隔离取证脚本结果、脱敏转发/readback/recovery 证据及转发 Gate PASS/BLOCKED 裁决。

**Exit gate:** 无论 PASS 或 BLOCKED 都必须有独立可复核记录；只有 G-WRITE-01/G-WRITE-02/G-SAFE 转发子项全部 PASS 才允许由 Create Story 生成新的 M-6 正式 Target Story，且不复用派发结论或启用生产 capability。

**Scope:** 转发场景独立合同/权限、无前置派发、可恢复原负责人样本、零写入负向路径、单次 mutation、readback、部分失败/未知恢复、脱敏证据和门禁更新。

**Non-goals:** 不接入生产 Workbench，不启用真实转发 capability，不因底层 endpoint 相同复用派发结论。

**Affected directories:** scripts/acceptance/、docs/fixtures/、docs/acceptance-records/、_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md。

**Acceptance commands:** bash scripts/acceptance/story_2_10_transfer_gate.sh（缺少环境授权、恢复方案或安全样本时必须 exit 2 且 mutation=0）；go test ./internal/einoapp/providers/...。

**Acceptance Criteria:**

**Given** 准备验证主动转发  
**When** 核对目标部署合同  
**Then** 独立证明当前 actor 对对象的权限、目标人员约束、负责人变更入口、错误语义和 readback 入口  
**And** 即使 endpoint 相同也必须使用转发专属 policy、verifier 与证据。

**Given** 一个当前负责人已知且可恢复的授权漏洞  
**When** 用户不经前置派发确认主动转发  
**Then** 只执行一次 mutation，独立回读目标负责人一致后才通过  
**And** 当前用户自己的漏洞也按实际 policy 验证，完成后恢复原负责人并回读证明。

**Given** 未确认类路径、部分失败、响应未知或结果不一致  
**When** 执行安全与恢复验证  
**Then** 未确认类真实 mutation 为 0，不确定请求不重发，逐条进入完成、待核验或待排查  
**And** 产品不触发派发、补偿动作或通知。

**Given** 验证记录生成  
**When** 更新转发门禁子项  
**Then** 仅客观、脱敏、可恢复证据可标记 PASS，无法验证的子项保持 BLOCKED  
**And** 生产转发入口继续关闭，Epic 3 不等待本 Story 或 Story 2.9 的裁决。

## Epic 3：管理员延时与运营误报处置

管理员可以对人工沟通确定的一条漏洞执行修复延时，运营可以对自己负责且已人工判断的漏洞标记误报；两类动作复用 Epic 2 的可信确认、执行、核验和操作记录，只有目标事实回读一致才显示成功。

### Story 3.1：验证并固定直接修复延时规则

As a 产品中心研发负责人，
I want 用目标部署证据固定直接修复延时的状态、期限和角色规则，
So that 产品不会猜测哪些漏洞可延时或接受什么新期限。

**WorkItemType:** decision / gate-evidence

**Requirements:** FR-11、FR-12、FR-13、FR-16；NFR-02、NFR-03、NFR-04；AR-26、AR-31、AR-34、AR-41；UX-DR-14、UX-DR-16、UX-DR-20、UX-DR-21；OQ-04、G-WRITE-03；M-5。

**Architecture decisions:** AD-07、AD-09、AD-13、AD-14、AD-15、AD-23、AD-27。

**Prerequisites:** Story 2.8；不依赖 Stories 2.9／2.10 或任何 M-6 Target Story。

**Inputs / Outputs:** 输入为目标部署直接延时合同、状态/期限/角色证据和产品决定；输出为版本化规则表、schema/policy/fixtures、OQ-04 与 G-WRITE-03 的 PASS/BLOCKED 裁决。

**Exit gate:** 必须形成零 mutation 的可复核规则证据；只有 OQ-04=CLOSED 且状态、期限、时区、管理员和 max_items=1 全部确定时才允许 Story 3.2。

**Scope:** 直接延时而非申请审批流程的合同证据、允许原状态、期限范围/粒度/时区/关系、管理员角色、原因可空、max_items=1、schema/policy/fixture 与 OQ-04 裁决。

**Non-goals:** 不执行真实 mutation，不启用生产延时 capability，不由 AI 决定是否延时、新期限或原因。

**Affected directories:** docs/schemas/、docs/fixtures/、internal/einoapp/{capabilities,facts,product,providers/fobrain}/、scripts/acceptance/、docs/acceptance-records/、_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md。

**Acceptance commands:** node scripts/eino_workbench_schema_validate.mjs --contract；go test ./internal/einoapp/...；bash scripts/acceptance/story_3_1_delay_rules_evidence.sh（证据不足时 exit 2）。

**Acceptance Criteria:**

**Given** FOBrain 可能同时存在直接延时和申请审批延时  
**When** 核对目标部署合同、角色与产品决定  
**Then** 必须证明本轮采用直接修复延时，不引入申请、审批、拒绝或催办流程  
**And** 历史源码只能作为线索，不能单独关闭 OQ-04。

**Given** 目标部署返回状态和角色信息  
**When** 固定对象与 actor policy  
**Then** 明确允许原状态集合及稳定映射，只有当前 FOBrain 身份解析出的管理员可准备和确认  
**And** 未知、缺失、冲突或不在允许集合时 fail closed。

**Given** 用户输入新修复期限和可选原因  
**When** 固定输入契约  
**Then** 明确最早/最晚值、时间粒度、规范时区、与当前时间及现期限关系，并覆盖边界 fixture  
**And** 新期限必须由人显式输入，原因允许为空且不作为成功条件，max_items 固定为 1。

**Given** 状态、期限、角色、schema、policy 和 fixture 证据完成  
**When** 产品负责人裁决 OQ-04  
**Then** 只有全部一致才关闭 OQ-04并同步门禁；任一未知则保持 BLOCKED  
**And** 本 Story 的真实 FOBrain mutation 调用数必须为 0。

### Story 3.2：复用控制面完成单条 Mock 修复延时

As a FOBrain 管理员，
I want 对人工沟通确定的一条漏洞输入新期限并二次确认 mock 延时，
So that 权限、人工参数、零写入路径和双事实核验可先被机器验证。

**WorkItemType:** user-value / mock

**Requirements:** FR-11、FR-12、FR-13、FR-16、FR-17、FR-18；NFR-02、NFR-03、NFR-04；AR-10、AR-12、AR-15、AR-16、AR-34、AR-41、AR-42、AR-45、AR-46；UX-DR-14、UX-DR-16、UX-DR-17、UX-DR-18、UX-DR-19、UX-DR-20、UX-DR-21、UX-DR-22、UX-DR-27、UX-DR-30；G-ARCH-V2、OQ-04、G-SAFE-01；M-5。

**Architecture decisions:** AD-07、AD-08、AD-09、AD-10、AD-11、AD-12、AD-15、AD-16、AD-21、AD-22、AD-23、AD-26。

**Prerequisites:** Story 3.1 且 OQ-04 closed；复用 Stories 2.3～2.8。

**Inputs / Outputs:** 输入为管理员 actor、单条授权漏洞、当前状态/期限、人工新期限和可选原因；输出为不可变延时 draft、Pending/Attempt/VerificationEvidence、双事实 ActionResult 和操作记录。

**Exit gate:** 原因空值、边界期限、角色/状态变化、零写入、双事实不一致、重复确认和 restart 测试通过，真实 mutation=0；否则 Story 3.4 保持阻塞。

**Scope:** delay capability/policy/verifier、单条人工期限/原因、不可变草案、确认/取消、mock mutation、状态＋期限双 readback、恢复与操作记录。

**Non-goals:** 不创建延时专用状态机，不连接真实 FOBrain mutation，不批量延时，不提供 AI 期限建议。

**Affected directories:** internal/einoapp/{capabilities,execution,facts,product}/、internal/einoapp/providers/mock/、web/eino-workbench/src/、docs/schemas/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** node scripts/eino_workbench_schema_validate.mjs --contract；go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_3_2_mock_delay_restart_smoke.sh。

**Acceptance Criteria:**

**Given** 管理员与相关人员已线下沟通  
**When** 管理员选择一条漏洞并显式输入新期限与可选原因  
**Then** 是否延时和新期限完全来自用户输入，AI 不推荐、不修正、不自动确认  
**And** 多条对象、非管理员、未知角色或无效期限均不生成草案。

**Given** PrepareAction 校验通过  
**When** 发布延时审批卡  
**Then** 冻结并完整展示 POC、IP、在线状态、业务系统、负责人、当前状态、当前期限、新期限和原因  
**And** 草案变化必须生成新 digest 并重新确认。

**Given** 未确认、取消、拒绝、过期、权限/状态/期限规则变化或重复提交  
**When** 共用控制面处理  
**Then** mock/真实 mutation 调用数为 0 或仅投影首次权威终态  
**And** 不创建延时专用补偿、重试或通知路径。

**Given** 管理员确认有效草案并执行 mock 延时  
**When** verifier 回读实际状态和修复期限  
**Then** 只有状态为“延时”且期限等于冻结目标时显示“延时成功”  
**And** 仅一个事实一致、provider accepted/2xx、事实缺失或回读不一致均不得成功。

**Given** 响应未知、失租约、回读不一致或预算耗尽  
**When** 恢复与 reconcile 运行  
**Then** 期限内显示待核验，耗尽或 control lost 显示待排查，sent attempt 不得重发  
**And** 同源操作记录可跨会话只读找到该项及 expected/observed 安全事实。

### Story 3.3：复用控制面完成人工误报 Mock 旅程

As a 产品中心运营，
I want 把自己负责且已人工判断的漏洞二次确认标记为误报，
So that 人工决定、对象权限和状态回读规则可以先在零真实写入下验证。

**WorkItemType:** user-value / mock

**Requirements:** FR-14、FR-15、FR-16、FR-17、FR-18；NFR-02、NFR-03、NFR-04；AR-10、AR-12、AR-15、AR-16、AR-40、AR-41、AR-42、AR-45、AR-46；UX-DR-14、UX-DR-16、UX-DR-17、UX-DR-18、UX-DR-19、UX-DR-20、UX-DR-21、UX-DR-22、UX-DR-27、UX-DR-30；G-ARCH-V2、G-SAFE-01；M-5。

**Architecture decisions:** AD-07、AD-08、AD-09、AD-10、AD-11、AD-12、AD-15、AD-16、AD-21、AD-22、AD-23、AD-26。

**Prerequisites:** Story 2.8；不依赖 Story 3.2、Stories 2.9／2.10 或任何 M-6 Target Story。

**Inputs / Outputs:** 输入为当前 actor、其负责对象的完整快照、人工误报结论和确认决定；输出为不可变误报 draft、逐条 Attempt/VerificationEvidence、三态 ActionResult 和操作记录。

**Exit gate:** 自己/非自己负责、超量、事实变化、取消、部分失败、readback 不一致、重复确认和 restart 测试通过，真实 mutation=0；否则 Story 3.5 保持阻塞。

**Scope:** false-positive capability/policy/verifier、当前 actor 负责对象校验、人工结论、不可变草案、确认/取消、mock mutation、状态 readback、部分结果、恢复与操作记录。

**Non-goals:** 不创建自动误报规则，不让 AI 推荐/判断误报，不连接真实 FOBrain mutation，不复制共用状态机。

**Affected directories:** internal/einoapp/{capabilities,execution,facts,product}/、internal/einoapp/providers/mock/、web/eino-workbench/src/、docs/schemas/、docs/fixtures/、scripts/acceptance/。

**Acceptance commands:** node scripts/eino_workbench_schema_validate.mjs --contract；go test ./internal/einoapp/...；npm --prefix web/eino-workbench test；bash scripts/acceptance/story_3_3_mock_false_positive_smoke.sh。

**Acceptance Criteria:**

**Given** 操作人查看权限内的漏洞并人工判断为误报  
**When** 发起误报标记  
**Then** 误报结论必须由操作人明确作出，当前负责人通过稳定 identity 与 actor 一致  
**And** AI 不推荐、不推断、不自动选择误报对象。

**Given** 用户选择一个或多个对象  
**When** PrepareAction 校验  
**Then** 对象必须来自完整、未过期的授权快照并逐条复核负责人、权限和允许原状态  
**And** 超量或任一不符合项必须明确拒绝／澄清，不得静默删除或分批。

**Given** 草案有效  
**When** 发布并确认误报审批卡  
**Then** 冻结每条 POC、IP、在线状态、业务系统、当前负责人、原状态和目标状态“误报”；多条显示前 5 条并可查看全部  
**And** 聊天文本或明细行选择不能替代二次确认。

**Given** 未确认、取消、过期、负责人/权限/状态变化或重复确认  
**When** 共用控制面处理  
**Then** 未批准路径 mock/真实 mutation 为 0，重复路径只投影首次终态  
**And** 不自动改为派发、转发或发送通知。

**Given** 确认后 mock adapter 返回逐条结果  
**When** verifier 回读状态  
**Then** 只有 provider accepted 且 conclusive readback 状态为“误报”的条目完成，其余为待核验／待排查  
**And** 成功项独立保留、sent 项不重发，跨会话操作记录展示同源逐条结果。

### Story 3.4：验证目标部署的延时权限、写入与双事实回读

As a 产品中心研发负责人，
I want 在获准环境中验证直接修复延时的管理员权限、单次写入和状态/期限回读，
So that 真实延时只有在两个目标事实都可证明时才可转换为 Target Story。

**WorkItemType:** gate-evidence

**Requirements:** FR-11、FR-12、FR-13、FR-16、FR-17、FR-18；NFR-02、NFR-03、NFR-04；AR-26、AR-30、AR-31、AR-34、AR-42、AR-45、AR-46；UX-DR-18、UX-DR-20、UX-DR-21；OQ-04、G-WRITE-01、G-WRITE-03、G-SAFE-01；M-5。

**Architecture decisions:** AD-09、AD-10、AD-11、AD-14、AD-15、AD-23、AD-27。

**Prerequisites:** Story 3.2；不依赖 Stories 2.9／2.10。

**Inputs / Outputs:** 输入为 OQ-04 规则、带日期授权、管理员/非管理员身份、单条可恢复状态/期限样本和写入上限；输出为隔离取证脚本结果、脱敏双事实/readback/recovery 证据及延时 Gate PASS/BLOCKED 裁决。

**Exit gate:** 无论 PASS 或 BLOCKED 都必须有可复核记录；只有 G-WRITE-01/G-WRITE-03/G-SAFE 延时子项全部 PASS 才允许由 Create Story 生成新的 M-6 正式 Target Story，且本 Story 永不启用生产 capability。

**Scope:** 目标部署直接延时合同、管理员/非管理员边界、带日期授权、单条可恢复样本、零写入负向路径、单次 mutation、状态＋期限 readback、未知/reconcile、恢复和脱敏门禁证据。

**Non-goals:** 不接入生产 Workbench，不启用真实延时 capability；Gate PASS 不等于产品已集成。

**Affected directories:** scripts/acceptance/、docs/fixtures/、docs/acceptance-records/、_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md。

**Acceptance commands:** bash scripts/acceptance/story_3_4_delay_gate.sh（缺少授权、恢复步骤或安全样本时 exit 2 且 mutation=0）；go test ./internal/einoapp/providers/...。

**Acceptance Criteria:**

**Given** 延时验证会改变状态和期限  
**When** 准备首次 mutation  
**Then** 必须具备带日期授权、验证窗口、原状态/原期限已知的单条样本、新期限、最大写入次数、恢复步骤和失败联系人  
**And** 任一条件缺失时 blocked/exit 2 且真实 mutation 为 0。

**Given** 管理员、非管理员、未知或冲突角色测试身份  
**When** 验证角色边界和未确认类路径  
**Then** 只有管理员的合法已确认草案可进入一次写入，其他路径 provider mutation transport 为 0  
**And** actor 不得来自前端、聊天声明或显示名称。

**Given** 一个符合 OQ-04 的样本确认执行直接延时  
**When** 独立只读回读状态和期限  
**Then** 只有状态为“延时”且期限等于冻结目标时通过  
**And** 完成后恢复原状态与期限并再次回读证明。

**Given** 响应丢失、仅一个目标事实一致或结果未知  
**When** 安全恢复验证运行  
**Then** 不重发 mutation，先回读并在预算内待核验、耗尽/control lost 后待排查  
**And** 证据不包含 raw payload、凭据或真实漏洞敏感信息。

**Given** 取证与恢复完成  
**When** 更新 G-WRITE-01/G-WRITE-03/G-SAFE 延时子项  
**Then** 仅客观双事实证据可标记 PASS，任一恢复失败或规则不明保持 BLOCKED  
**And** 生产延时入口继续关闭。

### Story 3.5：验证目标部署的误报权限、写入与状态回读

As a 产品中心研发负责人，
I want 在获准环境中验证误报的对象权限、单次写入和状态回读，
So that 产品不会允许用户修改不属于自己的漏洞或把接口接受误报为完成。

**WorkItemType:** gate-evidence

**Requirements:** FR-14、FR-15、FR-16、FR-17、FR-18；NFR-02、NFR-03、NFR-04；AR-26、AR-30、AR-31、AR-42、AR-45、AR-46；UX-DR-18、UX-DR-20、UX-DR-21；G-WRITE-01、G-WRITE-04、G-SAFE-01；M-5。

**Architecture decisions:** AD-09、AD-10、AD-11、AD-14、AD-15、AD-23、AD-27。

**Prerequisites:** Story 3.3；不依赖 Stories 2.9／2.10 或 Story 3.4 PASS。

**Inputs / Outputs:** 输入为误报合同、带日期授权、自己/非自己负责样本、可恢复原状态和写入上限；输出为隔离取证脚本结果、脱敏权限/mutation/readback/recovery 证据及误报 Gate PASS/BLOCKED 裁决。

**Exit gate:** 无论 PASS 或 BLOCKED 都必须有可复核记录；只有 G-WRITE-01/G-WRITE-04/G-SAFE 误报子项全部 PASS 才允许由 Create Story 生成新的 M-6 正式 Target Story，且本 Story 永不启用生产 capability。

**Scope:** 稳定负责人 identity、允许原/目标状态、带日期授权、自己负责/非自己负责边界、可恢复样本、零写入负向路径、单次/部分 mutation、readback、未知/reconcile、恢复和脱敏证据。

**Non-goals:** 不接入生产 Workbench，不启用真实误报 capability，不把显示名称或聊天声明当权限证明。

**Affected directories:** scripts/acceptance/、docs/fixtures/、docs/acceptance-records/、_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md。

**Acceptance commands:** bash scripts/acceptance/story_3_5_false_positive_gate.sh（缺少授权、恢复步骤或安全样本时 exit 2 且 mutation=0）；go test ./internal/einoapp/providers/...。

**Acceptance Criteria:**

**Given** 准备验证误报能力  
**When** 核对目标部署合同  
**Then** 明确“当前 actor 负责”的稳定 identity、允许原状态、目标状态“误报”、请求/错误语义和独立 readback  
**And** provider accepted 或 HTTP 2xx 不能单独判定成功。

**Given** 自己负责、非自己负责、负责人缺失、重名冲突或权限未知对象  
**When** 执行对象权限与未确认类测试  
**Then** 只有 stable owner identity 与 actor 一致且已确认的对象可写，其余路径 mutation transport 为 0  
**And** 不通过显示名称、前端传值或聊天声明扩大权限。

**Given** 操作人已人工判断一个获授权可恢复样本为误报  
**When** 确认后执行一次最小 mutation 并独立回读  
**Then** 只有实际状态为“误报”时通过，每个成功样本只写一次  
**And** 验证后恢复原状态并再次回读证明。

**Given** 获批的多对象样本产生成功、明确失败、响应未知或回读不一致  
**When** 逐条验证与恢复  
**Then** 结果独立映射为完成、待核验、待排查，未知写入不重发而先只读核验  
**And** 无法安全制造的部分失败子项保持 BLOCKED，不使用 mock 替代。

**Given** 所有证据和恢复步骤完成  
**When** 更新 G-WRITE-01/G-WRITE-04/G-SAFE 误报子项  
**Then** 仅客观、脱敏、可恢复证据可标记 PASS，任一恢复失败或权限不明保持 BLOCKED  
**And** 生产误报入口继续关闭，只有门禁后转换的 Target Story 才可接真实 mutation。
