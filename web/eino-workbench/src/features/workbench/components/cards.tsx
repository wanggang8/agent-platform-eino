import * as Collapsible from "@radix-ui/react-collapsible";
import { CheckCircle2, ChevronRight, CircleAlert, Hourglass, ShieldCheck } from "lucide-react";
import type { TimelineItem } from "../../../contracts/generated";
import { useWorkbenchUiStore } from "../state/useWorkbenchUiStore";
import { StructuredResultView } from "./StructuredResultView";

// ToolCard 展示工具调用的安全摘要和 StructuredResult，不展示工具内部 raw payload。
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

// ApprovalCard 展示写域审批等待态，前端不直接决定写域执行。
export function ApprovalCard({ item }: { readonly item: TimelineItem }) {
  const terminal = isTerminalPendingStatus(item.status);
  return (
    <article className={`pending-card approval${terminal ? " is-terminal" : ""}`} data-testid="approval-card">
      <div className="pending-icon"><ShieldCheck size={18} /></div>
      <div>
        <span>审批请求 · {statusLabel(item.status)}</span>
        <h3>{item.content}</h3>
        <p>{terminal ? "审批已进入只读记录，后续状态以运行事实为准。" : "未审批前，界面不会暗示写域操作已经执行。"}</p>
      </div>
      {terminal ? (
        <div className="pending-terminal" aria-label="审批终态">{statusLabel(item.status)}</div>
      ) : (
        <div className="pending-actions">
          <button type="button">取消</button>
          <button type="button">拒绝</button>
          <button type="button" className="primary">批准并提交</button>
        </div>
      )}
    </article>
  );
}

// ClarificationCard 展示澄清等待态，提交后的状态由后端 Product Facts 决定。
export function ClarificationCard({ item }: { readonly item: TimelineItem }) {
  const candidates = item.candidates ?? [];
  const terminal = isTerminalPendingStatus(item.status);
  return (
    <article className={`pending-card clarification${terminal ? " is-terminal" : ""}`} data-testid="clarification-card">
      <div className="pending-icon"><Hourglass size={18} /></div>
      <div>
        <span>需要澄清 · {statusLabel(item.status)}</span>
        <h3>{item.content}</h3>
        <p>{terminal ? "澄清已进入只读记录，后续状态以运行事实为准。" : "请确认候选对象后继续，提交后结果会变为只读记录。"}</p>
      </div>
      {terminal ? (
        <div className="pending-terminal" aria-label="澄清终态">{statusLabel(item.status)}</div>
      ) : (
        <div className="candidate-list" aria-label="候选对象">
          {candidates.length > 0 ? candidates.map((candidate) => (
            <button key={candidate.candidate_ref} type="button">
              <strong>{candidate.label}</strong>
              {candidate.description ? <span>{candidate.description}</span> : null}
            </button>
          )) : <button type="button">等待补充信息</button>}
        </div>
      )}
    </article>
  );
}

// isTerminalPendingStatus 判断 pending 是否已经不可再提交恢复动作。
function isTerminalPendingStatus(status: string | undefined) {
  return status === "submitted" || status === "approved" || status === "rejected" || status === "cancelled" || status === "expired" || status === "consumed";
}

// statusLabel 将工具/pending 状态映射为中文标签。
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
    waiting: "等待",
    submitted: "已提交",
    approved: "已批准",
    rejected: "已拒绝",
    expired: "已过期",
    consumed: "已恢复",
    pending_approval: "等待审批",
    provider_timeout: "服务超时"
  };
  return labels[status] ?? "未知";
}
