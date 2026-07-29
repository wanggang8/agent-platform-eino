import type { M1WorkbenchView as WorkbenchView } from "../../../contracts/generated";
import { Composer } from "./Composer";
import { Inspector } from "./Inspector";
import { WorkspaceSidebar } from "./WorkspaceSidebar";

type WorkbenchShellProps = {
  readonly view: WorkbenchView;
  readonly isSubmitting?: boolean;
  readonly onSubmit?: (content: string) => Promise<void> | void;
};

// WorkbenchShell 固定 nav/main/aside 的浅色桌面信息架构，不构造第二套事实状态。
export function WorkbenchShell({ view, isSubmitting = false, onSubmit = async () => undefined }: WorkbenchShellProps) {
  return (
    <div className="workbench-shell" data-testid="workbench-shell">
      <a className="skip-link" href="#conversation-main">跳到主要内容</a>
      <WorkspaceSidebar />
      <main id="conversation-main" className="conversation-main" tabIndex={-1}>
        <header className="conversation-header">
          <div>
            <p>安全运营工作台</p>
            <h1>新增漏洞查询</h1>
          </div>
          <span className={`query-state query-state-${view.status}`}>{statusLabel(view.status)}</span>
        </header>
        <section className="chat-timeline" data-testid="chat-timeline" aria-label="聊天记录" aria-live="polite">
          {view.messages.length === 0 ? <div className="chat-empty">输入查询后，这里会显示结果摘要。</div> : null}
          {view.messages.map((message) => (
            <article className={`chat-message chat-message-${message.role}`} key={message.message_id}>
              <span>{message.role === "user" ? "你" : "智能助手"}</span>
              <p>{message.content}</p>
            </article>
          ))}
          {view.result ? (
            <article className="query-result-card" data-testid="query-result-card" aria-label="查询结果摘要">
              <div>
                <span>查询 {view.result.query_sequence}</span>
                <time dateTime={view.result.observed_at}>{formatObservedAt(view.result.observed_at)}</time>
              </div>
              <h2>{view.result.summary}</h2>
              <p><strong>{view.result.count}</strong> 条</p>
            </article>
          ) : null}
          {view.safe_error ? <div className="query-error" role="alert">{view.safe_error}</div> : null}
        </section>
        <Composer disabled={isSubmitting} onSubmit={onSubmit} />
      </main>
      <Inspector view={view} />
    </div>
  );
}

function statusLabel(status: WorkbenchView["status"]) {
  const labels: Record<WorkbenchView["status"], string> = {
    idle: "等待查询",
    running: "查询中",
    resolved: "查询完成",
    empty: "查询完成",
    failed: "查询失败"
  };
  return labels[status];
}

function formatObservedAt(value: string) {
  return new Intl.DateTimeFormat("zh-CN", { dateStyle: "medium", timeStyle: "short", hour12: false }).format(new Date(value));
}
