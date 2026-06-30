import { contractSchemas } from "./generated";

const requiredSchemaIds = [
  "https://agent-platform-eino.local/schemas/eino_action_result.v1.schema.json",
  "https://agent-platform-eino.local/schemas/eino_product_facts.v1.schema.json",
  "https://agent-platform-eino.local/schemas/eino_workbench_stream_event.v1.schema.json",
  "https://agent-platform-eino.local/schemas/eino_workbench_view.v1.schema.json",
  "https://agent-platform-eino.local/schemas/fobrain/tool_result.v2.schema.json",
  "https://agent-platform-eino.local/schemas/tool.structured_result.v1.schema.json"
] as const;

const schemaIds = new Set(contractSchemas.map((schema) => schema.id));

for (const id of requiredSchemaIds) {
  if (!schemaIds.has(id)) {
    throw new Error(`missing generated contract schema ${id}`);
  }
}
