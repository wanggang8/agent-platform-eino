# Phase 8 Batch C Live Mapper

阶段：Phase 8.1 24 个只读工具矩阵  
日期：2026-07-06  
方案：Batch C “直接列表与统计”HTTP live mapper 最小切片  
执行人：Codex

## 环境

- Go：go1.23.12 darwin/arm64
- Node：v25.8.1
- npm：11.11.0
- OS：macOS 26.5.2 25F84
- 模型 provider：不使用
- Fobrain 环境：不使用，使用 `httptest`

## 命令

```bash
go test ./internal/einoapp/providers/fobrain -run 'DirectRead|BatchC' -count=1
go test ./internal/einoapp/providers/fobrain -count=1
go test ./cmd/eino-workbench -run 'Fobrain' -count=1
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
npm run eino-workbench:schema-test
go test ./...
npm run eino-workbench:contract-test
git diff --check
```

## 证据

- `HTTPClient` 已实现 `DirectReadClient`，因此 live 模式可执行 Batch C provider 调用链。
- `business_list` 映射到 `/api/business`，并兼容 `/api/v1/business`。
- `external_high_risk_assets` 映射到 `/api/external_ip_asset`，并兼容 `/api/v1/external_ip_asset`。
- `vulnerability_status_summary` 映射到 `/api/threat_center/count`，只保留安全指标。
- `pending_tickets` 映射到当前真实路径 `/api/ticket/pending`，保留 `/api/v1/ticket/pending` 作为旧环境 fallback；参数保留 `keyword`、`status[]`、`page/page_size`，不发送人员过滤参数，只保留安全工单行摘要。
- `ip_stats` / `vul_stats` 映射到 `/api/threat_center/relevance/ip_stats` 和 `/api/threat_center/relevance/vul_stats`，只保留安全指标。
- 业务错误、非法响应和认证错误只折叠为安全错误，不把 token、header、raw response 或 raw aggregation body 写入事实材料。

## 结论

- 通过 / 不通过 / skipped blocking：通过。
- 阻断 P0/P1/P2：本记录不单独声明 P0/P1/P2 通过。
- 阻断重构完成声明：仍阻断 Fobrain 24/24 完整恢复声明；Batch C 真实 live smoke、视觉、Action API 同源 smoke 尚未完成。
- 允许替换当前产品基线：不允许。
- report schema 校验结果：`npm run eino-workbench:schema-test` 通过。
- 不得声明的能力：Batch C live pass、Batch C 视觉通过、Action API 同源通过、Fobrain 24/24 只读恢复完成、Fobrain 能力可比、替换当前产品基线。
- 关联 ADR：`docs/adr/2026-07-02-fobrain-live-auth-and-batch-gates.md`。
