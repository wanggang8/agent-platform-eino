# 前端架构规范

本文固化 Workbench 前端技术栈、状态边界、目录结构和测试策略。Phase 2 之后不得在未新增 ADR 的情况下替换构建框架、状态分层、组件基础或样式体系。

## 技术栈

- 构建框架：Vite + React + TypeScript。
- 路由：React Router，使用 SPA/library mode，不启用 SSR 或 full-stack mode。
- 服务端状态：TanStack Query。
- UI/client 状态：Zustand。
- 复杂局部状态：组件内 `useReducer`。
- UI primitives：Radix UI Primitives。
- 样式：CSS Modules 或分层普通 CSS，不使用 Tailwind 作为默认样式体系。
- 图标：lucide-react。
- 测试：Vitest + Testing Library + Playwright。
- 类型边界：只从 `web/eino-workbench/src/contracts/generated.ts` 消费契约类型。

默认不使用 Next.js、Redux Toolkit、XState、Tailwind 或 shadcn 作为基础架构。引入这些能力必须新增 ADR，说明替代方案、收益、迁移成本和回滚条件。

## 架构边界

Product Facts 只存在于后端，是 Workbench 和 Action API 的唯一事实来源。前端只消费后端投影和 contract 类型，不重新建立事实模型。

TanStack Query 只管理 server state：

- Workbench view。
- Action result。
- Run snapshot。
- Replay view。
- Connector status。
- Provider credential binding。

Zustand 只管理 client/UI state：

- 当前 workspace/session 的 UI 选择。
- 选中的 Inspector tab。
- 展开的 tool card。
- sidebar 折叠状态。
- mobile panel 状态。
- composer draft 和本地输入状态。

组件内 reducer 只处理单个组件内的复杂交互，例如候选项选择、展开区域切换、临时表单校验。不得把 Product Facts、run lifecycle、approval/resume 状态机复制到前端 reducer 或 Zustand。

## 组件策略

Radix 只作为无样式交互基础，用于 Tabs、Dialog、Popover、Dropdown、Tooltip、Collapsible、ScrollArea 等复杂可访问交互。视觉样式由本项目 CSS 定义，不继承第三方默认视觉语言。

Workbench 必须保持 `docs/02-ux-visual-requirements.md` 定义的信息架构：

- Desktop：左导航、中间聊天、右 Inspector。
- Mobile：主聊天优先，Inspector 通过 tabs 或折叠入口进入。
- 时间线包含用户消息、assistant 文本、工具卡、审批卡、澄清卡和运行提示。
- Inspector 包含 evidence、structured、runtime、audit。

不得复制旧项目 DOM、CSS class、presenter 分支或接口字段。旧项目截图只能作为视觉验收基线。

## 建议目录

```text
web/eino-workbench/src/
  app/
    App.tsx
    router.tsx
    queryClient.ts
  contracts/
    generated.ts
  fixtures/
  features/
    workbench/
      api/
      components/
      state/
      views/
  components/
    primitives/
    layout/
  styles/
    tokens.css
    base.css
  test/
    render.tsx
```

`features/workbench` 承载业务工作台组合。`components/primitives` 只封装 Radix 和基础可访问控件，不依赖 Workbench 数据。`components/layout` 只处理布局骨架，不读取 provider 或 run 细节。

## 数据流

```text
React Router
  -> route params
  -> TanStack Query typed API
  -> generated contract types
  -> Workbench view projection
  -> components
  -> Zustand UI state for local selection only
```

Action mutation、approval resume、clarification resume 必须通过 typed API 调用后端。前端不得直接判断写域是否已经执行，不得跳过后端 policy、approval 或 audit。

SSE/stream 在 Phase 3 接入时，必须先归并为 contract 定义的 stream event，再更新 query cache 或局部 reducer。不得把 raw Eino event 或 raw provider payload 暴露给组件。

## 测试策略

- Vitest 覆盖 contract 消费、query key、reducer、UI state store 和安全渲染 helper。
- Testing Library 覆盖工具卡、审批卡、澄清卡、Inspector tabs、移动端 panel 的可访问交互。
- Playwright 覆盖 desktop 1440x900 和 mobile 390x844 的视觉 smoke。
- 动态字段必须 mask：时间、run id、随机 id、模型耗时、token 计数。
- 视觉测试通过结果必须来自新项目截图，旧项目截图只作为目标参考。

## 外部资料

- [React Build a React app from Scratch](https://react.dev/learn/build-a-react-app-from-scratch)
- [Vite Guide](https://vite.dev/guide/)
- [React Router](https://reactrouter.com/home)
- [TanStack Query](https://tanstack.com/query/latest)
- [Zustand](https://zustand.docs.pmnd.rs/)
- [Radix UI Primitives](https://www.radix-ui.com/primitives)
- [Playwright Visual Comparisons](https://playwright.dev/docs/test-snapshots)
