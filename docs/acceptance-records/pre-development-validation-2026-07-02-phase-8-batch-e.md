# Pre-development Validation: Phase 8 Batch E

日期：2026-07-02

范围：Fobrain Batch E（资产详情、漏洞详情、业务风险摘要、威胁关联资产）开工前复核。

## 外部资料复核

| 主题 | 官方资料 | 访问日期 | 结论 |
| --- | --- | --- | --- |
| Eino ChatModelAgent / Runner | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/ | 2026-07-02 | ChatModelAgent 仍以 ChatModel + Tools 驱动 ReAct loop；Batch E 应继续通过 capability metadata/tool call，不新增关键词路由。 |
| Eino HITL | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/ | 2026-07-02 | Batch E 是只读工具，不新增 approval；not_found/clarification 仍必须映射到 Product Facts。 |
| Eino Checkpoint / Interrupt | https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/ | 2026-07-02 | 详情查询结果必须通过 StructuredResult 进入 checkpoint/replay，不直接暴露 provider payload。 |
| Eino Callback | https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/ | 2026-07-02 | callback 仅用于 tracing/metrics，不作为 Workbench SSE 主来源。 |
| MCP lifecycle/tools | https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle、https://modelcontextprotocol.io/specification/2025-06-18/server/tools | 2026-07-02 | lifecycle、tools/list、tools/call、structuredContent 和 output schema 约束与现有 provider contract 不冲突。 |
| OpenAPI / JSON Schema | https://spec.openapis.org/oas/v3.1.2.html、https://json-schema.org/draft/2020-12/json-schema-core | 2026-07-02 | 继续使用 OpenAPI 3.1 / JSON Schema 2020-12；本切片只扩展 schema，不改变 dialect。 |
| Playwright / React / Vite | https://playwright.dev/docs/test-snapshots、https://react.dev/learn/start-a-new-react-project、https://vite.dev/guide/ | 2026-07-02 | 后续视觉证据仍按 Playwright snapshot；前端技术栈不变。 |

## 版本复核

```bash
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go list -m -versions github.com/eino-contrib/jsonschema
```

结论：

- 当前项目 `github.com/cloudwego/eino` 为 `v0.9.12`。
- 可见最新稳定仍为 `v0.9.12`，存在 `v0.10.0-alpha.*`，不采用 alpha。
- `eino-ext/components/model/openai` 可见到 `v0.1.13`，本切片不新增依赖。
- `github.com/eino-contrib/jsonschema` 可见到 `v1.0.3`，本切片不新增依赖。

## 旧项目与 Fobrain 源码复核

| 能力 | 旧项目证据 | Fobrain 源码证据 | 结论 |
| --- | --- | --- | --- |
| `get_asset_detail` | 旧 adapter `GetAssetDetail` 使用 `fobrainAssetDetailPath(network_type, asset_id)`，未知类型默认 internal。 | `routes/asset_center/asset_center.go` 暴露 `/internal_asset/:id`、`/external_ip_asset/:id`、`/device/:id`、`/domain_asset/:id`。 | 新 schema 增加可选 `network_type`；实现必须集中归一化类型、路径分流和 fallback。 |
| `get_vulnerability_detail` | 旧 adapter 调用 `/api/v1/threat_center/:id`。 | `threat_center.Show` 调用 `threat.Show`，成功响应为单个 `data`。 | 按当前私有 `/api/threat_center/:id` 优先，fallback `/api/v1/threat_center/:id`。 |
| `business_risk_summary` | 旧 adapter POST `/api/v1/threat_center/count`，聚合 `business.name.keyword`。 | `ThreatCountRequest` 支持 `count_name`、`aggregation_field`、`search_condition`、`data_range`。 | 只输出安全 metrics，不保留 POST body。 |
| `threat_relevance_list` | 旧 adapter GET `/api/v1/threat_center/relevance/list`，同时传 `keyword` 和 `vul_name`。 | `VulRelevanceListRequest` 支持 `keyword`、`ip`、`business_name`、`vul_name`。 | `vulnerability_name` 必须映射到 `keyword` 与 `vul_name`。 |

## 新设计确认

- Workbench 和 Action API 仍必须共用 Product Facts。
- Batch E 工具结果仍以 StructuredResult 为唯一事实材料。
- 不新增后端关键词路由；工具选择仍由 Eino native tool call + capability metadata 驱动。
- 详情 ID 和业务名是输入材料，不得把本地真实样本写入验收文档。
- raw provider payload、完整 ES hit、POST body、token、auth header 不得进入 Product Facts、Workbench、Action API、audit、replay 或 report。
- `get_asset_detail` 当前最大设计风险是资产类型路径分流；已通过 `network_type` 可选输入和文档门禁消除静默硬编码风险。
- Batch E live 详情 ID 不得从可提交 Batch D report 推导；必须来自 ignored local samples 或新的脱敏 discovery sidecar。

## Go / No-Go

允许进入 Batch E mock/catalog 和 HTTP mapper 设计实现；不得直接声明 Batch E live pass 或 Fobrain 24 只读完成。进入实现前必须先按 `docs/fobrain-batch-e-interface-plan.md` 编写测试，尤其覆盖资产详情类型归一化和路径分流、业务风险 POST body 脱敏、威胁关联参数映射、live 样本来源脱敏和 empty/not_found 状态。
