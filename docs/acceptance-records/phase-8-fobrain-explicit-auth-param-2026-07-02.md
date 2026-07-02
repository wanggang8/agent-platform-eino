# Phase 8 Fobrain Explicit Auth Param Acceptance Record

日期：2026-07-02

范围：Fobrain live read 认证参数名从隐式默认改为显式配置。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| 配置缺少 `credential.auth_param` 时失败 | `docs/fobrain-provider-config.md` | `TestFobrainConfigValidationRejectsInvalidFields/missing_auth_param` | `go test ./internal/einoapp/bootstrap -run FobrainConfig -count=1` | Pass | 不再由 bootstrap 默认补 `authorization`。 |
| Credential resolver 缺少 `AuthParam` 时失败 | `docs/adr/2026-07-02-fobrain-live-auth-and-batch-gates.md` | `TestStaticCredentialResolverRequiresExplicitAuthParam` | `go test ./internal/einoapp/providers/fobrain -run ExplicitAuthParam -count=1` | Pass | 返回安全 `credential_missing`，不泄漏 token。 |
| HTTP client 缺少 `AuthParam` 时不发请求 | `docs/06-security-and-projection.md` | `TestHTTPClientRequiresExplicitAuthParamBeforeRequest` | `go test ./internal/einoapp/providers/fobrain -run RequiresExplicitAuthParam -count=1` | Pass | 覆盖 Batch A、Batch D、Batch E live client 入口。 |
| 文档明确显式配置要求 | `docs/07-implementation-plan.md`、`docs/fobrain-provider-config.md`、ADR | 文档 diff | `git diff --check` | Pass | 仍要求当前真实环境配置 `auth_param: "authorization"`，但 provider 不再 hardcode 默认。 |

## 验证命令

```bash
go test ./internal/einoapp/bootstrap -run 'FobrainConfigValidationRejectsInvalidFields|FobrainConfig' -count=1
go test ./internal/einoapp/providers/fobrain -run 'ExplicitAuthParam|RequiresExplicitAuthParam|Credential|HTTPClient' -count=1
```

## 剩余风险

- 已存在的本地 live 配置必须显式包含 `fobrain.credential.auth_param: "authorization"`。
- 本记录不改变当前 workspace 共享 token 策略，也不新增 query/body 认证方式。
