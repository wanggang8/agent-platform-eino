# Eino-first 重开发文档包

本目录描述一套新的 Agent Workbench 全栈实现。文档按“产品要什么”和“技术怎么做”分层，避免把需求、架构、契约和计划写在同一个文件里。

核心原则：

```text
Eino-first, product-owned.
通用 Agent 执行能力优先交给 Eino。
产品契约、安全边界、展示事实和视觉体验由项目定义。
Runtime core 不耦合业务；业务能力通过 Skill / Tool / MCP / Capability Provider 接入。
```

完整重构目标是 P0 + P1 + P2 全部完成。P0/P1 证明新架构可行，P2 恢复既有业务能力；未完成 P2 不能声明重构完成或能力等价。

重构技术方案不是复制旧 runtime。新版本应基于 Eino-first 重新设计，去掉旧接口和旧运行时负债，同时保留既有产品能力和体验。

## 规范不是硬编码

本目录需要补齐的是新架构规范，不是旧实现迁移说明。

- 可以继承：产品能力、用户体验、安全边界、业务工具清单、验收场景。
- 不继承：旧 runtime 类型、旧执行链路、旧接口兼容层、旧 Workbench 静态结构、旧 provider raw payload。
- Fobrain 矩阵用于防止产品能力丢失，不表示复制旧代码结构。
- 字段级 schema、fixture、prompt id、mock/live 断言用于约束新系统输出，不是把旧内部模型硬编码进新实现。
- 新实现必须通过 Product Facts、StructuredResult、Safety Gate、Capability Provider 和 Eino 执行能力组合完成。

## 阅读顺序

1. [00-product-brief.md](./00-product-brief.md)：一句话说明要做什么产品，为什么重做。
2. [01-product-requirements.md](./01-product-requirements.md)：纯产品需求，不写框架、数据库、恢复机制等实现细节。
3. [02-ux-visual-requirements.md](./02-ux-visual-requirements.md)：前端视觉、交互和防跑偏标准。
4. [03-functional-scope.md](./03-functional-scope.md)：P0/P1/P2 产品能力范围。
5. [04-technical-architecture.md](./04-technical-architecture.md)：Eino-first 技术架构和模块边界。
6. [05-contract-design.md](./05-contract-design.md)：API、schema、fixtures、typed contract。
7. [06-security-and-projection.md](./06-security-and-projection.md)：StructuredResult、Safety Gate、安全投影。
8. [07-implementation-plan.md](./07-implementation-plan.md)：按依赖顺序拆分的实施任务。
9. [08-acceptance-plan.md](./08-acceptance-plan.md)：验收命令、截图、真实模型 smoke 和 live smoke。
10. [09-comparison-scorecard.md](./09-comparison-scorecard.md)：与当前实现和 B 方案对比的指标。

配套规范：

- [facts-contract.md](./facts-contract.md)：Product Facts 字段、状态、投影和幂等规则。
- [intent-and-capability-selection.md](./intent-and-capability-selection.md)：自然语言入口选择、Eino tool call、澄清和写域审批边界。
- [conversation-context.md](./conversation-context.md)：安全上下文投影、多轮引用和 context snapshot 规则。
- [run-lifecycle.md](./run-lifecycle.md)：Run、Tool、Pending 的 cancel、stop、timeout、retry 和终态幂等规则。
- [observability-and-budgets.md](./observability-and-budgets.md)：Eino callback、Product audit、telemetry 和预算门禁。
- [clarification-flow.md](./clarification-flow.md)：澄清 waiting/resume 状态机。
- [approval-flow.md](./approval-flow.md)：审批 waiting/resume 状态机。
- [capability-provider-contract.md](./capability-provider-contract.md)：工具、Skill、MCP、Connector 的新架构接入规范。
- [provider-policy-and-credentials.md](./provider-policy-and-credentials.md)：policy、permission、凭据绑定和 workspace scope。
- [fobrain-provider-config.md](./fobrain-provider-config.md)：Phase 5 Fobrain provider 本地配置、凭据和 connector 状态契约。
- [fobrain-source-api-reference.md](./fobrain-source-api-reference.md)：只读整理 `../fobrain` 源码中的接口、字段和 Batch D 样本发现流程。
- [fobrain-live-read-batch-plan.md](./fobrain-live-read-batch-plan.md)：Phase 8 Fobrain workspace 共享 token、live read 和 24 只读工具分批门禁。
- [llm-provider-safety.md](./llm-provider-safety.md)：模型 provider 配置、网络安全和错误脱敏。
- [frontend-architecture.md](./frontend-architecture.md)：Workbench 前端技术栈、状态边界、组件策略和测试策略。
- [phase-3-execution-and-facts-plan.md](./phase-3-execution-and-facts-plan.md)：Phase 3 Eino Runner、Product Facts、SQLite、SSE 和 Action API 的落地顺序。
- [pre-development-validation.md](./pre-development-validation.md)：开发前外部资料、版本和旧架构复核门禁。
- [legacy-acceptance-evidence.md](./legacy-acceptance-evidence.md)：旧项目最终验收证据索引，只读参考，不作为新项目通过结果。
- [fobrain-tool-matrix.md](./fobrain-tool-matrix.md)：P2/Phase 8 Fobrain 能力恢复矩阵。
- [tooling-and-reporting.md](./tooling-and-reporting.md)：前端工具链、视觉基线和报告规范。
- [visual-acceptance-matrix.md](./visual-acceptance-matrix.md)：Workbench block 级视觉验收矩阵。

## 文档规范

- `00-03` 只回答“我要什么”和“范围是什么”，不写框架、schema、测试命令、验收门禁或评分指标。
- `04-06` 才进入“怎么实现”。
- `07` 只写任务、文件、命令、任务级依赖和任务级检查；不得替代阶段验收或产品完成标准。
- `08` 只写验收方式、阶段门禁、报告规则和验收记录模板。
- `09` 只写对比指标、采集方式和评分记录；阶段是否完成以 `08` 为准。

ADR 放在 `docs/adr/YYYY-MM-DD-title.md`。Eino 版本升级、架构边界变化、破坏性契约变更、安全规则变化都必须新增 ADR。

ADR 模板：

```text
背景
决策
备选方案
影响
回滚条件
关联验收
```

## 外部资料

- [Eino Overview](https://www.cloudwego.io/docs/eino/overview/)
- [Eino ChatModelAgent](https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/)
- [Eino Agent HITL](https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_hitl/)
- [Eino Interrupt & CheckPoint](https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/checkpoint_interrupt/)
- [Eino Callback Manual](https://www.cloudwego.io/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/)
- [OpenAPI 3.1.2](https://spec.openapis.org/oas/v3.1.2.html)
- [JSON Schema 2020-12](https://json-schema.org/draft/2020-12/json-schema-core)
- [MCP Tools](https://modelcontextprotocol.io/specification/2025-06-18/server/tools)
- [MCP Lifecycle](https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle)
- [React Start a New Project](https://react.dev/learn/start-a-new-react-project)
- [Vite Guide](https://vite.dev/guide/)
- [Playwright Visual Comparisons](https://playwright.dev/docs/test-snapshots)
