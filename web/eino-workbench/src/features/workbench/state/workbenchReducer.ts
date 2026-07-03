import type { PendingInteraction, RunStatus, TimelineItem, WorkbenchStreamEvent, WorkbenchView } from "../../../contracts/generated";

const latestSequenceByView = new WeakMap<WorkbenchView, number>();

// applyWorkbenchStreamEvent 将后端 SSE 产品事件折叠进当前 WorkbenchView，不在前端重建 Product Facts。
export function applyWorkbenchStreamEvent(view: WorkbenchView, event: WorkbenchStreamEvent): WorkbenchView {
  if (event.run_id !== view.run_id) return view;
  if (isStaleEvent(view, event)) return view;
  if (event.type === "view.replaced" && event.view) return withLatestSequence(event.view, event.sequence);
  if (event.type === "run.updated" && event.run) {
    return withLatestSequence({
      ...view,
      status: event.run.status,
      inspector: {
        ...view.inspector,
        runtime: {
          ...view.inspector.runtime,
          status: event.run.status,
          ...(event.run.safe_error ? { safe_error: event.run.safe_error } : {})
        }
      }
    }, event.sequence);
  }
  if (event.type === "pending.updated" && event.pending) {
    return withLatestSequence(applyPendingPatch(view, event.pending), event.sequence);
  }
  return view;
}

// applyPendingPatch 更新或插入 pending 卡片；终态 pending 会从 runtime 摘要中移除。
function applyPendingPatch(view: WorkbenchView, pending: PendingInteraction): WorkbenchView {
  const item = timelineItemFromPending(pending);
  const existingIndex = view.timeline.findIndex((candidate) => candidate.pending_id === pending.pending_id);
  const timeline = existingIndex >= 0
    ? view.timeline.map((candidate, index) => (index === existingIndex ? { ...candidate, ...item } : candidate))
    : [...view.timeline, item];
  const nextStatus = pending.status === "waiting"
    ? "waiting"
    : view.status === "waiting" ? statusAfterTerminalPending(pending.status) : view.status;
  const runtime = pending.status === "waiting"
    ? { ...view.inspector.runtime, status: "waiting" as const, pending_id: pending.pending_id }
    : { ...runtimeWithoutPending(view), status: nextStatus };

  return {
    ...view,
    status: nextStatus,
    timeline,
    inspector: {
      ...view.inspector,
      runtime
    }
  };
}

// timelineItemFromPending 只投影产品展示字段，不把 resume_ref 或 checkpoint 细节放入时间线。
function timelineItemFromPending(pending: PendingInteraction): TimelineItem {
  const isApproval = pending.kind === "approval";
  const item: TimelineItem = {
    item_id: pending.pending_id,
    kind: isApproval ? "approval_card" : "clarification_card",
    pending_id: pending.pending_id,
    status: pending.status
  };
  const content = isApproval ? pending.risk_summary ?? pending.question : pending.question;
  return {
    ...item,
    ...(content ? { content } : {}),
    ...(pending.input_mode ? { input_mode: pending.input_mode } : {}),
    ...(pending.candidates ? { candidates: pending.candidates } : {})
  };
}

// runtimeWithoutPending 清理 runtime 中的等待态指针，保留其他运行摘要。
function runtimeWithoutPending(view: WorkbenchView) {
  const runtime = view.inspector.runtime;
  if (!runtime) return {};
  const { pending_id: _pendingID, ...rest } = runtime;
  return rest;
}

// isStaleEvent 使用模块内游标阻止旧 stream patch 回退终态；不把 sequence 写入 WorkbenchView contract。
function isStaleEvent(view: WorkbenchView, event: WorkbenchStreamEvent) {
  const latestSequence = latestSequenceByView.get(view) ?? 0;
  if (event.sequence <= latestSequence) return true;
  if (event.type === "pending.updated" && event.pending?.status === "waiting") {
    const existing = view.timeline.find((item) => item.pending_id === event.pending?.pending_id);
    return isTerminalPendingStatus(existing?.status);
  }
  return false;
}

function withLatestSequence(view: WorkbenchView, sequence: number) {
  latestSequenceByView.set(view, sequence);
  return view;
}

function isTerminalPendingStatus(status: string | undefined) {
  return status === "submitted" || status === "approved" || status === "rejected" || status === "cancelled" || status === "expired" || status === "consumed";
}

function statusAfterTerminalPending(status: PendingInteraction["status"]): RunStatus {
  if (status === "cancelled") return "cancelled";
  if (status === "rejected" || status === "expired") return "failed";
  return "running";
}
