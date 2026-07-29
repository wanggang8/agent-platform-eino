---
status: provisional
artifact_type: positioning-evidence
decision_state: dg_01_approved
owner_confirmation: dg_01_approved_2026_07_11
updated: 2026-07-11
---

# 产品定位证据与当前发现命题

## 0. 状态

早期的“通用 Agent Workbench”“AI 原生风险运营工作台”和“FOBrain AI 协作入口”均标记为 `SUPERSEDED_HISTORY`。它们不是当前推荐定位，也不得进入 Product Brief、PRD 或 Product Issues。

当前没有已成立的产品定位。DG-01 已确认的只是发现边界：

> FOBrain 仅作为第一个可验证业务样本，不是最终产品边界；先验证一个真实高价值任务，再由证据决定最终功能数量、多个工具／数据源、独立 Workbench、通用平台或“不做”。

## 1. Canonical evidence records

| ID | 关键主张 | Canonical evidence record | 证据与限制 |
| --- | --- | --- | --- |
| PE-01 | 历史需求主要列出 FOBrain 查询、详情、统计、澄清、工单状态更新和任务审计线索，而不是完整通用 Agent builder 产品面。 | `PROJECT_HISTORY / SYNTHESIS / NOT_APPLICABLE / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | [来源证据台账](../../planning-artifacts/product-blueprint/source-evidence-ledger.md)；本行只描述历史库存，目标产品适用性不在本主张中评价。 |
| PE-02 | FOBrain 官方公开声称覆盖多源资产／漏洞归一、资产关系、VPT 动态优先级和漏洞闭环。 | `VENDOR_SOURCE / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 即 `VENDOR_CLAIMED`；[FOBrain 官方产品页](https://www.huashunxinan.net/product-fobrain)。目标版本、配置、接口、采用和效果均未知。 |
| PE-03 | “独立风险运营工作台”与 FOBrain 的厂商声明范围存在能力名称重叠。 | `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 只用于阻止宣称独有差异；不能据此认定目标部署已解决任务，也不能删除真实用户问题。 |
| PE-04 | FOBrain 活跃用户在原生流程之后仍存在一个重要剩余任务。 | `NO_SOURCE / NONE / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 需要 DG-02 合规访问、DG-03 队列和近期真实任务观察。 |
| PE-05 | 该任务属于 A 高频效率型或 B 低频高损失型中的一条。 | `NO_SOURCE / NONE / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 采集前预注册两套准入规则；DG-04 在样本冻结后只选一路和一个任务。 |
| PE-06 | 某种 AI 或独立 Workbench 优于原生增强、非 AI 改进和不做。 | `NO_SOURCE / NONE / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | DG-05 先预注册，DG-06 才能按同一规则选择形态。 |
| PE-07 | 胜出形态应成为内部工具、FOBrain 标准功能、增值模块或独立产品。 | `NO_SOURCE / NONE / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | DG-04 后并行采集采用／买方证据，DG-07 决定；允许停止。 |

## 2. 只用于研究的候选角色与任务

以下不是 persona、首版范围或优先级：

- 漏洞／安全运营人员；
- 资产或业务责任人；
- 安全负责人／管理者；
- 审计、合规或风险监督角色。

任务候选仅作为关键事件访谈探针，例如责任范围核对、单资产调查、单漏洞影响确认、现有优先级解释、交接或关闭验证。正式研究不得询问“想要什么功能”，而应要求参与者重放最近一次真实任务。

## 3. 解法与产品关系保持中立

| 阶段 | 允许的结论 | 禁止提前形成的结论 |
| --- | --- | --- |
| DG-01～DG-03 | 发现边界、数据批准、研究队列 | 最终主用户、AI 必要性、产品形态、商业性质 |
| DG-04 | 一个任务 + A 或 B 唯一路径，或无合格任务 | 功能清单、入口形态、产品关系 |
| DG-05 | 预注册结果、基线、目标、保护和停止规则 | 看到结果后改指标 |
| DG-06 | 原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 改进或不做中的胜出项 | 用历史 UI 或技术偏好选形态 |
| DG-07 | 内部工具、FOBrain 标准功能、增值模块、独立产品或停止 | 用一线用户偏好代替采用／买方证据 |
| DG-08～DG-09 | 动作边界与适用边界 | 从技术实现反推责任或法规适用性 |

## 4. Forge 重新进入条件

- DG-07 后可以依据任务效果、最终主用户和采用／买方证据形成定位草案；
- DG-08、DG-09 及证据治理通过后才可形成 final forged idea；
- `VALIDATED_POSITIONING` 必须聚合独立的 `TASK_EFFECTIVE`、`ADOPTION_EVIDENCED` 和 `OWNER_DECISION_RECORDED` 记录；
- 任何条件未满足时，本文件保持 `provisional`。

## 5. 当前状态

DG-01 已通过。正式 BMAD Research 范围 `[C]` 与 DG-02 证据访问批准尚未完成；产品关系、主用户、核心结果、动作和适用范围仍必须留在各自后续 gate，不能从这次确认中推导。
