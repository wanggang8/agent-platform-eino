# 开发前复核记录模板

阶段：
日期：
执行人：

## 外部资料

| 主题 | 链接 | 访问日期 | 结论 | 对设计影响 |
| --- | --- | --- | --- | --- |
| Eino ChatModelAgent / Runner |  |  |  |  |
| Eino HITL |  |  |  |  |
| Eino Checkpoint / Interrupt |  |  |  |  |
| Eino Callback |  |  |  |  |
| MCP lifecycle |  |  |  |  |
| MCP tools |  |  |  |  |
| OpenAPI 3.1 |  |  |  |  |
| JSON Schema 2020-12 |  |  |  |  |
| Playwright visual comparisons |  |  |  |  |
| React / Vite |  |  |  |  |

## 版本命令摘要

```bash
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go list -m -versions github.com/eino-contrib/jsonschema
```

当前版本：

可见最新稳定版本：

未引入模块只记录 `-versions` 查询结果；不得要求当前项目 `go list -m` 未引入模块。

是否升级：

不升级或升级原因：

## 当前架构复核

| 材料 | 已读 | 结论 |
| --- | --- | --- |
| `docs/ARCHITECTURE.md` |  |  |
| `docs/STATUS.md` |  |  |
| `docs/ROADMAP.md` |  |  |
| `docs/PRODUCTION_CONNECTOR_GUIDE.md` |  |  |
| `docs/design/workbench-chat-contract.md` |  |  |
| `docs/design/workbench-chat-frontend-contract.md` |  |  |
| `docs/design/workbench-ui-acceptance-matrix.md` |  |  |

## 新设计确认

- Workbench 和 Action API 共用 Product Facts：
- 工具结果只有 StructuredResult 一份事实材料：
- JSON 不暴露可复用 resume token：
- Eino event 不直接暴露给前端或外部 API：
- HITL/checkpoint/resume 场景覆盖：
- Capability Provider 不暴露 raw provider payload：
- MCP provider 语义覆盖：
- Fobrain P2 能力已映射：
- 前端视觉基线已准备：
- 未复制旧 runtime 类型或旧接口兼容层：

## 风险与处理

| 风险 | 阻断 | 处理方案 | 负责人 |
| --- | --- | --- | --- |

## 结论

- 是否允许进入本 Phase：
- 需要更新的 ADR：
- 需要更新的 schema/fixture/report：
- 不得声明的能力：
