# UX 与视觉需求

本文只描述前端体验和视觉约束，不描述后端实现。

## 总体体验

Workbench 是工作台，不是营销页。第一屏就是可用的聊天工作区。

界面保持三栏结构：

- 左侧：workspace/session/navigation。
- 中间：主聊天时间线和输入框。
- 右侧：Inspector，展示证据、结构化数据、运行状态和审计。

## 时间线内容

主聊天默认只展示：

- 用户消息。
- assistant 过程文本和最终回答。
- 工具卡。
- 审批卡。
- 澄清卡。
- 运行提示。

默认不展示：

- raw provider JSON。
- raw prompt。
- secret、token、Authorization。
- raw credential。
- 内部恢复标识。
- 内部暂停标识。
- 内部 reason code。
- 注册工具 ID。
- model instruction。
- recovery policy。

## 工具卡

工具卡必须包含：

- 安全中文标题。
- 状态：pending、running、completed、failed、cancelled。
- 简短摘要。
- 展开区域。
- 结构化结果视图。
- 失败时的安全错误摘要。

工具卡不得从 raw tool ID 推断业务展示；业务展示必须来自安全结构化结果。

## 审批卡

审批卡必须包含：

- 操作名称。
- 风险摘要。
- 影响对象。
- approve/reject 操作。
- 状态变化：waiting、approved、rejected、expired。
- 重复提交后的稳定反馈。

未审批前，UI 不得暗示 mutation 已执行。

## 澄清卡

澄清卡必须包含：

- 需要用户补充的问题。
- 候选项或输入框。
- 提交、取消状态。
- 已提交后的只读结果。

## Inspector

Inspector 至少包含：

- evidence：安全事件列表。
- structured：选中工具的结构化结果。
- runtime：脱敏运行状态。
- audit：脱敏审计事件。

空态也必须稳定，不允许右栏布局塌陷。

## 响应式要求

Desktop 需要保持三栏工作台密度。

Mobile 需要优先保证：

- 主聊天可读。
- 输入框可用。
- 工具卡和审批卡不溢出。
- Inspector 可以折叠或通过 tabs 进入。

## 视觉防跑偏

视觉目标是保持当前工作台的专业、密集、聊天优先体验，不重新设计成营销首页或演示玩具。

视觉不可变项：

- 第一屏必须是可工作的聊天工作台，不出现 hero、营销说明或引导卡片。
- Desktop 保持左导航、中间聊天、右 Inspector 的三栏信息架构。
- 中间聊天区域是主焦点，工具卡和审批/澄清卡不得抢占输入区。
- Inspector 默认承载 evidence、structured、runtime、audit，不把 raw event 流作为主界面。
- 卡片层级保持克制，不能出现卡片套卡片或大面积装饰背景。

必须覆盖的可见状态：

- 桌面：空态、普通聊天、工具卡折叠、工具卡展开、审批等待、澄清等待、失败态、Inspector 多个标签。
- 移动端：主聊天、输入区、工具卡、审批卡、澄清卡。

视觉不跑偏要求：

- 三栏结构、聊天密度、工具卡层级、Inspector 信息组织保持稳定。
- 工具卡、审批卡、澄清卡在桌面和移动端都不能溢出、遮挡或挤压主聊天。
- 业务结果必须优先可读，装饰性元素不能压过工作流信息。
- 空态、失败态、等待态不能显得像未完成页面。
