# ADR: Pin Eino Version

状态：Accepted
日期：2026-06-30

## 背景

Eino-first 重构依赖 ChatModelAgent、Runner event、tool loop、interrupt、checkpoint 和 callback。HITL 与恢复语义会受 Eino API 版本影响，因此引入 runtime 前必须固定 module 版本，并记录升级门禁。

## 决策

Phase 1 先在 `go.mod` 固定：

- `github.com/cloudwego/eino v0.9.12`

版本固定通过 `internal/einoapp/execution/eino_version_pin.go` 维护。该文件使用 `eino_version_pin` build tag，不进入默认构建；Phase 3 实现执行层时必须移除该 pin 文件，并由真实 execution 代码直接引用 Eino。

暂不引入以下模块，等 Phase 3/模型 provider 接入时再按实际代码固定：

- `github.com/cloudwego/eino-ext/components/model/openai`
- `github.com/eino-contrib/jsonschema`

## 验证记录

2026-06-30 执行：

```bash
go list -m -versions github.com/cloudwego/eino
go list -m -versions github.com/cloudwego/eino-ext/components/model/openai
go get github.com/cloudwego/eino@v0.9.12
go list -f '{{.ImportPath}} {{.Name}}' github.com/cloudwego/eino
```

结果：

- `github.com/cloudwego/eino` 可见最新稳定版为 `v0.9.12`，同时存在 `v0.10.0-alpha.*`；本项目不采用 alpha。
- `github.com/cloudwego/eino-ext/components/model/openai` 可见最新到 `v0.1.13`，尚未引入。
- `github.com/cloudwego/eino` 根包可导入，包名为 `eino`。

## 约束

- 前端、Action API、audit、replay 不得直接消费 Eino event，只消费 Product Facts 投影。
- Checkpoint/interrupt id 不得暴露给 Workbench 或 Action API。
- Eino CheckPointStore 只依赖 bytes `Get` / `Set`；进程重启恢复、waiting 保留、checkpoint missing 安全错误和 duplicate resume 幂等由项目 store 与 Product Facts 实现。
- 升级 Eino 必须新增 ADR 或更新本 ADR，并重跑 HITL、checkpoint、tool loop、callback 验收。

## 备选方案

- 跟随 latest：拒绝，API 漂移会影响 HITL、checkpoint 和 event 投影。
- 采用 `v0.10.0-alpha.*`：拒绝，当前阶段不接受 alpha API 作为重构基线。
- 只在 `go.mod` 固定不写 ADR：拒绝，升级风险和验证命令不可追踪。

## 回滚条件

- Runner event 无法稳定生成 Product Facts。
- interrupt/checkpoint 无法满足进程重启恢复。
- tool loop 无法稳定接入 Capability Registry。
- callback 无法只作为 tracing、latency、metrics 和 internal diagnostics。

## 后续验收

Phase 3 实现 execution 后必须补齐并通过：

```bash
go test ./internal/einoapp/execution -run 'Chat|RunnerEvent|AgentToolLoop' -count=1
go test ./internal/einoapp/store/sqlite -run CheckPointStore -count=1
go test ./internal/einoapp/execution -run 'Interrupt|Resume|Callback' -count=1
```
