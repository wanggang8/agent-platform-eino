import { describe, expect, it } from "vitest";
import { fobrainVisualFixtures, fobrainVisualScenarios } from "./fobrainVisualFixtures";

const requiredRegions = ["main-chat", "fresh-main-chat", "process", "evidence", "audit", "internal-details"] as const;

describe("fobrain visual fixtures", () => {
  it("defines Batch A business-read scenarios with six visual regions", () => {
    // Batch A 视觉验收必须继承旧最终矩阵的六区域覆盖，但数据必须来自新项目 fixture。
    expect(fobrainVisualScenarios.map((scenario) => scenario.toolId)).toEqual([
      "tool.fobrain.current_user_context",
      "tool.fobrain.my_permissions"
    ]);

    for (const scenario of fobrainVisualScenarios) {
      expect(scenario.regions).toEqual(requiredRegions);
      expect(fobrainVisualFixtures[scenario.fixtureKey].schema_version).toBe("eino_workbench_view.v1");
      expect(fobrainVisualFixtures[scenario.fixtureKey].timeline.some((item) => item.kind === "tool_card")).toBe(true);
    }
  });

  it("keeps Fobrain visual fixtures free of unsafe material", () => {
    const encoded = JSON.stringify(fobrainVisualFixtures).toLowerCase();
    for (const forbidden of ["authorization", "bearer ", "api_token", "credential_ref", "raw provider", "raw body", "resume_token", "超级管理员"]) {
      expect(encoded.includes(forbidden), `fixture leaked ${forbidden}`).toBe(false);
    }
  });

  it("derives visible Fobrain copy from StructuredResult fields only", () => {
    // 该断言防止前端 fixture 在 StructuredResult 之外补写产品事实。
    for (const scenario of fobrainVisualScenarios) {
      const view = fobrainVisualFixtures[scenario.fixtureKey];
      const toolCard = view.timeline.find((item) => item.kind === "tool_card");
      const finalMessage = view.timeline.find((item) => item.kind === "assistant_message");
      const result = toolCard?.structured_result;

      expect(result).toBeDefined();
      expect(result?.schema_version).toBe("fobrain.tool_result.v2");
      if (!result || result.schema_version !== "fobrain.tool_result.v2") {
        throw new Error(`${scenario.fixtureKey} must use Fobrain StructuredResult`);
      }

      const expectedSummary = `${result.data.title}：${result.data.summary}`;
      expect(toolCard?.safe_summary).toBe(expectedSummary);
      expect(finalMessage?.content).toBe(`${expectedSummary}。`);
    }
  });
});
