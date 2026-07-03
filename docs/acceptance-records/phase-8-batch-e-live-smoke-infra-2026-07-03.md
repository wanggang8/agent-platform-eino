# Phase 8 Batch E Live Smoke Infrastructure Acceptance Record

日期：2026-07-03

范围：Fobrain Batch E 样本发现 sidecar、live smoke 脚本和脱敏 report schema。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| Discovery 本地 samples 可携带 Batch E 样本 | `docs/fobrain-batch-e-interface-plan.md` | `scripts/fobrain_sample_discovery/main_test.go` | `go test ./scripts/fobrain_sample_discovery -count=1` | Pass | 安全 report 只记录 presence booleans；真实 ID/业务名/漏洞名只写入 local samples。 |
| Discovery passed 必须覆盖 Batch E 样本完整性 | `docs/schemas/fobrain/sample_discovery_report.v1.schema.json` | `TestRunDiscoveryBlocksWhenBatchESamplesMissing` | `go test ./scripts/fobrain_sample_discovery -run BatchESamplesMissing -count=1` | Pass | Batch D 样本存在但 Batch E 样本缺失时 status=blocked，且不写 local samples。 |
| Discovery 不再默认认证参数 | `docs/fobrain-provider-config.md` | `TestDiscoveryClientRequiresExplicitAuthParam` | `go test ./scripts/fobrain_sample_discovery -run Discovery -count=1` | Pass | 缺 `auth_param` 时阻断，不自动补 `authorization`。 |
| Batch E live smoke 调用四个 provider capability | `docs/fobrain-live-read-batch-plan.md` | `scripts/fobrain_batch_e_smoke/main_test.go` | `go test ./scripts/fobrain_batch_e_smoke -count=1` | Pass | 通过 provider Invoke 和 StructuredResult safety gate，不读取 raw provider payload。 |
| Batch E blocked smoke 返回非零 | `docs/fobrain-live-read-batch-plan.md` | `TestRunBatchELiveSmokeBlocksWhenSamplesMissing` | `go test ./scripts/fobrain_batch_e_smoke -run BlocksWhenSamplesMissing -count=1` | Pass | 避免 CI 只看 exit code 时把 blocked 当 pass。 |
| Batch E report schema 可校验 | `docs/schemas/fobrain/batch_e_live_report.v1.schema.json` | schema validation | `npm run eino-workbench:schema-test` | Pass | schema 覆盖 passed/blocked/failed、四个 capability 和 blocks_claims 门禁。 |
| server smoke 支持 `fobrain-batch-e` skip | `scripts/eino_workbench_server_smoke.sh` | skip scenario output | `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-e --config /tmp/nonexistent-eino-workbench.yaml` | Pass | 无本地配置时生成 skip，不声明 live pass。 |
| 真实 discovery 缺 Batch E 业务样本时只阻断 E | local live discovery safe report | `test-results/eino-workbench-fobrain-sample-discovery-report.json` | `go run ./scripts/fobrain_sample_discovery --config configs/eino-workbench.local.yaml --output test-results/eino-workbench-fobrain-sample-discovery-report.json --samples-output .agent/tasks/fobrain-batch-e-live-smoke/batch-samples.local.json` | Blocked | 当前安全报告显示 owner/department/IP、资产详情、漏洞详情、威胁名 present，`business_present=false`；`blocks_claims` 只阻断 `fobrain-batch-e live pass` 和最终 24 只读验收。 |

## 验证命令

```bash
go test ./scripts/fobrain_sample_discovery -count=1
go test ./scripts/fobrain_batch_e_smoke -count=1
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-e --config /tmp/nonexistent-eino-workbench.yaml
go run ./scripts/fobrain_sample_discovery --config configs/eino-workbench.local.yaml --output test-results/eino-workbench-fobrain-sample-discovery-report.json --samples-output .agent/tasks/fobrain-batch-e-live-smoke/batch-samples.local.json
```

## 剩余风险

- 本记录不包含真实 Fobrain Batch E live pass；真实样本必须来自 ignored local samples 或人工只读确认。
- 当前真实 discovery 尚未发现稳定 `business` 样本，因此不能生成 Batch E local samples 文件，也不能执行完整 Batch E live pass。
- 本记录不包含 Workbench 视觉、Action API 同源展示、replay/audit 证据。
- 不得据此声明 Batch E 完成、Fobrain 24 个只读工具恢复完成或产品能力可比。
