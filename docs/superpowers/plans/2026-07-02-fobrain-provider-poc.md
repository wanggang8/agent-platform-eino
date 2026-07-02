# Fobrain Provider PoC Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a runnable Fobrain provider PoC that registers one read-only capability, invokes a mock/live-safe client, emits StructuredResult candidate, and passes `fobrain-poc` smoke without leaking credentials or raw payload.

**Architecture:** Implement `internal/einoapp/providers/fobrain` as a thin provider package over existing `capabilities.Provider`, `capabilities.Invoker`, `capabilities.EvaluatePolicy`, and `product.StructuredResultSafetyGate`. Keep HTTP/API, execution, Product Facts, and Workbench projection unchanged.

**Tech Stack:** Go standard library, existing project capability/product packages, YAML file config, JSON schema/fixture validation scripts, local smoke scripts.

---

## Reference Docs

- `docs/fobrain-provider-poc-design.md`
- `docs/fobrain-provider-config.md`
- `docs/adr/2026-07-02-fobrain-provider-poc-boundary.md`
- `docs/fobrain-tool-matrix.md`
- `docs/provider-policy-and-credentials.md`
- `docs/capability-provider-contract.md`
- `docs/facts-contract.md`
- Old read-only reference: `/Users/vick/Desktop/project/ai-agent/internal/business/fobrain`

## Task 1: Catalog Contract Test

- [ ] Create `internal/einoapp/providers/fobrain/catalog_test.go`.
- [ ] Assert `ListCapabilities()` returns exactly one Phase 5 PoC capability.
- [ ] Assert metadata:
  - `ID = tool.fobrain.current_user_context`
  - `ProviderID = fobrain`
  - `ResultSchema = fobrain.tool_result.v2`
  - `RiskLevel = low`
  - `SideEffect = read_external`
  - `PolicyRef = policy:fobrain:read:v1`
  - `CredentialBindingPolicy = required`
  - `ConnectorID = fobrain`
  - `ApprovalRequired = false`
- [ ] Run and confirm the test fails before implementation:

```bash
go test ./internal/einoapp/providers/fobrain -run ProviderCatalog -count=1
```

## Task 2: Provider Catalog Implementation

- [ ] Add `provider.go` and `catalog.go`.
- [ ] Implement `Provider.ID()` and `Provider.ListCapabilities()`.
- [ ] Keep all capability ids and policy refs as package constants in one file.
- [ ] Do not register provider in bootstrap until Task 7.
- [ ] Run:

```bash
go test ./internal/einoapp/providers/fobrain -run ProviderCatalog -count=1
```

## Task 3: Fobrain Config Contract Tests

- [ ] Add bootstrap tests for the `fobrain` YAML contract from `docs/fobrain-provider-config.md`.
- [ ] Test valid local config derives safe connector/client fields.
- [ ] Test invalid or missing `base_url`, `workspace_id`, `timeout`, `credential.status`, and `connector_status.mode`.
- [ ] Test redacted summary never includes `api_token`, Authorization, raw credential ref, URL query secret, or raw provider payload.
- [ ] Test `enabled=false` does not register Fobrain provider.
- [ ] Run and confirm failures:

```bash
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
```

## Task 4: Fobrain Config Implementation

- [ ] Extend `internal/einoapp/bootstrap/config.go` with Fobrain config structs.
- [ ] Validate connector id, workspace id, absolute base URL, timeout, credential status, owner scope, and connector mode.
- [ ] Derive safe `capabilities.CredentialBinding`, safe connector status, and provider client config.
- [ ] Extend `configs/eino-workbench.example.yaml` with a no-secret Fobrain example.
- [ ] Run:

```bash
go test ./internal/einoapp/bootstrap -run 'FobrainConfig|RedactedSummary|CredentialLeak' -count=1
```

## Task 5: Client and Credential Boundary Tests

- [ ] Add `client_test.go` and `credentials_test.go`.
- [ ] Define tests for:
  - missing credential returns `credential_missing` or `credential_scope_denied` safety reason.
  - workspace mismatch does not expose credential ref.
  - mock client is not called when credential/policy blocks execution.
  - provider errors are folded into stable safe reason codes.
- [ ] Run and confirm failures:

```bash
go test ./internal/einoapp/providers/fobrain -run 'Client|Credential' -count=1
```

## Task 6: Client and Credential Implementation

- [ ] Add `client.go`, `credentials.go`, and `errors.go`.
- [ ] Define `FobrainClient` interface with `CurrentUserContext(ctx, request)` returning typed safe result.
- [ ] Define `CredentialResolver` interface scoped by workspace; only provider client receives secret material.
- [ ] Implement mockable HTTP client shell without hardcoded live endpoint in business logic.
- [ ] Run:

```bash
go test ./internal/einoapp/providers/fobrain -run 'Client|Credential' -count=1
```

## Task 7: StructuredResult Mapper Tests

- [ ] Add `result_test.go`.
- [ ] Test successful current-user result maps to `product.StructuredResultCandidate`.
- [ ] Test `fobrain.tool_result.v2` is treated as business/display payload schema, while Product Facts uses `tool.structured_result.v1`.
- [ ] Test candidate passes `product.NewStructuredResultSafetyGate().Approve`.
- [ ] Test raw payload markers, token, Authorization, password, credential, checkpoint, interrupt are rejected.
- [ ] Run and confirm failures:

```bash
go test ./internal/einoapp/providers/fobrain -run 'ReadonlyPoC|StructuredResult|Unsafe' -count=1
```

## Task 8: Invoker Implementation

- [ ] Add `tools.go` and `result.go`.
- [ ] Implement `capabilities.Invoker` for the PoC capability.
- [ ] Dispatch by catalog operation, not by execution-layer tool name branches.
- [ ] Return safe error candidate for unknown capability; do not call client.
- [ ] Run:

```bash
go test ./internal/einoapp/providers/fobrain -run 'Provider|Client|ReadonlyPoC|StructuredResult|Unsafe' -count=1
```

## Task 9: Bootstrap Wiring

- [ ] Register Fobrain provider only from config/catalog path.
- [ ] Keep `configs/eino-workbench.local.yaml` ignored for live token.
- [ ] Add tests proving default startup does not hardcode Fobrain when config omits it.
- [ ] Run:

```bash
go test ./cmd/eino-workbench ./internal/einoapp/bootstrap ./internal/einoapp/capabilities -run 'Fobrain|Capability|Config' -count=1
```

## Task 10: Smoke and Report Contract Tests

- [ ] Add tests or script checks proving `fobrain-poc` produces `docs/schemas/fobrain/provider_poc_report.v1.schema.json` compatible output for mock provider mode.
- [ ] Add tests or script checks proving missing LLM or Fobrain live credential writes `docs/schemas/skip-report.schema.json` compatible output with `blocks_claims`.
- [ ] Add tests or script checks proving real-model report is generated only when both LLM and Fobrain local credentials are configured.
- [ ] Run and confirm failures:

```bash
npm run eino-workbench:schema-test
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config configs/eino-workbench.local.yaml
```

## Task 11: Smoke Scenario Implementation

- [ ] Implement `fobrain-poc` branch in `scripts/eino_workbench_server_smoke.sh`.
- [ ] Add report generation:
  - provider PoC always writes report matching `docs/schemas/fobrain/provider_poc_report.v1.schema.json`.
  - missing LLM or Fobrain live credential writes skip report matching `docs/schemas/skip-report.schema.json`.
  - fully configured real-model run writes report matching `docs/schemas/real-model-report.schema.json`.
- [ ] Assert provider PoC report includes provider mode, capability id, StructuredResult schema, Fobrain business result schema, result ref, safe summary, policy decision, credential binding status, connector status, redaction checks, failure category.
- [ ] Assert real-model report includes prompt, selected tool, sanitized args summary, tool card rendering status, assistant final answer, screenshot path, failure category.
- [ ] Run:

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config configs/eino-workbench.local.yaml
```

## Task 12: Full Verification and Docs Sync

- [ ] Run:

```bash
npm run eino-workbench:schema-test
go test ./... -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario mcp-mock
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config configs/eino-workbench.local.yaml
```

- [ ] Update `docs/07-implementation-plan.md` if any task order changes.
- [ ] Update acceptance records with passed/skipped status.
- [ ] Record remaining non-claims: no 24 read-only recovery, no write approval, no entity clarification, no replacement claim.
