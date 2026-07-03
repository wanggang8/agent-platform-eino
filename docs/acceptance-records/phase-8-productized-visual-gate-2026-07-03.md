# Phase 8 Productized Visual Gate Acceptance Record

日期：2026-07-03

范围：Workbench / Fobrain 视觉验收口径从技术调试截图收紧为产品验收截图。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| 产品截图不展示英文内部调试词 | 用户最新要求、`docs/visual-acceptance-matrix.md` | `web/eino-workbench/tests/fobrain-tool-visual.spec.ts` | Playwright visible text assertion | Pass | 断言 ASCII 英文调试词、tool id、schema、run id 不进入可见文本。 |
| 不把嵌套对象显示为 JSON | `docs/02-ux-visual-requirements.md` | `StructuredResultView` 单元测试 | Vitest | Pass | 对象值显示为“已整理”，不输出 `{}` 或字段名。 |
| Inspector 展示中文产品化运行/审计摘要 | `docs/visual-acceptance-matrix.md` | `Inspector.tsx`、视觉 baseline | Vitest + Playwright | Pass | `Product Facts`、`safe projection only`、`tool` 等内部词不再可见。 |
| Fobrain visual baseline 使用产品化中文展示 | Phase 8.2/8.3 | `web/eino-workbench/tests/__screenshots__/` | Playwright visual test | Pass | 重新生成受影响 desktop/mobile baseline。 |
| 产品化视觉门禁不替代同源事实 smoke | `docs/08-acceptance-plan.md` | 本记录 | 文档检查 | Pass | 本记录只证明可见产品展示；Workbench/Action API 同源 Product Facts 仍按后续真实服务 smoke 验收。 |

## 结论

后续视觉验收按产品界面验收：截图可见区域必须展示处理后的中文产品内容，不得展示工具原始 JSON、英文内部状态、tool id 或 schema 细节。

本记录不声明 Action API 与 Workbench 的真实服务同源事实门禁通过，也不替代 replay/audit 的后端 Product Facts smoke。
