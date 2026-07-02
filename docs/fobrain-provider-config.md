# Fobrain Provider 配置契约

本文定义 Phase 5 Fobrain provider PoC 与 Phase 8 live read 的本地配置形态。该契约只用于新项目实现，不复制旧项目配置结构。

## YAML Shape

实现必须在 `bootstrap.Config` 中增加 `fobrain` 根字段：

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
    api_token: "<ignored-local-token>"
  connector_status:
    mode: "mock"
    available: true
  tls:
    insecure_skip_verify: false
```

`configs/eino-workbench.example.yaml` 只能展示无 secret 示例；真实 `api_token` 只能写入已 ignored 的 `configs/eino-workbench.local.yaml`。当前真实 Fobrain 环境使用名为 `authorization` 的认证参数，值为原始 token；不得把 token 写入仓库。

## 字段要求

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `enabled` | 是 | 未启用时不得注册 Fobrain provider。 |
| `connector_id` | 是 | 固定安全标识，Phase 5 只允许 `fobrain`。 |
| `workspace_id` | 是 | 凭据解析必须按 workspace scope 校验。 |
| `base_url` | 是 | 真实 provider HTTP 入口；必须是绝对 URL，不能硬编码在业务代码。live mode 默认要求 HTTPS，仅本地 loopback HTTP 可用于验收代理。 |
| `timeout` | 是 | provider client 超时。 |
| `credential.status` | 是 | `missing` / `unbound` / `configured` / `bound`。 |
| `credential.display_ref` | 是 | 安全展示引用，不得包含 raw credential ref、token 或 URL。 |
| `credential.owner_scope` | 是 | Phase 5 固定 `workspace`。 |
| `credential.auth_param` | 是 | Fobrain HTTP client 使用的认证 header 名；当前真实环境为 `authorization`。只允许安全 header 参数名，不包含 token 值。 |
| `credential.api_token` | 本地可选 | 只允许 local ignored 配置使用；不得进入 RedactedSummary、Product Facts、日志、报告。 |
| `connector_status.mode` | 是 | Phase 5 只允许 `mock`；`live` 保留给真实 HTTP client 阶段。 |
| `connector_status.available` | 是 | policy 输入，不替代业务读取工具。 |
| `tls.insecure_skip_verify` | 否 | 只用于 Fobrain 私有证书环境；默认 `false`。开启后只影响 Fobrain HTTP client，不影响 LLM 或其它 provider。 |

`connector_status.mode=mock` 时可以缺少 `credential.api_token`。Phase 5 的 `fobrain-poc` provider report 仍只能使用 `mock` 并声明 `provider_mode=mock`；Phase 8 Batch A 起，服务配置允许 `live` mode 进入真实 HTTP client。真实模型工具选择或 Fobrain live read 未完成对应批次验收时只能写 blocking skip report，不能声明通过。

## Live Read 边界

Phase 8 live read 启用前必须先完成 `docs/fobrain-live-read-batch-plan.md`。当前凭据模型是 workspace 共享 token：

- `credential.api_token` 只代表当前 workspace 的 Fobrain token。
- 不从前端、Action API 请求或用户消息中接收 token。
- 不按单个工具、单个用户或单次 run 生成独立 token。
- live client 只接受显式配置的 `auth_param`，并把它当作 header 名发送；缺失时必须阻断请求，不能在 provider 内默认成 `authorization`；如果后续确认真实 API 使用 query/body 参数，必须先新增 ADR。
- Phase 8 live acceptance 配置必须设置 `auth_param: "authorization"`；其它认证参数名需要先新增 ADR。
- `connector_status.mode=live` 时启动/执行前必须确认 `api_token` 非空、`owner_scope=workspace`、`workspace_id` 匹配；workspace 不匹配时不得发起外部 HTTP 请求。
- `connector_status.mode=live` 时 `base_url` 必须使用 HTTPS；只有 `127.0.0.1`、`localhost`、`::1` 的 HTTP URL 可作为本地验收代理入口。
- 私有证书环境可设置 `tls.insecure_skip_verify: true`；该开关必须留在 ignored local 配置或受控部署配置中，验收报告只能记录布尔值，不能记录证书、token 或连接细节。
- 当前 Batch A 代码基础覆盖 `tool.fobrain.current_user_context` 和 `tool.fobrain.my_permissions` 的当前用户读取：标准路径优先使用 `/api/v1/user`，404/405 时兼容私有部署 `/api/user`。`connector.fobrain.security` 只展示安全配置、connector 可用性和凭据绑定摘要，不替代业务读工具。
- `my_permissions` 必须支持权限空态；当前用户接口未返回权限字段时输出安全摘要 `未返回权限字段`，不能把空态当执行失败。
- Batch A live/smoke 通过 `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-a --config configs/eino-workbench.local.yaml` 执行；通过报告写入 `test-results/eino-workbench-fobrain-batch-a-live-report.json`，schema 为 `docs/schemas/fobrain/batch_a_live_report.v1.schema.json`。

## 派生对象

实现阶段必须从该配置派生：

- `capabilities.CredentialBinding`：只包含 `schema_version`、`workspace_id`、`system=fobrain`、`status`、`display_ref`、`owner_scope`、`audit_ref`。
- `capabilities.ConnectorStatus`：只表达 available/unavailable。
- `FobrainClientConfig`：只在 provider client 内部持有 `base_url`、`timeout` 和 TLS 兼容开关；secret token 只由 `CredentialResolver` 解析成 `ResolvedCredential` 后进入单次请求。
- `auth_param`：只作为 provider client 认证参数名，当前 live client 应发送显式配置的 `authorization: <api_token>`；不要硬编码旧项目默认 header，也不要在缺失配置时自动补默认值。

## 脱敏规则

- `api_token`、Authorization、cookie、raw credential ref、连接串和 raw provider payload 不得进入 Product Facts、ActionResult、Workbench、SSE、audit、replay、report 或日志。
- RedactedSummary 只能展示 provider enabled、connector id、workspace id、base url host、timeout、credential status、auth parameter name、connector mode、available 和 TLS skip 布尔值。
- workspace 不匹配必须返回 `credential_scope_denied`，不能回显配置里的 workspace/token。

## 验收

实现必须配套配置和 provider 测试：

```bash
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
go test ./internal/einoapp/providers/fobrain -run 'HTTPClientCurrentUserContext|Credential' -count=1
go test ./cmd/eino-workbench -run 'Fobrain.*Live|FobrainOnlyWhenEnabled' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-a --config configs/eino-workbench.local.yaml
```

测试必须覆盖：最小合法配置、live mode workspace token、私有证书 TLS skip、缺必填字段、非法 URL、非法 credential status、非法 auth parameter name、secret 不出现在脱敏摘要、`enabled=false` 不注册 provider、workspace 不匹配不发起外部 HTTP。
