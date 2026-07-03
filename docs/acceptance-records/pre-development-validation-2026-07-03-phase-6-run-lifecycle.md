# 开发前复核记录

阶段：Phase 6.5 Run lifecycle controls  
日期：2026-07-03  
执行人：Codex

## 外部资料

| 主题 | 链接 | 访问日期 | 结论 | 对设计影响 |
| --- | --- | --- | --- | --- |
| Eino HITL | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/ | 2026-07-03 | HITL interrupt/resume 需要保留中断位置、展示信息和恢复数据。 | lifecycle cancel/timeout 只能关闭 Product Facts 的 pending 安全引用，不暴露 raw interrupt id。 |
| Eino Checkpoint / Interrupt | https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/ | 2026-07-03 | checkpoint/interrupt 是内部恢复能力，跨进程恢复依赖持久化 checkpoint。 | 本切片只处理 Product Facts lifecycle；完整 checkpoint resume 仍由 Phase 6.1-6.4 验收。 |
| Eino Callback | https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/ | 2026-07-03 | callback 适合作 tracing/metrics，不应成为产品 SSE 主来源。 | lifecycle API 仍从 Product Facts projection 返回，不直接返回 callback/event。 |
| OpenAPI 3.1 | https://spec.openapis.org/oas/v3.1.2.html | 2026-07-03 | API 描述可引用 JSON Schema 2020-12 语义。 | 新增 `/runs/{run_id}/lifecycle` path 并引用 `eino_run_lifecycle_request.v1`。 |
| JSON Schema 2020-12 | https://json-schema.org/draft/2020-12/json-schema-core | 2026-07-03 | `$id`、`$schema`、枚举和封闭 object 是版本化契约基础。 | 新 lifecycle request schema 必须封闭 `additionalProperties` 并纳入 schema/contract test。 |

## 当前架构复核

| 材料 | 已读 | 结论 |
| --- | --- | --- |
| `docs/run-lifecycle.md` | 是 | lifecycle 状态必须由 Product Facts 管理，终态幂等，retry 必须受安全策略限制。 |
| `docs/facts-contract.md` | 是 | Workbench、ActionResult、SSE、Replay、Audit 必须同源投影。 |
| `docs/07-implementation-plan.md` | 是 | Task 6.5 要求 cancel/stop/timeout/retry 和终态幂等 smoke。 |
| `docs/08-acceptance-plan.md` | 是 | `run-lifecycle` 不得继续返回 `exit 2`，但前端 RunNotice 可单独验收。 |
| `internal/einoapp/facts` | 是 | 已有 run/pending/tool 状态枚举，需要补事务式 lifecycle 迁移。 |
| `internal/einoapp/execution` | 是 | 命令层是 lifecycle 状态机入口，HTTP 不应直接写 facts。 |
| `internal/einoapp/httpapi` | 是 | 新 endpoint 必须解析请求后调用 execution，再由 product projection 返回。 |

## 新设计确认

- Workbench 和 Action API 共用 Product Facts：是。
- 工具结果只有 StructuredResult 一份事实材料：本切片不新增工具结果事实。
- JSON 不暴露可复用 resume token：是，lifecycle request 只接收 action 和 client_request_id。
- Eino event 不直接暴露给前端或外部 API：是。
- HITL/checkpoint/resume 场景覆盖：本切片只处理 pending timeout/cancel 终态，完整 clarification/approval 恢复仍由 Phase 6.2-6.4 覆盖。
- retry 安全策略：当前仅开放 `schema_invalid` retry；provider timeout retry 等 capability metadata 可证明只读后再开放。
- 未复制旧 runtime 类型或旧接口兼容层：是。

## 风险与处理

| 风险 | 阻断 | 处理方案 | 负责人 |
| --- | --- | --- | --- |
| lifecycle 写入 pending/run/audit 半状态 | 是 | SQLite 使用事务式 lifecycle transition；memory repository 用锁模拟原子迁移。 | Codex |
| retry 重复请求创建多个 run | 是 | 使用 `client_request_id` 进入 mutation idempotency 记录。 | Codex |
| provider timeout retry 误用于写域工具 | 是 | 本切片不开放 provider timeout retry。 | Codex |
| smoke 依赖未完成 HITL 创建 waiting 状态 | 否 | smoke 使用 SQLite seed 准备 running/waiting facts，并在验收记录中明确。 | Codex |
| 前端 RunNotice 未实现 | 否 | 后续 UI 任务验收，不在本后端切片声明完成。 | Codex |

## 结论

- 是否允许进入本 Phase：允许进入 Phase 6.5 后端 lifecycle 控制切片。
- 需要更新的 ADR：暂无。
- 需要更新的 schema/fixture/report：新增 `eino_run_lifecycle_request.v1` schema、OpenAPI path、contract index 和 Phase 6.5 acceptance record。
- 不得声明的能力：完整 HITL UI、前端 RunNotice、provider timeout retry、纯 API 驱动 waiting seed、持久化 event cursor。
