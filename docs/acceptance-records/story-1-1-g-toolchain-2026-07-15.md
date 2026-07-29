# Story 1.1 G-TOOLCHAIN 验收记录

日期：2026-07-15

阶段：M-0 / Story 1.1

裁决：`G-TOOLCHAIN=PASS`

Story 状态：`complete`

下一阶段：M-1

验证状态：`PASS`

## 裁决摘要

本记录只给出唯一 Story 门禁裁决，不产生“部分 PASS”。GitHub-only 静态契约、canonical
linux/amd64 image build/push、按 digest 验证的完整 baseline、批准的 71 张 desktop visual evidence
与 30 天 artifact 已在同一 clean commit 上取得真实远程证据。

run `29432627424` 对 commit `d174f79e5341dcce636d190aa0e510e4528a955a` 的 `toolchain-build`
与 `toolchain-verify` 均成功；GHCR 不可变 digest、desktop 17/17、contract smoke、末尾
tracked/staged/untracked clean gates 与可下载 artifact 同时满足 PASS 算法，因此裁决
`G-TOOLCHAIN=PASS`。

## 当前 canonical 契约

| 输入／产物 | 精确值 | 证据角色 |
| --- | --- | --- |
| GitHub repository | `https://github.com/wanggang8/agent-platform-eino` | 唯一公开远程身份 |
| commit identity | `d174f79e5341dcce636d190aa0e510e4528a955a` | workflow、tag、baseline clean identity |
| canonical image | `ghcr.io/wanggang8/agent-platform-eino/toolchain:d174f79e5341dcce636d190aa0e510e4528a955a@sha256:a68cf71dc506c21573dea7984c95a71ed5ff1837467481a62d4901e8cba7455f` | 真实 GHCR 不可变 image identity |
| workflow evidence | `https://github.com/wanggang8/agent-platform-eino/actions/runs/29432627424` | `toolchain-build` / `toolchain-verify` 串行 run |
| artifact | id `8350167050`，`toolchain-evidence-d174f79e5341dcce636d190aa0e510e4528a955a`，231272 bytes，expires_at `2026-08-14T16:35:24Z` | 未过期的 30 天 baseline log 与 desktop Playwright report |
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
| 最终 gate | clean GitHub Actions run、GHCR digest、完整 baseline 与 30 天 artifact | `G-TOOLCHAIN=PASS` | 无 |

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
| 迁移前 `git remote -v` | 无输出 / 历史事实 | 当时未配置公开 GitHub remote；现已解除 |
| GitHub Actions run `29432627424` | PASS | `toolchain-build` job `87411194696` success；`toolchain-verify` job `87411718738` success |
| 远程 canonical baseline | PASS | desktop 17/17、contract smoke PASS、末尾 tracked/staged/untracked clean checks PASS |
| 远程 artifact | PASS | id `8350167050`；`toolchain-evidence-d174f79e5341dcce636d190aa0e510e4528a955a`；`artifact_expired=false`；expires_at `2026-08-14T16:35:24Z` |

## PASS 算法逐项裁决

| 必须条件 | 状态 | 客观证据／缺口 |
| --- | --- | --- |
| exact declarations | PASS | 权威版本文件、manifest、lock、静态 verifier 与远程 baseline 通过 |
| verifier negative tests | PASS | 错版本、缺／短／大写 SHA、GitLab-only identity、workflow 漂移与 bypass 场景已覆盖 |
| canonical image build/push digest | PASS | `ghcr.io/wanggang8/agent-platform-eino/toolchain:d174f79e5341dcce636d190aa0e510e4528a955a@sha256:a68cf71dc506c21573dea7984c95a71ed5ff1837467481a62d4901e8cba7455f` |
| clean GitHub Actions run | PASS | `https://github.com/wanggang8/agent-platform-eino/actions/runs/29432627424`；HEAD SHA 与 commit identity 一致；build/verify 均 success |
| reproducible install/no dependency drift | PASS | 远程 canonical baseline PASS |
| schema/contract/OpenAPI | PASS | 远程 canonical baseline PASS |
| Go/test/race/vet/build/checkpoint/SQLite/boundary | PASS | 远程 canonical baseline PASS |
| TS/Vitest/stream/build | PASS | 远程 canonical baseline PASS |
| approved desktop visual baseline | PASS | Vick 于 2026-07-15 批准 71 张 snapshot；远程 desktop 17/17 PASS |
| contract smoke | PASS | 远程 contract smoke PASS |
| final tracked/staged/untracked clean checks | PASS | 远程 canonical baseline 末尾 clean gates PASS |
| `toolchain-evidence-$GITHUB_SHA` (30 days) | PASS | 真实 artifact 名称与 commit 一致；`artifact_expired=false`；expires_at `2026-08-14T16:35:24Z` |

所有条件已同时满足，整体门禁裁决为 `G-TOOLCHAIN=PASS`。静态与本地 preflight 仅作历史辅助证据，
最终裁决以上述真实远程 run 为准。

## 历史阻断解除

原 `BLOCKED_PENDING_GITHUB_RUN` 的直接阻断项已解除：公开 repository 已创建并推送；clean
GitHub Actions run 已完成 build/push、digest guard、pull 与唯一完整 baseline；GHCR digest 与
30 天 artifact 已产生。迁移前 GitLab 与本地 preflight 历史继续保留，但不再是现行阻断。

## 迁移首条真实远程证据（历史）

```text
repository=https://github.com/wanggang8/agent-platform-eino
commit=f975e258f0d5400fecbe81318083d0535eea4fb9
run_id=29427351390
run_url=https://github.com/wanggang8/agent-platform-eino/actions/runs/29427351390
toolchain_build_job=87393086944 success
toolchain_verify_job=87393589717 success
image=ghcr.io/wanggang8/agent-platform-eino/toolchain:f975e258f0d5400fecbe81318083d0535eea4fb9@sha256:1272821155d23cfb083cc4d0c1a7a77a96abd4ca00ae760d2339971a3bc6f080
artifact=toolchain-evidence-f975e258f0d5400fecbe81318083d0535eea4fb9
artifact_expired=false
artifact_expires_at=2026-08-14T15:22:11Z
desktop=17 passed
contract_smoke=passed
remote_clean_gates=passed
```

## 安全评审修复

整分支评审 C1/I1/I2/M1 修复将 GHCR 认证从会持久化 token 的 `docker login` 改为环境 credential
helper。Docker config 与 mode `0700` helper 副本只存在于 `$RUNNER_TEMP`，config 不含 token；actor 与
token 只注入实际 build/push、digest pull/run step，helper 固定拒绝 `store` / `erase`。token 作为
opaque secret 只检查非空与 ASCII 控制字符并做 JSON 安全转义，不解析 stateless JWT。

首个安全修复 run `29432152818` 对 commit `ff2e98f8fe27d81353ab983a225441c0a6b8b1fb`
在 helper configure 后因旧 `[A-Za-z0-9_]+` 假设拒绝 GitHub 2026 stateless installation token 而失败；
token 全程由 GitHub mask，失败记录只保留固定脱敏错误。依据 GitHub 官方格式变化，修复后的真实证据为：

```text
repository=https://github.com/wanggang8/agent-platform-eino
commit=d174f79e5341dcce636d190aa0e510e4528a955a
run_id=29432627424
run_url=https://github.com/wanggang8/agent-platform-eino/actions/runs/29432627424
toolchain_build_job=87411194696 success
toolchain_verify_job=87411718738 success
image=ghcr.io/wanggang8/agent-platform-eino/toolchain:d174f79e5341dcce636d190aa0e510e4528a955a@sha256:a68cf71dc506c21573dea7984c95a71ed5ff1837467481a62d4901e8cba7455f
artifact_id=8350167050
artifact=toolchain-evidence-d174f79e5341dcce636d190aa0e510e4528a955a
artifact_size_bytes=231272
artifact_expired=false
artifact_expires_at=2026-08-14T16:35:24Z
desktop=17 passed
contract_smoke=passed
remote_tracked_staged_untracked_clean_gates=passed
```

artifact 已下载到 ignored 验证目录复核，包含 `toolchain-baseline.log` 与 Playwright report；baseline
日志逐字记录 desktop `17 passed` 和 `contract smoke passed`。canonical baseline 在 smoke 后固定执行
tracked、staged、untracked 三类 clean gate；`Verify digest image` step success 证明三类静默 gate 均已通过。

## 未覆盖风险

- artifact 会在 `2026-08-14T16:35:24Z` 到期；当前 `artifact_expired=false`，到期后若需重新审计必须使用后续
  clean run 的新 artifact，不得伪造或回填已过期证据。
- 后续提交会触发新 run；新 run 必须保持 build/verify 全成功才能证明门禁未回归，但不为记录
  其 run ID 再制造自引用文档提交。

## 声明边界

- 通过/不通过/skipped blocking: 通过（G-TOOLCHAIN=PASS）
- 阻断 M1～M6 完成声明: 是，M1～M6 尚未实施/验收
- 阻断重构完成声明: 是
- 允许替换当前产品基线: 否
- report schema 校验结果: 不适用——本门禁无独立机器可读report schema；本次repository schema/contract/OpenAPI gates已通过
- 不得声明的能力：M1～M6 完成、27 个正式 Story 完成、真实 FOBrain 写动作已启用、
  生产环境可用或可替换当前产品基线。
- 关联 ADR：`docs/adr/2026-07-15-github-actions-ghcr-toolchain-gate.md`。

上述限制不阻断进入 M-1；它们只阻断超出 M-0 工具链门禁证据范围的完成与能力声明。

## 最终结论

`G-TOOLCHAIN=PASS`。Story 1.1 与 Sprint 1.1 已完成，允许进入 M-1。本结论仅表示 M-0
可复现工具链与开工门禁通过，不代表 M1～M6 或“产品实验完成”。
