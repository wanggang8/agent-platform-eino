import { MessageSquare, Plus, Search, Settings, UserRound } from "lucide-react";

type WorkspaceSidebarProps = {
  readonly activeRunId: string;
};

const conversations = [
  { title: "查询 10.10.11.69 关联资产与风险", meta: "已完成 · 3 个资产 · 1 个高危漏洞", time: "16:48" },
  { title: "漏洞工单状态更新审批流程", meta: "等待审批 · 工单 #INC-20250624-0178", time: "15:22" },
  { title: "资产基线合规检查", meta: "已完成 · 不合规 7 条", time: "14:10" },
  { title: "域名枚举与暴露面分析", meta: "已完成 · 发现域名 12 个", time: "11:35" }
];

export function WorkspaceSidebar({ activeRunId }: WorkspaceSidebarProps) {
  return (
    <aside className="workspace-sidebar" data-testid="workspace-sidebar" aria-label="Workspace navigation">
      <div className="brand-row">
        <div className="brand-mark">A</div>
        <div>
          <strong>智能任务台</strong>
          <span>AI Agent Workbench</span>
        </div>
      </div>
      <button className="new-chat" type="button"><Plus size={16} /> 新建会话</button>
      <label className="search-box">
        <Search size={16} />
        <input aria-label="搜索会话" placeholder="搜索会话或收藏内容" />
      </label>
      <nav className="conversation-list" aria-label="会话列表">
        {conversations.map((item, index) => (
          <button key={item.title} className={index === 0 ? "conversation-item is-active" : "conversation-item"} type="button">
            <span className="conversation-title">{item.title}</span>
            <span className="conversation-meta">{item.meta}</span>
            <span className="conversation-time">{item.time}</span>
            {index === 0 ? <MessageSquare size={14} aria-hidden /> : null}
          </button>
        ))}
      </nav>
      <div className="sidebar-footer">
        <button type="button"><Settings size={16} /> 设置与能力</button>
        <button type="button"><UserRound size={16} /> Admin · {activeRunId}</button>
      </div>
    </aside>
  );
}
