# ADR: Fobrain Provider PoC Boundary

状态：Accepted
日期：2026-07-02

## 背景

Phase 4 已完成 StructuredResult Safety Gate、capability registry、provider policy、mock MCP provider 和 tool loop 基础。Phase 5 需要开始接入 Fobrain，但旧项目 Fobrain 包含大量 domain、adapter、providerapi、runtime 和 Workbench 绑定。如果直接照搬旧结构，会破坏 Eino-first 重构目标。

## 决策

- Phase 5 只实现一个只读 Fobrain provider PoC capability，默认选择 `tool.fobrain.current_user_context`。
- Fobrain provider 必须实现新项目的 `capabilities.Provider` 与 `capabilities.Invoker`，不得复制旧 `providerapi.Provider` 或旧 runtime handler。
- capability 元数据全部集中在 `providers/fobrain/catalog.go`，execution、httpapi、product 不按 Fobrain 工具名硬编码分支。
- Fobrain connector 配置必须先按 `docs/fobrain-provider-config.md` 实现，base URL、workspace、timeout、connector 状态和本地 token 不得散落在业务代码或 smoke 脚本中。
- provider client 只在边界内处理 raw HTTP payload；进入产品层前必须转成 typed safe result，再转 `product.StructuredResultCandidate`。
- Product Facts 的唯一工具事实 schema 仍是 `tool.structured_result.v1`；`fobrain.tool_result.v2` 只作为 Fobrain 业务 payload/展示 schema，通过安全 result ref 关联。
- Workbench、Action API、audit、replay 继续只从 Product Facts 投影，不直接读取 provider result。
- `fobrain-poc` 先生成 provider PoC report；真实模型工具选择报告只在 LLM/Fobrain 凭据齐全时生成，缺凭据时使用 skip report 阻断对应声明。
- 写域 capability 只能保留在文档和矩阵中，Phase 6 approval/resume 完成前不得执行 mutation。

## 外部依据

- 2026-07-02 复核 CloudWeGo Eino 官方 ToolsNode 文档：tool 是模型可选择调用的外部能力，执行由工具节点负责。
- 2026-07-02 复核 CloudWeGo Eino ChatModelAgent 文档：配置 tools 后会进入 ReAct 工具循环。
- 2026-07-02 复核 MCP 2025-06-18 tools 规范：structuredContent 是工具结构化输出；本项目仍要求它先转 StructuredResult candidate，再经 Safety Gate。

## 后果

- Phase 5 的实现范围较小，便于验证 provider 边界、凭据脱敏、StructuredResult 和 Product Facts 同源链路。
- 24 个只读工具、connector 展示、实体消歧和写域审批继续留在 Phase 8/Phase 6 对应任务中。
- 后续新增 Fobrain 工具必须先扩展 schema/fixture/catalog 测试，再接 provider client operation。

## 验证

```bash
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
go test ./internal/einoapp/providers/fobrain -run 'Provider|Client|ReadonlyPoC|StructuredResult|Credential' -count=1
go test ./internal/einoapp/capabilities ./internal/einoapp/product ./internal/einoapp/execution -run 'Policy|StructuredResult|ToolSafety|ToolFacts' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config configs/eino-workbench.local.yaml
```
