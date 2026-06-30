# 意图识别与能力选择

本文定义用户自然语言、Action API、Capability Registry 和 Eino tool call 之间的选择边界。目标是避免把业务意图硬编码成后端关键词路由，同时保证工具选择可验收、可审计、可恢复。

## 设计原则

- 不新增传统 `IntentClassifier` 作为业务主路由。
- 不按自然语言关键词硬编码 Fobrain 工具选择。
- 顶层 selection 只决定入口，不决定具体业务工具。
- 业务工具选择优先由 Eino ChatModelAgent 的 native tool call 完成。
- 后端强制执行 schema、policy、credential、approval、clarification 和 Safety Gate。
- 所有选择结果必须进入 Product Facts，供 Workbench、ActionResult、Replay、Audit 和 real model report 同源读取。

## 两层选择

### 顶层入口选择

顶层选择只处理以下来源：

| 来源 | 行为 |
| --- | --- |
| `product_action` | 明确产品动作，按 action id 进入对应 capability 或 workflow。 |
| `capability_hint` / requested capability | 只作为候选入口，必须存在于 registry 且通过 policy。 |
| 受信 `intent_hint` | 只接受模型或系统产生的结构化 hint；低置信度或未知 reason code 忽略。 |
| 普通自然语言 | 默认进入 `llm.chat` / Eino ChatModelAgent。 |

顶层选择不得基于“资产”“漏洞”“工单”等词直接选择具体 Fobrain 工具。若用户输入不足或目标不明确，进入 clarification，而不是静默猜测。

### 业务工具选择

进入 ChatModelAgent 后，Capability Registry 把可用能力转换为 Eino tools：

```text
Capability metadata
  -> Eino tool name / description / input schema
  -> ChatModelAgent native tool call
  -> arguments schema validation
  -> provider policy / credential / approval
  -> provider execution
  -> StructuredResult candidate
  -> Safety Gate
  -> Product Facts
```

模型选择工具的材料只能来自：

- capability id 和稳定 tool name。
- `display_name_zh`。
- `description`。
- input schema。
- risk / side effect / approval metadata。
- safe prior facts 和 safe StructuredResult。

禁止把旧 runtime step id、旧 tool pair、旧 DOM 字段、raw provider payload、credential ref 或内部 reason code 作为模型选择材料。

## Capability 描述规则

Capability description 是模型工具选择的主要产品材料，必须写清：

- 适用自然语言问题。
- 不适用场景。
- 必填参数和可选参数。
- 参数不足时是否先读、澄清或拒绝。
- 与相邻工具的区别。
- 写域是否必须审批。
- 结果为空、失败或多候选时的安全行为。

示例：

```text
Use for owner asset lookup. Do not use for owner vulnerability lookup;
use tool.fobrain.list_vulnerabilities_by_owner when the user asks for vulnerabilities.
required parameter: person_name, unless person_staff_id is explicitly provided.
```

## 澄清与消歧

以下情况必须进入 `PendingInteraction(kind=clarification)`：

- 多个实体候选都可信，无法自动选择。
- 必填参数缺失，且不能从 safe prior facts 推导。
- 用户表达互相矛盾。
- capability_hint 与用户文本明显不一致。
- 工具结果返回 `entity_resolution` 且需要用户选择。

`tool.fobrain.query_normalize_and_resolve` 只能生成解析事实或候选，不得替代业务读工具，也不得输出下一步工具推荐。

## 写域和高风险能力

写域工具选择不是执行许可。即使模型选择了写域工具，也必须先经过：

```text
provider policy
  -> credential binding
  -> approval_required
  -> PendingInteraction(kind=approval)
  -> approve resume
  -> mutation once
```

未审批前不得执行真实 mutation。reject、cancel、expired 后不得继续执行。

## Product Facts

每次选择必须记录安全事实：

| 事实 | 要求 |
| --- | --- |
| top-level selection | selected capability、source、reason code、confidence、candidate count。 |
| tool selection | selected tool、tool call id、安全参数摘要、schema validation status。 |
| policy decision | allowed/denied/approval_required、policy ref、safe reason。 |
| clarification | question、candidate refs、resume ref。 |
| final result | StructuredResult ref、display type、status。 |

不得记录 raw prompt、raw provider body、credential ref、reusable resume token 或完整 secret。

## 验收要求

真实模型和 mock 验收必须覆盖：

- 普通自然语言默认进入 ChatModelAgent。
- Action API 显式 action 不被自然语言 selector 覆盖。
- capability_hint 不存在或无权限时安全拒绝或降级。
- 资产 vs 漏洞不误选。
- 单 IP 查询 vs IP 统计不误选。
- connector status 不替代业务读取。
- 写域工具只进入 approval waiting，不直接 mutation。
- 同名人员、同名部门或实体多候选进入 clarification。
- selected tool、safe args summary、policy decision 和 StructuredResult ref 写入 Product Facts，并出现在 audit/replay/real model report 中。

失败即阻断对应 Phase：P1 阻断主链可用声明，P2 阻断 Fobrain 能力可比声明。
