# Phase 8 Batch C Live Smoke Infra

阶段：Phase 8.1 24 个只读工具矩阵  
日期：2026-07-06  
方案：Batch C “直接列表与统计”live smoke/report 基础设施  
执行人：Codex

## 环境

- Go：go1.23.12 darwin/arm64
- Node：v25.8.1
- npm：11.11.0
- OS：macOS 26.5.2 25F84
- 模型 provider：不使用
- Fobrain 环境：不使用真实环境；单元测试使用 `httptest`

## 命令

```bash
go test ./scripts/fobrain_batch_c_smoke -count=1
go test ./internal/einoapp/providers/fobrain -run 'DirectRead|BatchC' -count=1
go test ./internal/einoapp/providers/fobrain -count=1
go test ./cmd/eino-workbench -run 'Fobrain' -count=1
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
npm run eino-workbench:schema-test
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-c --config /tmp/eino-missing-fobrain-live.yaml
go test ./...
npm run eino-workbench:contract-test
git diff --check
```

## 证据

- 新增 `scripts/fobrain_batch_c_smoke`，覆盖 `business_list`、`external_high_risk_assets`、`vulnerability_status_summary`、`pending_tickets`、`ip_stats`、`vul_stats` 六个 Batch C 工具。
- 新增 `docs/schemas/fobrain/batch_c_live_report.v1.schema.json`，要求 passed 报告中六个工具全部为 `resolved`、`allowed` 且 `item_count >= 1`。
- `scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-c` 已接入本地配置检查、skip report、live report 写入和 schema 校验。
- 报告只记录样本是否存在、StructuredResult schema/result_ref、item_count、policy decision 和稳定失败分类；不记录 token、authorization、raw provider payload、safe summary 或真实筛选值。
- 任一工具真实返回空结果时，runner 生成 `blocked` 报告并保留 `blocks_claims`，不能声明 Batch C live pass。

## 结论

- 通过 / 不通过 / skipped blocking：通过基础设施测试；真实 live smoke 未执行。
- 阻断 P0/P1/P2：本记录不单独声明 P0/P1/P2 通过。
- 阻断重构完成声明：仍阻断 Fobrain 24/24 完整恢复声明；Batch C 真实 live passed 报告、视觉、Action API 同源 smoke 尚未完成。
- 允许替换当前产品基线：不允许。
- report schema 校验结果：`npm run eino-workbench:schema-test` 通过。
- 不得声明的能力：Batch C live pass、Batch C 视觉通过、Action API 同源通过、Fobrain 24/24 只读恢复完成、Fobrain 能力可比、替换当前产品基线。
- 关联 ADR：`docs/adr/2026-07-02-fobrain-live-auth-and-batch-gates.md`。
