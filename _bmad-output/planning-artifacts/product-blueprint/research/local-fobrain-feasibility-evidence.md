---
status: active-mixed-evidence
evidence_type: current_live_integration_evidence_and_historical_reference
technical_research_scope_confirmation: current_live_evidence_reused
updated: '2026-07-12'
---

# FOBrain 当前 live 集成证据与历史可行性线索

本文件分别记录两层证据：

1. 当前仓库 `test-results/` 中已经存在的真实 FOBrain `provider_mode=live` 报告，可直接证明对应环境、时间和 capability 的调用／StructuredResult／脱敏结果；通过项达到 `INTEGRATION_AVAILABLE`。
2. 归档 `docs/` 中的历史材料，只能作为 `PROJECT_HISTORY / SOURCE_ONLY` 线索。

`INTEGRATION_AVAILABLE` 只证明工具真实可调用，不证明用户需要该工具、任务得到改善、产品形态正确或值得购买。

## 当前 live 证据（直接复用，不重复调研接口存在性）

| 批次 | 报告时间（UTC） | 结果 | capability 结论 | 证据成熟度与限制 |
| --- | --- | --- | --- | --- |
| Batch A | 2026-07-02 09:17 | `passed` | `current_user_context`、`my_permissions` 两个只读工具和 `connector.fobrain.security` 通过 | 对这 2 个工具与 connector：`EXPERIMENT_RESULT / DIRECT / PARTIAL / INTEGRATION_AVAILABLE`；不证明任务价值 |
| Batch B | 2026-07-06 08:17 | `blocked` | 6 个“我的／本部门范围”能力因空结果或 `missing_current_user_scope` 被阻塞 | 直接反证当前账号范围足以验收；阻止“24 个只读全部通过” |
| Batch C | 2026-07-06 07:32 | `passed` | `business_list`、`external_high_risk_assets`、`vulnerability_status_summary`、`ip_stats`、`vul_stats` 通过；`pending_tickets` 因功能不可用跳过 | 5 个工具达到 `INTEGRATION_AVAILABLE`；工单能力未证明 |
| Batch D | 2026-07-02 09:51 | `passed` | 按 owner／department／IP 查询资产和漏洞的 6 个工具通过 | 6 个工具达到 `INTEGRATION_AVAILABLE`；样本值被安全隔离，不证明最终用户范围 |
| Batch E | 2026-07-03 06:37 | `passed` | 资产详情、漏洞详情、业务风险汇总、威胁关联 4 个工具通过 | 4 个工具达到 `INTEGRATION_AVAILABLE`；不证明用户任务或闭环 |
| Sample discovery | 2026-07-03 06:29 | `passed` | live 样本发现通过且报告不暴露样本值 | 证明当前环境可产生安全样本，不证明其代表用户任务 |

总计：历史定义的 24 个只读工具中，**17 个 live passed、6 个 blocked、1 个 skipped**；connector 另行通过。所有报告的脱敏检查均包含凭据／Authorization／raw payload 禁出项，按各报告字段通过。

当前技术研究无需再从零确认上述 17 个工具是否能真实调用。后续只需：

- 对 blocked／skipped 能力按实际产品范围决定是否补证或删除；
- 把已通过工具组合进一个真实任务回放，验证用户结果；
- 采集现状流程、频率、成本、错误、输出和结束标准；
- 比较 FOBrain 原生、AI、非 AI 与不做，而不是重复执行接口盘点。

## 历史材料曾记录的线索

| ID | 历史陈述 | 本地来源 | 对发现工作的含义 |
| --- | --- | --- | --- |
| DF-01 | 归档矩阵曾列出当前用户与权限、按负责人／部门／IP 查询资产和漏洞、对象详情、业务风险汇总、外部高风险资产、威胁关联、状态／数量统计、待处理工单和“我的范围”等 24 个只读入口。 | `docs/fobrain-tool-matrix.md:58-85` | 可用于构造查询、详情和汇总类访谈探针；不能据此证明目标部署存在这些入口、24 项高频或值得进入首版。 |
| DF-02 | 归档资料曾只列出一个写域动作：更新工单状态，并陈述执行前审批、拒绝／取消稳定和批准后只执行一次。 | `docs/fobrain-tool-matrix.md:87-93,117-126` | 只能形成待核验的最小写域线索；不证明目标部署 API、真实样本、必要性或首版动作边界。 |
| DF-03 | 归档资料曾把资产、漏洞、人员／部门关系、业务、工单、统计与权限摘要列为安全事实候选，并禁止 raw payload、凭据和未投影数据进入用户输出。 | `docs/fobrain-source-api-reference.md:83-96`；`docs/fobrain-tool-matrix.md:11-39` | 可用于设计字段核验清单；具体字段存在性、完整性、时效、权限和业务含义仍需目标部署实测。 |
| DF-04 | 归档矩阵曾要求同名人员时澄清，不允许把候选结果当业务事实。 | `docs/fobrain-tool-matrix.md:64,87-92` | “不静默猜测对象”可作为安全研究候选；是否为真实高频问题和如何验收仍待用户与数据证据。 |

## 归档时历史材料曾记录的限制

| ID | 限制 | 本地来源 | 需求影响 |
| --- | --- | --- | --- |
| DL-01 | 归档时某一私有部署样本中的当前用户响应没有显式部门，`staff_ids` 为 `null`；当时的人员／部门发现路径未成立。 | `docs/fobrain-source-api-reference.md:55-69` | 只能形成目标部署身份映射核验项；不能断言当前环境仍不可用。 |
| DL-02 | 当前“我的资产／漏洞／业务”历史映射使用姓名 fallback，而不是稳定 staff id。 | `docs/fobrain-source-api-reference.md:61-69` | 重名、改名和跨部门人员可能造成范围错误；最终需求必须定义歧义、无映射和权限不足时的客观结果。 |
| DL-03 | 归档时的 live plan 曾记录当时没有待处理工单样本，并把该读取能力标为暂缓。 | `docs/fobrain-live-read-batch-plan.md:56-64` | 目标部署仍需重新取得代表性工单样本；历史缺样本不能证明当前有或没有。 |
| DL-04 | 历史“业务风险汇总”来源是聚合统计 bucket；工具矩阵同时期望“关键发现、风险判断”。 | `docs/fobrain-source-api-reference.md:76-81`；`docs/fobrain-tool-matrix.md:72` | 统计事实与 AI 风险判断必须分层；后者需给出证据、规则、未知和可复核验收，不能伪装成上游事实。 |
| DL-05 | 24 个历史能力以工具和实现恢复批次组织，并非以用户任务链、频率或结果组织。 | `docs/fobrain-tool-matrix.md:1-9,58-85`；`docs/fobrain-live-read-batch-plan.md` | 不能直接把“一工具一需求”投影到 Product Issues；必须先按真实任务重组和删减。 |
| DL-06 | 风险案件所需的责任确认、截止时间、评论／异议、风险例外、升级、修复证据和结果复核，不在当前 24 个只读入口或唯一写动作中得到证明。 | `docs/fobrain-tool-matrix.md:58-93` 与脑暴候选能力对照 | “完整风险案件”当前不可承诺；首版可先考虑只读调查视图、深链接或交接包，直到上游事实与动作能力得到验证。 |

## 仍必须补证的问题

1. 运营人员实际完成的任务是什么，其触发、过程、业务输出和结束标准是什么？
2. 该任务的频率、耗时、等待、返工、错误和失败影响是多少？
3. 17 个已通过工具中的哪些是该任务必需事实，哪些无关；6 个 blocked 和 1 个 skipped 是否真正阻断任务？
4. 任务所需字段是否完整、及时、语义一致，冲突／空结果／无权限／不存在能否可靠区分？
5. 哪个结果能证明任务完成，而不是只证明工具返回了数据？
6. FOBrain 原生操作、现有脚本或流程是否已经足够；AI 是否产生可测净增量？
7. 若最终需要写操作、人员范围、修复验证或跨系统交接，再补对应权限、动作和数据证据。

## 综合结论

- **查询与详情类能力：17 个只读工具已有当前 live 集成证据，不再重复调研其存在性；任务价值仍未知。**
- **责任范围：历史样本暴露身份与部门映射风险线索；当前状态待重新实测。**
- **证据化风险判断：历史材料只显示部分统计映射；判断规则、数据完整性和可解释性未知。**
- **完整风险案件与处置闭环：历史材料不足以支持任何当前可行性结论。**
- **受控工单状态更新：只有历史契约线索；目标部署 API、真实样本和用户必要性均未知。**

当前 live 结论使用 `EXPERIMENT_RESULT / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / OPEN / NOT_NORMATIVE`；历史条目仍保持 `PROJECT_HISTORY / SOURCE_ONLY / HYPOTHESIS`。两者均不能自动升级为 `TASK_EFFECTIVE`。
