# Acceptance Matrix

所有场景先使用 schema/fixture/mock；真实 FOBrain 写入只有 readiness gate 明确授权后才能执行。

| AC | Capability | Given | When | Then |
| --- | --- | --- | --- | --- |
| AC-01 | CAP-1 | 分别准备空数据库和非目标 schema epoch 数据库 | 服务启动并尝试创建新 Run | 空库直接建立目标 v2 schema 并可创建 Run；非目标 epoch 返回 `unsupported_schema_epoch`、readiness=false，且不创建兼容投影或草案 |
| AC-02 | CAP-1 | 一个安全 StructuredResult 含完整多行事实 | ToolResult 提交并在进程重启后读取 | schema/ref/digest/大小和完整安全 payload 一致，所有出口由同一值投影 |
| AC-03 | CAP-2 | 查询有多页结果且所有页成功 | 发布完整快照 | coverage 为 complete_set、稳定顺序和总数可验证，Workbench 每页100条不改变集合 |
| AC-04 | CAP-2 | 多页查询中任一页失败 | 尝试发布“全部”结果 | 只能形成 incomplete/partial 事实，系统拒绝全量 ActionDraft |
| AC-05 | CAP-3 | 当前上下文只有一个有效 result_ref | 用户说“这些漏洞” | resolver 唯一锁定该快照并记录引用；模型不能替换集合 |
| AC-06 | CAP-3 | 当前上下文有多个合理结果 | 用户使用模糊指代 | Run 进入 clarification，未创建 ActionDraft |
| AC-07 | CAP-3 | 摘要已压缩且保留 result_ref | 用户继续引用旧结果 | 系统从事实库重载并校验；过期、跨 actor/conversation 或不存在时拒绝 |
| AC-08 | CAP-4 | 有有效 complete_set 和人工动作参数 | PrepareAction | 生成一份不可变 draft、一个 Pending、一个 approval ref，并展示前5条与完整冻结集合入口 |
| AC-09 | CAP-4 | 用户未确认、拒绝、取消或草案过期 | Run 终结或等待超时 | provider mutation count 为0，audit/replay 保留明确原因 |
| AC-10 | CAP-4 | 确认 payload 的 draft/version/digest/actor 任一不匹配 | ConfirmAction | 返回安全 conflict/denied，approval 不推进，mutation count 为0 |
| AC-11 | CAP-5 | 派发/转发草案包含不同接收人或非全员列表候选 | prepare 或 confirm 校验 | 草案被拒绝或拆分，不能混合执行 |
| AC-12 | CAP-5 | 延时包含多条或 actor 非管理员；误报对象非当前 actor 负责 | prepare 或 confirm 校验 | fail closed 且不触达 provider mutation |
| AC-13 | CAP-6 | 批量动作包含成功、确定失败、预算内回读不一致和预算结束仍未知 | 执行并核验 | 每条独立结果；成功项为完成，确定失败为待排查，预算内未知进入 reconcile 并显示待核验，预算结束仍未知进入 manual_attention 并显示待排查；不整批伪报成功、不重复 mutation 且不发送通知 |
| AC-14 | CAP-6 | 外部写入可能成功但响应丢失 | executor 恢复 | 复用原 mutation/semantic claim，只回读不重写；deadline/预算内无法证明时保持 reconciling/待核验，能证明则按 reducer 终结，耗尽后才 manual_attention/待排查 |
| AC-15 | CAP-6 | 同一确认或 client request 重复提交 | 系统处理重复请求 | 相同 digest 返回既有投影，不同 digest conflict，成功 item 不重复写入 |
| AC-16 | CAP-7 | 相同有效 ActionDraft 分别由聊天与 Action API 创建/确认 | 两条路径执行安全 mock | 产生相同 policy、approval、item、attempt、verification 与产品投影；Action API 无伪 Eino checkpoint |
| AC-17 | CAP-7 | ChatModelAgent 配置 registry tools | 模型选择并调用工具 | 工具循环由 Eino Runner 推进，事件只经 EventMapper 写 facts，预算/迭代/timeout 可控 |
| AC-18 | CAP-8 | Run 正在 waiting approval | 服务重启并使用原 resume ref 确认 | continuation 可恢复，确认只消费一次，attempt 只预留一次 |
| AC-19 | CAP-8 | 旧 executor 失租约后迟到提交 | 新 executor 已取得更高 epoch | repository 拒绝旧 epoch，外部结果进入安全 reconcile 而非覆盖新事实 |
| AC-20 | CAP-9 | 客户端先读取 snapshot_sequence=N | N 后产生多个事实并断线重连 | 客户端只应用递增事件，无遗漏/重复；过旧游标收到 view.replaced 与新高水位 |
| AC-21 | CAP-9 | upsert、tombstone、数组 replace 和 keyed collection fixtures | 前端 reducer 处理重复和乱序事件 | 最终 view 与服务端同源投影一致，不出现幽灵 item 或迟到状态回退 |
| AC-22 | CAP-10 | 使用 Go 1.26.0/1.26.5、Node 24 LTS 和全新数据库 | 执行 Greenfield bootstrap 与验收套件 | schema、OpenAPI、Go、race、TS、stream、browser、backup restore 全通过 |
| AC-23 | CAP-10 | schema epoch/bootstrap、WAL/FK/busy timeout、integrity 或恢复扫描任一失败 | 服务启动 | readiness=false，不接业务流量，不创建 v2 Run |
| AC-24 | CAP-10 | G-ARCH/G-TOOLCHAIN 已通过但某写动作的 G-WRITE/G-SAFE 未通过 | 尝试启用对应真实 capability | capability 保持 disabled/blocked，mock 和只读事实不被误报为写域通过 |
| AC-25 | CAP-2 | FOBrain 存在多页新增漏洞、部分漏洞发现时间相同 | 执行精确新增查询 | 系统完整分页，按发现时间倒序及稳定身份 tie-breaker 冻结集合；查询成功空集、无权和 provider 失败产生三种不同事实 |
| AC-26 | CAP-2 | 查询结果包含字段齐全与字段缺失的漏洞 | 生成统一卡片、确认摘要、ActionResult、replay、audit 和模型上下文 | 所有出口同源展示 POC、IP、在线状态、业务系统、修复负责人；可用时展示业务/运维负责人，缺失值不被模型补造 |
| AC-27 | CAP-5 | 全员列表包含稳定身份；操作人逐条为多条漏洞选择接收人 | 创建派发或转发草案 | 同一人工选择的接收人可合并，不同接收人拆分；系统和模型均未自动推荐或替换接收人 |
| AC-28 | CAP-11 | 当前 actor 有六个月内终态记录、超过六个月但未关闭的记录，以及其他 actor 的记录 | 重新登录并打开固定“操作记录”入口，再进入一条详情 | 默认列表按时间倒序显示六个月内记录，未关闭记录无论时间仍可发现，其他 actor 记录不可见；详情来自原 Product Facts/ActionResult，只读且不调用 FOBrain、不创建 draft 或重试 |
| AC-29 | CAP-6 | 首次 claim 与从 definitive failed 转移的 retry claim 分别处于 `reserved_not_sent`；系统在 sent 事务前后崩溃，或 cancel/stop 与 sent CAS 并发 | 恢复 executor／处理取消 | sent CAS 获胜时只 Verify/Reconcile 且 cancel 返回 action_already_sent；cancel CAS 获胜时进入 cancelled_before_send、mutation count=0，首次 claim 变为 released_zero_write，retry claim 恢复 predecessor owner/status=failed，audit tombstone 保留；两者不可同时成功，deadline/预算/epoch 不伪造或重置 |
| AC-30 | CAP-6 | 分别注入 held+accepted+conclusive match、held+rejected+policy-approved non-application、预算内其他组合、预算耗尽其他组合和 control lost | verifier 收敛 | 唯一映射为 succeeded、failed、reconciling、manual_attention、manual_attention；Run/ActionResult/产品状态和 mutation_retry_allowed 与 AD-10 outcome table 完全一致 |
| AC-31 | CAP-4 | 两个普通 draft 并发确认同 semantic key；retry draft 包含 failed、succeeded、reconciling 或 manual_attention item；另执行 failed→retry reserved→cancelled_before_send 后提交普通 draft 与新 lineage retry | prepare/confirm | Confirm 事务的唯一 SemanticMutationClaim 只允许一个普通 draft 预留；只有逐条绑定 definitive failed 的新人工 retry 可事务转移 claim；取消 transferred retry 后普通 draft 仍因恢复的 failed claim 被拒绝，而新的合法 lineage retry 仍可转移；其余返回稳定 semantic_mutation_conflict/already-satisfied，provider mutation count 为0 |
| AC-32 | CAP-2 | 列表返回 opaque entity_ref 与独立 ProviderLocator | 打开资产或漏洞详情 | resolver 重验 scope/policy/version 后使用 locator；产品 ref 不等于且不能拆出 provider ID，错误资产类型、越权、过期或未知 locator 均 fail closed |
| AC-33 | CAP-11 | 历史记录跨六个日历月边界，包含 waiting/running、definitive failed open attention、manual_attention open attention、attention not_applicable 终态和权限被撤销批次 | 使用服务端 as_of 分页查询 | 后端按明确 Asia/Shanghai 月末 clamp 算法得到 inclusive cutoff，在同一 immutable OperationListSnapshot 上完成窗口与 lifecycle-open union、去重，cursor 绑定 operation_list_snapshot_ref；撤权/unknown 整批 fail closed 且不重算，前端不追加或重排；本 AC 不要求或复用 Run-local SSE snapshot_sequence |
| AC-34 | CAP-6 | 分别构造 initial/retry reserved_not_sent、两种 cancelled_before_send claim disposition、sent、verifying/reconciling、terminal attempt fixture | schema、Go、TS 与存储映射校验 | reserved_not_sent/cancelled_before_send 不含 sent-only 字段；cancelled 分支必含零 mutation、released_zero_write 或 restored_failed_predecessor 的 claim-disposition audit；sent+ 必填 immutable sent_at/deadline/policy/epoch/预算；核验后才允许 evidence/next poll；零值伪造与旧 `reserved` enum 均被拒绝 |
| AC-35 | CAP-6 | VerificationEvidence 覆盖 accepted+match、rejected+match、unknown+match、policy-approved rejection without observed、预算内 mismatch/unknown、耗尽 mismatch/unknown、预算未耗尽但 control lost | 唯一 reducer 执行 | evidence_kind union 条件校验 observed ref/digest；control lost 优先 manual_attention；仅 held+accepted+conclusive real-readback match 成功；仅 held+rejected+conclusive non-application 失败；其余分别 reconciling/manual_attention，manual_attention 聚合 outcome 只能 partial/none_succeeded |
| AC-36 | CAP-4 | 同一语义目标经过字段顺序、Unicode、时区、null/missing、set 顺序变化；另有真正不同目标 | 计算 smk.v1 | 等价输入产生相同 key，不同目标产生不同 key；schema/canonicalizer version 改变必须升级 key version，不能静默碰撞 |
| AC-37 | CAP-11 | 首次请求物化列表后，两个不同 Run 分别改变 lifecycle/attention，另有新增记录、cursor 篡改或跨 actor/filter 复用 | 查询下一页 | OperationListSnapshot 冻结完整有序安全行和 total，独立 ref 明确不是 Run-local sequence；所有页不受跨 Run 变化影响且无重复遗漏；snapshot 缺失/过期或非法 cursor 返回 `invalid_cursor` 且不泄露记录 |
| AC-38 | CAP-2 | 同一 provider 业务对象跨多次查询、workspace subject-derivation key 轮换和 locator encryption key 轮换生成不同 entity_ref，另构造 alias 冲突 | 去重、创建 semantic claim并打开详情 | active+retained fingerprint alias 始终命中同一 CSPRNG 256-bit immutable entity_subject_key、semantic claim 与 locator lineage；旧 key 退役受 alias migration/retention 门禁；多个 alias 指向不同 subject 时 entity_identity_conflict/fail closed；resolver 仅通过 LocatorRepository/Resolver |
| AC-39 | CAP-11 | ActionDraft/Waiting Pending/ActionResult 带 frozen history policy ref，另有缺失、损坏、当前 denied 情况 | 查询操作记录 | ref/version 原样传播；缺失/损坏整批返回 `history_policy_unavailable`，当前 denied/unknown 整批返回 `history_access_denied`，不按 item 过滤或重算 |

## 必须覆盖的失败注入

- staged continuation 后、Pending 发布前崩溃；
- Pending 发布后、Confirm 前重启；
- Confirm 事务提交前/后崩溃；
- provider 收到写入但响应丢失；
- `reserved_not_sent` 与 `sent` 事务前后分别崩溃；
- 回读最终一致窗口内/外变化；
- verification deadline／调用预算在重启前后保持一致；
- lease 过期与旧 executor 迟到提交；
- checkpoint 丢失、损坏、重复 resume；
- SSE snapshot 与订阅之间并发写入；
- 过旧/未知 Last-Event-ID；
- fresh database 重复 bootstrap、legacy epoch 拒绝、显式重建前备份恢复。
- 六个月边界、未关闭记录跨窗口可发现、actor/workspace 隔离和损坏历史引用。
- semantic mutation key 跨 draft 冲突、retry lineage 夹带非 failed item、opaque ref/locator 拆解攻击。
- 两个普通草案并发 claim、semantic canonicalization 等价/非等价输入、history cursor 篡改与跨页高水位变化。
- ActionAttempt phase union 的 sent-only 字段、VerificationEvidence discriminant 与 history policy ref 损坏。

## 验证命令基线

```bash
npm run eino-workbench:schema-test
npm run eino-workbench:contract-test
npm run eino-workbench:openapi-lint
GOTOOLCHAIN=local go test ./internal/einoapp/architecture ./internal/einoapp/facts ./internal/einoapp/execution ./internal/einoapp/store/sqlite -count=1
GOTOOLCHAIN=local go test -race ./internal/einoapp/execution ./internal/einoapp/store/sqlite -count=1
npm run eino-workbench:typecheck
npm run eino-workbench:test
npm run eino-workbench:stream-test
npm run eino-workbench:browser-test -- --project=desktop
```

未实现的阶段测试必须明确 skipped/blocked，不能作为通过证据。
