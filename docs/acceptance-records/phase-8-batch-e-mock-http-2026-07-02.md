# Phase 8 Batch E Mock HTTP Acceptance Record

日期：2026-07-02

范围：Fobrain Batch E 四个详情与风险关联只读工具的 catalog、输入校验、mock StructuredResult 和 HTTP mapper。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| Batch E 四个工具进入可选 provider catalog | `docs/fobrain-batch-e-interface-plan.md` | `internal/einoapp/providers/fobrain/batch_e_test.go` | `go test ./internal/einoapp/providers/fobrain -run BatchE -count=1` | Pass | 只有 client 实现 `DetailRiskClient` 时注册 Batch E。 |
| Batch E schema、矩阵和 catalog 输入口径一致 | `docs/schemas/fobrain/tool_inputs.v1.schema.json`、`docs/fixtures/fobrain/tool-matrix-24.json` | schema validation、catalog tests | `npm run eino-workbench:schema-test`、`go test ./internal/einoapp/providers/fobrain -run BatchE -count=1` | Pass | `business_risk_summary` 使用专用 `business_risk_query`；`threat_relevance_list` 支持可选 IP/业务过滤；`network_type` 允许兼容输入别名并在 provider 边界归一。 |
| 输入 mapper 覆盖必填、schema mismatch、空态和安全边界 | `docs/schemas/fobrain/tool_inputs.v1.schema.json` | provider invoke/missing argument/schema mismatch/unsafe tests | `go test ./internal/einoapp/providers/fobrain -run BatchE -count=1` | Pass | 非字符串必填参数和 unsafe 材料不会调用 client。 |
| mock 输出只形成 StructuredResult candidate | `docs/facts-contract.md` | mock provider success/empty/error folding tests | `go test ./internal/einoapp/providers/fobrain -run BatchE -count=1` | Pass | 不把 raw provider payload、POST body 或凭据写入 Product Facts；provider unknown error 会脱敏折叠。 |
| 资产详情 HTTP mapper 路径按 `network_type` 分流 | `docs/fobrain-source-api-reference.md` | `TestHTTPClientAssetDetailUsesCanonicalNetworkTypePath` | `go test ./internal/einoapp/providers/fobrain -run AssetDetail -count=1` | Pass | 支持内网、外网、device、domain 和未知类型默认 internal。 |
| 漏洞详情 HTTP mapper 支持非分页响应和 fallback | `docs/fobrain-batch-e-interface-plan.md` | `TestHTTPClientVulnerabilityDetailFallsBackToLegacyPath` | `go test ./internal/einoapp/providers/fobrain -run VulnerabilityDetail -count=1` | Pass | 当前私有 `/api/threat_center/:id` 优先，`/api/v1/...` 兜底。 |
| 业务风险 HTTP mapper 使用 count 聚合接口 | `docs/fobrain-batch-e-interface-plan.md` | `TestHTTPClientBusinessRiskSummaryPostsCountAggregation` | `go test ./internal/einoapp/providers/fobrain -run BusinessRisk -count=1` | Pass | 输出只保留安全 metric bucket，raw body 不进入 StructuredResult。 |
| 威胁关联 HTTP mapper 同步 `keyword` 和 `vul_name` | `docs/fobrain-batch-e-interface-plan.md` | `TestHTTPClientThreatRelevanceMapsVulnerabilityNameToKeywordAndVulName` | `go test ./internal/einoapp/providers/fobrain -run ThreatRelevance -count=1` | Pass | 可选 IP 和业务名作为查询过滤，不进入工具硬编码分支。 |
| Batch E HTTP 错误路径脱敏 | `docs/06-security-and-projection.md` | auth failure、business code error、malformed JSON、timeout tests | `go test ./internal/einoapp/providers/fobrain -run BatchE -count=1` | Pass | 错误摘要不包含 token、auth header、raw body、URL query 或 POST body 字段。 |

## 验证命令

```bash
go test ./internal/einoapp/providers/fobrain -run 'BatchE|AssetDetail|VulnerabilityDetail|BusinessRisk|ThreatRelevance' -count=1
go test ./internal/einoapp/providers/fobrain -count=1
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
```

## 剩余风险

- 本记录不包含 Batch E sample discovery、live smoke/report schema、真实 ID 脱敏 sidecar 或 live pass 报告。
- 本记录不包含 Workbench 视觉、Action API 同源展示、replay/audit 证据。
- 不得据此声明 Batch E 完成、Fobrain 24 个只读工具恢复完成或产品能力可比。
