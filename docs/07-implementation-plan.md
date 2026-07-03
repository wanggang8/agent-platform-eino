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
- `internal/einoapp/httpapi/routes.go`
- `internal/einoapp/httpapi/routes_test.go`
- `internal/einoapp/httpapi/actions.go`
- `internal/einoapp/execution/commands.go`
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
- `internal/einoapp/facts/memory_repository.go`
- `internal/einoapp/product/projection.go`
- `internal/einoapp/product/errors.go`
- `docs/adr/YYYY-MM-DD-backend-foundation-boundaries.md`

要求：

- 配置必须集中读取 server、database、LLM、security、observability、timeout、budget；配置日志必须脱敏，DSN、URL userinfo、URL query 和 fragment 不得进入摘要。
- 后端服务配置只来自配置文件；`configs/eino-workbench.local.yaml` 可包含本地密钥但必须 ignored，不使用环境变量覆盖服务配置。
- 成功响应保持 OpenAPI 中定义的业务 schema 直出；错误响应统一 `eino_error_envelope.v1`，并包含 request id、安全错误码和可展示摘要。
- SSE 编码必须统一处理 `id`、`event`、`data`、flush、content type、no-cache 和编码失败。
- `httpapi` 必须通过 `Dependencies` 注入 product projection 和 execution command 接口；handler 不得持有全局 projection 或自行制造 run 执行结果。
- `Dependencies` 必须成对提供 projection 和 commands；禁止只注入一侧后让另一侧静默 fallback，避免 split facts。
- `httpapi` 不得手写业务 DTO 第二套真相；请求解析结构必须覆盖对应 schema 字段并只负责 transport 边界，业务输出只能调用 product projection、facts query 和 execution command 接口。
- `httpapi` 必须解析 `message`、`action`、`resume` 请求并校验单个 JSON body；合法 schema 字段如 Action `attachments/context` 不得被误拒。
- `execution` 先提供 Message/Action/Resume command 边界；Phase 1 默认命令只能创建 accepted run fact，resume 必须确认目标 run fact 存在，不实现真实 Eino runner、checkpoint 或 stream reducer。
- `DefaultDependencies` 必须让 execution commands 和 product projection 共享同一个 facts repository；禁止默认服务路径出现执行接受结果和产品投影脱节。
- `llm` 只暴露 provider interface、config、mock provider、redacted error；真实 OpenAI-compatible provider 可在 Phase 4 扩展。
- `capabilities` 必须先提供 registry、provider interface、tool metadata、risk/approval policy；任何工具选择都必须通过 registry，不得写死工具名或自然语言关键词。
- `facts` 先定义最小 repository interface、事实模型和 Phase 1 内存 repository；Product Facts 的 SQLite 实现仍在 Phase 3.2 完成。
- `product` 先定义 Workbench、ActionResult、Replay、SSE projection 的接口边界和 facts-backed 空投影；实际完整映射仍在 Phase 3 完成。
- 不创建通用 `utils` 包；公共能力按职责放入 `bootstrap`、`httpapi`、`llm`、`capabilities`、`facts`、`product`、`observability`。

任务级检查：

```bash
go test ./internal/einoapp/bootstrap ./internal/einoapp/httpapi ./internal/einoapp/execution ./internal/einoapp/llm ./internal/einoapp/capabilities ./internal/einoapp/facts ./internal/einoapp/product -run 'Config|Response|SSE|Commands|Provider|Registry|Facts|Projection|Redaction' -count=1
go test ./internal/einoapp/httpapi -run 'MessageRequestParsesIntoExecutionCommand|ActionRequestParsesActionID|ResumeRequest|ResumeMissingRun|TrailingJSON|InvalidMessage|InvalidAction|InvalidResume|DefaultDependencies|PartialDependencies|InjectedProjection' -count=1
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

### Task 3.0 Phase 3 technical gate

读取：

- `docs/phase-3-execution-and-facts-plan.md`
- `docs/facts-contract.md`
- `docs/run-lifecycle.md`
- `docs/conversation-context.md`
- `docs/intent-and-capability-selection.md`
- `docs/observability-and-budgets.md`
- `docs/adr/2026-06-30-pin-eino-version.md`

要求：

- 开发前确认 Phase 3 仍按 facts model -> SQLite repository -> product projection -> execution event mapper -> ChatModelAgent Runner -> HTTP/SSE/Action API 的顺序执行。
- 如果文档、schema、fixture 或当前代码边界冲突，先修文档或新增 ADR，不得直接编码。
- Eino API 以 `go.mod` 锁定版本和本地 module source 为准；在线资料只作概念参考。

任务级检查：

```bash
rg -n 'phase-3-execution-and-facts-plan|Product Facts|ResumeWithParams|StructuredResult' docs
go list -m github.com/cloudwego/eino
go test ./internal/einoapp/architecture -run ImportBoundary -count=1
```

### Task 3.1 Pin Eino versions

修改：

- `go.mod`
- `go.sum`
- `docs/adr/YYYY-MM-DD-pin-eino-version.md`

任务级检查：

```bash
go list -m github.com/cloudwego/eino
go test ./internal/einoapp/execution -run 'EinoVersion|ChatModelVersion' -count=1
```

记录：

- 固定版本写入 ADR。
- `eino-ext` 模型 provider 只在实际引入 provider 代码时固定到 `go.mod`，不得要求未引入模块通过 `go list -m`。
- `eino-ext/components/model/openai` 版本调研和固定放到 Phase 4 真实 LLM provider 任务执行。
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

要求：

- 技术决策按 `docs/adr/2026-07-01-phase-3-facts-sqlite-decisions.md` 执行。
- 先补 facts 模型和 repository interface，再做 SQLite 实现。
- SQLite repository 必须覆盖 run/turn/tool/pending/audit/context/lifecycle/idempotency 的事务边界。
- memory repository 只能作为测试替身，不作为 Phase 3 产品事实来源。
- raw checkpoint id、interrupt id、credential ref、raw provider payload 和 reusable resume token 不进入 Product Facts；`checkpoint_ref` 只能是内部安全引用。

任务级检查：

```bash
go test ./internal/einoapp/facts -run 'RunStatus|StructuredResult|MemoryRepository' -count=1
go test ./internal/einoapp/store/sqlite -run 'RunTurnEvent|ToolCallResult|ContextSnapshot|LifecycleTransition|PendingResume|Idempotency|Unsafe' -count=1
```

### Task 3.3 ChatModel and Runner

创建：

- `internal/einoapp/execution/runner.go`
- `internal/einoapp/execution/events.go`
- `internal/einoapp/execution/selection.go`
- `internal/einoapp/execution/context.go`

要求：

- 先实现 execution event mapper 测试，再接真实 Runner。
- Runner event 只能写入 Product Facts，不能直接暴露给 Workbench、ActionResult 或 replay。
- assistant delta、tool call、tool result、run state 必须先落入 facts，再由投影层读取。
- Phase 3 可用 mock `llm.Provider` 适配 Eino `BaseChatModel` 跑通 ChatModelAgent Runner；真实 OpenAI-compatible provider 仍在 Phase 4。
- 引入真实 Eino Runner 代码后，必须移除 `internal/einoapp/execution/eino_version_pin.go` 占位导入，版本锁定由 `go.mod` 和测试保证。
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

要求：

- HTTP handler 只能调用 execution command 和 product projection，不读取 Eino event 或 provider payload。
- SSE event id 必须来自 Product Facts cursor/sequence。
- Action API 必须复用 Workbench 同一套命令、facts 和 projection 路径。
- 服务启动路径必须使用配置文件指定的 SQLite Product Facts repository；memory repository 只允许作为单元测试替身。
- Phase 3 capability registry 只能由配置文件或后续 provider 注册填充；不得在生产启动路径内置 smoke、Fobrain 或业务能力 ID。
- `chat-stream` 和 `action-basic` smoke 必须创建真实 run，读取 snapshot/stream，并断言输出来自 Product Facts。
- `capability-selection` smoke 必须验证显式 capability hint 经过配置驱动的 registry/policy，未知 hint 在创建 run 前被拒绝，需要审批的 hint 不触发 runner/context snapshot。
- `context-projection` smoke 必须验证模型输入前已写入 safe context snapshot，且不包含 raw provider payload、credential、token、checkpoint 或 interrupt 信息。
- Phase 3 当前只证明 mock ChatModelAgent assistant 路径、capability 选择门禁和 context projection；真实 tool loop、StructuredResult 转换、真实 LLM provider、HITL/checkpoint/resume 仍属于后续 Phase。

任务级检查：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario chat-stream
bash scripts/eino_workbench_server_smoke.sh --scenario action-basic
bash scripts/eino_workbench_server_smoke.sh --scenario capability-selection
bash scripts/eino_workbench_server_smoke.sh --scenario context-projection
```

## Phase 4：Tools、Provider、Safety

### Task 4.0 Phase 4 technical gate

读取：

- `docs/adr/2026-07-01-phase-4-tool-loop-before-real-llm.md`
- `docs/capability-provider-contract.md`
- `docs/facts-contract.md`
- `docs/llm-provider-safety.md`
- `docs/intent-and-capability-selection.md`
- 本地 Eino `v0.9.12` API：`components/tool/interface.go`、`schema/tool.go`、`adk/chatmodel.go`

要求：

- Phase 4 必须先证明本地 mock tool loop 的 Product Facts 链路，再接真实 LLM provider。
- 不得让真实 provider 网络、凭据或模型选择不稳定性成为 StructuredResult / Safety Gate / tool-card 的前置依赖。
- Eino tool 只作为执行编排入口；Workbench、ActionResult、SSE、Replay 和 Audit 仍只能从 Product Facts 投影。
- `tool-card` smoke 打开前，必须已有 StructuredResult Safety Gate、mock capability adapter 和 Eino tool loop 单元测试。

任务级检查：

```bash
go list -m github.com/cloudwego/eino
test -f docs/adr/2026-07-01-phase-4-tool-loop-before-real-llm.md
rg -n 'StructuredResult|Safety Gate|ToolInfo|InvokableTool|tool-card' docs/07-implementation-plan.md docs/capability-provider-contract.md docs/facts-contract.md
```

### Task 4.1 StructuredResult and Safety Gate

创建：

- `internal/einoapp/product/structured_result.go`
- `internal/einoapp/product/safety.go`
- `internal/einoapp/product/assistant_safety.go`
- `internal/einoapp/product/structured_result_test.go`
- `internal/einoapp/product/assistant_safety_test.go`

范围：

- 定义 StructuredResult candidate 与 Safety Gate 通过后的安全 StructuredResult。
- candidate 可以来自 provider/tool adapter；只有 Safety Gate 通过后的结果能写入 Product Facts。
- Safety Gate 必须拒绝 Authorization、API key、token、password、raw credential ref、raw provider body、checkpoint id、interrupt id。
- Assistant 输出进入 Product Facts 前也必须经过 assistant safety；`execution.EventMapper` 写入 `Turn` 和 `ToolResult` 前必须调用 product 安全门。
- 不把 provider raw payload、Eino raw event 或模型 raw message 放入 `facts.ToolResult`。

任务级检查：

```bash
go test ./internal/einoapp/product -run 'StructuredResult|Safety|AssistantSafety' -count=1
go test ./internal/einoapp/execution -run 'RunnerEventMapper' -count=1
go test ./internal/einoapp/facts ./internal/einoapp/store/sqlite -run 'StructuredResult|Unsafe' -count=1
```

完成状态：

- 已实现 `StructuredResultCandidate`、`StructuredResultSafetyGate`、`AssistantSafetyGate` 和 `ContainsUnsafeMaterial`。
- 已将安全门接入 `execution.EventMapper`，内存仓库和 SQLite 仓库路径都会拒绝 unsafe assistant、pending、StructuredResult schema/result_ref/summary 材料。
- 当前 Phase 4.1 candidate 只包含 schema、result_ref、safe_summary。若 Phase 4.2 需要增加 display payload、evidence、pagination 或嵌套结构，必须先补 schema、size、depth 和 raw error/provider payload 门禁测试。

### Task 4.2 Mock capability adapter and Eino tool loop

创建：

- `internal/einoapp/capabilities/tool_adapter.go`
- `internal/einoapp/capabilities/mock_provider.go`
- `internal/einoapp/capabilities/provider_contract_test.go`
- `internal/einoapp/execution/tools.go`
- `internal/einoapp/execution/tool_loop_test.go`

范围：

- 先用配置驱动的 mock read capability 跑通 Eino `tool.BaseTool` / `tool.InvokableTool`。
- 如果 mock adapter 增加 `StructuredResultCandidate` 的结构化字段，必须先扩展 Task 4.1 Safety Gate 的 schema/size/depth 校验。
- `Info()` 只能从 Capability Registry metadata 生成 `schema.ToolInfo`；tool name 不得按 Fobrain 名称、自然语言关键词或旧工具名分支。
- `InvokableRun(ctx, argumentsInJSON)` 必须校验参数、写 `ToolCall`、把 provider candidate 交给 Safety Gate、写 `ToolResult`，并只向 Eino 返回安全 tool message。
- tool call、tool result、safe args summary、safe audit 必须先进入 Product Facts，再由 product projection 生成 Workbench/ActionResult/SSE。
- 写域或 `approval_required=true` capability 在 Phase 6 前不能执行 mutation；Phase 4 只验证 read-only mock tool loop。
- `tool-card` smoke 必须启动服务，用临时配置注册 mock capability，生成 tool card、ActionResult result card、SSE tool patch 和 audit。

任务级检查：

```bash
go test ./internal/einoapp/capabilities -run 'ProviderContract|StructuredResultConversion|ToolSelectionMetadata|EinoToolAdapter' -count=1
go test ./internal/einoapp/execution -run 'AgentToolLoop|ToolFacts|ToolSafety' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario tool-card
```

完成状态：

- 已实现 registry metadata -> Eino `schema.ToolInfo` adapter、本地 mock read provider 和 `ToolLoopRunner`。
- Action API 显式选择只读 capability 时进入 mock Eino `InvokableTool`，写入 ToolCall、Safety Gate 后的 ToolResult、tool audit，并由 Product Facts 投影 Workbench tool card、ActionResult result card 和 SSE `tool.updated`。
- Phase 4.2 的参数装配按 `Capability.InputSchema.Properties` 选择字段，当前只支持单文本入口和轻量 `map[string]string` schema；完整 JSON Schema 2020-12 / `oneOf` / `$defs` adapter 必须在接真实 provider 前补齐。
- 写域或需要 approval 的 capability 仍在 Phase 6 前阻断执行，只保留选择审计。

### Task 4.3 Production LLM provider implementation

创建：

- `internal/einoapp/llm/openai_compatible.go`
- `internal/einoapp/llm/network_policy.go`
- `internal/einoapp/llm/openai_compatible_test.go`
- `scripts/eino_workbench_real_model_smoke.sh` 或扩展 `scripts/eino_workbench_server_smoke.sh`

范围：

- 基于 Phase 1.6 的 `llm.Provider` 和 `llm.Config` 接入 OpenAI-compatible ChatModel。
- 真实 provider 只替换模型边界，不改变 Product Facts、tool adapter、product projection 或 HTTP handler。
- API key、Authorization、base_url、model、timeout、network policy 只来自 ignored 本地配置文件。
- 禁止把 API key、Authorization、raw provider body 写入 Product Facts、Workbench、ActionResult、audit、replay、日志或验收记录。
- provider non-2xx、timeout、malformed response 和 network blocked 必须转换为安全错误。
- 无本地凭据时 real model smoke 只能生成 skipped report，不能作为 passed。
- `real-model-chat` 归属 P1 真实模型 smoke，不是 P0 `tool-card` 的前置条件。

任务级检查：

```bash
go test ./internal/einoapp/llm -run 'OpenAICompatible|NetworkPolicy|RedactedError' -count=1
# P1 真实模型 smoke；无 ignored 本地凭据时必须生成 skipped report，不能作为 P0 passed。
bash scripts/eino_workbench_server_smoke.sh --scenario real-model-chat --config configs/eino-workbench.local.yaml
```

完成状态：

- 已实现 OpenAI-compatible 非流式 `/chat/completions` provider，启动路径支持 `mock` 与 `openai_compatible`。
- API key 只进入 provider 私有边界；`llm.Config`、Product Facts、Workbench、ActionResult、audit 和报告只使用脱敏配置摘要。
- 网络策略已覆盖 userinfo 拒绝、host allowlist、非 HTTPS 阻断、redirect 阻断/重校验，以及 DNS/dial 后私网、loopback、link-local、metadata 目标阻断。
- `real-model-chat` smoke 已接入统一 smoke 脚本：无 ignored 本地配置或缺 `api_key` 时生成 `test-results/eino-workbench-skip-report.json`；有凭据时生成 `test-results/eino-workbench-real-model-provider-report.json`，成功写 `passed`，provider/action/snapshot/redaction 失败写 `failed`。
- smoke 运行配置由 `scripts/eino_workbench_config_prepare.go` 结构化生成，只覆写临时 server addr 和 SQLite DSN，并输出脱敏 LLM 摘要。
- 真实 streaming provider 尚未启用，`Stream` 当前返回安全错误；后续如接入 streaming，必须先补 streaming delta safety gate 和 SSE 映射验收。

### Task 4.4 Provider policy and credentials

创建：

- `internal/einoapp/capabilities/policy.go`
- `internal/einoapp/capabilities/credentials.go`
- `internal/einoapp/capabilities/provider_policy_test.go`

范围：

- risk、side_effect、workspace scope、permission、credential binding、approval_required 决策。
- credential missing、scope denied、connector unavailable 必须产生安全 ActionResult / Workbench 状态。
- 任何 `credential_ref`、secret 或 token hash 不得进入 JSON 输出。
- 该任务只增强 policy 与 credential 边界；Phase 6 前不实现真实 approval resume。

任务级检查：

```bash
go test ./internal/einoapp/capabilities -run 'ProviderPolicy|CredentialBinding|CredentialLeak' -count=1
```

完成状态：

- 已实现 provider policy decision：risk、side_effect、approval_required、workspace scope、credential binding policy 和 connector status 均进入统一决策；risk 使用 capability catalog 的 `none/low/medium/high` 枚举。
- policy reason 使用 `docs/schemas/provider_policy_decision.v1.schema.json` 中的稳定 reason code；写域和 `write_external` 在 Phase 6 前返回 `approval_required`，不会执行 mutation。
- 凭据输出只使用 `CredentialBinding` 安全摘要；`status` 必须归一到 `configured/missing/unbound/bound`，`credential_ref`、token、secret、Authorization 和 raw connector config 会被折叠为安全 display/audit 摘要。
- Action API 选择到缺凭据、workspace scope denied 或 connector unavailable 等非审批阻断时，会写入同源 Product Facts 的 failed run，`safe_error=provider_policy_blocked`，不会把 credential reason 明文写入 Product Facts。
- 配置文件 capability 元数据已要求声明 `side_effect`、`policy_ref`、`permission_scope`、`credential_binding_policy` 和 `idempotency_required`，并支持按 provider 需要声明 `connector_id`。

### Task 4.5 MCP adapter contract

创建：

- `internal/einoapp/capabilities/mcp_adapter.go`
- `internal/einoapp/capabilities/mcp_session.go`
- `internal/einoapp/capabilities/mcp_schema.go`
- `internal/einoapp/capabilities/mcp_mock_test.go`

范围：

- Phase 4 只完成 MCP adapter contract，不接生产 MCP server。
- MCP 具体生命周期、复杂 auth、多 server catalog 和 live MCP smoke 后移到 P1/P2。
- 先满足 `docs/capability-provider-contract.md`，再接入具体业务 provider。
- MCP 输出必须先转 StructuredResult candidate，再经过 Safety Gate。

MCP adapter contract 任务拆分：

- transport / server catalog config shape。
- lifecycle / initialize mock。
- tools/list pagination mock。
- inputSchema conversion。
- outputSchema / structuredContent / isError conversion。
- Safety Gate integration。
- annotations + project risk policy merge。

任务级检查：

```bash
go test ./internal/einoapp/capabilities -run 'MCPLifecycle|MCPToolListPagination|MCPSchemaConversion|MCPSafety|MCPRiskPolicy' -count=1
```

完成状态：

- 已实现 `mcp_mock_servers` 配置形态，支持 mock server catalog、tools/list metadata、tools/call 固定结果和安全摘要。
- 已实现 mock MCP initialize/session 状态、tools/list pagination、input schema 子集转换与 required/type 校验、structuredContent/isError 到 StructuredResult candidate 的转换。
- 已实现 MCP annotations 与项目可信 risk policy 合并；annotations 不可信，不能覆盖项目写域审批和高风险策略。
- MCP destructive annotation 触发写域默认策略时也会强制 `approval_required` 与 `idempotency_required`。
- 已将 mock MCP provider 接入启动 runtime：registry 和 invoker 通过 capability id 路由，不按工具名或 provider 名称硬编码执行分支。
- mock MCP `tools/call` 必须在 initialize 后执行；required 字段会传入 Eino ToolInfo，并用于 Action API 文本输入字段选择。
- `mcp_mock_servers[].tools[].project_policy` 已按 capability policy 枚举和写域幂等规则做启动校验。
- 已新增 `mcp-mock` smoke 场景；生产 MCP server、复杂 auth、多 server live catalog 和 live MCP smoke 仍后移。

## Phase 5：Fobrain provider PoC

设计与执行依据：

- `docs/fobrain-provider-poc-design.md`
- `docs/fobrain-provider-config.md`
- `docs/adr/2026-07-02-fobrain-provider-poc-boundary.md`
- `docs/superpowers/plans/2026-07-02-fobrain-provider-poc.md`

创建：

- `internal/einoapp/providers/fobrain/provider.go`
- `internal/einoapp/providers/fobrain/catalog.go`
- `internal/einoapp/providers/fobrain/client.go`
- `internal/einoapp/providers/fobrain/credentials.go`
- `internal/einoapp/providers/fobrain/tools.go`
- `internal/einoapp/providers/fobrain/result.go`
- `internal/einoapp/providers/fobrain/errors.go`
- `scripts/eino_workbench_server_smoke.sh` 的 `fobrain-poc` 分支和报告生成逻辑。

任务顺序：

1. 先补 provider catalog contract test，固定 `tool.fobrain.current_user_context` PoC capability 元数据。
2. 实现 `Provider/ListCapabilities`，所有 tool id、policy ref、schema ref 集中在 catalog。
3. 先补 Fobrain YAML 配置契约测试，再实现 `bootstrap.Config` 的 Fobrain 配置、派生安全摘要和 no-secret example。
4. 补 client/credential 边界测试，覆盖缺凭据、workspace scope denied、connector unavailable 和错误脱敏。
5. 实现 `FobrainClient`、`CredentialResolver` 和 provider safe error，不让 raw provider payload 离开 client 边界。
6. 补 StructuredResult mapper 测试，覆盖成功 candidate、Safety Gate 通过、`tool.structured_result.v1` / `fobrain.tool_result.v2` 关系和 unsafe material 拒绝。
7. 实现 `capabilities.Invoker`，通过 catalog operation 分发，不在 execution/httpapi 中按 Fobrain 工具名分支。
8. 通过配置/catalog 接入启动路径；默认配置未声明时不得内置 Fobrain provider。
9. 先补 `fobrain-poc` provider report / skip report schema 检查，再实现 smoke 分支。
10. `fobrain-poc` 必须生成 provider PoC report；真实模型工具选择仅在 LLM/Fobrain 凭据齐全时生成 real-model report，缺凭据时生成 skip report 并用 `blocks_claims` 阻断对应声明。

任务级检查：

```bash
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
go test ./internal/einoapp/providers/fobrain -run 'Provider|Client|ReadonlyPoC|StructuredResult|Credential|Unsafe' -count=1
go test ./cmd/eino-workbench ./internal/einoapp/bootstrap ./internal/einoapp/capabilities ./internal/einoapp/product ./internal/einoapp/execution -run 'Fobrain|Capability|Config|Policy|StructuredResult|ToolSafety|ToolFacts' -count=1
npm run eino-workbench:schema-test
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config configs/eino-workbench.local.yaml
```

完成后只允许声明 Fobrain provider PoC 已接入新能力链路；不得声明 24 只读恢复、live read 覆盖、实体消歧、写域审批可用或可替换旧项目。

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

- `internal/einoapp/execution/commands.go`
- `internal/einoapp/facts/model.go`
- `internal/einoapp/facts/repository.go`
- `internal/einoapp/store/sqlite/repository.go`
- `internal/einoapp/httpapi/actions.go`
- `internal/einoapp/httpapi/routes.go`
- `web/eino-workbench/src/components/RunNotice.tsx`（后续 UI 子任务）

覆盖：

- cancel waiting/running run。
- stop running run。
- provider timeout 和 pending timeout 的安全状态迁移。
- retry 只按 `run-lifecycle.md` 中的安全策略允许。
- 终态操作幂等。
- 当前后端切片已打开 `run-lifecycle` smoke：通过本地服务验证 cancel、stop、provider timeout、pending timeout、`schema_invalid` retry、终态幂等、replay/audit 投影和 raw payload 禁止项。smoke 使用 SQLite seed 准备 running/waiting 状态；provider timeout retry 需等 capability metadata 可证明只读后开放，Workbench `RunNotice` 仍按后续 UI 任务验收。

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
- 当前切片已打开 `action-consistency` smoke：同一 Action run 会校验 ActionResult、run snapshot、views/current、replay、SSE、audit refs 和 SQLite facts 同源，校验错 workspace 访问返回 404，且不暴露 raw provider payload 或凭据痕迹。

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
- 当前切片已打开 `replay` smoke：回放视图必须精确等于同 run 的 Workbench snapshot 投影，events 必须包含由 Product Facts snapshot 派生的产品事件和安全 audit event。持久化 facts cursor/event log 仍归属后续 replay 硬化任务。

任务级检查：

```bash
go test ./internal/einoapp/product -run Replay -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario replay
```

## Phase 8：Fobrain 产品能力恢复

Phase 8 是 P2 门禁，属于完整重构必做范围；未完成本阶段不得声明重构完成、能力等价或替换当前产品基线。

### Task 8.0 Live read 认证与分批门禁

创建/修改：

- `docs/adr/2026-07-02-fobrain-live-auth-and-batch-gates.md`
- `docs/fobrain-live-read-batch-plan.md`
- `docs/fobrain-provider-config.md`
- `docs/fobrain-tool-matrix.md`

要求：

- 固定 Fobrain live read 使用 workspace 共享 token，不实现 per-user/per-tool token。
- 固定 `auth_param` 来自配置文件，当前真实环境参数名为 `authorization`，token 值只来自 ignored local config。
- 明确 live client 只接受显式配置的 `auth_param`，并把它作为 header 名发送原始 token；缺失时必须阻断请求，不能在 provider 内默认；如需 query/body 参数必须新增 ADR。
- 24 个只读工具必须按 Batch A-E 小批次恢复，每批先 mock contract，再 live pass report，再视觉和脱敏证据；无 live 环境时只能生成 blocking skip report，不能声明通过。
- `docs/fixtures/fobrain/tool-matrix-24.json` 必须包含 `batch_gate`，schema 必须阻止未标批次的工具进入矩阵。
- 每批都必须证明 Workbench 和 Action API 使用同源 Product Facts，StructuredResult 是唯一事实材料。
- Batch E 开工前必须完成 `docs/fobrain-batch-e-interface-plan.md`，确认资产详情路径按 `network_type` 分流，业务风险走 `/threat_center/count`，威胁关联走 `/threat_center/relevance/list`，并明确 live 样本 ID 不进入提交物。

任务级检查：

```bash
npm run eino-workbench:schema-test
git diff --check
```

### Task 8.1 24 个只读工具矩阵

创建/修改：

- `docs/fobrain-tool-matrix.md`
- `docs/fobrain-live-read-batch-plan.md`
- `docs/fixtures/fobrain/*.json`
- `docs/schemas/fobrain/*.schema.json`
- `cmd/eino-workbench/main.go`
- `internal/einoapp/bootstrap/config.go`
- `internal/einoapp/providers/fobrain/http_client.go`
- `internal/einoapp/providers/fobrain/tools.go`
- `internal/einoapp/providers/fobrain/result.go`

前置要求：

- 24 个只读工具必须先补齐 `fobrain-tool-matrix.md` 的字段级恢复模板。
- 每个工具必须有 input schema、result schema、fixture、mock assertion、live assertion 和 screenshot state。
- 每个工具的 mock/live 断言必须覆盖 StructuredResult `display_type`、safe summary、evidence、pagination/empty/error 语义和敏感字段脱敏。
- 未补齐矩阵的工具不得进入实现。
- 工具实现顺序必须遵循 Batch A-E；每个批次完成前不得把后续批次标记为可验收。
- `batch_gate` 必须由 `tool_matrix.v1.schema.json` 校验，并与 `docs/fobrain-live-read-batch-plan.md` 一致。
- Batch A 代码基础必须覆盖 `connector.fobrain.security`、`tool.fobrain.current_user_context`、`tool.fobrain.my_permissions`；Batch A live/smoke 必须生成 `docs/schemas/fobrain/batch_a_live_report.v1.schema.json` 约束的报告。Batch A 最终完成声明还必须补 Workbench 视觉证据和验收记录。
- Batch E 当前切片已覆盖 `get_asset_detail`、`get_vulnerability_detail`、`business_risk_summary`、`threat_relevance_list` 的 catalog、输入 mapper、mock StructuredResult、HTTP mapper、sample discovery 本地样本扩展、`fobrain-batch-e` live smoke/report schema、真实环境脱敏 live pass 报告，以及 Workbench 视觉、fresh replay、audit evidence。该结论只覆盖 Batch E 四工具，不得把当前切片等同于 Fobrain 24 只读恢复完成。

任务级检查：

```bash
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
go test ./internal/einoapp/providers/fobrain -run 'HTTPClientCurrentUserContext|Provider|Client|Credential|StructuredResult|Unsafe' -count=1
go test ./internal/einoapp/providers/fobrain -run 'BatchE|AssetDetail|VulnerabilityDetail|BusinessRisk|ThreatRelevance' -count=1
go test ./scripts/fobrain_sample_discovery ./scripts/fobrain_batch_e_smoke -count=1
npm run eino-workbench:schema-test
go test ./cmd/eino-workbench -run 'Fobrain' -count=1
go test ./internal/einoapp/providers/fobrain -run ReadonlyToolMatrix -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-a --config configs/eino-workbench.local.yaml
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-e --config configs/eino-workbench.local.yaml --asset-id "<资产ID>" --vulnerability-id "<漏洞ID>" --business-name "<业务系统>" --vulnerability-name "<漏洞名>"
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-readonly
```

### Task 8.2 Fobrain presentation evidence matrix

创建/修改：

- `docs/fobrain-tool-matrix.md`
- `docs/fixtures/fobrain/tool-matrix-24.json`
- `web/eino-workbench/tests/fobrain-tool-visual.spec.ts`
- `test-results/eino-workbench-fobrain-tool-acceptance/`

当前进展：

- Batch A 业务只读视觉切片已覆盖 `tool.fobrain.current_user_context` 和 `tool.fobrain.my_permissions`，每个工具在 desktop/mobile 下生成六区域 Playwright screenshot baseline。
- Batch E 四个详情与风险关联只读工具已覆盖 Workbench 安全投影、StructuredResult 展示、fresh replay、audit evidence 和 desktop/mobile 六区域 baseline。
- 本任务仍需继续补齐 Batch B-D 视觉矩阵并执行 Action API 真实服务同源 smoke，不得据此声明 Fobrain 24/24 恢复完成。

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

当前进展：

- `connector.fobrain.security` 已补齐 Workbench 视觉证据，使用 `fobrain.tool_result.v2` 的 `connector_status` 安全投影展示 connector 状态和 workspace credential binding 摘要。
- connector 视觉基线已覆盖 desktop/mobile 的 `main-chat`、`fresh-main-chat`、`process`、`evidence`、`audit`、`internal-details` 六区域。
- 本任务只证明 connector 状态和凭据绑定安全展示；Batch B-E、24 个只读工具全量恢复、实体消歧和写域审批仍按后续任务验收。

任务级检查：

```bash
go test ./internal/einoapp/providers/fobrain -run 'ConnectorStatus|CredentialBinding' -count=1
go test ./internal/einoapp/... -run CredentialLeak -count=1
```

### Task 8.4 实体消歧澄清

创建/修改：

- `internal/einoapp/providers/fobrain/disambiguation.go`
- `internal/einoapp/execution/clarification_tool.go`
- `internal/einoapp/product/projection.go`
- `docs/schemas/eino_workbench_view.v1.schema.json`
- `web/eino-workbench/src/features/workbench/components/cards.tsx`

任务级检查：

```bash
go test ./internal/einoapp/providers/fobrain -run Disambiguation -count=1
go test ./internal/einoapp/execution -run Clarification -count=1
go test ./internal/einoapp/product -run ClarificationCandidates -count=1
npm run eino-workbench:contract-test
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
