# Story 1.2 M1 Walking Skeleton 验收记录

日期：2026-07-27

结论：`FIXED-TOOLCHAIN PREFLIGHT PASS / CLEAN GITHUB ACTIONS PENDING`

本记录证明 Story 1.2 的实现、定向测试和固定 Linux/amd64 工具链预检通过。由于变更尚未提交，当前结果不是 clean GitHub Actions 证据，Story 状态继续保持 `in-progress`，不得进入 Story 1.3。

## 范围与环境

| 项目 | 实际值 |
| --- | --- |
| 基线提交 | `b8f7e903c75bade747ae887225cacbe8b7db5a15` |
| 分支 | `codex/m1-readiness` |
| 固定工具链 | Story 1.1 的 `agent-platform-eino-toolchain:local`，Linux/amd64 |
| Go / Node / npm | `1.26.5 / 24.18.0 / 11.16.0`，由 `verify_toolchain.sh` 实际校验 |
| Eino | `v0.9.12` |
| modernc SQLite | driver `v1.46.2`；运行时 SQLite `>=3.51.3` 测试通过 |
| 产品范围 | 一个安全 fixture bundle、一个只读 capability、一个浅色桌面页面 |
| 禁止调用 | 真实 LLM=0、真实 FOBrain=0、外部 mutation=0 |

## 已执行命令

以下命令实际返回 0：

```bash
node scripts/eino_workbench_schema_validate.mjs
GOCACHE=<workspace-cache> go test ./... -count=1
npm --prefix web/eino-workbench run contract-test
npm --prefix web/eino-workbench run typecheck
npm --prefix web/eino-workbench test
npm --prefix web/eino-workbench run build
npm --prefix web/eino-workbench run browser-test
docker run --rm --platform linux/amd64 -e CI=false \
  -v <repository>:<repository> -w <story-worktree> \
  agent-platform-eino-toolchain:local bash scripts/run_toolchain_baseline.sh
```

固定工具链完整基线实际覆盖：

- schema、fixture、OpenAPI、generated contract 一致性；
- `go list -mod=readonly -m all`、`go mod verify`、全量 Go 测试；
- execution / SQLite race、M1 定向测试、`go vet ./...`、`go build ./...`、import boundary；
- TypeScript typecheck、5 个 Vitest 测试、Vite production build；
- Playwright desktop 1440×900：`resolved`、`empty`、`failed` 各 2 条，共 6 条真实纵向用例；
- `git diff --check`。

## AC 证据

| AC | 结果 | 机器证据 |
| --- | --- | --- |
| AC-01 | `PASS` | Eino `ChatModelAgent + Runner + ToolsConfig` 只选择 registry 中唯一只读 capability；StructuredResult、Snapshot、聚合和 FactEvent 原子持久化并同源投影 |
| AC-02 | `PASS` | 三个 fixture case 共用同一执行链；0 条精确显示“没有待派发漏洞”，失败精确显示“无法读取漏洞事实。此次请求不是空结果。” |
| AC-03 | `PASS` | desktop 三栏为 `232px / minmax(600px, 1fr) / 360px`；最小宽度 1180px；无 mobile/dark/write；操作记录 history 请求数为 0 |
| AC-04 | `PASS` | 真实关闭并重开同一 SQLite 文件后投影逐值一致，capability 调用仍为 1 |

## 安全与持久化证据

- SQLite 新 epoch、WAL 返回值、两个真实连接的 foreign keys / busy timeout、integrity / FK / 同连接回滚写探针均通过。
- StructuredResult、Snapshot、Aggregate、FactEvent 四个写点故障注入均完全回滚。
- malformed、secret、provider locator 与 fixture 尾随 JSON 在 Product Facts 前拒绝。
- HTTP、SSE、HTML 产品文本与 API body 扫描未发现 raw payload、token、credential、provider locator、internal result_ref 或旧 v1/runtime 类型。
- SQLite 二进制扫描未发现 raw payload、秘密、provider locator 或旧类型；opaque internal result_ref 只允许存在于专用元数据列。
- mock runner 计数：每次查询 mock model=1、tool=1、真实 LLM=0、真实 FOBrain=0、外部 mutation=0。

## 回归边界

- 当前 composition root 只装配 M1 config、Greenfield QueryRepository、单 fixture、单 capability、M1 service/projector/router。
- Action、resume、replay、history、MCP、真实模型和 FOBrain 未注册到 M1 router/OpenAPI/Workbench。
- 旧 Action/HITL/MCP/真实 provider/24 工具/mobile 测试已从当前 M1 suite 移除；未增加兼容、dual write、旧库迁移或 v1/v2 union。

## 未完成项

1. 变更尚未 commit/push，无法产生“当前实现对应的 clean GitHub Actions run”。
2. 本机 `gh` 的 GitHub 凭据失效；需重新认证后提交、推送并等待 `Toolchain Gate` 成功。
3. `npm ci` 报告 lockfile 中 4 个 high severity advisory；本 Story 未做前端依赖升级，该项不改变本次功能验收结论，但需独立依赖治理。

在新的 clean GitHub Actions run 成功并补充 run URL / commit SHA 前，本记录不能被解释为 Story `done` 或全局 `G-ARCH-V2=PASS`。
