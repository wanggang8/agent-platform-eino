#!/usr/bin/env bash
set -euo pipefail

fail() { printf 'toolchain lock error: invalid lock\n' >&2; exit 1; }
test "$#" = 1 || fail
lock=$1

# lock 属于不可信数据：只按 KEY=VALUE 读取，不允许执行、展开或回显原始内容。
platform= playwright_image= playwright_digest= playwright_version=
chromium_revision= chromium_version= base_os= font_policy=
go_sha256= node_sha256= docker_cli_image= docker_cli_digest=
docker_dind_image= docker_dind_digest=
seen_keys='|'
test -r "$lock" || fail
while IFS= read -r line || test -n "$line"; do
  [[ $line =~ ^([A-Z][A-Z0-9_]*)=([^[:space:]]+)$ ]] || fail
  key=${BASH_REMATCH[1]}
  value=${BASH_REMATCH[2]}
  [[ $value =~ ^[A-Za-z0-9._/:@+-]+$ ]] || fail
  [[ $seen_keys != *"|$key|"* ]] || fail
  seen_keys="${seen_keys}${key}|"
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
    *) fail ;;
  esac
done <"$lock"

for required in platform playwright_image playwright_digest playwright_version \
  chromium_revision chromium_version base_os font_policy go_sha256 node_sha256 \
  docker_cli_image docker_cli_digest docker_dind_image docker_dind_digest; do
  test -n "${!required:-}" || fail
done
[[ $platform == linux/amd64 ]] || fail
[[ $playwright_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
[[ $chromium_revision =~ ^[0-9]+$ ]] || fail
[[ $chromium_version =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
[[ $playwright_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail
[[ $docker_cli_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail
[[ $docker_dind_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail
[[ $go_sha256 =~ ^[0-9a-f]{64}$ ]] || fail
[[ $node_sha256 =~ ^[0-9a-f]{64}$ ]] || fail

# 固定字段顺序是 helper 与两个消费者之间唯一的可信接口，值域已禁止 tab/newline。
printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
  "$platform" "$playwright_image" "$playwright_digest" "$playwright_version" \
  "$chromium_revision" "$chromium_version" "$base_os" "$font_policy" \
  "$go_sha256" "$node_sha256" "$docker_cli_image" "$docker_cli_digest" \
  "$docker_dind_image" "$docker_dind_digest"
