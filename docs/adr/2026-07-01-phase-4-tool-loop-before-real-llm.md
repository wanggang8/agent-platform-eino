# ADR: Phase 4 Tool Loop Before Real LLM

状态：Accepted
日期：2026-07-01

## 背景

Phase 3 已完成 mock ChatModelAgent、Product Facts、Action API、SSE、capability selection 和 context projection 的主链路。Phase 4 需要引入 StructuredResult、Safety Gate、Eino tool adapter、真实 LLM provider 和 provider policy。

这些能力互相依赖，但风险不同：StructuredResult 和 tool facts 是产品事实链路，真实 LLM provider 是外部网络和凭据集成。如果先接真实 LLM，再验证 tool loop，失败原因会混杂在模型兼容、网络策略、凭据和产品事实之间，不利于定位。

## 决策

- Phase 4 先实现 StructuredResult 与 Safety Gate。
- 随后用本地 mock capability provider 跑通 Eino tool loop 和 `tool-card` smoke。
- 真实 OpenAI-compatible LLM provider 放在 tool loop 之后接入。
- MCP adapter 只在 Phase 4 完成 contract/mock，不阻塞 P0 tool-card 主链。
- 所有 provider raw payload 必须先转换为 StructuredResult candidate，经 Safety Gate 后才能写 Product Facts。
- Eino tool 的 `InvokableRun` 返回值只作为模型 tool message；Workbench、ActionResult、SSE、Replay 和 Audit 只能从 Product Facts 投影。
- 写域能力在 Phase 6 HITL/checkpoint/resume 完成前不得执行 mutation。

## Eino API 依据

本项目当前锁定 `github.com/cloudwego/eino v0.9.12`。

本地 API 复核结论：

- `components/tool.BaseTool` 通过 `Info(ctx)` 暴露 `schema.ToolInfo`。
- `components/tool.InvokableTool` 通过 `InvokableRun(ctx, argumentsInJSON)` 执行工具。
- `schema.ToolInfo` 支持 `ParamsOneOf`，可用 JSON Schema 参数描述。
- `adk.ChatModelAgentConfig.ToolsConfig` 持有工具列表；工具执行事件仍不能直接成为产品输出。

## 原因

- StructuredResult / Safety Gate 是跨 Workbench、Action API、Replay、Audit 的事实边界，必须先稳定。
- mock tool loop 能在无外部凭据、无真实模型不确定性的情况下证明 tool call、tool result、audit、SSE 和 ActionResult 同源。
- 真实 LLM provider 应只替换模型边界，不改变 tool facts、product projection 或 HTTP handler。
- MCP lifecycle、auth、pagination 和 listChanged 复杂度较高，不应阻塞 P0 tool-card 验收。

## 后果

- Phase 4 任务顺序调整为：
  1. StructuredResult and Safety Gate。
  2. Mock capability adapter and Eino tool loop。
  3. Production LLM provider implementation。
  4. Provider policy and credentials。
  5. MCP adapter contract。
- `tool-card` smoke 不依赖真实 LLM provider。
- `real-model-chat` smoke 归属 P1 真实模型门禁，需要本地 ignored 配置；无凭据只能生成 skipped report，不能作为通过信号。

## 验证

```bash
go test ./internal/einoapp/product -run 'StructuredResult|Safety|AssistantSafety' -count=1
go test ./internal/einoapp/capabilities -run 'ProviderContract|StructuredResultConversion|ToolSelectionMetadata|EinoToolAdapter' -count=1
go test ./internal/einoapp/execution -run 'AgentToolLoop|ToolFacts|ToolSafety' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario tool-card
go test ./internal/einoapp/llm -run 'OpenAICompatible|NetworkPolicy|RedactedError' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario real-model-chat --config configs/eino-workbench.local.yaml
```
