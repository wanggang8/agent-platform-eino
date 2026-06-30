import { Bot, Info, UserRound } from "lucide-react";
import type { TimelineItem } from "../../../contracts/generated";
import { ApprovalCard, ClarificationCard, ToolCard } from "./cards";

type TimelineProps = {
  readonly items: readonly TimelineItem[];
};

export function Timeline({ items }: TimelineProps) {
  return (
    <section className="chat-timeline" data-testid="chat-timeline" aria-label="对话时间线">
      {items.length === 0 ? (
        <div className="empty-state">
          <Bot size={28} />
          <h2>开始一个新的安全分析会话</h2>
          <p>发送问题后，工具结果会以 StructuredResult 的安全投影展示在这里。</p>
        </div>
      ) : (
        items.map((item) => <TimelineEntry key={item.item_id} item={item} />)
      )}
    </section>
  );
}

function TimelineEntry({ item }: { readonly item: TimelineItem }) {
  switch (item.kind) {
    case "user_message":
      return <MessageBubble role="user" content={item.content ?? ""} />;
    case "assistant_message":
      return <MessageBubble role="assistant" content={item.content ?? ""} />;
    case "tool_card":
      return <ToolCard item={item} />;
    case "approval_card":
      return <ApprovalCard item={item} />;
    case "clarification_card":
      return <ClarificationCard item={item} />;
    case "run_notice":
      return (
        <div className="run-notice">
          <Info size={16} />
          <span>{item.content}</span>
        </div>
      );
  }
}

function MessageBubble({ role, content }: { readonly role: "user" | "assistant"; readonly content: string }) {
  return (
    <article className={`message-row ${role}`}>
      <div className="avatar" aria-hidden>{role === "user" ? <UserRound size={16} /> : "A"}</div>
      <div className="message-bubble">
        <span>{role === "user" ? "用户" : "智能助手"}</span>
        <p>{content}</p>
      </div>
    </article>
  );
}
