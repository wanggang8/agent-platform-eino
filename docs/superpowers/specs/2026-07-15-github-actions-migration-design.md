# GitHub Actions 与 GHCR 迁移设计

日期：2026-07-15

状态：已实施并通过真实 GitHub Actions 验收

## 目标

把项目的唯一远程代码托管、CI 门禁和 canonical image registry 从 GitLab 迁移到 GitHub：

- 创建公开仓库 `wanggang8/agent-platform-eino`。
- 删除 `.gitlab-ci.yml` 和所有 GitLab 专用执行契约，不保留双 CI。
- 使用 GitHub Actions 构建 canonical linux/amd64 image，并推送到
  `ghcr.io/wanggang8/agent-platform-eino/toolchain:<40-char commit sha>`。
- 验证 job 必须使用构建 job 输出的 `tag@sha256:<64-hex>` image 执行唯一
  `scripts/run_toolchain_baseline.sh`。
- 保存 baseline log 和 Playwright report，形成可审计的 `G-TOOLCHAIN` 最终证据。

## 非目标

- 不修改产品功能、Eino runtime、Product Facts、Fobrain provider 或 Workbench UI。
- 不增加 GitLab/GitHub 双平台兼容层，不保留降级路径。
- 不为 fork pull request 发放 package write 权限。
- 不把本地 image ID、tag、GitHub-hosted runner 状态或静态 workflow 校验冒充 registry digest。

## 已选方案与替代方案

采用完整 GitHub 迁移。GitHub 成为唯一远程平台，GitHub Actions 成为唯一 CI，GHCR 成为唯一
canonical image registry。

未采用：

- 双 CI：会产生两套变量、DAG、门禁和故障语义，与“不要 GitLab”冲突。
- 仅推送代码：无法产生 registry digest、完整 baseline 和 artifacts，不能解除 `G-TOOLCHAIN`。

## 仓库与分支

- GitHub owner：`wanggang8`。
- repository：`agent-platform-eino`。
- visibility：public。
- 本地 remote：`origin=https://github.com/wanggang8/agent-platform-eino.git`。
- 首次推送保留当前开发分支 `codex/story-1-1-toolchain`；不会静默改写或强推其他分支。
- workflow 对仓库内任意 `push` 和 `workflow_dispatch` 运行，使首次开发分支即可产生门禁证据。
- 不配置 `pull_request` package-publish 路径；fork PR 不获得 `packages: write`。

## GitHub Actions 架构

创建唯一 workflow：`.github/workflows/toolchain.yml`，包含两个串行 job。

### toolchain-build

1. 运行于明确标签 `ubuntu-24.04`。
2. workflow/job 权限最小化：`contents: read`、`packages: write`，其他权限为 `none`。
3. checkout、Docker setup、Buildx setup 等 action 必须固定完整 40 字符 commit SHA，不能使用浮动
   `@main`、`@master` 或只固定 major tag。
4. Docker Engine 固定为现有批准版本 `29.4.0`；Docker Setup action 用 exact `version` 安装，避免依赖
   GitHub runner 的浮动预装 Docker。Buildx/BuildKit 的 exact 版本与 action SHA 进入
   `build/toolchain/toolchain.lock`，由静态 validator 统一校验。
5. 使用环境 credential helper 认证 `ghcr.io`：Docker config 与 mode `0700` helper 副本只写入
   `$RUNNER_TEMP`，config 不含 token；`${{ github.actor }}` 与 `${{ secrets.GITHUB_TOKEN }}` 只注入实际
   build/push step。helper 仅从当前进程环境响应 `get`，固定拒绝 `store` / `erase`。token 必须作为
   opaque secret 处理：非空、无 ASCII 控制字符并做 JSON 安全转义，但不解析或匹配前缀、长度、
   JWT 段数或内部字符结构。
6. canonical tag 必须由 `GITHUB_REPOSITORY` 与 `GITHUB_SHA` 唯一构造，并转为小写：
   `ghcr.io/wanggang8/agent-platform-eino/toolchain:<sha>`。
7. `scripts/build_toolchain_image.sh --push` 验证 GitHub identity、构建并 push；成功后输出唯一
   `TOOLCHAIN_IMAGE=<tag>@sha256:<digest>`。
8. 通过 `$GITHUB_OUTPUT` 传递 digest image 给下游 job，不使用可编辑的模拟值。

### toolchain-verify

1. `needs: toolchain-build`，只消费 build job 输出的 digest image。
2. 权限为 `contents: read`、`packages: read`。
3. checkout 同一 `GITHUB_SHA`，用相同环境 credential helper 认证 GHCR，并拒绝非 `@sha256:` image；
   actor/token 只注入 digest pull/run step。
4. pull digest image，随后执行：

   ```text
   docker run --rm --platform linux/amd64 \
     -e CI=true -e GITHUB_SHA="$GITHUB_SHA" \
     -v "$GITHUB_WORKSPACE:/workspace" -w /workspace \
     "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh
   ```

5. 无论成功失败，都上传 `test-results/toolchain-baseline.log` 与
   `test-results/eino-workbench-playwright-report/`，artifact retention 固定 30 天。
6. 只有 build、digest guard、pull、完整 baseline 和 artifact upload 的必要步骤均满足，workflow run
   才能作为 `G-TOOLCHAIN` 证据。

## 契约迁移

- `.gitlab-ci.yml` 删除，不保留兼容 stub。
- `scripts/validate_ci_config` 改为验证 `.github/workflows/toolchain.yml` 的 GitHub Actions DAG、
  permissions、runner、SHA-pinned actions、GHCR tag、job output、digest guard、baseline command 和
  artifact paths。
- `scripts/validate_ci_config.sh` 保持中立名称，但默认输入改为 GitHub workflow。
- `scripts/verify_toolchain.sh` 不再读取 `.gitlab-ci.yml`，改为拒绝 GitLab 文件存在、GitLab变量和
  GitLab专用文案，同时要求 GitHub workflow 通过 validator。
- `scripts/build_toolchain_image.sh` 的 push 身份改为 `GITHUB_REPOSITORY` + `GITHUB_SHA`；错误信息不再
  出现 GitLab。
- `scripts/run_toolchain_baseline.sh` 的 CI clean identity 改为 `GITHUB_SHA`，同时继续要求 CI 前后
  tracked、staged、untracked 均为 clean。
- `build/toolchain/toolchain.lock` 删除 GitLab DinD 专用 pins，增加 GitHub runner label、官方 action
  commit SHA、Docker Engine、Buildx 和 BuildKit exact pins。版本常量只能由 lock/validator 消费。

## 安全边界

- 仓库公开，但本地 `configs/eino-workbench.local.yaml`、token、cookie、真实 provider凭据仍由
  `.gitignore` 隔离，不得进入 Git 历史、workflow、artifact 或日志。
- workflow 显式声明最小 `GITHUB_TOKEN` permissions。GitHub 未声明的权限均为 `none`。
- token 不得写入 Docker config、`GITHUB_ENV`、`GITHUB_PATH`、仓库文件、日志、artifact 或报告；
  `$RUNNER_TEMP` 只允许保存不含 token 的 helper 配置与 helper 可执行副本。
- GitHub 自 2026-04-27 分阶段启用包含 JWT 分隔符的 stateless installation token，且范围包括 Actions
  `GITHUB_TOKEN`；安全校验不得依赖旧 `[A-Za-z0-9_]+` token 结构，也不得解析 JWT 内容。
- 只在仓库内 `push` 或人工 `workflow_dispatch` 发布 GHCR image；不在 fork PR 上运行 write token。
- action 必须使用完整 commit SHA；外部 action 升级必须修改 lock、validator、测试和验收记录。
- registry tag 只作地址，最终事实必须为 `tag@sha256:digest`。
- workflow 不执行来自未信任输入拼接的 shell；repository、SHA 和 image ref 必须按固定正则校验。

## 错误与证据

- build 失败：不产生 digest output；verify 不运行；保留 workflow 日志。
- digest 缺失或格式错误：verify fail-fast，不 pull tag。
- baseline 失败：上传现有 log/report，workflow 失败，`G-TOOLCHAIN` 保持
  `BLOCKED_PENDING_GITHUB_RUN`。
- artifact upload 使用 `if: always()`，但 upload 自身不能把前序失败改写为 PASS。
- 最终验收记录保存 repository URL、commit SHA、workflow run URL、GHCR digest、artifact 名称和完整
  baseline 结论；不保存 token。

## 测试与验收

实施必须先写失败测试，再迁移实现：

- validator 正向覆盖完整 GitHub workflow。
- 负向覆盖：浮动 action tag、错误 permissions、非 `ubuntu-24.04`、缺 packages write/read、错误 GHCR
  路径、tag-only 传递、非 40 字符 SHA、build/verify 无依赖、artifact path/retention 漂移、重新出现
  `.gitlab-ci.yml` 或 GitLab变量；重新出现 `docker login`、helper/config/path 漂移、token 提升到 job
  scope 或实际 registry step 缺少 step-local actor/token。
- credential helper 正向覆盖旧 opaque token、新 `ghs_APPID_JWT` 和一般 opaque secret 的 JSON 转义；
  空 token、CR/LF/control character、非法 actor/server、`store`、`erase` 与未知 action 必须脱敏失败。
- builder 覆盖 GitHub identity、大小写归一化、push target、digest 输出和敏感错误脱敏。
- baseline 覆盖 CI clean checks 使用 `GITHUB_SHA`。
- 本地运行 toolchain verifier/tests、Go validator tests、shell syntax、`git diff --check` 和禁止产品目录
  diff。
- 创建公开 GitHub repository、添加 origin、推送当前分支，观察真实 GitHub Actions run。
- 真实 run 必须产生 GHCR `tag@sha256`、完整 baseline PASS 和 30 天 artifacts，之后才能把
  `G-TOOLCHAIN` 改为 PASS 并允许进入 M-1。

## 文档与决策同步

- 新增 ADR，明确 GitHub Actions/GHCR 取代 GitLab/DinD，且不保留兼容层。
- 更新 `docs/07-implementation-plan.md`、`docs/08-acceptance-plan.md`、
  `docs/tooling-and-reporting.md` 和 Story 1.1 验收记录。
- 原 GitLab 阻断记录作为历史事实保留，但最终裁决明确由 GitHub workflow 解除。

## 权威依据

- [GitHub：发布 Docker image 到 GHCR](https://docs.github.com/en/actions/tutorials/publish-packages/publish-docker-images)
- [GitHub：workflow permissions 与 GITHUB_TOKEN](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax)
- [GitHub：workflow artifacts](https://docs.github.com/en/actions/concepts/workflows-and-actions/workflow-artifacts)
- [Docker：GitHub Actions 构建集成](https://docs.docker.com/build/ci/github-actions/)
- [Docker：固定 Buildx/BuildKit 版本](https://docs.docker.com/build/ci/github-actions/configure-builder/)
- [Docker：固定 Docker Engine 的 setup action](https://github.com/docker/setup-docker-action)
- [GitHub：2026 installation token format rollout](https://github.blog/changelog/2026-04-24-notice-about-upcoming-new-format-for-github-app-installation-tokens/)
