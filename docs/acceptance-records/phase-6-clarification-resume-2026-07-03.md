# Phase 6 Clarification Resume Acceptance

日期：2026-07-03

## 范围

本记录覆盖 Phase 6.3 clarification interrupt/resume 后端切片：安全候选/自由文本 submit、cancel、duplicate submit 幂等、expired 后拒绝、submit 前 checkpoint 存在性校验、cancel 不依赖 checkpoint、SQLite 重启后 submit，以及 Product Facts audit/replay 事实入库。

本记录不声明 Pending UI/SSE 视觉验收、真实 Eino graph resume、真实 Fobrain 多候选工具或前端澄清卡已完成。

## 环境

- Go：go1.23.12 darwin/arm64
- Node：v25.8.1
- npm：11.11.0
- OS：Darwin arm64
- Eino 版本：github.com/cloudwego/eino v0.9.12
- 模型 provider：mock / 单元测试替身
- Fobrain 环境：未使用

## 验收结论

| 门禁 | 证据 | 结论 |
| --- | --- | --- |
| clarification requested | `TestNewClarificationPendingEventBuildsSafeRunnerEvent` 覆盖安全 pending event | 通过 |
| 安全 submit | `TestClarificationSubmitConsumesPendingAndResumesRun` 覆盖候选 refs + free text、pending consumed、run running、audit | 通过 |
| duplicate submit | `TestClarificationDuplicateSubmitIsIdempotent` 覆盖相同 `client_request_id` 不重复追加 resume audit | 通过 |
| duplicate payload conflict | `TestClarificationDuplicateSubmitRejectsDifferentPayload` 覆盖同一 `client_request_id` 绑定不同 payload 返回冲突 | 通过 |
| cancel | `TestClarificationCancelClosesPendingWithoutCheckpoint` 覆盖无 checkpoint 时也能取消，并阻断后续 submit | 通过 |
| submit checkpoint missing | `TestClarificationSubmitCheckpointMissingDoesNotConsumePending` 覆盖 checkpoint 缺失时不消费 pending | 通过 |
| expired / unsafe data | `TestClarificationSubmitRejectsExpiredAndUnsafeResumeData` 覆盖 expired、未知 candidate、unsafe free text | 通过 |
| restart submit | `TestClarificationAfterRestartSubmitsWithCheckpoint` 使用 SQLite Product Facts + checkpoint store 重启后恢复 | 通过 |
| HTTP resume contract | `TestResumeRequestParsesClarificationCancelAndMixedInput`、schema fixtures 覆盖 cancel 和 mixed 输入 | 通过 |
| audit schema | `eino_audit_event.v1` enum 增加 `approval` / `clarification` 并通过 schema/contract 检查 | 通过 |
| 相关回归 | `go test ./internal/einoapp/execution -count=1`、`go test ./internal/einoapp/facts ./internal/einoapp/store/sqlite -count=1` | 通过 |

## 已执行命令

```bash
go test ./internal/einoapp/execution -run 'NewClarificationPendingEvent|ClarificationSubmit|ClarificationCancel|ClarificationDuplicate|ClarificationAfterRestart' -count=1
go test ./internal/einoapp/httpapi -run 'ResumeRequestParsesClarificationCancelAndMixedInput|InvalidResumeDecisionWithCandidate' -count=1
go test ./internal/einoapp/execution -run 'ClarificationDuplicateSubmitRejectsDifferentPayload|ClarificationSubmitCheckpointMissingDoesNotConsumePending' -count=1
go test ./internal/einoapp/execution -count=1
go test ./internal/einoapp/facts ./internal/einoapp/store/sqlite -count=1
npm run eino-workbench:contract-generate
npm run eino-workbench:schema-test && npm run eino-workbench:contract-test
go test ./...
git diff --check
```

## 剩余风险

- 当前切片只迁移 Product Facts 和安全 checkpoint 校验；真实 Eino graph resume 与 UI/SSE pending patch 仍在后续任务。
- submit 后 run 暂停态恢复为 `running`，后续执行结果仍需由 runner/tool loop 写入终态。
- clarification resume data 只进入安全 audit 摘要；不得把 raw provider candidate、内部 checkpoint id 或 reusable resume token 写入 Product Facts。
