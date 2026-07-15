#!/usr/bin/env bash
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
fail() { printf 'toolchain image error: %s\n' "$1" >&2; exit 1; }

# 共享 helper 是唯一 lock 安全边界；这里只绑定其固定顺序输出，不解释 lock 内容。
read_lock() {
  local output
  if ! output=$(bash "$root/scripts/toolchain_lock.sh" \
    "$root/build/toolchain/toolchain.lock" 2>/dev/null); then
    fail 'invalid toolchain lock'
  fi
  IFS=$'\t' read -r platform playwright_image playwright_digest playwright_version \
    chromium_revision chromium_version base_os font_policy ubuntu_snapshot apt_build_packages \
    go_sha256 node_sha256 github_runner actions_checkout_sha actions_upload_artifact_sha \
    docker_setup_docker_sha docker_setup_buildx_sha docker_engine_version \
    docker_buildx_version buildkit_image buildkit_digest <<<"$output" \
    || fail 'invalid toolchain lock'
}

read_version_file() {
  local name=$1 value
  test -r "$root/$name" || fail "invalid $name"
  IFS= read -r value <"$root/$name" || fail "invalid $name"
  [[ $value =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail "invalid $name"
  test "$(wc -l <"$root/$name" | tr -d '[:space:]')" = 1 || fail "invalid $name"
  printf '%s' "$value"
}

read_lock
go_version=$(read_version_file .go-version)
node_version=$(read_version_file .node-version)
test -r "$root/package.json" || fail 'invalid packageManager'
# builder 运行于最小 Docker CLI 环境，不能在 canonical 镜像建成前反向依赖宿主 Node。
if ! package_manager=$(sed -nE \
  's/^[[:space:]]*"packageManager"[[:space:]]*:[[:space:]]*"([^"]+)"[[:space:]]*,?[[:space:]]*$/\1/p' \
  "$root/package.json" 2>/dev/null); then
  fail 'invalid packageManager'
fi
[[ $package_manager =~ ^npm@[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail 'invalid packageManager'
npm_version=${package_manager#npm@}

case ${1:-} in
  --load)
    test "$#" = 1 || fail 'usage: --load or strict CI --push'
    image_ref=agent-platform-eino-toolchain:local
    output_flag=--load
    ;;
  --push)
    test "$#" = 4 || fail 'usage: --load or strict CI --push'
    test "${3:-}" = --env-file || fail 'usage: --load or strict CI --push'
    test "${CI:-}" = true || fail 'push requires CI'
    test -n "${GITHUB_REPOSITORY:-}" && test -n "${GITHUB_SHA:-}" || fail 'push requires GitHub identity'
    [[ $GITHUB_REPOSITORY =~ ^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$ ]] || fail 'invalid GitHub repository'
    [[ $GITHUB_SHA =~ ^[0-9a-f]{40}$ ]] || fail 'invalid GitHub commit SHA'
    repository=$(printf '%s' "$GITHUB_REPOSITORY" | tr '[:upper:]' '[:lower:]')
    expected_ref="ghcr.io/$repository/toolchain:$GITHUB_SHA"
    test "$2" = "$expected_ref" || fail 'push target must match GitHub identity'
    test "$4" = test-results/toolchain.env || fail 'dotenv path must be test-results/toolchain.env'
    image_ref=$2
    env_file="$root/$4"
    output_flag=--push
    ;;
  *) fail 'usage: --load or strict CI --push' ;;
esac

docker buildx build \
  --platform "$platform" \
  --provenance=false \
  --sbom=false \
  --build-arg "PLAYWRIGHT_IMAGE=$playwright_image" \
  --build-arg "PLAYWRIGHT_DIGEST=$playwright_digest" \
  --build-arg "GO_VERSION=$go_version" \
  --build-arg "NODE_VERSION=$node_version" \
  --build-arg "NPM_VERSION=$npm_version" \
  --build-arg "GO_SHA256=$go_sha256" \
  --build-arg "NODE_SHA256=$node_sha256" \
  --build-arg "CHROMIUM_REVISION=$chromium_revision" \
  --build-arg "CHROMIUM_VERSION=$chromium_version" \
  --build-arg "UBUNTU_SNAPSHOT=$ubuntu_snapshot" \
  --build-arg "APT_BUILD_PACKAGES=$apt_build_packages" \
  -f "$root/build/toolchain/Dockerfile" \
  -t "$image_ref" \
  "$output_flag" \
  "$root"

if test "$output_flag" = --load; then
  image_id=$(docker image inspect --format '{{.Id}}' "$image_ref" 2>/dev/null) || fail 'cannot inspect loaded image'
  [[ $image_id =~ ^sha256:[0-9a-f]{64}$ ]] || fail 'invalid loaded image ID'
  printf 'TOOLCHAIN_IMAGE_ID=%s\n' "$image_id"
else
  inspect=$(docker buildx imagetools inspect "$image_ref" 2>/dev/null) || fail 'cannot inspect pushed image'
  image_digest=$(awk '$1 == "Digest:" { print $2; exit }' <<<"$inspect")
  [[ $image_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail 'invalid pushed image digest'
  mkdir -p "$(dirname "$env_file")"
  umask 077
  temp_env="${env_file}.tmp.$$"
  trap 'rm -f "$temp_env"' EXIT
  printf 'TOOLCHAIN_IMAGE=%s@%s\n' "$image_ref" "$image_digest" >"$temp_env"
  chmod 0600 "$temp_env"
  mv "$temp_env" "$env_file"
  trap - EXIT
  printf 'TOOLCHAIN_IMAGE=%s@%s\n' "$image_ref" "$image_digest"
fi
