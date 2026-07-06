# Phase 8 Batch B Live Smoke Infra

日期：2026-07-06
方案：Batch B “我的范围”live smoke/report 基础设施
状态：Passed for local infra tests

## Scope

- 新增 `scripts/fobrain_batch_b_smoke`。
- 新增 `docs/schemas/fobrain/batch_b_live_report.v1.schema.json`。
- `scripts/eino_workbench_server_smoke.sh` 支持 `--scenario fobrain-batch-b`。

## Verification

```bash
go test ./scripts/fobrain_batch_b_smoke -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-b --config /tmp/nonexistent-fobrain-b.yaml
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
```

## Evidence

- 报告 schema 要求六个 Batch B capability 各出现一次。
- 前端 contract index 已同步新增 Batch B report schema，避免 Workbench 侧契约索引滞后。
- `passed` 报告必须为 live mode、`blocks_claims=[]`，且六个工具均有 `result_ref`、`item_count >= 1`、`policy_decision=allowed`。
- 任一工具空结果时，runner 生成 `blocked` report，`failure_category=empty_result`，不允许声明 live pass。
- 配置缺失时 server smoke 沿用现有批次约定：返回 0 并写 skip report；该结果只表示阻断报告生成成功，不作为 live pass 信号。
- 报告只记录 `keyword_present`、`result_ref`、`item_count` 和稳定失败分类；不记录当前用户、部门、token、auth header、raw provider payload 或 `safe_summary`。

## Explicit Non-Claims

- 不声明 Batch B 真实环境 live pass。
- 不声明 Batch B Workbench 视觉通过。
- 不声明 Batch B Action API 同源 smoke 通过。
- 不声明 Fobrain 24/24 只读恢复完成。

## Next Gate

使用真实 `configs/eino-workbench.local.yaml` 执行：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-b --config configs/eino-workbench.local.yaml
```

只有生成 `test-results/eino-workbench-fobrain-batch-b-live-report.json` 且 `status=passed` 时，才可声明 Batch B live pass。
