# Tooling 与报告规范

本文固定前端、契约、视觉和 smoke 报告的工具链规则，避免不同机器产生不可复现的验收结果。

## Canonical 工具链环境

Story 1.1 起，发布构建和 desktop visual 的唯一目标环境由
`build/toolchain/toolchain.lock` 与 `build/toolchain/Dockerfile` 固定。构建并取得本地 image ID：

```bash
image_id=$(bash scripts/build_toolchain_image.sh --load | awk -F= '/^TOOLCHAIN_IMAGE_ID=/{print $2}')
test -n "$image_id"
```

canonical identity 不只包含 runtime 与 browser：Node 使用官方 linux-x64 `tar.gz` 及其 SHA-256；
Noble `/etc/apt/sources.list.d/ubuntu.sources` 的每个 `Signed-By` stanza 都注入 lock 中的
`Snapshot: 20260708T000000Z`，再安装 lock 固定的 `build-essential`。Dockerfile 必须确认 snapshot
插入数量与 stanza 数量一致、清理 apt lists，并以 `command -v cc` 和临时 module 中的
`CGO_ENABLED=1 go test -race ./...` 固定 race detector 能力。以上值只能来自
`build/toolchain/toolchain.lock`，不得由环境变量或浮动 apt 状态覆盖。

进入同一环境进行诊断，或执行唯一的顺序 baseline：

```bash
docker run --rm -it --platform linux/amd64 -v "$PWD:/workspace" -w /workspace \
  "$image_id" bash
docker run --rm --platform linux/amd64 -v "$PWD:/workspace" -w /workspace \
  "$image_id" bash scripts/run_toolchain_baseline.sh
```

`run_toolchain_baseline.sh` 先验证工具链与 CI 配置，再执行 `npm ci`；随后按 contract、Go、前端、
desktop browser、service smoke 和 clean gate 的固定顺序运行。Go 命令必须显式使用
`GOTOOLCHAIN=local`，Playwright 门禁只运行 `--project=desktop`，mobile 不属于首版
`G-TOOLCHAIN` 通过条件。日志写入 `test-results/toolchain-baseline.log`。

本地 visual migration 与 Story/G-TOOLCHAIN 使用两条独立证据链：

- Visual migration chain：使用同一组 pinned external inputs 本地构建 canonical image，取得本地
  image ID 后运行唯一 `bash scripts/run_toolchain_baseline.sh`。首次运行可能在 desktop screenshot
  diff 处非零停止；保留 baseline log 与 Playwright actual/diff，并确认此前全部非视觉检查通过后，
  才能进入人工分类和 UX review。批准后只用 desktop-only 命令更新 snapshot，再在同一 image 完整
  重跑唯一 baseline，要求所有适用检查通过。此链不要求 GitLab remote、registry digest 或
  pipeline；结果始终只是 preflight/visual migration evidence，不能替代 clean GitLab pipeline，
  也不能单独形成 `G-TOOLCHAIN PASS`。
- Story/G-TOOLCHAIN chain：clean GitLab pipeline 必须 push canonical image，产出绑定 commit 的
  registry `tag@sha256` digest，并在该 image 执行完整 baseline 与 artifacts。只有该证据与已批准的
  visual evidence 同时存在，才可裁决 `G-TOOLCHAIN PASS`。

当前 pinned MCR base layer 传输超过 15 分钟无 layer 字节进展，人工终止后 exit 130；因此本地
image ID 和 pipeline image 都未产生，两条链均未前进。仓库没有 remote/pipeline 是 Story 门禁的
额外缺口，但不阻止未来在真实本地 canonical image 上生成和记录 visual evidence。任何阻断记录都
必须写明命令、退出码、最后一个可验证步骤、未产生的证据和解除条件；不得切换未批准镜像／宿主
环境，也不得用静态 CI PASS 代替真实 pipeline。当前唯一裁决记录见
`acceptance-records/story-1-1-g-toolchain-2026-07-15.md`。

## Node 与包管理

- Node 精确版本以根目录 `.node-version` 为权威源；npm 精确版本以根 `package.json` 的
  `packageManager` 为权威源。根与 workspace 的 `engines` 是由 verifier 核对的声明镜像。
- 使用 npm 和提交的 `package-lock.json`。
- CI 和本地都使用 `npm ci` 进行可复现安装。
- 根 `package.json` 只提供 wrapper scripts，不复制前端包内部逻辑。
- Chromium 与 OS/font layer 只随 pinned canonical image 构建；测试阶段不得运行浮动的
  `playwright install`、apt 安装或其他浏览器替换命令。
- Node archive 固定为官方 linux-x64 `tar.gz`；构建期 apt 只允许使用 lock 的 Ubuntu snapshot
  安装 `build-essential`，用于 image 内的 C compiler guard 与 Go race smoke，不得临时改用宿主编译器。

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

常规视觉基线的截图更新命令：

```bash
npm run eino-workbench:visual-test -- --update-snapshots
```

跨环境迁移是例外门禁。首次必须在本地 canonical image 运行唯一 baseline：

```bash
docker run --rm --platform linux/amd64 -v "$PWD:/workspace" -w /workspace \
  "$image_id" bash scripts/run_toolchain_baseline.sh
```

该次运行预期可能在 desktop screenshot diff 处非零停止。必须保存
`test-results/toolchain-baseline.log` 与 Playwright actual/diff，并确认 visual 之前的全部非视觉检查
通过；若在此前失败，不得进入迁移审批。人工完成差异分类且 UX 明确批准后，才可在同一 image 内
执行 desktop-only snapshot update：

```bash
docker run --rm --platform linux/amd64 -v "$PWD:/workspace" -w /workspace \
  "$image_id" npm run eino-workbench:browser-test -- --project=desktop --update-snapshots
```

更新后必须在同一 image 再完整执行 `bash scripts/run_toolchain_baseline.sh`，要求所有适用检查通过。
该 visual review 不依赖 GitLab remote 或 registry digest；本地结果只属于 preflight/visual migration
evidence，不替代 clean GitLab pipeline，也不能单独形成 `G-TOOLCHAIN PASS`。未产生 canonical
image 或仍为 `PENDING_UX_APPROVAL` 时，禁止运行任何 snapshot update。

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
