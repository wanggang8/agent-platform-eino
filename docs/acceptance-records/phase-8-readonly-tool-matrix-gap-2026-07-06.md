# Phase 8 Readonly Tool Matrix Gap Guard

阶段：Phase 8.1 24 个只读工具矩阵  
日期：2026-07-06  
方案：新增机器化矩阵与 provider catalog 对账测试  
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
go test ./internal/einoapp/providers/fobrain -run ReadonlyToolMatrix -count=1
npm run eino-workbench:schema-test
go test ./internal/einoapp/providers/fobrain -count=1
go test ./cmd/eino-workbench -run 'Fobrain' -count=1
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
go test ./...
npm run eino-workbench:contract-test
git diff --check
```

## 证据

- `docs/fixtures/fobrain/tool-matrix-24.json` 精确登记 24 个只读工具。
- Batch 分布为 A=2、B=6、C=6、D=6、E=4。
- 本记录创建时，provider catalog 只广告 `connector.fobrain.security`、Batch A、Batch D 和 Batch E。
- 本记录创建时，Batch B/C 十二个工具仍在矩阵中登记，但不得进入 provider catalog 或声明可用；后续 Batch B mock/catalog 已由 `phase-8-batch-b-mock-catalog-2026-07-06.md` 更新，Batch C mock/catalog 已由 `phase-8-batch-c-mock-catalog-2026-07-06.md` 更新。
- cmd 启动路径和 bootstrap Fobrain 配置回归测试通过。

## 结论

- 通过 / 不通过 / skipped blocking：通过。
- 阻断 P0/P1/P2：本记录不单独声明 P0/P1/P2 通过。
- 后续更新：`docs/acceptance-records/phase-8-batch-b-mock-catalog-2026-07-06.md` 已补 Batch B mock/catalog，`docs/acceptance-records/phase-8-batch-c-mock-catalog-2026-07-06.md` 已补 Batch C mock/catalog；当前 mock provider catalog 无只读工具缺口。
- 阻断重构完成声明：仍阻断 Fobrain 24/24 完整恢复声明；Batch B/C 视觉、live、Action API 同源 smoke 尚未完成。
- 允许替换当前产品基线：不允许。
- report schema 校验结果：`npm run eino-workbench:schema-test` 通过。
- 不得声明的能力：Batch B/C live 通过、Batch B/C 视觉通过、Action API 同源通过、Fobrain 24/24 只读恢复完成、Fobrain 能力可比、替换当前产品基线。
- 关联 ADR：`docs/adr/2026-07-02-fobrain-live-auth-and-batch-gates.md`。
