# Phase 8 Fobrain Batch A Live 验收记录

阶段：Phase 8 Fobrain Batch A live/smoke  
日期：2026-07-02  
执行人：Codex

## 变更范围

- 新增 `fobrain-batch-a` smoke 场景，使用 ignored local config 调用真实 Fobrain provider。
- 新增 Batch A live report schema 和 fixture，覆盖 `connector.fobrain.security`、`tool.fobrain.current_user_context`、`tool.fobrain.my_permissions`。
- Fobrain HTTP client 支持标准 `/api/v1/user` 优先，404/405 时回退到私有部署兼容 `/api/user`。
- `my_permissions` 支持权限空态：当前用户接口未返回权限字段时生成安全空态摘要，不把空态当 provider 执行失败。
- 真实 token、认证 header、raw provider payload 和连接细节未写入报告。

## 验证命令

| 命令 | 结果 | 说明 |
| --- | --- | --- |
| `go test ./internal/einoapp/providers/fobrain -run 'TestHTTPClientCurrentUserContextSupportsAPIBasePath\|TestHTTPClientCurrentUserContextFallsBackToAPIUserPath\|TestHTTPClientMyPermissionsAllowsEmptyPermissionFields\|TestHTTPClientMyPermissionsUsesCurrentUserEndpoint\|TestProviderInvokesMyPermissionsThroughStructuredResult' -count=1` | passed | 覆盖 `/api` base path、fallback、权限空态和 StructuredResult。 |
| `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-a --config configs/eino-workbench.local.yaml` | passed | 使用本机 ignored live 配置生成 Batch A live report。 |
| `node scripts/eino_workbench_report_validate.mjs --schema docs/schemas/fobrain/batch_a_live_report.v1.schema.json --report test-results/eino-workbench-fobrain-batch-a-live-report.json` | passed | `test-results/eino-workbench-fobrain-batch-a-live-report.json` 通过 Batch A live report schema。 |
| `go test ./... -count=1` | passed | 全量 Go 测试通过。 |
| `npm run eino-workbench:schema-test` | passed | `validated 32 schemas and fixture manifest`。 |
| `git diff --check` | passed | 无 whitespace error。 |

## 报告产物

- `test-results/eino-workbench-fobrain-batch-a-live-report.json`：`status=passed`，三项 Batch A capability 均为 `passed`，`provider_mode=live`。
- `tool.fobrain.current_user_context`：生成当前用户安全摘要。
- `tool.fobrain.my_permissions`：生成权限空态安全摘要 `未返回权限字段`。
- `connector.fobrain.security`：生成 connector 可用性和安全绑定摘要。

## 未声明能力

- 不声明 24 个 Fobrain 只读工具完整恢复。
- 不声明 Batch B-E 可用。
- 不声明 Workbench 视觉证据已完成。
- 不声明 live write、写域审批或 ticket mutation 可用。

## 风险与后续

- 当前 Batch A live 只证明 provider live read 基础路径和报告脱敏；Workbench 前端视觉证据仍需后续补齐。
- 当前真实环境未返回权限数组，`my_permissions` 以空态通过；后续如确认权限字段来源，需要补充字段映射和 live 断言。
- `/api/user` 是当前私有环境兼容路径；标准路径仍以 `/api/v1/user` 优先。
