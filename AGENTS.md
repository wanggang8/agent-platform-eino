# AGENTS.md

本项目是 Eino-first 全新重构实现。目标不是迁移旧代码，而是基于 `docs/` 重新开发一套新的 Agent Workbench。

## 目录角色

- 当前目录：新项目实现目录，可以修改。
- 旧项目目录：`/Users/vick/Desktop/project/ai-agent`，只读参考，不得修改。

## 开发原则

- 以 `docs/README.md` 为入口。
- 开发前必须执行 `docs/pre-development-validation.md`。
- 旧项目只能作为产品能力、验收基线和安全边界参考。
- 不复制旧 runtime 类型、旧执行链路、旧 Workbench 接口、旧 UI DOM 结构。
- 新实现必须使用 Eino-first 设计。
- Workbench 和 Action API 必须共用 Product Facts。
- 工具结果必须以 StructuredResult 作为唯一事实材料。
- Fobrain 能力恢复必须按 `docs/fobrain-tool-matrix.md`。

## 团队工程规范

### 任务开始前

- 每个任务开始前必须先确认所属 Phase、对应设计文档、实施计划、验收门禁和受影响目录。
- 如果需求、边界、数据来源、验收标准或旧项目参考含义不清楚，必须及时澄清；不得用猜测静默完成。
- 不得为了推进速度绕过 `docs/pre-development-validation.md`、schema、fixture、acceptance plan 或 import boundary。
- 旧项目 `/Users/vick/Desktop/project/ai-agent` 只能只读参考产品能力、验收基线和安全边界，不得复制旧 runtime 类型、旧执行链路、旧 Workbench 接口或旧 UI DOM 结构。

### 分层与职责

- 代码必须保持清晰分层：
  - `httpapi`：只处理 HTTP/API 边界、请求解析、响应封装。
  - `execution`：只处理 Eino 执行编排、run 生命周期、checkpoint/resume。
  - `facts`：只处理 Product Facts 的结构、存取和事实一致性。
  - `product`：只处理安全投影、Workbench/Action API 产品映射。
  - `capabilities`：只处理能力注册、能力选择、adapter 抽象。
  - `providers/*`：只处理具体业务 provider 接入。
  - `store/sqlite`：只处理持久化实现，不反向依赖业务、HTTP、LLM 或 provider。
- 不允许跨层偷懒调用。需要跨层协作时，通过明确接口、命令对象、事实对象或注册表连接。
- 不允许让 Workbench 展示、Action API、audit、replay、provider payload 互相直接耦合；所有产品出口必须从 Product Facts 投影。

### 禁止硬编码

- 不按工具名称、自然语言关键词、页面文案、旧 DOM 结构或 provider 私有字段硬编码业务逻辑。
- 工具选择必须基于能力注册、capability metadata、intent classification、policy decision 和 Product Facts，不得写死 `fobrain_xxx` 等工具分支。
- provider、模型、端口、凭据、外部 URL、开关、预算、超时、审批策略等部署相关值不得硬编码在业务代码中；必须来自配置、注册表、policy 或 fixture。
- schema version、capability id、tool id 可以作为契约常量存在，但必须集中定义并由测试覆盖，不得散落在业务分支中。
- 如果确实需要临时 hardcode，必须写明 Phase、移除条件、验收风险，并优先放入 fixture 或测试替身，不进入生产路径。

### Product Facts 与 StructuredResult

- Workbench 和 Action API 必须共用 Product Facts；不得各自维护事实模型。
- 工具结果必须以 StructuredResult 作为唯一事实材料；raw provider payload 只能在 provider 边界内处理，不能进入产品层、前端契约、audit 或 replay。
- LLM 上下文只能来自 safe Product Facts / StructuredResult projection，不得直接注入未投影 provider 数据。
- audit、replay、checkpoint、resume 必须引用同一组事实和 run 状态，不能产生第二套真相。

### 代码质量

- 新代码必须职责单一、命名稳定、可测试、可复用；不要提前抽象，但明显重复的契约映射、校验、错误封装应复用。
- Go 代码遵循 idiomatic Go：`gofmt`，小接口，清晰错误返回，错误字符串小写且不以标点结尾。
- TypeScript 代码以 contract 类型为边界，不手写与 schema 冲突的平行类型。
- 注释只解释不明显的设计原因、边界或风险；不要复述代码本身。
- 公共函数、导出类型、复杂状态机和安全边界必须有简短说明。
- 删除代码优于保留无用分支；保留兼容逻辑时必须说明触发条件和移除条件。

### 前端工程规范

- Workbench 前端架构遵循 `docs/frontend-architecture.md`，默认技术栈固定为 Vite + React + TypeScript + React Router + TanStack Query + Zustand + Radix UI Primitives + lucide-react。
- TanStack Query 只管理 server state；Zustand 只管理 UI/client state；不得在前端重建 Product Facts、run lifecycle、approval/resume 或 provider 事实模型。
- 前端组件只能消费 `web/eino-workbench/src/contracts/generated.ts` 的契约类型，不手写与 schema 冲突的 DTO。
- Radix 只作为无样式交互基础，视觉样式由本项目 CSS 定义；不得引入第三方默认视觉风格导致 Workbench 体验跑偏。
- 不默认使用 Next.js、Redux Toolkit、XState、Tailwind 或 shadcn；如需替换基础架构，必须先新增 ADR 并更新实施计划。
- Workbench UI 不得复制旧项目 DOM、CSS class、presenter 分支或接口字段；旧截图只能作为视觉验收基线。

### 测试与验收

- 行为变更必须配套测试；测试应验证事实、状态、错误和安全边界，而不是只验证实现细节。
- 每个任务完成后必须运行对应任务级检查，并在回复中列出实际执行的命令。
- schema、fixture、contract、OpenAPI、Go tests、TypeScript typecheck、smoke test 按受影响范围执行；不能只跑最小 happy path。
- 未到 Phase 的 smoke scenario 必须明确 `exit 2`，并标注不是通过信号。
- 修复 bug 时，先补能复现问题的测试或 fixture，再改实现。

### 文档同步

- 每次完成代码、契约、行为、目录结构或验收逻辑变更后，必须梳理并更新相关文档。
- 可能需要同步的文档包括：`docs/07-implementation-plan.md`、`docs/08-acceptance-plan.md`、schema、fixture、ADR、acceptance records、visual evidence、provider matrix。
- 如果实现改变了设计假设，必须更新设计文档或新增 ADR；不得让代码成为唯一事实来源。
- 如果发现文档与实现冲突，先修正冲突或明确记录偏差，不得继续在冲突基础上开发。

### 安全与配置

- 所有外部输入在可信边界校验；HTTP/API 入参、provider 返回、LLM 输出、resume payload 都视为不可信。
- 凭据、token、cookie、连接串不得写入代码、fixture、日志、截图或验收记录。
- provider 错误必须脱敏后进入 Product Facts / Workbench / Action API。
- 写域操作必须经过 policy、approval、audit；不得由 provider 或前端直接决定。
- 日志和 telemetry 只能记录必要上下文，不记录 raw payload、密钥、个人敏感信息或未投影安全数据。

### 评审与提交

- 变更应尽量小而完整；一个提交解决一个明确问题，避免把重构、格式化、功能开发和文档修复混在一起。
- 提交前自查：分层是否被破坏、是否新增硬编码、是否复用了契约、是否更新文档、是否跑过验收命令。
- 提交信息使用 Conventional Commits 风格，例如：
  - `feat(workbench): add fixture-driven shell`
  - `fix(facts): preserve run terminal state`
  - `docs(plan): clarify phase 2 acceptance gates`
  - `test(boundary): cover provider import rules`
- PR 或交付说明必须包含：改动摘要、验证命令、未覆盖风险、需要澄清的问题。
- 评审意见按事实和项目目标处理；不同意时说明原因和替代方案，不静默忽略。

### Agent 执行要求

- 每次任务完成后，必须用中文说明改动、验证结果、剩余风险和下一步建议。
- 有不确定、不清楚、不一致的问题，必须提出澄清；不要通过猜测扩大范围。
- 不得修改旧项目；不得回滚用户或其他 agent 的无关改动。
- 修改目录结构后，必须同步相关规范文档和实施计划。
- 修改代码后，必须检查分层清晰、结构清晰、职责清晰、复用合理，并运行相应测试。
