---
status: complete
audit_type: independent-evidence-and-logic-review
audited_at: 2026-07-11
reaudited_at: 2026-07-11
audit_versions:
  - v1
  - v2
  - v3
  - v4
  - v5
reviewed_artifacts:
  - market-research-draft.md
  - domain-research-draft.md
  - local-fobrain-feasibility-evidence.md
  - task-hypothesis-matrix.md
  - product-validation-evidence-plan.md
  - source-evidence-ledger.md
  - positioning-red-team.md
v2_additional_reviewed_artifacts:
  - evidence-status-model.md
  - research-synthesis-draft.md
  - assumptions-and-open-questions.md
  - owner-decision-gates.md
v3_reviewed_scope:
  - _bmad-output/planning-artifacts/product-blueprint/**/*.md
  - _bmad-output/forge/agent-workbench-product/*.md
v4_reviewed_scope:
  - _bmad-output/planning-artifacts/product-blueprint/**/*.md
  - _bmad-output/forge/agent-workbench-product/*.md
  - _bmad-output/brainstorming/brainstorm-security-risk-operations-capability-space-2026-07-11/brainstorm-intent.md
v5_reviewed_scope:
  - V4-L-01 status-and-version-closure
  - V4-L-02 historical-target-applicability-closure
  - discovery-A-B-preregistration-order
  - formal-BMAD-research-status
v1_verdict: NEEDS_WORK
v2_verdict: NEEDS_WORK
v3_verdict: NEEDS_WORK
v4_verdict: PASS_WITH_CONDITIONS
v5_verdict: PASS
verdict: PASS
---

# 研究主张独立审计

## 结论

**NEEDS WORK。** 未发现厂商自报效果被直接写成已独立验证的客户效果，也未发现正式 BMAD 步骤完成状态被伪造；数字样本量和通过阈值也普遍明确标为项目建议。但是，当前材料仍有 8 项会改变后续 Product Brief／PRD 结论的高风险问题：厂商能力声明被升级为“商品化／市场基线／直接重复”，行业通用痛点被赋予产品 P0/P1，产品关系与单任务发现的先后门禁互相冲突，技术可行性评级缺少实测，条件性法规义务被写成直接需求，以及补丁管理证据被扩大为全部漏洞关闭语义。

本审计只判断证据强度、逻辑保真和下游误用风险，不判断最终应建设哪种产品。

## CRITICAL

无。

## HIGH

### H-01 多个厂商“公开声称具备”不足以证明能力已经商品化或成为采购基线

- **位置：** `market-research-draft.md:53`、`market-research-draft.md:266-302`；`positioning-red-team.md:38`。
- **问题：** 竞品官网可以直接证明“厂商公开提供／宣称某能力”，不能证明客户实际采用深度、可靠性、效果、价格包含关系或买方将其视为必备基线。当前从多个营销页面直接推出“已经商品化”“这些是进入市场的基线组合”，证据跨越了一层。材料自身在 `market-research-draft.md:41-45,74-79,264-266` 已承认厂商资料不能证明效果，因此结论与证据规则不一致。
- **修正建议：** 将结论改为“多家相邻厂商已公开将这些能力纳入产品主张，故不能仅凭功能名称宣称差异化”；“商品化”“采购基线”必须增加独立客户采用、竞品实测、合同／套餐或买方访谈证据。竞争表的证据等级统一为“厂商声明，中”。

### H-02 FOBrain 官方产品页被同时当作公开定位证据和已验证业务能力边界

- **位置：** `market-research-draft.md:31-35,49,107,266,283`；`source-evidence-ledger.md:44-45`；`positioning-red-team.md:48-52`；对照 `domain-research-draft.md:64-72,88` 与 `local-fobrain-feasibility-evidence.md:21-30,42-50`。
- **问题：** FOBrain 官网确实直接声明资产融合、漏洞归一、VPT 动态优先级和闭环处置，但这只确认厂商公开产品范围。它不能证明目标客户版本已授权、已配置、API 已暴露、用户实际使用或效果成立。红队表中的“直接重复”容易被下游理解为已完成现场能力审计，而本地证据反而显示当前部署存在身份映射、工单样本和结果复核缺口。
- **修正建议：** 固定三层状态：`VENDOR_CLAIMED`（官网声明）、`DEPLOYMENT_AVAILABLE`（目标版本演示／合同／配置确认）、`TASK_EFFECTIVE`（真实用户任务实测）。当前只能据第一层阻止把能力名称当差异化，不能据此删除所有相关用户问题。把“直接重复”改成“声明范围直接重叠，部署与任务效果待核验”。

### H-03 行业通用痛点虽被声明“不是 FOBrain 缺口”，却仍以 P0/P1 和“必须”进入产品含义

- **位置：** `market-research-draft.md:49,55-59,150-164` 与 `market-research-draft.md:174-185`。
- **问题：** 前文正确指出 FOBrain 上线后的剩余痛点未知，但第 6 节仍把漏洞量大、数据分散、责任分歧、闭环断裂列为 P0/P1，并导出“必须解释”“必须成为显性状态”“完成必须有复核证据”等产品结论。这会让后续需求提炼器把行业背景误当 FOBrain 活跃用户已经验证的 Product Issue Background／Goal。
- **修正建议：** 第 6 节改名为“行业问题假设与 FOBrain 剩余性待验证项”；去掉 P0/P1，改用 `INDUSTRY_SUPPORTED / FOBRAIN_RESIDUAL_UNKNOWN` 双轴。只有任务观察证明仍存在的痛点才能获得产品优先级并进入 Product Issue。

### H-04 产品关系决策与单任务发现的先后门禁互相冲突

- **位置：** `source-evidence-ledger.md:91-99`、`positioning-red-team.md:96-103`、`product-validation-evidence-plan.md:23-32`，对照 `market-research-draft.md:55-59,325-335,431-466` 和 `domain-research-draft.md:50-55,356-369`。
- **问题：** 台账和红队要求先明确选择“FOBrain 增强／独立模块／独立平台／通用平台”才进入 Brief；市场与领域研究又明确要求先找到一个剩余任务、并行比较内嵌、独立和非 AI 解法，之后才选择产品形态。两套门禁不可能同时满足，会让团队在“定位未定不能研究”与“未经研究不能定定位”之间循环。
- **修正建议：** 固定唯一顺序：`发现边界确认（FOBrain 用户 + 不预设形态） → 单任务证据 → 解法形态对照 → 产品关系决策 → final Brief/PRD`。Owner 早期只确认“发现边界和禁止项”，不提前确认最终产品关系。同步改写 VH-01 和红队第 98 条。

### H-05 “技术可行性”高／中评级没有真实技术验证支撑

- **位置：** `domain-research-draft.md:309-328`；`local-fobrain-feasibility-evidence.md:10-18,21-30,42-50`。
- **问题：** `事实/推断/未知/冲突分离`、`实体歧义澄清`、`人类修改/否决`被评为“高可行”，但当前只有历史契约、候选 UI／输出规则和局部接口映射，没有目标版本 API 样本、端到端演示、模型错误率、权限测试、延迟／成本或用户可用性证据。这里评估的是“概念上可描述”或“实现依赖较少”，不是技术可行性。
- **修正建议：** 删除单一“可行性”评级，拆为 `数据可得性`、`API/权限可执行性`、`模型质量可验证性`、`产品交互可用性`、`运营与合规可接受性` 五列。无真实样本和最小原型测试时最高只能是 `UNKNOWN / PLAUSIBLE`，不能是“高”。

### H-06 法规适用性待确认，却把推导写成“直接转化”“必须”的候选需求

- **位置：** `domain-research-draft.md:198-217` 与 `domain-research-draft.md:219-227`。
- **问题：** 第 7.1 节正确声明适用性依部署主体、行业、数据和服务对象而定；第 7.3 节却使用“直接转化为候选需求”“任何外部模型或连接器启用前必须”等绝对措辞。部分条目是稳健的安全产品政策，但不是所引法规在所有部署中的直接法律义务。这样会把“建议的默认安全控制”和“经法务确认的适用义务”混成一类。
- **修正建议：** 拆成三类：`默认产品安全政策`、`条件性法律义务（附触发事实）`、`法务待判定`。每条法律要求记录适用主体、触发条件、法条、责任角色；未完成适用性评估前，不得在 PRD 中使用“合规必须”措辞。

### H-07 数据出境条目忽略豁免和阈值，容易被读成所有境外处理都必须走同一种监管评估

- **位置：** `domain-research-draft.md:215`、`domain-research-draft.md:221`。
- **问题：** “外部 LLM、境外遥测、跨境支持和备份必须进入数据流与出境评估”作为内部数据流盘点要求是合理的，但引用的《促进和规范数据跨境流动规定》同时设有多类豁免、10 万／100 万个人信息及 1 万敏感个人信息等不同路径，并区分关基运营者、重要数据和普通数据处理者。当前措辞没有区分“内部识别／法务评估”与“申报安全评估／标准合同／认证”。
- **修正建议：** 改成“所有境外接收路径先进入数据流清单；法务根据主体、数据类别、数量、接收方和豁免判断是否需申报安全评估、标准合同、认证或仅履行一般保护义务”。引用官方规定第 3—10 条及答记者问的豁免说明。

### H-08 补丁管理证据被扩大为全部漏洞／暴露关闭语义，T-08 输出状态不完整

- **位置：** `domain-research-draft.md:31-38,90-124,371-384`；`task-hypothesis-matrix.md:132-145`；`market-research-draft.md:156-164,180-185`。
- **问题：** NIST SP 800-40 Rev.4 直接支持的是补丁的识别、排序、获取、安装和验证，不足以独立定义所有攻击面、配置、补偿控制、误报、例外和风险接受的关闭语义。T-08 只给“已验证／仍存在／证据不足”三种结果，会遗漏“已缓解但未消除”“风险已接受”“不适用／误报”“延期且仍受控”等合法结局。
- **修正建议：** 把 NIST 结论限定为补丁／修复类任务；建立由“关闭原因 + 所需验证证据”驱动的状态模型。任何状态都不能只由工单关闭推出，但也不能假定所有关闭都应证明漏洞已消除。

## MEDIUM

### M-01 ISC2 数字本身准确，但“从概念进入采用阶段”扩大了样本含义

- **位置：** `market-research-draft.md:51,112-118`。
- **问题：** ISC2 原文是 28% 已集成、19% 正在测试、22% 处于早期评估，合计 69% “走向常态使用”；63% 是当前使用 AI 安全工具团队的自报生产力提升。把测试和早期评估一并概括为“进入采用阶段”，会让成熟度显得高于数据本身。
- **修正建议：** 直接保留三段比例，并写成“在 ISC2 样本中，评估／测试／集成已较广泛”；将 63% 明确为使用者自报，不能证明客观生产率、因果效果或本产品购买意愿。

### M-02 “首个合理进入点只能是 FOBrain 现有用户”是战略建议，不是研究结论

- **位置：** `market-research-draft.md:57,134-146`。
- **问题：** 本地证据使 FOBrain 活跃用户成为最低集成风险的首个发现样本，但没有证据排除尚未购买 FOBrain 的目标客户、FOBrain 售前账户或其他数据源用户。“只能”会提前限制市场发现。
- **修正建议：** 改为“推荐的第一发现队列是能提供 FOBrain 现状基线的活跃用户”；是否排除其他队列由研究成本和战略边界决定，不写成外部事实。

### M-03 对“任何可调用外部系统的 AI”统一要求人类批准，与文内 L0/L1 分级冲突

- **位置：** `domain-research-draft.md:38,40-46`，对照 `domain-research-draft.md:145-155`。
- **问题：** OWASP Excessive Agency 支持最小权限、下游鉴权，并对高影响动作要求人类批准；它不要求每个只读调用都人工审批。文内 L0 又明确按用户权限自动执行，形成内部矛盾。
- **修正建议：** 将通用约束写成“鉴权、最小权限、审计”；把显式批准限定为写操作、高影响、跨边界、不可逆或策略要求的动作。

### M-04 多来源拼接出的 NIST 结论被标为单一“事实（高）”

- **位置：** `market-research-draft.md:103-110`。
- **问题：** SP 800-40 直接支持补丁管理生命周期，IR 8286D 支持业务影响分析进入风险优先级；“把资产、漏洞、业务影响、资源和修复验证组合成持续流程”是合理综合推断，但不是任一引用直接逐字支持的单一事实。
- **修正建议：** 拆成两条直接事实，再单列“综合推断”；不要用一个“事实（高）”覆盖跨文献推演。

### M-05 购买门禁从一个厂商上线指南和自愿框架外推到所有企业

- **位置：** `market-research-draft.md:112-119,223-247`。
- **问题：** Microsoft 文档能证明其产品的权限、数据位置和容量上线要求，NIST AI RMF 是自愿风险管理框架；它们不足以单独证明所有目标买方“至少”采用同一八项顺序和角色链。
- **修正建议：** 将第 8 节标为“买方访谈假设清单”，不要标成已确认的通用采购流程；通过 FOBrain 客户采购／安全评审访谈确定顺序、否决项和责任角色。

### M-06 “一个任务”与同时选择 3—5 个任务的研究口径需要统一

- **位置：** `market-research-draft.md:31-35,55-59,251-258,399-408`；`domain-research-draft.md:40-46,50-55`。
- **问题：** 总策略要求只验证一个任务，但采用旅程和数据门禁又要求选 3—5 个任务。可以先发现多个候选再收敛一个，但当前没有明确阶段差异，容易导致首轮范围重新膨胀。
- **修正建议：** 明确“发现阶段观察多个近期样本／候选任务，实验阶段只选择一个任务；3—5 指同一任务的案例或场景，不是 3—5 个产品任务”。

## LOW

### L-01 证据等级口径在市场稿与领域稿之间不一致

- **位置：** `market-research-draft.md:37-45,49-53,103-118`；`domain-research-draft.md:57-68`。
- **问题：** 领域稿把厂商声明定义为“中”，市场稿却多次把包含厂商声明的能力结论标为“高”。同一证据可能因文档不同得到不同等级。
- **修正建议：** 全蓝图共用一张证据等级表，并把“来源真实性”“主张真实性”“目标客户适用性”分开评级。

### L-02 产品名与关系词缺少规范状态词典

- **位置：** `product-validation-evidence-plan.md:23-32`、`positioning-red-team.md:8-14,87-109`、`market-research-draft.md:25-35,304-335`。
- **问题：** “AI Workbench”“增值入口”“内嵌助手”“伴随式 AI”“独立 Workbench”“独立平台”交替出现。虽然多数被标为假设，但 VH-01 的“首版是”与红队的“暂停把入口当定位”在机器提炼时容易被读成冲突事实。
- **修正建议：** 建立关系枚举与状态：`DISCOVERY_BOUNDARY`、`SOLUTION_OPTION`、`OWNER_DECISION`、`VALIDATED_POSITIONING`；在决策前统一称“FOBrain 单任务 AI 增量机会”。

## 已通过的专项检查

1. **未把厂商自报效果直接当独立效果证据。** 市场稿和领域稿多处明确厂商资料只证明公开声明；需要修正的是从“声明存在”继续推到“商品化／市场基线／直接重复”的二次推断。
2. **数字阈值已明确标为建议。** `market-research-draft.md:366-429`、`product-validation-evidence-plan.md:9-12`、`positioning-red-team.md:73-85` 均说明样本量、比例、30%／95% 等是项目建议的预先决策规则，不是行业标准或 Owner 已确认值。
3. **法规没有被直接宣称全部必然适用。** `domain-research-draft.md:198-203` 和 `market-research-draft.md:337-349` 有适用性免责声明；问题集中在后续“直接转化／必须”的措辞没有继续保留条件标签。
4. **正式 BMAD 状态没有伪造。** 市场稿 `formal_steps_completed: []`、`scope_confirmation: pending`；领域稿 `formal_bmad_completion: false`；本地技术证据 `technical_research_scope_confirmation: pending`；任务矩阵、验证计划和红队分别保持 `draft`／`provisional`。这些状态与当前未获得逐步 `[C]` 确认的事实一致。
5. **本地部署限制没有被当成所有 FOBrain 客户的永久限制。** 本地可行性文档使用“当前私有部署”“当前证据”等限定语，且将身份、工单和闭环能力标为待补证。

## 关键来源抽查结果

- [FOBrain 官方产品页](https://www.huashunxinan.net/product-fobrain) 直接支持厂商公开声称的资产融合、漏洞归一、VPT 动态排序和闭环流程；不直接证明客户部署状态或实际效果。
- [ISC2 2025 Cybersecurity Workforce Study](https://www.isc2.org/Insights/2025/12/2025-ISC2-Cybersecurity-Workforce-Study) 直接支持 28% 已集成、19% 测试、22% 早期评估，以及使用者中 63% 自报显著生产力提升；不支持本产品的客观收益或购买意愿。
- [CISA BOD 23-01](https://www.cisa.gov/news-events/directives/bod-23-01-improving-asset-visibility-and-vulnerability-detection-federal-networks) 明确是对美国联邦行政部门民事机构（FCEB）的强制指令；其他组织只能把其内容作为参考实践，不能当作必然义务。
- [NIST SP 800-40 Rev.4](https://csrc.nist.gov/pubs/sp/800/40/r4/final) 直接支持企业补丁管理的识别、排序、获取、安装和验证；其直接范围是 patch management，不足以单独规定所有漏洞／暴露案件的关闭原因。
- [《生成式人工智能服务管理暂行办法》适用范围说明](https://www.cac.gov.cn/2023-07/13/c_1690898326863363.htm) 直接支持“未向境内公众提供生成式 AI 服务的企业内部研发／应用不适用该办法”；其他网络、数据和个人信息规则仍需分别判断。
- [《促进和规范数据跨境流动规定》](https://www.cac.gov.cn/2024-03/22/c_1712776611775634.htm) 明确存在豁免、数量阈值和不同合规路径，支持 H-07 的修正要求。
- [《国家网络安全事件报告管理办法》](https://www.cac.gov.cn/2025-09/15/c_1759583017717009.htm) 直接支持境内网络运营者对较大以上事件的分级报告规则；是否触发取决于主体和事件级别，AI 不应自行作最终法律判断。

## 修正优先级

1. 先修 H-03／H-04：阻止行业痛点直接进入需求，并统一发现到定位的门禁顺序。
2. 再修 H-01／H-02／H-05：统一证据等级，避免厂商声明和历史接口被当成已验证竞争／技术事实。
3. 再修 H-06／H-07／H-08：把条件性法规和补丁管理边界正确投影为候选约束。
4. 修完后重新做一次跨文档术语与状态扫描；在此之前，不应将 Market／Domain Research 标记为 final，也不应从这些草案自动生成 final Product Issues。

## 最终判定

**NEEDS WORK**

---

# V2 复审：更新后需求前置研究

## V2 范围与结论

本区保留上方 V1 历史结论不变，复核更新后的市场、领域、本地可行性、任务矩阵、验证计划、来源台账与定位红队，并新增审查：

- `research/evidence-status-model.md`
- `research/research-synthesis-draft.md`
- `assumptions-and-open-questions.md`
- `owner-decision-gates.md`

**V2 判定：`NEEDS_WORK`。** 第一轮 16 项 finding 中，14 项已实质消除，H-02 与 L-01 部分消除；未发现正式 BMAD Market／Domain／Technical Research 被标记为完成。更新后的材料已显著改善厂商声明、行业痛点、法规、阈值、关闭语义和技术可行性边界，但新证据状态与九步门禁仍有 5 项会改变下游结论的 HIGH 问题，不能据此关闭 OQ-B-00／OQ-B-08 或生成 final Brief／PRD／Product Issues。

## V2 对 V1 Finding 的逐项处置

| V1 ID | V2 disposition | 复核证据与剩余问题 |
| --- | --- | --- |
| H-01 | `RESOLVED` | `market-research-draft.md:53-59,275-310`、`source-evidence-ledger.md:42-45` 和 `research-synthesis-draft.md:46-58` 已把竞品能力限定为厂商公开主张，不再据此证明实际采用、效果或采购基线。 |
| H-02 | `PARTIALLY_RESOLVED` | 三／四层状态已由 `evidence-status-model.md:21-30,53-60` 固定，市场稿、台账和红队已区分部署／集成／任务效果；但 `domain-research-draft.md:36,76-83` 仍把厂商声明称为“既有业务底座”并直接规定 AI 应消费这些事实，`task-hypothesis-matrix.md:141` 仍写“直接重叠”，`product-validation-evidence-plan.md:59` 仍用“不与 FOBrain 既有能力重复”而非三层核验。剩余问题见 V2-M-01。 |
| H-03 | `RESOLVED` | `market-research-draft.md:176-200` 已改成 `INDUSTRY_SUPPORTED / UNKNOWN` 双轴并明确不赋予 P0/P1；`research-synthesis-draft.md:62-72,81` 禁止行业资料直接生成 Product Issue。 |
| H-04 | `RESOLVED` | 原“先选产品关系／先找任务”冲突已由 `owner-decision-gates.md:7-15`、`source-evidence-ledger.md:91-99`、`positioning-red-team.md:96-103` 和 `research-synthesis-draft.md:194-206` 统一为任务证据先于形态与产品关系。新九步门禁自身的依赖错误另见 V2-H-01～H-03。 |
| H-05 | `RESOLVED` | `domain-research-draft.md:313-332` 已改为 `PLAUSIBLE / UNVERIFIED` 或 `UNKNOWN`，明确不是可行性通过；`research-synthesis-draft.md:106-122,138-145` 要求目标部署、API、权限和真实样本实证。 |
| H-06 | `RESOLVED` | `domain-research-draft.md:221-231` 已分成 `DEFAULT_POLICY_CANDIDATE` 与 `LEGAL_REVIEW_TRIGGER`，明确不能直接成为 final PRD 的“合规必须”。 |
| H-07 | `RESOLVED` | `domain-research-draft.md:217,221-231` 已区分内部数据流盘点与安全评估／标准合同／认证等不同法律路径，并保留主体、数量和豁免判断。 |
| H-08 | `RESOLVED` | `market-research-draft.md:163`、`task-hypothesis-matrix.md:132-145`、`research-synthesis-draft.md:84,98,122,186` 及 `assumptions-and-open-questions.md:75` 已按关闭原因区分修复、缓解、接受、不适用／误报、延期、证据不足和仍存在，且把 NIST 限定在补丁／修复验证。 |
| M-01 | `RESOLVED` | `market-research-draft.md:53` 已准确拆出 ISC2 的 28% 已集成、19% 测试、22% 评估和 63% 使用者自报，并明确不能外推客观效果或购买意愿。 |
| M-02 | `RESOLVED` | `market-research-draft.md:59,148` 已改成“推荐的第一发现队列”，明确不排除售前账户或其他数据源用户。 |
| M-03 | `RESOLVED` | `domain-research-draft.md:38,145-155` 已把显式批准限定到写入、高影响、跨边界、不可逆或策略要求动作，只读调用使用确定性鉴权与最小权限。 |
| M-04 | `RESOLVED` | `market-research-draft.md:105-110` 已拆开 SP 800-40 与 IR 8286D 的直接事实，并把跨文献结论标为综合推断。 |
| M-05 | `RESOLVED` | `market-research-draft.md:214-250` 已把角色和决策标准明确标为“买方访谈假设”，不再声称所有企业采用同一顺序。 |
| M-06 | `RESOLVED` | `market-research-draft.md:407-416` 与 `research-synthesis-draft.md:124-132,155` 已明确 3～5 指同一任务的案例／场景，实验任一时点只允许一个任务。 |
| L-01 | `PARTIALLY_RESOLVED` | `evidence-status-model.md` 已建立统一来源／成熟度规则，但台账仍使用 `CONFIRMED/SUPPORTED/ASSUMPTION/CONFLICT/MISSING`，综合稿又使用 `CONFIRMED/SUPPORTED/HYPOTHESIS/CONTRADICTED/UNKNOWN`，且大小写与连字符不统一；缺少字段级映射。剩余问题见 V2-M-02 与 V2-L-01。 |
| L-02 | `RESOLVED` | `evidence-status-model.md:42-51` 已建立 `DISCOVERY_BOUNDARY / SOLUTION_OPTION / OWNER_DECISION / VALIDATED_POSITIONING`；`owner-decision-gates.md` 和综合稿也按该关系分阶段，不再把“首版是 Workbench”当已成立定位。 |

## V2 剩余 Findings

### CRITICAL

无。

### HIGH

#### V2-H-01 `VALIDATED_POSITIONING` 的升级条件遗漏采用／买方证据，并与 DG-06 形成循环依赖

- **位置：** `evidence-status-model.md:42-59`；`owner-decision-gates.md:11-12`；`assumptions-and-open-questions.md:42,60,79,107-108`。
- **问题：** 状态模型对 `VALIDATED_POSITIONING` 的定义包含采用证据（第 49 行），但升级规则只要求 `TASK_EFFECTIVE + Owner 决策`（第 59 行）。任务效果不能证明持续采用、组织归属、预算或采购。中央登记表又要求 DG-06 决定“标准功能／增值模块／独立产品”后，才按 OQ-H-12 条件性收集商业证据；同时 OQ-B-07 又要求 DG-06 决策包包含条件商业证据。这形成“先决定产品性质才解锁证据，但决定产品性质又需要该证据”的循环。
- **修正建议：** 增加 `ADOPTION_EVIDENCE`，并按产品性质区分内部采用、标准功能采用、增购／续费和独立购买证据。将 DG-06 拆成 `DG-06A 候选产品性质` 与 `DG-06B 经采用／买方证据验证的产品关系`，或把 OQ-H-12 提前为 DG-06 输入；`VALIDATED_POSITIONING` 至少需要 `TASK_EFFECTIVE + ADOPTION_EVIDENCE + OWNER_DECISION`。

#### V2-H-02 DG-05 要求实验选出形态，但实验主结果与判定规则直到 DG-08 才确认

- **位置：** `owner-decision-gates.md:10-14`；`research-synthesis-draft.md:126-130,155,196-205`；`assumptions-and-open-questions.md:40-44,114-124`；`product-validation-evidence-plan.md:29-32`。
- **问题：** DG-05 要在同任务实验后选出“胜出”形态，但唯一核心结果、基线、目标值、衡量周期和停止阈值被安排在 DG-08，即形态和产品关系之后。中央登记表虽在推荐步骤 7 插入 OQ-M-05 预注册阈值，却没有进入 canonical owner gate，且 VH-04 仍说目标值稍后确认。若没有预先确定主指标、保护指标、权重和胜出规则，DG-05 可以事后选择有利指标，无法客观判定内嵌、独立、非 AI 或不做谁胜出。
- **修正建议：** 在 DG-04 后增加强制的“实验契约／判定口径”门禁：主结果、基线、最小有意义改善、正确性／安全保护指标、样本、阈值和停止规则均在实验前签认。DG-08 可保留为基于实验结果确定 final 产品成功目标，但不能首次定义决定 DG-05 的度量。

#### V2-H-03 Product Brief 的 final 门禁在三份权威文件中不一致

- **位置：** `owner-decision-gates.md:12-15`；`research-synthesis-draft.md:134-151,210`；`assumptions-and-open-questions.md:103-110`。
- **问题：** Owner gate 表显示 DG-09 解锁 PRD 级非功能约束，DG-06／DG-08 解锁 Product Brief；中央登记表也把 DG-09 放到 PRD 必要条件，不列入 Brief。综合稿 PB-00、PB-09 和第 210 行却要求 DG-01～DG-09 全部通过后 Product Brief 才能 final。相同项目状态会被一份文件判为 Brief 可 final，另一份判为不可 final。
- **修正建议：** 确定一个 canonical 策略并同步三处。若 Brief 只需高层适用假设，将 DG-09 拆为 `Brief-level applicability boundary` 与 `PRD-level legal/control matrix`；若 Brief 必须完整通过 DG-09，则更新 owner gate 的“解锁”列和中央 final 门禁表。

#### V2-H-04 `CONTRADICTED` 混入“未证实／推理无效”，继续把缺少证据当成反证

- **位置：** `research-synthesis-draft.md:74-85`；对照 `evidence-status-model.md:21-30`。
- **问题：** X-06“聊天本身是独立产品价值”和 X-08“风险案件应成为首版核心对象”目前是 `HYPOTHESIS/UNKNOWN`，并没有反向实验或用户证据证明其为假；X-01、X-03、X-04 更接近“现有证据不足以推出该结论”或“推理规则无效”。把它们统一列为 `CONTRADICTED` 会把 absence of evidence 变成 evidence of absence，可能提前淘汰后来能被验证的合法选项。
- **修正建议：** 拆成 `CONTRADICTED_BY_EVIDENCE`、`UNSUPPORTED`、`INVALID_INFERENCE`、`PROHIBITED_UNTIL_VALIDATED`。只有存在直接反例或对照实验否定时使用 `CONTRADICTED_BY_EVIDENCE`；聊天与风险案件继续保持可证伪假设，而不是已否定产品方向。

#### V2-H-05 领域研究的“下一次正式 BMAD 动作”仍要求提前确认主用户与核心结果

- **位置：** `domain-research-draft.md:360-373,429-435`；对照 `owner-decision-gates.md:7-15`、`research-synthesis-draft.md:194-210` 与 `research-scope-proposal.md:49-59`。
- **问题：** canonical DG-01 明确只确认发现边界，不预设主用户、任务、形态或商业关系；领域稿结尾却要求产品所有者先确认“候选主用户、唯一核心结果和地域／行业范围”再选择 `[C]`。这会在正式 BMAD Domain Research 入口重新引入已经修正的前置定位，并一次请求多个未到达的决策。
- **修正建议：** 下一次正式动作只请求 DG-01 的发现边界与 BMAD scope `[C]`；研究队列在 DG-03 确认，单任务在 DG-04 选择，实验判据在新增门禁预注册，final 核心结果在相应后置门禁确认。第 13 节门禁和第 435 行必须与 canonical gate 同步。

### MEDIUM

#### V2-M-01 H-02 的声明升级仍残留在领域稿、任务矩阵和 PRD 准入规则

- **位置：** `domain-research-draft.md:36,76-83,313-317,375-380`；`task-hypothesis-matrix.md:141`；`product-validation-evidence-plan.md:50-60`。
- **问题：** 新模型已规定 `VENDOR_CLAIMED` 不能证明目标部署能力，但领域稿仍把官网能力称为“既有业务底座”，表中直接要求 AI 消费“现有事实”，T-08 仍写“直接重叠”，PRD 准入仍要求“不与 FOBrain 既有能力重复”。这些词可能让下游在未取得 `DEPLOYMENT_AVAILABLE/INTEGRATION_AVAILABLE` 前删除真实用户问题。
- **修正建议：** 统一改为“厂商声明范围重叠”；只有目标部署和任务基线确认后才能写“现有能力／既有底座／直接重复”。PRD 准入应要求“记录三层状态并说明任务增量”，不要求在厂商声明层证明不重复。

#### V2-M-02 所谓“统一证据状态”实际仍是三套未映射分类

- **位置：** `evidence-status-model.md:19-40`；`source-evidence-ledger.md:7-13`；`research-synthesis-draft.md:28-87`。
- **问题：** 新模型把来源类型、证据成熟度和认知状态放在一个 `状态` 列中；台账继续使用另一套 `CONFIRMED/SUPPORTED/ASSUMPTION/CONFLICT/MISSING`；综合稿又增加 `CONTRADICTED`。没有说明一条主张如何同时表达“厂商来源 + 目标适用性未知 + 当前是假设”，机器追踪和人工汇总仍会得到不同结果。
- **修正建议：** 把字段正交化为至少：`source_type`、`claim_support`、`target_applicability`、`maturity`、`decision_status`，并给旧值一张明确映射表。台账与综合稿逐行使用相同枚举，不只链接模型文件。

#### V2-M-03 “高频且高成本”与“低频但严重失败成本也可进入”的选择规则不一致

- **位置：** `owner-decision-gates.md:10`；`research-synthesis-draft.md:159,198-205`；`domain-research-draft.md:40-54`；`assumptions-and-open-questions.md:40,54`；`positioning-red-team.md:79`。
- **问题：** 多数门禁要求任务同时高频、高成本，ST-01 却允许“每周重复**或**具有显著失败成本”；红队阈值又要求每周发生，同时满足耗时或严重失败成本。当前规则无法判断低频但高损失、监管时限严格的任务是否合法进入实验。
- **修正建议：** 预先选择明确的准入逻辑，例如“达到最低频率且有可测成本”，或“高频效率型／低频高损失型二选一”，两类使用不同样本和成功指标。所有 gate、停止条件和访谈阈值使用同一布尔逻辑。

#### V2-M-04 `GO TO DISCOVERY` 与 DG-01 尚未通过的流程状态相冲突

- **位置：** `research-synthesis-draft.md:14-24`；`owner-decision-gates.md:17-25`；`assumptions-and-open-questions.md:37`。
- **问题：** 综合稿把当前状态写为 `GO TO DISCOVERY`，但 owner gate 明确 DG-01 仍为 `PENDING`，未通过时不得进入正式任务发现。前者看起来是已经获准执行，后者是等待授权。
- **修正建议：** 在 DG-01 前使用 `READY_TO_REQUEST_DG-01` 或 `DISCOVERY_RECOMMENDED — NOT AUTHORIZED`；只有 Owner 明确回答且正式 BMAD scope 选择 `[C]` 后才切换为 `DISCOVERY_AUTHORIZED`。

#### V2-M-05 `AUTHORITY_DIRECT` 没有记录规范效力，仍可能把指南写成强制要求

- **位置：** `evidence-status-model.md:21-30`；`research-synthesis-draft.md:42-49`；`domain-research-draft.md:233-247`。
- **问题：** 状态模型把法律、标准、政府指南合在 `AUTHORITY_DIRECT`，只处理“来源直接说了什么”，没有区分法律义务、有限主体的强制指令、自愿标准、框架和建议。综合稿 S-01 使用“流程要求”，领域稿表格使用“产品必须”，即使 NIST SP 800-40 和 AI RMF 主要是指南／自愿框架。
- **修正建议：** 增加 `normative_force = binding_law / scoped_directive / contractual / voluntary_standard / guidance / recommendation` 及适用主体；只有前三类在适用性确认后可直接形成合规 Must，其余先作为治理基线候选。

#### V2-M-06 真实用户／数据访问的最小合规审查被放到 DG-09，晚于发现和实验

- **位置：** `owner-decision-gates.md:8,15`；`research-synthesis-draft.md:141-149,198-206`；`assumptions-and-open-questions.md:38,45,58`。
- **问题：** DG-02 会取得目标版本、真实用户、近期任务和脱敏数据，DG-05 会运行实验；完整适用性与法务矩阵却锁到 DG-09。综合稿 DG-02 虽列出安全／法务责任人，但没有要求在访问前完成研究用途、授权、脱敏、保留和禁止外传的最小签认。
- **修正建议：** 在 DG-02 增加“研究数据使用批准”：研究目的、数据最小集、合法授权、脱敏、访问者、存储地域、保留／删除和不可进入模型项先通过；DG-09 保留为最终产品部署与行业法规适用性，不应承担首次研究数据审查。

#### V2-M-07 中央登记表对已更新材料的“当前事实”已经过期

- **位置：** `assumptions-and-open-questions.md:53,61,103-115`。
- **问题：** OQ-B-08 仍称“综合稿自身仍保留先选产品关系的冲突顺序”，但当前综合稿已经改为证据优先顺序；final 汇总仍只引用 V1 `NEEDS_WORK`，没有 v2 disposition。作为“唯一中央待决事项索引”，它不能同时保留已失效的当前事实而不标历史版本。
- **修正建议：** v2 审计后更新 OQ-B-00／OQ-B-08：逐项引用本区 disposition，删除已失效事实，只保留 V2-H／M 剩余项；记录审计版本和日期，避免“历史 finding”被当成“当前 finding”。

### LOW

#### V2-L-01 `USER_CONFIRMED` 实际表示 Owner 决策，名称与真实用户证据冲突

- **位置：** `evidence-status-model.md:21-29`；对照 `assumptions-and-open-questions.md:18-20`。
- **问题：** `USER_CONFIRMED` 被定义为产品所有者决定，而全蓝图另用 `USER-EVIDENCE` 表示真实一线用户证据。“用户”同时指 Owner 和最终用户，容易在自动追踪中误把 Owner 偏好当用户研究。
- **修正建议：** 改为 `OWNER_CONFIRMED` 或 `OWNER_DECISION_RECORDED`；`USER_EVIDENCE` 只保留给目标用户行为／任务证据。

#### V2-L-02 相同证据状态的大小写、连字符与斜杠写法仍不统一

- **位置：** `evidence-status-model.md:23-30`；`market-research-draft.md:31,39,51-55`；`research-synthesis-draft.md:46-58,108-113`。
- **问题：** 文档混用 `VENDOR_CLAIMED`、`vendor-claimed`、`厂商声明`、`deployment/integration availability`、`DEPLOYMENT_AVAILABLE` 和 `INTEGRATION_AVAILABLE`。对人可读，但不适合作为可机读追踪状态。
- **修正建议：** 正文可以使用中文解释，状态字段必须只使用 canonical 枚举；“deployment/integration”不得合并成一个状态，因为目标版本存在不等于接口可集成。

## V2 已通过检查

1. **正式 BMAD 状态未伪造。** `market-research-draft.md` 仍为 `scope_confirmation: pending`、`formal_steps_completed: []`；`domain-research-draft.md` 仍为 `formal_bmad_completion: false`；`local-fobrain-feasibility-evidence.md` 仍为 `technical_research_scope_confirmation: pending`；综合稿明确两类 Research 未完成。
2. **厂商效果未升级为实际效果。** 主要市场、台账、红队和综合结论已明确厂商页面只证明公开主张；剩余 V2-M-01 是局部“既有底座／直接重复”措辞和准入规则问题。
3. **行业痛点不再直接生成 Must。** 市场稿双轴与综合稿 Product Issue 准入已经阻止行业背景自动变成 FOBrain 剩余问题。
4. **建议阈值仍被正确标注。** 8～12 人、30%、95% 等继续标为待 Owner 预注册的工作建议，不是行业标准或已通过门槛。
5. **法规与跨境边界已正确分层。** 默认安全策略、条件性法律触发、法务判断及跨境不同路径已经分开；没有把公众生成式 AI 规则自动套用到内部场景。
6. **补丁证据范围与关闭分类已修正。** NIST 补丁验证不再被用于证明所有关闭类型，T-08 输出已覆盖主要合法结局。
7. **本地限制仍限定为当前部署／历史证据。** 未外推为所有 FOBrain 客户的永久事实。

## V2 修正顺序

1. 先修 V2-H-01／H-02：消除采用证据循环，并在 DG-05 前固定实验判定契约。
2. 再修 V2-H-03／H-05：统一 Brief final 门禁和正式 BMAD 下一步。
3. 修 V2-H-04 与 V2-M-02／M-05：正交化证据状态，纠正 `CONTRADICTED` 与规范效力。
4. 清理 V2-M-01／M-03／M-04／M-06／M-07 及两个 LOW 后，再执行 v3 独立复审。

## V2 最终判定

**NEEDS_WORK**

这不是正式 BMAD Research 完成信号；Market、Domain、Technical Research 继续保持草案／未完成状态。

---

# V3 复审：Canonical gate 与证据治理一致性

## V3 范围与结论

本轮完整复核 V2 的 14 项剩余 finding，并横向检查当前
`_bmad-output/planning-artifacts/product-blueprint` 与
`_bmad-output/forge/agent-workbench-product` 下的相关 Markdown。重点验证九步 gate、A／B
双准入路径、DG-05 实验前预注册、DG-06 的“不做”、DG-07 采用／买方证据、DG-02
研究数据批准、六字段证据记录、`maturity` 聚合、`normative_force`、历史库存映射和正式
BMAD 状态。

**V3 判定：`NEEDS_WORK`。** `owner-decision-gates.md`、中央开放问题表、研究综合稿和发现执行
工具包已经形成一致的九步定义；V2 的 8 项已完全消除，6 项只部分消除。当前仍有 5 项
HIGH、10 项 MEDIUM、1 项 LOW：Forge 旧定位材料仍要求先选 AI／高频效率／产品关系，市场研究的 M2
仍在 DG-05 之前执行形态对照，两个中央条目存在自锁，且若干“共同门禁”继续排除低频高损失路径。故不能关闭
OQ-B-00／OQ-B-08，也不能把 Product Brief、PRD 或 Product Issues 标记为 final。

## V3 对 V2 Finding 的逐项处置

| V2 ID | V3 disposition | 复核证据与剩余问题 |
| --- | --- | --- |
| V2-H-01 | `RESOLVED` | `evidence-status-model.md:46-51,79,89` 已要求独立聚合 `TASK_EFFECTIVE + ADOPTION_EVIDENCED + OWNER_DECISION_RECORDED`；`owner-decision-gates.md:12-13` 与 `assumptions-and-open-questions.md:43,60,79,120-123` 已把采用／买方取证提前到 DG-04 后并作为 DG-07 输入，不再循环解锁。 |
| V2-H-02 | `PARTIALLY_RESOLVED` | Canonical DG-05 已移到实验前（`owner-decision-gates.md:11-12`；`research-synthesis-draft.md:153-155,209-212`），但市场稿仍把形态对照放在 M2、数据检查和 DG-05 之前；见 V3-H-02。 |
| V2-H-03 | `RESOLVED` | `owner-decision-gates.md:15,19-21`、`research-synthesis-draft.md:142-159,218`、`assumptions-and-open-questions.md:103-110,126` 均明确 Product Brief final 必须通过 DG-01～DG-09。 |
| V2-H-04 | `PARTIALLY_RESOLVED` | 综合稿已删除 X-06／X-08 的错误反证并把聊天、风险案件保留为假设（`research-synthesis-draft.md:78-91`）；但 X-05 仍把证据层级推理无效写成 `CONTRADICTED`，见 V3-M-06。 |
| V2-H-05 | `RESOLVED` | `domain-research-draft.md:358-370,431-433` 下一次正式动作只请求 DG-01 与 `[C]`，主用户、任务、结果、关系均留在后续 gate。 |
| V2-M-01 | `PARTIALLY_RESOLVED` | 领域稿、T-08 和 PRD 准入已改为厂商声明范围；但市场稿仍称官网可“确定既有产品能力基线”，T-04 与 Forge 仍有未限定的“直接重叠／既有产品范围”推演，见 V3-M-01。 |
| V2-M-02 | `PARTIALLY_RESOLVED` | `evidence-status-model.md:19-70` 已正交化六字段并给出旧库存映射，综合稿也使用完整记录；但 Forge 关键主张没有六字段，A-008 又违反默认映射而未声明例外，见 V3-M-02／M-03。 |
| V2-M-03 | `PARTIALLY_RESOLVED` | 主控 gate、综合稿、中央登记表和执行工具包均已建立 A／B 双路径；但领域执行摘要、任务共同门禁、综合稿写操作条件和 Forge 命题仍只接受高频／效率改善，见 V3-H-01／H-03。 |
| V2-M-04 | `RESOLVED` | `research-synthesis-draft.md:8,14-24` 使用 `READY_TO_REQUEST_DG-01`，明确不是 `DISCOVERY_AUTHORIZED`；`owner-decision-gates.md:17-25` 仍保持 DG-01 `PENDING`。 |
| V2-M-05 | `PARTIALLY_RESOLVED` | 六字段已加入 `normative_force`，市场和领域稿已区分法律、指南及自愿标准；但综合稿把多来源主张压成一个 `GUIDANCE`，见 V3-M-05。 |
| V2-M-06 | `RESOLVED` | `owner-decision-gates.md:8`、`discovery-execution-kit.md:78-147` 已要求接触真实用户／数据前完成目的、授权、同意、脱敏、访问、地域、模型禁入、保留删除和事件处理批准。 |
| V2-M-07 | `RESOLVED` | `assumptions-and-open-questions.md:53,61,103-126` 已准确记录 v2 当前 finding、v3 待审计状态和最新九步顺序，不再把 V1 历史结论写成当前事实。 |
| V2-L-01 | `RESOLVED` | Canonical 名称已改为 `OWNER_DECISION_RECORDED`，`USER_EVIDENCE` 仅用于目标用户观察（`evidence-status-model.md:38-40`；`discovery-execution-kit.md:34`）。 |
| V2-L-02 | `RESOLVED` | 状态字段已统一为大写下划线枚举；“声明范围直接重叠”等中文仅作为判断正文，不再冒充成熟度值。 |

## V3 剩余 Findings

### CRITICAL

无。

### HIGH

#### V3-H-01 Forge 当前文件仍保留被 canonical 流程淘汰的定位与决策顺序

- **位置：** `../../../forge/agent-workbench-product/.memlog.md:11-22`；`../../../forge/agent-workbench-product/positioning-evidence.md:18-26,49-62,74-82`；`../../../forge/agent-workbench-product/user-problem-hypotheses.md:82-95`。
- **问题：** Forge memlog 仍把“AI 原生风险运营工作台／FOBrain AI 入口”记录为方向，并要求先确认产品关系、主用户和核心结果；定位证据又把“FOBrain AI 协作入口”称为当前最符合证据的首版候选，要求任务同时满足频率与失败成本、适合 AI 且证明效率改善，最后要求首先选择是否为 FOBrain AI 增强。用户问题稿也一次要求确认用户、结果、动作、产品关系、商业性质和三个任务。这与 DG-01 只确认发现边界、DG-04 才选 A 或 B、DG-06 才选形态、DG-07 才选产品关系直接冲突。`PROVISIONAL` 标签不能消除下游提炼器读取这些方向性陈述的风险。
- **修正建议：** 将旧方向显式标为 `SUPERSEDED_HISTORY`，把当前唯一命题替换为 DG-01 原文；删除“最符合证据的首版候选”、同时高频且高损失、预选 AI／效率和先选产品关系的要求。Forge 重新进入只能发生在 DG-06／DG-07 证据形成后。

#### V3-H-02 市场研究 M2 仍在数据验证和 DG-05 预注册之前完成解法形态对照

- **位置：** `market-research-draft.md:376-429`，尤其 `398-407` 对照 `409-429`；`research-synthesis-draft.md:153-155,209-212`。
- **问题：** 市场稿按 M0→M1→M2→M3→M4 排列；M2 明确要求“首先完成”原生、内嵌、独立和非 AI 的同任务对比，并以 AI 已改善时间／完整性／复核成本作为通过条件。数据可用性到 M3 才检查，DG-05 到 M4 才预注册。按该文本执行，团队会在没有验证数据和事前胜出规则时先宣布形态表现，重现 V2-H-02 的事后选指标风险。
- **修正建议：** 把 M2 限定为“现状替代与目标部署基线审计”，不得运行或判定候选解法；顺序固定为 M0/M1/M3 → DG-04 → DG-05 → 单任务对照 → DG-06。所有形态和“不做”只在 DG-05 后使用同一规则比较。

#### V3-H-03 多份共同准入规则仍把 B 低频高损失型强制改成高频效率型

- **位置：** `domain-research-draft.md:40-46`；`task-hypothesis-matrix.md:157-169`；`research-synthesis-draft.md:193`；`market-research-draft.md:327-333,420-429`；`../../../forge/agent-workbench-product/positioning-evidence.md:60-63`。
- **问题：** 领域执行摘要只允许“高频”调查；任务共同门禁要求所有任务的低保真实验同时证明效率改善；综合稿只有“真实高频价值”才允许写操作；市场主实验和 M4 只给更快／30% 时间改善；Forge 命题同时要求频率、失败成本和效率改善。这些是下游准入或范围规则，不是无害示例，会使已在 DG-04 合法选中的 B 路径无法进入实验、动作边界或需求。
- **修正建议：** 所有共同规则改为引用 DG-04 的唯一已选路径和 DG-05 契约。A 使用频率、净效率和质量保护；B 使用实际事件／受控演练、严重后果、客观任务结果与严重错误保护。只有明确标为 A 路径示例的地方才保留 30% 或高频措辞。

#### V3-H-04 两个中央证据条目被锁在自己必须解锁的 gate 上

- **位置：** `../assumptions-and-open-questions.md:58-59`，对照同文件 `25-29` 与 `owner-decision-gates.md:10,15`。
- **问题：** OQ-B-06 明确“阻止 DG-04”，当前状态却是 `LOCKED_BY_DG-04`；OQ-B-05 明确“阻止 DG-09”，状态却是 `LOCKED_BY_DG-09`。登记表自己规定 `LOCKED_BY_*` 表示受上游 gate 锁定，因此这两个条目永远不能在其 gate 前取得证据，而 gate 又永远不能在条目解决前通过。
- **修正建议：** OQ-B-06 应在 DG-03 观察后、DG-04 决策前处于可执行状态（例如 `LOCKED_BY_DG-03`／`PENDING_USER_EVIDENCE`）；OQ-B-05 应作为 DG-09 的输入在 DG-08 后解锁（例如 `LOCKED_BY_DG-08`／`PENDING_LEGAL_COMPLIANCE_REVIEW`）。

#### V3-H-05 执行计划允许 final Forge／Product Brief 绕过 DG-08 与 DG-09

- **位置：** `../workflow-plan.md:40-58,75-95,96-112`；对照 `../owner-decision-gates.md:13-21` 与 `../workflow-design.md:47-50`。
- **问题：** Task 2 的 re-entry gate 只等待到 DG-07 就允许输出 `forged-idea.md`；Task 4 只显式等待 DG-01～DG-03 和正式 Research，Task 5 的进入条件只写 PRFAQ verdict，没有写 DG-08／DG-09。主控文件要求 final Forge 与 Product Brief 都通过 DG-01～DG-09，执行计划却可按勾选顺序提前进入下游。
- **修正建议：** 将 Forge 的“方向硬化”分为 DG-07 后的产品关系草案与 DG-09 后的 final Forge；在 Task 5 入口和 Brief final 条件显式引用 DG-01～DG-09，不依赖读者跨文件推断。

### MEDIUM

#### V3-M-01 `VENDOR_CLAIMED` 的升级残留仍可把官网当目标能力基线

- **位置：** `market-research-draft.md:76-81`；`task-hypothesis-matrix.md:72-85`；`../../../forge/agent-workbench-product/positioning-evidence.md:18-26,49-55`。
- **问题：** 市场稿仍说 FOBrain 官方资料用于“确定既有产品能力基线”；T-04 使用未限定的“直接重叠”；Forge 把官网声明直接用于评定新工作台与“既有产品范围”高度重叠。官网只能确定厂商声明范围，不能确定目标版本、配置、接口、实际使用或任务结果。
- **修正建议：** 统一写成“确定 `VENDOR_CLAIMED` 范围／能力名称重叠”；只有 `DEPLOYMENT_AVAILABLE` 与任务基线证据形成后才使用“既有能力基线”“直接重复”。

#### V3-M-02 六字段规则尚未覆盖 Forge 的关键定位主张

- **位置：** `evidence-status-model.md:2-9,19-30`；`../../../forge/agent-workbench-product/positioning-evidence.md:18-45`；`../../../forge/agent-workbench-product/positioning-red-team.md:16-44`。
- **问题：** 证据模型明确适用于 Forge，并要求每条关键主张同时填写六个字段；两份 Forge 主文档仍只用“事实／推断／置信度／pending owner/user evidence”。这无法表达“厂商来源 + 目标适用性未知 + SOURCE_ONLY + 当前假设”，也使 Forge 可能绕过综合稿的 canonical records。
- **修正建议：** 为会影响定位、用户、解法和产品关系的 Forge 主张添加六字段记录或稳定引用综合稿中的对应 record ID；自由文本“事实／推断”只能作为解释，不能替代记录。

#### V3-M-03 历史库存映射存在未声明例外 A-008

- **位置：** `evidence-status-model.md:57-70`；`../source-evidence-ledger.md:47-58`。
- **问题：** 默认映射把全部 `A-* / ASSUMPTION` 解释为 `source_type=PROJECT_HISTORY`；A-008 的来源却明确是历史材料与 E-001 厂商声明的对照推断，实际应为 `MULTIPLE_SOURCES / SYNTHESIS`。规则要求例外显式填写完整字段，但该行没有，机器派生会产生错误来源类型。
- **修正建议：** 在 A-008 行显式记录六字段，至少改为 `source_type=MULTIPLE_SOURCES`、`claim_support=SYNTHESIS`、`target_applicability=UNKNOWN`、`maturity=SOURCE_ONLY`、`decision_status=HYPOTHESIS`、`normative_force=NOT_NORMATIVE`。

#### V3-M-04 A／B 路径的阈值预注册时间仍自相矛盾

- **位置：** `market-research-draft.md:376-396`；`discovery-execution-kit.md:24-28,375-397`；对照 `owner-decision-gates.md:10`。
- **问题：** 市场稿要求访谈／采集前就“选择一种路径”；执行工具包一方面要求所有阈值在采集前确认，另一方面要求 DG-04 在收集并聚类真实样本后才选择路径。若按前者执行，会在没有发现证据时预选 A/B；按后者执行，又可能看完数据再定有利阈值。
- **修正建议：** 在采集前同时预注册 A、B 两套候选准入定义和最低证据要求，但不选择路径；完成原始观察并冻结数据后，DG-04 按预先规则只选一条路径和一个任务。DG-05 再在实验前固定该路径的结果指标。

#### V3-M-05 `normative_force` 被错误聚合到多来源主张

- **位置：** `evidence-status-model.md:30,40,90`；`research-synthesis-draft.md:47-50`。
- **问题：** 模型说明 `normative_force` 评价单个法规、标准、合同或指南来源；S-05 却把 NIST 与 OWASP 合并为 `MULTIPLE_SOURCES / SYNTHESIS / ... / GUIDANCE`。不同来源可能分别是指南、标准或建议，一个聚合值会隐藏规范效力差异；S-02 已正确说明应逐项记录，S-05 与之不一致。
- **修正建议：** 将 S-05 的聚合 `normative_force` 设为 `UNKNOWN`／`NOT_NORMATIVE` 并附逐来源记录，或拆成每个 authority source 的完整 record，再把综合推断标为非规范性候选。

#### V3-M-06 `CONTRADICTED` 仍用于证据层级规则，而非直接反证

- **位置：** `evidence-status-model.md:26`；`research-synthesis-draft.md:82-90`。
- **问题：** X-05“官网声明等于目标部署可用且任务有效”没有直接反向现场证据；现有规则只能说明该推理无效、证据不足。标成 `CONTRADICTED / REJECTED` 仍把治理规则误写成实证反驳。X-02 的“独有”有厂商公开重叠反证，二者不应使用同一支持状态。
- **修正建议：** X-05 改为 `NONE / REJECTED`，原因码保留 `EVIDENCE_LAYER_COLLAPSE` 或 `INVALID_INFERENCE`；只有有直接反例的命题使用 `CONTRADICTED`。

#### V3-M-07 产品验证计划没有完整表达 DG-02 的证据治理责任

- **位置：** `product-validation-evidence-plan.md:15-23,29-31`；对照 `assumptions-and-open-questions.md:38` 与 `discovery-execution-kit.md:96-117`。
- **问题：** 验证计划声称列出证据责任类别，但只定义 Owner、用户、FOBrain API／数据和实验四类；VH-01B 的研究数据批准也没有 `FOBRAIN_PRODUCT_EVIDENCE`、`LEGAL_COMPLIANCE_REVIEW` 或等价的数据安全责任。执行工具包虽补齐了批准表，若单独读取验证计划仍可能把“已签认”理解为业务／API 人员即可关闭。
- **修正建议：** 复用中央登记表的完整 `responsibility_class` 枚举，并在 VH-01B 明确产品／部署、数据安全和法务／合规（按适用性）的签认职责；链接不能替代该行的责任边界。

#### V3-M-08 发现采集模板仍把 shorthand 与正交字段混作单一状态

- **位置：** `discovery-execution-kit.md:268-283,343-365`；对照 `evidence-status-model.md:19-46`。
- **问题：** 样本元数据把 `HYPOTHESIS`（`decision_status`）和 `UNKNOWN`（通常是 `target_applicability`）列在同一个“关联主张状态”值中；任务簇又在“能力成熟度”字段中列入 `VENDOR_CLAIMED`，但该词是四字段 shorthand，不是 `maturity` 枚举。实际采集表会再次压缩已正交化的状态。
- **修正建议：** 两处都直接给出六个字段；`maturity` 只能填写 `SOURCE_ONLY / DEPLOYMENT_AVAILABLE / INTEGRATION_AVAILABLE / TASK_EFFECTIVE / ADOPTION_EVIDENCED`，`VENDOR_CLAIMED` 只能作为完整字段组合的展示别名。

#### V3-M-09 采用／买方证据没有独立来源类型，也未进入核心验证假设

- **位置：** `evidence-status-model.md:23-30,46-51`；`product-validation-evidence-plan.md:25-37`；`../assumptions-and-open-questions.md:43,60,79`。
- **问题：** 模型新增 `ADOPTION_EVIDENCED` 成熟度，却没有能明确表示经济买方访谈、预算／采购证据、采用遥测或内部责任承诺的来源类型；中央表把这些全部归到 `USER_EVIDENCE`，而目标用户观察不能等同于买方／采购证据。验证计划的 VH 列表也没有一条 adoption/buyer 假设与通过条件，DG-07 的必要证据可能无法被一致采集和升级。
- **修正建议：** 增加明确来源类型（如 `TARGET_BUYER_OBSERVATION`、`ADOPTION_TELEMETRY`／`COMMERCIAL_RECORD`，或可证明等价的正交记录），并在验证计划新增 DG-07 采用／买方假设；按内部工具、标准功能、增值、独立产品分别定义证据，不能用一线用户偏好替代。

#### V3-M-10 局部历史证据仍使用“当前不可用／可行性”现时语态

- **位置：** `local-fobrain-feasibility-evidence.md:8-10,21-30,42-50`；`task-hypothesis-matrix.md:42-55,72-85,87-100`。
- **问题：** 文档 frontmatter 已说明来源是历史参考，但正文仍写“当前私有部署路径不可用”“查询与详情类有局部可行性证据”“统计事实可用”，任务矩阵也写“已确认的安全事实字段”“可用事实片段”。这些都来自已降级的 `docs/`，没有当前目标版本现场证据；现时语态容易被解析为 `DEPLOYMENT_AVAILABLE` 或 `INTEGRATION_AVAILABLE`。
- **修正建议：** 统一改为“归档时历史材料曾记录／曾提供线索”，并在每条结论显式保持 `PROJECT_HISTORY / ... / SOURCE_ONLY`；只有目标环境复核后才使用“当前不可用／可用／可行”。

### LOW

#### V3-L-01 市场稿的全局 `source_verification: true` 缺少可审计范围

- **位置：** `market-research-draft.md:1-23,486-519`。
- **问题：** Frontmatter 给出整体 `source_verification: true`，但来源索引没有逐项访问结果、验证日期、失败／重定向或“只验证来源存在、未验证主张适用性”的范围。该布尔值比正文实际可证明的状态更强。
- **修正建议：** 改为 `source_verification: partial`／`source_inventory_complete`，或增加逐来源核验记录并说明它只证明页面与引文匹配，不证明目标适用性、采用或效果。

## V3 已通过检查

1. **九步 canonical gate 定义正确。** 主控 gate、综合稿、流程设计与发现工具包均使用 DG-01 边界 → DG-02 数据批准 → DG-03 队列 → DG-04 单任务／单路径 → DG-05 预注册 → DG-06 形态／不做 → DG-07 产品性质／采用 → DG-08 动作 → DG-09 适用边界；中央条目的可执行状态仍需按 V3-H-04 修正。
2. **DG-06 明确保留“不做”。** 主控、综合、红队和发现合同均把原生增强、AI、独立、非 AI 与不做纳入同规则比较。
3. **采用／买方证据循环已消除。** 采用取证从 DG-04 后开始，在 DG-07 前完成；`VALIDATED_POSITIONING` 需要三类独立记录聚合。
4. **DG-02 已有可执行研究数据批准。** 目的、最小数据、授权／同意、脱敏、访问者、地域、模型／遥测禁入、保留删除、事件处理和批准人均有字段，未批准不得观察。
5. **法规规范效力主体已建立。** 法律表使用 `BINDING_LAW`，NIST 使用 `GUIDANCE`，ISO 使用 `VOLUNTARY_STANDARD`，且适用性继续待法务确认。
6. **当前状态未冒充正式 BMAD Research。** 市场稿仍为 `scope_confirmation: pending`、`formal_steps_completed: []`；领域稿为 `formal_bmad_completion: false`；综合稿为 `READY_TO_REQUEST_DG-01`；正式 `[C]` 尚未取得。
7. **Product Issue 契约保持正确。** Background、Goal、业务输入／输出、每个 Goal 至少一条客观 Given/When/Then，以及六字段证据记录均已进入模板和门禁。

## V3 最小修正清单

1. 重写 Forge 的 memlog、定位证据和用户问题稿：旧方向标记为已废弃，只保留 DG-01 当前命题与九步顺序。
2. 重排市场 M0～M4，并把所有共同实验／任务门禁改为按 A 或 B 的 DG-05 契约判定。
3. 修复 OQ-B-05／B-06 自锁，并在执行计划中把 final Forge／Brief 显式锁到 DG-09。
4. 为 Forge、采集模板和 A-008 补齐六字段；为采用／买方证据增加明确来源类型和验证假设。
5. 统一 A／B 阈值时序：采集前冻结两套准入规则，DG-04 冻结数据后只选一路，DG-05 再锁实验指标。
6. 清除“官网即能力基线／历史即当前可行性”，所有归档证据保持 `PROJECT_HISTORY / SOURCE_ONLY`。
7. 修正 S-05 的聚合 `normative_force`、X-05 的 `claim_support`、DG-02 责任类别和全局来源核验声明。

完成以上五组机械修正后再执行 v4；无需重新开展市场或领域资料搜索，但正式 BMAD Research 仍必须等待 DG-01、实际 `[C]`、DG-02 批准和真实任务证据。

## V3 最终判定

**NEEDS_WORK**

这不是正式 BMAD Research 完成信号；Market、Domain、Technical Research 继续保持草案／未完成状态。

---

# V4 复审：V3 修正后的稳定快照

## V4 范围与结论

本轮在 V3 全量范围上重新读取当前稳定快照，并补充检查脑暴快照是否仍会把候选解法直接注入下游。重点回归：Forge 重写、A／B 双规则预注册时点、M2 仅保留现状基线、DG-05 后的全形态／“不做”对照、OQ-B-05／B-06 解锁、final Forge／Brief 的 DG-09 门禁、采用／买方来源类型、历史证据语态以及多来源 `normative_force`。

**V4 判定：`PASS_WITH_CONDITIONS`。** V3 的 5 项 HIGH、10 项 MEDIUM、1 项 LOW 均已实质修正；重新分析后未发现新的 HIGH 或 MEDIUM。剩余 2 项均为 LOW：一项是 canonical 记录精度，一项是 V4 结束后必然要做的版本／状态回写。它们不会改变 DG-01～DG-09 顺序、任务准入、解法胜负、产品关系或 Product Brief final 结论。

本判定只表示当前草案可以在完成下列 LOW 清理后请求 DG-01 和正式 BMAD `[C]`。它**不是**正式 BMAD Market／Domain／Technical Research 完成信号，不授权访问真实用户／数据，不允许生成 final Forge、Product Brief、PRD 或 Product Issues。

## V4 对 V3 Finding 的逐项处置

| V3 ID | V4 disposition | 当前复核证据 |
| --- | --- | --- |
| V3-H-01 | `RESOLVED` | Forge memlog 已把旧方向标为 `superseded_history`，当前只保留一个重要剩余任务、A／B、全部形态与“不做”（`.memlog.md:7-14`）；定位证据和用户问题稿也只请求 DG-01，并把后续决定留在各自 gate（`positioning-evidence.md:13-29,42-62`；`user-problem-hypotheses.md:86-96`）。 |
| V3-H-02 | `RESOLVED` | 市场稿 M2 现在只审计原生流程、配置／培训、人工流程和已投入使用的替代，不制作或比较候选原型（`market-research-draft.md:401-410`）；全部形态只在 DG-05 后进入 M4（`market-research-draft.md:423-433`）。 |
| V3-H-03 | `RESOLVED` | 市场采用旅程、价值实验和风险停止规则均分别表达 A 效率结果与 B 高损失结果（`market-research-draft.md:262-271,330-348,364-376`）；领域摘要、任务共同门禁、综合稿和 Forge 红队均允许 B、非 AI 与“不做”，不再把高频／30% 时间改善设为共同门槛。 |
| V3-H-04 | `RESOLVED` | OQ-B-05 已在 DG-07／DG-08 后形成 DG-09 输入，不再由 DG-09 反向解锁；OQ-B-06 已在 DG-03 后解锁并作为 DG-04 输入（`../assumptions-and-open-questions.md:58-59`）。 |
| V3-H-05 | `RESOLVED` | 流程设计明确 DG-08／DG-09 后才 final Forge，Brief 必须通过 DG-01～DG-09（`../workflow-design.md:46-50`）；执行计划把 DG-07 后内容限定为定位草案，并把 Task 5 与进入 PRD 显式锁到 DG-09（`../workflow-plan.md:52-58,96-114`）。 |
| V3-M-01 | `RESOLVED` | 市场、任务矩阵和 Forge 都把官网限定为 `VENDOR_CLAIMED`／“厂商声明范围重叠”；目标版本、配置、集成、采用与任务效果继续为未知（`market-research-draft.md:32-40`；`task-hypothesis-matrix.md:15-17,36-37,51-52`；`positioning-evidence.md:23-29`）。 |
| V3-M-02 | `RESOLVED` | Forge 已为当前定位、任务、解法、产品关系和 Microsoft 厂商公开主张添加 PE-01～PE-07、PRT-01～PRT-05 六字段记录（`positioning-evidence.md:19-29`；`positioning-red-team.md:16-26`），并在正文稳定引用 `PRT-05 / VENDOR_CLAIMED`，不再让自由文本标签替代 canonical 记录。 |
| V3-M-03 | `RESOLVED` | A-008 已显式声明 `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE` 例外（`../source-evidence-ledger.md:58`）。 |
| V3-M-04 | `RESOLVED` | Owner gate、市场稿、发现工具包和综合稿均统一为：采集前同时冻结 A／B 两套规则但不选路径；样本冻结后 DG-04 只应用固定版本选一路和一个任务；DG-05 再锁实验结果规则（`../owner-decision-gates.md:10-11`；`market-research-draft.md:379-399`；`discovery-execution-kit.md:24-28,375-397`；`research-synthesis-draft.md:200-207`）。 |
| V3-M-05 | `RESOLVED` | 统一模型明确多来源综合主张的 `normative_force=NOT_NORMATIVE`（`evidence-status-model.md:90-92`）；综合稿 S-05 已按该规则修正，并把底层来源效力留在各自记录（`research-synthesis-draft.md:50`）。 |
| V3-M-06 | `RESOLVED` | X-05 已改为 `NONE / REJECTED`，明确只是无支持推理而非目标部署反证（`research-synthesis-draft.md:80-91`）。 |
| V3-M-07 | `RESOLVED` | 验证计划已列 Owner、用户、采用／买方、FOBrain 产品、API／数据、法务／合规和实验责任；VH-01B 要求访问前由对应责任人完整签认（`product-validation-evidence-plan.md:15-26,30-41`）。 |
| V3-M-08 | `RESOLVED` | 发现样本、附件和任务簇模板均要求逐字段填写 canonical record；`VENDOR_CLAIMED` 只作为 shorthand，`maturity` 仍填写合法枚举（`discovery-execution-kit.md:268-315,343-365`）。 |
| V3-M-09 | `RESOLVED` | 证据模型新增 `TARGET_STAKEHOLDER_OBSERVATION`、`ADOPTION_TELEMETRY`、`COMMERCIAL_RECORD` 与 `ADOPTION_BUYER_EVIDENCE` 规则，并要求独立聚合任务效果、采用与 Owner 决策（`evidence-status-model.md:23-52`）；VH-09 和 OQ-DG-07／OQ-H-12 已接入相同来源与责任类别。 |
| V3-M-10 | `RESOLVED` | 本地证据全文已限定为归档线索和 `PROJECT_HISTORY / SOURCE_ONLY`（`local-fobrain-feasibility-evidence.md:8-30,42-50`）；任务矩阵也把目标部署现状改为未知，并将混合来源任务记录为 `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE`（`task-hypothesis-matrix.md:11-23,37,67,82,112,127,142`）。 |
| V3-L-01 | `RESOLVED` | 市场稿 frontmatter 已改为 `source_verification_status: partial`，并明确只抽查关键官方来源、正式流程仍需逐项复核（`market-research-draft.md:20-23`）。 |

## V4 新 Findings

### CRITICAL

无。

### HIGH

无。

### MEDIUM

无。对 A／B、DG-05、形态／不做、采用证据、final 门禁和历史证据执行了一次反向路径重放，未发现会改变下游结论的残留。

### LOW

#### V4-L-01 V4 完成后的中央版本状态尚待机械回写

- **位置：** `../assumptions-and-open-questions.md:53,61,108-115`；`../workflow-plan.md:86-94`。
- **问题：** 中央表和执行计划仍保留“等待 v3／修正后的 v3”或“v4 尚未结束”等审计时点文字。这在 V4 落盘前是正确状态，但 V4 结束后会成为过期的行政状态。
- **条件：** 主流程在消费本审计后，把 OQ-B-00／OQ-B-08、final 门禁汇总、推荐顺序和 Task 4 审计步骤统一回写为 V4 `PASS_WITH_CONDITIONS`；不得借此把正式 BMAD Research 标为完成。

#### V4-L-02 三条历史范围记录的 `target_applicability=PARTIAL` 语义可更精确

- **位置：** `../../../forge/agent-workbench-product/positioning-evidence.md:23`；`research-synthesis-draft.md:49,52`。
- **问题：** 三条主张描述的是归档材料本身，却使用 `target_applicability=PARTIAL`；同一行又正确说明它们不证明目标产品／部署。由于 `source_type=PROJECT_HISTORY`、`maturity=SOURCE_ONLY`、`decision_status=HYPOTHESIS` 已阻止升级，这不会改变任何 gate，但机器读取时仍可能把 `PARTIAL` 理解为目标适用性已有局部实证。
- **条件：** 若主张对象只是历史库存，改为 `NOT_APPLICABLE`；若主张要投影到新产品，则改为 `UNKNOWN`。不得用 `PARTIAL` 表示“历史材料只覆盖一部分”。

## V4 已通过的专项检查

1. **Forge 已完成发现命题重写。** 旧 AI／Workbench 定位显式废弃，当前只请求 DG-01。
2. **A／B 时点无采样后挑规则。** 两套准入规则采集前冻结，路径在样本冻结后的 DG-04 选择。
3. **M2 不再提前实验。** M2 只做当前状态和替代方案基线。
4. **DG-05 先于任何方案胜负。** 核心结果、基线、目标、保护和停止线均先预注册。
5. **DG-06 保留完整结果集。** 原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 和“不做”均可胜出。
6. **B 路径没有被效率门槛覆盖。** B 使用事件／演练、严重后果、客观结果与严重错误保护；时限不是通用准入项。
7. **中央自锁已消除。** OQ-B-06 在 DG-03 后可执行，OQ-B-05 在 DG-07／DG-08 后形成 DG-09 输入。
8. **final Forge／Brief 不能绕过 DG-09。** 流程设计、计划、主控 gate 和综合稿一致。
9. **采用／买方证据可独立表达。** 用户观察、利益相关者观察、采用遥测、商业记录与 Owner 决策不再互相替代。
10. **Forge 竞品主张可追踪。** Microsoft 官方能力与风险提示已记录为 PRT-05／`VENDOR_CLAIMED`，不证明本项目效果或必须跟进。
11. **历史证据保持 `SOURCE_ONLY`。** 归档材料没有被写成目标部署当前可用、不可用或任务有效。
12. **多来源规范效力不再聚合。** 综合主张为 `NOT_NORMATIVE`，法律、指令、指南和标准分别保留效力。
13. **脑暴不会直接生成 Must。** 原 MoSCoW 已标为 `SUPERSEDED_BRAINSTORM_PRIORITY_SNAPSHOT`，仅经 DG-04～DG-09 的内容才可进入 final Brief。
14. **正式流程状态诚实。** Market `formal_steps_completed: []`、Domain `formal_bmad_completion: false`、Technical scope pending、综合稿 `READY_TO_REQUEST_DG-01` 均未伪造完成。
15. **Product Issue 契约保持完整。** Background、Goal、业务输入／输出、每个 Goal 至少一条客观 Given／When／Then 和六字段记录均仍是硬门禁。

## V4 最终判定

**PASS_WITH_CONDITIONS**

在完成 V4-L-01～L-02 的低风险机械清理后，可关闭本轮证据治理复审并请求 DG-01／正式 BMAD `[C]`。DG-01、DG-02、真实任务证据和后续 gate 仍全部未通过；正式 BMAD Research、final Forge、Product Brief、PRD 与 Product Issues 继续保持未完成。

---

# V5：V4 Conditions Closure

## 关闭复核

**V5 判定：`PASS`。** V4-L-01 与 V4-L-02 均已关闭，未发现新的残留：

- `V4-L-01 CLOSED`：OQ-B-00／OQ-B-08、final 汇总、推荐顺序、workflow plan 与 progress 已回写，并明确治理通过不等于正式 Research 或产品验证通过。
- `V4-L-02 CLOSED`：PE-01、S-04、S-07 均已改为 `target_applicability=NOT_APPLICABLE`，且正文明确这些记录只描述历史库存。
- 额外回归通过：发现工具包已把 A／B 双规则预注册列为任务观察前的独立步骤；采集前不选路径，样本冻结后 DG-04 才选一路和一个任务。

最小十点反向检查均无新 finding：两项条件原位、中央治理状态、final 门禁、推荐顺序、workflow plan、progress、发现执行顺序、历史字段全局残留、正式流程状态以及 audit 版本链均一致。

## V5 最终判定

**PASS**

本结论只关闭当前研究／Forge 草案的证据治理审计。正式 BMAD Market、Domain、Technical Research 仍未完成；DG-01、实际流程 `[C]`、DG-02 数据批准、真实任务证据和 DG-03～DG-09 均仍待执行。不得据此生成 final Forge、Product Brief、PRD 或 Product Issues。
