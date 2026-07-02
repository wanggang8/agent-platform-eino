# Phase 8.2 Fobrain Batch A 视觉验收记录

阶段：Phase 8.2 Fobrain presentation evidence matrix  
日期：2026-07-02  
范围：Batch A 业务只读工具视觉证据切片

## 覆盖范围

本记录覆盖：

- `tool.fobrain.current_user_context`
- `tool.fobrain.my_permissions`

不覆盖：

- `connector.fobrain.security`，按 Task 8.3 单独验收。
- Batch B-E、24 个只读工具全量恢复、live write、写域审批。

## 证据要求

每个工具必须由新项目生成以下六个视觉区域，不复用旧项目截图：

- `main-chat`
- `fresh-main-chat`
- `process`
- `evidence`
- `audit`
- `internal-details`

本次新增 Playwright baseline 共 24 张：2 个工具 × 6 个区域 × desktop/mobile。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| 两个 Batch A 业务只读工具均有 Workbench fixture | `docs/fobrain-live-read-batch-plan.md` | `web/eino-workbench/src/fixtures/fobrainVisualFixtures.ts` | Vitest fixture test | Pass | connector 不在本切片范围 |
| 每个工具覆盖六个视觉区域 | `docs/07-implementation-plan.md` Task 8.2 | `web/eino-workbench/tests/__screenshots__/{desktop,mobile}/fobrain*.png` | Playwright visual test | Pass | 24 张 baseline 已生成 |
| 视觉内容来自 StructuredResult 安全投影 | `docs/facts-contract.md`、`docs/fobrain-tool-matrix.md` | Fobrain fixture 引用 `docs/fixtures/fobrain/*-result.json` | Vitest 安全断言 + Playwright 可见文本断言 | Pass | raw provider payload 不进入前端 fixture |
| 不泄漏凭据或 raw provider 材料 | `docs/06-security-and-projection.md` | forbidden visible text 断言 | Vitest + Playwright | Pass | 覆盖 authorization、bearer、api_token、credential_ref、raw body 等关键词 |
| 不声明 24/24 或 connector 已完成 | `docs/08-acceptance-plan.md` | 本记录和计划文档限制说明 | 文档检查 | Pass | 后续仍需 Task 8.3 和 Batch B-E |
| 契约生成清单保持同步 | `docs/05-contract-design.md` | `web/eino-workbench/src/contracts/generated.ts` | contract-test | Pass | 补齐既有 Fobrain report schema 清单 |

## 执行命令

```bash
npm --workspace @agent-platform-eino/eino-workbench run test -- src/fixtures/fobrainVisualFixtures.test.ts --run
npm --workspace @agent-platform-eino/eino-workbench run visual-test -- --grep fobrain --update-snapshots
npm run eino-workbench:contract-generate
npm run eino-workbench:contract-test
npm --workspace @agent-platform-eino/eino-workbench run typecheck
npm --workspace @agent-platform-eino/eino-workbench run test -- --run
npm run eino-workbench:schema-test
npm --workspace @agent-platform-eino/eino-workbench run visual-test -- --grep fobrain
go test ./... -count=1
git diff --check
```

## 结论

Phase 8.2 Batch A 业务只读视觉证据切片通过。该结论只允许作为继续推进 `connector.fobrain.security` 视觉证据和 24 个只读工具矩阵的阶段性证据，不代表 Fobrain 全量恢复完成。
