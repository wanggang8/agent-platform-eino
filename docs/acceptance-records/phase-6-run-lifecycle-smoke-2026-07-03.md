# Phase 6 Run Lifecycle Smoke Acceptance

日期：2026-07-03

## 范围

本记录覆盖 Phase 6.5 后端 run lifecycle 控制切片：`cancel`、`stop`、`provider_timeout`、`pending_timeout`、`schema_invalid` retry、终态幂等、Product Facts 投影、replay/audit 和 SQLite pending 状态。它不声明前端 `RunNotice`、provider timeout retry、完整 approval/clarification UI 或持久化 event cursor 已完成。

## 验收结论

| 门禁 | 证据 | 结论 |
| --- | --- | --- |
| lifecycle API 按命令层执行 | `go test ./internal/einoapp/httpapi -run 'RunLifecycleEndpoint' -count=1` | 通过 |
| cancel/stop/timeout/retry 状态机 | `go test ./internal/einoapp/execution -run 'RunCancel|RunPendingTimeout|RunRetry|RunLifecycleRejects' -count=1` | 通过 |
| lifecycle API 冲突语义 | `go test ./internal/einoapp/httpapi -run 'RunLifecycleEndpoint' -count=1` 覆盖不允许操作和 facts 幂等冲突均返回 409 且不可重试 | 通过 |
| SQLite lifecycle 事务迁移 | `go test ./internal/einoapp/store/sqlite -run 'LifecycleTransition|ApplyLifecycle|CreateRetryRun' -count=1` | 通过 |
| 本地服务 smoke | `bash scripts/eino_workbench_server_smoke.sh --scenario run-lifecycle` 覆盖 cancel/stop active tool 关闭、provider timeout failed tool、pending expired、retry 合成 user message | 通过 |
| 产品出口无 raw provider payload | `run-lifecycle` smoke 检查 `authorization`、`api_key`、`api_token`、`provider_payload`、`raw_payload`、`raw provider`、`raw body` 和 bearer 痕迹 | 通过 |

## 已执行命令

```bash
go test ./internal/einoapp/execution -run 'RunCancel|RunPendingTimeout|RunRetry|RunLifecycleRejects' -count=1
go test ./internal/einoapp/httpapi -run 'RunLifecycleEndpoint' -count=1
go test ./internal/einoapp/store/sqlite -run 'LifecycleTransition|ApplyLifecycle|CreateRetryRun' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario run-lifecycle
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
go test ./...
```

## 剩余风险

- 当前 smoke 为了测试 running/waiting 状态，使用 SQLite seed 准备 Product Facts；后续 HITL/runner 完整接入后应改为纯 API 驱动。
- provider timeout retry 仍未开放；需要 capability metadata 能证明只读工具后再接入安全 retry policy。
- lifecycle transition 使用 run status compare-and-set；retry 使用仓储事务写入 idempotency、新 run、turn 和 audit。
- HTTP 将 lifecycle 不允许操作、CAS 冲突和 retry 幂等资源冲突统一为 409 类安全错误，不作为 500 可重试错误暴露。
- Workbench `RunNotice` 视觉与交互仍需在 UI 任务中验收。
- `clarification` smoke 仍按 Phase 6 后续任务保持 `exit 2`。
