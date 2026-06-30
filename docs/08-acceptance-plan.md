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

验收记录必须明确：

- 是否阻断 P0/P1/P2 声明。
- 是否阻断“重构完成”声明。
- 是否允许替换当前产品基线。
- 关联 report schema 校验是否通过。

## 阶段验收矩阵

| 阶段 | 范围 | 必跑门禁 | 可 skip 项 | 阻断条件 |
| --- | --- | --- | --- | --- |
| P0 | 最小产品闭环，实施 Phase 1-4 | 开发前复核、基础、Contract、前端、服务 smoke `contract/chat-stream/action-basic/tool-card`、安全门禁 | 无 | pre-development validation、contract、视觉、安全、Action API 任一失败 |
| P1 | 产品级运行能力，实施 Phase 5-7 | P0 全部、开发前复核更新、HITL、clarification、Fobrain PoC、real model smoke、projection/replay | 无 Fobrain/LLM 凭据时 real model smoke 可 skip 但必须记录 | approval、clarification、projection、Action/Workbench 同源任一失败 |
| P2 | 既有业务能力恢复，实施 Phase 8 | P1 全部、开发前复核更新、24 只读、connector、credential binding、disambiguation、write approval、live read/write | 无 live 凭据时 live read/write 可 skip 但不能声明能力可比 | Fobrain 恢复门禁未过时不能声明重构完成 |

完整重构目标是 P0 + P1 + P2。P0/P1 只能证明架构和主链可行；未完成 P2/Phase 8 不得声明重构完成、产品能力等价或替换当前产品基线。

## 基础门禁

开发前复核：

```bash
test -f docs/pre-development-validation.md
rg -n "pre-development-validation" docs/acceptance-records docs || true
go list -m github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go list -m -versions github.com/eino-contrib/jsonschema
```

通过标准：

- 本次 Phase 有对应 `acceptance-records/pre-development-validation-YYYY-MM-DD.md`。
- 外部资料访问日期、版本命令摘要、当前架构复核结论均已记录。
- 记录结论允许进入对应 Phase。

```bash
go test ./... -count=1
go vet ./...
go build ./...
```

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

脚本必须自动完成：

- 选择端口。
- 启动服务。
- ready check。
- 创建 run。
- 抽取 run_id。
- 校验 schema。
- 对 `capability-selection`，断言普通自然语言默认进入 ChatModelAgent，显式 action/capability_hint 通过 registry 和 policy，未知或无权限 hint 不会直接执行。
- 对 `context-projection`，断言模型输入上下文只来自 safe Product Facts / StructuredResult，并生成 context snapshot。
- 对 `run-lifecycle`，断言 cancel、stop、timeout、retry 和终态幂等符合 `run-lifecycle.md`。
- 对 `action-consistency`，提交 Action 后用同一 run_id 拉取 run、views/current、replay，并断言 assistant/tool/pending/audit 字段同源。
- 清理进程。

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
go test ./internal/einoapp/... -run 'Telemetry|Budget|ContextSnapshot|RunLifecycle' -count=1
bash scripts/eino_workbench_server_smoke.sh --scenario budget
```

通过标准：

- trace/run/workspace/model/tool count/latency/token usage 至少进入内部 telemetry 或安全报告。
- budget exceeded 产生安全 run notice 或 partial result，可 replay。
- callback 只用于 tracing、metrics、diagnostics，不作为主 SSE 来源。
- audit 与 replay 仍从 Product Facts 投影。

## 真实模型 smoke

环境变量：

```text
EINO_LLM_PROVIDER
EINO_LLM_BASE_URL
EINO_LLM_API_KEY
EINO_LLM_MODEL
FOBRAIN_BASE_URL
FOBRAIN_USER_API_TOKEN
ENABLE_FOBRAIN_CONNECTOR=true
```

命令：

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-poc
```

smoke 脚本负责启动服务、选择端口、运行真实模型 suite、清理进程。无凭据时脚本输出 skipped report，并在 acceptance 记录原因。有凭据时失败即 gate 失败。

报告必须包含：

- prompt。
- selected tool。
- sanitized args summary。
- tool card rendering status。
- assistant final answer。
- screenshot path。
- failure category。

真实模型工具选择必须符合 `intent-and-capability-selection.md`：不按后端关键词路由；资产/漏洞、单 IP 查询/IP 统计、connector status/业务读取、写域审批等场景不得误选。误选、未写入 selected tool 或未写入 safe args summary 均为 failed。

已执行报告必须通过 `schemas/real-model-report.schema.json`，状态只允许 passed/failed。无凭据时不生成 real model report，改为生成 skipped report；skipped report 必须通过 `schemas/skip-report.schema.json`，且 `blocks_claims` 明确标记阻断项。

## Fobrain live 门禁

```bash
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-readonly
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-live-read
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-write-approval
bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-live-write
```

无真环境凭据时必须记录 skip 原因。

通过标准：

- 24 个只读工具矩阵通过。
- `docs/fixtures/fobrain/tool-matrix-24.json` 精确覆盖 24 个只读工具，且已通过 schema-test。
- 工具矩阵 coverage 不低于旧项目最终通过证据中的 `24` 只读工具 + `1` connector、每场景六类截图区域和 `0 failed` 要求。
- connector 状态可展示。
- `credential_binding` 可展示，secret 不泄漏。
- 实体消歧澄清可恢复。
- 写域工单未审批不执行 mutation。
- live read 和 live write 生成 report；skip 时不能声明产品能力可比。
- `docs/fobrain-tool-matrix.md` 中每一行都有 mock/live 断言和截图证据。
- 已执行的 live read/write 报告必须分别通过 `schemas/live-read-report.schema.json`、`schemas/live-write-report.schema.json`，状态只允许 passed/failed。
- 无真环境凭据时不生成 live read/write report，改为生成 `schemas/skip-report.schema.json` 对应的 skipped report。
