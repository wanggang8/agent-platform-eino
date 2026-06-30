#!/usr/bin/env bash
set -euo pipefail

scenario="${1:-}"
if [[ "$scenario" == "--scenario" ]]; then
  scenario="${2:-}"
fi

not_implemented() {
  local phase="$1"
  echo "scenario ${scenario} is intentionally unavailable until ${phase}; this is not a passing signal." >&2
  exit 2
}

case "${scenario}" in
  contract)
    ;;
  capability-selection|context-projection)
    not_implemented "Phase 3"
    ;;
  chat-stream|action-basic)
    not_implemented "Phase 3"
    ;;
  tool-card)
    not_implemented "Phase 4"
    ;;
  run-lifecycle|clarification)
    not_implemented "Phase 6"
    ;;
  action-consistency|replay|budget)
    not_implemented "Phase 7"
    ;;
  fobrain-poc)
    not_implemented "Phase 5"
    ;;
  fobrain-readonly|fobrain-clarification|fobrain-write-approval|fobrain-live-read|fobrain-live-write)
    not_implemented "Phase 8"
    ;;
  "")
    echo "missing --scenario value" >&2
    exit 2
    ;;
  *)
    echo "unsupported scenario: ${scenario}" >&2
    echo "supported scenarios: contract, capability-selection, context-projection, chat-stream, action-basic, tool-card, run-lifecycle, clarification, action-consistency, replay, budget, fobrain-poc, fobrain-readonly, fobrain-clarification, fobrain-write-approval, fobrain-live-read, fobrain-live-write" >&2
    exit 2
    ;;
esac

cd "$(dirname "$0")/.."

port="${EINO_WORKBENCH_SMOKE_PORT:-18081}"
addr="127.0.0.1:${port}"
base_url="http://${addr}"
tmp_dir="$(mktemp -d)"
server_log="${tmp_dir}/server.log"

cleanup() {
  if [[ -n "${server_pid:-}" ]]; then
    kill "${server_pid}" >/dev/null 2>&1 || true
    wait "${server_pid}" >/dev/null 2>&1 || true
  fi
  rm -rf "${tmp_dir}"
}
trap cleanup EXIT

EINO_WORKBENCH_ADDR="${addr}" go run ./cmd/eino-workbench >"${server_log}" 2>&1 &
server_pid="$!"

for _ in {1..50}; do
  if curl -fsS "${base_url}/healthz" >/dev/null 2>&1; then
    break
  fi
  sleep 0.1
done

curl -fsS "${base_url}/healthz" >/dev/null

view_json="${tmp_dir}/view.json"
action_json="${tmp_dir}/action.json"
stream_headers="${tmp_dir}/stream.headers"

curl -fsS "${base_url}/api/workspaces/ws_smoke/views/current" -o "${view_json}"
python3 - "${view_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_workbench_view.v1"
assert body["workspace_id"] == "ws_smoke"
assert body["status"] == "created"
assert isinstance(body["timeline"], list)
assert isinstance(body["inspector"], dict)
PY

curl -fsS \
  -H "Content-Type: application/json" \
  -d '{}' \
  "${base_url}/api/workspaces/ws_smoke/agent/actions" \
  -o "${action_json}"
python3 - "${action_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["workspace_id"] == "ws_smoke"
assert body["status"] == "accepted"
assert isinstance(body["result_cards"], list)
assert isinstance(body["audit_refs"], list)
PY

curl -fsS -D "${stream_headers}" "${base_url}/api/workspaces/ws_smoke/runs/run_smoke/stream" >/dev/null
grep -qi '^Content-Type: text/event-stream' "${stream_headers}"

echo "contract smoke passed"
