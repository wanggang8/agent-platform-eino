import { defineConfig, devices } from "@playwright/test";

const baseURL = process.env.EINO_WORKBENCH_BASE_URL ?? "http://127.0.0.1:8081";
// visual 报告单独落目录，避免普通 smoke 和视觉验收产物互相覆盖。
const reportDir =
  process.env.EINO_WORKBENCH_REPORT_KIND === "visual"
    ? "../../test-results/eino-workbench-visual-report"
    : "../../test-results/eino-workbench-playwright-report";

export default defineConfig({
  testDir: "./tests",
  snapshotPathTemplate: "{testDir}/__screenshots__/{projectName}/{arg}{ext}",
  fullyParallel: true,
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 2 : 0,
  reporter: [
    ["list"],
    ["html", { outputFolder: reportDir, open: "never" }],
    ["json", { outputFile: `${reportDir}/results.json` }]
  ],
  use: {
    baseURL,
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
    video: "retain-on-failure"
  },
  expect: {
    toHaveScreenshot: {
      maxDiffPixelRatio: 0.002,
      threshold: 0.1
    }
  },
  webServer: {
    // Playwright 统一拉起 Vite dev server，保证本地和 CI 使用同一入口。
    command: "npm --workspace @agent-platform-eino/eino-workbench run dev",
    url: baseURL,
    reuseExistingServer: !process.env.CI,
    timeout: 120_000
  },
  projects: [
    // desktop/mobile 两个视口是 Workbench 首轮视觉门禁的固定基线。
    {
      name: "desktop",
      use: {
        ...devices["Desktop Chrome"],
        viewport: { width: 1440, height: 900 }
      }
    },
    {
      name: "mobile",
      use: {
        ...devices["Pixel 7"],
        viewport: { width: 390, height: 844 }
      }
    }
  ]
});
