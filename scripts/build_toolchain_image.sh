#!/usr/bin/env bash
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
fail() { printf 'toolchain image error: %s\n' "$1" >&2; exit 1; }

# lock 是外部公开摘要的只读输入；逐行白名单解析，绝不执行其中内容。
read_lock() {
  local line key value
  platform= playwright_image= playwright_digest= playwright_version=
  chromium_revision= chromium_version= base_os= font_policy=
  go_sha256= node_sha256= docker_cli_image= docker_cli_digest=
  docker_dind_image= docker_dind_digest=
  seen_lock_keys='|'
  test -r "$root/build/toolchain/toolchain.lock" || fail 'invalid toolchain lock'
  while IFS= read -r line || test -n "$line"; do
    [[ $line =~ ^([A-Z][A-Z0-9_]*)=([^[:space:]]+)$ ]] || fail 'invalid toolchain lock'
    key=${BASH_REMATCH[1]}
    value=${BASH_REMATCH[2]}
    [[ $value =~ ^[A-Za-z0-9._/:@+-]+$ ]] || fail 'invalid toolchain lock'
    [[ $seen_lock_keys != *"|$key|"* ]] || fail 'duplicate toolchain lock key'
    seen_lock_keys="${seen_lock_keys}${key}|"
    case "$key" in
      PLATFORM) platform=$value ;;
      PLAYWRIGHT_IMAGE) playwright_image=$value ;;
      PLAYWRIGHT_AMD64_DIGEST) playwright_digest=$value ;;
      PLAYWRIGHT_VERSION) playwright_version=$value ;;
      CHROMIUM_REVISION) chromium_revision=$value ;;
      CHROMIUM_VERSION) chromium_version=$value ;;
      BASE_OS) base_os=$value ;;
      FONT_POLICY) font_policy=$value ;;
      GO_LINUX_AMD64_SHA256) go_sha256=$value ;;
      NODE_LINUX_X64_SHA256) node_sha256=$value ;;
      DOCKER_CLI_IMAGE) docker_cli_image=$value ;;
      DOCKER_CLI_AMD64_DIGEST) docker_cli_digest=$value ;;
      DOCKER_DIND_IMAGE) docker_dind_image=$value ;;
      DOCKER_DIND_AMD64_DIGEST) docker_dind_digest=$value ;;
      *) fail 'unknown toolchain lock key' ;;
    esac
  done <"$root/build/toolchain/toolchain.lock"

  for key in platform playwright_image playwright_digest playwright_version \
    chromium_revision chromium_version base_os font_policy go_sha256 node_sha256 \
    docker_cli_image docker_cli_digest docker_dind_image docker_dind_digest; do
    test -n "${!key:-}" || fail 'missing toolchain lock key'
  done
  [[ $platform == linux/amd64 ]] || fail 'platform must be linux/amd64'
  [[ $playwright_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail 'invalid Playwright digest'
  [[ $docker_cli_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail 'invalid Docker CLI digest'
  [[ $docker_dind_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail 'invalid Docker DinD digest'
  [[ $go_sha256 =~ ^[0-9a-f]{64}$ ]] || fail 'invalid Go checksum'
  [[ $node_sha256 =~ ^[0-9a-f]{64}$ ]] || fail 'invalid Node checksum'
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
    test -n "${CI_REGISTRY_IMAGE:-}" && test -n "${CI_COMMIT_SHA:-}" || fail 'push requires GitLab identity'
    [[ $CI_REGISTRY_IMAGE =~ ^[A-Za-z0-9._/:@-]+$ ]] || fail 'invalid GitLab registry image'
    [[ $CI_COMMIT_SHA =~ ^[0-9a-f]{40}$ ]] || fail 'invalid GitLab commit SHA'
    expected_ref="$CI_REGISTRY_IMAGE/toolchain:$CI_COMMIT_SHA"
    test "$2" = "$expected_ref" || fail 'push target must match GitLab identity'
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
