import * as Collapsible from "@radix-ui/react-collapsible";
import { CheckCircle2, ChevronRight, CircleAlert, Hourglass, ShieldCheck } from "lucide-react";
import type { TimelineItem } from "../../../contracts/generated";
import { useWorkbenchUiStore } from "../state/useWorkbenchUiStore";
import { StructuredResultView } from "./StructuredResultView";

export function ToolCard({ item }: { readonly item: TimelineItem }) {
  const collapsed = useWorkbenchUiStore((state) => state.collapsedToolCards.has(item.item_id));
  const toggleToolCard = useWorkbenchUiStore((state) => state.toggleToolCard);

  return (
    <Collapsible.Root open={!collapsed} onOpenChange={() => toggleToolCard(item.item_id)}>
      <article className="tool-card" data-testid="tool-card">
        <Collapsible.Trigger className="tool-card-header">
          <span className={`status-dot ${item.status === "failed" ? "danger" : "success"}`}>
            {item.status === "failed" ? <CircleAlert size={16} /> : <CheckCircle2 size={16} />}
          </span>
          <div>
            <strong>工具调用</strong>
            <p>{item.safe_summary ?? "工具结果已生成安全结构化摘要。"}</p>
          </div>
          <span className={`status-pill ${item.status === "failed" ? "danger" : ""}`}>{statusLabel(item.status)}</span>
          <ChevronRight className={collapsed ? "" : "is-open"} size={18} />
        </Collapsible.Trigger>
        <Collapsible.Content>
          {item.structured_result ? <StructuredResultView result={item.structured_result} /> : null}
        </Collapsible.Content>
      </article>
    </Collapsible.Root>
  );
}

export function ApprovalCard({ item }: { readonly item: TimelineItem }) {
  return (
    <article className="pending-card approval" data-testid="approval-card">
      <div className="pending-icon"><ShieldCheck size={18} /></div>
      <div>
        <span>审批请求</span>
        <h3>{item.content}</h3>
        <p>未审批前，界面不会暗示写域操作已经执行。</p>
      </div>
      <div className="pending-actions">
        <button type="button">拒绝</button>
        <button type="button" className="primary">批准并提交</button>
      </div>
    </article>
  );
}

export function ClarificationCard({ item }: { readonly item: TimelineItem }) {
  return (
    <article className="pending-card clarification" data-testid="clarification-card">
      <div className="pending-icon"><Hourglass size={18} /></div>
      <div>
        <span>需要澄清</span>
        <h3>{item.content}</h3>
        <p>请确认候选对象后继续，提交后结果会变为只读记录。</p>
      </div>
      <div className="candidate-list" aria-label="候选对象">
        <button type="button">生产网段资产</button>
        <button type="button">办公网段资产</button>
      </div>
    </article>
  );
}

function statusLabel(status: string | undefined) {
  if (!status) return "未知";
  const labels: Record<string, string> = {
    pending: "等待中",
    queued: "已排队",
    running: "运行中",
    succeeded: "成功",
    completed: "完成",
    failed: "失败",
    cancelled: "已取消",
    waiting: "等待"
  };
  return labels[status] ?? status;
}
