# Provider Policy 与凭据契约

本文定义 Capability Provider、Connector Provider、权限策略和凭据绑定的产品契约。它只固定新架构行为，不复制旧实现结构。

## Policy 元数据

每个 capability 必须声明：

| 字段 | 说明 |
| --- | --- |
| `policy_ref` | 稳定策略引用，如 `policy:fobrain:read:v1` |
| `permission_scope` | `workspace` / `caller` / `system` |
| `side_effect` | `none` / `read_external` / `write_external` / `local_runtime` |
| `risk_level` | `none` / `low` / `medium` / `high` |
| `approval_required` | 写域或高风险必须为 true |
| `idempotency_required` | mutation 必须为 true |
| `credential_binding_policy` | `none` / `required` / `optional` |

## Permission Decision

机器契约见 `docs/schemas/provider_policy_decision.v1.schema.json`。

Provider 或 policy checker 返回统一决策：

| 字段 | 说明 |
| --- | --- |
| `allowed` | 是否允许执行 |
| `reason_code` | 稳定原因，不作为主聊天正文 |
| `policy_ref` | 命中的策略 |
| `audit_ref` | 安全审计引用 |
| `approval_required` | 是否需要审批后继续 |

允许的 reason code 至少包括：

```text
permission_denied
approval_required
credential_missing
credential_scope_denied
connector_auth_failure
connector_timeout
connector_transport_unavailable
connector_execution_failed
resource_conflict
schema_mismatch
```

## 凭据绑定

机器契约见 `docs/schemas/provider_credential_binding.v1.schema.json`。

真实凭据只在服务端解析。产品输出只能展示：

```text
credential_binding: configured | missing | unbound | bound:<system>:<binding>
```

禁止在 Workbench、ActionResult、SSE、replay、audit、catalog、registry projection 中展示：

- `credential_ref`
- `credential:<...>`
- API key / token / secret
- Authorization header
- raw connector config

## Workspace Scope

- credential lookup 必须带 `workspace_id`。
- workspace 不匹配返回 `credential_scope_denied`。
- connector status 可以展示绑定状态，但不能替代业务读工具。
- 凭据更新必须写 audit event，且 audit event 只保留安全摘要。

## Fobrain Policy

| 能力 | policy_ref | side_effect | approval |
| --- | --- | --- | --- |
| 24 个只读工具 | `policy:fobrain:read:v1` | `read_external` | 否 |
| `connector.fobrain.security` | `policy:fobrain:connector-read:v1` | `read_external` | 否 |
| `tool.fobrain.query_normalize_and_resolve` | `policy:fobrain:read:v1` | `read_external` | 否 |
| `tool.fobrain.update_ticket_status` | `policy:fobrain:ticket-mutation:v1` | `write_external` | 是 |

## 验收

- 缺少凭据时返回 connector 状态卡或安全错误，不泄漏内部 ref。
- workspace 不匹配返回 `credential_scope_denied`。
- 写域未审批前不执行 mutation。
- approve 后同一 idempotency key 最多执行一次 mutation。
- catalog、ActionResult、SSE、replay、audit 全部通过 credential leak regression。
