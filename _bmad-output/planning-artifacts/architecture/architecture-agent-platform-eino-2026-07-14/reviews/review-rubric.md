# Architecture Spine 独立 Rubric 评审

- 评审对象：`../ARCHITECTURE-SPINE.md`
- 评审方法：`bmad-architecture/references/reviewer-gate.md` 的 Good-spine checklist
- 评审日期：2026-07-14
- 评审性质：独立验证；未修改 Spine、PRD、UX、代码或其他文档

## Gate Verdict

**CHANGES REQUIRED（不通过交付门禁）。** Spine 已形成清晰的 Eino-first / Product Facts 主干，机械 lint 也通过，但仍有 1 个 Critical 和 4 个 High 级语义缺口：它没有继承当前“禁止外部写域实现”的实施门禁，使用了与既有 Product Facts 契约不一致的 Run 状态词汇，未将关键 PRD 动作约束落为可执行不变量，Eino 的六边形边界与当前代码现实不一致，并且 initiative 级运维/环境包络仍未闭合。上述问题会让下一级 epic/story 在都声称遵循本 Spine 时仍产生互不兼容或越过安全门禁的实现。

## 评审范围与证据

重点核对：

- `ARCHITECTURE-SPINE.md` 与伴随 `.memlog.md`
- `product-blueprint/prd.md`、`implementation-readiness-gate.md`、`requirements-traceability-matrix.md`
- 最终 UX `ux-designs/.../DESIGN.md` 与 desktop-only ADR
- `docs/README.md`、`04-technical-architecture.md`、`05-contract-design.md`、`06-security-and-projection.md`
- `facts-contract.md`、`run-lifecycle.md`、`approval-flow.md`、`conversation-context.md`
- `07-implementation-plan.md`、`08-acceptance-plan.md`、`provider-policy-and-credentials.md`、`observability-and-budgets.md`
- 当前 `internal/einoapp` 分层、import boundary、Product Facts 模型、capability metadata、`go.mod`、workspace `package.json` / `package-lock.json`

机械检查：

```text
uv run .agents/skills/bmad-architecture/scripts/lint_spine.py \
  --workspace _bmad-output/planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14
=> ok=true, total_findings=0

go test ./internal/einoapp/architecture -count=1
=> PASS
```

## Good-spine Checklist

| Checklist | 结论 | 说明 |
| --- | --- | --- |
| 固定本层真实分歧点且无重大遗漏 | **Fail** | 写域准入、动作基数/分组/授权、状态词汇及 Eino seam 仍可分叉。 |
| 每条 AD 的 Rule 可执行且确实阻止所述分歧 | **Partial** | AD-03、AD-05、AD-07～AD-11、AD-13～AD-16 较强；AD-12 状态词汇与现契约冲突，AD-18 的“可验证备份”未定义责任与门禁。 |
| Deferred 中没有会让下一级当前工作分叉的事项 | **Pass with caution** | 已列 Deferred 多数有明确触发条件；但数据保留、迁移/恢复、健康门禁并未进入 Deferred。 |
| 命名技术已验证为当前且适配 | **Partial** | 版本基本忠实反映 lockfile/go.mod，Eino 0.9.12 是当前最新稳定线；但部分前端 patch 已更新、modernc SQLite 明显落后当前稳定线，Radix 只写“1.x”，未留下统一的 ratification/升级策略。 |
| 忠实 ratify brownfield 代码库而不矛盾 | **Fail** | AD-12 与现有 Run schema 冲突；Eino 被画成独立 adapter，但当前多个核心/基础包直接依赖 Eino 类型。 |
| 覆盖驱动它的 PRD 能力 | **Fail** | F-01～F-06 有粗粒度映射，但 FR-02、FR-05/06/09/12/14 的结构性约束未落入 Rule。 |
| 不削弱继承的 parent spine | **N/A** | 本次未发现 parent spine 继承关系。 |
| initiative 所有维度均 decided/deferred/open | **Fail** | 部署形态已定，但数据生命周期、schema migration/restore、运行健康与启动恢复责任仍沉默。 |

## Critical / High Findings

### C-01 — Spine 未继承写域实施阻塞门禁

- **等级：Critical**
- **处置：Discuss，然后补为 Inherited Invariant / 明确 Gate；不可仅留给实施者自行理解。**
- **证据：** PRD 状态仍是 `BASELINE_APPROVED_PENDING_IMPLEMENTATION_GATE`；`implementation-readiness-gate.md` 明确“只允许进入只读事实链实现，不允许进入外部写入实现”，并要求 G-FACT-01、G-WRITE-01～04、G-SAFE-01 按动作通过。Spine 的 sources 未包含该门禁，却将 AD-07～AD-11、写域 sequence 与 F-02～F-05 映射全部标为 `[ADOPTED]`，没有任何规则说明这些是目标态设计而非当前实施授权。
- **为何会分叉：** 一个 story 团队可依据 Spine 直接实现 `ConfirmAction → external mutation`；另一个团队会依据 readiness gate 只实现只读事实链。两者都能声称遵循自己的上游材料，但前者越过项目的安全准入。
- **建议收口：** 在 Spine 明确继承 `implementation-readiness-gate.md`：架构可定义写域目标形态，但在每项适用的 G-FACT/G-WRITE/G-SAFE 通过前，禁止实现或启用外部 mutation；未授权环境只允许 schema、fixture、mock 与只读事实链。将门禁加入 sources/capability map，并说明通过状态由哪份记录裁决。

### H-01 — AD-12 的 Run 状态词汇与现有契约和代码冲突

- **等级：High**
- **处置：Discuss；二选一 ratify 现契约，或以破坏性契约变更 ADR 明确迁移。**
- **证据：** AD-12 要求持久化 `queued/running/waiting_approval/resuming/reconciling/terminal`。现有 `facts-contract.md`、`run-lifecycle.md`、`eino_product_facts.v1.schema.json`、`eino_run_snapshot.v1.schema.json` 和 `facts.RunStatus` 一致使用 `created/running/waiting/succeeded/failed/cancelled/stopped`；approval/clarification 通过 `PendingInteraction.kind/status` 表达，不拆成 Run enum。Spine 自己又在 Structural Seed 中将 `PENDING_INTERACTION` 建为独立事实，造成两套表达竞争。
- **为何会分叉：** 后端 story 可能新增 Run enum，前端/Schema story 继续消费 `waiting` + Pending，replay/SSE/Action API 随即出现不兼容状态机。
- **建议收口：** 最小变更是明确 AD-12 中的词是“恢复分类/阶段”而非 Run.status，并绑定现有 Run + Pending + ActionItem/Attempt 的组合表达；如确需扩展 Run.status，先新增契约 ADR、schema 版本迁移、repository migration 与前端生成契约计划，不能以 `[ADOPTED]` 静默覆盖现状。

### H-02 — PRD 的动作基数、分组、候选来源与对象级授权未成为可执行不变量

- **等级：High**
- **处置：Autofix 可行，但应先确认承载位置（AD Rule、capability policy/schema 或两者）。**
- **证据：** PRD 要求：批量确认默认前 5 条且确认前可查看完整冻结集合（FR-02）；接收人必须来自 FOBrain 全部人员列表（FR-05/09）；只有同一接收人才可合并，不同接收人必须拆分（FR-06/09）；直接延时每次只能一条且仅管理员执行（FR-12）；误报仅限操作人确认负责的漏洞（FR-14）。AD-07～AD-10 只规定通用 draft、人工输入、二阶段确认和逐条执行；没有约束 draft 的 `max_items`、`group_key`、recipient source、actor role/object predicate 或确认投影完整性。Capability map 只映射到 F 级，无法证明这些 FR 已被覆盖。
- **为何会分叉：** 两个能力实现都可满足现有 AD，却分别选择“一个草案混合多个接收人”和“按接收人拆草案”，或让延时批量执行；policy 实现也可能只校验 workspace 可见性而漏掉管理员/自己负责谓词。
- **建议收口：** 固定由 capability contract/policy 声明并由 Prepare + Confirm 双检的结构性字段，例如 `max_items`、`group_by/recipient_ref`、`recipient_source_capability`、`actor_role_predicate`、`object_permission_predicate`、`confirmation_projection=preview_5+full_frozen_set`；为四项动作明确规则映射和失败语义，未知权限继续 fail closed。

### H-03 — “Eino 是 adapter”的范式与当前依赖现实没有形成唯一 seam

- **等级：High**
- **处置：Discuss；必须选择并写清一种依赖形态。**
- **证据：** Design Paradigm 和图把 Eino 标为指向 `execution` 的独立 adapter，并声明依赖箭头只指向核心端口；Structural Seed 却没有 Eino adapter 包，反而写 `execution/ # Eino 编排`。当前 `execution` 直接 import `eino/adk`、callbacks、model、schema，`capabilities` 直接 import Eino tool/schema，`observability` 直接 import callbacks/model，`store/sqlite` 直接实现 Eino compose checkpoint 接口。现有 import boundary 只保证 `httpapi/facts/product/providers/fobrain` 不直接依赖 Eino，并未把 `execution` 本身隔离成纯核心。
- **为何会分叉：** 新 story 可能新建一套 Eino adapter/port，把现有 runner 迁出 execution；另一个 story 会沿用当前 convention 在 execution/capabilities/observability/store 中直接使用 Eino 公共类型。两者都会认为自己符合 Spine。
- **建议收口：** 若 ratify 当前现实，应将 Eino 定义为“受控 execution kernel dependency”，精确列出允许直接依赖 Eino 的包，并保持 Product Facts、product/httpapi/provider DTO Eino-free；图中不要再暗示不存在的独立 adapter。若坚持纯六边形 adapter，则必须命名 adapter 包、项目自有 port 和迁移边界，并解释 checkpoint/tool/callback 类型如何不越界。

### H-04 — initiative 级 operational/environmental envelope 仅部分决定

- **等级：High**
- **处置：可将明确项 Autofix；其余进入 Deferred/Open Items 并给出触发条件。**
- **证据：** AD-18 决定单实例 SQLite/WAL/备份，AD-12 决定持久化 run/租约，AD-19 决定日志与配置边界；但 Spine 未决定或 defer：数据库 schema migration 的所有者与启动失败策略、备份恢复验证门禁、不可变 facts/audit 与过期 snapshot/draft/checkpoint 的保留/清理策略、waiting checkpoint 的清理豁免、executor 启动恢复与 lease 失效的唯一所有者、health/readiness 对“存储可写/迁移完成/provider 不可用”的语义、开发/验收/受控实验环境的配置与写域开关边界。
- **为何会分叉：** store、execution、bootstrap 和运维 story 可各自实现 cleanup/recovery/health，造成 waiting approval 被误删、未迁移实例接流量、备份不可恢复或未知外部写被盲重试。
- **建议收口：** 至少固定 ownership 和不可违反的安全行为；具体保留时长可配置或 Deferred，但必须有“谁清、什么绝不能清、何时 readiness=false、恢复如何验证”的规则。生产 SLA 可继续列为非目标，不妨碍先封闭实验环境的安全运维语义。

## Medium Findings

### M-01 — Stack 对“verified-current / pinned”的表达不够一致

- **等级：Medium**
- **处置：Autofix。**
- **证据：** `go.mod` 与 lockfile 证实 Spine 中 Go/Eino/jsonschema/SQLite/React/Router/Query/Zustand/Vite/Vitest 版本是仓库实际解析版本。2026-07-14 通过 Go module proxy / npm registry 核验：Eino `v0.9.12` 仍是最新稳定（`v0.10.0-alpha.*` 为预发行），jsonschema `v1.0.3` 最新；React/Router/Query/Zustand 仍一致；Vite 已有 `8.1.4`、Vitest `4.1.10`，四个 Radix 包也有更新 patch；modernc SQLite 已到 `v1.53.0`。旧 pin 并非自动错误，但 Spine 只解释了 Eino 的保留理由，且 `Radix UI Primitives = package-pinned 1.x` 不是可复现版本，实际是四个不同包的 lockfile 版本。
- **建议收口：** 将 Stack 明确为“ratified repository baseline as of 2026-07-14”，以 `go.mod`/`package-lock.json` 为唯一精确版本源；Radix 列出包集合或直接引用 lockfile。对 modernc/Vite/Vitest 说明“有更新但本轮不顺带升级”的兼容验证/ADR 触发条件，避免把“当前”误解为“registry latest”。

### M-02 — AD-13 声明的 metadata 超出现有 registry/schema，迁移边界未说明

- **等级：Medium**
- **处置：Autofix 或明确为 target delta。**
- **证据：** AD-13 要求 metadata 声明 `freshness` 与 `verifier`；当前 `capabilities.Capability` 和 `capability_catalog.v1.schema.json` 已有 schema/risk/side effect/approval/timeout/credential/idempotency，但没有 freshness/verifier。该新增方向合理，却被 `[ADOPTED]` 表述成已存在的共同契约。
- **建议收口：** 明确这是需新增的 catalog/schema 版本与生成契约变更，并绑定迁移/验收；或把 freshness/verifier 作为稳定引用而非运行时函数，避免 provider 私有 verifier 穿透 registry。

## 通过项与值得保留的主干

- AD-03 将 StructuredResult → Safety Gate → Product Facts → 全部产品出口固定为单一事实链，和项目文档、现有分层一致。
- AD-05～AD-08 对 result snapshot、`complete_set`、不可变 ActionDraft、两阶段写协议的组合能有效防止分页、上下文与确认后重查造成的范围漂移。
- AD-10～AD-11 对逐条回读、部分成功、不确定写入进入 reconciling/manual attention 的规则符合 PRD 的“不得伪报成功”。
- AD-13～AD-16 对 registry、provider 不可信边界、单操作人 actor 来源和安全前端投影的方向清楚。
- AD-17 与最终 UX 和 desktop-only ADR 一致；AD-18 已避免实验期微服务化及多实例共享 SQLite 的错误假设。
- Deferred 大多给出合理触发条件，没有把自动派发、通知、多人身份或移动端偷渡回当前范围。

## 建议的最小再评审门槛

重新提交 Reviewer Gate 前至少应满足：

1. 写域 readiness gate 成为显式继承约束，且不会被 `[ADOPTED]` 目标态规则误读为实施授权。
2. Run/Pending/Action 状态词汇与现有 schema 统一，或存在正式破坏性迁移 ADR。
3. FR-02、FR-05/06/09/12/14 的结构性约束能从 AD/capability policy 直接追踪并自动验收。
4. Eino 依赖 seam 与当前 import boundary 只剩一种解释。
5. 数据生命周期、migration/restore、startup recovery/readiness 至少被决定或带触发条件地 Deferred。

完成以上项目后，再运行 deterministic lint、import-boundary test，并逐项复核 Good-spine checklist。

---

## Recheck — 2026-07-14

### 定向结论

**PASS。** 本次只复核原报告的 C-01、H-01、H-02、H-03、H-04；五项均已关闭。该结论不重新评估或扩展 Medium/Low 议题。

| 原发现 | 状态 | 关闭证据 |
| --- | --- | --- |
| C-01 写域实施阻塞门禁未继承 | **Closed** | Spine 已将 `implementation-readiness-gate.md` 加入 sources，并新增 `Current Implementation Gates`：明确架构目标态不构成实施授权，G-FACT-01、G-WRITE-01..04、G-SAFE-01 未 PASS 前不得实现或启用真实 mutation；`[ADOPTED TARGET]` 也被明确定义为受迁移 ADR/readiness gate 阻塞。Capability map 同步增加“写域实施准入”。迁移 ADR 第 12、116、142 行重复固定该边界；`approval-flow.md` 同步声明门禁前不得启用真实外部 mutation。 |
| H-01 Run 状态与现有契约冲突 | **Closed** | AD-12 已保留现有 `created/running/waiting/succeeded/failed/cancelled/stopped` 七态，并把 `reserved/sent/verifying/reconciling/manual_attention` 限定到 ActionItem/Attempt。AD-21 给出唯一聚合映射；迁移 ADR“状态映射”和 `run-lifecycle.md` 同步采用同一规则，避免新增竞争 Run enum。 |
| H-02 PRD 动作约束未成为不变量 | **Closed** | AD-09 已固定并要求 registry/policy 在 prepare/confirm 双检 `max_items`、`group_key`、`recipient_source_capability`、`actor_role_predicate`、`object_permission_predicate` 与确认投影；同时逐项写明全员列表来源、按接收人拆 draft、延时单条且管理员、误报“当前 actor 负责”、多条前 5 条 + 完整冻结集合。AD-13 与 Catalog v2 迁移说明为这些约束提供版本化承载边界。 |
| H-03 Eino adapter/seam 与代码现实冲突 | **Closed** | Paradigm 已改为“受控 Eino execution-kernel dependency”，图与文字不再声称存在独立 Eino adapter。AD-02 精确白名单 `execution/capabilities/observability/store/sqlite` 可直接 import Eino，并要求 `facts/product/httpapi/providers/*/web` Eino-free；该规则与当前 import boundary 和代码依赖现实一致。迁移 ADR只迁移自然语言 Runner/HITL 行为，不另造竞争 adapter 层。 |
| H-04 operational/environmental envelope 未闭合 | **Closed** | AD-22 固定 staged/retained checkpoint、确认事务顺序、attempt reservation 与 epoch/fencing lease；AD-25 固定 SQLite connection policy、启动 migration/integrity/可写/恢复扫描、readiness fail closed、active facts 禁止 cleanup、terminal facts/audit 首版不物理删除、备份 restore 验证。迁移 ADR进一步固定 WAL 备份范围、retention、`/health`/`/ready` 语义及开发/验收/受控实验环境隔离。 |

### Recheck Evidence

```text
uv run .agents/skills/bmad-architecture/scripts/lint_spine.py \
  --workspace _bmad-output/planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14
=> ok=true, total_findings=0

go test ./internal/einoapp/architecture -count=1
=> PASS
```

**仍未关闭的 C-01/H-01..H-04：无。**
