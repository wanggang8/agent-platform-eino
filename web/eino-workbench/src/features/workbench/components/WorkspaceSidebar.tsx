import { History, MessageSquareText } from "lucide-react";
import { useWorkbenchUiStore } from "../state/useWorkbenchUiStore";

// WorkspaceSidebar 只提供当前会话和固定操作记录入口，不伪造历史列表。
export function WorkspaceSidebar() {
  const historyOpen = useWorkbenchUiStore((state) => state.historyOpen);
  const toggleHistory = useWorkbenchUiStore((state) => state.toggleHistory);
  return (
    <nav className="workspace-sidebar" data-testid="workspace-sidebar" aria-label="工作区导航">
      <div className="brand-block">
        <span aria-hidden>智</span>
        <div><strong>智能任务台</strong><small>安全运营</small></div>
      </div>
      <div className="sidebar-section">
        <span className="sidebar-label">当前会话</span>
        <div className="current-conversation"><MessageSquareText size={18} /><span>新增漏洞查询</span></div>
      </div>
      <div className="sidebar-section sidebar-records">
        <button type="button" aria-expanded={historyOpen} onClick={toggleHistory}>
          <History size={18} />操作记录
        </button>
        {historyOpen ? <p data-testid="history-empty">操作记录尚未启用</p> : null}
      </div>
    </nav>
  );
}
