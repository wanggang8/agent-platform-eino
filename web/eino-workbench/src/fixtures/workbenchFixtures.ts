import type { ActionResult, WorkbenchView } from "../contracts/generated";
import actionApproval from "../../../../docs/fixtures/action-result-waiting-approval.json";
import actionClarification from "../../../../docs/fixtures/action-result-waiting-clarification.json";
import emptyView from "../../../../docs/fixtures/workbench-view-empty.json";
import failedView from "../../../../docs/fixtures/workbench-view-failed.json";
import successView from "../../../../docs/fixtures/workbench-view-success.json";
import { fobrainVisualFixtures, type FobrainVisualFixtureKey } from "./fobrainVisualFixtures";

// WorkbenchFixtureKey 枚举本地演示/视觉验收 fixture，不参与后端运行时事实选择。
export type WorkbenchFixtureKey = "success" | "empty" | "failed" | "approval" | "clarification" | FobrainVisualFixtureKey;

// approval/clarification fixture 复用 ActionResult 契约，保证等待态不手写平行 DTO。
const approvalResult = actionApproval as ActionResult;
const clarificationResult = actionClarification as ActionResult;

// approvalView 用于视觉验收写域审批卡，不表示写域操作已经执行。
const approvalView: WorkbenchView = {
  ...(successView as WorkbenchView),
  run_id: approvalResult.run_id,
  status: "waiting",
  timeline: [
    ...(successView as WorkbenchView).timeline,
    {
      item_id: "pending-approval-demo",
      kind: "approval_card",
      pending_id: approvalResult.approval_refs?.[0] ?? "approval:demo",
      status: "waiting",
      content: approvalResult.waiting?.question ?? "是否批准该高风险操作？"
    }
  ],
  inspector: {
    ...(successView as WorkbenchView).inspector,
    runtime: {
      status: "waiting",
      pending_id: approvalResult.approval_refs?.[0] ?? "approval:demo"
    }
  }
};

// clarificationView 用于视觉验收澄清卡，候选提交仍由后端 resume 决定。
const clarificationView: WorkbenchView = {
  ...(successView as WorkbenchView),
  run_id: clarificationResult.run_id,
  status: "waiting",
  timeline: [
    ...(successView as WorkbenchView).timeline,
    {
      item_id: "pending-clarification-demo",
      kind: "clarification_card",
      pending_id: clarificationResult.resume_refs?.[0] ?? "resume:demo",
      status: "waiting",
      content: clarificationResult.waiting?.question ?? "请选择需要继续处理的对象。",
      input_mode: clarificationResult.waiting?.input_mode ?? "single_choice",
      candidates: clarificationResult.waiting?.candidates ?? []
    }
  ],
  inspector: {
    ...(successView as WorkbenchView).inspector,
    runtime: {
      status: "waiting",
      pending_id: clarificationResult.resume_refs?.[0] ?? "resume:demo"
    }
  }
};

// workbenchFixtures 是 Phase 2 前端的唯一 fixture 数据入口。
export const workbenchFixtures = {
  success: successView as WorkbenchView,
  empty: emptyView as WorkbenchView,
  failed: failedView as WorkbenchView,
  approval: approvalView,
  clarification: clarificationView,
  ...fobrainVisualFixtures
} as const satisfies Record<WorkbenchFixtureKey, WorkbenchView>;

// fixtureLabels 是 fixture 切换器展示文案，不参与业务状态判断。
export const fixtureLabels: Record<WorkbenchFixtureKey, string> = {
  success: "工具完成",
  empty: "空态",
  failed: "失败态",
  approval: "审批等待",
  clarification: "澄清等待",
  fobrainConnectorSecurity: "Fobrain 连接器",
  fobrainCurrentUser: "Fobrain 当前用户",
  fobrainMyPermissions: "Fobrain 我的权限"
};
