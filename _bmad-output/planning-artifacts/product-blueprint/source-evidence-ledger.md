# 产品需求来源证据台账

## 使用规则

本台账只记录产品需求前置证据。历史文档中的架构、schema、接口、实施阶段和当前代码完成状态不进入新产品需求，除非它们揭示必须重新验证的产品约束。

本文件的 `状态` 是历史库存标签，不是完整证据状态。每个 ID 必须按 [research/evidence-status-model.md](research/evidence-status-model.md) 的“旧台账状态映射”派生 `source_type / claim_support / target_applicability / maturity / decision_status / normative_force`；下游工具不得只读取本列：

- `CONFIRMED`：仅用于 `U-*` 的 Owner 明确决定；外部来源不得用本值冒充 Owner 或真实用户证据。
- `SUPPORTED`：历史材料明确陈述且暂未发现冲突，但其目标适用性仍为 `UNKNOWN`，决策状态仍为 `HYPOTHESIS`。
- `ASSUMPTION`：单一历史来源、推断、愿望、继承性要求或方案性陈述。
- `CONFLICT`：来源对产品方向、边界或优先级给出不兼容答案。
- `MISSING`：完整需求蓝图必须回答、但当前 `source_type=NO_SOURCE`、`claim_support=NONE`。

## 已确认的流程与需求契约

| ID | 状态 | 陈述 | 来源 |
| --- | --- | --- | --- |
| U-001 | CONFIRMED | 当前 `docs/` 完整存档后只作为历史参考，不直接成为新需求真相。 | 用户确认，2026-07-11；项目 memlog |
| U-002 | CONFIRMED | 必须按照 BMAD 完成需求前置工作、Product Brief、PRFAQ、PRD 和整个产品需求蓝图。 | 用户确认，2026-07-11；项目 memlog |
| U-003 | CONFIRMED | 需求蓝图完成前不进入技术设计或代码实现。 | 用户确认，2026-07-11；项目 memlog |
| U-004 | CONFIRMED | 每个 Product Issue 必须包含 Background、Goal 和客观可验证的 Given/When/Then；每个业务目标至少映射一条 AC。 | 用户确认；项目 memlog |
| U-005 | CONFIRMED | 当前目标是实验性 Agent，接入部分 FOBrain 既有接口能力验证工作流；不构建完整漏洞运营平台。 | Owner 明确说明，2026-07-13 |
| U-006 | CONFIRMED | 本产品不做通知任何人的功能；FOBrain 自身自动通知不构成本产品能力、成功条件或验收范围。 | Owner 明确说明，2026-07-13 |
| U-007 | CONFIRMED | 本轮实验范围为 7 个代表性只读能力，以及手动派发、修复延时、误报标记、错派转发 4 个写动作；自动派发和其余工具不进入本轮。 | Owner 确认，2026-07-13；`experiment-scope-proposal.md` |

## 历史材料支持的产品方向

这些条目目前只能作为候选产品事实，不能直接进入 final PRD。

| ID | 状态 | 历史陈述 | 证据 | 需要验证 |
| --- | --- | --- | --- | --- |
| H-001 | SUPPORTED | 候选产品是聊天优先的 Agent Workbench，用户通过自然语言提交任务，系统可调用工具并返回结果。 | `docs/00-product-brief.md:3-5,17-20`；`docs/01-product-requirements.md:7-29` | 目标用户是谁、哪些任务值得使用 Agent 完成。 |
| H-002 | SUPPORTED | 候选产品同时面向 Web Workbench 用户和通过 API 提交任务的外部系统。 | `docs/00-product-brief.md:23-25`；`docs/01-product-requirements.md:5-13,52-57` | 两类入口是否属于同一产品、外部系统的真实使用者和业务场景。 |
| H-003 | SUPPORTED | 工具调用过程与结果需要以产品化卡片呈现，用户不应直接看到 raw provider 数据和内部技术信息。 | `docs/00-product-brief.md:20,27-33`；`docs/01-product-requirements.md:24-30`；`docs/02-ux-visual-requirements.md:15-59` | 用户需要多少过程透明度，哪些字段具有业务价值。 |
| H-004 | SUPPORTED | 高风险操作在执行前需要用户审批；信息不足时需要澄清，二者需要明确状态和可恢复结果。 | `docs/00-product-brief.md:21-22,31`；`docs/01-product-requirements.md:32-44` | 哪些业务操作属于高风险、谁有审批权、澄清对哪些实体和参数适用。 |
| H-005 | SUPPORTED | 用户可能需要查看一次任务的过程、证据、审计和回放。 | `docs/00-product-brief.md:23`；`docs/01-product-requirements.md:46-50`；`docs/02-ux-visual-requirements.md:83-92` | 哪些角色需要审计/回放、目的和保留范围是什么。 |
| H-006 | SUPPORTED | Workbench 与外部 API 应呈现同一次任务的一致、安全结果。 | `docs/00-product-brief.md:24,32`；`docs/01-product-requirements.md:11-13,52-57,88-93` | 这是用户需求、运营要求还是技术偏好；需要哪些跨入口场景。 |
| H-007 | SUPPORTED | 历史视觉目标是桌面三栏工作台：导航、主聊天、Inspector；移动端保持聊天和卡片可用。 | `docs/00-product-brief.md:25`；`docs/02-ux-visual-requirements.md:5-13,94-126` | 三栏结构是否来自用户研究，移动端是否真实需要。 |
| H-008 | SUPPORTED | Fobrain 安全运营数据查询与工单操作是候选业务领域；历史材料列出 24 个只读能力、连接器、实体澄清和工单状态更新。 | `docs/01-product-requirements.md:59-77`；`docs/fobrain-tool-matrix.md:58-93,117-126` | Fobrain 是核心产品、首个垂直场景还是验证样板；24 项能力是否全部仍有用户价值。 |
| H-009 | SUPPORTED | 历史材料排除了多租户、计费、完整 RBAC、worker 集群、高可用和全量生产迁移。 | `docs/01-product-requirements.md:79-86`；`docs/03-functional-scope.md:51-58` | 这些是永久非目标、首版非目标，还是旧重构项目的时间限制。 |

## 外部一手证据

| ID | 状态 | 外部事实 | 证据 | 对新产品的影响 |
| --- | --- | --- | --- | --- |
| E-001 | VENDOR_CLAIMED | FOBrain 官方把自身定义为网络资产攻击面管理平台，并公开声称覆盖多源资产与漏洞数据归一、资产关系、漏洞优先级排序和漏洞闭环处置。 | [华顺信安 FOBrain 官方产品页](https://www.huashunxinan.net/product-fobrain)，访问于 2026-07-11 | 这些名称不能被新项目直接宣称为差异化；目标部署是否授权、配置、可集成、被实际使用以及任务效果仍分别为 `UNKNOWN`。 |
| E-002 | VENDOR_CLAIMED | FOBrain 官方列出的目标痛点包括资产数据孤岛、漏洞数据处理耗费人力、线下通知无法闭环和攻击面可见性不足。 | [华顺信安 FOBrain 官方产品页](https://www.huashunxinan.net/product-fobrain)，访问于 2026-07-11 | 这只能证明厂商定位，可作为领域背景；不能自动成为新 AI 产品的用户问题，也不能证明目标客户在 FOBrain 上线后仍存在这些摩擦。 |

## 当前 live 技术证据

以下记录显式填写 canonical 字段，不使用旧库存标签映射。

| ID | 主张 | 直接证据 | Canonical 字段 | 产品边界 |
| --- | --- | --- | --- | --- |
| L-001 | Batch A 中当前用户、权限两个只读工具及 FOBrain connector 在 live provider 下通过。 | `test-results/eino-workbench-fobrain-batch-a-live-report.json`，2026-07-02 | `EXPERIMENT_RESULT / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / OPEN / NOT_NORMATIVE` | 证明调用、StructuredResult 和脱敏；不证明用户任务 |
| L-002 | Batch C 有 5 个查询／统计工具 live 通过，`pending_tickets` 跳过。 | `test-results/eino-workbench-fobrain-batch-c-live-report.json`，2026-07-06 | 通过项：`EXPERIMENT_RESULT / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / OPEN / NOT_NORMATIVE`；跳过项保持 `SOURCE_ONLY` | 工单读取仍不可宣称可用 |
| L-003 | Batch D 按 owner／department／IP 的 6 个资产／漏洞查询工具 live 通过。 | `test-results/eino-workbench-fobrain-batch-d-live-report.json`，2026-07-02 | `EXPERIMENT_RESULT / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / OPEN / NOT_NORMATIVE` | 参数化查询可用；不等于当前用户“我的范围” |
| L-004 | Batch E 的资产详情、漏洞详情、业务风险汇总、威胁关联 4 个工具 live 通过。 | `test-results/eino-workbench-fobrain-batch-e-live-report.json`，2026-07-03 | `EXPERIMENT_RESULT / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / OPEN / NOT_NORMATIVE` | 详情／汇总调用可用；不证明闭环或任务效果 |
| L-005 | Batch B 的 6 个“我的／本部门范围”工具在当前账号下因空结果或缺少 current-user scope 被阻塞。 | `test-results/eino-workbench-fobrain-batch-b-live-report.json`，2026-07-06 | `EXPERIMENT_RESULT / DIRECT / PARTIAL / DEPLOYMENT_AVAILABLE / CONFLICT / NOT_NORMATIVE` | 直接阻止 24 个只读全部通过；不否定参数化查询 |
| L-006 | live sample discovery 通过且安全报告不记录真实样本值。 | `test-results/eino-workbench-fobrain-sample-discovery-report.json`，2026-07-03 | `EXPERIMENT_RESULT / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / OPEN / NOT_NORMATIVE` | 可安全取得回放样本；代表性和使用授权另判 |
| L-007 | 当前 FOBrain provider mapper 可解析资产 `Status`，但 Product Facts 的 `StructuredResultRef` 只保留 `ResultRef` 与 `SafeSummary`；三个 FOBrain StructuredResult builder 的摘要只包含查询目标和返回条数，无法向产品出口提供行级主机状态。 | `internal/einoapp/providers/fobrain/http_client_batch_d.go:236-252`；`internal/einoapp/providers/fobrain/batch_c.go:97-112`；`internal/einoapp/providers/fobrain/batch_d.go:112-127`；`internal/einoapp/providers/fobrain/batch_e.go:160-178`；`internal/einoapp/product/structured_result.go:18-23,47-51` | `CURRENT_IMPLEMENTATION / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / CONFLICT / NOT_NORMATIVE` | 现有调用成功不等于能交付业务主机在线清单；这是产品需求输入，不在蓝图阶段决定技术实现。 |
| L-008 | FOBrain 源码定义资产 `status=1` 为在线、`status=2` 为离线；定时任务可按部署配置的最后数据源响应阈值将资产置为离线。业务系统 `running_state=1/0` 表示运行中／已下线，是不同语义。 | `/Users/vick/Desktop/project/fobrain/models/elastic/assets/assets.go:26-27,1046-1079`；`/Users/vick/Desktop/project/fobrain/fobrain/app/crontab/jobs_check_ds_response_at.go:17-69,91-93`；`/Users/vick/Desktop/project/fobrain/fobrain/app/services/asset_center/business_systems/business_systems_service.go:96` | `PROJECT_HISTORY / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 可作为 TS-001 的候选在线口径；必须确认目标部署版本、开关、阈值、状态更新时间和运营实际采用口径。 |
| L-009 | 当前 Eino Batch D 参数模型含 `Status`，但 owner／department／IP 三个资产查询分支未把资产状态作为查询条件发送；状态过滤只应用于漏洞分支。 | `internal/einoapp/providers/fobrain/batch_d.go:15-28,91-103`；`internal/einoapp/providers/fobrain/http_client_batch_d.go:49-80,95-107` | `CURRENT_IMPLEMENTATION / DIRECT / CONFIRMED / INTEGRATION_AVAILABLE / CONFLICT / NOT_NORMATIVE` | 现有资产查询 live pass 不能证明在线／离线筛选可用；不得据此写入正式 AC。 |
| L-010 | 当前 FOBrain capability catalog 包含业务系统列表、按负责人／部门／IP 查资产，却没有“按业务系统列出关联 IP 资产”的读取能力；而 TS-001 的实际步骤正是逐业务系统查询全部 IP 资产。 | `internal/einoapp/providers/fobrain/catalog.go:178,204-209`；TS-001 用户陈述，2026-07-12 | `MULTIPLE_SOURCES / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / CONFLICT / NOT_NORMATIVE` | 实际任务的核心读取步骤尚未被当前产品覆盖；不得以全局统计、负责人或部门查询冒充覆盖。 |
| L-011 | TS-001 的必交结果字段包含业务系统、主机总数、在线数、离线数、具体离线 IP、业务系统重要程度；而当前 Product Facts 不保留 IP／状态明细，`business_list` mapper 也未投影重要程度。 | TS-001 用户陈述，2026-07-12；`internal/einoapp/providers/fobrain/http_client_batch_c.go:229-241`；`internal/einoapp/product/structured_result.go:18-23,47-51` | `MULTIPLE_SOURCES / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / CONFLICT / NOT_NORMATIVE` | 现有调用成功不等于能交付老板所需结果表；这些字段是未来需求的必需输出，具体安全投影方案尚未决定。 |
| L-012 | FOBrain 历史前端和源码将漏洞状态 `0` 标记为“新增”；当前 Eino `liveThreatStatusCodes` 不接受数值 `0`，中文“新增”也未映射，`open` 则会筛选 `0/1/10/11/12/13/14/15/17` 多个状态。 | `/Users/vick/Desktop/project/fobrain_front/src/views/riskCenter/complianceRisks/detailColums.ts:188-198`；`/Users/vick/Desktop/project/fobrain/fobrain/app/repository/threat/threat.go:110`；`internal/einoapp/providers/fobrain/http_client_batch_d.go:382-398` | `MULTIPLE_SOURCES / DIRECT / PARTIAL / INTEGRATION_AVAILABLE / CONFLICT / NOT_NORMATIVE` | TS-002 需要“仅新增”漏洞；当前已通过查询不能精确表达该口径，目标部署状态值域也须验证。 |
| L-013 | FOBrain 源码存在自动派发配置、批量／全量异步派发、延时／催促／误报等漏洞状态操作、数据源新鲜度导致资产离线、以及合规风险派发闭环。 | `research/fobrain-source-derived-candidate-scenarios.md` 中列出的源码路径与行号 | `PROJECT_HISTORY / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 可用于主动生成候选场景；不得证明目标部署实际使用、用户价值或新产品范围。 |
| L-014 | FOBrain 源码同时存在直接“延时”状态操作和“申请延时→审批通过／拒绝”流程；二者的输入、权限和完成结果不同。 | `research/task-evidence-ts-005-repair-delay-disposition.md` 中列出的源码路径与行号 | `PROJECT_HISTORY / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | TS-005 必须先选择目标环境实际走的流程，不能合并为一个写域需求。 |
| L-015 | 直接延时支持批量漏洞、期限、备注和通知；仅允许特定原状态，逐条处理并可能部分失败，成功路径更新状态、期限和历史。 | `/Users/vick/Desktop/project/fobrain/fobrain/app/controller/threat_center/threat_distribute.go:204-355`；`/Users/vick/Desktop/project/fobrain/fobrain/app/services/threat_center/threat_center.go:259-384`；`/Users/vick/Desktop/project/fobrain/fobrain/app/repository/threat_history/threat_distribute.go:81-105` | `PROJECT_HISTORY / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 可形成 PI-005 草案的候选契约；不证明目标部署实际参数、授权或效果。 |
| L-016 | FOBrain 源码可在资产运维／业务负责人变更时按配置自动转交已派发漏洞；目标人员不存在、无修复负责人或状态不适用时跳过。 | `/Users/vick/Desktop/project/fobrain/fobrain/app/repository/threat_history/auto_transfer.go:17-129,245-258` | `PROJECT_HISTORY / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 可形成“责任人变更后的转交结果对账”候选；不证明目标部署已启用或存在人工痛点。 |
| L-017 | FOBrain 手动派发源码接口接受多个漏洞 ID 和 1～10 名去重接收人，逐条校验状态并更新；执行前可按配置过滤离线资产，过滤后无对象时直接返回成功。 | `/Users/vick/Desktop/project/fobrain/fobrain/app/controller/threat_center/threat_distribute.go:43-161`；`/Users/vick/Desktop/project/fobrain/fobrain/app/services/threat_center/threat_center.go:115-226,393-410` | `PROJECT_HISTORY / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 可作为 TS-002 的候选接口约束；不证明目标部署、真实运营批量方式或写入完成。产品必须回读负责人变更，不能以接口成功替代业务成功。 |
| L-018 | FOBrain 自动派发配置源码读取／写入单个 `poc_auto_distribute_config` 记录，配置模型没有业务系统 ID；自动派发执行时才从每条漏洞关联的业务系统负责人取人。 | `/Users/vick/Desktop/project/fobrain/fobrain/app/request/poc_auto_distribute_config/poc_auto_distribute_config.go:7-33`；`/Users/vick/Desktop/project/fobrain/models/mysql/poc_auto_distribute_config/poc_auto_distribute_config.go:7-24`；`/Users/vick/Desktop/project/fobrain/fobrain/app/repository/poc_auto_distribute_config/poc_auto_distribute_config.go:71-145,620-735`；Owner 澄清，2026-07-13 | `MULTIPLE_SOURCES / DIRECT / UNKNOWN / SOURCE_ONLY / DEFERRED / NOT_NORMATIVE` | “按业务系统自动派发”是本产品后续新增能力，不是当前 FOBrain 已有规则或本轮范围；源码只能作为未来接口支撑线索。 |
| L-019 | FOBrain 源码把误报作为批量状态操作，把责任人转发作为批量派发操作；二者允许的原状态、接收人、状态结果和通知对象均不同。 | `/Users/vick/Desktop/project/fobrain/fobrain/app/controller/threat_center/threat_distribute.go:43-161,180-351`；`/Users/vick/Desktop/project/fobrain/fobrain/app/services/threat_center/threat_center.go:230-410` | `PROJECT_HISTORY / DIRECT / UNKNOWN / SOURCE_ONLY / OPEN / NOT_NORMATIVE` | 只能用于建立 PI-003／PI-004 的候选边界；不证明真实触发、频率、权限、目标部署或用户价值。 |

合计：24 个历史只读工具中 17 个 `passed`、6 个 `blocked`、1 个 `skipped`；connector 另行 `passed`。这已经回答“哪些集成真实可用”，后续 Research 不再重复盘点已通过接口，只验证任务、价值、缺口和产品关系。

## 历史假设

| ID | 状态 | 假设 | 证据 | 风险 |
| --- | --- | --- | --- | --- |
| A-001 | ASSUMPTION | 产品重建的主要价值是练习、比较架构和降低技术负债。 | `docs/00-product-brief.md:7-15` | 这是研发动机，不是用户问题；不能作为产品价值主张。 |
| A-002 | ASSUMPTION | 旧产品体验和全部关键能力都应原样成为新产品范围。 | `docs/00-product-brief.md:11,35-43`；`docs/03-functional-scope.md:5-49` | 可能把历史实现偶然性和低价值能力带入新产品。 |
| A-003 | ASSUMPTION | 24 个 Fobrain 只读工具天然构成完整且正确的用户能力集合。 | `docs/fobrain-tool-matrix.md:58-85` | 工具清单按旧能力恢复组织，并非按用户任务、频率或价值验证。 |
| A-004 | ASSUMPTION | 三栏布局和密集工作台是所有目标用户的最佳体验。 | `docs/02-ux-visual-requirements.md:5-13,105-126` | 有视觉基线但没有用户研究、任务效率或可用性证据。 |
| A-005 | ASSUMPTION | Action API 是独立且必要的产品入口。 | `docs/00-product-brief.md:24`；`docs/01-product-requirements.md:11-13,52-57` | 没有外部调用方、调用频率、自动化场景和成功结果证据。 |
| A-006 | ASSUMPTION | 审计、回放和 Inspector 应全部直接面向普通 Workbench 用户。 | `docs/00-product-brief.md:23`；`docs/02-ux-visual-requirements.md:83-92` | 这些能力可能属于管理员、开发者或合规角色，而非同一用户。 |
| A-007 | ASSUMPTION | 旧 runtime、接口和调试成本严重到值得完整重做。 | `docs/00-product-brief.md:7-15` | 没有维护成本、故障、交付周期或复杂度基线；即使技术重做合理，也不能替代用户价值证明。 |
| A-008 | ASSUMPTION | 将 FOBrain 数据接入自然语言聊天就会产生独立产品价值。 | 历史材料与 E-001 的对照推断。**Canonical exception：** `MULTIPLE_SOURCES / SYNTHESIS / UNKNOWN / SOURCE_ONLY / HYPOTHESIS / NOT_NORMATIVE`。 | FOBrain 官方已声明覆盖底层聚合、排序与闭环，但目标部署和任务效果未知；必须验证候选解法能否改善一个通过 DG-04 准入的重要任务。 |

## 方向冲突

| ID | 状态 | 冲突 | 来源 | 需要决策 |
| --- | --- | --- | --- | --- |
| C-001 | CONFLICT | 历史材料一方面把产品描述为通用 Agent Workbench，另一方面把完整完成定义为恢复 Fobrain 全量能力。 | `docs/00-product-brief.md:5,33,35-57`；`docs/01-product-requirements.md:59-77` | 新产品究竟是通用 Agent 工作平台、Fobrain 安全运营助手，还是平台加首个垂直产品。 |
| C-002 | CONFLICT | 历史材料明确说“不是新产品探索”，用户现已要求把项目作为新产品重新完成需求蓝图。 | `docs/00-product-brief.md:11`；U-001/U-002 | 用户的新决定优先；旧范围只能作为候选证据，不能继续作为强制完成定义。 |
| C-003 | CONFLICT | 历史成功标准大多是功能链路可运行和视觉可比，没有用户结果或业务指标。 | `docs/00-product-brief.md:45-57`；`docs/01-product-requirements.md:88-94` | 必须重新定义用户成功、业务成功和可量化结果。 |
| C-004 | CONFLICT | Fobrain 被称为“业务能力样板、不能写进 runtime”，但其 24 项能力又被视为产品完成的必做范围。 | `docs/00-product-brief.md:33,57`；`docs/03-functional-scope.md:37-49` | 需要明确“技术样板”和“产品核心范围”的关系。 |
| C-005 | CONFLICT | “不直接替换当前主项目”与“P2 未完成前不得替换”暗示的最终去向不同，同时全量生产迁移又被延后。 | `docs/01-product-requirements.md:79-86`；`docs/03-functional-scope.md:49,51-58` | 需要重新定义产品的最终使用方式，不继承旧实施阶段结论。 |
| C-006 | CONFLICT | 历史产品材料声称只描述“要什么”，但混入 Eino、Product Facts、StructuredResult、runtime core 等具体解法。 | `docs/README.md:5-11,26-27,64-67`；`docs/00-product-brief.md:27-33` | 新 PRD 必须移除实现解法，只保留用户结果和可验证产品约束。 |
| C-007 | CONFLICT | “当前项目已经验证”没有明确指旧项目还是本仓库，而遗留证据索引明确旧证据不能作为新项目通过结果。 | `docs/00-product-brief.md:9`；`docs/legacy-acceptance-evidence.md:3-5,70-75` | 历史能力存在性与新需求正确性必须分别记录。 |
| C-008 | CONFLICT | 候选定位“独立资产与漏洞风险运营工作台”与 FOBrain 官方声明的产品范围高度重叠。 | E-001；`_bmad-output/forge/agent-workbench-product/positioning-evidence.md` 初始候选 | 先在目标 FOBrain 部署中发现一个真实剩余任务，再比较原生增强、内嵌 AI、独立 Workbench、非 AI 改进和不做；不预设 AI 协作层就是答案。 |
| C-009 | CONFLICT | FOBrain 官方声称具备完整漏洞闭环；归档 AI 集成矩阵曾只列出待处理工单读取和审批后更新工单状态，归档 live plan 还记录当时缺少待处理工单样本。 | E-001；`docs/fobrain-tool-matrix.md:77,87-93,117-126`；`docs/fobrain-live-read-batch-plan.md:56-64` | 必须区分厂商总体声明、归档时历史线索、目标部署数据、目标 API 暴露能力和新产品需求；不能假定任一者等同。 |

## 缺失信息

| ID | 状态 | 缺失问题 | 为什么阻断 |
| --- | --- | --- | --- |
| M-001 | MISSING | 主要用户的具体角色、职责、经验、权限和工作环境。 | 无法定义真实用户旅程、语言、风险和优先级。 |
| M-002 | MISSING | 用户当前最重要的工作任务、痛点、频率、现有替代方式和失败成本。 | 无法证明 Agent Workbench 或任何功能值得建设。 |
| M-003 | MISSING | 购买者、产品所有者、审批者、管理员、审计者和外部系统调用者是否为不同角色。 | 当前文档把多个角色压缩成“用户”。 |
| M-004 | MISSING | 产品级业务目标、用户结果、基线值、目标值和时间范围。 | 当前“能跑通”只能证明系统存在，不能证明产品成功。 |
| M-005 | MISSING | 产品定位：FOBrain 的 AI 协作入口、可替换 FOBrain 的独立风险运营产品、通用 Agent 平台或组合产品。 | 它阻止 final Brief 与范围，但不阻止发现；必须先获得单任务证据、DG-05 判定规则、DG-06 形态和采用／买方证据，再由 DG-07 决定。 |
| M-006 | MISSING | 首版范围的价值排序和删除标准。 | 历史 P0/P1/P2 按实现阶段组织，不等于产品 MVP。 |
| M-007 | MISSING | 数据隐私、合规、审计保留和人工责任边界。 | 历史资料只有展示脱敏规则，没有业务或法规要求来源。 |
| M-008 | MISSING | 市场替代品、直接竞品、内部替代流程和“不做”的机会成本。 | 无法判断差异化和必要性。 |
| M-009 | MISSING | AI 出错、工具误选、结果不完整和用户不信任时的产品补救策略。 | Agent 产品的核心风险尚未转化为用户需求。 |
| M-010 | MISSING | Action API、移动端、Inspector、回放和全量 Fobrain 能力的真实使用证据。 | 这些高成本范围目前主要来自继承性声明。 |
| M-011 | MISSING | 产品负责人、最终用户代表、Fobrain 业务负责人、数据负责人、安全/合规、运维和 API 消费方的决策权。 | 无法判断冲突需求由谁裁决，也无法规划需求验证。 |
| M-012 | MISSING | 性能、容量、可用性、无障碍、浏览器支持、数据保留、隐私合规、灾备和成本等产品级目标。 | 这些不能由架构自行决定，必须先形成可验证的产品约束。 |

## 第一批必须解决的依赖问题

1. 先由产品负责人确认**发现边界**：研究 FOBrain 活跃用户的剩余任务，明确不预设内嵌、独立、平台或通用 Agent 形态。
2. 解决 M-011：指定能提供目标 FOBrain 版本、真实用户、近期任务和脱敏数据的证据负责人。
3. 解决 M-001/M-002/M-003：通过真实任务证据明确执行角色、协作角色、任务、频率、现状成本和失败后果。
4. 对选定任务先解决 M-004：预注册一个核心结果、基线、目标值、保护指标、周期和停止条件。
5. 按同一判定规则比较 FOBrain 原生增强、内嵌／伴随式 AI、独立 Workbench、非 AI 改进和不做；同时收集采用责任与买方证据。
6. 证据形成后再解决 M-005/C-001/C-008：确定产品关系、商业性质、最终定位与首版范围。

其余历史能力在上述问题解决前均保持候选状态，不得进入 final PRD。
