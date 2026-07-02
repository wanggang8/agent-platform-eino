# Phase 8 Fobrain Sample Discovery Acceptance Record

日期：2026-07-02

范围：Fobrain Batch D 稳定样本发现脚本、脱敏 report schema 和源码接口参考文档。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| discovery report schema 可校验 | `docs/fobrain-source-api-reference.md` | `docs/schemas/fobrain/sample_discovery_report.v1.schema.json` | `npm run eino-workbench:schema-test` | Pass | schema 允许 blocked report 不带六项验证；passed 时强制六项验证非空。 |
| 脚本默认不记录真实样本值 | `docs/06-security-and-projection.md` | `scripts/fobrain_sample_discovery/main.go`、`main_test.go` | `go test ./scripts/fobrain_sample_discovery -count=1` | Pass | report 禁止 token、raw payload、`safe_summary` 和样本值；本地 samples 文件仅在显式 `--samples-output` 时写出。 |
| 源码接口参考已落文档 | `../fobrain` 只读源码 | `docs/fobrain-source-api-reference.md` | 文档检查 | Pass | 覆盖 `/api/v1/user`、`staff_list`、`personnel_departments`、`asset`、`threat_center` 和当前私有兼容路径。 |
| 当前真实配置可生成脱敏 report | `configs/eino-workbench.local.yaml` | `test-results/eino-workbench-fobrain-sample-discovery-report.json` | `go run ./scripts/fobrain_sample_discovery ...` + report validate | Pass | 当前 report 为 `blocked`，不是 live pass。 |
| blocked report 阻止最终声明 | `docs/fobrain-live-read-batch-plan.md` | report `blocks_claims` | schema validation | Pass | 阻止 `fobrain-batch-d live pass` 和 `Fobrain 24 readonly final acceptance`。 |

## 当前报告状态

- `status`: `blocked`
- `failure_category`: `sample_candidate_not_found`
- `owner_present`: `true`
- `department_present`: `true`
- `ip_present`: `false`
- `sample_file.written`: `false`

## 后续通过条件

当 discovery 找到 `ip_present=true` 后，可显式写出本地 samples 文件：

```bash
go run ./scripts/fobrain_sample_discovery \
  --config configs/eino-workbench.local.yaml \
  --output test-results/eino-workbench-fobrain-sample-discovery-report.json \
  --samples-output test-results/eino-workbench-fobrain-batch-d-samples.local.json
```

随后用本地 samples 文件中的值运行 Batch D smoke。samples 文件不得提交，不得复制到验收记录。
