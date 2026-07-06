# Phase 8 Batch C Mock Catalog

阶段：Phase 8.1 24 个只读工具矩阵  
日期：2026-07-06  
方案：Batch C “直接列表与统计”六个只读工具 mock/catalog 切片  
执行人：Codex

## 环境

- Go：go1.23.12 darwin/arm64
- Node：v25.8.1
- npm：11.11.0
- OS：macOS 26.5.2 25F84
- 模型 provider：不使用
- Fobrain 环境：不使用

## 命令

```bash
go test ./internal/einoapp/providers/fobrain -run 'BatchC|DirectRead|ReadonlyToolMatrix' -count=1
go test ./internal/einoapp/capabilities -run 'ToolAdapter' -count=1
go test ./internal/einoapp/providers/fobrain -count=1
go test ./cmd/eino-workbench -run 'Fobrain' -count=1
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
npm run eino-workbench:schema-test
go test ./...
npm run eino-workbench:contract-test
git diff --check
```

## 证据

- Batch C 六个工具进入 mock provider catalog：`business_list`、`external_high_risk_assets`、`vulnerability_status_summary`、`pending_tickets`、`ip_stats`、`vul_stats`。
- 输入只包含业务名、负责人、关键字、字段、严重级别、状态、时间范围、人员和分页等安全筛选字段。
- Batch C 显式空 `required` 会在 Eino ToolInfo 中保持全部可选，支持无筛选首读。
- provider 通过可选 `DirectReadClient` 暴露能力；未实现该接口的 client 不会广告 Batch C 工具。
- StructuredResult 摘要只使用本地元数据和已校验 query；client 返回的 raw title、raw query、Authorization、API key、credential ref 不进入 Product Facts。
- `ReadonlyToolMatrix` 当前对账为 connector + 24 个只读 mock catalog；该结论只覆盖 mock/catalog，不覆盖 live、视觉或 Action API 同源。

## 结论

- 通过 / 不通过 / skipped blocking：通过。
- 阻断 P0/P1/P2：本记录不单独声明 P0/P1/P2 通过。
- 阻断重构完成声明：仍阻断 Fobrain 24/24 完整恢复声明；Batch B/C live、视觉、Action API 同源 smoke 尚未完成。
- 允许替换当前产品基线：不允许。
- report schema 校验结果：`npm run eino-workbench:schema-test` 通过。
- 不得声明的能力：Batch C live pass、Batch C 视觉通过、Action API 同源通过、Fobrain 24/24 只读恢复完成、Fobrain 能力可比、替换当前产品基线。
- 关联 ADR：`docs/adr/2026-07-02-fobrain-live-auth-and-batch-gates.md`。
