---
status: active-governance-rule
applies_to:
  - research
  - forge
  - product-brief
  - prfaq
  - prd
  - product-issues
updated_at: 2026-07-11
---

# 需求蓝图证据状态模型

## 目的

同一个陈述必须分别回答三件事：来源是否真的这样说、目标环境是否真的具备、真实用户任务是否真的因此改善。任何下游文档不得把前一层自动升级为后一层。

## Canonical evidence record

每条关键主张必须同时填写以下正交字段。不得把来源类型、主张支持度、目标适用性、成熟度或 Owner 决策压缩成一个“已确认”状态。

| 字段 | Canonical 枚举 | 含义 |
| --- | --- | --- |
| `source_type` | `OWNER_DECISION` / `TARGET_USER_OBSERVATION` / `TARGET_STAKEHOLDER_OBSERVATION` / `ADOPTION_TELEMETRY` / `COMMERCIAL_RECORD` / `AUTHORITY_SOURCE` / `INDEPENDENT_RESEARCH` / `VENDOR_SOURCE` / `PROJECT_HISTORY` / `EXPERIMENT_RESULT` / `MULTIPLE_SOURCES` / `NO_SOURCE` | 谁或什么产生了证据；经济买方、内部采用负责人和设计伙伴访谈使用 `TARGET_STAKEHOLDER_OBSERVATION`，持续使用／资源投入使用 `ADOPTION_TELEMETRY`，预算／采购／续费／增购记录使用 `COMMERCIAL_RECORD`。 |
| `claim_support` | `DIRECT` / `SYNTHESIS` / `CONTRADICTED` / `NONE` | 来源是直接支持、跨来源推演、存在直接反证，还是没有证据。 |
| `target_applicability` | `CONFIRMED` / `PARTIAL` / `UNKNOWN` / `NOT_APPLICABLE` | 是否已证明适用于目标版本、角色、任务和部署。 |
| `maturity` | `SOURCE_ONLY` / `DEPLOYMENT_AVAILABLE` / `INTEGRATION_AVAILABLE` / `TASK_EFFECTIVE` / `ADOPTION_EVIDENCED` | 能力／机会已经走到哪一层；只能用实际证据升级。 |
| `decision_status` | `OWNER_DECISION_RECORDED` / `HYPOTHESIS` / `OPEN` / `CONFLICT` / `REJECTED` | 当前产品决策状态。Owner 决策不等于用户证据。 |
| `normative_force` | `NOT_NORMATIVE` / `BINDING_LAW` / `SCOPED_DIRECTIVE` / `CONTRACTUAL` / `VOLUNTARY_STANDARD` / `GUIDANCE` / `RECOMMENDATION` / `UNKNOWN` | 规范效力；只评价单个法规、标准、合同或指南来源。 |

### Canonical shorthand

下列简称只允许作为完整字段组合的可读别名：

| 简称 | 必须对应的字段 |
| --- | --- |
| `OWNER_DECISION_RECORDED` | `source_type=OWNER_DECISION`、`claim_support=DIRECT`、`decision_status=OWNER_DECISION_RECORDED` |
| `USER_EVIDENCE` | `source_type=TARGET_USER_OBSERVATION`、`claim_support=DIRECT`；不能由 Owner 偏好代替 |
| `ADOPTION_BUYER_EVIDENCE` | `source_type=TARGET_STAKEHOLDER_OBSERVATION / ADOPTION_TELEMETRY / COMMERCIAL_RECORD` 中与主张匹配的一种、`claim_support=DIRECT`、`maturity=ADOPTION_EVIDENCED`；必须注明角色、持续使用或承诺类型、资源／预算与可撤销条件，口头兴趣不算采用证据 |
| `AUTHORITY_DIRECT` | `source_type=AUTHORITY_SOURCE`、`claim_support=DIRECT`，并必须另填 `normative_force` 与适用主体 |
| `RESEARCH_DIRECT` | `source_type=INDEPENDENT_RESEARCH`、`claim_support=DIRECT` |
| `VENDOR_CLAIMED` | `source_type=VENDOR_SOURCE`、`claim_support=DIRECT`、`target_applicability=UNKNOWN`、`maturity=SOURCE_ONLY` |
| `DEPLOYMENT_AVAILABLE` | `maturity=DEPLOYMENT_AVAILABLE`；不代表接口可用或任务有效 |
| `INTEGRATION_AVAILABLE` | `maturity=INTEGRATION_AVAILABLE`；不代表用户结果改善 |
| `TASK_EFFECTIVE` | `maturity=TASK_EFFECTIVE`；只对预注册的目标任务、版本和队列有效 |
| `ADOPTION_EVIDENCED` | `maturity=ADOPTION_EVIDENCED`；需按产品性质提供内部采用、标准功能采用、增购／续费或独立购买证据 |

`maturity` 是**单条证据记录**的层级，一条记录只能填写一个值。`TASK_EFFECTIVE` 与
`ADOPTION_EVIDENCED` 不是同一条记录的连续自动升级：最终定位必须聚合至少一条
`TASK_EFFECTIVE` 结果记录、至少一条匹配产品性质的 `ADOPTION_EVIDENCED` 记录，以及
一条 `OWNER_DECISION_RECORDED` 决策记录；任何一条都不能替代另外两条。

### 置信度规则

`高／中／低` 只能写作“来源置信度”或“综合推断置信度”，不能写入上述 canonical 字段。任何推断无论置信度多高，`claim_support` 仍为 `SYNTHESIS`，`decision_status` 仍为 `HYPOTHESIS`，直到相应门禁通过。

### 旧台账状态映射

`source-evidence-ledger.md` 中的旧值保留为库存标签，但必须按 ID 前缀得到以下 canonical 字段；下游不得只读取旧值：

| ID 前缀／旧值 | `source_type` | `claim_support` | `target_applicability` | `maturity` | `decision_status` | `normative_force` |
| --- | --- | --- | --- | --- | --- | --- |
| `U-* / CONFIRMED` | `OWNER_DECISION` | `DIRECT` | `CONFIRMED` | `SOURCE_ONLY` | `OWNER_DECISION_RECORDED` | `NOT_NORMATIVE` |
| `H-* / SUPPORTED` | `PROJECT_HISTORY` | `DIRECT` | `UNKNOWN` | `SOURCE_ONLY` | `HYPOTHESIS` | `NOT_NORMATIVE` |
| `E-* / VENDOR_CLAIMED` | `VENDOR_SOURCE` | `DIRECT` | `UNKNOWN` | `SOURCE_ONLY` | `OPEN` | `NOT_NORMATIVE` |
| `A-* / ASSUMPTION` | `PROJECT_HISTORY` | `SYNTHESIS` | `UNKNOWN` | `SOURCE_ONLY` | `HYPOTHESIS` | `NOT_NORMATIVE` |
| `C-* / CONFLICT` | `MULTIPLE_SOURCES` | `CONTRADICTED` | `UNKNOWN` | `SOURCE_ONLY` | `CONFLICT` | `NOT_NORMATIVE` |
| `M-* / MISSING` | `NO_SOURCE` | `NONE` | `UNKNOWN` | `SOURCE_ONLY` | `OPEN` | `UNKNOWN` |

任何例外必须在对应行显式填写完整字段，不能静默偏离映射。

## 产品关系状态

| 状态 | 含义 |
| --- | --- |
| `DISCOVERY_BOUNDARY` | Owner 只确认研究谁、研究什么以及哪些结论暂不预设。 |
| `SOLUTION_OPTION` | 内嵌助手、原生增强、独立 Workbench、非 AI 改进或不做中的候选项。 |
| `OWNER_DECISION` | Owner 在证据形成后选择继续验证的关系与形态，仍不等于用户或市场验证。 |
| `VALIDATED_POSITIONING` | `TASK_EFFECTIVE`、`ADOPTION_EVIDENCED` 与 `OWNER_DECISION_RECORDED` 共同支持的最终定位。 |

当前状态为 `DISCOVERY_BOUNDARY: APPROVED（DG-01，2026-07-11）`；所有产品形态仍为 `SOLUTION_OPTION`。当前 live 技术报告可把对应 capability 升级为 `INTEGRATION_AVAILABLE`，但产品关系仍未进入 `OWNER_DECISION` 或 `VALIDATED_POSITIONING`。

## 下游升级规则

- `VENDOR_CLAIMED` 只能用于防止把能力名称宣称为独有差异，不能用于认定目标部署已经解决用户任务。
- `DEPLOYMENT_AVAILABLE` 不能自动升级为 `INTEGRATION_AVAILABLE` 或 `TASK_EFFECTIVE`。
- 历史文档和历史接口最多提供候选假设；未经当前目标部署验证，不得标记为 `INTEGRATION_AVAILABLE`。
- 行业普遍问题只有获得目标用户近期任务证据后，才能成为 Product Issue 的 Background。
- 只有聚合独立的 `TASK_EFFECTIVE`、`ADOPTION_EVIDENCED` 与 `OWNER_DECISION_RECORDED` 三类记录，才能支持 `VALIDATED_POSITIONING`。
- `AUTHORITY_DIRECT` 必须记录 `normative_force`、适用主体、地域、数据和触发条件。只有 `BINDING_LAW`、`SCOPED_DIRECTIVE` 或 `CONTRACTUAL` 在适用性确认后可以直接形成合规 Must；`VOLUNTARY_STANDARD`、`GUIDANCE` 和 `RECOMMENDATION` 只能先形成治理候选。
- `source_type=MULTIPLE_SOURCES` 且 `claim_support=SYNTHESIS` 的综合主张本身必须填 `normative_force=NOT_NORMATIVE`；各底层法规、标准、合同或指南的规范效力分别保留在各自来源记录中，不得聚合成一个值。

## Product Issue 准入

一条 Product Issue 至少需要：

1. Background 对应目标用户的近期任务证据，而不只是行业报告或厂商营销页；
2. Goal 对应已确认的用户结果，不包含解法；
3. 输入、输出、异常和 Given/When/Then 可由当前证据客观验证；
4. 依赖目标部署能力时，至少达到 `DEPLOYMENT_AVAILABLE`，进入开发前必须达到相应的 `INTEGRATION_AVAILABLE`；
5. 未达到 `TASK_EFFECTIVE` 时，只能标记为实验或发现项，不能标记为已验证产品需求。
