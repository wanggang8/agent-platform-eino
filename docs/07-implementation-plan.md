# 实施计划

本文只写任务、依赖、文件和任务级检查。产品需求见 `01-product-requirements.md`，技术设计见 `04-technical-architecture.md`。

权威阶段门禁在 `08-acceptance-plan.md`。本文中的命令用于任务完成时快速检查，不替代阶段、合并或完整重构验收。

## Phase 0：文档整理

目标：建立分层文档，避免需求和技术设计混写。

前置门禁：

- 执行 `pre-development-validation.md`，形成 `acceptance-records/pre-development-validation-YYYY-MM-DD.md`。
- 重新调研外部官方资料。
- 重新核对当前项目架构和关键代码边界。
- 确认新设计不需要复制旧 runtime 类型、旧执行链路或旧接口兼容层。

任务级检查：

```bash
find docs -maxdepth 1 -type f | sort
! rg -n "Eino|CheckPointStore|ResumeWithParams|OpenAPI|JSON Schema|MCP|Vite|React" docs/00-product-brief.md docs/01-product-requirements.md docs/02-ux-visual-requirements.md
test -f docs/pre-development-validation.md
```

预期：文件结构存在；产品文档不出现实现细节关键词。

## Phase 1：契约和服务骨架

### Task 1.1 Tooling bootstrap

创建：

- `package.json`
- `web/eino-workbench/package.json`
- `web/eino-workbench/tsconfig.json`
- `web/eino-workbench/src/contracts/.gitkeep`
- `docs/tooling-and-reporting.md`（完善/固化）

要求：

- 根 `package.json` 提供 `eino-workbench:*` wrapper。
- 前端包声明 Node 版本，并提交 lockfile。
- Playwright config 必须包含 webServer、baseURL、reporter、desktop/mobile viewport、截图阈值和动态内容 mask。
- scripts 至少包含：`schema-test`、`contract-generate`、`contract-test`、`typecheck`、`test`、`stream-test`、`browser-test`、`visual-test`、`build`、`send-smoke`。
- 安装命令包含 `npm ci` 和浏览器安装命令。

任务级检查：

```bash
npm ci
npx playwright install --with-deps
npm run eino-workbench:contract-test -- --help
```

### Task 1.2 OpenAPI、schema、fixtures

创建：

- `docs/api/eino-workbench.openapi.json`
- `docs/schemas/*.schema.json`
- `docs/fixtures/*.json`
- `docs/facts-contract.md`（完善/固化）
- `docs/clarification-flow.md`（完善/固化）
- `docs/approval-flow.md`（完善/固化）
- `docs/intent-and-capability-selection.md`（完善/固化）
- `docs/conversation-context.md`（完善/固化）
- `docs/run-lifecycle.md`（完善/固化）
- `docs/observability-and-budgets.md`（完善/固化）
- `docs/capability-provider-contract.md`（完善/固化）
- `docs/provider-policy-and-credentials.md`（完善/固化）
- `docs/llm-provider-safety.md`（完善/固化）
- `scripts/eino_workbench_schema_validate.mjs`

任务级检查：

```bash
npm run eino-workbench:schema-test
```

### Task 1.3 TypeScript contract

创建：

- `web/eino-workbench/src/contracts/generated.ts`
- `web/eino-workbench/src/contracts/contract.test.ts`

任务级检查：

```bash
npm run eino-workbench:contract-generate
npm run eino-workbench:contract-test
```

### Task 1.4 Package skeleton and import boundary

创建：

- `internal/einoapp/architecture/import_boundary_test.go`
- `internal/einoapp/facts/.gitkeep` 或实际 package 文件。
- `internal/einoapp/execution/.gitkeep` 或实际 package 文件。
- `internal/einoapp/product/.gitkeep` 或实际 package 文件。
- `internal/einoapp/capabilities/.gitkeep` 或实际 package 文件。
- `internal/einoapp/httpapi/.gitkeep` 或实际 package 文件。
- `internal/einoapp/store/sqlite/.gitkeep` 或实际 package 文件。
- `internal/einoapp/llm/.gitkeep` 或实际 package 文件。
- `internal/einoapp/observability/.gitkeep` 或实际 package 文件。
- `internal/einoapp/providers/fobrain/.gitkeep` 或实际 package 文件。

要求：

- 先固定目录角色和依赖方向，再写 HTTP 或执行层代码。
- `httpapi` 只能依赖 product/projection、facts 查询接口和 execution command 接口。
- `execution` 可以写 facts，但不得依赖 `httpapi` 或前端契约实现。
- `capabilities`、`llm`、`providers` 不得向外暴露 raw provider payload。
- `providers/fobrain` 不得被 runtime core 反向依赖。

任务级检查：

```bash
go test ./internal/einoapp/... -run 'ImportBoundary|PhaseOnePackageSkeleton' -count=1
```

### Task 1.5 HTTP skeleton

创建：

- `cmd/eino-workbench/main.go`
- `internal/einoapp/bootstrap/config.go`
- `internal/einoapp/httpapi/routes.go`
- `internal/einoapp/httpapi/actions.go`
- `scripts/eino_workbench_dev_server.sh`
- `scripts/eino_workbench_server_smoke.sh`

任务级检查：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario contract
```

### Task 1.6 Backend foundation hardening

目标：在进入 stream reducer、Eino chat 或 Action API 实现前，先固定后端基础设施，避免把配置、HTTP 响应、错误脱敏、LLM provider、能力注册和事实投影散落到业务 handler 中。

创建/修改：

- `internal/einoapp/bootstrap/config.go`
- `internal/einoapp/bootstrap/config_test.go`
- `internal/einoapp/httpapi/response.go`
- `internal/einoapp/httpapi/response_test.go`
- `internal/einoapp/httpapi/sse.go`
- `internal/einoapp/httpapi/sse_test.go`
- `internal/einoapp/llm/config.go`
- `internal/einoapp/llm/provider.go`
- `internal/einoapp/llm/redacted_error.go`
- `internal/einoapp/llm/provider_test.go`
- `internal/einoapp/capabilities/provider.go`
- `internal/einoapp/capabilities/registry.go`
- `internal/einoapp/capabilities/policy.go`
- `internal/einoapp/capabilities/registry_test.go`
- `internal/einoapp/facts/repository.go`
- `internal/einoapp/facts/model.go`
- `internal/einoapp/product/projection.go`
- `internal/einoapp/product/errors.go`
- `docs/adr/YYYY-MM-DD-backend-foundation-boundaries.md`

要求：

- 配置必须集中读取 server、database、LLM、security、observability、timeout、budget；配置日志必须脱敏。
- 成功响应保持 OpenAPI 中定义的业务 schema 直出；错误响应统一 `eino_error_envelope.v1`，并包含 request id、安全错误码和可展示摘要。
- SSE 编码必须统一处理 `id`、`event`、`data`、flush、content type、no-cache 和编码失败。
- `httpapi` 不得手写业务 DTO 第二套真相；只能调用 product projection、facts query 和 execution command 接口。
- `llm` 只暴露 provider interface、config、mock provider、redacted error；真实 OpenAI-compatible provider 可在 Phase 4 扩展。
- `capabilities` 必须先提供 registry、provider interface、tool metadata、risk/approval policy；任何工具选择都必须通过 registry，不得写死工具名或自然语言关键词。
- `facts` 先定义最小 repository interface 和事实模型；Product Facts 的 SQLite 实现仍在 Phase 3.2 完成。
- `product` 先定义 Workbench、ActionResult、Replay、SSE projection 的接口边界；实际完整映射仍在 Phase 3 完成。
- 不创建通用 `utils` 包；公共能力按职责放入 `bootstrap`、`httpapi`、`llm`、`capabilities`、`facts`、`product`、`observability`。

任务级检查：

```bash
go test ./internal/einoapp/bootstrap ./internal/einoapp/httpapi ./internal/einoapp/llm ./internal/einoapp/capabilities ./internal/einoapp/facts ./internal/einoapp/product -run 'Config|Response|SSE|Provider|Registry|Facts|Projection|Redaction' -count=1
go test ./internal/einoapp/architecture -run ImportBoundary -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario contract
```

## Phase 2：React Workbench

### Task 2.0 Frontend architecture decision

创建/修改：

- `docs/frontend-architecture.md`
- `docs/07-implementation-plan.md`
- `AGENTS.md`
- `web/eino-workbench/package.json`
- `web/eino-workbench/vite.config.ts`

要求：

- 固定前端基础架构：Vite + React + TypeScript + React Router + TanStack Query + Zustand + Radix UI Primitives + lucide-react。
- 固定样式策略：CSS Modules 或分层普通 CSS；不默认使用 Tailwind、shadcn 或重型视觉组件库。
- 固定状态边界：TanStack Query 只处理 server state，Zustand 只处理 UI/client state，Product Facts 仍只在后端。
- 固定 contract 边界：前端只从 `web/eino-workbench/src/contracts/generated.ts` 消费类型，不手写平行 DTO。
- 固定测试策略：Vitest + Testing Library + Playwright。
- 后续如替换构建框架、状态分层、组件基础或样式体系，必须新增 ADR。

任务级检查：

```bash
rg -n "Vite|React Router|TanStack Query|Zustand|Radix|Product Facts" docs/frontend-architecture.md docs/07-implementation-plan.md AGENTS.md
npm run eino-workbench:typecheck
```

### Task 2.1 React/Vite shell

创建：

- `web/eino-workbench/vite.config.ts`
- `web/eino-workbench/index.html`
- `web/eino-workbench/src/main.tsx`
- `web/eino-workbench/src/app/App.tsx`
- `web/eino-workbench/src/app/router.tsx`
- `web/eino-workbench/src/app/queryClient.ts`
- `web/eino-workbench/src/fixtures/workbenchFixtures.ts`
- `web/eino-workbench/src/features/workbench/api/workbenchQueries.ts`
- `web/eino-workbench/src/features/workbench/state/useWorkbenchUiStore.ts`
- `web/eino-workbench/src/features/workbench/views/WorkbenchPage.tsx`
- `web/eino-workbench/src/features/workbench/components/*.tsx`
- `web/eino-workbench/src/styles/*.css`
- `web/eino-workbench/src/test/*.tsx`
- `web/eino-workbench/tests/shell.spec.ts`
- `docs/visual-acceptance-matrix.md`（补 selector 和状态覆盖）

要求：

- 只使用 `docs/fixtures/` 和 `src/contracts/generated.ts` 渲染 fixture-only Workbench。
- Desktop 保持左侧会话、中间时间线、右侧 Inspector。
- Mobile 优先主聊天和输入区，Inspector 通过 panel/tabs 进入。
- 组件按 `timeline.kind`、`status`、`structured_result.display_type` 渲染，不得按 `tool_id` 或旧 DOM 分支。
- TanStack Query 只读取 fixture/server state；Zustand 只保存 UI 选择、Inspector tab、卡片展开和 mobile panel。

任务级检查：

```bash
npm run eino-workbench:build
npm run eino-workbench:browser-test -- --grep shell
```

### Task 2.2 Legacy visual evidence mapping

创建/修改：

- `docs/visual-acceptance-matrix.md`
- `docs/fixtures/visual-evidence-matrix.json`
- `web/eino-workbench/tests/visual-fixtures/*.json`

要求：

- 从 `legacy-acceptance-evidence.md` 读取旧项目最终验收截图和 contact sheet 路径。
- 为 `shell`、`sidebar`、`timeline`、`composer`、`tool-card`、`approval-card`、`clarification-card`、`inspector` 补目标截图来源、block crop、desktop/mobile viewport 和状态 fixture。
- 旧截图只能作为目标参考；不得复制旧 DOM、旧 CSS 类名、旧 presenter 分支或旧接口字段。
- 每个 block 必须声明对应状态、selector、fixture 和缺失时的阻断条件。

任务级检查：

```bash
npm run eino-workbench:schema-test
rg -n "legacy-acceptance-evidence|target|crop|selector|blocker" docs/visual-acceptance-matrix.md docs/fixtures/visual-evidence-matrix.json web/eino-workbench/tests/visual-fixtures
```

### Task 2.3 Stream reducer

创建：

- `web/eino-workbench/src/api/transport.ts`
- `web/eino-workbench/src/api/stream.ts`
- `web/eino-workbench/src/state/workbenchReducer.ts`

任务级检查：

```bash
npm run eino-workbench:stream-test
```

### Task 2.4 Visual baseline

创建：

- `web/eino-workbench/tests/visual.spec.ts`
- `web/eino-workbench/tests/__screenshots__/`
- `test-results/eino-workbench-visual-report/`

要求：

- Playwright 覆盖 desktop 1440x900 和 mobile 390x844。
- 覆盖空态、聊天态、工具卡、审批卡、澄清卡、失败态、Inspector。
- 动态字段必须 mask：时间、run id、随机 id、模型耗时、token 计数。
- 新项目截图必须重新生成，不能引用旧项目截图作为通过结果。

任务级检查：

```bash
npm run eino-workbench:visual-test
```

## Phase 3：Eino chat、facts、Action API

### Task 3.1 Pin Eino versions

修改：

- `go.mod`
- `go.sum`
- `docs/adr/YYYY-MM-DD-pin-eino-version.md`

任务级检查：

```bash
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go test ./internal/einoapp/execution -run 'EinoVersion|ChatModelVersion' -count=1
```

记录：

- 固定版本写入 ADR。
- `eino-ext` 模型 provider 只在实际引入 provider 代码时固定到 `go.mod`，不得要求未引入模块通过 `go list -m`。
- 升级版本必须新增 ADR，并重跑 Chat、Runner、tool loop 门禁；HITL、checkpoint 门禁在 Phase 6 引入后纳入升级回归。

### Task 3.2 Product facts

创建：

- `internal/einoapp/facts/run.go`
- `internal/einoapp/facts/turn.go`
- `internal/einoapp/facts/tool.go`
- `internal/einoapp/facts/pending.go`
- `internal/einoapp/facts/audit.go`
- `internal/einoapp/facts/context.go`
- `internal/einoapp/facts/lifecycle.go`
- `internal/einoapp/store/sqlite/repository.go`
- `internal/einoapp/store/sqlite/migrations.go`

任务级检查：

```bash
go test ./internal/einoapp/store/sqlite -run 'RunTurnEvent|ToolCallResult|ContextSnapshot|LifecycleTransition' -count=1
```

### Task 3.3 ChatModel and Runner

创建：

- `internal/einoapp/execution/agent.go`
- `internal/einoapp/execution/runner.go`
- `internal/einoapp/execution/events.go`
- `internal/einoapp/execution/selection.go`

要求：

- Runner event 只能写入 Product Facts，不能直接暴露给 Workbench、ActionResult 或 replay。
- assistant delta、tool call、tool result、run state 必须先落入 facts，再由投影层读取。
- 普通自然语言默认进入 ChatModelAgent；不得按关键词硬编码 Fobrain 工具路由。
- capability_hint / intent_hint 只作为候选入口，必须通过 registry 和 policy。
- 每次模型调用前必须使用 `conversation-context.md` 生成 safe context snapshot。
- 依赖 Phase 1.6 的 `llm.Provider`、`capabilities.Registry`、`facts.Repository` 和 `product.Projection` 接口，不得在 execution 中重新定义平行接口。

任务级检查：

```bash
go test ./internal/einoapp/execution -run 'Chat|RunnerEvent|FactsWrite|CapabilitySelection|ContextProjection' -count=1
```

### Task 3.4 Messages、SSE、Action API

修改：

- `internal/einoapp/httpapi/routes.go`
- `internal/einoapp/httpapi/sse.go`
- `internal/einoapp/httpapi/actions.go`

任务级检查：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario chat-stream
bash scripts/eino_workbench_server_smoke.sh --scenario action-basic
```

## Phase 4：Tools、Provider、Safety

### Task 4.1 StructuredResult and Safety Gate

创建：

- `internal/einoapp/product/structured_result.go`
- `internal/einoapp/product/safety.go`
- `internal/einoapp/product/assistant_safety.go`

任务级检查：

```bash
go test ./internal/einoapp/product -run 'StructuredResult|Safety|AssistantSafety' -count=1
```

### Task 4.2 Production LLM provider implementation

创建：

- `internal/einoapp/llm/openai_compatible.go`
- `internal/einoapp/llm/network_policy.go`
- `internal/einoapp/llm/openai_compatible_test.go`

范围：

- 基于 Phase 1.6 的 `llm.Provider` 和 `llm.Config` 接入 OpenAI-compatible ChatModel。
- 禁止把 API key、Authorization、raw provider body 写入 Product Facts、Workbench、ActionResult、audit 或 replay。
- provider non-2xx、timeout、malformed response 和 network blocked 必须转换为安全错误。

任务级检查：

```bash
go test ./internal/einoapp/llm -run 'OpenAICompatible|NetworkPolicy|RedactedError' -count=1
```

### Task 4.3 Capability adapters and Eino tool conversion

创建：

- `internal/einoapp/capabilities/tool_adapter.go`
- `internal/einoapp/capabilities/provider_contract_test.go`

范围：

- Local Tool Provider：P0/P1 必须实现。
- Connector Provider：P1 通过 Fobrain PoC 验证 metadata、risk、timeout、display name。
- Skill Provider：P1 可只保留接口，实际 skill registry 接入进入 P2。
- 所有 provider 原始结果都必须先转 StructuredResult candidate，再经过 Safety Gate。
- Capability description、input schema、risk、side effect 和 approval metadata 必须足以支撑 Eino native tool selection。

任务级检查：

```bash
go test ./internal/einoapp/capabilities -run 'ProviderContract|StructuredResultConversion|ToolSelectionMetadata|EinoToolAdapter' -count=1
```

### Task 4.4 MCP adapter

创建：

- `internal/einoapp/capabilities/mcp_adapter.go`
- `internal/einoapp/capabilities/mcp_session.go`
- `internal/einoapp/capabilities/mcp_schema.go`
- `internal/einoapp/capabilities/mcp_mock_test.go`

范围：

- P0 只完成 mock MCP adapter contract，不接生产 MCP server。
- P1 可接入受控 MCP server 做 adapter smoke。
- 生产 MCP server、复杂 auth 和多 server catalog 可在 P2 扩展。

契约要求：

- 先满足 `docs/capability-provider-contract.md`，再接入具体业务 provider。
- provider 原始结果不得直接进入 Workbench、ActionResult、audit 或 replay。
- MCP 输出必须先转 StructuredResult candidate，再经过 Safety Gate。

MCP adapter 任务拆分：

- transport / server catalog config。
- lifecycle / initialize。
- session auth / reconnect / shutdown。
- tools/list pagination。
- listChanged reload。
- inputSchema conversion。
- outputSchema / structuredContent / isError conversion。
- Safety Gate integration。
- annotations + project risk policy merge。
- minimal mock MCP server。

任务级检查：

```bash
go test ./internal/einoapp/capabilities -run 'ProviderRegistry|MCPLifecycle|MCPToolListPagination|MCPListChanged|MCPSchemaConversion|MCPSafety|MCPRiskPolicy' -count=1
```

### Task 4.5 Provider policy and credentials

创建：

- `internal/einoapp/capabilities/policy.go`
- `internal/einoapp/capabilities/credentials.go`
- `internal/einoapp/capabilities/provider_policy_test.go`

范围：

- risk、side_effect、workspace scope、permission、credential binding、approval_required 决策。
- credential missing、scope denied、connector unavailable 必须产生安全 ActionResult / Workbench 状态。
- 任何 `credential_ref`、secret 或 token hash 不得进入 JSON 输出。

任务级检查：

```bash
go test ./internal/einoapp/capabilities -run 'ProviderPolicy|CredentialBinding|CredentialLeak' -count=1
```

### Task 4.6 Eino tool loop

修改：

- `internal/einoapp/execution/tools.go`
- `internal/einoapp/execution/runner.go`

任务级检查：

```bash
go test ./internal/einoapp/execution -run AgentToolLoop -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario tool-card
```

## Phase 5：Fobrain provider PoC

创建：

- `internal/einoapp/providers/fobrain/provider.go`
- `internal/einoapp/providers/fobrain/client.go`
- `internal/einoapp/providers/fobrain/credentials.go`
- `internal/einoapp/providers/fobrain/tools.go`
- `internal/einoapp/providers/fobrain/result.go`
- `scripts/eino_workbench_real_model_tool_suite.mjs`

任务级检查：

```bash
go test ./internal/einoapp/providers/fobrain -run 'Provider|Client|ReadonlyPoC' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc
```

## Phase 6：HITL approval / clarification / resume

### Task 6.1 Checkpoint store

创建：

- `internal/einoapp/store/sqlite/checkpoint_store.go`
- `internal/einoapp/execution/checkpoint.go`

覆盖：

- waiting 前保存 checkpoint。
- checkpoint 丢失返回安全错误。
- 进程重启后可恢复。

任务级检查：

```bash
go test ./internal/einoapp/store/sqlite -run CheckPointStore -count=1
go test ./internal/einoapp/execution -run 'ResumeCheckpointMissing|ResumeAfterRestart' -count=1
```

### Task 6.2 Approval interrupt and resume

创建/修改：

- `internal/einoapp/execution/hitl.go`
- `internal/einoapp/execution/approval_tool.go`

覆盖：

- approval requested / approved / rejected / expired。
- 未审批不执行 mutation。
- approve 后继续执行。
- reject 后不可 approve。
- duplicate resume。
- pending 状态进入 audit/replay。

任务级检查：

```bash
go test ./internal/einoapp/execution -run 'ApprovalInterrupt|ResumeDuplicate|ResumeReject|ResumeAfterRestart' -count=1
```

### Task 6.3 Clarification interrupt and resume

创建/修改：

- `internal/einoapp/execution/clarification_tool.go`

覆盖：

- clarification requested / submitted / cancelled / expired。
- 候选消歧和自由文本澄清都必须转换为安全 resume data。
- duplicate submit 幂等。
- pending 状态进入 audit/replay。

任务级检查：

```bash
go test ./internal/einoapp/execution -run 'ClarificationRequested|ClarificationSubmit|ClarificationCancel|ClarificationDuplicate|ClarificationAfterRestart' -count=1
```

### Task 6.4 Pending UI and stream patches

创建/修改：

- `web/eino-workbench/src/components/ApprovalCard.tsx`
- `web/eino-workbench/src/components/ClarificationCard.tsx`
- `web/eino-workbench/src/state/workbenchReducer.ts`
- `internal/einoapp/httpapi/sse.go`

覆盖：

- approval 和 clarification waiting 卡片。
- submitted / approved / rejected / cancelled / expired 只读终态。
- SSE pending patch、ActionResult waiting 和 replay/audit 同源。

任务级检查：

```bash
npm run eino-workbench:stream-test -- --grep pending
bash scripts/eino_workbench_server_smoke.sh --scenario clarification
```

### Task 6.5 Run lifecycle controls

创建/修改：

- `internal/einoapp/execution/lifecycle.go`
- `internal/einoapp/httpapi/runs.go`
- `web/eino-workbench/src/components/RunNotice.tsx`

覆盖：

- cancel waiting/running run。
- stop running run。
- provider timeout 和 pending timeout 的安全状态迁移。
- retry 只按 `run-lifecycle.md` 中的安全策略允许。
- 终态操作幂等。

任务级检查：

```bash
go test ./internal/einoapp/execution -run 'RunCancel|RunStop|RunTimeout|RunRetry|TerminalIdempotency' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario run-lifecycle
```

## Phase 7：Projection / Audit / Replay

### Task 7.1 Product projection

创建：

- `internal/einoapp/product/projection.go`
- `internal/einoapp/httpapi/views.go`

覆盖：

- Workbench view/current、ActionResult、SSE patch、Replay View、Inspector 从 Product Facts 同源投影。
- 同一 run 的 assistant、tool、pending、audit 字段一致。
- 投影层不得读取 raw provider payload 或 Eino raw event。

任务级检查：

```bash
go test ./internal/einoapp/product -run 'Projection|Inspector' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario action-consistency
```

### Task 7.2 Audit callbacks, observability, and budgets

创建/修改：

- `internal/einoapp/facts/audit.go`
- `internal/einoapp/execution/callbacks.go`
- `internal/einoapp/execution/budgets.go`
- `internal/einoapp/observability/telemetry.go`

覆盖：

- run、turn、tool、pending、resume、provider policy 决策写入安全 audit event。
- audit event 不包含 raw prompt、secret、Authorization、raw provider body 或 reusable resume token。
- callback 只进入 telemetry/internal diagnostics，不驱动主 SSE。
- model call、tool call、run duration、context token 预算超限后安全停止或返回 partial result。

任务级检查：

```bash
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget' -count=1
```

### Task 7.3 Replay API

创建/修改：

- `internal/einoapp/httpapi/replay.go`

覆盖：

- replay 从 Product Facts 重建，不依赖前端状态。
- approval、clarification、tool result、失败态、cancel/timeout 都可回放。
- replay 输出必须通过 `eino_replay_view.v1` schema。

任务级检查：

```bash
go test ./internal/einoapp/product -run Replay -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario replay
```

## Phase 8：Fobrain 产品能力恢复

Phase 8 是 P2 门禁，属于完整重构必做范围；未完成本阶段不得声明重构完成、能力等价或替换当前产品基线。

### Task 8.1 24 个只读工具矩阵

创建/修改：

- `docs/fobrain-tool-matrix.md`
- `docs/fixtures/fobrain/*.json`
- `docs/schemas/fobrain/*.schema.json`
- `internal/einoapp/providers/fobrain/tools.go`
- `internal/einoapp/providers/fobrain/result.go`

前置要求：

- 24 个只读工具必须先补齐 `fobrain-tool-matrix.md` 的字段级恢复模板。
- 每个工具必须有 input schema、result schema、fixture、mock assertion、live assertion 和 screenshot state。
- 每个工具的 mock/live 断言必须覆盖 StructuredResult `display_type`、safe summary、evidence、pagination/empty/error 语义和敏感字段脱敏。
- 未补齐矩阵的工具不得进入实现。

任务级检查：

```bash
go test ./internal/einoapp/providers/fobrain -run ReadonlyToolMatrix -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-readonly
```

### Task 8.2 Fobrain presentation evidence matrix

创建/修改：

- `docs/fobrain-tool-matrix.md`
- `docs/fixtures/fobrain/tool-matrix-24.json`
- `web/eino-workbench/tests/fobrain-tool-visual.spec.ts`
- `test-results/eino-workbench-fobrain-tool-acceptance/`

要求：

- 每个 Fobrain 场景必须保留旧最终验收中的六类截图区域要求：`main-chat`、`fresh-main-chat`、`process`、`evidence`、`audit`、`internal-details`。
- 工具矩阵必须覆盖 24 个只读工具；connector 场景在 Task 8.3 覆盖。
- 截图证据必须由新项目运行生成，旧项目截图只作为 coverage 参考。
- 每个场景必须记录 prompt、selected tool、sanitized args summary、tool card rendering status、assistant final answer、failure category。

任务级检查：

```bash
npm run eino-workbench:visual-test -- --grep fobrain
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-readonly
```

### Task 8.3 Connector 状态与凭据绑定

创建/修改：

- `internal/einoapp/providers/fobrain/provider.go`
- `internal/einoapp/providers/fobrain/credentials.go`
- `web/eino-workbench/src/presenters/fobrain.ts`

任务级检查：

```bash
go test ./internal/einoapp/providers/fobrain -run 'ConnectorStatus|CredentialBinding' -count=1
go test ./internal/einoapp/... -run CredentialLeak -count=1
```

### Task 8.4 实体消歧澄清

创建/修改：

- `internal/einoapp/providers/fobrain/disambiguation.go`
- `internal/einoapp/execution/clarification_tool.go`
- `web/eino-workbench/src/components/ClarificationCard.tsx`

任务级检查：

```bash
go test ./internal/einoapp/providers/fobrain -run Disambiguation -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-clarification
```

### Task 8.5 写域工单审批

创建/修改：

- `internal/einoapp/providers/fobrain/write_tools.go`
- `internal/einoapp/providers/fobrain/ticket_result.go`
- `web/eino-workbench/src/presenters/fobrain.ts`

任务级检查：

```bash
go test ./internal/einoapp/providers/fobrain -run WriteApproval -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-write-approval
```

### Task 8.6 Live smoke reports

产物：

- `test-results/eino-workbench-real-model-tool-report.json`
- `test-results/eino-workbench-live-read-report.json`
- `test-results/eino-workbench-live-write-report.json`
- `test-results/eino-workbench-skip-report.json`
- `docs/schemas/real-model-report.schema.json`
- `docs/schemas/live-read-report.schema.json`
- `docs/schemas/live-write-report.schema.json`
- `docs/schemas/skip-report.schema.json`

任务级检查：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-live-read
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-live-write
```

## Phase 9：对比评估

更新：

- `docs/09-comparison-scorecard.md`

任务级检查：

```bash
rg -n "阶段：|结论：" docs/09-comparison-scorecard.md
```
