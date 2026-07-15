# 验收计划

本文只写怎么证明做对了。

本文是阶段、合并和完整重构的权威验收入口。`07-implementation-plan.md` 中的命令只作为任务级检查。

验收记录统一写入：

```text
docs/acceptance-records/
test-results/eino-*
```

skip 必须写入验收记录，包含命令、原因、缺少的环境变量或凭据、后续补跑条件。

验收记录模板见 `acceptance-records/TEMPLATE.md`。工具链、视觉基线和报告路径见 `tooling-and-reporting.md`。

旧项目最终验收证据只作为范围和视觉基线参考，索引见 `legacy-acceptance-evidence.md`。任何 P0/P1/P2 通过声明都必须来自新项目当次生成的 `test-results/eino-*` 报告，不能复用旧项目历史通过结论。

Story 1.1 起，本文后续列出的 contract、Go、TypeScript、browser、checkpoint、SQLite、boundary 和
smoke 命令是 `scripts/run_toolchain_baseline.sh` 的逻辑组成项；最终验收必须从 clean GitLab
pipeline 产出的 canonical linux/amd64 `tag@sha256` image 进入并顺序执行。宿主或本地运行只可记录为
preflight，不得拆分拼接为阶段 PASS。当前首发视觉门禁仅为 desktop，mobile 不属于
`G-TOOLCHAIN`。

验收记录必须明确：

- 是否阻断 P0/P1/P2 声明。
- 是否阻断“重构完成”声明。
- 是否允许替换当前产品基线。
- 关联 report schema 校验是否通过。

## 阶段验收矩阵

| 阶段 | 范围 | 必跑门禁 | 可 skip 项 | 阻断条件 |
| --- | --- | --- | --- | --- |
| M-0 | Story 1.1 固定可复现工具链 | exact declarations、负向 verifier、canonical image digest、clean pipeline、完整 baseline、批准的 desktop visual evidence | 无 | 任一证据缺失即 `G-TOOLCHAIN=BLOCKED`，不得进入 M-1 |
| P0 | 最小产品闭环，实施 Phase 1-4 | 开发前复核、基础、Contract、前端、服务 smoke `contract/chat-stream/action-basic/capability-selection/context-projection/tool-card`、安全门禁 | 无 | pre-development validation、contract、视觉、安全、Action API 任一失败 |
| P1 | 产品级运行能力，实施 Phase 5-7 | P0 全部、开发前复核更新、HITL、clarification、Fobrain PoC、real model smoke、projection/replay | 无 Fobrain/LLM 凭据时 real model smoke 可 skip 但必须记录 | approval、clarification、projection、Action/Workbench 同源任一失败 |
| P2 | 既有业务能力恢复，实施 Phase 8 | P1 全部、开发前复核更新、24 只读、connector、credential binding、disambiguation、write approval、live read/write | 无 live 凭据时 live read/write 可 skip 但不能声明能力可比 | Fobrain 恢复门禁未过时不能声明重构完成 |

完整重构目标是 P0 + P1 + P2。P0/P1 只能证明架构和主链可行；未完成 P2/Phase 8 不得声明重构完成、产品能力等价或替换当前产品基线。

## G-TOOLCHAIN 门禁

唯一裁决算法：

```text
PASS = exact declarations
    AND verifier negative tests
    AND canonical image build/push digest
    AND clean GitLab pipeline
    AND reproducible install/no dependency drift
    AND schema/contract/OpenAPI
    AND Go/test/race/vet/build/checkpoint/SQLite/boundary
    AND TS/Vitest/stream/build
    AND approved desktop visual baseline
    AND contract smoke
otherwise BLOCKED
```

canonical image identity 必须同时证明 Node linux-x64 tar.gz checksum、Ubuntu Noble snapshot、
source isolation、`build-essential`、`command -v cc` 和构建期最小
`CGO_ENABLED=1 go test -race ./...`。缺少任一 pinned 输入或真实 image 构建证据时，宿主 race
PASS 与静态 Dockerfile 检查都不能替代门禁条件。

本地 visual migration 必须先从同一 pinned inputs 取得真实 image ID，在该 image 运行唯一完整
baseline；若首次停在 desktop screenshot diff，必须保留 baseline log 和 actual/diff，并确认所有前置
非视觉检查通过。UX 明确批准后才可在同一 image 更新 desktop snapshot，并再次完整运行唯一 baseline。
本地结果始终只是 preflight/visual migration evidence。

Story 最终证据必须来自包含全部变更的 clean commit：GitLab pipeline push canonical image，输出
`CI_COMMIT_SHA`、`CI_PIPELINE_URL`、`TOOLCHAIN_IMAGE=tag@sha256`，并保留
`test-results/toolchain-baseline.log` 与 `test-results/eino-workbench-playwright-report/`。缺 remote、
pipeline、registry digest、artifact 或 UX approval 时不得填写模拟值，也不得形成中间“部分 PASS”。

2026-07-15 裁决为 `G-TOOLCHAIN=BLOCKED`：本地 canonical image、非视觉 baseline 与 71 张 desktop
候选已产生，但 UX approval 和批准后的完整 baseline 尚未发生；仓库无 Git remote，未产生 GitLab
lint/pipeline/clean `CI_COMMIT_SHA`、registry `tag@sha256` 或 artifacts。Story 1.1 和 Sprint 1.1
必须保持 `in-progress`，M-1 不得开始。详见
`acceptance-records/story-1-1-g-toolchain-2026-07-15.md`。

## 基础门禁

开发前复核：

```bash
test -f docs/pre-development-validation.md
rg -n "pre-development-validation" docs/acceptance-records docs || true
go list -m github.com/cloudwego/eino
```

通过标准：

- 本次 Phase 有对应 `acceptance-records/pre-development-validation-YYYY-MM-DD.md`。
- 外部资料访问日期、版本命令摘要、当前架构复核结论均已记录。
- 记录结论允许进入对应 Phase。
- 尚未引入的集成模块只在对应 Phase 前调研版本：`eino-ext/components/model/openai` 放在 Phase 4，`github.com/eino-contrib/jsonschema` 放在 contract generator 首次真实引入前。

```bash
go test ./... -count=1
go vet ./...
go build ./...
```

## Backend Foundation 门禁

Phase 2.3 stream reducer、Phase 3 Eino chat、Action API、真实 LLM provider 和业务 provider 实现前，必须先通过后端基础门禁：

```bash
go test ./internal/einoapp/bootstrap ./internal/einoapp/httpapi ./internal/einoapp/llm ./internal/einoapp/capabilities ./internal/einoapp/facts ./internal/einoapp/product -run 'Config|Response|SSE|Provider|Registry|Facts|Projection|Redaction' -count=1
go test ./internal/einoapp/architecture -run ImportBoundary -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario contract
```

通过标准：

- 配置覆盖 server、database、LLM、security、observability、timeout、budget，并且脱敏输出不包含 token、credential、Authorization、API key 或连接密钥。
- HTTP 成功响应保持 OpenAPI 业务 schema 直出；错误响应统一 `eino_error_envelope.v1`，包含 request id、安全错误码和可展示摘要。
- SSE 编码统一处理 `id`、`event`、`data`、flush、content type、no-cache 和编码失败。
- LLM 只通过 `llm.Provider` 接口接入 execution；provider 错误必须脱敏。
- 工具注册和选择只通过 `capabilities.Registry`、metadata 和 policy；不得按自然语言关键词或工具名硬编码生产分支。
- Workbench、ActionResult、SSE、Replay、Inspector 只能从 Product Facts / product projection 读取事实，不直接消费 execution event 或 provider payload。
- 不存在通用 `utils` 包承载跨层逻辑。

## Phase 3 技术门禁

进入 Eino chat、Product Facts 持久化、SSE 或 Action API 编码前，必须先通过 Phase 3 技术门禁：

```bash
rg -n 'phase-3-execution-and-facts-plan|Product Facts|ResumeWithParams|StructuredResult' docs
go list -m github.com/cloudwego/eino
go test ./internal/einoapp/architecture -run ImportBoundary -count=1
```

通过标准：

- `docs/phase-3-execution-and-facts-plan.md` 已明确 facts model -> SQLite repository -> product projection -> execution event mapper -> ChatModelAgent Runner -> HTTP/SSE/Action API 的顺序。
- Eino 版本、Runner API、resume/checkpoint 边界和 callback 诊断边界与 `go.mod` 锁定版本一致。
- 当前 import boundary 仍阻止 HTTP、product、facts、Fobrain provider 层直接依赖 Eino 包；provider raw payload 外泄由 schema、projection、redaction 和 provider 安全测试覆盖。

## Contract 门禁

```bash
npm run eino-workbench:schema-test
npm run eino-workbench:contract-generate
npm run eino-workbench:contract-test
```

通过标准：

- `docs/api/eino-workbench.openapi.json` 可 lint。
- `docs/schemas/` 和 `docs/schemas/fobrain/` 中所有 schema 是合法 JSON Schema 2020-12。
- 所有 fixtures 通过 JSON Schema 2020-12 校验。
- TypeScript contract 与 schema 等价。
- server 成功响应、错误响应、SSE payload 均可校验。
- `eino_action_result.v1` 覆盖 approval_refs/resume_refs 矩阵，且禁止 `resume_token`、`credential_ref` 和 raw provider body。

## 前端门禁

```bash
npm run eino-workbench:typecheck
npm run eino-workbench:test
npm run eino-workbench:stream-test
npm run eino-workbench:browser-test
npm run eino-workbench:visual-test
```

通过标准：

- 三栏 shell 存在。
- reducer 处理重复、乱序、reconnect、tool patch、pending patch。
- desktop/mobile screenshot baseline 通过。
- baseline 更新必须有 diff、mask 说明和验收记录。
- `visual-acceptance-matrix.md` 中 relevant block 都有 selector、fixture screenshot、desktop/mobile evidence。
- 视觉 coverage 不低于 `legacy-acceptance-evidence.md` 中最终通过的场景、状态和截图区域要求。
- visual report 输出到 `test-results/eino-workbench-visual-report/`。
- `docs/fixtures/visual-evidence-matrix.json` 必须通过 schema-test，且 8 个必选 block 和状态 fixture 引用完整。
- 产品视觉截图不得出现英文调试文案、tool id、schema 名、raw JSON、run id 或 provider 原始字段；可见内容必须是处理后的中文产品展示。
- 覆盖 desktop 1440x900 与 mobile 390x844 的空态、聊天、工具卡、审批卡、澄清卡、失败态、Inspector。
- 动态字段必须 mask：时间、run id、随机 id、模型耗时、token 计数。

## 服务 smoke

未到对应 Phase 的 scenario 必须由 smoke 脚本返回 `exit 2` 和明确的 intentionally unavailable 提示；这种结果不是通过信号，也不得写入验收记录作为 pass。

Phase 1 当前只要求：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario contract
```

后续 Phase 逐步打开：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario chat-stream
bash scripts/eino_workbench_server_smoke.sh --scenario action-basic
bash scripts/eino_workbench_server_smoke.sh --scenario capability-selection
bash scripts/eino_workbench_server_smoke.sh --scenario context-projection
bash scripts/eino_workbench_server_smoke.sh --scenario run-lifecycle
bash scripts/eino_workbench_server_smoke.sh --scenario action-consistency
bash scripts/eino_workbench_server_smoke.sh --scenario tool-card
bash scripts/eino_workbench_server_smoke.sh --scenario clarification
bash scripts/eino_workbench_server_smoke.sh --scenario replay
```

Phase 3 完成后，`chat-stream`、`action-basic`、`capability-selection` 和 `context-projection` 不得再返回 `exit 2`；必须启动服务、创建 run、读取 snapshot/stream 或 SQLite Product Facts，并验证消息、Action、capability 选择和模型上下文都从 Product Facts 同源读写。

脚本必须自动完成：

- 选择端口。
- 启动服务。
- ready check。
- 创建 run。
- 抽取 run_id。
- 校验 schema。
- 对 `capability-selection`，断言普通自然语言默认进入 ChatModelAgent，显式 action/capability_hint 通过配置驱动的 registry 和 policy，未知 hint 在创建 run 前被拒绝，需要审批的 hint 不会触发 runner 或 context snapshot。
- 对 `context-projection`，断言模型输入上下文只来自 safe Product Facts / StructuredResult，并生成不含 raw、credential、token、checkpoint、interrupt 的 context snapshot。
- 对 `tool-card`，断言 mock read capability 经 Eino tool loop 写入 ToolCall、Safety Gate 后的 ToolResult、Workbench 工具卡、ActionResult result card、SSE tool patch 和 audit。
- 对 `run-lifecycle`，断言 cancel、stop、provider timeout、pending timeout、`schema_invalid` retry 和终态幂等符合 `run-lifecycle.md`；当前后端 smoke 使用 SQLite seed 准备 running/waiting Product Facts，验证 lifecycle API、ActionResult、run snapshot、replay/audit、SQLite pending 状态和 retry 幂等，不声明前端 RunNotice 或 provider timeout retry 已完成。
- 对 `action-consistency`，提交 Action 后用同一 run_id 拉取 ActionResult、run、views/current、replay、SSE、audit refs 和 SQLite facts，并断言当前 mock read 场景中的 tool/audit/result card 字段同源、错 workspace 访问返回 404，且产品出口不泄漏 raw provider payload。
- 对 `replay`，断言 replay view 精确等于同 run 的 Workbench snapshot 投影，replay events 同时包含由 Product Facts snapshot 派生的产品事件和安全 audit event。
- 清理进程。

Phase 4/P0 完成后，`tool-card` 不得再返回 `exit 2`；它必须使用本地 mock capability provider，不依赖真实模型凭据。`real-model-chat` 归属 P1 真实模型 smoke：只有在 Phase 4.3 引入 OpenAI-compatible provider 后才能打开；无本地凭据时必须生成 skipped report，不能作为通过信号。

Phase 7 当前进展：`action-consistency`、`replay` 和 `budget` 已不再允许返回 `exit 2`；它们必须启动真实本地服务并通过同源 Product Facts 断言。`budget` 当前覆盖最小 `budget_exceeded` lifecycle：安全失败、active tool 取消、waiting pending 过期、audit 写入和 replay 投影；内部 telemetry/counter、ChatModel Eino callback handler、OpenAI-compatible token usage 映射、只统计型 cost estimate、内部 run summary 和 telemetry summary JSON report 已接入安全 sink。OTel exporter、workspace quota、rate limit 和限制型 cost budget 仍由后续硬化门禁覆盖。assistant/pending 的完整同源验收仍由 `chat-stream`、HITL、clarification 和后续 replay 硬化门禁共同覆盖。

## HITL 门禁

```bash
go test ./internal/einoapp/execution -run 'ApprovalInterrupt|ResumeDuplicate|ResumeReject|ResumeCheckpointMissing|ResumeAfterRestart' -count=1
go test ./internal/einoapp/execution -run 'ClarificationRequested|ClarificationSubmit|ClarificationCancel|ClarificationDuplicate|ClarificationAfterRestart' -count=1
go test ./internal/einoapp/store/sqlite -run CheckPointStore -count=1
npm run eino-workbench:stream-test -- --grep pending
bash scripts/eino_workbench_server_smoke.sh --scenario clarification
```

通过标准：

- 未审批不执行 mutation。
- approve 后继续执行。
- reject 后不可 approve。
- duplicate resume 幂等。
- checkpoint 丢失返回安全错误。
- 进程重启后可恢复。
- approval 状态机、ActionResult waiting、SSE pending patch 和 replay/audit 必须符合 `approval-flow.md`。
- clarification request、submit、cancel、duplicate、restart 后恢复均可测。
- 澄清状态机、ActionResult waiting、SSE pending patch 和 replay/audit 必须符合 `clarification-flow.md`。

Phase 6.2 当前进展：后端 approval interrupt/resume 已覆盖 requested、approved、rejected、cancelled、duplicate approve/cancel、reject/cancel 后不可 approve、checkpoint missing 和 restart resume；pending UI/SSE 的产品化展示、clarification resume 和完整 replay 视觉验收仍必须在后续任务单独通过，不能由 Phase 6.2 自动声明完成。

Phase 6.3 当前进展：后端 clarification submit/cancel 已覆盖安全 resume data、duplicate submit、cancel、expired 后拒绝、checkpoint missing 防消费和 restart submit；Pending UI/SSE 的产品化展示与完整 replay 视觉验收仍必须在 Phase 6.4/Phase 7 单独通过。

Phase 6.4 当前进展：已实现 `pending.updated` 前端 reducer、approval/clarification waiting 与终态只读中文卡片、终态 pending SSE 不暴露旧 `resume_ref`、以及 `clarification` smoke 门禁。真实 UI 按钮提交/取消请求绑定、Playwright 视觉截图验收和完整 replay 视觉验收仍需后续任务覆盖。

## 安全门禁

```bash
go test ./internal/einoapp/... -run 'Safety|Leak|Redaction' -count=1
npm run eino-workbench:stream-test -- --grep safety
```

通过标准：

- assistant 不泄漏 raw prompt、tool id、reason code、model instruction。
- tool result 不泄漏 Authorization、API key、raw credential_ref。
- ActionResult 不泄漏 resume token。
- audit/replay 不泄漏 raw provider body。
- provider config、network blocked、provider non-2xx、malformed response 和 timeout 必须符合 `llm-provider-safety.md`。
- credential missing、credential scope denied、connector status 和 mutation approval 必须符合 `provider-policy-and-credentials.md`。
- context snapshot 不包含 raw prompt、raw provider body、credential、checkpoint id、interrupt id。
- telemetry/callback 不直接进入 Workbench SSE，且不记录 secret 或 raw provider payload。

## Observability / Budget 门禁

```bash
go test ./scripts/telemetry_summary_report ./internal/einoapp/observability -run 'TelemetrySummary|Report|Summary|Cost' -count=1
go run ./scripts/telemetry_summary_report --output test-results/eino-workbench-telemetry-usage-summary-report.json
node scripts/eino_workbench_report_validate.mjs --schema docs/schemas/telemetry_usage_summary_report.v1.schema.json --report test-results/eino-workbench-telemetry-usage-summary-report.json
go test ./internal/einoapp/... -run 'Telemetry|Budget|ContextSnapshot|RunLifecycle' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario budget
```

通过标准：

- trace/run/workspace/model/tool count/latency/token usage、只统计型 `estimated_cost_microunits`、run summary 和 telemetry summary report 至少进入内部 telemetry 或安全报告；telemetry summary report 必须由 `scripts/telemetry_summary_report` 生成并通过 `docs/schemas/telemetry_usage_summary_report.v1.schema.json` 校验。
- budget exceeded 必须通过 execution 层写入 Product Facts，产生 `budget_exceeded` 安全失败摘要，取消 active tool，关闭 waiting pending，并可 replay。
- replay events 必须包含安全 `event_type=budget` audit event，且不得泄漏 raw prompt、secret、Authorization、API key、provider payload 或 reusable resume token。
- 内部 telemetry 事件、run summary 和 telemetry summary report 不得进入 Product Facts、Workbench SSE、Action API、Replay 或 audit；字段只能包含安全 label、latency、token/tool count、只统计型 cost estimate 和 failure category。
- 预算阈值必须来自配置或 policy，评估器只能返回安全 decision，不能直接绕过 lifecycle 写 run 状态。
- callback 只用于 tracing、metrics、diagnostics，不作为主 SSE 来源。
- audit 与 replay 仍从 Product Facts 投影。

Phase 7 closeout 通过条件：

- `action-consistency`、`replay`、`budget` 三个本地 smoke 必须全部通过。
- telemetry summary report 必须生成并通过 schema 校验。
- 任务级检查必须覆盖 Product projection/replay Go 测试、observability/report Go 测试、安全边界 Go 测试、`go test ./...`、contract-test 和 `git diff --check`。
- 不得把 OpenTelemetry exporter、workspace quota、rate limit、限制型 cost budget 或持久化 facts cursor/event log replay 硬化声明为已完成。

Phase 7 closeout 权威验收命令：

```bash
go test ./internal/einoapp/product -run 'Projection|Inspector|Replay' -count=1
go test ./scripts/telemetry_summary_report ./internal/einoapp/observability -run 'TelemetrySummary|Report|Summary|Cost' -count=1
go run ./scripts/telemetry_summary_report --output test-results/eino-workbench-telemetry-usage-summary-report.json
node scripts/eino_workbench_report_validate.mjs --schema docs/schemas/telemetry_usage_summary_report.v1.schema.json --report test-results/eino-workbench-telemetry-usage-summary-report.json
go test ./internal/einoapp/... -run 'Audit|Safety|Leak|Redaction|Telemetry|Budget|ContextSnapshot|RunLifecycle|ImportBoundary|Usage|Cost|Summary|Report' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario action-consistency
bash scripts/eino_workbench_server_smoke.sh --scenario replay
bash scripts/eino_workbench_server_smoke.sh --scenario budget
go test ./...
npm run eino-workbench:contract-test
git diff --check
```

## 真实模型与 Fobrain PoC smoke

真实模型和 Fobrain smoke 使用本地 ignored 配置文件，不使用环境变量作为服务配置来源：

```text
configs/eino-workbench.local.yaml
```

配置文件必须包含最小 LLM provider 字段、Fobrain connector、timeout 和 budget。Fobrain live 配置必须包含认证参数名，当前真实环境为 `credential.auth_param=authorization`，真实 token 只放在 ignored local 配置的 `credential.api_token`。LLM `credential_binding`、`model_label` 和 `network_safety` 由系统从 `provider/base_url/api_key/model/timeout` 派生；该文件可包含本地密钥，但不得提交；验收记录只能写脱敏摘要。

命令：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc --config configs/eino-workbench.local.yaml
bash scripts/eino_workbench_server_smoke.sh --scenario real-model-chat --config configs/eino-workbench.local.yaml
```

smoke 脚本负责启动服务、选择端口、运行真实模型 suite、清理进程。无凭据时脚本输出 skipped report，并在 acceptance 记录原因。有凭据时失败即 gate 失败。`real-model-chat` 只验证模型 provider 边界和错误脱敏；Fobrain 工具选择仍由 `fobrain-poc` / Phase 5+ 门禁覆盖。

`real-model-chat` provider-only 报告必须包含：

- prompt。
- provider kind、model label 和脱敏配置摘要。
- assistant final answer。
- provider error category 或 success status。
- redaction checks。

`fobrain-poc` provider PoC 报告必须通过 `schemas/fobrain/provider_poc_report.v1.schema.json`，且必须包含：

- provider mode。
- capability id。
- StructuredResult schema。
- Fobrain business result schema。
- result ref。
- safe summary。
- policy decision。
- credential binding status。
- connector status。
- redaction checks。
- failure category。

`fobrain-poc` 真实模型工具选择报告只在 LLM 和 Fobrain 本地凭据齐全时生成，必须包含：

- prompt。
- selected tool。
- sanitized args summary。
- tool card rendering status。
- assistant final answer。
- screenshot path。
- failure category。

`real-model-chat` 已执行报告必须通过 `schemas/real-model-provider-report.schema.json`，状态只允许 passed/failed。

Fobrain 真实模型工具选择必须符合 `intent-and-capability-selection.md`：不按后端关键词路由；资产/漏洞、单 IP 查询/IP 统计、connector status/业务读取、写域审批等场景不得误选。误选、未写入 selected tool 或未写入 safe args summary 均为 failed。

`fobrain-poc` provider PoC 报告状态只允许 passed/failed。真实模型工具选择报告必须通过 `schemas/real-model-report.schema.json`，状态只允许 passed/failed。无 LLM 或 Fobrain 凭据时不生成 real model report，改为生成 skipped report；skipped report 必须通过 `schemas/skip-report.schema.json`，且 `blocks_claims` 明确标记阻断项。

## Fobrain live 门禁

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-batch-a --config configs/eino-workbench.local.yaml
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-readonly
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-live-read
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-write-approval
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-live-write
```

无真环境凭据时必须记录 skip 原因。

通过标准：

- 24 个只读工具矩阵通过。
- `docs/fixtures/fobrain/tool-matrix-24.json` 精确覆盖 24 个只读工具，且已通过 schema-test。
- Batch A live/smoke 报告必须通过 `schemas/fobrain/batch_a_live_report.v1.schema.json`，覆盖 `connector.fobrain.security`、`tool.fobrain.current_user_context`、`tool.fobrain.my_permissions` 三项能力。
- 工具矩阵 coverage 不低于旧项目最终通过证据中的 `24` 只读工具 + `1` connector、每场景六类截图区域和 `0 failed` 要求。
- connector 状态可展示。
- `credential_binding` 可展示，secret 不泄漏。
- 实体消歧澄清可恢复。
- 写域工单未审批不执行 mutation。
- live read 和 live write 生成 report；skip 时不能声明产品能力可比。
- `docs/fobrain-tool-matrix.md` 中每一行都有 mock/live 断言和截图证据。
- 已执行的 live read/write 报告必须分别通过 `schemas/live-read-report.schema.json`、`schemas/live-write-report.schema.json`，状态只允许 passed/failed。
- 无真环境凭据时不生成 live read/write report，改为生成 `schemas/skip-report.schema.json` 对应的 skipped report。

阶段性记录：

- Phase 8.2 Batch A 业务只读视觉切片已为 `tool.fobrain.current_user_context` 和 `tool.fobrain.my_permissions` 生成新项目 Playwright baseline：2 个工具 × 6 个区域 × desktop/mobile = 24 张截图。
- Phase 8.3 已为 `connector.fobrain.security` 生成新项目 Playwright baseline：1 个 connector × 6 个区域 × desktop/mobile = 12 张截图。
- 当前记录合计覆盖 Batch A 三项能力和 Batch E 四个只读工具的 Workbench 视觉证据；24 个只读工具全量视觉矩阵、Batch B-D 视觉、live read/write、实体消歧和写域审批仍按 Phase 8 后续任务验收。
- Batch E 必须额外验收 `docs/fobrain-batch-e-interface-plan.md`：资产详情 `network_type` 归一化和路径分流、漏洞详情非分页响应、业务风险 count POST body、威胁关联 relevance list 参数映射、live 样本 ID 来源脱敏和 StructuredResult 唯一事实边界。
- Batch E mock/catalog/HTTP mapper 代码切片记录见 `docs/acceptance-records/phase-8-batch-e-mock-http-2026-07-02.md`。该记录不允许单独作为 live pass、视觉通过或 24 只读恢复完成声明。
- Batch E live smoke 记录见 `docs/acceptance-records/phase-8-batch-e-live-smoke-infra-2026-07-03.md`。该记录已覆盖 discovery sidecar、report schema、smoke 脚本和真实脱敏 live pass。
- Batch E 视觉、fresh replay 和 audit evidence 记录见 `docs/acceptance-records/phase-8-batch-e-visual-replay-audit-2026-07-03.md`。该记录只允许声明 Batch E 四个只读工具的 fixture-based 视觉切片通过，不代表 Batch B-D、Action API 真实服务同源 smoke 或 Fobrain 24/24 恢复完成。
- Phase 8.1 只读矩阵初始对账记录见 `docs/acceptance-records/phase-8-readonly-tool-matrix-gap-2026-07-06.md`。Batch B/C mock/catalog 完成后，mock provider catalog 已覆盖 connector + 24 个只读工具；live、视觉和 Action API 同源仍未完成。
- Batch B mock/catalog 记录见 `docs/acceptance-records/phase-8-batch-b-mock-catalog-2026-07-06.md`。该记录只允许声明 Batch B provider mock/catalog 可用，不代表 live、视觉、Action API 同源 smoke 或 Fobrain 24/24 恢复完成。
- Batch B HTTP mapper 记录见 `docs/acceptance-records/phase-8-batch-b-live-mapper-2026-07-06.md`。该记录只允许声明 Batch B live mapper 有本地 focused tests，不代表真实分页 live pass、视觉、Action API 同源 smoke 或 Fobrain 24/24 恢复完成。
- Batch B live smoke/report 基础设施记录见 `docs/acceptance-records/phase-8-batch-b-live-smoke-infra-2026-07-06.md`。该记录只允许声明 Batch B report schema、smoke runner 和 server smoke 入口可用；没有真实 passed 报告前，不代表 Batch B live pass、视觉、Action API 同源 smoke 或 Fobrain 24/24 恢复完成。
- Batch B 真实 live smoke 阻断记录见 `docs/acceptance-records/phase-8-batch-b-live-smoke-blocked-2026-07-06.md`。当前真实配置下，当前用户上下文缺少可解析部门，且“我的范围”查询返回空结果，因此不得声明 Batch B live pass。
- Batch C mock/catalog 记录见 `docs/acceptance-records/phase-8-batch-c-mock-catalog-2026-07-06.md`。该记录只允许声明 Batch C provider mock/catalog 可用，不代表 live、视觉、Action API 同源 smoke 或 Fobrain 24/24 恢复完成。
- Batch C HTTP mapper 记录见 `docs/acceptance-records/phase-8-batch-c-live-mapper-2026-07-06.md`。该记录只允许声明 Batch C live mapper 有本地 focused tests，不代表 Batch C 5 工具 live pass、视觉、Action API 同源 smoke 或 Fobrain 24/24 恢复完成。
- Batch C live smoke/report 基础设施记录见 `docs/acceptance-records/phase-8-batch-c-live-smoke-infra-2026-07-06.md`。该记录只允许声明 Batch C report schema、smoke runner 和 server smoke 入口可用；没有真实 live passed 报告前，不代表 Batch C 5 工具 live pass、视觉、Action API 同源 smoke 或 Fobrain 24/24 恢复完成。
- Batch C 真实 live smoke 历史阻断记录见 `docs/acceptance-records/phase-8-batch-c-live-smoke-blocked-2026-07-06.md`。该记录证明当前真实环境 `pending_tickets` 返回空结果；后续用户已确认先跳过工单读取。
- Batch C 工单跳过记录见 `docs/acceptance-records/phase-8-batch-c-skip-pending-ticket-2026-07-06.md`。用户确认当前真实环境没有待处理工单功能数据，`pending_tickets` 在 Batch C live smoke 中标记为 `skipped`，Batch C 5 工具 live pass 只覆盖其余 5 个工具；Fobrain 24/24 恢复完成声明仍被工单读取暂缓阻断。
