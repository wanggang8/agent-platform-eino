# Phase 6 Approval Resume Acceptance

日期：2026-07-03

## 范围

本记录覆盖 Phase 6.2 approval interrupt/resume 后端切片：approval-required capability 创建 waiting pending、审批前不执行 mutation、安全 `resume_ref:` / `checkpoint_ref:`、内部 continuation checkpoint、approve 后恢复并执行一次 capability、reject 后不可 approve、duplicate resume 幂等、进程重启后继续执行，以及 Product Facts audit/replay 事实入库。

本记录不声明 clarification resume、Pending UI/SSE 视觉验收、真实写域 provider 或 Fobrain write approval 已完成。

## 环境

- Go：go1.23.12 darwin/arm64
- Node：v25.8.1
- npm：11.11.0
- OS：Darwin 25.5.0 arm64
- Eino 版本：github.com/cloudwego/eino v0.9.12
- 模型 provider：mock / 单元测试替身
- Fobrain 环境：未使用

## 验收结论

| 门禁 | 证据 | 结论 |
| --- | --- | --- |
| approval requested | `TestApprovalInterruptCreatesWaitingPendingAndDoesNotExecuteMutation` 覆盖 waiting pending、audit 和 checkpoint binding | 通过 |
| 审批前不执行 mutation | 同一测试断言 tool runner 未被调用 | 通过 |
| approve 后继续执行 | `TestApprovalInterruptResumeDuplicateApproveDoesNotRerunMutation` 覆盖 approve 后 capability 执行 | 通过 |
| duplicate resume | 同一测试重复相同 `client_request_id`，断言只执行一次 | 通过 |
| reject 后不可 approve | `TestApprovalInterruptResumeRejectCannotApprove` 覆盖 rejected、`approval_rejected` 和后续 approve 拒绝 | 通过 |
| reject 不依赖 checkpoint | `TestApprovalInterruptResumeRejectDoesNotRequireCheckpoint` 覆盖 checkpoint 丢失时仍可安全关闭审批 | 通过 |
| expired 后不可 approve | `TestApprovalInterruptResumeExpiredCannotApprove` 覆盖 pending timeout 后旧 resume_ref 不可恢复 | 通过 |
| restart resume | `TestApprovalInterruptResumeAfterRestartContinuesMutation` 使用 SQLite Product Facts + checkpoint store 重启后恢复 | 通过 |
| approval 后仍执行 provider policy | `TestToolLoopApprovedCapabilityStillEnforcesProviderPolicy` 覆盖缺凭据、跨 workspace 和 connector unavailable 不执行 | 通过 |
| approval 后 policy 允许才执行 | `TestToolLoopApprovedCapabilityRunsAfterPolicyAllows` 覆盖已审批且凭据/scope/connector 通过后执行 | 通过 |
| checkpoint missing 回归 | `go test ./internal/einoapp/execution -run 'ResumeCheckpointMissing|ResumeAfterRestart' -count=1` | 通过 |
| 全量 Go 回归 | `go test ./...` | 通过 |

## 已执行命令

```bash
go test ./internal/einoapp/execution -run 'ApprovalInterrupt|ResumeDuplicate|ResumeReject|ResumeAfterRestart' -count=1
go test ./internal/einoapp/execution -run 'ToolLoopApproved|ApprovalInterrupt|ResumeDuplicate|ResumeReject|ResumeAfterRestart' -count=1
go test ./internal/einoapp/capabilities -run 'ProviderPolicy' -count=1
go test ./internal/einoapp/execution -count=1
go test ./internal/einoapp/store/sqlite -run 'Pending|Approval|Checkpoint' -count=1
go test ./...
```

## 剩余风险

- Phase 6.2 只完成后端 approval resume；Pending UI/SSE patch、clarification submit/cancel、真实 Eino graph interrupt resume 和视觉验收仍需后续任务完成。
- 当前 continuation payload 保存在内部 checkpoint store；Product Facts、audit、ActionResult 和 replay 只能看到安全引用与安全摘要，不能展示内部 checkpoint id。
- reject 后 run 使用 `failed` + `safe_error=approval_rejected`；如产品需要单独的 user-declined 视觉状态，应在投影层新增安全映射，不应改变 provider 或 checkpoint 事实。
