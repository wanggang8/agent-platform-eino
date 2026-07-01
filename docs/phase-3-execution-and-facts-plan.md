# Phase 3 Execution and Facts Plan

> **For agentic workers:** implement Phase 3 only after this plan, `docs/07-implementation-plan.md`, `docs/facts-contract.md`, `docs/run-lifecycle.md`, `docs/conversation-context.md`, and `docs/intent-and-capability-selection.md` agree. If they conflict, stop and clarify before coding.

**Goal:** wire Eino ChatModelAgent execution to durable Product Facts, then expose the same facts through Workbench SSE and Action API.

**Architecture:** Eino Runner is an execution input, not a product contract. Runner events are mapped into Product Facts in a single backend path; Workbench, Action API, audit, replay, and model context all read product projections from those facts.

**Tech Stack:** Go, CloudWeGo Eino `v0.9.12`, SQLite store, existing `llm.Provider`, `capabilities.Registry`, `facts.Repository`, and `product.Projection` boundaries.

---

## Source Basis

- Pinned module: `github.com/cloudwego/eino v0.9.12`.
- Runner APIs: `adk.NewRunner`, `adk.RunnerConfig`, `Runner.Run`, `Runner.ResumeWithParams`, `adk.ResumeParams.Targets`.
- Agent APIs: `adk.ChatModelAgentConfig`, `adk.ToolsConfig`, `adk.WithAfterToolCallsHook`.
- HITL APIs: `adk.InterruptInfo`, `adk.InterruptCtx`, checkpoint-backed `ResumeInfo`.
- Callback APIs are diagnostics only; callback output must not become Product Facts.

Use local module source and pinned ADR as implementation truth. Online docs are secondary concept references.

## Mandatory Order

### 1. Version Gate

Confirm `go.mod` still pins Eino `v0.9.12`. Remove `internal/einoapp/execution/eino_version_pin.go` only when real execution code imports Eino in the same change. Do not introduce `eino-ext` until a real provider implementation needs it.

### 2. Product Facts Model

Expand `internal/einoapp/facts` before touching Runner integration:

- run lifecycle: `created`, `running`, `waiting`, `succeeded`, `failed`, `cancelled`, `stopped`
- turn/message records
- tool call/result records
- pending approval/clarification records
- audit events
- context snapshots
- idempotency records for resume, stop, cancel, retry, and mutation requests

Facts writes must be append-only or explicit state migrations. Raw provider payload, raw checkpoint id, interrupt id, reusable resume token, credential ref, and secret material cannot be stored as product facts. `checkpoint_ref` is allowed only as an internal, non-displayable safe reference and must not be the raw Eino checkpoint id.

### 3. SQLite Repository

Implement `internal/einoapp/store/sqlite` before Runner integration. The repository owns transactions for:

- creating runs and turns
- appending mapped events
- advancing lifecycle state
- consuming `resume_ref` exactly once
- writing context snapshots before model calls
- writing checkpoint references as internal-only records

SQLite is the Phase 3 source of truth. The existing memory repository remains a test fixture only.

### 4. Product Projection

Extend `internal/einoapp/product` after repository tests pass. Projections produce:

- Workbench view
- run snapshot
- stream event cursor view
- Action API result
- replay view
- safe context projection

SSE event ids must come from Product Facts sequence/cursor values, not Eino event ids or frontend counters.

### 5. Execution Event Mapper

Create an event mapper that converts Eino events to facts:

- assistant message delta/final -> message facts
- tool call request -> tool call fact with capability metadata
- tool result -> StructuredResult-compatible fact only
- interrupt -> pending interaction fact with product `resume_ref`
- error -> redacted failure fact
- lifecycle transition -> run state fact

The mapper is the only package that may understand Eino `AgentEvent` shape. HTTP, Workbench, Action API, replay, and audit must not import or switch on Eino event structs.

### 6. ChatModelAgent Runner

Integrate Eino only after facts, repository, projection, and mapper tests exist.

- Ordinary natural language enters ChatModelAgent.
- `product_action`, `capability_hint`, and trusted `intent_hint` can narrow candidate capabilities but cannot directly pick a provider branch.
- Capability Registry generates tool metadata and schemas; no `fobrain_*` keyword routing.
- Before every model call, build and store a safe context snapshot from Product Facts.
- `WithAfterToolCallsHook` may trigger context refresh, but it must not write product facts directly.

Phase 3 may keep using mock `llm.Provider`. Real OpenAI-compatible provider implementation belongs to Phase 4 unless a new ADR changes that boundary.

### 7. Resume and Checkpoint Boundary

Product layer exposes only one-time `resume_ref`. Internally it maps to checkpoint id, interrupt address, pending id, and expected resume schema.

Phase 3 records the facts/store/projection boundaries needed by resume, but full durable CheckPointStore, approval interrupt, clarification interrupt, and after-restart resume behavior are implemented in Phase 6. Until Phase 6, resume/checkpoint tests may use controlled fakes and must not claim production HITL recovery.

Resume flow:

```text
resume_ref
  -> transaction: validate pending + consume resume_ref
  -> load checkpoint mapping
  -> build typed resume data
  -> Runner.ResumeWithParams(ctx, checkpoint_id, &adk.ResumeParams{Targets: map[string]any{interrupt_address: resume_data}})
  -> map resumed events back into Product Facts
```

Duplicate resume returns the existing terminal pending state and must not re-run tools or mutations. Missing or corrupt checkpoint returns a redacted safe error and preserves audit evidence.

### 8. HTTP, SSE, Action API

HTTP handlers call execution commands and product projections only:

- `/messages` starts a run/turn and returns accepted or initial snapshot.
- SSE reads Product Facts cursor projections and never streams raw Runner events.
- Action API uses the same command path and projection path as Workbench.
- Errors use the existing unified error envelope.

## Required Test Shape

Run Phase 3 tasks in this order:

1. facts model tests
2. SQLite repository tests
3. product projection tests
4. event mapper tests with fake Eino events
5. runner integration tests with mock model/tools
6. HTTP/SSE/Action smoke tests

Minimum commands before claiming all of Phase 3 complete:

```bash
go test ./internal/einoapp/facts ./internal/einoapp/store/sqlite ./internal/einoapp/product ./internal/einoapp/execution -count=1
go test ./internal/einoapp/httpapi -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario chat-stream
bash scripts/eino_workbench_server_smoke.sh --scenario action-basic
```

If a scenario is not implemented yet, the script must exit `2` and the task cannot be marked complete for that scenario.

## Clarification Triggers

Stop and ask before coding if any of these become unclear:

- whether a field belongs in Product Facts or only in internal execution state
- whether an Eino event should be projected to users
- whether a resume payload can be retried
- whether a tool result is StructuredResult-compliant
- whether an intent hint is trusted enough to narrow capabilities
- whether a provider error contains raw payload or secrets
- whether a task needs real LLM behavior before Phase 4
