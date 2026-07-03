import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(scriptDir, "..");
const schemasDir = path.join(root, "docs/schemas");
const outFile = path.join(root, "web/eino-workbench/src/contracts/generated.ts");
const checkOnly = process.argv.includes("--check");

// walk 收集 schema 文件，保证 contract 生成顺序稳定。
function walk(dir) {
  const out = [];
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) out.push(...walk(full));
    if (entry.isFile() && entry.name.endsWith(".json")) out.push(full);
  }
  return out.sort();
}

// readSchema 只读取 schema 的公开标识，不把 schema 内容复制进前端。
function readSchema(file) {
  const schema = JSON.parse(fs.readFileSync(file, "utf8"));
  if (!schema.$id || !schema.title) {
    throw new Error(`${path.relative(root, file)} must define $id and title`);
  }
  return {
    id: schema.$id,
    title: schema.title,
    path: path.relative(root, file).replaceAll(path.sep, "/")
  };
}

const schemas = walk(schemasDir).map(readSchema);
// typeDefinitions 是前端 contract 类型的集中生成模板，避免组件手写平行 DTO。
const typeDefinitions = `
export type JsonPrimitive = string | number | boolean | null;
export type JsonValue = JsonPrimitive | { readonly [key: string]: JsonValue } | readonly JsonValue[];

export type RunStatus = "created" | "running" | "waiting" | "succeeded" | "failed" | "cancelled" | "stopped";
export type ActionStatus = "accepted" | "running" | "waiting" | "completed" | "failed" | "cancelled" | "stopped" | "denied" | "blocked";
export type RunLifecycleAction = "cancel" | "stop" | "provider_timeout" | "pending_timeout" | "budget_exceeded" | "retry";
export type TimelineKind = "user_message" | "assistant_message" | "tool_card" | "approval_card" | "clarification_card" | "run_notice";
export type InspectorTab = "evidence" | "structured" | "runtime" | "audit";
export type PendingKind = "approval" | "clarification";
export type PendingStatus = "waiting" | "submitted" | "approved" | "rejected" | "cancelled" | "expired" | "consumed";
export type StreamEventType = "message.delta" | "message.updated" | "tool.updated" | "pending.updated" | "run.updated" | "view.replaced" | "audit.updated";
export type AuditEventType = "message" | "tool" | "pending" | "resume" | "approval" | "clarification" | "lifecycle" | "budget" | "safety" | "error";
export type DisplayType = "entity_collection" | "entity_detail" | "metrics_summary" | "operation_result" | "connector_status" | "entity_resolution";
export type StructuredStatus = "resolved" | "waiting" | "pending_approval" | "not_found" | "empty" | "failed" | "partial";
export type Tone = "neutral" | "info" | "success" | "warning" | "danger";
export type InputMode = "single_choice" | "multi_choice" | "free_text" | "mixed";

export type PendingCandidateField = {
  readonly label: string;
  readonly value: string;
};

export type PendingCandidate = {
  readonly candidate_ref: string;
  readonly label: string;
  readonly description?: string;
  readonly entity_type: "person" | "asset" | "vulnerability" | "department" | "business";
  readonly safe_fields?: readonly PendingCandidateField[];
};

export type DisplayField = {
  readonly key?: string;
  readonly label: string;
  readonly value: JsonValue;
  readonly type?: "text" | "number" | "status" | "severity" | "code" | "datetime" | "duration" | "badge";
  readonly tone?: Tone;
  readonly copyable?: boolean;
};

export type StructuredMetric = {
  readonly label: string;
  readonly value: JsonValue;
  readonly sub_label?: string;
  readonly tone?: Tone;
  readonly icon_hint?: string;
};

export type StructuredColumn = {
  readonly key: string;
  readonly label: string;
  readonly type?: DisplayField["type"];
  readonly width?: "compact" | "normal" | "wide" | "fill";
  readonly align?: "left" | "center" | "right";
  readonly copyable?: boolean;
};

export type StructuredItem = {
  readonly entity_ref?: string;
  readonly row_ref?: string;
  readonly display_name?: string;
  readonly masked_ip?: string;
  readonly asset_type?: string;
  readonly business_system?: string;
  readonly owner_name?: string;
  readonly vulnerability_ref?: string;
  readonly title?: string;
  readonly severity?: string;
  readonly status?: string;
  readonly affected_assets?: number;
  readonly ticket_ref?: string;
  readonly target_status?: string;
  readonly risk_summary?: string;
  readonly summary?: string;
  readonly badges?: readonly { readonly label: string; readonly tone?: Tone }[];
};

export type StructuredAction = {
  readonly id: string;
  readonly label: string;
  readonly kind?: "open" | "download" | "export" | "retry" | "resume" | "approve" | "reject" | "select_candidate";
  readonly role?: "open" | "download" | "export" | "retry" | "resume" | "approve" | "reject" | "confirm" | "cancel";
  readonly tone?: Tone;
  readonly enabled?: boolean;
};

export type GenericStructuredResult = {
  readonly schema_version: "tool.structured_result.v1";
  readonly status: Exclude<StructuredStatus, "not_found">;
  readonly data: {
    readonly summary?: string;
    readonly facts?: readonly DisplayField[];
  };
  readonly metadata: {
    readonly safe: true;
    readonly result_ref?: string;
    readonly source?: string;
  };
};

export type FobrainStructuredResult = {
  readonly schema_version: "fobrain.tool_result.v2";
  readonly tool_id: string;
  readonly display_type: DisplayType;
  readonly entity_type?: "asset" | "vulnerability" | "ticket" | "business_system" | "user" | "permission" | "risk" | "connector" | "entity";
  readonly status: StructuredStatus;
  readonly query: Record<string, JsonValue>;
  readonly data: {
    readonly title?: string;
    readonly summary?: string;
    readonly metrics?: readonly StructuredMetric[];
    readonly columns?: readonly StructuredColumn[];
    readonly items?: readonly StructuredItem[];
    readonly facts?: readonly DisplayField[];
    readonly actions?: readonly StructuredAction[];
    readonly resume_refs?: readonly string[];
    readonly candidate_refs?: readonly string[];
    readonly resolved_entity?: PendingCandidate;
    readonly candidates?: readonly PendingCandidate[];
  };
  readonly metadata: {
    readonly safe: true;
    readonly source?: "fobrain";
    readonly result_ref?: string;
    readonly attention_reason?: string;
    readonly observed_at?: string;
  };
};

export type StructuredResult = GenericStructuredResult | FobrainStructuredResult;

export type AuditEvent = {
  readonly schema_version: "eino_audit_event.v1";
  readonly audit_id: string;
  readonly run_id: string;
  readonly event_type: AuditEventType;
  readonly safe_summary: string;
  readonly actor: string;
  readonly created_at: string;
};

export type PendingInteraction = {
  readonly schema_version: "eino_workbench_pending_interaction.v1";
  readonly pending_id: string;
  readonly run_id: string;
  readonly kind: PendingKind;
  readonly status: PendingStatus;
  readonly question?: string;
  readonly operation_name?: string;
  readonly risk_summary?: string;
  readonly target_summary?: string;
  readonly input_mode?: InputMode;
  readonly candidates?: readonly PendingCandidate[];
  readonly resume_ref?: string;
  readonly expires_at?: string;
};

export type TimelineItem = {
  readonly item_id: string;
  readonly kind: TimelineKind;
  readonly tool_call_id?: string;
  readonly pending_id?: string;
  readonly content?: string;
  readonly structured_result?: StructuredResult;
  readonly status?: string;
  readonly safe_summary?: string;
  readonly input_mode?: InputMode;
  readonly candidates?: readonly PendingCandidate[];
};

export type WorkbenchInspector = {
  readonly tabs: readonly InspectorTab[];
  readonly evidence?: readonly {
    readonly label: string;
    readonly value: string;
    readonly target_ref?: string;
    readonly observed_at?: string;
  }[];
  readonly structured?: {
    readonly tool_call_id?: string;
    readonly result_ref?: string;
    readonly structured_result?: StructuredResult;
  };
  readonly runtime?: {
    readonly status?: RunStatus;
    readonly model_label?: string;
    readonly safe_error?: string;
    readonly pending_id?: string;
  };
  readonly audit?: readonly AuditEvent[];
};

export type WorkbenchView = {
  readonly schema_version: "eino_workbench_view.v1";
  readonly workspace_id: string;
  readonly run_id: string;
  readonly status: RunStatus;
  readonly timeline: readonly TimelineItem[];
  readonly inspector: WorkbenchInspector;
};

export type ActionResult = {
  readonly schema_version: "eino_action_result.v1";
  readonly workspace_id: string;
  readonly action_id: string;
  readonly run_id: string;
  readonly status: ActionStatus;
  readonly final_answer?: string;
  readonly result_cards: readonly {
    readonly card_id: string;
    readonly tool_call_id?: string;
    readonly title: string;
    readonly status: "queued" | "running" | "succeeded" | "failed" | "cancelled";
    readonly safe_summary?: string;
    readonly structured_result?: StructuredResult;
  }[];
  readonly approval_refs?: readonly string[];
  readonly resume_refs?: readonly string[];
  readonly waiting?: {
    readonly kind: "approval" | "clarification";
    readonly question: string;
    readonly risk_summary?: string;
    readonly target_summary?: string;
    readonly approval_refs?: readonly string[];
    readonly resume_refs?: readonly string[];
    readonly input_mode?: InputMode;
    readonly candidates?: readonly PendingCandidate[];
  };
  readonly audit_refs: readonly string[];
};

export type RunLifecycleRequest = {
  readonly schema_version: "eino_run_lifecycle_request.v1";
  readonly action: RunLifecycleAction;
  readonly client_request_id: string;
};

export type MessagePatch = {
  readonly message_id: string;
  readonly role: "user" | "assistant" | "system_notice";
  readonly content_delta?: string;
  readonly content?: string;
  readonly status?: "streaming" | "completed" | "failed";
};

export type ToolPatch = {
  readonly tool_call_id: string;
  readonly tool_id?: string;
  readonly display_name?: string;
  readonly status: "queued" | "running" | "succeeded" | "failed" | "cancelled";
  readonly safe_summary?: string;
  readonly structured_result?: StructuredResult;
};

export type RunPatch = {
  readonly status: RunStatus;
  readonly safe_error?: string;
  readonly updated_at?: string;
};

export type WorkbenchStreamEvent = {
  readonly schema_version: "eino_workbench_stream_event.v1";
  readonly event_id: string;
  readonly run_id: string;
  readonly type: StreamEventType;
  readonly sequence: number;
  readonly created_at: string;
  readonly message?: MessagePatch;
  readonly tool?: ToolPatch;
  readonly pending?: PendingInteraction;
  readonly run?: RunPatch;
  readonly view?: WorkbenchView;
  readonly audit?: AuditEvent;
};
`;
const body = `// Generated by scripts/eino_workbench_contract_generate.mjs. Do not edit by hand.
export type ContractSchema = {
  readonly id: string;
  readonly path: string;
  readonly title: string;
};

export const contractSchemas = ${JSON.stringify(schemas, null, 2)} as const satisfies readonly ContractSchema[];
${typeDefinitions.trimEnd()}
`;

// --check 用于 CI/验收确认生成文件没有落后。
if (checkOnly) {
  const current = fs.existsSync(outFile) ? fs.readFileSync(outFile, "utf8") : "";
  if (current !== body) {
    throw new Error(`${path.relative(root, outFile)} is out of date; run npm run eino-workbench:contract-generate`);
  }
  console.log(`contract index is current for ${schemas.length} schemas`);
} else {
  // 非 check 模式才写入 generated.ts，生成文件本身不手工编辑。
  fs.writeFileSync(outFile, body);
  console.log(`generated ${path.relative(root, outFile)} from ${schemas.length} schemas`);
}
