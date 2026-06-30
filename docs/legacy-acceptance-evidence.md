# 旧项目最终验收证据索引

本文记录旧项目中可作为产品能力、视觉体验和安全边界参考的最终验收证据。旧项目路径固定为 `/Users/vick/Desktop/project/ai-agent`，只读，不得修改。

这些证据不是新项目的通过结果。新项目必须重新生成自己的 `test-results/eino-*` 报告和截图；旧证据只用于防止能力丢失、视觉跑偏和验收范围缩水。

## 权威聚合入口

最终验收以旧项目聚合报告为准：

```text
/Users/vick/Desktop/project/ai-agent/test-results/workbench-acceptance-audit/workbench-acceptance-audit.json
```

该报告生成于 `2026-06-28T09:59:38.602Z`，结论为 `6/6 passed`、`failed=0`。只有该报告引用的 artifact 才视为最终验收链路的一部分。

## 最终通过证据

| 证据 | 旧项目路径 | 最低保留结论 |
| --- | --- | --- |
| Workbench 全场景回放 | `test-results/workbench-ui-rewrite/17-full-scenario-acceptance/full-scenario-acceptance.json` | `25` 场景，`25 passed`，`0 failed`，`275` 截图，静态边界通过 |
| Workbench 图像审计 | `test-results/workbench-ui-rewrite/17-full-scenario-acceptance/image-audit.json` | `275` 截图，`failedCount=0`，`11` contact sheets |
| Task15 full shell DOM | `test-results/workbench-ui-rewrite/15-visual-acceptance/full-shell-dom-metrics.json` | `25` 场景，`25` full-shell 截图，无 failed records |
| Task15 视觉指标 | `test-results/workbench-ui-rewrite/15-visual-acceptance/visual-metrics.json` | `25 passed`，`0 failed`，`150` 报告截图，产品和 DOM 门禁全通过 |
| Fobrain 工具矩阵 | `test-results/fobrain-tool-acceptance/fobrain-tool-acceptance-matrix.json` | `24` 只读工具 + `1` connector，`25 passed`，`0 failed` |
| Fobrain live smoke | `test-results/fobrain-live-smoke/fobrain-live-smoke.json` | `7 passed`，`0 failed`，读域和写域审批 gate 通过 |

Fobrain 工具矩阵的源报告为：

```text
/Users/vick/Desktop/project/ai-agent/test-results/workbench-real-model-tool-report.json
```

每个工具场景必须至少保留旧验收中的六类截图区域要求：`main-chat`、`fresh-main-chat`、`process`、`evidence`、`audit`、`internal-details`。

## 视觉参考材料

当前产品视觉基线从旧项目读取：

```text
/Users/vick/Desktop/project/ai-agent/docs/design/assets/workbench-target-ui-v1-2026-06-26.png
/Users/vick/Desktop/project/ai-agent/docs/design/workbench-ui-acceptance-matrix.md
/Users/vick/Desktop/project/ai-agent/test-results/workbench-ui-rewrite/15-visual-acceptance/
/Users/vick/Desktop/project/ai-agent/test-results/workbench-ui-rewrite/17-full-scenario-acceptance/screenshots/
/Users/vick/Desktop/project/ai-agent/test-results/workbench-ui-rewrite/17-full-scenario-acceptance/contact-sheets/
```

新项目可以参考这些截图的信息层级、布局密度、状态覆盖和验收粒度，但不得复制旧 DOM、旧 CSS 类名、旧前端 presenter 分支或旧接口结构。

## 不作为最终通过证据

以下材料只能作为历史参考，不能单独证明最终验收通过：

- `test-results/ui-redesign-audit/`
- `test-results/workbench-design-audit-2026-06-24/`
- 零散的 `workbench-complex-context-*` 截图。
- 未被 `workbench-acceptance-audit.json` 引用的历史报告或截图。

## 新项目使用规则

- Phase 1 只引用本文作为范围基线，不声明能力已恢复。
- Phase 2 必须用新项目 Playwright 和视觉报告重建截图证据。
- Phase 8 必须用新项目 Fobrain mock/live 报告重新证明 `24 + 1` 覆盖。
- 无 live 凭据时必须生成 skip report，且不得声明产品能力可比。
- 若旧证据与新文档冲突，以新项目 `docs/` 中的 Product Facts、StructuredResult、Safety Gate 和 Eino-first 架构为准。
