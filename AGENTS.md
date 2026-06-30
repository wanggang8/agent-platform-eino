# AGENTS.md

本项目是 Eino-first 全新重构实现。目标不是迁移旧代码，而是基于 `docs/` 重新开发一套新的 Agent Workbench。

## 目录角色

- 当前目录：新项目实现目录，可以修改。
- 旧项目目录：`/Users/vick/Desktop/project/ai-agent`，只读参考，不得修改。

## 开发原则

- 以 `docs/README.md` 为入口。
- 开发前必须执行 `docs/pre-development-validation.md`。
- 旧项目只能作为产品能力、验收基线和安全边界参考。
- 不复制旧 runtime 类型、旧执行链路、旧 Workbench 接口、旧 UI DOM 结构。
- 新实现必须使用 Eino-first 设计。
- Workbench 和 Action API 必须共用 Product Facts。
- 工具结果必须以 StructuredResult 作为唯一事实材料。
- Fobrain 能力恢复必须按 `docs/fobrain-tool-matrix.md`。

## 每个任务的工程规范

- 先确认任务所属 Phase、对应设计文档、实施计划和验收门禁，再改代码。
- 代码分层必须清晰：`httpapi` 只做 API 边界，`execution` 只做 Eino 执行编排，`facts` 只做 Product Facts，`product` 只做安全投影和产品映射，`capabilities` 只做能力注册和 adapter，`providers/*` 只做业务 provider。
- 新代码必须职责单一、命名稳定、可测试；优先复用已有 schema、fixture、Product Facts、StructuredResult 和 Safety Gate，不新增平行事实模型。
- 不把业务 provider、Workbench 展示、Action API 响应、audit、replay 互相耦合；所有产品出口必须从 Product Facts 投影。
- 不按自然语言关键词硬编码工具路由；意图识别和能力选择遵循 `docs/intent-and-capability-selection.md`。
- 模型上下文只能来自 safe Product Facts / StructuredResult，遵循 `docs/conversation-context.md`。
- cancel、stop、timeout、retry 和终态幂等遵循 `docs/run-lifecycle.md`。
- callback、telemetry、budget 不得绕过 Product Facts 或 Safety Gate，遵循 `docs/observability-and-budgets.md`。
- 改动代码、契约或行为时，同步更新相关设计文档、schema、fixture、ADR、验收计划或验收记录；不得让实现与 `docs/` 脱节。
- 代码完成后的交付说明、验收记录和面向协作者的文档更新使用中文；代码标识符保持清晰英文命名。
- 代码、测试、脚本完成后必须跑对应任务级检查；最终回复使用中文说明改动、验证结果和剩余风险。
