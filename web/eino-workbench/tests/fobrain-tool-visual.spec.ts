import { expect, test, type Locator, type Page } from "@playwright/test";

const fobrainVisualScenarios = [
  {
    fixtureKey: "fobrainConnectorSecurity",
    prompt: "查看 Fobrain 连接器状态",
    label: "Fobrain 连接器"
  },
  {
    fixtureKey: "fobrainCurrentUser",
    prompt: "查看当前 Fobrain 用户信息",
    label: "Fobrain 当前用户"
  },
  {
    fixtureKey: "fobrainMyPermissions",
    prompt: "查看我的 Fobrain 权限范围",
    label: "Fobrain 我的权限"
  },
  {
    fixtureKey: "fobrainAssetDetail",
    prompt: "查看 Fobrain 资产详情",
    label: "Fobrain 资产详情"
  },
  {
    fixtureKey: "fobrainVulnerabilityDetail",
    prompt: "查看 Fobrain 漏洞详情",
    label: "Fobrain 漏洞详情"
  },
  {
    fixtureKey: "fobrainBusinessRiskSummary",
    prompt: "汇总 Fobrain 业务风险",
    label: "Fobrain 业务风险"
  },
  {
    fixtureKey: "fobrainThreatRelevanceList",
    prompt: "查看 Fobrain 威胁关联资产",
    label: "Fobrain 威胁关联"
  }
] as const;

const forbiddenVisibleText = [
  "authorization",
  "bearer ",
  "api_token",
  "credential_ref",
  "raw provider",
  "raw body",
  "resume_token"
] as const;

test.describe("@visual Fobrain Batch A/E connector and business read evidence", () => {
  for (const scenario of fobrainVisualScenarios) {
    test(`@visual ${scenario.fixtureKey} six evidence regions`, async ({ page }) => {
      await page.goto("/workspaces/ws-demo");
      await selectFobrainFixture(page, scenario.label);
      await expect(page.getByTestId("tool-card")).toBeVisible();
      await expect(page.getByText(scenario.prompt)).toBeVisible();
      await assertNoUnsafeVisibleText(page);

      await captureRegion(page.getByTestId("chat-timeline"), `${scenario.fixtureKey}-main-chat.png`);
      await captureRegion(page.locator(".run-status"), `${scenario.fixtureKey}-process.png`);

      await showInspectorIfNeeded(page);
      await page.getByRole("tab", { name: "证据链" }).click();
      await captureRegion(page.getByTestId("inspector"), `${scenario.fixtureKey}-evidence.png`);

      await page.getByRole("tab", { name: "结构化结果" }).click();
      await captureRegion(page.getByTestId("inspector"), `${scenario.fixtureKey}-internal-details.png`);

      await page.getByRole("tab", { name: "审计" }).click();
      await captureRegion(page.getByTestId("inspector"), `${scenario.fixtureKey}-audit.png`);

      // fresh-main-chat 模拟新项目重新加载后的同一事实投影，不能复用旧项目截图。
      await page.reload();
      await selectFobrainFixture(page, scenario.label);
      await captureRegion(page.getByTestId("chat-timeline"), `${scenario.fixtureKey}-fresh-main-chat.png`);
    });
  }
});

async function selectFobrainFixture(page: Page, label: string) {
  // 使用 fixture 切换器进入场景，不依赖旧项目 DOM 或 provider 私有字段。
  await page.getByRole("button", { name: label }).click();
}

async function showInspectorIfNeeded(page: Page) {
  // 移动端 Inspector 通过面板按钮进入；桌面端该按钮不存在。
  const mobileInspectorButton = page.getByRole("button", { name: "证据" });
  if (await mobileInspectorButton.isVisible()) {
    await mobileInspectorButton.click();
  }
}

async function captureRegion(locator: Locator, name: string) {
  await expect(locator).toBeVisible();
  await expect(locator).toHaveScreenshot(name, {
    animations: "disabled",
    mask: [locator.locator(".conversation-time"), locator.locator(".sidebar-footer")],
    maskColor: "#06141e"
  });
}

async function assertNoUnsafeVisibleText(page: Page) {
  // 安全断言覆盖主聊天和 Inspector 可见文本，避免视觉证据里出现 raw provider 材料。
  const visibleText = (await page.getByTestId("workbench-shell").innerText()).toLowerCase();
  for (const forbidden of forbiddenVisibleText) {
    expect(visibleText.includes(forbidden), `visible text leaked ${forbidden}`).toBe(false);
  }
}
