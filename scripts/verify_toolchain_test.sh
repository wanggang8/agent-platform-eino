#!/usr/bin/env bash
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

mkdir -p "$tmp/repo/scripts" "$tmp/repo/build/toolchain" "$tmp/repo/web/eino-workbench" \
  "$tmp/repo/.github/workflows" "$tmp/bin"
cp "$root/scripts/verify_toolchain.sh" "$root/scripts/toolchain_lock.sh" \
  "$root/scripts/build_toolchain_image.sh" "$root/scripts/run_toolchain_baseline.sh" \
  "$tmp/repo/scripts/"
cp "$root/build/toolchain/toolchain.lock" "$root/build/toolchain/Dockerfile" \
  "$tmp/repo/build/toolchain/"
cp "$root/.go-version" "$root/.node-version" "$root/go.mod" \
  "$root/package.json" "$root/package-lock.json" "$tmp/repo/"
cp "$root/.github/workflows/toolchain.yml" "$tmp/repo/.github/workflows/"
cp "$root/web/eino-workbench/package.json" "$tmp/repo/web/eino-workbench/"

# fixture 使用专用 stub 证明生产 verifier 总会调用静态 CI validator，不保留绕过开关。
cat >"$tmp/repo/scripts/validate_ci_config.sh" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' called >>"$CI_VALIDATOR_LOG"
SH
chmod +x "$tmp/repo/scripts/validate_ci_config.sh"

real_node=$(command -v node)
real_npm=$(command -v npm)
export REAL_NODE="$real_node" REAL_NPM="$real_npm"

# PATH shim 仅存在于临时仓库，用真实 Node 执行 JSON 契约校验，不给生产脚本留绕过开关。
cat >"$tmp/bin/go" <<'SH'
#!/usr/bin/env bash
case "$*" in
  'env GOVERSION') printf '%s\n' 'go1.26.5' ;;
  'version') printf '%s\n' 'go version go1.26.5 linux/amd64' ;;
  *) exit 64 ;;
esac
SH
cat >"$tmp/bin/node" <<'SH'
#!/usr/bin/env bash
if [[ ${1:-} == '--version' ]]; then printf '%s\n' 'v24.18.0'; else exec "$REAL_NODE" "$@"; fi
SH
cat >"$tmp/bin/npm" <<'SH'
#!/usr/bin/env bash
if [[ ${1:-} == '--version' ]]; then printf '%s\n' '11.16.0'; else exec "$REAL_NPM" "$@"; fi
SH
chmod +x "$tmp/bin/go" "$tmp/bin/node" "$tmp/bin/npm"

export CI_VALIDATOR_LOG="$tmp/ci-validator.log"
PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"
grep -Fxq called "$CI_VALIDATOR_LOG"

printf '%s\n' 'stages: [toolchain, verify]' >"$tmp/repo/.gitlab-ci.yml"
if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"; then
  echo 'expected GitLab CI residue to fail' >&2
  exit 1
fi
rm "$tmp/repo/.gitlab-ci.yml"

# 四个生产入口必须 fail-closed；缺失入口只能返回固定摘要且不得泄露临时绝对路径。
assert_unreadable_ci_entry_failure() {
  local entry=$1
  local stdout_file="$tmp/unreadable-entry.stdout"
  local stderr_file="$tmp/unreadable-entry.stderr"
  local actual rc

  mv "$tmp/repo/$entry" "$tmp/repo/$entry.good"
  if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh" \
    >"$stdout_file" 2>"$stderr_file"; then
    echo "expected unreadable CI entry to fail: $entry" >&2
    exit 1
  else
    rc=$?
  fi
  mv "$tmp/repo/$entry.good" "$tmp/repo/$entry"

  actual=$(<"$stderr_file")
  [[ $rc -ne 0 ]] || { echo "expected non-zero unreadable entry failure: $entry" >&2; exit 1; }
  [[ ! -s $stdout_file ]] || { echo "expected empty stdout for unreadable entry: $entry" >&2; exit 1; }
  [[ $actual == 'toolchain mismatch: GitHub CI production entry is unreadable' ]] || {
    printf 'unexpected unreadable entry error for %s: %q\n' "$entry" "$actual" >&2
    exit 1
  }
  [[ $actual != *"$tmp"* ]] || { echo 'unreadable entry error leaked temporary path' >&2; exit 1; }
}

assert_unreadable_ci_entry_failure scripts/build_toolchain_image.sh
assert_unreadable_ci_entry_failure scripts/run_toolchain_baseline.sh

# CI 不得复制三项 authoritative sources，避免形成第二组 runtime 常量。
cp "$tmp/repo/.github/workflows/toolchain.yml" "$tmp/repo/.github/workflows/toolchain.yml.good"
printf '%s\n' '# duplicated Go version: 1.26.5' >>"$tmp/repo/.github/workflows/toolchain.yml"
if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"; then
  echo 'expected duplicated CI Go version to fail' >&2
  exit 1
fi
mv "$tmp/repo/.github/workflows/toolchain.yml.good" "$tmp/repo/.github/workflows/toolchain.yml"

# 逐项替换实际版本或声明，证明 Go 1.23、Node 25、npm 错版和声明漂移都会被拒绝。
cp "$tmp/bin/go" "$tmp/bin/go.good"
cat >"$tmp/bin/go" <<'SH'
#!/usr/bin/env bash
case "$*" in
  'env GOVERSION') printf '%s\n' 'go1.23.12' ;;
  'version') printf '%s\n' 'go version go1.23.12 linux/amd64' ;;
  *) exit 64 ;;
esac
SH
chmod +x "$tmp/bin/go"
if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"; then
  echo 'expected Go 1.23 to fail' >&2
  exit 1
fi
mv "$tmp/bin/go.good" "$tmp/bin/go"

cp "$tmp/bin/node" "$tmp/bin/node.good"
cat >"$tmp/bin/node" <<'SH'
#!/usr/bin/env bash
if [[ ${1:-} == '--version' ]]; then printf '%s\n' 'v25.8.1'; else exec "$REAL_NODE" "$@"; fi
SH
chmod +x "$tmp/bin/node"
if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"; then
  echo 'expected Node 25 to fail' >&2
  exit 1
fi
mv "$tmp/bin/node.good" "$tmp/bin/node"

sed -i.bak 's/24\.18\.0/25.8.1/' "$tmp/repo/.node-version"
if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"; then
  echo 'expected Node declaration drift to fail' >&2
  exit 1
fi
mv "$tmp/repo/.node-version.bak" "$tmp/repo/.node-version"

sed -i.bak 's/toolchain go1\.26\.5/toolchain go1.26.4/' "$tmp/repo/go.mod"
if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"; then
  echo 'expected Go toolchain drift to fail' >&2
  exit 1
fi
mv "$tmp/repo/go.mod.bak" "$tmp/repo/go.mod"

cp "$tmp/bin/npm" "$tmp/bin/npm.good"
cat >"$tmp/bin/npm" <<'SH'
#!/usr/bin/env bash
if [[ ${1:-} == '--version' ]]; then printf '%s\n' '11.15.0'; else exec "$REAL_NPM" "$@"; fi
SH
chmod +x "$tmp/bin/npm"
if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"; then
  echo 'expected npm version mismatch to fail' >&2
  exit 1
fi
mv "$tmp/bin/npm.good" "$tmp/bin/npm"

"$REAL_NODE" - "$tmp/repo/package.json" <<'NODE'
const fs = require('node:fs');
const path = process.argv[2];
const value = JSON.parse(fs.readFileSync(path, 'utf8'));
value.engines.node = '>=20';
fs.writeFileSync(path, JSON.stringify(value, null, 2) + '\n');
NODE
if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"; then
  echo 'expected manifest drift to fail' >&2
  exit 1
fi

# 解析失败只能返回稳定的相对文件摘要，禁止泄露临时绝对路径、原始异常或 stack。
assert_safe_input_failure() {
  local expected=$1
  local stdout_file="$tmp/verifier.stdout"
  local stderr_file="$tmp/verifier.stderr"
  local actual rc

  if PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh" >"$stdout_file" 2>"$stderr_file"; then
    echo "expected safe input failure: $expected" >&2
    exit 1
  else
    rc=$?
  fi
  actual=$(<"$stderr_file")
  [[ $rc -ne 0 ]] || { echo "expected non-zero input failure: $expected" >&2; exit 1; }
  [[ ! -s $stdout_file ]] || { echo "expected empty stdout: $expected" >&2; exit 1; }
  [[ $actual == "$expected" ]] || {
    printf 'expected safe error %q, got %q\n' "$expected" "$actual" >&2
    exit 1
  }
  [[ $actual != *"$tmp"* ]] || { echo 'safe error leaked temporary path' >&2; exit 1; }
  [[ $actual != *' at '* ]] || { echo 'safe error leaked stack trace' >&2; exit 1; }
}

cp "$root/package.json" "$tmp/repo/package.json"

mv "$tmp/repo/go.mod" "$tmp/repo/go.mod.good"
assert_safe_input_failure 'toolchain mismatch: invalid go.mod'
mv "$tmp/repo/go.mod.good" "$tmp/repo/go.mod"

cp "$tmp/repo/go.mod" "$tmp/repo/go.mod.good"
: >"$tmp/repo/go.mod"
assert_safe_input_failure 'toolchain mismatch: invalid go.mod'
mv "$tmp/repo/go.mod.good" "$tmp/repo/go.mod"

cp "$tmp/repo/package.json" "$tmp/repo/package.json.good"
printf '%s\n' '{invalid json' >"$tmp/repo/package.json"
assert_safe_input_failure 'toolchain mismatch: invalid package.json'
mv "$tmp/repo/package.json.good" "$tmp/repo/package.json"
