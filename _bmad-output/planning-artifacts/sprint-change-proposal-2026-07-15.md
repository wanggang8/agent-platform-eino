---
title: Agent Platform Eino 实施准备纠偏提案
status: PLANNING_CORRECTION_COMPLETED_SPRINT_PLANNED
date: '2026-07-15'
mode: Batch
changeScope: Major
approvedBy: Vick
approvedAt: '2026-07-15'
trigger:
  - _bmad-output/planning-artifacts/implementation-readiness-report-2026-07-15.md
---

# Sprint Change Proposal

> 本文是纠偏提案，不是开发授权。提案获批并完成文档修订、再次通过 Implementation Readiness 前，不得运行 Sprint Planning，也不得启用任何真实 FOBrain 写入能力。

## 1. Issue Summary

### 1.1 触发原因

项目尚未进入开发。实施准备审查发现，现有 PRD、UX、Architecture、SPEC 与 Epics 已覆盖产品目标，但还不能稳定地被 Sprint Planning 和开发 Agent 消费。问题属于“已确认需求没有被完整、统一地操作化”，不是新需求、市场转向或代码实现失败。

### 1.2 核心问题

当前计划存在八项需要在开发前修正的问题：

1. 6 个 Target Story Candidate 使用标准 `### Story x.y` 标题，会被 Sprint Planning 自动解析为正式 Story。
2. Epic 1 的 Story 1.3、1.6、1.7、1.8、1.9、1.10、1.17 同时跨越多个层、契约和验收面，不能作为一次可完成的开发单元。
3. Story 2.6 已决定“操作记录 + 六个月”，但 UX 和 Architecture 仍保留 OQ-05，未形成统一产品决策。
4. UX 的 ST-19 把所有回读不一致直接映射为“待排查”，而 Architecture、SPEC 与 Epics 存在 `reconciling/待核验` 过渡态。
5. Epics 使用 FR-19～FR-25 表达 7 个已批准的代表性只读能力，但 PRD 正式编号只到 FR-18。
6. 37 个当前 Story 中有 31 个缺少显式 requirements metadata，无法稳定追踪 FR/NFR/AR/AD/UX/Gate。
7. Epic 3 声称不依赖 Epic 2，但两个写域共同依赖的 OQ-05 决策只放在 Story 2.6。
8. 真实动作 Candidate 同时包含 provider、policy、UX 和 E2E；即使门禁通过，也不能原样转换成单个正式 Story。

### 1.3 客观证据

- `bmad-sprint-planning` 明确按 `### Story 1.1: ...` 形式提取 Story；当前 Candidate 标题与该模式相同。
- Implementation Readiness 对 PRD FR-01～FR-18 的覆盖率为 100%，说明无需删减既定业务目标。
- `EXPERIENCE.md` 的 OQ-05 仍为开放问题，而 `epics.md` Story 2.6 已规定固定入口、六个月范围和跨会话访问。
- `EXPERIENCE.md` ST-19、`ARCHITECTURE-SPINE.md` AD-21/AD-23、SPEC state machine 对不一致结果的收敛时点不一致。
- 当前没有 sprint-status、已完成 Story、生产代码或已启用写动作，不存在代码回滚成本。

## 2. Impact Analysis

### 2.1 Epic Impact

| Epic | 目标是否保留 | 必要调整 |
| --- | --- | --- |
| Epic 1：可信、可恢复的 FOBrain 事实工作台 | 保留 | 将 7 个过大基础 Story 拆为单职责、可验收 Story；承接跨 Epic 共用的“操作记录”能力。 |
| Epic 2：人工选择接收人的派发与主动转发 | 保留 | 移出 Story 2.6；将 2.7～2.9 改为非 Sprint Candidate；真实能力门禁后按三层切片转换。 |
| Epic 3：人工判断后的延时与误报处置 | 保留 | 明确只依赖 Epic 1 的共用控制面与历史结果能力；将 3.6～3.8 改为非 Sprint Candidate；门禁后拆分转换。 |

不新增、不删除、不重排 Epic。MVP 目标和四类动作均不改变。

### 2.2 Story Impact

- Epic 1 由 20 个 Story 重排为 37 个更小的 Story；旧 ID 不能继续作为稳定引用，必须一次性更新追踪矩阵和文档链接。
- Story 2.6 移到 Epic 1，作为所有写域共用的操作历史投影能力。
- Epic 2 和 Epic 3 各保留 5 个当前可规划 Story。
- 6 个 Target Story Candidate 不再使用 Story 编号或 Story 标题；它们是门禁后的目标需求容器，不是当前 backlog。
- 所有正式 Story 增加统一 requirements metadata；Candidate 增加 Target Requirements metadata。

### 2.3 Artifact Conflicts

| Artifact | 冲突 | 修正 |
| --- | --- | --- |
| PRD | FR-19～FR-25 缺少正式来源；FR-17 未区分暂时待核验与最终待排查。 | 正式加入 F-07 代表性只读能力；修订 FR-17；为跨功能边界补稳定 NFR ID。 |
| Architecture | 未固定历史入口；回读不一致的时间边界不完整；来源未包含 UX 主文档。 | 新增共用操作历史决策；收紧 AD-21/AD-23；同步 Source Map。 |
| UX EXPERIENCE | OQ-05 未关闭；ST-19 与 ST-27 冲突。 | 新增历史模式与结果入口决策；统一 ST-19/ST-27/ST-28；关闭 OQ-05。 |
| UX DESIGN | 只有会话导航示例，缺少固定操作记录入口及历史模式状态。 | 补导航、筛选、只读历史和三态结果视觉规则；更新相关 mockup。 |
| Epics | Candidate 会被解析；基础 Story 过大；追踪元数据缺失；Epic 3 依赖表述矛盾。 | 机器隔离 Candidate、重切 Epic 1、移动 Story 2.6、增加 metadata、更新依赖和映射。 |
| Runtime SPEC | `reconciling`、`manual_attention` 与产品文案映射含混。 | 修订 state machine 和 acceptance matrix，固定过渡态与终态。 |
| Traceability/Gates | 只追踪 FR-01～18；当前 readiness 仍为阻塞。 | 扩展到 FR-25 和正式 NFR；同步 gate；纠偏后重新运行 readiness。 |

### 2.4 Technical Impact

当前没有需要回滚的实现、数据库或部署。纠偏会改变未来的 schema、状态投影、SSE、操作历史查询和 Story 执行边界，因此必须在代码开始前完成。已存在的 Implementation Readiness 报告作为历史证据保留，不覆盖修改；完成修订后生成新的报告。

## 3. Recommended Approach

### 3.1 选项比较

| 选项 | 可行性 | 工作量 | 风险 | 结论 |
| --- | --- | --- | --- | --- |
| 1. 直接调整现有计划 | 可行 | 中高 | 低至中 | 推荐。业务范围不变，在开发前统一基线和 backlog。 |
| 2. 回滚已有工作 | 不适用 | 低 | 中 | 当前没有已实施 Story 或代码可回滚，不能解决跨文档冲突。 |
| 3. 缩减或重定义 MVP | 可行但无必要 | 中 | 高 | 18 条原始 FR 已 100% 覆盖，删范围不能解决机器解析和状态语义问题。 |

### 3.2 选择

采用“直接调整 + 基线规范化”。变更按 BMAD 分类为 **Major**，原因是需要 PM、Architect、UX 和 backlog 共同修订；它不是产品范围扩大，也不需要重新做产品发现。

### 3.3 MVP 与进度影响

- MVP：不变。
- 真实写域：继续阻塞，门禁不降低。
- 当前可实施范围：纠偏文档、契约/fixture 规划和 readiness 复核；在复核通过前不进入开发。
- 预计影响：增加一次文档重排和一次 Implementation Readiness 复核，避免开发后跨层返工。

## 4. Detailed Change Proposals

### 4.1 让 Candidate 无法被 Sprint Planning 解析

**Artifact:** `epics.md` Target Story Candidate 标题与引用

**OLD**

```markdown
### Story 2.7（Target Story Candidate）：接入真实手动派发能力
### Story 2.8（Target Story Candidate）：接入真实主动转发能力
### Story 2.9（Target Story Candidate）：完成派发与转发端到端验收
### Story 3.6（Target Story Candidate）：接入真实单条修复延时能力
### Story 3.7（Target Story Candidate）：接入真实人工误报标记能力
### Story 3.8（Target Story Candidate）：完成延时与误报端到端验收
```

**NEW**

```markdown
## Target Candidates（非 Sprint 工作项）

### Target Candidate TC-2.1：接入真实手动派发能力
### Target Candidate TC-2.2：接入真实主动转发能力
### Target Candidate TC-2.3：派发与转发目标验收包

### Target Candidate TC-3.1：接入真实单条修复延时能力
### Target Candidate TC-3.2：接入真实人工误报标记能力
### Target Candidate TC-3.3：延时与误报目标验收包
```

**转换规则**

1. Candidate 不使用 `Story` 词、`Epic.Story` 数字或正式 Story 状态。
2. Sprint Planning 前运行结构检查，断言 `TC-*` 不出现在 sprint status 中。
3. 门禁通过后不原地改名；由 Create Story 生成新的正式 Story ID 和独立 Story 文件。
4. TC-2.1、TC-2.2、TC-3.1、TC-3.2 分别按三层转换：
   - provider/capability + policy/verifier；
   - Product Facts projection + 受控 UX 流程；
   - 获授权环境的 integration/E2E + gate evidence。
5. TC-2.3、TC-3.3 只保存跨动作验收目标和证据清单，不作为一个包含全部实现的 Story。

**理由：** 结构上阻止自动排期，同时避免门禁通过后仍产生横跨 provider、execution、product、web 和 E2E 的巨型 Story。

### 4.2 将 Epic 1 重切为可完成单元

**Artifact:** `epics.md` Epic 1 Story inventory

**OLD**

- Story 1.3 同时定义 Action、Catalog、Run、Pending、Continuation、Verification、HTTP、SSE、OpenAPI 和 generated contract。
- Story 1.6 同时定义全部 v2 聚合、事件、幂等、租约、恢复和清理。
- Story 1.7～1.10 分别跨越后端投影、前端 reducer、Eino、LLM、审批、执行、租约、回读和 reconcile。
- Story 1.17 同时完成多来源映射、统一行事实、快照、result index、引用解析和恢复。

**NEW：Epic 1 规范 Story inventory**

| 新 ID | Story 目标 | 来源 |
| --- | --- | --- |
| 1.1 | 固定可复现工具链与验收基线 | 保留旧 1.1 |
| 1.2 | 建立安全 StructuredResult 与查询快照 v2 契约 | 保留旧 1.2 |
| 1.3 | 建立 Catalog 与 Action domain schema/fixture | 拆旧 1.3 |
| 1.4 | 建立 Run、Pending、Continuation 与 Verification domain schema/fixture | 拆旧 1.3 |
| 1.5 | 固定 Action API、Resume、Result、OpenAPI 与 generated contract | 拆旧 1.3 |
| 1.6 | 固定 SSE event/patch schema 与 generated contract | 拆旧 1.3 |
| 1.7 | 建立 SQLite additive migration 与启动 readiness | 旧 1.4 |
| 1.8 | 持久化完整安全 StructuredResult 并兼容读取 v1 | 旧 1.5 |
| 1.9 | 持久化 QuerySnapshot 与 Action 聚合 repository | 拆旧 1.6 |
| 1.10 | 固定 FactEvent sequence 与聚合事务一致性 | 拆旧 1.6 |
| 1.11 | 持久化 continuation、attempt、lease 与启动恢复材料 | 拆旧 1.6 |
| 1.12 | 建立后端 Product Facts 投影与 v1/v2 union | 拆旧 1.7 |
| 1.13 | 建立持久化 SSE view、游标与重连投影 | 拆旧 1.7 |
| 1.14 | 建立 LLM safe Product Facts context assembler | 拆旧 1.7 |
| 1.15 | 建立前端 SSE reducer、乱序保护与重连恢复 | 拆旧 1.7 |
| 1.16 | 建立 Capability Registry 与 Eino ToolsConfig adapter | 拆旧 1.8 |
| 1.17 | 建立 Eino Runner 与受控工具循环 | 拆旧 1.8 |
| 1.18 | 建立 LLM provider interface、网络策略与错误脱敏 | 拆旧 1.8 |
| 1.19 | 建立 observability budget、redaction 与 bootstrap wiring | 拆旧 1.8 |
| 1.20 | 建立不可变 PrepareAction、policy 与 ActionDraft | 拆旧 1.9 |
| 1.21 | 建立 ConfirmAction、approval 消费与双入口 continuation | 拆旧 1.9 |
| 1.22 | 建立 mutation 幂等与 attempt reservation | 拆旧 1.10 |
| 1.23 | 建立 executor lease、fencing 与崩溃接管 | 拆旧 1.10 |
| 1.24 | 建立 verifier、reconcile 与 ActionResult 终态投影 | 拆旧 1.10 |
| 1.25 | 关闭 G-READ/G-FACT 数据可得性证据 | 旧 1.11 |
| 1.26 | 识别当前用户与 FOBrain 权限范围 | 旧 1.12 |
| 1.27 | 按 IP 查询资产并查看资产详情 | 旧 1.13 |
| 1.28 | 按 IP 查询漏洞并查看漏洞详情 | 旧 1.14 |
| 1.29 | 查询全部业务系统 | 旧 1.15 |
| 1.30 | 精确查询全部新增漏洞 | 拆旧 1.16 |
| 1.31 | 查询并安全选择 FOBrain 全员 | 拆旧 1.16 |
| 1.32 | 建立统一漏洞事实、多来源映射与冻结快照 | 拆旧 1.17 |
| 1.33 | 建立 Conversation Result Index、指代解析与恢复 | 拆旧 1.17 |
| 1.34 | 交付只读浅色桌面聊天 Workbench | 旧 1.18 |
| 1.35 | 交付百条分页明细与右侧事实面板 | 旧 1.19 |
| 1.36 | 交付共用操作记录与跨会话结果入口 | 移入旧 2.6 |
| 1.37 | 完成 READ-01 端到端验收与门禁记录 | 旧 1.20 |

每个新 Story 必须只包含一个主要职责、明确受影响目录、输入/输出契约、测试命令和门禁结果；任何 Story 不允许依赖尚未完成的后续 Story。

**理由：** 让 Create Story 和 Dev Story 能在一个上下文中实现、验证和回滚单一职责，同时保持 AGENTS.md 的层级边界。

### 4.3 把操作历史移为共用基础能力并关闭 OQ-05

**Artifacts:** `epics.md`、`EXPERIENCE.md`、`DESIGN.md`、`ARCHITECTURE-SPINE.md`

**OLD**

```markdown
### Story 2.6：固定跨会话操作记录与待处理结果入口
OQ-05：离开会话或重新登录后，用户从哪里重新找到历史逐条结果与待排查对象？
Epic 3 ... 不依赖 Epic 2 ... OQ-05 关闭前不能完成。
```

**NEW**

```markdown
### Story 1.36：交付共用操作记录与跨会话结果入口

RD-03 / OQ-05 CLOSED：
- 左侧导航提供系统固定入口“操作记录”，它不是聊天会话。
- 入口复用 IA-01 主表面的历史模式，不新增第三个产品表面。
- 选择记录后进入既有 IA-02 漏洞操作工作区，展示只读逐条结果。
- 默认展示最近六个月并按操作时间倒序；待核验、待排查和 manual_attention 对象在关闭前不受六个月默认窗口限制。
- 历史只从 Product Facts/ActionResult 投影；不从聊天摘要、raw provider payload 或重新调用 FOBrain 重建。
- 每次访问重新校验 actor/workspace；历史不自动生成再次执行或重试能力。
- 六个月是默认产品可发现窗口，不是物理删除策略；物理保留服从 Architecture 的事实/audit 策略。
```

Architecture 新增 `AD-26 [ADOPTED TARGET] — 操作记录是 Product Facts 的安全投影`，并由 Story 1.36 实现。Epic 2 和 Epic 3 的实施说明统一为“直接依赖 Epic 1；彼此不依赖”。

**理由：** 这是四类动作共同依赖的结果发现能力，应由共用事实层提供，不应让 Epic 3 间接依赖 Epic 2。

### 4.4 统一待核验与待排查的状态边界

**Artifacts:** PRD FR-17、Architecture AD-21/AD-23、EXPERIENCE ST-19/ST-27/ST-28、SPEC state machine/acceptance matrix、全部写域 AC

**OLD**

```markdown
FR-17：写入失败、部分失败或回读不一致均逐条显示未完成／待排查。
ST-19：只要回读目标未达成，该条仍为未完成／待排查。
Epics：未知结果或回读不一致显示待核验并进入 reconcile。
```

**NEW**

| 事实条件 | 内部 item 状态 | Run/ActionResult | 产品文案 | 是否允许再次 mutation |
| --- | --- | --- | --- | --- |
| control certainty 已丢失（包括失租约） | `manual_attention` | Run=`failed`、ActionResult=`blocked` | 待排查 | 否 |
| control held + provider accepted + conclusive 真实 readback 与 expected 一致 | `succeeded` | 按全部 item 汇总 | 完成 | 否 |
| control held + provider rejected + policy 接受 conclusive non-application | `failed` | 按全部 terminal item 汇总 `partial/none_succeeded` | 待排查 | 否；只允许重新发起新的人工决策流程 |
| 其他组合且仍在 verification policy 的 settle deadline/调用预算内 | `reconciling` | Run=`running` | 待核验 | 否 |
| 其他组合在 deadline/预算结束后仍无法证明目标 | `manual_attention` | Run=`failed`、ActionResult=`blocked` | 待排查 | 否 |

PRD FR-17 修订为：

> 任一写动作必须逐条区分完成、待核验和待排查。当前无法确认且仍可沿同一执行事实链继续核验时显示“待核验”；明确失败，或核验期限/预算结束后仍不一致、未知时显示“待排查”。成功项可独立完成；任何未完成项不得伪报成功，也不得通过重复 mutation 核验。

ST-19 改为“核验期限内不一致进入 ST-27 待核验；达到 policy 收敛条件仍不一致才进入待排查”。SPEC AC-13 同步区分 transient reconcile 和 terminal manual attention。

**理由：** 同时满足“接口成功不等于业务成功”“结果未知不能重复写入”和“最终异常必须可排查”。

### 4.5 将 7 个代表性只读能力正式写回 PRD

**Artifact:** `product-blueprint/prd.md`

**OLD**

```markdown
## 6. 本轮范围与非目标
范围内：7 个代表性只读能力……
```

PRD 没有列出这 7 个能力的正式 FR，Epics 自行新增 FR-19～FR-25。

**NEW**

在 F-06 后新增：

```markdown
### F-07：代表性只读能力

- FR-19：读取当前 FOBrain 用户的安全身份、部门和角色上下文，不展示 token、raw account 或凭据。
- FR-20：读取并安全展示当前用户的 FOBrain 权限和数据范围；无权限使用明确状态，不展示 raw policy payload。
- FR-21：按 IP 查询资产列表及安全在线状态。
- FR-22：从资产列表进入单个资产的安全详情；safe ref 不冒充 provider 真实 asset ID。
- FR-23：按 IP 查询关联漏洞列表并支持资产—漏洞联查。
- FR-24：从漏洞列表进入单个漏洞的安全详情，为人工处置提供事实。
- FR-25：查询当前 actor 获授权的业务系统列表。
```

同时把 PRD 的 4 条跨功能边界正式编号为 NFR-01～NFR-04，并保持现有语义；Architecture/UX 派生约束继续以 AR/AD/UX ID 追踪，不伪装成新的业务 FR。

**影响：** 不扩展 MVP；只是把 Owner 已批准且 Epics 已使用的范围放回规范需求源。

### 4.6 为所有工作项增加可解析的 requirements metadata

**Artifact:** `epics.md`

**OLD**

大多数 Story 只有 user story 和 Acceptance Criteria，需求来源需要从段落推断。

**NEW：正式 Story 模板**

```yaml
Requirements:
  FR: [FR-01]
  NFR: [NFR-01, NFR-02]
  AR: [AR-01]
  Architecture: [AD-01]
  UX: [UX-DR-01, ST-01]
  Gates: [G-READ-01]
  Milestone: [READ-01]
```

**NEW：Target Candidate 模板**

```yaml
TargetRequirements:
  FR: [FR-07]
  Architecture: [AD-20, AD-23]
  UX: [ST-12, ST-14, ST-27]
  BlockingGates: [G-WRITE-01, G-SAFE-01]
  BlockingDecisions: [OQ-02]
  Conversion: provider-policy-verifier / product-ux / authorized-e2e
```

所有正式 Story 都必须填充该 metadata；空维度使用 `[]`，不删除字段。Requirements Coverage Map 和 `requirements-traceability-matrix.md` 按新 ID 同步。

**理由：** 让 Create Story、readiness 和后续审查基于显式来源，不靠 AI 猜测文档关系。

### 4.7 修正 Epic 依赖声明

**Artifact:** `epics.md` Epic 2/3 implementation notes

**OLD**

```markdown
Epic 3 直接依赖 Epic 1，不依赖 Epic 2；OQ-05 关闭前不能声明 Epic 3 完成。
```

但 OQ-05 的关闭 Story 位于 Epic 2。

**NEW**

```markdown
Epic 2 与 Epic 3 均直接依赖 Epic 1 的共用 Action 控制面、verifier、结果投影和 Story 1.36 操作记录；Epic 2 与 Epic 3 彼此不依赖。
```

**理由：** 消除隐式前向依赖，允许两个写域在各自门禁通过后独立推进。

### 4.8 同步 Architecture、UX 与 Runtime SPEC

**Architecture edits**

- AD-21：加入 `reconciling -> 待核验`、terminal failed/manual_attention -> 待排查的唯一映射。
- AD-23：verification policy 明确 settle deadline、预算耗尽和 mismatch 收敛规则。
- 新增 AD-26：操作记录只从 Product Facts/ActionResult 安全投影，actor/workspace 访问时重验，历史只读。
- Source Map 增加 EXPERIENCE/DESIGN，并关联 RD-03、ST-19、ST-27、Story 1.36。

**UX edits**

- IA：左侧增加系统固定“操作记录”；打开后仍是 IA-01 历史模式，详情进入 IA-02。
- ST-19：改为回读不一致的时限分流；ST-27 保持同一执行记录继续核验；ST-28 支持跨会话权威恢复。
- OQ-05：标记 CLOSED，引用 RD-03/AD-26。
- Components：定义操作记录行、筛选、加载/空/无权限/损坏引用状态；不提供重试写入。
- DESIGN/mockup：补固定导航入口、历史筛选和只读明细视觉；保持首版浅色桌面端。

**SPEC edits**

- `state-machines.md`：定义 reconcile deadline 前后的状态迁移。
- `acceptance-matrix.md` AC-13：未知项先投影待核验，收敛后才投影待排查。
- fixtures：至少覆盖 deadline 内恢复成功、明确失败、deadline 后 mismatch、失租约未知、部分成功和重启恢复。

### 4.9 同步辅助文档与门禁

批准后同步：

1. `product-blueprint/requirements-traceability-matrix.md`：覆盖 FR-01～25、NFR-01～04、Story 新 ID。
2. `product-blueprint/implementation-readiness-gate.md`：记录 Candidate 非 Sprint 规则、OQ-05 CLOSED 证据和仍阻塞的真实写门禁。
3. `docs/07-implementation-plan.md`、`docs/08-acceptance-plan.md`：采用新 Story 顺序和状态验收口径。
4. 受影响 schema/fixture/ADR 索引：只更新计划与契约，不提前实现真实 provider mutation。
5. 保留 2026-07-15 readiness report，完成纠偏后生成新报告，不改写历史结论。
6. 当前不存在 `sprint-status.yaml`，无需迁移；纠偏和 readiness 通过后首次生成。

## 5. Implementation Handoff

### 5.1 Scope Classification

**Major — Fundamental replan required.** 原因是 backlog 需要重切并同时修订 PRD、Architecture、UX 和 SPEC；业务目标与 MVP 不变。

### 5.2 执行顺序与责任

| 顺序 | 责任角色 | 交付物 | 完成门禁 |
| --- | --- | --- | --- |
| 1 | Product Manager | PRD FR-17、FR-19～25、NFR 编号；确认 RD-03/OQ-05 CLOSED | PRD 内无悬空只读范围或状态冲突 |
| 2 | Solution Architect | AD-21/23/26、状态映射、操作历史、Source Map；同步 Runtime SPEC | Architecture、SPEC、PRD 对状态和历史入口一致 |
| 3 | UX Designer | EXPERIENCE/DESIGN 的 IA、ST-19/27/28、历史状态和 mockup | OQ-05 CLOSED；两表面架构不增加第三表面 |
| 4 | Product Manager / PO | 重切 Epic 1、移动 Story 2.6、隔离 Candidate、补齐 metadata 与依赖 | Sprint parser 只识别正式 Story；无巨型 Candidate |
| 5 | Technical Writer / Developer | 同步 traceability、readiness gate、implementation/acceptance plan | 所有引用和 ID 一致，`git diff --check` 通过 |
| 6 | Readiness Reviewer | 重新运行 Implementation Readiness | 结论为 READY，或形成新的明确阻塞清单 |
| 7 | PO / Developer | 仅在 READY 后运行 Sprint Planning、Create Story，再进入开发 | Candidate 不在 Sprint；真实写域仍受各自 gate 阻塞 |

本项目由当前 Codex/BMAD 工作流顺序执行上述文档修订；无需用户逐项确认。只有本提案本身和任何新增业务范围/安全边界变化需要用户批准。

### 5.3 Success Criteria

1. Sprint Planning 输出不包含任何 `TC-*`。
2. PRD 正式需求为 FR-01～25，并与追踪矩阵、Epics 完全一致。
3. 所有正式 Story 都有完整 requirements metadata。
4. Epic 1 的基础 Story 均为单职责、可单独验收，不跨越不必要的层。
5. OQ-05 被 RD-03/AD-26 客观关闭，Epic 2/3 不再存在隐式依赖。
6. PRD、UX、Architecture、SPEC 对“待核验/待排查”的映射完全一致。
7. 新的 Implementation Readiness 结论为 READY，才允许 Sprint Planning。
8. 所有真实写动作仍只在对应门禁和环境授权通过后转换为正式 Story。

## 6. Checklist Record

### Section 1 — Trigger and Context

- [x] 1.1 触发项：Implementation Readiness 审查；关联 Story 1.3、1.6～1.10、1.17、2.6 和 6 个 Candidate。
- [x] 1.2 问题类型：原需求的规划表达和操作化不完整。
- [x] 1.3 证据：8 项审查发现、Sprint parser 规则和跨文档原文已核对。

### Section 2 — Epic Impact

- [x] 2.1 三个 Epic 目标均可完成。
- [x] 2.2 修改 Story 结构，不新增/删除 Epic。
- [x] 2.3 Epic 2/3 共同依赖移入 Epic 1。
- [x] 2.4 无 Epic 失效，无新 Epic 必需。
- [x] 2.5 Epic 顺序和业务优先级不变。

### Section 3 — Artifact Impact

- [x] 3.1 PRD：正式补 FR-19～25、NFR ID、FR-17 状态边界；MVP 不变。
- [x] 3.2 Architecture：AD-21/23/26、历史投影、状态机和来源同步。
- [x] 3.3 UX：IA、ST-19/27/28、OQ-05、组件和 mockup 同步。
- [x] 3.4 其他：Runtime SPEC、traceability、gates、implementation/acceptance plan；无代码/部署回滚。

### Section 4 — Path Forward

- [x] 4.1 Direct Adjustment：可行；工作量中高；风险低至中。
- [N/A] 4.2 Rollback：当前没有已实施工作可回滚。
- [x] 4.3 MVP Review：MVP 可实现，不需缩减或重定义。
- [x] 4.4 推荐 Direct Adjustment + baseline normalization。

### Section 5 — Proposal and Handoff

- [x] 5.1 问题摘要已记录。
- [x] 5.2 Epic 和 artifact 影响已记录。
- [x] 5.3 路径、替代方案和理由已记录。
- [x] 5.4 MVP 影响、顺序和依赖已记录。
- [x] 5.5 Major handoff 角色、责任和成功标准已记录。

### Section 6 — Final Review

- [x] 6.1 适用检查项已完成；无未记录的 Action-needed。
- [x] 6.2 提案已做一致性和可执行性复核。
- [x] 6.3 用户于 2026-07-15 明确批准完整提案。
- [x] 6.4 首次 Sprint Planning 已生成 `implementation-artifacts/sprint-status.yaml`，且只含 47 个正式 Story。
- [x] 6.5 已按 Major 变更路由至 Product Manager、Solution Architect、UX Designer、Product Owner/Developer 和 Readiness Reviewer；由当前 Codex/BMAD 工作流顺序执行。

## 7. Approval

当前状态：`PLANNING_CORRECTION_COMPLETED_SPRINT_PLANNED`。

批准记录：Vick 于 2026-07-15 明确回复“同意”。

按第 5.2 节执行文档修订。任何真实写域能力仍需独立门禁和环境授权，本提案批准不等于真实写入授权。

## 8. Handoff Record

- Change scope：Major。
- Product Manager：修订 PRD 与需求编号。
- Solution Architect：修订状态映射、操作历史决策与 Runtime SPEC。
- UX Designer：修订信息架构、状态模式和历史入口。
- Product Owner / Developer：重切 Epic、隔离 Candidate、补 requirements metadata。
- Readiness Reviewer：修订完成后重新运行 Implementation Readiness。
- Sprint status：Implementation Readiness=READY 后已生成；3 个 Epic、47 个正式 Story 均为 backlog，3 个 retrospective 为 optional。

## 9. Execution Record

- 2026-07-15：PRD、Architecture Spine、Runtime SPEC、UX EXPERIENCE/DESIGN/mockup 已完成纠偏。
- 2026-07-15：Epic 1 已重切为 Story 1.1～1.37；Epic 2/3 各保留 Story 1～5；6 个真实写入目标改为 `TC-*`，21 个旧详细拆解改为 `SB-*`，均不会被 Sprint parser 当作 Story。
- 2026-07-15：47 个正式 Story 均具备完整 Requirements metadata；FR-01～25、规范 NFR-01～04、UX-DR-01～30 覆盖检查通过，无同 Epic 前向依赖；Architecture/UX 派生约束继续使用 AR/AD/UX ID。
- 2026-07-15：追踪矩阵、implementation readiness gate、实施计划、验收计划、run lifecycle、迁移 ADR 与 FOBrain tool matrix 已同步。
- 2026-07-15：Architecture reviewer gate 最终为 Reality PASS、Good-spine PASS（High=0、Medium=0）、Incompatible Units IMPLEMENTATION_SAFE。
- 2026-07-15：新的 Implementation Readiness 报告结论为 READY；允许运行 Sprint Planning，但不解除任何真实写域门禁。
- 2026-07-15：Sprint Planning 已生成 3 个 Epic、47 个正式 Story 与 3 个 retrospective；6 个 TC 与 21 个 SB 未进入 sprint-status，真实写域仍保持关闭。
- 下一步：按 M-0→M-5 从 Story 1.1 运行 Create Story；开发前必须执行 `docs/pre-development-validation.md`。
