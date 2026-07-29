---
stepsCompleted:
  - step-01-document-discovery
  - step-02-prd-analysis
  - step-03-epic-coverage-validation
  - step-04-ux-alignment
  - step-05-epic-quality-review
  - step-06-final-assessment
includedDocuments:
  - _bmad-output/planning-artifacts/product-blueprint/prd.md
  - _bmad-output/planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14/ARCHITECTURE-SPINE.md
  - _bmad-output/planning-artifacts/epics.md
  - _bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/DESIGN.md
  - _bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/EXPERIENCE.md
---

# Implementation Readiness Assessment Report

**Date:** 2026-07-15
**Project:** Agent Platform Eino

## Document Discovery

### PRD

- `_bmad-output/planning-artifacts/product-blueprint/prd.md`（8,497 bytes，2026-07-14 14:20:07）

### Architecture

- `_bmad-output/planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14/ARCHITECTURE-SPINE.md`（28,555 bytes，2026-07-14 16:23:07）

### Epics & Stories

- `_bmad-output/planning-artifacts/epics.md`（160,598 bytes，2026-07-15 11:37:56）

### UX Design

- `_bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/DESIGN.md`（22,543 bytes，2026-07-14 14:59:16）
- `_bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/EXPERIENCE.md`（32,824 bytes，2026-07-14 14:59:16）

### Discovery Result

- 四类必需文档均已找到。
- 未发现 whole document 与 `index.md` 分片版本并存的重复问题。
- UX 的 DESIGN 与 EXPERIENCE 分别承载视觉规范和交互体验，作为互补主文档共同纳入。
- Architecture/UX 的 review、mockup 与 validation report 作为辅助证据，不替代上述主文档。

## PRD Analysis

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

FR-17：任一写动作的失败、部分失败或回读不一致必须逐条显示为未完成／待排查；成功的同批对象可独立完成。

FR-18：本产品不得发送消息、邮件、webhook、催办或通知中心通知；FOBrain 自身通知不构成本产品功能或成功条件。

**PRD 明确功能需求总数：18。**

### Non-Functional Requirements

PRD 未使用 NFR 编号；以下 4 条跨功能边界按原文作为可验证的非功能／安全要求提取：

NFR-PRD-01：外部原始返回不能直接作为产品展示、审计或模型上下文事实；产品只消费安全、结构化的业务事实。

NFR-PRD-02：产品不得扩大 FOBrain 权限、代替 FOBrain 管理权限，或将查询不到的对象视为可操作对象。

NFR-PRD-03：所有数据修改均有执行人二次确认、客观回读和失败可见性；这些是产品安全边界，不是可选体验。

NFR-PRD-04：不把目标部署 API 已验证等同于产品已集成；不对写入结果、写域状态值域或权限作静默假设。

**PRD 明确编号 NFR：0；从跨功能边界提取的可验证 NFR：4。**

PRD 未承诺生产 SLA、吞吐、并发或商业化可用性；生产 SLA 明确属于本轮范围外。

### Additional Requirements

#### 用户与决策边界

- 产品中心／运营只在 FOBrain 内置权限可见范围内查看和处置漏洞。
- 管理员是执行直接修复延时的唯一授权角色。
- FOBrain 系统负责人负责排查写入失败、回读不一致及异常。
- AI 不是接收人、误报或延时的决策主体。
- 本产品不是通知服务、自动派发器、FOBrain 权限管理器或漏洞修复执行方。

#### 统一业务规则

- 统一漏洞卡片包含 POC、IP、IP 在线／离线状态、业务系统和修复负责人，并作为列表、确认和结果的共同事实展示。
- 二次确认必须绑定待执行对象集合和目标变化；单条展示完整对象，多条默认展示前 5 条且确认前可查看完整冻结集合。
- 写入成功必须通过外部事实回读核验，不能只使用接口返回成功。
- 写入失败、部分失败或回读不一致保留为逐条待排查状态，由 FOBrain 系统负责人排查。
- 产品只使用 FOBrain 接口返回给当前用户的可见对象范围，不扩大或自行推断权限。

#### 范围约束

- 范围内：7 个代表性只读能力、四项受控写动作、统一卡片、人工判断、二次确认、逐条回读和失败呈现。
- 范围外：自动派发及配置、自动决定接收人／误报／延时、其余 FOBrain 工具恢复、通知、完整漏洞管理、工单与修复闭环、公司主机在线盘点独立交付、生产 SLA 和商业化定位。

#### 实验判定

- 取消或未确认时没有外部写入记录。
- 每个成功项具有动作目标的回读事实。
- 每个失败或回读不一致项逐条保留待排查状态。
- 接收人、误报、延时及新期限均由操作人决定。
- 没有自动派发、通知或未纳入能力进入本轮。

#### 实施前阻塞项

1. 在当前产品中实现并 smoke 验证目标部署已通过的“全部新增漏洞”精确查询、倒序和空结果语义。
2. 在当前产品中实现并 smoke 验证 FOBrain 全部人员列表及接口内置可见范围。
3. 为派发、转发、延时、误报分别验证动作权限、写入、回读、部分失败及错误脱敏。
4. 验证统一卡片所需行级事实可以安全提供，并对真实缺失字段稳定显示为空／未分配。
5. 明确直接延时的新期限约束与目标部署可延时状态集合。

### PRD Completeness Assessment

- PRD 对实验目标、角色、四条写域旅程、统一事实、人工控制、成功判定、失败语义、范围和实施阻塞项定义清晰。
- 18 条 FR 均具有可客观核对的业务结果。
- PRD 的非功能要求没有独立编号，后续追踪必须同时覆盖跨功能边界和实施阻塞项，不能只追踪 FR-01～18。
- 具体工具链、持久化、运行时恢复、可访问性和 UI 细节不在 PRD 内定义，应由 Architecture、UX、SPEC 和 Additional Requirements 补足。
- PRD 状态为 `BASELINE_APPROVED_PENDING_IMPLEMENTATION_GATE`，明确不构成真实写域实施授权。

## Epic Coverage Validation

### Coverage Matrix

| FR | PRD Requirement | Epic / Story Coverage | Status |
| --- | --- | --- | --- |
| FR-01 | 产品必须只使用当前用户已授权的 FOBrain 事实展示漏洞及关联上下文。 | Epic 1，Story 1.12、1.17；全部后续动作继承 actor/scope 双检。 | ✓ Covered |
| FR-02 | 产品在漏洞列表、二次确认和结果展示中必须呈现统一漏洞事实。单条直接展示完整对象；多条时确认卡默认展示前 5 条并提供“查看全部”，使执行人在确认前可以核对完整冻结集合；不得仅用数量替代待写对象明细。 | Epic 1，Story 1.9、1.17、1.19；Epic 2/3 各动作 Story 复用同一冻结事实与确认投影。 | ✓ Covered |
| FR-03 | 产品必须区分真实空结果、无权限、数据不足和外部失败；不得由 AI 虚构“没有待处理漏洞”等结果。 | Epic 1，Story 1.2、1.7、1.12～1.17、1.20。 | ✓ Covered |
| FR-04 | 运营必须能够查看全部新增漏洞，按发现时间倒序；无结果时展示“没有待派发漏洞”。 | Epic 1，Story 1.11、1.16、1.17、1.20。 | ✓ Covered |
| FR-05 | 运营必须先查看漏洞、IP 在线状态、业务系统、业务系统负责人和运维负责人，再从 FOBrain 全部人员列表选择接收人；AI 不得替代该决定。 | Epic 1，Story 1.16、1.17、1.19；Epic 2，Story 2.1、2.2，目标实现 Story 2.7 Candidate。 | ✓ Covered；真实派发 gated |
| FR-06 | 仅已逐条判断且最终选择同一接收人的漏洞可合并派发；不同接收人不得混合。 | Epic 2，Story 2.2；目标实现 Story 2.7 Candidate。 | ✓ Covered；真实派发 gated |
| FR-07 | 派发仅在二次确认后写入；每条漏洞只有接口成功且回读修复负责人为指定人员时才完成。 | Epic 2，Story 2.2 mock、2.4 Gate/Evidence、2.7/2.9 Candidate。 | ✓ Covered；真实派发 gated |
| FR-08 | 当前用户可主动对其内置权限范围内的漏洞发起转发；不要求系统识别错派或先发生其他动作。 | Epic 2，Story 2.3 mock、2.5 Gate/Evidence、2.8/2.9 Candidate。 | ✓ Covered；真实转发 gated |
| FR-09 | 转发接收人必须从 FOBrain 全部人员列表选择；同一接收人可合并、不同接收人必须拆分。 | Epic 2，Story 2.1、2.3；目标实现 Story 2.8 Candidate。 | ✓ Covered；真实转发 gated |
| FR-10 | 转发仅在二次确认后写入；逐条回读负责人为目标人员时才完成。 | Epic 2，Story 2.3 mock、2.5 Gate/Evidence、2.8/2.9 Candidate。 | ✓ Covered；真实转发 gated |
| FR-11 | 运营与相关人员人工沟通后确定新修复期限；产品不得自动决定是否延时或期限。 | Epic 3，Story 3.1、3.2；目标实现 Story 3.6 Candidate。 | ✓ Covered；真实延时 gated |
| FR-12 | 直接延时一次只处理一条漏洞，延时原因可不填；只有管理员完成二次确认后可执行。 | Epic 3，Story 3.1、3.2、3.4；目标实现 Story 3.6 Candidate。 | ✓ Covered；真实延时 gated |
| FR-13 | 只有回读状态为“延时”且修复期限等于指定新期限时，结果才显示“延时成功”。 | Epic 3，Story 3.2 mock、3.4 Gate/Evidence、3.6/3.8 Candidate。 | ✓ Covered；真实延时 gated |
| FR-14 | 操作人可对自己确认负责的漏洞发起误报标记；误报判断必须由人作出，不能扩展为自动规则。 | Epic 3，Story 3.3 mock、3.5 Gate/Evidence、3.7/3.8 Candidate。 | ✓ Covered；真实误报 gated |
| FR-15 | 误报标记仅在二次确认后写入；只有回读状态为“误报”时才显示标记成功。 | Epic 3，Story 3.3 mock、3.5 Gate/Evidence、3.7/3.8 Candidate。 | ✓ Covered；真实误报 gated |
| FR-16 | 任一写动作在未完成二次确认或取消确认时，不得向外部系统写入。 | Epic 1，Story 1.9、1.10；Epic 2/3 的全部 mock、Gate/Evidence 和 Candidate 写域 Story。 | ✓ Covered；真实动作 gated |
| FR-17 | 任一写动作的失败、部分失败或回读不一致必须逐条显示为未完成／待排查；成功的同批对象可独立完成。 | Epic 1，Story 1.10；Epic 2，Story 2.2～2.9；Epic 3，Story 3.2～3.8。 | ✓ Covered；真实动作 gated |
| FR-18 | 本产品不得发送消息、邮件、webhook、催办或通知中心通知；FOBrain 自身通知不构成本产品功能或成功条件。 | Epic 2/3 的 mock、目标部署取证、Candidate 与 E2E Story 均含零产品通知断言。 | ✓ Covered；真实动作 gated |

### Epic FR Coverage Extracted

- Epic 1：FR-01、FR-02、FR-03、FR-04，以及 Epic 文档新增的 FR-19～25。
- Epic 2：FR-05、FR-06、FR-07、FR-08、FR-09、FR-10、FR-16、FR-17、FR-18。
- Epic 3：FR-11、FR-12、FR-13、FR-14、FR-15、FR-16、FR-17、FR-18。

### Missing Requirements

未发现 PRD FR-01～18 的 Epic/Story 覆盖缺口。

### Requirements Present in Epics but Not Numbered in PRD

Epic 文档新增 FR-19～25：当前用户身份、权限范围、按 IP 查资产、资产详情、按 IP 查漏洞、漏洞详情和业务系统列表。这 7 条与 PRD 的“7 个代表性只读能力”、AR-32、FOBrain 工具矩阵及 Owner 已批准实验范围一致，因此不是无依据扩项；但它们没有进入 PRD 的正式 FR 编号体系。

**Traceability risk:** PRD 声称是需求基线，但 Epic 以 FR-19～25 表达额外功能需求。进入 Sprint Planning 前应选择一个规范来源：将 FR-19～25 回写 PRD，或在追踪矩阵中明确标为 `AR/Scope-derived FR` 并记录批准来源，避免后续把 Epic 自身当作需求来源。

### Coverage Statistics

- PRD FR 总数：18
- 已在 Epics/Stories 中覆盖：18
- 缺失：0
- PRD FR 覆盖率：100%
- Epic 中额外编号 FR：7（FR-19～25）
- 覆盖不等于实施授权：FR-07～18 的真实写域完成路径仍受 Target Story Candidate 和 readiness gate 约束。

## UX Alignment Assessment

### UX Document Status

**Found.** 本次共同采用：

- `DESIGN.md`：浅色桌面视觉 token、组件视觉规则、状态色和视觉可访问性。
- `EXPERIENCE.md`：两表面信息架构、组件行为、ST-01～29、交互、可访问性和关键旅程。

两份文档状态均为 `final`，且明确互补，不是重复版本。

### UX ↔ PRD Alignment

一致项：

- UX 的 READ-01 与 UJ-01～UJ-05 对应 PRD 的只读链、派发、转发、延时、误报和异常处置旅程。
- 聊天是唯一动作意图入口；接收人、误报判断、是否延时和新期限均由人决定。
- 统一漏洞事实、完整冻结集合、前 5 条确认预览、逐条回读、部分失败隔离、零产品通知均与 FR-01～18 一致。
- 空结果、失败、无权限、数据不足和缺失字段使用确定性中文语义，与 PRD 安全边界一致。
- 浅色桌面、两表面、百条分页、键盘/读屏和状态语义属于 PRD 之外但与产品范围兼容的具体 UX 约束。

未发现 UX 擅自加入自动派发、AI 业务判断、通知、工单闭环、移动端或暗色主题。

### UX ↔ Architecture Alignment

一致项：

- Vite + React + TypeScript + React Router + TanStack Query + Zustand + Radix Primitives + lucide-react 与架构栈一致。
- 232px / `minmax(600px, 1fr)` / 360px 浅色桌面三栏和 IA-01/IA-02 两表面由 AD-17 支撑。
- 只读 IA-02、聊天审批卡、Product Facts 同源投影、generated contracts、SSE 恢复、不可变草案、确认、幂等、verifier 和逐条结果均有对应架构不变量。
- 628 条级样例通过冻结快照和每页 100 条分页处理；PRD 明确不承诺生产 SLA，因此没有缺失一个已承诺的响应时间指标。

### Alignment Issues

#### UX-A01 — OQ-05 的目标决定未回写 UX/Architecture（High）

- `epics.md` Story 2.6 已规定左侧导航固定入口“操作记录”、最近六个月可发现窗口和跨会话重新打开逐条结果。
- `EXPERIENCE.md` 仍将 OQ-05 标为未决，并只定义 IA-01/IA-02 两张表面；ST-29 只说从会话导航重新打开历史结果，未定义独立“操作记录”入口及六个月范围。
- Architecture AD-23/AD-25 只定义事实保留和不自动物理删除，没有定义六个月产品可发现窗口，也未说明“操作记录”如何保持两表面约束。

**Impact:** Story 2.6 实施者可能新增第三张中栏表面、把历史列表塞入不合适区域，或采用与 Product Facts 保留策略冲突的删除逻辑。

**Required resolution before Story 2.6 implementation:** 更新 EXPERIENCE/UX-DR、Architecture 或 ADR，明确“操作记录”是左栏导航分组还是 IA-01 的历史模式、选择记录后如何进入 IA-02、六个月只影响产品发现还是物理保留，并正式关闭 OQ-05。

#### UX-A02 — 回读不一致的“待核验/待排查”转换边界不一致（High）

- Architecture AD-10/AD-23：短期回读未一致进入 `reconciling`，达到 verification deadline 后才进入 `manual_attention` 或确定失败。
- Epic Candidate 2.7、2.8、3.6、3.7：超时、未知或回读不一致先显示“待核验”并进入 reconcile。
- UX ST-19 与 PRD 统一词汇：回读不一致直接显示“未完成／待排查”。UX ST-27 只覆盖结果未知，没有明确暂时不一致的核验窗口。

**Impact:** 同一事实可能在 Workbench、Action API、历史结果和验收脚本中得到不同产品状态，违反 Product Facts 单一投影原则。

**Required resolution before真实写域:** 在 Product Facts 状态映射、UX ST-19/ST-27 和动作 AC 中统一：verification deadline 前的暂时不一致如何显示，deadline 后如何从“待核验”收口为“待排查”或 `manual_attention`，并固定中文文案和机器状态映射。

### Warnings and Intended Gates

- OQ-02 与 OQ-04 仍在 UX 中开放，但 Epic 2.1 与 3.1 已将其作为先行决策/门禁 Story；这不会阻塞当前 READ-01，只阻塞相应真实动作。
- OQ-08 未关闭时技术信息整体隐藏，Epics 与 Architecture 均保持该边界。
- OQ-06/07 属于后续离线与效率增强，不在当前可启用范围内。
- Architecture Spine 的 `sources` 显式列出 DESIGN 但未列出 EXPERIENCE；正文实际采用了 Experience 规则。建议后续同步来源清单，避免审计时误判 EXPERIENCE 未纳入架构输入。

## Epic Quality Review

### Structure Summary

- Epic：3
- Story headings：37
- Acceptance Criteria：298，全部具有 Given/When/Then。
- 正式/Decision/Gate/Evidence Story：31。
- Target Story Candidate：6（2.7～2.9、3.6～3.8）。
- Story 内显式出现 FR 编号：6；未显式出现 FR 编号：31。

### Epic Compliance

| Epic | User Value | Independence | Main Quality Result |
| --- | --- | --- | --- |
| Epic 1：可信、可恢复的 FOBrain 事实工作台 | 明确交付可独立使用的只读查询与事实核对。 | 可独立完成，不需要 Epic 2/3。 | Epic 合格；部分 foundation Story 过大。 |
| Epic 2：运营人工派发与责任人转发 | 目标用户价值明确。 | 只依赖 Epic 1，不依赖 Epic 3。 | 真实用户价值位于 Candidate；下游必须可靠排除 Candidate。 |
| Epic 3：管理员延时与运营误报处置 | 目标用户价值明确。 | 不依赖未来 Epic；依赖 Epic 1，并实际复用 Story 2.6 的全局 OQ-05 决策。 | 依赖描述需收口；真实用户价值位于 Candidate。 |

### 🔴 Critical Violations

#### EQ-C01 — Candidate 使用标准 Story heading，存在被 Sprint Planning 自动纳入的风险

6 个不可执行 Candidate 仍使用 `### Story N.M` 标题和标准 AC 结构。虽然正文明确“不得进入 Sprint”，但下游工具若按 Story heading 扫描，可能把 2.7～2.9、3.6～3.8 当作正式待开发 Story。

**Impact:** 直接破坏用户要求的阶段门禁，使真实 mutation 在 G-WRITE/G-SAFE/OQ 未通过时进入开发队列。

**Required remediation:** Sprint Planning 前将 Candidate 移到机器可识别的独立区段/状态字段，或给 sprint workflow 增加经过验证的排除规则；生成 sprint-status 后必须断言这 6 项不存在。Candidate 转正式 Story 时再创建正式 Story 文件并记录证据、日期和批准人。

#### EQ-C02 — 多个 Foundation Story 超出单一开发 Agent 的合理切片

以下 Story 在一个交付中同时跨越多个独立组件、契约或故障域：

- Story 1.3：Action、Catalog、Run、Pending/Continuation、Verification、Action/Resume/Result、SSE、OpenAPI、Go/TypeScript generated contracts。
- Story 1.6：全部 v2 聚合 persistence、FactEvent、幂等、lease/fencing、retention、恢复扫描和 migration。
- Story 1.7：后端 product projection、HTTP view/SSE、v1/v2 adapter、LLM context 和前端 stream reducer。
- Story 1.8：Eino Registry/ToolsConfig、Agent loop、LLM adapter、observability、bootstrap 和安全测试。
- Story 1.9：Prepare/Confirm、Eino continuation、Action API continuation、事务确认和崩溃恢复。
- Story 1.10：mutation 幂等、executor lease/fencing、reconcile、Verifier 和结果聚合。
- Story 1.17：统一漏洞多源映射、QuerySnapshot、Conversation Result Index、Reference Resolver 和跨重启恢复。

**Impact:** 单 Story 难以在一个开发 Agent 会话内实现、验证和评审；失败时难以定位是契约、存储、runtime、projection 还是恢复问题，也会造成大范围文件修改。

**Required remediation:** Sprint Planning 前按可独立验收的纵向/基础切片拆分；保持 M-0→M-5 顺序和现有不变量，不把公共基础拆回四个业务动作中。每个拆分 Story 必须有自己的 schema/fixture/test gate，并且只依赖更早切片。

### 🟠 Major Issues

#### EQ-M01 — Story 级需求追踪不足

37 个 Story 中只有 6 个在正文显式引用 FR 编号，31 个主要依赖 Epic 级覆盖表、AR/M/Gate 名称或语义匹配。

**Impact:** 人工可以理解，但 `create-story`、验收报告和变更影响分析无法稳定证明每条 Story 实现哪些 FR/NFR/AR/UX-DR。

**Required remediation:** 为每个正式 Story 增加机器可读的 `Requirements:` 或等价字段，至少列出适用 FR、NFR、AR/AD、UX-DR 和 gate。技术 foundation Story 可以不绑定业务 FR，但必须绑定明确 NFR/AR/AD/Milestone。

#### EQ-M02 — Epic 3 的依赖声明与实际全局决策位置不一致

Epic 3 实施说明声称“不依赖 Epic 2”，但真实延时/误报 Candidate 要求 OQ-05 关闭，而 OQ-05 的固定入口与六个月决策仅在 Story 2.6 中定义。

**Impact:** 若团队并行或跳过 Epic 2 mutation，可能误以为 Epic 3 不需要 Story 2.6 的跨会话结果决策。

**Required remediation:** 明确 Epic 3“不依赖 Epic 2 的派发/转发 mutation，但依赖全局 OQ-05 决策产物”；或者把 Story 2.6 移到共享 foundation Epic/Story。

#### EQ-M03 — Candidate 转正式时仍需要重新切片

Story 2.7、2.8、3.6、3.7 各自同时包含生产 provider mutation、policy、verifier、Workbench/Action API 投影、历史入口、错误恢复和批准环境 E2E。

**Impact:** 即使门禁通过，直接“转换标题状态”仍会产生过大 Story。

**Required remediation:** Candidate 是目标 AC 容器，不应原样交给开发 Agent。转换时至少拆为 provider/policy、执行与回读、产品投影/UX、批准环境验收等可独立评审切片，并确保前一个切片本身可验证。

### 🟡 Minor Concerns

#### EQ-m01 — Foundation Story 的用户价值表达偏技术

Story 1.1、1.4 等使用项目维护者作为用户并交付工具链/migration readiness。它们有明确安全价值且是架构门禁需要，因此不构成技术 Epic，但在 Sprint 说明中应明确它们是实现用户价值的必要 enabling Story，不以“基础设施完成”冒充产品价值完成。

#### EQ-m02 — 显式未来 Story 引用需要保持“禁止性”解释

Story 1.4 对 1.5/1.6 的引用仅说明延后建表；Story 2.6、3.1 对未来 Candidate 的引用仅用于阻止提前开发。当前没有发现真正依赖未来输出的前向依赖，但拆分和生成 Story 文件时应保留这种含义，避免转写成 `depends_on`。

### Passed Best-Practice Checks

- 三个 Epic 标题与目标均表达用户结果，不是纯技术里程碑。
- Epic 1 可独立交付只读价值；Epic 2 不依赖 Epic 3；没有循环依赖。
- 数据库创建时机正确：Story 1.4 只建立 migration/readiness，StructuredResult 与聚合表分别在 Story 1.5/1.6 首次需要时创建。
- Architecture 未指定必须从某个 starter template 初始化；现有仓库已存在，Story 1.1 负责工具链、CI 和可复现安装，符合当前项目形态。
- 所有 Story 都有 As a / I want / So that 和可测试 BDD AC，且普遍覆盖错误、安全、恢复与边界场景。
- 公共 Action 状态机、确认、幂等、reconcile 和 verifier 集中在 Epic 1，没有在 Epic 2/3 重复创建。

## Summary and Recommendations

### Overall Readiness Status

## NOT READY

当前规划材料已经形成完整产品蓝图，PRD FR 覆盖率为 100%，架构与 UX 主方向一致；但它尚不适合直接进入完整 Sprint Planning 或把全部 Story 交给开发 Agent。

允许的下一步是修正规划缺口、继续门禁/证据工作，以及按唯一 readiness gate 已明确许可的基础工作。不得把本报告、100% FR 覆盖或 Candidate 的完整 AC 解释为真实写域实施授权。

### Critical Issues Requiring Immediate Action

1. **Candidate 可能被自动排入 Sprint。** 2.7～2.9、3.6～3.8 使用标准 Story heading；必须建立机器可验证的排除机制，并证明 sprint-status 不包含它们。
2. **Foundation Story 过大。** Story 1.3、1.6、1.7、1.8、1.9、1.10、1.17 不满足单一开发 Agent 的合理切片，应在 Sprint Planning 前拆分。
3. **OQ-05 决定未跨文档收口。** “操作记录 + 六个月”只存在于 Story 2.6；UX、Architecture 和两表面落位仍未更新。
4. **回读不一致状态映射冲突。** UX 的“待排查”与 Architecture/Epics 的“待核验 → reconcile → 终态”缺少统一转换规则。

### Major Issues Requiring Resolution

5. **FR-19～25 不在 PRD 正式编号体系。** 必须回写 PRD，或在追踪矩阵中明确为 Owner 批准的 scope/AR-derived requirements。
6. **Story 级追踪不足。** 31/37 Story 没有显式 requirements metadata，技术与产品验收难以自动追踪。
7. **Epic 3 依赖描述不准确。** 它不依赖 Epic 2 的 mutation，但依赖目前位于 Story 2.6 的全局 OQ-05 决策。
8. **真实动作 Candidate 转正式时仍然过大。** 不能只删除 Candidate 标签；必须按 provider/policy、执行/回读、产品投影和批准环境验收重新切片。

### Recommended Next Steps

1. **先修 Sprint 安全边界：** 将 Candidate 改为机器可识别的非 Story 结构，或为 Sprint Planning 增加并验证硬排除规则；生成一个测试 sprint-status 证明 6 个 Candidate 均未进入。
2. **重拆 Foundation Story：** 对 1.3、1.6～1.10、1.17 分别按单一职责和独立验收证据拆分，保持 M-0→M-5 顺序，不改变既有架构不变量。
3. **补齐 Story metadata：** 为全部正式 Story 标注 FR/NFR/AR/AD/UX-DR/Gate；Foundation Story 至少绑定明确 NFR/AR/AD/Milestone。
4. **关闭跨文档 UX/状态冲突：** 更新 EXPERIENCE、Architecture/ADR、Product Facts 状态映射和相关 AC，固定 OQ-05 的两表面落位、六个月发现窗口，以及 `待核验 → 待排查/manual_attention` 的时间边界。
5. **统一需求基线：** 将 FR-19～25 纳入 PRD 或追踪矩阵的正式派生需求区，并记录 Owner 批准来源。
6. **修正 Epic 3 依赖说明与 Architecture source 清单。** 明确共享决策依赖，并把 EXPERIENCE 纳入架构来源。
7. **重新运行 Implementation Readiness。** 只有 Critical/High 问题关闭且报告转为 READY/受限 READY 后，才运行 Sprint Planning。
8. **真实写域继续独立受门禁约束。** 即使规划问题全部修复，派发、转发、延时和误报仍必须分别通过 G-WRITE/G-SAFE/OQ、环境授权和 Candidate 转换流程。

### Final Note

本次评估识别出 **8 个需要在 Sprint Planning 前处理的关键/高/主要问题**，分布在需求追踪、UX/Architecture 对齐和 Epic/Story 质量三个类别；另有若干轻微文档与表达问题。当前材料足以指导修正工作和受限基础门禁工作，但不足以授权完整产品实现或任何真实 FOBrain mutation。

**Assessment date:** 2026-07-15  
**Assessor:** Codex / BMAD Implementation Readiness
