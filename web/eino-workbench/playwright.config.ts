import { defineConfig, devices } from "@playwright/test";

const baseURL = process.env.EINO_WORKBENCH_BASE_URL ?? "http://127.0.0.1:8081";

export default defineConfig({
  testDir: "./tests",
  fullyParallel: true,
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 2 : 0,
  reporter: [
    ["list"],
    ["html", { outputFolder: "../../test-results/eino-workbench-playwright-report", open: "never" }],
    ["json", { outputFile: "../../test-results/eino-workbench-playwright-report/results.json" }]
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
    command: "npm --workspace @agent-platform-eino/eino-workbench run dev",
    url: baseURL,
    reuseExistingServer: !process.env.CI,
    timeout: 120_000
  },
  projects: [
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
