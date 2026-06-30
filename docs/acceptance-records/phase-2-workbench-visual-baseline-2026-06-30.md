# Phase 2 Workbench Visual Baseline 验收记录

阶段：Phase 2.2 / Phase 2.4
日期：2026-06-30
方案：Eino-first Workbench fixture-only 前端视觉 baseline
执行人：Codex

## 环境

- Go：go1.23.12 darwin/arm64
- Node：v25.8.1
- npm：11.11.0
- OS：macOS 26.5.1
- Eino 版本：见 `go.mod`
- 模型 provider：未接入，本阶段只验证 fixture-only Workbench
- Fobrain 环境：未接入，本阶段只使用旧验收截图作为视觉目标参考

## 命令

```bash
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
npm run eino-workbench:typecheck
npm run eino-workbench:test
npm run eino-workbench:build
npm run eino-workbench:browser-test -- --grep @shell
npm run eino-workbench:visual-test -- --update-snapshots
npm run eino-workbench:visual-test
```

## 报告

- schema/contract：通过，29 个 schema 和 fixture manifest 校验通过，contract index current。
- Playwright：通过，`browser-test -- --grep @shell` 覆盖 desktop/mobile shell smoke，报告本地输出到 `test-results/eino-workbench-playwright-report/results.json`。
- visual：通过，`visual-test` 覆盖 18 个 desktop/mobile visual tests，报告本地输出到 `test-results/eino-workbench-visual-report/results.json`。
- real model：未执行，Phase 3+ 后接入。
- live read：未执行，Phase 8 后接入。
- live write：未执行，Phase 8 后接入。
- skip：real model/live read/live write 不属于本阶段通过声明。

## 截图

- desktop baseline：`web/eino-workbench/tests/__screenshots__/desktop/{state_id}-{block_id}.png`
- mobile baseline：`web/eino-workbench/tests/__screenshots__/mobile/{state_id}-{block_id}.png`
- legacy target reference：`docs/assets/legacy/workbench/`
- diff：本次为 baseline 初始化，无旧新项目 baseline diff；后续 baseline 更新必须附 diff 和 mask 说明。

## skip

- command：无本阶段必跑命令 skip。
- missing_env：无。
- credential_scope：无。
- reason：真实模型、Fobrain live 能力未到对应 Phase。
- rerun_condition：进入 Phase 3/8 后按 `docs/08-acceptance-plan.md` 补跑。
- blocks_claims：阻断 P0/P1/P2 完成声明，阻断“重构完成”和“可替换当前产品基线”声明。
- expires_at：进入 Phase 3 前需要更新前端验收记录。

## 失败分类

- contract：无。
- safety：无可见 raw tool id、result ref、credential、resume token 或 raw provider payload 泄漏。
- visual：无，本阶段 baseline 已建立。
- model/tool selection：未覆盖。
- HITL/resume：只覆盖审批/澄清 fixture 视觉，不覆盖真实 resume。
- Fobrain live：未覆盖。

## 结论

- 通过 / 不通过 / skipped blocking：Phase 2 fixture-only Workbench 视觉 baseline 通过。
- 阻断 P0/P1/P2：不单独阻断 Phase 2；仍不能声明 P0/P1/P2 完成。
- 阻断重构完成声明：是，后续 Phase 未完成。
- 允许替换当前产品基线：否。
- report schema 校验结果：schema-test 通过。
- 不得声明的能力：真实模型、Eino stream、Action API 一致性、HITL resume、Fobrain live、P0/P1/P2 完成。
- 关联 ADR：无。
