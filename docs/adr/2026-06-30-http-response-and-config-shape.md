# ADR: HTTP Response and Config Shape

状态：Accepted
日期：2026-06-30

## 背景

Phase 1.6 开始补齐后端基础骨架。需要先固定 HTTP 响应形态和配置来源，避免后续 Workbench、Action API、SSE、LLM provider 和业务 provider 各自形成不同约定。

## 决策

- 配置文件是后端服务唯一配置来源。
- 默认本地配置路径为 `configs/eino-workbench.local.yaml`，该文件必须被 `.gitignore` 忽略。
- 示例配置提交为 `configs/eino-workbench.example.yaml`。
- 启动服务必须通过 `--config <path>` 指定配置；测试和 smoke 使用临时配置文件。
- 本地配置文件可以包含密钥，但日志、错误、验收记录和截图必须脱敏。
- 配置摘要必须清理 DSN、URL userinfo、URL query 和 fragment；LLM base URL 只允许保留安全的 scheme/host/path。
- Phase 3 允许配置文件声明 capability metadata，用于填充 registry；生产启动路径不得内置 smoke、Fobrain 或业务 capability ID。
- 不使用环境变量覆盖服务配置；环境变量只允许用于前端测试工具自身的运行参数，不作为后端服务配置来源。
- HTTP 成功响应保持业务 schema 直出，例如 `eino_workbench_view.v1`、`eino_action_result.v1`。
- HTTP 错误响应统一使用 `eino_error_envelope.v1`。
- `X-Request-Id` 作为响应 header 返回；错误响应 body 也包含同一个 `request_id`。
- SSE 使用 HTTP/SSE 自身的 `id`、`event`、`data` 结构，`data` 保持 `eino_workbench_stream_event.v1`，不再额外套 HTTP 成功 envelope。

## 原因

- WorkbenchView、ActionResult、ReplayView 和 StreamEvent 本身已经是版本化产品投影文档。
- Product Facts projection 输出的业务 schema 应该直接成为成功响应，避免 HTTP 层产生第二套事实包装。
- 错误响应有共同安全需求：request id、safe detail、retryable、错误码和脱敏，因此需要统一 envelope。
- SSE 已经有协议级 envelope，再包一层 API envelope 会增加 reducer 和 replay 复杂度。
- 配置文件比环境变量更适合本项目需要的多段配置、credential binding、capability metadata、provider policy、预算和安全开关审查。

## 备选方案

- 所有成功响应统一 `{data, request_id}` envelope：拒绝。它会让业务 schema 与 Product Facts projection 多一层包装，并且不适合 SSE。
- 使用环境变量作为主配置来源：拒绝。多 provider、凭据绑定、审批策略和预算配置会变得不可审查且不利于本地复现。
- 引入 Viper/Cobra 配置体系：暂不采用。Phase 1.6 只需要明确 YAML struct、校验和脱敏摘要。

## 约束

- OpenAPI、fixtures、前端 contract 和 HTTP handler 必须保持成功响应业务 schema 直出。
- `httpapi` 不得使用 `http.Error` 输出产品错误。
- provider raw error、API key、Authorization、DSN secret、URL query secret、raw prompt 和 raw response body 不得进入日志或产品响应。
- `configs/eino-workbench.local.yaml` 不得提交。

## 验证

```bash
go test ./internal/einoapp/bootstrap ./internal/einoapp/httpapi -run 'Config|Response|SSE|RequestID' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario contract
```
