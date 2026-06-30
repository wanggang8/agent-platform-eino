# Conversation Context

本文定义模型输入上下文如何从 Product Facts 安全投影。它不是长期记忆系统，也不允许把旧 runtime 材料或 provider raw payload 回流到模型。

## 目标

开放式对话由 Eino ChatModelAgent 执行。项目只负责构造安全上下文：

```text
Product Facts
  -> safe context selector
  -> context snapshot
  -> model messages / tool observations
  -> ChatModelAgent
```

P0/P1 只做同一 session 内的 safe prior facts。跨 session 长期记忆、向量检索和用户画像不进入当前重构目标。

## 允许来源

| 来源 | 可进入模型上下文 | 规则 |
| --- | --- | --- |
| User Turn | 是 | 只使用用户可见文本；过 Safety Gate。 |
| Assistant Turn | 是 | 只使用已通过 Safety Gate 的最终回答或必要摘要。 |
| ToolResult | 是 | 只使用 StructuredResult 中允许字段、safe summary、result_ref。 |
| PendingInteraction | 是 | 只使用安全问题、候选 safe refs、用户选择。 |
| AuditEvent | 限制 | 只可用于安全摘要和故障恢复，不作为业务事实来源。 |
| Provider raw payload | 否 | 永不进入上下文。 |
| Checkpoint / interrupt id | 否 | 只用于内部恢复。 |
| Credential / token / raw prompt | 否 | 永不进入上下文。 |

## Context Snapshot

每次调用模型前必须生成 `context_snapshot`，作为内部事实或 audit 安全摘要记录：

| 字段 | 要求 |
| --- | --- |
| `snapshot_id` | 稳定 id。 |
| `run_id` / `turn_id` | 关联当前运行和轮次。 |
| `included_fact_refs` | 使用的 safe fact refs。 |
| `excluded_reasons` | 被排除事实的安全原因。 |
| `token_budget` | 本次上下文预算。 |
| `truncation` | 是否截断、截断策略。 |
| `safe_summary` | 可审计摘要，不含 raw prompt/provider body。 |

`context_snapshot` 不直接展示给用户，但可在 Inspector audit/runtime 中显示安全摘要。

## 多轮引用规则

多轮追问只能引用 safe prior facts：

- “这个 IP” 可解析为上一轮 StructuredResult 中唯一 IP。
- “刚才那个人” 可解析为上一轮 clarification 用户选择或唯一 person candidate。
- “上一个漏洞” 可解析为上一轮唯一 vulnerability safe ref 或真实业务 ID。
- 若存在多个候选，必须进入 clarification。
- 若上轮只有摘要没有结构化实体，不得猜测参数。

## 截断策略

上下文按以下优先级保留：

1. 当前用户消息。
2. 当前 pending/resume 数据。
3. 直接相关的 StructuredResult refs。
4. 最近 assistant final answer。
5. 安全运行摘要。

超出预算时必须截断低优先级材料，并在 `context_snapshot.truncation` 中记录。截断不得改变 Product Facts。

## 禁止事项

- 不保存或注入 raw provider payload。
- 不把 audit 当作业务事实替代 StructuredResult。
- 不从前端 view 反推模型上下文。
- 不让模型看到 checkpoint id、interrupt id、resume token。
- 不使用旧项目 memory、runtime step、tool pair 结构作为新上下文契约。

## 验收

- 多轮 IP、资产、漏洞、人员候选复用可测。
- 多候选时进入 clarification。
- context snapshot 记录 included/excluded refs 和 token budget。
- 安全扫描确认上下文不含 raw provider body、credential、checkpoint id、interrupt id。
- Workbench、ActionResult、Replay 与上下文引用同一 Product Facts。

参考资料：

- Eino ChatModelAgent：`https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/`
- Eino ChatModel：`https://www.cloudwego.io/docs/eino/core_modules/components/chat_model_guide/`
