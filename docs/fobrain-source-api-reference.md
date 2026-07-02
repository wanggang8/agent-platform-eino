# Fobrain Source API Reference

本文只读整理 `../fobrain` 源码中与 Phase 8 live read 相关的接口、字段和样本发现流程。它用于后续实现和验收参考，不表示复制旧源码结构、旧 runtime 或 raw payload。

## 源码入口

标准接口由 Gin 路由挂在 `/api/v1` 下：

- 路由根：`../fobrain/fobrain/common/middleware/init.go`
- 路由注册：`../fobrain/fobrain/routes/init.go`
- 资产路由：`../fobrain/fobrain/routes/asset_center/asset_center.go`
- 漏洞路由：`../fobrain/fobrain/routes/threat_center/threat_center.go`
- 当前用户路由：`../fobrain/fobrain/routes/user/user.go`
- 人员/部门路由：`../fobrain/fobrain/routes/user_access/user_access.go`、`../fobrain/fobrain/routes/personnel_departments/personnel_departments.go`

响应包装来自 `../fobrain/fobrain/common/response/response.go`：成功响应为 `{code,message,data}`；分页数据在 `data.items`，分页字段为 `page`、`per_page`、`total`。

## 样本发现接口

稳定 `owner / department / ip` 样本不应手写猜测，应从只读接口中挑选能同时命中资产和漏洞的值。

| 目的 | 标准接口 | 关键字段 | 说明 |
| --- | --- | --- | --- |
| 当前用户 | `GET /api/v1/user` | `username`、`role`、`rule_info` | 可作为 owner 候选；源码模型不保证返回部门。当前私有部署可兼容 `GET /api/user`。 |
| 人员样本 | `GET /api/v1/user_access/staff_list?page=1&per_page=20` | `items[].name`、`items[].department[]` | `name` 可作为 owner 候选，`department[]` 可作为 department 候选。 |
| 部门样本 | `GET /api/v1/personnel_departments?page=1&per_page=20` | `items[].name` | 可作为 department 候选，但仍需资产/漏洞验证。 |
| 资产样本 | `GET /api/v1/asset?page=1&per_page=20` | `items[].ip`、`items[].oper_info[].name`、`items[].business_department[].name` | 最可靠的 Batch D 样本来源。当前私有部署可兼容 `GET /api/asset`。 |
| 漏洞验证 | `GET /api/v1/threat_center?page=1&per_page=20&data_range=4` | `items[].ip`、`items[].person_info[].name`、`items[].person_department[].name` | 用于验证同一 owner/department/ip 是否能命中漏洞。当前私有部署可兼容 `GET /api/threat_center`。 |

推荐流程：

1. 调用资产列表，找一条 `ip` 非空、`oper_info[].name` 非空、`business_department[].name` 非空的资产。
2. 使用该资产的 `ip` 调用漏洞列表：`GET /api/v1/threat_center?page=1&per_page=20&data_range=4&search_condition={"ip":["<IP>"],"operation_type_string":"=="}`。
3. 如果漏洞列表也有数据，再确认 `person_info[].name` 或 `person_department[].name` 与候选 owner/department 一致或可关联。
4. 只有三类样本都稳定存在时，才运行 Batch D live pass smoke。

## Batch D 查询映射

新项目 provider 当前先尝试私有部署路径 `/api/asset`、`/api/threat_center`，再 fallback 到旧源码标准路径 `/api/v1/asset`、`/api/v1/threat_center`。部署配置的 `base_url` 可指向服务根、`/api` 或 `/api/v1`，client 会归一化路径。

| 能力 | 标准接口 | 查询参数 | 源码字段证据 |
| --- | --- | --- | --- |
| `tool.fobrain.list_assets_by_owner` | `GET /api/v1/asset` | `page`、`per_page`、`search_condition={"oper_info.name":["<owner>"],"operation_type_string":"=="}` | 资产模型 `oper_info[].name` |
| `tool.fobrain.list_vulnerabilities_by_owner` | `GET /api/v1/threat_center` | `page`、`per_page`、`data_range=4`、`search_condition={"person_info.name":["<owner>"],"operation_type_string":"=="}` | 漏洞模型 `person_info[].name` |
| `tool.fobrain.list_assets_by_department` | `GET /api/v1/asset` | `page`、`per_page`、`search_condition={"business_department.name.keyword":["<department>"],"operation_type_string":"=="}` | 资产模型 `business_department[].name` |
| `tool.fobrain.list_vulnerabilities_by_department` | `GET /api/v1/threat_center` | `page`、`per_page`、`data_range=4`、`search_condition={"person_department.name.keyword":["<department>"],"operation_type_string":"=="}` | 漏洞模型 `person_department[].name` |
| `tool.fobrain.list_assets_by_ip` | `GET /api/v1/asset` | `page`、`per_page`、`keyword=<ip>`、`search_condition={"ip":["<ip>"],"operation_type_string":"=="}` | 资产模型 `ip` |
| `tool.fobrain.list_vulnerabilities_by_ip` | `GET /api/v1/threat_center` | `page`、`per_page`、`data_range=4`、`search_condition={"ip":["<ip>"],"operation_type_string":"=="}` | 漏洞模型 `ip` |

`search_condition` 在旧源码中是字符串数组，由 `ParseQueryConditions` 解析。HTTP query 中可重复传入 `search_condition=<json>`；当前新项目 Batch D client 一次只传一个 JSON 字符串。

注意：旧源码 `ThreatListRequest.Ip` 会让 controller 在 `len(params.Ip)>0` 时强制 `data_range=1`，这是 IP 画像场景的历史行为。Batch D 验收要求查询非回收站完整范围，因此漏洞按 IP 查询不得发送 `ip=<ip>` query 参数，必须使用 `search_condition` 约束 `ip` 字段。

## Batch E 查询映射

Batch E 详情和风险关联的字段级计划见 `docs/fobrain-batch-e-interface-plan.md`。实现时必须遵循以下已确认接口证据：

| 能力 | 标准接口 | 参数 | 说明 |
| --- | --- | --- | --- |
| `tool.fobrain.get_asset_detail` | `GET /api/v1/internal_asset/:id`、`/api/v1/external_ip_asset/:id`、`/api/v1/device/:id`、`/api/v1/domain_asset/:id` | path `id`，可选 `network_type` 决定路径 | 旧 adapter 默认未知类型走 internal；新实现不得只硬编码一种资产类型而不记录 fallback。 |
| `tool.fobrain.get_vulnerability_detail` | `GET /api/v1/threat_center/:id` | path `id` | 成功响应为 `{code,message,data}`，不是分页列表。 |
| `tool.fobrain.business_risk_summary` | `POST /api/v1/threat_center/count` | body 数组，`count_name=business_risk`、`aggregation_field=business.name.keyword`、`data_range=4`、`search_condition` | 输出只保留安全统计 bucket，不保留 raw request body。 |
| `tool.fobrain.threat_relevance_list` | `GET /api/v1/threat_center/relevance/list` | `page`、`per_page`、`keyword`、`vul_name`，可选 `ip`、`business_name` | `vulnerability_name` 需同时映射为 `keyword` 与 `vul_name`。 |

## 字段与安全边界

只允许把以下内容作为新项目 `StructuredResult` 的安全事实材料：

- 资产：`id`、`ip`、`hostname`、`status`、`network_type`、`ip_type`、安全摘要字段。
- 漏洞：`id`、`name`、`level`、`status`/`statusCode`、`person_info[].name`、`risk_num`、安全摘要字段。
- 样本报告：只记录样本是否存在、capability id、`result_ref`、失败分类和脱敏检查结果。

禁止进入 Product Facts、Workbench、Action API、audit、replay 或验收报告：

- `authorization`、token、cookie、凭据绑定真实值。
- raw provider payload、完整 ES hit、数据库字段全集。
- 真实本地配置文件内容。
- 未经 `StructuredResultSafetyGate` 审批的 provider 返回。

## Smoke 使用

先生成脱敏 discovery report，并按需写出本地样本文件：

```bash
go run ./scripts/fobrain_sample_discovery \
  --config configs/eino-workbench.local.yaml \
  --output test-results/eino-workbench-fobrain-sample-discovery-report.json \
  --samples-output test-results/eino-workbench-fobrain-batch-d-samples.local.json
```

`sample_discovery_report` 不记录真实 owner、department、IP；`*.local.json` 才包含可手工使用的样本值，只能保留在本地 ignored 目录，不得提交或粘贴到验收记录。

稳定样本确认后运行：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-d \
  --config configs/eino-workbench.local.yaml \
  --owner "<负责人>" \
  --department "<部门>" \
  --ip "<IP>"
```

通过条件：

- `test-results/eino-workbench-fobrain-batch-d-live-report.json` 的 `status` 为 `passed`。
- `sample_inputs.owner_present`、`department_present`、`ip_present` 都为 `true`。
- 六个 Batch D capability 都为 `passed`，且 `blocks_claims=[]`。
- 报告不包含样本值、token、auth header、credential ref、raw payload 或 `safe_summary`。

缺少任一样本时，报告必须是 `blocked`，不得声明 Batch D live pass 或 Fobrain 24 只读完成。
