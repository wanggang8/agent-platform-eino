# Phase 8.4 Fobrain Clarification Acceptance Record

日期：2026-07-02

范围：Fobrain 实体消歧安全候选、clarification PendingInteraction 同源投影、Workbench 候选渲染。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| 多候选实体进入 clarification | `docs/clarification-flow.md`、`docs/fobrain-live-read-batch-plan.md` | `internal/einoapp/providers/fobrain/disambiguation_test.go` | `go test ./internal/einoapp/providers/fobrain -run Disambiguation -count=1` | Pass | 本阶段只实现候选分类 helper，不声明 24 个真实工具全部接入。 |
| 候选不泄露 raw/provider/凭据/手机号/邮箱 | `docs/facts-contract.md` | `facts.UnsafePendingCandidates` 和 provider 测试 | `go test ./internal/einoapp/providers/fobrain -run Disambiguation -count=1` | Pass | raw provider payload 仍只允许在 provider 边界内处理。 |
| PendingInteraction 保存 input_mode/candidates | `docs/facts-contract.md` | facts model、memory repository、SQLite repository | `go test ./...` | Pass | SQLite 以 JSON 保存安全候选，并覆盖 snapshot/resume round-trip；已有本地库会补列。 |
| Workbench、ActionResult、SSE 同源投影候选 | `docs/05-contract-design.md` | `internal/einoapp/product/projection_test.go` | `go test ./internal/einoapp/product -run ClarificationCandidates -count=1` | Pass | 候选来自同一 PendingInteraction。 |
| 前端 ClarificationCard 不写死候选 | `docs/02-ux-visual-requirements.md` | `web/eino-workbench/src/features/workbench/components/cards.tsx` | `npm run eino-workbench:test`、`npm run eino-workbench:typecheck` | Pass | 真实提交动作仍属于后续 resume 阶段。 |
| schema/fixture/contract 同步 | `docs/05-contract-design.md` | Workbench View schema、ActionResult fixture、generated TS | `npm run eino-workbench:contract-test` | Pass | schema 数量保持 32。 |

## 审查修复

- 子 agent 审查发现候选安全门禁、SSE pending `run_id` 和 `input_mode` 枚举校验缺口。
- 已补充格式化手机号、raw provider id、中文敏感字段、schema 外 entity/input mode 的拒绝规则。
- 已补充 SSE pending patch `run_id` 和 SQLite candidates/input_mode round-trip 回归测试。

## 验证命令

```bash
go test ./...
npm run eino-workbench:typecheck
npm run eino-workbench:test
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-clarification
```

## 剩余风险

- 真实 Fobrain Batch D 查询工具尚未把 resolver 接入 HTTP client；后续必须按 `docs/fobrain-tool-matrix.md` 逐工具恢复。
- clarification submit/cancel/restart 后 resume 状态机仍属于 Phase 6/8 后续任务，本记录不声明完成。
