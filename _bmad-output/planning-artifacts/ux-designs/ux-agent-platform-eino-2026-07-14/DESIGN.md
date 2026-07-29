---
name: FOBrain 精密指挥台
description: 面向安全运营人员的浅色桌面端 Agent Workbench 视觉契约。
status: final
sources:
  - ../../product-blueprint/product-brief.md
  - ../../product-blueprint/prd.md
  - ../../product-blueprint/implementation-readiness-gate.md
  - ../../product-blueprint/requirements-traceability-matrix.md
  - ../../product-blueprint/product-issues/PI-002-daily-new-vulnerability-manual-dispatch-draft.md
  - ../../product-blueprint/product-issues/PI-003-false-positive-marking-draft.md
  - ../../product-blueprint/product-issues/PI-004-responsible-person-transfer-draft.md
  - ../../product-blueprint/product-issues/PI-005-direct-repair-delay-draft.md
  - ../../product-blueprint/research/fobrain-target-live-read-evidence-2026-07-14.md
updated: '2026-07-15'
colors:
  surface-canvas: '#F3F5F7'
  surface-subtle: '#F8FAFC'
  surface-raised: '#FFFFFF'
  surface-sidebar: '#18212D'
  surface-sidebar-active: '#263545'
  text-primary: '#111827'
  text-secondary: '#475569'
  text-muted: '#64748B'
  text-inverse: '#FFFFFF'
  text-on-sidebar: '#B8C5D3'
  text-muted-on-sidebar: '#8B9BAD'
  border-default: '#D8DEE6'
  border-strong: '#C9D2DD'
  border-subtle: '#E3E8EE'
  border-on-sidebar: '#2D3947'
  focus-ring: '#111827'
  focus-ring-inverse: '#FFFFFF'
  query: '#0F6CBD'
  query-soft: '#E8F2FC'
  query-border: '#BFDCF3'
  success: '#146C43'
  success-soft: '#E9F7EF'
  pending-write: '#A15C00'
  pending-write-soft: '#FFF6E5'
  pending-write-border: '#E3B35D'
  failure: '#B42318'
  failure-soft: '#FDECEC'
typography:
  heading:
    fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif'
    fontSize: '16px'
    fontWeight: '700'
    lineHeight: '1.4'
    letterSpacing: '0'
  card-title:
    fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif'
    fontSize: '14px'
    fontWeight: '700'
    lineHeight: '1.45'
    letterSpacing: '0'
  body:
    fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif'
    fontSize: '13px'
    fontWeight: '400'
    lineHeight: '1.55'
    letterSpacing: '0'
  body-compact:
    fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif'
    fontSize: '12px'
    fontWeight: '400'
    lineHeight: '1.45'
    letterSpacing: '0'
  label:
    fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif'
    fontSize: '11px'
    fontWeight: '600'
    lineHeight: '1.4'
    letterSpacing: '0.02em'
  meta:
    fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif'
    fontSize: '10px'
    fontWeight: '400'
    lineHeight: '1.4'
    letterSpacing: '0.04em'
  metric:
    fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif'
    fontSize: '24px'
    fontWeight: '700'
    lineHeight: '1'
    letterSpacing: '-0.02em'
  data:
    fontFamily: 'ui-monospace, "SFMono-Regular", Consolas, "Liberation Mono", monospace'
    fontSize: '11px'
    fontWeight: '500'
    lineHeight: '1.4'
    letterSpacing: '0'
rounded:
  xs: '3px'
  sm: '4px'
  md: '5px'
  lg: '6px'
  xl: '8px'
  shell: '12px'
  full: '9999px'
spacing:
  '1': '2px'
  '2': '4px'
  '3': '6px'
  '4': '8px'
  '5': '10px'
  '6': '12px'
  '7': '14px'
  '8': '16px'
  '9': '18px'
  '10': '22px'
  '11': '24px'
  '12': '28px'
  '13': '32px'
components:
  工作台骨架:
    background: '{colors.surface-canvas}'
    minWidth: '1180px'
    columns: '232px minmax(600px, 1fr) 360px'
    dividerColor: '{colors.border-default}'
    dividerWidth: '1px'
    dividerStyle: 'solid'
    radius: '{rounded.shell}'
  会话导航:
    background: '{colors.surface-sidebar}'
    foreground: '{colors.text-on-sidebar}'
    activeForeground: '{colors.text-inverse}'
    activeBackground: '{colors.surface-sidebar-active}'
    activeIndicatorColor: '{colors.query}'
    activeIndicatorWidth: '3px'
    dividerColor: '{colors.border-on-sidebar}'
    dividerWidth: '1px'
    itemRadius: '{rounded.md}'
  操作记录入口:
    background: '{colors.surface-sidebar}'
    foreground: '{colors.text-on-sidebar}'
    activeForeground: '{colors.text-inverse}'
    activeBackground: '{colors.surface-sidebar-active}'
    activeIndicatorColor: '{colors.query}'
    dividerColor: '{colors.border-on-sidebar}'
    dividerWidth: '1px'
    itemRadius: '{rounded.md}'
  操作记录列表:
    background: '{colors.surface-raised}'
    foreground: '{colors.text-primary}'
    secondaryForeground: '{colors.text-secondary}'
    borderColor: '{colors.border-subtle}'
    selectedIndicatorColor: '{colors.query}'
    selectedBackground: '{colors.query-soft}'
    rowMinHeight: '52px'
  聊天时间线:
    background: '{colors.surface-canvas}'
    contentMaxWidth: '760px'
    paddingInline: '{spacing.12}'
    paddingBlock: '{spacing.10}'
    itemGap: '{spacing.9}'
  消息气泡:
    background: '{colors.surface-raised}'
    foreground: '{colors.text-primary}'
    borderColor: '{colors.border-default}'
    borderWidth: '1px'
    userBackground: '{colors.surface-subtle}'
    radius: '{rounded.lg}'
    paddingBlock: '{spacing.6}'
    paddingInline: '{spacing.7}'
  查询结果卡:
    background: '{colors.surface-raised}'
    foreground: '{colors.text-primary}'
    borderColor: '{colors.query-border}'
    borderWidth: '1px'
    accent: '{colors.query}'
    accentBackground: '{colors.query-soft}'
    radius: '{rounded.xl}'
    padding: '{spacing.7}'
  漏洞明细工作区:
    background: '{colors.surface-raised}'
    queueMinWidth: '600px'
    detailWidth: '360px'
    headerHeight: '52px'
    dividerColor: '{colors.border-default}'
    dividerWidth: '1px'
    selectedIndicatorColor: '{colors.query}'
    selectedIndicatorWidth: '3px'
    pageSize: '100'
  漏洞事实行:
    background: '{colors.surface-raised}'
    foreground: '{colors.text-primary}'
    minHeight: '40px'
    dividerColor: '{colors.border-subtle}'
    dividerWidth: '1px'
    paddingBlock: '{spacing.4}'
    paddingInline: '{spacing.5}'
    dataFontFamily: '{typography.data.fontFamily}'
    summaryLines: '1'
  审批卡:
    background: '{colors.surface-raised}'
    foreground: '{colors.text-primary}'
    pendingBackground: '{colors.pending-write-soft}'
    pendingForeground: '{colors.pending-write}'
    borderColor: '{colors.pending-write-border}'
    borderWidth: '1px'
    radius: '{rounded.xl}'
    padding: '{spacing.7}'
  澄清卡:
    background: '{colors.surface-subtle}'
    foreground: '{colors.text-primary}'
    secondaryForeground: '{colors.text-secondary}'
    borderColor: '{colors.border-strong}'
    borderWidth: '1px'
    focusRingColor: '{colors.focus-ring}'
    focusRingWidth: '2px'
    radius: '{rounded.lg}'
    padding: '{spacing.7}'
  人员选择卡:
    background: '{colors.surface-raised}'
    foreground: '{colors.text-primary}'
    secondaryForeground: '{colors.text-secondary}'
    borderColor: '{colors.border-strong}'
    borderWidth: '1px'
    selectedBorderColor: '{colors.query}'
    focusRingColor: '{colors.focus-ring}'
    focusRingWidth: '2px'
    radius: '{rounded.lg}'
    padding: '{spacing.7}'
  执行结果卡:
    background: '{colors.surface-raised}'
    foreground: '{colors.text-primary}'
    borderColor: '{colors.border-strong}'
    borderWidth: '1px'
    success: '{colors.success}'
    successBackground: '{colors.success-soft}'
    failure: '{colors.failure}'
    failureBackground: '{colors.failure-soft}'
    verificationPending: '{colors.text-secondary}'
    verificationPendingBackground: '{colors.surface-subtle}'
    radius: '{rounded.xl}'
    padding: '{spacing.7}'
  事实与执行面板:
    background: '{colors.surface-raised}'
    foreground: '{colors.text-primary}'
    width: '360px'
    dividerColor: '{colors.border-default}'
    dividerWidth: '1px'
    activeTab: '{colors.query}'
    panelPadding: '{spacing.8}'
  消息输入框:
    background: '{colors.surface-raised}'
    foreground: '{colors.text-primary}'
    placeholder: '{colors.text-muted}'
    borderColor: '{colors.border-strong}'
    borderWidth: '1px'
    focusRingColor: '{colors.focus-ring}'
    focusRingWidth: '2px'
    radius: '{rounded.xl}'
    minHeight: '76px'
    paddingBlock: '{spacing.5}'
    paddingInline: '{spacing.6}'
  状态标记:
    neutralForeground: '{colors.text-secondary}'
    neutralBackground: '{colors.surface-subtle}'
    queryForeground: '{colors.query}'
    queryBackground: '{colors.query-soft}'
    pendingForeground: '{colors.pending-write}'
    pendingBackground: '{colors.pending-write-soft}'
    successForeground: '{colors.success}'
    successBackground: '{colors.success-soft}'
    failureForeground: '{colors.failure}'
    failureBackground: '{colors.failure-soft}'
    verificationForeground: '{colors.text-secondary}'
    verificationBackground: '{colors.surface-subtle}'
    radius: '{rounded.sm}'
    paddingBlock: '{spacing.1}'
    paddingInline: '{spacing.2}'
---

## Brand & Style

FOBrain 精密指挥台是一套冷静、专业、高密度的安全运营视觉系统。它首先是一件可靠的工作工具：信息边界硬朗、事实排列紧凑、状态含义确定，聊天仍是主工作面，但任何操作范围、确认与结果都以可核对的结构呈现。第一屏必须像指挥台而不是营销页；不使用英雄插画、装饰性渐变、玻璃拟态或无业务含义的彩色点缀。

本契约只覆盖浅色桌面端 Web Workbench，不提供暗色 token、主题切换、移动断点或触屏专用变体。交互基础使用无样式的 Radix UI Primitives，所有视觉由项目 CSS 实现；图标使用既定的 lucide-react。字体仅使用操作系统字体栈和系统等宽字体，不加载外部字体。不得引入 Tailwind、shadcn 或第三方组件的默认视觉层。

视觉层不重建事实语义。界面只呈现安全投影后的 Product Facts；raw JSON、provider 私有字段、内部 ID、token、密钥和英文内部状态码均不得成为视觉材料。

视觉方向以[精密指挥台](mockups/direction-precision-command.html)为基线；它只说明三栏构图、密度和状态分区，不是 DOM／交互模板。当前首版以[当前只读聊天](mockups/key-current-readonly-chat.html)为实现参考；门禁后的写态只能参考[目标审批与部分失败](mockups/key-target-write-approval-partial-failure.html)。任何 mock 与本 spine 冲突时，以本 spine 为准。

## Colors

颜色是状态契约，不是装饰。基础界面由 `{colors.surface-canvas}`、`{colors.surface-raised}`、`{colors.text-primary}` 与 `{colors.border-default}` 构成；深色仅用于 `{colors.surface-sidebar}`，不构成暗色主题。

- **蓝色**：`{colors.query}`、`{colors.query-soft}` 与 `{colors.query-border}` 只表示查询、查询结果、当前导航或可导航目标。蓝色不得表示“通用主要操作”，也不得给写操作确认、执行成功或失败染色。
- **琥珀色**：`{colors.pending-write}`、`{colors.pending-write-soft}` 与 `{colors.pending-write-border}` 只表示等待用户二次确认且尚未发生外部写入的操作。缺失字段、一般提醒、澄清和只读异常不得使用琥珀色。
- **红色**：`{colors.failure}` 与 `{colors.failure-soft}` 只表示已经发生的失败、被拒绝的执行结果或无法完成的失败项。不得用红色表示待处理、风险等级、离线或缺失数据。
- **绿色**：`{colors.success}` 与 `{colors.success-soft}` 只表示经过回读或等价证据核验的成功。请求已发送、工具已返回或正在执行都不具备使用绿色的资格。
- **中性色**：缺失事实用确定性文案配合 `{colors.text-secondary}` 或 `{colors.surface-subtle}`，例如“未分配”“未提供”；不以空白或语义色替代文字。在线/离线等事实状态默认使用中性文本，除非其本身是已核验成功或失败结论。

浅色表面上的键盘焦点使用 `{colors.focus-ring}` 的深色 2px 轮廓；深色会话导航和深色按钮使用 `{colors.focus-ring-inverse}` 的浅色 2px 轮廓。两者都保留 2px 间距，并与相邻表面达到至少 3:1 的非文本对比。文本与背景组合须达到 WCAG 2.2 AA；不得只用颜色传达状态，必须同时有确定性文案或图标与文本。

## Typography

所有界面文字使用系统字体栈。中文优先由 PingFang SC 或 Microsoft YaHei 渲染，拉丁字符回退到系统 UI 字体；IP、时间、计数和稳定标识使用 `{typography.data}`，以等宽数字支持纵向核对。禁止加载外部字体，也不以品牌字体制造展示性标题。

`{typography.heading}` 只用于页面级标题；卡片标题使用 `{typography.card-title}`；正文使用 `{typography.body}`；表格、按钮和高密度辅助文本使用 `{typography.body-compact}`；短标签与状态名称使用 `{typography.label}`；时间戳和来源说明使用 `{typography.meta}`。查询总量使用 `{typography.metric}`，数字采用 tabular numerals，不通过超大字号制造仪表盘式噪声。

标题、标签和状态文案以中文为主，不使用全大写英文内部状态。行高必须保证中文不拥挤；文本缩放至 200% 时不截断关键事实、确认对象和结果状态。

## Layout & Spacing

`{components.工作台骨架.columns}` 是桌面首版的固定三栏骨架：232px 会话区、最小 600px 的主工作面、360px 的事实与执行区。视口最小宽度为 `{components.工作台骨架.minWidth}`；低于该宽度时允许页面级横向滚动或明确提示桌面宽度不足，不折叠为移动布局，不把右栏变成抽屉。浏览器文字缩放至 200% 时，聊天、澄清、审批和结果卡必须在中栏内换行且不遮挡动作；只有真正需要二维阅读的漏洞明细表可以在自身容器内水平滚动，关键确认与错误摘要不得依赖整页横向往返阅读。

聊天模式中，中栏保持会话标题、聊天时间线和消息输入框的纵向结构；漏洞明细模式保留左侧会话导航，中栏切换为高密度队列，右栏显示当前选中项详情。两种模式共用同一骨架和分割线，不通过浮层叠加出第二套工作台。

间距以 2px 为基础微单位，常用密度为 4 / 8 / 12 / 16 / 22 / 28 / 32px。表格行、标签和元数据使用 `{spacing.2}` 至 `{spacing.5}`；卡片内边距使用 `{spacing.6}` 至 `{spacing.8}`；主栏水平留白使用 `{spacing.12}`。高密度不等于贴边：任何可点击目标保留清晰边界和至少 32px 的桌面指针命中高度，主要按钮建议 36px。

## Elevation & Depth

层级首先依靠色调、1px 分割线和留白建立。基础面板无阴影；消息气泡只允许 `0 1px 2px rgba(15, 23, 42, 0.03)`；查询结果卡、审批卡和执行结果卡可使用 `0 4px 14px rgba(15, 23, 42, 0.06)`，但同一视口不得层层叠加阴影。

首版不授权未在两份 spine 中命名的通用 Dialog、Popover 或 Tooltip。正式组件确需使用 Radix 覆盖层时，最大阴影为 `0 12px 32px rgba(15, 23, 42, 0.14)`，并配合 `{colors.border-strong}` 边框；其关闭、焦点回归和读屏关系由对应组件契约负责。不得用大面积投影、发光或彩色阴影表示状态；状态仍由确定性文字和规定的语义色承担。

## Shapes

形状保持工具感：`{rounded.xs}` 用于紧凑标签，`{rounded.sm}` 至 `{rounded.md}` 用于行内控件和导航项，`{rounded.lg}` 至 `{rounded.xl}` 用于消息与结构化卡片，`{rounded.shell}` 只用于独立演示或宿主容器。产品实际全屏工作台外框不需要圆角。

`{rounded.full}` 只用于头像、单一状态点等本质为圆形的元素。按钮、状态标记、标签页和卡片不得做成胶囊；同一组件的内外圆角应相差 2px 左右，保持精密、可装配的边界关系。

## Components

- **工作台骨架**：使用 `{components.工作台骨架.columns}` 三栏网格与 `{colors.surface-canvas}` 底色；左、中、右之间只有 1px 中性分割线。首屏占满可用桌面高度，栏宽不随内容跳动；浏览器示意外框不属于产品 UI。
- **会话导航**：固定 232px 深色侧栏，主文字用 `{colors.text-on-sidebar}`，弱化信息用 `{colors.text-muted-on-sidebar}`。系统固定“操作记录”入口与聊天会话使用独立分组和分割线，不能伪装成某个会话；当前目标使用 `{colors.surface-sidebar-active}` 和 3px 蓝色内侧导航指示条，蓝色仅因其导航含义而使用。导航项不加投影，层级由缩进、分组标签和分割线完成；键盘焦点必须使用 `{colors.focus-ring-inverse}`，不得沿用深色焦点环。
- **操作记录列表**：复用 IA-01 中栏，使用白色列表和 1px 中性分隔线。每行按动作、操作时间、目标安全显示和“完成／待核验／待排查”计数形成稳定层级；待核验使用中性灰底，不使用审批琥珀色，待排查只在状态文字／图标使用红色，不整行染红。时间、动作与状态筛选保持紧凑；“待核验／待排查”快捷筛选是导航控件，不是动作按钮。选中记录只使用蓝色导航指示，不表示批准或成功。
- **聊天时间线**：使用 `{colors.surface-canvas}` 平面底色和最大 760px 的事实阅读宽度，消息间距为 `{spacing.9}`。卡片沿统一左轴排列，用户消息可右对齐但不得挤压结构化卡片；百条级对象不能在时间线全量展开。
- **消息气泡**：Agent 与用户消息都使用中性表面，用户消息以 `{colors.surface-subtle}` 区分，不默认使用蓝色，因为消息也可能表达写意图。正文使用 `{typography.body}`；发言者、时间和来源使用 `{typography.meta}`。气泡不承载成功、失败或待确认状态色。
- **查询结果卡**：白色表面、蓝色 1px 边框和局部 `{colors.query}` 导航/查询强调；总量使用 `{typography.metric}`，摘要指标以分栏和中性分割线组织。“查看全部明细”可使用蓝色链接，因为它是从查询摘要进入明细的导航。卡片不得内嵌百条级完整列表。
- **漏洞明细工作区**：保留三栏骨架，中栏显示最小 600px 的队列，右栏固定 360px 详情。表头使用 `{colors.surface-subtle}`，选中行仅以蓝色左边线和淡蓝底表示当前导航位置；选择不等于批准。查询快照固定按每页 `{components.漏洞明细工作区.pageSize}` 条分页，页码、总页数、总数与已加载状态始终可见；分页控件保持紧凑，不用彩色热力图替代事实字段。布局参考[分页漏洞明细](mockups/key-paginated-vulnerability-details.html)。
- **漏洞事实行**：最小高度 40px，以 1px `{colors.border-subtle}` 分隔；POC、IP、时间和计数按列对齐，稳定数据使用 `{typography.data}`。表格单元格只显示 `{components.漏洞事实行.summaryLines}` 行摘要：超长 POC 使用视觉省略，多个 IP 显示“首个 IP + 其余 N 个”；单元格的可访问名称保留完整安全值，当前行的右栏按一项一行展示完整 POC 与全部 IP。缺失负责人显示“未分配”，缺失业务系统或负责人显示“未提供”，全部使用中性色。行不得因风险、在线状态或缺失事实整行染色。
- **审批卡**：只用于尚未写入且等待二次确认的动作，使用琥珀色标题带、边框和确认按钮。卡片必须在视觉上呈现动作、目标负责人、总数和前 5 条对象，并提供进入全部明细的中性或导航入口；确认与取消留在卡内。确认后卡片不继续显示为待确认，应转为中性执行中状态或由执行结果卡接管。
- **澄清卡**：使用中性灰白表面与 `{colors.border-strong}`，不使用琥珀或红色，因为澄清不是待确认写入，也不是失败。问题、可选项和“为何需要澄清”的说明分层排列；当前键盘选项使用深色焦点环。人员无法唯一确定、引用歧义等阻断须用明确文字说明。
- **人员选择卡**：只在 PI-002／PI-004 写域门禁通过后出现，使用中性白色表面、搜索输入、候选列表和已选人员复核区。候选项使用姓名与经批准的安全消歧字段，不显示 AI 推荐、推荐排序或内部人员 ID；当前键盘候选使用 `{colors.focus-ring}`，最终选中只使用蓝色导航边界，不使用绿色成功语义。加载、空、失败、同名和不可选状态都必须保留文字。
- **执行结果卡**：卡片外框保持中性；逐条结果中，只有回读核验完成项使用绿色，只有明确失败或核验期限／预算结束仍未知的待排查项使用红色。核验期限／预算内的待核验使用 `{components.执行结果卡.verificationPending}` 与中性底，不使用审批琥珀色或失败红色。混合结果不得整卡染绿或染红，标题必须同时给出完成、待排查与待核验数量。正在执行和等待回读不得提前显示绿色。
- **事实与执行面板**：固定 360px 白色右栏，以标签页承载事实和执行记录；技术信息只有在 OQ-08 的授权与字段 allowlist 明确后才可增加。当前标签可使用蓝色下划线，因为标签切换属于导航；事实键值保持左右对齐，选中漏洞的超长 POC 和全部 IP 在右栏完整换行展示，不能继续省略，也不能只靠 hover 查看。聊天模式显示通用证据与运行摘要，明细模式优先显示选中漏洞的业务事实。
- **消息输入框**：白色面、1px 强边框、`{rounded.xl}` 圆角与至少 76px 高度；发送按钮使用 `{colors.text-primary}` 中性深色，不使用蓝色或琥珀，以免把任意自然语言预先解释为查询或写入。占位文案说明可查询或处理，辅助文案强调写范围由结果引用锁定；键盘焦点使用深色 2px 轮廓。
- **状态标记**：采用小圆角矩形和“图标/符号 + 确定性文字”，不做纯色圆点。中性变体用于缺失、执行中、待核验与普通事实；蓝色变体仅用于查询或当前导航；琥珀变体仅用于等待二次确认；绿色变体仅用于已核验成功；红色变体仅用于明确失败或最终待排查。所有变体须在无颜色条件下仍可由文字区分。

Radix UI Primitives 只提供键盘交互、焦点管理、弹层关系和 ARIA 基础，不携带视觉默认值。项目 CSS 必须直接消费上述 token；不得通过 Tailwind utility、shadcn theme 或外部字体间接改写这些视觉规则。

## Do's and Don'ts

| Do | Don't |
|---|---|
| 用中性表面、1px 边界和紧凑对齐建立专业密度 | 用大卡片、留白海报或装饰插画稀释工作台信息 |
| 蓝色只用于查询、查询结果和导航位置 | 把发送、确认、成功或一般主要操作统一涂蓝 |
| 琥珀色只用于尚未写入的二次确认 | 用琥珀表示缺失、风险、提醒或普通阻断 |
| 绿色只用于回读或等价证据核验成功 | 工具返回、请求已发出或执行中就显示成功 |
| 红色只用于失败，并同时给出失败文案 | 用红色表示离线、待处理、风险等级或未提供 |
| 缺失事实明确显示“未分配”或“未提供” | 留空、用 `--` 掩盖，或让 AI 补全缺失事实 |
| 使用系统字体、Radix 无样式基础和项目 CSS | 加载外部字体，或引入 Tailwind、shadcn 默认视觉 |
| 保持 232px / minmax(600px, 1fr) / 360px 桌面三栏 | 折叠成移动布局、抽屉 Inspector 或触屏专用导航 |
| 让审批范围和混合执行结果逐条可核对 | 整卡染色、只报总数，或把部分失败包装成全量成功 |
| 用固定“操作记录”入口和只读列表重新找到历史事实 | 把操作记录伪装成聊天会话，或从历史卡直接重试写入 |
| 所有状态同时提供文字、必要图标和可见焦点 | 只靠颜色、阴影、动效或内部状态码表达含义 |
| 深色侧栏和深色按钮使用浅色反向焦点环 | 在深色表面继续使用不可见的深色焦点环 |
