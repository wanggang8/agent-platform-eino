import { describe, expect, it } from "vitest";
import type { WorkbenchStreamEvent, WorkbenchView } from "../../../contracts/generated";
import { applyWorkbenchStreamEvent } from "./workbenchReducer";

const baseView: WorkbenchView = {
  schema_version: "eino_workbench_view.v1",
  workspace_id: "ws-demo",
  run_id: "run-pending",
  status: "running",
  timeline: [],
  inspector: {
    tabs: ["evidence", "structured", "runtime", "audit"],
    runtime: { status: "running" }
  }
};

describe("workbenchReducer pending stream", () => {
  it("applies pending.updated by inserting an approval card from generated contract types", () => {
    // reducer 只能把 SSE patch 投影为 WorkbenchView，不在前端重建 Product Facts。
    const event: WorkbenchStreamEvent = {
      schema_version: "eino_workbench_stream_event.v1",
      event_id: "run-pending:000002",
      run_id: "run-pending",
      type: "pending.updated",
      sequence: 2,
      created_at: "2026-07-03T00:00:00Z",
      pending: {
        schema_version: "eino_workbench_pending_interaction.v1",
        pending_id: "pending-approval-1",
        run_id: "run-pending",
        kind: "approval",
        status: "waiting",
        question: "是否批准更新工单状态？",
        operation_name: "更新工单状态",
        risk_summary: "将修改外部工单状态",
        target_summary: "工单 T-1001",
        resume_ref: "resume_ref_approval_1"
      }
    };

    const next = applyWorkbenchStreamEvent(baseView, event);

    expect(next.status).toBe("waiting");
    expect(next.inspector.runtime?.pending_id).toBe("pending-approval-1");
    expect(next.timeline).toHaveLength(1);
    expect(next.timeline[0]).toMatchObject({
      kind: "approval_card",
      pending_id: "pending-approval-1",
      status: "waiting",
      content: "将修改外部工单状态"
    });
  });

  it("updates an existing pending card to terminal status without duplicating timeline items", () => {
    // 终态 patch 必须更新原卡片，避免 replay/stream 中出现两张同一 pending 卡。
    const waitingView = applyWorkbenchStreamEvent(baseView, {
      schema_version: "eino_workbench_stream_event.v1",
      event_id: "run-pending:000002",
      run_id: "run-pending",
      type: "pending.updated",
      sequence: 2,
      created_at: "2026-07-03T00:00:00Z",
      pending: {
        schema_version: "eino_workbench_pending_interaction.v1",
        pending_id: "pending-clarification-1",
        run_id: "run-pending",
        kind: "clarification",
        status: "waiting",
        question: "请选择人员",
        resume_ref: "resume_ref_clarification_1",
        input_mode: "single_choice",
        candidates: [{
          candidate_ref: "candidate:fobrain:person:1",
          label: "张三",
          entity_type: "person"
        }]
      }
    });

    const terminal = applyWorkbenchStreamEvent(waitingView, {
      schema_version: "eino_workbench_stream_event.v1",
      event_id: "run-pending:000003",
      run_id: "run-pending",
      type: "pending.updated",
      sequence: 3,
      created_at: "2026-07-03T00:00:01Z",
      pending: {
        schema_version: "eino_workbench_pending_interaction.v1",
        pending_id: "pending-clarification-1",
        run_id: "run-pending",
        kind: "clarification",
        status: "cancelled",
        question: "请选择人员",
        resume_ref: "resume_ref_clarification_1",
        input_mode: "single_choice"
      }
    });

    expect(terminal.timeline).toHaveLength(1);
    expect(terminal.timeline[0]).toMatchObject({
      kind: "clarification_card",
      pending_id: "pending-clarification-1",
      status: "cancelled"
    });
    expect(terminal.status).toBe("cancelled");
    expect(terminal.inspector.runtime?.pending_id).toBeUndefined();
    expect(terminal.inspector.runtime?.status).toBe("cancelled");
  });

  it("ignores stream events for another run and older pending sequences", () => {
    // reducer 必须隔离 run_id，并拒绝旧序号覆盖新状态，避免切换会话后串流污染当前视图。
    const waitingView = applyWorkbenchStreamEvent(baseView, {
      schema_version: "eino_workbench_stream_event.v1",
      event_id: "run-pending:000010",
      run_id: "run-pending",
      type: "pending.updated",
      sequence: 10,
      created_at: "2026-07-03T00:00:00Z",
      pending: {
        schema_version: "eino_workbench_pending_interaction.v1",
        pending_id: "pending-approval-1",
        run_id: "run-pending",
        kind: "approval",
        status: "waiting",
        question: "是否批准？",
        operation_name: "更新工单",
        risk_summary: "写域操作需要审批",
        target_summary: "工单 T-1001",
        resume_ref: "resume_ref_approval_1"
      }
    });
    const terminal = applyWorkbenchStreamEvent(waitingView, {
      schema_version: "eino_workbench_stream_event.v1",
      event_id: "run-pending:000011",
      run_id: "run-pending",
      type: "pending.updated",
      sequence: 11,
      created_at: "2026-07-03T00:00:01Z",
      pending: {
        schema_version: "eino_workbench_pending_interaction.v1",
        pending_id: "pending-approval-1",
        run_id: "run-pending",
        kind: "approval",
        status: "cancelled",
        question: "是否批准？",
        operation_name: "更新工单",
        risk_summary: "写域操作需要审批",
        target_summary: "工单 T-1001"
      }
    });
    const crossRun = applyWorkbenchStreamEvent(terminal, {
      schema_version: "eino_workbench_stream_event.v1",
      event_id: "run-other:000012",
      run_id: "run-other",
      type: "pending.updated",
      sequence: 12,
      created_at: "2026-07-03T00:00:02Z",
      pending: {
        schema_version: "eino_workbench_pending_interaction.v1",
        pending_id: "pending-other",
        run_id: "run-other",
        kind: "approval",
        status: "waiting",
        question: "是否批准其他会话？",
        operation_name: "其他操作",
        risk_summary: "其他风险",
        target_summary: "其他目标",
        resume_ref: "resume_ref_other"
      }
    });
    const olderWaiting = applyWorkbenchStreamEvent(crossRun, {
      schema_version: "eino_workbench_stream_event.v1",
      event_id: "run-pending:000009",
      run_id: "run-pending",
      type: "pending.updated",
      sequence: 9,
      created_at: "2026-07-03T00:00:03Z",
      pending: {
        schema_version: "eino_workbench_pending_interaction.v1",
        pending_id: "pending-approval-1",
        run_id: "run-pending",
        kind: "approval",
        status: "waiting",
        question: "是否批准？",
        operation_name: "更新工单",
        risk_summary: "写域操作需要审批",
        target_summary: "工单 T-1001",
        resume_ref: "resume_ref_approval_1"
      }
    });

    expect(crossRun).toEqual(terminal);
    expect(olderWaiting).toEqual(terminal);
  });

  it("does not let terminal pending override a completed run patch", () => {
    // 后端 snapshot stream 先发 run.updated，再发 pending.updated；pending 终态不能回退已完成的 run 状态。
    const completed = applyWorkbenchStreamEvent(baseView, {
      schema_version: "eino_workbench_stream_event.v1",
      event_id: "run-pending:000001",
      run_id: "run-pending",
      type: "run.updated",
      sequence: 1,
      created_at: "2026-07-03T00:00:00Z",
      run: {
        status: "succeeded",
        updated_at: "2026-07-03T00:00:00Z"
      }
    });

    const afterPending = applyWorkbenchStreamEvent(completed, {
      schema_version: "eino_workbench_stream_event.v1",
      event_id: "run-pending:000002",
      run_id: "run-pending",
      type: "pending.updated",
      sequence: 2,
      created_at: "2026-07-03T00:00:01Z",
      pending: {
        schema_version: "eino_workbench_pending_interaction.v1",
        pending_id: "pending-approval-1",
        run_id: "run-pending",
        kind: "approval",
        status: "approved",
        question: "是否批准？",
        operation_name: "更新工单",
        risk_summary: "写域操作需要审批",
        target_summary: "工单 T-1001"
      }
    });

    expect(afterPending.status).toBe("succeeded");
    expect(afterPending.inspector.runtime?.status).toBe("succeeded");
    expect(afterPending.timeline[0]).toMatchObject({
      pending_id: "pending-approval-1",
      status: "approved"
    });
  });
});
