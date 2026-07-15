# Story 1.1 G-TOOLCHAIN 验收记录

日期：2026-07-15

阶段：M-0 / Story 1.1

裁决：`BLOCKED_PENDING_GITHUB_RUN`

Story 状态：`in-progress`

下一阶段：不得进入 M-1

验证状态：`NEEDS_MORE_EVIDENCE`

## 裁决摘要

本记录只给出唯一 Story 门禁裁决，不产生“部分 PASS”。GitHub-only 静态契约已实现并通过本地静态
验证；canonical linux/amd64 local image、批准的 71 张 desktop visual evidence 与批准后完整本地
baseline 是有效 preflight 历史，但 local image ID 不能代替 GHCR digest。

`BLOCKED_PENDING_GITHUB_RUN`：GitHub-only 静态契约已实现，但尚未创建／推送公开仓库，未产生真实
GitHub Actions run URL、GHCR digest 与 30 天 artifacts。全部远程证据产生前，不得把静态或本地结果
改写为 `G-TOOLCHAIN PASS`。

## 当前 canonical 契约

| 输入／产物 | 精确值 | 证据角色 |
| --- | --- | --- |
| GitHub repository | `wanggang8/agent-platform-eino` | 唯一公开远程身份 |
| commit identity | 40 位小写十六进制 `GITHUB_SHA` | workflow、tag、baseline clean identity |
| canonical image | `ghcr.io/wanggang8/agent-platform-eino/toolchain:$GITHUB_SHA@sha256:$IMAGE_DIGEST` | 唯一最终 image identity；必须是真实 GHCR digest |
| workflow evidence | GitHub Actions run URL | 证明 build/verify 串行 run 来源 |
| artifact | `toolchain-evidence-$GITHUB_SHA` (30 days) | 包含 baseline log 与 desktop Playwright report |
| Go module language / toolchain | `1.26.0` / `1.26.5` | `.go-version` 与 `go.mod` exact declaration |
| Go linux/amd64 archive SHA-256 | `5c2c3b16caefa1d968a94c1daca04a7ca301a496d9b086e17ad77bb81393f053` | canonical image 下载校验 |
| Node.js / npm | `24.18.0` / `11.16.0` | `.node-version` 与根 manifest exact declaration |
| Node type declarations | `@types/node@24.13.3` | 本 Story 唯一批准的 exact dependency 变更 |
| Node linux-x64 tar.gz SHA-256 | `783130984963db7ba9cbd01089eaf2c2efb055c7c1693c943174b967b3050cb8` | canonical image 下载校验 |
| Playwright / Chromium | `1.61.1` / revision `1228`, version `149.0.7827.55` | desktop browser identity |
| Playwright MCR base | `mcr.microsoft.com/playwright:v1.61.1-noble@sha256:cf0daee9b994042e011bc29f20cdff1a9f682a039b43fcd738f7d8a9d3bcd9d6` | pinned external base；不是最终 GHCR digest |
| OS / platform / font policy | `ubuntu-24.04-noble` / `linux/amd64` / `playwright-v1.61.1-noble-bundled` | release-build 与 visual identity |
| Ubuntu archive snapshot | `20260708T000000Z` | 固定 apt 仓库状态与 source isolation |
| Build prerequisites | `build-essential` | C compiler 与 Go race baseline |

## GitHub-only 实现提交范围

在设计／计划提交之后，Tasks 1-4 的已评审运行实现为：

```text
ed5d540 build(toolchain): pin GitHub Actions runtime
c413744 build(toolchain): bind images to GitHub identity
eefe23d fix(toolchain): enforce GitHub baseline identity
a53ced3 test(ci): validate GitHub Actions toolchain gate
f37b0ae fix(ci): require complete toolchain lock schema
06cff08 ci(github): replace GitLab toolchain pipeline
67a5297 fix(ci): fail closed on unreadable entries
```

这些提交只证明静态契约实现范围，不是 clean GitHub Actions run、远程 commit、GHCR digest 或
artifact 证据。本次文档迁移不修改上述运行实现。

## 迁移前 GitLab 阻断（历史事实）

以下内容只保留为 2026-07-15 迁移前历史事实，已被 Accepted ADR 删除，不是现行命令、门禁变量或
解除条件：

- 当时仓库没有 remote，未产生 GitLab lint、clean pipeline、`CI_COMMIT_SHA`、pipeline URL、
  registry `tag@sha256` 或 pipeline artifacts，因此原门禁保持 BLOCKED。
- 当时 lock 包含 Docker CLI 与 Docker DinD digest，pipeline 使用 DinD service 和 GitLab dotenv
  artifact 传递 image；这些字段和语义现已删除。
- 当时静态 validator 通过只证明 `.gitlab-ci.yml` 契约，不证明 runner、push/pull 或 clean pipeline；
  该结论不能转用为当前 GitHub Actions PASS。

历史失败仍证明“静态配置和 local image 不能替代真实远程 run”这一安全边界；它不再要求恢复旧平台。

## 证据分级

| 级别 | 已有证据 | 能支持的结论 | 不能支持的结论 |
| --- | --- | --- | --- |
| 静态 | GitHub workflow、exact declarations、lock/Dockerfile、validator、负向 drift/bypass tests | GitHub-only 仓库契约会拒绝已覆盖的版本、identity、digest 和旧平台残留漂移 | image 已远程构建、GHCR push/pull 或 Actions run 已通过 |
| 本地 canonical | linux/amd64 image ID、构建期 race smoke、批准的 71 张 desktop baseline、完整 baseline log | image 可本地构建；visual migration preflight 已闭合 | GHCR digest、GitHub Actions run URL、30 天 artifacts |
| 最终 gate | 未产生 | 无 | `G-TOOLCHAIN PASS`、进入 M-1 |

## 实际运行命令与结果

| 命令／程序 | 结果 | 解释 |
| --- | --- | --- |
| `bash scripts/verify_toolchain_test.sh` | PASS | 正向 shim 通过，并拒绝版本、声明、旧平台残留和不安全输入 |
| `bash scripts/build_toolchain_image_test.sh` | PASS | builder、lock、GitHub identity、GHCR digest ref、0600 env file、baseline clean identity 与负向用例通过 |
| 显式 Go 1.26.5 `GOTOOLCHAIN=local go test ./scripts/validate_ci_config -count=1` | PASS | GitHub Actions validator 正向／负向测试通过；只属于静态 preflight |
| 显式 Go 1.26.5 `GOTOOLCHAIN=local bash scripts/validate_ci_config.sh` | PASS | 输出 `GitHub Actions config validated`；不代替真实 Actions run |
| `bash -n` Story 工具链 shell scripts | PASS | shell 语法通过 |
| 首次 `bash scripts/build_toolchain_image.sh --load` | exit 130 / 历史事实 | pinned MCR base layer 长时间无字节进展后人工终止；相同 pins 后续重试成功 |
| source isolation 加固后 `bash scripts/build_toolchain_image.sh --load` | PASS / 历史 preflight | 产生 local image ID `sha256:41f8317a3cc392bca3eedcf390e70b1c6ae4be0b4548cc4d7d3c4f200a29ea93`；它不能代替 GHCR digest |
| 批准后挂载 Git common dir 的完整 local baseline | PASS / 历史 preflight | desktop 17/17、contract smoke 与末尾 clean gate 通过；只闭合本地 visual migration 链 |
| `git remote -v` | 无输出 | 未配置公开 GitHub remote，未触发真实 Actions run |

## PASS 算法逐项裁决

| 必须条件 | 状态 | 客观证据／缺口 |
| --- | --- | --- |
| exact declarations | SATISFIED_STATIC | 权威版本文件、manifest、lock 和 verifier tests |
| verifier negative tests | SATISFIED_STATIC | 错版本、缺／短／大写 SHA、GitLab-only identity、workflow 漂移与 bypass 场景已覆盖 |
| canonical image build/push digest | BLOCKED_PENDING_GITHUB_RUN | local image ID 已产生；真实 GHCR digest 未产生 |
| clean GitHub Actions run | BLOCKED_PENDING_GITHUB_RUN | 公开 repository 尚未创建／推送，GitHub Actions run URL 未产生 |
| reproducible install/no dependency drift | SATISFIED_LOCAL / BLOCKED_PENDING_GITHUB_RUN | local canonical baseline 通过；clean GitHub Actions run 尚未执行 |
| schema/contract/OpenAPI | SATISFIED_LOCAL / BLOCKED_PENDING_GITHUB_RUN | local canonical baseline 通过；clean GitHub Actions run 尚未执行 |
| Go/test/race/vet/build/checkpoint/SQLite/boundary | SATISFIED_LOCAL / BLOCKED_PENDING_GITHUB_RUN | local canonical baseline 通过；clean GitHub Actions run 尚未执行 |
| TS/Vitest/stream/build | SATISFIED_LOCAL / BLOCKED_PENDING_GITHUB_RUN | local canonical baseline 通过；clean GitHub Actions run 尚未执行 |
| approved desktop visual baseline | SATISFIED_LOCAL / BLOCKED_PENDING_GITHUB_RUN | Vick 于 2026-07-15 批准；71 张 snapshot 已更新，local baseline 17/17 PASS |
| contract smoke | SATISFIED_LOCAL / BLOCKED_PENDING_GITHUB_RUN | local canonical baseline 通过；clean GitHub Actions run 尚未执行 |
| `toolchain-evidence-$GITHUB_SHA` (30 days) | BLOCKED_PENDING_GITHUB_RUN | baseline log 与 Playwright report 的远程 artifact 未产生 |

任一条件缺失即整体 `BLOCKED_PENDING_GITHUB_RUN`；`SATISFIED_STATIC` 与 `SATISFIED_LOCAL` 都不是
中间门禁 PASS。

## 直接阻断项与解除条件

1. `BLOCKED_PENDING_GITHUB_RUN`：GitHub-only 静态契约已实现，但尚未创建／推送公开仓库，未产生
   真实 GitHub Actions run URL、GHCR digest 与 30 天 artifacts。
2. 解除条件：在包含全部变更的 clean commit 上创建／推送公开 repository，取得真实 `GITHUB_SHA`，
   让 clean GitHub Actions run 成功完成 build/push、digest guard、pull 与唯一完整 baseline，并保存
   `ghcr.io/wanggang8/agent-platform-eino/toolchain:$GITHUB_SHA@sha256:$IMAGE_DIGEST`、GitHub Actions
   run URL 和 `toolchain-evidence-$GITHUB_SHA` (30 days)。

## 未产生的最终引用

- `GITHUB_SHA`：未产生远程 clean commit 证据。
- GitHub Actions run URL：未产生。
- `ghcr.io/wanggang8/agent-platform-eino/toolchain:$GITHUB_SHA@sha256:$IMAGE_DIGEST`：未产生。
- `toolchain-evidence-$GITHUB_SHA` (30 days)：未产生。
- artifact 内 `test-results/toolchain-baseline.log` 与 desktop Playwright report：未产生。

## 未覆盖风险

- pinned MCR、Go、Node tar.gz、snapshot apt、C compiler 与 race smoke 已在 local linux/amd64 image
  运行；GitHub-hosted `ubuntu-24.04` runner 的同一路径尚未证明。
- GHCR login、push、digest inspect、跨 job output、pull 与最小 `GITHUB_TOKEN` permissions 尚未真实运行。
- artifact 的 `if: always()` 上传、精确命名与 30 days retention 尚未由真实 run 证明。
- 当前工作区静态 clean 不等于 clean GitHub Actions run；任何后续提交都必须重新跑完整证据链。

## 最终结论

`BLOCKED_PENDING_GITHUB_RUN`。Story 1.1 与 Sprint 1.1 必须保持 `in-progress`，不得进入 M-1。

下一步只允许补齐唯一直接 blocker 的真实证据：创建／推送公开 GitHub repository，在包含全部变更的
clean commit 上取得成功的 GitHub Actions run URL、GHCR digest 和 30 天 artifacts。全部 PASS 算法
条件同时满足前，不得降低 AC 或把 local／static preflight 改写为 PASS。
