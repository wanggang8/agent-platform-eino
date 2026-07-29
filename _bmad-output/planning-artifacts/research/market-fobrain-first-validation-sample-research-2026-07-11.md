---
stepsCompleted: [1, 2, 3, 4, 5, 6]
inputDocuments: []
workflowType: 'research'
lastStep: 1
research_type: 'market'
research_topic: '以 FOBrain 为首个可验证业务样本的 AI 协作产品机会与最终产品边界'
research_goals:
  - '验证首个 FOBrain 样本中是否存在真实高价值剩余任务'
  - '分析用户、替代方案、竞争格局、采用与购买条件'
  - '由证据判断最终产品是否应扩展到多个工具、多个数据源、独立 Workbench、通用平台或不做'
user_name: 'Vick'
date: '2026-07-11'
web_research_enabled: true
source_verification: true
---

# Market Research：以 FOBrain 为首个可验证业务样本的产品机会

**Date:** 2026-07-11
**Author:** Vick
**Research Type:** Market Research

---

## Executive Summary

企业安全运营面对持续的漏洞利用、数据分散、技能与预算约束，AI 安全工具也已进入广泛评估阶段：ISC2 2025 研究中，28% 已集成、19% 正测试、22% 处于早期评估。但这些行业信号只能证明“值得验证”，不能证明 FOBrain 客户存在同样问题，也不能证明 AI、独立 Workbench 或通用平台是正确答案。[ISC2 2025 Workforce Study](https://www.isc2.org/Insights/2025/12/2025-ISC2-Cybersecurity-Workforce-Study)

竞争已经覆盖原生增强、Exposure 平台、综合 SecOps、ITSM、MSSP、非 AI 内部方案和通用 Agent。成熟厂商拥有原生数据、权限、工作流、渠道和合同，且多源、攻击图、自然语言、Agent 与自动化均已商品化。因此当前不存在可验证的产品定位、市场份额、买方、定价或 TAM；FOBrain 只能作为第一个真实样本。

研究的战略结论是 `PIVOT / UNVALIDATED`：下一步不是写功能清单，而是通过 DG-02 批准最小数据与研究访问，随后冻结真实样本，选择一个可证伪的高价值任务，并比较 FOBrain 原生增强、嵌入式 AI、独立入口、非 AI 和不做。只有任务效果、持续采用和所有者决策同时成立，才允许形成 Product Brief 与定位。

## Table of Contents

1. [Research Initialization](#research-initialization)
2. [Customer Behavior and Segments](#customer-behavior-and-segments)
3. [Customer Pain Points and Needs](#customer-pain-points-and-needs)
4. [Customer Decision Processes and Journey](#customer-decision-processes-and-journey)
5. [Competitive Landscape](#competitive-landscape)
6. [Research Synthesis and Strategic Recommendations](#research-synthesis-and-strategic-recommendations)

---

## Research Initialization

### Research Understanding Confirmed

**Topic**：以 FOBrain 为首个可验证业务样本的 AI 协作产品机会与最终产品边界  
**Goals**：验证真实高价值任务；分析用户、替代方案、竞争、采用与购买条件；由证据决定最终产品范围。  
**Research Type**：Market Research  
**Date**：2026-07-11

### Research Scope

**市场分析重点：**

- 首个 FOBrain 样本中的用户任务、现状成本、替代方案和未满足需求；
- 安全运营 AI、内嵌助手、独立 Workbench、非 AI 改进及通用 Agent 平台的竞争与替代关系；
- 目标组织的采用、信任、安全评审、预算和购买条件；
- 从单一样本外推到多个工具／数据源或通用平台所需的额外证据与停止条件。

**明确边界：**

- FOBrain 只是首个验证样本，不是最终产品边界；
- 不预设 AI 必要、功能数量、产品形态、商业关系或“不做”；
- 市场规模只在产品类别和购买者被证据确定后估算，避免为未定义产品制造数字。

**研究方法：**

- 使用当前公开资料并逐条验证来源；
- 关键主张尽量使用多个独立来源；
- 区分权威来源、独立研究、厂商声明、项目历史和真实用户证据；
- 对未知、适用性和置信度显式标记。

### Next Steps

1. ✅ 初始化并提出范围（当前步骤）
2. 用户与任务行为分析
3. 竞争与替代方案分析
4. 市场机会、采用／购买证据和战略综合

**Research Status**：范围已由用户于 2026-07-11 选择 `[C]` 确认；进入客户与任务行为分析。

**Scope Confirmed by User:** 2026-07-11

## Customer Behavior and Segments

### Customer Behavior Patterns

本研究是企业级安全产品，不适合套用消费品的年龄、收入和生活方式分群。可观察的客户行为单位应是“岗位 × 组织环境 × 任务 × 决策权”。公开证据显示，安全团队对 AI 的态度是**务实评估而非普遍成熟采用**：ISC2 2025 AI Pulse Survey 的 436 名受访者中，30% 已集成 AI 安全工具，42% 正在评估或测试；组织选择工具时主要依靠内部研究（72%）、培训教育（56%）、厂商建议（45%）和政府／公共政策信息（33%）。这意味着本项目必须以可复现实验、培训成本和治理证据支持采用，不能只靠演示或厂商品牌。[ISC2 2025 AI Pulse Survey](https://www.isc2.org/insights/2025/07/2025-isc2-ai-pulse-survey)

采用者首先期待改善耗时且可验证的工作。ISC2 受访者认为短期 AI 影响最大的方向包括网络监控／入侵检测（60%）、终端检测响应（56%）和漏洞管理（50%）；已使用者中 70% 报告团队效果有积极改善。但这些是自报调查，不能证明本项目的具体任务有效。[ISC2 调查说明](https://www.isc2.org/Insights/2025/07/ISC2-Research-Cybersecurity-Teams-Cautious-on-AI-Adoption)

公开证据也包含反向信号：SANS 2024 SOC Survey 的 403 人样本中，AI/ML 分析技术满意度 GPA 从 2.17 降至 1.99，生成式 AI/ML 为 1.80，计划实施比例从 21% 降至 11%。该样本北美占比较高且受访身份未经独立验证，但它足以否定“市场兴趣自然等于满意采用”的假设。[SANS 2024 SOC Survey](https://files.abnormalsecurity.com/production/files/SANS-2024-SOC-Survey.pdf?dm=1720793408)

厂商实验提供了“任务级验证”而非“产品级成功”的参考：微软在 147 名有经验安全专业人员的随机对照试验中报告平均快 22%、准确率高 7%、97% 愿意再次使用；研究由产品厂商实施，外部适用性有限，只能说明速度、准确性和再次使用意愿可以被同任务对照测量。[Microsoft RCT infographic](https://www.microsoft.com/content/dam/microsoft/final/en-us/microsoft-brand/documents/EconStudy-Speed-accuracy-and-sentiment-Infographic-Final.pdf)

_Behavior Drivers：减少重复调查、提高准确性、扩展技能，而不是“拥有 AI”。_  
_Interaction Preferences：优先在真实任务中查看证据、回到来源并对比不同候选形态。_  
_Decision Habits：先内部研究／试用，再结合培训、治理、安全评审和预算决定。_  
_Confidence：中；有公开调查和厂商 RCT，但没有目标客户观察。_

### Demographic Segmentation

本项目采用组织特征和职责分群，明确不推断个人年龄、性别、收入或教育背景。

| 分群维度 | 公开信号 | 对首个样本的含义 | 证据限制 |
| --- | --- | --- | --- |
| 组织规模 | ISC2 调查中，超过 10,000 人组织的已使用率为 37%；2,500–9,999 人和 100–499 人组织均为 33%；500–2,499 人和 1–99 人组织均为 20%。 | 规模不是线性采用预测变量；首批样本需同时记录安全团队规模、现有平台成熟度和治理资源。 | 436 人全球样本，自报，不代表中国 FOBrain 客户总体。 |
| 行业 | 同一调查中，工业企业 38%、IT 服务 36%、专业服务 34% 领先；金融 21%、公共部门 16% 较低。 | 行业监管强度不能直接推导采用意愿；需要分别验证数据、采购和审批约束。 | 行业样本基数未在摘要中完整公开。 |
| 岗位 | NIST 将企业补丁管理定义为识别、排序、获取、安装和验证，并指出业务／任务负责人和安全／技术管理者之间可能存在价值分歧。[NIST SP 800-40 Rev.4](https://www.nist.gov/publications/guide-enterprise-patch-management-planning-preventive-maintenance-technology) | 至少分开任务执行者、风险／业务责任人、平台／数据负责人和治理／审批者。 | 这是领域流程指导，不是目标客户岗位分布。 |

### Psychographic Profiles

这里的“心理特征”仅指与工作决策有关的风险态度，不推断个人性格：

1. **务实评估者**：愿意测试 AI，但要求在熟悉任务中看到速度、准确性或工作量改善。Microsoft 自有 RCT 对 147 名有经验的安全专业人员报告平均快 22%、准确性高 7%、97% 愿意在同类任务再次使用；这是特定产品和特定任务证据，只能证明“可通过任务实验测量”，不能外推本项目。[Microsoft RCT 摘要](https://www.microsoft.com/content/dam/microsoft/final/en-us/microsoft-brand/documents/EconStudy-Speed-accuracy-and-sentiment-Infographic-Final.pdf)
2. **证据与控制优先者**：关心来源、权限、错误边界、人工复核和责任。NIST AI RMF 要求定义具体任务、人机角色、监督、部署情境和测试证据；OWASP 把过度功能、权限和自主性列为 Excessive Agency 根因。[NIST AI RMF Core](https://airc.nist.gov/airmf-resources/airmf/5-sec-core/)、[OWASP LLM06:2025](https://genai.owasp.org/llmrisk/llm062025-excessive-agency/)
3. **资源受限采用者**：希望减少重复工作，但对额外平台、集成、培训和预算更敏感。行业调查显示技能和预算约束同时存在，不能假定“功能更多”就是更好。[ISC2 Workforce Study](https://www.isc2.org/Insights/2025/12/2025-ISC2-Cybersecurity-Workforce-Study)
4. **治理责任者**：需要知道谁批准、谁执行、数据去了哪里、失败如何停止和追责；他们可能阻止一个一线用户喜欢但缺乏治理证据的方案。[NIST AI RMF](https://www.nist.gov/itl/ai-risk-management-framework)

_Confidence：中等；这些是基于职责和权威治理框架的候选行为型，不是经目标客户验证的 persona。_

### Customer Segment Profiles

| 候选分群 | 主要工作结果 | 采用判断 | 不能由公开资料确定的事项 |
| --- | --- | --- | --- |
| 一线安全／漏洞运营人员 | 更快完成调查、优先级解释、交接或验证，同时保持事实正确 | 任务时间、返工、错误、证据可追溯性 | 哪一个任务最重要、频率和当前成本 |
| 资产／业务责任人 | 理解“为什么是我、为什么现在、应做什么、如何证明完成” | 业务语言是否减少澄清和往返 | 是否愿意进入新界面、真实结束证据 |
| 安全负责人／风险所有者 | 在有限资源下决定优先级并监督结果 | 风险结果、治理、预算、责任边界 | 经济买方、预算来源、采用承诺 |
| 平台／数据／安全管理员 | 保证身份、权限、数据、插件和审计边界 | 集成成本、最小权限、数据治理、可停用性 | FOBrain 目标部署的真实契约与限制 |
| 审计／法务／合规角色 | 复核依据、授权、数据处理与动作记录 | 是否满足具体组织与法律要求 | 目标行业、地域和数据下的实际适用性 |

这些分群只形成 DG-03 的候选队列，不是最终用户排序。CISA 的 NICE Framework 也明确 Work Role 是个人或团队负责并承担问责的一组工作，不等同于职位名称；因此后续样本必须按真实工作而不是 HR 职称纳入。[NICE Workforce Framework](https://niccs.cisa.gov/workforce-development/nice-framework)

### Behavior Drivers and Influences

- **理性驱动**：任务是否更快、更准、更少返工；是否可回到来源；是否减少跨系统检索；是否有可客观验证的结束证据。
- **风险驱动**：错误建议、数据泄漏、权限扩大、未经批准的动作和第二套事实系统会降低采用意愿。NIST 要求对部署相似条件进行评价并记录一般化限制；OWASP 建议高影响动作采用最小权限和人工批准。[NIST AI RMF Core](https://airc.nist.gov/airmf-resources/airmf/5-sec-core/)、[OWASP LLM06:2025](https://genai.owasp.org/llmrisk/llm062025-excessive-agency/)
- **组织影响**：ISC2 调查显示，组织在选择 AI 安全工具时主要依赖内部研究（72%）、培训教育（56%）、厂商建议（45%）和政府／公共政策（33%）。这支持“可验证试点 + 培训 + 治理材料”的采用路径，而不是仅靠产品演示。[ISC2 AI Adoption Pulse Survey](https://www.isc2.org/insights/2025/07/2025-isc2-ai-pulse-survey)
- **经济影响**：预算和人员约束会同时促进自动化需求并压低额外平台、座席和集成的可接受成本；商业模型必须在 DG-07 用真实预算记录判断。

### Customer Interaction Patterns

公开市场已经出现三种交互供给：现有安全产品内嵌助手、独立安全协作入口、通过插件／Agent 扩展的多工具平台。Microsoft 官方文档明确同时提供独立和嵌入式体验，并允许插件接入第三方服务；管理员可以控制插件可用范围。[Microsoft Security Copilot](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot)、[插件说明](https://learn.microsoft.com/en-us/copilot/security/plugin-overview)

这只能证明形态**存在**，不能证明用户偏好。合理的目标用户研究顺序是：重放最近一次真实任务 → 记录原有系统切换和证据 → 用同任务比较原生增强、内嵌、独立、非 AI 与不做 → 再观察持续使用。购买决策则应另外验证安全评审、平台管理员、产品负责人、经济买方和最终责任人。

证据支持的候选购买路径是：定义任务和当前损失 → 检查现有平台能力 → 冻结历史任务基线 → 只读／低风险试点 → 审查数据、身份、权限、模型、日志、供应链和合同 → 培训用户 → 按任务结果扩大、维持或停止。NIST CSF 2.0 允许组织用期望安全结果评价产品和服务，并强调从治理层到执行者及供应商的沟通；它支持多角色评审，但不规定本项目的具体采购流程。[NIST CSF 2.0](https://www.nist.gov/publications/nist-cybersecurity-framework-csf-20)

_Research and Discovery：内部研究、受控测试、培训和权威材料共同影响。_  
_Purchase Decision Process：一线价值、技术可行、安全／合规、预算与商业关系是不同门禁。_  
_Post-Purchase Behavior：当前没有本项目采用遥测；任何持续使用结论均为未知。_  
_Loyalty and Retention：只能由重复真实任务使用、资源承诺、续费／增购或替代支出证明。_

### Step 2 Quality Assessment

- **高置信度**：安全运营与漏洞管理存在多角色、跨系统、资源与技能约束；AI 安全工具处于采用／评估并存阶段；治理和人工责任必须显式定义。
- **中置信度**：效率、证据化输出和现有流程集成是重要采用驱动。
- **低置信度／未知**：本项目首个主用户、首个高价值任务、交互偏好、预算、持续采用和 FOBrain 样本能否外推到通用平台。
- **限制**：ISC2 为行业样本调查；Microsoft 结果属于厂商自有产品实验；两者均不能替代 DG-02 后的目标用户证据。
- **相互矛盾信号**：ISC2 报告采用者自报积极效果，而 SANS 报告 AI/ML 满意度下降；这要求本项目把“评估兴趣、已采用、满意、任务有效、持续使用”分别测量。

## Customer Pain Points and Needs

### Customer Challenges and Frustrations

公开资料显示的痛点集中在“任务系统无法协同”，而不是“缺少聊天框”：

1. **跨工具可见性与集成不足。** SANS 2026 SOC Survey 的公开摘要称，24% 的安全运营领导者把企业级可见性不足列为 SOC 有效性的最大单项障碍，并指出很多组织不是没有工具，而是工具之间缺少有效集成。[SANS 2026 SOC Survey 公告](https://www.sans.org/press/announcements/24-of-cyber-leaders-cite-lack-of-enterprise-wide-visibility-as-the-biggest-barrier-to-soc-effectiveness-the-2026-sans-soc-survey-finds)
2. **人员、自动化和上下文同时不足。** SANS 2024 调查中，缺少自动化／编排是最高单项障碍；将高人员需求和缺少技能人员合计后，人员相关问题更突出，同时还存在企业可见性和上下文不足。该调查不等同于漏洞运营用户总体，但证明这些问题经常共同出现。[SANS 2024 SOC Survey](https://files.abnormalsecurity.com/production/files/SANS-2024-SOC-Survey.pdf?dm=1720793408)
3. **漏洞排序不能只看单一严重度。** FIRST 明确 CVSS Base 描述系统无关的固有严重性；EPSS 只估计未来 30 天在野利用概率；CISA KEV 只表示已知在野利用。真实优先级还需要资产、暴露、业务影响、补偿控制和组织风险容忍度。[FIRST CVSS v4 Consumer Guide](https://www.first.org/cvss/v4.0/implementation-guide)、[FIRST EPSS User Guide](https://www.first.org/epss/user-guide)、[CISA KEV](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)
4. **从发现到修复和验证存在责任断点。** NIST 把企业补丁管理定义为识别、确定优先级、获取、安装和验证，并指出业务／任务负责人和安全／技术管理之间可能存在价值分歧；仅产生列表或工单不代表风险已经消除。[NIST SP 800-40 Rev.4](https://www.nist.gov/publications/guide-enterprise-patch-management-planning-preventive-maintenance-technology)
5. **AI 增加新的复核负担。** Google 官方文档直接提示生成式 AI 可能产生看似合理但事实错误的输出，并要求使用前验证；这意味着“生成更快”可能被核验时间抵消。[Google Gemini Security Command Center](https://docs.cloud.google.com/security-command-center/docs/gemini-overview)

KPMG 2024 对 200 名安全领导者的调查提供了交叉信号：数据质量／完整性和低质量告警疲劳各为 30%，29% 指出工具缺乏集成，32% 指出判断威胁和漏洞严重程度困难。该样本规模有限，但与 SANS、FIRST 和 CISA 的方向一致。[KPMG SOC Survey 2024](https://kpmg.com/kpmg-us/content/dam/kpmg/pdf/2024/time-to-transform-now.pdf)

_Primary Frustrations：工具分散、上下文不足、优先级解释困难、跨团队协调和结果验证。_  
_Frequency Analysis：行业资料证明问题存在，但不能证明 FOBrain 样本中的频率或成本。_  
_Confidence：问题类别高；本项目目标用户和优先级未知。_

### Unmet Customer Needs

以下只能作为待验证需求，不可直接进入 PRD：

| 候选未满足需要 | 客观结果 | 公开证据支持 | 必须补充的目标证据 |
| --- | --- | --- | --- |
| 统一任务上下文 | 在同一任务中看到对象、来源、时间、权限和缺口 | SANS 可见性／集成问题；多平台均提供数据聚合 | 最近真实任务的系统切换、耗时与遗漏 |
| 可解释的优先级 | 保留 CVSS、EPSS、KEV、资产和控制的独立语义 | FIRST、CISA、NIST | 用户是否需要解释、哪些事实改变决定 |
| 跨角色交接 | 安全人员、资产责任人和审批者共享可复核输入／输出 | NIST 多角色流程 | 真实责任链、等待、返工和结束证据 |
| 可验证的 AI 辅助 | 输出包含来源、未知、冲突和最小补证，而非只给答案 | NIST AI RMF、Google 错误提示 | 同任务正确率、复核时间和失败样本 |
| 受控动作 | 人员决定、审批、执行、回执和复核分开 | OWASP Excessive Agency、NIST AI RMF | DG-08 动作必要性、权限、补救和真实回执 |
| 可量化运营结果 | 用任务时间、错误、返工、积压或结果验证衡量 | SANS 指标与资源管理背景 | DG-05 预注册基线、目标、保护和停止线 |

真正的市场缺口不是“把所有数据放在一个页面”，而是能否在不建立第二套事实的前提下，把一个高价值任务从触发推进到可复核结束。该主张仍是研究综合，不是用户事实。

### Barriers to Adoption

- **价格／预算障碍**：SANS 2024 调查中 38.2% 的受访者不知道 SOC 预算；这表明一线痛点与预算控制可能分离，但不能推断本项目支付意愿。[SANS 2024 SOC Survey](https://files.abnormalsecurity.com/production/files/SANS-2024-SOC-Survey.pdf?dm=1720793408)
- **技术障碍**：多源实体不一致、身份与权限、数据新鲜度、接口覆盖、部署方式、日志和错误语义会造成集成与持续运维成本。
- **信任障碍**：生成错误、来源错配、不确定性隐藏、模型越权和无法复现会提高复核成本。NIST 要求在类似部署条件下测量并记录一般化限制；OWASP 建议限制功能、权限和自主性。[NIST AI RMF Core](https://airc.nist.gov/airmf-resources/airmf/5-sec-core/)、[OWASP LLM06:2025](https://genai.owasp.org/llmrisk/llm062025-excessive-agency/)
- **便利障碍**：独立 Workbench 会增加第二界面和上下文切换；内嵌形态可能受单一平台边界约束；两者必须用同任务对照判断。
- **治理障碍**：数据用途、存储地域、模型接收、保留删除、第三方供应链、人工监督和责任划分未明确时，正式试点不能开始。
- **治理成熟度差距**：英国政府 2025/2026 调查中，在使用、采用中或考虑 AI 的组织里，只有 24% 企业和 27% 慈善机构已有 AI 风险管理流程；31% 和 28% 明确没有实施计划。它不是安全运营产品专项调查，但直接支持“采用意愿不能替代治理准备度”。[UK Cyber Security Breaches Survey 2025/2026](https://www.gov.uk/government/statistics/cyber-security-breaches-survey-20252026/cyber-security-breaches-survey-20252026)
- **替代障碍**：FOBrain 原生配置、培训、保存视图、报表、脚本或非 AI 工作流可能以更低成本解决问题；若成立，应选择替代方案或停止。

### Service and Support Pain Points

安全产品的服务问题主要发生在实施后：连接器和权限配置、数据质量、规则与模型更新、培训、故障定位、动作失败和审计取证。Qualys、Rapid7、Microsoft、Palo Alto 等公开产品均强调平台、连接器或插件，这反向说明持续配置和生态支持是产品的一部分，而不是一次性交付。[Microsoft 插件说明](https://learn.microsoft.com/en-us/copilot/security/plugin-overview)、[Rapid7 Exposure Command](https://www.rapid7.com/products/command/exposure-management/)、[Qualys 10-K](https://www.sec.gov/Archives/edgar/data/1107843/000110784326000008/qlys-20251231.htm)

当前没有 FOBrain 目标客户的工单、实施周期、支持响应或升级记录，因此不能声称具体服务问题。DG-02 后需记录：谁维护连接、谁解释数据、谁处理错误、升级时谁负责、不可用时如何回退以及用户多快得到可行动答复。

实施支持本身有公开证据：ISC2 调查中只有 30% 的组织为 AI 安全工具培训配置预算，24% 依赖厂商教育；英国网络安全技能研究报告 AI 日常使用与正式培训覆盖存在差距。由此可把用例定义、效果评价、数据边界、培训和运行治理列为采用条件，但仍没有可靠独立证据证明“传统客服响应慢”是首要痛点。[ISC2 AI Pulse](https://www.isc2.org/insights/2025/07/2025-isc2-ai-pulse-survey)、[UK Cyber Security Skills 2025](https://assets.publishing.service.gov.uk/media/68be9b4c11b4ded2da19ff1b/Cyber_security_skills_in_the_UK_labour_market_2025-_findings_report.pdf)

### Customer Satisfaction Gaps

市场信号存在明确落差：ISC2 2025 调查中，已采用 AI 安全工具的受访者多数自报积极效果；SANS 2024 则记录 AI/ML 分析满意度下降、生成式 AI/ML 得分更低。差异可能来自样本、任务、产品、成熟度和评价方法，不能简单取平均。[ISC2 AI Pulse](https://www.isc2.org/insights/2025/07/2025-isc2-ai-pulse-survey)、[SANS 2024 SOC Survey](https://files.abnormalsecurity.com/production/files/SANS-2024-SOC-Survey.pdf?dm=1720793408)

因此必须分别评价：

1. 是否愿意评估；
2. 是否完成部署；
3. 是否在特定任务中更快／更准；
4. 是否减少总复核与返工；
5. 是否持续使用；
6. 是否愿意投入资源或付费。

任何一项都不能替代下一项。

### Emotional Impact Assessment

公开调查把安全运营与压力、人员不足、工具维护和离职意向关联起来。Splunk 2025 厂商调查称 46% 受访者花更多时间维护工具而非防御，52% 表示工作压力曾使其考虑离开网络安全；这是厂商调查，只能说明潜在影响，不能用作 FOBrain 用户 Background。[Splunk State of Security 2025 公告](https://newsroom.cisco.com/c/r/newsroom/en/us/a/y2025/m05/global-state-of-security-report-reveals-critical-need-for-connected-security-operations.html)

ISC2 2025 的 16,029 人研究提供更广的独立信号：48% 因持续追踪新威胁和技术感到疲惫，47% 经常被工作量压垮，32% 因人员或技能短缺而过度工作。它证明工作负担问题存在，但不能证明 AI 会减负；AI 也可能增加验证、监督和异常处理。[ISC2 Workforce Study 2025](https://www.isc2.org/Insights/2025/12/2025-ISC2-Cybersecurity-Workforce-Study)

产品研究应避免用“焦虑”“疲劳”等情绪词代替可观察事实。真正可进入需求的证据应是任务积压、等待、返工、遗漏、错误、加班、升级和人员流失记录，并确认其与候选任务存在因果关系。

### Pain Point Prioritization

公开资料只能形成候选优先级：

| 候选级别 | 痛点 | 进入 DG-04 前的最低证据 |
| --- | --- | --- |
| 高候选 | 跨工具上下文重建、优先级证据组合、责任交接、结果验证 | 近期真实任务；可测时间／返工／错误或高损失事件；现有方案差距 |
| 中候选 | 报告草拟、自然语言查询、知识辅助、建议生成 | 证明不是仅“方便”，且总复核时间下降、准确性不降低 |
| 低候选 | 新聊天入口、通用平台、全部 24 个历史工具、自动高风险处置 | 当前无任务／采用证据；不得进入首版范围 |

最终只允许 DG-04 根据冻结样本和预注册 A/B 规则选择一个任务，也允许“没有合格任务”。

### Step 3 Quality Assessment

- **高置信度**：可见性、集成、人员／技能、上下文、优先级、跨团队处置和 AI 可信度是行业级问题类别。
- **中置信度**：统一任务上下文和证据化 AI 可能改善其中一部分工作。
- **未知**：FOBrain 样本的具体问题、频率、损失、替代方案、预算和产品形态。
- **反证门禁**：若原生配置、培训、非 AI 改进或现有平台已经解决任务，停止对应产品机会。

---

## Customer Decision Processes and Journey

### Customer Decision-Making Processes

当前公开资料不能直接观察 FOBrain 目标客户的采购流程，因此以下是由企业安全采购、AI 安全工具评估和成熟厂商销售结构综合出的**待验证决策链**，不是已验证的目标客户事实：

1. **确认问题与责任人**：用一个真实任务说明当前耗时、错误、返工、等待或损失，并确认使用者、风险所有者和预算影响者。
2. **内部研究与替代核查**：检查 FOBrain 原生能力、现有安全平台、流程配置、脚本、培训、MSSP 和“不做”是否已足够。
3. **形成可比较的评价条件**：冻结任务样本、基线、目标、保护指标、数据边界和停止规则，而不是比较演示效果。
4. **低风险试点**：优先只读、建议或可回退场景，验证任务时间、准确率、复核成本和采用行为。
5. **技术与供应商尽调**：审查数据、身份、权限、部署、模型、日志、供应链、合同、退出和事故响应。
6. **联合决策**：一线使用者评价可用性，安全／业务负责人评价结果，平台与数据负责人评价集成，合规／采购评价风险和合同，经济买方决定资源。
7. **扩大、维持或停止**：只有净任务收益、治理条件和持续采用同时成立才扩大；否则回到原生增强、非 AI 改进或停止。

ISC2 2025 的 436 人调查显示，组织评估 AI 安全工具主要依靠内部研究（72%）、培训教育（56%）、厂商建议（45%）和政府／公共政策信息（33%），说明决策不是单一销售触点。[ISC2 AI Pulse Survey](https://www.isc2.org/insights/2025/07/2025-isc2-ai-pulse-survey) NIST CSF 2.0 则明确可用于比较产品／服务提供商并沟通供应商要求。[NIST CSF FAQ](https://www.nist.gov/cyberframework/faqs)、[NIST SP 1305](https://csrc.nist.gov/pubs/sp/1305/final)

_Decision Stages：问题确认 → 内部研究 → 替代核查 → 任务级试点 → 尽调 → 联合决策 → 扩大／维持／停止。_  
_Decision Timelines：目标客户未知；成熟厂商以年度订阅为主不能证明本项目销售周期。_  
_Complexity Levels：中高；至少涉及任务、数据、身份、治理、采购和持续运营。_  
_Evaluation Methods：冻结同任务样本，比较总耗时、正确性、遗漏、返工、复核时间、治理风险和总成本。_  
_Confidence：决策结构中等；具体角色、顺序、周期和签字权低，必须由 DG-02 后真实访谈／观察验证。_

### Decision Factors and Criteria

| 因素 | 应验证的问题 | 通过信号 | 延迟／否决信号 |
| --- | --- | --- | --- |
| 任务净收益 | 是否减少一个任务的总成本，而非只缩短生成时间 | 总时间、返工或损失下降且保护指标不恶化 | 复核成本抵消生成收益 |
| 事实与可追溯性 | 输出能否回到权威来源、时间和权限上下文 | 结论、缺口与来源可客观复核 | 第二套事实、来源错配或无法复现 |
| 数据与部署 | 数据是否允许进入候选模型／服务路径 | 数据流、存储、保留、删除和责任明确 | 敏感数据、跨境或内网边界不清 |
| 身份与动作 | 是否继承原权限并限制自主性 | 最小权限、审批、审计、失败补救明确 | 越权、不可回退或责任不清 |
| 集成与转换 | 是否必须新增界面、连接器和运维 | 单任务所需最小集成可持续运行 | 第二界面和维护成本高于价值 |
| 供应商与合同 | 是否满足安全、服务、退出和持续披露要求 | SLA、漏洞披露、事故响应和退出可约定 | 锁定、责任转嫁或供应链不可见 |
| 采用与培训 | 真实用户是否持续使用且能正确监督 | 无强制情况下持续使用，培训成本可接受 | 演示认可但实际绕开、误用或弃用 |
| 经济性 | 谁承担预算，价值是否超过采购与运营总成本 | 经济买方、预算科目和续用条件明确 | 痛点方无预算，或现有方案更便宜 |

NIST CSF 2.0 要求供应商尽调与关系风险、关键性和复杂度相称，并可在合同中明确安全、SLA、信息共享和持续披露要求。[NIST CSF 2.0 Reference Tool](https://csrc.nist.gov/Projects/Cybersecurity-Framework/Filters) 英国政府 2025/2026 调查显示，整体只有 22% 企业在购买软件时“很大程度”考虑网络安全，但大型企业为 69%；审查直接供应商风险的比例整体为 15%、大型企业为 48%，要求供应商认证的比例整体为 11%、大型企业为 41%。这说明组织规模与治理成熟度会改变权重，不能把大型企业尽调流程外推到所有客户。[UK Cyber Security Breaches Survey 2025/2026](https://www.gov.uk/government/statistics/cyber-security-breaches-survey-20252026/cyber-security-breaches-survey-20252026)

_Primary Decision Factors：同任务净收益、事实正确、数据／权限边界、现有方案差距。_  
_Secondary Decision Factors：产品形态、品牌、生态、培训、渠道和合同期限。_  
_Weighing Analysis：高影响行业／大组织更可能提高供应链、认证与治理权重；资源受限组织更可能提高部署与总成本权重。_  
_Evolution Patterns：试点前重可行性和风险；试点中重任务效果；采购与续用阶段重治理、运营、总成本和持续采用。_

### Customer Journey Mapping

| 阶段 | 客户要完成的决定 | 必要证据／触点 | 本项目门禁 |
| --- | --- | --- | --- |
| Awareness | 是否存在值得解决的剩余任务 | 真实事件、积压、错误、损失、使用者陈述 | DG-02／DG-03 |
| Consideration | AI、非 AI、原生增强、独立入口或不做哪个更合理 | 现状流程、替代方案、数据与部署可行性 | DG-04／DG-06 |
| Trial | 候选方案是否在同任务中产生净收益 | 冻结样本、预注册指标、失败与停止记录 | DG-05／DG-06 |
| Decision | 组织是否愿意承担治理、集成和预算 | 安全评审、供应商尽调、经济买方与合同 | DG-07／DG-09 |
| Purchase／Adoption | 是否能在真实环境受控运行 | 权限、培训、SLA、支持、退出和回退 | DG-08／DG-09 |
| Post-Purchase | 是否持续产生结果，是否扩大或停止 | 持续采用、质量、事故、成本和续用证据 | 后续产品门禁；不在当前公开研究中预判 |

成熟安全厂商的公开文件说明市场常见年度订阅和直销／渠道／MSSP 并行：Tenable 主要按一年订阅销售，并通过直销、分销商、经销商和 MSSP 获取机会；Qualys 2025 年 49% 收入经渠道伙伴产生。这只能证明成熟市场触达结构存在，不能证明本项目应采用相同商业模式。[Tenable 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1660280/000166028026000005/tenb-20251231.htm)、[Qualys 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1107843/000110784326000008/qlys-20251231.htm)

_Awareness Stage：从真实任务问题开始，不从“需要 AI”开始。_  
_Consideration Stage：保留原生增强、嵌入式 AI、独立 Workbench、非 AI、托管服务和不做。_  
_Decision Stage：任务效果、治理可行性、采用和经济买方必须同时成立。_  
_Purchase Stage：合同与部署只是采用开始，不是产品验证完成。_  
_Post-Purchase Stage：用持续任务结果、事故、成本和实际使用决定续用／扩大／停止。_

### Touchpoint Analysis

- **数字触点**：内部检索、产品文档、安全与隐私说明、标准／监管资料、技术演示、试点环境、支持门户和任务结果仪表。
- **组织内触点**：一线用户、任务负责人、CISO／风险负责人、平台／数据／身份管理员、架构、安全评审、法务、隐私、采购和经济买方。
- **外部触点**：厂商售前／交付、渠道／集成商、MSSP、同行、专业组织、监管与标准机构、独立测试和客户参考。
- **高价值触点**：真实任务回放、可复现实验、数据流图、权限矩阵、失败演示、退出／回退方案和合同安全条款。

ISC2 的信息来源分布证明内部研究、教育、厂商和公共政策共同影响判断；NIST SP 1305 进一步把供应商要求和生命周期沟通列入采购能力。[ISC2 AI Pulse Survey](https://www.isc2.org/insights/2025/07/2025-isc2-ai-pulse-survey)、[NIST SP 1305](https://csrc.nist.gov/pubs/sp/1305/final)

_Digital Touchpoints：官方文档、独立研究、试点与结果记录。_  
_Offline Touchpoints：工作观察、跨角色评审、采购／合规会议和设计伙伴回顾。_  
_Information Sources：内部事实优先，其次是权威标准／监管、独立研究、厂商资料。_  
_Influence Channels：真实任务负责人和治理否决者的影响高于通用营销内容。_

### Information Gathering Patterns

客户需要收集的信息至少分为六包：问题与基线、候选方案、数据／权限、任务效果、供应商／合同、持续运营。每包都必须标明来源、日期、适用范围和未知项。信息优先级应为：目标组织真实记录与观察 > 官方合同／技术证据 > 权威监管与标准 > 独立研究 > 厂商声明 > 社区线索。

公开资料不能给出目标客户的研究时长。ISC2 只能证明多来源尽调普遍存在；UK 调查同时显示不少组织很少正式审查供应商，因此研究深度必须由组织规模、行业、数据敏感度和任务风险验证，不能预设统一流程。[ISC2 AI Pulse Survey](https://www.isc2.org/insights/2025/07/2025-isc2-ai-pulse-survey)、[UK Cyber Security Breaches Survey 2025/2026](https://www.gov.uk/government/statistics/cyber-security-breaches-survey-20252026/cyber-security-breaches-survey-20252026)

_Research Methods：任务访谈、现场观察、历史样本回放、方案对照、供应商尽调和受控试点。_  
_Information Sources Trusted：目标组织事实、官方监管／标准、独立验证；厂商资料只证明供给。_  
_Research Duration：未知；DG-02 后记录每一阶段的实际等待和返工。_  
_Evaluation Criteria：证据质量、适用性、可复现性、更新时间和是否存在利益冲突。_

### Decision Influencers

企业安全采购不存在消费品式“亲友影响”。应按职责记录影响与否决权：

| 影响者 | 关注点 | 可能的否决原因 |
| --- | --- | --- |
| 一线使用者 | 任务速度、准确性、可解释和额外负担 | 复核更慢、改变现有流程、错误不可定位 |
| 任务／业务负责人 | 风险结果、SLA、责任和中断 | 无法证明结果或动作风险过高 |
| 平台／数据／身份负责人 | 集成、权限、数据质量和运维 | 接口不稳定、越权、数据外流、维护不可持续 |
| CISO／治理角色 | 安全、审计、供应链、事故与责任 | 治理不可证明或第三方风险不可接受 |
| 法务／隐私／采购 | 法律适用、合同、保留、退出和价格 | 条款、数据路径、锁定或责任不清 |
| 经济买方 | 净收益、优先级、预算和替代成本 | 无预算、价值归属不清或现有方案更优 |
| 同行／专家／标准机构 | 风险框架、实践与可信度 | 只能影响判断，不能替代目标任务证据 |

_Peer Influence：同岗位、设计伙伴和客户参考可降低感知风险，但不能替代本地验证。_  
_Expert Influence：架构、安全、隐私、监管和标准专家可形成条件与否决项。_  
_Media Influence：用于发现候选方案，证据权重低。_  
_Social Proof Influence：客户案例可支持进入评估，不能证明目标部署效果。_

### Purchase Decision Factors

**立即进入试点的驱动**应是：一个责任明确且损失可量化的任务、现有方案仍有缺口、可取得合规样本、能先只读验证、结果指标和停止规则明确。**造成延迟或停止的因素**包括：任务不重要、样本不可用、原生／非 AI 方案已足够、数据或权限边界未批准、无法形成同任务对照、复核成本过高、买方缺失、合同／部署风险和持续运营责任不清。

成熟厂商的年度或多年期订阅只说明行业常见合同形态；本项目在 DG-07 前不得据此确定订阅、按资产、按用量、附加模块或服务模式。[Tenable 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1660280/000166028026000005/tenb-20251231.htm)

_Immediate Purchase Drivers：高价值任务、同任务证据、治理可行、经济买方和可控部署。_  
_Delayed Purchase Drivers：不清楚的任务／数据／责任、复杂集成、供应链风险、预算和替代方案。_  
_Brand Loyalty Factors：历史效果、支持、集成和转换成本；不能把品牌当作正确性证明。_  
_Price Sensitivity：必须比较采购、集成、模型、培训、监督、支持、退出和机会成本的总和。_

### Customer Decision Optimizations

本项目不应优化“尽快成交”，而应优化“尽快形成正确的继续／停止决定”：

1. 用一页任务证据卡固定问题、角色、输入、输出、损失和当前方案。
2. 公开保留 AI、非 AI、原生、独立和不做，不将产品形态写入问题定义。
3. 用冻结历史样本演示正确、错误、缺失和回退，不只展示 happy path。
4. 提前提供数据流、权限、模型、日志、保留、跨境、退出和事故处理材料。
5. 把一线采用、治理通过和经济购买分成三个独立信号。
6. 预先定义扩大、维持、转向和停止条件，避免试点因沉没成本自动延长。

_Friction Reduction：一次收集可复用的任务、数据、供应商和合同证据。_  
_Trust Building：来源可追溯、限制公开、权限最小、失败可见、动作可审批／回退。_  
_Conversion Optimization：只优化合格问题进入试点的速度，不把不合格线索转成需求。_  
_Loyalty Building：持续证明任务结果与服务可靠性，而不是增加功能数量。_

### Step 4 Quality Assessment

- **高置信度**：企业 AI 安全工具评估依赖多来源研究；供应商风险、数据、权限、合同和培训会影响采用；大型组织通常进行更严格尽调。
- **中置信度**：问题／基线—替代核查—低风险试点—联合尽调—扩大／停止是合理的候选决策链。
- **低置信度／未知**：FOBrain 目标客户的实际角色、顺序、销售周期、预算、签字权、渠道和续用条件。
- **门禁结论**：Step 4 只能定义要验证的决策旅程，不能证明客户会购买，也不能确定 FOBrain 内嵌、独立 Workbench、通用平台或任何商业模式。

---

## Competitive Landscape

### Key Market Players

竞争不是一个统一的“AI 安全助手市场”，而是多类供给争夺同一安全任务：

| 竞争层 | 代表玩家／方案 | 结构性优势 | 对本项目的含义 |
| --- | --- | --- | --- |
| FOBrain 原生增强 | FOBrain | 首个样本已有资产、漏洞、优先级和闭环事实 | 必须先证明原生配置／增强仍无法解决任务，不能重复包装现有价值。[FOBrain](https://huashunxinan.net/product-fobrain) |
| 漏洞／Exposure 平台 | Tenable、Qualys、Rapid7 | 传感器、资产基础、研究数据、连接器、优先级、攻击路径和修复工作流 | “多源聚合、风险排序、攻击图、AI”均已商品化。[Tenable One](https://www.tenable.com/products/tenable-one)、[Qualys ETM](https://docs.qualys.com/en/etm/latest/about_etm.htm)、[Rapid7 Exposure Command](https://www.rapid7.com/products/command/exposure-management/) |
| 综合 SecOps 平台 | Microsoft、Google、Palo Alto Networks、CrowdStrike | SIEM／XDR／SOAR／端点／身份数据、原生权限、渠道与合同 | 已同时提供嵌入式、独立和 Agent 形态；单一 AI 入口容易被复制。[Microsoft Security Copilot](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot)、[Google SecOps](https://cloud.google.com/security/products/security-operations/investigate)、[Cortex AgentiX](https://docs-cortex.paloaltonetworks.com/r/Cortex-AgentiX/Cortex-AgentiX-Documentation/Get-Started-with-Cortex-AgentiX)、[CrowdStrike Charlotte AI](https://www.crowdstrike.com/en-us/platform/charlotte-ai/) |
| ITSM／企业工作流 | ServiceNow／Armis | 工单、审批、资产、变更、企业责任链和跨部门分发 | 发现到责任协作的价值可能由工作流平台吸收。[ServiceNow／Armis](https://newsroom.servicenow.com/press-releases/details/2026/ServiceNow-completes-Armis-acquisition-closing-the-gap-between-asset-visibility-and-cyber-risk/default.aspx) |
| 托管服务 | MSSP／MDR／集成商 | 产品、人员、流程和 SLA 一起交付 | 对资源不足客户可能比新软件更具替代性。 |
| 非 AI／内部方案 | 现有页面、报表、脚本、规则、培训和人工流程 | 可预测、低治理成本、贴合现有权限责任 | 必须作为最低成本对照，而不是“落后方案”。 |
| 通用 Agent 平台 | Agent builder、插件／MCP 生态 | 工具扩展和工作流编排速度 | 缺少安全语义、权威事实、买方和责任边界；首个样本不能证明平台机会。 |

Tenable 和 Rapid7 的 2025 年报均把竞争描述为快速变化且会因 AI 加剧；Tenable 明确列出 Qualys、Rapid7、综合安全厂商、CrowdStrike、Palo Alto、云厂商和点工具。这支持“跨类别竞争”，但属于卖方披露，不是独立市场排名。[Tenable 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1660280/000166028026000005/tenb-20251231.htm)、[Rapid7 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1560327/000156032726000008/rp-20251231.htm)

### Market Share Analysis

当前不存在公开、统一、可审计的“Exposure Management × AI Security Operations × Agent Workbench”市场份额。原因包括：

- 类别边界重叠，厂商把漏洞、云、身份、SIEM、SOAR、AI、服务和平台收入混合披露；
- 同一客户可能同时使用多家产品，收入、资产、席位和数据摄入无法直接换算；
- 公开“领导者”或份额主张常来自厂商引用的付费报告，底层方法和数据不可核验；
- FOBrain 的目标客户、经济买方和产品类别尚未确定，任何份额分母都不成立。

可核验的公司收入只能说明商业规模，不能当作目标细分份额：Qualys 2025 年收入约 6.691 亿美元，Rapid7 约 8.598 亿美元；两家公司均跨多个产品类别。[Qualys 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1107843/000110784326000008/qlys-20251231.htm)、[Rapid7 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1560327/000156032726000008/rp-20251231.htm)

_Market Share：未知，不生成伪精确排名。_  
_Position Proxy：只记录产品范围、存量数据、渠道、收入和客户结构等可验证指标。_  
_Re-estimation Gate：DG-07 确认产品类别、买方、收费边界和目标客户后再估算 TAM／SAM／SOM。_

### Competitive Positioning

| 候选定位 | 已有竞争优势 | 本项目潜在切口 | 必须证伪的问题 |
| --- | --- | --- | --- |
| FOBrain 原生助手 | 零／低额外分发，继承原事实和权限 | 降低单任务操作与解释成本 | FOBrain 原生增强是否已经足够？ |
| 跨源独立 Workbench | 可跨系统组织任务和证据 | 若任务天然跨 FOBrain、ITSM、身份或其他来源，可能减少切换 | 第二界面、第二套事实和集成成本是否超过收益？ |
| 安全运营 AI 能力层 | 可被多个产品嵌入 | 若同一能力在多个任务／客户重复出现，可能形成复用 | 是否存在共同买方、共同事实与共同验收？ |
| 通用 Agent 平台 | 扩展到任意工具和流程 | 只有多个独立领域样本才能支持 | 安全领域证据如何外推到通用需求？ |
| 托管／混合服务 | 人与产品共同承担结果 | 对技能不足客户降低落地门槛 | 是否有服务交付能力、毛利和责任边界？ |
| 非 AI 流程改进 | 成本、确定性和合规优势 | 可能以最低成本解决首个任务 | 若同样有效，为什么需要 AI 产品？ |

当前没有任何候选定位达到 `VALIDATED_POSITIONING`。可接受的临时表达只有：**以 FOBrain 为首个样本，寻找一个原生流程之后仍存在、可客观复核且值得解决的任务；最终关系和形态由 DG-04–DG-07 决定。**

### Strengths and Weaknesses

以下 SWOT 评估的是“从 FOBrain 样本开始做证据化任务协作”这一路径，而不是尚未成立的产品：

| 维度 | 内容 | 证据状态 |
| --- | --- | --- |
| Strengths | 有现成业务样本、历史文档、FOBrain 产品语境和可建立任务回放的工程基础 | 项目／公开证据；真实访问仍待 DG-02 |
| Strengths | 研究已把 AI、原生、独立、非 AI 和不做放在同一门禁比较 | 已建立研究治理，不代表客户价值 |
| Weaknesses | 无目标用户访谈、任务样本、采购记录、预算、采用遥测和真实部署数据 | 已确认缺口 |
| Weaknesses | FOBrain 已声称多源归一、动态优先级和闭环，差异空间被压缩 | 厂商供给证据，目标部署待核实 |
| Opportunities | 某个跨角色／跨系统任务若仍有高频成本或低频高损失，可能形成任务切口 | 假设，待 DG-03／DG-04 |
| Opportunities | 中国企业内部、私有／本地、安全审计与受控动作的组合可能影响采用 | 监管与供给信号，需求待验证 |
| Threats | 成熟平台用原生数据、权限、Agent、图谱和渠道吸收表层功能 | 高置信度供给趋势 |
| Threats | 非 AI 配置、流程优化或托管服务以更低成本解决任务 | 必须保留的反证 |
| Threats | 数据、漏洞、跨境、CII 或产品认证约束阻断试点／商业化 | 条件性监管风险 |

### Market Differentiation

研究排除以下伪差异：聊天界面、接入工具数量、Agent 数量、通用模型、攻击图、自然语言查询、自动生成报告、平台标签或“端到端”口号。成熟供应商已公开覆盖其中多数能力。

只有以下差异候选值得进入真实验证：

1. **任务结果差异**：在同一真实任务中稳定减少总耗时、错误、遗漏、返工或损失，而不是只缩短生成时间。
2. **事实一致性差异**：不建立第二套资产、漏洞、风险、工单或状态，所有结论能回到权威来源。
3. **验证成本差异**：来源、冲突、未知和适用范围清晰，使人工复核净减少。
4. **控制与责任差异**：继承原身份／权限，高影响动作可预览、审批、审计、回退和证明完成。
5. **部署适配差异**：若目标客户确有内网、境内、私有或断连要求，能以可接受成本满足；不能只把私有部署当口号。

这些不是定位声明，而是 DG-05／DG-06 的可测假设。若成熟平台或非 AI 方案同样满足，应停止对应差异主张。

### Competitive Threats

- **原生平台挤压**：Microsoft、Google、Palo Alto、CrowdStrike、Tenable、Qualys、Rapid7 和 ServiceNow 均可用存量数据与分发复制表层 AI 功能。
- **FOBrain 自身替代**：首个任务可能通过配置、培训、报表、规则或原生增强解决。
- **类别拥挤**：Exposure、CTEM、ASM、AI SOC、Agentic SOC 和安全 Workbench 术语重叠，客户难以理解新增购买理由。
- **集成税**：连接器、实体一致性、权限、数据新鲜度、版本和支持成本可能持续吞噬收益。
- **信任与责任**：错误排序、虚构依据、越权动作和第二套事实会造成比手工流程更高的风险。
- **渠道和转换成本**：成熟平台已有年度／多年合同、代理、数据历史、培训和渠道，新增独立产品需要额外评审。
- **监管分化**：数据、漏洞、跨境、公众／内部、CII 和行业规则会限制统一 SaaS 形态。
- **预算错位**：受痛点影响的一线用户可能没有预算，经济买方可能更重视平台整合或服务结果。

### Opportunities

| 机会假设 | 进入条件 | 停止条件 |
| --- | --- | --- |
| FOBrain 单任务辅助 | 最近真实任务存在可测剩余损失；原生／非 AI 仍不足 | 无合格任务或净收益不成立 |
| 跨源任务协作 | 完成任务必须跨至少两个权威系统，且切换／交接是主要损失 | 单源原生能力即可解决，或第二界面成本更高 |
| 受控动作闭环 | 建议到执行之间存在高损失等待，且动作可回退／审计 | 无必要动作、风险不可接受或真实回执不可得 |
| 内部私有部署方案 | 目标客户明确要求境内／本地／断连且愿承担成本 | 云／原生方案已满足，或本地运维成本超过价值 |
| 多客户／多工具产品 | 多个独立样本出现共同任务、共同买方、共同事实与验收 | 需求只属于 FOBrain 或单一客户定制 |
| 托管／混合交付 | 客户缺技能但愿意为结果与 SLA 付费 | 无服务能力、责任不可承受或经济模型不成立 |

市场进入应遵循“先验证、后定位”。美国 SBA 的官方指南区分公开资料与直接客户研究：公开资料适合回答一般、可量化问题，直接访谈／调查用于理解特定客户；竞争分析还应覆盖直接、间接竞争、进入窗口和壁垒。[SBA Market Research and Competitive Analysis](https://www.sba.gov/business-guide/plan-your-business/market-research-competitive-analysis) 这与当前 DG-02 后先取得目标任务证据、再决定产品的顺序一致。

### Step 5 Quality Assessment

- **高置信度**：竞争跨越原生平台、Exposure、SecOps、ITSM、服务、内部方案和通用 Agent；成熟厂商已覆盖多源、图谱、AI、工作流和 Agent。
- **中置信度**：若存在跨角色／跨系统且可测的剩余任务，“证据化协作且不建立第二套事实”可能形成差异。
- **低置信度／未知**：目标细分、市场份额、经济买方、定价、渠道、FOBrain 原生差距及独立产品可赢性。
- **战略结论**：现在不能“进入一个市场”；只能进入 DG-02–DG-06 的问题验证。任何 GTM、定价或平台扩展都必须等待 DG-07。

---

## Research Synthesis and Strategic Recommendations

### 1. Market Research Introduction and Methodology

本研究回答的不是“安全 AI 市场是否增长”，而是“当前项目是否有足够证据定义一个可开发、可购买的产品”。研究范围覆盖客户行为、行业痛点、替代方案、决策旅程、竞争、监管交叉条件和市场进入门禁；时间重点为 2024–2026 年，地理上以全球供给和研究信号为背景、以中国企业部署约束为首个验证语境。

资料按证据等级处理：目标组织事实与真实任务最高，其次是法律／监管、权威标准、独立研究、上市公司披露、厂商功能声明和社区线索。当前只有后五类公开资料，**没有目标客户访谈、任务记录、采购数据、采用遥测或获批真实样本**。SBA 官方指南也区分公开资料与直接客户研究：前者适合一般市场问题，后者用于理解特定客户。[SBA Market Research Guide](https://www.sba.gov/business-guide/plan-your-business/market-research-competitive-analysis)

**研究目标完成度：**

| 原始目标 | 公开研究结果 | 状态 |
| --- | --- | --- |
| 验证 FOBrain 是否存在真实高价值剩余任务 | 形成候选问题、样本规则和停止条件；没有真实任务证据 | 未验证，进入 DG-02 |
| 分析用户、替代、竞争、采用和购买 | 已形成行业级角色、决策链、竞争类别和采用条件 | 公开研究完成；目标客户适用性未知 |
| 判断是否扩展为多工具／独立／通用平台 | 建立所需外推证据和反证；没有任何形态胜出 | 保持开放 |

### 2. Market Analysis and Dynamics

市场驱动来自漏洞利用速度、攻击面复杂度、技能／预算约束和跨系统协调。Verizon 2026 DBIR 报告漏洞利用占事件入口 31%，第三方相关事件上升，说明修复速度和供应链事实持续重要；它不能证明本产品需求。[Verizon 2026 DBIR](https://www.verizon.com/about/news/breach-industry-wide-dbir-finds)

市场分类同时包含 Vulnerability Management、ASM、Exposure／CTEM、SecOps、AI SOC、ITSM 和服务，公开估值口径相差一个数量级以上。故：

- `TAM = UNKNOWN`，DG-07 前禁止选择单一商业报告数字；
- 市场增长不等于可进入机会；
- 真正分析单位是“角色 × 任务 × 组织环境 × 决策权”，不是泛化的安全团队；
- 商业模式可能是订阅、平台附加、用量、托管或内部能力，目前均未验证。

### 3. Customer Insights and Behavior

行业证据支持的行为是“积极评估、谨慎投产”：客户先内部研究、培训和厂商／政策咨询，再通过试点、安全评审、合同和预算决定。目标角色至少包括一线运营、任务／业务责任人、平台／数据／身份负责人、CISO／风险、法务／隐私／采购和经济买方。

候选痛点包括跨工具上下文重建、优先级证据组合、跨角色交接、结果验证和治理；但它们只能进入访谈和样本观察，不能直接进入 PRD。客户旅程必须从真实问题和基线开始，经替代核查、受控试点、尽调和联合决策，最终允许扩大、维持、转向或停止。

### 4. Competitive Landscape and Positioning

竞争优势主要来自存量事实、身份权限、工作流、渠道和持续运营，不来自模型或界面。FOBrain 已公开声明多源资产／漏洞归一、VPT 动态优先级和闭环；Tenable、Qualys、Rapid7、Microsoft、Google、Palo Alto、CrowdStrike 与 ServiceNow 又覆盖 Exposure、图谱、SecOps、Agent 和修复工作流。因此以下不能作为差异：

- “接入多工具”；
- “AI 对话／自然语言查询”；
- “攻击图／风险评分”；
- “Agent 自动化”；
- “统一平台／端到端”；
- “私有部署”但没有目标客户和成本证据。

唯一可检验的差异候选是：在不建立第二套事实的前提下，使一个真实任务的总耗时、错误、遗漏、返工或损失显著下降，同时降低而非增加复核与治理成本。

### 5. Strategic Market Recommendations

**当前机会判断：** 不进入正式市场，不发布定位，不承诺产品形态；进入一轮有停止条件的任务发现。

**建议顺序：**

1. DG-02 批准研究人员、数据字段、环境、访问、保留和禁止项。
2. DG-03 形成覆盖最近成功、失败、边界和替代方案的冻结样本。
3. DG-04 只选一个高频效率型或低频高损失型任务，也允许无合格任务。
4. DG-05 预注册结果、基线、目标、保护指标和停止规则。
5. DG-06 比较原生、嵌入式、独立、非 AI 和不做。
6. DG-07 只有采用、买方和资源承诺成立时才形成定位与商业假设。

### 6. Market Entry and Growth Strategy

市场进入被重新定义为“设计伙伴验证”，不是销售扩张：首个伙伴只提供一个任务样本；不要求购买，也不把试点等同采用。渠道、定价和规模化必须等待共同任务、共同买方、共同事实和共同验收在多个独立样本中重复出现。

| 阶段 | 进入条件 | 输出 | 停止条件 |
| --- | --- | --- | --- |
| 任务发现 | DG-02 数据／人员获批 | 任务卡、基线、替代、损失 | 无合格任务 |
| 方案对照 | 样本和指标冻结 | 原生／AI／非 AI／不做结果 | 净收益不成立 |
| 采用验证 | 任务有效 | 无强制情况下的持续使用与资源承诺 | 用户绕开、误用或弃用 |
| 买方验证 | 采用成立 | 经济买方、预算、合同与续用条件 | 无买方或总成本不可接受 |
| 多样本扩展 | 多个独立样本重复 | 细分、定位、商业模式候选 | 仅单客户定制 |

### 7. Risk Assessment and Mitigation

风险按 ISO 31000 的识别、分析、评价、处置、监测和沟通思路管理；ISO 31000 是指导标准，不能用于认证。[ISO 31000:2018](https://www.iso.org/standard/65694.html)

| 风险 | 领先指标 | 处置／备用路径 |
| --- | --- | --- |
| 伪需求 | 只有愿望，无最近任务、损失或替代 | 停止，不生成需求 |
| 原生替代 | 配置／培训即可解决 | 保留为 FOBrain 增强或结束 |
| AI 净收益为负 | 复核、返工和错误抵消速度 | 改用非 AI 或停止 |
| 平台挤压 | 竞品已有同任务原生能力 | 比较迁移／集成总成本，不靠功能数量竞争 |
| 第二套事实 | 状态、优先级或对象与源系统冲突 | 所有产品输出回到权威事实；独立入口失败则内嵌 |
| 无采用／买方 | 演示认可但无持续使用或预算 | 不进入定位与商业化 |
| 数据／监管阻断 | 样本、跨境、漏洞、CII 或权限不获批 | 使用脱敏／合成数据仅做可行性；不冒充任务验证 |

### 8. Evidence Roadmap and Success Metrics

研究阶段不使用 DAU、消息数或模型得分作为产品成功。最小指标组为：

- **任务结果**：完成率、总时间、正确性、遗漏、错误、返工、等待和损失；
- **保护指标**：来源正确、权限一致、隐私／数据边界、不可逆动作、事故和回退；
- **采用**：代表性用户的重复真实使用、覆盖／撤销和培训／支持负担；
- **经济**：采购、集成、推理、监督、运维和退出的总成本；
- **决策**：继续／转向／停止的证据和所有者签字。

NIST AI RMF 要求把评价连接到部署情境、记录测试集和方法、测量安全可靠性，并建立明确的 go／no-go 决策。[NIST AI RMF Core](https://airc.nist.gov/airmf-resources/airmf/5-sec-core/)

### 9. Future Market Outlook

近 1–2 年，Agent、Exposure Graph、MCP／工具生态、私有部署和 Agent 身份会继续成为成熟平台能力；表层功能更快商品化。3–5 年是否出现独立跨平台安全工作台，取决于客户是否愿意承担第二界面与治理成本，以及现有平台是否无法解决跨系统任务。5 年以上不作产品预测，因为模型、协议、监管和平台并购变化过快。

研究建议保持技术可替换、证据可复用、门禁可停止，不以任何厂商路线图作为需求事实。

### 10. Methodology and Source Verification

- **主要权威来源**：中国人大／国务院／网信办／工信部、NIST、CISA、FIRST、ISO、SEC 文件。
- **独立／专业研究**：ISC2、SANS、Verizon DBIR、Mandiant、同行评审论文。
- **供给证据**：FOBrain、Microsoft、Google、Palo Alto、CrowdStrike、Tenable、Qualys、Rapid7、ServiceNow 官方文档。
- **质量规则**：来源日期、主体、利益关系、适用范围和置信度逐项标记；厂商声明只证明供给；商业报告只作低置信度方向线索。
- **限制**：没有目标客户的一手证据，不能声称需求、定位、市场份额、购买意愿或产品效果已验证。

### 11. Appendices and Resources

本文件前述客户、痛点、旅程和竞争表即为详细附录；领域监管、标准与技术证据见配套 Domain Research。后续真实研究使用 `product-blueprint/research/discovery-execution-kit.md`、`product-validation-evidence-plan.md` 和 `owner-decision-gates.md`，不得从综合摘要跳过门禁。

### Market Research Conclusion

**已确认：** 行业问题类别、竞争拥挤、采用与治理条件、任务级验证方法和停止逻辑。  
**未确认：** FOBrain 真实剩余任务、用户、损失、输入输出、买方、定位、产品形态、市场规模和商业模式。  
**下一步：** DG-02 只批准最小研究访问；没有批准就不接触真实用户／数据，不形成 Product Brief。  
**Research Completion Date:** 2026-07-12  
**Research Status:** `PUBLIC_RESEARCH_COMPLETE / PRODUCT_UNVALIDATED`

---

<!-- Content will be appended sequentially through research workflow steps -->
