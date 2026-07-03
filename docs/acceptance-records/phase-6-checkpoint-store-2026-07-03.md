# Phase 6 Checkpoint Store Acceptance

日期：2026-07-03

## 范围

本记录覆盖 Phase 6.1 checkpoint store 基础切片：SQLite Eino `CheckPointStore`、安全 `checkpoint_ref` 到内部 checkpoint id 的映射、resume 前 checkpoint 缺失检查、进程重启后 checkpoint lookup，以及 HTTP `checkpoint_missing` 安全错误。它不声明 approval interrupt、clarification interrupt、真实 Runner resume 或 Pending UI 已完成。

## 验收结论

| 门禁 | 证据 | 结论 |
| --- | --- | --- |
| SQLite checkpoint 持久化 | `go test ./internal/einoapp/store/sqlite -run CheckPointStore -count=1` | 通过 |
| checkpoint 安全引用边界 | `CheckPointStore` 测试覆盖 `checkpoint_ref:` 前缀、前后空白拒绝、`CheckpointRef != CheckpointID`、run/pending binding 校验 | 通过 |
| checkpoint binding 不可变 | `CheckPointStore` 测试覆盖相同绑定幂等、不同 checkpoint/run/pending 重绑返回冲突 | 通过 |
| resume 缺失 checkpoint 安全错误 | `go test ./internal/einoapp/execution -run 'ResumeCheckpointMissing|ResumeAfterRestart|ResumeRejects' -count=1` | 通过 |
| resume 状态门禁 | execution 测试覆盖非 waiting run 和终态 pending 返回 `ErrResumeNotAllowed` | 通过 |
| HTTP resume 安全错误映射 | `go test ./internal/einoapp/httpapi -run 'Resume.*Checkpoint|ResumeNotAllowed|ResumeRequest|ResumeMissingRun' -count=1` | 通过 |
| 服务装配启用 checkpoint resolver | `go test ./cmd/eino-workbench -run ServiceCommandsInjectCheckpointResolver -count=1` | 通过 |
| 全量 Go 回归 | `go test ./...` | 通过 |
| schema/contract 回归 | `npm run eino-workbench:schema-test && npm run eino-workbench:contract-test` | 通过 |
| lifecycle smoke 回归 | `bash scripts/eino_workbench_server_smoke.sh --scenario run-lifecycle` | 通过 |
| diff 格式 | `git diff --check` | 通过 |

## 已执行命令

```bash
go test ./internal/einoapp/store/sqlite -run CheckPointStore -count=1
go test ./internal/einoapp/execution -run 'ResumeCheckpointMissing|ResumeAfterRestart|ResumeRejects' -count=1
go test ./internal/einoapp/httpapi -run 'Resume.*Checkpoint|ResumeNotAllowed|ResumeRequest|ResumeMissingRun' -count=1
go test ./cmd/eino-workbench -run ServiceCommandsInjectCheckpointResolver -count=1
go test ./...
npm run eino-workbench:schema-test && npm run eino-workbench:contract-test
bash scripts/eino_workbench_server_smoke.sh --scenario run-lifecycle
git diff --check
```

## 剩余风险

- 当前切片只确认 checkpoint 存在性和持久化；真实 Eino Runner resume、approval decision、clarification submit 和 pending UI 仍在 Phase 6.2-6.4。
- `checkpoint_ref` 只允许 `checkpoint_ref:` 前缀的安全引用，并绑定 run/pending；它不应直接显示给用户，前端视觉验收仍需确认不出现 checkpoint/resume 标识。
- 同一 `checkpoint_ref` 首次绑定后不可变；重复相同绑定幂等，不同绑定返回冲突。
- 当前 checkpoint store 保留 checkpoint bytes，后续如需要清理策略，应实现 Eino `CheckPointDeleter` 或 TTL 清理任务。
