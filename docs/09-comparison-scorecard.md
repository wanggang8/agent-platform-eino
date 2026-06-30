# 对比评分表

本文记录 A/B/C 三种方案的真实实验结果。

## 方案

- A：当前自研 Runtime 基线。
- B：当前项目渐进替换 Eino。该方案是后续对照项，入口设计需要单独补充 ADR 或对照计划。
- C：Eino-first 全栈重开发。

## 指标

| 指标 | 采集方式 |
| --- | --- |
| 后端核心 LOC | `find internal/einoapp cmd/eino-workbench -name '*.go' -print0 | xargs -0 wc -l` |
| 前端 LOC | `find web/eino-workbench/src \\( -name '*.ts' -o -name '*.tsx' -o -name '*.css' \\) -print0 | xargs -0 wc -l` |
| Contract/脚本 LOC | `find docs scripts -path '*eino*' -type f | xargs wc -l` |
| 核心概念数量 | 人工列举 |
| 普通聊天闭环文件数 | 人工列举 |
| 工具调用闭环文件数 | 人工列举 |
| HITL 闭环文件数 | 人工列举 |
| Action API 闭环文件数 | 人工列举 |
| 测试数量 | `rg "func Test|test\\("` |
| 真实模型工具成功率 | real model report |
| Fobrain 产品恢复率 | 24 只读、connector、clarification、write approval、live smoke |
| 安全回退数 | safety tests |
| Contract 漂移次数 | schema/fixture 变更记录 |
| 视觉漂移次数 | Playwright visual baseline |
| 调试入口数量 | 验收记录 |

## 阶段记录模板

```text
阶段：
日期：
方案：
命令：
结果：
后端 LOC：
前端 LOC：
Contract/脚本 LOC：
核心概念：
闭环文件数：
测试数量：
失败点：
调试路径：
安全问题：
视觉问题：
结论：
```

## Phase 0

尚未执行。

## Phase 1

尚未执行。

## Phase 2

尚未执行。

## Phase 3

尚未执行。

## Phase 4

尚未执行。

## Phase 5

尚未执行。

## Phase 6

尚未执行。

## Phase 7

尚未执行。

## Phase 8

尚未执行。

## Phase 9

尚未执行。

## 判断规则

09 只记录评分和对比结果，不判定阶段是否完成。P0/P1/P2 是否通过、C 方案是否可声明重构完成或能力可比，以 `08-acceptance-plan.md` 的验收结论为准。

P1 后可以根据评分判断 C 架构是否需要调整，但不改变 P2/Phase 8 作为完整重构必做门禁的地位；未完成 P2 不得声明重构完成或能力可比。

B 方案只作为后续对照项。若要正式实施 B，必须先补充 ADR 或对照计划，说明如何避免 Eino 状态和旧自研状态双写复杂度。
