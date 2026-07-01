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
  contract|chat-stream|action-basic|capability-selection|context-projection|tool-card)
    ;;
  real-model-chat)
    not_implemented "Phase 4.3"
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
    echo "supported scenarios: contract, capability-selection, context-projection, chat-stream, action-basic, tool-card, real-model-chat, run-lifecycle, clarification, action-consistency, replay, budget, fobrain-poc, fobrain-readonly, fobrain-clarification, fobrain-write-approval, fobrain-live-read, fobrain-live-write" >&2
    exit 2
    ;;
esac

cd "$(dirname "$0")/.."

port="$(python3 - <<'PY'
import socket
with socket.socket() as s:
    s.bind(("127.0.0.1", 0))
    print(s.getsockname()[1])
PY
)"
addr="127.0.0.1:${port}"
base_url="http://${addr}"
tmp_dir="$(mktemp -d)"
server_log="${tmp_dir}/server.log"
config_file="${tmp_dir}/eino-workbench.yaml"

cleanup() {
  if [[ -n "${server_pid:-}" ]]; then
    kill "${server_pid}" >/dev/null 2>&1 || true
    wait "${server_pid}" >/dev/null 2>&1 || true
  fi
  rm -rf "${tmp_dir}"
}
trap cleanup EXIT

cat >"${config_file}" <<YAML
server:
  addr: "${addr}"
  read_timeout: "10s"
  write_timeout: "30s"
database:
  driver: "sqlite"
  dsn: "${tmp_dir}/eino-workbench.db"
llm:
  provider: "mock"
  base_url: "http://127.0.0.1/mock-llm"
  model: "mock-chat"
  timeout: "30s"
security:
  redact_secrets: true
  allow_private_network: false
observability:
  log_level: "debug"
  enable_request_log: true
budgets:
  default_timeout: "60s"
  max_tool_timeout: "120s"
capabilities:
  - id: "cap.smoke.read"
    provider_id: "phase3-smoke"
    tool_name: "phase3_smoke_read"
    display_name: "Phase 3 read smoke"
    description: "Smoke-only read capability registered from temporary config"
    result_schema: "tool.structured_result.v1"
    risk_level: "read_only"
    approval_required: false
    timeout: "5s"
  - id: "cap.smoke.write"
    provider_id: "phase3-smoke"
    tool_name: "phase3_smoke_write"
    display_name: "Phase 3 write smoke"
    description: "Smoke-only write capability registered from temporary config"
    result_schema: "tool.structured_result.v1"
    risk_level: "write"
    approval_required: true
    timeout: "5s"
YAML

go run ./cmd/eino-workbench --config "${config_file}" >"${server_log}" 2>&1 &
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
capability_json="${tmp_dir}/capability.json"
capability_write_json="${tmp_dir}/capability-write.json"
capability_error_json="${tmp_dir}/capability-error.json"
message_json="${tmp_dir}/message.json"
snapshot_json="${tmp_dir}/snapshot.json"
stream_headers="${tmp_dir}/stream.headers"
stream_body="${tmp_dir}/stream.body"

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
  -d '{"schema_version":"eino_action_request.v1","action_id":"action-smoke","client_request_id":"client-smoke-action","input":{"text":"hello"}}' \
  "${base_url}/api/workspaces/ws_smoke/agent/actions" \
  -o "${action_json}"
python3 - "${action_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["workspace_id"] == "ws_smoke"
assert body["status"] in ("accepted", "completed")
assert isinstance(body["result_cards"], list)
assert isinstance(body["audit_refs"], list)
PY

contract_run_id="$(python3 - "${action_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
print(body["run_id"])
PY
)"
curl -fsS -D "${stream_headers}" "${base_url}/api/workspaces/ws_smoke/runs/${contract_run_id}/stream" >/dev/null
grep -qi '^Content-Type: text/event-stream' "${stream_headers}"

if [[ "${scenario}" == "chat-stream" ]]; then
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_workbench_message_request.v1","message":"hello smoke","client_request_id":"client-smoke-message"}' \
    "${base_url}/api/workspaces/ws_smoke/messages" \
    -o "${message_json}"
  run_id="$(python3 - "${message_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_workbench_message_response.v1"
assert body["status"] == "accepted"
assert body["run_id"]
print(body["run_id"])
PY
)"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${run_id}" -o "${snapshot_json}"
python3 - "${snapshot_json}" "${run_id}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_workbench_view.v1"
assert body["run_id"] == sys.argv[2]
assert any(item["kind"] == "user_message" and item["content"] == "hello smoke" for item in body["timeline"])
assert any(item["kind"] == "assistant_message" and item["content"] == "已收到请求。" for item in body["timeline"])
PY
  curl -fsS -D "${stream_headers}" "${base_url}/api/workspaces/ws_smoke/runs/${run_id}/stream" -o "${stream_body}"
  grep -qi '^Content-Type: text/event-stream' "${stream_headers}"
  grep -q "event: message.updated" "${stream_body}"
  grep -q "hello smoke" "${stream_body}"
  grep -q "已收到请求。" "${stream_body}"
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    count = conn.execute("select count(*) from context_snapshots where run_id = ?", (run_id,)).fetchone()[0]
assert count >= 1
PY
  echo "chat-stream smoke passed"
  exit 0
fi

if [[ "${scenario}" == "action-basic" ]]; then
  run_id="$(python3 - "${action_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["run_id"]
assert body["status"] == "completed"
assert body["final_answer"] == "已收到请求。"
assert body["snapshot_url"].endswith("/runs/" + body["run_id"])
assert body["events_url"].endswith("/runs/" + body["run_id"] + "/stream")
print(body["run_id"])
PY
)"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${run_id}" -o "${snapshot_json}"
  python3 - "${snapshot_json}" "${run_id}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_workbench_view.v1"
assert body["run_id"] == sys.argv[2]
assert any(item["kind"] == "user_message" and item["content"] == "hello" for item in body["timeline"])
assert any(item["kind"] == "assistant_message" and item["content"] == "已收到请求。" for item in body["timeline"])
assert body["inspector"]["runtime"]["status"] == "succeeded"
PY
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    count = conn.execute("select count(*) from context_snapshots where run_id = ?", (run_id,)).fetchone()[0]
assert count >= 1
PY
  echo "action-basic smoke passed"
  exit 0
fi

if [[ "${scenario}" == "capability-selection" ]]; then
  python3 - "${action_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["status"] == "completed"
assert body["final_answer"] == "已收到请求。"
PY
  python3 - "${tmp_dir}/eino-workbench.db" "${contract_run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    rows = conn.execute("select safe_summary from audit_events where run_id = ?", (run_id,)).fetchall()
assert not any("capability selected:" in summary for (summary,) in rows), rows
PY
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-capability-smoke","client_request_id":"client-smoke-capability","capability_hint":"cap.smoke.read","input":{"text":"capability smoke"}}' \
    "${base_url}/api/workspaces/ws_smoke/agent/actions" \
    -o "${capability_json}"
  run_id="$(python3 - "${capability_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["status"] == "completed"
assert len(body["result_cards"]) == 1, body
assert body["result_cards"][0]["safe_summary"] == "Phase 3 read smoke 完成", body
print(body["run_id"])
PY
)"
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    rows = conn.execute(
        "select event_type, safe_summary from audit_events where run_id = ? order by created_at",
        (run_id,),
    ).fetchall()
    tool_count = conn.execute("select count(*) from tool_calls where run_id = ?", (run_id,)).fetchone()[0]
    result_count = conn.execute(
        "select count(*) from tool_results tr join tool_calls tc on tr.tool_call_id = tc.tool_call_id where tc.run_id = ?",
        (run_id,),
    ).fetchone()[0]
assert any(event_type == "tool" and "capability selected: cap.smoke.read" in summary and "allowed" in summary for event_type, summary in rows), rows
assert tool_count == 1, tool_count
assert result_count >= 1, result_count
PY
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-capability-write","client_request_id":"client-smoke-capability-write","capability_hint":"cap.smoke.write","input":{"text":"write smoke"}}' \
    "${base_url}/api/workspaces/ws_smoke/agent/actions" \
    -o "${capability_write_json}"
  write_run_id="$(python3 - "${capability_write_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["status"] == "accepted"
assert body.get("final_answer", "") == ""
print(body["run_id"])
PY
)"
  python3 - "${tmp_dir}/eino-workbench.db" "${write_run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    rows = conn.execute(
        "select event_type, safe_summary from audit_events where run_id = ? order by created_at",
        (run_id,),
    ).fetchall()
    context_count = conn.execute("select count(*) from context_snapshots where run_id = ?", (run_id,)).fetchone()[0]
assert any(event_type == "tool" and "capability selected: cap.smoke.write" in summary and "write_requires_approval" in summary for event_type, summary in rows), rows
assert context_count == 0, context_count
PY
  run_count_before_missing="$(python3 - "${tmp_dir}/eino-workbench.db" <<'PY'
import sqlite3, sys
with sqlite3.connect(sys.argv[1]) as conn:
    print(conn.execute("select count(*) from runs").fetchone()[0])
PY
)"
  status_code="$(curl -sS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-capability-missing","client_request_id":"client-smoke-capability-missing","capability_hint":"cap.missing","input":{"text":"capability smoke"}}' \
    -o "${capability_error_json}" \
    -w "%{http_code}" \
    "${base_url}/api/workspaces/ws_smoke/agent/actions")"
  [[ "${status_code}" == "400" ]]
  python3 - "${capability_error_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_error_envelope.v1"
assert body["error"]["code"] == "capability_not_registered"
PY
  run_count_after_missing="$(python3 - "${tmp_dir}/eino-workbench.db" <<'PY'
import sqlite3, sys
with sqlite3.connect(sys.argv[1]) as conn:
    print(conn.execute("select count(*) from runs").fetchone()[0])
PY
)"
  [[ "${run_count_after_missing}" == "${run_count_before_missing}" ]]
  echo "capability-selection smoke passed"
  exit 0
fi

if [[ "${scenario}" == "tool-card" ]]; then
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-tool-card","client_request_id":"client-smoke-tool-card","capability_hint":"cap.smoke.read","input":{"text":"tool card smoke"}}' \
    "${base_url}/api/workspaces/ws_smoke/agent/actions" \
    -o "${capability_json}"
  run_id="$(python3 - "${capability_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["status"] == "completed", body
assert len(body["result_cards"]) == 1, body
card = body["result_cards"][0]
assert card["tool_call_id"], card
assert card["title"] == "Phase 3 read smoke", card
assert card["safe_summary"] == "Phase 3 read smoke 完成", card
assert card["structured_result"]["schema_version"] == "tool.structured_result.v1", card
print(body["run_id"])
PY
)"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${run_id}" -o "${snapshot_json}"
  python3 - "${snapshot_json}" "${run_id}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_workbench_view.v1"
assert body["run_id"] == sys.argv[2]
tool_cards = [item for item in body["timeline"] if item["kind"] == "tool_card"]
assert len(tool_cards) == 1, body["timeline"]
card = tool_cards[0]
assert card["status"] == "succeeded", card
assert card["safe_summary"] == "Phase 3 read smoke 完成", card
assert body["inspector"]["structured"]["structured_result"]["schema_version"] == "tool.structured_result.v1"
PY
  curl -fsS -D "${stream_headers}" "${base_url}/api/workspaces/ws_smoke/runs/${run_id}/stream" -o "${stream_body}"
  grep -qi '^Content-Type: text/event-stream' "${stream_headers}"
  grep -q "event: tool.updated" "${stream_body}"
  grep -q "Phase 3 read smoke 完成" "${stream_body}"
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    tool_count = conn.execute("select count(*) from tool_calls where run_id = ?", (run_id,)).fetchone()[0]
    result_rows = conn.execute(
        "select tr.structured_schema_version, tr.safe_summary from tool_results tr join tool_calls tc on tr.tool_call_id = tc.tool_call_id where tc.run_id = ?",
        (run_id,),
    ).fetchall()
    audit_rows = conn.execute("select event_type, safe_summary from audit_events where run_id = ?", (run_id,)).fetchall()
assert tool_count == 1, tool_count
assert any(schema == "tool.structured_result.v1" and summary == "Phase 3 read smoke 完成" for schema, summary in result_rows), result_rows
assert any(event_type == "tool" and "tool completed: cap.smoke.read" in summary for event_type, summary in audit_rows), audit_rows
PY
  echo "tool-card smoke passed"
  exit 0
fi

if [[ "${scenario}" == "context-projection" ]]; then
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_workbench_message_request.v1","message":"context projection smoke","client_request_id":"client-smoke-context"}' \
    "${base_url}/api/workspaces/ws_smoke/messages" \
    -o "${message_json}"
  run_id="$(python3 - "${message_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_workbench_message_response.v1"
assert body["status"] == "accepted"
print(body["run_id"])
PY
)"
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    rows = conn.execute(
        "select safe_summary from context_snapshots where run_id = ? order by created_at",
        (run_id,),
    ).fetchall()
assert rows, "missing context snapshot"
joined = "\n".join(summary for (summary,) in rows)
assert "safe context messages=1" in joined, joined
for forbidden in ("Authorization", "credential", "raw", "provider payload", "token", "checkpoint", "interrupt"):
    assert forbidden not in joined, joined
PY
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${run_id}" -o "${snapshot_json}"
  python3 - "${snapshot_json}" "${run_id}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_workbench_view.v1"
assert body["run_id"] == sys.argv[2]
assert body["inspector"]["runtime"]["status"] == "succeeded"
assert any(item["kind"] == "assistant_message" and item["content"] == "已收到请求。" for item in body["timeline"])
PY
  echo "context-projection smoke passed"
  exit 0
fi

echo "contract smoke passed"
