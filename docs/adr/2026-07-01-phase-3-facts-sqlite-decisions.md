# Phase 3 Facts and SQLite Decisions

## 背景

Phase 3.2 要把 Product Facts 从 Phase 1 内存替身推进到可持久化实现。该实现会支撑 Workbench、Action API、SSE、audit、replay 和后续 resume/checkpoint 边界，因此需要先固定 SQLite、migration、facts 存储和事务策略。

## 决策

- SQLite driver 使用纯 Go 实现，优先避免 CGO 依赖造成本地和 CI 环境差异。
- migration 先由项目内 `internal/einoapp/store/sqlite/migrations.go` 和 schema version 表管理，不引入外部 migration 框架。
- facts 表结构采用核心索引字段列化、扩展详情 JSON 化；raw provider payload、raw checkpoint id、secret、credential ref 和 reusable resume token 不入库。
- SSE cursor 使用同一 run 内单调递增 sequence；SSE event id 从稳定 fact id 和 sequence 派生。
- repository API 按业务事务暴露，不提供通用 `Save(any)`。
- repository 必须提供按 run 读取完整 Product Facts snapshot 的接口，供后续 Workbench、Action API、replay 和 context projection 同源读取。
- resume 消费和 idempotency record 写入必须支持同一 SQLite transaction；mutation idempotency key 使用同一 idempotency 表。
- repository 对 `safe_*`、`args_preview`、`checkpoint_ref` 等入库字段执行最小不安全材料检测，拒绝明显 secret、credential、Authorization、raw provider payload 或 raw checkpoint marker。
- ID 和时间由调用层生成；需要稳定测试时使用显式时间和 id，不在 repository 内生成随机值。
- Phase 3.2 只记录 `checkpoint_ref` 内部安全引用和 pending/resume 事务边界；durable CheckPointStore、approval/clarification interrupt、重启后恢复留到 Phase 6。
- 测试必须真实执行 migration；SQLite repository 测试使用临时数据库，不只测 mock。

## 备选方案

- 使用 CGO SQLite driver：性能成熟，但增加本地和 CI 构建差异，暂不采用。
- 引入 goose/atlas 等 migration 工具：功能完整，但 Phase 3.2 迁移规模很小，暂不采用。
- 全 JSON event store：扩展简单，但查询、幂等和事务断言更弱，暂不采用。
- repository 内部生成 id/time：调用简单，但测试稳定性和审计可解释性更差，暂不采用。

## 影响

- Phase 3.2 可以先形成稳定事实和事务边界，再接 Eino Runner。
- 后续 provider、SSE、audit 和 replay 读取同一 Product Facts，不需要重复事实模型。
- 若事实查询复杂度明显上升，可以在保持 repository 契约不变的前提下新增索引或专用表。

## 回滚条件

- 纯 Go SQLite driver 无法满足目标平台或事务行为要求。
- 项目内 migration 无法支撑 schema 演进审计。
- facts 查询性能或 SQL 复杂度证明当前列化/JSON 混合结构不合适。

## 关联验收

- `go test ./internal/einoapp/facts -count=1`
- `go test ./internal/einoapp/store/sqlite -count=1`
- `go test ./internal/einoapp/architecture -run 'ImportBoundary|ProductLayers' -count=1`
