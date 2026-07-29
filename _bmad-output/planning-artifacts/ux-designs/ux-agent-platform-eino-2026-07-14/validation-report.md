# Validation Report — FOBrain 漏洞处置实验 Agent

- **DESIGN.md:** `DESIGN.md`
- **EXPERIENCE.md:** `EXPERIENCE.md`
- **Run at:** 2026-07-14T14:14:52+08:00

## Overall verdict

整体为 **thin**。视觉 token、四条目标写流程、核心安全原则和两表面信息架构已经形成稳定骨架，但还不能作为无歧义的实施合同：当前唯一可启用的只读流程没有独立闭环，628 条浏览、长字段和人员选择仍缺少确定交互，审批恢复与结果未知状态尚未闭合，选定视觉稿也未进入正式引用链。

三位评审均未发现 critical。它们合计报告 22 个 high、15 个 medium、4 个 low 观察项；其中大量指向同一根因。去重后，本报告收敛为 14 个 high、11 个 medium、3 个 low 待处理主题。

## Category verdicts

- Flow coverage — adequate
- Token completeness — strong
- Component coverage — thin
- State coverage — adequate
- Visual reference coverage — thin
- Bloat & overspecification — strong
- Inheritance discipline — adequate
- Shape fit — adequate

## Findings by severity

### Critical (0)

无。

### High (14)

**[Flow coverage] — 当前只读版本没有独立 Key Flow**（`EXPERIENCE.md` Foundation、Key Flows）  
现有四条流程都以门禁后的写入与回读为高潮，无法直接抽取“查询—摘要—查看全部—写意图阻断”的当前版本验收旅程。  
Fix：新增具名的当前只读流程，以客观零写入的 ST-11 结束。

**[Flow coverage] — 子集分组和多查询合并缺少可见契约**（`EXPERIENCE.md` Interaction Primitives、PI-002）  
纯聊天入口下，不同负责人子集、多个查询来源、重叠对象和最终冻结总数尚未形成确定交互。  
Fix：规定聊天中按查询序号和行级安全事实形成子集；审批显示全部来源、去重后总数和最终冻结集合。

**[State coverage / Accessibility] — 628 条确定性浏览与键盘模型未定**（OQ-01）  
分页、虚拟化、位置恢复、局部失败以及 628 个 Tab stop 的处理都会影响当前只读实现。  
Fix：在只读实现前固定浏览模型、稳定位置、语义表格/roving focus 方案和加载状态。

**[State coverage] — 长 POC、多 IP、日期时间和完整值入口未定**（OQ-03）  
当前无法客观验收高密度事实是否完整可读。  
Fix：固定摘要、展开/详情、复制、时区和日期格式规则，不能只依赖省略号或 hover。

**[Component coverage / Accessibility] — 243 人选择没有独立组件**（OQ-02）  
人员选择被塞进消息气泡，缺少搜索、同名消歧、键盘模型、状态和最终复核。  
Fix：新增“人员选择卡”到两份 spine；继续作为 PI-002/PI-004 写域门禁。

**[Visual reference / Safety] — 选定视觉稿混合当前只读态与目标写态**（`.working/direction-precision-command.html`）  
页头表示只读已验证，却同屏展示可确认写入，且状态色与正式语义色冲突。  
Fix：明确标注目标态，并补当前 ST-11 关键屏；提升前按 DESIGN.md 语义色重绘。

**[Accessibility] — 深色左栏焦点环不可见，焦点顺序存在冲突**（`DESIGN.md` Colors；`EXPERIENCE.md` Accessibility Floor）  
深色焦点环与左栏约 1.09:1；输入框在 DOM 中属于中区，但文档把右栏排在输入框之前。  
Fix：新增暗背景焦点 token；采用 DOM 顺序，禁止正 tabindex，并规定新审批卡的可达策略。

**[State coverage / Safety] — 审批过期、重复确认和重启恢复未闭合**（State Patterns）  
陈旧审批可能继续看似可确认，重复提交与重启后的 waiting 卡缺少权威可见状态。  
Fix：增加已过期、提交锁定/幂等回显、Product Facts 恢复后复核三类状态。

**[Recovery] — 确认后断线缺少“结果待核验”状态**（State Patterns）  
确认请求已发出后不能断言零写入或失败，也不能允许再次执行。  
Fix：保留冻结范围、禁用再次确认，恢复后从同一执行事实链对账并逐条结算。

**[Accessibility] — 批量 live region 缺少播报协议**（Accessibility Floor）  
628 条逐行播报会轰炸读屏，只更新颜色又会漏掉关键失败。  
Fix：进度使用节流后的批次摘要；最终只播报完成/待排查/待核验计数，逐条事实进入只读明细。

**[Safety] — 执行记录没有最小安全审计字段**（事实与执行面板）  
用户无法核对谁确认了哪个冻结范围；直接显示内部引用又违反脱敏边界。  
Fix：固定动作、确认人、时间、来源查询、对象数量、目标变化、阶段、结算计数和安全错误摘要。

**[Recovery] — 历史结果和待排查对象没有重开入口**（OQ-05）  
部分失败和回读不一致可能只存在于瞬时卡片。  
Fix：明确跨会话保留范围和只读入口；历史结果不得自动恢复为可确认草案。

**[Inheritance discipline] — 产品 AC 与 UX 的“前 5 条 + 查看全部”不一致**（PI-002、PRD、追踪矩阵）  
来源仍要求确认时逐条展示，UX 已记录为卡片默认前 5 条、全部对象可在确认前核对。  
Fix：同步产品来源与验收追踪，避免两个权威合同。

**[Visual reference coverage] — 选定方向尚未正式提升和内联引用**（`.working/direction-precision-command.html`）  
下游消费者目前无法确认用户选中的视觉方向。  
Fix：修正后提升到 `mockups/`，两份 spine 内联说明其用途，并注明 spine 优先。

### Medium (11)

1. **UJ-05 与 UJ-01～UJ-04 的追踪不够直接。** 在 Key Flow 标题写出 UJ ID/名称，并增加跨动作异常处置流程或完整映射。
2. **Dialog/Popover/Tooltip 只被允许、未被规范。** 首版不用则删除允许性描述；使用时必须补视觉与行为契约。
3. **IA-02 和人员列表异步状态不完整。** 补首次加载、空、局部失败、刷新后变化，并禁止把旧事实当新事实。
4. **200% 文本缩放尚未闭合。** 卡片应换行且不遮挡按钮；只有明细表自身允许二维滚动。
5. **“管理员或排障”不是可执行授权规则。** OQ-08 解决前隐藏技术信息入口。
6. **失败态没有已确认的人工下一步。** 稳定提示“请联系 FOBrain 系统负责人排查”，不新增通知或重试。
7. **三栏缺少 landmark 与跳到主内容。** 规定命名 `nav`、`main`、`aside` 和 skip link。
8. **关键事实变化与恢复后的状态分流偏薄。** 显示变化类别、受影响数量和零写入结论；恢复只依据 Product Facts。
9. **查询与后续页失败缺少局部恢复。** 区分整次查询失败与明细后续加载失败，不承诺离线写入。
10. **取消、未确认、过期、待核验的用户语义混合。** 分开文案和权威状态；结算统一为完成/待排查/待核验。
11. **缺少 Inspiration & Anti-patterns。** 简要记录选择精密指挥台、拒绝暖色编辑感和暗色监控感的原因。

### Low (3)

1. **视觉稿的非语义 DOM、小字对比和前 3 条示意不能作为交互实现。** 正式稿使用原生/Radix 语义，并与“前 5 条”一致。
2. **派发、转发、误报成功文案未统一。** 仅在回读成功后使用“派发完成/转发完成/已标记误报”。
3. **人员状态字段是否展示尚未决定。** 并入 OQ-02，不静默隐藏其影响，也不默认全部可选。

## Reviewer files

- `review-rubric.md` — 0 critical / 7 high / 5 medium / 0 low
- `review-accessibility-safety.md` — 0 critical / 8 high / 5 medium / 2 low
- `review-edge-cases-recovery.md` — 0 critical / 7 high / 5 medium / 2 low
