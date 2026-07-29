---
status: draft
scope_confirmation: confirmed_2026_07_11
formal_steps_completed: []
draft_sections_completed:
  - customer-segments-and-jtbd
  - pain-points-and-alternatives
  - buying-and-adoption-decisions
  - competitive-landscape
  - opportunity-and-risk-synthesis
workflow_type: research
research_type: market
research_topic: FOBrain 现有用户的单一重要剩余任务及候选解法的增量价值
research_goals:
  - 验证 FOBrain 用户在既有产品能力之后仍存在的高成本任务与剩余痛点
  - 识别现有替代方案、竞争基线和可能的产品切口
  - 明确企业购买、部署和采用 AI 安全运营产品的决策条件
  - 判断某个已验证任务更适合 FOBrain 原生增强、内嵌 AI、独立 Workbench、非 AI 改进还是不做
  - 为 Product Brief 和 PRD 提供可追溯的市场证据与待验证假设
geographic_scope: 全球产品模式与中国企业市场语境
research_date: 2026-07-11
source_verification_status: partial
source_verification_scope: 已抽查关键官方来源；正式 BMAD Research 仍需逐项复核全部来源、适用范围与失效链接
---

# FOBrain 资产与漏洞风险运营的候选解法增量价值：市场研究草案

## 0. 文档状态与使用边界

本文件是用于继续需求发现的**研究草案**，不是正式完成的 BMAD Market Research。BMAD 市场研究流程要求产品所有者逐步确认研究范围并选择 `[C]` 后才能记录正式步骤完成；当前尚未获得这些逐步确认，因此 `formal_steps_completed` 保持为空，`scope_confirmation` 保持 `pending`。

FOBrain 官方资料证明厂商把自身定位为已承载业务能力的平台，而不是只提供原始数据的来源；厂商公开声称覆盖多源资产收集与归一、字段冲突处理、多源漏洞去重、动态优先级，以及从采集、验证到闭环处置的流程。[FOBrain 官方产品页](https://www.huashunxinan.net/product-fobrain)（`VENDOR_CLAIMED`；目标部署、集成可用性和任务效果均未验证）。因此，暂定研究对象收窄为：

> 验证 FOBrain 活跃用户是否存在一个未被现有产品充分解决、可走高频效率型或低频高损失型准入路径的单一任务，以及某种候选解法是否能在不增加关键错误和复核负担的前提下改善该任务。AI 与交付形态均不预设。

“数据源可扩展的独立风险运营平台”只保留为高风险的远期假设：在证明 FOBrain 增量入口有价值、并获得第二种数据源和第二类客户需求前，不得把它写成首版定位。该定位、主用户、购买场景和核心结果仍未由产品所有者最终确认。本文只能为后续决策提供证据，不能直接作为 final PRD 的产品事实。

### 证据记录写法

本稿只采用全蓝图统一的[证据状态模型](evidence-status-model.md)。目标能力仍必须依次验证 `VENDOR_CLAIMED → DEPLOYMENT_AVAILABLE → INTEGRATION_AVAILABLE → TASK_EFFECTIVE`，前一层不能替代后一层。

| 写法 | 所属字段／组合 | 本稿中的完整展开规则 |
| --- | --- | --- |
| `AUTHORITY_DIRECT` | Canonical shorthand | `AUTHORITY_SOURCE / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / <逐来源 normative_force>`；必须另记适用主体。 |
| `RESEARCH_DIRECT` | Canonical shorthand | `INDEPENDENT_RESEARCH / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE`；不能外推到所有客户或本产品效果。 |
| `VENDOR_CLAIMED` | Canonical shorthand | `VENDOR_SOURCE / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE`；不证明目标部署、集成、使用或效果。 |
| `HYPOTHESIS` | `decision_status` 单字段 | 不能单独充当证据状态；关键主张必须另填来源、支持、适用性、成熟度和规范效力。 |
| `UNKNOWN` | 通常为 `target_applicability` 单字段 | 不能与 `HYPOTHESIS` 互换；缺证时还需分别填写 `claim_support=NONE`、`maturity=SOURCE_ONLY`、`decision_status=OPEN` 及责任人。 |

正文为可读性使用 shorthand 或单字段时，权威记录以[研究综合草案](./research-synthesis-draft.md)的 C/S/H/X/U 六字段表为准；正文标签不得被下游单独解析为完整证据记录。

来源／综合置信度可另记为高、中、低，但它不是证据状态；即使综合置信度为高，研究推断仍保持 `HYPOTHESIS`。

## 1. 执行摘要

1. **问题空间真实存在，但 FOBrain 之后的剩余问题尚未被证明。** NIST 指南描述补丁管理的识别、排序、获取、安装和验证，并讨论优先级、测试、时限及业务／安全价值判断；CISA KEV 提供真实利用信号作为优先级输入。[NIST SP 800-40 Rev.4](https://csrc.nist.gov/pubs/sp/800/40/r4/final)、[NIST SP 1800-31](https://csrc.nist.gov/pubs/sp/1800/31/final)、[CISA KEV](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)（`AUTHORITY_DIRECT`；`normative_force=GUIDANCE`；来源置信度高）。FOBrain 官方声称覆盖多源归一、动态优先级和闭环处置（`VENDOR_CLAIMED`）；尚无目标部署或客户证据说明实际使用后仍有哪些高成本工作。

2. **ISC2 的调查样本显示 AI 安全工具已进入采用、测试或评估。** 该研究中 28% 的受访团队已集成、19% 正在测试、22% 处于早期评估；当前使用者中 63% 自报显著生产力提升，受访者也预期安全运营和漏洞管理会受影响。[ISC2 2025 Cybersecurity Workforce Study](https://www.isc2.org/Insights/2025/12/2025-ISC2-Cybersecurity-Workforce-Study)（`RESEARCH_DIRECT`；来源置信度中；不能外推为所有企业，不能证明客观效果或独立 Workbench 的购买意愿）。

3. **“统一数据 + 动态优先级 + 闭环处置”已被多个竞争产品公开宣称，FOBrain 也把它们列为既有范围。** FOBrain、Tenable、Qualys、Microsoft、CrowdStrike、Rapid7、ServiceNow，以及中国的新华三、华为、启明星辰均在官方资料中覆盖其中多项（`VENDOR_CLAIMED`；来源置信度中；不证明目标客户已部署或实际效果）。因此，这些名称不能在缺少客户与竞品实测时被宣称为独有差异；“查询 24 个工具”也不能单独成为重建理由。

4. **当前可研究的单位只能是一个剩余任务，不能先锁定产品形态。** 自然语言调用、带来源／时间／冲突／未知的结果，以及 AI 建议、人员决定与 FOBrain 真实变更的分离，都是可放入任务实验的候选解法（`HYPOTHESIS`；综合置信度低）。公开资料不能证明 FOBrain 或相邻竞品存在这些缺口；必须通过 FOBrain 用户任务观察和同任务对照验证。

5. **推荐的第一发现队列是已使用 FOBrain 且能展示现状基线的组织。** 研究先回答“FOBrain 用户是否仍有一个重要剩余任务”，而不是假定连接 FOBrain 或使用 AI 就有价值（`HYPOTHESIS`；综合置信度低）。这不排除售前账户或其他数据源用户；若剩余工作可由原生配置、培训或小幅非 AI 增强解决，则对应产品机会应停止。

6. **当前状态是 `DG-01_APPROVED_PENDING_FORMAL_C`，正式发现尚未获授权。** DG-01 已确认 FOBrain 只是第一个可验证样本；BMAD 范围 `[C]` 与 DG-02 通过后，才允许建立目标部署基线并发现一个重要任务。Product Brief 只能保持草案；在取得真实任务、同规则对照和采用证据前，不应确认任何产品形态、范围或商业价值。

## 2. 研究范围与方法

### 2.1 本轮回答的问题

- 谁实际承担资产与漏洞风险运营工作，谁协作、谁购买、谁控制风险？
- 他们需要完成哪些 Jobs-to-be-Done，现有工作为什么困难？
- FOBrain 已经如何完成这些任务，用户在其后还剩哪些符合 A 或 B 准入的重要步骤？
- 若某种候选解法胜出，企业如何评估、采用、采购、部署和扩大它？
- 全球与中国市场已有产品覆盖什么，哪些能力已成为基线？
- 原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 改进或“不做”哪一种按预注册规则胜出；若继续，什么产品关系有采用／买方证据？

### 2.2 来源层级

1. 政府、标准机构和行业协会：NIST、CISA、中国网信办、ISC2。
2. 厂商官方产品页与文档：只用于确认公开定位、能力声明和公开治理主张；不证明目标部署、集成、采用或效果。
3. FOBrain 官方资料：只用于确定 `VENDOR_CLAIMED` 范围与能力名称重叠；目标版本的既有能力基线必须由部署演示、配置和任务观察建立。
4. 项目历史资料：只作为候选 AI 交互和工具能力线索，不作为市场需求证明。

本轮不引用付费报告中的市场规模、CAGR 或无法核验的市场份额，也不以厂商“领导者”称号排序竞争力。

### 2.3 局限

- 未完成目标客户访谈、工作观察、采购访谈和竞品动手测试。
- 未获得 FOBrain 的客户分布、合同、部署方式、数据质量、现有 AI 能力和真实使用行为。
- 厂商资料只能证明官方公开提供或宣称某能力，不能证明目标客户已授权、已配置、实际使用深度或效果。
- 中国企业的安全、数据、采购和部署要求高度依行业及客户性质变化，必须由法律、安全和业务负责人逐案确认。

## 3. 市场边界与变化方向

### 3.1 本项目所在的交叉市场

这不是一个单一、边界稳定的传统品类，而是四类产品的交叉位置：

| 相邻品类 | 已有价值中心 | 对本项目的直接含义 |
| --- | --- | --- |
| 漏洞／暴露面管理 | 资产发现、风险信号汇聚、排序、修复与度量 | 最直接的功能竞争；客户可能已在现有平台内获得相同能力。 |
| SOC／安全运营平台 | 告警、调查、编排、响应、报告 | AI 工作台可能被作为现有 SOC 的一个交互层，而非独立采购。 |
| ITSM／工单与工作流平台 | 责任分派、SLA、跨团队流程、审计 | “闭环”并非独特能力；与既有系统协作比重建一套工作流更可能被接受。 |
| 安全 AI Copilot／Agent | 自然语言调查、摘要、建议和受控动作 | 多个公开产品已提供或宣称这些形态；实际采用、效果和采购基线仍需独立证据。 |

对本项目而言，最重要的产品边界不是上述品类，而是 **FOBrain 自身**：它已经处在漏洞／暴露面管理与处置工作流交叉位置。若单任务实验成立，应先比较 FOBrain 原生增强、内嵌／伴随式 AI 和独立 Workbench；只有后者在同任务中占优，才继续研究独立入口。独立产品还必须额外证明跨 FOBrain 客户／数据源的可复用价值。

### 3.2 从“漏洞列表”向“暴露面与风险运营”收敛

- NIST SP 800-40 Rev.4 描述补丁管理的识别、排序、获取、安装和验证；NIST IR 8286D 说明业务影响分析可为风险优先级提供信息。[NIST SP 800-40 Rev.4](https://csrc.nist.gov/pubs/sp/800/40/r4/final)、[NIST IR 8286D](https://www.nist.gov/news-events/news/2022/11/nist-releases-ir-8286d-using-business-impact-analysis-inform-risk)（`AUTHORITY_DIRECT`；`normative_force=GUIDANCE`；来源置信度高）。综合推断仍为 `HYPOTHESIS`，也不等于所有漏洞必须采用同一关闭流程。
- CISA KEV 提供已知在野利用漏洞信号，可作为优先级框架输入；它不单独决定所有组织处理顺序。[CISA KEV](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)（`AUTHORITY_DIRECT`；`normative_force=GUIDANCE`；来源置信度高）。
- FOBrain 官方公开声称提供多源资产关联融合、字段冲突处理、多源漏洞标准化去重、基于多维漏洞与资产信息的动态排序，以及漏洞闭环处置。[FOBrain 官方产品页](https://www.huashunxinan.net/product-fobrain)（`VENDOR_CLAIMED`；目标适用性未知）。
- Tenable、Qualys、Microsoft、Rapid7、CrowdStrike 和 ServiceNow 的官方资料均强调跨来源数据、真实利用、资产／业务上下文与修复编排，说明相似能力主张已广泛出现，而不是明显空白。[Tenable prioritization](https://www.tenable.com/products/vulnerability-management/use-cases/prioritization)、[Qualys ETM](https://docs.qualys.com/en/etm/latest/)、[Microsoft Defender Vulnerability Management](https://learn.microsoft.com/en-us/defender-vulnerability-management/defender-vulnerability-management)、[Rapid7 Exposure Command](https://www.rapid7.com/products/command/exposure-management/)、[CrowdStrike Falcon Exposure Management](https://www.crowdstrike.com/en-us/resources/data-sheets/falcon-exposure-management/)、[ServiceNow Unified Security Exposure Management](https://www.servicenow.com/uk/products/vulnerability-response.html)（`VENDOR_CLAIMED`；来源置信度中；不证明实际采用或效果）。

**`HYPOTHESIS`（综合置信度高）**：首版若仍按“恢复 24 个查询工具”或“补齐归一、排序、闭环”定义完成，很可能在 FOBrain 的厂商声明范围内重复建设；真正需要验证的是 AI 是否改善了目标部署用户尚未被解决的调查、解释、决策或协作任务。

### 3.3 AI 进入安全运营，但治理同时成为采用条件

- ISC2 的样本显示安全团队正采用或评估 AI，并把安全运营和漏洞管理视为近期高影响领域。[ISC2 2025 Workforce Study](https://www.isc2.org/Insights/2025/12/2025-ISC2-Cybersecurity-Workforce-Study)（`RESEARCH_DIRECT`；来源置信度中）。
- Microsoft Security Copilot 将权限继承、角色分配、数据位置、隐私设置和容量监控纳入正式上线步骤，且产品以用户身份访问数据，不自动获得更高权限。[Security Copilot authentication](https://learn.microsoft.com/en-us/copilot/security/authentication)、[privacy and data security](https://learn.microsoft.com/en-us/copilot/security/privacy-data-security)、[onboarding](https://learn.microsoft.com/en-us/copilot/security/get-started-security-copilot)（`VENDOR_CLAIMED`；来源真实性高；目标客户适用性未知）。
- NIST AI RMF 把 AI 测试、评估、验证和风险管理作为全生命周期活动，而不是只验证模型能回答问题。[NIST AI RMF](https://www.nist.gov/itl/ai-risk-management-framework)、[NIST AI Resource Center](https://airc.nist.gov/)（`AUTHORITY_DIRECT`；`normative_force=GUIDANCE`；来源置信度高）。

**`HYPOTHESIS`（综合置信度高）**：目标企业的采用决策可能同时评估任务效果、数据与权限、审计责任、成本和可运营性；具体否决项和顺序必须由买方访谈确认。

## 4. 客户、角色与分群

### 4.1 角色不是同一个“用户”

| 角色 | 候选职责与任务 | 价值诉求 | 关键风险 | 当前证据与置信度 |
| --- | --- | --- | --- | --- |
| 漏洞／安全运营分析员（候选主用户） | 识别责任范围、判断优先级、补齐证据、协调处置、复核结果 | 少切换系统、少重建上下文、快速形成可解释决策 | AI 错误、数据陈旧、归属错误、权限越界 | NIST 流程、CISA 优先级和历史 FOBrain 能力共同支持；未访谈（中）。 |
| 资产／业务／IT 责任人（协作用户） | 评估业务影响、安排变更、执行修复、提供完成证据 | 清楚“为什么是我、为什么现在、该做什么、怎么证明完成” | 错误分派、业务中断、不可执行建议 | NIST 指出业务／任务负责人和安全团队存在价值分歧；H3C 公开提供任务下发、审核、复验（中）。[NIST SP 800-40 Rev.4](https://csrc.nist.gov/pubs/sp/800/40/r4/final)、[H3C CSAP-iSOC](https://www.h3c.com/cn/Products_And_Solution/Proactive_Security/Product_Series/Security_Management/NSSA/NSSA/SecCenter_CSAP/) |
| 安全负责人／CISO（经济买方候选） | 确定风险偏好、资源和时限，监督风险下降 | 可衡量的风险下降、效率和责任，不是更多告警或聊天 | 投资回报不明、黑箱评分、供应商锁定 | NIST CSF 2.0 强调治理；厂商普遍提供领导视图（中）。[NIST CSF 2.0](https://www.nist.gov/news-events/news/2024/02/nist-releases-version-20-landmark-cybersecurity-framework) |
| 安全平台／数据／IT 管理员（技术把关者） | 连接数据源、配置身份权限、监控使用和成本 | 可集成、最小权限、稳定、可观测、可撤回 | 数据泄露、权限放大、部署和运维负担 | Microsoft 官方上线与权限文档直接支持（高）。 |
| 合规、审计、法务、采购（控制角色） | 审查数据、合同、保留、供应链与问责 | 可证明的控制、日志、责任和供应商承诺 | 法律适用错误、数据跨境、供应链风险 | CISA/NIST 采购指南与中国法规支持（高，但具体适用未知）。 |

**研究假设**：第一批任务发现优先观察漏洞／安全运营分析员；资产／业务责任人、安全负责人和控制角色作为协作或决策参与者同时取证。该假设未确认，不得提前写成首版用户模型。

### 4.2 按组织条件分群

人口统计对企业安全软件价值很低，应该按工作和环境分群：

| 分群 | 可观察条件 | 初步吸引力 | 进入风险 |
| --- | --- | --- | --- |
| A. 已使用 FOBrain、仍存在可观察人工任务的组织 | 用户能展示 FOBrain 原生流程之后仍需跨对象解释、补证、决策或协作的步骤 | 最高：能直接测量 AI 相对 FOBrain 原生流程的增量价值 | 剩余痛点可能很小，或可由 FOBrain 原生配置、培训和小幅增强解决。 |
| B. 多扫描器／多安全平台的大中型企业 | 数据分散，已有 CMDB／ITSM，需要跨源统一上下文 | 潜在价值高，问题强 | 成熟暴露面平台竞争最激烈，集成与替换成本高。 |
| C. 安全团队小、资产规模中等的组织 | 人手有限，希望减少重复分析 | AI 效率价值可能明显 | 预算、数据成熟度和实施资源不足，可能偏好托管服务。 |
| D. 高监管或关键业务组织 | 强权限、审计、本地化和责任要求 | 证据化和受控动作可能有价值 | 销售周期、合规、定制和部署成本高。 |
| E. 已全面采用单一大型安全生态的组织 | Microsoft、CrowdStrike、Tenable、Qualys 或 ServiceNow 覆盖完整 | 市场教育成本低 | 增量产品易被现有许可证和集成生态替代。 |

**候选 ICP（`HYPOTHESIS`）**：已使用 FOBrain，且能展示满足 DG-04 任一准入路径的真实任务，并愿意提供目标部署基线与研究证据的组织。高频效率型和低频高损失型必须分开采样、评价和决定；若两者都不存在，当前方向应停止。

## 5. Jobs-to-be-Done 与行为链路

### 5.1 核心工作链路

下列阶段描述行业工作，不代表 FOBrain 没有覆盖。其官方资料已经覆盖范围收集、漏洞归一、动态排序和闭环处置；访谈必须逐阶段寻找“原生能力之后仍发生的人工步骤”，不能把整条链路重新包装成新需求。

| 阶段 | 当…… | 我想要…… | 以便…… | 公开证据 |
| --- | --- | --- | --- | --- |
| 1. 定责与范围 | 新漏洞、暴露变化或管理要求出现 | 确认哪些资产、业务和责任人受影响 | 不遗漏关键对象，也不把无关工作分给错误的人 | NIST 强调准确、及时、完整的资产信息是有效补丁流程前提。[NIST SP 1800-31](https://csrc.nist.gov/pubs/sp/1800/31/final) |
| 2. 排序 | 待处理项远超可用资源 | 结合利用证据、暴露、业务影响和现有控制排序 | 把有限资源放在真实风险上 | [CISA KEV](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)、[NIST IR 8286D](https://www.nist.gov/news-events/news/2022/11/nist-releases-ir-8286d-using-business-impact-analysis-inform-risk) |
| 3. 调查与补证 | 排名或结论不清楚 | 看到依据、来源、时间、冲突与未知 | 决定是否相信和采取行动 | Microsoft 的 Copilot 输出仍由用户 review and assess；CrowdStrike 公开强调允许分析员询问排序原因。[Microsoft Security Copilot](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot)、[CrowdStrike Exposure Prioritization Agent](https://www.crowdstrike.com/en-us/blog/built-for-scale-powered-by-ai-innovation-driving-falcon-exposure-management/) |
| 4. 决策与交接 | 风险需要处置 | 把原因、影响、责任、时限和建议交给执行人 | 避免反复解释与责任丢失 | ServiceNow、Rapid7、H3C 均公开提供跨团队修复或工单流程。[ServiceNow USEM](https://www.servicenow.com/uk/products/vulnerability-response.html)、[Rapid7 Automation workflows](https://docs.rapid7.com/insightconnect/workflows)、[H3C CSAP-iSOC](https://www.h3c.com/cn/Products_And_Solution/Proactive_Security/Product_Series/Security_Management/NSSA/NSSA/SecCenter_CSAP/) |
| 5. 受控执行 | 系统建议改变外部状态 | 知道动作、影响和授权人，并可修改、拒绝或取消 | AI 不越权，人对高风险结果负责 | Microsoft 的 on-behalf-of 权限模型和产品角色说明权限控制已是上线条件。[Security Copilot authentication](https://learn.microsoft.com/en-us/copilot/security/authentication) |
| 6. 结果复核 | 工单声称完成 | 按关闭原因核对相应事实和授权证据 | 不把状态关闭自动解释为风险消除，也不把风险接受或缓解误写成修复 | NIST 将安装验证列为补丁管理正式阶段；其他关闭原因仍需组织规则定义。[NIST SP 800-40 Rev.4](https://csrc.nist.gov/pubs/sp/800/40/r4/final) |
| 7. 汇报与学习 | 需要管理、审计或改进 | 解释风险为何变化、谁做了什么、哪里仍未知 | 证明结果并改进流程 | NIST CSF 2.0 增加 Govern；厂商普遍提供风险趋势与报告。[NIST CSF 2.0](https://www.nist.gov/news-events/news/2024/02/nist-releases-version-20-landmark-cybersecurity-framework) |

**真正待验证的剩余任务与候选增量**：用户是否仍需跨对象重建上下文、是否难以核对现有排序依据、数据冲突和未知是否被隐藏、判断是否在聊天／会议／其他系统中反复重建，以及哪种解法能改善已选路径的结果。公开资料无法回答这些问题，AI 也不是预设答案。

### 5.2 行为特征

- **理性驱动**：真实利用、资产关键性、业务影响、可执行修复、控制覆盖和处置成本共同影响顺序（`事实／推断`，高）。
- **风险厌恶**：错误处置可能导致服务中断，错误权限可能泄露安全数据，因此用户需要证据和控制，而不仅是更快答案（`事实／推断`，高）。[NIST SP 1800-31](https://csrc.nist.gov/pubs/sp/1800/31/final)
- **跨角色依赖**：安全分析员通常不能独立完成所有修复，必须与业务、IT、平台、管理和合规角色协作（`事实／推断`，高）。
- **生态依赖**：客户已有扫描器、终端、云、CMDB、ITSM 和身份体系，新产品通常需要叠加而非一开始替换全部系统（`事实／推断`，高）。ServiceNow、Rapid7、Qualys 和 Tenable 都把第三方数据或既有工具接入作为公开产品能力。
- **自然语言不是最终记录**：聊天适合提出问题，但责任、证据、审批、状态和复核可能需要可持久化结构（`HYPOTHESIS`；综合置信度中；需用户研究验证）。

## 6. 行业问题假设与 FOBrain 剩余性待验证项

### 6.1 Canonical 证据记录与目标剩余性

下表只说明行业资料支持这些问题被关注；它不证明问题在目标 FOBrain 用户中仍然存在，也不赋予本产品 P0/P1 优先级。只有真实任务观察确认“FOBrain 原生流程之后仍存在”，才可把对应约束带入需求。

| Canonical evidence record | 候选问题 | 证据判断 | 若真实任务命中时的实验约束 |
| --- | --- | --- | --- |
| `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 待处理漏洞／暴露多，单一严重度不足以决定顺序 | CISA KEV 和 NIST 补丁管理资料支持行业背景；综合推断置信度高，但目标剩余性未知。底层来源的规范效力分别记录。 | 若任务涉及优先级，测试“为什么现在处理”的可复核解释，不生成第二套不透明分数。 |
| `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 资产、业务、身份、漏洞、工单事实可能分散或质量不一 | NIST 与多个厂商来源只支持行业关注；来源规范效力分别记录，目标剩余性未知。 | 若任务依赖多源事实，展示覆盖、新鲜度、冲突和归属。 |
| `AUTHORITY_SOURCE / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / GUIDANCE` | 安全判断与业务／IT 处置之间可能存在责任和价值分歧 | NIST 资料直接讨论角色间分歧；是否发生于目标队列仍未知。 | 测试安全理由、业务影响、责任和执行条件是否减少沟通返工。 |
| `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 交接、时限、执行和复核可能断裂 | NIST 流程与厂商工作流只形成行业线索；底层来源规范效力分别记录，目标剩余性未知。 | 若任务覆盖完成声明，按关闭原因要求对应验证证据；不预设新工作流。 |
| `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | AI 可能产生虚假确定性、越权或泄露敏感数据 | NIST AI RMF、OWASP 与厂商控制资料支持风险候选；本综合主张不具有规范效力，具体部署风险仍待确认。 | 任何 AI 实验区分事实、推断、未知，继承最小权限并保留人类控制。 |
| `VENDOR_SOURCE / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 管理层可能难以从风险分理解可行动原因与结果 | 主要从厂商报告形态推演；综合推断置信度低，缺少独立用户证据。 | 仅在管理任务被真实观察后测试摘要，不预设管理视图。 |

### 6.2 不能直接宣称的痛点

以下陈述目前缺少目标客户证据，不得写成 final PRD 事实：

- “用户每天要切换 N 个系统”以及具体耗时；
- “当前漏判率、错派率、逾期率很高”；
- “聊天是用户最喜欢的入口”；
- “FOBrain 的 24 项历史工具都是高频、必要能力”；
- “客户愿意为独立 AI 工作台新增预算”；
- “现有竞品不能展示证据、冲突或不确定性”。

## 7. 当前替代方案

| 替代方案 | 客户为何使用 | 优点 | 主要限制假设 | 新产品必须证明什么 |
| --- | --- | --- | --- | --- |
| 扫描器／暴露面平台原生控制台 | 数据已在系统中，风险排序和报表开箱可用 | 数据接近、生态完整、采购关系已存在 | 跨系统和本地业务上下文可能不足 | 不替换原平台也能显著缩短可信决策时间。 |
| SIEM／SOC／XDR 平台 | 已覆盖威胁调查、响应、自动化和报告 | 统一遥测、自动响应、成熟运营体系 | 漏洞和业务责任工作可能只是其中一部分 | 一个具体 FOBrain 剩余任务的候选解法必须优于通用 SOC 流程。 |
| CMDB + ITSM／工单 | 责任、SLA、状态和审计成熟 | 跨团队接受度高，是记录系统 | 安全证据和优先级上下文需要人工搬运 | 把安全上下文带入既有流程，而不是重建另一个工单系统。 |
| 表格、脚本、群聊、邮件和人工报告 | 灵活、低新增采购成本、适合临时流程 | 易修改，专家可控制细节 | 难以保持来源、状态和复核一致；但缺少本项目客户实证 | 用真实任务证明减少重复劳动且不降低判断质量。 |
| MSSP／顾问和内部专家 | 数据或人员成熟度不足时直接购买结果 | 有经验、能适应复杂组织 | 成本、可扩展性、知识沉淀和响应时间可能受限 | 产品是增强专家，不是用未经验证的 AI 替换责任人。 |
| 通用企业 Copilot／自建 Agent | 可利用已有 AI 平台和连接器 | 快速定制、许可可能已拥有 | 安全领域证据、权限、评估和工作流需要自行建设 | 专业任务上下文和安全治理必须显著降低自建成本。 |
| FOBrain 原生能力／小幅增强 | 官方已声明多源资产和漏洞归一、动态排序与闭环处置；目标客户的实际版本、配置和使用深度未知。[官方产品页](https://www.huashunxinan.net/product-fobrain) | 可能拥有最低新增平台、迁移和集成成本 | 是否仍缺少重要任务未知 | 所有候选解法必须相对原生流程按 DG-05 产生任务收益；若小幅增强即可解决，则其胜出或当前方向停止。 |

## 8. 买方访谈假设：购买与采用决策

本节是用于访谈 FOBrain 客户的假设清单，不是已确认的通用采购流程。角色、顺序、否决项和责任人都必须从目标客户近期采购或安全评审实例中确认。

### 8.1 触发购买或试点的事件

以下是待客户验证的候选触发器：

- 新漏洞或活跃利用事件导致现有团队无法及时确定影响范围；
- 审计、监管或管理要求证明漏洞处置时限和责任链；
- 多个扫描器、云、终端或业务系统造成重复和冲突数据；
- 安全团队人力受限，管理层要求用 AI 提升效率；
- FOBrain 用户能展示原生流程之后仍需人工完成的跨对象解释、补证、决策或交接；
- 现有暴露面平台无法适配本地业务关系、私有部署或中国生态。

### 8.2 决策参与者

**`HYPOTHESIS`（综合置信度中）**：参考软件采购与治理资料，目标客户的购买可能涉及以下角色；实际角色组合和决策权必须逐客户访谈：

- 安全负责人提出风险与结果目标；
- 漏洞／安全运营负责人定义任务和验收；
- IT／业务责任人确认处置流程不会破坏业务；
- 身份、数据、安全架构和运维团队审查权限、集成与部署；
- 法务、合规、采购审查数据处理、合同、供应链和成本；
- 财务或高管批准预算。

CISA 的软件采购指南面向 IT 决策者、采购人员和供应商，并建议把软件保障和供应商风险纳入采购生命周期；NIST CSF 2.0 也覆盖治理和供应链。[CISA Software Acquisition Guide](https://www.cisa.gov/news-events/news/cisa-unveils-tool-boost-procurement-software-supply-chain-security)、[NIST SP 1305](https://csrc.nist.gov/pubs/sp/1305/final)（`AUTHORITY_DIRECT`；`normative_force=GUIDANCE`；来源置信度高；目标买方流程仍需验证）。

### 8.3 待验证的决策标准

以下顺序仅用于组织访谈，不代表所有企业采用同一流程；目标客户可以调整、合并或否决这些项目。

| 顺序 | 标准 | 必须提供的证据 | 当前状态 |
| --- | --- | --- | --- |
| 1 | 目标任务有效性 | 与现状同任务比较：时间、完整性、错误、复核负担 | 未验证 |
| 2 | 数据适配与质量 | 目标数据覆盖率、新鲜度、实体匹配、冲突和来源可见 | 未验证 |
| 3 | 权限、隐私与责任边界 | 用户身份继承、最小权限、敏感字段控制、审批、审计、删除／保留 | 历史设计有候选约束，未形成产品要求 |
| 4 | 既有流程集成 | FOBrain、身份、CMDB、ITSM、工单和报告的真实集成演示 | FOBrain 是首版基线，其他未知 |
| 5 | 可靠性和可运营性 | 降级、失败恢复、结果一致、使用与成本监控 | 未验证 |
| 6 | 部署与合规 | 云／本地选择、数据位置、行业要求、合同和供应链审查 | 未确认 |
| 7 | 总成本与价值 | 许可、模型、集成、运维、培训成本相对节省与风险下降 | 无数据 |
| 8 | 供应商与服务能力 | 支持、更新、漏洞响应、退出和数据可携带性 | 无数据 |

### 8.4 采用旅程

建议把采用设计为可停止的门禁，而不是一次全量上线：

1. **建立现状基线**：观察多个近期案例并形成候选任务池；按采集前预注册规则只选一个符合 A 高频效率型或 B 低频高损失型的重要且可复核任务，使用同任务案例记录相应成本／后果、来源、错误和结束证据。
2. **完成数据与权限准备**：确认数据源、实体、时效、责任关系、敏感字段和用户可见范围。
3. **影子／沙箱方案验证**：DG-05 后按同一规则比较原生增强、AI、独立 Workbench、非 AI 和不做；任何候选都不直接写入外部系统，由分析员评估所选路径的结果与保护指标。
4. **有限团队试点**：只开放已验证任务，保留人工决策和外部系统原始记录。
5. **受审批动作试点**：只有 DG-08 明确允许一个动作，且所选 A/B 路径的任务价值、下游鉴权、审批、补救和复核证据成立后，才评价该动作。
6. **度量与扩大**：以任务结果、复核负担、采用率、事故和总成本决定扩大、修正或停止。

Microsoft Security Copilot 的正式上线文档先区分许可和容量，再配置角色、数据位置、插件、设置与用量监控；这直接证明 Microsoft 自身产品公开规定这些步骤，可作为目标客户访谈提示，不能外推为所有企业的统一采购顺序。[Security Copilot onboarding](https://learn.microsoft.com/en-us/copilot/security/get-started-security-copilot)、[Deploy and Operate Security Copilot](https://learn.microsoft.com/en-us/training/paths/deploy-operate-security-copilot/)（`VENDOR_CLAIMED`；来源真实性高；目标客户适用性未知）。

## 9. 竞争格局

### 9.1 阅读规则

下表只比较厂商公开资料中可核验的产品方向，不比较真实准确率、客户满意度、市场份额或价格。所有条目当前最多为 `VENDOR_CLAIMED`；目标客户的 `DEPLOYMENT_AVAILABLE`、`INTEGRATION_AVAILABLE` 和 `TASK_EFFECTIVE` 均待验证。FOBrain 不是普通候选数据源，而是单任务实验必须纳入对照的**原生替代方案**。

### 9.2 全球直接与相邻竞争者

| 产品／生态 | 官方公开能力 | 对本项目的威胁 | 可研究但未证实的切口 |
| --- | --- | --- | --- |
| Tenable One / Vulnerability Management | 统一第一方和第三方风险，结合可利用性、资产关键性、访问路径、业务影响排序，并向修复团队解释原因和衡量趋势。[官方说明](https://www.tenable.com/products/vulnerability-management/use-cases/prioritization) | 厂商已公开主张排序、数据汇聚和管理度量；若目标客户部署可用且任务有效，独立工作台可能成为重复界面。 | 比统一分数更显式的证据、冲突和未知，但必须实测，不能由营销页缺失推断。 |
| Qualys Enterprise TruRisk Management | 聚合资产与风险信号、加入业务上下文、排序和编排风险处置；可按业务实体运行排序计划。[ETM overview](https://docs.qualys.com/en/etm/latest/)、[prioritization workflow](https://docs.qualys.com/en/etm/latest/risk_management/prioritization_workflow.htm) | 厂商公开主张覆盖完整暴露面与风险运营；真实功能广度和客户效果待核验。 | 作为既有 Qualys 之上的协作层是否仍有价值，必须实测。 |
| Microsoft Defender Vulnerability Management + Security Copilot | Defender VM 用威胁情报、利用概率、业务上下文和设备评估排序；Security Copilot 提供自然语言、插件、组织数据上下文和受角色控制的操作。[Defender VM](https://learn.microsoft.com/en-us/defender-vulnerability-management/defender-vulnerability-management)、[Security Copilot](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot) | 厂商公开提供在 Microsoft 许可、身份和安全生态内的 AI 方案；目标客户授权、配置与效果待核验。 | 非 Microsoft／本地私有生态、多数据源案件模型，但必须证明集成成本可接受。 |
| CrowdStrike Falcon Exposure Management + Charlotte AI | 持续暴露面、威胁与环境上下文排序、自然语言询问排序原因、受控 Agent 和自动化工作流。[Exposure Management](https://www.crowdstrike.com/en-us/resources/data-sheets/falcon-exposure-management/)、[prioritization agent](https://www.crowdstrike.com/en-us/blog/built-for-scale-powered-by-ai-innovation-driving-falcon-exposure-management/)、[Charlotte AI](https://www.crowdstrike.com/en-us/platform/charlotte-ai/agentic-security-workforce/) | 厂商已公开主张“AI 解释为什么排序 + 自动行动”；不能据此认定客户成熟使用或任务有效。 | 本项目不能仅靠可解释排序区分；需比较具体证据粒度、冲突处理和跨系统责任记录。 |
| Rapid7 Exposure Command | 聚合原生与第三方数据，以利用可能性、可达性、严重度和业务上下文排序；含修复指导、自动化、SLA，并可结合 Rapid7 Automation 工作流。[产品页](https://www.rapid7.com/products/command/exposure-management/)、[Automation workflow docs](https://docs.rapid7.com/insightconnect/workflows) | 厂商已公开主张可插拔数据、统一资产、工单和自动化工作流；目标部署与任务效果未知。 | 风险案件的人机决策和结果复核体验需做同任务对比。 |
| ServiceNow Unified Security Exposure Management | 汇聚扫描器、云、容器和代码结果，映射 CMDB，应用威胁情报和 AI 排序，并通过 SLA 编排跨 IT 团队修复。[官方产品页](https://www.servicenow.com/uk/products/vulnerability-response.html) | 厂商公开主张其产品结合工作流、CMDB 和协作；客户实际部署优势待核验。 | 是否更适合作为安全证据层接入既有 ServiceNow，必须由同任务比较决定。 |

### 9.3 中国市场直接与相邻竞争者

| 产品／厂商 | 官方公开能力 | 对本项目的含义 |
| --- | --- | --- |
| FOBrain | 厂商公开声称提供多源资产收集关联与字段冲突处理、多源漏洞归一去重、融合漏洞与资产信息的动态优先级、从漏洞采集和验证到闭环处置。[官方产品页](https://www.huashunxinan.net/product-fobrain) | 单任务实验的原生对照；目标版本、配置、API、真实使用和任务效果必须分别核验。任何新需求都必须指出目标部署原生流程之后的剩余任务。 |
| 新华三 AI SOC / CSAP-iSOC | 厂商公开声称提供自然语言资产排查、告警分析、报告，以及资产、风险、多源数据、工单下发、审核、状态跟踪和复验。[AI SOC](https://www.h3c.com/cn/d_202510/2659785_30008_0.htm)、[CSAP-iSOC](https://www.h3c.com/cn/Products_And_Solution/Proactive_Security/Product_Series/Security_Management/NSSA/NSSA/SecCenter_CSAP/) | 相似能力名称在中国厂商主张中已经出现；本地化或功能名称不能单独构成差异。 |
| 华为星河 AI 网络安全 Agentic SOC / HiSec NG-SIEM | 厂商公开声称提供多源数据标准化、AI 分析、SOAR 编排、人工确认／自动模式、资产风险和自动化处置。[HiSec NG-SIEM](https://e.huawei.com/marketingcloud/pep/asset/20000001/Material/3c95eefd807f469eb9b233632dfd364f/M3T1A590N1250415141188292742/%E5%8D%8E%E4%B8%BAHiSec%20NG-SIEM%E4%B8%8B%E4%B8%80%E4%BB%A3%E5%AE%89%E5%85%A8%E8%BF%90%E8%90%A5%E5%B9%B3%E5%8F%B0%E5%BD%A9%E9%A1%B5.pdf)、[华为 Agentic SOC](https://e.huawei.com/cn/news/2026/solutions/enterprise-network/launches-network-security-agentic-soc) | 能力广度上的竞争风险需要真实客户与竞品实测确认，不能仅据厂商材料排序。 |
| 启明星辰安星／AISOP | 厂商公开声称多智能体覆盖日志、告警、漏洞定级、响应、风险和情报，并提供自定义智能体、知识库和可视化工作流。[官方说明](https://www.venustech.com.cn/new_type/cpdt/20250804/28677.html) | 中国厂商也公开提出垂直 Agent 与构建能力；通用 Agent builder 只能保留为未验证选项。 |
| 安天 XDR／安全运营平台 | 厂商公开声称覆盖多渠道资产融合、漏洞风险、上下文调查、响应编排和持续运营。[官方产品页](https://www.antiy.cn/Security_Product/XDR.html) | 非生成式 AI 自动化是必须纳入同任务比较的替代选项；实际强弱待实测。 |

### 9.4 已被多个产品公开提供或宣称的能力

以下均不能单独作为产品差异化：

- 自然语言查询资产、漏洞、告警或风险；
- 聚合多个第一方／第三方数据源；
- 结合威胁情报、可利用性、资产关键性和业务影响排序；
- 风险分、趋势、仪表盘和管理报告；
- 工单创建、责任分派、SLA、通知和进度跟踪；
- 自动化／编排修复流程；
- 解释“为什么优先”或提供修复指导；
- 角色权限、审计、数据控制和人工确认。

公开资料不能证明这些能力的实际效果、普及率或采购基线，但足以说明它们不能在缺少客户和竞品实测时被当作独有市场空白。

## 10. FOBrain 单任务候选解法与差异化假设

### 10.1 候选切口

| 假设 | 候选价值 | 为什么可能成立 | 为什么可能不成立 | 置信度 |
| --- | --- | --- | --- | --- |
| FOBrain AI 协作／决策入口 | 用自然语言调用既有能力，串联跨对象调查，解释结果并保留人机决策记录 | 不重建数据、排序和闭环，最小化范围；历史工具能力可以作为原型输入 | FOBrain 可能已有或可低成本增加同类能力；用户可能不需要自然语言 | 中低 |
| 证据与不确定性优先 | 每个结论显示来源、时间、适用范围、数据冲突、未知和推断 | 高风险决策需要复核；公开产品普遍强调“上下文”，但不一定把不确定性作为核心体验 | 成熟产品可能已有同等能力，只是营销页未显示；用户可能只需要排序结果 | 中低 |
| 风险案件是稳定工作单位 | 把范围、证据、决定、责任、动作和复核放进一个可追踪对象 | 解决聊天、工单和报表之间上下文丢失 | ServiceNow 等已有记录和流程模型；新对象可能增加重复录入 | 中低 |
| 数据源中立的独立产品 | FOBrain 首接入，未来接其他来源；结果回写原记录系统 | 企业不愿替换已投资的安全与流程平台 | 当前没有第二数据源、第二客户或独立采购证据；连接器维护成本高，巨头也支持第三方数据 | 低；高风险远期假设 |
| 明确 AI／人的责任边界 | 分离模型建议、人员决定、审批和外部变更，并验证动作结果 | 权限、问责和信任是 AI 采用门槛 | 这也正成为成熟 AI 平台基线 | 中 |
| 中国私有环境和本地业务关系适配 | 适配本地部署、中文安全语义、组织／业务／责任关系和 FOBrain | 全球平台在部分客户环境可能难以落地 | 国内大厂已经强覆盖；适配需求和付费意愿未知 | 低 |

### 10.2 推荐的价值主张实验

不建议验证“AI 能不能回答问题”，建议验证：

> 对同一个经 DG-04 选择的真实任务，原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 改进或“不做”中，哪一个按 DG-05 预注册的 A 或 B 结果与保护指标胜出。

该实验把差异化从功能数量转为 DG-05 预注册的可观察结果：A 可使用净效率与质量保护，B 使用高损失任务结果与严重错误保护。它仍是建议的验证方法，不是已成立的产品价值。

### 10.3 交付形态必须同时比较

即使找到真实任务，也不能直接推出“独立 Workbench”。同一任务必须公平比较以下候选及“不做”：

| 解法 | 适用条件 | 主要优势 | 否决条件 |
| --- | --- | --- | --- |
| FOBrain 原生产品增强 | 任务可由现有对象、页面、规则或原生功能直接改善 | 最低切换、采购和治理成本 | 按 DG-05 仍不能改善所选路径结果。 |
| FOBrain 内嵌／伴随式 AI | 任务需要语言／非结构化辅助，同时高度依赖 FOBrain 页面上下文与权限 | 保留原工作上下文，新增入口成本较低 | AI 复核／风险抵消收益，或原生／非 AI 解法按同一规则更优。 |
| 独立 Workbench | 同一重要任务稳定跨越多个系统，且独立工作对象按 DG-05 规则明显优于原产品导航 | 可统一人机协作体验和跨系统证据 | 只有 FOBrain 单一场景，或用户不愿维护第二工作界面。 |
| 非 AI 流程／配置改进 | 痛点来自字段、导航、报表、规则、培训或责任流程，而非推理和语言交互 | 更确定、成本低、风险小 | 无法处理频繁变化、非结构化上下文或高成本调查。 |

“不做”不是附加项，而是 DG-06 的正式候选结果：如果其他解法相对 FOBrain 原生流程都没有通过预注册规则，就没有新产品需求。

## 11. 中国企业市场语境

### 11.1 已确认的法规背景

- 2025 年修订的《中华人民共和国网络安全法》已于 2026 年 1 月 1 日施行。[中国网信网全文](https://www.cac.gov.cn/2025-12/29/c_1768735112911946.htm)、[主席令](https://www.cac.gov.cn/2025-10/29/c_1763461514706581.htm)（`AUTHORITY_DIRECT`；`normative_force=BINDING_LAW`；具体主体与项目适用性待法务确认）。
- 《国家网络安全事件报告管理办法》自 2025 年 11 月 1 日起施行，规定报告主体、流程、时限和内容。[中国网信网](https://www.cac.gov.cn/2025-09/15/c_1759583017563621.htm)（`AUTHORITY_DIRECT`；`normative_force=BINDING_LAW`；触发主体与事件级别待确认）。
- 《生成式人工智能服务管理暂行办法》适用于向中国境内公众提供生成式 AI 服务；企业等研发或应用生成式 AI、但不向境内公众提供服务的，不适用该办法。[中国网信网](https://www.cac.gov.cn/2023-07/13/c_1690898326795531.htm)（`AUTHORITY_DIRECT`；`normative_force=BINDING_LAW`；本项目服务对象待确认）。

### 11.2 对产品的谨慎推断

- 如果产品仅供企业内部安全团队使用，不能自动假定适用公众生成式 AI 服务规则；但网络安全、数据安全、个人信息、行业监管、合同和客户内部制度仍可能适用（`HYPOTHESIS`；综合置信度中）。
- 产品可能处理人员、部门、资产、漏洞、业务关系、日志和工单，其中部分数据可能涉及个人信息、重要数据或高敏安全信息；数据分类、最小必要、保留、部署位置和模型提供方需要逐客户确认（`HYPOTHESIS`；综合置信度高，但具体数据类别为 `UNKNOWN`）。
- “符合中国法规”不能写成一个通用验收标准；必须由法律／合规负责人按部署模式、用户范围、行业、数据和模型调用方式出具适用性判断（`UNKNOWN`）。

## 12. 市场与产品风险

| 风险 | 发生可能性 | 影响 | 早期证据／触发条件 | 缓解或停止条件 |
| --- | --- | --- | --- | --- |
| 重复 FOBrain 原生能力 | 高 | 极高 | 需求仍是多源归一、动态排序、漏洞闭环或既有 24 项查询的重新实现 | 每项需求先写“FOBrain 现状／剩余任务／候选解法增量”；无法写出则删除。 |
| 被现有安全平台功能覆盖 | 高 | 高 | 客户已有 Tenable／Qualys／Microsoft／CrowdStrike／Rapid7／ServiceNow 或国内 AI SOC | 必须用真实任务证明增量价值；不能则停止独立产品方向。 |
| 把 FOBrain 增量能力误包装成独立平台 | 高 | 高 | 所有高价值任务只依赖 FOBrain 私有对象，且买方、合同和部署仍属于 FOBrain | 由 DG-07 根据采用／买方证据选择标准功能、增值模块、内部工具、独立产品或停止；无第二数据源和独立采购证据时禁止平台化。 |
| 数据质量使 AI 结论不可用 | 高 | 高 | 资产归属、人员、业务关系、漏洞和状态陈旧或冲突 | 把覆盖和新鲜度作为试点前置门禁；未达门槛不进入自动建议。 |
| 用户不信任或复核成本更高 | 中高 | 高 | 每次回答都需重新查原系统，证据不能定位 | 影子模式按 DG-05 保护指标评价；A 若净时间不改善则暂停，B 若关键结果／严重错误保护不通过则暂停。 |
| AI 越权或错误动作 | 中 | 极高 | 权限继承错误、审批绕过、动作不可撤回或无法验证 | DG-08 前只做只读／影子实验；若 DG-08 允许写操作，必须独立授权、确认、审计、补救和结果复核。 |
| 集成和维护成本失控 | 高 | 高 | 连接器、字段、身份、模型和流程高度客户定制 | 首版限制为 FOBrain 原生接口；只有第二数据源需求验证后才抽象通用连接器。 |
| 价值无法度量 | 高 | 高 | 只统计对话量、工具调用或功能覆盖 | 上线前固定任务基线和结果指标；无法测量则不扩大。 |
| 产品定位持续摇摆 | 高 | 高 | 同时要求通用 Agent builder、SOC、漏洞平台、数据源平台和 FOBrain 全量恢复 | 产品所有者先锁定“只验证一个重要剩余任务”的发现边界和准入路径；通过门禁后才选择形态与关系。 |
| 合规与采购周期被低估 | 中高 | 高 | 数据位置、模型供应商、日志保留、供应链材料不完整 | 试点前完成数据流、角色、保留和供应商清单；未通过不接生产数据。 |

## 13. 可证伪的验证计划

以下门禁是研究建议，不是行业事实。任务采集阶段只记录原始频率、成本、后果和可评价性，不预选 A/B；产品所有者必须在采集前同时固定 A/B 两套准入定义与最低证据要求。原始样本冻结后，DG-04 按预注册规则选择唯一任务和唯一准入路径。实验指标则必须在 DG-05、任何方案实验之前固定。

### Gate M0：FOBrain 产品基线

先获得目标客户所用 FOBrain 版本、配置、角色、数据源、原生页面／流程、接口和使用记录，针对候选任务逐项记录：`原生已完成`、`可通过配置完成`、`需要人工补充`、`明确缺失`。

通过条件：每个候选任务都能指出目标部署原生输入、原生输出、剩余人工步骤、失败成本和待测增量。无法取得产品基线时，不得把历史 24 项工具或官方营销页的缺口直接转成需求。

### Gate M1：问题存在

以下样本量只是待 Owner 采用的建议。Owner 在采集前分别签认两条路径的阈值和最低证据要求，但不选择路径；采集时同时保留频率／可测成本与事件／严重后果／可评价性原始证据，不向参与者预设 A 或 B。样本冻结后，DG-04 依据预注册规则只选择一种路径；两条路径不得混成一个通过率。

高频效率型建议条件：

- 至少 8/12 将“FOBrain 原生流程之后仍需重建上下文、解释或协调决定”列为每周发生且高成本的问题；
- 至少 6 次观察中有 4 次出现来源追溯、数据冲突、责任确认或重复交接；
- 能获取现状任务时长、参与角色、输入、输出和失败结果，而不只有主观意见。

低频高损失型不使用“每周一次”条件，必须另有已发生事件或经批准的受控演练、明确严重后果，以及足够建立金标准和安全对照的案例；时限只有在真实任务本身存在时才记录，不是通用准入条件。若两条路径均无证据，或任务主要由单系统低成本解决，则停止当前定位。

### Gate M2：现状替代方案基线

此阶段只审计真实现状：FOBrain 原生流程、已有配置／培训／保存视图、当前人工流程和已经投入使用的替代工具。可以核验相邻产品公开主张与目标客户实际可用性，但不得制作或比较内嵌 AI、独立 Workbench 等候选原型，也不得在 DG-05 前宣布任何形态胜出。

建议通过条件：

- 至少 3 个组织愿意展示现有流程和缺口；
- 缺口不能通过 FOBrain 的简单配置、培训或小幅非 AI 增强低成本解决；
- 现状缺口具有可复核证据，并能进入 DG-04 单一任务选择；
- 若原生配置、培训或小幅非 AI 改进已足够，则停止相应产品机会，不进入方案实验。

### Gate M3：数据可用

实验只选择一个任务；使用该任务的 3–5 个代表性案例／场景，检查 FOBrain 的对象覆盖、主键、归属、新鲜度、冲突、权限和关闭原因所需验证证据。3–5 指案例，不是 3–5 个产品任务。

建议通过条件：

- 每个任务所需事实、来源和时间均可定位；
- 关键缺失和冲突能被检测并呈现，而非由模型猜测；
- 用户权限能够映射到数据和动作范围；
- 至少有一种方式复核处置结果。

### Gate M4：任务收益

先由 DG-05 在查看结果前签认主结果、基线、目标、保护指标、样本、胜出规则和停止线，再在只读／影子模式对同一任务的代表性案例比较 FOBrain 原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 改进和“不做”。所有形态使用同一规则，任何一项都可胜出。

建议通过条件：

- A 路径仅在 Owner 于 DG-05 采用时，才可使用“中位任务总时长比现状下降至少 30%”等效率建议值；
- B 路径必须使用已预注册的高损失任务结果、严重错误保护和受控案例规则，不以时间改善替代；
- 两条路径的关键事实完整性和错误率均不得违反各自保护指标；
- 复核者无需回到原系统重建大部分证据；
- 没有权限越界、无证据确定性结论或未审批准写入。

### Gate M5：采用与购买意愿

建议通过条件：

- 至少 3 个设计伙伴愿意提供真实数据、用户时间和安全评审资源；
- 至少 2 个具有预算影响力的负责人能说明这是 FOBrain 增购／升级、独立采购还是内部能力，并说明预算来源和替代支出；
- 产品的部署、数据和责任方案能通过至少一个目标客户的安全／合规预审。

如果只能获得“看起来不错”的反馈，而没有数据、时间、流程或预算承诺，则不认为需求已验证。

## 14. Owner gate 与后续客户研究问题

### 14.1 产品所有者决策

DG-01 已于 2026-07-11 通过：

> FOBrain 仅作为第一个可验证业务样本，不是最终产品边界；先验证一个真实高价值任务，再由证据决定最终功能、工具／数据源和产品形态。

下一步先在正式 BMAD Market Research 中取得范围 `[C]`，再处理 DG-02。其余问题不得合并：研究岗位留到 DG-03，采集前 A/B 双规则与 DG-04 任务／路径分开记录，结果指标留到 DG-05，形态留到 DG-06，产品关系留到 DG-07，动作和适用性分别留到 DG-08／DG-09。

### 14.2 DG-02／DG-03 通过后的客户验证问题

1. 最近一次确定某漏洞是否应优先处理，从触发到决定经历了什么？请展示实际输入、输出和系统切换。
2. 哪一步最耗时、最常返工或最容易出错？失败造成了什么结果？
3. 谁能决定优先级，谁能批准动作，谁实际执行，谁验证完成？
4. 你怎样判断现有结果可信、完整并可行动？哪些证据缺失时必须回到来源系统？
5. FOBrain 原生流程已经完成了什么、哪里需要离开 FOBrain 或人工补充？为什么不能通过配置、培训或小幅增强解决？
6. 哪些数据不能离开本地或不能进入模型？数据需要保留多久？
7. 这次真实任务由谁用什么结果判断成功或失败，现状是否记录调查、分派、修复、暴露、逾期、复发或审计成本？
8. 仅在 DG-04 选定任务后，再向采用／买方角色追问：谁拥有资源或预算，什么触发采用／采购，必须替代什么成本或风险才能批准？

## 15. 对下游需求蓝图的约束

在完成上述门禁前，Product Brief 和 PRD 应遵守：

- 主用户、核心结果、FOBrain 角色和动作边界均标为 `PROVISIONAL`；
- FOBrain AI 入口、独立 Workbench 和数据源平台均不得标记为已成立定位；
- 在范围扩展前，只允许把一个经真实观察、通过 DG-04 单一路径和 DG-05 预注册判定规则的任务作为首个需求候选；
- 必须按 DG-05 同一规则记录并比较 FOBrain 原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 改进和“不做”，不能从“适合 AI”直接跳到“需要独立产品”；
- 不把 24 个历史工具直接转成 24 个产品需求；
- 不把自然语言、风险分、工单、仪表盘或多数据源接入写成差异化；
- 每个产品目标必须对应可观察的用户结果和可测验收，而非“链路跑通”；
- 把证据来源、数据新鲜度、冲突、未知、权限和结果复核写成产品需求候选；
- 保留既有扫描器、FOBrain、CMDB 和 ITSM 的记录系统地位，除非客户证据证明需要替换；
- 通用 Agent builder、无人监督修复、全量多租户平台和无证据综合风险分保持首版非目标。

## 16. 研究结论

市场和标准证据足以确认：企业资产与漏洞风险运营是一个持续、跨角色、需要优先级和结果验证的真实工作；行业调查与厂商官方资料显示 AI 正被评估，并已有产品公开提供或宣称数据汇聚、业务上下文、排序、工作流和自然语言等能力。FOBrain 也公开声明覆盖多源归一、动态优先级和闭环处置。这些材料不证明实际采用效果、目标部署可用性或本项目客户的剩余问题。

这些证据不能支持一个简单结论：本项目不能因为拥有聊天、工具调用、风险卡片和 FOBrain 查询就认为找到了产品机会，也不能把“FOBrain AI 协作入口”当作已成立定位。DG-01 只确认 FOBrain 是首个样本；取得正式范围 `[C]` 与 DG-02 批准后，先找到**一个**通过预注册准入路径的重要任务，再按预注册结果比较原生增强、AI、独立 Workbench、非 AI、多工具／数据源或通用平台方向和不做。

目前没有足够证据选择产品形态。任何“可信决策层”“风险案件”或“数据源可扩展平台”都只是待测解法。能否形成产品，取决于三个尚未获得的证据：FOBrain 用户真实工作观察、相对原生与非 AI 替代方案的任务收益、以及设计伙伴的数据与采用承诺。在这三项通过前，本研究状态必须保持 `draft`。

## 17. 核心来源索引

### 标准、政府与行业研究

- [NIST SP 800-40 Rev.4：企业补丁管理规划](https://csrc.nist.gov/pubs/sp/800/40/r4/final)
- [NIST SP 1800-31：改进企业补丁管理](https://csrc.nist.gov/pubs/sp/1800/31/final)
- [NIST IR 8286D：用业务影响分析支持风险排序](https://www.nist.gov/news-events/news/2022/11/nist-releases-ir-8286d-using-business-impact-analysis-inform-risk)
- [NIST Cybersecurity Framework 2.0](https://www.nist.gov/news-events/news/2024/02/nist-releases-version-20-landmark-cybersecurity-framework)
- [NIST AI Risk Management Framework](https://www.nist.gov/itl/ai-risk-management-framework)
- [CISA Known Exploited Vulnerabilities Catalog](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)
- [CISA Software Acquisition Guide 工具公告](https://www.cisa.gov/news-events/news/cisa-unveils-tool-boost-procurement-software-supply-chain-security)
- [ISC2 2025 Cybersecurity Workforce Study](https://www.isc2.org/Insights/2025/12/2025-ISC2-Cybersecurity-Workforce-Study)
- [中华人民共和国网络安全法（2025 修订）](https://www.cac.gov.cn/2025-12/29/c_1768735112911946.htm)
- [国家网络安全事件报告管理办法](https://www.cac.gov.cn/2025-09/15/c_1759583017563621.htm)
- [生成式人工智能服务管理暂行办法](https://www.cac.gov.cn/2023-07/13/c_1690898326795531.htm)

### 全球产品官方资料

- [Tenable vulnerability prioritization](https://www.tenable.com/products/vulnerability-management/use-cases/prioritization)
- [Qualys Enterprise TruRisk Management](https://docs.qualys.com/en/etm/latest/)
- [Microsoft Defender Vulnerability Management](https://learn.microsoft.com/en-us/defender-vulnerability-management/defender-vulnerability-management)
- [Microsoft Security Copilot](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot)
- [CrowdStrike Falcon Exposure Management](https://www.crowdstrike.com/en-us/resources/data-sheets/falcon-exposure-management/)
- [Rapid7 Exposure Command](https://www.rapid7.com/products/command/exposure-management/)
- [ServiceNow Unified Security Exposure Management](https://www.servicenow.com/uk/products/vulnerability-response.html)

### 中国产品官方资料

- [FOBrain 网络资产攻击面管理平台](https://www.huashunxinan.net/product-fobrain)
- [华为星河 AI 网络安全 Agentic SOC](https://e.huawei.com/cn/news/2026/solutions/enterprise-network/launches-network-security-agentic-soc)
- [新华三新一代 AI SOC](https://www.h3c.com/cn/d_202510/2659785_30008_0.htm)
- [新华三 CSAP-iSOC](https://www.h3c.com/cn/Products_And_Solution/Proactive_Security/Product_Series/Security_Management/NSSA/NSSA/SecCenter_CSAP/)
- [启明星辰安星人工智能安全运营系统](https://www.venustech.com.cn/new_type/cpdt/20250804/28677.html)
- [安天鉴形 XDR](https://www.antiy.cn/Security_Product/XDR.html)
