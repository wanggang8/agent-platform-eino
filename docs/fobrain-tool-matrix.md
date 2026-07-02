# Fobrain 工具恢复矩阵

本文定义 P2/Phase 8 必须恢复的 Fobrain 产品能力。它是产品能力矩阵，不要求复制旧实现代码；新实现仍以 Eino-first 架构、Capability Provider、StructuredResult 和安全投影完成。

本矩阵不是硬编码旧实现。它只固定产品可见能力、输入输出契约、安全展示和验收断言；具体实现必须服从新架构的 provider、schema、Product Facts 和 Safety Gate。

Phase 8 是完整重构必做门禁。未完成本矩阵，不得声明重构完成、Fobrain 能力可比或替换当前产品基线。

Phase 8 必须按 `docs/fobrain-live-read-batch-plan.md` 分批恢复：Batch A 先完成 connector、当前用户和权限，再进入我的范围、直接列表/统计、参数化查询、详情与风险关联。不得跳过批次直接声明 24/24 完成。

## 通用要求

机器可读矩阵固定在：

```text
docs/schemas/fobrain/tool_matrix.v1.schema.json
docs/schemas/fobrain/tool_inputs.v1.schema.json
docs/fixtures/fobrain/tool-matrix-24.json
```

`tool-matrix-24.json` 必须始终登记在 `docs/fixtures/manifest.json`，确保 24 个只读工具的 tool id、输入契约引用、结果契约、mock/live 断言和截图状态可被 contract test 覆盖。

每个工具必须定义：

- tool id。
- 中文展示名。
- 输入 schema。
- StructuredResult schema。
- 分批门禁。
- 典型真实模型 prompt。
- mock 断言。
- live 断言。
- 工具卡截图要求。
- 敏感字段规则。
- 是否纳入 24/24 只读矩阵。

所有输出必须进入 `tool.structured_result.v1`，Fobrain 业务结果使用 `fobrain.tool_result.v2`。不得新增独立 `result_summary`、`display_items` 或 `user_safe_summary` 事实字段。

每个工具进入实现前必须补齐字段级迁移信息：

| 字段 | 要求 |
| --- | --- |
| `input_schema_ref` | 指向 `docs/schemas/` 中的输入 schema |
| `result_schema_ref` | 固定为 `fobrain.tool_result.v2` 或更具体的业务结果 schema |
| `required_fields` | 必填输入字段和默认值 |
| `batch_gate` | 固定为 `A` / `B` / `C` / `D` / `E`，与 live read 分批计划一致 |
| `display_type` | table / detail / metrics / narrative / connector_status / approval |
| `entity_type` | person / asset / vulnerability / department / business / ticket / connector |
| `result_status` | resolved / waiting / pending_approval / not_found / empty / failed / partial |
| `fixture` | 对应 fixture 文件名 |
| `real_model_prompt_id` | 真实模型工具选择用例 id |
| `mock_required` | 是否必须 mock 覆盖 |
| `live_required` | 是否必须 live 覆盖；无凭据时必须生成 blocking skip |
| `screenshot_state` | collapsed / expanded / detail / approval / clarification |

字段级恢复信息不得引用旧代码包名、旧 runtime step id、旧 resume token、旧 UI DOM 结构或 raw provider payload 字段。需要表达旧能力时，只能转换成新 schema、fixture、StructuredResult 字段和验收断言。

## 24 个只读工具

| # | Tool id | 中文展示名 | 输入范围 | StructuredResult 要求 | 典型 prompt | Mock / live 断言 | 截图 | 敏感规则 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | `tool.fobrain.current_user_context` | 读取当前用户信息 | 当前用户 | user context、部门、角色安全摘要 | 查看当前 Fobrain 用户信息 | 返回当前用户安全字段；live 用户名非空 | 工具卡展开 | 不展示 token/account raw |
| 2 | `tool.fobrain.my_permissions` | 读取我的权限范围 | 当前用户 | 权限、数据范围安全摘要 | 我的权限范围是什么 | 权限列表可读；无权限返回空态 | 工具卡展开 | 不展示 raw policy payload |
| 3 | `tool.fobrain.list_assets_by_owner` | 查询负责人资产 | person_name/person_staff_id、分页 | 资产表格、总数、分页 | 查张三负责的资产 | 同名人员触发澄清或返回资产 | 工具卡折叠/展开 | 不展示 raw credential |
| 4 | `tool.fobrain.list_vulnerabilities_by_owner` | 查询负责人漏洞 | person_name/person_staff_id、分页、可选严重度 | 漏洞表格、风险摘要 | 查张三负责的漏洞 | 结果包含漏洞名称/等级/状态 | 工具卡展开 | 不展示 provider body |
| 5 | `tool.fobrain.list_assets_by_department` | 查询部门资产 | department_name、分页 | 资产表格、部门摘要 | 查安全部有哪些资产 | 部门名可直接首读 | 工具卡展开 | 部门路径按业务值展示 |
| 6 | `tool.fobrain.list_vulnerabilities_by_department` | 查询部门漏洞 | department_name、分页、可选严重度 | 漏洞表格、部门风险摘要 | 查安全部有哪些漏洞 | 部门漏洞可展示 | 工具卡展开 | 不展示 raw query |
| 7 | `tool.fobrain.list_assets_by_ip` | 查询 IP 资产 | ip、可选分页/筛选 | IP 资产表格 | 查 10.10.11.226 的资产 | IP-only 直接调用 | 工具卡展开 | IP 可完整展示 |
| 8 | `tool.fobrain.list_vulnerabilities_by_ip` | 查询 IP 漏洞 | ip、可选严重度/状态/分页 | IP 漏洞表格 | 查 10.10.11.226 的漏洞 | 可复用上一轮 IP | 工具卡展开 | IP/CVE 可完整展示 |
| 9 | `tool.fobrain.get_asset_detail` | 查询资产详情 | asset_id | 资产详情 facts/sections | 查询资产 asset-1 详情 | real asset id 直接查询 | 详情卡 | safe ref 不当作 real id |
| 10 | `tool.fobrain.get_vulnerability_detail` | 查询漏洞详情 | vulnerability_id | 漏洞详情 facts/sections | 查询漏洞 vuln-1 详情 | real vuln id 直接查询 | 详情卡 | CVE 可展示 |
| 11 | `tool.fobrain.business_risk_summary` | 汇总业务风险 | business_name、time_range、分页 | 指标、关键发现、风险判断 | 汇总支付系统业务风险 | 风险摘要可读 | 指标+表格 | 不展示 raw prompt |
| 12 | `tool.fobrain.business_list` | 查询业务系统 | owner/keyword、分页 | 业务系统表格 | 查询业务系统列表 | 无筛选可首读 | 表格卡 | 不展示 raw provider payload |
| 13 | `tool.fobrain.external_high_risk_assets` | 查询外部高风险资产 | severity、分页 | 高风险资产表格 | 查外网严重资产 | 不调用 connector 状态替代业务读 | 表格卡 | 风险业务值可展示 |
| 14 | `tool.fobrain.threat_relevance_list` | 查询威胁关联资产 | vulnerability_name/vulnerability_id、分页 | 关联资产/漏洞表格 | 查 Log4j 关联资产 | 命名威胁直接调用 | 表格卡 | 不展示 raw search body |
| 15 | `tool.fobrain.vulnerability_status_summary` | 汇总漏洞状态 | severity、time_range | 状态分布指标 | 汇总漏洞状态 | 不误用漏洞数量统计 | 指标卡 | 状态码需转安全标签 |
| 16 | `tool.fobrain.pending_tickets` | 查询待处理工单 | person、status、分页 | 工单表格 | 查询我的待处理工单 | 工单 id/status 可展示 | 工单表格 | mutation 仍需审批 |
| 17 | `tool.fobrain.ip_stats` | 统计 IP 资产 | field、time_range | IP 资产统计指标 | 统计 IP 资产情况 | 不用于单 IP 查询 | 指标卡 | 不展示 raw aggregation body |
| 18 | `tool.fobrain.vul_stats` | 统计漏洞情况 | field、time_range | 漏洞数量统计指标 | 统计漏洞情况 | 不用于状态分布问题 | 指标卡 | 不展示 raw aggregation body |
| 19 | `tool.fobrain.my_assets` | 查询我的资产 | keyword、分页 | 当前用户资产表格 | 我的资产有哪些 | 内部解析当前用户 | 表格卡 | 不展示 credential |
| 20 | `tool.fobrain.my_department_assets` | 查询本部门资产 | keyword、分页 | 本部门资产表格 | 查我部门资产 | 不追问部门名 | 表格卡 | 部门业务值可展示 |
| 21 | `tool.fobrain.my_vulnerabilities` | 查询我的漏洞 | keyword、分页 | 当前用户漏洞表格 | 我的漏洞有哪些 | 不先追问严重度 | 表格卡 | 不展示 raw query |
| 22 | `tool.fobrain.my_department_vulnerabilities` | 查询本部门漏洞 | keyword、分页 | 本部门漏洞表格 | 查我部门漏洞 | 不追问部门名/严重度 | 表格卡 | 不展示 raw query |
| 23 | `tool.fobrain.my_business_systems` | 查询我的业务系统 | keyword、分页 | 业务系统表格 | 查询我的业务系统 | 当前用户上下文可内部解析 | 表格卡 | 不展示 raw user payload |
| 24 | `tool.fobrain.my_important_business_systems` | 查询我的重要业务系统 | keyword、分页 | 重要业务系统表格 | 查询我的重要业务系统 | important_only 生效 | 表格卡 | 不展示 raw user payload |

## 辅助能力

| 能力 | 是否纳入 24/24 | 要求 |
| --- | --- | --- |
| `tool.fobrain.query_normalize_and_resolve` / `CapabilityQueryNormalizeAndResolve` | 否 | 仅用于实体识别和候选生成；候选结果进入 clarification，不作为业务读结果替代 |
| `connector.fobrain.security` | 否 | 只回答连接器健康、配置、凭据绑定；不能替代业务数据读取 |
| `tool.fobrain.update_ticket_status` | 否 | 写域 mutation，必须走审批；未审批不得调用真实更新接口 |

## 字段级恢复模板

Phase 8 开始前，每个工具必须按以下模板补齐一行；未补齐不得进入实现：

```text
tool_id:
display_name_zh:
input_schema_ref:
required_fields:
result_schema_ref:
batch_gate:
display_type:
entity_type:
result_status:
fixture:
real_model_prompt_id:
mock_assertions:
live_assertions:
screenshot_state:
sensitive_fields:
```

## 写域审批要求

`tool.fobrain.update_ticket_status` 必须满足：

- 输入包含 ticket id/ref、目标状态、影响摘要。
- `fixed`、`ignored`、`in_progress` 等状态必须标准化。
- 审批卡展示操作名称、工单、目标状态、风险摘要。
- approve 后才执行真实 mutation。
- reject/cancel/duplicate/restart 后行为稳定。
- live write report 必须证明未审批不执行、审批后只执行一次。
