#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root="$(cd "${script_dir}/.." && pwd)"
temp_dir="$(mktemp -d)"
config_path="${temp_dir}/eino-workbench.m1.yaml"
fixture_case="${M1_FIXTURE_CASE:-resolved}"

case "${fixture_case}" in
  resolved|empty|failed) ;;
  *) echo "unsupported M1 fixture case" >&2; exit 2 ;;
esac

cleanup() {
  if [[ -n "${backend_pid:-}" ]]; then
    kill "${backend_pid}" 2>/dev/null || true
  fi
  rm -rf "${temp_dir}"
}
trap cleanup EXIT INT TERM

sed -e "s#data/eino-workbench-m1.db#${temp_dir}/facts.db#" -e "s#case: \"resolved\"#case: \"${fixture_case}\"#" "${root}/configs/eino-workbench.m1.example.yaml" > "${config_path}"
(cd "${root}" && go run ./cmd/eino-workbench --config "${config_path}") &
backend_pid=$!

for _ in {1..100}; do
  if curl -fsS "http://127.0.0.1:8080/readyz" >/dev/null 2>&1; then
    break
  fi
  if ! kill -0 "${backend_pid}" 2>/dev/null; then
    wait "${backend_pid}"
  fi
  sleep 0.1
done

curl -fsS "http://127.0.0.1:8080/readyz" >/dev/null
npm --prefix "${root}/web/eino-workbench" run dev
