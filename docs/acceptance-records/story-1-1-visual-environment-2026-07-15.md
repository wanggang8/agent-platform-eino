# Story 1.1 desktop visual 环境迁移记录

日期：2026-07-15
阶段／门禁：M-0 / G-TOOLCHAIN
Verification status：`NEEDS_MORE_EVIDENCE`
迁移状态：`PENDING_CANONICAL_IMAGE_BUILD`
UX 状态：`PENDING_UX_APPROVAL`

## 裁决摘要

本记录只证明旧 desktop baseline 已被识别、目标环境已被固定，以及首次 canonical image 构建的
外部阻断已被如实保留。canonical image 没有构建完成，因而没有 image ID／最终 image digest，
没有在目标环境运行 baseline，也没有产生任何 actual、diff 或像素差。当前不能区分字体／栅格
环境差异与产品回归，不能更新 screenshot，也不能将本记录计为 `G-TOOLCHAIN PASS`。

## 旧环境与目标环境

| 项目 | 旧历史 baseline | 新 canonical 目标 |
| --- | --- | --- |
| 操作系统 | macOS 26.5.1 | Ubuntu 24.04 Noble |
| 架构 | arm64 | linux/amd64 |
| Node.js | Node 25 历史环境 | 24.18.0 |
| Playwright | 仓库 lockfile 1.61.1 | 1.61.1 |
| Chromium | 未作为旧 baseline 的可复现发布环境固定 | revision 1228 / 149.0.7827.55 |
| 环境身份 | 历史截图，只作迁移输入 | canonical image 最终 digest：未产生 |

`build/toolchain/toolchain.lock` 已固定 MCR Playwright Noble base digest，但该 base digest 不是本项目
最终 canonical image digest，不能冒充构建产物证据。

## Canonical image 构建阻断证据

首次执行：

```text
bash scripts/build_toolchain_image.sh --load
```

构建成功解析 Dockerfile frontend，并按已固定 digest 解析
`mcr.microsoft.com/playwright:v1.61.1-noble`。随后停在 base layer 的 `FROM` 步骤超过 15 分钟，
没有任何 layer 字节进展；观察期间 Build Cache 保持 3.126GB 不变。为避免无边界等待，终止该
唯一尝试，结果为：

```text
#6 CANCELED
ERROR: failed to build: failed to solve: Canceled: context canceled
exit 130
```

随后检查本地 tag 返回 `No such image`。因此：

- `TOOLCHAIN_IMAGE_ID`：未产生。
- canonical image 最终 digest：未产生。
- `scripts/run_toolchain_baseline.sh` 的 canonical 容器运行：未执行。
- desktop actual／diff：未产生。
- 像素差与回归分类：不可裁决。

本 Task 遵守已知边界，没有重试 MCR build，没有切换镜像或宿主环境，也没有伪造 image
identity、actual、diff 或审批材料。

## Desktop snapshot 逐项迁移状态

旧路径根目录为 `web/eino-workbench/tests/__screenshots__/desktop/`。下面列出当前全部 71 个旧
baseline。对每一项，canonical actual 路径均为“未产生”，canonical diff 路径均为“未产生”，
像素差均为“不可裁决”，分类均为“无法判定”。这里不填写推测路径，避免把不存在的 artifact
写成证据。

```text
approval-waiting-approval-card.png
approval-waiting-composer.png
approval-waiting-inspector.png
approval-waiting-timeline.png
chat-composer.png
chat-inspector.png
chat-shell.png
chat-sidebar.png
chat-timeline.png
chat-tool-card.png
clarification-waiting-clarification-card.png
clarification-waiting-composer.png
clarification-waiting-inspector.png
clarification-waiting-timeline.png
empty-composer.png
empty-inspector.png
empty-shell.png
empty-sidebar.png
empty-timeline.png
failed-inspector.png
failed-timeline.png
failed-tool-card.png
fobrainAssetDetail-audit.png
fobrainAssetDetail-evidence.png
fobrainAssetDetail-fresh-main-chat.png
fobrainAssetDetail-internal-details.png
fobrainAssetDetail-main-chat.png
fobrainAssetDetail-process.png
fobrainBusinessRiskSummary-audit.png
fobrainBusinessRiskSummary-evidence.png
fobrainBusinessRiskSummary-fresh-main-chat.png
fobrainBusinessRiskSummary-internal-details.png
fobrainBusinessRiskSummary-main-chat.png
fobrainBusinessRiskSummary-process.png
fobrainConnectorSecurity-audit.png
fobrainConnectorSecurity-evidence.png
fobrainConnectorSecurity-fresh-main-chat.png
fobrainConnectorSecurity-internal-details.png
fobrainConnectorSecurity-main-chat.png
fobrainConnectorSecurity-process.png
fobrainCurrentUser-audit.png
fobrainCurrentUser-evidence.png
fobrainCurrentUser-fresh-main-chat.png
fobrainCurrentUser-internal-details.png
fobrainCurrentUser-main-chat.png
fobrainCurrentUser-process.png
fobrainMyPermissions-audit.png
fobrainMyPermissions-evidence.png
fobrainMyPermissions-fresh-main-chat.png
fobrainMyPermissions-internal-details.png
fobrainMyPermissions-main-chat.png
fobrainMyPermissions-process.png
fobrainThreatRelevanceList-audit.png
fobrainThreatRelevanceList-evidence.png
fobrainThreatRelevanceList-fresh-main-chat.png
fobrainThreatRelevanceList-internal-details.png
fobrainThreatRelevanceList-main-chat.png
fobrainThreatRelevanceList-process.png
fobrainVulnerabilityDetail-audit.png
fobrainVulnerabilityDetail-evidence.png
fobrainVulnerabilityDetail-fresh-main-chat.png
fobrainVulnerabilityDetail-internal-details.png
fobrainVulnerabilityDetail-main-chat.png
fobrainVulnerabilityDetail-process.png
inspector-tabs-inspector.png
shell-fixture.png
tool-collapsed-timeline.png
tool-collapsed-tool-card.png
tool-expanded-inspector.png
tool-expanded-timeline.png
tool-expanded-tool-card.png
```

本次没有修改 `web/eino-workbench/tests/__screenshots__/desktop/**`，也没有修改 mobile screenshot。

## UX 审批门禁

| 字段 | 当前值 |
| --- | --- |
| 审批人 | 未指定 |
| 审批时间 | 未发生 |
| 可供审查的 actual/diff | 未产生 |
| 结论 | `PENDING_UX_APPROVAL` |

没有 actual/diff 时不能进入 checkpoint/human review，也不能把差异预先归类为纯字体／栅格环境
差异。任何 DOM、文案、尺寸或交互变化都必须先按产品回归调查，不能仅凭跨 OS 迁移解释。

## 解除条件

1. MCR 传输恢复后，只重跑同一 pinned `bash scripts/build_toolchain_image.sh --load`，成功取得并验证 image ID；clean GitLab pipeline 还必须产出绑定 commit 的最终 image digest。
2. 在该 image 的 linux/amd64 环境执行 `bash scripts/run_toolchain_baseline.sh`，保留每个 desktop snapshot 的真实 actual/diff 与像素差。
3. 逐项检查功能 DOM、文案、尺寸和交互，再区分环境栅格差异、产品回归或仍无法判定。
4. 通过 checkpoint/human review 展示完整 diff，记录 UX 审批人、时间和明确结论。
5. 只有批准后才可运行 desktop-only `--update-snapshots`；更新后重跑完整 canonical baseline，并确认 mobile 与业务实现目录仍无 diff。

在上述条件全部满足前，Story 1.1 保持 in-progress，`G-TOOLCHAIN` 保持 `BLOCKED`。
