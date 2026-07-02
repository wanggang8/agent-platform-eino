import fs from "node:fs";
import process from "node:process";
import Ajv2020 from "ajv/dist/2020.js";
import addFormats from "ajv-formats";

function usage() {
  console.error("usage: node scripts/eino_workbench_report_validate.mjs --schema <schema.json> --report <report.json>");
  process.exit(2);
}

let schemaPath = "";
let reportPath = "";
for (let i = 2; i < process.argv.length; i += 1) {
  if (process.argv[i] === "--schema") {
    schemaPath = process.argv[i + 1] ?? "";
    i += 1;
  } else if (process.argv[i] === "--report") {
    reportPath = process.argv[i + 1] ?? "";
    i += 1;
  } else {
    usage();
  }
}
if (!schemaPath || !reportPath) usage();

const ajv = new Ajv2020({ allErrors: true, strict: true, validateFormats: true });
addFormats(ajv);
const schema = JSON.parse(fs.readFileSync(schemaPath, "utf8"));
const report = JSON.parse(fs.readFileSync(reportPath, "utf8"));
const validate = ajv.compile(schema);
if (!validate(report)) {
  console.error(ajv.errorsText(validate.errors, { separator: "\n" }));
  process.exit(1);
}
console.log(`${reportPath} valid for ${schemaPath}`);
