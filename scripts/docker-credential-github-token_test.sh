#!/usr/bin/env bash
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
helper="$root/scripts/docker-credential-github-token"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

run_failure() {
  local action=$1 input=$2 expected=$3 stdout_file="$tmp/stdout" stderr_file="$tmp/stderr" rc
  : >"$stdout_file"
  : >"$stderr_file"
  if printf '%s' "$input" | HOME="$tmp/home" DOCKER_CONFIG="$tmp/docker" \
    GHCR_ACTOR="${TEST_ACTOR-}" GHCR_TOKEN="${TEST_TOKEN-}" \
    "$helper" "$action" >"$stdout_file" 2>"$stderr_file"; then
    echo "expected credential helper failure: $action" >&2
    exit 1
  else
    rc=$?
  fi
  [[ $rc -ne 0 ]]
  [[ ! -s $stdout_file ]]
  [[ $(<"$stderr_file") == "$expected" ]]
  [[ $(<"$stderr_file") != *'fake_token_123'* ]]
  [[ $(<"$stderr_file") != *"$tmp"* ]]
}

# 正向路径只允许从当前进程环境读取凭据，且不得触碰 HOME 或 Docker 配置目录。
mkdir -p "$tmp/home" "$tmp/docker"
output=$(printf '%s' ghcr.io | HOME="$tmp/home" DOCKER_CONFIG="$tmp/docker" \
  GHCR_ACTOR=github-actions GHCR_TOKEN=fake_token_123 "$helper" get)
[[ $output == '{"Username":"github-actions","Secret":"fake_token_123"}' ]]
[[ -z $(find "$tmp/home" "$tmp/docker" -mindepth 1 -print -quit) ]]

unset TEST_ACTOR TEST_TOKEN
run_failure get ghcr.io 'credential helper: credentials are unavailable'
TEST_ACTOR=github-actions
run_failure get ghcr.io 'credential helper: credentials are unavailable'
TEST_TOKEN=fake_token_123
TEST_ACTOR='invalid_actor'
run_failure get ghcr.io 'credential helper: credentials are invalid'
TEST_ACTOR=github-actions
TEST_TOKEN='invalid-token'
run_failure get ghcr.io 'credential helper: credentials are invalid'
TEST_TOKEN=fake_token_123
run_failure get registry.example 'credential helper: server is not allowed'
run_failure store '{}' 'credential helper: persistence is disabled'
run_failure erase ghcr.io 'credential helper: persistence is disabled'
run_failure unknown ghcr.io 'credential helper: action is not supported'

list_output=$(HOME="$tmp/home" DOCKER_CONFIG="$tmp/docker" "$helper" list)
[[ $list_output == '{}' ]]
[[ -z $(find "$tmp/home" "$tmp/docker" -mindepth 1 -print -quit) ]]
