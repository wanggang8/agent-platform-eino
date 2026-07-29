import { defineConfig, devices } from "@playwright/test";

const baseURL = process.env.EINO_WORKBENCH_BASE_URL ?? "http://127.0.0.1:8081";

export default defineConfig({
  testDir: "./tests",
  fullyParallel: false,
  workers: 1,
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 1 : 0,
  reporter: [["list"], ["html", { outputFolder: "../../test-results/eino-workbench-playwright-report", open: "never" }]],
  use: {
    baseURL,
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
    video: "retain-on-failure"
  },
  webServer: {
    command: "bash ../../scripts/m1_e2e_server.sh",
    url: baseURL,
    reuseExistingServer: false,
    timeout: 120_000
  },
  projects: [{
    name: "desktop",
    use: { ...devices["Desktop Chrome"], viewport: { width: 1440, height: 900 } }
  }]
});
