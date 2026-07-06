# Phase 8 Batch C Skip Pending Ticket

阶段：Phase 8.1 24 个只读工具矩阵  
日期：2026-07-06  
方案：Batch C live smoke 暂跳过 `pending_tickets`  
执行人：Codex

## 背景

- 真实接口直测确认 `/api/ticket/pending` 可访问，HTTP 200、业务 code 0，但当前环境 `total=0` 且 `list=null`。
- `/api/v1/ticket/pending` 在当前真实环境返回 404，只作为旧环境 fallback 保留。
- 用户确认：“没有待处理功能就先把工单的跳过”。

## 实现口径

- `tool.fobrain.pending_tickets` 在 `fobrain-batch-c` live smoke 中写入报告，但状态为 `skipped`。
- skipped 项固定为：
  - `result_state=not_run`
  - `item_count=0`
  - `failure_category=skipped_unavailable_feature`
  - `sample_source=unavailable_feature`
- Batch C live smoke 的 passed 口径只要求其余 5 个工具非空且通过 StructuredResult safety gate。
- 该口径不删除 provider catalog 中的 `pending_tickets`，只调整 Batch C live smoke 验收门禁。

## 影响

- 允许声明 Batch C 除工单外 5 个工具 live smoke 通过。
- 不允许声明 Fobrain 24/24 只读恢复完成。
- 后续恢复工单读取时，必须重新启用 `pending_tickets` live smoke，并补真实非空或明确空态验收口径。

## 验证

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-c --config configs/eino-workbench.local.yaml
node scripts/eino_workbench_report_validate.mjs --schema docs/schemas/fobrain/batch_c_live_report.v1.schema.json --report test-results/eino-workbench-fobrain-batch-c-live-report.json
```

脱敏结果：

- `status=passed`
- `tool.fobrain.business_list`：passed，20 条
- `tool.fobrain.external_high_risk_assets`：passed，20 条
- `tool.fobrain.vulnerability_status_summary`：passed，5 条
- `tool.fobrain.pending_tickets`：skipped，0 条，`skipped_unavailable_feature`
- `tool.fobrain.ip_stats`：passed，517 条
- `tool.fobrain.vul_stats`：passed，522 条

## 不得声明

- Fobrain 24/24 只读恢复完成。
- `pending_tickets` 真实 live 非空通过。
- 工单读取 Workbench 视觉通过。
- Action API 对工单读取同源通过。
