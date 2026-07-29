---
id: PI-XXX
title: ""
status: draft
primary_user: ""
upstream_business_goals: []
upstream_prd_requirements: []
upstream_capabilities: []
evidence_refs: []
evidence_records: []
assumptions: []
---

# PI-XXX：标题

## Product Issue

用一句话说明要解决的独立业务增量，不写实现方式。

## 证据记录

| 主张／字段 | source_type | claim_support | target_applicability | maturity | decision_status | normative_force | 证据引用 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Background／Goal：... | ... | ... | ... | ... | ... | ... | ... |

## 背景（Background）

### 当前问题

...

### 受影响者

...

### 为什么需要解决

...

### 不解决的影响

...

## 目标（Goal）

### G-01

...

> 只描述业务结果、用户能力和成功后的变化；不描述数据库、接口、框架或技术方案。

## 产品输入／输出契约

### 输入与触发

...

### 正常输出

...

### 非正常与边界输出

- 数据不足：...
- 对象歧义：...
- 无权限：...
- 真实空结果：...
- 外部失败：...
- 事实／推断／未知：...

## 验收标准（Acceptance Criteria）

### AC-01

**映射目标：** G-01

**Given：**

...

**When：**

...

**Then：**

...

### AC-02

**映射目标：** G-01

**Given：**

...

**When：**

...

**Then：**

...

## 非目标

- ...

## 依赖

- ...

## 未决项

- 无；或列出负责人、影响和解决门槛。

## 附件（可选）

- PRD：...
- 原型图：...
- UI 设计稿：...
- 流程图：...
- 接口文档：...
- 历史 Issue：...
- 相关业务文档：...

## 就绪检查

- [ ] 每个 Goal 至少映射一条 AC。
- [ ] 每条 AC 可以客观验证。
- [ ] 覆盖适用的异常、权限和不确定性边界。
- [ ] 上游引用完整且可追溯。
- [ ] 每条关键主张具有完整 canonical evidence record，且未跨层升级。
- [ ] 未泄漏实现方式。
- [ ] 高影响未决项已解决。
