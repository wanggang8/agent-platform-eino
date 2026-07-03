# Phase 8 Batch E Visual Replay Audit Acceptance Record

日期：2026-07-03

范围：Fobrain Batch E 四个详情与风险关联只读工具的 Workbench 视觉、fresh replay 和 audit evidence。

## 覆盖范围

本记录覆盖：

- `tool.fobrain.get_asset_detail`
- `tool.fobrain.get_vulnerability_detail`
- `tool.fobrain.business_risk_summary`
- `tool.fobrain.threat_relevance_list`
- desktop/mobile 六区域截图：`main-chat`、`fresh-main-chat`、`process`、`evidence`、`audit`、`internal-details`

不覆盖：

- Batch B-D 视觉证据补齐。
- Action API 真实服务同源 smoke。
- Fobrain 24 个只读工具全量恢复声明。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| Batch E 四工具拥有 Workbench fixture | `docs/fobrain-batch-e-interface-plan.md` | `web/eino-workbench/src/fixtures/fobrainVisualFixtures.ts` | Vitest | Pass | fixture 直接消费 `docs/fixtures/fobrain/*-result.json` 的 StructuredResult。 |
| 每个工具覆盖六个视觉区域 | `docs/07-implementation-plan.md` Task 8.2 | `web/eino-workbench/tests/__screenshots__/{desktop,mobile}/fobrain*Detail-*.png` 和 `fobrain*Summary-*.png` | Playwright visual test | Pass | 4 个工具 x 6 区域 x desktop/mobile = 48 张 baseline。 |
| replay 视图与主聊天同源 | `docs/08-acceptance-plan.md` | `fresh-main-chat` screenshots | Playwright reload/fresh render | Pass | fresh view 重新从同一 WorkbenchView fixture 投影，不读取前端临时状态。 |
| audit 视图只展示安全摘要 | `docs/06-security-and-projection.md` | `*-audit.png` screenshots | Playwright forbidden text assertion | Pass | audit tab 展示 selected tool、safe args summary、evidence refs 和 run 状态。 |
| 不泄漏 token、auth、raw provider payload 或 credential ref | `docs/fobrain-live-read-batch-plan.md` | forbidden text assertions | Vitest + Playwright | Pass | 覆盖 `authorization`、`bearer`、`api_token`、`credential_ref`、`raw body` 等关键词。 |
| 不声明 24/24 或 Batch B-D 完成 | `docs/08-acceptance-plan.md` | 本记录和计划文档 | 文档检查 | Pass | Batch E 视觉切片完成不等于 Fobrain 全量恢复完成。 |

## 执行命令

```bash
npm --workspace @agent-platform-eino/eino-workbench run test -- src/fixtures/fobrainVisualFixtures.test.ts --run
npm --workspace @agent-platform-eino/eino-workbench run visual-test -- --grep fobrain --update-snapshots
npm --workspace @agent-platform-eino/eino-workbench run visual-test -- --grep fobrain
find web/eino-workbench/tests/__screenshots__ -type f | rg 'fobrain(AssetDetail|VulnerabilityDetail|BusinessRiskSummary|ThreatRelevanceList)' | sort | wc -l
```

## 结论

Batch E 四个只读工具的 Workbench 视觉、fresh replay 和 audit evidence 已通过当前 fixture-based 门禁。该结论只允许声明 Batch E 视觉/replay/audit 切片完成；Fobrain 24 个只读工具最终恢复仍需 Batch B-D 视觉、Action API 同源 smoke、实体消歧和写域审批等后续门禁。
