# Story 1.2 开发前复核记录

日期：2026-07-27

结论：`ALLOW`

## 1. 工作项锁定

| 项目 | 裁决 |
| --- | --- |
| Milestone / Epic / Story | M1 / Epic 1 / Story 1.2 |
| 标题 | 用单一安全 Fixture 跑通只读 Walking Skeleton |
| WorkItemType | user-value |
| Sprint 状态 | `ready-for-dev` |
| 前置 Story | Story 1.1=`done` |
| Implementation Readiness | 2026-07-20=`READY` |
| G-TOOLCHAIN | `PASS` |
| G-ARCH-V2 | 全局 `BLOCKED`；Story 1.2 是获准建立 M1 read subset 证据的实现 Story，不得据此提前置全局 PASS |
| 其他 Gate | G-READ/G-FACT/G-SAFE/G-WRITE 保持当前状态；不构成单 fixture mock 只读链的前置，但禁止扩大到真实读写能力 |
| OQ | OQ-08 未授权技术 Inspector；Story 右栏仅保留“事实／执行记录”安全空态。没有其他开放 OQ 阻断 M1 |
| mutation 上限 | 真实 LLM=0；真实 FOBrain read=0；真实 FOBrain mutation=0；其他外部 mutation=0 |

Requirements、Architecture decisions、Inputs / Outputs、Scope、Non-goals、Affected directories、Acceptance commands、Acceptance Criteria 与 Exit gate 已按 [`Story 1.2`](../../_bmad-output/implementation-artifacts/1-2-用单一安全-fixture-跑通只读-walking-skeleton.md) 逐项锁定，与 [`epics.md`](../../_bmad-output/planning-artifacts/epics.md#story-12用单一安全-fixture-跑通只读-walking-skeleton) 一致。

## 2. 范围与输入输出

### 允许

- 一个版本化安全 fixture bundle，内部含 `resolved`、`empty`、`failed` 三个测试 case。
- 一个注册的新增漏洞只读 capability。
- Eino v0.9.12 mock ChatModel tool call 与单次 tool loop。
- Product Facts v2 的最薄 StructuredResult、QueryResultSnapshot、SnapshotItem、FactEvent/SSE 子集。
- Greenfield SQLite schema epoch、WAL、连接策略、readiness 与单一结果重启恢复。
- 一个浅色桌面页面、一张查询结果卡、固定但不查数据的“操作记录”入口和右栏安全空态。

### 禁止

- 真实 FOBrain、真实 LLM、MCP、ActionDraft、approval/resume、写入、完整历史、完整分页、完整引用恢复。
- 第二个 capability、24 工具齐套、移动端、暗色主题、自动派发、通知、工单或旧系统兼容。
- v1/v2 union、dual write、旧库自动迁移、旧 runtime/API/DOM adapter 或并行 walking-skeleton runtime。

### 固定产物位置

| 产物 | 目标位置 |
| --- | --- |
| StructuredResult v2 schema | `docs/schemas/tool.structured_result.v2.schema.json` |
| QueryResultSnapshot schema | `docs/schemas/eino_query_result_snapshot.v2.schema.json` |
| Product Facts v2 schema | `docs/schemas/eino_product_facts.v2.schema.json` |
| Workbench view v2 schema | `docs/schemas/eino_workbench_view.v2.schema.json` |
| SSE event v2 schema | `docs/schemas/eino_workbench_stream_event.v2.schema.json` |
| fixture bundle | `docs/fixtures/new-vulnerability-walking-skeleton.v1.json` |
| OpenAPI | `docs/api/eino-workbench.openapi.json` |
| generated TS contract | `web/eino-workbench/src/contracts/generated.ts` |
| Story E2E | `web/eino-workbench/tests/` 与 Story 专用启动/验收脚本 |
| 完成证据 | `docs/acceptance-records/story-1-2-m1-walking-skeleton-YYYY-MM-DD.md` |

上述名称沿用现有 schema 命名方式。实现若必须改变文件名或拆分 schema，需先同步 Story、manifest、OpenAPI 和本记录，不能静默产生第二套合同。

## 3. 官方资料与版本复核

### 实际命令

2026-07-27 执行：

```bash
go version
node --version
npm --version
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino github.com/cloudwego/eino-ext/components/model/openai github.com/eino-contrib/jsonschema modernc.org/sqlite
go list -m modernc.org/sqlite
# Story Task 1 升级后由真实连接执行 SELECT sqlite_version() 并断言 >= 3.51.3
```

结果：

| 组件 | 当前/可见状态 | 决策 |
| --- | --- | --- |
| Go | 本机 `go1.26.5 darwin/arm64`；module language `1.26.0` | 与固定 toolchain 一致，可做 Go 本地预检 |
| Node / npm | 本机 `25.8.1 / 11.11.0`；固定值 `24.18.0 / 11.16.0` | 本机结果不能作为 PASS；最终使用 Story 1.1 canonical linux/amd64 环境 |
| Eino | 当前 `v0.9.12`；模块索引可见稳定 `v0.9.13` 和 `v0.10 alpha` | 保持架构固定 `v0.9.12`；M1 所需 `ChatModelAgent + Runner + ToolsConfig` 已支持，不混入普通升级 |
| eino-ext OpenAI | 索引最新 `v0.1.13` | M1 禁止真实 LLM，不安装、不升级、不阻断 |
| eino-contrib/jsonschema | 当前/索引最新 `v1.0.3` | 保持 |
| modernc SQLite | 当前 `v1.34.5`；索引可见 `v1.46.2` 和更新版本 | 按 WAL 安全 ADR，Story 首项精确升级 `v1.46.2`；不直接追随最新 `v1.54.0` |
| OpenAPI / JSON Schema | OpenAPI 3.1.2 / JSON Schema 2020-12 官方合同 | 保持 3.1/2020-12；不引入 OAS 3.0 nullable 或 OAS 3.2 迁移 |
| React / Vite | lockfile 当前 React 19.2.7、Vite 8.1.0 | 官方当前主线支持现有 Vite 客户端构建；本 Story 不升级前端依赖 |
| Playwright | lockfile 1.61.1；canonical desktop 1440×900 | 只保留 desktop project；移除 mobile 作为 M1 验收项 |

### 官方依据及影响

- [Eino v0.9.12 release](https://github.com/cloudwego/eino/releases/tag/v0.9.12) 与 [v0.9.12 ADK source](https://github.com/cloudwego/eino/tree/v0.9.12/adk)：确认 M1 可使用 ChatModelAgent、Runner、ToolsConfig；Eino event 只作为 Product Facts 输入。
- [SQLite WAL](https://www.sqlite.org/wal.html)：WAL 必须检查返回 `wal`；官方记录 SQLite 3.7.0～3.51.2 的 WAL-reset 缺陷，M1 不允许在受影响版本启用 WAL。
- [SQLite foreign keys](https://www.sqlite.org/foreignkeys.html)：foreign keys 默认不保证启用，必须对每个真实连接验证。
- [modernc SQLite changelog](https://gitlab.com/cznic/sqlite/-/blob/master/CHANGELOG.md)：`v1.46.2` 对应 SQLite 3.51.3；采用 [`WAL 安全基线 ADR`](../adr/2026-07-20-modernc-sqlite-wal-safety-floor.md)。
- [OpenAPI 3.1.2](https://spec.openapis.org/oas/v3.1.2.html) 与 [JSON Schema 2020-12](https://json-schema.org/draft/2020-12/json-schema-core)：schema/OpenAPI/generated contract 必须同源。
- [React creating an app](https://react.dev/learn/creating-a-react-app)、[Vite guide](https://vite.dev/guide/) 与 [Playwright visual comparisons](https://playwright.dev/docs/test-snapshots)：维持 Vite 客户端工作台与固定桌面截图环境。

当前 `go.mod` 仍是 modernc v1.34.5，这是明确的待实现差距而不是替代方案：在 Story Task 1 完成升级、由真实连接回读 `sqlite_version()` 且测试通过前，不得启用 WAL writer 或把服务标记 ready。

## 4. 数据、fixture 与 provider 合同

| 项目 | 复核结论 |
| --- | --- |
| 数据来源 | 仅版本化合成安全 fixture，不读取真实 FOBrain，不需要 token/URL/credential |
| fixture 数量 | 一个 bundle；三个 case 共用一个 schema、一个 adapter 和一个 capability |
| `resolved` | 至少两条不同 `discovered_at`，并含同时间 tie-breaker；固定 `discovered_at DESC, snapshot_item_ref ASC` |
| `empty` | 成功且 `coverage=complete_set`、items=[]；产品文案“没有待派发漏洞” |
| `failed` | 脱敏 typed error；不得生成 QueryResultSnapshot 或空结果文案 |
| SnapshotItem | 封闭类型：内部 opaque `snapshot_item_ref`、安全 `display_label`、`discovered_at`；不使用开放 `map[string]any` |
| 时间 | `captured_at` 来自注入时钟；`expires_at` 由版本化 freshness policy/config 计算，M1 不硬编码 TTL |
| 敏感边界 | raw fixture/provider shape 停留在 adapter；产品出口无 raw payload、secret、locator、internal result_ref 或旧类型 |
| 内部引用 | opaque result_ref 必须存在于 facts/SQLite 专用元数据列，不嵌入 StructuredResult 安全 JSON、产品投影、模型上下文、日志或截图 |
| 外部调用 | mock LLM/tool + local SQLite；真实 LLM、FOBrain read、FOBrain mutation 均为 0 |

fixture 的具体字符串与时间值可以在实现时选择，但必须无真实公司/IP/人员数据、可由 schema 严格校验，并覆盖上述排序、空、失败和 Safety Gate 负向断言。

## 5. 架构逐项复核

下表评估 Story 目标路径。当前 v1/Phase 8 原型的偏差已在 Story 的“现有代码现实与文件策略”中逐项定位，并属于本 Story 必须替换的活动路径；任何偏差若在 runtime writer 或产品入口启用时仍存在，Story 立即失败。

| 问题 | 目标裁决 | 约束/证据 |
| --- | --- | --- |
| Workbench 与未来 Action API 是否只消费同一 Product Facts？ | 是 | M1 只启用 Workbench，但 v2 facts/projector 是唯一事实源；不建立 Workbench 专用事实 |
| 工具结果是否只以 StructuredResult 进入产品层？ | 是 | candidate→Safety Gate→按值 StructuredResult→facts；摘要不得反向拼装结果 |
| raw provider payload 是否停留在 provider/adapter 边界？ | 是 | fixture adapter 是 M1 唯一 provider-like 边界；负测覆盖 raw/secret/locator |
| LLM 上下文是否只来自安全投影？ | 是 | mock Eino 输入只允许 safe Product Facts/context；不注入 raw fixture |
| 工具选择是否基于 registry/metadata/intent/policy？ | 是 | 只注册一个 read capability，由 Eino tool call 选择；禁止 query/关键词/tool id 分支 |
| 分层依赖是否清晰？ | 是 | 复用并扩展 AST import boundary；composition root 是唯一具体装配点 |
| Eino event 是否只作为 Product Facts 输入？ | 是 | Web/HTTP/SSE 不暴露 Eino event/type |
| checkpoint/resume/SSE/audit/history 是否共用事实？ | 是/按阶段 | M1 只实现 facts-owned 最小 SSE；checkpoint/resume/audit/history 不启用、不造第二事实源 |
| 不可信输入是否校验并脱敏？ | 是 | HTTP、mock model/tool candidate、fixture case 与 SSE data 均在边界校验 |
| 是否不需要复制或兼容旧实现？ | 是 | Greenfield epoch、无旧 run、无旧 API/DOM、无 v1 adapter/dual write |

## 6. 当前代码事实与实施可行性

已按真实实现追踪入口、数据和状态：

- `cmd/eino-workbench/main.go` 当前可装配真实 LLM/MCP/FOBrain；M1 composition root 必须只装配 mock LLM、fixture adapter 和唯一 read capability。
- `internal/einoapp/execution/runner.go` 当前 ChatModelAgent 未设置 ToolsConfig；M1 需改为 Eino tool loop。
- `internal/einoapp/facts/model.go`、`repository.go` 和 `product/projection.go` 当前只保存/拼装摘要引用；M1 必须由完整安全 StructuredResult 与 QueryResultSnapshot 替代。
- `internal/einoapp/store/sqlite/migrations.go` 当前 schema v1 且有 `addColumnIfMissing`；M1 必须直接建立目标 epoch并拒绝旧 epoch。
- `internal/einoapp/store/sqlite/repository.go` 当前只启用 foreign keys；M1 必须补 WAL、busy timeout、多连接验证、integrity/writeability/readiness 和事务故障注入。
- `internal/einoapp/product/projection.go` 当前 SSE sequence 查询时重排；M1 必须使用持久 FactEvent sequence 的最小子集。
- `internal/einoapp/httpapi/routes.go` 当前暴露 Action/resume/lifecycle/replay；M1 活动 router/OpenAPI 不得注册这些后续入口。
- Web 当前从静态 fixture 读取，包含暗色/移动/技术 Inspector/虚假会话；M1 必须改为真实 Go API、浅色桌面三栏和安全空态。
- 现有 import boundary、统一错误、SSE encoder、capability registry、SQLite driver 和 React/Vite 工程骨架可复用；不得复制一套平行 runtime。

上述差距均有明确目标合同、目录、任务和验收断言；没有需要猜测的业务含义或未选择的架构分支。

## 7. 旧项目只读核对

Story 1.2 是合成 fixture walking skeleton，不需要旧项目产品行为或 FOBrain 字段映射。未读取 `/Users/vick/Desktop/project/ai-agent`；旧 runtime、API、数据库、DOM/CSS 和旧验收结果均不构成本 Story 输入或通过证据。

## 8. 验收命令与失败语义

正式命令均已存在：

```bash
node scripts/eino_workbench_schema_validate.mjs
go test ./...
npm --prefix web/eino-workbench run typecheck
npm --prefix web/eino-workbench test
npm --prefix web/eino-workbench run browser-test
```

补充命令：

```bash
go list -m modernc.org/sqlite
go test ./internal/einoapp/facts ./internal/einoapp/product ./internal/einoapp/store/sqlite ./internal/einoapp/execution ./internal/einoapp/httpapi -count=1
go test -race ./internal/einoapp/store/sqlite ./internal/einoapp/execution -count=1
go test ./internal/einoapp/architecture -run 'ImportBoundary|ProductLayers' -count=1
npm --prefix web/eino-workbench run contract-test
npm --prefix web/eino-workbench run build
bash scripts/run_toolchain_baseline.sh
```

失败语义：

- 任一当前 Story 命令非零退出、断言不一致、产品泄漏或三类真实调用计数不为 0，结论均为 `FAIL`。
- 未到阶段的 Action/MCP/real-provider/FOBrain/mobile 场景必须从当前 suite 移除或显式 `exit 2/SKIPPED`；它们不是 PASS。
- 本机 Node/npm 版本不符，因此本机前端结果只能是 preflight；最终 PASS 必须来自新 clean canonical GitHub Actions run。
- Story 1.1 的历史 run 只证明固定工具链存在，不能证明 Story 1.2 实现。

当前脚本仍验证旧 24 工具、Action 与 mobile 原型，这是 Story 已明确要求收口的测试基线差距，不造成开工语义不明。完成时若仍存在则 Story 失败。

## 9. 最终裁决

`ALLOW`

理由：

1. Story、前置、范围、数据来源、目标合同、分层、版本和验收标准已确定。
2. Story 1.2 是 G-ARCH-V2 M1 read subset 的授权实现/取证 Story；全局 Gate 的 `BLOCKED` 不授权后续能力，也不阻止在隔离范围内建立该子集。
3. modernc SQLite 的已知 WAL 风险已通过 ADR 固定到明确安全版本，升级是首项任务；升级完成前 readiness 必须 fail closed。
4. 唯一 fixture、唯一 capability、浅色桌面、同源 facts、三态、重启恢复和零真实调用均有可执行断言。
5. 当前 v1 原型差距已定位到具体入口和文件，Story 明确要求替换唯一活动路径，不需要兼容、降级或复制旧实现。

本裁决只允许开始 Story 1.2。它不表示 Story 已完成、不解除全局 G-ARCH-V2/G-READ/G-FACT/G-SAFE/G-WRITE，也不允许真实 LLM、FOBrain 或 mutation。
