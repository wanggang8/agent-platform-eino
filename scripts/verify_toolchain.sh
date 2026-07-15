#!/usr/bin/env bash
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
fail() { printf 'toolchain mismatch: %s\n' "$1" >&2; exit 1; }
read_one() {
  local file=$1 relative value
  relative=${file#"$root"/}
  test -f "$file" || fail "missing $relative"
  if ! value=$(tr -d '[:space:]' 2>/dev/null <"$file"); then
    fail "invalid $relative"
  fi
  printf '%s' "$value"
}

# go.mod 读取异常与缺失指令都在 shell 边界收敛，避免 awk 泄露本机路径。
read_go_directive() {
  local key=$1 value
  test -r "$root/go.mod" || fail 'invalid go.mod'
  if ! value=$(awk -v key="$key" '
    $1 == key { count++; value = $2 }
    END {
      if (count != 1 || value == "") exit 1
      print value
    }
  ' "$root/go.mod" 2>/dev/null); then
    fail 'invalid go.mod'
  fi
  printf '%s' "$value"
}

# 共享 helper 集中维护 lock 白名单与摘要格式；verifier 只绑定固定安全输出。
read_container_lock() {
  local output
  if ! output=$(bash "$root/scripts/toolchain_lock.sh" \
    "$root/build/toolchain/toolchain.lock" 2>/dev/null); then
    fail 'invalid toolchain.lock'
  fi
  IFS=$'\t' read -r container_platform playwright_image playwright_digest playwright_version \
    chromium_revision chromium_version base_os font_policy ubuntu_snapshot apt_build_packages \
    go_linux_amd64_sha256 node_linux_x64_sha256 github_runner actions_checkout_sha \
    actions_upload_artifact_sha docker_setup_docker_sha docker_setup_buildx_sha \
    docker_engine_version docker_buildx_version buildkit_image buildkit_digest \
    <<<"$output" || fail 'invalid toolchain.lock'
}

# Dockerfile 只能消费传入参数，并在镜像内验证 pinned Chromium，不读取宿主安装产物。
verify_container_dockerfile() {
  local dockerfile="$root/build/toolchain/Dockerfile" required forbidden
  test -r "$dockerfile" || fail 'invalid toolchain Dockerfile'
  for required in \
    'FROM ${PLAYWRIGHT_IMAGE}@${PLAYWRIGHT_DIGEST}' \
    'ARG CHROMIUM_REVISION' \
    'ARG CHROMIUM_VERSION' \
    'ARG UBUNTU_SNAPSHOT' \
    'ARG APT_BUILD_PACKAGES' \
    'node-v${NODE_VERSION}-linux-x64.tar.gz' \
    'tar -C /usr/local --strip-components=1 -xzf /tmp/node.tar.gz' \
    'Snapshot: ${UBUNTU_SNAPSHOT}' \
    'test "$snapshot_count" = "$signed_by_count"' \
    "! -name 'ubuntu.sources' -delete" \
    'Dir::Etc::sourcelist="$sources"' \
    'Dir::Etc::sourceparts="-"' \
    'apt-get "${apt_snapshot_options[@]}" update' \
    'apt-get "${apt_snapshot_options[@]}" install -y --no-install-recommends "$APT_BUILD_PACKAGES"' \
    'rm -rf /var/lib/apt/lists/*' \
    'command -v cc' \
    'go mod init toolchain-race-smoke' \
    'CGO_ENABLED=1 go test -race ./...' \
    'find "/ms-playwright/chromium-${CHROMIUM_REVISION}"' \
    'chromium_actual=$("$chromium_binary" --version)' \
    '[[ "$chromium_actual" == *"${CHROMIUM_VERSION}"* ]]'; do
    grep -Fq "$required" "$dockerfile" 2>/dev/null || fail 'invalid toolchain Dockerfile'
  done
  for forbidden in "$go_version" "$node_version" "$npm_version" "$chromium_revision" \
    "$chromium_version" "$playwright_digest" "$ubuntu_snapshot" "$apt_build_packages" \
    'node_modules/playwright-core/browsers.json' 'tar.xz' 'xJf'; do
    if grep -Fq "$forbidden" "$dockerfile" 2>/dev/null; then
      fail 'invalid toolchain Dockerfile'
    fi
  done
}

# CI 必须复用 lock 与公共 validator，禁止复制三项 runtime 权威版本形成第二套真相。
verify_ci_config() {
  local ci="$root/.github/workflows/toolchain.yml" validator="$root/scripts/validate_ci_config.sh" forbidden target grep_status
  test ! -e "$root/.gitlab-ci.yml" || fail 'GitLab CI residue is forbidden'
  for target in "$ci" "$root/scripts/build_toolchain_image.sh" \
    "$root/scripts/run_toolchain_baseline.sh" "$root/scripts/validate_ci_config.sh" \
    "$root/scripts/docker-credential-github-token"; do
    test -r "$target" || fail 'GitHub CI production entry is unreadable'
    if grep -Eq 'CI_(COMMIT_SHA|REGISTRY|PROJECT_DIR)|GitLab|\.gitlab-ci' "$target" 2>/dev/null; then
      fail 'GitLab CI residue is forbidden'
    else
      grep_status=$?
      test "$grep_status" = 1 || fail 'cannot scan GitHub CI production entry'
    fi
  done
  for forbidden in "$go_version" "$node_version" "$npm_version"; do
    if grep -Fq "$forbidden" "$ci" 2>/dev/null; then
      fail 'CI duplicates an authoritative runtime version'
    fi
  done
  if ! bash "$validator"; then
    fail 'invalid GitHub Actions workflow'
  fi
}

# 版本文件与根 packageManager 是唯一权威源，其余声明只做一致性镜像。
go_version=$(read_one "$root/.go-version")
node_version=$(read_one "$root/.node-version")
go_language="${go_version%.*}.0"
read_container_lock

# 固定 local 可阻止宿主 Go 自动下载目标 toolchain 后伪装为环境合规。
actual_go_line=$(GOTOOLCHAIN=local go version 2>/dev/null || true)
actual_go=$(awk '{print $3}' <<<"$actual_go_line")
actual_node=$(node --version 2>/dev/null || true)
actual_npm=$(npm --version 2>/dev/null || true)

[[ $actual_go == "go$go_version" ]] || fail "go expected=go$go_version actual=${actual_go:-missing}"
[[ $actual_node == "v$node_version" ]] || fail "node expected=v$node_version actual=${actual_node:-missing}"

declared_go=$(read_go_directive go)
declared_toolchain=$(read_go_directive toolchain)
[[ $declared_go == "$go_language" ]] || fail "go.mod language expected=$go_language actual=${declared_go:-missing}"
[[ $declared_toolchain == "go$go_version" ]] || fail "go.mod toolchain expected=go$go_version actual=${declared_toolchain:-missing}"

node - "$root" "$node_version" "$actual_npm" "$playwright_version" <<'NODE'
const fs = require('node:fs');
const path = require('node:path');
const [root, nodeVersion, actualNpm, playwrightVersion] = process.argv.slice(2);
const fail = (message) => { console.error(`toolchain mismatch: ${message}`); process.exit(1); };

// 仓库声明也视为不可信输入，解析异常只暴露稳定的相对文件名。
const read = (name) => {
  try {
    return JSON.parse(fs.readFileSync(path.join(root, name), 'utf8'));
  } catch {
    fail(`invalid ${name}`);
  }
};
const rootPackage = read('package.json');
const webPackage = read('web/eino-workbench/package.json');
const lock = read('package-lock.json');

// npm 版本只从根 packageManager 派生，避免 verifier 引入第二个工具链常量。
if (!/^npm@\d+\.\d+\.\d+$/.test(String(rootPackage.packageManager || ''))) fail('packageManager must be an exact npm version');
const npmVersion = rootPackage.packageManager.slice('npm@'.length);
if (actualNpm !== npmVersion) fail(`npm expected=${npmVersion} actual=${actualNpm || 'missing'}`);
for (const [name, pkg] of [['package.json', rootPackage], ['web/eino-workbench/package.json', webPackage]]) {
  if (pkg.engines?.node !== nodeVersion) fail(`${name} engines.node expected=${nodeVersion} actual=${pkg.engines?.node || 'missing'}`);
  if (pkg.engines?.npm !== npmVersion) fail(`${name} engines.npm expected=${npmVersion} actual=${pkg.engines?.npm || 'missing'}`);
}

// Node 类型版本是批准的依赖契约常量，manifest 与 lockfile 必须精确一致。
if (rootPackage.devDependencies?.['@types/node'] !== '24.13.3') fail('@types/node must equal 24.13.3');
const lockedRoot = lock.packages?.[''];
if (lockedRoot?.devDependencies?.['@types/node'] !== '24.13.3') fail('lock root @types/node must equal 24.13.3');
if (lock.packages?.['node_modules/@types/node']?.version !== '24.13.3') fail('resolved @types/node must equal 24.13.3');
for (const [name, pkg] of [['lock root', lockedRoot], ['lock workspace', lock.packages?.['web/eino-workbench']]]) {
  if (pkg?.engines?.node !== nodeVersion) fail(`${name} engines.node expected=${nodeVersion} actual=${pkg?.engines?.node || 'missing'}`);
  if (pkg?.engines?.npm !== npmVersion) fail(`${name} engines.npm expected=${npmVersion} actual=${pkg?.engines?.npm || 'missing'}`);
}
for (const name of ['@playwright/test', 'playwright', 'playwright-core']) {
  const actual = lock.packages?.[`node_modules/${name}`]?.version;
  if (actual !== playwrightVersion) fail(`${name} expected=${playwrightVersion} actual=${actual || 'missing'}`);
}
NODE

npm_version=$actual_npm
verify_container_dockerfile
verify_ci_config

printf 'toolchain verified: go=%s node=%s npm=%s\n' "$go_version" "$node_version" "$actual_npm"
