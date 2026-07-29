---
status: draft
owner_confirmation: dg_01_approved
artifact_type: cross-source-research-synthesis
formal_bmad_market_research_complete: false
formal_bmad_domain_research_complete: false
formal_bmad_continue_confirmation: confirmed_2026_07_11
decision_state: formal_scope_confirmed_pending_dg_02
updated: 2026-07-11
---

# 需求前置研究综合草案

## 0. 状态与结论

**当前判定：`FORMAL_SCOPE_CONFIRMED_PENDING_DG_02`，不是 `DATA_ACCESS_AUTHORIZED`，更不是 `GO_TO_BUILD`。**

本文件只收敛市场、领域、本地历史证据、定位红队和任务假设中相互一致的结论，不替代任何原始研究报告。产品负责人尚未逐步选择 BMAD 的 `[C]`，因此：

- 正式 BMAD Market Research **未完成**；
- 正式 BMAD Domain Research **未完成**；
- 本文件不得伪造 `stepsCompleted`，不得被当作正式研究签署结果；
- Product Brief、PRFAQ 和 PRD 目前只能预建带证据状态的草案骨架；Product Brief 只有在第 3 节前置证据门禁全部通过后才能标记为 final；
- DG-01 已确认 FOBrain 仅为第一个可验证业务样本，正式 BMAD Market／Domain Research 范围 `[C]` 也已取得；DG-02 通过前仍不得接触真实用户或敏感数据。获准后先验证一个重要剩余任务，不预设 AI、产品形态、功能数量、工具／数据源范围或开发顺序。

流程状态见[研究范围提案](./research-scope-proposal.md)、[市场研究草案](./market-research-draft.md)和[领域研究草案](./domain-research-draft.md)。

## 1. 跨来源证据结论

### 1.1 Owner 决策与当前项目状态

| ID | 事项 | Canonical evidence record | 边界与来源 |
| --- | --- | --- | --- |
| C-01 | 现有 `docs/` 已降级为历史参考，不能直接成为新产品需求真相。 | `OWNER_DECISION / DIRECT / CONFIRMED / SOURCE_ONLY / OWNER_DECISION_RECORDED / NOT_NORMATIVE` | 见[来源证据台账](../source-evidence-ledger.md)。 |
| C-02 | 需求蓝图完成前不进入架构设计或代码实现。 | `OWNER_DECISION / DIRECT / CONFIRMED / SOURCE_ONLY / OWNER_DECISION_RECORDED / NOT_NORMATIVE` | 见[来源证据台账](../source-evidence-ledger.md)。 |
| C-03 | 每个 Product Issue 必须有 Background、Goal 和客观可验证的 Given/When/Then；每个业务目标至少对应一条 AC。 | `OWNER_DECISION / DIRECT / CONFIRMED / SOURCE_ONLY / OWNER_DECISION_RECORDED / NOT_NORMATIVE` | 见[来源证据台账](../source-evidence-ledger.md)。 |
| C-04 | 当前没有产品负责人对正式市场／领域研究步骤的 `[C]` 确认。 | `PROJECT_HISTORY / DIRECT / CONFIRMED / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 三份研究文档 frontmatter 一致；这不是产品事实。 |
| C-05 | 当前八个候选任务都没有真实用户证据，均未达到已验证需求。 | `PROJECT_HISTORY / DIRECT / CONFIRMED / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | [单任务假设矩阵](./task-hypothesis-matrix.md)与[验证计划](./product-validation-evidence-plan.md)一致。 |

字段顺序固定为：`source_type / claim_support / target_applicability / maturity / decision_status / normative_force`。

### 1.2 有来源支持、但未形成产品决策的结论

| ID | 有支持的结论 | Canonical evidence record | 证据边界 |
| --- | --- | --- | --- |
| S-01 | NIST 描述企业补丁管理／修复的识别、确定优先级、获取、安装、验证和持续监控流程。 | `AUTHORITY_SOURCE / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / GUIDANCE` | [NIST SP 800-40 Rev.4](https://csrc.nist.gov/pubs/sp/800-40/r4/final)；是指南，不是本项目法定义务，也不能外推整个漏洞闭环或用户痛点。 |
| S-02 | CVSS Base、EPSS、KEV、SSVC 与组织风险不是同一个概念，不能混成无来源的单一权威风险分。 | `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 各来源规范效力必须逐项记录；本行是跨来源综合，本身不具有规范效力。 |
| S-03 | FOBrain 官方公开声明覆盖多源资产／漏洞归一、动态优先级和闭环处置。 | `VENDOR_SOURCE / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 即 `VENDOR_CLAIMED`；目标部署、集成和任务效果必须分别取证。 |
| S-04 | 历史本地材料提供部分资产、漏洞、业务、责任范围、汇总、工单查询与一次写操作线索。 | `PROJECT_HISTORY / DIRECT / NOT_APPLICABLE / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 本行只描述历史库存；不是 `DEPLOYMENT_AVAILABLE` 或生产可行性证明，也不评价目标产品适用性。 |
| S-05 | NIST／OWASP 材料支持把来源、权限、人工覆盖、监控、事件处理和停用作为治理候选。 | `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 本行是跨来源综合，本身不具有规范效力；NIST 的 `GUIDANCE` 与 OWASP 建议必须在各自来源记录中分别保留，不是所有部署的强制法律 Must。 |
| S-06 | 相邻厂商公开主张提供自然语言、安全数据 grounding、风险解释、工作流和受控动作。 | `VENDOR_SOURCE / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 即 `VENDOR_CLAIMED`；不证明部署、集成、采用、效果或本项目必须跟进。 |
| S-07 | 历史能力更像 FOBrain 查询／交互层，而不是完整的通用 Agent 构建平台。 | `PROJECT_HISTORY / SYNTHESIS / NOT_APPLICABLE / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 本行只描述历史库存；不评价目标产品适用性，换赛道仍需重新研究。 |

所有 FOBrain 和竞品结论统一使用以下证据梯级，完整定义见[证据状态模型](./evidence-status-model.md)：

1. `VENDOR_CLAIMED`；
2. `DEPLOYMENT_AVAILABLE`；
3. `INTEGRATION_AVAILABLE`；
4. `TASK_EFFECTIVE`；
5. 形成最终定位还必须取得 `ADOPTION_EVIDENCED` 和 `OWNER_DECISION_RECORDED`。

任何结论都不得跳过中间成熟度；产品定位还必须额外满足采用证据和 Owner 决策。

### 1.3 HYPOTHESIS

以下只能用于研究设计，不得写成 Product Issue 的 Background 或 Goal：

| ID | 待验证假设 | Canonical evidence record | 最小证据 |
| --- | --- | --- | --- |
| H-01 | 漏洞／安全运营分析员是候选主用户。 | `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 角色任务分布、近期真实工作样本、频率、权限和失败成本比较。 |
| H-02 | FOBrain 原生流程之后仍存在效率型高频任务，或低频但高损失且可复现的任务。 | `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 现场观察、原生基线，并按 DG-04 选择一种准入路径。 |
| H-03 | 可追溯 AI 能改善单任务结果，且不会把收益转移为更高核验、纠错和信任成本。 | `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 预注册的同任务对照和盲评。 |
| H-04 | 事实／推断／未知／冲突分层与证据链是用户采用 AI 所必需。 | `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 对照概念测试与真实任务观察。 |
| H-05 | “风险案件”比会话、FOBrain 原有对象或工单更适合承载跨时段工作。 | `PROJECT_HISTORY / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 对象、状态、责任和事实归属差距审计。 |
| H-06 | 某种 AI 形态优于 FOBrain 原生增强、非 AI 改进和不做。 | `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 同任务、同数据、同用户、同判定规则的形态对照。 |
| H-07 | 客户或内部采用方愿意投入数据、用户时间、安全评审、责任资源和相应预算。 | `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` | 与候选产品性质匹配的采用／买方证据。 |

### 1.4 REJECTED AS CONFIRMED CLAIM

下表不是说所有候选命题都已被事实证伪，而是说明它们**不能以当前证据升级为已确认产品主张**。每行明确区分“直接矛盾”“无证据”“错误推理”或“解法未验证”。

| ID | 不得作为已确认事实的表述 | `claim_support / decision_status` | 原因码 | 为什么不能继续作为确认项 |
| --- | --- | --- | --- | --- |
| X-01 | “24 个历史工具就是 24 个产品需求或首版完整范围。” | `NONE / REJECTED` | `UNSUPPORTED_SCOPE_EQUIVALENCE` | 工具按旧能力恢复组织，没有用户任务、频率、成本或结果证据。 |
| X-02 | “资产归一、漏洞归一、动态排序和闭环处置是新 AI 产品的独有差异化。” | `CONTRADICTED / REJECTED` | `VENDOR_CLAIM_OVERLAP` | 仅“独有”被 FOBrain 的 `VENDOR_CLAIMED` 直接反驳；目标部署和真实任务是否仍有问题保持 `UNKNOWN`。 |
| X-03 | “当前历史范围足以定义一个通用 Agent 平台。” | `NONE / REJECTED` | `INSUFFICIENT_SCOPE` | 历史范围没有平台构建、发布、生命周期治理与运营产品面；若换赛道须重新研究。 |
| X-04 | “行业普遍存在漏洞运营痛点，所以本项目对应能力应进入 Must。” | `NONE / REJECTED` | `INVALID_INFERENCE` | 行业资料不能证明 FOBrain 用户使用现有产品后的剩余任务。 |
| X-05 | “FOBrain／竞品官网声明等于目标部署已经可用，并且任务有效。” | `NONE / REJECTED` | `UNSUPPORTED_INFERENCE` | 官网没有提供目标部署与任务效果证据；这项等同推理无支持，但当前也没有目标部署反证。 |
| X-07 | “工单关闭等于风险已经消除。” | `CONTRADICTED / REJECTED` | `INVALID_STATE_EQUIVALENCE` | 对声称已修复的补丁／修复，指南要求新验证事实；其他关闭处置也有独立语义。 |

“聊天本身有独立价值”和“风险案件应成为首版对象”继续保留为 H-03／H-05 的可证伪 `HYPOTHESIS`，不属于已被反证的命题。

### 1.5 OPEN / UNKNOWN

下表统一使用 `NO_SOURCE / NONE / UNKNOWN / SOURCE_ONLY / OPEN / UNKNOWN`，直到每项取得相应证据。

| ID | 未知事项 | 解除方式 |
| --- | --- | --- |
| U-01 | 新项目与 FOBrain 的产品关系。 | 先完成 DG-05 判定规则、DG-06 解法形态对照和采用／买方取证，再由产品负责人在 DG-07 做唯一选择。 |
| U-02 | 唯一主用户、协作用户、买方、审批人和风险责任人。 | 产品负责人定研究对象，随后用真实角色证据校正。 |
| U-03 | FOBrain 用户最重要的一个剩余任务及其频率、时长、失败成本和当前替代。 | 真实任务观察。 |
| U-04 | 首要用户结果、业务结果、基线、目标值和衡量周期。 | 产品负责人预注册口径，用户实验测量。 |
| U-05 | 目标客户当前 FOBrain 版本、配置、页面、AI 能力、数据源和实际使用深度。 | 产品演示、使用记录和部署审计。 |
| U-06 | 稳定实体 ID、责任／部门映射、来源、时间、错误语义和权限是否足以支撑候选任务。 | API 合同、真实响应、异常样本和权限测试。 |
| U-07 | FOBrain VPT 是否提供可解释因素、规则版本和时间语义。 | 官方演示和真实数据验证。 |
| U-08 | 工单、补丁／修复、缓解、风险接受、不适用／误报、延期、证据不足和风险仍存在是否都有可核对的状态、理由、责任人和证据。 | 真实工单、授权决定和前后事实样本。 |
| U-09 | 内嵌、独立 Workbench、非 AI 改进或“不做”哪个结果更优。 | 同任务对照。 |
| U-10 | Action API、移动端、Inspector、回放和管理报告是否有真实用户。 | 角色与场景证据。 |
| U-11 | 数据类别、个人信息、跨境、行业监管和生成式 AI 规则在目标部署中的具体适用性。 | 法务／合规按客户、地域、行业、数据与服务对象确认；通用研究不是法律意见。 |
| U-12 | 独立采购、增购／续费、预算来源和设计伙伴承诺。 | 客户与经济买方证据。 |

## 2. 八个任务假设摘要（不排序）

### 2.1 评级说明

- “证据”沿用[单任务假设矩阵](./task-hypothesis-matrix.md)的 `L2/L1/L0`，没有任何 `L3/已验证`。
- “本地技术证据摘要”只描述局部历史线索、概念可探索或关键证据不存在，不创建第二套状态，也不是可行性评级。正式状态仍使用 `SOURCE_ONLY / DEPLOYMENT_AVAILABLE / INTEGRATION_AVAILABLE / TASK_EFFECTIVE`。
- “重叠”只衡量与 [FOBrain 的 `VENDOR_CLAIMED` 能力](https://www.huashunxinan.net/product-fobrain)的公开表述重叠，不代表目标部署／集成可用，更不代表任务效果。
- 表格只摘要历史覆盖与待证缺口，不表示研究顺序、DG-04 分数、开发范围或开发优先级。

| 任务 | 历史材料覆盖度 | 本地历史线索摘要（非可行性评级） | 与 FOBrain `VENDOR_CLAIMED` 重叠 | 待证事项 |
| --- | --- | --- | --- | --- |
| T-01 核对“我／本部门负责什么风险” | L1 | 存在查询线索；归档样本的部门与人员稳定映射不足 | 高 | 用户价值、目标身份／权限和实际后果。 |
| T-02 单 IP／资产可复核调查包 | L2 | 查询和详情有局部历史线索；跨对象 ID、来源、时间和责任关系未证实 | 高 | 目标原生流程、A/B 准入与客观结果。 |
| T-03 单漏洞／CVE 影响范围与责任人 | L2 | 漏洞详情和关联有局部历史线索；责任链、时效和身份映射未证实 | 很高 | 跨对象关联是否仍为真实剩余任务。 |
| T-04 解释既有优先级“为什么现在处理” | L1 | 归档资料没有 VPT 因素与规则关键证据 | 声明范围高重叠 | 解释数据、用户问题与正确决策结果。 |
| T-05 关键业务系统风险简报 | L1 | 有聚合线索；缺趋势、业务影响、口径和消费角色 | 高 | 制作者／消费者、实际决策与结果标准。 |
| T-06 调查结果交接给责任人 | L1 | 最多有交接内容线索；缺创建、分派、评论、确认和状态证据 | 高 | 目标闭环中是否确有上下文丢失。 |
| T-07 人工确认后更新既有工单状态 | L1 | 归档契约曾列状态更新；归档计划记录当时缺代表性工单样本 | 声明范围完全重叠 | 目标样本、必要性、动作安全与任务净收益。 |
| T-08 验证关闭结果 | L0 | 归档资料缺补丁／修复、缓解、风险接受、误报／不适用、延期及复测证据 | 声明范围高重叠 | 真实关闭分类、反例和客观验证证据。 |

### 2.2 候选任务不预排序

T-01～T-08 都只是访谈探针，当前没有用户证据支持预选 T-02、T-03、T-04 或任何其他任务。DG-03 后先采集并冻结近期真实任务样本，再由 DG-04 按采集前预注册的 A/B 规则只选择**一个任务**，或选择“无合格任务”。同一任务可以使用多个真实案例，但多个案例不是多个任务，也不得把不同任务案例合并成一个效果数字。

若选定任务被目标部署原生能力、配置、培训或更简单的非 AI 改进按预注册规则完成，就停止该任务，不按预设顺序自动顺延成另一个任务或开发需求；重新选择必须重新通过 DG-04。

## 3. Product Brief 入场门禁

Product Brief 可以提前预建草案骨架，但这不等于进入或完成 Product Brief。**只有以下前置证据门禁全部通过后，Product Brief 才能标记为 final**；任一项未通过都必须保留 `draft` 和证据状态。

| Gate | 必须交付的证据 | 当前状态 |
| --- | --- | --- |
| PB-00 正式流程 | 产品负责人按第 6 节 DG-01 至 DG-09 留下逐项决定，并在 BMAD Market／Domain Research 的实际步骤中选择 `[C]`；本表不能代替 `[C]`。 | 部分通过：DG-01 与正式范围 `[C]` 已通过，DG-02～DG-09 未通过 |
| PB-01 发现边界与证据访问 | 已固定“FOBrain 是第一个可验证样本而非最终产品边界”，并取得目标部署、接口、真实样本、用户和责任人的授权访问。 | 部分通过：发现边界已确认，证据访问未批准 |
| PB-02 样本与单一任务 | 先完成代表性样本采集，再只选择一个任务；实验案例全部属于该任务，不能混入第二任务。 | 未通过 |
| PB-03 Canonical evidence record | 每项主张分别记录六个 canonical 字段；能力成熟度依次区分 `VENDOR_CLAIMED`、`DEPLOYMENT_AVAILABLE`、`INTEGRATION_AVAILABLE`、`TASK_EFFECTIVE` 和 `ADOPTION_EVIDENCED`。 | 未通过 |
| PB-04 技术可行性实证 | 单一任务所需对象 ID、来源、时间、权限、空结果、错误、冲突和真实样本均由目标部署验证；本地历史线索不能代替。 | 未通过 |
| PB-05 核心结果与判定规则 | 在实验前确定一个用户可观察结果、现状基线、目标值、衡量周期、保护指标和停止阈值；不能用对话量或工具调用量替代。 | 未通过 |
| PB-06 单任务净收益与形态 | 同一任务、多个案例按 PB-05 的同一规则比较原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 和不做，只记录胜出形态或“不做”；关键事实错误、复核负担和越权风险不得恶化。 | 未通过 |
| PB-07 产品关系与采用证据 | 在决定内部工具、FOBrain 标准功能、增值模块、独立产品或停止前，已取得与各选项相匹配的使用责任、设计伙伴、经济买方、预算或内部采用证据。 | 未通过 |
| PB-08 动作边界 | 已确定只读／写边界；若任务涉及处置结果，分别为 verified remediation、mitigated、risk accepted、not applicable/false positive、deferred、insufficient evidence、still present 定义责任与证据。 | 未通过 |
| PB-09 适用性 | 人机责任、数据流、部署地域、保留、模型接收方和法规适用已由对应责任人审查；受限适用指南不得写成所有客户义务。 | 未通过 |

**当前 Product Brief 入场结果：`NOT READY — OWNER / USER / DATA / EXPERIMENT EVIDENCE MISSING`。**

## 4. 证伪与停止条件

数值门槛必须在实验前由产品负责人确认。现有研究中的 `30%` 时间改善、`95%` 关键事实引用正确率、真实用户／组织样本数等都只是建议工作门槛，不是行业事实。一次实验只评估一个任务，可以包含该任务的多个案例；若要更换任务，必须结束当前实验并重新通过 DG-04。

| ID | 停止信号 | 停止什么 |
| --- | --- | --- |
| ST-01 | 按 DG-04 预注册的效率型或高损失型准入路径，都找不到足够真实证据的剩余任务。 | 停止当前 FOBrain AI 产品机会，不转而凭空扩大为通用平台。 |
| ST-02 | FOBrain 原生配置、培训、保存视图、报表或小幅非 AI 增强可用相近成本与质量完成任务。 | 停止对应 AI 任务。 |
| ST-03 | AI 节省的查询时间被补问、复核、纠错和回到原系统的时间抵消。 | 停止自然语言作为该任务主要入口。 |
| ST-04 | 关键事实错误增加、来源错配、未知被写成事实，或高风险建议无法稳定复现。 | 停止该 AI 结果用于决策；回退到导航／只读查询或取消。 |
| ST-05 | 任务所需稳定 ID、来源、时间、权限或错误语义不可取得。 | 停止“证据化调查／决策”承诺，不用生成内容掩盖数据缺口。 |
| ST-06 | 当前用户、部门或责任范围只能依赖不可靠姓名 fallback。 | 停止自动定责和“我的范围”产品承诺。 |
| ST-07 | 独立 Workbench 不优于 FOBrain 内嵌形态，或用户拒绝第二工作界面。 | 停止独立 Workbench 定位，保留内嵌或不做结论。 |
| ST-08 | 只有 FOBrain 单一数据源、单一购买关系，且没有第二数据源／客户的独立需求。 | 停止“数据源中立平台”与独立平台声明。 |
| ST-09 | 客户只有口头兴趣，不愿提供负责人、真实数据、用户时间、安全评审或预算路径。 | 停止独立商业产品声明。 |
| ST-10 | 任何写操作出现未经授权执行、审批绕过、范围错选、重复执行或无法核对真实回执。 | 立即停止写路径，恢复只读。 |
| ST-11 | 目标部署无法满足数据分类、最小权限、地域／跨境、日志保留或行业监管要求。 | 停止生产数据试点，直到责任人确认可接受方案。 |
| ST-12 | 四周试点中新鲜感消退后用户回到原流程，或人工接管率不降。 | 停止扩大范围，重新判断任务是否存在。 |

## 5. 对 PRD 的硬约束

1. **证据状态先于范围。** 只有通过 Product Brief 门禁的用户任务才能进入 final PRD 的 Must；其余统一标记为 `HYPOTHESIS` 或移入研究附录。
2. **每个 Product Issue 可追溯。** 必须从真实任务证据映射到 Background、业务 Goal、输入／输出／异常与 Given/When/Then；每个 Goal 至少一条客观 AC。
3. **行业背景不能替代本产品问题。** NIST、CISA、FIRST、OWASP 和厂商资料只能提供领域／治理约束，不能单独生成用户 Must。
4. **厂商声明不得升级。** PRD 必须逐项填写 canonical evidence record，并分开使用 `VENDOR_CLAIMED`、`DEPLOYMENT_AVAILABLE`、`INTEGRATION_AVAILABLE`、`TASK_EFFECTIVE` 与 `ADOPTION_EVIDENCED`；前一层不能证明后一层。
5. **不按工具拆需求。** 禁止把 24 个历史工具或接口逐项转换为需求；只能把必要能力映射到已验证任务链。
6. **不重复 FOBrain 事实系统。** 不另建资产、漏洞、优先级或闭环真相；若提出新持久对象，必须证明不会形成第二套状态与责任。
7. **模型不是事实来源。** 每个关键结论必须携带对象、范围、来源、时间、数据质量以及事实／推断／未知／冲突状态。
8. **信号语义不可混淆。** CVSS、EPSS、KEV、SSVC 和 FOBrain VPT 保留原意；禁止无证据综合风险分或静默重算权威排序。
9. **不确定性必须是可验收结果。** 实体不唯一、无权限、数据缺失、来源失败、条件错误、事实冲突和真实不存在必须客观区分；正确结果可以是“无法确定 + 最小补证清单”。
10. **权限由确定性边界执行。** 用户可见范围继承源系统，AI 不得提升读取或写入权限；外部内容视为不可信。
11. **人机责任必须分开。** AI 建议、人员决定、审批、外部系统实际变更和真实回执分别记录；风险接受、法务判断、监管报告和高影响动作不得委托给 AI。
12. **首个实验只评价一个重要任务。** 任务可由 A 高频效率型或 B 低频高损失型准入，允许多个属于同一任务的案例，不允许把第二个任务混入同一实验；DG-08 前所有外部动作只做影子／沙箱评价。写操作只有在该路径的真实任务价值、下游强鉴权、影响预览、独立审批、幂等／补救、真实回执和处置复核全部存在时才能进入后续需求候选。
13. **关闭原因必须显式分类。** 关闭或处置结果至少区分：`verified remediation`（有新的补丁／修复验证事实）、`mitigated`（有补偿控制与残余风险）、`risk accepted`（由授权人决定并有期限／复审）、`not applicable/false positive`（有不适用／误报证据）、`deferred`（有授权理由、责任人和期限）、`insufficient evidence`（当前不能判断）、`still present`（风险事实仍存在）。只有 `verified remediation` 可以在相应证据充分时表述为已修复；工单状态本身不能证明风险消除。NIST 补丁证据只用于补丁／修复验证约束。[NIST SP 800-40 Rev.4](https://csrc.nist.gov/pubs/sp/800/40/r4/final)
14. **法规必须逐部署适用。** PRD 先记录数据类别、目的、地域、接收方、保留和责任人；中国网络、数据、个人信息、漏洞与生成式 AI 规则由法务按实际场景确认。[《网络安全法》](https://www.cac.gov.cn/2025-12/29/c_1768735112911946.htm)、[《数据安全法》](https://www.cac.gov.cn/2021-06/11/c_1624994566919140.htm)、[《个人信息保护法》](https://www.miit.gov.cn/jgsj/zfs/fl/art/2022/art_515a4b20c12f430eab54bb4f56d89f56.html)。
15. **产品形态不预设。** 内嵌助手、独立 Workbench、非 AI 改进和不做都是合法结论；历史三栏 UI、Action API、移动端、Inspector 和回放不能因旧文档存在就进入 Must。
16. **明确首版非目标。** 通用 Agent builder、无人监督自主修复、无证据综合风险分、完整多租户平台和“聊天即唯一案件记录”保持非目标，除非重新完成相应产品研究。
17. **产品约束不由架构补写。** 性能、容量、可用性、数据保留、部署、成本、无障碍和支持范围必须由产品负责人给出可验证目标；PRD 不写 Eino、数据库、接口结构或其他实现解法。

## 6. 产品负责人必须逐项确认的决策门禁

统一使用以下顺序。每项必须有一条明确决定和证据引用；不得跳过、并行锁定后项，或用后项反推前项。这个决策表也不代表 BMAD 流程已经取得 `[C]`。

| 顺序 | 决策门禁 | 必须形成的唯一答案 |
| ---: | --- | --- |
| DG-01 | **发现边界（已通过）** | FOBrain 仅作为第一个可验证业务样本，不是最终产品边界；先验证一个真实高价值任务，再由证据决定最终功能、工具／数据源和产品形态。A／B 路径、主用户、任务、AI、商业关系和“不做”均未预设。正式 BMAD 步骤仍需另行取得 `[C]`。 |
| DG-02 | **证据访问与研究数据批准** | 指定产品／部署、数据、用户研究、安全／法务责任人；确认目标版本、原生流程、API、真实数据、异常样本和用户；在访问前签认研究目的、数据最小集、授权、脱敏、访问者、存储地域、保留／删除和禁止进入模型项。 |
| DG-03 | **研究队列** | 选择第一批观察的具体岗位、职责、权限、工作环境与样本纳入／排除条件；明确研究队列不等于最终主用户。 |
| DG-04 | **首个实验任务** | 引用任务采集前已冻结的 A/B 两套准入规则版本；样本冻结后只应用固定规则选择一个可客观复核任务和 A 或 B 唯一路径，也可选择“无合格任务”。实验可含多个同任务案例，不得混入第二任务。 |
| DG-05 | **核心结果与判定规则** | 在实验前唯一确定一个用户可观察结果、现状基线、目标值、衡量周期、保护指标和停止阈值；看到实验结果后不得修改规则。 |
| DG-06 | **解法形态** | 按 DG-05 的同一规则，对 DG-04 的同一任务比较 FOBrain 原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 改进和不做；只选择证据支持的结果。 |
| DG-07 | **产品关系与性质** | 结合 DG-06 胜出形态和已在 DG-04 后收集的采用／买方证据，唯一确定内部工具、FOBrain 标准功能、可选增值模块、独立销售产品、多个工具／数据源的产品、通用 Agent 平台或停止；若超出首个样本可支持的外推范围，必须补做对应研究，不能直接继承结论。 |
| DG-08 | **动作边界** | 唯一确定首版只读建议，或允许一个明确的受控外部动作；若允许写操作，逐项确认责任人、审批人、下游鉴权、影响预览、补救、真实回执和复核。若涉及处置结果，按关闭原因分别定义证据。 |
| DG-09 | **适用边界** | 由责任人按目标客户、地域、行业、数据、模型、服务对象和部署方式确认数据／隐私／跨境／漏洞／生成式 AI／行业规则的适用性与组织规则；受限指南不外推。 |

## 7. 允许的下一步

DG-01 与正式 Market／Domain Research 范围 `[C]` 已通过。下一步一次只处理 DG-02 及后续门禁；在 DG-02 前，正式研究仅使用公开资料，不访问真实用户或敏感数据。只有 DG-01 至 DG-09 与 Product Brief 前置证据门禁全部通过后，Product Brief 才能标记为 final；此后才可把经验证的单一任务转成产品蓝图和 Product Issues。

内部证据入口：

- [市场研究草案](./market-research-draft.md)
- [领域研究草案](./domain-research-draft.md)
- [FOBrain 本地可行性证据](./local-fobrain-feasibility-evidence.md)
- [单任务假设矩阵](./task-hypothesis-matrix.md)
- [产品需求证据验证计划](./product-validation-evidence-plan.md)
- [研究范围提案](./research-scope-proposal.md)
- [来源证据台账](../source-evidence-ledger.md)
- [定位证据](../../../forge/agent-workbench-product/positioning-evidence.md)
- [定位红队](../../../forge/agent-workbench-product/positioning-red-team.md)
- [用户问题假设](../../../forge/agent-workbench-product/user-problem-hypotheses.md)
- [脑暴意图](../../../brainstorming/brainstorm-security-risk-operations-capability-space-2026-07-11/brainstorm-intent.md)
