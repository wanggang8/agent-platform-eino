# 产品需求蓝图重建执行计划

> **执行方式：** 当前 Codex 任务持续执行；每个 BMAD workflow 单独加载并以落盘产物传递上下文。

**Goal:** 将历史文档降级为参考证据，按照 BMAD 完成从问题定义到完整 Product Issue 蓝图的全部需求工作。

**Approach:** 采用“证据基线 → BMAD Analysis → Product Brief/PRFAQ → PRD → Product Issue 投影 → 追踪与 readiness”的单向流程。任何下游产物只能引用已确认的上游事实，并保留假设、冲突和来源。

**Tools:** BMAD Method 6.10.0、Codex project skills、Web research、Markdown、SHA-256 manifest。

## 全局约束

- 不修改当前 `docs/`、现有架构文档或代码。
- 不从实现细节反推产品需求。
- 不把旧项目内部结构复制成新产品需求。
- 所有调研事实附来源；推断与假设必须显式标记。
- 所有用户确认写入 BMAD memlog。
- 一个 BMAD workflow 完成并通过检查后才能进入下一个。
- 不执行 `bmad-quick-dev`、`bmad-dev-auto`、架构或开发 workflow。

---

### Task 1：建立历史证据基线

**Create:**

- `source-archive/docs-pre-bmad-2026-07-11.tar.gz`
- `source-archive/docs-pre-bmad-2026-07-11.sha256`
- `source-archive/archive-metadata.md`
- `source-evidence-ledger.md`

**Steps:**

- [x] 对当前 `docs/` 创建包含未提交内容的完整压缩快照。
- [x] 对快照和全部源文件生成 SHA-256 清单。
- [x] 记录 HEAD、分支、工作区状态、文件数量和快照时间。
- [x] 按 `CONFIRMED/SUPPORTED/ASSUMPTION/CONFLICT/MISSING` 提取产品相关陈述。
- [x] 检查台账中每条陈述均有来源路径，推断不得伪装为事实。

### Task 2：重新锻造产品方向

**Create:**

- `forge/.memlog.md`
- `forge/forged-idea.md`
- `forge/forge-report.html`

**Steps:**

- [x] 使用 `bmad-forge-idea` 读取证据台账中的相关摘录，形成待验证定位与用户问题假设。
- [x] 对关键假设进行攻击与防守测试，确认当前定位只能是 provisional。
- [x] 将“先定产品关系”修正为 DG-01～DG-07 的证据优先顺序：DG-06 只选解法形态，DG-07 才结合采用／买方证据决定产品关系。
- [ ] 从 DG-03／DG-04 的真实任务证据明确任务执行者候选、核心问题和现状成本；DG-03 研究队列不得直接写成最终主用户。
- [ ] 在 DG-05 预注册核心结果，DG-06 只确定胜出形态，结合最终主用户与采用／买方证据在 DG-07 决定产品价值和关系；动作与适用边界分别留到 DG-08／DG-09。
- [ ] 将无法从证据确定的方向性问题逐个提交用户确认。
- [ ] 仅在方向被硬化后输出 `forged-idea.md`。

**Re-entry gate:** 本任务的假设与红队阶段已完成；DG-07 后可以补全定位草案，但 `forged-idea.md` 只有在 DG-08 动作边界、DG-09 适用边界及证据治理复审也通过后才能标记为 final。该门禁不阻止先执行 Task 3／Task 4。

### Task 3：问题空间与能力候选发散

**Create:**

- `brainstorming/.memlog.md`
- `brainstorming/brainstorm-intent.md`
- `brainstorming/brainstorm-report.html`

**Steps:**

- [x] 使用 `bmad-brainstorming` 围绕用户任务、失败场景和替代方案发散。
- [x] 将候选能力与用户问题连接，不讨论数据库、接口或框架。
- [x] 聚类、去重并标记证据强度。
- [x] 只保留进入研究验证的候选方向，不直接承诺范围。

### Task 4：完成需求前置研究

**Create:**

- `research/market-research.md`
- `research/domain-research.md`
- `research/technical-feasibility-constraints.md`
- `research/research-synthesis.md`

**Steps:**

- [ ] 使用 `bmad-market-research` 验证用户、痛点、竞品和替代方案。
- [ ] 使用 `bmad-domain-research` 提取领域角色、业务流程、规则、术语和风险。
- [ ] 使用 `bmad-technical-research` 只提取影响产品承诺的可行性和外部依赖。
- [x] 形成公开市场、领域、本地可行性和研究综合草案；保持正式 BMAD `[C]` 未完成状态。
- [x] 为外部事实保存直接来源、访问日期，并统一 `VENDOR_CLAIMED`、`DEPLOYMENT_AVAILABLE`、`INTEGRATION_AVAILABLE`、`TASK_EFFECTIVE` 状态。
- [x] 综合研究草案，更新证据台账中的支持、反证、停止条件和缺口。
- [x] 通过 V5 独立证据治理复审：V4 两个 LOW 条件均关闭，无新增残留，最终判定 `PASS`。此项不代表正式 BMAD Research 完成。
- [x] 通过 DG-01：Owner 于 2026-07-11 确认 FOBrain 仅为第一个可验证业务样本，不是最终产品边界。
- [x] 在正式 BMAD Market／Domain Research 工作流中取得范围 `[C]`（2026-07-11）。
- [ ] 通过 DG-02～DG-03 后，先在任务采集前冻结 A/B 两套准入规则但不选路径，再用合规取得的真实用户任务证据完成正式 BMAD Research；样本冻结后由 DG-04 只选一路和一个任务。

### Task 5：形成 Product Brief 并通过 PRFAQ 挑战

**Entry gate:** DG-01～DG-09、final Forge 和研究证据治理门禁全部通过；否则只能保留骨架，不得把 Brief 或 PRFAQ 标记为 final。

**Create:**

- `product-brief/.memlog.md`
- `product-brief/brief.md`
- `product-brief/addendum.md`
- `prfaq/prfaq.md`
- `prfaq/verdict.md`

**Steps:**

- [ ] 使用 `bmad-product-brief` 创建产品简介，覆盖问题、方案、用户、差异、成功、范围和愿景。
- [ ] 所有未确认内容使用明确的假设标记。
- [ ] 使用 `bmad-prfaq` 从未来用户结果和常见质疑反向检查产品方向。
- [ ] 解决 Product Brief 与 PRFAQ 的方向性冲突。
- [ ] 只有 final Product Brief 已满足 DG-01～DG-09，且 PRFAQ verdict 明确允许继续时，才能进入 PRD。

### Task 6：创建并验证完整 PRD

**Create:**

- `prd/.memlog.md`
- `prd/prd.md`
- `prd/addendum.md`
- `prd/validation-report.md`
- `prd/validation-report.html`

**Steps:**

- [ ] 使用 `bmad-prd` Create 模式，输入 Brief、PRFAQ、研究和证据台账摘录。
- [ ] 明确用户、JTBD、用户旅程、范围、非目标、成功指标、功能需求和产品级非功能约束。
- [ ] 为每个需求定义可测试结果，不写技术实现。
- [ ] 使用 `bmad-prd` Validate 模式检查完整性、一致性、可测试性和范围边界。
- [ ] 修复所有阻断项并将 PRD 状态置为 final。

### Task 7：生成能力地图与 Product Issue 清单

**Create:**

- `capability-map.md`
- `product-issues/PI-*.md`
- `assumptions-and-open-questions.md`

**Steps:**

- [ ] 从 final PRD 提取用户能力，不从代码或架构补能力。
- [ ] 按用户价值和业务目标组织能力层级。
- [ ] 为每个可独立验收的业务增量创建稳定 Product Issue。
- [ ] 每个 Issue 完整填写 Background、Goal、Given/When/Then、非目标、依赖和上游引用。
- [ ] 将低影响未决项集中登记；高影响未决项必须逐个获得用户确认。

### Task 8：建立追踪矩阵并执行最终门禁

**Create:**

- `requirements-traceability.md`
- `blueprint-readiness-report.md`

**Steps:**

- [ ] 建立业务目标 → PRD 需求 → 产品能力 → Product Issue → AC 的双向追踪。
- [ ] 检查孤立目标、孤立需求、重复 Issue、不可验证 AC 和实现泄漏。
- [ ] 检查所有外部结论均有来源，所有假设均被标记。
- [ ] 对照 `workflow-design.md` 的完成门禁生成 READY/NEEDS WORK/NOT READY 结论。
- [ ] 如果不是 READY，回到对应上游产物修正并重新检查；只有 READY 才交付最终需求蓝图。

## 自检

- 计划覆盖存档、证据提取、BMAD Analysis、Brief、PRFAQ、PRD、Issue 和 readiness。
- 输出路径职责单一，没有修改现有 `docs/` 或代码的步骤。
- 每个阶段都有明确输入、输出和进入下一阶段的条件。
- 最终完成定义与 `workflow-design.md` 一致。
