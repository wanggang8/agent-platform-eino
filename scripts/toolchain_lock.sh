#!/usr/bin/env bash
set -euo pipefail

fail() { printf 'toolchain lock error: invalid lock\n' >&2; exit 1; }
test "$#" = 1 || fail
lock=$1

# lock 属于不可信数据：只按 KEY=VALUE 读取，不允许执行、展开或回显原始内容。
platform= playwright_image= playwright_digest= playwright_version=
chromium_revision= chromium_version= base_os= font_policy=
ubuntu_snapshot= apt_build_packages=
go_sha256= node_sha256=
github_runner= actions_checkout_sha= actions_upload_artifact_sha=
docker_setup_docker_sha= docker_setup_buildx_sha=
docker_engine_version= docker_buildx_version= buildkit_image= buildkit_digest=
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
    UBUNTU_SNAPSHOT) ubuntu_snapshot=$value ;;
    APT_BUILD_PACKAGES) apt_build_packages=$value ;;
    GO_LINUX_AMD64_SHA256) go_sha256=$value ;;
    NODE_LINUX_X64_SHA256) node_sha256=$value ;;
    GITHUB_RUNNER) github_runner=$value ;;
    ACTIONS_CHECKOUT_SHA) actions_checkout_sha=$value ;;
    ACTIONS_UPLOAD_ARTIFACT_SHA) actions_upload_artifact_sha=$value ;;
    DOCKER_SETUP_DOCKER_SHA) docker_setup_docker_sha=$value ;;
    DOCKER_SETUP_BUILDX_SHA) docker_setup_buildx_sha=$value ;;
    DOCKER_ENGINE_VERSION) docker_engine_version=$value ;;
    DOCKER_BUILDX_VERSION) docker_buildx_version=$value ;;
    BUILDKIT_IMAGE) buildkit_image=$value ;;
    BUILDKIT_DIGEST) buildkit_digest=$value ;;
    *) fail ;;
  esac
done <"$lock"

for required in platform playwright_image playwright_digest playwright_version \
  chromium_revision chromium_version base_os font_policy ubuntu_snapshot apt_build_packages \
  go_sha256 node_sha256 \
  github_runner actions_checkout_sha actions_upload_artifact_sha \
  docker_setup_docker_sha docker_setup_buildx_sha docker_engine_version \
  docker_buildx_version buildkit_image buildkit_digest; do
  test -n "${!required:-}" || fail
done
[[ $platform == linux/amd64 ]] || fail
[[ $playwright_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
[[ $chromium_revision =~ ^[0-9]+$ ]] || fail
[[ $chromium_version =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
[[ $ubuntu_snapshot =~ ^[0-9]{8}T[0-9]{6}Z$ ]] || fail
[[ $apt_build_packages =~ ^[a-z0-9][a-z0-9+.-]*$ ]] || fail
[[ $playwright_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail
[[ $go_sha256 =~ ^[0-9a-f]{64}$ ]] || fail
[[ $node_sha256 =~ ^[0-9a-f]{64}$ ]] || fail
[[ $github_runner == ubuntu-24.04 ]] || fail
for action_sha in "$actions_checkout_sha" "$actions_upload_artifact_sha" \
  "$docker_setup_docker_sha" "$docker_setup_buildx_sha"; do
  [[ $action_sha =~ ^[0-9a-f]{40}$ ]] || fail
done
[[ $docker_engine_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
[[ $docker_buildx_version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
[[ $buildkit_image =~ ^moby/buildkit:v[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail
[[ $buildkit_digest =~ ^sha256:[0-9a-f]{64}$ ]] || fail

# 固定字段顺序是 helper 与两个消费者之间唯一的可信接口，值域已禁止 tab/newline。
printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' \
  "$platform" "$playwright_image" "$playwright_digest" "$playwright_version" \
  "$chromium_revision" "$chromium_version" "$base_os" "$font_policy" \
  "$ubuntu_snapshot" "$apt_build_packages" \
  "$go_sha256" "$node_sha256" "$github_runner" "$actions_checkout_sha" \
  "$actions_upload_artifact_sha" "$docker_setup_docker_sha" "$docker_setup_buildx_sha" \
  "$docker_engine_version" "$docker_buildx_version" "$buildkit_image" "$buildkit_digest"
