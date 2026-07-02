# Fobrain Provider 配置契约

本文定义 Phase 5 Fobrain provider PoC 的本地配置形态。该契约只用于新项目实现，不复制旧项目配置结构。

## YAML Shape

Phase 5 实现必须在 `bootstrap.Config` 中增加 `fobrain` 根字段：

```yaml
fobrain:
  enabled: true
  connector_id: "fobrain"
  workspace_id: "default"
  base_url: "https://fobrain.example.local"
  timeout: "10s"
  credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "authorization"
    api_token: "local-only-secret"
  connector_status:
    mode: "mock"
    available: true
```

`configs/eino-workbench.example.yaml` 只能展示无 secret 示例；真实 `api_token` 只能写入已 ignored 的 `configs/eino-workbench.local.yaml`。当前真实 Fobrain 环境使用名为 `authorization` 的认证参数，值为原始 token；不得把 token 写入仓库。

## 字段要求

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `enabled` | 是 | 未启用时不得注册 Fobrain provider。 |
| `connector_id` | 是 | 固定安全标识，Phase 5 只允许 `fobrain`。 |
| `workspace_id` | 是 | 凭据解析必须按 workspace scope 校验。 |
| `base_url` | 是 | 真实 provider HTTP 入口；必须是绝对 URL，不能硬编码在业务代码。 |
| `timeout` | 是 | provider client 超时。 |
| `credential.status` | 是 | `missing` / `unbound` / `configured` / `bound`。 |
| `credential.display_ref` | 是 | 安全展示引用，不得包含 raw credential ref、token 或 URL。 |
| `credential.owner_scope` | 是 | Phase 5 固定 `workspace`。 |
| `credential.auth_param` | 是 | Fobrain HTTP client 使用的认证参数名；当前真实环境为 `authorization`。只允许安全 header/query 参数名，不包含 token 值。 |
| `credential.api_token` | 本地可选 | 只允许 local ignored 配置使用；不得进入 RedactedSummary、Product Facts、日志、报告。 |
| `connector_status.mode` | 是 | Phase 5 只允许 `mock`；`live` 保留给真实 HTTP client 阶段。 |
| `connector_status.available` | 是 | policy 输入，不替代业务读取工具。 |

`connector_status.mode=mock` 时可以缺少 `credential.api_token`。Phase 5 服务启动配置不接受 `live` mode，避免真实配置静默降级成 mock 结果；`fobrain-poc` provider report 的 `provider_mode` 只能是 `mock`。真实模型工具选择或 Fobrain live read 在 Phase 5 只能写 skip report，不能声明通过。

## 派生对象

实现阶段必须从该配置派生：

- `capabilities.CredentialBinding`：只包含 `schema_version`、`workspace_id`、`system=fobrain`、`status`、`display_ref`、`owner_scope`、`audit_ref`。
- `capabilities.ConnectorStatus`：只表达 available/unavailable。
- `FobrainClientConfig`：只在 provider client 内部持有 `base_url`、`timeout`、secret token。
- `auth_param`：只作为 provider client 认证参数名，当前 live client 应发送 `authorization: <api_token>`；不要硬编码旧项目默认 header。

## 脱敏规则

- `api_token`、Authorization、cookie、raw credential ref、连接串和 raw provider payload 不得进入 Product Facts、ActionResult、Workbench、SSE、audit、replay、report 或日志。
- RedactedSummary 只能展示 provider enabled、connector id、workspace id、base url host、timeout、credential status、auth parameter name、connector mode 和 available。
- workspace 不匹配必须返回 `credential_scope_denied`，不能回显配置里的 workspace/token。

## 验收

实现前必须先补配置测试：

```bash
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
```

测试必须覆盖：最小合法配置、缺必填字段、非法 URL、非法 credential status、非法 auth parameter name、secret 不出现在脱敏摘要、`enabled=false` 不注册 provider。
