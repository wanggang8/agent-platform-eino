---
status: draft
discovery_boundary_confirmation: approved_2026_07_11
updated: 2026-07-11
---

# 产品需求证据验证计划

本计划用于回答“哪些功能真的需要”。历史工具数量、现有页面和 AI 能做到什么，都不能替代用户任务证据。

所有证据升级遵守[统一证据状态模型](./evidence-status-model.md)，尤其不得把 `VENDOR_CLAIMED` 自动升级为 `DEPLOYMENT_AVAILABLE`、`INTEGRATION_AVAILABLE` 或 `TASK_EFFECTIVE`。

> 表中的样本数、比例和改善幅度是本轮需求发现的**建议工作门槛**，用于防止凭单个案例定需求，不是外部行业标准，也不代表产品负责人已经确认。

## 证据责任类别

下列是 `responsibility_class`，不是 canonical evidence status：

- `OWNER_DECISION`：产品关系、商业边界、权限和成功口径，由产品负责人决定并记录。
- `USER_EVIDENCE`：任务频率、耗时、痛点、替代方式和信任要求，由目标用户观察确认。
- `ADOPTION_BUYER_EVIDENCE`：内部采用责任、标准功能持续使用、增购／续费或独立采购证据，由相应采用负责人、设计伙伴或经济买方提供；口头兴趣不构成通过。
- `FOBRAIN_PRODUCT_EVIDENCE`：目标版本原生功能、配置、培训和产品边界，由 FOBrain 产品／实施负责人及目标部署演示确认。
- `FOBRAIN_API_DATA_EVIDENCE`：数据字段、权限、时效、错误和动作能力，由 FOBrain 契约与真实样本确认。
- `LEGAL_COMPLIANCE_REVIEW`：研究数据与最终部署的适用性、授权和限制，由法务／合规／数据安全责任人确认。
- `EXPERIMENT_RESULT`：候选解法是否改善预注册结果，由同任务对照验证。
- 任何一项缺少必要证据时，只能保留为假设，不能升级为 final PRD 的 Must。

## 第一批关键假设

| ID | 假设 | 需要的最小证据 | 通过条件 | 失败后的处理 |
| --- | --- | --- | --- | --- |
| VH-01 | FOBrain 仅作为第一个可验证业务样本，不是最终产品边界；先验证其中的重要剩余任务，暂不预设 AI、产品形态、功能数量或工具／数据源范围。 | `OWNER_DECISION`：Owner 于 2026-07-11 确认 DG-01。 | **已通过**：先观察任务，再由证据比较所有解法、扩展边界与不做。 | 若 Owner 撤销或改变发现锚点，重新确认研究范围。 |
| VH-01B | 能合规取得产品基线和真实任务证据。 | `OWNER_DECISION + FOBRAIN_PRODUCT_EVIDENCE + FOBRAIN_API_DATA_EVIDENCE + USER_EVIDENCE + LEGAL_COMPLIANCE_REVIEW`：分别指定产品／部署、数据、用户研究、数据安全与法务／合规（按适用性）责任人；在访问前签认目的、最小数据、授权／同意、脱敏、访问者、地域、保留删除和模型／遥测禁入项。 | 至少一个可观察队列，且对应责任人全部签认研究数据批准；业务或 API 人员不能单独关闭本项。 | 公开资料只维持草案。 |
| VH-02 | 某个具体岗位是选定任务的最终主用户。 | `USER_EVIDENCE`：样本量与协作角色覆盖由 Owner 在研究前确认。 | 该岗位实际执行 DG-04 任务、拥有核心结果责任；其他角色关系有证据。 | 改主用户并重做旅程。 |
| VH-03 | FOBrain 原生流程之后存在满足 DG-04 准入的重要剩余任务。 | `OWNER_DECISION + USER_EVIDENCE + FOBRAIN_API_DATA_EVIDENCE`：采集前同时预注册 A/B 两套最低证据要求但不选路径；冻结真实任务样本后由 DG-04 只选一路。 | A 达到最低频率与可测成本；或 B 有事件／受控演练、严重后果和足够评价案例；原生／非 AI 小改进不能按同一规则解决。 | 选择“无合格任务”，修正现有流程或停止任务。 |
| VH-04 | 某个候选解法改善单一任务结果且不恶化保护指标。 | `EXPERIMENT_RESULT`：DG-05 在看见结果前签认完整判定规则，再比较所有形态。 | 只按 DG-05 事前规则判定；未预先采用的建议值不构成通过。 | 停止失败解法或选择“不做”。 |
| VH-05 | 证据来源、时效和不确定性是采用候选解法的必要条件。 | `USER_EVIDENCE`：比较不同证据呈现。 | 目标用户按预注册规则需要且能正确使用证据。 | 降级复杂 UI，保留安全披露。 |
| VH-06 | 首版需要受审批约束的写操作。 | `OWNER_DECISION + USER_EVIDENCE + FOBRAIN_API_DATA_EVIDENCE`：确认动作、责任、补救、审批和 API。 | 一个动作具有明确任务价值、权限模型、失败补救和真实回执；否则不通过。 | 保持只读建议。 |
| VH-07 | 风险案件视图比单次聊天更适合承载任务。 | `USER_EVIDENCE + FOBRAIN_API_DATA_EVIDENCE`：用两个概念完成跨时段任务。 | 按预注册指标减少重查且不建立第二套事实。 | 退回会话 + 深链接。 |
| VH-08 | 24 项历史工具应进入首版。 | `USER_EVIDENCE + FOBRAIN_API_DATA_EVIDENCE`：逐项映射真实任务与替代。 | 只有支持 VH-03/VH-04 的必要能力进入范围。 | 未映射工具删除或延期。 |
| VH-09 | 胜出形态具有与候选产品性质相匹配的持续采用或购买依据。 | `OWNER_DECISION + ADOPTION_BUYER_EVIDENCE`：DG-04 后即并行收集 `TARGET_STAKEHOLDER_OBSERVATION`、`ADOPTION_TELEMETRY` 或 `COMMERCIAL_RECORD`；分别覆盖内部责任与资源、标准功能持续采用、增购／续费预算或独立采购／替代支出。 | DG-07 前至少有一组与候选性质匹配的可追溯承诺、实际使用／资源或预算记录和停止条件；也允许证据支持“不做／停止”。 | 不得形成 `VALIDATED_POSITIONING`；停止或调整产品关系。 |

## 真实任务采集字段

每个样本至少记录：

- 触发事件与期望业务结果；
- 执行角色、协作角色和权限；
- 使用的系统、页面、查询、导出和人工沟通；
- 输入对象、必要事实、数据缺口和歧义；
- 完成步骤、等待时间、总耗时和返工；
- 用户如何判断结果可信、完整和可行动；
- 最终决定、外部动作、审批、失败和补救；
- 任务结束的客观证据。

## 需求进入 PRD 的门槛

一个候选能力只有同时满足以下条件，才可进入 final PRD：

1. 明确映射到一个已验证用户任务和问题；
2. 有基线或可观察的现状成本；
3. 定义用户可感知的输入、输出和异常结果；
4. 明确事实来源、人机责任和失败补救；
5. 至少有一条可客观验证的成功标准；
6. 对依赖的 FOBrain 能力分别记录 `VENDOR_CLAIMED`、`DEPLOYMENT_AVAILABLE`、`INTEGRATION_AVAILABLE` 和 `TASK_EFFECTIVE`，并说明该任务的可验证增量；厂商声明层的范围重叠不能单独删除真实用户问题；
7. 不泄漏数据库、接口、框架或内部实现方案。

## 当前结论

目前仅 VH-01 达到 Owner 决策通过条件；其余 VH 均未通过。现有文档和外部研究足以形成候选蓝图草案，但不足以把主用户、核心任务、产品形态、AI 增量价值或写操作范围标为 `TASK_EFFECTIVE`、`ADOPTION_EVIDENCED` 或对应的 Owner 决策结论。
