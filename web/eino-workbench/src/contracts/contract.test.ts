import { contractSchemas } from "./generated";
import { describe, expect, it } from "vitest";

// requiredSchemaIds 是前端必须能消费的运行期契约集合，防止生成脚本漏掉核心 schema。
const requiredSchemaIds = [
  "https://agent-platform-eino.local/schemas/eino_action_result.v1.schema.json",
  "https://agent-platform-eino.local/schemas/eino_product_facts.v1.schema.json",
  "https://agent-platform-eino.local/schemas/eino_workbench_stream_event.v1.schema.json",
  "https://agent-platform-eino.local/schemas/eino_workbench_view.v1.schema.json",
  "https://agent-platform-eino.local/schemas/fobrain/tool_result.v2.schema.json",
  "https://agent-platform-eino.local/schemas/tool.structured_result.v1.schema.json"
] as const;

const schemaIds = new Set(contractSchemas.map((schema) => schema.id));

describe("generated contracts", () => {
  it("includes Workbench runtime schema entries", () => {
    // 这里验证生成物覆盖范围，不直接校验 schema 内容；内容由 schema validate 脚本负责。
    for (const id of requiredSchemaIds) {
      expect(schemaIds.has(id), `missing generated contract schema ${id}`).toBe(true);
    }
  });
});
