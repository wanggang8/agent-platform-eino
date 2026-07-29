---
stepsCompleted: [1, 2, 3, 4, 5, 6]
inputDocuments: []
workflowType: 'research'
lastStep: 1
research_type: 'domain'
research_topic: '以 FOBrain 为首个可验证业务样本的企业资产与漏洞风险运营 AI 协作领域约束'
research_goals:
  - '识别真实业务角色、任务、责任链、事实语义和结束证据'
  - '研究法规、标准、安全治理和人机动作边界'
  - '判断单一样本结论外推到多工具、多数据源或通用平台时需要补充的领域证据'
user_name: 'Vick'
date: '2026-07-11'
web_research_enabled: true
source_verification: true
---

# Domain Research：企业资产与漏洞风险运营 AI 协作

**Date:** 2026-07-11
**Author:** Vick
**Research Type:** Domain Research

---

## Research Overview

本研究以 FOBrain 为首个可验证样本，分析企业资产与漏洞风险运营中的行业结构、事实语义、角色责任、监管、技术趋势和安全边界。它不把 FOBrain 预设为最终产品，也不把历史功能、厂商路线或公开市场趋势当作用户需求。

公开证据显示，漏洞利用和跨系统协调持续形成运营压力，Exposure、SecOps 与 Agent 平台也正在融合；同时，数据、个人信息、未公开漏洞、跨境、CII、身份权限和自主动作对不同部署形成条件性约束。最重要的领域不变量是：权威事实必须留在来源系统，模型输出不是事实，高影响动作必须与权限、审批、审计和回退绑定。

完整综合见本文件末尾的 [Domain Research Synthesis](#domain-research-synthesis)。当前状态是公开领域研究完成、产品未验证；下一步是 DG-02 批准最小研究数据与访问，而不是进入架构或实现。

## Executive Summary

- **行业结论**：价值链从事实发现、归一／Exposure、任务协作、动作记录到治理服务；成熟平台正跨层整合。
- **技术结论**：平台原生 Agent、Exposure Graph、工具协议、Agent 身份、任务级 TEVV 和主权部署是明确供给趋势，但都不是自动需求。
- **监管结论**：纯内部 AI 通常不直接落入面向公众的生成式 AI 规章，但网络安全、数据、个人信息、漏洞与等级保护义务仍按实际活动适用。
- **产品结论**：只允许先验证一个真实任务；FOBrain 内嵌、独立 Workbench、多工具平台、非 AI 和不做仍然并列。

## Table of Contents

1. [Industry Analysis](#industry-analysis)
2. [Competitive Landscape](#competitive-landscape)
3. [Regulatory Requirements](#regulatory-requirements)
4. [Technical Trends and Innovation](#technical-trends-and-innovation)
5. [Recommendations](#recommendations)
6. [Domain Research Synthesis](#domain-research-synthesis)

## Domain Research Scope Confirmation

**Research Topic:** 以 FOBrain 为首个可验证业务样本的企业资产与漏洞风险运营 AI 协作领域约束

**Research Goals:**

- 识别真实业务角色、任务、责任链、事实语义和结束证据；
- 研究法规、标准、安全治理和人机动作边界；
- 判断单一样本结论外推到多工具、多数据源或通用平台时需要补充的领域证据。

**Domain Research Scope:**

- 行业分析：风险运营、漏洞管理、Exposure Management 与 AI 安全运营的市场结构和参与者；
- 监管与标准：适用法规、权威框架、组织政策和证据要求；
- 技术趋势：AI 辅助调查、证据引用、受控动作与人机协作；
- 经济因素：市场发展、采用驱动、成本和价值形成；
- 价值链：数据源、风险平台、AI 协作层、工单／变更系统及服务生态。

**明确边界:** FOBrain 仅为首个可验证样本，不是最终产品边界；产品是否扩展到多个工具、多个数据源、独立 Workbench 或通用平台仍待证。

**Research Methodology:**

- 所有时效性主张使用当前公开网页验证；
- 关键主张进行多来源交叉验证；
- 区分权威要求、独立研究、厂商声明、项目历史和真实用户证据；
- 对不确定性、适用性和研究缺口标记置信度。

**Scope Confirmed:** 2026-07-11（用户选择 `[C]`）

## Industry Analysis

### Market Size and Valuation

公开市场规模数字无法形成一个可审计的单一 TAM。原因是研究机构使用的类别不同：

| 公开估计 | 当前公开基准 | 预测 | 解释限制 |
| --- | ---: | ---: | --- |
| 全球安全总支出（母市场） | 2026 年约 3,080 亿美元 | 2029 年约 4,300 亿美元 | IDC 公开新闻稿，底层 Spending Guide 付费；只能说明母市场，不是 Exposure 或本项目 TAM。[IDC Security Spending Guide](https://www.idc.com/resource-center/press-releases/wwsecuritysg/) |
| 中国网络安全产业（广义） | 2023 年约 2,200 亿元人民币 | 未用于本项目预测 | 工信部转载行业发言；与全球安全支出和 Exposure 分类不可直接比较。[工信部转载](https://www.miit.gov.cn/xbymdz/xwfb/mtbd/twbd/art/2024/art_ad1b65e8b4b142f487ea8a81f3fd8d72.html) |
| Security & Vulnerability Management | 约 155.3 亿美元 | 2030 年约 223 亿美元，CAGR 7.5% | 宽口径，可能包含扫描、配置、补丁、服务等；不能等同本项目。[Research and Markets 摘要](https://www.researchandmarkets.com/reports/4591794/security-and-vulnerability-management-market) |
| Exposure Management | 2024 年约 33 亿美元 | 2030 年约 109.1 亿美元，2025–2030 CAGR 22.9% | 新兴聚合类别，定义与传统漏洞管理交叉。[Grand View Research 摘要](https://www.grandviewresearch.com/industry-analysis/exposure-management-market-report) |
| Attack Surface Management | 2025 年约 10.3 亿美元 | 2034 年约 50 亿美元，CAGR 21.03% | 更窄的攻击面发现／管理口径。[Fortune Business Insights 摘要](https://www.fortunebusinessinsights.com/attack-surface-management-market-110386) |

同一相邻市场的 2025 估计从约 10 亿到 155 亿美元，说明分类和收入边界差异巨大。当前正确结论不是选择最大数字，而是：**产品类别、经济买方和收费边界在 DG-07 前未知，因此 TAM 不可定稿。** 这些商业报告摘要仅作为低置信度方向线索，不是权威统计。

_Total Market Size：不可形成单一可信值。_  
_Growth Rate：传统漏洞管理摘要约为个位数增长；攻击面／Exposure Management 摘要约为 20% 以上增长，但口径不一致。_  
_Economic Impact：只能确认企业持续为发现、优先级、处置与验证投入；本项目可获得的价值份额未知。_  
_Confidence：低；待产品类别与买方确定后重新估算。_

### Market Dynamics and Growth

1. **从工具清单转向持续暴露管理流程。** Gartner 公开摘要指出，企业采用孤立、工具中心的方法难以降低暴露，CTEM 是面向过程的持续框架；这支持按任务与结果研究，而不是“一工具一功能”。[Gartner CTEM 摘要](https://www.gartner.com/en/documents/6735134)
2. **漏洞数量与资源不对称推动优先级。** CISA 维护 KEV 作为已在野利用漏洞的权威来源，并建议组织把它作为漏洞优先级输入；FIRST 明确 EPSS 只是未来 30 天利用概率，不能当作完整风险分。[CISA KEV](https://www.cisa.gov/known-exploited-vulnerabilities-catalog)、[FIRST EPSS User Guide](https://www.first.org/epss/user-guide)
3. **业务、技术与责任人之间存在协调成本。** NIST 把企业补丁管理定义为识别、确定优先级、获取、安装和验证，并指出业务／任务负责人和安全／技术管理之间常有价值分歧。[NIST SP 800-40 Rev.4](https://www.nist.gov/publications/guide-enterprise-patch-management-planning-preventive-maintenance-technology)
4. **AI 采用由效率和技能约束推动，但治理滞后。** ISC2 2025 调查显示 30% 已使用、42% 正评估／测试；NIST AI RMF 要求任务、角色、监督、测试和部署适用性可记录。[ISC2 AI Pulse](https://www.isc2.org/insights/2025/07/2025-isc2-ai-pulse-survey)、[NIST AI RMF Core](https://airc.nist.gov/airmf-resources/airmf/5-sec-core/)
5. **漏洞利用与修复速度持续形成运营压力。** Verizon 2026 DBIR 报告漏洞利用占已知初始访问路径的 31%，关键 KEV 完全修复比例为 26%，完全解决中位时间为 43 天；Mandiant 2026 报告漏洞利用连续第六年为最常见初始感染方式，占 32%。两者均来自各自事件／客户样本，存在选择偏差，但方向一致地支持“发现—处置—验证”的运营缺口，不能直接证明本产品需求。[Verizon 2026 DBIR](https://www.verizon.com/business/resources/T1e0/reports/2026-dbir-data-breach-investigations-report.pdf)、[Mandiant M-Trends 2026](https://cloud.google.com/blog/topics/threat-intelligence/m-trends-2026)

_Growth Drivers：攻击面复杂度、资源稀缺、优先级与证据需求、平台整合、AI 辅助。_  
_Growth Barriers：数据孤岛、权限、集成成本、可信度、治理、预算和既有平台替代。_  
_Market Maturity：传统漏洞管理成熟；Exposure Management 正在聚合；AI 协作处于早期采用与快速产品化阶段。_

### Market Structure and Segmentation

行业价值链可以按职责而不是厂商名称拆为五层：

1. **事实与发现层**：资产、漏洞、配置、身份、威胁情报、业务上下文和处置状态来源。
2. **归一与暴露层**：对象关联、去重、覆盖、攻击面、优先级和风险状态。FOBrain 已公开声称覆盖多源归一、VPT 动态优先级和漏洞闭环，因此这些能力不能被新产品直接宣称为独有价值。[FOBrain 官方产品页](https://huashunxinan.net/product-fobrain)
3. **任务协作层**：调查、解释、澄清、交接、报告和建议；AI 只是候选实现形态。
4. **动作与记录层**：工单、变更、审批、执行、回执和验证；最终事实必须留在权威系统。
5. **治理与服务层**：身份权限、数据治理、审计、模型风险、实施、托管服务和培训。

地理和行业差异主要来自法规、部署、数据敏感度、人员与预算，而不是一个通用功能清单。首个 FOBrain 样本只能验证第一个组织语境；扩展到多个工具／数据源必须证明跨系统任务与买方共同存在。

### Industry Trends and Evolution

- **平台化与第三方扩展。** Microsoft Security Copilot 同时提供独立、嵌入式和插件扩展形态；Tenable 和 Qualys 均公开把产品扩展为 Exposure Management／第三方平台。它们是厂商供给证据，说明竞争会来自既有平台扩展，而不证明用户已经获得效果。[Microsoft Security Copilot](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot)、[Tenable One](https://www.tenable.com/products/tenable-one)、[Qualys 2025 业绩说明](https://investor.qualys.com/static-files/fa509928-0dc3-4c39-84cf-dc7106b4a9a0)
- **平台整合通过并购加速。** ServiceNow 2026 年完成对 Armis 的收购，把资产可见性和网络风险工作流整合；Tenable 2025 年收购 Vulcan Cyber，补足第三方数据、优先级和修复运营。交易金额和战略表述属于公司公告／SEC 证据，只能证明供给侧整合，不证明客户效果。[ServiceNow 公告](https://newsroom.servicenow.com/press-releases/details/2026/ServiceNow-completes-Armis-acquisition-closing-the-gap-between-asset-visibility-and-cyber-risk/default.aspx)、[Tenable 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1660280/000166028026000005/tenb-20251231.htm)
- **从严重度到情境化优先级。** FIRST 建议把 CVSS Threat/Environmental 信息加入部署情境；EPSS、KEV、资产价值、暴露和补偿控制各有不同语义，不能被静默合并成未经验证的新分数。[FIRST CVSS v4 Consumer Guide](https://www.first.org/cvss/v4.0/implementation-guide)、[FIRST EPSS Guide](https://www.first.org/epss/user-guide)
- **从回答到受控 Agent。** 插件和 Agent 增加可调用能力，也增加过度功能、权限和自主性的风险；OWASP 建议最小功能、最小权限、用户上下文和高影响动作批准。[OWASP LLM06:2025](https://genai.owasp.org/llmrisk/llm062025-excessive-agency/)
- **评价从演示转向任务证据。** NIST AI RMF 要求在类似部署条件下测量性能，记录一般化限制，并定义人机监督；这与本蓝图的单任务、预注册结果和停止条件一致。[NIST AI RMF Core](https://airc.nist.gov/airmf-resources/airmf/5-sec-core/)

### Competitive Dynamics

竞争不是单一“AI 安全助手”类别，而是五类替代：

| 替代类别 | 结构性优势 | 结构性限制／待验证点 |
| --- | --- | --- |
| FOBrain 原生增强 | 已有数据、身份、流程和客户关系 | 目标部署能力、可扩展性和任务效果未知 |
| 大型安全平台内嵌 AI | 数据与工作流深度、已有分发、可嵌入 | 平台锁定、成本、跨栈覆盖和本地适用性 |
| 独立安全 Workbench | 可跨系统组织任务和证据 | 第二界面、集成、权限和第二套事实风险 |
| 非 AI 流程改进 | 可预测、低治理风险、易验收 | 对复杂跨对象推理和语言协作的增量有限，需实测 |
| 通用 Agent 平台 | 扩展性和生态潜力 | 与专用安全语义、权限、证据和买方距离更远；首个样本不能证明通用需求 |

进入壁垒主要是：取得权威事实与稳定标识、继承源系统权限、证明同任务净收益、建立可追踪与安全治理、获得设计伙伴和预算路径。Gartner 还公开指出 Exposure Management 供应商数量过多、能力趋同，意味着仅增加“AI”“平台”或“暴露管理”标签不构成差异。[Gartner Exposure Management 摘要](https://www.gartner.com/en/documents/6664234)

供给侧还包括传统 VM、综合安全平台、Exposure 专业厂商、云安全平台、ITSM／工作流、MSSP、系统集成商以及企业自研。Tenable 在其 10-K 中把市场描述为碎片化、竞争激烈且持续演化，并列举多类平台和内部方案作为竞争；这是厂商陈述，但与 Gartner 类别拥挤的公开信号一致。[Tenable 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1660280/000166028026000005/tenb-20251231.htm)、[Gartner EAP 分类](https://www.gartner.com/reviews/market/exposure-assessment-platforms)

### Step 2 Quality Assessment

- **高置信度**：行业正在从孤立工具转向持续过程；漏洞优先级必须结合多类证据；既有平台已同时布局内嵌、独立和扩展生态。
- **中置信度**：AI 协作可能在调查、总结、解释和跨系统上下文中产生任务价值，但必须逐任务验证。
- **低置信度**：市场规模、产品类别、目标买方、本项目可占市场和通用平台机会。
- **停止规则**：在 DG-07 前不得使用任何单一商业报告数字写 TAM，不得把厂商供给形态写成用户需求。

## Competitive Landscape

### Key Players and Market Leaders

该领域没有一个可验证的统一“AI 风险运营”排行榜。竞争发生在多个重叠层：

| 竞争组 | 代表供给 | 已公开形态 | 证据边界 |
| --- | --- | --- | --- |
| FOBrain 原生平台 | FOBrain | 多源资产／漏洞归一、VPT 动态优先级、漏洞闭环 | 仅为厂商声明；目标部署、API、用户采用和效果未知。[FOBrain](https://huashunxinan.net/product-fobrain) |
| 传统 VM 向 Exposure 扩展 | Tenable、Qualys、Rapid7 | 资产／漏洞、第三方数据、优先级、修复工作流、AI | SEC 文件和官方文档证明供给范围，不证明效果。[Tenable 10-K](https://www.sec.gov/Archives/edgar/data/1660280/000166028026000005/tenb-20251231.htm)、[Qualys 10-K](https://www.sec.gov/Archives/edgar/data/1107843/000110784326000008/qlys-20251231.htm)、[Rapid7 10-K](https://www.sec.gov/Archives/edgar/data/1560327/000156032726000008/rp-20251231.htm) |
| 综合 SecOps 平台 | Microsoft、Google、Palo Alto Networks、CrowdStrike | 数据平台、SIEM／XDR／SOAR、嵌入式或独立 AI、Agent、托管服务 | 强分发和数据优势；厂商结果需独立验证。[Microsoft](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot)、[Google](https://cloud.google.com/security/products/security-operations/investigate)、[Palo Alto](https://docs-cortex.paloaltonetworks.com/r/Cortex-XSIAM/Cortex-XSIAM-3.x-Documentation/Cortex-XSIAM-architecture)、[CrowdStrike](https://www.crowdstrike.com/content/dam/crowdstrike/marketing/en-us/documents/pdfs/data-sheets/crowdstrike-charlotte-ai-datasheet.pdf) |
| 工作流／ITSM 平台 | ServiceNow／Armis | 资产可见性、Exposure、工单、变更、审批和执行协作 | 擅长跨团队流程；安全事实深度取决于集成。[ServiceNow Vulnerability Response](https://www.servicenow.com/docs/r/security-management/vulnerability-response/vuln-landing-page.html) |
| 非 AI 与内部方案 | 现有页面、报表、脚本、SIEM／ITSM 配置、人工流程 | 低风险、可预测、与现有责任一致 | 经常是最低成本基线，不能被忽略。 |
| 通用 Agent 平台 | 通用 Agent builder 与插件生态 | 多工具编排、可扩展工作流 | 缺少安全领域事实、权限、验证和独立买方证据；首个 FOBrain 样本不能证明。 |

“市场领导者”只能按具体类别和公开口径讨论。Tenable 在 SEC 文件中引用 IDC 称其连续多年位列设备漏洞管理份额第一，但底层 IDC 报告付费且该陈述来自公司披露；不能据此建立整个 Exposure 或 AI 协作市场份额表。[Tenable 2025 年报](https://www.sec.gov/Archives/edgar/data/1660280/000166028026000017/a2025annualreport_bmkxweb.pdf)

### Market Share and Competitive Positioning

当前没有公开、统一且可审计的全球 Exposure Management、CTEM 或 AI Security Operations 份额数据。原因包括：类别边界重叠、平台收入不单独披露、厂商按不同资产／模块／服务计费，以及付费分析报告方法不公开。

因此本研究使用“定位与结构优势”而非伪造份额：

- **Tenable／Qualys／Rapid7**：从漏洞扫描和资产基础扩展到第三方数据、上下文优先级、修复和 Exposure 平台。
- **Microsoft**：独立与嵌入式 Security Copilot、插件和既有 Microsoft 安全栈；按 Security Compute Unit 容量计费并逐步与 E5/E7 权益结合。[Microsoft Pricing](https://www.microsoft.com/en-us/security/pricing/microsoft-security-copilot/)
- **Google／Palo Alto／CrowdStrike**：把 AI 放入原生 SecOps 数据、案件、响应或托管服务流程，而不是只提供独立聊天。[Google SecOps](https://cloud.google.com/security/products/security-operations/investigate)、[Cortex XSIAM](https://docs-cortex.paloaltonetworks.com/p/XSIAM)、[CrowdStrike Charlotte AI](https://www.crowdstrike.com/content/dam/crowdstrike/marketing/en-us/documents/pdfs/data-sheets/crowdstrike-charlotte-ai-datasheet.pdf)
- **ServiceNow／Armis**：把资产发现、Exposure 和企业工作流结合，强化从发现到责任协作和执行记录的链路。[ServiceNow 收购 Armis](https://newsroom.servicenow.com/press-releases/details/2026/ServiceNow-completes-Armis-acquisition-closing-the-gap-between-asset-visibility-and-cyber-risk/default.aspx)
- **FOBrain**：在中国企业资产攻击面和漏洞运营语境中拥有首个样本所需的既有事实与流程线索，但不能从单一厂商页面推断市场地位或最终产品边界。

### Competitive Strategies and Differentiation

1. **平台整合**：把 SIEM、XDR、SOAR、ASM、Exposure、云和身份数据放入共享数据层，降低跨产品切换。Palo Alto XSIAM 官方架构明确支持第三方遥测和扩展 Exposure 等能力。[Cortex XSIAM Architecture](https://docs-cortex.paloaltonetworks.com/r/Cortex-XSIAM/Cortex-XSIAM-3.x-Documentation/Cortex-XSIAM-architecture)
2. **从原生数据建立优势**：Tenable 的 Nessus 社区与传感数据、Microsoft 安全栈、Google SecOps、CrowdStrike Falcon、FOBrain 自身数据都形成上下文和分发壁垒。
3. **第三方生态**：Tenable 收购 Vulcan 后强调 100 多个安全产品的数据整合；Rapid7 宣称 500 多个集成；Microsoft 支持插件。数量是供给指标，真正壁垒是实体一致性、权限、新鲜度和错误语义。[Tenable 收购文件](https://www.sec.gov/Archives/edgar/data/1660280/000166028025000007/a1292025ex991.htm)、[Rapid7 Exposure Command](https://www.rapid7.com/products/command/exposure-management/)、[Microsoft Plugins](https://learn.microsoft.com/en-us/copilot/security/plugin-overview)
4. **从发现延伸到修复**：Qualys、Rapid7、ServiceNow 和 Tenable 都公开强调修复工作流、补丁、工单或自动化，竞争焦点已不只是“找到漏洞”。
5. **嵌入式、独立和 Agent 并行**：Microsoft 同时提供独立／嵌入式体验；Google 把 Gemini 放入案件调查；CrowdStrike 强调原生 Agent；说明形态竞争尚未收敛。[Microsoft Experiences](https://learn.microsoft.com/en-us/copilot/security/experiences-security-copilot)、[Google Triage Agent](https://docs.cloud.google.com/chronicle/docs/secops/triage-investigation-agent)
6. **托管与专家监督**：MSSP／MDR 把产品、流程和人员作为结果服务；对资源不足组织可能比独立软件更有竞争力。

### Business Models and Value Propositions

| 模式 | 常见计费／销售 | 价值主张 | 对本项目的挑战 |
| --- | --- | --- | --- |
| 订阅平台 | 按资产、模块、用户、期限或套餐 | 单一平台与持续更新 | 需要证明替换／增购理由 |
| 用量容量 | 按计算单元、数据摄入或调用量 | 灵活扩容 | 成本不可预测、需监控使用；Microsoft SCU 是公开实例 |
| 平台附加模块 | Exposure、AI、ASM、身份、云、自动化 add-on | 利用既有数据和合同 | 独立产品分发处于劣势 |
| 托管服务 | 按组织、资产、SLA 或服务包 | 以人员和结果弥补技能缺口 | 多租户、责任、数据隔离和毛利复杂 |
| 内部能力 | 开源、脚本、现有平台配置和人工流程 | 成本／控制可控 | 可能已足够解决单一任务 |

Tenable 10-K 说明其方案主要按一年期订阅销售，并通过直销、分销商、经销商和 MSSP 渠道触达；这显示成熟厂商不仅竞争产品，还竞争渠道、续费和交叉销售。[Tenable 10-K](https://www.sec.gov/Archives/edgar/data/1660280/000166028026000005/tenb-20251231.htm)

### Competitive Dynamics and Entry Barriers

- **市场拥挤**：Tenable 与 Qualys 均在 SEC 文件中称市场高度碎片化、竞争激烈，并把传统 VM、综合安全、云安全、端点、点工具和企业自研列为竞争者。
- **整合加速**：ServiceNow 收购 Armis；Tenable 收购 Vulcan Cyber 和 Apex；平台通过并购获得数据、优先级、修复和 AI 攻击面能力。
- **存量数据与权限壁垒**：成熟平台已有代理、遥测、身份、客户合同和管理员信任；新入口必须承担连接、权限和数据治理成本。
- **转换成本**：部署代理／连接器、建立资产基线、调优规则、迁移历史、重建工作流、培训、审计和重新评审供应商都会增加切换成本。
- **可信结果壁垒**：安全产品错误可能导致漏报、错误排序或生产中断；同任务效果、可追溯性和失败补救比模型能力更难复制。
- **第二套事实风险**：独立 Workbench 如果保存平行对象、状态、优先级或工单，会与 FOBrain／ITSM 冲突。

### Ecosystem and Partnership Analysis

生态控制点包括：

1. 资产／漏洞／身份／云／OT 数据来源；
2. 实体归一、风险与威胁语义；
3. 用户身份、权限和租户边界；
4. ITSM／工单／变更／补丁与回执；
5. 模型、插件、Agent 与评测；
6. 渠道、MSSP、集成商和客户安全评审。

Microsoft 通过 Entra、Azure、SCU、插件和商业市场控制一部分生态；Palo Alto、Google、CrowdStrike 通过原生数据平台和安全产品组合控制工作流；Tenable、Qualys、Rapid7 通过传感器、研究、连接器和 Exposure 平台扩展；ServiceNow 控制企业工作流。FOBrain 首个样本的优势或限制必须在 DG-02 后通过真实部署证明。

对新产品最关键的生态问题不是“能接多少工具”，而是：是否能在用户原权限下取得完成一个任务所需的事实，是否能返回来源系统执行和验证，是否有独立买方愿意承担连接与治理成本。

### Step 3 Quality Assessment

- **高置信度**：竞争边界重叠、平台化和并购加速；现有平台同时向数据、优先级、修复、AI 和工作流扩张。
- **中置信度**：跨平台任务协作可能存在切口，但会面对存量平台、ITSM、MSSP 和内部方案。
- **低置信度／未知**：统一市场份额、独立 Workbench 的可赢细分、FOBrain 相对竞争地位、买方和商业模型。
- **不可外推**：厂商功能页、排名或客户案例不能证明目标部署可用、用户采用或任务有效。

---

## Regulatory Requirements

> 本节是截至 2026-07-12 的产品研究，不构成法律意见。法规适用必须由目标客户法务根据实际用户、数据、部署、模型、动作、地区和行业重新确认。

### Applicable Regulations

首个验证场景若限定为“境内部署、仅授权员工使用、不向公众开放、不调用境外模型、只读建议并由人工审批”，其主要监管基础不是单一“AI 法”，而是网络安全、数据、个人信息、漏洞和客户行业义务的组合：

| 制度 | 触发条件 | 首个内部场景的候选适用性 | 边界 |
| --- | --- | --- | --- |
| 《网络安全法》修正版 | 在境内建设、运营、维护或使用网络 | 直接相关；2026-01-01 起施行的修法已生效 | 包括等级保护、运行安全、日志、漏洞补救和事件处置；不能因是内部 AI 而排除。[国家法律法规数据库](https://flk.npc.gov.cn/detail?fileId=&id=021e7d7684474107b8f3febbb1c4f8b5&title=%E4%B8%AD%E5%8D%8E%E4%BA%BA%E6%B0%91%E5%85%B1%E5%92%8C%E5%9B%BD%E7%BD%91%E7%BB%9C%E5%AE%89%E5%85%A8%E6%B3%95&type=) |
| 《数据安全法》 | 在境内开展数据处理 | 直接相关 | 数据分类分级、全流程安全和风险监测；重要数据有附加义务。[全国人大公报](https://wb.flk.npc.gov.cn/flfg/PDF/8a19eb7aaa1e463cb21ff7cc2bac99d1.pdf) |
| 《网络数据安全管理条例》 | 境内网络数据处理 | 直接相关；2025-01-01 起施行 | 外部模型、受托处理、数据提供、日志和安全措施须按真实数据流判断。[国家网信办](https://www.cac.gov.cn/2024-09/30/c_1729384452307680.htm) |
| 《个人信息保护法》 | 数据可单独或结合其他信息识别自然人 | 条件适用 | 账号、员工号、IP、设备、登录和操作记录可能构成个人信息；必须逐字段判断。[全国人大公报](https://wb.flk.npc.gov.cn/flfg/PDF/f67af9f12e1b4c83a998cf5a876ce0e4.pdf) |
| 《网络产品安全漏洞管理规定》 | 网络产品提供／运营及漏洞发现、收集、发布 | 高度相关 | 未公开漏洞不得向产品提供者之外的境外组织或个人提供；自有产品漏洞报告义务不能泛化到所有漏洞。[国家网信办](https://www.cac.gov.cn/2021-07/13/c_1627761607640342.htm) |
| 关键信息基础设施制度 | 客户已被主管部门认定并通知为 CII | 仅条件适用 | 不能仅凭金融、能源、政务等行业名称自行认定 CII。[司法部](https://xzfg.moj.gov.cn/front/law/detail?LawID=683&Query=%E4%BF%A1%E6%81%AF%E5%8C%96) |
| 《生成式人工智能服务管理暂行办法》 | 向境内公众提供生成式 AI 服务 | 纯内部工具通常不直接适用 | 企业内部研发／应用且未向公众提供服务的，官方明确不适用；开放给客户、多租户或公众后必须重判。[办法](https://www.cac.gov.cn/2023-07/13/c_1690898327029107.htm)、[答记者问](https://www.cac.gov.cn/2023-07/13/c_1690898326863363.htm) |
| 算法推荐、深度合成与生成内容标识 | 提供相应互联网信息服务／公众生成传播服务 | 纯内部运营工具通常不自动触发 | “使用 LLM”本身不等于提供算法推荐或深度合成互联网信息服务。[算法推荐规定](https://www.cac.gov.cn/2022-01/04/c_1642894606364259.htm)、[深度合成规定](https://www.cac.gov.cn/2022-12/11/c_1672221949354811.htm)、[标识办法](https://www.cac.gov.cn/2025-03/14/c_1743654684782215.htm) |

法规适用不是一次性标签。集团员工、外包人员、供应商、客户、自助注册用户和不特定公众之间的边界，以及模型是否由境外主体接收数据，都必须按实际服务关系和数据流确认。

### Industry Standards and Best Practices

- **等级保护**：GB/T 22239-2019 是网络安全等级保护基本要求的推荐性国家标准，用于支持法定等级保护义务落地；具体对象、定级、备案和测评要求取决于系统边界、等级和主管要求。[国家标准信息公共服务平台](https://std.samr.gov.cn/gb/search/gbDetailedCNF?id=88F4E6DA63434198E05397BE0A0ADE2D)
- **中国 AI 治理参考**：《人工智能安全治理框架》可用于建立风险分类、治理和持续监测，但不能自动替代具体法律义务。[国家网信办](https://www.cac.gov.cn/rootimages/uploadimg/1727568303900999/1727568303900999.pdf)
- **NIST AI RMF／GenAI Profile、NIST CSF 2.0**：是自愿指导，可用于任务、监督、测试、供应链与安全结果映射，不是中国法定认证。[NIST AI RMF](https://www.nist.gov/itl/ai-risk-management-framework)、[NIST GenAI Profile](https://www.nist.gov/publications/artificial-intelligence-risk-management-framework-generative-artificial-intelligence)、[NIST CSF 2.0](https://www.nist.gov/publications/nist-cybersecurity-framework-csf-20)
- **ISO/IEC 42001、27001、23894**：分别提供 AI 管理体系、信息安全管理体系和 AI 风险管理参考；除非合同、采购或监管明确要求，否则不能写成强制认证。[ISO/IEC 42001](https://www.iso.org/standard/42001)、[ISO/IEC 27001](https://www.iso.org/standard/27001)、[ISO/IEC 23894](https://www.iso.org/standard/77304.html)
- **OWASP GenAI／Agentic 风险清单**：适合检查提示注入、工具滥用、身份权限、供应链和过度代理权，但属于社区安全实践而非法律。[OWASP GenAI Security Project](https://genai.owasp.org/)

### Compliance Frameworks

进入真实研究和试点前应采用“适用性决策表”，而不是把全部法规变成功能：

1. **主体**：谁是网络运营者、数据处理者、个人信息处理者、模型提供者、受托人和动作执行者。
2. **用户**：单一法人内部员工、集团、外包／供应商、客户还是公众。
3. **数据**：一般数据、个人信息、敏感个人信息、重要／核心数据、未公开漏洞、工作秘密或国家秘密。
4. **数据流**：采集、模型输入、存储、日志、遥测、远程运维、导出、删除及境外访问。
5. **系统与客户身份**：等保对象及等级，客户是否被正式认定为 CII，是否有行业专门规则。
6. **AI 行为**：只读检索／摘要、建议、自动化决策，还是封禁、删除、隔离、修复等写域动作。
7. **服务关系**：内部工具、客户专属部署、多租户 SaaS 或面向公众的信息服务。

`DG-02` 只批准研究数据和访问，不代表最终产品适用性已经确定；完整适用性必须在 `DG-08` 动作边界和 `DG-09` 部署／客户边界确认后关闭。

### Data Protection and Privacy

安全数据不能一概认定为个人信息，也不能默认不包含个人信息。账号、员工号、IP、设备标识、邮箱、登录行为、操作日志、调查和处置记录应逐字段判断能否识别自然人，以及是否会用于对员工产生重大影响的判断。

候选合规条件包括：合法性基础、目的限定、最小必要、保存期限、访问控制、受托处理合同与监督；在处理敏感个人信息、自动化重大影响决定、向其他处理者提供、公开或跨境前，需判断是否进行个人信息保护影响评估。处理超过 1000 万人个人信息的处理者至少每两年开展一次合规审计；监管部门还可在较大风险或重大事件等条件下要求专项审计。[个人信息保护合规审计管理办法](https://www.cac.gov.cn/2025-02/14/c_1741233507681519.htm)

数据出境不能只按服务器所在地判断。调用境外 LLM、海外 SaaS、境外运维或允许境外主体远程访问均可能构成出境。现行规则针对 CII、重要数据以及不同规模的一般／敏感个人信息规定安全评估、标准合同、认证或豁免路径；数量豁免不解除一般数据保护义务，也不解除未公开漏洞的独立限制。[促进和规范数据跨境流动规定](https://www.cac.gov.cn/2024-03/22/c_1712776612187994.htm)、[官方答记者问](https://www.cac.gov.cn/2024-03/22/c_1712776611649184.htm)

### Licensing and Certification

- 自用内部软件不能仅因名为“安全运营 Workbench”就断言需要产品认证；商业化销售或提供时，应核对其是否实质落入现行网络关键设备／网络安全专用产品目录。
- 等保定级、备案和测评针对具体网络系统边界，不是“AI 产品证书”。
- ISO/IEC 42001、27001 等通常是自愿认证或采购／合同条件，不得在 PRD 中写成普遍法定义务。
- 若向公众提供具有舆论属性或社会动员能力的生成式 AI／算法服务，可能涉及安全评估、算法备案及生成内容标识；纯内部场景不能据此提前制造备案需求。[生成式人工智能服务备案公告](https://www.cac.gov.cn/2024-04/02/c_1713729983803145.htm)

### Implementation Considerations

在 DG-02 尚未取得前，首个验证应采用可逆的研究安全假设，而不是永久产品范围：

- 境内受控环境，明确授权用户，不向公众开放；
- 只使用获批的最小数据集，默认禁止把未公开漏洞、PoC、攻击路径和内部资产详情发给境外模型／人员；
- 模型、遥测、日志、远程支持和再委托路径逐一登记；
- 默认只读建议，来源可追溯；封禁、删除、隔离、修复、权限变更及事件等级／报告决定由人工审批并留审计；
- 不让 AI 延迟、压制或替代法定事件报告；
- 预留删除、回退、停止、供应商退出和事故响应；
- 对外开放、多租户化、跨境或服务 CII 客户前重新通过适用性门禁。

2026 年 8 月 20 日起，《网络数据安全风险评估办法》将施行：重要数据处理者应每年开展风险评估，一般数据处理者被鼓励至少每三年一次，且报告、报送和整改存在条件性要求。当前产品蓝图必须记录这一即将生效的变化，但在未确认重要数据处理者身份前不能把年度报送写成所有客户的需求。[国家网信办第 24 号令](https://www.cac.gov.cn/2026-06/18/c_1783525609815499.htm)

### Risk Assessment

| 优先级 | 风险 | 蓝图处理 |
| --- | --- | --- |
| 极高 | 未公开漏洞、PoC、攻击路径或内部资产详情发送给境外模型／人员 | DG-02 前默认禁止；法务和数据所有者明确批准后方可改变 |
| 极高 | 未识别个人信息、重要数据或 CII 状态即接入外部模型／云 | 阻断真实数据试点，先完成字段、数据流和客户身份判定 |
| 高 | 内部工具扩展给客户／公众后仍沿用“内部不适用”结论 | 服务关系变化自动重开 DG-09 |
| 高 | AI 自动执行高影响动作或改变事件报告 | DG-08 要求最小权限、人工批准、回退和审计 |
| 高 | AI 对员工或自然人作出重大影响决定 | 重新确认合法基础、影响评估、解释、申诉和人工复核 |
| 中高 | 商业化产品落入安全专用产品目录但未分类 | 商业化前取得产品分类与法务意见 |
| 中高 | 供应商合同未覆盖数据、再委托、保密、删除、审计和事件响应 | 形成供应商证据包和退出条件 |
| 中 | 把自愿框架或推荐标准错误写成强制法律 | 所有要求标注“法律／监管／合同／自愿实践”来源等级 |

### Step 4 Quality Assessment

- **高置信度**：纯内部 AI 不自动适用面向公众的生成式 AI 规则；网络、数据、个人信息和漏洞义务仍按实际处理活动适用；未公开漏洞外发和未经判断的数据跨境是高风险边界。
- **条件性**：个人信息、重要数据、CII、等保等级、网络安全专用产品和公众服务身份均需目标部署证据或主管／法务判断。
- **时间敏感**：《网络数据安全风险评估办法》已公布但将于 2026-08-20 施行；EU AI Act 及其 2026 修订进程也需在任何欧盟部署前重新核验最终文本。
- **门禁结论**：监管研究定义安全的研究起点和必须澄清的条件，不得直接决定最终产品形态或生成一套未经适用性验证的“合规功能清单”。

---

## Technical Trends and Innovation

### Emerging Technologies

当前最明显的技术变化不是“聊天框变强”，而是安全平台把 AI 从问答扩展到**有身份、有上下文、能调用工具并受治理的任务代理**。Microsoft Security Copilot 已提供可配置权限的安全 Agent；CrowdStrike 提供在 Falcon 内构建、测试和管理安全 Agent 的 AgentWorks；Palo Alto Cortex AgentiX 把 Agent、SOAR playbook、集成和 MCP 支持放在同一运营平台；Google 则以 Agentic SOC 描述取证、分析和结论链路，同时保留最终人工监督。[Microsoft Security Copilot Agents](https://learn.microsoft.com/en-us/copilot/security/agents-overview)、[CrowdStrike Charlotte AI](https://www.crowdstrike.com/en-us/platform/charlotte-ai/)、[Cortex AgentiX](https://docs-cortex.paloaltonetworks.com/r/Cortex-AgentiX/Cortex-AgentiX-Documentation/Get-Started-with-Cortex-AgentiX)、[Google Agentic SOC](https://cloud.google.com/security/resources/agentic-soc)

第二个变化是**从孤立漏洞分数转向关系图和攻击路径**。Tenable、Qualys 和 Microsoft 都公开以资产、身份、配置、漏洞、云和外部来源的关系建立攻击路径或暴露图。这说明“跨对象上下文”正在成为成熟平台的基础能力，新产品不能仅凭“把更多数据放在一起”形成差异。[Tenable Attack Path](https://docs.tenable.com/exposure-management/Content/attack-path/attack-path.htm)、[Qualys Attack Path Analysis](https://docs.qualys.com/en/etm-identity/latest/attack_path_analysis/qualys_attack_path_analysis.htm)、[Microsoft Security Exposure Management](https://learn.microsoft.com/en-gb/security-exposure-management/cross-workload-attack-surfaces)

第三个变化是**工具接入协议标准化**。MCP 使 Agent 以统一方式访问外部数据和工具，并已加入传输层授权及企业集中授权扩展；但协议只解决一部分互操作问题，不自动解决工具可信、业务语义、数据最小化、授权决策和结果正确性。[MCP Authorization](https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization)、[Enterprise-Managed Authorization](https://blog.modelcontextprotocol.io/posts/enterprise-managed-auth/)

第四个变化是**Agent 身份成为独立治理对象**。Google 2026 年公开将 Agent Identity 作为区别于人和普通服务账号的一等主体；这反映出代理需要独立身份、最小权限、生命周期和审计，而不能借用共享管理员凭据。[Google Cloud IAM 更新](https://cloud.google.com/blog/products/identity-security/whats-new-in-iam-security-governance-and-runtime-defense)

_Confidence：高；多家成熟平台的官方产品和文档方向一致，但厂商供给不等于目标客户需求或实际效果。_

### Digital Transformation

安全运营的数字化方向正在从“更多控制台”转向三种收敛：

1. **事实层收敛**：资产、身份、漏洞、配置、威胁和业务上下文通过图或统一对象关联。
2. **任务层收敛**：检测、调查、优先级、交接、修复建议和验证被组织为可持续工作流。
3. **治理层收敛**：Agent、连接器、权限、审批、审计、成本和模型风险进入统一管理。

成熟供应商正在同时扩展这些层。Tenable 将 Tenable One 定位为跨 IT、云、容器、Web、身份、第三方连接和 AI 资产的 Exposure 平台；Rapid7 将 Exposure Management 与检测响应合并叙述；ServiceNow 完成 Armis 收购后，把资产可见性、风险与企业工作流进一步整合。它们证明竞争边界正在收敛，也意味着单点功能更容易被平台吸收。[Tenable 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1660280/000166028026000005/tenb-20251231.htm)、[Rapid7 2025 Form 10-K](https://www.sec.gov/Archives/edgar/data/1560327/000156032726000008/rp-20251231.htm)、[ServiceNow／Armis 公告](https://newsroom.servicenow.com/press-releases/details/2026/ServiceNow-completes-Armis-acquisition-closing-the-gap-between-asset-visibility-and-cyber-risk/default.aspx)

对本项目而言，这一趋势不能推导出“必须做通用平台”。相反，它提高了验证门槛：只有当一个任务在现有 FOBrain、现有工作流和非 AI 改进之后仍有净价值，并且跨系统协作确实是必要原因，才有理由扩展产品边界。

### Innovation Patterns

| 创新模式 | 市场信号 | 对蓝图的含义 | 主要反证 |
| --- | --- | --- | --- |
| 平台原生 Agent | Microsoft、Google、CrowdStrike、Palo Alto 均把 Agent 放入已有数据与控制面 | 原生上下文、身份和分发是强优势 | 新产品可能只是存量平台可复制的附加功能 |
| Exposure／Security Graph | 多家平台关联资产、身份、漏洞和攻击路径 | 应验证关系信息是否真正改变一个任务决定 | 图更大不等于判断更准，实体错误可能被放大 |
| Agent-to-tool 协议 | MCP 等协议降低工具连接门槛 | 工具生态可更开放、可替换 | 接入数量不等于业务语义、权限和结果可靠 |
| 任务级评测 | NIST 推动 TEVV 和结构化 GenAI 评估 | 用真实任务、用户和环境比较人／AI／不同方案 | 静态 benchmark 可能与真实部署脱节 |
| 主权／私有／断连 AI | 厂商提供本地推理、私有云及断连环境 | 高敏客户可能要求模型与数据边界可选 | 本地化带来模型、算力、升级和运维成本 |
| Agent 身份与治理 | Agent 作为独立主体配置权限和审计 | 自主程度必须与身份、权限和责任绑定 | 仅增加审批界面不能解决错误授权与人类误批 |

NIST 的 TEVV 项目和 GenAI 评估计划强调测量、验证和人与 AI 的比较；这比通用“模型排行”更接近本项目需要的任务证据。[NIST TEVV](https://www.nist.gov/ai-test-evaluation-validation-and-verification-tevv)、[NIST GenAI Evaluation](https://ai-challenges.nist.gov/genai)

### Future Outlook

未来 12–36 个月最可能持续的方向是：

- **Agent 功能快速商品化**：主流安全平台会继续增加调查、检测工程、威胁狩猎、优先级和响应 Agent；“有 Agent”本身不构成定位。
- **平台与工作流进一步整合**：Exposure、SecOps、云、身份、ITSM 和托管服务的边界继续重叠，分发和既有事实将比单一模型更重要。
- **评测从静态题集转向动态真实环境**：2026 年公开研究显示，真实威胁狩猎任务中前沿模型仍可能表现很差，说明能力提升与可靠部署之间存在间隙。[Cyber Defense Benchmark](https://arxiv.org/abs/2604.19533)
- **协议生态扩大，同时安全面扩大**：MCP 降低工具集成成本，但恶意工具、身份混淆、越权、上下文／记忆污染和跨 Agent 传播成为新风险。NSA 2026 年也发布了 MCP 安全设计考虑。[NSA MCP Security Design Considerations](https://media.defense.gov/2026/Jun/02/2003943289/-1/-1/0/CSI_MCP_SECURITY.PDF)
- **主权与部署选择成为产品条件**：Microsoft 已公开支持客户控制边界内的本地推理和完全断连场景，说明高敏行业会要求云、私有、本地或断连的选择；但具体需求必须由目标客户验证。[Microsoft Sovereign Cloud](https://blogs.microsoft.com/blog/2026/02/24/microsoft-sovereign-cloud-adds-governance-productivity-and-support-for-large-ai-models-securely-running-even-when-completely-disconnected/)

这些是方向性判断，不是发布时间承诺。模型、协议和厂商路线变化快，任何进入 PRD 的技术依赖都必须绑定当前证据日期和替代路径。

### Implementation Opportunities

当前只支持五类**待验证机会**，不支持直接立项全部能力：

1. **单任务证据化协作**：把来源、时间、对象、冲突、缺口和建议放在同一任务内，验证是否减少总复核成本。
2. **原权限下的最小工具协作**：只连接完成该任务必需的来源，验证跨系统是否真是价值原因，而不是功能展示。
3. **人机职责分级**：区分只读检索、摘要、建议、拟执行和真实执行，按风险逐级验证。
4. **真实任务评测资产**：冻结历史样本、基线、失败样本、结果和停止条件，使模型或实现可以替换而不丢失验收基准。
5. **部署与模型可替换性假设**：在 DG-09 记录目标客户是否需要境内云、私有、本地或断连，以及其愿意承担的成本。

只有当同一问题跨多个任务／系统／客户重复出现、具有共同买方且现有平台仍未解决时，才可把这些机会外推为独立 Workbench 或通用平台。

### Challenges and Risks

- **可靠性错觉**：语言流畅不等于事实正确；静态 benchmark 不能替代目标任务。
- **过度代理**：OWASP 2026 Agentic Top 10 把目标劫持、工具滥用、身份／权限、记忆与上下文污染、非预期执行和级联故障列为关键风险。[OWASP Agentic Top 10](https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/)
- **人类误批**：人工在环降低风险但不消除风险；Google 的 MCP 安全文档明确指出人工可能错误批准 Agent 建议。[Google Cloud MCP Security](https://docs.cloud.google.com/mcp/ai-security-safety)
- **第二套事实**：独立产品若复制资产、漏洞、优先级、工单和状态，会与权威系统冲突。
- **实体与关系污染**：错误的资产归一、身份关系或陈旧数据会让攻击路径和 AI 解释同时失真。
- **协议与供应链**：工具、插件、MCP Server、模型和 Agent 都可能成为第三方供应链入口。
- **成本与运维**：推理、连接器、数据更新、评测、监督、审计和本地部署可能抵消任务收益。
- **平台挤压**：成熟厂商已有数据、权限、渠道和合同，可以快速复制表层 AI 功能。
- **监管与部署分化**：公众／内部、境内／跨境、一般／重要数据、普通／CII 客户会形成不同条件，不能用一个默认部署覆盖全部。

## Recommendations

### Technology Adoption Strategy

采用“任务先于技术、证据先于扩展”的顺序：

1. DG-02 只批准最小研究数据与访问。
2. DG-03 冻结真实任务样本和角色。
3. DG-04 只选择一个高价值任务，也允许没有合格任务。
4. DG-05 预注册结果、保护指标和停止条件。
5. DG-06 同时比较原生增强、嵌入式 AI、独立入口、非 AI 和不做。
6. DG-08／DG-09 再决定自主程度、数据、部署和适用性。

不要因 Agent、MCP、图谱或私有模型是趋势而提前写入需求；只有它们被证明是完成已确认任务的必要条件时，才能进入 PRD。

### Innovation Roadmap

- **阶段 A：只读事实辅助**——检索、聚合、摘要和证据回链；目标是验证事实正确和复核时间。
- **阶段 B：受控建议**——形成优先级解释、缺失信息和下一步建议；目标是验证决策质量和采用。
- **阶段 C：拟执行与审批**——只在动作确有价值时生成预览、影响、审批和回退材料。
- **阶段 D：有限执行**——仅对已证明可回退、可审计且获批的动作开放最小权限。
- **阶段 E：跨任务／平台扩展**——只有多个独立样本证明共同任务、共同买方和共同事实模型后进入。

每一阶段都允许停止、回到非 AI 方案或保留为 FOBrain 原生增强。

### Risk Mitigation

- 为每个任务保留权威来源、权限、时间、新鲜度、未知和冲突；不把模型输出当作事实。
- Agent 使用独立身份和最小权限；工具与数据源进入受控清单，记录版本和责任人。
- 对提示、工具、记忆、上下文、模型和供应商执行威胁建模、回归评测和异常监控。
- 高影响动作必须预览、审批、审计、回执和回退；人工批准本身也要通过清晰证据降低误批。
- 使用真实历史样本和失败样本持续 TEVV，不以厂商 benchmark 或单次演示替代。
- 为模型、协议、部署和供应商保留替代路径，避免让技术选择成为产品事实。

### Step 5 Quality Assessment

- **高置信度**：Agent 平台化、Exposure Graph、工具协议、Agent 身份、任务级评测和主权部署均已形成明确供给趋势。
- **中置信度**：这些趋势会提高跨系统任务协作和受控自动化的可行性。
- **未知**：目标客户是否需要其中任何一项、愿意承担多少集成／治理成本，以及它们能否产生净任务收益。
- **门禁结论**：技术趋势只能形成可检验选项，不能绕过 DG-02–DG-09，也不能把独立 Workbench、通用平台或 Agent 自主执行写成既定方向。

---

## Domain Research Synthesis

### 1. Research Introduction and Methodology

安全运营 AI 的关键矛盾不是模型能否生成答案，而是答案是否来自正确事实、能否在原权限下被复核、是否改善真实任务，以及错误时谁负责。研究综合使用 2024–2026 年法律／监管、国家与国际标准、独立行业调查、事件报告、上市公司披露和厂商文档；对中国企业内部场景做重点适用性分析，同时保留全球供给与框架作为比较。

研究采用五个交叉维度：事实与角色、任务与结果、技术与竞争、数据与监管、采用与买方。所有结论区分 `VERIFIED_PUBLIC_FACT`、`INFERENCE`、`HYPOTHESIS` 和 `UNKNOWN`，未取得目标部署证据的事项不得进入 PRD。

### 2. Industry Overview and Market Dynamics

行业价值链包含：

1. 资产、漏洞、身份、配置、威胁、业务和工单等权威事实；
2. 实体归一、Exposure、优先级和攻击路径；
3. 调查、解释、交接和建议等任务协作；
4. 工单、变更、执行、回执和验证；
5. 身份、数据、模型、审计、服务和监管治理。

传统漏洞管理成熟，Exposure／CTEM 正在聚合，SecOps 与 Agent 又向事实、工作流和动作层扩展。市场数字因分类重叠无法形成单一 TAM；对本项目有意义的经济单位只能是一个任务当前损失与候选方案总成本。

### 3. Technology Landscape and Innovation

供给趋势形成六个方向：平台原生 Agent、关系图／攻击路径、工具协议、Agent 独立身份、任务级评测和主权／私有部署。中国 2026 年《智能体规范应用与创新发展实施意见》也强调应用牵引、先易后难、用户最终决策权、权限边界、行为可验证／可追溯和全周期供应链安全。[国家网信办](https://www.cac.gov.cn/2026-05/08/c_1779979789523320.htm)

这些方向提高了可行性，同时扩大了攻击面。OWASP Agentic Top 10、NSA MCP 安全考虑和 Google MCP 安全文档都表明工具、身份、记忆、上下文、授权和人类误批需要独立控制。[OWASP](https://genai.owasp.org/resource/owasp-top-10-for-agentic-applications-for-2026/)、[NSA](https://media.defense.gov/2026/Jun/02/2003943289/-1/-1/0/CSI_MCP_SECURITY.PDF)、[Google Cloud](https://docs.cloud.google.com/mcp/ai-security-safety)

### 4. Regulatory Framework and Compliance

首个内部验证的临时安全边界是：境内受控用户、最小数据、只读建议、人工审批、来源追溯和完整审计。这是研究假设，不是最终产品范围。

适用性必须逐部署回答：主体是谁、用户是否为公众、字段是否含个人／敏感个人信息或重要数据、客户是否被认定为 CII、数据是否出境、是否涉及未公开漏洞、系统等保边界、AI 是否产生重大影响决定或执行高风险动作。纯内部场景不自动适用面向公众的生成式 AI 规章，但不排除网络安全、数据、个人信息和漏洞义务。[生成式 AI 办法答记者问](https://www.cac.gov.cn/2023-07/13/c_1690898326863363.htm)、[网络数据安全管理条例](https://www.cac.gov.cn/2024-09/30/c_1729384452307680.htm)

### 5. Competitive Landscape and Ecosystem

FOBrain 原生事实与流程、Exposure 厂商、综合 SecOps、ITSM、MSSP、内部方案和通用 Agent 平台共同竞争。平台整合和并购使单点功能更容易被吸收；连接器数量不是壁垒，真正控制点是实体一致性、身份权限、权威事实、新鲜度、任务结果、分发和合同。

生态伙伴包括数据／工具提供方、身份与策略、模型与 Agent、ITSM／执行系统、评测／审计、渠道／MSSP 和监管／标准机构。任何合作都必须明确数据角色、再委托、漏洞与事件、权限、退出和责任。

### 6. Strategic Insights and Domain Opportunities

跨域综合只支持一个战略：**以一个真实、可证伪的任务验证“事实—判断—责任—结束证据”链路。**

候选机会按证据顺序排列：

1. 只读事实聚合与来源回链；
2. 缺口、冲突和适用性提示；
3. 证据化优先级／下一步建议；
4. 经批准的拟执行与回退材料；
5. 只有前述均有效才考虑有限动作；
6. 只有多个独立样本重复才考虑跨工具／平台扩展。

每一级都允许被 FOBrain 原生增强、非 AI 流程或不做取代。

### 7. Implementation Considerations and Risk Assessment

这里的“实施”指研究与验证，不是软件架构：

- 先冻结样本、输入、输出、角色、来源、基线和结束证据；
- 对每个字段记录数据分类、权限、保留、跨境和模型可见性；
- 把模型输出与 StructuredResult／Product Facts 等安全投影概念区分，模型不成为事实源；
- 将只读、建议、拟执行和真实执行分级；
- 记录失败、人工覆盖、撤销、事件和回退；
- 使用独立评审和真实用户反馈校验任务结果。

最高风险是未公开漏洞／资产详情外发、错误实体关系被放大、模型越权、法定事件被延迟、第二套事实、人工误批以及本地／私有部署成本被低估。

### 8. Future Outlook and Strategic Planning

近 1–2 年，Agent 和图谱能力会快速进入既有安全平台，监管与标准也会进一步细化；新的表层 AI 功能窗口将缩短。3–5 年，是否形成跨平台独立产品取决于协议、身份与治理成熟度，以及客户是否愿意承担第二入口。长期不做确定预测。

战略规划应保持三个可替换性：模型可替换、工具／协议可替换、部署可替换；保持两个不变量：权威事实不复制、验收围绕任务结果。

### 9. Research Methodology and Source Verification

| 来源层 | 代表来源 | 允许支持的主张 |
| --- | --- | --- |
| 法律／监管 | 全国人大、国务院、网信办、工信部 | 法律状态、适用条件和监管方向 |
| 权威标准／框架 | NIST、CISA、FIRST、ISO、国家标准平台 | 术语、风险、评价和控制参考 |
| 独立行业证据 | ISC2、SANS、Verizon、Mandiant | 行业行为、事件和趋势，不能外推目标客户 |
| 公司披露 | SEC 年报、收购公告 | 商业范围、竞争、渠道和风险 |
| 厂商文档 | FOBrain 及国内外平台 | 功能供给，不证明效果 |
| 项目历史 | 已存档 docs | 假设、验收基线和安全边界参考，不是新需求真相 |

质量限制仍是没有目标用户、任务、部署和买方证据。研究的高置信度只覆盖公开行业、监管和供给事实；产品结论保持未知。

### 10. Appendices and Additional Resources

详细市场与价值链表见 Industry Analysis；竞品与生态见 Competitive Landscape；适用性矩阵见 Regulatory Requirements；技术趋势、风险和阶段策略见 Technical Trends and Innovation／Recommendations。后续所有真实证据必须登记到产品蓝图 evidence ledger，并保留来源、日期、权限、适用范围和置信度。

### Research Conclusion

**领域不变量：** 权威系统保留事实；AI 只消费安全投影；高影响动作与权限、审批、审计、回执和回退绑定。  
**尚未验证：** 真实任务、数据可得性、用户责任链、任务损失、产品形态、跨工具必要性和商业关系。  
**下一门禁：** DG-02 研究数据与访问批准。  
**Research Completion Date:** 2026-07-12  
**Research Status:** `PUBLIC_DOMAIN_RESEARCH_COMPLETE / PRODUCT_UNVALIDATED`

---

<!-- Content will be appended sequentially through research workflow steps -->
