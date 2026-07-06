#!/usr/bin/env bash
set -euo pipefail

scenario=""
provided_config=""
batch_d_owner=""
batch_d_department=""
batch_d_ip=""
batch_c_keyword=""
batch_c_business_name=""
batch_c_severity=""
batch_c_status=""
batch_c_person=""
batch_c_field=""
batch_c_time_range=""
batch_e_asset_id=""
batch_e_asset_network_type=""
batch_e_vulnerability_id=""
batch_e_business_name=""
batch_e_vulnerability_name=""

require_option_value() {
  local option="$1"
  local value="${2:-}"
  if [[ -z "${value}" || "${value}" == --* ]]; then
    echo "missing value for ${option}" >&2
    exit 2
  fi
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --scenario)
      require_option_value "$1" "${2:-}"
      scenario="${2:-}"
      shift 2
      ;;
    --config)
      require_option_value "$1" "${2:-}"
      provided_config="${2:-}"
      shift 2
      ;;
    --owner)
      require_option_value "$1" "${2:-}"
      batch_d_owner="${2:-}"
      shift 2
      ;;
    --department)
      require_option_value "$1" "${2:-}"
      batch_d_department="${2:-}"
      shift 2
      ;;
    --ip)
      require_option_value "$1" "${2:-}"
      batch_d_ip="${2:-}"
      shift 2
      ;;
    --keyword)
      require_option_value "$1" "${2:-}"
      batch_c_keyword="${2:-}"
      shift 2
      ;;
    --severity)
      require_option_value "$1" "${2:-}"
      batch_c_severity="${2:-}"
      shift 2
      ;;
    --status)
      require_option_value "$1" "${2:-}"
      batch_c_status="${2:-}"
      shift 2
      ;;
    --person)
      require_option_value "$1" "${2:-}"
      batch_c_person="${2:-}"
      shift 2
      ;;
    --field)
      require_option_value "$1" "${2:-}"
      batch_c_field="${2:-}"
      shift 2
      ;;
    --time-range)
      require_option_value "$1" "${2:-}"
      batch_c_time_range="${2:-}"
      shift 2
      ;;
    --asset-id)
      require_option_value "$1" "${2:-}"
      batch_e_asset_id="${2:-}"
      shift 2
      ;;
    --asset-network-type)
      require_option_value "$1" "${2:-}"
      batch_e_asset_network_type="${2:-}"
      shift 2
      ;;
    --vulnerability-id)
      require_option_value "$1" "${2:-}"
      batch_e_vulnerability_id="${2:-}"
      shift 2
      ;;
    --business-name)
      require_option_value "$1" "${2:-}"
      batch_c_business_name="${2:-}"
      batch_e_business_name="${2:-}"
      shift 2
      ;;
    --vulnerability-name)
      require_option_value "$1" "${2:-}"
      batch_e_vulnerability_name="${2:-}"
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
  contract|chat-stream|action-basic|capability-selection|context-projection|tool-card|run-lifecycle|action-consistency|replay|mcp-mock|real-model-chat|fobrain-poc|fobrain-batch-a|fobrain-batch-c|fobrain-batch-d|fobrain-batch-e|fobrain-clarification|clarification|budget)
    ;;
  fobrain-readonly|fobrain-write-approval|fobrain-live-read|fobrain-live-write)
    not_implemented "Phase 8"
    ;;
  "")
    echo "missing --scenario value" >&2
    exit 2
    ;;
  *)
    echo "unsupported scenario: ${scenario}" >&2
    echo "supported scenarios: contract, capability-selection, context-projection, chat-stream, action-basic, tool-card, mcp-mock, real-model-chat, run-lifecycle, clarification, action-consistency, replay, budget, fobrain-poc, fobrain-batch-a, fobrain-batch-c, fobrain-batch-d, fobrain-batch-e, fobrain-readonly, fobrain-clarification, fobrain-write-approval, fobrain-live-read, fobrain-live-write" >&2
    exit 2
    ;;
esac

cd "$(dirname "$0")/.."

if [[ "${scenario}" == "fobrain-clarification" ]]; then
  # Phase 8.4 先验证实体消歧和 clarification 同源投影，不声明 24 个真实只读工具已恢复。
  go test ./internal/einoapp/providers/fobrain -run Disambiguation -count=1
  go test ./internal/einoapp/execution -run Clarification -count=1
  go test ./internal/einoapp/product -run ClarificationCandidates -count=1
  npm run eino-workbench:contract-test
  echo "fobrain-clarification smoke passed"
  exit 0
fi

if [[ "${scenario}" == "clarification" ]]; then
  # Phase 6.4 验证 HITL clarification/pending 的后端事实、SSE patch 和前端 reducer，不声明真实 Fobrain 多候选工具已完成。
  go test ./internal/einoapp/execution -run 'NewClarificationPendingEvent|ClarificationSubmit|ClarificationCancel|ClarificationDuplicate|ClarificationAfterRestart' -count=1
  go test ./internal/einoapp/product -run 'ClarificationCandidates|TerminalPending' -count=1
  npm run eino-workbench:stream-test -- --grep pending
  npm run eino-workbench:contract-test
  echo "clarification smoke passed"
  exit 0
fi

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

write_fobrain_skip_report() {
  local reason="$1"
  local category="${2:-missing_live_prerequisite}"
  mkdir -p test-results
  python3 - "${reason}" "${provided_config:-configs/eino-workbench.local.yaml}" "${category}" <<'PY'
import datetime, json, sys
reason, config_path, category = sys.argv[1], sys.argv[2], sys.argv[3]
missing_env = []
rerun_condition = "Implement and verify the later Fobrain live HTTP client phase."
if category != "phase_unsupported":
    missing_env = ["local_llm_or_fobrain_live_credential"]
    rerun_condition = "Create ignored configs/eino-workbench.local.yaml with LLM api_key and fobrain live api_token."
report = {
    "schema_version": "eino.skip_report.v1",
    "command": f"bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config {config_path}",
    "missing_env": missing_env,
    "credential_scope": "fobrain-real-model",
    "reason": reason,
    "rerun_condition": rerun_condition,
    "blocks_claims": ["fobrain-poc real-model tool selection", "Fobrain live read"],
    "expires_at": (datetime.datetime.now(datetime.timezone.utc) + datetime.timedelta(days=7)).isoformat().replace("+00:00", "Z"),
}
with open("test-results/eino-workbench-fobrain-skip-report.json", "w", encoding="utf-8") as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
PY
  echo "fobrain-poc real-model selection skipped: ${reason}"
}

write_fobrain_batch_a_skip_report() {
  local reason="$1"
  mkdir -p test-results
  python3 - "${reason}" "${provided_config:-configs/eino-workbench.local.yaml}" <<'PY'
import datetime, json, sys
reason, config_path = sys.argv[1], sys.argv[2]
report = {
    "schema_version": "eino.skip_report.v1",
    "command": f"bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-a --config {config_path}",
    "missing_env": ["local_fobrain_live_config"],
    "credential_scope": "fobrain-workspace",
    "reason": reason,
    "rerun_condition": "Create ignored configs/eino-workbench.local.yaml with fobrain.connector_status.mode=live and workspace credential.",
    "blocks_claims": ["fobrain-batch-a live pass", "Fobrain Batch A final acceptance"],
    "expires_at": (datetime.datetime.now(datetime.timezone.utc) + datetime.timedelta(days=7)).isoformat().replace("+00:00", "Z"),
}
with open("test-results/eino-workbench-fobrain-batch-a-skip-report.json", "w", encoding="utf-8") as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
PY
  echo "fobrain-batch-a live smoke skipped: ${reason}"
}

write_fobrain_batch_d_skip_report() {
  local reason="$1"
  mkdir -p test-results
  python3 - "${reason}" "${provided_config:-configs/eino-workbench.local.yaml}" <<'PY'
import datetime, json, sys
reason, config_path = sys.argv[1], sys.argv[2]
report = {
    "schema_version": "eino.skip_report.v1",
    "command": f"bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-d --config {config_path}",
    "missing_env": ["local_fobrain_live_config_or_batch_d_samples"],
    "credential_scope": "fobrain-workspace",
    "reason": reason,
    "rerun_condition": "Create ignored configs/eino-workbench.local.yaml with Fobrain live credential and pass --owner/--department/--ip when required.",
    "blocks_claims": ["fobrain-batch-d live pass", "Fobrain 24 readonly final acceptance"],
    "expires_at": (datetime.datetime.now(datetime.timezone.utc) + datetime.timedelta(days=7)).isoformat().replace("+00:00", "Z"),
}
with open("test-results/eino-workbench-fobrain-batch-d-skip-report.json", "w", encoding="utf-8") as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
PY
  echo "fobrain-batch-d live smoke skipped: ${reason}"
}

write_fobrain_batch_c_skip_report() {
  local reason="$1"
  mkdir -p test-results
  python3 - "${reason}" "${provided_config:-configs/eino-workbench.local.yaml}" <<'PY'
import datetime, json, sys
reason, config_path = sys.argv[1], sys.argv[2]
report = {
    "schema_version": "eino.skip_report.v1",
    "command": f"bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-c --config {config_path}",
    "missing_env": ["local_fobrain_live_config"],
    "credential_scope": "fobrain-workspace",
    "reason": reason,
    "rerun_condition": "Create ignored configs/eino-workbench.local.yaml with Fobrain live credential; optional filters can be passed with --keyword/--business-name/--severity/--status/--person/--field/--time-range.",
    "blocks_claims": ["fobrain-batch-c live pass", "Fobrain 24 readonly final acceptance"],
    "expires_at": (datetime.datetime.now(datetime.timezone.utc) + datetime.timedelta(days=7)).isoformat().replace("+00:00", "Z"),
}
with open("test-results/eino-workbench-fobrain-batch-c-skip-report.json", "w", encoding="utf-8") as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
PY
  echo "fobrain-batch-c live smoke skipped: ${reason}"
}

write_fobrain_batch_e_skip_report() {
  local reason="$1"
  mkdir -p test-results
  python3 - "${reason}" "${provided_config:-configs/eino-workbench.local.yaml}" <<'PY'
import datetime, json, sys
reason, config_path = sys.argv[1], sys.argv[2]
report = {
    "schema_version": "eino.skip_report.v1",
    "command": f"bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-e --config {config_path}",
    "missing_env": ["local_fobrain_live_config_or_batch_e_samples"],
    "credential_scope": "fobrain-workspace",
    "reason": reason,
    "rerun_condition": "Create ignored configs/eino-workbench.local.yaml with Fobrain live credential and pass --asset-id/--vulnerability-id/--business-name/--vulnerability-name.",
    "blocks_claims": ["fobrain-batch-e live pass", "Fobrain 24 readonly final acceptance"],
    "expires_at": (datetime.datetime.now(datetime.timezone.utc) + datetime.timedelta(days=7)).isoformat().replace("+00:00", "Z"),
}
with open("test-results/eino-workbench-fobrain-batch-e-skip-report.json", "w", encoding="utf-8") as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
PY
  echo "fobrain-batch-e live smoke skipped: ${reason}"
}

if [[ "${scenario}" == "fobrain-batch-a" ]]; then
  config_for_fobrain="${provided_config:-configs/eino-workbench.local.yaml}"
  if [[ ! -f "${config_for_fobrain}" ]]; then
    write_fobrain_batch_a_skip_report "本地 Fobrain live 配置文件不存在"
    exit 0
  fi
  if ! go run ./scripts/eino_workbench_config_prepare.go --check-fobrain-live-credential --source "${config_for_fobrain}" >/dev/null 2>&1; then
    write_fobrain_batch_a_skip_report "本地 Fobrain live 工作区凭据未配置"
    exit 0
  fi
  report_path="test-results/eino-workbench-fobrain-batch-a-live-report.json"
  go run ./scripts/fobrain_batch_a_smoke \
    --config "${config_for_fobrain}" \
    --output "${report_path}"
  node scripts/eino_workbench_report_validate.mjs \
    --schema docs/schemas/fobrain/batch_a_live_report.v1.schema.json \
    --report "${report_path}" >/dev/null
  echo "fobrain-batch-a live smoke passed"
  exit 0
fi

if [[ "${scenario}" == "fobrain-batch-c" ]]; then
  config_for_fobrain="${provided_config:-configs/eino-workbench.local.yaml}"
  if [[ ! -f "${config_for_fobrain}" ]]; then
    write_fobrain_batch_c_skip_report "本地 Fobrain live 配置文件不存在"
    exit 0
  fi
  if ! go run ./scripts/eino_workbench_config_prepare.go --check-fobrain-live-credential --source "${config_for_fobrain}" >/dev/null 2>&1; then
    write_fobrain_batch_c_skip_report "本地 Fobrain live 工作区凭据未配置"
    exit 0
  fi
  report_path="test-results/eino-workbench-fobrain-batch-c-live-report.json"
  args=(--config "${config_for_fobrain}" --output "${report_path}")
  if [[ -n "${batch_c_keyword}" ]]; then
    args+=(--keyword "${batch_c_keyword}")
  fi
  if [[ -n "${batch_c_business_name}" ]]; then
    args+=(--business-name "${batch_c_business_name}")
  fi
  if [[ -n "${batch_c_severity}" ]]; then
    args+=(--severity "${batch_c_severity}")
  fi
  if [[ -n "${batch_c_status}" ]]; then
    args+=(--status "${batch_c_status}")
  fi
  if [[ -n "${batch_c_person}" ]]; then
    args+=(--person "${batch_c_person}")
  fi
  if [[ -n "${batch_c_field}" ]]; then
    args+=(--field "${batch_c_field}")
  fi
  if [[ -n "${batch_c_time_range}" ]]; then
    args+=(--time-range "${batch_c_time_range}")
  fi
  set +e
  go run ./scripts/fobrain_batch_c_smoke "${args[@]}"
  batch_c_exit=$?
  set -e
  if [[ ! -f "${report_path}" ]]; then
    echo "fobrain-batch-c live smoke failed before report was written" >&2
    exit "${batch_c_exit}"
  fi
  node scripts/eino_workbench_report_validate.mjs \
    --schema docs/schemas/fobrain/batch_c_live_report.v1.schema.json \
    --report "${report_path}" >/dev/null
  status="$(python3 - "${report_path}" <<'PY'
import json, sys
print(json.load(open(sys.argv[1], encoding="utf-8"))["status"])
PY
)"
  if [[ "${status}" == "passed" ]]; then
    echo "fobrain-batch-c live smoke passed"
  else
    echo "fobrain-batch-c live smoke ${status}; report blocks final claims"
  fi
  if [[ "${batch_c_exit}" -ne 0 ]]; then
    exit "${batch_c_exit}"
  fi
  exit 0
fi

if [[ "${scenario}" == "fobrain-batch-d" ]]; then
  config_for_fobrain="${provided_config:-configs/eino-workbench.local.yaml}"
  if [[ ! -f "${config_for_fobrain}" ]]; then
    write_fobrain_batch_d_skip_report "本地 Fobrain live 配置文件不存在"
    exit 0
  fi
  if ! go run ./scripts/eino_workbench_config_prepare.go --check-fobrain-live-credential --source "${config_for_fobrain}" >/dev/null 2>&1; then
    write_fobrain_batch_d_skip_report "本地 Fobrain live 工作区凭据未配置"
    exit 0
  fi
  report_path="test-results/eino-workbench-fobrain-batch-d-live-report.json"
  args=(--config "${config_for_fobrain}" --output "${report_path}")
  if [[ -n "${batch_d_owner}" ]]; then
    args+=(--owner "${batch_d_owner}")
  fi
  if [[ -n "${batch_d_department}" ]]; then
    args+=(--department "${batch_d_department}")
  fi
  if [[ -n "${batch_d_ip}" ]]; then
    args+=(--ip "${batch_d_ip}")
  fi
  set +e
  go run ./scripts/fobrain_batch_d_smoke "${args[@]}"
  batch_d_exit=$?
  set -e
  if [[ ! -f "${report_path}" ]]; then
    echo "fobrain-batch-d live smoke failed before report was written" >&2
    exit "${batch_d_exit}"
  fi
  node scripts/eino_workbench_report_validate.mjs \
    --schema docs/schemas/fobrain/batch_d_live_report.v1.schema.json \
    --report "${report_path}" >/dev/null
  status="$(python3 - "${report_path}" <<'PY'
import json, sys
print(json.load(open(sys.argv[1], encoding="utf-8"))["status"])
PY
)"
  if [[ "${status}" == "passed" ]]; then
    echo "fobrain-batch-d live smoke passed"
  else
    echo "fobrain-batch-d live smoke ${status}; report blocks final claims"
  fi
  if [[ "${batch_d_exit}" -ne 0 ]]; then
    exit "${batch_d_exit}"
  fi
  exit 0
fi

if [[ "${scenario}" == "fobrain-batch-e" ]]; then
  config_for_fobrain="${provided_config:-configs/eino-workbench.local.yaml}"
  if [[ ! -f "${config_for_fobrain}" ]]; then
    write_fobrain_batch_e_skip_report "本地 Fobrain live 配置文件不存在"
    exit 0
  fi
  if ! go run ./scripts/eino_workbench_config_prepare.go --check-fobrain-live-credential --source "${config_for_fobrain}" >/dev/null 2>&1; then
    write_fobrain_batch_e_skip_report "本地 Fobrain live 工作区凭据未配置"
    exit 0
  fi
  report_path="test-results/eino-workbench-fobrain-batch-e-live-report.json"
  args=(--config "${config_for_fobrain}" --output "${report_path}")
  if [[ -n "${batch_e_asset_id}" ]]; then
    args+=(--asset-id "${batch_e_asset_id}")
  fi
  if [[ -n "${batch_e_asset_network_type}" ]]; then
    args+=(--asset-network-type "${batch_e_asset_network_type}")
  fi
  if [[ -n "${batch_e_vulnerability_id}" ]]; then
    args+=(--vulnerability-id "${batch_e_vulnerability_id}")
  fi
  if [[ -n "${batch_e_business_name}" ]]; then
    args+=(--business-name "${batch_e_business_name}")
  fi
  if [[ -n "${batch_e_vulnerability_name}" ]]; then
    args+=(--vulnerability-name "${batch_e_vulnerability_name}")
  fi
  set +e
  go run ./scripts/fobrain_batch_e_smoke "${args[@]}"
  batch_e_exit=$?
  set -e
  if [[ ! -f "${report_path}" ]]; then
    echo "fobrain-batch-e live smoke failed before report was written" >&2
    exit "${batch_e_exit}"
  fi
  node scripts/eino_workbench_report_validate.mjs \
    --schema docs/schemas/fobrain/batch_e_live_report.v1.schema.json \
    --report "${report_path}" >/dev/null
  status="$(python3 - "${report_path}" <<'PY'
import json, sys
print(json.load(open(sys.argv[1], encoding="utf-8"))["status"])
PY
)"
  if [[ "${status}" == "passed" ]]; then
    echo "fobrain-batch-e live smoke passed"
  else
    echo "fobrain-batch-e live smoke ${status}; report blocks final claims"
  fi
  if [[ "${batch_e_exit}" -ne 0 ]]; then
    exit "${batch_e_exit}"
  fi
  exit 0
fi

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
    risk_level: "low"
    side_effect: "read_external"
    policy_ref: "policy:smoke:read:v1"
    permission_scope: "workspace"
    credential_binding_policy: "none"
    approval_required: false
    idempotency_required: false
    timeout: "5s"
  - id: "cap.smoke.write"
    provider_id: "phase3-smoke"
    tool_name: "phase3_smoke_write"
    display_name: "Phase 3 write smoke"
    description: "Smoke-only write capability registered from temporary config"
    result_schema: "tool.structured_result.v1"
    risk_level: "high"
    side_effect: "write_external"
    policy_ref: "policy:smoke:write:v1"
    permission_scope: "workspace"
    credential_binding_policy: "none"
    approval_required: true
    idempotency_required: true
    timeout: "5s"
mcp_mock_servers:
  - server_id: "mock"
    tools:
      - name: "asset_lookup"
        title: "MCP 资产查询"
        description: "Smoke-only MCP mock capability"
        input_schema:
          type: "object"
          properties:
            query: "string"
          required: ["query"]
        output_schema:
          type: "object"
        annotations:
          read_only_hint: true
    results:
      asset_lookup:
        structured_content:
          result_ref: "result:mcp:asset_lookup"
          safe_summary: "MCP 资产查询完成"
YAML
if [[ "${scenario}" == "fobrain-poc" ]]; then
cat >>"${config_file}" <<YAML
fobrain:
  enabled: true
  connector_id: "fobrain"
  workspace_id: "ws_smoke"
  base_url: "https://fobrain.example.local/api"
  timeout: "10s"
  credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
  connector_status:
    mode: "mock"
    available: true
YAML
fi
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

for _ in {1..150}; do
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
current_view_json="${tmp_dir}/current-view.json"
replay_json="${tmp_dir}/replay.json"
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

if [[ "${scenario}" == "run-lifecycle" ]]; then
  run_lifecycle_message() {
    local run_id="$1"
    local text="$2"
    local output_file="$3"
    curl -fsS \
      -H "Content-Type: application/json" \
      -d "{\"schema_version\":\"eino_workbench_message_request.v1\",\"message\":\"${text}\",\"client_request_id\":\"client-${run_id}\",\"run_id\":\"${run_id}\"}" \
      "${base_url}/api/workspaces/ws_smoke/messages" \
      -o "${output_file}"
  }

  run_lifecycle_action() {
    local run_id="$1"
    local action="$2"
    local output_file="$3"
    curl -fsS \
      -H "Content-Type: application/json" \
      -d "{\"schema_version\":\"eino_run_lifecycle_request.v1\",\"action\":\"${action}\",\"client_request_id\":\"client-${run_id}-${action}\"}" \
      "${base_url}/api/workspaces/ws_smoke/runs/${run_id}/lifecycle" \
      -o "${output_file}"
  }
  run_lifecycle_seed_status() {
    local run_id="$1"
    local status="$2"
    python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" "${status}" <<'PY'
import sqlite3, time, sys
db_path, run_id, status = sys.argv[1], sys.argv[2], sys.argv[3]
now = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
with sqlite3.connect(db_path) as conn:
    conn.execute("update runs set status = ?, safe_error = '', updated_at = ? where run_id = ?", (status, now, run_id))
PY
  }
  run_lifecycle_seed_tool() {
    local run_id="$1"
    local tool_call_id="$2"
    python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" "${tool_call_id}" <<'PY'
import sqlite3, time, sys
db_path, run_id, tool_call_id = sys.argv[1], sys.argv[2], sys.argv[3]
now = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
with sqlite3.connect(db_path) as conn:
    conn.execute(
        "insert into tool_calls(tool_call_id, run_id, tool_id, display_name, status, args_hash, args_preview, created_at, ended_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
        (tool_call_id, run_id, "cap.smoke.read", "只读查询", "running", "", "scope=smoke", now, ""),
    )
PY
  }

  run_lifecycle_message "run-life-cancel" "cancel smoke" "${message_json}"
  run_lifecycle_seed_status "run-life-cancel" "running"
  run_lifecycle_seed_tool "run-life-cancel" "call-life-cancel-1"
  run_lifecycle_action "run-life-cancel" "cancel" "${capability_json}"
  python3 - "${capability_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
payload = json.dumps(body, ensure_ascii=False).lower()
for forbidden in ("authorization", "api_key", "api_token", "provider_payload", "raw_payload", "raw provider", "raw body", "bearer "):
    assert forbidden not in payload, forbidden
assert body["schema_version"] == "eino_action_result.v1"
assert body["run_id"] == "run-life-cancel"
assert body["status"] == "cancelled", body
PY
  run_lifecycle_action "run-life-cancel" "cancel" "${capability_json}"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/run-life-cancel" -o "${snapshot_json}"
  python3 - "${snapshot_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["run_id"] == "run-life-cancel"
assert body["status"] == "cancelled", body
assert body["inspector"]["runtime"]["safe_error"] == "user_cancelled", body["inspector"]["runtime"]
tool_cards = [item for item in body["timeline"] if item["kind"] == "tool_card"]
assert len(tool_cards) == 1 and tool_cards[0]["status"] == "cancelled", tool_cards
PY

  run_lifecycle_message "run-life-stop" "stop smoke" "${message_json}"
  run_lifecycle_seed_status "run-life-stop" "running"
  run_lifecycle_seed_tool "run-life-stop" "call-life-stop-1"
  run_lifecycle_action "run-life-stop" "stop" "${capability_json}"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/run-life-stop/replay" -o "${replay_json}"
  python3 - "${capability_json}" "${replay_json}" <<'PY'
import json, sys
action = json.load(open(sys.argv[1]))
replay = json.load(open(sys.argv[2]))
assert action["status"] == "stopped", action
assert replay["view"]["status"] == "stopped", replay["view"]
assert any(event.get("event_type") == "lifecycle" and event.get("safe_summary") == "run lifecycle: stopped" for event in replay["events"] if isinstance(event, dict)), replay["events"]
tool_cards = [item for item in replay["view"]["timeline"] if item["kind"] == "tool_card"]
assert len(tool_cards) == 1 and tool_cards[0]["status"] == "cancelled", tool_cards
PY

  run_lifecycle_message "run-life-timeout" "timeout smoke" "${message_json}"
  run_lifecycle_seed_status "run-life-timeout" "running"
  run_lifecycle_seed_tool "run-life-timeout" "call-life-timeout-1"
  run_lifecycle_action "run-life-timeout" "provider_timeout" "${capability_json}"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/run-life-timeout" -o "${snapshot_json}"
  python3 - "${capability_json}" "${snapshot_json}" <<'PY'
import json, sys
action = json.load(open(sys.argv[1]))
snapshot = json.load(open(sys.argv[2]))
payload = json.dumps([action, snapshot], ensure_ascii=False).lower()
for forbidden in ("authorization", "api_key", "api_token", "provider_payload", "raw_payload", "raw provider", "raw body", "bearer "):
    assert forbidden not in payload, forbidden
assert action["status"] == "failed", action
assert snapshot["status"] == "failed", snapshot
assert snapshot["inspector"]["runtime"]["safe_error"] == "provider_timeout", snapshot["inspector"]["runtime"]
tool_cards = [item for item in snapshot["timeline"] if item["kind"] == "tool_card"]
assert len(tool_cards) == 1 and tool_cards[0]["status"] == "failed", tool_cards
PY

  run_lifecycle_message "run-life-retry" "retry smoke" "${message_json}"
  python3 - "${tmp_dir}/eino-workbench.db" <<'PY'
import sqlite3, time, sys
db_path = sys.argv[1]
now = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
with sqlite3.connect(db_path) as conn:
    conn.execute("update runs set status = ?, safe_error = ?, updated_at = ? where run_id = ?", ("failed", "schema_invalid", now, "run-life-retry"))
PY
  run_lifecycle_action "run-life-retry" "retry" "${capability_json}"
  retry_run_id="$(python3 - "${capability_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["run_id"] != "run-life-retry", body
assert body["status"] == "accepted", body
print(body["run_id"])
PY
  )"
  run_lifecycle_action "run-life-retry" "retry" "${capability_json}"
  duplicate_retry_run_id="$(python3 - "${capability_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["status"] == "accepted", body
print(body["run_id"])
PY
  )"
  [[ "${duplicate_retry_run_id}" == "${retry_run_id}" ]]
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${retry_run_id}" -o "${snapshot_json}"
  python3 - "${snapshot_json}" "${retry_run_id}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["run_id"] == sys.argv[2]
assert body["status"] == "created", body
user_messages = [item.get("content") for item in body["timeline"] if item.get("kind") == "user_message"]
assert user_messages == ["retry run run-life-retry"], user_messages
assert "retry smoke" not in json.dumps(body, ensure_ascii=False), body
PY

  run_lifecycle_message "run-life-pending" "pending timeout smoke" "${message_json}"
  python3 - "${tmp_dir}/eino-workbench.db" <<'PY'
import sqlite3, time, sys
db_path = sys.argv[1]
now = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
with sqlite3.connect(db_path) as conn:
    conn.execute("update runs set status = ?, updated_at = ? where run_id = ?", ("waiting", now, "run-life-pending"))
    conn.execute(
        "insert into pending_interactions(pending_id, run_id, kind, status, resume_ref, checkpoint_ref, question, operation_name, risk_summary, target_summary, input_mode, candidates_json, expires_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
        ("pending-life-1", "run-life-pending", "clarification", "waiting", "resume-safe-life-1", "checkpoint_ref:life-1", "请选择实体", "", "", "", "single_choice", "[]", ""),
    )
PY
  run_lifecycle_action "run-life-pending" "pending_timeout" "${capability_json}"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/run-life-pending/replay" -o "${replay_json}"
  python3 - "${replay_json}" "${tmp_dir}/eino-workbench.db" <<'PY'
import json, sqlite3, sys
replay = json.load(open(sys.argv[1]))
assert replay["view"]["status"] == "failed", replay["view"]
assert replay["view"]["inspector"]["runtime"]["safe_error"] == "pending_timeout", replay["view"]["inspector"]["runtime"]
pending_cards = [item for item in replay["view"]["timeline"] if item["kind"] == "clarification_card"]
assert len(pending_cards) == 1 and pending_cards[0]["status"] == "expired", pending_cards
with sqlite3.connect(sys.argv[2]) as conn:
    pending_status = conn.execute("select status from pending_interactions where pending_id = ?", ("pending-life-1",)).fetchone()
assert pending_status == ("expired",), pending_status
PY

  echo "run-lifecycle smoke passed"
  exit 0
fi

if [[ "${scenario}" == "budget" ]]; then
  # Phase 7.2 最小预算门禁：预算超限必须落 Product Facts、audit 和 replay，不由 telemetry 直接驱动产品出口。
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_workbench_message_request.v1","message":"budget smoke","client_request_id":"client-run-budget","run_id":"run-budget-smoke"}' \
    "${base_url}/api/workspaces/ws_smoke/messages" \
    -o "${message_json}"
  python3 - "${tmp_dir}/eino-workbench.db" <<'PY'
import sqlite3, time, sys
db_path = sys.argv[1]
now = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
with sqlite3.connect(db_path) as conn:
    conn.execute("update runs set status = ?, safe_error = '', updated_at = ? where run_id = ?", ("running", now, "run-budget-smoke"))
    conn.execute(
        "insert into tool_calls(tool_call_id, run_id, tool_id, display_name, status, args_hash, args_preview, created_at, ended_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
        ("call-budget-1", "run-budget-smoke", "cap.smoke.read", "只读查询", "running", "", "scope=smoke", now, ""),
    )
PY
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_run_lifecycle_request.v1","action":"budget_exceeded","client_request_id":"client-budget-exceeded"}' \
    "${base_url}/api/workspaces/ws_smoke/runs/run-budget-smoke/lifecycle" \
    -o "${capability_json}"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/run-budget-smoke/replay" -o "${replay_json}"
  python3 - "${capability_json}" "${replay_json}" <<'PY'
import json, sys
action = json.load(open(sys.argv[1]))
replay = json.load(open(sys.argv[2]))
payload = json.dumps([action, replay], ensure_ascii=False).lower()
for forbidden in ("authorization", "api_key", "api_token", "provider_payload", "raw_payload", "raw prompt", "raw provider", "raw body", "bearer "):
    assert forbidden not in payload, forbidden
assert action["schema_version"] == "eino_action_result.v1"
assert action["status"] == "failed", action
assert replay["view"]["status"] == "failed", replay["view"]
assert replay["view"]["inspector"]["runtime"]["safe_error"] == "budget_exceeded", replay["view"]["inspector"]["runtime"]
assert any(event.get("event_type") == "budget" and event.get("safe_summary") == "budget exceeded" for event in replay["events"] if isinstance(event, dict)), replay["events"]
tool_cards = [item for item in replay["view"]["timeline"] if item["kind"] == "tool_card"]
assert len(tool_cards) == 1 and tool_cards[0]["status"] == "cancelled", tool_cards
PY

  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_workbench_message_request.v1","message":"budget pending smoke","client_request_id":"client-run-budget-pending","run_id":"run-budget-pending-smoke"}' \
    "${base_url}/api/workspaces/ws_smoke/messages" \
    -o "${message_json}"
  python3 - "${tmp_dir}/eino-workbench.db" <<'PY'
import sqlite3, time, sys
db_path = sys.argv[1]
now = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
with sqlite3.connect(db_path) as conn:
    conn.execute("update runs set status = ?, safe_error = '', updated_at = ? where run_id = ?", ("waiting", now, "run-budget-pending-smoke"))
    conn.execute(
        "insert into pending_interactions(pending_id, run_id, kind, status, resume_ref, checkpoint_ref, question, operation_name, risk_summary, target_summary, input_mode, candidates_json, expires_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
        ("pending-budget-1", "run-budget-pending-smoke", "approval", "waiting", "resume-safe-budget-smoke", "checkpoint_ref:budget-smoke", "是否批准？", "预算验收操作", "预算等待态验收", "目标已脱敏", "", "[]", ""),
    )
PY
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_run_lifecycle_request.v1","action":"budget_exceeded","client_request_id":"client-budget-pending-exceeded"}' \
    "${base_url}/api/workspaces/ws_smoke/runs/run-budget-pending-smoke/lifecycle" \
    -o "${capability_json}"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/run-budget-pending-smoke/replay" -o "${replay_json}"
  python3 - "${tmp_dir}/eino-workbench.db" "${capability_json}" "${replay_json}" <<'PY'
import json, sqlite3, sys
db_path, action_path, replay_path = sys.argv[1], sys.argv[2], sys.argv[3]
action = json.load(open(action_path))
replay = json.load(open(replay_path))
payload = json.dumps([action, replay], ensure_ascii=False).lower()
for forbidden in ("authorization", "api_key", "api_token", "provider_payload", "raw_payload", "raw prompt", "raw provider", "raw body", "bearer ", "checkpoint_ref:"):
    assert forbidden not in payload, forbidden
assert action["run_id"] == "run-budget-pending-smoke", action
assert action["status"] == "failed", action
assert replay["view"]["status"] == "failed", replay["view"]
assert replay["view"]["inspector"]["runtime"]["safe_error"] == "budget_exceeded", replay["view"]["inspector"]["runtime"]
pending_cards = [item for item in replay["view"]["timeline"] if item["kind"] == "approval_card" and item.get("pending_id") == "pending-budget-1"]
assert len(pending_cards) == 1 and pending_cards[0]["status"] == "expired", pending_cards
with sqlite3.connect(db_path) as conn:
    pending_status = conn.execute("select status from pending_interactions where pending_id = ?", ("pending-budget-1",)).fetchone()[0]
assert pending_status == "expired", pending_status
assert any(event.get("event_type") == "budget" and event.get("safe_summary") == "budget exceeded" for event in replay["events"] if isinstance(event, dict)), replay["events"]
PY
  echo "budget smoke passed"
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
assert any(event_type == "tool" and "capability selected: cap.smoke.write" in summary and "approval_required" in summary for event_type, summary in rows), rows
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

if [[ "${scenario}" == "action-consistency" ]]; then
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-consistency-smoke","client_request_id":"client-smoke-action-consistency","capability_hint":"cap.smoke.read","input":{"text":"same facts smoke"}}' \
    "${base_url}/api/workspaces/ws_smoke/agent/actions" \
    -o "${capability_json}"
  run_id="$(python3 - "${capability_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
payload = json.dumps(body, ensure_ascii=False).lower()
# 同源事实验收不能允许 provider 原始材料或凭据痕迹进入 Action API。
for forbidden in ("authorization", "api_key", "api_token", "credential_ref", "provider_payload", "raw_payload", "raw provider", "raw body", "bearer "):
    assert forbidden not in payload, forbidden
assert body["schema_version"] == "eino_action_result.v1"
assert body["workspace_id"] == "ws_smoke"
assert body["status"] == "completed", body
assert body["snapshot_url"].endswith(f"/runs/{body['run_id']}")
assert body["stream_ref"].endswith(f"/runs/{body['run_id']}/stream")
assert len(body["result_cards"]) == 1, body
card = body["result_cards"][0]
assert card["tool_call_id"], card
assert card["title"] == "Phase 3 read smoke", card
assert card["status"] == "succeeded", card
assert card["safe_summary"] == "Phase 3 read smoke 完成", card
assert card["structured_result"]["schema_version"] == "tool.structured_result.v1", card
assert body["audit_refs"] == [f"/api/workspaces/ws_smoke/runs/{body['run_id']}/replay"], body["audit_refs"]
print(body["run_id"])
PY
  )"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${run_id}" -o "${snapshot_json}"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/views/current" -o "${current_view_json}"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${run_id}/replay" -o "${replay_json}"
  wrong_snapshot_status="$(curl -sS -o /dev/null -w "%{http_code}" "${base_url}/api/workspaces/ws_other/runs/${run_id}")"
  wrong_replay_status="$(curl -sS -o /dev/null -w "%{http_code}" "${base_url}/api/workspaces/ws_other/runs/${run_id}/replay")"
  wrong_stream_status="$(curl -sS -o /dev/null -w "%{http_code}" "${base_url}/api/workspaces/ws_other/runs/${run_id}/stream")"
  [[ "${wrong_snapshot_status}" == "404" ]]
  [[ "${wrong_replay_status}" == "404" ]]
  [[ "${wrong_stream_status}" == "404" ]]
  python3 - "${capability_json}" "${snapshot_json}" "${current_view_json}" "${replay_json}" "${run_id}" <<'PY'
import json, sys
action = json.load(open(sys.argv[1]))
snapshot = json.load(open(sys.argv[2]))
current = json.load(open(sys.argv[3]))
replay = json.load(open(sys.argv[4]))
run_id = sys.argv[5]

def reject_raw(name, value):
    payload = json.dumps(value, ensure_ascii=False).lower()
    # 这些检查只验证产品投影，不检查 provider 内部存储。
    for forbidden in ("authorization", "api_key", "api_token", "credential_ref", "provider_payload", "raw_payload", "raw provider", "raw body", "bearer "):
        assert forbidden not in payload, f"{name} leaked {forbidden}"

for name, value in (("action", action), ("snapshot", snapshot), ("current", current), ("replay", replay)):
    reject_raw(name, value)

assert snapshot["schema_version"] == "eino_workbench_view.v1"
assert current["schema_version"] == "eino_workbench_view.v1"
assert replay["schema_version"] == "eino_replay_view.v1"
assert action["run_id"] == snapshot["run_id"] == current["run_id"] == replay["run_id"] == run_id
assert current == snapshot, "views/current must project the same latest run snapshot"
assert replay["view"] == snapshot, "replay view must be rebuilt from the same Product Facts snapshot"

result_card = action["result_cards"][0]
snapshot_cards = [item for item in snapshot["timeline"] if item["kind"] == "tool_card"]
current_cards = [item for item in current["timeline"] if item["kind"] == "tool_card"]
replay_cards = [item for item in replay["view"]["timeline"] if item["kind"] == "tool_card"]
assert len(snapshot_cards) == len(current_cards) == len(replay_cards) == 1
for card in (snapshot_cards[0], current_cards[0], replay_cards[0]):
    assert card["tool_call_id"] == result_card["tool_call_id"], card
    assert card["status"] == result_card["status"], card
    assert card["safe_summary"] == result_card["safe_summary"], card
    assert card["structured_result"]["schema_version"] == result_card["structured_result"]["schema_version"], card

# assistant 和 pending 在当前 mock read 场景中通常为空；一旦出现，也必须来自同一 timeline。
assistant_messages = [item for item in snapshot["timeline"] if item["kind"] == "assistant_message"]
if action.get("final_answer"):
    assert any(item.get("content") == action["final_answer"] for item in assistant_messages), assistant_messages
pending_items = [item for item in snapshot["timeline"] if item["kind"] in ("approval_card", "clarification_card")]
if action.get("waiting") is not None:
    assert any(item.get("pending_id") for item in pending_items), pending_items
else:
    assert not pending_items, pending_items

structured = snapshot["inspector"]["structured"]
assert structured["tool_call_id"] == result_card["tool_call_id"], structured
assert structured["structured_result"]["schema_version"] == "tool.structured_result.v1", structured
audit = snapshot["inspector"]["audit"]
assert any(row["event_type"] == "tool" and row["safe_summary"] == "tool completed: cap.smoke.read" for row in audit), audit
events = replay["events"]
assert any(row.get("type") == "tool.updated" and row.get("tool", {}).get("safe_summary") == result_card["safe_summary"] for row in events), events
assert any(row.get("event_type") == "tool" and row.get("safe_summary") == "tool completed: cap.smoke.read" for row in events), events
PY
  curl -fsS -D "${stream_headers}" "${base_url}/api/workspaces/ws_smoke/runs/${run_id}/stream" -o "${stream_body}"
  grep -qi '^Content-Type: text/event-stream' "${stream_headers}"
  grep -q "event: tool.updated" "${stream_body}"
  grep -q "Phase 3 read smoke 完成" "${stream_body}"
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    tool_rows = conn.execute("select tool_call_id, tool_id, status from tool_calls where run_id = ?", (run_id,)).fetchall()
    result_rows = conn.execute(
        "select tr.structured_schema_version, tr.safe_summary from tool_results tr join tool_calls tc on tr.tool_call_id = tc.tool_call_id where tc.run_id = ?",
        (run_id,),
    ).fetchall()
    audit_rows = conn.execute("select event_type, safe_summary from audit_events where run_id = ?", (run_id,)).fetchall()
assert len(tool_rows) == 1, tool_rows
assert tool_rows[0][1] == "cap.smoke.read" and tool_rows[0][2] == "succeeded", tool_rows
assert result_rows == [("tool.structured_result.v1", "Phase 3 read smoke 完成")], result_rows
assert any(event_type == "tool" and summary == "tool completed: cap.smoke.read" for event_type, summary in audit_rows), audit_rows
PY
  echo "action-consistency smoke passed"
  exit 0
fi

if [[ "${scenario}" == "replay" ]]; then
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-replay-smoke","client_request_id":"client-smoke-replay","capability_hint":"cap.smoke.read","input":{"text":"replay smoke"}}' \
    "${base_url}/api/workspaces/ws_smoke/agent/actions" \
    -o "${capability_json}"
  run_id="$(python3 - "${capability_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["status"] == "completed", body
assert len(body["result_cards"]) == 1, body
print(body["run_id"])
PY
  )"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${run_id}" -o "${snapshot_json}"
  curl -fsS "${base_url}/api/workspaces/ws_smoke/runs/${run_id}/replay" -o "${replay_json}"
  python3 - "${capability_json}" "${snapshot_json}" "${replay_json}" "${run_id}" <<'PY'
import json, sys
action = json.load(open(sys.argv[1]))
snapshot = json.load(open(sys.argv[2]))
replay = json.load(open(sys.argv[3]))
run_id = sys.argv[4]
for name, value in (("action", action), ("snapshot", snapshot), ("replay", replay)):
    payload = json.dumps(value, ensure_ascii=False).lower()
    # replay 只能暴露 Product Facts 投影，不能包含 provider 原始响应或凭据。
    for forbidden in ("authorization", "api_key", "api_token", "credential_ref", "provider_payload", "raw_payload", "raw provider", "raw body", "bearer "):
        assert forbidden not in payload, f"{name} leaked {forbidden}"

assert replay["schema_version"] == "eino_replay_view.v1"
assert replay["workspace_id"] == "ws_smoke"
assert replay["run_id"] == snapshot["run_id"] == action["run_id"] == run_id
assert replay["view"] == snapshot, "replay view must exactly match run snapshot projection"
assert replay["view"]["inspector"]["structured"]["structured_result"]["schema_version"] == "tool.structured_result.v1"
event_ids = [row.get("event_id") for row in replay["events"] if isinstance(row, dict) and "event_id" in row]
assert event_ids and all(event_id.startswith(f"{run_id}:") for event_id in event_ids), event_ids
assert any(row.get("type") == "tool.updated" for row in replay["events"] if isinstance(row, dict)), replay["events"]
assert any(row.get("schema_version") == "eino_audit_event.v1" and row.get("event_type") == "tool" for row in replay["events"] if isinstance(row, dict)), replay["events"]
result_card = action["result_cards"][0]
tool_cards = [item for item in replay["view"]["timeline"] if item["kind"] == "tool_card"]
assert len(tool_cards) == 1
assert tool_cards[0]["safe_summary"] == result_card["safe_summary"] == "Phase 3 read smoke 完成"
PY
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    run_status = conn.execute("select status from runs where run_id = ?", (run_id,)).fetchone()
    result_count = conn.execute(
        "select count(*) from tool_results tr join tool_calls tc on tr.tool_call_id = tc.tool_call_id where tc.run_id = ?",
        (run_id,),
    ).fetchone()[0]
    audit_count = conn.execute("select count(*) from audit_events where run_id = ? and event_type = 'tool'", (run_id,)).fetchone()[0]
assert run_status == ("succeeded",), run_status
assert result_count == 1, result_count
assert audit_count >= 1, audit_count
PY
  echo "replay smoke passed"
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

if [[ "${scenario}" == "mcp-mock" ]]; then
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-mcp-mock","client_request_id":"client-smoke-mcp-mock","capability_hint":"mcp.mock.tool.asset_lookup","input":{"text":"mcp asset smoke"}}' \
    "${base_url}/api/workspaces/ws_smoke/agent/actions" \
    -o "${capability_json}"
  run_id="$(python3 - "${capability_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["status"] == "completed", body
assert len(body["result_cards"]) == 1, body
card = body["result_cards"][0]
assert card["title"] == "MCP 资产查询", card
assert card["safe_summary"] == "MCP 资产查询完成", card
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
assert tool_cards[0]["safe_summary"] == "MCP 资产查询完成", tool_cards[0]
PY
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    rows = conn.execute(
        "select tc.tool_id, tr.safe_summary from tool_results tr join tool_calls tc on tr.tool_call_id = tc.tool_call_id where tc.run_id = ?",
        (run_id,),
    ).fetchall()
assert any(tool_id == "mcp.mock.tool.asset_lookup" and summary == "MCP 资产查询完成" for tool_id, summary in rows), rows
PY
  echo "mcp-mock smoke passed"
  exit 0
fi

if [[ "${scenario}" == "fobrain-poc" ]]; then
  curl -fsS \
    -H "Content-Type: application/json" \
    -d '{"schema_version":"eino_action_request.v1","action_id":"action-fobrain-poc","client_request_id":"client-smoke-fobrain-poc","capability_hint":"tool.fobrain.current_user_context","input":{"text":"查看当前 Fobrain 用户信息"}}' \
    "${base_url}/api/workspaces/ws_smoke/agent/actions" \
    -o "${capability_json}"
  run_id="$(python3 - "${capability_json}" <<'PY'
import json, sys
body = json.load(open(sys.argv[1]))
assert body["schema_version"] == "eino_action_result.v1"
assert body["status"] == "completed", body
assert len(body["result_cards"]) == 1, body
card = body["result_cards"][0]
assert card["title"] == "读取当前 Fobrain 用户信息", card
assert "Fobrain 当前用户" in card["safe_summary"], card
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
assert "Fobrain 当前用户" in tool_cards[0]["safe_summary"], tool_cards[0]
encoded = json.dumps(body, ensure_ascii=False).lower()
for forbidden in ("authorization", "bearer ", "api_key", "api_token", "token", "credential_ref", "raw provider", "raw body", "checkpoint-raw", "interrupt-raw"):
    assert forbidden not in encoded, forbidden
PY
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import sqlite3, sys
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    rows = conn.execute(
        "select tc.tool_id, tr.structured_schema_version, tr.safe_summary from tool_results tr join tool_calls tc on tr.tool_call_id = tc.tool_call_id where tc.run_id = ?",
        (run_id,),
    ).fetchall()
assert any(tool_id == "tool.fobrain.current_user_context" and schema == "tool.structured_result.v1" and "Fobrain 当前用户" in summary for tool_id, schema, summary in rows), rows
PY
  mkdir -p test-results
  python3 - "${tmp_dir}/eino-workbench.db" "${run_id}" <<'PY'
import datetime, json, sys
import sqlite3
db_path, run_id = sys.argv[1], sys.argv[2]
with sqlite3.connect(db_path) as conn:
    row = conn.execute(
        "select tr.result_ref, tr.safe_summary from tool_results tr join tool_calls tc on tr.tool_call_id = tc.tool_call_id where tc.run_id = ? and tc.tool_id = ?",
        (run_id, "tool.fobrain.current_user_context"),
    ).fetchone()
assert row, "missing fobrain tool result"
result_ref, safe_summary = row
report = {
    "schema_version": "eino.fobrain_provider_poc_report.v1",
    "scenario": "fobrain-poc",
    "status": "passed",
    "provider_mode": "mock",
    "capability_id": "tool.fobrain.current_user_context",
    "structured_result_schema": "tool.structured_result.v1",
    "business_result_schema": "fobrain.tool_result.v2",
    "result_ref": result_ref,
    "safe_summary": safe_summary,
    "policy_decision": "allowed",
    "credential_binding_status": "bound",
    "connector_status": "available",
    "redaction_checks": ["no_token", "no_auth_header", "no_secret_ref", "no_raw_body", "no_checkpoint", "no_interrupt"],
    "failure_category": "none",
    "run_id": run_id,
    "report_created_at": datetime.datetime.now(datetime.timezone.utc).isoformat().replace("+00:00", "Z"),
}
with open("test-results/eino-workbench-fobrain-provider-poc-report.json", "w", encoding="utf-8") as f:
    json.dump(report, f, ensure_ascii=False, indent=2)
PY
  if [[ ! -f "${provided_config:-configs/eino-workbench.local.yaml}" ]]; then
    write_fobrain_skip_report "本地 LLM/Fobrain live 配置文件不存在"
  elif ! go run ./scripts/eino_workbench_config_prepare.go --check-llm-api-key --source "${provided_config:-configs/eino-workbench.local.yaml}" >/dev/null 2>&1; then
    write_fobrain_skip_report "本地 LLM api_key 未配置"
  elif ! go run ./scripts/eino_workbench_config_prepare.go --check-fobrain-live-credential --source "${provided_config:-configs/eino-workbench.local.yaml}" >/dev/null 2>&1; then
    write_fobrain_skip_report "本地 Fobrain live 凭据未配置"
  else
    write_fobrain_skip_report "Phase 5 不支持 Fobrain live read 或真实模型工具选择声明" "phase_unsupported"
  fi
  echo "fobrain-poc smoke passed"
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
