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

当前代码基础已覆盖 Batch A 三项能力：`connector.fobrain.security` 从安全配置和凭据摘要生成状态，`current_user_context` 与 `my_permissions` 从 `/api/v1/user` 读取安全字段。Batch A 最终完成声明仍必须补齐 live pass report、Workbench 视觉证据和验收记录。

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
- live 环境可用时生成 live pass report；不可用时只能生成 blocking skip report，且 `blocks_claims` 阻止通过声明。
- Workbench 和 Action API 消费同一组 Product Facts。
- replay、audit、SSE 和报告不包含 token、raw provider payload 或 credential ref。
- 前端可用后补齐工具卡折叠/展开、证据、审计和 fresh replay 截图。

## 最终门禁

只有 24 个只读工具和 `connector.fobrain.security` 全部通过 mock、live pass、视觉和脱敏门禁后，才能声明 Fobrain 只读能力恢复。blocking skip report 只能说明环境不可验收，不能作为最终通过依据。写域 `tool.fobrain.update_ticket_status` 仍属于审批门禁，不计入只读 24 工具完成声明。
