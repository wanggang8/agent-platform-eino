import { expect, test, type Locator, type Page } from "@playwright/test";
import visualFixtures from "./visual-fixtures/workbench-blocks.json" assert { type: "json" };

type VisualBlockId = keyof typeof visualFixtures.blocks;
type VisualState = (typeof visualFixtures.states)[number];

// 视觉用例从验收 fixture 驱动，避免把旧项目 DOM 或截图尺寸硬编码进测试逻辑。
test.describe("@visual Workbench block baselines", () => {
  for (const state of visualFixtures.states) {
    test(`@visual ${state.state_id} required blocks`, async ({ page }, testInfo) => {
      await page.goto("/workspaces/ws-demo");
      await selectFixture(page, state);
      await prepareState(page, state, testInfo.project.name);

      for (const blockId of state.required_blocks as VisualBlockId[]) {
        const block = visualFixtures.blocks[blockId];
        if (block.desktop_only && testInfo.project.name !== "desktop") {
          continue;
        }
        if (blockId === "inspector" && testInfo.project.name !== "desktop") {
          await page.getByRole("button", { name: "证据" }).click();
        }

        const locator = page.locator(block.selector);
        await expect(locator, `${state.state_id}:${blockId} missing. ${state.blocker}`).toBeVisible();
        await expect(locator).toHaveScreenshot(`${state.state_id}-${blockId}.png`, {
          mask: screenshotMasks(page, locator),
          maskColor: "#06141e",
          animations: "disabled"
        });
      }
    });
  }
});

async function selectFixture(page: Page, state: VisualState) {
  // fixture label 是验收配置，不参与业务判断。
  await page.getByRole("button", { name: state.fixture_label }).click();
}

async function prepareState(page: Page, state: VisualState, projectName: string) {
  // 只做截图前的 UI 展开/切换，不修改 Product Facts fixture。
  if (state.state_id === "tool-collapsed") {
    await page.getByTestId("tool-card").locator(".tool-card-header").click();
  }
  if (state.state_id === "inspector-tabs") {
    if (projectName !== "desktop") {
      await page.getByRole("button", { name: "证据" }).click();
    }
    await page.getByRole("tab", { name: "运行详情" }).click();
  }
}

function screenshotMasks(page: Page, scope: Locator) {
  // mask 动态时间和页脚信息，保留真正的视觉结构差异。
  return [
    scope.locator(".conversation-time"),
    scope.locator(".sidebar-footer"),
    page.locator(".header-actions span").first()
  ];
}
