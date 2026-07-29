# 产品需求蓝图重建流程设计

## 目标

按照 BMAD Analysis 与 Planning 流程，把当前项目重新视为一个尚未完成需求定义的新产品。现有 `docs/` 仅作为历史输入，不直接成为新产品需求；最终形成经过证据、用户确认和完整性门禁的产品需求蓝图。

## 范围

本流程只回答“为什么做、为谁做、解决什么问题、提供哪些能力、成功如何验证”。

本流程不产出或修订：

- 技术架构；
- 数据库、接口和内部类型设计；
- 实施计划和代码任务；
- 当前代码完成度判断；
- 旧项目迁移或兼容方案。

技术研究只用于识别产品可行性约束，不选择实现方案。

## 历史资料处理

执行前对当前 `docs/` 做完整快照，包含已提交和未提交内容。快照、SHA-256 清单、Git 基线和工作区状态共同构成历史证据基线。

历史材料先保留以下五种**库存标签**；它们不是完整证据状态，所有下游材料必须按
[证据状态模型](research/evidence-status-model.md)派生六个 canonical 字段：

1. `CONFIRMED`：只用于 `U-*` 中产品负责人已经明确记录的决定；外部资料和真实用户观察不得使用本标签冒充 Owner 决策。
2. `SUPPORTED`：历史材料直接陈述且没有发现冲突，但目标适用性仍为 `UNKNOWN`，产品决策仍为 `HYPOTHESIS`。
3. `ASSUMPTION`：单一来源、推断、愿望或方案性陈述。
4. `CONFLICT`：两个来源对同一产品问题给出不同答案。
5. `MISSING`：形成完整产品需求所必需、但现有资料没有回答的问题。

权威来源、独立研究、厂商来源和真实用户观察直接使用六字段记录，不压缩成上述库存标签。代码、schema、测试和现有界面只能证明“当前存在什么”，不能单独证明“新产品应该需要什么”。

## BMAD 执行顺序

1. 来源存档、证据提取和统一证据状态。
2. 用 `bmad-forge-idea` 形成**待验证**的问题与定位假设，并由红队找出最危险假设；此时不得定稿。
3. 用 `bmad-brainstorming` 发散候选用户任务、替代方案和可能结果；此时不得生成 Must 功能。
4. 用公开市场、领域和本地可行性研究建立问题空间、厂商声明与停止条件；这仍不是目标用户证据。
5. 依次通过 DG-01 发现边界、DG-02 证据访问与研究数据使用批准、DG-03 研究队列；DG-02 未批准时不得接触真实用户／数据。
6. DG-03 后先由 Owner 在采集前同时预注册 A/B 两套准入规则但不选路径；再在正式 `bmad-market-research`、`bmad-domain-research` 和约束型 `bmad-technical-research` 中观察近期真实任务，冻结样本后通过 DG-04 只选择一个任务和一种准入路径：高频效率型或低频高损失型。
7. 通过 DG-05 预注册单一任务的核心结果、基线、目标、保护指标和停止规则；规则固定后才允许开始实验。
8. 对同一任务按同一规则比较 FOBrain 原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 改进和不做；由实验通过 DG-06 选择形态。
9. 在 DG-04 后同步收集用户采用、内部责任或经济买方证据；结合胜出形态通过 DG-07 决定产品关系与性质，只形成 Forge／定位草案。
10. 通过 DG-08 动作边界和 DG-09 适用边界后，才完成 final Forge。
11. 用 `bmad-product-brief` 形成产品战略边界和成功定义；只有 DG-01～DG-09 全部通过才可 final。
12. 用 `bmad-prfaq` 从用户结果和发布承诺反向挑战 Product Brief；未通过则回到对应上游门禁。
13. 用 `bmad-prd` 形成完整、可验证、无实现细节的产品需求并验证。
14. 只从 final PRD 投影产品能力地图和符合契约的 Product Issue。
15. 建立目标、需求、能力、Issue、验收标准之间的双向追踪矩阵。
16. 执行需求蓝图完整性门禁；只有 `READY` 才结束需求阶段。

以上顺序的关键约束是：**先确认发现边界，再取得真实任务证据；先选择单一任务并固定胜负规则，再比较解法形态；先取得形态与采用证据，再决定产品关系。**

## 澄清规则

- 可从本地材料或权威来源确定的内容由 AI 主动查证，不询问用户。
- 低影响未知项保留在假设与开放问题清单，不打断主流程。
- 会改变目标用户、核心问题、产品边界、成功标准或合规边界的问题必须由用户决定。
- 每次只提出一个高影响问题，并附当前证据、推荐答案和不同选择的影响。
- 未获得用户确认的内容不得升级为 `CONFIRMED`。

## 最终产物

```text
_bmad-output/planning-artifacts/product-blueprint/
├── workflow-design.md
├── workflow-plan.md
├── source-archive/
├── source-evidence-ledger.md
├── forge/
├── brainstorming/
├── research/
├── product-brief/
├── prfaq/
├── prd/
├── capability-map.md
├── product-issues/
├── requirements-traceability.md
├── assumptions-and-open-questions.md
└── blueprint-readiness-report.md
```

## Product Issue 契约

每个 Product Issue 必须包含：

- Product Issue 标题和稳定 ID；
- Background：当前业务问题、受影响者、解决原因、不解决的影响；
- Goal：业务结果、用户获得的能力、成功后的变化，不描述实现方式；
- 一条或多条 Given / When / Then 验收标准；
- 上游目标、PRD 需求和产品能力引用；
- 明确的非目标、依赖和未决项。

每个 Goal 至少映射一条可客观验证的验收标准。任何无法客观验证的 Issue 不得进入最终蓝图。

## 完成门禁

只有同时满足以下条件，需求蓝图才可标记为完成：

- 核心用户、问题、价值、范围和非目标已确认；
- 关键市场与领域结论有可追溯来源；
- Product Brief 与 PRFAQ 无未解决的方向性冲突；
- PRD 中所有功能需求和非功能产品约束均可验证；
- 每个业务目标映射到至少一个产品能力和 Product Issue；
- 每个 Product Issue 符合 Background / Goal / Given-When-Then 契约；
- 追踪矩阵不存在孤立目标、孤立需求或孤立 Issue；
- 所有高影响开放问题已由用户确认；
- readiness 报告结论为 `READY`。
