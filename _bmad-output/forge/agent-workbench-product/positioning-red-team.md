---
status: provisional
review_type: adversarial-product-positioning
reviewed_at: 2026-07-11
owner_evidence: pending
---

# FOBrain 单一重要剩余任务定位假设红队审查

## 审查结论

**临时结论：PIVOT。** 不停止对 FOBrain 用户剩余任务的研究，但应暂停把 AI 或“协作与决策入口”当成已经成立的产品定位；先收窄为一个可证伪假设：**针对一个通过高频效率型或低频高损失型准入路径的重要 FOBrain 任务，按预注册规则比较原生增强、AI、独立 Workbench、非 AI 和不做。**

这是推断，不是用户事实。当前没有 FOBrain 活跃用户访谈、任务观察、产品使用数据、客户购买证据或对照实验；产品负责人也尚未确认产品关系、首要用户和唯一核心结果。现有研究范围已明确这些材料只能保持 `DRAFT`（`_bmad-output/planning-artifacts/product-blueprint/research/research-scope-proposal.md:49-56`）。因此，本结论为 **pending owner/user evidence**，不得据此把 Product Brief、PRFAQ 或 PRD 标为 final。

### Canonical evidence records

| ID | 主张 | Canonical evidence record | 稳定来源 |
| --- | --- | --- | --- |
| PRT-01 | FOBrain 官方公开声称覆盖多源归一、VPT 动态优先级和漏洞闭环。 | `VENDOR_SOURCE / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | [FOBrain 官方页](https://www.huashunxinan.net/product-fobrain)；即 `VENDOR_CLAIMED`。 |
| PRT-02 | 目标 FOBrain 用户仍存在一个重要剩余任务。 | `NO_SOURCE / NONE / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | [研究综合草案 C/H/U 记录](../../planning-artifacts/product-blueprint/research/research-synthesis-draft.md)；等待 DG-02～DG-04。 |
| PRT-03 | AI 或独立 Workbench 是胜出解法。 | `NO_SOURCE / NONE / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 等待 DG-05／DG-06；原生增强、非 AI 和不做同样合法。 |
| PRT-04 | 产品应成为 FOBrain 标准功能、增值模块、内部工具或独立产品。 | `NO_SOURCE / NONE / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 等待 `TASK_EFFECTIVE`、`ADOPTION_EVIDENCED` 与 DG-07 Owner 决策。 |
| PRT-05 | Microsoft 官方文档公开描述 Security Copilot 的自然语言、组织数据 grounding、插件、权限控制、人工复核及输出可能不准确／不完整／过时等产品主张。 | `VENDOR_SOURCE / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 即 `VENDOR_CLAIMED`；只证明 Microsoft 的公开产品说明，不证明本项目目标部署、采用、效果或必须跟进。 |

以下正文中的“事实／推断／置信度”只作解释；影响定位的结论必须以上述记录或[统一证据模型](../../planning-artifacts/product-blueprint/research/evidence-status-model.md)中的完整六字段为准。

## 已确认的证据边界

- **事实：** FOBrain 官方把产品定义为网络资产攻击面管理平台，并声称已经提供多源资产关联融合、完整可信资产清单、漏洞数据去重归一、动态优先级排序以及漏洞闭环处置。[FOBrain 官方产品页](https://www.huashunxinan.net/product-fobrain)
- **事实：** 历史项目目标是重构已有聊天优先 Workbench、恢复既有能力并比较架构收益，而不是基于用户研究探索新产品（`docs/00-product-brief.md:3-15`）。
- **事实：** 历史需求把连续对话、工具卡、审批、澄清、审计、回放、Action API 和三栏界面列为产品目标（`docs/00-product-brief.md:17-33`），但把成功定义为这些链路能跑通和视觉可比（`docs/00-product-brief.md:45-57`；`docs/01-product-requirements.md:88-94`）。
- **事实：** 历史 FOBrain 范围主要是 24 个只读工具、连接状态、凭据展示、实体澄清、工单状态更新审批及真环境验收（`docs/01-product-requirements.md:59-77`）。
- **事实：** 当前用户、任务、购买者、基线指标、目标值、动作权限和产品关系仍被明确记录为缺失或待确认（`_bmad-output/planning-artifacts/product-blueprint/source-evidence-ledger.md:73-98`；`_bmad-output/forge/agent-workbench-product/user-problem-hypotheses.md:86-96`）。
- **厂商公开主张（PRT-05 / `VENDOR_CLAIMED`）：** Microsoft Security Copilot 官方文档描述自然语言、组织数据 grounding、跨产品插件、优先风险、报告、受权限约束的动作和人工复核，并提示 AI 输出可能不准确、不完整或过时。[Microsoft Security Copilot](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot)；[Microsoft Security Copilot agents application card](https://learn.microsoft.com/en-us/copilot/security/security-copilot-application-card-agents)

## 对抗性发现

- “AI 协作与决策入口”描述的是交互形态和系统位置，不是可验证的用户结果；当前仍未选定单一重要任务、准入路径、基线、目标值和衡量周期。**推断：** 在这些变量缺失时，“入口”无法证明是独立产品、付费模块，还是 FOBrain 的普通界面增强。
- FOBrain 已把资产关联融合、漏洞降噪、动态优先级和闭环处置作为自身价值与能力。[FOBrain 官方产品页](https://www.huashunxinan.net/product-fobrain) **推断：** 当前九个候选能力簇中的职责范围、数据质量、优先级、变化待办、协作交接、受控动作和结果验证均可能重做 FOBrain 已有产品面，而不是形成 AI 增量（`_bmad-output/brainstorming/brainstorm-security-risk-operations-capability-space-2026-07-11/brainstorm-intent.md:19-49`）。
- 早期材料曾把“逐页筛选、人工拼接上下文和线下解释”作为机会，但没有 FOBrain 用户观察或使用数据证明这些行为在 FOBrain 部署后仍然发生；当前问题稿已将其降级为待证任务（`_bmad-output/forge/agent-workbench-product/user-problem-hypotheses.md:86-90`）。**推断：** 这可能只是从通用安全运营痛点倒推出来的故事，而不是 FOBrain 的剩余问题。
- “跨资产、漏洞、业务、人员和工单”目前是跨对象，不等于跨系统；历史范围只明确了 FOBrain 工具和一个外部 Action API（`docs/01-product-requirements.md:52-75`）。**推断：** 在没有第二个业务事实源和真实跨系统流程之前，“跨系统上下文重建”属于过度表述。
- 历史 24 个只读工具按能力恢复门禁组织，而不是按任务频率、用户价值或失败成本组织（`docs/01-product-requirements.md:59-77`；`_bmad-output/planning-artifacts/product-blueprint/source-evidence-ledger.md:49-58`）。**推断：** 用这些工具拼出聊天体验，只能证明集成存在，不能证明产品需求正确。
- “风险案件”被建议为核心产品单位，但 FOBrain 官方已经声称存在漏洞闭环处置。[FOBrain 官方产品页](https://www.huashunxinan.net/product-fobrain) **推断：** 在没有逐字段、逐状态、逐责任边界的对象对比前，新建风险案件极可能制造第二套状态和第二个事实源。
- 历史写域能力仅明确到“工单状态更新审批”（`docs/01-product-requirements.md:67-75`）。**推断：** 这不足以支撑“处置决策协作”“结果验证”或完整人机行动闭环；把一个状态变更描述成风险处置会高估实际用户结果。
- 自然语言并不自动优于筛选器、保存视图、批量操作或固定报表；当前材料没有同任务对照数据（`_bmad-output/planning-artifacts/product-blueprint/research/research-scope-proposal.md:15-30`）。**推断：** 对重复且结构化的运营任务，聊天可能增加表达不确定性、等待时间和复核成本。
- “证据可追溯、事实/推断/未知区分、审计回放”是重要可信约束，但历史需求已经包含工具证据、审计和回放（`docs/00-product-brief.md:19-33`；`docs/01-product-requirements.md:24-50`）。**推断：** 这些能力能降低 AI 风险，却不能单独证明客户会采用或付费。
- 候选用户同时包含漏洞运营人员、资产责任人、安全负责人和审计角色，但没有任何一类被确认（`_bmad-output/forge/agent-workbench-product/user-problem-hypotheses.md:20-64`）。**推断：** 四类角色的语言、权限、频率和成功标准不同，把它们放进同一个首版产品会导致范围膨胀并削弱首要任务。
- 当前“决策”没有定义决策权、输入事实、决策规则、可接受错误、升级条件或最终责任人（`_bmad-output/brainstorming/brainstorm-security-risk-operations-capability-space-2026-07-11/brainstorm-intent.md:66-76`）。**推断：** 产品可能只能生成解释性文本，却被误认为改善了风险判断质量。
- PRT-05／`VENDOR_CLAIMED` 描述 Microsoft 已把自然语言调查、数据 grounding、优先风险、报告、插件扩展和人工监督动作纳入 Security Copilot。[Microsoft Security Copilot](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot) **推断：** “自然语言 + 上下文 + 证据 + 人在回路”是相邻品类的公开能力组合，不足以构成 FOBrain 新模块的差异化。
- PRT-05／`VENDOR_CLAIMED` 也记录 Microsoft 对 AI 输出可能不准确、不完整或过时的官方提示。[Microsoft Security Copilot agents application card](https://learn.microsoft.com/en-us/copilot/security/security-copilot-application-card-agents) **推断：** 如果 FOBrain 的现有结构化排序已可完成任务，引入 AI 可能增加复核成本；实际净结果必须按 A/B 路径验证。
- 当前材料没有证明 FOBrain API/工具能够返回决策所需的责任关系、业务重要性、暴露状态、时间序列、修复后状态和错误语义；这些均被列为待验证（`_bmad-output/planning-artifacts/product-blueprint/research/research-scope-proposal.md:39-47`；`_bmad-output/brainstorming/brainstorm-security-risk-operations-capability-space-2026-07-11/brainstorm-intent.md:66-76`）。**推断：** 数据不足时，模型只能生成更流畅但未必更可靠的解释。
- 历史成功条件是功能链路“能跑通”，不是用户完成任务更快、更正确或更少升级（`docs/01-product-requirements.md:88-94`）。**推断：** 如果继续沿用旧验收方式，团队可能交付一个工程上完整、产品上不可证伪的 Workbench。
- 产品商业关系未定义：没有购买者、预算来源、部署模式、模块定价、续费理由或 FOBrain 产品线负责人承诺（`_bmad-output/planning-artifacts/product-blueprint/source-evidence-ledger.md:77-88`）。**推断：** 即使个别用户喜欢聊天，也不代表它能成为独立产品或可售增值模块。
- “入口”可能把用户从 FOBrain 的现有工作上下文中拉到独立 Workbench；PRT-05 的厂商公开说明展示了相邻形态线索，但入口仍需按任务测试。[Microsoft Security Copilot](https://learn.microsoft.com/en-us/copilot/security/microsoft-security-copilot) **推断：** 内嵌、独立、非 AI 和不做都应由 DG-06 决定。
- 九个能力簇包含权限、证据、数据质量、优先级、待办、协作、动作、验证和学习报告（`_bmad-output/brainstorming/brainstorm-security-risk-operations-capability-space-2026-07-11/brainstorm-intent.md:19-57`）。**推断：** 这是一个完整运营平台范围，不是验证 AI 增量的最小楔子；在首要任务未证实时继续细化会把假设固化成需求。

## 重复能力审计

| 候选能力 | FOBrain 官方声明／历史证据 | 红队判断 | 进入需求前必须证明 |
| --- | --- | --- | --- |
| 资产与漏洞归一、资产关系 | FOBrain 官方声称多源资产关联融合、资产清单和漏洞去重归一。[来源](https://www.huashunxinan.net/product-fobrain) | **声明范围直接重叠。** 目标部署可用性与任务效果待核验；不能把名称本身当差异化。 | 目标版本是否可用、用户是否有效使用；AI 只消费既有事实，不新建平行资产/漏洞事实。 |
| 风险动态优先级 | FOBrain 官方声称通过 VPT 融合漏洞与资产信息并动态排序。[来源](https://www.huashunxinan.net/product-fobrain) | **声明范围直接重叠。** “AI 可解释”也不能暗示重新计算更优排序。 | 目标部署是否提供该结果；用户是否因现有排序不可解释而做额外人工工作；解释是否改善正确决策。 |
| 漏洞闭环、责任与状态跟踪 | FOBrain 官方声称完整闭环；历史 AI 集成仅明确工单状态查询与审批更新。[来源](https://www.huashunxinan.net/product-fobrain)；`docs/01-product-requirements.md:67-75` | **声明范围高度重叠，但当前集成边界未核实。** 新风险案件可能形成第二套真相。 | 目标版本的闭环对象、状态、评论、证据和责任字段，以及用户实际缺口。 |
| 自然语言复杂查询 | 历史 Workbench 已支持聊天与工具调用（`docs/00-product-brief.md:3-5,17-23`）。 | **实现已存在，价值未验证。** | 与 FOBrain 当前页面、报表和保存视图的同任务时间、错误率和学习成本对照。 |
| 实体歧义澄清 | 历史 P2 已列实体消歧（`docs/01-product-requirements.md:67-75`）。 | **可能是增量 UX，也可能只是工具参数缺陷的补丁。** | 真实用户任务中歧义发生率、失败成本及现有 FOBrain 的处理方式。 |
| 证据解释与不确定性 | 历史需求已有工具证据、结构化卡片、审计和回放（`docs/01-product-requirements.md:24-50`）。 | **可信门禁，不是独立价值证明。** | 来源完整率、用户判断正确率和复核时间是否改善。 |
| 受控动作 | 历史只确认工单状态更新审批（`docs/01-product-requirements.md:67-75`）。 | **能力窄于定位承诺。** | 哪个真实动作、谁有权批准、失败如何恢复、外部系统如何确认结果。 |
| 风险案件与决策记录 | 当前仅为 brainstorming 中心洞察（`_bmad-output/brainstorming/brainstorm-security-risk-operations-capability-space-2026-07-11/brainstorm-intent.md:11-17`）。 | **未经验证且可能重复闭环对象。** | 用户是否需要独立对象；为什么 FOBrain 工单/漏洞处置记录不足。 |
| 审计与回放 | 历史产品已经要求审计和回放（`docs/00-product-brief.md:21-33`）。 | **治理能力，不是首要任务。** | 哪个角色以何频率使用、法律/内控依据、保留范围和可接受成本。 |

## 最危险的假设

| 假设 | 为什么危险 | 必须获得的证据 |
| --- | --- | --- |
| FOBrain 上线后仍有一个满足预注册准入路径的重要剩余任务 | FOBrain 官方已经声称解决资产孤岛、漏洞处理人力和闭环问题。[来源](https://www.huashunxinan.net/product-fobrain) | 活跃客户的真实工作观察；效率型记录频率与成本，高损失型记录已发生／演练事件、严重后果和可客观评价性。 |
| 自然语言比现有页面、报表或脚本更快且不更易错 | 目前只有交互愿望，没有同任务基线（`research-scope-proposal.md:15-30`）。 | 交叉对照任务实验；时间、完成率、错误率和复核负担。 |
| FOBrain 数据足以支撑可信决策解释 | 关键字段、时间序列、权限上下文和错误语义仍待验证（`research-scope-proposal.md:39-47`）。 | 真实 API/数据字典、覆盖率、时效、缺失率、冲突率和权限测试。 |
| AI 解释提高判断质量，而不只是提高主观信心 | PRT-05／`VENDOR_CLAIMED` 记录 Microsoft 的风险提示；不证明本项目效果。[来源](https://learn.microsoft.com/en-us/copilot/security/security-copilot-application-card-agents) | 独立专家金标准、盲测、证据引用准确率、错误建议和过度信任观察。 |
| 漏洞/安全运营人员是首要用户 | 当前只是候选，尚未由产品所有者或用户证据确认（`user-problem-hypotheses.md:20-64`）。 | 角色任务分布、使用频率、权限、预算影响力和首要痛点比较。 |
| “风险案件”不是 FOBrain 闭环对象的重复品 | FOBrain 已声称提供漏洞闭环。[来源](https://www.huashunxinan.net/product-fobrain) | 两者对象、状态、责任、证据、审批和审计模型的差距清单。 |
| 客户愿意为 AI 入口单独采用或付费 | 当前没有购买者或商业证据（`source-evidence-ledger.md:77-88`）。 | 设计伙伴承诺、预算/采购路径、付费或续费意向与反对理由。 |
| 独立 Workbench 是正确入口 | 历史三栏 UI 来自重构基线，不是用户研究（`docs/00-product-brief.md:17-25`）。 | 嵌入式、独立式和现有界面的任务对照及长期采用数据。 |

## 可证伪产品实验

以下样本数和阈值是本项目建议的**预先决策规则**，不是行业事实。实验开始前应由产品负责人和 FOBrain 业务负责人确认，避免看到结果后移动门槛。

| 实验 | 做法 | 暂定通过门槛 | 证伪/停止信号 |
| --- | --- | --- | --- |
| E-01 剩余任务发现 | 样本量由 Owner 事前确认；要求活跃用户展示最近完成的真实任务，不先展示解法。采集前同时预注册 A/B 最低证据要求，但不选路径。 | 样本冻结后由 DG-04 按预注册规则只选一路：A 达到最低频率且有可测成本；或 B 有已发生／受控演练、严重后果和足够评价案例。具体阈值均为 `[OWNER-TBD]`。 | 两条路径均无足够证据，或问题已由原生页面／流程、配置、培训低成本解决。 |
| E-02 能力重叠审计 | 让 FOBrain 产品负责人逐项演示候选任务，并核对产品说明、数据字典、接口和闭环对象。 | 明确找到至少一个无法由当前 FOBrain 在合理步骤内完成、且用户证据支持的任务缺口。 | 候选工作流已能在 FOBrain 内以相近时间和质量完成；AI 只换了输入方式。 |
| E-03 同任务交叉对照 | DG-05 先签认主结果、基线、目标、保护指标、样本、权重、胜出和停止规则；随后同一批用户按交换顺序比较所有候选形态。 | 只按 DG-05 事前规则判定；30% 等建议值未被 Owner 事前采用时不构成通过信号。 | 无形态通过、复核抵消收益、错误／越权增加，或原生／非 AI 方案胜出。 |
| E-04 证据与安全盲测 | 用至少 50 个脱敏历史案例建立双人复核金标准，测试事实引用、未知识别、澄清、排序解释和动作建议。 | 关键事实引用正确率至少 95%，高风险无依据建议为 0，所有关键缺失信息均被显式标识或触发澄清。 | 任何未经批准的高风险动作、来源错配、把未知写成事实，或无法稳定复现关键结论。 |
| E-05 入口形态对照 | 对同一任务按 DG-05 规则比较 FOBrain 原生增强、页面内嵌／伴随式 AI、独立 Workbench、非 AI 改进和“不做”。 | 只记录符合所选 A/B 路径主结果与保护指标的胜出项；任一项或“不做”均可胜出。若独立 Workbench 不占优，则取消“独立入口”承诺。 | 无候选形态通过，或原生／非 AI／不做胜出；用户不愿切换上下文。 |
| E-06 真实试点 | 仅开放通过 E-03/E-04 的单一任务，在真实权限下记录使用、人工接管、失败和复核。 | A 效率型按预注册持续使用与净收益判断；B 高损失型按事件／演练就绪度、正确性和响应结果判断；两者不得共用采用指标。 | 采用／就绪证据不足、用户绕回原流程、人工接管居高不下或出现安全事故。 |
| E-07 商业承诺测试 | 向至少 5 个现有/目标客户账户展示已测量结果，不问泛泛“是否喜欢”，而要求指定负责人、数据接入和试点成功指标。 | 至少 2 个账户承诺有负责人和真实数据的试点，其中至少 1 个明确预算来源或与续费/增购挂钩。 | 只有口头兴趣，无人负责、无数据、无预算，或客户认为应是 FOBrain 免费基础功能。 |

## 定位选项比较

| 选项 | 可独立成立的价值假设 | 与 FOBrain 重复风险 | 当前证据 | 红队建议 |
| --- | --- | --- | --- | --- |
| A. FOBrain 页面内的单任务 AI 助手 | 对一个已验证重要任务改善预注册结果，同时保留 FOBrain 为事实与闭环系统。 | 低到中；更像产品增强，不一定能单独收费。 | 与历史工具能力匹配，但没有用户与对照实验（`docs/01-product-requirements.md:59-77`）。 | **候选解法。** 不再预设它优先于原生或非 AI 改进。 |
| B. FOBrain 的独立 AI 协作与决策 Workbench | 多步、跨对象、跨角色任务需要持续上下文、证据、交接和审批，现有页面无法承载。 | 中到高；风险案件、优先级和闭环可能重复 FOBrain。[FOBrain 官方页](https://www.huashunxinan.net/product-fobrain) | 只有假设和历史 UI/运行链路，没有剩余任务、采用或购买证据。 | **当前不可继续按产品定位锁定。** 只有 E-01 至 E-07 证明独立入口优于嵌入式，才升级为候选产品。 |
| C. 数据源中立的独立风险运营平台 | 用统一风险案件跨多个安全/业务系统完成调查、决策与闭环。 | 极高；与 FOBrain 和暴露管理产品的厂商声明范围高度重叠。 | 没有第二数据源、独立购买者、替换意图或迁移证据（`positioning-evidence.md`）。 | **暂缓。** 仅作为长期战略选项；先证明多源客户需求和 FOBrain 之外的事实源。 |
| D. 通用 Agent Workbench 平台 | 让组织创建、部署和治理任意 Agent。 | 与 FOBrain 重复低，但与成熟 Agent 平台竞争且偏离现有能力。 | 历史需求没有 builder、版本、测试、发布和运营治理（`positioning-evidence.md:28-36`）。 | **首版停止。** 除非产品负责人明确换赛道并重新做用户与市场研究。 |

## 继续进入 Product Brief 前的门禁

1. 产品负责人先确认**发现边界**：研究 FOBrain 活跃用户的剩余任务，并允许内嵌、独立、非 AI 和不做，不提前选择最终产品关系。
2. 指定证据负责人，并在接触真实用户／数据前批准研究目的、最小数据、授权、脱敏、访问者、地域、保留／删除和模型禁入项。
3. 完成 E-01 与 E-02，从真实样本中只选一个任务和一种准入路径；记录角色、触发、步骤、基线、成本／后果、替代和结束证据。
4. 在查看实验结果前通过 DG-05，固定主结果、目标、保护指标、样本、胜出和停止规则。
5. 完成 E-03／E-04／E-05，按同一规则决定原生增强、内嵌、独立 Workbench、非 AI 改进或不做。
6. 从 DG-04 后即收集内部采用、设计伙伴、经济买方与预算证据；结合胜出形态决定 FOBrain 标准功能、增值模块、内部工具、独立产品或停止，不能先选性质再找证据。
7. final Product Brief 前还必须明确动作边界和适用边界；详细控制进入 PRD，但高层部署／数据／行业边界不能留空。

在以上门禁通过前，最准确的表述不是“我们正在开发 FOBrain 的 AI 协作与决策入口”，而是：

> **我们准备请求验证：FOBrain 活跃用户是否存在一个未被现有产品充分解决的重要任务，以及哪种解法（包括不做）按预注册规则胜出。**

最终状态：**PIVOT — pending owner/user evidence**。
