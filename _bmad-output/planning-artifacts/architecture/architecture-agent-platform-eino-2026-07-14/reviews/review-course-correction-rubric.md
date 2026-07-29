# Course Correction Good-spine 独立复评（最终重跑）

评审对象：`ARCHITECTURE-SPINE.md` 及当前 PRD、Runtime SPEC、UX、Epics、readiness gate、traceability 与实施／验收文档  
评审日期：2026-07-15  
评审意图：Architecture Reviewer Gate / Good-spine rubric walker  
本轮重点：前次 `reserved_not_sent` cancel/stop 阻断的关闭情况，以及此前已收口的 semantic claim、verification、operation history、entity identity/locator/history policy 与 Story metadata。

## Gate verdict

**PASS。** 前次唯一 High 已通过 draft 级全有或全无事务闭合，并传播到 Architecture、Runtime SPEC domain union/state machines/acceptance、ADR、lifecycle 与正式 Story。当前剩余 **0 High、0 Medium**；未发现仍会让独立实现单元产生不兼容状态、写入或产品投影的 Good-spine 阻断。

该 PASS 表示架构材料已经具备一致的下游构建基线，不表示真实写域门禁已自动放行。真实 capability 仍必须由适用 `G-WRITE-*` 与共同 `G-FACT-01/G-SAFE-01/G-ARCH-V2/G-TOOLCHAIN` 的实际 PASS 证据裁决。

## 自动检查结果

```text
lint_spine.py: ok=true, total_findings=0
unresolved source/companion paths: 0

formal Story headings: 47
duplicate formal Story ids: 0
Target Candidate headings: 6
Source Bundle headings: 21
missing/duplicate Requirements fields: 0
FR coverage: 25/25
NFR coverage: 4/4 (NFR-01..04)
UX-DR coverage: 30/30
metadata namespace errors: 0
AD-19 formal Story binding: present
```

定点检索未发现旧的 deadline／预算措辞、NFR-05..10 metadata、`Architecture` 字段中的 G-* 引用，或仍只允许 `reserved_not_sent → sent` 的孤立权威状态图。

## 前次 High 关闭证据

### H-01 — `reserved_not_sent` 在 Run cancel/stop 下没有合法退出路径：CLOSED

当前唯一契约为：

1. Confirm 后的 cancel/stop 是 **draft 级全有或全无事务**，不是逐项部分取消。
2. 只有当该 draft 的全部 attempts 仍为 `reserved_not_sent` 时，execution 才能在同一事务逐项 phase CAS 为 `cancelled_before_send`、记录 `mutation_count=0`、释放全部 active claims，并保留 audit tombstone/lineage。
3. 用户 cancel 固定 `Run/ActionResult=cancelled/cancelled`；系统或用户 stop 固定 `stopped/stopped`；两者 outcome 均为 `none_succeeded`。
4. 任一 attempt 的 sent CAS 已提交时，整个 cancel/stop 返回 `action_already_sent`，不修改任何 item/claim，不伪报 Run 终态；后续只执行既有发送或 Verify/Reconcile。
5. `reserved_not_sent → sent` 与 `reserved_not_sent → cancelled_before_send` 使用互斥 CAS；SQLite 事务与 fencing token 保证只有一边提交。

传播与可验收性均已满足：

- AD-11 固定事务边界、claim release 例外、Run/ActionResult/outcome 和 sent+ 拒绝语义。
- AD-23 与 Runtime SPEC `ActionAttempt` 封闭 union 增加 `cancelled_before_send`；该 phase 禁止 sent-only/provider/evidence 字段，必含取消决定、零 mutation 与 claim release audit ref。
- Runtime SPEC item/claim 状态图、聚合表与 Confirm crash-recovery 图采用相同迁移；不存在 reserved orphan 或伪造 sent 字段的合法路径。
- `docs/run-lifecycle.md` 与迁移 ADR区分未确认 cancel、pre-send cancel/stop、sent+ `action_already_sent` 和受控进程退出释放 lease。
- AC-29 验证 cancel CAS 与 sent CAS 竞态、claim release/tombstone、零 mutation 与不可同时成功；AC-34 验证 phase union；Story 1.22/1.23 绑定同一竞态与崩溃恢复要求。

## 其余定点复核

| 检查项 | 结论 | 说明 |
| --- | --- | --- |
| 真实写域实施门禁 | **PASS** | Spine 引用适用 `G-WRITE-*` 和全部共同门禁；Epics 引用各 Candidate 的完整 `ConversionGates`，目标态不替代 readiness 证据。 |
| 待核验／待排查时间边界 | **PASS** | Architecture outcome table、PRD、SPEC reducer/AC、UX journey/state patterns、Story、ADR 与 lifecycle 都以 deadline/预算和 conclusiveness 得到唯一映射。 |
| Semantic mutation identity | **PASS** | `smk.v1` 固定 JCS canonical payload、NFC、UTC RFC3339 纳秒、missing/null、set 排序与版本升级规则；claim 有唯一主键、并发获取、failed retry lineage、pre-send release 例外与 retention。 |
| ActionAttempt phase union | **PASS** | `reserved_not_sent`、`cancelled_before_send`、`sent`、核验与 sent-terminal 字段均由 phase 判别；禁止零值伪造和旧 `reserved` enum。 |
| VerificationEvidence reducer | **PASS** | evidence kind、control certainty、provider acceptance、observed 条件字段、conclusiveness、版本、调用序号与唯一 reducer 已固定；provider adapter 不决定产品 outcome。 |
| Operation history snapshot | **PASS** | 首次请求物化不可变安全行；snapshot 独立于 Run-local sequence；六个月 month-clamp、授权整批 fail closed、稳定排序及防篡改 cursor 已固定，AC-33/37/39 提供 oracle。 |
| Entity subject / locator / history policy | **PASS** | random immutable subject、active/retained fingerprint alias、opaque ref、LocatorRepository/Resolver、key rotation、policy ref/version、损坏与越权 fail closed、retention 均有唯一规则，AC-38/39 可验收。 |
| Story metadata | **PASS** | 47 个正式 Story 的七字段、FR-01..25、NFR-01..04、UX-DR-01..30、namespace 与 AD-19 binding 均通过；Candidate/Source Bundle 不会被正式 Story parser 捕获。 |
| OQ-05 / operation history ownership | **PASS** | OQ-05 已 CLOSED；共用记录能力属于 Epic 1 Story 1.36，Epic 2/3 只依赖 Epic 1。 |
| Source Map | **PASS** | Sprint Change Proposal、Runtime SPEC、Epics 与所有 companion 相对路径均可解析。 |

## Good-spine checklist

| 检查项 | 结论 | 说明 |
| --- | --- | --- |
| 固定下一级真实分叉点 | **PASS** | write boundary、claim、cancel/sent 竞态、verification、history、identity 与 locator 均有唯一合同和失败路径。 |
| AD Rule 可执行并阻止 stated divergence | **PASS** | 关键规则包含事务边界、封闭枚举、字段条件、错误码、聚合映射和 acceptance oracle。 |
| Deferred 不让单元隐式分叉 | **PASS** | deferred 范围均有触发条件，未把当前关键决策推给实现单元。 |
| 命名技术与版本治理 | **PASS WITH IMPLEMENTATION GATES** | 当前兼容线、lockfile 与迁移目标分离；真实实施继续受 readiness gate 约束。 |
| Ratify brownfield 现实 | **PASS** | `[ADOPTED TARGET]`、迁移 ADR 与 `G-ARCH-V2` 清楚区分目标态和当前实现。 |
| 覆盖 PRD/SPEC/UX/Epics | **PASS** | 核心能力、状态与安全边界已同步到下游正式验收材料。 |
| Parent spine 继承 | **N/A** | initiative altitude，无 inherited parent spine。 |
| Initiative 维度完整 | **PASS** | 部署、环境、存储、恢复、备份、身份、安全、API、前端、运维与迁移均有决定或 defer。 |

## 结论与后续门禁

Architecture Good-spine reviewer gate 可记录为 **PASS**。下一步应按流程重跑 Implementation Readiness，并继续以 gate registry 的实际状态决定哪些 Candidate 可以转正式 Sprint Story 或启用真实写域；不得把本报告当作 `G-WRITE-*`、`G-SAFE-01` 或环境授权的替代证据。

本复评仅更新本报告，未修改 PRD、Architecture Spine、SPEC、UX、Epics、gate、trace、代码或其他主文档。
