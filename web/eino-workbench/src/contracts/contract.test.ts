import { contractSchemas } from "./generated";
import { describe, expect, it } from "vitest";

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
    for (const id of requiredSchemaIds) {
      expect(schemaIds.has(id), `missing generated contract schema ${id}`).toBe(true);
    }
  });
});
