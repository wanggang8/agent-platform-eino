import * as Tabs from "@radix-ui/react-tabs";
import type { M1WorkbenchView as WorkbenchView } from "../../../contracts/generated";
import { useWorkbenchUiStore } from "../state/useWorkbenchUiStore";

// Inspector 在 M1 只是右侧安全操作区空态，不展示技术证据或内部运行信息。
export function Inspector({ view }: { readonly view: WorkbenchView }) {
  const activeRightTab = useWorkbenchUiStore((state) => state.activeRightTab);
  const setActiveRightTab = useWorkbenchUiStore((state) => state.setActiveRightTab);
  return (
    <aside className="inspector" data-testid="inspector" aria-label="查询操作区">
      <Tabs.Root value={activeRightTab} onValueChange={(value) => setActiveRightTab(value as "事实" | "执行记录")}>
        <Tabs.List className="inspector-tabs" aria-label="查询信息">
          {view.right_panel.tabs.map((tab) => <Tabs.Trigger key={tab} value={tab}>{tab}</Tabs.Trigger>)}
        </Tabs.List>
        <Tabs.Content value="事实" className="inspector-empty"><strong>暂无可展示事实</strong><p>{view.right_panel.empty_message}</p></Tabs.Content>
        <Tabs.Content value="执行记录" className="inspector-empty"><strong>暂无执行记录</strong><p>完成一次查询后显示安全执行摘要。</p></Tabs.Content>
      </Tabs.Root>
    </aside>
  );
}
