# Phase 7 Budget Telemetry Acceptance Record

## Scope

本记录覆盖 Phase 7.2 的最小 budget exceeded 切片，不声明完整 Eino callback、token/cost 估算、workspace quota 或 rate limit 已完成。

## Change Summary

- 新增 `budget_exceeded` lifecycle action，由 execution 层统一更新 Product Facts。
- `budget_exceeded` 只允许作用于 `running` / `waiting` run，结果为安全失败 `safe_error=budget_exceeded`。
- active tool 会被取消，waiting pending 会过期，避免旧 resume ref 在预算终态后继续使用。
- replay events 包含 `event_type=budget`、`safe_summary=budget exceeded` 的安全 audit event。
- audit event contract 已同时覆盖既有 `lifecycle` 和新增 `budget`，避免 replay 中已有生命周期事件与 schema/TS 契约不一致。
- `budget` smoke 从 Phase 7 占位改为真实本地服务验收。
- `budget` smoke 覆盖两条路径：`running + active tool` 取消，以及 `waiting + approval pending` 在 SQLite 与 replay 中过期。

## Verification

```bash
go test ./internal/einoapp/execution -run Budget -count=1
npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario budget
bash scripts/eino_workbench_server_smoke.sh --scenario run-lifecycle
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget|RunLifecycle' -count=1
go test ./...
git diff --check
```

结果：以上命令均通过。

## Safety Checks

`budget` smoke 会断言 ActionResult 和 Replay payload 不包含 `authorization`、`api_key`、`api_token`、`provider_payload`、`raw_payload`、`raw prompt`、`raw provider`、`raw body` 或 `bearer `。

## Remaining Risks

- 预算判断层尚未接入真实 Eino callback、token usage、latency、cost estimate 或 workspace quota。
- 当前 `budget_exceeded` 是已判定超限后的安全终态入口；后续仍需实现预算计数、阈值配置和 telemetry 报告。
