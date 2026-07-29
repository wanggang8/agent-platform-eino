import { contractSchemas } from "./generated";
import { describe, expect, it } from "vitest";

// requiredSchemaIds 是前端必须能消费的运行期契约集合，防止生成脚本漏掉核心 schema。
const requiredSchemaIds = [
  "https://agent-platform-eino.local/schemas/eino_product_facts.v2.schema.json",
  "https://agent-platform-eino.local/schemas/eino_query_result_snapshot.v2.schema.json",
  "https://agent-platform-eino.local/schemas/eino_workbench_stream_event.v2.schema.json",
  "https://agent-platform-eino.local/schemas/eino_workbench_view.v2.schema.json",
  "https://agent-platform-eino.local/schemas/tool.structured_result.v2.schema.json"
] as const;

const schemaIds = new Set<string>(contractSchemas.map((schema) => schema.id));

describe("generated contracts", () => {
  it("只包含 M1 walking skeleton 的 v2 运行期契约", () => {
    // M1 不允许用 v1/v2 union 或后续 Action 契约掩盖新事实链缺失。
    for (const id of requiredSchemaIds) {
      expect(schemaIds.has(id), `missing generated contract schema ${id}`).toBe(true);
    }
    expect(schemaIds.has("https://agent-platform-eino.local/schemas/eino_product_facts.v1.schema.json")).toBe(false);
    expect(schemaIds.has("https://agent-platform-eino.local/schemas/eino_action_result.v1.schema.json")).toBe(false);
  });
});
