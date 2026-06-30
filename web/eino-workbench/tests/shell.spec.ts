import { expect, test } from "@playwright/test";

test("@shell renders responsive workbench regions", async ({ page }, testInfo) => {
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
});

test("@visual captures fixture shell", async ({ page }) => {
  await page.goto("/workspaces/ws-demo");
  await expect(page.getByTestId("workbench-shell")).toHaveScreenshot("shell-fixture.png", {
    mask: [page.getByText(/run-/), page.getByText(/\d{2}:\d{2}/)],
    maskColor: "#06141e"
  });
});
