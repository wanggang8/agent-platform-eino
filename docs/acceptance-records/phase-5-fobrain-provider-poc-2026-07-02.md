# Phase 5 Fobrain Provider PoC 验收记录

阶段：Phase 5 Fobrain provider PoC  
日期：2026-07-02  
执行人：Codex

## 变更范围

- 新增 `internal/einoapp/providers/fobrain` provider PoC package。
- 新增 Fobrain YAML 配置解析、校验和 RedactedSummary。
- 将 Fobrain provider 按配置接入 capability registry / InvokerMux。
- 将配置派生的安全 policy context 注入 Action selection 和 ToolLoopRunner。
- 打开 `fobrain-poc` smoke，生成 mock provider PoC report；缺少 live 前置条件或 Phase 5 不支持 live 时生成 skip report。
- Phase 5 服务启动只允许 Fobrain `mock` mode；`live` mode 在真实 HTTP client 阶段前会被配置校验拒绝。

## 验证命令

| 命令 | 结果 | 说明 |
| --- | --- | --- |
| `go test ./internal/einoapp/providers/fobrain -run 'Provider\|Client\|ReadonlyPoC\|StructuredResult\|Credential\|Unsafe' -count=1` | passed | 覆盖 catalog、凭据、client、StructuredResult mapper 和 unsafe 拒绝。 |
| `go test ./internal/einoapp/bootstrap -run 'FobrainConfig\|RedactedSummary\|CredentialLeak' -count=1` | passed | 覆盖 Fobrain 配置契约和脱敏摘要。 |
| `go test ./cmd/eino-workbench -run 'Fobrain\|Capability\|Config' -count=1` | passed | 覆盖配置启用时 provider 注册，默认不内置 Fobrain。 |
| `go test ./internal/einoapp/execution -run 'DefaultPolicyContext\|ToolLoopUsesDefaultPolicyContext\|PolicyContext\|ToolLoop' -count=1` | passed | 覆盖配置派生 policy context 和 tool loop 执行。 |
| `go test ./internal/einoapp/execution -run 'DifferentWorkspace\|DefaultPolicyContext' -count=1` | passed | 覆盖请求/run workspace 与配置凭据 workspace 不一致时拒绝执行。 |
| `go test ./internal/einoapp/providers/fobrain ./internal/einoapp/bootstrap ./internal/einoapp/execution ./cmd/eino-workbench -count=1` | passed | 受影响 Go package 组合检查。 |
| `npm run eino-workbench:schema-test` | passed | `validated 31 schemas and fixture manifest`。 |
| `bash scripts/eino_workbench_server_smoke.sh --scenario mcp-mock` | passed | 回归 MCP mock smoke。 |
| `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config configs/eino-workbench.local.yaml` | passed | provider PoC report 生成；本机无 live 配置时生成 skip report。 |
| `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config configs/eino-workbench.example.yaml` | passed | 配置文件存在但缺少 LLM/Fobrain live 凭据时仍生成 skip report。 |
| `bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config <temp-live-looking-config>` | passed | LLM/Fobrain live 凭据看似齐全时仍生成 Phase 5 unsupported skip report。 |

## 报告产物

- `test-results/eino-workbench-fobrain-provider-poc-report.json`：provider PoC passed，`provider_mode` 固定为 `mock`，`result_ref` 来自 Product Facts 中实际 `tool_results`。
- `test-results/eino-workbench-fobrain-skip-report.json`：缺少本地 LLM/Fobrain live 配置、live 凭据，或 Phase 5 尚不支持 live 时，阻断真实模型工具选择和 Fobrain live read 声明。

## 未声明能力

- 不声明 24 个 Fobrain 只读工具完整恢复。
- 不声明 Fobrain live read / live write 可用。
- 不声明实体消歧、clarification、approval resume 或 ticket mutation 可用。
- 不声明可替换旧项目 Fobrain 能力。

## 风险与后续

- 当前 PoC 使用 mock Fobrain client；真实 HTTP client、24 只读矩阵和 live 证据仍归属后续 Phase。
- 真实模型工具选择和 Fobrain live read 在 Phase 5 只允许生成 skip report；即使本机配置齐全，也不能把 mock provider PoC 解释为 live 证据。
