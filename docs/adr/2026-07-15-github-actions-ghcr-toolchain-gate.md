# ADR：GitHub Actions 与 GHCR 工具链门禁

日期：2026-07-15

Status: Accepted

## 背景

Story 1.1 的 canonical 工具链门禁需要由唯一远程平台、唯一 CI 和唯一 registry 产生可审计证据。
原平台契约依赖另一套变量、runner service 和 artifact 传递语义，无法与 GitHub Actions/GHCR
共存而不形成双重事实来源，因此本次迁移是不兼容的平台替换。

## 决策

- GitHub 是唯一远程平台；GitHub Actions 是唯一 CI；GHCR 是唯一 canonical registry。
- `.gitlab-ci.yml`、GitLab变量、DinD service 和 GitLab dotenv artifact 语义全部删除。
- clean GitHub Actions run 必须以 `GITHUB_SHA` 构造并验证
  `ghcr.io/wanggang8/agent-platform-eino/toolchain:$GITHUB_SHA@sha256:$IMAGE_DIGEST`。
- 最终证据必须包含 GitHub Actions run URL，以及 `toolchain-evidence-$GITHUB_SHA` (30 days)
  artifact；artifact 包含 baseline log 与 Playwright report。
- GHCR 认证使用环境 credential helper；无 token 的 Docker config 与 mode `0700` helper 可落在
  `$RUNNER_TEMP`，token 只允许存在于实际 registry step 的当前进程环境。helper 的 `store` / `erase`
  必须拒绝，禁止凭据进入 Docker config、GitHub environment/path 文件、artifact 或报告。
- `GITHUB_TOKEN` 作为 opaque secret 处理：只要求非空并拒绝 ASCII 控制字符，不匹配前缀、长度、JWT
  分段或其他内部格式；helper 对 JSON 所需字符做安全转义。GitHub 自 2026-04-27 分阶段把 GitHub App
  installation token（包括 Actions `GITHUB_TOKEN`）迁移为 stateless `ghs_APPID_JWT`，因此依赖旧
  `[A-Za-z0-9_]+` 格式会错误拒绝官方 token。
- Product code、Eino runtime、Product Facts 和 Workbench UI 不受该 ADR 影响。
- 只有新的 ADR 才能改变平台；不得在故障时静默恢复双 CI 或 tag-only 验证。

## 后果

- workflow、builder、baseline identity、静态 validator、文档门禁和验收记录只消费 GitHub/GHCR
  契约，平台漂移由负向测试拒绝。
- 本地 linux/amd64 image ID 继续作为 preflight 与 visual migration 历史证据，但不能代替 GHCR
  digest、真实 Actions run 或 30 天 artifacts。
- 在真实 run URL、GHCR digest 和 artifacts 全部产生前，Story 1.1 必须保持
  `BLOCKED_PENDING_GITHUB_RUN`，不得填写模拟 PASS；这些证据现已产生，当前门禁为 PASS，下一步为 M-1。

## 拒绝方案

- 双 CI：拒绝。它会恢复两套变量、DAG、故障语义和门禁事实来源。
- tag-only 验证：拒绝。可变 tag 不能证明下游 baseline 使用的是构建 job 产出的不可变 image。
- 仅保留本地 image ID：拒绝。它不能证明远程 registry push/pull、GitHub-hosted runner 或 artifact
  retention。

## 回滚条件

只有新的 Accepted ADR 可以改变唯一远程平台、CI 或 canonical registry。GitHub Actions 或 GHCR
短期故障只允许保持门禁阻断并重试；不得静默恢复双 CI、旧平台变量、DinD service、dotenv artifact
或 tag-only 验证。若确需更换平台，新 ADR 必须同时定义迁移范围、不可变 image identity、证据保留、
安全权限和新的验收解除条件。

## 关联验收

- `docs/08-acceptance-plan.md` 的 `G-TOOLCHAIN` 唯一 PASS 算法。
- `docs/acceptance-records/story-1-1-g-toolchain-2026-07-15.md` 的当前阻断裁决。
- `.github/workflows/toolchain.yml` 与 `scripts/run_toolchain_baseline.sh` 的 GitHub-only 静态契约。

## 外部依据

- [GitHub Docs：GitHub token formats](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-authentication-to-github#githubs-token-formats)
- [GitHub Changelog：2026 installation token format rollout](https://github.blog/changelog/2026-04-24-notice-about-upcoming-new-format-for-github-app-installation-tokens/)
