# Phase 8 Batch D Live Report Acceptance Record

日期：2026-07-02

范围：Fobrain Batch D 六个参数化只读工具的 live smoke/report schema 与脚本门禁。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| Batch D live report schema 固定覆盖六个工具 | `docs/tooling-and-reporting.md`、`docs/fobrain-live-read-batch-plan.md` | `docs/schemas/fobrain/batch_d_live_report.v1.schema.json` | `npm run eino-workbench:schema-test` | Pass | schema 允许 `passed/failed/blocked`，但固定 6 条 capability result；passed 时每项 `item_count >= 1`。 |
| smoke 支持 owner/department/IP 样本并脱敏落盘 | `docs/06-security-and-projection.md` | `scripts/fobrain_batch_d_smoke/main.go`、`main_test.go` | `go test ./scripts/fobrain_batch_d_smoke -count=1` | Pass | 报告只记录样本是否存在、StructuredResult 出口和非敏感 `item_count`，不记录 token/raw payload。 |
| shell 场景接入 `fobrain-batch-d` | `docs/tooling-and-reporting.md` | `scripts/eino_workbench_server_smoke.sh` | `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-d --config configs/eino-workbench.local.yaml --owner <local> --department <local> --ip <local>` | Pass | 当前本地真实配置输出 `passed`。样本值来自 ignored local samples 文件，未写入本文档。 |
| 稳定样本通过六项 live capability | `docs/fobrain-live-read-batch-plan.md` | `test-results/eino-workbench-fobrain-batch-d-live-report.json` | report schema validation | Pass | owner/department/IP 六项均为 passed，且每项 `item_count > 0`，`blocks_claims=[]`。 |
| 敏感信息不进入报告 | `docs/06-security-and-projection.md` | report grep | `rg "api_token|credential_ref|bearer|raw provider|raw body|workspace-token|Authorization|authorization" test-results/eino-workbench-fobrain-batch-d-live-report.json` | Pass | grep 无命中。 |

## 当前报告状态

- `status`: `passed`
- `failure_category`: `none`
- `owner_present`: `true`
- `department_present`: `true`
- `ip_present`: `true`
- `capability_results`: 6 passed, each `item_count > 0`
- `blocks_claims`: `[]`

## 后续通过条件

重新运行时提供稳定样本：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-d --config configs/eino-workbench.local.yaml --owner "<负责人>" --department "<部门>" --ip "<IP>"
```

只有报告 `status=passed` 且 `blocks_claims=[]` 时，才能声明 Batch D live pass。
