# UX 独立审查——可访问性与安全写操作

## 结论

整体安全主干是可靠的：聊天作为唯一动作入口、写前二次确认、冻结对象、权限复核、逐条回读、部分失败不伪装成功、Product Facts／StructuredResult 单一事实链，以及 raw provider 数据不进入界面，均已明确写入契约。

当前不存在必须推翻信息架构的 **critical** 缺陷，但有 8 项 **high** 缺口会让下游实现产生不可见焦点、陈旧审批仍可点击、重复确认状态不确定、批量行键盘操作失控，或直接照着视觉稿实现出与门禁相冲突的写入口。它们应在 UX 定稿或对应实现 Phase 开始前关闭。

严重度计数：**critical 0 · high 8 · medium 5 · low 2**。

## 审查范围与证据

- `DESIGN.md`、`EXPERIENCE.md`、`.memlog.md`。
- `.working/direction-precision-command.html`、`.working/result-reference-architecture-handoff.md`、`.working/source-extract.md`。
- 两份 spine 的 9 个 frontmatter sources，重点核对 PI-002～PI-005、实施就绪门禁和需求追踪矩阵。
- 项目契约：`docs/approval-flow.md`、`docs/facts-contract.md`、`docs/provider-policy-and-credentials.md`。

## High

### H-01：所选视觉稿同时表达“只读已验证”和“可确认写入”，会误导门禁实现

所选方向的页头显示“只读查询已验证”，同一屏却提供可点击的“确认执行”写操作卡（`.working/direction-precision-command.html:98,107-112`）。这与当前 UX 契约“写域尚未启用时不得提供可操作审批卡”（`EXPERIENCE.md:29-32,77,100,138`）以及实施门禁“当前不允许进入外部写入实现”（`../../product-blueprint/implementation-readiness-gate.md:9-24`）冲突。虽然 spine 声明优先于 mock，但该稿一旦作为实现视觉参考，极易被误读为首版当前态。

*Fix：* 在视觉稿显著标注“门禁后的目标态，不代表当前可用能力”，并另留一个当前态关键屏：收到写意图时显示 ST-11，无确认按钮。写域关键屏只有在适用 `G-WRITE-*`、`G-FACT-01`、`G-SAFE-01` 通过后才可作为实现验收参考。

### H-02：全局深色焦点环在深色会话导航上不可见

`{colors.focus-ring}` 为 `#111827`，会话导航背景为 `#18212D`（`DESIGN.md:12-17,24-27`）；两者对比度约 **1.09:1**，远低于可见焦点图形通常需要的 3:1。`DESIGN.md:259` 又将这一深色 2px 轮廓作为统一键盘焦点样式，因此左栏会成为焦点不可见区域。

*Fix：* 增加暗背景专用焦点 token（例如 `focus-ring-inverse`），在会话导航、深色按钮等暗表面使用可达 3:1 的浅色双层或外轮廓；保留 `outline-offset`，并在视觉验收中覆盖左栏每类可交互项。

### H-03：焦点顺序与 DOM／组件归属冲突，新审批卡对键盘用户难以到达

Accessibility Floor 规定“会话导航 → 中区主内容 → 事实与执行面板 → 消息输入框”（`EXPERIENCE.md:128`），但消息输入框属于 IA-01 中区（`EXPERIENCE.md:81`），所选稿 DOM 也是 `nav → center（含 composer）→ aside`（`.working/direction-precision-command.html:86-122`）。若实现者为满足文档而使用正 `tabindex`，会制造脆弱的非 DOM 焦点顺序。更关键的是，用户在输入框提交写意图后，新审批卡插入输入框之前；契约一方面禁止动态内容抢焦点（`EXPERIENCE.md:134`），另一方面没有提供“转到待确认操作”的可访问路径，用户只能反向遍历时间线寻找高风险确认。

*Fix：* 明确 DOM 顺序即焦点顺序，禁止正 `tabindex`。对用户主动提交后生成的查询／澄清／审批结果，定义唯一焦点策略：要么把焦点移至卡片标题并播报，要么保持输入框焦点但提供紧邻输入框的“转到待确认操作”控制；后台更新只用 live region，不抢焦点。同步消除 `EXPERIENCE.md:128` 的矛盾顺序。

### H-04：628 条明细的“整行可聚焦 + 所有交互 Tab 到达”会形成超长 Tab 序列

漏洞事实行被定义为整行可聚焦和选择（`EXPERIENCE.md:76`），所有可交互元素又要求由 `Tab` 到达（`EXPERIENCE.md:116`）。对实测 628 条结果，这会让键盘用户经过数百个 Tab 停靠点；OQ-01 只讨论分页／虚拟滚动，没有确定行级键盘模型（`EXPERIENCE.md:186`）。这会直接阻碍用户核对冻结范围和部分失败项。

*Fix：* 在 OQ-01 定案时同时固定可访问集合模式：优先使用语义表格且行本身不进入 Tab 序列，提供单一“查看行详情”动作；如确需交互式 grid，则采用 roving tabindex，并定义方向键、Home/End、PageUp/PageDown、当前行播报、`aria-rowcount`／`aria-rowindex` 与虚拟化焦点保留。不得让 628 行全部成为 Tab stop。

### H-05：243 人选择被塞进“消息气泡”，缺少独立的安全选择组件与读屏语义

目标负责人必须从 FOBrain 243 人全量列表中人工选择，不能由 AI 推荐或手填；当前仅在“消息气泡”行为中描述为中性选择（`EXPERIENCE.md:73`），没有独立组件、焦点模型、同名消歧或选中结果复核。OQ-02 已确认这些规则未决（`EXPERIENCE.md:187`）。这是 PI-002／PI-004 写域的关键输入，模糊实现会导致选错人员或读屏无法确认最终接收人。

*Fix：* 将 OQ-02 保持为写域硬门禁。定案后新增独立“人员选择控件”组件契约，明确 combobox/listbox（或等价 Radix Primitive）的搜索、键盘操作、同名消歧字段、无结果、加载失败与最终选中复核；审批卡必须再次显示该安全人员事实。不得用消息气泡承担 243 项交互控件。

### H-06：审批过期、重复确认和重启恢复缺少明确 UI 终态

UX 已覆盖取消、执行中、引用失效和一般会话恢复（`EXPERIENCE.md:91,101-110`），但未单独规定 approval `expired`、重复 approve／已消费、重启后 waiting 审批重新可用的界面状态。项目审批契约明确要求这些状态并规定重复提交不得再次 mutation（`docs/approval-flow.md:44-75,112-126`），事实契约也要求 resume 单次消费和幂等事务（`docs/facts-contract.md:148-164`）。缺失 UI 规则会让陈旧卡仍像可确认、双击后结果不确定，或重启后误把历史摘要当可执行审批。

*Fix：* 在 State Patterns 增加三类确定状态：① 已过期／已失效：禁用确认并明确“未写入，需重新准备”；② 提交中／重复提交：首次提交立即锁定控件，重复响应投影当前终态，不再执行；③ 重启恢复：只从 Product Facts 恢复 waiting 卡，并在权限、草案摘要和有效期复核后才启用确认。视觉上不能继续使用“等待确认”的琥珀态伪装终态。

### H-07：live region 只有原则，没有批量写入的播报粒度，可能沉默或轰炸读屏

`EXPERIENCE.md:129` 只要求使用“适当的 live region”，未定义查询中、待确认、执行中、回读中、部分失败的 politeness、原子性和批量策略。628 条逐条回读若逐行播报会轰炸读屏；若只更新颜色或右栏则关键失败可能完全沉默。错误后强制焦点（`EXPERIENCE.md:133`）又可能与 live region 重复播报。

*Fix：* 固定播报协议：进度和待确认用单一 `role=status`／polite 摘要；用户提交导致的阻断或最终失败使用一次明确错误摘要；批量回读只播报节流后的批次计数与最终“完成 X、待排查 Y”，逐条结果由用户进入 IA-02 阅读。规定焦点移动和 live announcement 二选一，避免同一事件重复播报。

### H-08：“执行记录”没有规定安全审计字段，无法从界面核对谁确认了哪个冻结范围

右栏只规定显示“事实和执行记录”（`EXPERIENCE.md:80`），没有列出写操作最小可见字段。项目契约要求审批请求、审批结果、mutation 结果进入 audit/replay（`docs/approval-flow.md:72-76`），AuditEvent 提供 actor、safe summary 和 created_at（`docs/facts-contract.md:125-146`）。如果只显示泛化的“已完成”，操作人无法核对确认者、确认时间、冻结对象摘要、部分失败和回读结论；若实现者直接展示内部 audit/ref，又会违反脱敏边界。

*Fix：* 为当前上下文的执行记录固定安全字段：业务动作、实际确认人安全显示名／角色、确认或取消时间、来源查询序号与观察时间、冻结对象数量、目标变化摘要、执行／回读阶段、完成／待排查计数及安全错误摘要。明确不展示 `result_ref`、`action_draft_ref`、`resume_ref`、checkpoint、raw args 或 provider 字段。跨会话查找仍保留在 OQ-05，不扩张首版范围。

## Medium

### M-01：整页固定最小宽度与 WCAG 2.2 AA／200% 文本缩放承诺尚未闭合

设计允许低于 1180px 时整页横向滚动（`DESIGN.md:271`），同时承诺 200% 文本缩放不截断关键事实（`DESIGN.md:267`）和 WCAG 2.2 AA（`EXPERIENCE.md:127`）。这不足以说明聊天、审批卡和错误摘要在缩放后不会同时产生横纵双向滚动。桌面限定不等于可以忽略浏览器缩放。

*Fix：* 增加 200% 缩放验收：聊天、审批、澄清和结果卡在中栏内换行且不遮挡按钮；仅真正需要二维布局的明细表允许自身水平滚动，不能让关键确认和错误信息依赖整页横向来回滚动。这是桌面无障碍行为，不要求新增移动端布局。

### M-02：技术信息的授权主体“管理员或排障”不可执行

`EXPERIENCE.md:80` 使用“管理员或排障”描述可见范围，但“排障”不是可判定角色；OQ-08 也承认角色和脱敏字段未定（`EXPERIENCE.md:193`）。若照此实现，可能扩大 IP、运行与诊断信息暴露。

*Fix：* 保持 OQ-08 为“技术信息实现前”硬门禁；在角色、policy decision 和字段 allowlist 明确前隐藏整个标签，而不是展示空壳或根据前端角色名推断。无权限状态不得泄露该标签是否含目标对象。

### M-03：失败态缺少已确认的人工排查下一步

PI-002、PI-003、PI-004、PI-005 都要求写入失败或回读不一致后由操作人联系 FOBrain 系统负责人排查（例如 `../../product-blueprint/product-issues/PI-002-daily-new-vulnerability-manual-dispatch-draft.md:61-65`、`PI-004-responsible-person-transfer-draft.md:59-69`）。UX 只显示“待排查”和脱敏原因（`EXPERIENCE.md:106-108`），没有告诉用户排查归属。

*Fix：* 在最终失败／部分失败摘要中加入稳定下一步：“未完成，请联系 FOBrain 系统负责人排查。”不要新增通知、工单或重试按钮；这只是已确认流程的可见指引。

### M-04：缺少 landmark／跳过重复导航契约

三栏工作台有固定、重复的会话导航，但 Accessibility Floor 没有规定 `nav`、`main`、补充面板 landmark 或“跳到主内容”。对频繁切换会话的键盘和读屏用户，每次都需穿过左栏。

*Fix：* 在工作台骨架中规定唯一可命名的 `nav`、`main`、补充 `aside`，并提供首个可聚焦的“跳到主内容”。表面切换后主标题作为 `main` 的可编程焦点落点；不增加新业务表面。

### M-05：明细长值“可完整读取”仍被 OQ-03 阻塞

右栏长值截断时要求有可访问完整值入口（`DESIGN.md:301`），明细也要求长内容可完整读取（`EXPERIENCE.md:131`），但 OQ-03 尚未确定超长 POC、IP 列表、时间格式和完整值查看规则（`EXPERIENCE.md:188`）。依赖仅 hover tooltip 会排除键盘、触控板和读屏。

*Fix：* 维持 OQ-03 为只读实现门禁；定案时要求完整值可由键盘聚焦／展开并被读屏读取，截断文本保留 accessible name，复制行为不暴露当前权限外事实。不要把 `title` 属性或仅悬停 tooltip 作为唯一入口。

## Low

### L-01：视觉稿中的状态用色与主契约相反，且存在小字对比失败

所选稿把用户气泡和发送按钮染蓝，把在线染绿、离线染红、缺失染琥珀（`.working/direction-precision-command.html:44,65,69`），而 `DESIGN.md:253-257,294,297,302-303` 明确要求蓝色仅用于查询／导航、在线离线和缺失使用中性表达、发送按钮使用中性深色。实测视觉稿小字也有失败：`#74879B`／`#18212D` 约 4.39:1，active session 的 `#778A9D`／`#263545` 约 3.52:1，`#64748B`／`#F3F5F7` 约 4.35:1，均低于普通文本 4.5:1。

*Fix：* 在视觉稿晋升为 mockup 前按 spine 替换颜色并重新计算实际文本／背景组合；不要只验证 token 各自存在。

### L-02：视觉稿不是可交互原型，但控件形态可能被误复制为无语义元素

“新建会话”、会话、查看全部、标签和发送在 HTML 中使用 `div`／`span`，没有原生按钮／链接／tab 语义（`.working/direction-precision-command.html:88-94,104,114,117`）。作为纯视觉方向可接受，但若后续直接 image-to-code 或复制 DOM，会失去键盘和读屏行为；同时稿中只展示前 3 条，而主契约要求前 5 条（`.working/direction-precision-command.html:110-111`；`EXPERIENCE.md:77,101`）。

*Fix：* 晋升时注明“仅视觉参考，不是 DOM／交互参考”；最终关键屏若承担交互验收，应使用原生 button/link 或 Radix 对应 primitive，并将示例改为前 5 条或明确为视觉截断。

## 已通过的关键安全项

- 聊天是查询与动作的唯一入口，IA-02 明确只读，行选择不改变冻结范围（`EXPERIENCE.md:23,38-43,75,119`）。
- 所有写动作必须二次确认；确认／取消只在审批卡，确认前与取消后均为零写入（`EXPERIENCE.md:77,101-103,120`）。
- 自然语言指代由 AI 理解，但程序校验安全引用并冻结 ActionDraft；批准后不得重新解释范围（`.working/result-reference-architecture-handoff.md:13-20,77-88`）。
- 权限、引用有效性和关键事实变化会使草案失效并要求重新查询、重新确认（`EXPERIENCE.md:97,110,122`）。
- 成功只由回读事实决定；部分失败和回读不一致逐条保留，不以接口成功或批次成功覆盖（`EXPERIENCE.md:104-109,149-180`）。
- 缺失负责人／业务事实有确定性文案，AI 不补全；负责人由 FOBrain 人员列表人工选择，不做推荐（`EXPERIENCE.md:54-56,73,95-98`）。
- Product Facts／StructuredResult 是唯一事实材料，技术引用、凭据、raw payload 和 provider 私有字段不进入主界面（`EXPERIENCE.md:23,32,41,80`；`DESIGN.md:247`）。
- 状态不只靠颜色表达，审批确认与取消要求完整可辨名称，表面切换和错误具有焦点恢复原则（`EXPERIENCE.md:125-134`）。

