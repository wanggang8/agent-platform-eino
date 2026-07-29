---
baseline_commit: b8f7e903c75bade747ae887225cacbe8b7db5a15
---

# Story 1.2：用单一安全 Fixture 跑通只读 Walking Skeleton

Status: done

WorkItemType: user-value

Milestone: M1

## Story

As a 产品中心运营，
I want 在浅色桌面 Workbench 中发起一次固定的新增漏洞查询并看到摘要，
so that 我能尽早验证从聊天入口到安全事实展示的完整产品链路。

## Requirements

- Functional：FR-01、FR-02、FR-03、FR-04。
- Non-functional：NFR-01、NFR-03。
- Architecture：AR-02、AR-03、AR-04、AR-05、AR-06、AR-07、AR-08、AR-17、AR-19、AR-20、AR-21、AR-22、AR-23、AR-24、AR-25、AR-28、AR-29、AR-30、AR-31、AR-36、AR-37。
- UX：UX-DR-01～UX-DR-07、UX-DR-19～UX-DR-25、UX-DR-28、UX-DR-29。
- Decisions：AD-01～AD-05、AD-13、AD-16～AD-18、AD-20、AD-24、AD-25，以及 SQLite WAL 安全基线 ADR。
- Gate / milestone：G-ARCH-V2 的 M1 子集；M-1。

## Prerequisites

- Story 1.1=`done`，`G-TOOLCHAIN=PASS`。
- Implementation Readiness 2026-07-20=`READY`。
- 开发前复核已于 2026-07-27 裁决 `ALLOW`，记录见 [`pre-development-validation-1-2-2026-07-27.md`](../../docs/acceptance-records/pre-development-validation-1-2-2026-07-27.md)；该裁决只授权本 Story 的限定范围。
- 本 Story 的真实 LLM、真实 FOBrain 和外部 mutation 上限均为 `0`。

## Inputs / Outputs

- 输入：一个版本化的新增漏洞安全 fixture bundle、全新 SQLite 数据库、固定聊天查询。
- 输出：一个按值持久化的 StructuredResult、一个 QueryResultSnapshot、同源 Product Facts / Workbench view / SSE 投影，以及一张浅色桌面查询结果卡。
- fixture bundle 只描述一个能力，包含 `resolved`、`empty`、`failed` 三个命名 case；一次运行只选择一个 case，不得注册三个 capability。`resolved` 至少包含两条不同发现时间的安全对象，并包含同一发现时间的 tie-breaker 样本，用于证明全局 `discovered_at DESC + snapshot_item_ref ASC` 稳定排序。

## Acceptance Criteria

### AC-01：贯通唯一只读事实链

**Given** 全新数据库和唯一批准的新增漏洞安全 fixture

**When** 用户在桌面聊天入口提交固定查询

**Then** mock Eino loop 只能选择一个注册的只读 capability，结果经 Safety Gate 成为唯一 StructuredResult 与 QueryResultSnapshot 后持久化

**And** Workbench 从同一 Product Facts 投影显示查询序号、摘要、数量和观察时间。

### AC-02：空结果与失败确定性区分

**Given** fixture 分别表示有结果、0 条和失败

**When** 查询完成

**Then** 0 条显示“没有待派发漏洞”，失败显示“无法读取漏洞事实。此次请求不是空结果。”，两者不得混淆

**And** raw payload、token、provider locator、内部 result_ref 和旧 runtime 类型不得进入产品出口。

### AC-03：浅色桌面三栏与零写入

**Given** 用户查看首版桌面页面

**When** 浏览三栏结构与键盘焦点

**Then** 页面仅使用浅色桌面布局，左栏固定提供“操作记录”入口，中栏显示聊天，右栏仅显示安全“事实／执行记录”空态

**And** 操作记录入口不查询历史数据，不出现写入入口，真实 FOBrain mutation 调用数为 0。

### AC-04：从持久事实恢复同一结果

**Given** 服务已经完成一次查询并持久化 Product Facts

**When** 服务重启后重新打开该结果

**Then** 同一结果仍由持久事实安全投影，查询序号、摘要、数量和观察时间保持一致

**And** 不重新调用 fixture/capability，不依赖聊天摘要、日志或 Eino event 重建，也不要求完整 SSE 恢复、上下文引用或历史操作列表。

## Tasks / Subtasks

- [x] Task 1：先固定 M1 安全依赖与机器契约（AC: 01、02、04）
  - [x] 将 `modernc.org/sqlite` 精确升级到 `v1.46.2`，确认 `sqlite_version() >= 3.51.3`；只做该安全升级，不混入其他依赖维护。
  - [x] 新建 M1 所需的 StructuredResult v2、Product Facts v2 / QueryResultSnapshot、Workbench view 与最小 SSE schema；不得生成 v1/v2 union。
  - [x] 新建一个版本化新增漏洞 fixture bundle，固定 `resolved`、`empty`、`failed` case 的安全字段、数量、观察时间和确定性错误；`resolved` 覆盖多条、同时间 tie-breaker 与全局发现时间倒序。
  - [x] 同步 OpenAPI 3.1、schema manifest、contract generator 与 `web/eino-workbench/src/contracts/generated.ts`；TypeScript 不得手写平行 DTO。
  - [x] 收窄 schema validator，使当前 M1 只以本 Story 契约为通过条件；不得继续把 FOBrain 24 工具、Action 或旧视觉矩阵作为 M1 前置。

- [x] Task 2：建立最薄 Product Facts v2 领域与端口（AC: 01、02、04）
  - [x] 在 `facts` 定义按值不可变 StructuredResult：安全 JSON、schema version、opaque internal ref、SHA-256 内容摘要、字节大小。
  - [x] 定义 QueryResultSnapshot 及 SnapshotItem，至少冻结 workspace/conversation/actor scope、来源 run/tool result、查询语义、`captured_at/expires_at`、`coverage=complete_set`、稳定排序和完整安全行事实。
  - [x] M1 SnapshotItem 使用明确类型，最小字段只包含内部 opaque `snapshot_item_ref`、安全 `display_label` 与 `discovered_at`；按 `discovered_at DESC + snapshot_item_ref ASC` 排序。POC、IP、在线状态、业务系统和负责人等统一行事实留到 Stories 1.5/1.8/1.9，不得用 `map[string]any` 提前伪造或扩范围。
  - [x] `captured_at` 来自可注入的受控 clock，`expires_at` 由版本化 freshness policy/config 计算；不得硬编码 TTL。M1 重启恢复只验证未过期快照，过期与引用解析留到后续 Story。
  - [x] 定义 M1 所需的当前聚合与最小 FactEvent/sequence；状态更新与事件追加必须共享事务语义，但不实现完整 Event Sourcing。
  - [x] 更新 repository port，支持创建查询事实、读取当前结果和重启恢复；禁止通用 `Save(any)`、摘要反向拼装结果或第二套页面 DTO。
  - [x] 对全部安全材料做边界校验；外部失败只保存脱敏 typed error，不创建伪空快照。

- [x] Task 3：用 Greenfield SQLite 实现唯一持久事实源（AC: 01、02、04）
  - [x] 全新数据库直接创建目标 schema epoch；删除活动路径中的 `addColumnIfMissing` 等旧库兼容逻辑。
  - [x] 非目标 epoch 返回 `unsupported_schema_epoch`、readiness=false；不得自动迁移、删库、dual write 或读取旧 v1 run。
  - [x] 统一 database handle / connection policy，验证 `journal_mode=wal`、每连接 `foreign_keys=on`、配置化 busy timeout、固定 pool/transaction、integrity/foreign-key check 和可写探针。
  - [x] 把 `sqlite_version()` 解析为数值元组比较，覆盖 `<3.51.3` 与 malformed 值的 readiness fail-closed；至少获取两个真实 `sql.Conn` 分别断言 foreign keys / busy timeout，并验证 WAL 返回值，不能用单连接测试替代连接策略。
  - [x] 在一次事务中持久化 StructuredResult、QueryResultSnapshot、聚合状态和对应 FactEvent sequence；任何一步失败均不得发布部分事实。
  - [x] 在 StructuredResult、QueryResultSnapshot、聚合状态和 FactEvent 各写点做故障注入，断言事务完全回滚、无可见 product projection，失败 sequence 不得成为已发布事实。
  - [x] 用关闭并重新打开同一临时数据库的测试证明结果恢复；恢复路径不得重新执行 capability。

- [x] Task 4：以 Eino-first 方式执行唯一 mock 只读 capability（AC: 01、02）
  - [x] 复用 capability registry 基础结构，只注册一个 `read-only/no-side-effect/no-approval` 的新增漏洞能力；schema、risk、timeout 和 fixture adapter 均由 metadata/端口连接。
  - [x] 为 mock ChatModel 提供确定性 tool call，使 `adk.ChatModelAgent + Runner + ToolsConfig` 完成一次工具循环；不得使用 `if query == ...`、关键词匹配或 Action API `capability_hint` 代替 Eino 选择。
  - [x] fixture adapter 只产生已脱敏 candidate，Safety Gate 校验 schema 和安全材料后才能进入 facts；raw fixture/provider 结构不得进入 product、LLM context、HTTP、SSE 或 Web。
  - [x] 测试内构造 malformed、secret 和 provider-locator candidate（不增加第四个 fixture case），断言 Safety Gate 拒绝且 facts/SQLite/product/SSE 均无落库或泄漏。
  - [x] 验证三个 case 都通过同一个 capability 和同一映射链；失败 case 不生成“没有待派发漏洞”。
  - [x] composition root 固定使用 mock LLM 与 fixture adapter；不装配真实 OpenAI-compatible provider、FOBrain、MCP 或任何写 capability。

- [x] Task 5：从同一 Product Facts 投影 HTTP、Workbench 与最小 SSE（AC: 01、02、04）
  - [x] `product` 只从 repository 生成查询编号、摘要、数量、观察时间和确定性状态；禁止从摘要重建 StructuredResult。
  - [x] `httpapi` 只保留 M1 所需的 message、current/run snapshot、stream 与 health/readiness 边界，统一验证输入和返回 error envelope。
  - [x] Action、resume、lifecycle、replay 及真实 provider 装配不得成为活动 M1 路由；后续代码即使仍在仓库中也不能被注册或出现在生成契约。
  - [x] 复用统一 SSE encoder，事件只携带安全 product projection；使用持久 sequence 生成稳定 event id，M1 只验证同次查询与重启后读回，不提前实现 Story 1.11 的完整重连代数。
  - [x] Workbench view、HTTP snapshot 和 SSE 对同一结果的查询序号、摘要、数量、观察时间必须逐值一致。

- [x] Task 6：实现唯一浅色桌面 Workbench 页面（AC: 01、02、03、04）
  - [x] 前端通过 TanStack Query 调用 Go API，Composer 真正提交固定查询；删除活动页面的 fixture switcher 和静态 fixture 直读。
  - [x] 使用 `nav` / `main` / `aside` 构成 `232px / minmax(600px, 1fr) / 360px`、最小宽度 `1180px` 的浅色三栏；提供 skip link 和清晰焦点。
  - [x] 左栏提供与会话分组独立的固定“操作记录”入口，但点击后只显示当前未启用/空态，不请求历史；不得硬编码虚假会话数据。
  - [x] 中栏只显示聊天、查询过程和一张结果卡；结果卡显示查询序号、摘要、数量、观察时间，不显示内部 ref 或技术标签。
  - [x] 右栏只提供“事实”“执行记录”两个安全空态；删除/隐藏证据链、结构化结果、运行详情、审计等未授权技术 Inspector。
  - [x] 删除移动面板、移动断点、暗色 token、主题切换、审批/澄清/人员选择/写按钮、完整分页与“查看全部”入口。
  - [x] 验证 DOM 顺序、Tab/Enter/Space、可见焦点、1440×900 和 200% 文本缩放；低于最小宽度只能横向滚动或提示桌面宽度不足，不折叠为移动布局。
  - [x] 浏览器 spy 与后端计数共同证明点击“操作记录”时 history request 数为 0。

- [x] Task 7：建立单 fixture 纵向 E2E 与负向安全断言（AC: 01、02、03、04）
  - [x] 新增浏览器→Go message API→Eino mock tool loop→Safety Gate→SQLite→product projection→SSE→React 的真实纵向测试，不能只渲染前端 fixture。
  - [x] 分别执行 `resolved`、`empty`、`failed`，验证精确产品文案和状态机器可区分。
  - [x] 对 `resolved` 逐层断言 `discovered_at DESC + snapshot_item_ref ASC`，保证 fixture、StructuredResult、QueryResultSnapshot、SQLite 读回和投影使用同一稳定顺序。
  - [x] 关闭/重启后端并重新打开相同 run，验证投影逐值不漂移、工具调用计数不增加。
  - [x] 断言真实 LLM 调用数、真实 FOBrain 调用数、外部 mutation 调用数均为 0；写意图只能 fail closed，不能生成假审批。
  - [x] 扫描 HTTP、SSE、HTML、日志和截图，拒绝 raw payload、Authorization/token/cookie/credential、provider locator、内部 result_ref 和旧 runtime 类型。SQLite 的专用元数据列必须保存 opaque internal result_ref，但 StructuredResult 安全 JSON、业务投影列及其他可展示材料不得嵌入该引用；数据库泄漏扫描仍拒绝 raw payload、秘密、provider locator 和旧类型。

- [x] Task 8：收口旧原型活动路径与回归边界（AC: 01、02、03、04）
  - [x] 不创建并行 `walking_skeleton` runtime；M1 v2 必须替换唯一活动 facts/repository/projection/runtime 路径。
  - [x] 复用并加强 import-boundary、统一错误、SSE encoder、registry、React/Vite/TanStack Query/Zustand/Radix 工程骨架；不复制旧项目或现有原型的错误状态模型。
  - [x] 将旧 Action/HITL/MCP/真实 provider/24 工具/mobile 测试从当前 M1 suite 和产品出口移除或改为明确未到阶段；不得为了保持旧测试绿色添加兼容分支。
  - [x] 保证 `httpapi`、`execution`、`facts`、`product`、`capabilities`、`llm`、`observability`、`store/sqlite` 依赖方向符合 AD-01/AD-02。

- [ ] Task 9：执行门禁、记录证据并同步文档（AC: 01；AC: 02；AC: 03；AC: 04）
  - [x] 运行本 Story 的定向测试、五条正式 acceptance commands、SQLite/execution race 和 import boundary。
  - [x] 在固定 Go/Node/canonical linux/amd64 工具链产生新的 clean GitHub Actions 证据；Story 1.1 的历史 run 不能证明本次实现。
  - [x] 新建 `docs/acceptance-records/story-1-2-m1-walking-skeleton-YYYY-MM-DD.md`，记录命令、环境、三态、重启恢复、泄漏扫描和三类调用计数。
  - [x] 同步 schema、fixture、OpenAPI、generated contract、ADR、实施/验收计划、G-ARCH-V2 的 M1 子集证据和 sprint status；Story 1.2 最多记录 `M1 read subset=PASS/evidence added`，不得直接把全局 G-ARCH-V2 置为 PASS（其完整解除仍受 Story 1.12 与后续 Action 子集约束）。只有全部 AC/命令通过后才标记本 Story `done`。

## Dev Notes

### 实施边界

- 本 Story 只有“一 fixture bundle、一只读 capability、一页面、一结果”。`resolved/empty/failed` 是同一能力的测试 case，不是三项产品能力。
- FR-04 的“全部”仅表示 fixture 中的完整 `complete_set`，不表示已经实现真实 FOBrain 全量分页。
- M1 只建立 Product Facts、StructuredResult、QueryResultSnapshot 和最小 FactEvent/SSE 的只读子集。ActionDraft、ActionItem、ActionAttempt、VerificationEvidence、approval/resume/result、人员列表、写入、完整操作记录均留到 M4 以后。
- Story 1.9 才实现每页 100 条的完整明细工作区；Story 1.10 才实现多查询 result_ref 指代；Story 1.11 才实现完整 SSE 断线恢复。
- OQ-08 未关闭，产品出口不得显示技术信息；右栏在 M1 只显示“事实／执行记录”安全空态。
- 写意图不需要关键词分类器：M1 registry 中不存在写 capability，任何无法安全选择的意图都不执行外部操作。

### 现有代码现实与文件策略

当前仓库是旧 v1/Phase 8 原型，不是本 Story 的部分完成实现：

- `facts.ToolResult` 只保存摘要引用，`product` 会从摘要拼装结果；不存在生产 `QueryResultSnapshot`。
- `execution.ChatModelRunner` 未配置 Eino `ToolsConfig`，聊天路径不会选择工具；旧工具路径依赖 Action API hint。
- SQLite 是 schema version 1，并通过 `addColumnIfMissing` 兼容旧库；未启用 WAL、epoch/readiness 或完整结果持久化。
- SSE sequence 在读取时重排，不是持久事实游标。
- Web 直接读静态 fixture，存在暗色 token、移动面板、虚假会话、开发 fixture switcher 和技术 Inspector。
- HTTP/runtime 暴露 Action、resume、replay、MCP、真实模型和真实 FOBrain 等超出 M1 的入口。

开发时按以下策略处理：

- **复用框架**：`internal/einoapp/architecture/import_boundary_test.go`、`capabilities/registry.go`、`httpapi/response.go`、`httpapi/sse.go`、SQLite driver、React/Vite/TanStack Query/Zustand/Radix 项目结构、schema/OpenAPI/contract generator 脚本框架。
- **重写活动契约与路径**：`facts` model/repository、`store/sqlite` migration/repository、Eino mock tool loop、`product` projection、message/snapshot/SSE HTTP、composition root、Web API/Composer/Shell/Sidebar/right panel、generated contract、Playwright setup。
- **新建**：M1 v2 schemas、单一新增漏洞 fixture bundle、Greenfield epoch/readiness、专用 walking-skeleton E2E 和 Story 验收记录。
- **停用而非兼容**：v1 contract、旧库迁移、Action/Resume/Replay 活动路由、真实 provider 装配、24 工具当前门禁、静态 fixture 页面、暗色/移动项目和后续阶段视觉用例。

不得新建第二套并行 runtime 来绕开旧原型，也不得为了让旧测试继续通过而保留 v1/v2 union、dual write、旧路由 adapter 或 mobile fallback。

### Architecture Guardrails

- Eino 只允许由 `execution`、`capabilities`、`observability` 和 checkpoint adapter import；facts/product/httpapi/web 保持 Eino-free。
- `execution` 拥有 run 命令和状态迁移；`facts` 拥有领域值与 repository port；`product` 只读投影；`httpapi` 只做传输边界。
- StructuredResult 是唯一事实材料。raw provider/fixture payload 只能存在于 adapter 边界，不能进入产品层、LLM context、audit、replay、HTTP、SSE 或前端。
- opaque internal result_ref 是 facts/SQLite 的权威元数据，只允许程序解析；不得嵌入 StructuredResult 安全 JSON、产品投影、日志、截图或模型上下文。
- Workbench、HTTP snapshot、SSE 和重启恢复必须读取同一 Product Facts projector，不得各自维护 DTO 或状态映射。
- 固定查询是验收输入，不是业务分支。能力选择必须来自 Eino tool call 与 capability registry metadata。
- schema version、capability id、tool id 可以集中为契约常量；provider URL、凭据、端口、timeout、budget、fixture case 等来自配置或测试装配。
- 所有外部输入、Eino/tool 输出和 resume-like payload 均视为不可信；M1 不存在真实 resume 或 mutation 路径。

### Library / Version Notes

- Go module language `1.26.0`、toolchain `1.26.5`；Node `24.18.0`、npm `11.16.0`；Eino `v0.9.12`。
- Eino `v0.9.12` 官方 ADK 支持 `ChatModelAgent + Runner + ToolsConfig`，由框架处理 tool loop；项目仍负责事实、安全和产品投影。
- `modernc.org/sqlite v1.34.5` 内含 SQLite 3.46.0，处于 SQLite 官方 WAL-reset 缺陷影响范围。按 ADR 升级到 `v1.46.2`（SQLite 3.51.3）后才允许启用 M1 WAL。
- React 19.2、Vite 8、React Router、TanStack Query、Zustand 与 Radix 版本继续由 lockfile 固定；本 Story 不做前端依赖升级。
- OpenAPI 3.1 Schema Object 基于 JSON Schema 2020-12；不得使用 OAS 3.0 `nullable` 或开放对象逃逸严格契约。
- 当前本机 Node/npm 不是 Story 1.1 固定版本，只可做预检；最终 PASS 必须来自 canonical linux/amd64 工具链。

### Testing Standards

正式验收命令：

```bash
node scripts/eino_workbench_schema_validate.mjs
go test ./...
npm --prefix web/eino-workbench run typecheck
npm --prefix web/eino-workbench test
npm --prefix web/eino-workbench run browser-test
```

实现期间还必须运行：

```bash
go list -m modernc.org/sqlite
go test ./internal/einoapp/facts ./internal/einoapp/product ./internal/einoapp/store/sqlite ./internal/einoapp/execution ./internal/einoapp/httpapi -count=1
go test -race ./internal/einoapp/store/sqlite ./internal/einoapp/execution -count=1
go test ./internal/einoapp/architecture -run 'ImportBoundary|ProductLayers' -count=1
npm --prefix web/eino-workbench run contract-test
npm --prefix web/eino-workbench run build
bash scripts/run_toolchain_baseline.sh
```

- SQLite 测试使用 `t.TempDir()`；重启恢复必须真正关闭并重新打开同一数据库文件。
- Playwright 只运行 canonical desktop 1440×900，单 worker；不得生成 mobile project/baseline。失败保留 trace/screenshot/video，产品截图不得包含内部 ref 或秘密。
- browser-test 必须启动 Go + Vite 的真实纵向环境，不能只验证静态 fixture。
- 测试重点是事实、状态、错误、安全边界与调用次数，不只覆盖 happy path 或 DOM class。
- 最终证据需包含 tracked、staged、untracked clean gate；既有无关 worktree 修改不能被回滚或混入交付。

### Exit Gate

只有以下条件同时满足才可进入 Story 1.3：

1. AC-01～AC-04 全部有机器证据。
2. 五条正式 acceptance commands 和适用定向/race/boundary 检查全部 PASS。
3. `resolved/empty/failed` 三态与重启恢复通过。
4. 真实 LLM、真实 FOBrain、外部 mutation 调用数均为 0。
5. 产品出口泄漏扫描为 0，且没有旧 runtime/v1 contract/mobile/dark/Action 活动入口。
6. 新 clean canonical GitHub Actions 证据、验收记录、文档和 sprint 状态一致。

## Project Structure Notes

- 正式 affected directories：`docs/schemas/`、`docs/fixtures/`、`docs/api/`、`internal/einoapp/{bootstrap,httpapi,execution,facts,product,capabilities,llm,observability,store/sqlite}/`、`cmd/eino-workbench/`、`configs/`、`web/eino-workbench/`、`scripts/`、`docs/acceptance-records/`。
- `providers/fobrain/` 不是 M1 实现目录；如测试需要证明未装配，只能做 import/registration 负向断言，不修改或调用真实 provider。
- contract 类型只在 schema/OpenAPI 生成链和 `web/eino-workbench/src/contracts/generated.ts` 定义，不在 feature 目录复制。
- frontend server state 只放 TanStack Query；Zustand 只保留当前 tab 等纯 UI 状态，不保存 facts、run lifecycle 或结果快照。
- 旧项目 `/Users/vick/Desktop/project/ai-agent` 本 Story 无需读取；不得复制其 runtime、API、数据库、DOM 或 CSS。

## References

- Story、范围、AC 与命令：[epics.md](../planning-artifacts/epics.md#story-12用单一安全-fixture-跑通只读-walking-skeleton)
- M1 分层与推进边界：[docs/07-implementation-plan.md](../../docs/07-implementation-plan.md#m1story-12-后端基础不变量)
- M1 验收重点：[docs/08-acceptance-plan.md](../../docs/08-acceptance-plan.md#m1-验收重点)
- 开发前门禁：[docs/pre-development-validation.md](../../docs/pre-development-validation.md)
- 架构 AD-01～AD-05、AD-13、AD-16～AD-18、AD-20、AD-24、AD-25：[ARCHITECTURE-SPINE.md](../planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14/ARCHITECTURE-SPINE.md)
- 桌面三栏与视觉约束：[DESIGN.md](../planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/DESIGN.md)
- IA、状态和可访问性：[EXPERIENCE.md](../planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/EXPERIENCE.md)
- SQLite WAL 安全基线：[2026-07-20-modernc-sqlite-wal-safety-floor.md](../../docs/adr/2026-07-20-modernc-sqlite-wal-safety-floor.md)
- Story 1.1 工具链结果：[1-1-固定可复现工具链与开工门禁.md](./1-1-固定可复现工具链与开工门禁.md)
- Eino v0.9.12 release：<https://github.com/cloudwego/eino/releases/tag/v0.9.12>（2026-07-20 复核）
- Eino ADK v0.9.12 source：<https://github.com/cloudwego/eino/tree/v0.9.12/adk>（2026-07-20 复核）
- SQLite WAL：<https://www.sqlite.org/wal.html>（2026-07-20 复核）
- SQLite foreign keys：<https://www.sqlite.org/foreignkeys.html>（2026-07-20 复核）
- modernc SQLite changelog：<https://gitlab.com/cznic/sqlite/-/blob/master/CHANGELOG.md>（2026-07-20 复核）
- OpenAPI 3.1.2：<https://spec.openapis.org/oas/v3.1.2.html>（2026-07-20 复核）
- JSON Schema 2020-12：<https://json-schema.org/draft/2020-12/json-schema-core>（2026-07-20 复核）
- React 19.2：<https://react.dev/learn/creating-a-react-app>（2026-07-20 复核）
- Vite 8：<https://vite.dev/guide/>（2026-07-20 复核）
- Playwright visual comparisons：<https://playwright.dev/docs/test-snapshots>（2026-07-20 复核）

## Dev Agent Record

### Agent Model Used

Codex（GPT-5）

### Debug Log References

- `docs/acceptance-records/pre-development-validation-1-2-2026-07-27.md`
- `docs/acceptance-records/story-1-2-m1-walking-skeleton-2026-07-27.md`
- `test-results/toolchain-baseline.log`（固定 Linux/amd64 工具链预检，本地生成且不提交）

### Completion Notes List

- 已实现唯一 M1 活动链：chat message → Eino ADK mock tool loop → Safety Gate → StructuredResult / QueryResultSnapshot → SQLite Product Facts → HTTP/SSE → React。
- 已覆盖 resolved / empty / failed、稳定排序、原子回滚、重启恢复、产品泄漏扫描和三类真实调用计数为 0。
- 已收口 M1 router、OpenAPI、contract、Workbench 和当前测试套件中的旧 Action/HITL/MCP/真实 provider/24 工具/mobile 路径；未加入兼容或降级分支。
- 固定 Linux/amd64 工具链完整预检通过；clean GitHub Actions 因尚未提交且 GitHub 凭据失效仍待完成，Story 保持 `in-progress`。

### File List

- 契约与数据：`docs/schemas/*.v2.schema.json`、`docs/schemas/new_vulnerability_fixture_bundle.v1.schema.json`、`docs/fixtures/manifest.json`、`docs/fixtures/new-vulnerability-walking-skeleton.v1.json`、`docs/api/eino-workbench.openapi.json`、`web/eino-workbench/src/contracts/generated.ts`。
- 后端活动链：`cmd/eino-workbench/main.go`、`configs/eino-workbench.m1.example.yaml`、`internal/einoapp/{bootstrap/m1_config.go,capabilities/new_vulnerability_fixture.go,execution/m1_runner.go,execution/m1_service.go,facts/query_facts.go,facts/query_repository.go,httpapi/m1_routes.go,product/m1_projection.go,store/sqlite/query_repository_v2.go}`。
- 后端测试：上述 M1 文件对应测试、`internal/einoapp/execution/m1_restart_integration_test.go`、`internal/einoapp/store/sqlite/version_pin_test.go`、`internal/einoapp/architecture/import_boundary_test.go`。
- Workbench：`web/eino-workbench/src/app/router.tsx`、`src/features/workbench/{api,components,state,views}/`、`src/styles/`、`vite.config.ts`、`playwright.config.ts`、`tests/shell.spec.ts`。
- 门禁脚本：`scripts/eino_workbench_contract_generate.mjs`、`scripts/eino_workbench_schema_validate.mjs`、`scripts/m1_e2e_server.sh`、`scripts/m1_browser_test.sh`、`scripts/run_toolchain_baseline.sh`、根与 Workbench `package.json`。
- 文档与记录：`docs/{README.md,07-implementation-plan.md,08-acceptance-plan.md,pre-development-validation.md}`、`docs/adr/2026-07-20-modernc-sqlite-wal-safety-floor.md`、本 Story 两份验收记录、implementation readiness gate、sprint status。
- 移除：旧 M1 当前套件中的 Action/HITL/MCP/真实 provider/24 工具测试，以及静态 fixture、技术 Inspector、dark/mobile/旧视觉前端与 Playwright 用例。

### Change Log

- 2026-07-29：完成 Story 1.2 M1 实现、固定工具链预检并生成 clean GitHub Actions 后，完成 Task 9 并转为 `done`。
