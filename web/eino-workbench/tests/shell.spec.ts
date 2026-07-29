import { expect, test } from "@playwright/test";

test("真实纵向查询显示同源摘要且操作记录零请求", async ({ page }) => {
  const fixtureCase = process.env.M1_FIXTURE_CASE ?? "resolved";
  const historyRequests: string[] = [];
  const apiBodies: Promise<string>[] = [];
  let messagePosts = 0;
  page.on("request", (request) => {
    if (request.url().toLowerCase().includes("history")) historyRequests.push(request.url());
    if (request.method() === "POST" && request.url().includes("/messages")) messagePosts += 1;
  });
  page.on("response", (response) => {
    if (response.url().includes("/api/") && response.status() < 400) apiBodies.push(response.text().catch(() => ""));
  });
  await page.goto("/workspaces/ws-workbench");
  await expect(page.getByTestId("workbench-shell")).toBeVisible();
  await page.getByRole("button", { name: "操作记录" }).click();
  await expect(page.getByTestId("history-empty")).toHaveText("操作记录尚未启用");
  expect(historyRequests).toHaveLength(0);

  await page.getByLabel("输入查询").fill("查询新增的漏洞");
  await page.getByRole("button", { name: "发送" }).click();
  if (fixtureCase === "resolved") {
    await expect(page.getByTestId("query-result-card")).toContainText("发现 3 条新增漏洞");
    await expect(page.getByTestId("query-result-card")).toContainText("3 条");
  } else if (fixtureCase === "empty") {
    await expect(page.getByTestId("query-result-card")).toContainText("没有待派发漏洞");
    await expect(page.getByTestId("query-result-card")).toContainText("0 条");
  } else {
    await expect(page.getByTestId("query-result-card")).toHaveCount(0);
    await expect(page.getByRole("alert")).toContainText("无法读取漏洞事实。此次请求不是空结果。");
  }
  await expect.poll(() => apiBodies.length).toBeGreaterThanOrEqual(3);
  await assertNoUnsafeProductText(page);
  for (const body of await Promise.all(apiBodies)) assertNoUnsafeMaterial(body);
  const screenshot = await page.screenshot();
  expect(screenshot.byteLength).toBeGreaterThan(1000);

  // 刷新后从持久事实恢复；不能再次 POST，也不能触发第二次 capability。
  await page.reload();
  if (fixtureCase === "failed") await expect(page.getByRole("alert")).toBeVisible();
  else await expect(page.getByTestId("query-result-card")).toBeVisible();
  expect(messagePosts).toBe(1);
});

test("固定桌面三栏、低宽度与文本放大均不折叠成移动布局", async ({ page }) => {
  await page.goto("/workspaces/ws-workbench");
  const shell = page.getByTestId("workbench-shell");
  const columns = await shell.evaluate((element) => getComputedStyle(element).gridTemplateColumns);
  expect(columns.startsWith("232px ")).toBe(true);
  expect(columns.endsWith(" 360px")).toBe(true);
  await expect(page.getByTestId("workspace-sidebar")).toBeVisible();
  await expect(page.getByTestId("inspector")).toBeVisible();
  await expect(page.getByRole("button", { name: /移动|主题|暗色/ })).toHaveCount(0);

  await page.setViewportSize({ width: 1100, height: 900 });
  const dimensions = await page.evaluate(() => ({ body: document.body.scrollWidth, viewport: window.innerWidth }));
  expect(dimensions.body).toBeGreaterThanOrEqual(1180);
  expect(dimensions.body).toBeGreaterThan(dimensions.viewport);
  await expect(page.getByTestId("workspace-sidebar")).toBeVisible();
  await expect(page.getByTestId("inspector")).toBeVisible();

  await page.evaluate(() => { document.body.style.fontSize = "28px"; });
  await expect(page.getByLabel("输入查询")).toBeVisible();
  await expect(page.getByRole("tab", { name: "事实" })).toBeVisible();
});

async function assertNoUnsafeProductText(page: import("@playwright/test").Page) {
  const text = (await page.getByTestId("workbench-shell").innerText()).toLowerCase();
  assertNoUnsafeMaterial(text);
  expect(text).not.toMatch(/schema_version|run_/);
}

function assertNoUnsafeMaterial(value: string) {
  expect(value.toLowerCase()).not.toMatch(/result_ref|snapshot_item_ref|authorization|bearer|token|credential|provider_locator|raw_payload|actiondraft|tool\.structured_result\.v1/);
}
