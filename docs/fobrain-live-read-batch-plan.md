# Fobrain Live Read 分批计划

本文定义 Phase 8 进入 24 个 Fobrain 只读工具恢复前的分批节奏。目标是小批量恢复、逐批验收，避免一次性实现导致事实、展示或安全边界丢失。

## 认证与凭据

- 认证参数名来自 `fobrain.credential.auth_param`，当前为 `authorization`。
- `auth_param` 必须在配置文件中显式声明；provider、credential resolver 和 HTTP client 不得自动补默认值。
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

当前代码进展：Batch B 六个“我的范围”只读工具已进入 mock provider catalog，并具备 `keyword/page/page_size` 安全入参、当前用户/本部门/重要业务范围语义标记、mock client 和 StructuredResult 安全摘要 mapper。HTTP live mapper 已完成 focused tests，覆盖当前用户/部门派生、`/api` 到 `/api/v1` fallback、业务系统 `person_base.name` owner scope + `business_name` 收窄过滤、重要性 `assets_attribute.important_types` 过滤和错误脱敏。`fobrain-batch-b` live smoke/report 基础设施已补齐，报告通过 `docs/schemas/fobrain/batch_b_live_report.v1.schema.json` 约束；无真实配置时只生成 skip report，真实返回任一空结果时生成 blocking report。2026-07-06 真实配置下，Batch B live smoke 因当前用户上下文缺少可解析部门且“我的范围”查询为空而 blocked。Workbench 视觉证据和 Action API 同源 smoke 仍需后续任务补齐；未生成真实 passed 报告前，不得声明 Batch B live pass 或 Fobrain 24/24 恢复完成。

### Batch C：直接列表与统计

范围：

- `tool.fobrain.business_list`
- `tool.fobrain.external_high_risk_assets`
- `tool.fobrain.vulnerability_status_summary`
- `tool.fobrain.ip_stats`
- `tool.fobrain.vul_stats`
- `tool.fobrain.pending_tickets`

验收重点：无筛选首读、指标卡、表格卡、状态标签安全映射、pending ticket 只读展示。

当前代码进展：Batch C 六个“直接列表与统计”只读工具已进入 mock provider catalog，并具备安全筛选入参、无筛选默认范围、列表/指标两类 mock client 返回和 StructuredResult 安全摘要 mapper。HTTP live mapper 已按当前真实 Fobrain 路由补齐 focused tests：业务系统列表走 `/api/business`，外部高风险资产走 `/api/external_ip_asset`，漏洞状态汇总走 `/api/threat_center/count`，待处理工单走 `/api/ticket/pending` 且只发送 `keyword`、`status[]`、`page/page_size`，不暴露人员过滤语义，IP/漏洞统计走 `/api/threat_center/relevance/ip_stats` 和 `/api/threat_center/relevance/vul_stats`。`fobrain-batch-c` live smoke/report 基础设施已补齐，报告通过 `docs/schemas/fobrain/batch_c_live_report.v1.schema.json` 约束；无真实配置时只生成 skip report，真实返回任一空结果时生成 blocking report。2026-07-06 真实配置确认 `pending_tickets` 当前真实环境没有待处理数据，已按用户确认在 Batch C live smoke 中标记为 `skipped`，不再阻断 Batch C 其他 5 个工具的 live pass；但工单读取能力暂缓，仍不得据此声明 Fobrain 24/24 恢复完成。Workbench 视觉证据和 Action API 同源 smoke 仍需后续任务补齐。

### Batch D：参数化查询

范围：

- `tool.fobrain.list_assets_by_owner`
- `tool.fobrain.list_vulnerabilities_by_owner`
- `tool.fobrain.list_assets_by_department`
- `tool.fobrain.list_vulnerabilities_by_department`
- `tool.fobrain.list_assets_by_ip`
- `tool.fobrain.list_vulnerabilities_by_ip`

验收重点：输入 schema、实体消歧、IP 完整展示、分页和 clarification 分支。

当前代码进展：Batch D 六个参数化只读工具已进入 mock/参数化 client 可执行 catalog，并具备 mock 参数解析、必填字段校验和 `StructuredResult` 安全摘要 mapper。HTTP live mapper 已按当前真实接口返回重新归一化：资产优先 `/api/asset`，漏洞优先 `/api/threat_center`，并兼容旧 `/api/v1/...` fallback；真实返回包装为 `{code,data,message}`，列表位于 `data.items`。Batch D 专用 live report 已覆盖六项工具、脱敏和非空 `item_count` 门禁；实体消歧接入、真实分页批量验收和 Workbench 视觉证据仍需后续批次验收。

Batch D live smoke/report 命令：

```bash
go run ./scripts/fobrain_sample_discovery --config configs/eino-workbench.local.yaml --output test-results/eino-workbench-fobrain-sample-discovery-report.json --samples-output test-results/eino-workbench-fobrain-batch-d-samples.local.json
```

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-d --config configs/eino-workbench.local.yaml --owner "<负责人>" --department "<部门>" --ip "<IP>"
```

未提供稳定 `owner`、`department`、`ip` 样本，或任一稳定样本查询返回空结果时，Batch D smoke 必须生成 `blocked` report，并通过 `blocks_claims` 阻止 Batch D live pass 声明。样本发现流程、源码接口证据和 discovery report 规则见 `docs/fobrain-source-api-reference.md`。当前本地 discovery report 已找到 owner/department/IP 三类稳定样本，`test-results/eino-workbench-fobrain-batch-d-live-report.json` 的六个 Batch D capability 均为 `passed`，且每项 `item_count > 0`。

### Batch E：详情与风险关联

范围：

- `tool.fobrain.get_asset_detail`
- `tool.fobrain.get_vulnerability_detail`
- `tool.fobrain.business_risk_summary`
- `tool.fobrain.threat_relevance_list`

验收重点：详情页事实分区、风险摘要、威胁关联表格、证据引用和 replay 一致性。

Batch E 开发前必须先完成 `docs/fobrain-batch-e-interface-plan.md` 的接口矩阵和样本策略。资产详情不能只按 `asset_id` 猜测路径：`network_type` 可选，缺省按旧 adapter 走 internal fallback；实现必须集中归一化 `1/2`、中英文内外网和 device/domain alias。详情 ID 只能来自 ignored local samples 或新的脱敏 discovery sidecar，不能从可提交 Batch D live report 反推，也不能回读 raw provider payload。Batch E live report 必须区分 resolved、empty/not_found、provider failure，并且不得把真实详情 ID、raw payload、POST body、auth header 或本地配置写入可提交报告。

当前代码进展：Batch E 四个工具已进入可选 `DetailRiskClient` catalog 和 provider 调用链，mock client 可返回 StructuredResult 候选，HTTP client 已实现当前私有 `/api/...` 优先、`/api/v1/...` fallback 的详情和聚合 mapper。已覆盖的安全边界包括 `network_type` 中央归一、必填参数校验、敏感字段拒绝、业务风险 raw POST body 不进入 StructuredResult、威胁关联参数映射。样本发现已扩展 Batch E 本地样本字段，并从 `/api/v1/business`/`/api/business` 补充业务系统样本；安全 discovery report 只记录 presence booleans。当前真实 discovery 已找到 owner/department/IP、资产详情、漏洞详情、业务系统和威胁名样本；`fobrain-batch-e` live smoke 已生成 passed 报告，四个工具均从 provider Invoke 输出 StructuredResult。Workbench 视觉、fresh replay 和 audit evidence 已生成 fixture-based desktop/mobile 六区域 baseline，记录见 `docs/acceptance-records/phase-8-batch-e-visual-replay-audit-2026-07-03.md`。该结论不覆盖 Batch B-D 视觉、Action API 真实服务同源 smoke、实体消歧或写域审批，因此不得声明 Fobrain 24 只读最终验收完成。

Batch E discovery 与 smoke 命令：

```bash
go run ./scripts/fobrain_sample_discovery \
  --config configs/eino-workbench.local.yaml \
  --output test-results/eino-workbench-fobrain-sample-discovery-report.json \
  --samples-output test-results/eino-workbench-fobrain-batch-e-samples.local.json
```

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-e \
  --config configs/eino-workbench.local.yaml \
  --asset-id "<资产ID>" \
  --asset-network-type "<internal|external|device|domain>" \
  --vulnerability-id "<漏洞ID>" \
  --business-name "<业务系统>" \
  --vulnerability-name "<漏洞名>"
```

上述样本值只能来自 ignored local samples 或人工只读确认，不得写入可提交报告。

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
