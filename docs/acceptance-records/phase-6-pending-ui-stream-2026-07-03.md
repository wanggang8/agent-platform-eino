# Phase 6 Pending UI And Stream Acceptance

日期：2026-07-03

## 范围

本记录覆盖 Phase 6.4 最小可验收切片：Workbench 前端 reducer 消费 `pending.updated` SSE patch，approval/clarification pending 卡片按 waiting 与终态做中文产品化展示，终态 pending 只读，Product Facts 投影中的终态 pending SSE 不继续暴露旧 `resume_ref`。

本记录中的 stream 验收指 Product Facts projection + 前端 reducer，不声明真实 HTTP SSE transport 长连接已验收。本记录不声明真实 UI 按钮已调用 resume API、不声明 Playwright 视觉截图验收已完成、不声明真实 Fobrain 多候选工具已完成。

## 验收结论

| 门禁 | 证据 | 结论 |
| --- | --- | --- |
| pending stream reducer | `workbenchReducer.test.ts` 覆盖 approval 插入、clarification 终态更新、runtime pending 清理、run 隔离、旧 sequence 保护、pending 终态不回退已完成 run | 通过 |
| 终态只读卡片 | `WorkbenchShell.test.tsx` 覆盖 cancelled approval 不显示批准/拒绝动作 | 通过 |
| 终态不暴露旧 resume_ref | `TestFactsProjectionProjectsTerminalPendingAsReadonly` 覆盖 terminal pending SSE patch 无 `resume_ref` | 通过 |
| clarification smoke | `bash scripts/eino_workbench_server_smoke.sh --scenario clarification` 覆盖 execution/product/stream/contract | 通过 |
| contract 类型 | `npm run eino-workbench:contract-test` 确认 generated contract current | 通过 |

## 已执行命令

```bash
npm run eino-workbench:contract-generate
go test ./internal/einoapp/product -run 'TerminalPending|FactsProjection' -count=1
npm run eino-workbench:stream-test -- --grep pending
bash scripts/eino_workbench_server_smoke.sh --scenario clarification
npm run eino-workbench:typecheck
npm --workspace @agent-platform-eino/eino-workbench run test -- src/features/workbench/components/WorkbenchShell.test.tsx
npm run eino-workbench:test -- --run
npm run eino-workbench:contract-test
go test ./...
git diff --check
```

## 剩余风险

- 真实按钮点击调用 resume API 仍未绑定；当前卡片只验证产品化展示和只读状态。
- 视觉截图验收尚未执行，不能声明 Workbench 最终视觉通过。
- 当前 stream reducer 只覆盖 pending/run/view 相关基础路径；message/tool/audit 的增量消费和 HTTP SSE transport 长连接验收仍需后续流式 UI 任务继续补齐。
