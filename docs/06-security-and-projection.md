# 安全与投影

本文描述哪些数据可以展示，哪些必须留在内部，以及产品输出如何从统一事实生成。

## 核心规则

- raw provider payload 不进入产品事实。
- secret、Authorization、API key、raw credential_ref 不进入 Workbench 或 ActionResult。
- checkpoint_id、interrupt_id 不作为产品恢复凭证。
- tool arguments 只保存 hash 和安全摘要。
- 工具结果必须先转成 StructuredResult 候选，再过 Safety Gate。
- assistant 文本也必须过 Safety Gate。

## 允许展示

- 用户输入的业务值。
- IP、CVE、资产 ID、工单 ID 等业务对象。
- `credential_binding`。
- `approval_refs`。
- `resume_refs`。
- 安全摘要。
- StructuredResult 中允许的字段。

## 禁止展示

- secret。
- Authorization。
- API key。
- raw credential_ref。
- raw prompt。
- raw provider body。
- checkpoint_id。
- interrupt_id。
- internal reason_code 作为主聊天正文。
- registered tool id 作为主聊天正文。
- model_instruction。
- recovery_policy。
- fact/ref 技术串。

## StructuredResult

工具、MCP、provider 返回的数据只是 StructuredResult 候选，不等于可展示数据。

进入产品输出前必须检查：

- schema_version 是否登记。
- 字段是否符合 schema。
- payload 大小。
- 文本长度。
- 数组长度。
- 对象深度。
- 敏感字段。
- raw error。

## Assistant Safety Gate

assistant delta 和最终回答必须检查：

- raw prompt。
- raw tool args。
- provider body。
- credential。
- checkpoint/interrupt。
- internal tool id。
- reason code。
- model instruction。
- recovery policy。

不合格内容转成安全 run notice 或安全错误摘要。

## Projection

Product Facts 是唯一事实来源：

```text
Run / Turn / ToolCall / ToolResult / PendingInteraction / AuditEvent
  -> Workbench View
  -> Workbench Stream Event
  -> ActionResult
  -> Replay View
  -> Inspector evidence
  -> Inspector structured
  -> Audit list
```

禁止 Workbench、Action API、audit、replay 各自复制工具结果事实。

## 审计和回放

审计和回放展示的是安全事实，不展示 raw provider payload。

必须覆盖：

- 用户消息。
- assistant 文本。
- 工具调用。
- 工具结果。
- 审批请求。
- 审批结果。
- 澄清请求。
- 澄清结果。
- 失败和取消。

## 泄漏回归测试

必须测试：

- assistant 泄漏 raw prompt。
- assistant 泄漏 reason_code。
- tool result 泄漏 Authorization。
- tool result 泄漏 raw credential_ref。
- ActionResult 泄漏 resume token。
- replay/audit 泄漏 raw provider body。
