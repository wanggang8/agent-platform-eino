# Phase 8 Batch B Live Smoke Blocked

日期：2026-07-06
方案：Batch B “我的范围”真实 live smoke
状态：Blocked

## Command

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-b --config configs/eino-workbench.local.yaml
```

## Report

本地 ignored 报告：

```text
test-results/eino-workbench-fobrain-batch-b-live-report.json
```

报告已通过：

```bash
node scripts/eino_workbench_report_validate.mjs \
  --schema docs/schemas/fobrain/batch_b_live_report.v1.schema.json \
  --report test-results/eino-workbench-fobrain-batch-b-live-report.json
```

## Result Summary

- `status=blocked`
- `failure_category=missing_current_user_scope`
- `tool.fobrain.my_department_assets`：`blocked`，`missing_current_user_scope`
- `tool.fobrain.my_department_vulnerabilities`：`blocked`，`missing_current_user_scope`
- 其余四个“我的范围”工具返回空结果，均为 `blocked`，`empty_result`
- 报告不包含 token、auth header、credential ref、raw provider payload、当前用户值、部门值或 `safe_summary`

## Decision

当前真实环境不能声明 Batch B live pass。原因不是 smoke 基础设施缺失，而是当前用户上下文不足以派生本部门范围，且当前用户范围查询返回空结果。后续需要确认当前用户接口是否有可安全解析的部门字段，或由产品确认当前账号本身没有可见的“我的范围”数据。

## Explicit Non-Claims

- 不声明 Batch B live pass。
- 不声明 Batch B Workbench 视觉通过。
- 不声明 Batch B Action API 同源 smoke 通过。
- 不声明 Fobrain 24/24 只读恢复完成。
