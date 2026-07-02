# Fobrain Live 认证与分批门禁

## 背景

Phase 5 已完成 Fobrain provider PoC，但仍只允许 mock connector mode。进入真实 Fobrain 只读工具恢复前，需要固定认证参数、凭据作用域和 24 个只读工具的分批验收方式。

旧项目可作为产品能力和验收基线参考；旧实现中的默认 header 名称不作为新项目契约来源。当前真实环境的认证参数名由新项目配置显式声明。

## 决策

- Fobrain live client 使用配置文件中的 `fobrain.credential.auth_param` 作为认证参数名，当前目标值为 `authorization`。
- `auth_param` 只是参数名；真实 token 只允许写入已忽略的 `configs/eino-workbench.local.yaml` 的 `fobrain.credential.api_token`。
- 当前阶段只支持 workspace 共享 token，不实现 per-user、per-tool 或 OAuth delegation。
- workspace 请求必须与凭据绑定的 `workspace_id` 一致；不一致返回 `credential_scope_denied`。
- live HTTP client 只接受显式配置的 `auth_param`，并把它作为请求 header 名发送原始 token 值，不添加 `Bearer` 前缀；缺失时必须阻断请求，不能在 provider 内默认成 `authorization`。
- token、raw provider payload、Authorization 值不得进入 Product Facts、StructuredResult、Workbench、Action API、audit、replay、日志或验收报告。
- 24 个只读工具按批次恢复；每批必须先通过 mock contract，再生成 live pass report。无 live 环境时只能生成 blocking skip report，最终 24/24 + connector 全部 live pass 后才可声明 Fobrain 产品能力恢复。

## 备选方案

- 每用户 token：更接近细粒度权限，但当前没有明确绑定和轮换契约，会拖慢只读能力恢复。
- 继续沿用旧项目默认 header 名称：会把旧实现细节带入新架构，且与当前真实环境不一致。
- 一次性实现 24 个工具：验收面太大，容易丢失 StructuredResult、事实同源和可视证据约束。

## 影响

- `configs/eino-workbench.example.yaml` 只能展示无 secret 示例。
- live read 实现必须在 provider 边界内完成认证、HTTP 调用、raw payload 解析和脱敏。
- Product 层只能消费 StructuredResult，不允许读取 raw Fobrain response。
- 工具矩阵必须用 `batch_gate` 字段固化 Batch A-E，避免越批实现或越批验收。
- 后续如果确认 Fobrain 要求 query/body 参数而不是 header，必须新增 ADR 并同步配置文档、测试和 live smoke。

## 回滚条件

- 真实 Fobrain 服务确认不接受 `authorization` 参数名。
- workspace 共享 token 无法满足最小只读验收。
- live client 出现 token 或 raw payload 泄漏。

## 关联验收

- `docs/fobrain-live-read-batch-plan.md`
- `docs/fobrain-provider-config.md`
- `docs/fobrain-tool-matrix.md`
- `docs/07-implementation-plan.md` Phase 8
