# 开发前复核记录

阶段：Phase 5 Fobrain provider PoC  
日期：2026-07-02  
执行人：Codex

## 外部资料

| 主题 | 链接 | 访问日期 | 结论 | 对设计影响 |
| --- | --- | --- | --- | --- |
| Eino ToolsNode / Tool Guide | https://www.cloudwego.io/docs/eino/core_modules/components/tools_node_guide/ | 2026-07-02 | Tool 是模型可选择的外部能力，执行由工具节点承担。 | Fobrain 必须作为 capability/tool provider 接入，不在 HTTP 或 execution 中硬编码业务分支。 |
| Eino ChatModelAgent / Runner | https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/ | 2026-07-02 | 配置 tools 后进入 ReAct tool call 循环。 | Phase 5 不改模型循环，只补 provider invoker 和 StructuredResult candidate。 |
| MCP tools | https://modelcontextprotocol.io/specification/2025-06-18/server/tools | 2026-07-02 | tool 有名称、元数据、输入 schema 和 structured output。 | Fobrain provider 输出仍须先转 StructuredResult candidate，不能把 structuredContent/raw payload 当 Product Facts。 |
| Provider policy / credential 本地契约 | `docs/provider-policy-and-credentials.md` | 2026-07-02 | Fobrain 只读能力使用 `policy:fobrain:read:v1`，缺凭据和 scope denied 必须安全折叠。 | PoC capability 必须声明 credential binding required，并覆盖泄漏回归。 |

## 版本命令摘要

Phase 5 不引入或升级 Eino、MCP、OpenAPI、JSON Schema、React/Vite 依赖；本次只复核 provider 接入边界。版本复核继续沿用 `docs/acceptance-records/pre-development-validation-2026-07-02-phase-4-5.md` 的结果。若实现过程中新增依赖，必须重新执行：

```bash
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go list -m -versions github.com/eino-contrib/jsonschema
```

## 当前架构复核

| 材料 | 已读 | 结论 |
| --- | --- | --- |
| `docs/07-implementation-plan.md` | 是 | Phase 5 原计划粒度不足，已新增 PoC 技术方案和可执行计划。 |
| `docs/fobrain-tool-matrix.md` | 是 | 24 只读完整恢复属于 Phase 8；Phase 5 只能做单工具 PoC。 |
| `docs/provider-policy-and-credentials.md` | 是 | Phase 4.4 已有 policy 和安全 CredentialBinding，Phase 5 应复用。 |
| `docs/facts-contract.md` | 是 | provider 只能产 StructuredResult candidate，Product Facts 仍是唯一事实来源。 |
| `docs/fobrain-provider-config.md` | 是 | Phase 5 需要先补 Fobrain connector 本地配置契约，避免实现期硬编码。 |
| 旧项目 `internal/business/fobrain` | 是 | 可参考能力矩阵、adapter/client、凭据和结果语义，不复制旧 providerapi/runtime。 |
| 旧项目 `internal/registry/capability`、`internal/registry/providerapi` | 是 | 旧 registry 经验只作为产品边界参考，新项目使用 `internal/einoapp/capabilities`。 |

## 新设计确认

- Workbench 和 Action API 共用 Product Facts：是，Phase 5 不新增平行事实模型。
- 工具结果只有 StructuredResult 一份事实材料：是，provider 只产 candidate。
- JSON 不暴露可复用 resume token：是，Phase 5 不涉及 resume token。
- Eino event 不直接暴露给前端或外部 API：是，仍经 Product Facts 映射。
- HITL/checkpoint/resume 场景覆盖：Phase 5 不实现，写域能力继续阻断到 Phase 6/8。
- Capability Provider 不暴露 raw provider payload：是，已列为 mapper/client 边界测试。
- MCP provider 语义覆盖：Phase 5 不改 MCP；沿用 Phase 4.5 mock MCP 结论。
- Fobrain P2 能力已映射：是，`docs/fobrain-tool-matrix.md` 覆盖 24 只读、connector、消歧、写域审批；Phase 5 不声明完整恢复。
- 前端视觉基线已准备：Phase 5 smoke 只要求 PoC 工具卡；完整截图矩阵在 Phase 8。
- 未复制旧 runtime 类型或旧接口兼容层：是，设计明确禁止复制旧 providerapi/runtime/workbench 接口。

## 风险与处理

| 风险 | 阻断 | 处理方案 | 负责人 |
| --- | --- | --- | --- |
| Phase 5 被误解为 24 只读恢复 | 否 | 在设计和计划中明确不可声明能力。 | Codex |
| catalog 变成工具名硬编码分支 | 是 | 所有 metadata 集中在 catalog，execution/httpapi 不按工具名分支，增加测试。 | Codex |
| 缺少 Fobrain 配置契约导致 endpoint/token 硬编码 | 是 | 已新增 `docs/fobrain-provider-config.md`，实现前先写 bootstrap config 测试。 | Codex |
| provider PoC 被真实模型凭据误挡 | 否 | 已拆分 provider PoC report 与 real-model report；缺凭据时用 skip report 阻断真实模型声明。 | Codex |
| raw provider payload 或 secret 进入 Product Facts | 是 | client/result mapper/Safety Gate 三层测试。 | Codex |
| 无本地 Fobrain 凭据导致 smoke 不稳定 | 否 | 无凭据时生成 skip report，并用 `blocks_claims` 阻断 live 能力声明。 | Codex |

## 结论

- 是否允许进入本 Phase：允许进入 Phase 5 设计后的实现，但必须按 `docs/superpowers/plans/2026-07-02-fobrain-provider-poc.md` 分步执行。
- 需要更新的 ADR：已新增 `docs/adr/2026-07-02-fobrain-provider-poc-boundary.md`。
- 需要更新的 schema/fixture/report：已新增 `docs/schemas/fobrain/provider_poc_report.v1.schema.json` 和 `docs/fixtures/fobrain/provider-poc-report.json`；实现阶段需生成对应 report，真实模型缺凭据时生成 skipped report。
- 不得声明的能力：24 只读完整恢复、live read 覆盖、写域审批可用、实体消歧可用、可替换旧项目。
