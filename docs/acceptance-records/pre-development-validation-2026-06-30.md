# 开发前复核记录：2026-06-30

## 范围

- 目标 Phase：Phase 1，契约和服务骨架。
- 新项目路径：`/Users/vick/Desktop/project/agent-platform-eino`。
- 旧项目路径：`/Users/vick/Desktop/project/ai-agent`，只读参考，未修改。
- 结论：允许进入 Phase 1；Eino module 已允许固定版本，但不允许进入 Phase 3 runtime、Phase 6 或 Phase 8。

## 外部资料

| 主题 | 访问日期 | 结论 | 对本阶段影响 |
| --- | --- | --- | --- |
| Eino ChatModelAgent / Runner | 2026-06-30 | Phase 1 只固定 `github.com/cloudwego/eino v0.9.12`，不接入 runtime；Phase 3 前必须重新核对 Runner event、tool call 和 stream 行为。 | 不阻断 Phase 1。 |
| Eino HITL / Checkpoint / Callback | 2026-06-30 | Phase 1 只定义契约；Phase 6 前必须用实际版本验证 interrupt、checkpoint、callback。 | 不阻断 Phase 1。 |
| MCP lifecycle / tools | 2026-06-30 | 新契约已要求 initialize、session、tools/list、listChanged、annotations、structuredContent/isError。 | Phase 4 前重新复核。 |
| OpenAPI 3.1 / JSON Schema 2020-12 | 2026-06-30 | 本阶段采用 OpenAPI 3.1 与 JSON Schema 2020-12，schema 必须有 `$schema`、`$id` 和 `schema_version`。 | 进入 Phase 1 的基础。 |
| React / Vite / Playwright | 2026-06-30 | 本阶段只固化工具链和视觉验收矩阵；真实脚手架在 Phase 2 实现。 | 不阻断 Phase 1。 |

## 版本命令摘要

当前新项目已创建 `go.mod`，Phase 1 只固定 Eino 核心模块。未引入模块只做可见版本查询，不要求 `go list -m` 当前依赖通过：

```bash
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go list -m -versions github.com/eino-contrib/jsonschema
```

当前固定：`github.com/cloudwego/eino v0.9.12`。Phase 3 引入 Eino runtime 前必须重新执行以上命令，并复核 `docs/adr/2026-06-30-pin-eino-version.md`。

## 当前架构复核

已读取新项目全部分层文档、旧项目架构文档、Workbench 契约、UI 验收矩阵、presentation、capability/providerapi、Fobrain、LLM 与 Eino adapter 相关路径。

复核结论：

- 新项目目标是 Eino-first 全新实现，不复制旧 runtime 类型、旧 step runner、旧 tool pair、旧 Workbench DOM 或旧接口兼容层。
- 旧项目只作为产品能力、安全边界、视觉基线和验收覆盖参考。
- 已复核旧项目最终验收聚合报告：`/Users/vick/Desktop/project/ai-agent/test-results/workbench-acceptance-audit/workbench-acceptance-audit.json`，结论为 `6/6 passed`、`failed=0`。该证据已登记到 `legacy-acceptance-evidence.md`，只能作为范围和视觉基线参考，不能作为新项目通过结果。
- Workbench 与 Action API 必须从 Product Facts 投影。
- 工具结果只能以 StructuredResult 作为事实材料。
- Fobrain P2 需要 24 只读工具、connector、凭据绑定、实体消歧、写域审批和 live smoke 完整恢复。

## 新设计确认

- Workbench 和 Action API 共用 Product Facts：是，见 `facts-contract.md` 和 `05-contract-design.md`。
- 工具结果只有 StructuredResult 一份事实材料：是，见 `06-security-and-projection.md`、`tool.structured_result.v1`、`fobrain.tool_result.v2`。
- JSON 不暴露可复用 resume token：是，见 `approval-flow.md`、`clarification-flow.md`、`eino_action_result.v1`。
- Eino event 只作为 Product Facts 输入：是，见 `04-technical-architecture.md`。
- HITL/checkpoint/resume 场景覆盖：Phase 1 只定义契约；Phase 6 前必须实现并验证。
- Capability Provider 禁止 raw provider payload 进入产品输出：是，见 `capability-provider-contract.md`。
- MCP provider 覆盖 lifecycle 和 tool result 安全转换：设计覆盖，Phase 4 前需重新复核。
- Fobrain P2 能力已映射：类别已映射；字段级 schema、fixture、mock/live 断言必须在 Phase 8 Task 8.1 前补齐。
- 前端视觉基线已准备：已有 `visual-acceptance-matrix.md` 和 `legacy-acceptance-evidence.md`；Phase 2 必须从旧最终证据提取目标截图、block crop、desktop/mobile evidence，并重新生成新项目 Playwright baseline。
- 旧项目最终通过证据没有被当成新项目通过结果：是，`08-acceptance-plan.md` 要求 P0/P1/P2 通过声明只能来自新项目当次 `test-results/eino-*` 报告。

## 风险与处理

| 风险 | 处理 |
| --- | --- |
| Eino 版本 API 与设计不一致 | Phase 3 前重新执行版本复核、更新 ADR 并补齐 runtime 验收。 |
| Fobrain 字段级恢复不足 | Phase 8 前逐工具补齐 schema、fixture、mock/live 断言。 |
| Workbench 视觉只满足大方向 | Phase 2 必须补目标截图、block crop、desktop/mobile evidence。 |
| ActionResult 与 Workbench 事实分叉 | Phase 1 schema、fixture 和 `action-consistency` smoke 必须校验同源投影。 |

## 结论

- 是否允许进入 Phase 1：是。
- 是否允许进入 Phase 3：否，需正式 Eino runtime API 复核和执行层验收。
- 是否允许进入 Phase 6：否，需 HITL/checkpoint 实现设计复核。
- 是否允许进入 Phase 8：否，需 Fobrain 字段级矩阵补齐。
- 是否允许声明 P0/P1/P2 完成：否，需按 `08-acceptance-plan.md` 执行阶段验收。
