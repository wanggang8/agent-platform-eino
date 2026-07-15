# Story 1.1 G-TOOLCHAIN 验收记录

日期：2026-07-15

阶段：M-0 / Story 1.1

裁决：`G-TOOLCHAIN=BLOCKED`

Story 状态：`in-progress`

下一阶段：不得进入 M-1

验证状态：`NEEDS_MORE_EVIDENCE`

## 裁决摘要

本记录只给出唯一 Story 门禁裁决，不产生“部分 PASS”。版本声明、负向 verifier、canonical
linux/amd64 image 和首次本地 baseline 已形成真实开发 preflight；但最终 PASS 算法中的 registry
image digest、clean GitLab pipeline、pipeline artifacts 和批准的 desktop visual evidence 仍缺失。
因此整体必须为 `BLOCKED`，本地 image、宿主临时 Go 1.26.5 或静态 CI 结果不得替代最终证据。

## 公开版本与 pinned 输入

| 输入 | 精确值 | 证据角色 |
| --- | --- | --- |
| Go module language / toolchain | `1.26.0` / `1.26.5` | `.go-version` 与 `go.mod` exact declaration |
| Go linux/amd64 archive SHA-256 | `5c2c3b16caefa1d968a94c1daca04a7ca301a496d9b086e17ad77bb81393f053` | canonical image 下载校验 |
| Node.js / npm | `24.18.0` / `11.16.0` | `.node-version` 与根 manifest exact declaration |
| Node type declarations | `@types/node@24.13.3` | 本 Story 唯一批准的 exact dependency 变更 |
| Node linux-x64 tar.gz SHA-256 | `783130984963db7ba9cbd01089eaf2c2efb055c7c1693c943174b967b3050cb8` | canonical image 下载校验；不依赖 base 中不存在的 xz-utils |
| Playwright / Chromium | `1.61.1` / revision `1228`, version `149.0.7827.55` | desktop browser identity |
| Playwright MCR base | `mcr.microsoft.com/playwright:v1.61.1-noble@sha256:cf0daee9b994042e011bc29f20cdff1a9f682a039b43fcd738f7d8a9d3bcd9d6` | pinned external base；不是最终 image digest |
| OS / platform / font policy | `ubuntu-24.04-noble` / `linux/amd64` / `playwright-v1.61.1-noble-bundled` | release-build 与 visual identity |
| Ubuntu archive snapshot | `20260708T000000Z` | Noble deb822 每个 Signed-By stanza 的固定 apt 仓库状态；构建删除其他 source 并显式只消费该文件 |
| Build prerequisites | `build-essential` | 提供 Go linux/amd64 race detector 所需 C compiler |
| Race capability smoke | `command -v cc` + `CGO_ENABLED=1 go test -race ./...` | image 构建期编译并运行最小 race test，成功后清理临时 module |
| Docker CLI | `docker:29.4.0-cli@sha256:bb21349a52c00b206ad8b5c03fa52023c741c4cf11f269d40d40b9ebaac73d96` | GitLab job image |
| Docker DinD | `docker:29.4.0-dind@sha256:4d2c6e334de4b26d492c0a8cc5438e3dbf1a02eee899fc0d4d39b96202c943a7` | GitLab service image |
| Eino | `v0.9.12` | 保持稳定 pin，本 Story 不升级 |

以上值来自仓库权威声明和 `build/toolchain/toolchain.lock`。本地 image ID 已产生，但最终 registry
digest 尚未产生；不得把 local image ID 或 MCR base digest 写成 registry 产物 digest。

## 实现提交范围

实现基线为 `f847458`，Story 1.1 工具链实现提交集合为 `f847458..8bc724b`（首个实现提交
`422a6b5`）：

```text
422a6b5 build(toolchain): pin greenfield runtime versions
2cb27cf fix(toolchain): sanitize verifier input errors
8634428 build(toolchain): add canonical release image
f61f4da fix(toolchain): centralize lock validation
bc47c53 ci(toolchain): run baseline in pinned image
d134f58 fix(ci): enforce canonical toolchain jobs
cfd4fb6 docs(toolchain): record blocked visual migration
e53cf18 docs(toolchain): separate visual and gate evidence
1ec4273 docs(toolchain): require full visual baselines
d60cc70 docs(toolchain): record g-toolchain blocked verdict
8f0af13 docs(toolchain): align gate command evidence
65aace9 fix(toolchain): pin build prerequisites
ade9c54 docs(toolchain): record pinned build prerequisites
3d4dac8 fix(toolchain): isolate snapshot package sources
4a6093e docs(toolchain): record canonical visual checkpoint
8bc724b fix(toolchain): stabilize desktop baseline
```

该范围只说明已审查的实现输入，不是 clean pipeline commit 证据。clean `CI_COMMIT_SHA` 未产生。

## 证据分级

| 级别 | 已有证据 | 能支持的结论 | 不能支持的结论 |
| --- | --- | --- | --- |
| 静态 | exact declarations、lock/Dockerfile、CI DAG validator、baseline sequencing、负向 drift/bypass tests | 仓库契约会拒绝已覆盖的版本／CI 漂移 | image 可构建、runner 可用、pipeline 已通过 |
| 开发 preflight | 官方 Go 1.26.5 darwin/arm64 归档校验后运行 validator tests/static | validator 在目标 Go 版本可编译并通过 | canonical linux/amd64、发布环境或 clean pipeline PASS |
| 本地 canonical | linux/amd64 image ID、构建期 race smoke、已批准的 71 张 desktop baseline、完整 baseline log | image 可构建；本地 visual migration A 链完整通过 | registry digest、clean pipeline |
| 最终 gate | 未产生 | 无 | `G-TOOLCHAIN PASS`、进入 M-1 |

## 实际运行命令与结果

| 命令／程序 | 结果 | 解释 |
| --- | --- | --- |
| `bash scripts/verify_toolchain_test.sh` | PASS | 目标 shim 正向通过，并拒绝 Go 1.23、Node 25、npm 错版、声明／CI／lock 漂移和不安全输入 |
| `bash scripts/build_toolchain_image_test.sh` | PASS | builder、16 字段 pinned inputs、tar.gz、snapshot、C compiler/race smoke、push/load 契约、dotenv、单 worker desktop baseline 顺序与末尾 clean gate 的静态／行为测试通过 |
| Node `SHASUMS256.txt` 精确查询 | PASS | 官方结果为 `783130984963db7ba9cbd01089eaf2c2efb055c7c1693c943174b967b3050cb8  node-v24.18.0-linux-x64.tar.gz` |
| Ubuntu snapshot Noble `InRelease` HEAD | PASS / HTTP 200 | `https://snapshot.ubuntu.com/ubuntu/20260708T000000Z/dists/noble/InRelease` 可访问；只证明公开输入存在，不证明 image 已构建 |
| 首次 `bash scripts/build_toolchain_image.sh --load` | exit 130 | MCR base metadata 按 digest 解析；base layer 超过 15 分钟无字节进展后人工终止，历史阻断保留 |
| source 隔离加固后重跑 `bash scripts/build_toolchain_image.sh --load` | PASS | 输出 local image ID `sha256:41f8317a3cc392bca3eedcf390e70b1c6ae4be0b4548cc4d7d3c4f200a29ea93`；实际依赖下载来自固定 Ubuntu snapshot，版本断言、cc、race smoke、Chromium identity 真实通过 |
| `docker image inspect` + 镜像内版本回读 | PASS | `linux/amd64`；Go `1.26.5`、Node `24.18.0`、npm `11.16.0`、`/usr/bin/cc`、Chromium `149.0.7827.55` |
| canonical `bash scripts/run_toolchain_baseline.sh` 首跑 | PARTIAL / browser exit 1 | browser 前全部非视觉检查通过；desktop 1/17 PASS、16/17 FAIL，产生 10 组 actual/diff；contract smoke 与最终 clean check 未执行，因此不是 baseline PASS |
| 临时 detached worktree 单 worker candidate generation | PASS / 17 tests | 首个成功构建的同 pins image 下 17/17 desktop tests 通过并生成全部 71 张 Linux candidate；未修改当前分支 screenshot，仍待 UX approval |
| 最新 hardened image 单 worker candidate recheck | PASS / 17 tests | 不运行 snapshot update，既有 71 张候选逐项匹配；确认 source 隔离重建未改变候选 |
| 固定单 worker 后 canonical baseline 重跑 | PARTIAL / browser exit 1 | browser 前全部非视觉检查再次通过；desktop 1/17 PASS、16/17 全为 screenshot diff、零 target crash；正式 baseline 未迁移，contract smoke 与最终 clean check 未执行 |
| 批准后的 canonical baseline 首次重跑 | PARTIAL / exit 128 | desktop 17/17 与 contract smoke 通过；worktree Git common dir 未挂载，末尾 clean gate 无法解析 `.git` 绝对指针，不计 PASS |
| 挂载 Git common dir 后批准 baseline 完整重跑 | PASS / exit 0 | 干净提交 `e455c3c`、同一 image；全部前置检查、desktop 17/17、contract smoke 与末尾 `git diff --check` 通过，本地 A 链闭合 |
| `GOROOT=<Go 1.26.5 toolchain root> PATH=<Go 1.26.5 bin> GOTOOLCHAIN=local go test ./scripts/validate_ci_config -count=1` | PASS preflight | 官方 darwin/arm64 archive SHA-256 为 `efb87ff28af9a188d0536ef5d42e63dd52ba8263cd7344a993cc48dd11dedb6a`；不是 canonical linux/amd64 证据 |
| `GOROOT=<Go 1.26.5 toolchain root> PATH=<Go 1.26.5 bin> GOTOOLCHAIN=local bash scripts/validate_ci_config.sh` | PASS preflight | 只证明静态 GitLab CI 配置契约，不是 GitLab lint 或 pipeline |
| `bash -n scripts/verify_toolchain.sh scripts/verify_toolchain_test.sh scripts/build_toolchain_image.sh scripts/build_toolchain_image_test.sh scripts/toolchain_lock.sh scripts/run_toolchain_baseline.sh scripts/validate_ci_config.sh` | PASS | Task 5 report 中列出的 Story shell scripts 语法通过 |
| `git diff --cached --check` | PASS / exit 0 | 检查当时已 staged 的 Task 5 文档 diff；不扩大为未执行的 unstaged 或历史提交范围 |
| `git diff --cached --name-only -- internal cmd web/eino-workbench/src` | PASS / 无输出 | staged Task 5 diff 未修改业务实现目录 |
| `git diff --name-only 1ec4273 -- internal cmd web/eino-workbench/src` | PASS / 无输出 | Task 5 相对 task base 的业务实现目录零 diff |
| `git diff --name-only f847458..1ec4273 -- internal cmd web/eino-workbench/src` | PASS / 无输出 | 业务实现目录零 diff |
| `git remote -v` | 无输出 | 仓库没有 remote；未调用无目标的 `glab ci lint`，也未触发 pipeline |
| `bash scripts/verify_toolchain.sh` | exit 1 | 非 canonical 宿主按预期 fail-fast：`node expected=v24.18.0 actual=v25.8.1`；不是 PASS，不改变整体 `BLOCKED` |

## PASS 算法逐项裁决

| 必须条件 | 状态 | 客观证据／缺口 |
| --- | --- | --- |
| exact declarations | SATISFIED_STATIC | 权威版本文件、manifest、lock 和 verifier tests |
| verifier negative tests | SATISFIED_STATIC | 错版本、声明漂移、CI 漂移、错误脱敏和 bypass 场景已覆盖 |
| canonical image build/push digest | PARTIAL_LOCAL / BLOCKED_FINAL | local linux/amd64 image ID 已产生并验证；registry `tag@sha256` 尚未产生 |
| clean GitLab pipeline | BLOCKED | 无 Git remote；GitLab lint、pipeline、clean `CI_COMMIT_SHA` 与 pipeline URL 均未产生 |
| reproducible install/no dependency drift | SATISFIED_LOCAL / BLOCKED_FINAL | canonical `npm ci` 未改变 lock/module files；clean pipeline 尚未执行 |
| schema/contract/OpenAPI | SATISFIED_LOCAL / BLOCKED_FINAL | canonical 首跑通过；未在 clean pipeline 重跑 |
| Go/test/race/vet/build/checkpoint/SQLite/boundary | SATISFIED_LOCAL / BLOCKED_FINAL | canonical 首跑通过；未在 clean pipeline 重跑 |
| TS/Vitest/stream/build | SATISFIED_LOCAL / BLOCKED_FINAL | canonical 首跑通过；未在 clean pipeline 重跑 |
| approved desktop visual baseline | SATISFIED_LOCAL / BLOCKED_FINAL | Vick 于 2026-07-15 批准；71 张正式 desktop snapshot 已更新，批准后的同 image 完整 baseline 中 17/17 PASS |
| contract smoke | SATISFIED_LOCAL / BLOCKED_FINAL | 批准后的本地完整 baseline 通过；clean GitLab pipeline 尚未执行 |

任一条件缺失即整体 `BLOCKED`；上表中的 `SATISFIED_STATIC` 不是中间门禁 PASS。

## 直接阻断项与解除条件

1. `BLOCKED_NO_REMOTE_OR_PIPELINE`：仓库没有 Git remote；未产生 GitLab lint、clean pipeline、
   `CI_COMMIT_SHA`、pipeline URL、registry `tag@sha256`、`toolchain-baseline.log` 或 Playwright report。
   配置 remote/runner 后，必须在包含全部变更的 clean commit 上运行真实 pipeline。

## 未产生的最终引用

- `CI_COMMIT_SHA`：未产生。
- GitLab pipeline URL：未产生。
- final `TOOLCHAIN_IMAGE=tag@sha256`：未产生。
- `toolchain-baseline.log` pipeline artifact：未产生。
- desktop Playwright report artifact：未产生。
- pipeline 保留的 canonical desktop actual/diff 与 UX approval：未产生；本地审查材料已产生，不能替代该最终证据。

## 未覆盖风险

- pinned MCR、Go、Node tar.gz、隔离后的 snapshot apt、`build-essential`、`cc` guard 与最小 race smoke 已在本地
  linux/amd64 image 真实运行；GitLab runner 中的相同端到端路径尚未证明。
- privileged DinD runner、registry push/pull、dotenv artifact 传递和 commit-bound digest 尚未实测。
- canonical `npm ci`、Go/TypeScript/contract、desktop 17/17、contract smoke 与最终 clean check 已在
  同一本地完整 baseline 通过；clean GitLab pipeline 尚未执行。
- 71 个 desktop Linux baseline 已获 human UX 批准并更新，本地 visual migration A 链已闭合。
- 当前仓库静态 clean 不等于 clean GitLab pipeline；任何后续提交都必须重新跑完整证据链。

## 最终结论

`G-TOOLCHAIN=BLOCKED`。Story 1.1 与 Sprint 1.1 必须保持 `in-progress`，不得进入 M-1。

下一步只允许补齐唯一直接 blocker 的真实证据：在包含全部变更的 clean commit 上运行真实 GitLab
pipeline 并保存 registry digest 与 artifacts。全部 PASS 算法条件同时满足前，不得降低 AC 或把
本地／静态 preflight 结果改写为 PASS。
