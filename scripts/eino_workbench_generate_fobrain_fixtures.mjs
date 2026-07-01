import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(scriptDir, "..");
const matrixPath = path.join(root, "docs/fixtures/fobrain/tool-matrix-24.json");
const matrix = JSON.parse(fs.readFileSync(matrixPath, "utf8"));

// dataFor 根据工具展示类型生成安全 fixture 数据，不使用旧项目 raw payload。
function dataFor(tool) {
  const title = tool.display_name_zh;
  if (tool.display_type === "entity_collection") {
    return {
      title,
      summary: `${title} fixture`,
      columns: [
        {"key": "display_name", "label": "名称", "type": "text"},
        {"key": "status", "label": "状态", "type": "status"}
      ],
      items: [
        {
          "entity_ref": `${tool.entity_type}:fobrain:fixture-1`,
          "display_name": `${title} 示例`,
          "status": "resolved"
        }
      ],
      pagination: {"page": 1, "page_size": 20, "total": 1, "has_more": false}
    };
  }
  if (tool.display_type === "entity_detail") {
    return {
      title,
      summary: `${title} fixture`,
      facts: [
        {"key": "display_name", "label": "名称", "value": `${title} 示例`, "type": "text"},
        {"key": "status", "label": "状态", "value": "resolved", "type": "status"}
      ]
    };
  }
  if (tool.display_type === "metrics_summary") {
    return {
      title,
      summary: `${title} fixture`,
      metrics: [
        {"label": "总数", "value": 1, "tone": "info"}
      ],
      narrative: {
        "key_findings": [
          {"label": "摘要", "summary": `${title} 示例指标`, "tone": "info"}
        ]
      }
    };
  }
  return {
    title,
    summary: `${title} fixture`
  };
}

// fixtureFor 生成符合 fobrain.tool_result.v2 的 StructuredResult fixture。
function fixtureFor(tool) {
  return {
    "schema_version": "fobrain.tool_result.v2",
    "tool_id": tool.tool_id,
    "display_type": tool.display_type,
    "entity_type": tool.entity_type,
    "status": tool.result_status,
    "query": {},
    "data": dataFor(tool),
    "metadata": {
      "safe": true,
      "source": "fobrain",
      "result_ref": `result:fobrain:${tool.tool_id.replace(/^tool\\.fobrain\\./, "").replaceAll("_", "-")}:fixture`
    }
  };
}

let written = 0;
// 只补缺失 fixture，避免覆盖已经人工验收过的样例。
for (const tool of matrix.tools) {
  const target = path.join(root, tool.fixture);
  if (fs.existsSync(target)) continue;
  fs.mkdirSync(path.dirname(target), { recursive: true });
  fs.writeFileSync(target, `${JSON.stringify(fixtureFor(tool), null, 2)}\n`);
  written += 1;
}

console.log(`generated ${written} missing Fobrain result fixtures`);
