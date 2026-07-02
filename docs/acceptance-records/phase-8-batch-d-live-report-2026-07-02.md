# Phase 8 Batch D Live Report Acceptance Record

日期：2026-07-02

范围：Fobrain Batch D 六个参数化只读工具的 live smoke/report schema 与脚本门禁。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| Batch D live report schema 固定覆盖六个工具 | `docs/tooling-and-reporting.md`、`docs/fobrain-live-read-batch-plan.md` | `docs/schemas/fobrain/batch_d_live_report.v1.schema.json` | `npm run eino-workbench:schema-test` | Pass | schema 允许 `passed/failed/blocked`，但固定 6 条 capability result。 |
| smoke 支持 owner/department/IP 样本并脱敏落盘 | `docs/06-security-and-projection.md` | `scripts/fobrain_batch_d_smoke/main.go`、`main_test.go` | `go test ./scripts/fobrain_batch_d_smoke -count=1` | Pass | 报告只记录样本是否存在和 StructuredResult 出口，不记录 token/raw payload。 |
| shell 场景接入 `fobrain-batch-d` | `docs/tooling-and-reporting.md` | `scripts/eino_workbench_server_smoke.sh` | `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-d --config configs/eino-workbench.local.yaml` | Pass | 当前本地真实配置输出 `blocked`，不是 live pass。 |
| 缺样本不伪造通过 | `docs/fobrain-live-read-batch-plan.md` | `test-results/eino-workbench-fobrain-batch-d-live-report.json` | report schema validation | Pass | 当前 owner 两项 passed；department/IP 四项 blocked，`blocks_claims` 阻止最终声明。 |
| 敏感信息不进入报告 | `docs/06-security-and-projection.md` | report grep | `rg "api_token|credential_ref|bearer|raw provider|raw body|workspace-token|Authorization" test-results/eino-workbench-fobrain-batch-d-live-report.json` | Pass | grep 无命中。 |

## 当前报告状态

- `status`: `blocked`
- `failure_category`: `missing_sample_input`
- `owner_present`: `true`
- `department_present`: `false`
- `ip_present`: `false`
- `blocks_claims`: `fobrain-batch-d live pass`、`Fobrain 24 readonly final acceptance`

## 后续通过条件

重新运行时提供稳定样本：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-d --config configs/eino-workbench.local.yaml --owner "<负责人>" --department "<部门>" --ip "<IP>"
```

只有报告 `status=passed` 且 `blocks_claims=[]` 时，才能声明 Batch D live pass。
