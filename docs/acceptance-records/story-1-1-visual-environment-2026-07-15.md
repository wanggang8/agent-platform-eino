# Story 1.1 desktop visual 环境迁移记录

日期：2026-07-15
阶段／门禁：M-0 / G-TOOLCHAIN
Verification status：`LOCAL_VISUAL_MIGRATION_PASS`
迁移状态：`COMPLETE`
UX 状态：`APPROVED`
当前 Story 门禁：`G-TOOLCHAIN=PASS`；M-0 已完成，下一步为 M-1

## 裁决摘要

本地 canonical image 已在相同 pinned inputs 下成功构建，snapshot apt 安装、Go/Node/npm、C
compiler、最小 race smoke 与 Chromium identity 均在真实 linux/amd64 image layer 通过。首次完整
baseline 的所有非视觉检查通过，desktop browser 因旧 macOS screenshot 与 Linux candidate 差异
停止；并发首跑为 1/17 PASS、16/17 FAIL，其中 3 项为 target crash。唯一 baseline 固定为单
worker 后重跑仍为 1/17 PASS、16/17 FAIL，但 16 项全部是旧 screenshot 差异，target crash 为零。
随后仅在临时 detached
worktree 中使用首个成功构建的同 pins image、单 worker 和 `--update-snapshots` 生成审查候选，
17/17 PASS、71/71 候选齐全；在删除浮动 apt source 后重建的最新 image 中又以零更新、单 worker
方式 17/17 PASS，
证明 71 张候选完全一致。Vick 随后批准迁移，正式 desktop snapshot 已在相同最新 image 中使用
`--project=desktop --update-snapshots --workers=1` 更新，17/17 PASS；mobile snapshot 未修改。

UX checkpoint 已完成；批准后的唯一完整 baseline 已从包含本次 snapshot 与审批记录的干净提交
在同一 image 重跑，desktop 17/17、contract smoke 与末尾 `git diff --check` 全部通过，A 链闭合。
本地 visual migration PASS 本身不能替代远程 registry digest 或 clean CI；迁移前因此不能单独形成
`G-TOOLCHAIN PASS`。GitHub/GHCR 远程链现已另行闭合，当前 Story 门禁为 PASS。

## 旧环境与目标环境

| 项目 | 旧历史 baseline | 新 canonical 目标 |
| --- | --- | --- |
| 操作系统 | macOS 26.5.1 | Ubuntu 24.04 Noble |
| 架构 | arm64 | linux/amd64 |
| Go | 未作为旧 visual baseline 身份记录 | 1.26.5 |
| Node.js | Node 25 历史环境 | 24.18.0 |
| npm | 未作为旧 visual baseline 身份记录 | 11.16.0 |
| Playwright | 仓库 lockfile 1.61.1 | 1.61.1 |
| Chromium | 未作为旧 baseline 的可复现发布环境固定 | revision 1228 / 149.0.7827.55 |
| 环境身份 | 历史截图，只作迁移输入 | local image ID：`sha256:41f8317a3cc392bca3eedcf390e70b1c6ae4be0b4548cc4d7d3c4f200a29ea93`；registry digest 未产生 |

`build/toolchain/toolchain.lock` 已固定 MCR Playwright Noble base digest，但该 base digest 不是本项目
最终 canonical image digest，不能冒充构建产物证据。

## Canonical image 与首次 baseline 证据

首次执行曾在 pinned MCR base layer 长时间无字节进展后人工终止，exit 130；该历史阻断没有被
改写。修复 Node archive 与 Linux race build prerequisites 后，在完全相同的公开 pins 下重新执行：

```text
bash scripts/build_toolchain_image.sh --load
```

进一步删除基底镜像附带的浮动 NodeSource source，并显式限制 apt 只消费注入 snapshot 的
`ubuntu.sources` 后，最新重建成功，结果为：

```text
TOOLCHAIN_IMAGE_ID=sha256:41f8317a3cc392bca3eedcf390e70b1c6ae4be0b4548cc4d7d3c4f200a29ea93
```

`docker image inspect` 确认 `os=linux arch=amd64`；镜像内直接回读为 Go `1.26.5`、Node
`24.18.0`、npm `11.16.0`、`/usr/bin/cc` 与 Chromium `149.0.7827.55`。Dockerfile build layer 中的
最小 `CGO_ENABLED=1 go test -race ./...` 真实通过。随后执行唯一完整 baseline：

```text
docker run --rm --platform linux/amd64 -v "$PWD:/workspace" -w /workspace \
  sha256:41f8317a...29ea93 bash scripts/run_toolchain_baseline.sh
```

在 browser 之前，toolchain verifier、CI validator、`npm ci`、39 项 schema/contract/OpenAPI、全部
Go tests、Linux race、checkpoint/SQLite、vet/build/import boundary、typecheck、15 项 Vitest、4 项
stream test 与 Vite build 全部通过。desktop browser 首跑产生 10 组 actual/diff 后停止，最终为
1 PASS、16 FAIL；未执行其后的 contract smoke 与最终 `git diff --check`，因此该 baseline 不是
整体 PASS。临时 detached worktree 随后在最新 image 中不更新任何文件、以单 worker 对 71 张
候选复验，17/17 PASS；并发 target crash 归类为本机 amd64 模拟资源噪声，不改变截图差异仍需
人工审批的结论。唯一 baseline 随后固定 `--workers=1` 并再次运行，browser 前检查全部通过，
desktop 1/17 PASS、16/17 screenshot diff、零 target crash；仍因正式截图尚未批准迁移而在 browser
处停止。

批准并提交正式 baseline 与审批记录后，从干净提交 `e455c3c` 使用同一 image 重跑。首次运行已通过
desktop 17/17 与 contract smoke，但因容器只挂载 worktree、无法解析指向主仓库的 `.git` 绝对指针，
末尾 clean gate exit 128，不计 PASS。补充只读挂载主仓库 Git common dir，并用最小 smoke 前后
`git status`/`git diff --check` 验证后，完整重跑 exit 0：全部前置检查、desktop 17/17、contract
smoke 与末尾 `git diff --check` 均通过。该挂载只恢复 worktree Git 元数据，不改变镜像、源码、依赖
或测试参数。

## Desktop snapshot 逐项迁移状态

路径根目录为 `web/eino-workbench/tests/__screenshots__/desktop/`。下面列出迁移前的 71 个
baseline；临时 detached worktree 的同名路径曾保存 71 张 Linux candidate。71/71 文件均发生像素
变化；46 张尺寸相同，25 张只有高度 `+1` 或 `-1` 像素，最大绝对尺寸差为 `0x1`。ImageMagick
PHASH 归一化差异均值为 `0.064117`，主要高值集中于高度很小的运行状态条和工具卡文字抗锯齿；
联系表与代表性 old/new 对照保存在忽略目录 `test-results/story-1-1-visual-*.png` 供本地 checkpoint。

逐项候选生成时，17 个 desktop tests 在单 worker 下全部通过可见性、中文产品文案和敏感文本安全
断言。审查发现布局层级、控件、状态颜色、交互区域和主要尺寸保持一致；差异主要来自跨 OS 字体／
栅格与 1 像素高度变化。`shell-fixture.png` 中资产展示值从历史 `prod-web-01` 对齐为当前安全中文
投影“生产网站一号”，与当前 fixture 和禁止 ASCII 产品文案断言一致，不属于本 Story 代码改动。
最终分类与是否批准仍由 UX reviewer 决定。

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

批准后已更新 `web/eino-workbench/tests/__screenshots__/desktop/**` 的 71 张正式 baseline；没有修改
mobile screenshot。

## UX 审批门禁

| 字段 | 当前值 |
| --- | --- |
| 审批人 | Vick |
| 审批时间 | 2026-07-15 |
| 可供审查的 actual/diff | 71 张临时 Linux candidate、10 组首跑 actual/diff、old/new 联系表与代表性对照 |
| 结论 | `APPROVED`；允许把 71 张 Linux candidate 迁移为正式 desktop baseline |

checkpoint/human review 已批准迁移。正式 snapshot update 已严格使用审查时的同一 canonical
image 和单 worker 执行；任何后续 DOM、文案、尺寸或交互变化仍必须按产品回归处理。

## 两条独立证据链

### A. Visual migration chain

1. 已使用同一组 pinned external inputs 重跑 `bash scripts/build_toolchain_image.sh --load`，成功取得并
   验证本地 image ID。
2. 在该本地 image 的 linux/amd64 环境运行唯一
   `bash scripts/run_toolchain_baseline.sh`。首次运行可能在 desktop screenshot diff 处非零停止；必须
   保存 `test-results/toolchain-baseline.log` 与 Playwright actual/diff，并从日志确认此前全部非视觉
   检查通过。若在 visual 之前失败，不得进入迁移审批。
3. 已生成 71 张临时候选并完成机器辅助的尺寸、像素与代表性视觉检查。
4. Vick 已通过 checkpoint/human review 明确批准；已在同一本地 image 运行 desktop-only
   `--update-snapshots --workers=1`，17/17 PASS。
5. 已从包含正式 snapshot 与审批记录的干净提交，在同一 image 完整重跑
   `bash scripts/run_toolchain_baseline.sh`；schema/contract/OpenAPI、Go/test/race/vet/build、
   checkpoint/SQLite/boundary、TS/Vitest/stream/build、desktop 17/17、contract smoke 与末尾 clean
   gate 全部通过。

本链只依赖真实本地 canonical image ID 与相同的 pinned inputs，不以远程平台、CI run 或 registry
digest 为前置。A 链结果始终只是 preflight/visual migration evidence，不替代 B 链 clean GitHub
Actions run，也不能单独形成 `G-TOOLCHAIN PASS`。A 链已完成。

### B. Story / G-TOOLCHAIN chain

以下原结论明确属于迁移前历史快照：当时要求 clean GitLab pipeline 产生 registry digest 与 release
artifacts；当时仓库缺少 remote/pipeline，canonical registry digest 未产生，因此 Story 1.1 保持
in-progress、`G-TOOLCHAIN` 保持 `BLOCKED`。这段历史不能作为现行平台或阻断条件。

Accepted GitHub/GHCR 迁移已用唯一公开 GitHub remote 取代上述平台。安全评审修复后的 clean GitHub Actions run
`29432627424` 对 commit `d174f79e5341dcce636d190aa0e510e4528a955a` 完成 canonical GHCR
`tag@sha256` build/push 与 digest pull/run；`toolchain-build`、`toolchain-verify`、desktop 17/17、contract
smoke 和末尾 clean gates 全部通过，并产生 30 天 `toolchain-evidence-$GITHUB_SHA` artifact。A/B 两链
均已闭合，当前唯一裁决为 `G-TOOLCHAIN=PASS`，M-0 已完成，下一步为 M-1；精确 digest、run URL、
artifact expiry 与逐项裁决见 `story-1-1-g-toolchain-2026-07-15.md`。
