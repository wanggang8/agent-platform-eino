import { expect, test } from "@playwright/test";

test("@shell renders responsive workbench regions", async ({ page }, testInfo) => {
  // 响应式验收只检查稳定区域和安全文案，不绑定旧项目 DOM 结构。
  await page.goto("/workspaces/ws-demo");
  await expect(page.getByTestId("workbench-shell")).toBeVisible();
  await expect(page.getByTestId("chat-timeline")).toBeVisible();
  await expect(page.getByTestId("chat-composer")).toBeVisible();
  if (testInfo.project.name === "desktop") {
    await expect(page.getByTestId("workspace-sidebar")).toBeVisible();
    await expect(page.getByTestId("inspector")).toBeVisible();
  } else {
    await expect(page.getByTestId("workspace-sidebar")).toBeHidden();
    await expect(page.getByRole("button", { name: "证据" })).toBeVisible();
  }
  await expect(page.getByText("tool.fobrain")).toHaveCount(0);
  await assertProductVisibleText(page);
});

test("@visual captures fixture shell", async ({ page }) => {
  // 视觉截图屏蔽动态 run/time，避免非产品差异造成基线漂移。
  await page.goto("/workspaces/ws-demo");
  await assertProductVisibleText(page);
  await expect(page.getByTestId("workbench-shell")).toHaveScreenshot("shell-fixture.png", {
    mask: [page.getByText(/run-/), page.getByText(/\d{2}:\d{2}/)],
    maskColor: "#06141e"
  });
});

async function assertProductVisibleText(page: import("@playwright/test").Page) {
  // shell 级验收禁止产品首屏出现内部英文、tool id 或 JSON 结构。
  const visibleText = (await page.getByTestId("workbench-shell").innerText()).toLowerCase();
  expect(visibleText).not.toMatch(/[a-z]/);
  expect(visibleText).not.toMatch(/[{}]|tool\.|schema_version|display_type|structuredresult|product facts|safe projection|run-/);
}
