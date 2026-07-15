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

# canonical lock 只允许公开摘要字段；解析过程不执行 lock 内容，也不回显不可信值。
read_container_lock() {
  local line key value required
  container_platform= playwright_image= playwright_digest= playwright_version=
  chromium_revision= chromium_version= base_os= font_policy=
  go_linux_amd64_sha256= node_linux_x64_sha256=
  docker_cli_image= docker_cli_digest= docker_dind_image= docker_dind_digest=
  seen_container_keys='|'
  test -r "$root/build/toolchain/toolchain.lock" || fail 'invalid toolchain.lock'
  while IFS= read -r line || test -n "$line"; do
    [[ $line =~ ^([A-Z][A-Z0-9_]*)=([^[:space:]]+)$ ]] || fail 'invalid toolchain.lock'
    key=${BASH_REMATCH[1]}
    value=${BASH_REMATCH[2]}
    [[ $value =~ ^[A-Za-z0-9._/:@+-]+$ ]] || fail 'invalid toolchain.lock'
    [[ $seen_container_keys != *"|$key|"* ]] || fail 'invalid toolchain.lock'
    seen_container_keys="${seen_container_keys}${key}|"
    case "$key" in
      PLATFORM) container_platform=$value ;;
      PLAYWRIGHT_IMAGE) playwright_image=$value ;;
      PLAYWRIGHT_AMD64_DIGEST) playwright_digest=$value ;;
      PLAYWRIGHT_VERSION) playwright_version=$value ;;
      CHROMIUM_REVISION) chromium_revision=$value ;;
      CHROMIUM_VERSION) chromium_version=$value ;;
      BASE_OS) base_os=$value ;;
      FONT_POLICY) font_policy=$value ;;
      GO_LINUX_AMD64_SHA256) go_linux_amd64_sha256=$value ;;
      NODE_LINUX_X64_SHA256) node_linux_x64_sha256=$value ;;
      DOCKER_CLI_IMAGE) docker_cli_image=$value ;;
      DOCKER_CLI_AMD64_DIGEST) docker_cli_digest=$value ;;
      DOCKER_DIND_IMAGE) docker_dind_image=$value ;;
      DOCKER_DIND_AMD64_DIGEST) docker_dind_digest=$value ;;
      *) fail 'invalid toolchain.lock' ;;
    esac
  done <"$root/build/toolchain/toolchain.lock"

  for required in container_platform playwright_image playwright_digest playwright_version \
    chromium_revision chromium_version base_os font_policy go_linux_amd64_sha256 \
    node_linux_x64_sha256 docker_cli_image docker_cli_digest docker_dind_image docker_dind_digest; do
    test -n "${!required:-}" || fail 'invalid toolchain.lock'
  done
  [[ $container_platform == linux/amd64 ]] || fail 'invalid toolchain.lock'
  [[ $playwright_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail 'invalid toolchain.lock'
  [[ $chromium_revision =~ ^[0-9]+$ ]] || fail 'invalid toolchain.lock'
  [[ $chromium_version =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail 'invalid toolchain.lock'
  [[ $playwright_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail 'invalid toolchain.lock'
  [[ $docker_cli_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail 'invalid toolchain.lock'
  [[ $docker_dind_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail 'invalid toolchain.lock'
  [[ $go_linux_amd64_sha256 =~ ^[0-9a-f]{64}$ ]] || fail 'invalid toolchain.lock'
  [[ $node_linux_x64_sha256 =~ ^[0-9a-f]{64}$ ]] || fail 'invalid toolchain.lock'
}

# Dockerfile 只能消费传入参数，并在镜像内验证 pinned Chromium，不读取宿主安装产物。
verify_container_dockerfile() {
  local dockerfile="$root/build/toolchain/Dockerfile" required forbidden
  test -r "$dockerfile" || fail 'invalid toolchain Dockerfile'
  for required in \
    'FROM ${PLAYWRIGHT_IMAGE}@${PLAYWRIGHT_DIGEST}' \
    'ARG CHROMIUM_REVISION' \
    'ARG CHROMIUM_VERSION' \
    'find "/ms-playwright/chromium-${CHROMIUM_REVISION}"' \
    'chromium_actual=$("$chromium_binary" --version)' \
    '[[ "$chromium_actual" == *"${CHROMIUM_VERSION}"* ]]'; do
    grep -Fq "$required" "$dockerfile" 2>/dev/null || fail 'invalid toolchain Dockerfile'
  done
  for forbidden in "$go_version" "$node_version" "$npm_version" "$chromium_revision" \
    "$chromium_version" "$playwright_digest" 'node_modules/playwright-core/browsers.json'; do
    if grep -Fq "$forbidden" "$dockerfile" 2>/dev/null; then
      fail 'invalid toolchain Dockerfile'
    fi
  done
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

printf 'toolchain verified: go=%s node=%s npm=%s\n' "$go_version" "$node_version" "$actual_npm"
