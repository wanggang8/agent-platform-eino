# Greenfield Bootstrap and Gates

## 当前事实

- 当前过渡 schema 只有 Run、Turn、ToolCall、ToolResult、PendingInteraction、AuditEvent，且 `additionalProperties: false`；它不是目标兼容基线。
- 当前 SQLite 只保存 StructuredResult schema/ref/摘要，不能在重启后恢复完整列表。
- 当前自然语言 ChatModelAgent 未配置完整工具循环；显式 capability runner 直接调用 Eino tool interface。
- 当前审批 continuation 不是 Eino interrupt/resume。
- 当前 SSE sequence 按查询快照重新编号，不是持久化事实游标。
- 当前 SQLite 未形成统一、可验证的 WAL/每连接 FK/busy-timeout policy。

因此所有 v2 决策均为已批准 target，不是当前实现声明。

## 迁移切片

| 顺序 | 必须交付 | 完成门禁 |
| --- | --- | --- |
| M-0 工具链 | Go 1.26.0/1.26.5、Node 24 LTS、CI/版本文件、`npm ci` | 现有 schema/Go/TS/browser 基线在新工具链通过 |
| M-1 契约 | Product Facts v2、Catalog v2、Action/Resume/Result/SSE v2、含 cancelled-before-send 的 ActionAttempt union、判别式 VerificationEvidence、SemanticMutationClaim、EntitySubject/Fingerprint/ProviderLocator、OperationListSnapshot schema、fixtures、OpenAPI、generated types | schema/contract/OpenAPI 全通过；无业务 writer、无旧契约 union |
| M-2 存储 | Greenfield SQLite schema bootstrap、safe StructuredResult value、v2 aggregates、FactEvent sequence、semantic claim unique constraint、subject/fingerprint/locator repositories、OperationListSnapshot 物化 repository、统一 connection policy | fresh bootstrap、legacy epoch 拒绝、并发 claim/cancel CAS、secret rotation identity、跨 Run 列表快照、FK/WAL/busy/backup restore 测试通过 |
| M-3 投影 | 唯一 v2 projection、Run-local SSE snapshot high-water、SSE upsert/tombstone、replay/audit/context、独立 OperationListSnapshot 的操作记录列表与详情 | 重复/乱序/重连/过旧游标、跨 Run 翻页、cursor 篡改/TTL、六个月 month-end clamp、open 补集、actor/history-policy 隔离与安全泄漏测试通过 |
| M-4 Runtime | Eino ToolsConfig/ReAct、真实 chat interrupt/resume、project Action continuation、lease/fencing/reconcile | Eino/continuation/崩溃点/重复确认/预算测试通过 |
| M-5 只读事实链 | 精确新增查询、全员列表、七项代表性能力、统一漏洞卡片、QueryResultSnapshot、entity ref/locator resolver | G-READ-01/02/03 与 G-FACT-01 所需产品证据通过 |
| M-6 写域准入 | 按动作补权限、provider mutation、verifier、部分失败和安全证据 | 仅适用 G-WRITE/G-SAFE、共同 G-ARCH-V2/G-TOOLCHAIN/G-FACT-01 全 PASS 后创建真实写域 Story |

不得把 M-1 至 M-5 分散进四项 FOBrain 写动作 Story；基础能力先独立完成。

## Greenfield schema 规则

- 空库直接创建目标 schema epoch，并允许在相同目标 epoch 上重复启动。
- 非目标 epoch 返回 `unsupported_schema_epoch`、readiness=false；服务不自动迁移、删库或提供兼容读取。
- 操作者显式重建前必须可执行受控备份；备份恢复到临时库并通过 integrity、foreign key 和 schema version 检查才有效。
- Product Facts v2 是唯一 writer 和 reader；不允许旧事实参与目标运行时。
- v2 schema/fixture/OpenAPI/generated contract 全部通过后才能启用 writer。
- Greenfield 基线后的未来 schema 变化仍必须前向、可重复且有备份恢复测试。

## SQLite 启动门禁

`bootstrap` 在监听业务端口前必须完成：

1. 目标 schema bootstrap 与 epoch 校验；
2. `journal_mode=WAL` 并验证返回 `wal`；
3. 每连接 `foreign_keys=ON` 和配置化 `busy_timeout`；
4. integrity/foreign-key check；
5. 可写探针；
6. continuation、lease 和未知写入恢复扫描。

任何一步失败：`/health` 可表示进程存活，但 `/ready` 为 false 且不接业务流量。FOBrain/LLM 不可用只让相关 capability degraded/fail closed，不使事实读取服务失去核心 readiness。

## Retention

- waiting Pending/continuation、running/reconciling Run、`attention_status=open` 及其 snapshot/draft 禁止 cleanup。
- 首版不自动物理删除 terminal Product Facts 或 audit。
- 操作记录默认六个月只限制产品发现窗口；Run waiting/running 或 attention open 的记录始终可发现，不因窗口清理。首版没有 attention 关闭命令。
- snapshot/draft 可逻辑过期；物理清理必须确认无 active ref 并满足配置化 retention。OperationListSnapshot 只保留到配置化分页 TTL，过期 cursor 稳定返回 `invalid_cursor`。
- semantic claim、EntitySubject、EntityIdentityFingerprint、EntityRefMapping 与 ProviderLocator 至少保留到 draft/audit retention、reconcile 窗口、引用材料和 open attention 中最晚者。旧 subject-derivation key 必须等全部保留主体完成新 alias 且上述 retention 结束后才能退役；locator encryption key 可独立轮换但不得改变 subject。

## 实施授权门禁

| Gate | 当前 | 解锁内容 |
| --- | --- | --- |
| G-ARCH-V2 | BLOCKED | v2 writer、稳定 SSE 和 v2 runtime |
| G-TOOLCHAIN | BLOCKED | 下一实现阶段的有效验收环境 |
| G-READ-01/02 | 目标 API PASS、产品 capability BLOCKED | 精确新增查询和全员列表产品能力 |
| G-READ-03 | 七项代表性只读接口已有集成证据；字段完整性、opaque ref/locator 与 v2 产品投影 BLOCKED | FR-19～FR-25 产品能力及其三态/详情 resolver |
| G-FACT-01 | BLOCKED | 统一行级事实和可冻结快照 |
| G-WRITE-01..04 | BLOCKED | 对应真实 mutation 与 verifier |
| G-SAFE-01 | BLOCKED | 写域交付声明 |

门禁通过必须有 `implementation-readiness-gate.md` 指定的机器报告和验收记录；代码存在、接口返回 2xx、mock 通过或 Spine 标记 target 都不能替代。

## 回滚与停止条件

- fresh database 无法稳定创建目标 epoch，或 legacy epoch 无法 fail closed：停止实施。
- 安全 StructuredResult 持久化不能满足大小或脱敏限制：停止快照能力。
- Eino interrupt/resume 不能通过重启、重复确认和 checkpoint 损坏测试：聊天写域保持禁用。
- SSE sequence 不能证明 snapshot-to-stream 无缝衔接：Workbench 只允许轮询完整 view，不声明断线续传。
- 不确定外部写入没有可靠回读 verifier：对应写 capability 不得启用。

## 权威输入

- `docs/schemas/eino_product_facts.v1.schema.json`
- `docs/schemas/eino_action_request.v1.schema.json`
- `docs/schemas/eino_action_result.v1.schema.json`
- `docs/schemas/eino_workbench_resume_request.v1.schema.json`
- `docs/schemas/eino_workbench_stream_event.v1.schema.json`
- `docs/schemas/capability_catalog.v1.schema.json`
- `docs/adr/2026-07-14-product-facts-v2-eino-runtime-migration.md`
- `docs/adr/2026-07-15-greenfield-no-legacy-compatibility.md`

上列 `.v1` schema 是当前实现事实输入，不是目标兼容要求；目标 Story 必须按新契约决定删除、替换或保留，不能据此生成 runtime adapter。
