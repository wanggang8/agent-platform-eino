#!/usr/bin/env bash
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
lock="$root/build/toolchain/toolchain.lock"
dockerfile="$root/build/toolchain/Dockerfile"
lock_helper="$root/scripts/toolchain_lock.sh"

test -f "$lock" || { printf 'missing toolchain.lock\n' >&2; exit 1; }
test -f "$dockerfile" || { printf 'missing toolchain Dockerfile\n' >&2; exit 1; }
test -f "$lock_helper" || { printf 'missing shared toolchain lock helper\n' >&2; exit 1; }
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
grep -Fxq 'NODE_LINUX_X64_SHA256=783130984963db7ba9cbd01089eaf2c2efb055c7c1693c943174b967b3050cb8' "$lock"
test "$(bash "$lock_helper" "$lock" | awk -F '\t' '{ print NF }')" = 21
grep -Fq 'FROM ${PLAYWRIGHT_IMAGE}@${PLAYWRIGHT_DIGEST}' "$dockerfile"
! grep -Eq 'go1\.26\.5|node-v24\.18\.0' "$dockerfile"
! grep -Eq 'tar\.xz|xJf' "$dockerfile"
for expected in \
  'ARG TARGETARCH' \
  'ARG UBUNTU_SNAPSHOT' \
  'ARG APT_BUILD_PACKAGES' \
  'RUN test "$TARGETARCH" = "amd64"' \
  'node-v${NODE_VERSION}-linux-x64.tar.gz' \
  'tar -C /usr/local --strip-components=1 -xzf /tmp/node.tar.gz' \
  'Snapshot: ${UBUNTU_SNAPSHOT}' \
  'signed_by_count=$(grep -c' \
  "'^Signed-By:'" \
  'test "$snapshot_count" = "$signed_by_count"' \
  "! -name 'ubuntu.sources' -delete" \
  'Dir::Etc::sourcelist="$sources"' \
  'Dir::Etc::sourceparts="-"' \
  'apt-get "${apt_snapshot_options[@]}" update' \
  'apt-get "${apt_snapshot_options[@]}" install -y --no-install-recommends "$APT_BUILD_PACKAGES"' \
  'rm -rf /var/lib/apt/lists/*' \
  'test "$(go env GOVERSION)" = "go${GO_VERSION}"' \
  'test "$(node --version)" = "v${NODE_VERSION}"' \
  'test "$(npm --version)" = "${NPM_VERSION}"' \
  'command -v cc' \
  'go mod init toolchain-race-smoke' \
  'CGO_ENABLED=1 go test -race ./...' \
  'rm -rf "$race_dir"' \
  'find "/ms-playwright/chromium-${CHROMIUM_REVISION}"' \
  'chromium_actual=$("$chromium_binary" --version)' \
  '[[ "$chromium_actual" == *"${CHROMIUM_VERSION}"* ]]'; do
  grep -Fq "$expected" "$dockerfile"
done
for consumer in "$root/scripts/build_toolchain_image.sh" "$root/scripts/verify_toolchain.sh"; do
  grep -Fq 'scripts/toolchain_lock.sh' "$consumer"
  ! grep -Fq 'case "$key" in' "$consumer"
done
bash -n "$root/scripts/toolchain_lock.sh" "$root/scripts/build_toolchain_image.sh"
node - "$root/package.json" <<'NODE'
const pkg = require(process.argv[2]);
if (pkg.scripts?.['build:toolchain-image'] !== 'bash scripts/build_toolchain_image.sh --load') {
  console.error('missing canonical image build script');
  process.exit(1);
}
if (pkg.scripts?.['test:toolchain-image'] !== 'bash scripts/build_toolchain_image_test.sh') {
  console.error('missing canonical image test script');
  process.exit(1);
}
if (pkg.scripts?.['validate:ci'] !== 'bash scripts/validate_ci_config.sh') {
  console.error('missing static CI validation script');
  process.exit(1);
}
if (pkg.scripts?.['test:ci-config'] !== 'GOTOOLCHAIN=local go test ./scripts/validate_ci_config -count=1') {
  console.error('missing static CI validator test script');
  process.exit(1);
}
NODE

# baseline 顺序是 canonical image 的验收契约，安装必须先于 desktop browser。
baseline="$root/scripts/run_toolchain_baseline.sh"
test -f "$baseline"
npm_ci_line=$(grep -n '^npm ci$' "$baseline" | cut -d: -f1)
browser_line=$(grep -n 'browser-test.*--project=desktop' "$baseline" | cut -d: -f1)
test "$npm_ci_line" -lt "$browser_line"
grep -Fxq 'npm run eino-workbench:browser-test -- --project=desktop --workers=1' "$baseline"
! grep -Eq '^(go test|go vet|go build|go list|go mod)' "$baseline"
! grep -Eq -- '--project=mobile' "$baseline"
test "$(grep -Fxc '  git diff --quiet' "$baseline")" = 2
test "$(grep -Fxc '  git diff --cached --quiet' "$baseline")" = 2
test "$(grep -Fxc '  test -z "$(git ls-files --others --exclude-standard)"' "$baseline")" = 2
smoke_line=$(grep -n '^bash scripts/eino_workbench_server_smoke.sh --scenario contract$' "$baseline" | cut -d: -f1)
final_clean_line=$(grep -n '^if \[\[ \${CI:-false} == true \]\]; then$' "$baseline" | tail -1 | cut -d: -f1)
test "$final_clean_line" -gt "$smoke_line"
test "$(tail -n 5 "$baseline")" = 'if [[ ${CI:-false} == true ]]; then
  git diff --quiet
  git diff --cached --quiet
  test -z "$(git ls-files --others --exclude-standard)"
fi'

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# 用最小 shim 证明测试产生的 tracked/staged/untracked 污染都会被末尾 clean check 捕获。
baseline_repo="$tmp/baseline-repo"
baseline_bin="$tmp/baseline-bin"
baseline_mutation="$tmp/baseline-mutation"
mkdir -p "$baseline_repo/scripts" "$baseline_bin"
cp "$baseline" "$baseline_repo/scripts/"
for file in package-lock.json go.mod go.sum; do
  printf '%s\n' fixture >"$baseline_repo/$file"
done
for script in verify_toolchain.sh validate_ci_config.sh eino_workbench_server_smoke.sh; do
  cat >"$baseline_repo/scripts/$script" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
exit 0
SH
done
cat >"$baseline_bin/npm" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
if [[ ${1:-} == ci && -n ${MUTATION_KIND:-} ]]; then
  printf '%s\n' "$MUTATION_KIND" >"$BASELINE_MUTATION"
fi
exit 0
SH
cat >"$baseline_bin/go" <<'SH'
#!/usr/bin/env bash
exit 0
SH
cat >"$baseline_bin/sha256sum" <<'SH'
#!/usr/bin/env bash
for file in "$@"; do
  printf 'fixture  %s\n' "$file"
done
SH
cat >"$baseline_bin/git" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
kind=
if [[ -f $BASELINE_MUTATION ]]; then
  kind=$(<"$BASELINE_MUTATION")
fi
case "$*" in
  'diff --quiet') test "$kind" != tracked ;;
  'diff --cached --quiet') test "$kind" != staged ;;
  'ls-files --others --exclude-standard')
    if [[ $kind == untracked ]]; then printf '%s\n' generated.file; fi
    ;;
  'diff --check') exit 0 ;;
  *) exit 64 ;;
esac
SH
chmod +x "$baseline_bin/npm" "$baseline_bin/go" "$baseline_bin/sha256sum" "$baseline_bin/git"
export BASELINE_MUTATION="$baseline_mutation"
for kind in tracked staged untracked; do
  rm -f "$baseline_mutation"
  if CI=true GITHUB_SHA=0123456789abcdef0123456789abcdef01234567 \
    MUTATION_KIND=$kind PATH="$baseline_bin:$PATH" \
    bash "$baseline_repo/scripts/run_toolchain_baseline.sh" >"$tmp/baseline-$kind.out" 2>&1; then
    printf 'expected final clean check to reject %s mutation\n' "$kind" >&2
    exit 1
  fi
done
rm -f "$baseline_mutation"
CI=true GITHUB_SHA=0123456789abcdef0123456789abcdef01234567 PATH="$baseline_bin:$PATH" \
  bash "$baseline_repo/scripts/run_toolchain_baseline.sh" >"$tmp/baseline-clean.out" 2>&1

# baseline 自身必须固定拒绝无效 identity，且错误不能回显不可信 SHA。
assert_baseline_identity_failure() {
  local case_name=$1
  shift
  if env -u GITHUB_SHA -u CI_COMMIT_SHA CI=true "$@" PATH="$baseline_bin:$PATH" \
    bash "$baseline_repo/scripts/run_toolchain_baseline.sh" \
      >"$tmp/baseline-identity.out" 2>&1; then
    printf 'expected baseline identity failure: %s\n' "$case_name" >&2
    exit 1
  fi
  local attempt=0
  while test ! -s "$tmp/baseline-identity.out" && test "$attempt" -lt 100; do
    attempt=$((attempt + 1))
    sleep 0.01
  done
  test "$(<"$tmp/baseline-identity.out")" = 'invalid GitHub commit SHA'
  ! grep -Fq "$tmp" "$tmp/baseline-identity.out"
}

assert_baseline_identity_failure missing
assert_baseline_identity_failure short GITHUB_SHA=0123456789abcdef
assert_baseline_identity_failure uppercase \
  GITHUB_SHA=ABCDEF0123456789ABCDEF0123456789ABCDEF01
assert_baseline_identity_failure gitlab-only \
  CI_COMMIT_SHA=0123456789abcdef0123456789abcdef01234567

make_repo() {
  local repo=$1
  mkdir -p "$repo/scripts" "$repo/build/toolchain" "$repo/web/eino-workbench"
  cp "$root/scripts/toolchain_lock.sh" "$root/scripts/build_toolchain_image.sh" \
    "$root/scripts/verify_toolchain.sh" "$repo/scripts/"
  cp "$root/build/toolchain/toolchain.lock" "$root/build/toolchain/Dockerfile" "$repo/build/toolchain/"
  cp "$root/.go-version" "$root/.node-version" "$root/go.mod" \
    "$root/package.json" "$root/package-lock.json" "$repo/"
  cp "$root/.gitlab-ci.yml" "$repo/"
  cp "$root/web/eino-workbench/package.json" "$repo/web/eino-workbench/"
  # image fixture 不启动 Go validator；该 stub 只验证 verifier 的强制调用边界。
  cat >"$repo/scripts/validate_ci_config.sh" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
exit 0
SH
  chmod +x "$repo/scripts/validate_ci_config.sh"
}

repo="$tmp/repo"
bin="$tmp/bin"
make_repo "$repo"
mkdir -p "$bin"
cat >"$bin/docker" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
case "${1:-} ${2:-}" in
  'buildx build')
    {
      printf '%s\n' 'CALL buildx build'
      printf '%s\n' "$@"
    } >>"$DOCKER_LOG"
    ;;
  'image inspect')
    printf '%s\n' 'sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
    ;;
  'buildx imagetools')
    {
      printf '%s\n' 'CALL buildx imagetools'
      printf '%s\n' "$@"
    } >>"$DOCKER_LOG"
    printf '%s\n' 'Name: registry.example/team/project/toolchain:commit'
    printf '%s\n' 'Digest: sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'
    ;;
  *) exit 64 ;;
esac
SH
chmod +x "$bin/docker"
cat >"$bin/node" <<'SH'
#!/usr/bin/env bash
exit 127
SH
chmod +x "$bin/node"

export DOCKER_LOG="$tmp/docker.log"
# 环境中的同名值不得覆盖三个权威源或公开摘要锁。
GO_VERSION=0.0.0 NODE_VERSION=0.0.0 NPM_VERSION=0.0.0 UBUNTU_SNAPSHOT=20990101T000000Z \
APT_BUILD_PACKAGES=unfixed-package \
PLAYWRIGHT_AMD64_DIGEST=sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc \
  PATH="$bin:$PATH" bash "$repo/scripts/build_toolchain_image.sh" --load >"$tmp/load.out"
grep -Fxq 'TOOLCHAIN_IMAGE_ID=sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' "$tmp/load.out"
for expected in \
  '--platform' 'linux/amd64' '--provenance=false' '--sbom=false' '--load' \
  'GO_VERSION=1.26.5' 'NODE_VERSION=24.18.0' 'NPM_VERSION=11.16.0' \
  'PLAYWRIGHT_DIGEST=sha256:cf0daee9b994042e011bc29f20cdff1a9f682a039b43fcd738f7d8a9d3bcd9d6' \
  'CHROMIUM_REVISION=1228' 'CHROMIUM_VERSION=149.0.7827.55' \
  'UBUNTU_SNAPSHOT=20260708T000000Z' 'APT_BUILD_PACKAGES=build-essential'; do
  grep -Fxq -- "$expected" "$DOCKER_LOG"
done
! grep -Eq 'GO_VERSION=0\.0\.0|NODE_VERSION=0\.0\.0|NPM_VERSION=0\.0\.0|sha256:c{64}|20990101T000000Z|unfixed-package' "$DOCKER_LOG"

assert_builder_failure() {
  local expected=$1
  shift
  if PATH="$bin:$PATH" bash "$repo/scripts/build_toolchain_image.sh" "$@" >"$tmp/fail.out" 2>"$tmp/fail.err"; then
    printf 'expected builder failure: %s\n' "$expected" >&2
    exit 1
  fi
  test ! -s "$tmp/fail.out"
  test "$(<"$tmp/fail.err")" = "toolchain image error: $expected"
  ! grep -Fq "$tmp" "$tmp/fail.err"
  ! grep -Fq 'SHOULD_NOT_LEAK' "$tmp/fail.err"
}

# GitHub identity 用例先清空宿主 CI 变量，避免本地或 runner 环境污染失败边界。
assert_builder_failure_env() {
  local expected=$1
  shift
  local -a command_env=()
  while test "$#" -gt 0 && test "$1" != --; do
    command_env+=("$1")
    shift
  done
  test "${1:-}" = --
  shift
  if env -u CI -u GITHUB_REPOSITORY -u GITHUB_SHA \
    -u CI_REGISTRY_IMAGE -u CI_COMMIT_SHA \
    "${command_env[@]}" PATH="$bin:$PATH" \
    bash "$repo/scripts/build_toolchain_image.sh" "$@" >"$tmp/fail.out" 2>"$tmp/fail.err"; then
    printf 'expected builder failure: %s\n' "$expected" >&2
    exit 1
  fi
  test ! -s "$tmp/fail.out"
  test "$(<"$tmp/fail.err")" = "toolchain image error: $expected"
  ! grep -Fq "$tmp" "$tmp/fail.err"
  ! grep -Fq 'SHOULD_NOT_LEAK' "$tmp/fail.err"
}

cp "$repo/build/toolchain/toolchain.lock" "$tmp/lock.good"
printf '%s\n' 'PLATFORM=linux/amd64' >>"$repo/build/toolchain/toolchain.lock"
SECRET=SHOULD_NOT_LEAK assert_builder_failure 'invalid toolchain lock' --load
cp "$tmp/lock.good" "$repo/build/toolchain/toolchain.lock"
sed -i.bak '/^GO_LINUX_AMD64_SHA256=/d' "$repo/build/toolchain/toolchain.lock"
go_sha256=dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd \
  assert_builder_failure 'invalid toolchain lock' --load
mv "$repo/build/toolchain/toolchain.lock.bak" "$repo/build/toolchain/toolchain.lock"
sed -i.bak '/^UBUNTU_SNAPSHOT=/d' "$repo/build/toolchain/toolchain.lock"
ubuntu_snapshot=20990101T000000Z assert_builder_failure 'invalid toolchain lock' --load
mv "$repo/build/toolchain/toolchain.lock.bak" "$repo/build/toolchain/toolchain.lock"
printf '%s\n' 'FONT_POLICY=bad value' >>"$repo/build/toolchain/toolchain.lock"
assert_builder_failure 'invalid toolchain lock' --load
cp "$tmp/lock.good" "$repo/build/toolchain/toolchain.lock"
printf '%s\n' 'PLAYWRIGHT_VERSION=$(touch${IFS}/tmp/toolchain-lock-must-not-run)' >>"$repo/build/toolchain/toolchain.lock"
rm -f /tmp/toolchain-lock-must-not-run
assert_builder_failure 'invalid toolchain lock' --load
test ! -e /tmp/toolchain-lock-must-not-run
cp "$tmp/lock.good" "$repo/build/toolchain/toolchain.lock"

# GitHub runner、action 与 Docker/BuildKit pin 都必须是不可漂移的完整值。
sed -i.bak 's/GITHUB_RUNNER=ubuntu-24\.04/GITHUB_RUNNER=ubuntu-latest/' "$repo/build/toolchain/toolchain.lock"
assert_builder_failure 'invalid toolchain lock' --load
mv "$repo/build/toolchain/toolchain.lock.bak" "$repo/build/toolchain/toolchain.lock"
for key in ACTIONS_CHECKOUT_SHA ACTIONS_UPLOAD_ARTIFACT_SHA \
  DOCKER_SETUP_DOCKER_SHA DOCKER_SETUP_BUILDX_SHA; do
  sed -i.bak "s/^${key}=.*/${key}=abcdef/" "$repo/build/toolchain/toolchain.lock"
  assert_builder_failure 'invalid toolchain lock' --load
  mv "$repo/build/toolchain/toolchain.lock.bak" "$repo/build/toolchain/toolchain.lock"
done
for mutation in \
  'DOCKER_ENGINE_VERSION=latest' \
  'DOCKER_BUILDX_VERSION=0.35' \
  'BUILDKIT_IMAGE=moby/buildkit:latest'; do
  key=${mutation%%=*}
  sed -i.bak "s#^${key}=.*#${mutation}#" "$repo/build/toolchain/toolchain.lock"
  assert_builder_failure 'invalid toolchain lock' --load
  mv "$repo/build/toolchain/toolchain.lock.bak" "$repo/build/toolchain/toolchain.lock"
done
sed -i.bak 's/^BUILDKIT_DIGEST=.*/BUILDKIT_DIGEST=sha256:abcdef/' "$repo/build/toolchain/toolchain.lock"
assert_builder_failure 'invalid toolchain lock' --load
mv "$repo/build/toolchain/toolchain.lock.bak" "$repo/build/toolchain/toolchain.lock"

assert_builder_failure 'usage: --load or strict CI --push'

commit_sha=0123456789abcdef0123456789abcdef01234567
image_ref="ghcr.io/wanggang8/agent-platform-eino/toolchain:$commit_sha"
different_sha=1123456789abcdef0123456789abcdef01234567
assert_builder_failure_env 'push requires CI' CI=false -- \
  --push "$image_ref" --env-file test-results/toolchain.env
assert_builder_failure_env 'push requires GitHub identity' \
  CI=true GITHUB_SHA=$commit_sha -- \
  --push "$image_ref" --env-file test-results/toolchain.env
assert_builder_failure_env 'invalid GitHub repository' \
  CI=true GITHUB_REPOSITORY=wanggang8/agent-platform-eino/extra GITHUB_SHA=$commit_sha -- \
  --push "$image_ref" --env-file test-results/toolchain.env
assert_builder_failure_env 'invalid GitHub commit SHA' \
  CI=true GITHUB_REPOSITORY=wanggang8/agent-platform-eino GITHUB_SHA=0123456789abcdef -- \
  --push "$image_ref" --env-file test-results/toolchain.env
assert_builder_failure_env 'invalid GitHub commit SHA' \
  CI=true GITHUB_REPOSITORY=wanggang8/agent-platform-eino \
  GITHUB_SHA=ABCDEF0123456789ABCDEF0123456789ABCDEF01 -- \
  --push "$image_ref" --env-file test-results/toolchain.env
assert_builder_failure_env 'push target must match GitHub identity' \
  CI=true GITHUB_REPOSITORY=wanggang8/agent-platform-eino GITHUB_SHA=$commit_sha -- \
  --push "registry.example/wanggang8/agent-platform-eino/toolchain:$commit_sha" \
  --env-file test-results/toolchain.env
assert_builder_failure_env 'push target must match GitHub identity' \
  CI=true GITHUB_REPOSITORY=wanggang8/agent-platform-eino GITHUB_SHA=$commit_sha -- \
  --push "ghcr.io/wanggang8/agent-platform-eino/toolchain:$different_sha" \
  --env-file test-results/toolchain.env
assert_builder_failure_env 'dotenv path must be test-results/toolchain.env' \
  CI=true GITHUB_REPOSITORY=wanggang8/agent-platform-eino GITHUB_SHA=$commit_sha -- \
  --push "$image_ref" --env-file toolchain.env
assert_builder_failure_env 'push requires GitHub identity' \
  CI=true CI_REGISTRY_IMAGE=registry.example/team/project CI_COMMIT_SHA=$commit_sha -- \
  --push "registry.example/team/project/toolchain:$commit_sha" \
  --env-file test-results/toolchain.env

: >"$DOCKER_LOG"
CI=true GITHUB_REPOSITORY=WangGang8/Agent-Platform-Eino GITHUB_SHA=$commit_sha \
  PATH="$bin:$PATH" bash "$repo/scripts/build_toolchain_image.sh" \
    --push "$image_ref" --env-file test-results/toolchain.env >"$tmp/push.out"
grep -Fxq -- '--push' "$DOCKER_LOG"
grep -Fxq 'CALL buildx imagetools' "$DOCKER_LOG"
test "$(grep -Fxc "$image_ref" "$DOCKER_LOG")" = 2
awk -v ref="$image_ref" '
  $0 == "-t" { getline; if ($0 == ref) tagged = 1 }
  END { exit(tagged ? 0 : 1) }
' "$DOCKER_LOG"
awk -v ref="$image_ref" '
  $0 == "CALL buildx imagetools" {
    getline first; getline second; getline action; getline target
    if (first == "buildx" && second == "imagetools" && action == "inspect" && target == ref) inspected = 1
  }
  END { exit(inspected ? 0 : 1) }
' "$DOCKER_LOG"
grep -Fxq "TOOLCHAIN_IMAGE=$image_ref@sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" \
  "$repo/test-results/toolchain.env"
grep -Fxq "TOOLCHAIN_IMAGE=$image_ref@sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" \
  "$tmp/push.out"
file_mode() {
  local path=$1 mode
  if mode=$(stat -c '%a' "$path" 2>/dev/null); then
    printf '%s' "$mode"
  elif mode=$(stat -f '%Lp' "$path" 2>/dev/null); then
    printf '%s' "$mode"
  else
    return 1
  fi
}
test "$(file_mode "$repo/test-results/toolchain.env")" = 600

# verifier 必须在没有 node_modules 的 clean checkout 中只使用 lock/package-lock/Dockerfile 证据。
real_node=$(command -v node)
real_npm=$(command -v npm)
export REAL_NODE="$real_node" REAL_NPM="$real_npm"
cat >"$bin/go" <<'SH'
#!/usr/bin/env bash
case "$*" in
  'version') printf '%s\n' 'go version go1.26.5 linux/amd64' ;;
  *) exit 64 ;;
esac
SH
cat >"$bin/node" <<'SH'
#!/usr/bin/env bash
if [[ ${1:-} == '--version' ]]; then printf '%s\n' 'v24.18.0'; else exec "$REAL_NODE" "$@"; fi
SH
cat >"$bin/npm" <<'SH'
#!/usr/bin/env bash
if [[ ${1:-} == '--version' ]]; then printf '%s\n' '11.16.0'; else exec "$REAL_NPM" "$@"; fi
SH
chmod +x "$bin/go" "$bin/node" "$bin/npm"
PATH="$bin:$PATH" bash "$repo/scripts/verify_toolchain.sh" >/dev/null

assert_verifier_rejects() {
  local reason=$1
  if PATH="$bin:$PATH" bash "$repo/scripts/verify_toolchain.sh" >"$tmp/verify.out" 2>"$tmp/verify.err"; then
    printf 'expected verifier to reject %s\n' "$reason" >&2
    exit 1
  fi
  test ! -s "$tmp/verify.out"
  ! grep -Fq "$tmp" "$tmp/verify.err"
}

sed -i.bak 's/PLAYWRIGHT_VERSION=1\.61\.1/PLAYWRIGHT_VERSION=1.61.0/' "$repo/build/toolchain/toolchain.lock"
assert_verifier_rejects 'Playwright lock drift'
mv "$repo/build/toolchain/toolchain.lock.bak" "$repo/build/toolchain/toolchain.lock"
printf '%s\n' 'PLATFORM=linux/amd64' >>"$repo/build/toolchain/toolchain.lock"
assert_verifier_rejects 'duplicate lock key'
cp "$tmp/lock.good" "$repo/build/toolchain/toolchain.lock"
sed -i.bak '/^PLAYWRIGHT_VERSION=/d' "$repo/build/toolchain/toolchain.lock"
playwright_version=1.61.1 assert_verifier_rejects 'environment-backed missing lock key'
mv "$repo/build/toolchain/toolchain.lock.bak" "$repo/build/toolchain/toolchain.lock"
sed -i.bak 's/UBUNTU_SNAPSHOT=20260708T000000Z/UBUNTU_SNAPSHOT=floating/' "$repo/build/toolchain/toolchain.lock"
assert_verifier_rejects 'Ubuntu snapshot drift'
mv "$repo/build/toolchain/toolchain.lock.bak" "$repo/build/toolchain/toolchain.lock"
cp "$repo/build/toolchain/Dockerfile" "$repo/build/toolchain/Dockerfile.good"
sed -i.bak 's#/ms-playwright/chromium-#/other/chromium-#' "$repo/build/toolchain/Dockerfile"
assert_verifier_rejects 'Dockerfile Chromium evidence drift'
mv "$repo/build/toolchain/Dockerfile.good" "$repo/build/toolchain/Dockerfile"
cp "$repo/build/toolchain/Dockerfile" "$repo/build/toolchain/Dockerfile.good"
sed -i.bak 's/CGO_ENABLED=1 go test -race/CGO_ENABLED=0 go test -race/' "$repo/build/toolchain/Dockerfile"
assert_verifier_rejects 'Dockerfile race smoke drift'
mv "$repo/build/toolchain/Dockerfile.good" "$repo/build/toolchain/Dockerfile"
cp "$repo/build/toolchain/Dockerfile" "$repo/build/toolchain/Dockerfile.good"
sed -i.bak 's/apt-get "${apt_snapshot_options\[@\]}" update/apt-get update/' "$repo/build/toolchain/Dockerfile"
assert_verifier_rejects 'Dockerfile snapshot update drift'
mv "$repo/build/toolchain/Dockerfile.good" "$repo/build/toolchain/Dockerfile"
