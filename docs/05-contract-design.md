# 契约设计

本文描述 API、schema、fixtures 和 typed contract。

字段级 Product Facts、投影矩阵、幂等规则见 `facts-contract.md`。上下文投影见 `conversation-context.md`。运行生命周期见 `run-lifecycle.md`。澄清链路状态机、请求响应和恢复语义见 `clarification-flow.md`。

## Product API

```text
GET  /app/
GET  /api/workspaces/{workspace_id}/views/current
POST /api/workspaces/{workspace_id}/messages
GET  /api/workspaces/{workspace_id}/runs/{run_id}/stream
GET  /api/workspaces/{workspace_id}/runs/{run_id}
GET  /api/workspaces/{workspace_id}/runs/{run_id}/replay
POST /api/workspaces/{workspace_id}/runs/{run_id}/resume
POST /api/workspaces/{workspace_id}/agent/actions
```

旧 `/workbench/`、`/chat`、`/v2/...` 不作为新项目接口。

## API 通用规则

- 所有 JSON 响应必须包含 `schema_version`。
- 错误响应使用统一 envelope：`error.code`、`error.message`、`error.safe_detail`、`error.retryable`、`request_id`。
- 2xx 只表示请求被接受或处理成功；业务等待状态通过 body 中的 `status=waiting` 表达。
- 4xx 表示请求、鉴权、schema 或幂等错误；5xx 表示服务端不可恢复错误。
- mutation、resume、approval、clarification submit 必须支持 `client_request_id` 幂等。
- 时间字段使用 RFC 3339 / ISO 8601。
- SSE reconnect 使用 `Last-Event-ID`，服务端从稳定 `event_id` 后继续发送。

## 必须定义的 schema

- `eino_workbench_view.v1`
- `eino_workbench_stream_event.v1`
- `eino_workbench_message_request.v1`
- `eino_workbench_message_response.v1`
- `eino_workbench_resume_request.v1`
- `eino_workbench_pending_interaction.v1`
- `eino_product_facts.v1`
- `eino_action_request.v1`
- `eino_action_result.v1`
- `eino_run_snapshot.v1`
- `eino_replay_view.v1`
- `eino_audit_event.v1`
- `eino_error_envelope.v1`
- `capability_catalog.v1`
- `provider_policy_decision.v1`
- `provider_credential_binding.v1`
- `connector_status.v1`
- `llm_provider_config.v1`
- `provider_redacted_error.v1`
- `tool.structured_result.v1`
- `fobrain.tool_result.v2`
- `fobrain.tool_matrix.v1`
- `fobrain.tool_inputs.v1`
- `visual_evidence_matrix.v1`

当前骨架文件位于：

```text
docs/api/eino-workbench.openapi.json
docs/schemas/
docs/schemas/fobrain/
docs/fixtures/
docs/fixtures/fobrain/
```

每个 schema 必须包含：

- `$schema: "https://json-schema.org/draft/2020-12/schema"`。
- 稳定 `$id`。
- `schema_version`。

OpenAPI 使用 3.1。validator 和 generator 必须明确支持 OpenAPI 3.1 与 JSON Schema 2020-12。可空字段使用 JSON Schema 的 `type: ["string", "null"]` 等 null type 表达，不使用 OAS 3.0 `nullable` 关键字。

## Schema 命名约定

schema 文件使用以下约定：

```text
docs/schemas/
  eino_workbench_view.v1.schema.json
  eino_workbench_stream_event.v1.schema.json
  eino_product_facts.v1.schema.json
  eino_action_result.v1.schema.json
  eino_workbench_pending_interaction.v1.schema.json
  tool.structured_result.v1.schema.json
  fobrain/
    tool_result.v2.schema.json
    <tool-id-without-prefix>.input.v1.schema.json
```

`$id` 必须与文件名稳定对应。schema 名称不得引用旧 Go 包名、旧 runtime 类型名或旧 UI DOM 名称。

## 必须覆盖的 fixtures

- empty。
- success。
- error。
- assistant delta。
- tool started。
- tool failed。
- stream duplicate。
- stream out-of-order。
- stream reconnect。
- run failed。
- run stopped。
- view replaced。
- approval waiting。
- approval approved。
- approval rejected。
- approval timeout。
- approve after reject。
- mutation not executed before approval。
- mutation executed once after approval。
- clarification requested。
- clarification submitted。
- action result success。
- action result waiting。
- action result error。
- replay view。
- audit event。
- clarification cancel。
- clarification timeout。
- resume duplicate。
- resume after restart。

当前最小 fixture 已包含：

- `docs/fixtures/workbench-view-success.json`
- `docs/fixtures/action-result-waiting-approval.json`
- `docs/fixtures/action-result-waiting-clarification.json`
- `docs/fixtures/fobrain/my-assets-result.json`

## TypeScript contract

前端类型必须由 schema/OpenAPI 生成，或由 contract test 证明与 schema fixture 等价。

前端 reducer、transport、presenter 引用生成类型，不复制字段定义。

## SSE 规则

SSE 事件必须有稳定 `event_id`。

前端 reducer 必须处理：

- 重复事件。
- 乱序事件。
- reconnect。
- tool patch。
- pending patch。
- run failed/stopped。
- view replacement。

SSE 事件必须包含：

- `event_id`。
- `run_id`。
- `type`。
- `schema_version`。
- `sequence`。
- `created_at`。

## ActionResult

`ActionResult` 给非 Workbench 客户端使用，必须包含：

- run id。
- status。
- safe assistant answer。
- structured tool results。
- presentation fields derived from StructuredResult。
- approval/resume refs。
- safe error。

它必须和 Workbench View 使用同一套 Product Facts。不得新增独立的工具结果摘要事实；展示摘要只能从 StructuredResult 派生。

`ActionResult` 必须覆盖 accepted、running、waiting、completed、failed、stopped、denied、blocked。waiting 时：

- approval 只暴露 `approval_refs`。
- clarification/user input 只暴露 `resume_refs`。
- 不得出现 `resume_token`、`credential_ref`、raw provider body 或 token hash。

详细审批状态机见 `approval-flow.md`，澄清状态机见 `clarification-flow.md`。

## 报告 schema

真实模型、live smoke、skip 记录必须有机器可读 schema：

- `schemas/real-model-provider-report.schema.json`
- `schemas/real-model-report.schema.json`
- `schemas/live-read-report.schema.json`
- `schemas/live-write-report.schema.json`
- `schemas/skip-report.schema.json`

real model provider、real model tool、live read、live write report 只表达已执行后的 `passed` / `failed`。未执行场景统一写 `skip-report.schema.json`，不得把 skipped 混入 real/live report。

skip report 必须包含 command、missing_env、credential_scope、reason、rerun_condition、blocks_claims、expires_at。

## 当前校验入口

```bash
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
npm run eino-workbench:openapi-lint
```

校验内容：

- 编译 `docs/schemas/**/*.json`。
- 校验 `docs/fixtures/manifest.json` 中登记的 fixture。
- lint `docs/api/eino-workbench.openapi.json` 的 OpenAPI 版本、operationId 和 `$ref` 路径。
- 扫描 fixture 中禁止出现的敏感标记。
- 禁止 schema / fixture 中出现 `additionalProperties: true`。

## 契约变更规则

- 破坏性 schema 变更必须新增 ADR。
- schema 变更必须同时更新 fixtures、TypeScript contract、reducer tests。
- server handler 的成功响应、错误响应、SSE payload 都必须跑 schema validator。

## 外部资料

- [OpenAPI 3.1.2](https://spec.openapis.org/oas/v3.1.2.html)
- [JSON Schema 2020-12](https://json-schema.org/draft/2020-12/json-schema-core)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
