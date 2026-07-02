import type { StructuredResult, WorkbenchView } from "../contracts/generated";
import currentUserResult from "../../../../docs/fixtures/fobrain/current-user-context-result.json";
import myPermissionsResult from "../../../../docs/fixtures/fobrain/my-permissions-result.json";

// FobrainVisualFixtureKey 固定当前 Batch A 业务只读视觉 fixture，connector 仍由 Task 8.3 单独覆盖。
export type FobrainVisualFixtureKey = "fobrainCurrentUser" | "fobrainMyPermissions";

// FobrainVisualRegion 对齐旧最终验收六区域要求，但截图必须由新项目重新生成。
export type FobrainVisualRegion =
  | "main-chat"
  | "fresh-main-chat"
  | "process"
  | "evidence"
  | "audit"
  | "internal-details";

type FobrainVisualScenario = {
  readonly fixtureKey: FobrainVisualFixtureKey;
  readonly toolId: "tool.fobrain.current_user_context" | "tool.fobrain.my_permissions";
  readonly prompt: string;
  readonly label: string;
  readonly regions: readonly FobrainVisualRegion[];
};

const requiredRegions = ["main-chat", "fresh-main-chat", "process", "evidence", "audit", "internal-details"] as const;

// JSON fixture 先经过 schema-test 校验，这里只在前端契约边界收窄为 StructuredResult。
const currentUserStructuredResult = currentUserResult as StructuredResult;
const myPermissionsStructuredResult = myPermissionsResult as StructuredResult;

// fobrainVisualScenarios 是视觉测试和验收记录的场景清单，不参与运行时工具选择。
export const fobrainVisualScenarios = [
  {
    fixtureKey: "fobrainCurrentUser",
    toolId: "tool.fobrain.current_user_context",
    prompt: "查看当前 Fobrain 用户信息",
    label: "Fobrain 当前用户",
    regions: requiredRegions
  },
  {
    fixtureKey: "fobrainMyPermissions",
    toolId: "tool.fobrain.my_permissions",
    prompt: "查看我的 Fobrain 权限范围",
    label: "Fobrain 我的权限",
    regions: requiredRegions
  }
] as const satisfies readonly FobrainVisualScenario[];

// fobrainVisualFixtures 使用后端契约 fixture 组装 WorkbenchView，避免前端手写平行事实模型。
export const fobrainVisualFixtures = {
  fobrainCurrentUser: buildFobrainView({
    scenario: fobrainVisualScenarios[0],
    runId: "run-fobrain-current-user-visual",
    toolCallId: "call-fobrain-current-user-visual",
    result: currentUserStructuredResult
  }),
  fobrainMyPermissions: buildFobrainView({
    scenario: fobrainVisualScenarios[1],
    runId: "run-fobrain-my-permissions-visual",
    toolCallId: "call-fobrain-my-permissions-visual",
    result: myPermissionsStructuredResult
  })
} as const satisfies Record<FobrainVisualFixtureKey, WorkbenchView>;

function buildFobrainView({
  scenario,
  runId,
  toolCallId,
  result
}: {
  readonly scenario: FobrainVisualScenario;
  readonly runId: string;
  readonly toolCallId: string;
  readonly result: StructuredResult;
}): WorkbenchView {
  // WorkbenchView 只承载 StructuredResult 的安全投影，raw provider payload 不进入前端 fixture。
  const safeSummary = buildStructuredResultSummary(result);
  return {
    schema_version: "eino_workbench_view.v1",
    workspace_id: "ws-demo",
    run_id: runId,
    status: "succeeded",
    timeline: [
      {
        item_id: `${scenario.fixtureKey}-user`,
        kind: "user_message",
        content: scenario.prompt
      },
      {
        item_id: `${scenario.fixtureKey}-tool`,
        kind: "tool_card",
        tool_call_id: toolCallId,
        status: "succeeded",
        safe_summary: safeSummary,
        structured_result: result
      },
      {
        item_id: `${scenario.fixtureKey}-assistant-final`,
        kind: "assistant_message",
        content: buildStructuredResultAnswer(result)
      }
    ],
    inspector: {
      tabs: ["evidence", "structured", "runtime", "audit"],
      evidence: [
        { label: "场景", value: scenario.label },
        { label: "工具卡", value: "已从 StructuredResult 渲染" }
      ],
      structured: {
        tool_call_id: toolCallId,
        ...(result.metadata.result_ref ? { result_ref: result.metadata.result_ref } : {}),
        structured_result: result
      },
      runtime: {
        status: "succeeded",
        model_label: "fixture-safe"
      },
      audit: [
        {
          schema_version: "eino_audit_event.v1",
          audit_id: `audit-${scenario.fixtureKey}-tool-completed`,
          run_id: runId,
          event_type: "tool",
          safe_summary: `${scenario.label} tool completed from StructuredResult`,
          actor: "tool",
          created_at: "2026-07-02T00:00:00Z"
        }
      ]
    }
  };
}

function buildStructuredResultSummary(result: StructuredResult): string {
  // 摘要只从 StructuredResult 的安全字段生成，避免前端引入第二套事实。
  const title = result.schema_version === "fobrain.tool_result.v2" ? (result.data.title ?? "Fobrain 工具结果") : "Fobrain 工具结果";
  const summary = result.data.summary ?? `状态：${result.status}`;
  return `${title}：${summary}`;
}

function buildStructuredResultAnswer(result: StructuredResult): string {
  // 最终回答复用同一组安全字段，保持工具卡和聊天文本事实同源。
  return `${buildStructuredResultSummary(result)}。`;
}
