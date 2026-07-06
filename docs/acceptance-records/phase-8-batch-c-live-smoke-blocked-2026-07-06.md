# Phase 8 Batch C Live Smoke Blocked

阶段：Phase 8.1 24 个只读工具矩阵  
日期：2026-07-06  
方案：Batch C 真实 Fobrain live smoke 首次运行  
执行人：Codex

## 环境

- 配置：`configs/eino-workbench.local.yaml`（ignored，本记录不包含 token 或 URL）
- Fobrain 模式：live
- 私有证书：按本地配置处理
- 模型 provider：不使用

## 命令

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-c --config configs/eino-workbench.local.yaml
node scripts/eino_workbench_report_validate.mjs --schema docs/schemas/fobrain/batch_c_live_report.v1.schema.json --report test-results/eino-workbench-fobrain-batch-c-live-report.json
```

## 脱敏结果

- 报告文件：`test-results/eino-workbench-fobrain-batch-c-live-report.json`
- 报告 schema：通过 `docs/schemas/fobrain/batch_c_live_report.v1.schema.json`
- 总状态：`blocked`
- 失败分类：`empty_result`
- `tool.fobrain.business_list`：passed，20 条
- `tool.fobrain.external_high_risk_assets`：passed，20 条
- `tool.fobrain.vulnerability_status_summary`：passed，5 条
- `tool.fobrain.pending_tickets`：blocked，0 条
- `tool.fobrain.ip_stats`：passed，517 条
- `tool.fobrain.vul_stats`：passed，522 条

## 处理结论

- 本次运行不能声明 Batch C live pass。
- `pending_tickets` 已按当前真实接口证据修正为优先 `/api/ticket/pending`，保留 `/api/v1/ticket/pending` 作为旧环境 fallback；参数保留 `keyword`、`status[]`、`page/page_size`，不暴露人员过滤语义。
- 真实直测结论：`/api/v1/ticket/pending` 当前返回 404；`/api/ticket/pending` 返回 HTTP 200、业务 code 0，但 `total=0` 且 `list=null`。多组状态码和 30/365 天时间窗口均未返回数据。
- 当前阻断不是凭据、TLS、schema 或 transport 问题，而是 `pending_tickets` 真实业务结果为空。
- 后续若要恢复 `pending_tickets` live 验收，需要提供一个能让 `/api/ticket/pending` 返回非空的真实状态/时间窗口样本，或调整验收口径允许该只读增量接口空态作为通过；调整验收口径前不得静默改成 passed。
- 后续处理：用户已确认当前没有待处理工单功能数据，先跳过工单读取；新的验收口径见 `docs/acceptance-records/phase-8-batch-c-skip-pending-ticket-2026-07-06.md`。
