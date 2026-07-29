# Spine Pair Review — FOBrain 漏洞处置实验 Agent

## Overall verdict

整体为 **thin**：两份 spine 已把核心安全语义、四条目标动作、13 个产品组件、23 个状态和浅色桌面视觉语言组织成清晰骨架，token、来源路径和核心组件名称也能稳定解析。但它还不能作为无歧义的下游实施合同：当前唯一获准的只读范围没有独立闭环，人员选择与不同负责人分组缺少可实施交互，两个只读阻塞问题仍留在 Open Questions，且用户选定的视觉方向尚未进入正式引用链并与 spine 的语义色规则冲突。

## 1. Flow coverage — adequate

检查了 9 个 frontmatter 来源中的 `UJ-01`～`UJ-05`、`PI-002`～`PI-005` 与 `FR-01`～`FR-18`。`PI-002`～`PI-005` 均有具名主人公、编号步骤、Climax 和失败路径；四条动作流程合计覆盖人工判断、二次确认、逐条写入、回读和部分失败。

### Findings

- **[high] 当前唯一获准实施的“精确查询—安全摘要—完整只读明细—写意图阻断”没有自己的 Key Flow 与成功高潮；现有四条 Key Flow 全部被声明为门禁后的目标体验，却继续写到真实写入和回读。下游 story-dev 无法直接抽取当前版本的完整验收旅程。** (`EXPERIENCE.md:27-32`, `EXPERIENCE.md:136-185`, `implementation-readiness-gate.md:9-11,33-35`). *Fix:* 增加一条明确标为“当前可启用”的只读主人公流程，以 ST-11 的客观零写入结论收束；目标写流程继续保留并显式标为后续门禁范围。
- **[high] 手动派发要求逐条判断并允许不同接收人拆分，但 spine 只规定“不同负责人形成不同审批卡”，没有定义用户如何在纯聊天入口中指明不同子集、程序如何形成派生范围；IA-02 又明确禁止选择写入对象。该缺口使 UJ-01 的一般情况无法落地。** (`prd.md:60-65`, `EXPERIENCE.md:75-77,119-122,140-151`). *Fix:* 在 Key Flow 和 Interaction Primitives 中提交一个不依赖 IA-02 写选择的确定性子集引用／分组流程，并规定歧义、重叠和未分配对象的结果。
- **[medium] 来源中的 `UJ-05：处置异常` 没有以原名出现为独立 Key Flow。四条动作的失败路径覆盖了大部分异常规则，但没有具名操作者、编号步骤和“识别待排查对象并交给 FOBrain 系统负责人”的跨动作高潮。** (`prd.md:34-40`, `EXPERIENCE.md:140-185`). *Fix:* 保留现有每条失败路径，同时新增 `UJ-05` 异常处置流程，或在四条标题旁显式标出 UJ 映射并说明 UJ-05 由哪些步骤完整覆盖。

## 2. Token completeness — strong

YAML 中所有颜色 token 都是六位十六进制值；浅色单主题与用户决定一致。全部 `{path.to.token}` 引用均能解析到已定义 token 或已定义的 typography/component 对象。核心正常文字和四种语义色组合均达到 WCAG AA：例如 muted/raised 4.76:1、query/query-soft 4.75:1、pending/pending-soft 4.84:1、success/success-soft 5.84:1、failure/failure-soft 5.76:1。`DESIGN.md:259` 还明确承诺 WCAG 2.2 AA 与非颜色线索。

### Findings

无。

## 3. Component coverage — thin

两份 spine 的 13 个正式组件名称完全一致，且每个组件同时具有视觉规范和行为规范：工作台骨架、会话导航、聊天时间线、消息气泡、查询结果卡、漏洞明细工作区、漏洞事实行、审批卡、澄清卡、执行结果卡、事实与执行面板、消息输入框、状态标记。

### Findings

- **[high] FOBrain 全部人员选择是当前只读准入能力和 PI-002／PI-004 的关键交互，但没有对应的正式组件。它被塞进“消息气泡”的一句行为规则，视觉契约没有定义 243 人列表、候选身份／状态、搜索或同名消歧的承载形态；OQ-02 又确认这些规则未定。** (`implementation-readiness-gate.md:17-19,33-35`, `EXPERIENCE.md:68-82,145,168,194`). *Fix:* 新增同名的“人员选择器／人员选择卡”组件到两份 spine，并在启用写域前提交候选读取、加载、空结果、同名消歧、键盘选取和选中结果规则；若当前只读版本仅验证人员列表读取，也要定义只读呈现组件。
- **[medium] `Dialog`、`Popover`、`Tooltip` 被允许作为真实覆盖层，但 Radix 被明确设为无样式基础，当前既没有这些覆盖层的正式组件行，也没有 EXPERIENCE 行为约束。下游实现会缺少关闭、焦点回归、遮罩和嵌套规则。** (`DESIGN.md:245,279-281,305`, `EXPERIENCE.md:116-123`). *Fix:* 若首版实际使用这些覆盖层，将其作为同名组件加入两份 spine；若不使用，删去允许性描述，避免产生未规范的实现入口。

## 4. State coverage — adequate

IA-01 与 IA-02 已覆盖待命、恢复、查询、空结果、字段缺失、澄清、无权限、数据不足、外部失败、写域不可用、待确认、取消、执行、回读、状态不允许、失败、部分失败、回读不一致、成功、引用失效、离线及焦点等 23 个状态。错误、离线、权限和焦点均未与空结果混淆。

### Findings

- **[high] IA-02 的 628 条级读取方式和超长事实展示仍分别由 OQ-01、OQ-03 标为“只读实现前”阻塞项；这不是可留待实现者自行选择的装饰细节，会决定分页／虚拟化状态、当前位置、返回恢复、表格语义和完整值读取。** (`EXPERIENCE.md:75-76,131,193-195`). *Fix:* 在只读实施前选定确定性浏览方式、定位／返回规则、时间格式与完整值入口，并补充相应 cold-load、增量加载、末页、加载失败和恢复状态。
- **[medium] 首屏恢复状态只标注 IA-01，人员列表只有失败／数据不足状态；IA-02 首次进入和人员候选读取都没有加载中、成功为空或刷新后变化的状态契约。** (`EXPERIENCE.md:90-112`, `EXPERIENCE.md:193-195`). *Fix:* 将异步状态按表面补齐，并明确加载期间不能把历史快照或旧人员列表当成新事实。

## 5. Visual reference coverage — thin

`imports/` 为空，`mockups/` 与 `wireframes/` 不存在，因此两份 spine 当前没有任何正式内联视觉引用。`.working/` 中有三个设计方向，memlog 已记录用户选择 `direction-precision-command.html`，但该文件仍是工作材料，不在正式引用链。

### Findings

- **[high] 用户选定的“精密指挥台”没有被提升到 `mockups/`，两份 spine 也没有在相关 IA／组件段落内联链接并说明它展示什么。下游消费者只能从 prose 重建构图，无法确认选定方向。** (`.memlog.md:29-32`, `.working/direction-precision-command.html:1-126`, `EXPERIENCE.md:19-43`). *Fix:* 将修订后的选定方向提升为正式 mockup，在 Information Architecture／Components 附近内联链接并写明“聊天查询与待确认英雄状态”；保留一次“spines win on conflict”声明。两个未选方向可继续留在 `.working/`，无需成为正式引用。
- **[high] 选定 HTML 的状态色与正式 DESIGN 冲突：用户消息和发送按钮使用蓝色，在线使用绿色、离线使用红色、缺失值使用琥珀色；正式 spine 明确蓝色只用于查询／导航，在线／离线和缺失事实均为中性，发送按钮也必须为中性。若直接按该稿实现，会破坏安全状态语义。** (`.working/direction-precision-command.html:44,65,69,100-120`, `DESIGN.md:251-259,294,297,302-303,312-316`). *Fix:* 正式提升前先按 DESIGN 重绘这些状态；蓝色只保留查询与导航，琥珀只保留待确认，绿／红只保留回读成功／实际失败。

## 6. Bloat & overspecification — strong

两份 spine 的篇幅与任务复杂度匹配。DESIGN 的 token、组件视觉规则与 Do/Don't 分工清楚；EXPERIENCE 以表格承载 IA、组件和状态，叙述主要用于安全边界和流程，不存在大段 persona、FR 或市场背景复述。固定三栏尺寸虽同时出现在 token 与 prose 中，但 prose 给出了低宽度行为和产品理由，仍具有下游价值。

### Findings

无。

## 7. Inheritance discipline — adequate

两份 frontmatter 的 9 个来源均可解析；13 个正式组件名称在两份 spine 内一致；Product Facts、回读核验、待排查、未分配／未提供等关键词汇总体与来源一致。EXPERIENCE 没有悬空的 DESIGN token 引用。

### Findings

- **[high] 正式来源仍要求二次确认时“逐条展示每条漏洞”，而 UX 根据后续用户决定改为审批卡默认前 5 条、通过“查看全部”进入 IA-02。该决定记录在 memlog 和 spine 中，却没有同步到 PI-002、PRD／追踪矩阵；架构或验收消费者会面对两个权威合同。** (`PI-002-daily-new-vulnerability-manual-dispatch-draft.md:73-77`, `.memlog.md:28`, `DESIGN.md:298`, `EXPERIENCE.md:77,101,147`). *Fix:* 在 PI-002 及共同确认要求中明确“卡片默认前 5 条 + 确认前全部对象可进入只读明细核对”是否满足逐条可审阅，并同步追踪矩阵；不要只让 UX 文件隐式覆盖产品 AC。
- **[medium] 来源的 UJ 名称没有在 Key Flow 标题中原样继承；标题只使用 PI 编号，导致 `UJ-01`～`UJ-04` 需要人工通过追踪矩阵推导，`UJ-05` 则没有直接映射。** (`prd.md:34-40`, `requirements-traceability-matrix.md:11-18`, `EXPERIENCE.md:136-185`). *Fix:* 在每个 Key Flow 标题中同时写出原 UJ ID／名称和 PI ID，保持来源到体验 spine 的一跳追踪。

## 8. Shape fit — adequate

DESIGN 的章节顺序符合规范：Brand & Style → Colors → Typography → Layout & Spacing → Elevation & Depth → Shapes → Components → Do's and Don'ts。EXPERIENCE 包含 Foundation、Information Architecture、Voice and Tone、Component Patterns、State Patterns、Interaction Primitives、Accessibility Floor、Key Flows；Open Questions 对当前阻塞项有实际用途。单一桌面表面不需要 Responsive & Platform，且 Foundation 已明确仅桌面、仅浅色。

### Findings

- **[medium] memlog 明确记录了三个视觉方向、一个选定方向和两个拒绝方向，但 EXPERIENCE 缺少触发后应出现的 Inspiration & Anti-patterns，导致“为什么采用精密指挥台、为什么拒绝纸本审阅台／信号矩阵”只存在于工作记忆。** (`.memlog.md:29-32`, `EXPERIENCE.md:19-200`). *Fix:* 增加精简的 Inspiration & Anti-patterns，记录保留的高密度／硬边界／状态语义，以及拒绝的暖色编辑感与暗色监控感；不要重复 DESIGN 的完整品牌叙述。

## Mechanical notes

- Frontmatter：`DESIGN.md` 与 `EXPERIENCE.md` 均包含 `name`、`status`、9 个可解析 `sources` 和 `updated`；当前 `status: review` 与 Reviewer Gate 阶段一致。
- Token：所有颜色为 `#RRGGBB`；所有引用路径均指向已定义 token 或对象；未发现暗色 token 或主题切换。
- Components：两份 spine 的 13 个正式组件集合和中文名称完全一致。
- Flows：`PI-002`～`PI-005` 均有具名主人公、编号步骤、Climax 与失败路径；`UJ-05` 未原名出现。
- States：共 23 个显式状态；IA-02 加载／大结果浏览和人员读取的异步状态仍不完整。
- Visual files：正式 `imports/` 0 个文件，`mockups/` 0 个，`wireframes/` 0 个；`.working/` 有 3 个 HTML 方向，其中 precision-command 已被用户选择但尚未正式引用。
- Cross-reference：未发现 Mermaid 图，因此无 Mermaid 语法问题；未发现断裂 Markdown 文件链接，但也没有正式视觉链接。
- Severity totals：critical 0，high 7，medium 5，low 0。
