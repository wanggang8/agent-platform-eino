# Eino-first Agent Workbench 文档入口

本项目按当前产品蓝图重新开发，不迁移或兼容旧 runtime。正式范围是 FOBrain 漏洞处置实验 Agent 的 27 个纵向 Story；首版只做浅色桌面端。

## 当前权威顺序

发生冲突时按以下顺序裁决：

1. [`product-blueprint/prd.md`](../_bmad-output/planning-artifacts/product-blueprint/prd.md)：产品需求与范围。
2. [`ARCHITECTURE-SPINE.md`](../_bmad-output/planning-artifacts/architecture/architecture-agent-platform-eino-2026-07-14/ARCHITECTURE-SPINE.md)：架构决策与安全不变量。
3. [`DESIGN.md`](../_bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/DESIGN.md) 与 [`EXPERIENCE.md`](../_bmad-output/planning-artifacts/ux-designs/ux-agent-platform-eino-2026-07-14/EXPERIENCE.md)：浅色桌面 UX 规范。
4. [`_bmad-output/planning-artifacts/epics.md`](../_bmad-output/planning-artifacts/epics.md)：27 个正式 Story、前置、输入输出、命令与验收标准。
5. [`07-implementation-plan.md`](./07-implementation-plan.md)：M0～M6 推进顺序。
6. [`08-acceptance-plan.md`](./08-acceptance-plan.md)：证据等级、里程碑门禁与完成裁决。
7. [`pre-development-validation.md`](./pre-development-validation.md)：每个 Story 的开工复核。

Sprint 状态以 [`sprint-status.yaml`](../_bmad-output/implementation-artifacts/sprint-status.yaml) 为准，门禁状态以 [`implementation-readiness-gate.md`](../_bmad-output/planning-artifacts/product-blueprint/implementation-readiness-gate.md) 为准。

## 当前里程碑

| 里程碑 | Story | 范围 |
| --- | --- | --- |
| M0 | 1.1 | 固定工具链；已完成 |
| M1 | 1.2 | 单 fixture 只读 walking skeleton；固定工具链预检通过，clean GitHub Actions 待完成 |
| M2 | 1.3～1.8 | 代表性真实读取与统一行事实 |
| M3 | 1.9～1.12 | 冻结快照、引用、恢复、READ-01 |
| M4 | 2.1～2.8 | 人员选择与 mock Action 控制面 |
| M5 | 2.9～2.10、3.1～3.5 | 规则关闭、mock 延时/误报与隔离 Gate Evidence |
| M6 | 后续转换的新 Story | 仅已过门禁动作的真实产品接入 |

## 不可变原则

- Eino-first；Eino 负责通用执行，产品事实与安全边界由项目掌握。
- Workbench、Action API、audit、replay 共用 Product Facts。
- StructuredResult 是唯一事实材料，raw provider payload 不出 provider 边界。
- 工具选择由 capability metadata、intent 和 policy 驱动，不按工具名或关键词硬编码。
- 写动作必须经过 policy、人工确认、audit、唯一写入、独立回读和恢复门禁。
- Greenfield 实现，不复制或兼容旧 runtime、数据库、API、执行链或 UI DOM。

## 本轮明确不做

- 移动端与深色模式。
- 自动派发配置、产品新增通知、工单流程。
- FOBrain 24 工具齐套或旧产品能力等价声明。
- M5 Gate Evidence 未转换为新 Story 前的生产写入口。

## 历史与参考文档

`docs/00-*`～`docs/06-*`、`docs/09-*`、旧 Phase/P0/P1/P2 计划、旧 FOBrain 24 工具矩阵以及 2026-06/07 的 Phase 验收记录，均是前期架构探索或旧实现验收基线。它们可以只读用于追溯设计原因和安全边界，但不是当前 Story inventory、实施顺序或完成裁决。

旧项目 `/Users/vick/Desktop/project/ai-agent` 同样只读参考，不得修改，不得把其通过记录当成新项目通过结果。

## 工作规则

- 开发前必须先取得 Implementation Readiness=`READY`，再按 Story 执行开发前复核。
- 每次行为变更同步 schema、fixture、OpenAPI、generated contract、实施/验收计划和验收记录。
- 未到阶段、缺授权或缺安全样本的验收命令必须 `exit 2`；它不是 PASS。
- ADR 存放在 `docs/adr/YYYY-MM-DD-title.md`；关键架构、范围或安全边界变化必须先形成 ADR。
