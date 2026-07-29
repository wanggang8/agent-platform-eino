---
title: Agent Platform Eino Implementation Readiness Assessment
date: '2026-07-15'
status: SUPERSEDED
assessment: PRE_GREENFIELD_DECISION
readiness: REVALIDATION_REQUIRED
stepsCompleted:
  - step-01-document-discovery
  - step-02-prd-analysis
  - step-03-epic-coverage-validation
  - step-04-ux-alignment
  - step-05-epic-quality-review
  - step-06-final-assessment
inputDocuments:
  prd:
    - _bmad-output/planning-artifacts/product-blueprint/prd.md
  architecture:
    - _bmad-output/planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14/ARCHITECTURE-SPINE.md
    - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/SPEC.md
    - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/domain-contracts.md
    - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/state-machines.md
    - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/migration-and-gates.md
    - _bmad-output/specs/spec-product-facts-v2-runtime-foundation/acceptance-matrix.md
  epics:
    - _bmad-output/planning-artifacts/epics.md
  ux:
    - _bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/EXPERIENCE.md
    - _bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/DESIGN.md
  governance:
    - _bmad-output/planning-artifacts/product-blueprint/requirements-traceability-matrix.md
    - _bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md
    - _bmad-output/planning-artifacts/sprint-change-proposal-2026-07-15.md
---

# Implementation Readiness Assessment Report

> 本报告在 2026-07-15 Greenfield 无旧实现兼容决策之前生成，已被 `docs/adr/2026-07-15-greenfield-no-legacy-compatibility.md` 取代为历史证据。Story 1.1 工具链工作仍可继续；M-1 开始前必须基于更新后的 Architecture、SPEC、Epics 和 Gate 重新执行 Implementation Readiness，不得引用本报告的 READY 结论。

**Date:** 2026-07-15  
**Project:** Agent Platform Eino  
**Assessment:** Approved Course Correction 后复核

## 1. Document Discovery

| 类型 | 采用文档 | 大小 | 发现结果 |
| --- | --- | ---: | --- |
| PRD | `product-blueprint/prd.md` | 10,901 bytes | 唯一 canonical whole document |
| Architecture | `architecture/.../ARCHITECTURE-SPINE.md` + Runtime SPEC companions | 41,262 bytes（Spine） | 唯一 canonical Spine；SPEC 是机器契约 companion，不是重复 Architecture |
| Epics & Stories | `epics.md` | 200,524 bytes | 唯一 canonical whole document |
| UX | `ux-designs/.../EXPERIENCE.md` + `DESIGN.md` | 35,193 + 24,416 bytes | 两份分别承担体验/行为与视觉 token，不是 whole/sharded 重复 |

未发现 whole + sharded 的重复格式，也未发现缺失的必需文档。已有 `implementation-readiness-report-2026-07-15.md` 是纠偏前历史结论，本报告使用 `-post-correction` 新文件名保留历史，不覆盖原报告。

## 2. PRD Analysis

### Functional Requirements

- **FR-01：** 产品必须只使用当前用户已授权的 FOBrain 事实展示漏洞及关联上下文。
- **FR-02：** 产品在漏洞列表、二次确认和结果展示中必须呈现统一漏洞事实。单条直接展示完整对象；多条时确认卡默认展示前 5 条并提供“查看全部”，使执行人在确认前可以核对完整冻结集合；不得仅用数量替代待写对象明细。
- **FR-03：** 产品必须区分真实空结果、无权限、数据不足和外部失败；不得由 AI 虚构“没有待处理漏洞”等结果。
- **FR-04：** 运营必须能够查看全部新增漏洞，按发现时间倒序；无结果时展示“没有待派发漏洞”。
- **FR-05：** 运营必须先查看漏洞、IP 在线状态、业务系统、业务系统负责人和运维负责人，再从 FOBrain 全部人员列表选择接收人；AI 不得替代该决定。
- **FR-06：** 仅已逐条判断且最终选择同一接收人的漏洞可合并派发；不同接收人不得混合。
- **FR-07：** 派发仅在二次确认后写入；每条漏洞只有接口成功且回读修复负责人为指定人员时才完成。
- **FR-08：** 当前用户可主动对其内置权限范围内的漏洞发起转发；不要求系统识别错派或先发生其他动作。
- **FR-09：** 转发接收人必须从 FOBrain 全部人员列表选择；同一接收人可合并、不同接收人必须拆分。
- **FR-10：** 转发仅在二次确认后写入；逐条回读负责人为目标人员时才完成。
- **FR-11：** 运营与相关人员人工沟通后确定新修复期限；产品不得自动决定是否延时或期限。
- **FR-12：** 直接延时一次只处理一条漏洞，延时原因可不填；只有管理员完成二次确认后可执行。
- **FR-13：** 只有回读状态为“延时”且修复期限等于指定新期限时，结果才显示“延时成功”。
- **FR-14：** 操作人可对自己确认负责的漏洞发起误报标记；误报判断必须由人作出，不能扩展为自动规则。
- **FR-15：** 误报标记仅在二次确认后写入；只有回读状态为“误报”时才显示标记成功。
- **FR-16：** 任一写动作在未完成二次确认或取消确认时，不得向外部系统写入。
- **FR-17：** 任一写动作必须逐条区分完成、待核验和待排查。当前无法确认且仍可沿同一执行事实链继续核验时显示“待核验”；明确失败，或核验期限／调用预算结束后仍不一致、未知时显示“待排查”。成功项可独立完成；任何未完成项不得伪报成功，也不得通过重复外部写入完成核验。
- **FR-18：** 本产品不得发送消息、邮件、webhook、催办或通知中心通知；FOBrain 自身通知不构成本产品功能或成功条件。
- **FR-19：** 产品必须能够读取当前 FOBrain 用户接口提供的安全身份、部门和角色上下文，用于解释凭据主体和权限范围；部门等源字段未提供时必须显示稳定的“源字段不可用／未提供”，不得推断或伪报已解析，也不得展示 token、原始账号对象或凭据。
- **FR-20：** 产品必须能够读取并安全展示当前用户接口提供的 FOBrain 权限和数据范围，并明确区分 `resolved`、`confirmed_empty/no_permission` 与 `source_field_unavailable/unknown`；未知不得折叠为无权限或通过，不得展示原始 policy payload。
- **FR-21：** 产品必须能够按 IP 查询资产列表及其安全在线状态，作为统一漏洞事实和资产核对上下文。
- **FR-22：** 产品必须能够从资产列表进入单个资产的安全详情；安全引用不得冒充 provider 的真实资产 ID。
- **FR-23：** 产品必须能够按 IP 查询关联漏洞列表，展示安全漏洞事实并支持资产—漏洞联查。
- **FR-24：** 产品必须能够从漏洞列表进入单个漏洞的安全详情，为人工处置前判断提供事实。
- **FR-25：** 产品必须能够通过本轮已批准的 `business_list` 能力查询 FOBrain 全部业务系统安全列表，为统一漏洞卡片和业务上下文提供事实；可见范围只由 FOBrain 接口内置权限裁剪，不得替换为“我的业务系统”等已排除能力，也不得展示原始 provider payload。

**Total FRs: 25**

### Non-Functional Requirements

- **NFR-01（安全事实）：** 外部原始返回不能直接作为产品展示、审计或模型上下文事实；产品只消费安全、结构化的业务事实。
- **NFR-02（权限边界）：** 产品不得扩大 FOBrain 权限、代替 FOBrain 管理权限，或将查询不到的对象视为可操作对象。
- **NFR-03（写入完整性）：** 所有数据修改均须绑定执行人二次确认、客观回读和逐条失败可见性；未确认、取消、拒绝或过期时外部写入数必须为零。
- **NFR-04（证据边界）：** 不把目标部署 API 已验证等同于产品已集成；不得对写入结果、写域状态值域或权限作静默假设。

**Total NFRs: 4**

### Additional Requirements

- 本轮只验证“读—人工判断—确认—写—回读”，不建设完整漏洞管理、自动派发、通知或工单闭环。
- 接收人、误报、是否延时和新期限始终由人决定；AI 只组织事实、理解指代和选择获准能力。
- 统一漏洞卡片、冻结全集、前 5 条确认摘要、全部对象只读入口和逐条结果是共同验收边界。
- 前两项只读 capability、七项代表性只读能力和统一行级事实可先行；四类真实 mutation 继续受独立门禁阻塞。

### PRD Completeness Assessment

PRD 的业务范围、角色、关键旅程、25 条 FR、4 条跨功能边界、非目标和实施前阻塞项均明确且编号连续。FR-17 已消除待核验／待排查歧义，FR-19～FR-25 已回到规范需求源，FR-25 明确绑定 `business_list`。PRD 完整度足以进入覆盖验证；它本身不授权任何真实写入。

## 3. Epic Coverage Validation

覆盖只从每个正式 Story 的独立 `Requirements.FR` 行解析，避免把 `NFR-*` 误判为 `FR-*`；`TC-*` 和 `SB-*` 不参与覆盖计算。

### Coverage Matrix

| FR | 需求摘要 | 正式 Story 覆盖 | 状态 |
| --- | --- | --- | --- |
| FR-01 | 仅使用当前用户授权事实 | 1.2、1.8、1.12、1.14、1.16、1.17、1.25、1.26、1.34、1.37 | ✓ Covered |
| FR-02 | 统一漏洞事实与完整冻结集合 | 1.6、1.8、1.9、1.12～1.15、1.25、1.32～1.35、1.37 | ✓ Covered |
| FR-03 | 区分空、无权、数据不足和失败 | 1.2、1.6、1.8、1.12～1.15、1.17～1.18、1.25、1.30、1.33～1.34、1.37 | ✓ Covered |
| FR-04 | 全部新增漏洞倒序与空态 | 1.30、1.34～1.35、1.37 | ✓ Covered |
| FR-05 | 查事实并从全员列表人工选人 | 1.20、1.31～1.32、2.1～2.2、2.4 | ✓ Covered |
| FR-06 | 同接收人合并、不同人拆分 | 1.20、2.2、2.4 | ✓ Covered |
| FR-07 | 派发确认后写入并回读负责人 | 1.24、2.2、2.4 | ✓ Covered |
| FR-08 | 主动转发可见漏洞 | 1.20、2.3、2.5 | ✓ Covered |
| FR-09 | 转发使用全员列表并按人拆分 | 1.20、1.31、2.1、2.3、2.5 | ✓ Covered |
| FR-10 | 转发确认后逐条回读 | 1.24、2.3、2.5 | ✓ Covered |
| FR-11 | 延时与期限由人工沟通决定 | 1.20、3.1～3.2、3.4 | ✓ Covered |
| FR-12 | 管理员单条延时、原因可空 | 3.1～3.2、3.4 | ✓ Covered |
| FR-13 | 状态与期限双事实核验 | 1.24、3.1～3.2、3.4 | ✓ Covered |
| FR-14 | 人工判断自己负责漏洞误报 | 1.20、3.3、3.5 | ✓ Covered |
| FR-15 | 误报确认后状态回读 | 1.24、3.3、3.5 | ✓ Covered |
| FR-16 | 未确认／取消零外部写入 | 1.3～1.5、1.9、1.11、1.20～1.23、2.2～2.5、3.2～3.5 | ✓ Covered |
| FR-17 | 完成／待核验／待排查逐条分流 | 1.3～1.5、1.10～1.11、1.22～1.24、1.36、2.2～2.5、3.2～3.5 | ✓ Covered |
| FR-18 | 产品不发送通知 | 1.3、2.2～2.3、3.2～3.3 | ✓ Covered |
| FR-19 | 当前用户安全身份三态 | 1.25～1.26、1.37 | ✓ Covered |
| FR-20 | 当前权限与数据范围三态 | 1.25～1.26、1.37 | ✓ Covered |
| FR-21 | 按 IP 查询资产与在线状态 | 1.25、1.27、1.37 | ✓ Covered |
| FR-22 | 资产 opaque ref 详情 | 1.25、1.27、1.37 | ✓ Covered |
| FR-23 | 按 IP 查询关联漏洞 | 1.25、1.28、1.37 | ✓ Covered |
| FR-24 | 漏洞 opaque ref 详情 | 1.25、1.28、1.37 | ✓ Covered |
| FR-25 | `business_list` 全部业务系统 | 1.25、1.29、1.37 | ✓ Covered |

### Missing Requirements

无。没有 Epics 独有但 PRD 不存在的 FR，也没有 PRD FR 缺少正式 Story。

### Coverage Statistics

- Total PRD FRs: 25
- FRs covered in formal stories: 25
- Missing FRs: 0
- Coverage: **100%**

## 4. UX Alignment Assessment

### UX Document Status

**Found.** `EXPERIENCE.md` 是体验、旅程、状态与交互契约，`DESIGN.md` 是浅色桌面视觉系统与组件契约；四张关键屏 mockup 提供当前只读、冻结快照、目标写入确认和混合结果的可视验收材料。它们的职责互补，不构成重复 UX 文档。

### UX ↔ PRD Alignment

- UJ-01～UJ-05 覆盖新增漏洞查询、人工派发、主动转发、单条修复延时、人工误报和异常恢复；角色、人工决定、二次确认、逐条回读、统一漏洞事实与不发送通知均与 FR-01～FR-18 一致。
- 当前用户与权限三态、按 IP 资产／漏洞联查、业务系统列表及 opaque 详情引用进入 IA-01／IA-02 的只读事实流，与 FR-19～FR-25 一致。
- 多条操作的前 5 条摘要、“查看全部”、冻结全集、同接收人合并和不同接收人拆分与 PRD 一致；聊天只表达意图，程序使用 result_ref／draft 锁定事实。
- 写入结果统一按时间和预算分流：预算内 unknown/mismatch 为“待核验”；明确未应用或预算耗尽仍不能证明为“待排查”。UX 不允许从任一结果文案触发重复 mutation。
- 首版只支持桌面浅色，与 PRD 非目标和已批准 desktop-only ADR 一致；没有移动端、暗色、自动决策、通知或完整漏洞管理的隐性范围扩张。

### UX ↔ Architecture Alignment

- IA-01 聊天工作台、IA-02 漏洞操作工作区和固定“操作记录”入口均由 `product` 从同一 Product Facts／ActionResult 投影；前端不重建事实、审批、Run 或 provider DTO。
- AD-10／AD-21／AD-23 支持 ST-19／ST-27 的完成、待核验、待排查和部分成功；持久化 SSE、snapshot sequence 与 `view.replaced` 支持刷新、断线和上下文压缩后的确定恢复。
- AD-11 的不可变草案、SemanticMutationClaim 与 post-sent 只核验，以及 AD-26 的 actor-scoped 历史、服务端六个月窗口和快照一致 cursor，支持 UX 的确认与跨会话找回而不形成第二套真相。
- AD-27 的 `entity_ref → entity_subject_key → ProviderLocator` 支持列表进入安全详情；内部身份、locator、checkpoint、raw payload 和凭据不会进入 UI。
- 前端技术边界固定为 Vite + React + TypeScript、React Router、TanStack Query、Zustand、Radix primitives 与生成契约，支持 1440×900 浅色桌面三栏、100 条分页和可访问性状态播报。

### Alignment Issues

**无阻断性不一致。** 纠偏前“未知／不一致立即待排查”和 operation history 时间锚点分叉已在 UX、Architecture、SPEC 与正式 Story 中统一。

### Warnings

- OQ-02（真实人员写域规则）和 OQ-04（真实修复延时期限定规）仍是对应 Target Candidate 的转换门禁；它们不阻断当前正式 read-only／mock／证据 backlog，但阻断真实写 Story 转正。
- OQ-06～OQ-08 是后续恢复细节、快捷键和技术信息开放项；进入相关 Story 前必须按既有 gate 关闭，不能由开发者静默选择。
- UX 中的真实写入屏是门禁后的目标态，不表示当前 capability 可用；当前真实 mutation 继续关闭。

## 5. Epic Quality Review

### Epic Structure

| Epic | 独立用户价值 | 依赖检查 | 结论 |
| --- | --- | --- | --- |
| Epic 1：可信、可恢复的 FOBrain 事实工作台 | 运营可独立完成授权只读查询、冻结快照、百条核对、引用恢复和跨会话历史查看；真实写入关闭 | 不依赖 Epic 2/3 | 通过 |
| Epic 2：运营人工派发与责任人转发 | 在 Epic 1 事实和控制面之上提供人工选择、确认、负责人写入与回读证据 | 只依赖 Epic 1，不依赖 Epic 3 | 通过；真实能力仍由 Candidate 转换门禁控制 |
| Epic 3：管理员延时与运营误报处置 | 在 Epic 1 控制面之上提供人工延时／误报规则、mock 与目标部署证据 | 只依赖 Epic 1，不依赖 Epic 2 | 通过；真实能力仍由 Candidate 转换门禁控制 |

Epic 1 内含 M-0～M-5 技术使能 Story，但 Epic 标题、目标和完成信号是可独立验收的运营工作台，不是“建数据库／开发 API”式技术 Epic。每个使能 Story 都绑定用户或交付角色、可验证产物、架构决策和门禁；公共控制面只建立一次，Epic 2/3 禁止复制。

### Story Structure and Sizing

- 正式 Story：**47**（Epic 1=37、Epic 2=5、Epic 3=5）。
- Target Candidate：**6**；Source Bundle：**21**。两类均不使用 `Story` 标题，不进入 Sprint parser。
- 47/47 具备唯一 `As a / I want / So that`、Requirements metadata 和 Acceptance Criteria。
- 共验证 **117** 个 Given／When／Then 场景；37 个基础切片各有一个聚焦场景，Epic 2/3 的 10 个业务／证据 Story 各有 8 个场景，覆盖成功、安全失败、部分失败、回读、恢复和门禁记录。
- 每个正式 Story 为一个契约、repository、投影、恢复接缝、只读 capability、界面切片或单一业务动作证据；未发现把完整 Epic 塞入单个 Story 的情况。

### Dependency and Persistence Review

- 机械解析未发现任何正式 Story 指向更高编号的未来 Story；Epic 2/3 只引用已完成的 Epic 1 公共能力。
- M-0→M-5 顺序与 `docs/07-implementation-plan.md` 一致。Story 1.7 只建立 v2 additive migration、备份恢复和启动 readiness；Story 1.8～1.11 随首次使用分别交付安全结果、repository、事务和恢复材料，没有“先创建所有未来业务表”的无界基础设施 Story。
- Architecture 未要求外部 starter template；现有仓库是 Eino-first 新实现承载体。Story 1.1 已覆盖受支持 Go/Node 工具链、CI／本地基线和开工门禁。

### Traceability Integrity

- 质量审查发现并已修复一项阻断性元数据漂移：Epics 曾自行派生 NFR-05～10，并使 NFR-01～04 与 PRD 同号异义。
- 修复后 Epic inventory、47 个正式 Story、6 个 Candidate、追踪矩阵与纠偏执行记录只使用 PRD 规范 NFR-01～04；四条均至少有一个正式 Story，Architecture/UX 派生约束继续使用 AR/AD/UX ID。
- 当前正式 Story 覆盖 FR-01～25、NFR-01～04、UX-DR-01～30；Architecture metadata 只包含 `AD-*`／`M-*`，Gate metadata 独立保存。

### Findings by Severity

#### 🔴 Critical Violations

无。

#### 🟠 Major Issues

无未解决项。NFR 同号异义已在本轮审查中修复并通过重新解析。

#### 🟡 Minor Concerns / Scheduling Controls

- Epic 1 有 37 个小型顺序切片，Sprint Planning 必须按 M-0～M-5 依赖和验收门禁分批，不应一次性把全部标记为 in-progress。
- Story 2.4、2.5、3.4、3.5 是目标部署 Gate/Evidence Story；缺少明确环境授权、可恢复样本或失败联系人时应保持 blocked／不排入执行，不能用 mock PASS 替代。
- 六个 `TC-*` 即使文本完整也不是 Story；只有 ConversionGates 全部 PASS、重新运行 readiness 且生成全新正式 Story ID 后才能排入 Sprint。

### Best-Practice Verdict

**PASS。** Epic 具备独立用户价值，正式 Story 结构完整、大小可控、无前向依赖，数据库与公共控制面按首次需要建立，需求追踪已恢复为单一规范编号体系。

## Summary and Recommendations

### Overall Readiness Status

**READY — 可以进入 Sprint Planning。**

该结论表示 PRD、UX、Architecture、Runtime SPEC 与 47 个正式 Story 已达到一致、可追踪和可实施的规划状态；不表示真实 FOBrain 写域已获授权，也不替代任何 G-WRITE、G-FACT、G-SAFE、G-ARCH-V2、G-TOOLCHAIN 或环境验收记录。

### Critical Issues Requiring Immediate Action

无未解决的 Critical 或 Major 问题。

本轮曾发现的 NFR 编号漂移、Verifier evidence reducer、跨 Run Operation History 快照、entity subject 密钥轮换、cancel/sent 竞态以及 failed-retry claim 取消恢复等缺陷均已写入规范文档并闭合。最终 Architecture reviewer gate 的客观结论为：

- Reality review：**PASS**；目标架构与当前仓库现实、门禁和 Candidate 转换条件一致。
- Good-spine rubric：**PASS**，High=0、Medium=0。
- Incompatible-units review：**IMPLEMENTATION_SAFE**；不存在会令独立实现产生竞争权威状态或分页契约的未决接缝。

### Recommended Next Steps

1. 运行 Sprint Planning，只解析 `epics.md` 中 47 个正式 Story；6 个 `TC-*` 与 21 个 `SB-*` 必须排除。
2. 按 M-0→M-5 和 Story 依赖生成 backlog；首次开发只从满足门禁的 M-0/M-1 基础切片开始，不把 47 个 Story 一次性标记为进行中。
3. Story 2.4、2.5、3.4、3.5 在缺少目标环境授权、可恢复样本或失败联系人时保持 backlog，不以 mock、HTTP 2xx 或接口存在替代通过证据。
4. 每个开发任务开始前执行 `docs/pre-development-validation.md`，按影响范围补 schema、fixture、migration、contract、并发 CAS 与验收记录。
5. 任何真实写能力都必须先通过对应 Candidate ConversionGates，再生成全新正式 Story ID；不得原地把 Candidate 转为开发任务。

### Final Note

本次终态在文档发现、需求覆盖、UX 对齐、Epic 质量和架构兼容性五类检查中有 **0 项未解决阻断**。规划材料可安全进入 Sprint Planning；真实外部 mutation 继续保持关闭。

**Assessment date:** 2026-07-15  
**Assessor:** Codex / BMAD Implementation Readiness Reviewer
