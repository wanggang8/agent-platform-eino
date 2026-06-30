# LLM Provider Safety

本文定义模型 provider 配置、网络访问和错误脱敏要求。它适用于 Eino ChatModelAgent、OpenAI-compatible model adapter 和后续模型 provider。

## 配置要求

机器契约见 `docs/schemas/llm_provider_config.v1.schema.json`。

模型配置必须包含：

- `provider`
- `base_url`
- `api_key`
- `model`
- `timeout_ms`
- `network_safety`

`model_label` 可展示；`base_url`、`api_key`、Authorization 和 provider raw error 不得进入产品输出。

## Network Safety

Provider 初始化和请求必须校验：

- URL scheme 只能是 `https`，本地开发显式允许时才可用 `http`。
- URL 不得包含 userinfo。
- host 必须非空。
- redirect 必须重新执行安全校验。
- 默认阻断 metadata、loopback、link-local、private network 目标；测试环境可显式 allowlist。
- DNS 解析和 dial 目标都必须检查。

失败 reason code：

```text
provider_network_blocked
provider_host_not_allowed
provider_insecure_scheme_blocked
provider_redirect_blocked
provider_config_invalid
provider_config_missing
provider_timeout
provider_unavailable
provider_malformed_response
```

## Tool Call Safety

- Eino / provider tool call arguments 必须先过 JSON schema。
- invalid arguments 变成安全错误，不进入 raw prompt 或 assistant 正文。
- streaming tool call chunks 可以修复，但修复后的 arguments 仍必须 schema validate。
- tool name 只能来自 Capability Registry 转换结果。

## Error Redaction

机器契约见 `docs/schemas/provider_redacted_error.v1.schema.json`。

所有 provider error 进入 Product Facts 前必须脱敏：

- URL query secret。
- API key、token、secret、password。
- Authorization / Bearer。
- raw prompt。
- provider raw response body。
- internal stack trace。

Workbench、ActionResult、SSE、replay、audit 只能展示 safe category 和 safe summary。

## 验收

- invalid config 返回安全错误。
- network blocked 返回 typed safe error。
- non-2xx provider error 不泄漏 token、Authorization、raw prompt。
- malformed JSON 返回 `provider_malformed_response`。
- timeout 返回 `provider_timeout`。
- assistant final answer 和 delta 仍经过 Assistant Safety Gate。
