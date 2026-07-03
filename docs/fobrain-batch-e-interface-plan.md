# Fobrain Batch E Interface Plan

本文只读整理 Phase 8 Batch E 四个工具的接口证据、输入约束和实现顺序。旧项目和 Fobrain 源码只作为产品能力、验收基线和安全边界参考；新实现仍必须走 Eino-first provider、StructuredResult 和 Product Facts。

## 范围

| 工具 | 输入 | 展示类型 | 目标 |
| --- | --- | --- | --- |
| `tool.fobrain.get_asset_detail` | `asset_id`，可选 `network_type` | `entity_detail` | 读取单个资产安全详情。 |
| `tool.fobrain.get_vulnerability_detail` | `vulnerability_id` | `entity_detail` | 读取单个漏洞安全详情。 |
| `tool.fobrain.business_risk_summary` | `business_name` | `metrics_summary` | 汇总业务系统关联漏洞风险。 |
| `tool.fobrain.threat_relevance_list` | `vulnerability_name`，可后续扩展 IP/业务过滤 | `entity_collection` | 查询指定漏洞/威胁关联资产。 |

## 旧源码接口证据

| 工具 | 标准路径 | 方法 | 参数形态 | 源码证据 | 当前实现建议 |
| --- | --- | --- | --- | --- | --- |
| `get_asset_detail` | `/api/v1/internal_asset/:id`、`/api/v1/external_ip_asset/:id`、`/api/v1/device/:id`、`/api/v1/domain_asset/:id` | GET | path `id`；`network_type` 决定路径 | `routes/asset_center/asset_center.go`；`internal_asset.Show`、`external_ip_asset.Show`、`device_asset.Show` | `network_type` 缺省先按旧 adapter 使用 internal；实现必须集中归一化 `1/2`、中英文内外网和 device/domain alias。 |
| `get_vulnerability_detail` | `/api/v1/threat_center/:id` | GET | path `id` | `routes/threat_center/threat_center.go`、`threat_center.Show`、`threat.Show` | 当前私有路径优先 `/api/threat_center/:id`，fallback `/api/v1/threat_center/:id`。 |
| `business_risk_summary` | `/api/v1/threat_center/count` | POST | body 数组：`count_name=business_risk`、`aggregation_field=business.name.keyword`、`data_range=4`、`search_condition` | 旧 adapter `BusinessRiskSummary`；Fobrain `threat_center.Count`、`threat.Count` | 不走自然语言总结；只把统计 bucket 归一为安全 metrics。 |
| `threat_relevance_list` | `/api/v1/threat_center/relevance/list` | GET | `page`、`per_page`、`keyword`、`vul_name`，可选 `ip`、`business_name` | `routes/threat_center/threat_center.go`、`relevance.go`、`VulRelevanceListRequest` | `vulnerability_name` 同时写入 `keyword` 和 `vul_name`，与旧 adapter 保持行为。 |

## 字段边界

- 资产详情允许字段：`id`、`ip`、`hostname`、`status`、`network_type`、`ip_type`、业务系统名、负责人名、风险摘要计数。
- 漏洞详情允许字段：`id`、`name`、`level`、`status`/`statusCode`、`cve`/`cnvd`/`cnnvd`、`person_info[].name`、`risk_num`、安全描述摘要。
- 业务风险允许字段：统计 label、count、aggregation bucket；不得把 POST body、raw ES query 或 raw prompt 写入 Product Facts。
- 威胁关联允许字段：漏洞名、等级、关联数量、CVE/CNVD/CNNVD、安全摘要和关联项计数。

## 实现顺序

1. 扩展 Batch E catalog/input mapper 和 mock StructuredResult，先覆盖四工具成功、空态、缺参和敏感字段拒绝。
2. 实现 HTTP live mapper，所有路径先尝试当前私有 `/api/...`，再 fallback `/api/v1/...`。
3. 增加 Batch E sample discovery：详情 ID 只能来自 ignored local samples 或新的脱敏 discovery sidecar；不得从可提交 Batch D live report 推导真实 ID，也不得回读 raw provider payload。
4. 增加 Batch E live smoke/report schema，passed 必须证明四工具返回 resolved 或 documented not_found，且 raw payload/token 不出现。
5. 补 Workbench 视觉和 replay/audit 证据后，才允许声明 Batch E 完成。

当前代码进展：步骤 1-4 的基础设施已完成 provider catalog、输入 mapper、mock StructuredResult、HTTP mapper、sample discovery 本地样本扩展、Batch E live smoke 脚本和 `batch_e_live_report` schema。HTTP mapper 已覆盖资产详情 `network_type` 中央归一、漏洞详情非分页响应、业务风险 count POST body、威胁关联 `keyword/vul_name` 双参数映射和 `/api` 到 `/api/v1` fallback。真实环境 live pass 报告和步骤 5 仍未完成，因此不得声明 Batch E 完成或 Fobrain 24 只读恢复完成。

输入契约进展：`asset_detail.network_type` schema 允许当前真实输入中的数字、中文和英文别名，但 provider 边界必须归一为 `internal`、`external`、`device`、`domain` 后再进入路径选择；`business_risk_summary` 使用专用 `business_risk_query`，避免影响 Batch C 的 `business_list`；`threat_relevance_list` 明确支持可选 `ip` 与 `business_name` 过滤。

## 阻塞条件

- `get_asset_detail` 无法获得真实资产 ID 或资产类型时，只能生成 blocking report；不得猜测 ID。
- `network_type` 必须先归一化为 `internal`、`external`、`device` 或 `domain`，不能把 provider 原始数字、中文或资产类别字段散落到业务分支。
- `business_risk_summary` 若真实部署不支持 `/threat_center/count` 的 `business.name.keyword` 聚合，必须记录 `connector_execution_failed`，不得由模型编造摘要。
- `threat_relevance_list` 允许安全空态，但 live pass report 必须区分 empty 与 transport/provider failure。
- 任一工具出现 raw provider payload、凭据、auth header、完整 ES hit 或本地配置内容，必须失败。
