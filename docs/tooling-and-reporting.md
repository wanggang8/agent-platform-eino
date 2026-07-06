# Tooling 与报告规范

本文固定前端、契约、视觉和 smoke 报告的工具链规则，避免不同机器产生不可复现的验收结果。

## Node 与包管理

- Node 版本由 `web/eino-workbench/package.json` 的 `engines.node` 固定。
- 默认使用 npm 和 `package-lock.json`。
- CI 和本地都使用 `npm ci` 进行可复现安装。
- 根 `package.json` 只提供 wrapper scripts，不复制前端包内部逻辑。

当前 Phase 1 文档契约已提供根级 wrapper：

```bash
npm run eino-workbench:schema-test
npm run eino-workbench:contract-generate
npm run eino-workbench:contract-test
npm run eino-workbench:openapi-lint
npm run eino-workbench:typecheck
npm run eino-workbench:test
npm run eino-workbench:send-smoke
```

这些命令使用 `scripts/eino_workbench_schema_validate.mjs`、Ajv 2020 和 `ajv-formats` 校验 schema、fixture manifest、OpenAPI `$ref`、敏感标记和开放对象；`contract-generate` 由 JSON Schema 生成前端可导入的契约索引。

`schema-test` 还会校验 `docs/fixtures/fobrain/tool-matrix-24.json` 精确覆盖 24 个 Fobrain 只读工具，不允许重复或遗漏。

以下命令已预留，但在对应实现阶段前必须失败退出，不得作为通过信号：

```bash
npm run eino-workbench:stream-test
npm run eino-workbench:browser-test
npm run eino-workbench:visual-test
```

## Playwright

`web/eino-workbench/playwright.config.ts` 必须定义：

- `webServer` 自动启动前端或整合服务。
- `baseURL` 从 `EINO_WORKBENCH_BASE_URL` 读取，默认 `http://127.0.0.1:8081`。
- desktop viewport：1440x900。
- mobile viewport：390x844。
- reporter 输出到 `test-results/eino-workbench-playwright-report/`。
- screenshot/video/trace 失败时保留。
- visual snapshot threshold 和 mask 规则。

唯一截图更新命令：

```bash
npm run eino-workbench:visual-test -- --update-snapshots
```

## 视觉基线

基线来源：

- 当前产品参考截图；旧项目最终可用路径见 `legacy-acceptance-evidence.md`。
- 固定 fixture 数据。
- `docs/visual-acceptance-matrix.md` 中的 block selector 和状态清单。
- desktop/mobile 固定 viewport。
- 动态字段使用 mask：时间、run id、随机 id、模型耗时、token 计数。

旧项目截图只用于建立新项目视觉目标和覆盖范围。新项目验收必须重新输出 Playwright baseline、visual report、真实服务截图和当次验收记录。

基线更新必须在验收记录中写明：

- 更新原因。
- diff 路径。
- 影响的状态。
- 是否改变产品信息层级。

## 报告路径

| 报告 | 路径 |
| --- | --- |
| 真实模型 provider 报告 | `test-results/eino-workbench-real-model-provider-report.json` |
| 真实模型工具报告 | `test-results/eino-workbench-real-model-tool-report.json` |
| live read 报告 | `test-results/eino-workbench-live-read-report.json` |
| live write 报告 | `test-results/eino-workbench-live-write-report.json` |
| skipped 报告 | `test-results/eino-workbench-skip-report.json` |
| telemetry summary 报告 | `test-results/eino-workbench-telemetry-usage-summary-report.json` |
| 视觉报告 | `test-results/eino-workbench-visual-report/` |

telemetry summary 报告由独立内部脚本生成，不连接 Workbench、Action API、Product Facts 或真实 provider：

```bash
go run ./scripts/telemetry_summary_report --output test-results/eino-workbench-telemetry-usage-summary-report.json
node scripts/eino_workbench_report_validate.mjs --schema docs/schemas/telemetry_usage_summary_report.v1.schema.json --report test-results/eino-workbench-telemetry-usage-summary-report.json
```

## skip report

skip report 必须机器可读：

```json
{
  "schema_version": "eino.skip_report.v1",
  "command": "bash scripts/eino_workbench_server_smoke.sh --scenario fobrain-live-read",
  "missing_env": ["FOBRAIN_USER_API_TOKEN"],
  "credential_scope": "fobrain-live-read",
  "reason": "live credential unavailable",
  "rerun_condition": "provide Fobrain live token",
  "blocks_claims": ["fobrain-capability-comparable", "rewrite-complete"],
  "expires_at": "2026-07-07T00:00:00Z"
}
```

任何 `blocks_claims` 非空的 skip 都不得被写成通过。
