# 技术架构

本文描述怎么实现。产品需求见 `01-product-requirements.md`。

## 架构原则

- Eino 负责通用 Agent 执行能力。
- 项目负责产品契约、安全投影、展示事实和业务能力接入边界。
- Runtime core 不 import 具体业务 provider。
- Workbench 和 Action API 共用 Product Facts。
- 不以兼容旧 runtime 类型、旧执行链路或旧 Workbench 接口为目标。
- 不复制旧 step runner、pending、tool pair 技术模型；Eino Runner、interrupt、checkpoint 是新的执行事实来源，项目只维护产品事实、安全投影和业务接入边界。
- 旧项目只能作为产品能力参考和验收基线，不能作为新模块结构、状态模型或接口兼容目标。

## 总览

```text
React/Vite Workbench
  -> Typed client
  -> Workbench Product API / Action API
  -> Product Mapper / Safety Gate
  -> Eino Execution
      -> ChatModelAgent + Runner
      -> Graph / Workflow
      -> Tools
      -> Interrupt / Checkpoint
      -> Callback / Tracing
  -> Capability Registry
      -> Local Tool Provider
      -> MCP Provider
      -> Skill / Connector Provider
  -> Product Facts
  -> Projection / Audit / Replay
  -> External Systems
```

## Eino 使用边界

开放式对话 Agent 默认使用 Eino ADK `ChatModelAgent + Runner`。Runner event 是生成 Product Facts 的输入来源；前端和外部 API 只消费 Product Facts 投影，不直接消费 Eino 内部事件。

Graph 用于需要分支、循环或可视化的确定性编排。Workflow 只用于 DAG 固定流程和 data mapping，不用于循环 ReAct。

Callback 用于 tracing、latency、metrics、internal diagnostics，不作为主 SSE 来源。

Checkpoint/interrupt 用于内部暂停恢复，但 checkpoint ID 和 interrupt ID 不暴露给前端或外部 API。

模型输入上下文必须从 Product Facts 安全投影，规则见 `conversation-context.md`。运行生命周期、cancel/stop/timeout/retry 规则见 `run-lifecycle.md`。callback、telemetry 和预算门禁见 `observability-and-budgets.md`。

实现以 `go.mod` 锁定版本和对应 tag 源码为准；在线文档只作概念参考。升级 Eino 必须新增 `docs/adr/YYYY-MM-DD-title.md`，并重跑 HITL、checkpoint、tool loop、callback 验收。

## 建议目录

```text
cmd/eino-workbench/
web/eino-workbench/
internal/einoapp/
  architecture/
  bootstrap/
  httpapi/
  product/
  capabilities/
  execution/
  facts/
  store/sqlite/
  llm/
  observability/
  providers/fobrain/
```

## 后端基础边界

在实现 stream reducer、Eino chat、Action API 或业务 provider 前，必须先完成后端基础边界：

- `bootstrap` 统一加载 server、database、LLM、security、observability、timeout、budget 配置，并提供脱敏后的配置摘要。
- `httpapi` 统一处理请求解析、request id、错误响应、JSON 编码和 SSE 编码。成功响应保持 OpenAPI 业务 schema 直出；错误响应统一 `eino_error_envelope.v1`。
- `llm` 只暴露 provider interface、config、mock provider、redacted error 和 network policy；真实模型 provider 不得绕过该接口进入 execution。
- `capabilities` 统一维护 provider interface、registry、tool metadata、risk policy、approval policy 和 Eino tool adapter。工具选择必须经过 registry 和 policy，不得按工具名或自然语言关键词硬编码。
- `facts` 定义 Product Facts repository interface 和事实模型；持久化实现放在 `store/sqlite`。
- `product` 定义 Workbench、ActionResult、SSE、Replay、Inspector 的安全投影接口。HTTP、Action API 和前端不得直接消费 execution event 或 provider payload。
- 不设通用 `utils` 包；复用能力必须按职责归属到上述 package，避免跨层隐式依赖。

## Capability Registry

Capability Registry 接收业务能力并转成 Eino 可调用工具。

意图识别与能力选择边界见 `intent-and-capability-selection.md`。顶层 selection 只决定入口；普通自然语言默认进入 ChatModelAgent，由模型基于 registry 暴露的 tool description 和 input schema 进行 native tool call。不得在后端按自然语言关键词硬编码业务工具路由。

Provider 接入规则见 `capability-provider-contract.md`。该契约优先于具体业务 provider 的内部实现。

能力来源：

- Local Tool Provider。
- MCP Provider。
- Skill Provider。
- Connector Provider。

Registry 维护：

- tool name。
- display name。
- description。
- input schema。
- result schema。
- risk level。
- required approval。
- timeout。

Fobrain 是 Connector Provider 样板，不是 runtime core 的一部分。

业务 provider 可以参考旧产品能力清单，但必须输出新架构定义的 input schema、StructuredResult、Product Facts 和安全投影。不得把旧 provider raw DTO、旧 runtime step id、旧 resume token 或旧 UI 展示字段作为新契约。

## MCP Adapter

MCP 不是普通本地函数列表。MCP adapter 必须处理：

- transport：stdio、HTTP/SSE 或 streamable HTTP 必须显式配置，不允许隐式信任本机任意命令。
- server catalog：server id、transport、command/url、auth、workspace scope、enabled、timeout、risk defaults。
- initialize。
- protocol version / capabilities negotiation。
- session lifecycle：启动、重连、关闭、超时、认证失败。
- `tools/list` 分页。
- `listChanged` 后刷新工具，并使 tool cache 失效。
- `inputSchema` 到 Eino tool schema。
- `outputSchema`、`structuredContent`、`isError` 到 StructuredResult 候选。
- text/image/resource content 的安全转换；resource/image 只落地为安全引用或受控附件，不直接进入主聊天。
- timeout。
- MCP annotations 与项目 risk policy 的合并判断；冲突时项目 risk policy 优先。
- 最小 mock MCP server 验收：initialize、tools/list、tool call、listChanged、structuredContent、isError、timeout、安全脱敏。

## Product Facts

最小事实模型：

- Run。
- Turn。
- ToolCall。
- ToolResult。
- PendingInteraction。
- AuditEvent。

Workbench view、SSE patch、ActionResult、Replay View、Inspector evidence/structured/runtime/audit 都从这些事实投影。

字段级定义、状态枚举、投影矩阵和幂等规则见 `facts-contract.md`。实现时不得绕过 Product Facts 直接为 Workbench、ActionResult、audit 或 replay 单独拼装工具结果。

## HITL 与恢复

产品层只暴露 `resume_ref`、`approval_refs`、`resume_refs` 等不可复用安全引用。

内部映射：

```text
resume_ref
  -> PendingInteraction
  -> checkpoint_id + interrupt_id + resume_data_schema
  -> validate request
  -> typed resume data
  -> adk.TypedRunner[M].ResumeWithParams(ctx, checkpoint_id, &adk.ResumeParams{Targets: map[string]any{interrupt_id: typedResumeData}})
```

Eino CheckPointStore 只提供 checkpoint bytes 的 `Get` / `Set` 接口；持久化、cleanup、丢失 checkpoint 的安全错误和重复 resume 幂等都由本项目实现。P0/P1 默认使用 SQLite CheckPointStore 支撑产品级恢复；未来可以替换为其他持久化实现，但必须保持相同恢复语义：

- checkpoint 可持久化。
- 进程重启后可恢复。
- checkpoint 损坏或丢失返回安全错误。
- waiting 状态 checkpoint 不被 cleanup 删除。

## 外部资料

- [Eino ChatModelAgent](https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/)
- [Eino Agent HITL](https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/)
- [Eino Interrupt & CheckPoint](https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/)
- [Eino Callback Manual](https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/)
- [MCP Tools](https://modelcontextprotocol.io/specification/2025-06-18/server/tools)
- [MCP Lifecycle](https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle)
