import type { StructuredResult, WorkbenchView } from "../contracts/generated";
import assetDetailResult from "../../../../docs/fixtures/fobrain/get-asset-detail-result.json";
import businessRiskSummaryResult from "../../../../docs/fixtures/fobrain/business-risk-summary-result.json";
import connectorSecurityResult from "../../../../docs/fixtures/fobrain/connector-security-result.json";
import currentUserResult from "../../../../docs/fixtures/fobrain/current-user-context-result.json";
import myPermissionsResult from "../../../../docs/fixtures/fobrain/my-permissions-result.json";
import threatRelevanceListResult from "../../../../docs/fixtures/fobrain/threat-relevance-list-result.json";
import vulnerabilityDetailResult from "../../../../docs/fixtures/fobrain/get-vulnerability-detail-result.json";

// FobrainVisualFixtureKey 固定当前 Batch A/E 视觉 fixture，包含 connector 和详情风险工具。
export type FobrainVisualFixtureKey =
  | "fobrainConnectorSecurity"
  | "fobrainCurrentUser"
  | "fobrainMyPermissions"
  | "fobrainAssetDetail"
  | "fobrainVulnerabilityDetail"
  | "fobrainBusinessRiskSummary"
  | "fobrainThreatRelevanceList";

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
  readonly toolId:
    | "connector.fobrain.security"
    | "tool.fobrain.current_user_context"
    | "tool.fobrain.my_permissions"
    | "tool.fobrain.get_asset_detail"
    | "tool.fobrain.get_vulnerability_detail"
    | "tool.fobrain.business_risk_summary"
    | "tool.fobrain.threat_relevance_list";
  readonly prompt: string;
  readonly label: string;
  readonly regions: readonly FobrainVisualRegion[];
};

const requiredRegions = ["main-chat", "fresh-main-chat", "process", "evidence", "audit", "internal-details"] as const;

// JSON fixture 先经过 schema-test 校验，这里只在前端契约边界收窄为 StructuredResult。
const connectorSecurityStructuredResult = connectorSecurityResult as StructuredResult;
const currentUserStructuredResult = currentUserResult as StructuredResult;
const myPermissionsStructuredResult = myPermissionsResult as StructuredResult;
const assetDetailStructuredResult = assetDetailResult as StructuredResult;
const vulnerabilityDetailStructuredResult = vulnerabilityDetailResult as StructuredResult;
const businessRiskSummaryStructuredResult = businessRiskSummaryResult as StructuredResult;
const threatRelevanceListStructuredResult = threatRelevanceListResult as StructuredResult;

// fobrainVisualScenarios 是视觉测试和验收记录的场景清单，不参与运行时工具选择。
export const fobrainVisualScenarios = [
  {
    fixtureKey: "fobrainConnectorSecurity",
    toolId: "connector.fobrain.security",
    prompt: "查看 Fobrain 连接器状态",
    label: "Fobrain 连接器",
    regions: requiredRegions
  },
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
  },
  {
    fixtureKey: "fobrainAssetDetail",
    toolId: "tool.fobrain.get_asset_detail",
    prompt: "查看 Fobrain 资产详情",
    label: "Fobrain 资产详情",
    regions: requiredRegions
  },
  {
    fixtureKey: "fobrainVulnerabilityDetail",
    toolId: "tool.fobrain.get_vulnerability_detail",
    prompt: "查看 Fobrain 漏洞详情",
    label: "Fobrain 漏洞详情",
    regions: requiredRegions
  },
  {
    fixtureKey: "fobrainBusinessRiskSummary",
    toolId: "tool.fobrain.business_risk_summary",
    prompt: "汇总 Fobrain 业务风险",
    label: "Fobrain 业务风险",
    regions: requiredRegions
  },
  {
    fixtureKey: "fobrainThreatRelevanceList",
    toolId: "tool.fobrain.threat_relevance_list",
    prompt: "查看 Fobrain 威胁关联资产",
    label: "Fobrain 威胁关联",
    regions: requiredRegions
  }
] as const satisfies readonly FobrainVisualScenario[];

// fobrainVisualFixtures 使用后端契约 fixture 组装 WorkbenchView，避免前端手写平行事实模型。
export const fobrainVisualFixtures = {
  fobrainConnectorSecurity: buildFobrainView({
    scenario: fobrainVisualScenarios[0],
    runId: "run-fobrain-connector-security-visual",
    toolCallId: "call-fobrain-connector-security-visual",
    result: connectorSecurityStructuredResult
  }),
  fobrainCurrentUser: buildFobrainView({
    scenario: fobrainVisualScenarios[1],
    runId: "run-fobrain-current-user-visual",
    toolCallId: "call-fobrain-current-user-visual",
    result: currentUserStructuredResult
  }),
  fobrainMyPermissions: buildFobrainView({
    scenario: fobrainVisualScenarios[2],
    runId: "run-fobrain-my-permissions-visual",
    toolCallId: "call-fobrain-my-permissions-visual",
    result: myPermissionsStructuredResult
  }),
  fobrainAssetDetail: buildFobrainView({
    scenario: fobrainVisualScenarios[3],
    runId: "run-fobrain-asset-detail-visual",
    toolCallId: "call-fobrain-asset-detail-visual",
    result: assetDetailStructuredResult
  }),
  fobrainVulnerabilityDetail: buildFobrainView({
    scenario: fobrainVisualScenarios[4],
    runId: "run-fobrain-vulnerability-detail-visual",
    toolCallId: "call-fobrain-vulnerability-detail-visual",
    result: vulnerabilityDetailStructuredResult
  }),
  fobrainBusinessRiskSummary: buildFobrainView({
    scenario: fobrainVisualScenarios[5],
    runId: "run-fobrain-business-risk-summary-visual",
    toolCallId: "call-fobrain-business-risk-summary-visual",
    result: businessRiskSummaryStructuredResult
  }),
  fobrainThreatRelevanceList: buildFobrainView({
    scenario: fobrainVisualScenarios[6],
    runId: "run-fobrain-threat-relevance-list-visual",
    toolCallId: "call-fobrain-threat-relevance-list-visual",
    result: threatRelevanceListStructuredResult
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
