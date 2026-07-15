# GitHub Actions 与 GHCR 迁移实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 Story 1.1 的唯一远程、CI 和 canonical registry 从 GitLab 完整迁移到公开 GitHub 仓库、GitHub Actions 与 GHCR，并用真实 digest image、完整 baseline 和 30 天 artifacts 解除 `G-TOOLCHAIN` 阻断。

**Architecture:** 保留现有 `scripts/run_toolchain_baseline.sh` 作为唯一完整验收入口，由 GitHub Actions 的 `toolchain-build` job 构建并推送 commit-bound GHCR image，再由串行 `toolchain-verify` job 只按 `tag@sha256` 执行 baseline。`build/toolchain/toolchain.lock` 继续作为外部执行环境的唯一 pin 来源，静态 Go validator 与 shell 负向测试共同拒绝 workflow、权限、版本、digest 和 GitLab 残留漂移。

**Tech Stack:** GitHub Actions、GHCR、Docker Engine 29.4.0、Buildx 0.35.0、BuildKit 0.31.1、Bash、Go 1.26.5、`gopkg.in/yaml.v3`、GitHub CLI。

## Global Constraints

- 所属阶段：Story 1.1／`G-TOOLCHAIN`；门禁保持 `BLOCKED`，直到真实 GitHub Actions run 全部通过。
- 设计依据：`docs/superpowers/specs/2026-07-15-github-actions-migration-design.md`。
- 仓库固定为公开 `wanggang8/agent-platform-eino`，remote 固定为 `origin=https://github.com/wanggang8/agent-platform-eino.git`。
- 只保留 GitHub Actions 与 GHCR；删除 `.gitlab-ci.yml`，不保留 GitLab 兼容、双 CI 或降级路径。
- workflow 只响应仓库内任意 `push` 与 `workflow_dispatch`；不增加可向 fork PR 发放 `packages: write` 的路径。
- canonical image 固定为 `ghcr.io/wanggang8/agent-platform-eino/toolchain:$GITHUB_SHA@sha256:$IMAGE_DIGEST`，其中 SHA 必须为 40 位小写十六进制，digest 必须为 64 位小写十六进制。
- runner 固定 `ubuntu-24.04`；Docker Engine `29.4.0`、Buildx `0.35.0`、BuildKit `moby/buildkit:v0.31.1` 均必须精确固定。
- action 只能使用本文列出的完整 40 字符 commit SHA；禁止浮动 tag、major tag、`main` 或 `master`。
- workflow 顶层 `permissions: {}`；build job 只允许 `contents: read`、`packages: write`，verify job 只允许 `contents: read`、`packages: read`。
- Workbench 和产品代码零改动：禁止修改 `internal/`、`cmd/`、`web/eino-workbench/src/`、schema、Product Facts、provider 或 Eino runtime。
- token、cookie、本地 FOBrain 配置、raw provider payload 不得写入代码、fixture、日志、artifact 或 Git 历史。
- 每个实现任务遵循测试先行；每次提交只包含该任务的脚本／测试／文档变更。

---

## 文件职责图

| 文件 | 操作 | 单一职责 |
| --- | --- | --- |
| `build/toolchain/toolchain.lock` | 修改 | 固定 runner、官方 action、Docker、Buildx、BuildKit 与既有 runtime 外部输入 |
| `scripts/toolchain_lock.sh` | 修改 | 将不可信 lock 解析为固定 21 字段安全接口 |
| `scripts/build_toolchain_image.sh` | 修改 | 从 GitHub identity 构造、校验并推送唯一 GHCR tag，输出 digest ref |
| `scripts/run_toolchain_baseline.sh` | 修改 | 在 CI 中绑定 40 字符 `GITHUB_SHA` 并执行前后 clean gate |
| `scripts/build_toolchain_image_test.sh` | 修改 | 覆盖 lock、builder、GitHub identity、digest、baseline clean 负向行为 |
| `.github/workflows/toolchain.yml` | 新建 | 定义 build→verify 的唯一 GitHub Actions 门禁 |
| `.gitlab-ci.yml` | 删除 | 移除旧 GitLab/DinD CI 契约 |
| `scripts/validate_ci_config/main.go` | 重写 | 静态验证 GitHub workflow、lock pins、DAG、权限、脚本引用和 artifacts |
| `scripts/validate_ci_config/main_test.go` | 重写 | 正向 fixture 与逐项负向漂移测试 |
| `scripts/validate_ci_config.sh` | 修改 | 默认验证 `.github/workflows/toolchain.yml` |
| `scripts/verify_toolchain.sh` | 修改 | 拒绝 GitLab 残留并强制调用 GitHub workflow validator |
| `scripts/verify_toolchain_test.sh` | 修改 | 验证 GitHub workflow fixture 与 GitLab 残留拒绝策略 |
| `docs/adr/2026-07-15-github-actions-ghcr-toolchain-gate.md` | 新建 | 记录 GitHub/GHCR 单平台架构决策与不可兼容边界 |
| `docs/07-implementation-plan.md` | 修改 | 将最终门禁、变量和证据语义改为 GitHub |
| `docs/08-acceptance-plan.md` | 修改 | 将 PASS 算法和解除条件改为 Actions/GHCR |
| `docs/tooling-and-reporting.md` | 修改 | 将 canonical CI 证据链和报告来源改为 GitHub |
| `docs/acceptance-records/story-1-1-g-toolchain-2026-07-15.md` | 修改 | 保留本地历史事实并记录 GitHub 最终裁决证据 |
| `docs/superpowers/specs/2026-07-15-github-actions-migration-design.md` | 修改 | 将设计状态更新为已实施或已验收 |

## 固定外部输入

```text
GITHUB_RUNNER=ubuntu-24.04
ACTIONS_CHECKOUT_SHA=9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0
ACTIONS_UPLOAD_ARTIFACT_SHA=043fb46d1a93c77aae656e7c1c64a875d1fc6a0a
DOCKER_SETUP_DOCKER_SHA=6d7cfa65f60a9dda7b46e5513fa982536f3c9877
DOCKER_SETUP_BUILDX_SHA=bb05f3f5519dd87d3ba754cc423b652a5edd6d2c
DOCKER_ENGINE_VERSION=29.4.0
DOCKER_BUILDX_VERSION=0.35.0
BUILDKIT_IMAGE=moby/buildkit:v0.31.1
BUILDKIT_DIGEST=sha256:6b59b7df63a8cb9902736f9ddf7fcff8261613d3e7449b8ea8b7537fc399c03a
```

---

### Task 1: 将 toolchain lock 迁移为 GitHub runner 与 action pins

**Files:**
- Modify: `build/toolchain/toolchain.lock`
- Modify: `scripts/toolchain_lock.sh`
- Modify: `scripts/build_toolchain_image_test.sh`

**Interfaces:**
- Consumes: 既有 12 个 runtime/browser/snapshot lock 字段。
- Produces: `scripts/toolchain_lock.sh build/toolchain/toolchain.lock` 的 21 个 tab 分隔字段，顺序为既有 12 字段后接 `GITHUB_RUNNER`、两个 GitHub action SHA、两个 Docker action SHA、Docker Engine、Buildx、BuildKit image、BuildKit digest。

- [ ] **Step 1: 先修改 lock 契约测试并验证失败**

  将 `scripts/build_toolchain_image_test.sh` 顶部的必需 key 列表改为：

  ```bash
  for key in PLATFORM PLAYWRIGHT_IMAGE PLAYWRIGHT_AMD64_DIGEST PLAYWRIGHT_VERSION \
    CHROMIUM_REVISION CHROMIUM_VERSION GO_LINUX_AMD64_SHA256 NODE_LINUX_X64_SHA256 \
    GITHUB_RUNNER ACTIONS_CHECKOUT_SHA ACTIONS_UPLOAD_ARTIFACT_SHA \
    DOCKER_SETUP_DOCKER_SHA DOCKER_SETUP_BUILDX_SHA DOCKER_ENGINE_VERSION \
    DOCKER_BUILDX_VERSION BUILDKIT_IMAGE BUILDKIT_DIGEST \
    UBUNTU_SNAPSHOT APT_BUILD_PACKAGES; do
    grep -Eq "^${key}=[^[:space:]]+$" "$lock" || {
      printf 'missing canonical lock key: %s\n' "$key" >&2
      exit 1
    }
  done
  test "$(bash "$lock_helper" "$lock" | awk -F '\t' '{ print NF }')" = 21
  ```

  删除对 `DOCKER_CLI_IMAGE`、`DOCKER_CLI_AMD64_DIGEST`、`DOCKER_DIND_IMAGE`、`DOCKER_DIND_AMD64_DIGEST` 的正向断言，并增加完整 action SHA、semver、BuildKit digest 的错误 lock mutation。

  Run: `bash scripts/build_toolchain_image_test.sh`

  Expected: FAIL，首个错误是缺少 `GITHUB_RUNNER` 或字段数不是 21。

- [ ] **Step 2: 替换 lock 中 GitLab/DinD pins**

  删除 lock 末尾四个 Docker CLI/DinD 字段，原样加入“固定外部输入”中的九行。不得修改已有 Playwright、Chromium、snapshot、Go 或 Node 值。

- [ ] **Step 3: 将 lock helper 扩展为 21 字段安全接口**

  在 `scripts/toolchain_lock.sh` 中增加变量和 case 分支：

  ```bash
  github_runner= actions_checkout_sha= actions_upload_artifact_sha=
  docker_setup_docker_sha= docker_setup_buildx_sha=
  docker_engine_version= docker_buildx_version= buildkit_image= buildkit_digest=
  ```

  必需字段检查与值域必须包含：

  ```bash
  [[ $github_runner == ubuntu-24.04 ]] || fail
  for action_sha in "$actions_checkout_sha" "$actions_upload_artifact_sha" \
    "$docker_setup_docker_sha" "$docker_setup_buildx_sha"; do
    [[ $action_sha =~ ^[0-9a-f]{40}$ ]] || fail
  done
  [[ $docker_engine_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
  [[ $docker_buildx_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
  [[ $buildkit_image =~ ^moby/buildkit:v[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
  [[ $buildkit_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail
  ```

  `printf` 保持固定顺序并输出 21 个字段；不得 `source` lock 或回显无效原始值。

- [ ] **Step 4: 运行 lock 与 shell 测试**

  Run: `bash -n scripts/toolchain_lock.sh scripts/build_toolchain_image_test.sh && bash scripts/build_toolchain_image_test.sh`

  Expected: PASS；恶意 `$(...)`、重复 key、缺 key、浮动版本、非 40 字符 action SHA 和错误 digest 均被拒绝且不泄露临时路径或 secret。

- [ ] **Step 5: 提交 lock 迁移**

  ```bash
  git add build/toolchain/toolchain.lock scripts/toolchain_lock.sh scripts/build_toolchain_image_test.sh
  git commit -m "build(toolchain): pin GitHub Actions runtime"
  ```

---

### Task 2: 将 builder 与 baseline identity 迁移到 GitHub

**Files:**
- Modify: `scripts/build_toolchain_image.sh`
- Modify: `scripts/run_toolchain_baseline.sh`
- Modify: `scripts/build_toolchain_image_test.sh`

**Interfaces:**
- Consumes: `CI=true`、`GITHUB_REPOSITORY=wanggang8/agent-platform-eino`、40 位小写十六进制 `GITHUB_SHA`。
- Produces: `TOOLCHAIN_IMAGE=ghcr.io/wanggang8/agent-platform-eino/toolchain:$GITHUB_SHA@sha256:$IMAGE_DIGEST`，同时以 mode `0600` 原子写入 `test-results/toolchain.env`。

- [ ] **Step 1: 先写 GitHub identity 与 baseline 失败测试**

  将所有 baseline fixture 的 `CI_COMMIT_SHA=0123456789abcdef` 改成：

  ```bash
  GITHUB_SHA=0123456789abcdef0123456789abcdef01234567
  ```

  将 push 正向场景改成：

  ```bash
  commit_sha=0123456789abcdef0123456789abcdef01234567
  image_ref="ghcr.io/wanggang8/agent-platform-eino/toolchain:$commit_sha"
  CI=true GITHUB_REPOSITORY=WangGang8/Agent-Platform-Eino GITHUB_SHA=$commit_sha \
    PATH="$bin:$PATH" bash "$repo/scripts/build_toolchain_image.sh" \
      --push "$image_ref" --env-file test-results/toolchain.env >"$tmp/push.out"
  ```

  添加负向断言，分别覆盖：缺 `GITHUB_REPOSITORY`、repo 不含单个 `/`、非 40 字符 SHA、带大写 SHA、非 `ghcr.io` target、tag 与 SHA 不一致、env file 非 canonical path、GitLab变量不能替代 GitHub变量。

  Run: `bash scripts/build_toolchain_image_test.sh`

  Expected: FAIL，现有实现报告 `push requires GitLab identity` 或无法接受 GHCR target。

- [ ] **Step 2: 更新 builder 的固定 lock 字段绑定**

  `read_lock` 使用 Task 1 的 21 字段顺序，后九个字段只用于验证 lock 完整性，不作为外部环境覆盖 build args：

  ```bash
  IFS=$'\t' read -r platform playwright_image playwright_digest playwright_version \
    chromium_revision chromium_version base_os font_policy ubuntu_snapshot apt_build_packages \
    go_sha256 node_sha256 github_runner actions_checkout_sha actions_upload_artifact_sha \
    docker_setup_docker_sha docker_setup_buildx_sha docker_engine_version \
    docker_buildx_version buildkit_image buildkit_digest <<<"$output" \
    || fail 'invalid toolchain lock'
  ```

- [ ] **Step 3: 用 GitHub identity 取代 GitLab push 分支**

  `--push` 分支采用以下完整校验逻辑：

  ```bash
  test "${CI:-}" = true || fail 'push requires CI'
  test -n "${GITHUB_REPOSITORY:-}" && test -n "${GITHUB_SHA:-}" || fail 'push requires GitHub identity'
  [[ $GITHUB_REPOSITORY =~ ^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$ ]] || fail 'invalid GitHub repository'
  [[ $GITHUB_SHA =~ ^[0-9a-f]{40}$ ]] || fail 'invalid GitHub commit SHA'
  repository=${GITHUB_REPOSITORY,,}
  expected_ref="ghcr.io/$repository/toolchain:$GITHUB_SHA"
  test "$2" = "$expected_ref" || fail 'push target must match GitHub identity'
  ```

  其余 buildx build、digest inspect、临时文件、`0600` 和脱敏错误行为保持不变。

- [ ] **Step 4: 将 baseline CI identity 改为严格 GitHub SHA**

  在 baseline 的首个 CI 分支中用以下断言替换非空 GitLab变量检查：

  ```bash
  [[ ${GITHUB_SHA:-} =~ ^[0-9a-f]{40}$ ]]
  ```

  前后 tracked、staged、untracked clean checks 不变。

- [ ] **Step 5: 运行 builder 与 baseline 行为测试**

  Run: `bash -n scripts/build_toolchain_image.sh scripts/run_toolchain_baseline.sh scripts/build_toolchain_image_test.sh && bash scripts/build_toolchain_image_test.sh`

  Expected: PASS；正向输出 lowercase GHCR digest ref，所有 GitLab identity 和错误 target 均 fail-fast。

- [ ] **Step 6: 提交 identity 迁移**

  ```bash
  git add scripts/build_toolchain_image.sh scripts/run_toolchain_baseline.sh scripts/build_toolchain_image_test.sh
  git commit -m "build(toolchain): bind images to GitHub identity"
  ```

---

### Task 3: 重写 GitHub Actions 静态 validator

**Files:**
- Modify: `scripts/validate_ci_config/main_test.go`
- Modify: `scripts/validate_ci_config/main.go`
- Modify: `scripts/validate_ci_config.sh`

**Interfaces:**
- Consumes: GitHub workflow YAML、Task 1 的九个 CI pin 字段、仓库根目录。
- Produces: `validateConfig(data []byte, lock lockValues) error`、`loadLock(path string) (lockValues, error)` 和成功文本 `GitHub Actions config validated`。

- [ ] **Step 1: 用 GitHub workflow fixture 重写正向测试**

  在 `main_test.go` 定义 `completeTestLock()`，字段与 Task 1 完全一致；正向 fixture 必须包含：

  ```yaml
  name: Toolchain Gate
  on:
    push:
    workflow_dispatch:
  permissions: {}
  jobs:
    toolchain-build:
      runs-on: ubuntu-24.04
      permissions:
        contents: read
        packages: write
      outputs:
        toolchain_image: ${{ steps.build.outputs.TOOLCHAIN_IMAGE }}
    toolchain-verify:
      needs: toolchain-build
      runs-on: ubuntu-24.04
      permissions:
        contents: read
        packages: read
  ```

  fixture 中还必须提供完整 checkout/setup-docker/setup-buildx/login/build、digest guard/pull/baseline/upload-artifact steps。测试断言 `validateConfig` 返回 `nil`，`validateScriptReferences` 只能解析 `scripts/build_toolchain_image.sh` 与 `scripts/run_toolchain_baseline.sh`。

  Run: `GOTOOLCHAIN=local go test ./scripts/validate_ci_config -count=1`

  Expected: FAIL，现有 GitLab parser/validator 不接受 GitHub jobs 与 workflow triggers。

- [ ] **Step 2: 写完逐项负向测试矩阵**

  每个测试只改变一个事实，并断言稳定错误摘要：

  ```text
  missing push/workflow_dispatch
  unexpected pull_request trigger
  non-empty top-level permissions
  missing or extra job
  runner not ubuntu-24.04
  build permissions missing packages:write
  verify permissions missing packages:read
  action uses v7/main instead of 40-char SHA
  action SHA differs from lock
  Docker/Buildx/BuildKit version differs from lock
  build output not sourced from steps.build.outputs.TOOLCHAIN_IMAGE
  verify lacks needs:toolchain-build
  login does not use GHCR/GITHUB_TOKEN stdin
  builder target is not ghcr.io/$repository/toolchain:$GITHUB_SHA
  digest guard accepts tag-only image
  pull/baseline command differs
  artifact paths differ or retention-days != 30
  artifact step lacks always()
  absolute/../symlink script escape
  duplicate/missing/malformed lock key
  GitLab variable or .gitlab-ci text reappears
  ```

  Run: `GOTOOLCHAIN=local go test ./scripts/validate_ci_config -count=1`

  Expected: FAIL until the new validator is implemented。

- [ ] **Step 3: 定义新的 lock 与 workflow validator 边界**

  用以下类型替换 GitLab `lockValues`：

  ```go
  // lockValues 仅承载 GitHub Actions 执行边界需要核对的固定输入。
  type lockValues struct {
      GitHubRunner              string
      ActionsCheckoutSHA        string
      ActionsUploadArtifactSHA  string
      DockerSetupDockerSHA      string
      DockerSetupBuildxSHA      string
      DockerEngineVersion       string
      DockerBuildxVersion       string
      BuildKitImage             string
      BuildKitDigest            string
  }
  ```

  `loadLock` 继续逐行读取 `KEY=VALUE`，拒绝重复 key、空白和值域外字符；四个 action SHA 用 `^[0-9a-f]{40}$`，digest 用 `^sha256:[0-9a-f]{64}$`，版本用 exact semver，runner 必须为 `ubuntu-24.04`。

- [ ] **Step 4: 实现 exact workflow 验证**

  保留 `yaml.v3`，但删除 GitLab stage/template/extends/DinD merge 逻辑。新增下列小函数并各自只负责一个边界：

  ```go
  func validateTriggers(config map[string]any) error
  func validateTopLevelPermissions(config map[string]any) error
  func validateBuildJob(job map[string]any, lock lockValues) error
  func validateVerifyJob(job map[string]any, lock lockValues) error
  func validateExactPermissions(jobName string, actual any, expected map[string]string) error
  func validateActionStep(step map[string]any, expectedUses string, expectedWith map[string]string) error
  func validateRunStep(step map[string]any, expectedID, expectedRun string) error
  func validateArtifactStep(step map[string]any, lock lockValues) error
  func canonicalScriptReferences(config map[string]any) ([]string, error)
  ```

  所有 `run` 值先 `strings.TrimSpace` 再与 canonical multiline 常量精确比较；steps 数量、顺序、`id`、`uses`、`with`、`env` 和 permissions 都不得容忍额外执行入口。脚本引用仍通过 `EvalSymlinks`、`filepath.Rel` 和 `scripts/` 前缀防止逃逸。

- [ ] **Step 5: 将 CLI 默认输入和成功文本迁移到 GitHub**

  ```go
  configPath := flag.String("config", ".github/workflows/toolchain.yml", "GitHub Actions workflow path")
  data, err := os.ReadFile(*configPath)
  if err != nil {
      exitInvalid(errors.New("cannot read CI config"))
  }
  lock, err := loadLock(*lockPath)
  if err != nil {
      exitInvalid(err)
  }
  if err := validateConfig(data, lock); err != nil {
      exitInvalid(err)
  }
  config, err := parseConfig(data)
  if err != nil {
      exitInvalid(err)
  }
  if err := validateScriptReferences(*root, config); err != nil {
      exitInvalid(err)
  }
  fmt.Println("GitHub Actions config validated")
  ```

  `scripts/validate_ci_config.sh` 改成：

  ```bash
  exec env GOTOOLCHAIN=local go run ./scripts/validate_ci_config \
    -config .github/workflows/toolchain.yml -lock build/toolchain/toolchain.lock -root "$root"
  ```

- [ ] **Step 6: 运行 validator 单元测试和静态检查**

  Run: `gofmt -w scripts/validate_ci_config/main.go scripts/validate_ci_config/main_test.go && GOTOOLCHAIN=local go test ./scripts/validate_ci_config -count=1 && GOTOOLCHAIN=local go vet ./scripts/validate_ci_config`

  Expected: PASS；所有正向/负向测试通过，无 GitLab validator 类型或错误文案残留。

- [ ] **Step 7: 提交 validator 重写**

  ```bash
  git add scripts/validate_ci_config/main.go scripts/validate_ci_config/main_test.go scripts/validate_ci_config.sh
  git commit -m "test(ci): validate GitHub Actions toolchain gate"
  ```

---

### Task 4: 创建唯一 GitHub workflow 并删除 GitLab CI

**Files:**
- Create: `.github/workflows/toolchain.yml`
- Delete: `.gitlab-ci.yml`
- Modify: `scripts/verify_toolchain.sh`
- Modify: `scripts/verify_toolchain_test.sh`
- Modify: `scripts/build_toolchain_image_test.sh`

**Interfaces:**
- Consumes: Task 1 pins、Task 2 builder/baseline、Task 3 validator。
- Produces: `toolchain-build.outputs.toolchain_image`，供 `toolchain-verify` 作为唯一 `TOOLCHAIN_IMAGE`。

- [ ] **Step 1: 先把 shell fixtures 切到 GitHub workflow 并验证失败**

  两个 shell test 的临时仓库先创建 `.github/workflows/`，复制 `.github/workflows/toolchain.yml`，不再复制 `.gitlab-ci.yml`。在 `verify_toolchain_test.sh` 添加：

  ```bash
  printf '%s\n' 'stages: [toolchain, verify]' >"$tmp/repo/.gitlab-ci.yml"
  if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"; then
    echo 'expected GitLab CI residue to fail' >&2
    exit 1
  fi
  rm "$tmp/repo/.gitlab-ci.yml"
  ```

  Run: `bash scripts/verify_toolchain_test.sh`

  Expected: FAIL，因为 workflow 尚不存在且 verifier 仍读取 `.gitlab-ci.yml`。

- [ ] **Step 2: 创建 canonical GitHub Actions workflow**

  `.github/workflows/toolchain.yml` 必须实现以下执行向量：

  ```yaml
  name: Toolchain Gate

  on:
    push:
    workflow_dispatch:

  permissions: {}

  jobs:
    toolchain-build:
      runs-on: ubuntu-24.04
      permissions:
        contents: read
        packages: write
      outputs:
        toolchain_image: ${{ steps.build.outputs.TOOLCHAIN_IMAGE }}
      steps:
        - name: Checkout
          uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0
          with:
            fetch-depth: 0
        - name: Setup Docker
          uses: docker/setup-docker-action@6d7cfa65f60a9dda7b46e5513fa982536f3c9877
          with:
            version: v29.4.0
        - name: Setup Buildx
          uses: docker/setup-buildx-action@bb05f3f5519dd87d3ba754cc423b652a5edd6d2c
          with:
            version: v0.35.0
            driver-opts: image=moby/buildkit:v0.31.1@sha256:6b59b7df63a8cb9902736f9ddf7fcff8261613d3e7449b8ea8b7537fc399c03a
        - name: Login GHCR
          env:
            GHCR_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          run: printf '%s' "$GHCR_TOKEN" | docker login ghcr.io --username "$GITHUB_ACTOR" --password-stdin
        - name: Build and push canonical image
          id: build
          run: |
            repository=${GITHUB_REPOSITORY,,}
            image_ref="ghcr.io/$repository/toolchain:$GITHUB_SHA"
            bash scripts/build_toolchain_image.sh --push "$image_ref" --env-file test-results/toolchain.env
            cat test-results/toolchain.env >> "$GITHUB_OUTPUT"

    toolchain-verify:
      needs: toolchain-build
      runs-on: ubuntu-24.04
      permissions:
        contents: read
        packages: read
      env:
        TOOLCHAIN_IMAGE: ${{ needs.toolchain-build.outputs.toolchain_image }}
      steps:
        - name: Checkout
          uses: actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0
          with:
            fetch-depth: 0
        - name: Setup Docker
          uses: docker/setup-docker-action@6d7cfa65f60a9dda7b46e5513fa982536f3c9877
          with:
            version: v29.4.0
        - name: Login GHCR
          env:
            GHCR_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          run: printf '%s' "$GHCR_TOKEN" | docker login ghcr.io --username "$GITHUB_ACTOR" --password-stdin
        - name: Verify digest image
          run: |
            [[ "$TOOLCHAIN_IMAGE" =~ ^ghcr\.io/[a-z0-9._/-]+/toolchain:[0-9a-f]{40}@sha256:[0-9a-f]{64}$ ]]
            docker pull "$TOOLCHAIN_IMAGE"
            docker run --rm --platform linux/amd64 -e CI=true -e GITHUB_SHA="$GITHUB_SHA" -v "$GITHUB_WORKSPACE:/workspace" -w /workspace "$TOOLCHAIN_IMAGE" bash scripts/run_toolchain_baseline.sh
        - name: Upload toolchain evidence
          if: ${{ always() }}
          uses: actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a
          with:
            name: toolchain-evidence-${{ github.sha }}
            path: |
              test-results/toolchain-baseline.log
              test-results/eino-workbench-playwright-report/
            if-no-files-found: error
            retention-days: 30
  ```

- [ ] **Step 3: 让 verifier 强制 GitHub-only 契约**

  `verify_ci_config` 改成读取 `.github/workflows/toolchain.yml`，并只扫描生产执行入口、先拒绝：

  ```bash
  test ! -e "$root/.gitlab-ci.yml" || fail 'GitLab CI residue is forbidden'
  for target in "$ci" "$root/scripts/build_toolchain_image.sh" \
    "$root/scripts/run_toolchain_baseline.sh" "$root/scripts/validate_ci_config.sh" \
    "$root/scripts/validate_ci_config/main.go"; do
    if grep -Eq 'CI_COMMIT_SHA|CI_REGISTRY|CI_PROJECT_DIR|GitLab|\.gitlab-ci' "$target"; then
      fail 'GitLab CI residue is forbidden'
    fi
  done
  ```

  继续拒绝 workflow 复制 Go、Node、npm 三个权威版本，再调用 `scripts/validate_ci_config.sh`；错误改为 `invalid GitHub Actions workflow`。

- [ ] **Step 4: 删除 GitLab CI 并运行所有 CI 契约测试**

  删除 `.gitlab-ci.yml`，然后运行：

  ```bash
  bash scripts/verify_toolchain_test.sh
  bash scripts/build_toolchain_image_test.sh
  GOTOOLCHAIN=local go test ./scripts/validate_ci_config -count=1
  bash scripts/validate_ci_config.sh
  ```

  Expected: 全部 PASS；`find . -name .gitlab-ci.yml -print` 无输出，workflow validator 输出 `GitHub Actions config validated`。

- [ ] **Step 5: 提交 CI cutover**

  ```bash
  git add .github/workflows/toolchain.yml .gitlab-ci.yml \
    scripts/verify_toolchain.sh scripts/verify_toolchain_test.sh scripts/build_toolchain_image_test.sh
  git commit -m "ci(github): replace GitLab toolchain pipeline"
  ```

---

### Task 5: 同步 ADR、实施计划、验收规则与历史记录

**Files:**
- Create: `docs/adr/2026-07-15-github-actions-ghcr-toolchain-gate.md`
- Modify: `docs/07-implementation-plan.md`
- Modify: `docs/08-acceptance-plan.md`
- Modify: `docs/tooling-and-reporting.md`
- Modify: `docs/acceptance-records/story-1-1-g-toolchain-2026-07-15.md`
- Modify: `docs/superpowers/specs/2026-07-15-github-actions-migration-design.md`

**Interfaces:**
- Consumes: 已通过静态测试的 GitHub-only 契约。
- Produces: 与代码一致的 `G-TOOLCHAIN` PASS 算法；在真实 run 前仍明确为 `BLOCKED_PENDING_GITHUB_RUN`。

- [ ] **Step 1: 新增不可兼容迁移 ADR**

  ADR 必须包含 `Status: Accepted`、背景、决策、后果、拒绝方案和回滚条件，并明确：

  ```text
  GitHub 是唯一远程平台；GitHub Actions 是唯一 CI；GHCR 是唯一 canonical registry。
  .gitlab-ci.yml、GitLab变量、DinD service 和 GitLab dotenv artifact 语义全部删除。
  Product code、Eino runtime、Product Facts 和 Workbench UI 不受该 ADR 影响。
  只有新的 ADR 才能改变平台；不得在故障时静默恢复双 CI 或 tag-only 验证。
  ```

- [ ] **Step 2: 统一替换门禁术语与命令**

  在三份规范文档中统一使用：

  ```text
  clean GitHub Actions run
  GITHUB_SHA
  GitHub Actions run URL
  ghcr.io/wanggang8/agent-platform-eino/toolchain:$GITHUB_SHA@sha256:$IMAGE_DIGEST
  toolchain-evidence-$GITHUB_SHA (30 days)
  ```

  保留本地 image ID 作为 preflight 历史，但明确它不能代替 GHCR digest；删除 GitLab runner、pipeline、registry、DinD 和 `CI_PROJECT_DIR` 的现行要求。

- [ ] **Step 3: 更新 Story 1.1 验收记录但暂不伪造 PASS**

  将原 GitLab 阻断改写为历史事实；新增当前直接阻断：

  ```text
  BLOCKED_PENDING_GITHUB_RUN：GitHub-only 静态契约已实现，但尚未创建／推送公开仓库，未产生真实 Actions run URL、GHCR digest 与 30 天 artifacts。
  ```

  设计文档状态改为“设计已批准；实施完成后等待真实 GitHub run 验收”。

- [ ] **Step 4: 做文档一致性扫描**

  Run:

  ```bash
  rg -n 'GitLab|CI_COMMIT_SHA|CI_REGISTRY|CI_PROJECT_DIR|\.gitlab-ci' \
    docs/07-implementation-plan.md docs/08-acceptance-plan.md docs/tooling-and-reporting.md \
    docs/acceptance-records/story-1-1-g-toolchain-2026-07-15.md \
    docs/superpowers/specs/2026-07-15-github-actions-migration-design.md
  ```

  Expected: 只允许 ADR/验收记录中标记为“历史事实”或“已删除方案”的句子；不得出现现行 GitLab 执行指令。

- [ ] **Step 5: 提交文档迁移**

  ```bash
  git add docs/adr/2026-07-15-github-actions-ghcr-toolchain-gate.md \
    docs/07-implementation-plan.md docs/08-acceptance-plan.md docs/tooling-and-reporting.md \
    docs/acceptance-records/story-1-1-g-toolchain-2026-07-15.md \
    docs/superpowers/specs/2026-07-15-github-actions-migration-design.md
  git commit -m "docs(ci): adopt GitHub Actions gate"
  ```

---

### Task 6: 在本地 canonical image 中执行迁移后的完整回归

**Files:**
- Modify only if a test exposes a defect: files already listed in Tasks 1-5
- Evidence output, not committed: `test-results/toolchain-baseline.log`
- Evidence output, not committed: `test-results/eino-workbench-playwright-report/`

**Interfaces:**
- Consumes: Tasks 1-5 的 clean commit 与已批准本地 image ID `sha256:41f8317a3cc392bca3eedcf390e70b1c6ae4be0b4548cc4d7d3c4f200a29ea93`。
- Produces: GitHub-only 静态/行为测试 PASS 与本地完整 baseline preflight；不产生最终 gate PASS。

- [ ] **Step 1: 运行任务级静态和行为检查**

  ```bash
  bash -n scripts/toolchain_lock.sh scripts/build_toolchain_image.sh \
    scripts/run_toolchain_baseline.sh scripts/verify_toolchain.sh \
    scripts/verify_toolchain_test.sh scripts/build_toolchain_image_test.sh \
    scripts/validate_ci_config.sh
  bash scripts/verify_toolchain_test.sh
  bash scripts/build_toolchain_image_test.sh
  GOTOOLCHAIN=local go test ./scripts/validate_ci_config -count=1
  GOTOOLCHAIN=local go vet ./scripts/validate_ci_config
  bash scripts/validate_ci_config.sh
  git diff --check
  ```

  Expected: 全部 exit 0。

- [ ] **Step 2: 证明产品目录零 diff 与 GitLab 执行契约清零**

  ```bash
  test -z "$(git diff --name-only 6ebf465..HEAD -- internal cmd web/eino-workbench/src)"
  test ! -e .gitlab-ci.yml
  ! rg -n 'CI_COMMIT_SHA|CI_REGISTRY|CI_PROJECT_DIR' \
    .github build/toolchain scripts \
    --glob '!validate_ci_config/main_test.go' \
    --glob '!verify_toolchain_test.sh'
  ```

  Expected: 全部 exit 0；产品实现目录没有变更。

- [ ] **Step 3: 用已批准 canonical image 重跑唯一 baseline**

  ```bash
  image_id=sha256:41f8317a3cc392bca3eedcf390e70b1c6ae4be0b4548cc4d7d3c4f200a29ea93
  git_common_dir=$(git rev-parse --path-format=absolute --git-common-dir)
  docker run --rm --platform linux/amd64 \
    -v "$PWD:/workspace" -v "$git_common_dir:$git_common_dir:ro" \
    -w /workspace "$image_id" bash scripts/run_toolchain_baseline.sh
  ```

  Expected: exit 0；desktop 17/17、contract smoke 和末尾 clean gate PASS。该结果仍标记 `evidence_mode=preflight`，不能解除远程门禁。

- [ ] **Step 4: 若发现缺陷，按最小测试循环修正并提交**

  每个失败先确认是测试、validator、workflow 还是环境根因；只修改 Tasks 1-5 已列文件，重跑对应失败命令及 Step 1-3。修复提交使用：

  ```bash
  git add build/toolchain/toolchain.lock .github/workflows/toolchain.yml \
    scripts/toolchain_lock.sh scripts/build_toolchain_image.sh \
    scripts/run_toolchain_baseline.sh scripts/verify_toolchain.sh \
    scripts/verify_toolchain_test.sh scripts/build_toolchain_image_test.sh \
    scripts/validate_ci_config.sh scripts/validate_ci_config/main.go \
    scripts/validate_ci_config/main_test.go
  git commit -m "fix(ci): close GitHub toolchain gate regression"
  ```

  Expected: 工作区 clean，所有本地证据 PASS。

---

### Task 7: 创建公开 GitHub 仓库、推送当前分支并验证真实 Actions run

**Files:**
- External state: GitHub repository `wanggang8/agent-platform-eino`
- External state: GHCR package `ghcr.io/wanggang8/agent-platform-eino/toolchain`
- External evidence: GitHub Actions run and artifacts

**Interfaces:**
- Consumes: Task 6 的 clean branch `codex/story-1-1-toolchain` 与已登录 `gh` 账号 `wanggang8`。
- Produces: `origin` remote、真实 run URL、commit-bound GHCR digest、`toolchain-evidence-$GITHUB_SHA` artifact。

- [ ] **Step 1: 在创建远程前做 secret 和仓库身份检查**

  ```bash
  git status --short
  git remote -v
  gh auth status --hostname github.com
  git ls-files configs/eino-workbench.local.yaml
  rg -n --hidden --glob '!.git/**' --glob '!test-results/**' \
    'FOBRAIN_USER_API_TOKEN=|(?i)(token|secret|password)[[:space:]]*[:=][[:space:]]*[^$[:space:]]' .
  ```

  Expected: 工作区 clean；无 remote；账号为 `wanggang8`；本地配置未 tracked；已知 token 与凭据赋值零匹配。若发现凭据，停止创建仓库并先清理 Git 历史/文件。

- [ ] **Step 2: 创建公开空仓库并添加 origin**

  ```bash
  gh repo create wanggang8/agent-platform-eino \
    --public \
    --description "Eino-first Agent Workbench"
  git remote add origin https://github.com/wanggang8/agent-platform-eino.git
  gh repo view wanggang8/agent-platform-eino --json nameWithOwner,visibility,url
  ```

  Expected: `nameWithOwner=wanggang8/agent-platform-eino`、`visibility=PUBLIC`，remote 精确匹配设计。若仓库已存在，先用 `gh repo view` 验证 owner、visibility 与空/兼容状态，不覆盖未知内容。

- [ ] **Step 3: 推送当前开发分支，不强推、不改写本地分支**

  ```bash
  git push -u origin codex/story-1-1-toolchain
  ```

  Expected: push 成功并触发 `Toolchain Gate`；当前本地分支名保持不变。

- [ ] **Step 4: 找到并等待真实 workflow run**

  ```bash
  run_id=$(gh run list --repo wanggang8/agent-platform-eino \
    --branch codex/story-1-1-toolchain --workflow toolchain.yml \
    --limit 1 --json databaseId --jq '.[0].databaseId')
  test -n "$run_id"
  gh run watch "$run_id" --repo wanggang8/agent-platform-eino --exit-status
  gh run view "$run_id" --repo wanggang8/agent-platform-eino \
    --json url,headSha,status,conclusion,jobs
  ```

  Expected: build 和 verify job 均 `success`；`headSha` 等于本地 `git rev-parse HEAD`。

- [ ] **Step 5: 若 workflow 失败，收集日志并回到最小修复循环**

  ```bash
  gh run view "$run_id" --repo wanggang8/agent-platform-eino --log-failed \
    >test-results/github-actions-failed.log
  ```

  先定位首个失败步骤；补失败测试后修改 Tasks 1-5 范围内的文件，运行 Task 6，单独提交并正常 push。每次 push 都重新取得最新 `run_id`；禁止手工伪造 output、放宽 digest guard、使用浮动 action 或切回 GitLab。

- [ ] **Step 6: 验证 GHCR digest 与 artifacts**

  ```bash
  head_sha=$(git rev-parse HEAD)
  gh run download "$run_id" --repo wanggang8/agent-platform-eino \
    --name "toolchain-evidence-$head_sha" \
    --dir "test-results/github-run-$run_id"
  test -s "test-results/github-run-$run_id/toolchain-baseline.log"
  test -d "test-results/github-run-$run_id/eino-workbench-playwright-report"
  gh api "/repos/wanggang8/agent-platform-eino/actions/runs/$run_id/artifacts" \
    --jq '.artifacts[] | select(.name == "toolchain-evidence-'"$head_sha"'") | {name,expired,expires_at}'
  docker buildx imagetools inspect \
    "ghcr.io/wanggang8/agent-platform-eino/toolchain:$head_sha"
  ```

  Expected: artifact 名称精确、`expired=false`、到期日约为创建后 30 天；镜像输出匹配 `sha256:[0-9a-f]{64}`，并与 build job 输出和 verify job pull 的 digest 相同。

---

### Task 8: 形成最终验收记录并关闭 `G-TOOLCHAIN`

**Files:**
- Modify: `docs/acceptance-records/story-1-1-g-toolchain-2026-07-15.md`
- Modify: `docs/07-implementation-plan.md`
- Modify: `docs/08-acceptance-plan.md`
- Modify: `docs/superpowers/specs/2026-07-15-github-actions-migration-design.md`

**Interfaces:**
- Consumes: Task 7 的 repository URL、head SHA、run URL、GHCR digest、artifact 名称/expiry 和完整 baseline PASS。
- Produces: `G-TOOLCHAIN=PASS` 的可审计裁决，并允许实施计划进入 M-1。

- [ ] **Step 1: 写入第一条真实远程证据，逐项裁决 PASS 算法**

  验收记录必须写入实际值而不是变量名：

  ```text
  repository URL
  commit SHA
  GitHub Actions run URL
  toolchain-build conclusion
  toolchain-verify conclusion
  GHCR tag@sha256 digest
  artifact name
  artifact expires_at
  desktop 17/17
  contract smoke PASS
  final tracked/staged/untracked clean checks PASS
  ```

  原本地 image、UX approval 与阻断历史保留；最终结论只有在所有行均有真实证据时改为 `G-TOOLCHAIN=PASS`。

- [ ] **Step 2: 更新里程碑状态和设计状态**

  `docs/07-implementation-plan.md` 将 Story 1.1／Sprint 1.1 标为完成，并把下一步设为 M-1；`docs/08-acceptance-plan.md` 记录最终 run 的裁决；设计状态改为“已实施并通过真实 GitHub Actions 验收”。

- [ ] **Step 3: 提交验收记录并推送**

  ```bash
  git add docs/acceptance-records/story-1-1-g-toolchain-2026-07-15.md \
    docs/07-implementation-plan.md docs/08-acceptance-plan.md \
    docs/superpowers/specs/2026-07-15-github-actions-migration-design.md
  git commit -m "docs(toolchain): record GitHub gate evidence"
  git push
  ```

- [ ] **Step 4: 等待验收文档提交触发的最终 clean run**

  ```bash
  final_sha=$(git rev-parse HEAD)
  final_run_id=$(gh run list --repo wanggang8/agent-platform-eino \
    --branch codex/story-1-1-toolchain --workflow toolchain.yml \
    --limit 1 --json databaseId,headSha \
    --jq '.[] | select(.headSha == "'"$final_sha"'") | .databaseId')
  test -n "$final_run_id"
  gh run watch "$final_run_id" --repo wanggang8/agent-platform-eino --exit-status
  gh run view "$final_run_id" --repo wanggang8/agent-platform-eino \
    --json url,headSha,status,conclusion,jobs
  ```

  Expected: 最终文档提交所在 clean SHA 的 build/verify 均成功。此 run 的 URL、digest 与 artifact 信息记录在最终交付说明中；不再制造仅为记录自身 run ID 的无限文档提交链。

- [ ] **Step 5: 完成最终一致性与安全检查**

  ```bash
  git status --short
  git diff --check
  test -z "$(git diff --name-only 6ebf465..HEAD -- internal cmd web/eino-workbench/src)"
  test ! -e .gitlab-ci.yml
  git remote get-url origin
  gh repo view wanggang8/agent-platform-eino --json visibility,url
  ```

  Expected: 工作区 clean；产品实现零 diff；GitLab CI 文件不存在；origin 为 GitHub；仓库为 PUBLIC；最终 Actions run 成功。

---

## 最终交付门禁

只有以下条件同时满足才能汇报“完成”：

- 本地 shell、Go validator、工具链测试和完整 canonical baseline 均 exit 0。
- `.gitlab-ci.yml`、GitLab变量和现行 GitLab文档语义已清零。
- GitHub 仓库真实存在、公开、remote 正确，当前分支已正常推送。
- 最终 clean commit 的 `toolchain-build` 与 `toolchain-verify` 均成功。
- GHCR image 同时具备 commit tag 与不可变 `sha256` digest，verify job 按 digest 使用它。
- `toolchain-baseline.log` 与 Playwright report artifact 可下载且 retention 为 30 天。
- 验收记录包含第一条真实远程证据，最终交付说明包含最终文档提交的 run URL/digest/artifact。
- `internal/`、`cmd/`、`web/eino-workbench/src/` 相对计划基线零 diff。
