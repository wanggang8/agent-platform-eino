---
status: architecture_handoff_candidate
date: '2026-07-14'
source: UX Discovery
---

# 多轮结果引用与写操作事实锁定——架构交接

## 已确认的快照分页决定

- IA-02 必须基于同一 `QueryResultSnapshot` 做快照分页，固定每页 100 条；不得在翻页时重新解释查询或混入后续新增对象。
- 后端需要提供稳定全局行序号、固定总数／总页数、稳定排序 tie-breaker，以及页码到冻结对象集合的确定性映射。
- 客户端恢复页码、全局当前行和滚动位置；分页失败只能重读同一快照页，不能隐式创建新查询结果。

## 已确认的安全展示决定

- 表格摘要与右栏完整值必须投影自同一行级 Product Fact；视觉省略不能产生第二套字符串事实。
- POC 表格单行摘要、右栏完整换行；多个 IP 在表格显示首个和剩余数量，右栏按一项一行显示全部安全 IP。
- 用户可见时间统一投影为 `YYYY-MM-DD HH:mm:ss（UTC+8）`；架构仍需保留权威时间语义，不能由浏览器本地时区决定排序或事实相等性。

## 要解决的问题

用户在聊天中可能连续执行多次查询，然后使用“这些漏洞”“刚才业务系统 A 的漏洞”“第二次和第五次结果”等自然语言引用历史结果。AI 可以理解语义，但外部写入不能依赖模型记忆、压缩后的聊天文本或重新生成的对象集合。

## 已确认的责任边界

- AI 负责把自然语言指代映射为一个或多个安全 `result_ref`。
- 程序不按“这些”“刚才”等关键词硬编码业务逻辑，只校验 AI 提交的结构化引用。
- 程序负责验证引用、解析对象集合、冻结动作草案并绑定审批。
- 用户确认的是冻结后的 `action_draft_ref`，不是重新解释的聊天消息。
- 写工具只能执行已批准的动作草案，不能自行扩大、替换或重新查询目标集合。
- 上下文可以压缩；查询结果集合、动作草案、审批和回读事实不能依赖聊天文本恢复。

## 建议增加的 Product Facts

### QueryResultSnapshot

至少包含：

- `result_ref`
- `workspace_id`
- `actor_ref`
- `conversation_ref`／`run_id`
- `entity_type`
- `query_summary`
- `query_fingerprint`
- `item_refs`
- 安全行级事实引用
- `item_count`
- `observed_at`
- `expires_at` 或失效策略
- `permission_scope_hash`
- `items_digest`

每次查询生成独立、不可变的结果快照。连续五次查询对应五个不同 `result_ref`；系统可以另外维护 `active_result_ref`，但不能覆盖历史快照。

### ActionDraft

至少包含：

- `action_draft_ref`
- `source_result_refs`
- 合并、去重后的冻结 `item_refs`
- `target_action`
- 目标负责人／状态／期限等安全参数引用
- `items_digest`
- `created_by`
- `created_at`
- `expires_at`
- `approval_ref`
- 幂等键
- 当前状态

动作草案一旦进入审批，目标集合和目标参数不可变。任何变化都必须生成新草案并重新确认。

## 上下文压缩与恢复

模型上下文只注入最近相关的结果索引，例如：

```json
[
  {"ref":"result:3","summary":"业务系统 A，19 条"},
  {"ref":"result:5","summary":"今天发现的漏洞，11 条","active":true}
]
```

完整对象集合保存在 Product Facts repository，不进入模型上下文。压缩或进程重启后，context builder 根据当前 pending、最近安全结果引用和用户消息重新装配安全上下文；不得从聊天摘要重新提取对象 ID。

## 动作生命周期

1. 查询工具写入 `QueryResultSnapshot`。
2. AI 从安全结果索引选择 `source_result_ref`，调用动作准备 capability。
3. 后端校验引用的存在性、归属、实体类型、权限、有效期和完整性。
4. 多结果集合并时，后端确定性去重并生成新的派生快照或直接生成冻结动作草案。
5. 后端生成 `ActionDraft` 和审批事实。
6. 聊天审批卡展示查询摘要、对象数量、目标变化及完整对象入口。
7. 用户拒绝、取消或审批过期时，外部写入数为零。
8. 用户批准后，后端按 `action_draft_ref` 加载冻结集合；不得再次由 AI 解释范围。
9. 写入前重新校验权限和易变业务状态；若集合或关键事实已变化，草案失效并要求重新查询／确认。
10. 写入后逐条回读并写入同一 Product Facts／audit／replay 事实链。

## 歧义规则

- 紧邻单一查询结果的“这些漏洞”可映射到 `active_result_ref`。
- 用户明确引用结果序号、查询摘要或多个结果时，AI提交对应引用；后端不解析自然语言。
- AI未提交引用、引用不存在、候选不唯一、引用失效或实体类型不匹配时，不生成动作草案，返回结构化 clarification。
- clarification 只展示安全摘要、数量和时间，不展示内部 UUID 或 provider 原始数据。

## UX 必须呈现的行为

- 聊天查询结果显示用户可理解的结果序号、摘要、数量和时间。
- 写操作确认卡显示所引用的查询结果、冻结对象数量和目标变化。
- 百条级完整明细可从聊天摘要／审批卡进入专门明细表面，但动作仍由聊天发起。
- 引用歧义时由 AI 发起澄清；引用失效或事实变化时明确要求重新确认。
- 技术 `result_ref`、`action_draft_ref` 和内部 UUID 不直接显示给用户。

## 当前实现缺口

当前 `facts.StructuredResultRef` 仅保存 schema version、引用和安全摘要；`StructuredResultCandidate.ItemCount` 不进入 Product Facts。现状不足以持久化多轮行级结果集合，也无法让审批绑定不可变动作草案。

## 架构阶段必须落档

1. 新增或扩展 ADR：多轮结果快照、动作草案与审批绑定。
2. 更新 Product Facts 模型、repository interface、SQLite schema 和幂等事务边界。
3. 更新 StructuredResult／safe projection schema，禁止 raw provider payload 进入快照。
4. 更新 conversation context、approval、audit、replay 和 checkpoint/resume 契约。
5. 更新 capability input schema：写动作准备必须携带 `source_result_ref`，执行只能携带 `action_draft_ref`／审批引用。
6. 增加上下文压缩、五次查询歧义、多结果合并、引用失效、权限变化、事实变化、重复审批和进程重启恢复测试。
