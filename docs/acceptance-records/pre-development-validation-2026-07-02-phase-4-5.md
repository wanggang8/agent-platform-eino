# Phase 4.5 Pre-development Validation

日期：2026-07-02
目标任务：Phase 4.5 MCP adapter contract + 可运行 mock MCP provider。

## 外部资料复核

| 主题 | 访问日期 | 结论 | 对本项目影响 |
| --- | --- | --- | --- |
| MCP Lifecycle | 2026-07-02 | `initialize` 必须先进行 protocol/capability negotiation，成功后 client 发送 `notifications/initialized`，operation 阶段只能使用协商能力。 | mock MCP session 必须显式记录 initialized 状态和 tools/listChanged 能力。 |
| MCP Tools | 2026-07-02 | tools/list 支持 pagination；tool 定义包含 name/title/description/inputSchema/outputSchema/annotations；tools/call 可返回 `structuredContent` 和 `isError`；客户端应校验结果并执行安全审计。 | Phase 4.5 只实现 mock server catalog、分页工具列表、structuredContent/isError 到 StructuredResult candidate 的安全转换。 |
| Eino Overview | 2026-07-02 | Eino 仍作为通用 agent orchestration 能力，产品事实和安全投影由本项目持有。 | MCP adapter 只作为 capability provider 来源，不改变 Product Facts / StructuredResult 边界。 |

参考链接：

- `https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle`
- `https://modelcontextprotocol.io/specification/2025-06-18/server/tools`
- `https://www.cloudwego.io/docs/eino/overview/`

## 版本复核

执行命令：

```bash
go list -m github.com/cloudwego/eino
go list -m github.com/cloudwego/eino-ext/components/model/openai
go list -m github.com/eino-contrib/jsonschema
go list -m -versions github.com/cloudwego/eino
```

结果摘要：

- `github.com/cloudwego/eino` 当前项目版本为 `v0.9.12`。
- `github.com/cloudwego/eino-ext/components/model/openai` 和 `github.com/eino-contrib/jsonschema` 当前未作为已知依赖引入。
- 当前 shell 网络受限，`go list -m -versions` 无法访问 `proxy.golang.org`，沿用 2026-07-01 Phase 4 复核结论；本任务不升级 Eino，不引入新外部依赖。

## 当前代码复核

- `capabilities.Provider` 只负责能力注册，`capabilities.Invoker` 只返回 `product.StructuredResultCandidate`。
- `product.StructuredResultSafetyGate` 是工具结果进入 Product Facts 的唯一门禁。
- Phase 4.4 已实现 provider policy、CredentialBinding 安全摘要和 Action API policy blocked failed run。
- 现有 `MockProvider` 是本地只读能力 mock，不包含 MCP lifecycle、tools/list pagination、structuredContent/isError 语义。

## 设计确认

| 问题 | 结论 |
| --- | --- |
| Workbench 和 Action API 是否仍共用 Product Facts？ | 是，本任务只新增 provider adapter，不改变 facts repository 或 product projection。 |
| MCP `structuredContent` 是否能直接进入产品输出？ | 否，只能转为 StructuredResult candidate，并通过 Safety Gate。 |
| MCP annotations 是否可信？ | 否，只能作为与项目 policy 合并的输入，不能覆盖本项目写域审批和风险策略。 |
| Phase 4.5 是否接生产 MCP server？ | 否，只接可运行 mock provider；生产 lifecycle、auth、多 server catalog 和 live smoke 后移。 |
| 是否需要复制旧项目 runtime 或 Workbench 接口？ | 否。 |

## 结论

- 是否允许进入 Phase 4.5：是。
- 是否允许接生产 MCP server：否。
- 是否允许绕过 StructuredResult Safety Gate：否。
