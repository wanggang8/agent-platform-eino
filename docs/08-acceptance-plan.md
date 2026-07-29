# 验收计划

本文定义当前 27 个正式 Story 的证据等级、里程碑门禁和最终裁决。每个 Story 的 Given / When / Then 与具体命令以 [`epics.md`](../_bmad-output/planning-artifacts/epics.md) 为唯一明细来源。

## 当前状态

- Story 1.1：`done`。
- `G-TOOLCHAIN`：`PASS`，权威证据见 [`story-1-1-g-toolchain-2026-07-15.md`](./acceptance-records/story-1-1-g-toolchain-2026-07-15.md)。
- Story 1.2：实现、三态纵向 E2E、重启恢复、安全扫描和固定 Linux/amd64 工具链预检已通过；证据见 [`story-1-2-m1-walking-skeleton-2026-07-27.md`](./acceptance-records/story-1-2-m1-walking-skeleton-2026-07-27.md)；已生成 clean GitHub Actions（run: https://github.com/wanggang8/agent-platform-eino/actions/runs/30446096645，commit: `62d91a0c1174b1f8de11f70b06993de56f746af0`）。Story 1.2 标记 `done`。
- `G-READ-01/02` 的目标部署 API 证据已存在，但当前产品 capability 仍未完成。
- 真实写域门禁保持 `BLOCKED`；任何 mock 或 Gate Evidence 通过都不自动开启生产入口。

## 证据与裁决规则

| 结果 | 含义 | 是否可推进 |
| --- | --- | --- |
| `PASS` | 在 Story 指定环境运行全部适用命令，AC、事实、安全和退出门禁均通过 | 可以进入下一 Story |
| `FAIL` | 断言失败、数据不一致、泄漏、越权或发生非预期写入 | 不可以 |
| `BLOCKED` | 缺少前置 Story、门禁、授权、可恢复样本或必要外部能力 | 不可以 |
| `exit 2` / `SKIPPED` | 当前阶段未到或环境不具备安全执行条件 | 不是 PASS，不可以作为解除门禁证据 |

共同规则：

1. 最终证据必须来自固定工具链；非 canonical 本地结果只能作为 preflight。
2. 命令、日期、环境摘要和脱敏输出必须落到 `docs/acceptance-records/`。
3. HTTP 2xx、provider accepted、模型自然语言“成功”都不能替代 Product Facts 的独立回读事实。
4. StructuredResult 是唯一事实材料；raw provider payload、凭据和未投影字段不得进入 Workbench、Action API、模型上下文、audit 或 replay。
5. 空结果、无权限、数据不足、部分结果和外部失败必须机器可区分。
6. 修改类路径在未确认、取消、拒绝、过期、门禁未过或权限未知时，外部 mutation 数必须为 0。
7. 所有批量动作逐项产生完成、待核验或待排查；单项失败不得覆盖同批其他结果。

## 里程碑验收矩阵

| 里程碑 | Story | 必须证明 | 退出信号 |
| --- | --- | --- | --- |
| M0 | 1.1 | Go/Node/CI/canonical image 一致，基线可复现 | `G-TOOLCHAIN=PASS`；已完成 |
| M1 | 1.2 | 单 fixture 贯通浅色桌面聊天、Eino mock、StructuredResult、Product Facts、SQLite、SSE 和最小重启恢复 | Story 1.2 所有命令与 AC `PASS`；真实 FOBrain/LLM/Action mutation=0 |
| M2 | 1.3～1.8 | 七项代表性读取、精确新增、全部业务系统和统一行事实具有安全字段与稳定三态 | `G-READ-01/03`、`G-FACT-01` 的适用产品证据通过 |
| M3 | 1.9～1.12 | 冻结分页、确定性引用、SSE/刷新/重启恢复和完整只读旅程 | `READ-01=PASS` |
| M4 | 2.1～2.8 | 全员安全选人、OQ-02、不可变草案、确认/取消、mock 幂等/lease/verifier/reconcile、操作记录 | 全部 M4 Story `PASS`；未批准与真实 provider mutation=0 |
| M5 | 2.9～2.10、3.1～3.5 | 四类动作各自规则、权限、最小写入、独立回读、部分失败与恢复证据 | 每个门禁独立标记 `PASS` 或继续 `BLOCKED`；不产生生产入口 |
| M6 | 转换后的新 Story | 已过门禁动作的 provider/policy/verifier、Product Facts/UX 与授权 E2E | 新 Story 全部 AC 通过后才可暴露相应真实 capability |

## M1 验收重点

Story 1.2 是最薄纵向切片，不是静态原型。至少验证：

- 一个固定安全 fixture 只经过 provider adapter → StructuredResult → Product Facts → product projection。
- Workbench 与 HTTP 读出口消费同一投影，不建立第二套 DTO 真相。
- SQLite 为全新 schema epoch，foreign keys、WAL、bootstrap/readiness fail closed；driver 精确为 `modernc.org/sqlite v1.46.2` 且运行时 SQLite 不低于 `3.51.3`。
- SSE 使用统一 encoder，cursor/reconnect 语义可测试；服务重启后能恢复该单一结果。
- 浅色桌面布局在当前 Story 的目标视口可用；不验移动端、深色或完整 Action UI。
- Eino 使用 mock ChatModel / tool；不访问真实 LLM、真实 FOBrain 或任何 mutation。
- import boundary、配置加载、统一错误 envelope、错误脱敏和日志安全通过。

## M2 / M3 读取门禁

读取通过不能只看接口能返回数据，还必须具有产品链证据：

- `current_user_context` 与 `my_permissions` 使用 stable identity，并区分 resolved、confirmed empty/no permission、source unavailable/unknown。
- 精确新增漏洞查询覆盖完整、空、部分和失败，且按发现时间倒序。
- 资产、漏洞明细只能通过 opaque ref→provider locator 解析，前端与模型不可构造 raw locator。
- `business_list` 不得替换为其他近似能力。
- 统一漏洞行事实至少覆盖当前 Story 规定的 POC、IP、在线状态、业务系统、负责人等字段；缺失值保留稳定三态。
- 查询快照冻结，多次查询通过程序持久化的 result_ref 确定性选择；有歧义必须澄清。
- READ-01 同时覆盖分页 100 条、刷新、SSE 断线和服务重启恢复。

`G-READ-01/02` 的既有目标 API 证据不等于上述产品能力通过。

## M4 Mock Action 门禁

M4 只允许 mock mutation adapter，但必须验证真实控制面语义：

- 全部人员列表只能使用安全字段和 stable identity；重名消歧后选出唯一对象。
- ActionDraft 绑定冻结 result_ref、逐条事实、actor、recipient、policy 与有效期，确认后不可变。
- 未确认、取消、拒绝、过期、权限未知均为零写入。
- 单条只写一次；批量逐项记录 attempt、lease、readback 和 reconcile。
- 响应未知先读后判，不自动重发。
- 派发、主动转发、延时、误报共用 Product Facts 和状态机；不同动作只增加专属参数、角色和 verifier。
- 操作记录从固定入口读取同源逐条结果；OQ-05 由 Story 2.8 实现。

Mock 通过只证明控制面行为，不解除任何 `G-WRITE-*`。

## M5 Gate Evidence 门禁

真实取证必须同时具备：

- 带日期的明确授权、目标环境、动作类型和最大写入数量；
- 可恢复、非生产关键的最小样本；
- 写前基线、唯一 mutation、独立 readback、恢复原状态和恢复后 readback；
- stable actor/owner/recipient identity 与实际权限证据；
- 脱敏错误、部分失败和响应未知的安全处理；
- 每个成功样本 mutation 次数为 1；未授权路径为 0。

Gate Evidence Story 仅可修改 `scripts/acceptance/`、`docs/fixtures/`、`docs/acceptance-records/` 和 `_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md`。无法安全制造的子项保持 `BLOCKED`，不得用 mock 代替目标部署证据。

## 当前门禁总表

| 门禁 | 当前状态 | 解除位置 |
| --- | --- | --- |
| G-TOOLCHAIN | `PASS` | Story 1.1 |
| G-ARCH-V2 | 全局 `BLOCKED`；M1 read subset 固定工具链预检 `PASS`，clean CI 已补齐 | Stories 1.2、1.12 的契约与产品证据 |
| G-READ-01 | 目标 API `PASS`；产品 `BLOCKED` | Stories 1.3、1.5、1.12 |
| G-READ-02 | 目标 API `PASS`；产品 `BLOCKED` | Stories 2.1、2.2 |
| G-READ-03 | 接口证据存在；产品 `BLOCKED` | Stories 1.3、1.4、1.6～1.8、1.12 |
| G-FACT-01 | `BLOCKED` | Stories 1.3、1.8～1.12 |
| G-SAFE-01 | `BLOCKED` | M4 mock 证明控制面；M5 按动作取证 |
| G-WRITE-01/02 | `BLOCKED` | Stories 2.9、2.10、3.4、3.5 |
| G-WRITE-03 | `BLOCKED` | Stories 3.1、3.4 |
| G-WRITE-04 | `BLOCKED` | Story 3.5 |

## 桌面与安全验收

- 首版只验浅色桌面端；最小支持宽度与具体视觉断言以 UX 设计和 Story 要求为准。
- 视觉证据只覆盖当前 Story 已实现的状态，不把旧项目截图当新项目通过结果。
- 键盘操作、焦点、可辨识错误、加载、空、部分结果和恢复状态按当前可见组件验收。
- 日志、SSE、HTTP 错误、截图与验收记录不得包含 token、cookie、连接串、raw provider payload 或个人敏感数据。

## 最终完成规则

只有满足以下条件才能宣布本轮产品实验完成：

1. 27 个正式 Story 按依赖完成，或明确记录哪些 M5 Gate Evidence 因安全条件保持 `BLOCKED`。
2. Epic 1 的 `READ-01=PASS`，M4 的 mock Action 安全门禁通过。
3. 任何真实动作只在其 M6 转换 Story 完成后开放；未转换动作保持关闭。
4. schema、fixture、OpenAPI、generated contract、Go tests、TypeScript typecheck、browser/E2E 和 acceptance records 在受影响范围内一致。
5. 追踪矩阵、实施计划、Sprint 状态、门禁文件和实际代码没有相互矛盾。
