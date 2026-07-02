# Fobrain Live Read 分批计划

本文定义 Phase 8 进入 24 个 Fobrain 只读工具恢复前的分批节奏。目标是小批量恢复、逐批验收，避免一次性实现导致事实、展示或安全边界丢失。

## 认证与凭据

- 认证参数名来自 `fobrain.credential.auth_param`，当前为 `authorization`。
- token 来自 ignored 本地配置 `configs/eino-workbench.local.yaml` 的 `fobrain.credential.api_token`。
- 当前只支持 workspace 共享 token；不做单用户 token、单工具 token 或前端传 token。
- live client 只在 provider 边界持有 token；所有产品出口只能展示安全摘要。
- live `base_url` 默认使用 HTTPS；本地验收代理可以使用 loopback HTTP。

## 批次顺序

### Batch A：连接器与当前用户

范围：

- `connector.fobrain.security`
- `tool.fobrain.current_user_context`
- `tool.fobrain.my_permissions`

验收重点：认证成功/失败、connector 状态、workspace scope、当前用户安全摘要、权限空态和错误脱敏。

当前代码基础已覆盖 Batch A 三项能力：`connector.fobrain.security` 从安全配置和凭据摘要生成状态，`current_user_context` 与 `my_permissions` 从当前用户接口读取安全字段。标准路径优先使用 `/api/v1/user`，404/405 时兼容私有部署 `/api/user`；权限字段缺失时 `my_permissions` 输出权限空态安全摘要。Batch A live/smoke 验收命令为：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-a --config configs/eino-workbench.local.yaml
```

配置齐全且真实环境可用时，该命令必须生成 `test-results/eino-workbench-fobrain-batch-a-live-report.json`，并符合 `docs/schemas/fobrain/batch_a_live_report.v1.schema.json`。配置缺失时只能生成 blocking skip report，不能声明 Batch A live pass。Batch A 最终完成声明仍必须补齐 Workbench 视觉证据和验收记录。

视觉证据进展：`connector.fobrain.security`、`tool.fobrain.current_user_context` 和 `tool.fobrain.my_permissions` 已生成新项目 desktop/mobile Playwright baseline，每个能力覆盖 `main-chat`、`fresh-main-chat`、`process`、`evidence`、`audit`、`internal-details` 六区域。该结论只覆盖 Batch A 视觉证据，不代表 Batch B-E 或 24 个只读工具全量恢复完成。

### Batch B：我的范围

范围：

- `tool.fobrain.my_assets`
- `tool.fobrain.my_department_assets`
- `tool.fobrain.my_vulnerabilities`
- `tool.fobrain.my_department_vulnerabilities`
- `tool.fobrain.my_business_systems`
- `tool.fobrain.my_important_business_systems`

验收重点：当前用户上下文复用、无需追问部门/负责人、分页、空态、表格 StructuredResult 和敏感字段脱敏。

### Batch C：直接列表与统计

范围：

- `tool.fobrain.business_list`
- `tool.fobrain.external_high_risk_assets`
- `tool.fobrain.vulnerability_status_summary`
- `tool.fobrain.ip_stats`
- `tool.fobrain.vul_stats`
- `tool.fobrain.pending_tickets`

验收重点：无筛选首读、指标卡、表格卡、状态标签安全映射、pending ticket 只读展示。

### Batch D：参数化查询

范围：

- `tool.fobrain.list_assets_by_owner`
- `tool.fobrain.list_vulnerabilities_by_owner`
- `tool.fobrain.list_assets_by_department`
- `tool.fobrain.list_vulnerabilities_by_department`
- `tool.fobrain.list_assets_by_ip`
- `tool.fobrain.list_vulnerabilities_by_ip`

验收重点：输入 schema、实体消歧、IP 完整展示、分页和 clarification 分支。

当前代码进展：Batch D 六个参数化只读工具已进入 mock/参数化 client 可执行 catalog，并具备 mock 参数解析、必填字段校验和 `StructuredResult` 安全摘要 mapper。HTTP live mapper 已按当前真实接口返回重新归一化：资产优先 `/api/asset`，漏洞优先 `/api/threat_center`，并兼容旧 `/api/v1/...` fallback；真实返回包装为 `{code,data,message}`，列表位于 `data.items`。该切片只证明 capability metadata、输入契约、HTTP mapper 和 Product Facts 结果边界；Batch D 专用 live report、真实分页批量验收、实体消歧接入和视觉证据仍需后续批次验收。

Batch D live smoke/report 命令：

```bash
go run ./scripts/fobrain_sample_discovery --config configs/eino-workbench.local.yaml --output test-results/eino-workbench-fobrain-sample-discovery-report.json --samples-output test-results/eino-workbench-fobrain-batch-d-samples.local.json
```

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-d --config configs/eino-workbench.local.yaml --owner "<负责人>" --department "<部门>" --ip "<IP>"
```

未提供稳定 `owner`、`department`、`ip` 样本时，Batch D smoke 必须生成 `blocked` report，并通过 `blocks_claims` 阻止 Batch D live pass 声明。样本发现流程、源码接口证据和 discovery report 规则见 `docs/fobrain-source-api-reference.md`。当前本地 discovery report 显示 owner/department 已找到，IP 未找到可同时命中资产和漏洞的样本，因此 `test-results/eino-workbench-fobrain-batch-d-live-report.json` 仍是 blocking report。

### Batch E：详情与风险关联

范围：

- `tool.fobrain.get_asset_detail`
- `tool.fobrain.get_vulnerability_detail`
- `tool.fobrain.business_risk_summary`
- `tool.fobrain.threat_relevance_list`

验收重点：详情页事实分区、风险摘要、威胁关联表格、证据引用和 replay 一致性。

## 每批完成标准

每批必须同时满足：

- capability catalog、input schema、result schema、fixture 和 StructuredResult mapper 已补齐。
- `docs/fixtures/fobrain/tool-matrix-24.json` 中对应工具的 `batch_gate` 与本计划一致。
- mock test 覆盖成功、空态、provider 错误、schema mismatch 和敏感字段拒绝。
- live 环境可用时生成对应批次 live pass report；不可用时只能生成 blocking skip report，且 `blocks_claims` 阻止通过声明。
- Workbench 和 Action API 消费同一组 Product Facts。
- replay、audit、SSE 和报告不包含 token、raw provider payload 或 credential ref。
- 前端可用后补齐工具卡折叠/展开、证据、审计和 fresh replay 截图。

## 最终门禁

只有 24 个只读工具和 `connector.fobrain.security` 全部通过 mock、live pass、视觉和脱敏门禁后，才能声明 Fobrain 只读能力恢复。blocking skip report 只能说明环境不可验收，不能作为最终通过依据。写域 `tool.fobrain.update_ticket_status` 仍属于审批门禁，不计入只读 24 工具完成声明。
