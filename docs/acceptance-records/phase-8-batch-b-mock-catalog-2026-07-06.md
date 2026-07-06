# Phase 8 Batch B Mock Catalog

阶段：Phase 8.1 24 个只读工具矩阵  
日期：2026-07-06  
方案：Batch B “我的范围”六个只读工具 mock/catalog 切片  
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
go test ./internal/einoapp/providers/fobrain -run 'BatchB|MyScope|ReadonlyToolMatrix' -count=1
go test ./internal/einoapp/providers/fobrain -count=1
go test ./cmd/eino-workbench -run 'Fobrain' -count=1
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
npm run eino-workbench:schema-test
go test ./...
npm run eino-workbench:contract-test
git diff --check
```

## 证据

- Batch B 六个工具进入 mock provider catalog：`my_assets`、`my_department_assets`、`my_vulnerabilities`、`my_department_vulnerabilities`、`my_business_systems`、`my_important_business_systems`。
- 输入只包含 `keyword/page/page_size`，不要求用户提供 owner 或 department。
- provider 在调用前设置当前用户范围、本部门范围或重要业务范围语义；StructuredResult 摘要只使用本地元数据和已校验 query。
- Mock client 可返回安全 `MyScopeResult`；raw user payload、Authorization、API key、credential ref 不进入 Product Facts。
- 本记录创建时，`ReadonlyToolMatrix` 已更新为当前 provider 广告 connector + Batch A/B/D/E，Batch C 仍是未实现缺口；后续 Batch C mock/catalog 已由 `phase-8-batch-c-mock-catalog-2026-07-06.md` 更新。

## 结论

- 通过 / 不通过 / skipped blocking：通过。
- 阻断 P0/P1/P2：本记录不单独声明 P0/P1/P2 通过。
- 阻断重构完成声明：仍阻断 Fobrain 24/24 完整恢复声明；Batch B/C live、视觉、Action API 同源 smoke 尚未完成。
- 允许替换当前产品基线：不允许。
- report schema 校验结果：`npm run eino-workbench:schema-test` 通过。
- 不得声明的能力：Batch B/C live pass、Batch B/C 视觉通过、Action API 同源通过、Fobrain 24/24 只读恢复完成、Fobrain 能力可比、替换当前产品基线。
- 关联 ADR：`docs/adr/2026-07-02-fobrain-live-auth-and-batch-gates.md`。
