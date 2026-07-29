# Product Issue 契约与门禁

本目录收录从 PRD 范围基线投影出的、可独立验收的业务增量。历史工具、页面、接口或实现任务不能直接成为 Product Issue。当前 PI-002～PI-005 已完成业务需求发现，但统一处于 `EXPERIMENT_REQUIREMENT_PENDING_DEPLOYMENT_CAPABILITIES`：它们可用于架构前的需求审阅，不能据此进入外部写入实施。

## 进入条件

候选 Issue 必须先具备：

1. 一个已确认的目标用户和真实任务；
2. 一个上游业务目标与稳定 PRD Requirement ID；
3. 明确的现状问题、影响者和不解决后果；
4. 产品输入、产品输出、空结果、错误、权限和不确定性边界；
5. 至少一个不依赖技术实现的业务结果；
6. 每个业务结果至少一条客观可执行的 Given / When / Then；
7. 明确非目标、依赖、证据来源和仍未解决的问题。
8. 每条关键 Background／Goal 主张都有六字段 canonical evidence record；厂商声明不得代替目标用户或任务证据。

任一项缺失时，候选项留在能力地图或开放问题清单，不创建正式 `PI-*` 文件。

## 字段规则

### Background

必须回答：

- 当前发生了什么业务问题；
- 谁受到影响；
- 为什么现在需要解决；
- 不解决会带来什么可观察影响。

不得用“系统缺少某页面／接口／表”代替业务问题。

### Goal

必须描述：

- 希望实现的业务结果；
- 用户最终获得的能力；
- 成功后用户工作或业务状态的变化。

每个结果使用稳定 `G-01`、`G-02` 编号。不得描述数据库、接口、框架、组件或内部算法。

### Input / Output Contract

这里描述用户和业务可见的输入输出，不是 API 契约：

- 触发条件和用户提供的信息；
- 产品可使用的已授权事实范围；
- 用户获得的结果、证据和下一步；
- 数据不足、对象歧义、无权限、真实空结果和失败时的结果；
- 哪些内容是事实、推断或未知。

### Acceptance Criteria

- 使用稳定 `AC-01`、`AC-02` 编号；
- 每条 AC 显式映射一个或多个 Goal ID；
- Given 描述可构造的前置事实和权限；
- When 描述用户可观察的行为或事件；
- Then 描述可客观检查的结果、状态、证据或禁止行为；
- 不允许“正确”“友好”“快速”“智能”等无阈值形容词；
- 正常、空结果、歧义、无权限、外部失败和高风险动作按适用性覆盖。

## 完成门禁

一个 Product Issue 只有同时满足以下条件，才可标为 `ready`：

- Goal → AC 覆盖率为 100%；
- 每条 AC 可由与实现无关的观察者判定通过或失败；
- 上游目标、PRD 需求和能力引用均存在；本轮采用 [需求追踪矩阵](../requirements-traceability-matrix.md) 显式维护 Goal／PRD／AC 映射；
- `source_type / claim_support / target_applicability / maturity / decision_status / normative_force` 均可追溯，且没有跨层升级；
- 没有未标记假设或相互冲突的要求；
- 没有技术实现泄漏；
- 依赖和高影响开放问题已解决；
- 与其他 Issue 不重复，也不把一个独立业务增量拆成工具级碎片。

## 文件命名

```text
PI-001-short-business-outcome.md
PI-002-short-business-outcome.md
```

ID 一经进入追踪矩阵不得复用；删除时保留 tombstone 记录。
