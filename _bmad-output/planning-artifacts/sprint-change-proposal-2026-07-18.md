---
title: Agent Platform Eino 纵向实施切片纠偏提案
status: APPROVED_HANDOFF_READY
date: '2026-07-18'
mode: Batch
changeScope: Moderate
approvedBy: Vick
approvedAt: '2026-07-18'
trigger:
  - _bmad-output/planning-artifacts/implementation-readiness-report-2026-07-16.md
supersedesPlanningClaimsFrom:
  - _bmad-output/planning-artifacts/sprint-change-proposal-2026-07-15.md
---

# Sprint Change Proposal

> 本提案只调整实施顺序、Story 边界、门禁状态和辅助计划，不改变已批准的 PRD、UX 核心体验或 Architecture 不变量。提案获批、文档修订并重新通过 Implementation Readiness 前，不进入 M-1 开发；真实 FOBrain 写入继续关闭。

## 1. Issue Summary

### 1.1 触发原因

2026-07-18 完成的 Implementation Readiness 复核确认：FR-01～FR-25 已 100% 映射，PRD、UX 与 Architecture 主线一致，但当前 Epics 不能稳定转化为可独立、可按序执行的开发 Story。

触发项不是某个代码失败 Story，而是以下两个规划缺陷：

1. Epic 1 在一个 Epic 中连续安排 37 个 Story，Story 1.1～1.24 主要按工具链、schema、数据库、repository、SSE、Eino、LLM 和 Action 控制面横向分层，用户可见价值直到后段才形成。
2. Story 1.31 已要求搜索、筛选、同名消歧和键盘选择，并把 OQ-02 作为 Gate；后续 Epic 2 Story 2.1 才负责定义并关闭 OQ-02，同时又依赖 Epic 1 的人员列表，形成跨 Epic 循环。

问题类型：**Failed planning approach requiring a different implementation slicing strategy**。产品目标没有误解，失败的是把正确蓝图操作化为开发顺序的方法。

### 1.2 客观证据

- `implementation-readiness-report-2026-07-16.md` 结论为 `NOT READY`，记录 CQ-01 与 CQ-02 两个 Critical。
- Epic 1 明确承载 M-0～M-5、37 个正式 Story；其中 1.1～1.24 是长串水平依赖。
- Story 1.31 的 `Gates: [G-READ-02, OQ-02]` 和 AC 已包含人员选择规则；Story 2.1 AC-2.1.8 才允许关闭 OQ-02。
- Epic 1 的正式 Story 通常只有一个合并 GWT 场景，详细负向 AC 位于明确声明“不得被开发 Agent 当作任务”的 Archived Source Bundles。
- `docs/07-implementation-plan.md` 与 `docs/08-acceptance-plan.md` 仍保留旧 Phase 0～9、移动端验收、24 工具恢复和工单写域；这些与当前“浅色桌面、7 项代表性只读、四类受控动作实验”的正式蓝图冲突。
- `implementation-readiness-gate.md` 仍把 G-TOOLCHAIN 标为 BLOCKED；权威验收记录和 `docs/07`、`docs/08` 已裁决 PASS。
- `sprint-status.yaml` 仍把 Story 1.1 标为 `in-progress`；权威验收记录已声明 Story 1.1 complete。

## 2. Impact Analysis

### 2.1 Epic Impact

| Epic | 用户目标 | 当前问题 | 调整 |
|---|---|---|---|
| Epic 1：可信、可恢复的 FOBrain 事实工作台 | 保留 | 37 个水平 Story；包含与只读价值无关的完整写控制面；人员选择前向依赖 | 重写为 12 个纵向 Story；保留查询事实链，把人员选择和 Action 控制面移到首次消费它们的 Epic 2 |
| Epic 2：运营人工派发与责任人转发 | 保留 | OQ-02 关闭与 Epic 1 循环；共用控制面被提前放在 Epic 1；操作记录归属不自然 | 重写为 10 个递进 Story；先交付人员读取/选择，再以 mock 派发纵向建立共用 Action 控制面，最后交付转发、操作记录与 Gate Evidence |
| Epic 3：管理员延时与运营误报处置 | 保留 | 为避免依赖 Epic 2，被迫要求 Epic 1 提前建设所有写控制面 | 调整为正常依赖 Epic 2 已交付的共用 Action 控制面和操作记录；自身 Story 保持延时/误报专属规则与证据 |

不新增、不删除产品 Epic，不改变四类动作和只读能力范围。Epic 依赖改为正常顺序：`Epic 1 → Epic 2 → Epic 3`。

### 2.2 Story Impact

- Story 1.1 保留原 ID 并标记 `done`。
- 旧 Story 1.2～1.37 不再作为后续稳定 ID；新 Epic 1 使用 1.2～1.12。
- 旧 Story 1.31 的“人员读取”和“人员选择”拆到 Epic 2 的 2.1、2.2，彻底解除 OQ-02 循环。
- 完整 ActionDraft／Confirm／Attempt／Lease／Verifier 不再作为只读 Epic 的水平前置；它们通过 mock 派发的可见旅程在 Epic 2 逐步建立。
- Story 1.36 操作记录移到 Epic 2，在共用 ActionResult 已产生后实施；Epic 3 正常复用 Epic 2，不再人为保持“彼此不依赖”。
- Gate/Evidence 工作项继续可排期，但增加 `WorkItemType: gate-evidence`，与用户价值 Story 分开统计。
- Target Candidate 保持 `TC-*`、不可进入 Sprint、不可原地改名。
- Archived Source Bundles 移出可执行 `epics.md`；有效 AC 必须进入对应正式 Story，自身只作为历史来源保存。

### 2.3 Artifact Impact

| Artifact | 是否改变产品语义 | 必须调整 |
|---|---:|---|
| PRD | 否 | FR/NFR、范围、目标与非目标保持不变；只更新 Story traceability 引用 |
| Architecture Spine | 否 | AD-01～AD-27 保持；只补“架构能力按纵向切片落地，门禁只在全部所需证据完成后 PASS”的实施映射，并更新 Epic 3 依赖说明 |
| UX DESIGN / EXPERIENCE | 否 | 两表面、浅色桌面、聊天唯一动作入口、只读 IA-02 均不变；OQ-02 继续由新 Story 2.2 关闭 |
| Runtime SPEC | 否 | CAP-1～CAP-11 和状态不变量不变；更新 acceptance-to-story 映射，明确 query slice 与 action slice 的落地顺序 |
| Epics | 是，实施结构 | 重写 Story inventory、完整 AC、依赖、WorkItemType、Gate 和 Target Candidate 转换映射 |
| Traceability Matrix | 否 | 全量替换 Story ID；保持 FR 25/25、NFR 4/4 |
| Readiness Gate | 否 | G-TOOLCHAIN 改为 PASS；更新正式 backlog 结构、当前允许范围与新 Story ID |
| `docs/07-implementation-plan.md` | 是，执行计划 | 用 M-0～M-6 纵向 Story 计划替换旧 Phase 0～9；删除移动端、24 工具全量恢复、工单写域和旧兼容描述 |
| `docs/08-acceptance-plan.md` | 是，验收计划 | 按新 Story/Gate 重写验收矩阵；只验浅色桌面、7 项代表性只读和本轮四类动作的 mock/evidence/target gate |
| `docs/pre-development-validation.md` | 是，门禁入口 | 从旧 Phase 节点改为新 M 节点；旧项目继续只读参考，不再要求 24 工具、移动端和工单恢复 |
| `sprint-status.yaml` | 是，backlog 状态 | Story 1.1=`done`；删除旧 1.2～3.5 条目并按批准后的新 inventory 重建；TC/SB 不进入 |
| 历史 readiness / change proposal | 否 | 原样保留，不改写历史结论；生成新的 readiness 报告 |

### 2.4 Technical and Delivery Impact

- 不回滚任何 Story 1.1 工具链实现、数据库或产品代码。
- M-1 尚未开始，没有已实施的 Product Facts v2 代码需要丢弃。
- 新顺序要求第一个开发切片是 fixture-driven walking skeleton；它以一个安全新增漏洞 fixture 贯穿 schema、SQLite、Product Facts、registry、mock Eino tool loop、SSE 和桌面 Workbench，证明端到端路径后再接真实 FOBrain 能力。
- G-ARCH-V2 可以在多个纵向 Story 中累积证据，但不得在 query/action 所需契约未全部完成前标记 PASS。
- 所有真实 mutation 仍受 G-FACT、适用 G-WRITE、G-SAFE 和环境授权阻塞。

## 3. Path Forward Evaluation

| 选项 | 可行性 | 工作量 | 风险 | 结论 |
|---|---|---|---|---|
| Option 1：直接调整现有计划 | 可行 | 中 | 低至中 | **推荐**。范围和架构不变，只重组 backlog 与辅助计划 |
| Option 2：回滚 | 不适用 | 低 | 中 | Story 1.1 已客观完成且不造成问题；M-1 未开始，无可带来收益的回滚 |
| Option 3：缩减／重定义 MVP | 无必要 | 中高 | 高 | 25 条 FR 已完整且一致；删功能不能解决 Story 切片和文档漂移 |

选择：**Option 1 — Direct Adjustment**。

范围分类：**Moderate**。这是较大的 backlog 重组，需要 Product Owner / Developer 主导并由 Architecture／UX 做一致性复核，但不构成产品或架构根本重做。

进度影响：增加一次 Epics/计划同步、Sprint status 重建和 Implementation Readiness 复跑；避免进入开发后再拆除横向底座和循环依赖。

## 4. Detailed Change Proposals

### 4.1 重写 Epic 1 为纵向只读切片

**OLD**

```text
Epic 1 = M-0～M-5 + 37 Stories
1.1～1.24：按工具链/schema/DB/repository/SSE/Eino/LLM/Action 水平建设
1.25～1.37：最后接入真实查询、统一事实和 Workbench
```

**NEW**

| ID | 类型 | 用户可验证结果 | 主要技术落点 |
|---|---|---|---|
| 1.1 | enabler / done | 使用唯一可复现工具链产生可信证据 | G-TOOLCHAIN |
| 1.2 | user-value | 在浅色桌面 Workbench 中用固定安全 fixture 完成一次新增漏洞查询并看到摘要 | 最小 query schema、Greenfield SQLite、Product Facts、registry、mock Eino loop、projection、SSE、桌面 UI |
| 1.3 | gate-evidence | 证明精确新增查询和 7 项代表性能力所需字段的 resolved/empty/unavailable 状态 | G-READ-01/03、G-FACT 数据可得性；零 mutation |
| 1.4 | user-value | 查看当前用户与权限三态，明确事实属于谁及可见范围 | current_user_context、my_permissions、安全投影 |
| 1.5 | user-value | 查询真实全部新增漏洞，正确区分有结果、0 条、部分分页和失败 | exact-new capability、complete_set、倒序、空态 |
| 1.6 | user-value | 按 IP 查询资产并打开安全资产详情 | asset list/detail、opaque locator |
| 1.7 | user-value | 按 IP 查询漏洞并打开安全漏洞详情 | vulnerability list/detail、opaque locator |
| 1.8 | user-value | 查询全部业务系统并补齐统一漏洞行事实 | business_list、多源字段来源/缺失语义 |
| 1.9 | user-value | 在冻结快照中按每页 100 条查看完整统一事实 | snapshot coverage、stable order、IA-02、右栏详情 |
| 1.10 | user-value | 对多次查询的“这些/第 N 次”进行确定性引用或澄清 | Result Index、resolver、context compression |
| 1.11 | user-value | 刷新、SSE 断线或服务重启后继续查看同一只读事实 | persistent FactEvent/SSE cursor、read recovery |
| 1.12 | acceptance | 完成 READ-01 端到端验收并保持所有写域禁用 | schema/Go/TS/browser/smoke/security evidence |

Story 1.2 是刻意设计的薄 walking skeleton，只允许一个版本化 fixture 和 read-only capability；不得顺带建立 ActionDraft、mutation、verifier、操作记录或真实写入口。

### 4.2 解除 Story 1.31／Story 2.1 的 OQ-02 循环

**OLD — Story 1.31**

```markdown
Gates: [G-READ-02, OQ-02]
When 用户搜索、筛选、键盘选择或遇到同名人员
Then 按已固定规则展示、消歧并只允许人工单选
```

**OLD — Story 2.1**

```markdown
Given Epic 1 已具备完整人员列表安全投影
...
Then OQ-02 才能标记为关闭
```

**NEW**

```markdown
Story 2.1：读取 FOBrain 全部人员安全列表
- 只验收完整快照、稳定人员 identity、安全显示字段、来源顺序、加载/空/失败。
- 不包含搜索、排序、同名消歧、可选状态或 ActionDraft。
- Gate: G-READ-02；OQ-02 仍 open。

Story 2.2：固定并交付人工选人规则
- 基于 Story 2.1 的冻结人员快照定义搜索、筛选、同名消歧、可选状态和键盘单选。
- 目标部署状态 allowlist 与安全消歧字段有证据后关闭 OQ-02。
- 未关闭时生产人员选择和真实写入继续隐藏。
```

理由：读取事实与使用事实作人工选择是两个不同能力；后者不能反向阻塞只读事实链。

### 4.3 在 Epic 2 的派发旅程中纵向建立共用 Action 控制面

**OLD**

```text
Epic 1 Story 1.3～1.24 提前建立完整 Action schema、Confirm、attempt、lease、verifier、reconcile。
Epic 2 只增加动作专属规则。
```

**NEW**

| ID | 类型 | 用户可验证结果 | 建立的共用能力 |
|---|---|---|---|
| 2.1 | user-value | 查看完整人员安全列表 | G-READ-02 产品 capability |
| 2.2 | decision/user-value | 人工搜索、消歧并选择唯一接收人 | OQ-02、人员选择卡 |
| 2.3 | user-value/mock | 基于冻结新增漏洞和人工接收人看到不可变派发草案 | Catalog v2、PrepareAction、ActionDraft、policy 双检、前 5 条+全部 |
| 2.4 | user-value/mock | 确认或取消派发；未确认/取消/过期均证明零 mutation | Pending、ConfirmAction、双 Continuation、digest、一次性 approval |
| 2.5 | user-value/mock | 单条 mock 派发只执行一次并以负责人回读显示完成 | attempt reservation、SemanticMutationClaim、verifier |
| 2.6 | user-value/mock | 多条 mock 派发逐条显示完成、待核验、待排查并可恢复 | lease/fencing、partial、reconcile、restart |
| 2.7 | user-value/mock | 主动转发复用同一控制面，不复制状态机 | recipient grouping、transfer policy/verifier |
| 2.8 | user-value | 跨会话从“操作记录”找到同源逐条结果 | AD-26、OperationListSnapshot、六个月/open union |
| 2.9 | gate-evidence | 验证目标部署派发权限、写入、回读、部分失败和恢复 | G-WRITE-01/02、G-SAFE 子项 |
| 2.10 | gate-evidence | 验证目标部署转发权限、写入、回读、部分失败和恢复 | G-WRITE-01/02、G-SAFE 子项 |

Story 2.3～2.7 只连接版本化 mock provider／mock verifier；正式 AC 必须明确真实 FOBrain mutation 调用数为 0。Target Candidate 只有在 Gate Evidence PASS 后转换。

### 4.4 调整 Epic 3 依赖

**OLD**

```markdown
Epic 3 直接依赖 Epic 1，不依赖 Epic 2；复用 Epic 1 的统一写控制面和 Story 1.36 操作记录。
```

**NEW**

```markdown
Epic 3 依赖 Epic 1 的统一只读事实链，并依赖 Epic 2 已交付的共用 Action 控制面、状态投影和 Story 2.8 操作记录。
Epic 3 不复制 ActionDraft、Confirm、幂等、lease、verifier、reconcile 或历史投影；只增加延时/误报专属角色、基数、目标参数和 verifier 规则。
```

Epic 3 正式 inventory 保留 5 个当前工作项：

1. 3.1 `gate-evidence`：关闭 OQ-04 并固定直接延时规则。
2. 3.2 `user-value/mock`：单条修复延时 mock 旅程。
3. 3.3 `user-value/mock`：人工误报 mock 旅程。
4. 3.4 `gate-evidence`：目标部署延时权限、双事实写入/回读/恢复。
5. 3.5 `gate-evidence`：目标部署误报对象权限、写入/回读/恢复。

TC-3.1～TC-3.3 保持非 Sprint 工作项。

### 4.5 正式 Story 必须自包含，归档材料退出执行入口

**OLD**

```markdown
正式 Story：一个合并 GWT。
Archived Source Bundle：包含详细 AC，但明确不得被开发 Agent 当作任务。
```

**NEW**

- `epics.md` 中每个正式 Story 至少包含：用户结果、WorkItemType、FR/NFR/AD/UX/Gate、前置 Story、范围、非目标、完整 GWT AC、受影响目录和验收命令。
- 当前仍有效的负向、安全、恢复 AC 直接进入正式 Story。
- Source Bundles 移到 `planning-artifacts/archive/epic-source-bundles-2026-07-15.md`，只保留历史追踪链接，不参与 Create Story。
- Create Story 生成的 Story 文件不得把 Source Bundle 当隐式验收来源。

### 4.6 修正门禁和 Sprint 状态的客观漂移

**OLD**

```text
implementation-readiness-gate.md: G-TOOLCHAIN = BLOCKED
sprint-status.yaml: Story 1.1 = in-progress
docs/07、docs/08、权威 acceptance record: G-TOOLCHAIN PASS，Story 1.1 complete
```

**NEW**

```text
G-TOOLCHAIN = PASS
Story 1.1 = done
Epic 1 = in-progress（只表示 Story 1.1 已完成；不表示允许跳过新的 readiness）
```

门禁文档引用 `docs/acceptance-records/story-1-1-g-toolchain-2026-07-15.md` 的最终 PASS 裁决；历史失败只保留在验收记录，不继续作为当前状态。

### 4.7 用当前蓝图替换旧实施/验收计划

**OLD**

- `docs/07`：Phase 0～9，包含移动端 Workbench、24 只读工具恢复、MCP 市场式扩展、写域工单审批。
- `docs/08`：P0/P1/P2、desktop/mobile 截图、24 工具 + connector 等价恢复和 live ticket write。
- `pre-development-validation.md`：要求复核旧 24 工具、移动端和工单恢复。

**NEW**

| Milestone | 内容 | 完成信号 |
|---|---|---|
| M-0 | Story 1.1 工具链 | 已 PASS |
| M-1 | Story 1.2 read-only walking skeleton | 单 fixture 端到端、Greenfield DB、桌面 UI、零 mutation |
| M-2 | Story 1.3～1.8 真实代表性读取 | G-READ-01/03 产品证据、统一事实字段来源明确 |
| M-3 | Story 1.9～1.12 快照、引用、恢复、READ-01 | G-FACT/read E2E、桌面验收 |
| M-4 | Story 2.1～2.8 人员与 mock Action 控制面 | OQ-02 closed、G-ARCH-V2、mock 零写入/逐条三态/恢复 |
| M-5 | Story 2.9～2.10、3.1～3.5 Gate Evidence | 各写域客观证据；仍不等于产品集成 |
| M-6 | 转换后的 Target Story | 仅适用门禁全部 PASS 后进入 |

明确删除当前范围外门禁：移动端、暗色、24 工具等价恢复、工单审批、通知、自动派发和旧 runtime 兼容。保留旧项目只读参考和安全基线，不把旧项目范围变成新项目 MVP。

### 4.8 更新追踪和再次就绪检查

批准并实施后：

1. `requirements-traceability-matrix.md` 使用新 Story ID，FR-01～25 与 NFR-01～04 覆盖率必须仍为 100%。
2. `implementation-readiness-gate.md` 更新 G-TOOLCHAIN、正式 backlog 和 M-1 允许范围。
3. `sprint-status.yaml` 按新 inventory 重建；TC、SB 不得出现。
4. 运行 Epics 结构检查：无跨 Epic 前向依赖；Epic 1 不引用 OQ-02；Epic 3 只依赖既有 Epic 1/2 输出。
5. 重新执行 Implementation Readiness，生成新报告；结论必须为 READY 才能 Create Story 1.2。

## 5. Implementation Handoff

### 5.1 Scope Classification

**Moderate — backlog reorganization required.**

不需要重新做 PRD、UX 或架构设计；需要 PO/Developer 重写 Epics 和计划，Architecture/UX 做一致性审阅，Readiness Reviewer 重新裁决。

### 5.2 Execution Order

| 顺序 | 责任角色 | 交付物 | 门禁 |
|---|---|---|---|
| 1 | Product Owner / Developer | 按 4.1～4.5 重写 `epics.md`，迁出 Source Bundles | 无前向依赖；正式 Story 自包含 |
| 2 | Technical Writer / Developer | 更新 traceability、readiness gate、`docs/07`、`docs/08`、pre-development validation | 当前范围、Story ID、Gate 状态一致 |
| 3 | Product Owner | 重建 sprint-status；Story 1.1 done，TC/SB 排除 | YAML 与 Epics 精确一致 |
| 4 | Architect / UX Reviewer | 确认 AD/UX 语义未变化，纵向切片没有绕过安全门禁 | 无新增架构/UX 冲突 |
| 5 | Readiness Reviewer | 重跑 Implementation Readiness | READY 才能进入 Story 1.2 |
| 6 | Developer | Create Story 1.2，执行开发前复核后开发 | 只做 fixture read-only walking skeleton |

### 5.3 Success Criteria

1. Epic 1 从 37 个水平 Story 收敛为 12 个纵向工作项，Story 1.1 保持 done。
2. Epic 1 不再引用 OQ-02；人员读取和人员选择在 Epic 2 内按 2.1→2.2 顺序完成。
3. Epic 3 正常依赖 Epic 2 共用控制面，不复制状态机。
4. 正式 Story 的负向、安全和恢复 AC 自包含，不依赖 Archived Source Bundles。
5. `docs/07`、`docs/08`、pre-development validation 不再包含移动端、24 工具等价恢复、工单写域或旧兼容目标。
6. G-TOOLCHAIN 在所有当前状态文档中为 PASS；Story 1.1 在 sprint status 中为 done。
7. FR 25/25、NFR 4/4 追踪保持完整。
8. 新 Implementation Readiness 结论为 READY 后才允许 Create Story 1.2。
9. 所有真实写动作继续由适用 G-WRITE、G-SAFE、G-FACT 和环境授权独立阻塞。

## 6. Checklist Record

### Section 1 — Trigger and Context

- [x] 1.1 触发项：2026-07-18 Implementation Readiness；核心关联 Epic 1 inventory、Story 1.31 和 Story 2.1。
- [x] 1.2 问题类型：实施切片方法失败，不是新需求或战略转向。
- [x] 1.3 证据：2 Critical、4 Major、2 Minor；Epics 原文、门禁、状态和辅助计划冲突已核对。

### Section 2 — Epic Impact

- [x] 2.1 当前 Epic 1 不能按原计划安全完成，必须重写 Story inventory。
- [x] 2.2 保留三个 Epic 目标，修改 Story 边界和 Epic 依赖。
- [x] 2.3 Epic 2 承接人员选择和共用 Action 控制面；Epic 3 正常依赖 Epic 2。
- [x] 2.4 无 Epic 失效，无新产品 Epic 必需。
- [x] 2.5 顺序明确为 Epic 1→2→3；真实写域优先级不提前。

### Section 3 — Artifact Impact

- [x] 3.1 PRD：无需求或 MVP 变更，只更新 Story 追踪。
- [x] 3.2 Architecture：无核心决策变更，只更新落地映射与依赖说明。
- [x] 3.3 UX：无表面/流程变更；OQ-02 路由到新 Story 2.2。
- [x] 3.4 其他：Epics、SPEC mapping、traceability、readiness gate、docs/07、docs/08、pre-development validation、sprint status 必须同步；无代码回滚。

### Section 4 — Path Forward

- [x] 4.1 Direct Adjustment：可行；工作量中；风险低至中。
- [N/A] 4.2 Rollback：Story 1.1 正确且 M-1 未开始，无合理回滚对象。
- [x] 4.3 MVP Review：MVP 仍可实现，无需缩减。
- [x] 4.4 选择 Direct Adjustment，范围 Moderate。

### Section 5 — Proposal and Handoff

- [x] 5.1 问题摘要与证据已记录。
- [x] 5.2 Epic、Story 和 artifact 影响已记录。
- [x] 5.3 推荐路径、替代方案、工作量和风险已记录。
- [x] 5.4 MVP 不变；M-0～M-6 新顺序已定义。
- [x] 5.5 PO/Developer、Architect/UX Reviewer、Readiness Reviewer 和 Developer handoff 已定义。

### Section 6 — Final Review

- [x] 6.1 所有适用检查项已完成；当前无未记录 Action-needed。
- [x] 6.2 提案已与 PRD、Architecture、UX、SPEC、Epics、门禁和 sprint status 交叉核对。
- [x] 6.3 Vick 于 2026-07-18 明确回复 `yes`，批准完整提案。
- [!] 6.4 交接实施时重建 `sprint-status.yaml`；本 Correct Course 工作流不提前修改 backlog 状态。
- [x] 6.5 已按 Moderate 范围路由至 Product Owner / Developer，并要求 Architecture/UX 一致性复核与 Readiness Reviewer 收口。

## 7. Approval

当前状态：`APPROVED_HANDOFF_READY`。

批准记录：Vick 于 2026-07-18 明确回复 `yes`。

本提案批准不等于真实 FOBrain 写入授权。只有完成第 5 节交接、重新获得 Implementation Readiness=`READY`，才允许 Create Story 1.2；所有真实写动作仍需独立通过适用门禁。

## 8. Handoff Record

- Change trigger：Implementation Readiness 发现 Epic 1 水平大底座和 OQ-02 跨 Epic 循环。
- Change scope：Moderate。
- Primary recipient：Product Owner / Developer。
- Review recipients：Solution Architect、UX Reviewer、Readiness Reviewer。
- Required deliverables：重写 Epics、迁出 Source Bundles、更新追踪／门禁／实施与验收计划、重建 Sprint status、重新运行 Implementation Readiness。
- Implementation boundary：不修改 PRD 业务范围，不改变 AD-01～AD-27，不启用真实 mutation。

## 9. Workflow Execution Log

- 2026-07-18：以 Batch 模式完成触发、Epic、artifact、路径和 handoff 分析。
- 2026-07-18：生成完整 Sprint Change Proposal，分类为 Moderate Direct Adjustment。
- 2026-07-18：Vick 明确批准提案。
- 2026-07-18：Correct Course 工作流完成；下一步路由到 backlog 重组，不直接进入开发。
