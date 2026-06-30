import process from "node:process";

const [command = "unknown", phase = "future phase"] = process.argv.slice(2);

console.error(
  `eino-workbench:${command} is intentionally unavailable until ${phase}. ` +
    "This command exists to reserve the documented interface and must not be treated as a passing check."
);
process.exit(2);
