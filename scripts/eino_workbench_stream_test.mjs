import { spawnSync } from "node:child_process";
import process from "node:process";

const grepIndex = process.argv.indexOf("--grep");
const grep = grepIndex >= 0 ? process.argv[grepIndex + 1] : "";
const testNameArgs = grep ? ["--testNamePattern", grep] : [];

// stream-test 是 Phase 6.4 的专用入口，当前聚焦前端 reducer 是否正确消费 pending SSE patch。
const result = spawnSync("npm", [
  "--workspace",
  "@agent-platform-eino/eino-workbench",
  "run",
  "test",
  "--",
  "src/features/workbench/state/workbenchReducer.test.ts",
  ...testNameArgs
], {
  stdio: "inherit",
  cwd: process.cwd()
});

process.exit(result.status ?? 1);
