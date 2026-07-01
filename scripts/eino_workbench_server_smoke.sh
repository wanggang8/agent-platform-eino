#!/usr/bin/env bash
set -euo pipefail

scenario=""
provided_config=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --scenario)
      scenario="${2:-}"
      shift 2
      ;;
    --config)
      provided_config="${2:-}"
      shift 2
      ;;
    *)
      if [[ -z "${scenario}" ]]; then
        scenario="$1"
        shift
      else
        echo "unsupported argument: $1" >&2
        exit 2
      fi
      ;;
  esac
done

not_implemented() {
  local phase="$1"
  echo "scenario ${scenario} is intentionally unavailable until ${phase}; this is not a passing signal." >&2
  exit 2
}

case "${scenario}" in
  contract|chat-stream|action-basic|capability-selection|context-projection|tool-card|real-model-chat)
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

write_skip_report() {
  local reason="$1"
  mkdir -p test-results
  python3 - "${reason}" "${provided_config:-configs/eino-workbench.local.yaml}" <<'PY'
import datetime, json, sys
reason, config_path = sys.argv[1], sys.argv[2]
report = {
    "schema_version": "eino.skip_report.v1",
    "command": f"bash scripts/eino_workbench_server_smoke.sh --scenario real-model-chat --config {config_path}",
    "missing_env": ["local_llm_config_or_api_key"],
    "credential_scope": "llm",
    "reason": reason,
    "rerun_condition": "Create ignored configs/eino-workbench.local.yaml with llm.provider=openai_compatible and llm.api_key.",
    "blocks_claims": ["real-model-chat", "P1 real model provider"],
    "expires_at": (datetime.datetime.now(datetime.timezone.utc) + datetime.timedelta(days=7)).isoformat().replace("+00:00", "Z"),
}
with open("test-results/eino-workbench-skip-report.json", "w", encoding="utf-8") as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
PY
  echo "real-model-chat smoke skipped: ${reason}"
}

if [[ "${scenario}" == "real-model-chat" ]]; then
  config_for_real="${provided_config:-configs/eino-workbench.local.yaml}"
  if [[ ! -f "${config_for_real}" ]]; then
    write_skip_report "本地 LLM 配置文件不存在"
    exit 0
  fi
  if ! go run ./scripts/eino_workbench_config_prepare.go --check-llm-api-key --source "${config_for_real}" >/dev/null 2>&1; then
    write_skip_report "本地 LLM api_key 未配置"
    exit 0
  fi
fi

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
real_model_source_config="${provided_config:-configs/eino-workbench.local.yaml}"
real_model_summary_json="${tmp_dir}/real-model-summary.json"

cleanup() {
  if [[ -n "${server_pid:-}" ]]; then
    kill "${server_pid}" >/dev/null 2>&1 || true
    wait "${server_pid}" >/dev/null 2>&1 || true
  fi
  rm -rf "${tmp_dir}"
}
trap cleanup EXIT

if [[ "${scenario}" != "real-model-chat" ]]; then
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
fi

if [[ "${scenario}" == "real-model-chat" ]]; then
  go run ./scripts/eino_workbench_config_prepare.go \
    --source "${real_model_source_config}" \
    --target "${config_file}" \
    --addr "${addr}" \
    --dsn "${tmp_dir}/eino-workbench.db" \
    --summary "${real_model_summary_json}"
  base_url="http://${addr}"
fi

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

action_status="200"
if [[ "${scenario}" == "real-model-chat" ]]; then
  action_status="$(curl -sS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-smoke","client_request_id":"client-smoke-action","input":{"text":"hello"}}' \
    "${base_url}/api/workspaces/ws_smoke/agent/actions" \
    -o "${action_json}" \
    -w "%{http_code}")"
else
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-smoke","client_request_id":"client-smoke-action","input":{"text":"hello"}}' \
    "${base_url}/api/workspaces/ws_smoke/agent/actions" \
    -o "${action_json}"
fi
if [[ "${scenario}" != "real-model-chat" ]]; then
python3 - "${action_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["workspace_id"] == "ws_smoke"
assert body["status"] in ("accepted", "completed")
assert isinstance(body["result_cards"], list)
assert isinstance(body["audit_refs"], list)
PY
fi

if [[ "${scenario}" != "real-model-chat" ]]; then
contract_run_id="$(python3 - "${action_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
print(body["run_id"])
PY
)"
curl -fsS -D "${stream_headers}" "${base_url}/api/workspaces/ws_smoke/runs/${contract_run_id}/stream" >/dev/null
grep -qi '^Content-Type: text/event-stream' "${stream_headers}"
fi

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

if [[ "${scenario}" == "real-model-chat" ]]; then
  mkdir -p test-results
  run_id="$(python3 - "${action_json}" "${real_model_summary_json}" "${action_status}" <<'PY'
import datetime, json, sys

action_path, summary_path, action_status = sys.argv[1:]
summary = json.load(open(summary_path, encoding="utf-8"))
try:
    body = json.load(open(action_path, encoding="utf-8"))
except Exception:
    body = {}

def write_report(status, provider_error_category, failure_category, assistant_final_answer="", run_id=""):
    report = {
        "schema_version": "eino.real_model_provider_report.v1",
        "scenario": "real-model-chat",
        "status": status,
        "prompt_id": "real-model-smoke-basic",
        "prompt": "hello",
        "provider_kind": summary["provider"],
        "model_label": summary["model"],
        "redacted_config_summary": summary,
        "assistant_final_answer": assistant_final_answer,
        "provider_error_category": provider_error_category,
        "redaction_checks": [
            "authorization_not_present",
            "api_key_not_present",
            "token_not_present",
            "raw_prompt_not_present",
            "raw_provider_body_not_present"
        ],
        "run_id": run_id,
        "report_created_at": datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00", "Z"),
        "failure_category": failure_category
    }
    with open("test-results/eino-workbench-real-model-provider-report.json", "w", encoding="utf-8") as f:
        json.dump(report, f, ensure_ascii=False, indent=2)

encoded = json.dumps(body, ensure_ascii=False).lower()
for forbidden in ("authorization", "bearer ", "api_key", "raw provider", "raw body", "token"):
    if forbidden in encoded:
        write_report("failed", "redaction_failed", "redaction_failed")
        sys.exit(1)

if action_status != "200":
    code = body.get("error", {}).get("code") or f"http_{action_status}"
    write_report("failed", code, "provider_error")
    sys.exit(1)

if body.get("schema_version") != "eino_action_result.v1" or body.get("status") != "completed" or not body.get("final_answer", "").strip():
    code = body.get("error", {}).get("code", "unexpected_action_result")
    write_report("failed", code, "unexpected_action_result", body.get("final_answer", ""), body.get("run_id", ""))
    sys.exit(1)

print(body["run_id"])
PY
  )" || {
    echo "real-model-chat smoke failed; report written to test-results/eino-workbench-real-model-provider-report.json"
    exit 1
  }
  if ! curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${run_id}" -o "${snapshot_json}"; then
    python3 - "${action_json}" "${real_model_summary_json}" "${run_id}" <<'PY'
import datetime, json, sys
action_path, summary_path, run_id = sys.argv[1:]
body = json.load(open(action_path, encoding="utf-8"))
summary = json.load(open(summary_path, encoding="utf-8"))
# snapshot 失败也必须写 failed report，避免真实模型验收只留下 shell 退出码。
report = {
    "schema_version": "eino.real_model_provider_report.v1",
    "scenario": "real-model-chat",
    "status": "failed",
    "prompt_id": "real-model-smoke-basic",
    "prompt": "hello",
    "provider_kind": summary["provider"],
    "model_label": summary["model"],
    "redacted_config_summary": summary,
    "assistant_final_answer": body.get("final_answer", ""),
    "provider_error_category": "snapshot_unavailable",
    "redaction_checks": ["authorization_not_present", "api_key_not_present", "token_not_present", "raw_prompt_not_present", "raw_provider_body_not_present"],
    "run_id": run_id,
    "report_created_at": datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00", "Z"),
    "failure_category": "snapshot_unavailable"
}
with open("test-results/eino-workbench-real-model-provider-report.json", "w", encoding="utf-8") as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
PY
    echo "real-model-chat smoke failed; report written to test-results/eino-workbench-real-model-provider-report.json"
    exit 1
  fi
  python3 - "${action_json}" "${snapshot_json}" "${real_model_summary_json}" "${run_id}" <<'PY'
import datetime, json, sys

action_path, snapshot_path, summary_path, run_id = sys.argv[1:]
body = json.load(open(action_path, encoding="utf-8"))
snapshot = json.load(open(snapshot_path, encoding="utf-8"))
summary = json.load(open(summary_path, encoding="utf-8"))
failure_category = "none"
provider_error_category = "none"
status = "passed"
encoded = json.dumps({"action": body, "snapshot": snapshot}, ensure_ascii=False).lower()
for forbidden in ("authorization", "bearer ", "api_key", "raw provider", "raw body", "token"):
    if forbidden in encoded:
        status = "failed"
        failure_category = "redaction_failed"
        provider_error_category = "redaction_failed"
if snapshot.get("schema_version") != "eino_workbench_view.v1" or snapshot.get("run_id") != run_id:
    status = "failed"
    failure_category = "snapshot_mismatch"
    provider_error_category = "snapshot_mismatch"
if not any(item.get("kind") == "assistant_message" and item.get("content", "").strip() for item in snapshot.get("timeline", [])):
    status = "failed"
    failure_category = "missing_assistant_message"
    provider_error_category = "missing_assistant_message"
report = {
    "schema_version": "eino.real_model_provider_report.v1",
    "scenario": "real-model-chat",
    "status": status,
    "prompt_id": "real-model-smoke-basic",
    "prompt": "hello",
    "provider_kind": summary["provider"],
    "model_label": summary["model"],
    "redacted_config_summary": summary,
    "assistant_final_answer": body.get("final_answer", ""),
    "provider_error_category": provider_error_category,
    "redaction_checks": [
        "authorization_not_present",
        "api_key_not_present",
        "token_not_present",
        "raw_prompt_not_present",
        "raw_provider_body_not_present"
    ],
    "run_id": run_id,
    "report_created_at": datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00", "Z"),
    "failure_category": failure_category
}
with open("test-results/eino-workbench-real-model-provider-report.json", "w", encoding="utf-8") as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
if status != "passed":
    sys.exit(1)
PY
  echo "real-model-chat smoke passed"
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
