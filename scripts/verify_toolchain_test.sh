#!/usr/bin/env bash
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

mkdir -p "$tmp/repo/scripts" "$tmp/repo/web/eino-workbench" "$tmp/bin"
cp "$root/scripts/verify_toolchain.sh" "$tmp/repo/scripts/"
cp "$root/.go-version" "$root/.node-version" "$root/go.mod" \
  "$root/package.json" "$root/package-lock.json" "$tmp/repo/"
cp "$root/web/eino-workbench/package.json" "$tmp/repo/web/eino-workbench/"

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

PATH="$tmp/bin:$PATH" bash "$tmp/repo/scripts/verify_toolchain.sh"

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
