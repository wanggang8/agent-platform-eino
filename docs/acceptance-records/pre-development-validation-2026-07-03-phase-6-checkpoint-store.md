# 开发前复核记录

阶段：Phase 6.1 Checkpoint store  
日期：2026-07-03  
执行人：Codex

## 外部资料

| 主题 | 链接 | 访问日期 | 结论 | 对设计影响 |
| --- | --- | --- | --- | --- |
| Eino Checkpoint / Interrupt | https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/ | 2026-07-03 | `CheckPointStore` 是 `string -> []byte` KV；恢复依赖稳定 checkpoint id 和一致的编排/CallOption。 | SQLite 需要实现 Eino KV 接口；Product Facts 只保存 `checkpoint_ref:` 安全引用，真实 checkpoint id 保持内部映射。 |
| Eino HITL | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/ | 2026-07-03 | interrupt/resume 需要保存中断信息并用 resume data 恢复。 | Phase 6.1 只打 checkpoint/resume 基础，不声明 approval/clarification 完整恢复完成。 |
| Eino ChatModelAgent / Runner | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/ | 2026-07-03 | Runner 事件是执行输入，不能直接作为 Workbench 契约。 | checkpoint 缺失错误进入 execution safe error，不绕过 Product Facts 投影。 |
| Eino Callback | https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/ | 2026-07-03 | callback 适合 tracing/metrics，不是产品 SSE 事实源。 | checkpoint store 不新增 callback 到产品输出的旁路。 |

## 版本复核

```bash
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go list -m -versions github.com/eino-contrib/jsonschema
```

结论：

- 当前项目使用 `github.com/cloudwego/eino v0.9.12`。
- 2026-07-03 可见最新稳定仍为 `v0.9.12`，`v0.10.0-alpha.*` 不采用。
- `eino-ext/components/model/openai` 可见最新 `v0.1.13`，本切片不引入。
- `eino-contrib/jsonschema` 可见最新 `v1.0.3`，本切片不引入。

## 当前项目架构复核

| 材料 | 已读 | 结论 |
| --- | --- | --- |
| `docs/07-implementation-plan.md` Task 6.1 | 是 | 需要 SQLite checkpoint store 和 execution checkpoint 边界，覆盖 missing/restart。 |
| `docs/run-lifecycle.md` | 是 | checkpoint missing 必须返回安全错误，不能泄漏 raw checkpoint/interrupt。 |
| `docs/phase-3-execution-and-facts-plan.md` | 是 | `resume_ref -> pending -> checkpoint mapping -> Runner.Resume` 是后续链路，Phase 6.1 先落 mapping/store。 |
| `internal/einoapp/store/sqlite` | 是 | 已有 Product Facts repository；checkpoint store 应保持内部表，不进入产品投影。 |
| `internal/einoapp/execution` | 是 | `Resume` 当前只校验 run；需增加可选 checkpoint resolver，完整 HITL resume 留给 6.2/6.3。 |

## 旧项目只读参考

| 路径 | 结论 |
| --- | --- |
| `/Users/vick/Desktop/project/ai-agent/docs/STATUS.md` | 旧项目已具备 resume/recovery/checkpoint/projection/audit/replay；新项目不能复制旧 runtime，只保留能力边界。 |
| `/Users/vick/Desktop/project/ai-agent/docs/design/workbench-chat-contract.md` | approval/clarification 卡片使用安全 pending refs，禁止可复用 resume token、raw provider payload 和内部 checkpoint 细节。 |
| `/Users/vick/Desktop/project/ai-agent/internal/runtime/resume_run.go` | 旧实现有复杂 token/hash/step 绑定；新实现只参考服务端校验要求，不复制旧类型。 |

## 新设计确认

- Workbench 和 Action API 继续共用 Product Facts：是。
- 工具结果继续只有 StructuredResult 一份事实材料：本切片不新增工具结果。
- JSON 输出禁止可复用 `resume_token`：是。
- Eino checkpoint id 不进入 Product Facts：是，使用安全 `checkpoint_ref` 映射内部 checkpoint id。
- checkpoint missing 覆盖：本切片实现安全错误和 after-restart lookup。
- approval/clarification 完整恢复：否，留给 Phase 6.2/6.3，不在本切片声明。

## 风险与处理

| 风险 | 阻断 | 处理方案 |
| --- | --- | --- |
| raw Eino checkpoint id 写入 Product Facts | 是 | SQLite checkpoint table 保存内部 id；Product Facts 只保存 `checkpoint_ref:` 前缀安全引用，并拒绝 `CheckpointRef == CheckpointID`。 |
| checkpoint_ref 误绑定其他 pending | 是 | resolver 必须带 run_id / pending_id 校验，不能只按 ref 查。 |
| resume 时先消费 pending 后发现 checkpoint 缺失 | 是 | execution 先 resolve checkpoint，再由后续 HITL 切片消费 resume。 |
| 只实现 Eino KV，缺少 safe ref 映射 | 是 | checkpoint store 同时提供 `BindCheckpointRef` / `ResolveCheckpointID` 内部接口。 |

## 结论

- 是否允许进入本 Phase：允许进入 Phase 6.1 checkpoint store。
- 需要更新的 ADR：暂无。
- 需要更新的 schema/fixture/report：新增 Phase 6.1 acceptance record；当前不新增外部 API schema。
- 不得声明的能力：完整 approval resume、clarification resume、pending UI、纯 API 驱动 HITL smoke。
