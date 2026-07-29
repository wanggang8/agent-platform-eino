# Architecture Course Correction 版本／现实／门禁一致性复评

评审对象：`ARCHITECTURE-SPINE.md` 及 2026-07-15 纠偏后的直接依赖  
评审范围：FR-19/20 三态与 G-READ-03、FR-25 `business_list`、FR-21～24 opaque ref→locator、operation history open attention、Candidate 非 Sprint、真实写门禁、命名版本与仓库现实  
复评日期：2026-07-15  
结论：**PASS — 初审四项现实偏差与最后一项写域门禁摘要偏差均已关闭；当前 Spine 在本评审范围内与唯一 readiness gate、Candidate 转换条件、仓库现实和已记录版本来源一致。**

## 核验方法与边界

- 重新运行 Architecture Spine deterministic lint：通过，0 finding。
- 对照当前 PRD、唯一 `implementation-readiness-gate.md`、Spine、Runtime SPEC、Epics、已批准纠偏提案、live report 与当前仓库代码／lockfile。
- 检查 Sprint heading 结构：正式 `Story` 47 项，`Target Candidate` 6 项；不存在 `### Story ... TC-*` 或原 2.7～2.9／3.6～3.8 Story heading。
- 未重新联网研究无关依赖；版本判断复用既有一手来源记录，并用当前 `go.mod`、`package-lock.json`、本机工具链复核。
- 本次只更新本报告，不修改主文档、代码、schema、fixture 或门禁状态。

## 执行摘要

初审发现的四项问题已经按正确层级修复：

1. FR-19/20 不再把字段缺失冒充已解析；PRD、G-READ-03、Architecture Map 和 Epic Story 已固定 `resolved`、`confirmed_empty/no_permission`、`source_field_unavailable/unknown` 三态。
2. FR-25 已恢复并显式绑定已批准、live passed 的 `tool.fobrain.business_list`，同时禁止替换为当前 blocked 且被范围排除的 `my_business_systems`。
3. AD-27 已固定 opaque `entity_ref` 与独立 `ProviderLocator`，并要求 resolver 重验 workspace/actor/capability/version/object permission/expiry；现有接口 pass 被正确降级为 integration evidence。
4. AD-21/25/26 与 SPEC 已新增 `attention_status=open`，definitive failed 与 manual_attention 均创建 open attention；首版没有关闭命令，操作记录按六个月窗口与 open lifecycle/attention 做后端 union，不再漏掉待排查项。

Candidate 也已从 Sprint parser 结构中隔离：六项真实写能力使用 `### Target Candidate TC-*`，不是 `### Story N.M`；每项均有 `ConversionGates`、重新 readiness、环境授权和“转换为新正式 Story 后才实施”的约束。正式 Story 2.4/2.5/3.4/3.5 被限定为 Gate/Evidence Story，只允许获准环境的最小可恢复取证，不接入生产 Workbench 或启用生产 capability。

最终快速复核确认，Spine 的 “Current Implementation Gates” 收口句现已与唯一 readiness gate 同义：适用 G-WRITE、共同 G-FACT-01、G-SAFE-01、G-ARCH-V2 与 G-TOOLCHAIN 必须全部 PASS，目标部署接口存在或 `[ADOPTED]` 标记都不能替代门禁。具体 Candidate 的 `ConversionGates` 仍包含同一共同前置，因此不存在未来部分门禁转绿时的错误放行路径。

## 最终复核

### C-06 — 真实写域放行条件：Closed

**证据**

- 唯一权威 `implementation-readiness-gate.md:31` 规定：每项写动作必须同时通过适用 `G-WRITE-*`、共同 `G-FACT-01`、`G-SAFE-01`、`G-ARCH-V2`、`G-TOOLCHAIN` 后才可开始实现。
- 同文件 `:46` 继续声明真实写域仍需全部适用门禁、环境授权和可回滚样本。
- Candidate metadata 与转换 AC 正确包含这些共同门禁，例如 TC-2.1/2.2 的 `ConversionGates` 为 G-TOOLCHAIN、G-ARCH-V2、G-FACT-01、G-WRITE-01/02、G-SAFE-01、OQ-02；TC-3.1/3.2 也包含对应全集。
- Spine `ARCHITECTURE-SPINE.md:93` 现已明确：适用 G-WRITE、共同 G-FACT-01、G-SAFE-01、G-ARCH-V2 与 G-TOOLCHAIN 均有 PASS 验收记录后，才能进入真实写域实现。
- `epics.md` 全局规则仍使用“对应门禁 + 重新 readiness”摘要，但六个 Candidate 的机器可消费 `ConversionGates` 均列出 G-TOOLCHAIN、G-ARCH-V2、G-FACT-01、适用 G-WRITE、G-SAFE 与 OQ；转换 AC 又要求重新执行唯一 readiness，因此没有形成实际放行缺口。

**结论**

当 G-WRITE/G-FACT/G-SAFE 先转为 PASS、但迁移或工具链尚未通过时，Spine、唯一 readiness 和 Candidate 现在都会阻止实施。原 H-01 已关闭，无剩余真实写门禁阻断。

## 已关闭的初审发现

### C-01 — FR-19/20 三态与 G-READ-03：Closed

- PRD 明确：部门等源字段未提供时显示稳定 unavailable，不得推断为 resolved；权限／数据范围区分 resolved、confirmed empty/no permission、source unavailable/unknown。
- `implementation-readiness-gate.md` 新增 G-READ-03，状态为“接口证据存在、产品链 BLOCKED”，通过证据包含 schema、fixture、live/mock assertion、locator contract、Product Facts 投影与产品 smoke。
- Spine Current Gates 与 Capability Map 都引用 G-READ-03；不再把 Batch A `passed` 的接口调用等同字段完整或产品完成。
- Epic Story 1.25/1.26 具有对应字段级机器证据和三态 AC。

当前 live report 仍真实反映“当前用户无部门字段、权限字段未返回”，但新基线已把它作为 unknown/unavailable 的输入证据而不是通过信号，因此现实一致。

### C-02 — FR-25 `business_list` 范围：Closed

- PRD FR-25、Spine Capability Map、Epics FR/AR/Story 1.29 均固定 `tool.fobrain.business_list`。
- 文档明确可见范围来自 FOBrain 接口内置权限，不声称“我的业务系统”；禁止替换为 `my_business_systems`。
- 这与批准的 `experiment-scope-proposal.md` 和当前证据一致：Batch C `business_list` passed，Batch B `my_business_systems` blocked/empty 且不在本轮范围。

### C-03 — FR-21～24 opaque ref→locator：Closed

- AD-03/05/06 已把 F-07 纳入 Binds；AD-27 专门固定引用边界。
- 产品可见 `entity_ref` 必须不可逆、不可拆解；真实 provider id、资产类型和路由材料进入独立 `ProviderLocator`。
- resolver 必须重新校验 workspace、actor、source capability、locator version、对象权限与有效期；不得从 ref 字符串恢复 locator。
- SPEC `domain-contracts.md`、AC-32、Epics AR-44 与 Story 1.27/1.28 已同步要求 ref 与 provider ID 不相等、错误类型/越权/过期/unknown fail closed、保留来源与采集时间。
- G-READ-03 保持产品链 BLOCKED，正确承认当前代码仍把清洗后的 provider id 放入 EntityRef；现有 Batch D/E pass 只证明 integration，不冒充目标态已实现。

### C-04 — Operation history open attention：Closed

- AD-21 把 definitive failed 与 manual_attention 都投影为待排查并创建 `attention_status=open`；执行终态不等于 attention 关闭。
- AD-25 禁止 cleanup open attention 引用。
- AD-26 将查询定义为六个月窗口与 Run waiting/running 或 attention open 的去重 union；首版 failed/manual_attention 均持续可发现且没有关闭命令。
- SPEC 进一步固定 `attention_status=not_applicable|open|closed`、首版无 closure command、Asia/Shanghai 六个日历月 inclusive cutoff、稳定后端游标和当前授权整批 fail closed；AC-28/33 覆盖跨窗口、open attention、撤权与分页。

该修正已消除“definitive failed 超过六个月被隐藏”和“terminal manual_attention 何时关闭”两项歧义。

### C-05 — Candidate 非 Sprint 与真实取证边界：Closed

- 结构检查得到 47 个正式 `### Story N.M：`，6 个 `### Target Candidate TC-*：`；原 2.7～2.9、3.6～3.8 不再存在 Story heading。
- Backlog gate 明确 Sprint Planning 只解析 Story heading；TC/SB 永远不是 Sprint 工作项，不能进入 sprint status、Create Story 或开发 Agent。
- Candidate 不允许原地改名，必须重新 readiness 后生成新的正式 Story，并按 provider/policy/verifier、Product Facts/UX、授权 integration/E2E 拆分。
- 所有 Candidate 当前都列出完整 `ConversionGates` 并要求环境授权、可恢复样本、失败联系人和批准记录。
- Gate/Evidence Story 的真实调用仅用于获准环境的最小可恢复取证；任一授权、样本或恢复条件缺失时 `exit 2`、mutation=0，且生产入口/capability 始终关闭。

因此 Candidate parser 风险和“取证等同生产启用”的风险已经关闭。

## 版本与仓库现实复核

### PASS

- Go module language 1.22、Eino 0.9.12、jsonschema 1.0.3、modernc SQLite 1.34.5 与 `go.mod` 一致。
- React 19.2.7、React Router DOM 7.18.1、TanStack Query 5.101.2、Zustand 5.0.14、Vite 8.1.0、Vitest 4.1.9 与根 `package-lock.json` 解析值一致；Radix 仍由 committed lockfile + `npm ci` 固定。
- Go 1.26.x 与 Node 24 LTS 继续明确标为迁移目标并受 G-TOOLCHAIN 阻塞；当前本机 Go 1.23.12 / Node 25.8.1 没有被误报为新阶段通过证据。
- 本轮没有新增未经既有来源核验的第三方技术或版本声明。

## 最终判定

1. Spine 第 93 行已补齐 `G-ARCH-V2` 与 `G-TOOLCHAIN`，并与唯一 readiness gate 同义。
2. 六个 Candidate 的 `ConversionGates` 与重新 readiness 约束保持完整，非 Sprint 结构检查通过。
3. `lint_spine.py` 再次通过，0 finding。

**最终 verdict：PASS。** 当前仍为 BLOCKED 的 G-READ-03、G-FACT、G-WRITE、G-SAFE、G-ARCH-V2 和 G-TOOLCHAIN 是被正确登记的实施门禁，不是架构现实矛盾；不得把本 PASS 解释为真实写入授权或门禁已通过。

## 实际执行的核验命令

```bash
uv run .agents/skills/bmad-architecture/scripts/lint_spine.py --workspace _bmad-output/planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14
rg -c '^### Story [0-9]+\.[0-9]+：' _bmad-output/planning-artifacts/epics.md
rg -c '^### Target Candidate TC-[0-9]+\.[0-9]+：' _bmad-output/planning-artifacts/epics.md
rg -n '^### Story .*TC-|^### Story (2\.[789]|3\.[678])' _bmad-output/planning-artifacts/epics.md
rg -n 'FR-(19|20|21|22|23|24|25)|G-READ-03|business_list|my_business_systems|ProviderLocator|attention_status|ConversionGates|G-WRITE|G-ARCH-V2|G-TOOLCHAIN' _bmad-output docs internal test-results
go version
node --version
sed -n '1,80p' go.mod
jq -r '.packages["node_modules/react"].version, .packages["node_modules/react-router-dom"].version, .packages["node_modules/@tanstack/react-query"].version, .packages["node_modules/zustand"].version, .packages["node_modules/vite"].version, .packages["node_modules/vitest"].version' package-lock.json
```
