# Architecture Spine 版本与现实核验

评审对象：`ARCHITECTURE-SPINE.md`  
评审日期：2026-07-14  
评审范围：命名技术与版本、Eino 能力边界、当前仓库/既有设计现实、brownfield 兼容性  
结论：**REJECT — 当前 Spine 不可直接作为 build substrate；修正 P0 现实差距并形成迁移/版本 ADR 后再审。**

## 执行摘要

Spine 的核心方向——Eino v0.9.12 稳定线、Product Facts 控制面、StructuredResult 安全边界、桌面优先 UX、OpenAPI 3.1 / JSON Schema 2020-12——大体有一手资料与仓库基础支持。尤其 Eino v0.9.12 确实提供 `Runner`、`ResumeWithParams`、interrupt/checkpoint 和 callback 能力，拒绝采用 v0.10 alpha 也与官方 release 状态一致。

但 Spine 把大量**目标态**写成 `[ADOPTED]` 和当前不变量，而 brownfield 仍是另一套可运行契约：自然语言 Runner 没有工具；审批恢复不是 Eino interrupt/resume；Product Facts 不保存完整结构化结果；没有 QueryResultSnapshot、ActionDraft、ActionItem/Attempt；Run 状态和 Action API 仍使用旧契约；SSE sequence 是查询时重排而非持久化事实序列。若按 Spine 直接开发，会在契约、数据库、恢复语义和 API 上发生未声明的破坏性替换。

运行时版本也不能按现文接受：`Node.js >=20` 既允许已 EOL 的 Node 20/25，又允许 Vite 8 不支持的 Node 20.0–20.18；Go 1.22 早已越过官方支持窗口；本机/既有验收记录使用的 Go 1.23.12 和 Node 25.8.1 在评审日同样已失去官方支持。SQLite 的 WAL、busy timeout 与连接级 foreign-key 保证在当前代码中也未落实。

## 核验基线

### 官方/一手资料

- Eino 官方 release 显示 v0.9.12 是 latest stable，v0.10.0-alpha.11 是 pre-release：[cloudwego/eino releases](https://github.com/cloudwego/eino/releases)。
- v0.9.12 源码包含 `Runner.Query`、`Runner.ResumeWithParams` 和 checkpoint store 接口：[adk/runner.go@v0.9.12](https://github.com/cloudwego/eino/blob/v0.9.12/adk/runner.go)；interrupt 事件与状态保存见 [adk/interrupt.go@v0.9.12](https://github.com/cloudwego/eino/blob/v0.9.12/adk/interrupt.go)。
- Eino compose v0.9.12 的 `Interrupt` / `ResumeWithData` / `StatefulInterrupt` 语义由官方 Go package 文档确认：[compose package](https://pkg.go.dev/github.com/cloudwego/eino/compose)。
- Go 官方仅支持最新两个 major release；评审日稳定线为 Go 1.25/1.26，因此 1.22/1.23 均已退出支持：[Go release policy/history](https://go.dev/doc/devel/release)。
- Node 官方表明 v20 于 2026-03-24 EOL，v25 于 2026-03-31 EOL；生产应用应使用受支持的 LTS/maintenance 线：[Node release schedule](https://nodejs.org/en/about/previous-releases)。
- npm registry 对 Vite 8.1.0 声明的 Node engine 为 `^20.19.0 || >=22.12.0`；仓库的 `>=20` 更宽。可复核命令见文末，包一手记录见 [vite@8.1.0 registry record](https://registry.npmjs.org/vite/8.1.0)。
- SQLite 默认 journal mode 是 DELETE；WAL 需显式启用，busy handler 默认为空；foreign keys 必须对每个 connection 启用：[WAL](https://sqlite.org/wal.html)、[busy handler](https://www.sqlite.org/c3ref/busy_handler.html)、[foreign keys](https://www.sqlite.org/foreignkeys.html)。
- `modernc.org/sqlite` 当前发布线为 v1.53.0；仓库仍是 v1.34.5：[modernc SQLite package](https://pkg.go.dev/modernc.org/sqlite)。

### 仓库事实

- `go.mod` 固定 Eino v0.9.12、modernc SQLite v1.34.5，module language 为 Go 1.22。
- `web/eino-workbench/package.json` 使用 caret ranges，Node engine 为 `>=20`；lockfile 才固定具体解析版本。
- 现有 Product Facts schema 只有 Run/Turn/ToolCall/ToolResult/PendingInteraction/AuditEvent；Run 枚举为 `created/running/waiting/succeeded/failed/cancelled/stopped`。
- 现有 Action API 是 `POST /api/workspaces/{workspace_id}/agent/actions` 加统一 resume endpoint，没有 Prepare/Confirm API 或 draft digest contract。
- 现有 Phase 6 验收记录明确承认“真实 Eino graph interrupt resume”尚未完成。

## 发现

### 1. [P0] Spine 把破坏性的目标事实模型伪装成已采用现实

Spine AD-05/07/10/11 和 ER 图要求 `QueryResultSnapshot`、`ActionDraft`、`ActionItem`、`ActionAttempt`、coverage、item refs、digest、attempt lineage（Spine 91–131、239–253）。当前 `facts.Snapshot`、`eino_product_facts.v1` 和 SQLite schema 都没有这些实体；全仓搜索也没有对应生产类型。现 schema 又是 `additionalProperties: false`，不能以兼容附加字段方式逐步落地。

这不是普通实现缺口，而是 Product Facts v1、SQLite schema v1、generated TypeScript contracts、OpenAPI、fixtures、replay/audit 和 migration 的联合破坏性变更。Spine 没有 companion ADR、schema v2 计划或旧数据迁移/双读边界，却把规则标为 `[ADOPTED]`。

要求：把这些决定标成 target/proposed，或先新增正式 ADR，定义 Product Facts v2、迁移顺序、旧 run 的读取策略、contract versioning 和验收门禁。

### 2. [P0] 当前 Product Facts 丢失 StructuredResult 数据，无法支撑冻结查询与逐条动作

Spine AD-03/05 声称 Product Facts 保存可冻结、可引用的完整安全事实。当前 Go 模型 `facts.StructuredResultRef` 只保存 `SchemaVersion/ResultRef/SafeSummary`；SQLite `tool_results` 也只保存这三类摘要。产品投影随后合成 `facts: []` 的空 StructuredResult，而不是恢复真实工具结构化数据（`internal/einoapp/product/projection.go:615-632`）。

因此当前持久化事实不能在刷新/重启后重建漏洞行、稳定排序、完整 item refs 或确认范围。Spine 必须明确这是一次存储模型替换，并定义安全 structured payload 的持久化格式、大小/索引策略、schema migration 和历史数据语义；否则 AD-03/05 不成立。

### 3. [P0] “自然语言工具编排必须进入 Eino Runner”与当前执行代码相反

Eino v0.9.12 有能力实现 AD-02，但当前 `ChatModelRunner` 创建 `ChatModelAgent` 时未提供 `ToolsConfig`，且 `MaxIterations: 1`（`internal/einoapp/execution/runner.go:88-102`），不会形成 ReAct tool loop。当前 `ToolLoopRunner` 接收显式 capability id 后直接调用 `InvokableRun`（`internal/einoapp/execution/tools.go:48-96`）；它只是使用 Eino tool interface，不是 Eino Runner 的模型工具选择/循环。

这使 Spine sequence diagram 的 `Eino Runner -> 选择 capability` 和 AD-02 成为未实现目标。要求给出从当前 explicit-hint runner 到真正 ChatModelAgent tools loop 的阶段迁移、事件映射、MaxIterations/budget、失败与回归门禁；在此之前不得称为现行不变量。

### 4. [P0] 审批“checkpoint”当前不是 Eino interrupt/checkpoint，Spine 混淆两种恢复语义

当前 approval 代码手工把 `{capability_id,input_text}` JSON 写入实现了 `compose.CheckPointStore` 的表，并生成名为 `eino_checkpoint:` 的 id；approve 后反序列化 payload，再直接调用 `RunApprovedCapability`（`internal/einoapp/execution/approval_tool.go:29-69,124-178`）。没有 Eino interrupt event、interrupt address、Runner checkpoint state 或 `Runner.ResumeWithParams`。

仓库自己的验收记录也明确写明“真实 Eino graph interrupt resume”仍未完成（`docs/acceptance-records/phase-6-approval-resume-2026-07-03.md:53-56`）。因此 AD-02、AD-08 和 sequence diagram 把当前 continuation store 描述成 Eino HITL 是事实错误。要求二选一并形成 ADR：迁移到真实 Eino interrupt/resume；或承认 deterministic Action approval 是项目 continuation protocol，不再把它称作 Eino checkpoint。

### 5. [P0] 持久化单调 sequence 不存在，现 SSE ID 会随快照内容重排

Spine AD-04/12/16 要求事实追加时获得持久化、单调 sequence，并用它支持 Last-Event-ID。当前只有 Turn 表持久化 sequence；tool/pending/audit/run 迁移没有全局 fact sequence。`streamEventsFromSnapshot` 每次查询都从 1 重新编号，并按“run、所有 turns、所有 tools、所有 pendings”的当前集合位置生成 event id（`internal/einoapp/product/projection.go:525-583`）。

新增一个较早排序的事实或改变集合数量会让既有 tool/pending 的 event id 改变；状态更新也复用同一位置 id。它不是可订阅的不可变增量日志，无法可靠实现 Last-Event-ID、去重或断线续传。要求先定义并迁移统一 fact/event sequence 与不可变记录表，再宣称 AD-04/12/16 已采用。

### 6. [P0] AD-18 声称的 SQLite WAL/外键/busy-timeout 不变量未被当前连接模型保证

Spine 说 SQLite“启用 WAL、外键、busy timeout、事务和可验证备份”（169–173）。当前 facts repository 和 checkpoint store 分别 `sql.Open` 同一 DSN；只各执行一次 `PRAGMA foreign_keys=ON`，没有 `journal_mode=WAL`、`busy_timeout`、连接初始化 hook 或连接池限制（`internal/einoapp/store/sqlite/repository.go:24-38`、`checkpoint_store.go:23-36`）。

`database/sql` 是连接池，SQLite 官方又规定 foreign keys 必须对每个 connection 启用，所以一次 `db.Exec` 不能证明后续池连接均受约束。默认 busy handler 为空，两个 pool 并发写还会增加 `SQLITE_BUSY` 风险。要求把连接参数变成可验证 DSN/hook，确认 WAL 返回值、设置 busy timeout/池策略，并增加多连接 FK 与并发事务测试；备份机制也需给出真实实现/验收证据。

### 7. [P1] `Node.js >=20` 同时违反 Vite engine 和受支持运行时政策

Spine/两个 package manifest 都声明 `>=20`。它允许 Node 20.0–20.18，而 Vite 8.1.0 明确要求 `^20.19.0 || >=22.12.0`；也允许奇数线 Node 21/23/25。评审日 Node 20 和 Node 25 均已 EOL，现有验收记录与本机使用 Node 25.8.1，不能作为继续验收的受支持环境。

要求固定受支持 LTS 范围，例如 `^22.12.0 || >=24`（是否允许 future majors 需另决策），并增加 `.nvmrc`/Volta/toolchain 或 CI matrix 作为实际门禁。仅写 `>=20` 不足以复现或保证安全支持。

### 8. [P1] Go 1.22 已退出支持，现有 Go 1.23.12 验收环境也已过时

Spine 列 `Go module language 1.22`，现有 acceptance records 使用 Go 1.23.12。Go 官方只支持最新两个 major；2026-07-14 支持线是 1.25/1.26。继续以 1.22 language line 和 1.23 toolchain 作为架构/验收基线，会漏掉已不再回补的安全与工具链修复。

Eino v0.9.12 的最低 Go 要求仅为 1.18，因此升级项目 toolchain 不被 Eino 阻挡。要求明确区分 language compatibility 与构建 toolchain，至少把 CI/发布工具链固定到受支持版本并重跑 CGO-free SQLite、Eino HITL、race/backup 验收。

### 9. [P1] Action API 二阶段协议与现有 OpenAPI/handler 不兼容

AD-08 和 Mutation convention 定义 `PrepareAction -> ConfirmAction`，确认绑定 draft version/digest。现有 OpenAPI 只有一个 `POST .../agent/actions` 和通用 `POST .../runs/{run_id}/resume`；`eino_action_request.v1` 没有 result_ref/draft，resume schema 没有 draft version/digest/actor binding（`docs/api/eino-workbench.openapi.json:85-126`、`docs/schemas/eino_action_request.v1.schema.json`、`eino_workbench_resume_request.v1.schema.json`）。

要求明确新协议是替换还是扩展：endpoint、HTTP status、approval_required body、draft retrieval/expiry、幂等冲突、旧 resume_ref 兼容期都需要先进入 OpenAPI/Schema ADR。否则“Workbench 与 Action API 共用同一协议”没有可实现合同。

### 10. [P1] Spine 的 Run 状态词汇会被现有封闭 schema 拒绝

AD-12 要持久化 `queued/running/waiting_approval/resuming/reconciling/terminal`；当前 Run 枚举是 `created/running/waiting/succeeded/failed/cancelled/stopped`，ActionResult 又是另一套 `accepted/running/waiting/completed/...`。现有 `eino_product_facts.v1` 封闭枚举不会接受 Spine 新值。

`terminal` 还是状态类别而非可解释终态，若真作为单值会丢失 succeeded/failed/cancelled。要求建立一份明确状态映射：哪些是 Run.status，哪些是 ActionItem/Attempt.status，哪些只是分类；通过 schema version 和迁移实现，不能把两套枚举并列写成已采用事实。

### 11. [P1] Capability Registry 契约不足以派生 Spine 宣称的全部门禁

AD-13 要 registry metadata 声明 freshness 与 verifier，并让 Action prepare/policy/provider routing 全部派生。现 `capability_catalog.v1` 只有 schema refs、risk、side effect、approval、timeout、credential policy、idempotency 等字段，没有 freshness/expiry policy、verifier/read-back capability、actor policy、batch/coverage 语义；Go `Registry.Register` 只校验 id/provider/tool name（`internal/einoapp/capabilities/registry.go:18-33`）。

因此 AD-05/09/10 的过期、二次校验和逐条回读无法从当前 registry 派生。要求先设计 catalog v2 及 validator，并说明 verifier 是 capability ref、策略对象还是代码 port；未完成前应把相应规则标为 proposed。

### 12. [P1] 单操作人 actor 绑定当前仍是文档愿望

AD-15 要 actor 只能来自 FOBrain credential 解析的 current user，并绑定 draft/approval/attempt/audit。当前 approval transition 把 audit actor 硬编码为字符串 `user`（`internal/einoapp/execution/approval_tool.go:212-218`），Pending/Action 请求也没有不可伪造 actor 绑定，现有 facts model 更没有 draft/attempt actor 字段。

即使 current-user read capability 已存在，它也不等于写入口已经把解析身份和权限快照原子绑定。要求在进入任何写域前补 actor/credential-context contract、prepare/confirm 双次校验与审计测试；否则 fail-closed 声明没有仓库证据。

### 13. [P2] 版本表是某次 lock 快照，不是可执行的版本治理策略

React 19.2.7、React Router 7.18.1、TanStack Query 5.101.2、Zustand 5.0.14 与 Eino 0.9.12 在评审日可由 primary registry/release 证实；但 Vite 8.1.0、Vitest 4.1.9 和多个 Radix 包已有更新 patch，modernc SQLite 已从 1.34.5 前进到 1.53.0。版本旧本身未证明有漏洞，但 Spine 没有给出选择旧 patch 的兼容/安全证据。

更重要的是，前端 manifest 全部使用 `^`，与“任何前端基础栈升级必须形成 ADR”和“package-pinned 1.x”不一致：`npm install` 可更新 minor/patch，只有 `npm ci` + committed lockfile 才能复现当前解析版本。要求把版本治理写成可执行规则：package manifest range 策略、lockfile/`npm ci` 门禁、patch 更新是否需要 ADR、定期更新与回滚证据；不要把“1.x”称作 pinned。

## 已证实、无需否定的决定

- **Eino v0.9.12 pin：支持。** 它是评审日最新 stable；v0.10 是 alpha。问题是当前代码没有使用它的完整能力，不是版本不存在。
- **Eino 能力边界：框架能力上可行。** Runner、targeted resume、interrupt/stateful interrupt、checkpoint store 和 callback 在 v0.9.12 源码中存在。
- **OpenAPI 3.1 / JSON Schema 2020-12：支持。** 当前 OpenAPI 文件实际为 3.1.2，schema 使用 2020-12 dialect。
- **React 19.2.7 / React Router 7.18.1 / TanStack Query 5.101.2 / Zustand 5.0.14：版本存在且 peer range 自洽。** 仍需用 lock/CI 固定实际解析版本。
- **桌面浅色三栏、1180px 最小宽度：有 UX Spine 与 accepted ADR 证据。** 这是产品范围决定，不属于版本现实问题。

## 重新评审前的最小门禁

1. 新增一份 migration ADR：Product Facts v2、QuerySnapshot/Draft/Item/Attempt、状态映射、API versioning、SQLite migration 和旧 run 读取策略。
2. 对 Eino 执行作出真实选择并给证据：自然语言 ReAct tools loop；approval 是真实 `Runner.ResumeWithParams`，或明确是项目 deterministic continuation，禁止混称。
3. 持久化统一 fact sequence，证明 SSE Last-Event-ID 在新增/更新事实后仍稳定。
4. 固化 SQLite connection policy（WAL、每连接 FK、busy timeout、pool/transaction、backup）并做多连接/重启测试。
5. 把 Go/Node 基线升级到评审日受支持运行时，修正 Vite engine，使用 CI + lockfile 重跑全套契约、Go、TS、browser 和 checkpoint tests。
6. 在写域准入前补 actor binding、catalog verifier/freshness、逐条回读和不确定写入恢复 contract。

## 可复核命令

```bash
go version
node --version
npm --version
go list -m -u -json github.com/cloudwego/eino github.com/eino-contrib/jsonschema modernc.org/sqlite
go list -m -versions github.com/cloudwego/eino github.com/eino-contrib/jsonschema modernc.org/sqlite
npm view vite@8.1.0 version engines --json
npm view vitest@4.1.9 version engines peerDependencies --json
npm view react-router-dom@7.18.1 version engines peerDependencies --json
rg -n 'QueryResultSnapshot|ActionDraft|action_draft|complete_set|partial_page|reconciling|manual_attention' docs internal web
rg -n 'PRAGMA|journal_mode|busy_timeout|foreign_keys|SetMaxOpenConns' internal/einoapp/store/sqlite
```

本评审未修改 Spine、代码、schema 或其他文档。

---

## 2026-07-14 P0 Recheck

### 结论：CHANGES_REQUIRED

本次只复核原 P0 范围：目标/现实区分、StructuredResult v2、Eino ReAct/HITL 迁移、持久化 SSE sequence、SQLite policy，以及 Go/Node target gate。更新后的 Spine 与迁移 ADR 已解决大部分原始阻断：Product Facts v2 和安全 StructuredResult JSON 已明确为目标态并受联合契约迁移约束；自然语言路径明确迁移到 Eino `ChatModelAgent + Runner + ToolsConfig`；真实聊天确认明确采用 Eino interrupt/checkpoint/`Runner.ResumeWithParams`；FactEvent sequence、SQLite connection policy、Go 1.26.x 与 Node 24 LTS 也已有明确目标。

仍有以下三个阻断项：

1. **[P0] 新增架构与工具链 gate 没有进入其声明的唯一权威门禁文档。** Spine 第 73 行声明 `implementation-readiness-gate.md` 是实施授权的唯一裁决来源，却只在 Spine 第 81–82 行新增 `G-ARCH-V2`、`G-TOOLCHAIN`。权威门禁文件的表格仍只包含 `G-READ-*`、`G-FACT-01`、`G-WRITE-*`、`G-SAFE-01`（第 15–24 行），通过规则也没有要求这两个 gate（第 26–31 行）。因此当前没有一份一致、可执行的 target gate 能阻止在 v2/toolchain 前置条件未满足时开始下一阶段。需要把两个 gate 的定义、状态、证据和通过规则同步到权威门禁文档，或撤回“唯一裁决来源”的声明并指定唯一真相来源。

2. **[P0] SSE 的当前态/目标态仍自相矛盾。** AD-24 已正确把持久化 sequence、`snapshot_sequence` 和稳定重连标为 `[ADOPTED TARGET]`，但 AD-16 仍标为 `[ADOPTED]`，并断言“SSE 使用稳定 event id/sequence/Last-Event-ID”（Spine 第 180–184 行）；这与当前查询时重新编号的实现现实冲突。Deferred 表又称“当前状态 + 不可变事实记录满足审计与 replay”（第 377 行），但不可变 FactEvent 也是 v2 目标而非当前事实。需要把 AD-16 的 SSE 子句和 Deferred 理由改成目标态表述，或拆分当前已采用的安全投影与尚未采用的持久化游标契约。

3. **[P0] HITL checkpoint 协议无条件覆盖 ActionDraft，与确定性 Action API 决策冲突。** 迁移 ADR 第 65 行明确确定性 Action API 使用项目 continuation、不强制伪造 Eino Agent；但第 70 行又规定“一个 ActionDraft”进入等待态前一律持久化 staged Eino checkpoint，Spine AD-22 第 220 行同样无条件要求 checkpoint。这样确定性 API 要么被迫制造无对应 Runner 的 Eino checkpoint，要么无法满足审批事务不变量。需要把协议拆清：聊天/Eino 路径使用 staged Eino checkpoint；确定性 Action API 使用等价的项目 continuation/attempt reservation 持久态，但不得伪造 Eino checkpoint；两者再共享同一 draft、approval、idempotency、execution 与 verification 事务约束。

### 已关闭的原 P0

- **StructuredResult v2：关闭。** ADR 明确安全 JSON 值、schema version、result_ref、SHA-256 和大小的不可变持久化，并禁止用摘要替代事实。
- **Eino ReAct 目标：关闭。** 自然语言路径已明确使用 Eino v0.9.12 Runner/ToolsConfig，当前手工 continuation 仅保留为 v1 历史兼容。
- **持久化 SSE sequence 的目标设计：关闭。** 同事务 FactEvent、Run 内单调 sequence、稳定 event id、`snapshot_sequence`、upsert/tombstone 与过期游标回退均已定义；剩余问题仅是当前/目标标签冲突。
- **SQLite policy：关闭。** 统一 handle/connection policy、WAL 验证、每连接 FK/busy timeout、pool/transaction、启动检查、恢复与备份验收均已进入 ADR。
- **Go/Node 目标版本：关闭。** Go 1.26.x CI/release 与 Node 24 LTS 已明确为迁移目标；剩余问题是 gate 尚未进入权威门禁文件。

本次 recheck 仅追加本报告，未修改 Spine、ADR、代码、schema 或其他文档。

---

## 2026-07-14 Final P0 Recheck

### 结论：PASS

本次仅复核上一轮三个剩余阻断项，现已全部关闭：

1. **权威 gate 已同步。** `implementation-readiness-gate.md` 已定义 `G-ARCH-V2` 与 `G-TOOLCHAIN` 的适用范围、`BLOCKED` 状态和通过证据，并在通过规则中明确两者是写动作实施前置条件；与 Spine 的唯一裁决来源声明一致。
2. **SSE/FactEvent 当前态与目标态已区分。** AD-16 已标为 `[ADOPTED TARGET]`，明确当前 SSE 查询时重排不满足稳定游标，迁移后才采用 AD-24；Deferred 也明确不可变 FactEvent 属于受 `G-ARCH-V2` 阻塞的 v2 目标。
3. **HITL continuation backend 已拆分。** Spine AD-22 与迁移 ADR 一致规定：聊天路径使用 staged Eino checkpoint，确定性 Action API 使用项目 continuation record 且不创建伪 Eino checkpoint；两者通过统一 `ContinuationRef` port 共享审批与事务不变量。

在限定复核范围内无剩余阻断项。本次 final recheck 仅追加本报告，未修改 Spine、ADR、readiness gate、代码、schema 或其他文件。
