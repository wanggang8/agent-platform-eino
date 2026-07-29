---
title: FOBrain 漏洞处置实验 Agent 实施就绪门禁
status: STORY_1_2_CLEAN_CI_PASS
updated: '2026-07-27'
---

# 实施就绪门禁

## 当前结论

Story 1.1 与工具链门禁已完成。Story 1.2 的单 fixture 只读实现、三态纵向 E2E、重启恢复、安全扫描和固定 Linux/amd64 工具链预检已通过，证据见 [`story-1-2-m1-walking-skeleton-2026-07-27.md`](../../../docs/acceptance-records/story-1-2-m1-walking-skeleton-2026-07-27.md)。Story 1.2 已有对应 clean GitHub Actions run（30446096645）与 commit `62d91a0c1174b1f8de11f70b06993de56f746af0`，可标记为 `done`。

真实 FOBrain 写入保持关闭。M5 的规则与 Gate Evidence 只用于客观裁决，不自动注册生产 capability；真实动作必须在适用门禁通过后另建 M6 Target Story。

## 门禁状态

| 门禁 ID | 必须证实的事实 | 当前状态 | 正式 Story / 通过证据 |
| --- | --- | --- | --- |
| G-TOOLCHAIN | 固定 Go 1.26.x、Node 24 LTS、CI 与 canonical image 可复现 | `PASS` | Story 1.1 `done`；[`story-1-1-g-toolchain-2026-07-15.md`](../../../docs/acceptance-records/story-1-1-g-toolchain-2026-07-15.md) |
| G-ARCH-V2 | Greenfield schema epoch、StructuredResult、Product Facts、projection、SSE、恢复与 import boundary 形成唯一 v2 链路 | 全局 `BLOCKED`；M1 read subset 固定工具链预检 `PASS`，clean CI 待补 | Stories 1.2、1.12；M4 再扩展 Action 控制面 |
| G-READ-01 | 全部新增漏洞精确查询、发现时间倒序，空结果与失败可区分 | 目标部署 API `PASS`；产品 `BLOCKED` | Stories 1.3、1.5、1.12；[`fobrain-target-live-read-evidence-2026-07-14.md`](./research/fobrain-target-live-read-evidence-2026-07-14.md) |
| G-READ-02 | 读取 FOBrain 全部人员安全列表，按内置权限选择 stable identity | 目标部署 API `PASS`；产品 `BLOCKED` | Stories 2.1、2.2；同上只读证据 |
| G-READ-03 | 七项代表性读取具有安全字段、三态缺失、opaque ref→locator 与产品投影 | 接口证据存在；产品 `BLOCKED` | Stories 1.3、1.4、1.6～1.8、1.12 |
| G-FACT-01 | 统一漏洞卡片的行级事实足以支持核对与后续动作 | `BLOCKED` | Stories 1.3、1.8～1.12 |
| G-SAFE-01 | 未确认、取消、拒绝、过期、权限未知均零写入，产品不新增通知 | `BLOCKED` | Stories 2.4～2.7、2.9、2.10、3.2～3.5 |
| G-WRITE-01 | 每个动作的权限、唯一写入、独立回读与错误脱敏 | `BLOCKED` | Stories 2.9、2.10、3.4、3.5 |
| G-WRITE-02 | 批量派发/转发逐条区分完成、待核验、待排查 | `BLOCKED` | Stories 2.9、2.10 |
| G-WRITE-03 | 延时允许状态、期限、管理员权限和双事实回读 | `BLOCKED` | Stories 3.1、3.4 |
| G-WRITE-04 | 误报的 stable owner 权限、目标状态、写入和回读 | `BLOCKED` | Story 3.5 |

## 通过规则

1. 目标部署 API 证据不等于当前产品 capability 通过；必须补齐 StructuredResult、Product Facts、safe projection 和产品 smoke。
2. 门禁失败时只允许补证据、补契约或调整范围，不得硬编码、模拟成功或跳过回读。
3. 缺授权、可恢复样本、独立回读或恢复步骤时，M5 命令必须 `exit 2`、`mutation=0` 并保持 `BLOCKED`。
4. 每个动作独立裁决；派发/转发未通过不阻塞 Epic 3 的 mock 旅程或独立 Gate Evidence。
5. 真实动作只有在其全部适用 `G-WRITE-*`、`G-FACT-01`、`G-SAFE-01`、`G-ARCH-V2`、`G-TOOLCHAIN` 通过并获得授权后，才可转换为新的 M6 正式 Story。
6. M5 Gate Evidence 不得修改或注册生产 FOBrain mutation runtime，只能修改批准的验收脚本、fixture、验收记录和本门禁文件。
7. 任何 v1 adapter、dual write、fallback、自动旧库迁移、旧 API/DOM 兼容均视为架构偏差，不是解除 G-ARCH-V2 的方案。

## 正式 Backlog 与里程碑

| 里程碑 | 正式 Story | 状态 / 允许范围 |
| --- | --- | --- |
| M0 | 1.1 | `done`；G-TOOLCHAIN `PASS` |
| M1 | 1.2 | Story=`done`；实现与固定工具链预检、clean GitHub Actions 均通过，G-ARCH-V2 全局仍 BLOCKED |
| M2 | 1.3～1.8 | 代表性真实读取与统一行事实 |
| M3 | 1.9～1.12 | 快照、引用、恢复与 READ-01 |
| M4 | 2.1～2.8 | 人员选择与 mock Action 控制面；真实 mutation=0 |
| M5 | 2.9～2.10、3.1～3.5 | 规则关闭、mock 延时/误报与隔离 Gate Evidence |
| M6 | 后续转换的新 Story | 仅适用门禁 PASS 后进入真实产品集成 |

当前正式 inventory 为 Epic 1 的 1.1～1.12、Epic 2 的 2.1～2.10、Epic 3 的 3.1～3.5，共 27 个 Story。

## 非执行项门禁

- 历史 Target Candidate：TC-2.1、TC-2.2、TC-2.3、TC-3.1、TC-3.2、TC-3.3。
- 历史 Source Bundle：SB-1.1～SB-1.20、SB-2.6。
- 它们只保留在 [`archive/epics-before-vertical-reslice-2026-07-18.md`](../archive/epics-before-vertical-reslice-2026-07-18.md) 及课程修正记录中，不得进入 sprint status、Create Story 或开发 Agent。
- 门禁通过后不得把 TC 原地改名；必须按 provider/policy/verifier、Product Facts/UX、授权 integration/E2E 的边界生成新 Story。

## 已关闭决策

- OQ-02：由 Story 2.2 使用 FOBrain 全员列表、stable identity、人工搜索与重名消歧关闭。
- OQ-05：由 RD-03 / AD-26 关闭；固定“操作记录”入口归属 Story 2.8，不再归属归档的旧横向 Story。
- Epic 3 只依赖 Epic 2 的共用控制面至 Story 2.8，不依赖 Stories 2.9、2.10 或任何 Target Story。

## 下一步

为 Story 1.2 的当前变更创建提交并推送到 GitHub，已完成 clean `Toolchain Gate` run，已补充 commit SHA / run URL，Story 标记为 `done`。全局 G-ARCH-V2 继续保持 `BLOCKED`，真实 LLM、真实 FOBrain 和所有外部 mutation 继续为 0。
