# 实施计划

本文把产品蓝图转换为唯一可执行的 M0～M6 纵向交付顺序。正式需求、故事输入输出、验收标准和逐 Story 命令以 [`epics.md`](../_bmad-output/planning-artifacts/epics.md) 为准；本文只定义先后关系、允许范围和退出门禁。

## 当前裁决

- M0 / Story 1.1 已完成，`G-TOOLCHAIN=PASS`。
- Implementation Readiness 已于 2026-07-20 给出 `READY`；Story 1.2 的实现和固定 Linux/amd64 工具链预检已于 2026-07-27 通过，clean GitHub Actions 证据仍待提交后生成，因此 Story 保持 `in-progress`。
- 真实 FOBrain 写入入口保持关闭。Stories 2.9、2.10、3.1、3.4、3.5 只产生规则或 Gate Evidence，不接入生产 mutation runtime。
- 首版只做浅色桌面端，不做移动端、深色模式、自动派发、通知、工单、24 工具齐套或旧系统兼容。

## 唯一推进规则

一个 Story 只有同时满足以下条件才可开始：

1. 前置 Story 全部 `done`。
2. `epics.md` 中列出的 Requirements、Architecture decisions、Inputs / Outputs、Scope、Non-goals 和 Affected directories 已核对。
3. 适用门禁不是 `BLOCKED`；取证 Story 可在隔离范围内验证门禁，但不得把取证视为生产能力已启用。
4. 已按 [`pre-development-validation.md`](./pre-development-validation.md) 生成该 Story 的开工记录。
5. 验收命令存在且失败语义明确；缺授权、缺安全样本或对应阶段未到时必须 `exit 2`，不能伪报通过。

一个 Story 只有全部 Given / When / Then、退出门禁和实际命令通过，文档与 Sprint 状态同步后，才可标记 `done` 并进入下一项。

## M0～M6 路线

| 里程碑 | 正式 Story | 用户可见结果 | 退出门禁 |
| --- | --- | --- | --- |
| M0 工具链 | 1.1 | 所有后续证据来自固定 Go、Node、CI 与 canonical image | `G-TOOLCHAIN=PASS`；已完成 |
| M1 Walking Skeleton | 1.2 | 单一安全 fixture 从聊天入口贯通 Eino、StructuredResult、Product Facts、SSE 和浅色桌面摘要 | Story 1.2 全部 AC 通过；禁止真实 FOBrain、真实 LLM 和 Action |
| M2 代表性真实读取 | 1.3～1.8 | 完成数据可得性、身份权限、新增漏洞、资产、漏洞、业务系统和统一行事实 | `G-READ-01/03`、`G-FACT-01` 的适用读取部分有产品证据 |
| M3 可核对与可恢复读取 | 1.9～1.12 | 百条冻结快照、确定性 result ref、断线/刷新/重启恢复和 READ-01 | `READ-01=PASS`，Epic 1 完成 |
| M4 Mock Action 控制面 | 2.1～2.8 | 安全选人、消歧、不可变草案、确认/取消、单批 Mock 派发/转发及操作记录 | OQ-02、OQ-05 关闭；未批准路径 mutation=0；真实写入口仍关闭 |
| M5 目标规则与 Gate Evidence | 2.9～2.10、3.1～3.5 | 分别验证派发、转发、延时、误报的规则和目标部署证据；延时/误报先完成 Mock 旅程 | 各动作门禁独立裁决；任何 BLOCKED 不阻塞无依赖的其他取证 Story |
| M6 Target Story 转换 | 后续新建 Story | 仅对已通过门禁且获得授权的动作接入生产 provider/policy/verifier、Product Facts/UX 和 integration/E2E | 新 Story 自身验收通过；不得原地改名或把 Gate Evidence 当实现 Story |

## M1：Story 1.2 后端基础不变量

Story 1.2 虽然是最薄 walking skeleton，但必须同时建立后续纵向切片复用的边界，不能用页面 mock 绕过架构：

- `httpapi` 只负责 HTTP 输入校验、统一错误响应、SSE 编码和响应封装。
- `execution` 只负责 Eino 执行、run 生命周期及 checkpoint/resume 端口。
- `facts` 定义 StructuredResult、Product Facts 和 repository interface；它们是唯一事实材料。
- `product` 只从 Product Facts 生成 Workbench / Action API 安全投影。
- `capabilities` 提供 registry、metadata、选择与 adapter 抽象，不按工具名或自然语言关键词硬编码。
- `llm` 固定 provider interface、配置、网络策略、mock provider 和错误脱敏；M1 只使用 mock。
- `observability` 只记录脱敏后的 tracing、metrics、budget 和诊断。
- `store/sqlite` 只实现持久化端口，不反向依赖业务、HTTP、LLM 或 provider。
- SQLite WAL 必须先落实 [`modernc SQLite WAL 安全基线 ADR`](./adr/2026-07-20-modernc-sqlite-wal-safety-floor.md)：driver 精确版本为 `v1.46.2`，运行时 SQLite 不低于 `3.51.3`。
- `providers/*` 只处理 provider 边界；raw payload 不得越过该边界。
- 配置、统一错误响应、SSE encoder、LLM provider interface、capability registry、Product Facts repository 和 product projection interface 必须在进入真实 Eino chat、Action API、真实 LLM 或业务 provider 前固定。

M1 只允许一个 fixture、一个只读 capability、一个桌面页面和一个结果。不得提前实现生产 provider、真实写入、通用 Action 状态机或多能力目录。

## M2：代表性读取切片

按 1.3→1.8 顺序逐项交付，每个 Story 都必须把 provider 返回转换为 StructuredResult，再进入同一 Product Facts 投影：

1. Story 1.3：先验证七项代表性读取、精确新增队列和统一行事实的数据可得性，缺失字段固定三态语义。
2. Story 1.4：当前用户与权限；不得根据聊天声明或显示名扩大权限。
3. Story 1.5：全部新增漏洞，按发现时间倒序，区分完整、空、部分和失败。
4. Story 1.6：按 IP 查询资产，并通过 opaque ref 查看安全详情。
5. Story 1.7：按 IP 查询漏洞，并通过 opaque ref 查看安全详情。
6. Story 1.8：查询 `business_list` 并补齐统一漏洞行事实；不得替换为 `my_business_systems`。

M2 不承诺 24 个 FOBrain 工具齐套，不创建真实 Action 草案或 mutation 入口。

## M3：冻结、引用、恢复与 READ-01

- Story 1.9 固定快照边界和每页 100 条的完整核对体验。
- Story 1.10 由程序保存并解析 `result_ref → snapshot/query/selection`，模型只表达意图；歧义时必须澄清，不能猜。
- Story 1.11 让刷新、SSE 断线和服务重启恢复同一安全事实与 run 状态。
- Story 1.12 对 Epic 1 的全部代表性读取、投影、引用、恢复和浅色桌面体验做 READ-01 E2E。

只有 `READ-01=PASS` 才可进入 Epic 2。

## M4：Mock Action 控制面

M4 只使用安全 mock mutation adapter，按 2.1→2.8 建立可复用控制面：

- 人员列表和唯一 stable identity 选人；重名必须消歧。
- 基于冻结事实生成不可变草案；批量只合并人工确认属于同一接收人的对象。
- 所有修改先展示逐条事实并要求人工确认；取消、拒绝、过期和未确认均为零写入。
- 单条与批量逐项结算为完成、待核验、待排查；响应未知先读后判，不自动重发。
- 派发和主动转发共用状态机、审计、verifier、reconcile 和操作记录，不复制第二套事实。
- 固定“操作记录”入口在 Story 2.8 激活；其数据仍来自同一 Product Facts。

## M5：隔离取证

Stories 2.9、2.10、3.1、3.4、3.5 的写域活动只允许修改：

- `scripts/acceptance/`
- `docs/fixtures/`
- `docs/acceptance-records/`
- `_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md`

它们不得修改或注册生产 FOBrain mutation runtime。缺少授权、可恢复样本、独立回读或恢复步骤时，命令必须 `exit 2` 且 `mutation=0`。

Stories 3.2、3.3 仅使用 M4 控制面和 mock adapter 完成延时、误报旅程。Epic 3 的前置只到 Story 2.8，不依赖 2.9、2.10 或任何 Target Story。

## M6：真实动作转换

每种动作独立满足适用门禁后，重新运行 Create Story，生成新的正式 Story ID，并至少拆分为：

1. provider / policy / verifier；
2. Product Facts / Workbench 投影；
3. 获授权 integration / E2E。

未转换前，生产 capability registry、Workbench 和 Action API 不得暴露真实派发、转发、延时或误报入口。

## 每项交付记录

每个 Story 完成时必须同步：

- `_bmad-output/implementation-artifacts/sprint-status.yaml`
- 该 Story 的实现记录
- 受影响 schema、fixture、OpenAPI、ADR 和矩阵
- `docs/acceptance-records/` 下的实际命令与结果
- 若门禁状态变化，更新 `implementation-readiness-gate.md` 和验收计划

开发 Agent 不得根据本文跳过 `epics.md` 的逐 Story Acceptance Criteria。
