# Phase 8.3 Fobrain Connector 视觉验收记录

阶段：Phase 8.3 Connector 状态与凭据绑定  
日期：2026-07-02  
范围：`connector.fobrain.security` Workbench 视觉证据

## 覆盖范围

本记录覆盖：

- `connector.fobrain.security`
- connector status 安全展示
- workspace credential binding 安全摘要展示

不覆盖：

- Batch B-E 和 24 个只读工具全量恢复。
- 实体消歧、写域审批、live write。
- 真实 token 或 raw connector 配置展示。

## 验收矩阵

| Criterion | Source | Evidence | Verification method | Status | Notes/Risks |
| --- | --- | --- | --- | --- | --- |
| connector fixture 使用 `fobrain.tool_result.v2` | `docs/schemas/fobrain/tool_result.v2.schema.json` | `docs/fixtures/fobrain/connector-security-result.json` | schema-test + Vitest | Pass | display_type 为 `connector_status` |
| Workbench 展示 connector 状态和凭据绑定安全摘要 | `docs/provider-policy-and-credentials.md` | `web/eino-workbench/src/fixtures/fobrainVisualFixtures.ts` | Vitest + visual-test | Pass | 只展示 status/workspace/owner scope |
| 六区域视觉证据由新项目生成 | `docs/07-implementation-plan.md` Task 8.3 | `web/eino-workbench/tests/__screenshots__/{desktop,mobile}/fobrainConnectorSecurity-*.png` | Playwright visual test | Pass | 12 张 baseline |
| 不泄漏凭据或 raw provider 材料 | `docs/06-security-and-projection.md` | forbidden text 断言 | Vitest + Playwright | Pass | 覆盖 authorization、bearer、api_token、credential_ref、raw body 等关键词 |
| 后端 connector/credential 安全摘要门禁通过 | `docs/07-implementation-plan.md` Task 8.3 | `internal/einoapp/providers/fobrain/batch_a_test.go` | Go test | Pass | connector 不调用业务读取 client |
| 不声明 24/24 或后续能力完成 | `docs/08-acceptance-plan.md` | 本记录和计划文档 | 文档检查 | Pass | 后续仍需 Batch B-E |

## 执行命令

```bash
npm --workspace @agent-platform-eino/eino-workbench run test -- src/fixtures/fobrainVisualFixtures.test.ts --run
npm --workspace @agent-platform-eino/eino-workbench run visual-test -- --grep fobrain --update-snapshots
npm --workspace @agent-platform-eino/eino-workbench run visual-test -- --grep fobrain
npm --workspace @agent-platform-eino/eino-workbench run typecheck
npm --workspace @agent-platform-eino/eino-workbench run test -- --run
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
go test ./internal/einoapp/providers/fobrain -run 'ConnectorStatus|CredentialBinding' -count=1
go test ./internal/einoapp/... -run CredentialLeak -count=1
go test ./... -count=1
git diff --check
```

## 结论

`connector.fobrain.security` 的 Workbench 视觉证据通过。该结论只允许作为 Batch A connector 状态和凭据绑定安全展示证据，不代表 Fobrain 全量恢复完成。
