import fs from "node:fs";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";
import Ajv2020 from "ajv/dist/2020.js";
import addFormats from "ajv-formats";

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(scriptDir, "..");
const forbiddenMarkers = [
  "resume_token",
  "credential_ref",
  "authorization",
  "bearer ",
  "api_key",
  "raw_provider_body",
  "raw provider payload"
];
const fobrainReadonlyToolIds = [
  "tool.fobrain.current_user_context",
  "tool.fobrain.my_permissions",
  "tool.fobrain.list_assets_by_owner",
  "tool.fobrain.list_vulnerabilities_by_owner",
  "tool.fobrain.list_assets_by_department",
  "tool.fobrain.list_vulnerabilities_by_department",
  "tool.fobrain.list_assets_by_ip",
  "tool.fobrain.list_vulnerabilities_by_ip",
  "tool.fobrain.get_asset_detail",
  "tool.fobrain.get_vulnerability_detail",
  "tool.fobrain.business_risk_summary",
  "tool.fobrain.business_list",
  "tool.fobrain.external_high_risk_assets",
  "tool.fobrain.threat_relevance_list",
  "tool.fobrain.vulnerability_status_summary",
  "tool.fobrain.pending_tickets",
  "tool.fobrain.ip_stats",
  "tool.fobrain.vul_stats",
  "tool.fobrain.my_assets",
  "tool.fobrain.my_department_assets",
  "tool.fobrain.my_vulnerabilities",
  "tool.fobrain.my_department_vulnerabilities",
  "tool.fobrain.my_business_systems",
  "tool.fobrain.my_important_business_systems"
];
const visualBlockIds = [
  "shell",
  "sidebar",
  "timeline",
  "composer",
  "tool-card",
  "approval-card",
  "clarification-card",
  "inspector"
];

function readJSON(filePath) {
  return JSON.parse(fs.readFileSync(path.join(root, filePath), "utf8"));
}

function walk(dir, predicate = () => true) {
  const out = [];
  for (const entry of fs.readdirSync(path.join(root, dir), { withFileTypes: true })) {
    const rel = path.join(dir, entry.name);
    if (entry.isDirectory()) out.push(...walk(rel, predicate));
    if (entry.isFile() && predicate(rel)) out.push(rel);
  }
  return out.sort();
}

function schemaKey(filePath) {
  return path.relative("docs/schemas", filePath).replaceAll(path.sep, "/");
}

function loadAjv() {
  const ajv = new Ajv2020({ allErrors: true, strict: true, validateFormats: true });
  addFormats(ajv);
  const schemaFiles = walk("docs/schemas", (rel) => rel.endsWith(".json"));
  for (const file of schemaFiles) {
    const schema = readJSON(file);
    ajv.addSchema(schema, schemaKey(file));
  }
  return { ajv, schemaFiles };
}

function assertNoOpenObjects(schemaFiles) {
  const failures = [];
  for (const file of schemaFiles) {
    const text = fs.readFileSync(path.join(root, file), "utf8");
    if (text.includes('"additionalProperties": true')) {
      failures.push(file);
    }
  }
  if (failures.length > 0) {
    throw new Error(`open additionalProperties found:\n${failures.join("\n")}`);
  }
}

function assertNoForbiddenMarkers(filePath, value) {
  const encoded = JSON.stringify(value).toLowerCase();
  const hit = forbiddenMarkers.find((marker) => encoded.includes(marker));
  if (hit) {
    throw new Error(`${filePath} contains forbidden marker ${hit}`);
  }
}

function validateFixtures(ajv) {
  const manifest = readJSON("docs/fixtures/manifest.json");
  if (manifest.schema_version !== "eino_fixture_manifest.v1") {
    throw new Error("fixture manifest schema_version mismatch");
  }
  for (const item of manifest.fixtures) {
    const schema = readJSON(item.schema);
    const fixture = readJSON(item.fixture);
    const validate = ajv.getSchema(schema.$id) ?? ajv.compile(schema);
    if (!validate(fixture)) {
      throw new Error(`${item.fixture} failed ${item.schema}:\n${ajv.errorsText(validate.errors, { separator: "\n" })}`);
    }
    assertNoForbiddenMarkers(item.fixture, fixture);
  }
}

function validateFobrainToolMatrix(ajv) {
  const matrix = readJSON("docs/fixtures/fobrain/tool-matrix-24.json");
  const expected = new Set(fobrainReadonlyToolIds);
  const seen = new Set();
  const resultSchema = readJSON("docs/schemas/fobrain/tool_result.v2.schema.json");
  const validateResult = ajv.getSchema(resultSchema.$id) ?? ajv.compile(resultSchema);
  for (const tool of matrix.tools ?? []) {
    if (seen.has(tool.tool_id)) {
      throw new Error(`duplicate Fobrain tool matrix entry ${tool.tool_id}`);
    }
    seen.add(tool.tool_id);
    if (!fs.existsSync(path.join(root, tool.fixture))) {
      throw new Error(`Fobrain tool matrix fixture missing for ${tool.tool_id}: ${tool.fixture}`);
    }
    const fixture = readJSON(tool.fixture);
    if (!validateResult(fixture)) {
      throw new Error(`${tool.fixture} failed Fobrain result schema:\n${ajv.errorsText(validateResult.errors, { separator: "\n" })}`);
    }
    if (fixture.tool_id !== tool.tool_id) {
      throw new Error(`${tool.fixture} tool_id ${fixture.tool_id} does not match matrix ${tool.tool_id}`);
    }
    if (fixture.display_type !== tool.display_type) {
      throw new Error(`${tool.fixture} display_type ${fixture.display_type} does not match matrix ${tool.display_type}`);
    }
  }
  const missing = [...expected].filter((toolID) => !seen.has(toolID));
  const extra = [...seen].filter((toolID) => !expected.has(toolID));
  if (missing.length > 0 || extra.length > 0) {
    throw new Error(`Fobrain 24 tool matrix mismatch\nmissing: ${missing.join(", ")}\nextra: ${extra.join(", ")}`);
  }
}

function validateVisualEvidenceMatrix() {
  const matrix = readJSON("docs/fixtures/visual-evidence-matrix.json");
  const expected = new Set(visualBlockIds);
  const seen = new Set();
  for (const block of matrix.blocks ?? []) {
    if (seen.has(block.block_id)) {
      throw new Error(`duplicate visual block matrix entry ${block.block_id}`);
    }
    seen.add(block.block_id);
    if (block.target_crop_path && !fs.existsSync(path.join(root, block.target_crop_path))) {
      throw new Error(`visual block target reference missing for ${block.block_id}: ${block.target_crop_path}`);
    }
  }
  const missing = [...expected].filter((blockID) => !seen.has(blockID));
  const extra = [...seen].filter((blockID) => !expected.has(blockID));
  if (missing.length > 0 || extra.length > 0) {
    throw new Error(`visual block matrix mismatch\nmissing: ${missing.join(", ")}\nextra: ${extra.join(", ")}`);
  }
  for (const state of matrix.states ?? []) {
    if (!fs.existsSync(path.join(root, state.fixture))) {
      throw new Error(`visual state fixture missing for ${state.state_id}: ${state.fixture}`);
    }
  }
}

function lintOpenAPI() {
  const doc = readJSON("docs/api/eino-workbench.openapi.json");
  if (doc.openapi !== "3.1.2") throw new Error("OpenAPI version must be 3.1.2");
  const operationIds = new Set();
  for (const [apiPath, methods] of Object.entries(doc.paths ?? {})) {
    for (const [method, operation] of Object.entries(methods)) {
      if (!operation.operationId) throw new Error(`${method.toUpperCase()} ${apiPath} missing operationId`);
      if (operationIds.has(operation.operationId)) throw new Error(`duplicate operationId ${operation.operationId}`);
      operationIds.add(operation.operationId);
      const encoded = JSON.stringify(operation);
      for (const match of encoded.matchAll(/"\$ref":"([^"]+)"/g)) {
        const ref = match[1].split("#")[0];
        if (!ref) continue;
        const refPath = path.normalize(path.join("docs/api", ref));
        if (!fs.existsSync(path.join(root, refPath))) {
          throw new Error(`${operation.operationId} references missing schema ${ref}`);
        }
      }
    }
  }
}

function main() {
  const { ajv, schemaFiles } = loadAjv();
  for (const file of schemaFiles) {
    const schema = readJSON(file);
    if (!ajv.getSchema(schema.$id)) {
      ajv.compile(schema);
    }
  }
  assertNoOpenObjects(schemaFiles);
  validateFixtures(ajv);
  validateFobrainToolMatrix(ajv);
  validateVisualEvidenceMatrix();
  lintOpenAPI();
  console.log(`validated ${schemaFiles.length} schemas and fixture manifest`);
}

main();
