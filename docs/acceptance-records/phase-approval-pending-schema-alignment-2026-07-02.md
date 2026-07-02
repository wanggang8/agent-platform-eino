# Approval Pending Schema Alignment Acceptance Record

日期：2026-07-02

范围：补齐 approval PendingInteraction 的 `operation_name` / `target_summary` Product Facts、SQLite 持久化和 SSE pending patch。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| approval pending facts 保存操作名和目标摘要 | `docs/approval-flow.md`、`docs/facts-contract.md` | `facts.PendingInteraction`、SQLite repository | `go test ./internal/einoapp/store/sqlite -count=1` | Pass | 只补字段同源，不实现完整 approval resume 流程。 |
| SSE pending patch 符合 pending schema 必填字段 | `docs/schemas/eino_workbench_pending_interaction.v1.schema.json` | `internal/einoapp/product/projection_test.go` | `go test ./internal/einoapp/product -count=1` | Pass | ActionResult 不输出 schema 外 `operation_name`。 |
| 既有 clarification/Fobrain 契约不回退 | `docs/clarification-flow.md` | Phase 8.4 smoke | `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-clarification` | Pass | 保持 candidates/input_mode 投影。 |

## 验证命令

```bash
go test ./...
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
npm run eino-workbench:typecheck
npm run eino-workbench:test
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-clarification
```
