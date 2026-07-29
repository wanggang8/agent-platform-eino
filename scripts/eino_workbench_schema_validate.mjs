import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import Ajv2020 from "ajv/dist/2020.js";
import addFormats from "ajv-formats";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const m1Schemas = [
  "docs/schemas/tool.structured_result.v2.schema.json",
  "docs/schemas/eino_query_result_snapshot.v2.schema.json",
  "docs/schemas/eino_product_facts.v2.schema.json",
  "docs/schemas/eino_workbench_view.v2.schema.json",
  "docs/schemas/eino_workbench_stream_event.v2.schema.json",
  "docs/schemas/new_vulnerability_fixture_bundle.v1.schema.json"
];
const forbiddenMarkers = [
  "authorization",
  "bearer ",
  "api_key",
  "token",
  "cookie",
  "credential",
  "raw_provider",
  "provider_url",
  "provider_locator",
  "resume_ref",
  "actiondraft",
  "tool.structured_result.v1"
];

// readJSON 统一从仓库根目录读取 JSON，避免校验器依赖调用位置。
function readJSON(file) {
  return JSON.parse(fs.readFileSync(path.join(root, file), "utf8"));
}

// schemaKey 与本地 $ref 使用同一相对命名空间。
function schemaKey(file) {
  return path.relative("docs/schemas", file).replaceAll(path.sep, "/");
}

function buildValidator() {
  const ajv = new Ajv2020({ allErrors: true, strict: true, validateFormats: true });
  addFormats(ajv);
  for (const file of m1Schemas) {
    const schema = readJSON(file);
    if (JSON.stringify(schema).includes('"additionalProperties":true')) {
      throw new Error(`${file} opens additionalProperties`);
    }
    ajv.addSchema(schema, schemaKey(file));
  }
  return ajv;
}

function validateFixture(ajv) {
  const manifest = readJSON("docs/fixtures/manifest.json");
  if (manifest.schema_version !== "eino_fixture_manifest.v1" || manifest.fixtures.length !== 1) {
    throw new Error("M1 fixture manifest must contain exactly one versioned bundle");
  }
  const entry = manifest.fixtures[0];
  if (entry.fixture !== "docs/fixtures/new-vulnerability-walking-skeleton.v1.json") {
    throw new Error("M1 fixture bundle path mismatch");
  }
  const fixture = readJSON(entry.fixture);
  const schema = readJSON(entry.schema);
  const validate = ajv.getSchema(schema.$id) ?? ajv.compile(schema);
  if (!validate(fixture)) {
    throw new Error(`${entry.fixture} failed schema:\n${ajv.errorsText(validate.errors, { separator: "\n" })}`);
  }
  const encoded = JSON.stringify(fixture).toLowerCase();
  const hit = forbiddenMarkers.find((marker) => encoded.includes(marker));
  if (hit) throw new Error(`${entry.fixture} contains forbidden marker ${hit}`);

  const { resolved, empty, failed } = fixture.cases;
  if (resolved.data.count !== resolved.data.items.length) throw new Error("resolved count/items mismatch");
  if (empty.data.summary !== "没有待派发漏洞" || empty.data.count !== 0 || empty.data.items.length !== 0) {
    throw new Error("empty case must be a complete zero result");
  }
  if (failed.message !== "无法读取漏洞事实。此次请求不是空结果。") {
    throw new Error("failed case must remain distinct from empty");
  }
  const sorted = [...resolved.data.items].sort((a, b) => {
    const byTime = b.discovered_at.localeCompare(a.discovered_at);
    return byTime || a.snapshot_item_ref.localeCompare(b.snapshot_item_ref);
  });
  if (JSON.stringify(sorted) !== JSON.stringify(resolved.data.items)) {
    throw new Error("resolved items must use discovered_at DESC + snapshot_item_ref ASC");
  }
}

function lintOpenAPI() {
  const doc = readJSON("docs/api/eino-workbench.openapi.json");
  if (doc.openapi !== "3.1.2") throw new Error("OpenAPI version must be 3.1.2");
  const expectedPaths = [
    "/api/workspaces/{workspace_id}/messages",
    "/api/workspaces/{workspace_id}/runs/{run_id}",
    "/api/workspaces/{workspace_id}/runs/{run_id}/stream",
    "/api/workspaces/{workspace_id}/views/current",
    "/healthz",
    "/readyz"
  ];
  if (JSON.stringify(Object.keys(doc.paths ?? {}).sort()) !== JSON.stringify(expectedPaths.sort())) {
    throw new Error("M1 OpenAPI path set drifted");
  }
  const encoded = JSON.stringify(doc);
  for (const marker of ["action", "resume", "replay", "lifecycle", ".v1.schema.json"]) {
    if (encoded.toLowerCase().includes(marker)) throw new Error(`M1 OpenAPI contains out-of-scope marker ${marker}`);
  }
  for (const methods of Object.values(doc.paths ?? {})) {
    for (const operation of Object.values(methods)) {
      if (!operation.operationId) throw new Error("OpenAPI operation missing operationId");
      for (const match of JSON.stringify(operation).matchAll(/"\$ref":"([^"]+)"/g)) {
        const ref = match[1].split("#")[0];
        if (!ref) continue;
        if (!fs.existsSync(path.resolve(root, "docs/api", ref))) throw new Error(`missing OpenAPI ref ${ref}`);
      }
    }
  }
}

const ajv = buildValidator();
for (const file of m1Schemas) {
  const schema = readJSON(file);
  if (!ajv.getSchema(schema.$id)) ajv.compile(schema);
}
validateFixture(ajv);
lintOpenAPI();
console.log(`validated ${m1Schemas.length} M1 schemas and one fixture bundle`);
