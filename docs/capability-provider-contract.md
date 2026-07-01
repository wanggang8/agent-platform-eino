# Capability Provider 契约

本文定义业务能力、工具、Skill、MCP 和 Connector 在新架构中的接入规范。它不描述旧项目迁移方式。

## 目标

Capability Provider 负责把业务能力注册为 Eino 可调用工具，并把执行结果转换为项目拥有的安全结构化结果。

Provider 可以来自：

- Local Tool Provider。
- MCP Provider。
- Skill Provider。
- Connector Provider。

无论来源是什么，进入产品输出前都必须经过同一套 schema、StructuredResult、Safety Gate 和 Product Facts 投影。

## 注册元数据

机器契约见 `docs/schemas/capability_catalog.v1.schema.json`。每个 capability 必须提供：

| 字段 | 说明 |
| --- | --- |
| `capability_id` | 稳定能力 id，不直接作为主聊天展示标题 |
| `display_name_zh` | 中文产品展示名 |
| `description` | 给模型选择工具的描述 |
| `input_schema_ref` | 输入 schema |
| `result_schema_ref` | 结果 schema |
| `risk_level` | none / low / medium / high |
| `side_effect` | none / read_external / write_external / local_runtime |
| `approval_required` | 是否需要审批 |
| `timeout_ms` | 调用超时 |
| `credential_binding_policy` | 是否需要凭据绑定 |
| `idempotency_required` | mutation 是否要求幂等 |

## 执行结果

Provider 返回的原始结果不是产品输出。执行结果必须转换为：

```text
provider raw result
  -> StructuredResult candidate
  -> Safety Gate
  -> Product Facts
  -> Workbench / ActionResult / Replay / Audit
```

Provider 不得直接写 Workbench view、ActionResult、audit 或 replay。

## 安全职责

Provider 必须：

- 不返回 secret、Authorization、API key、raw credential_ref。
- 不把 raw provider payload 放入 StructuredResult。
- 对错误返回 safe error code 和 safe summary。
- 对 mutation 提供 idempotency key 或业务安全引用。
- 对外部系统对象使用业务安全引用，不暴露内部恢复 token。

Safety Gate 仍然是最终产品输出门禁。Provider 声明“安全”不代表可以绕过 Safety Gate。

## Eino Tool Adapter

Capability Registry 负责把 provider metadata 转成 Eino tool。

转换规则：

- tool name 来自 registry 中稳定的 `tool_name` 元数据；`capability_id` 必须写入 ToolInfo extra 和 Product Facts，不能靠工具名反查业务分支。
- tool description 使用 provider 描述和风险提示。
- input schema 来自 `input_schema_ref`。
- Phase 4.2 mock tool adapter 只支持轻量字段 schema，并映射到 Eino `schema.NewParamsOneOfByParams`；完整 JSON Schema 2020-12、`oneOf`、`$defs` 等复杂输入必须在真实 provider/MCP adapter 前补齐，并优先使用 Eino 的 JSON Schema 参数入口。
- tool result 先进入 StructuredResult candidate。
- Eino event 只能作为生成 Product Facts 的输入，不直接暴露给前端或外部 API。

## MCP Provider

MCP Provider 必须额外处理：

- server catalog。
- initialize/session lifecycle。
- tool list cache 和 listChanged。
- MCP annotations 与项目 risk policy 合并。
- text/image/resource content 到安全结构化结果的转换。

MCP 返回的 `structuredContent` 仍然只是 StructuredResult candidate，必须经过 Safety Gate。

## Connector Provider

Connector Provider 负责真实外部系统适配。

connector 产品状态必须符合 `docs/schemas/connector_status.v1.schema.json`，凭据展示必须符合 `docs/schemas/provider_credential_binding.v1.schema.json`。

要求：

- 凭据只通过 credential binding 进入产品展示。
- connector status 使用安全健康状态，不展示真实 token 或 raw credential ref。
- 业务读工具和 connector 健康检查必须分开，不能用 connector status 替代业务读取。
- 写域工具必须显式声明 `approval_required=true` 和 `idempotency_required=true`。

## 禁止事项

- 不把旧业务 DTO 直接作为产品 schema。
- 不把旧 runtime step id、pending id、resume token 作为新契约。
- 不在 provider 内直接拼 Workbench UI 字段。
- 不把 MCP raw content 或 provider raw payload 当作已安全结果。
- 不用注册工具 id 作为主聊天标题。
