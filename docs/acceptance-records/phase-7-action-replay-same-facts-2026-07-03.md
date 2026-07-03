# Phase 7 Action / Replay Same Facts Acceptance

日期：2026-07-03

## 范围

本记录只覆盖 Phase 7 中 Action API、Workbench 当前视图、run snapshot、SSE、replay、audit refs 和 SQLite Product Facts 的同源事实 smoke。它不声明 `budget`、完整 run lifecycle、Fobrain 24/24 或写域审批已完成。

## 验收结论

| 门禁 | 证据 | 结论 |
| --- | --- | --- |
| ActionResult 与 Workbench 同源 | `bash scripts/eino_workbench_server_smoke.sh --scenario action-consistency` | 通过 |
| views/current 指向同一 run facts | `action-consistency` 校验 ActionResult、run snapshot、views/current 的同一 tool card、StructuredResult 和 audit | 通过 |
| replay 从 Product Facts 重建 | `bash scripts/eino_workbench_server_smoke.sh --scenario replay` | 通过 |
| workspace 隔离 | `action-consistency` 使用同一 run_id 访问其他 workspace 的 snapshot、replay、stream 并要求 404 | 通过 |
| SSE 与 replay 产品事件同源 | `action-consistency` 校验 `tool.updated` SSE；`replay` 校验由同一 Product Facts snapshot 派生的产品事件 | 通过 |
| raw provider payload 不进入产品出口 | 两个 smoke 均检查 `authorization`、`api_key`、`api_token`、`credential_ref`、`provider_payload`、`raw_payload`、`raw provider`、`raw body` 和 bearer 痕迹 | 通过 |

## 已执行命令

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario action-consistency
bash scripts/eino_workbench_server_smoke.sh --scenario replay
```

## 剩余风险

- `budget` smoke 仍按 Phase 7 后续任务保持 `exit 2`。
- cancel、timeout、retry、approval/clarification 的 replay 分支仍需在对应阶段门禁中补齐。
- 本记录使用本地 mock capability，不代表真实 Fobrain provider 全量验收。
- 当前 replay/SSE 事件序号由 Product Facts snapshot 派生，不代表已完成持久化 facts cursor/event log。
