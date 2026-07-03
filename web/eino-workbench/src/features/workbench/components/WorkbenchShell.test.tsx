import { cleanup, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import type { StructuredResult } from "../../../contracts/generated";
import { fobrainVisualFixtures } from "../../../fixtures/fobrainVisualFixtures";
import { workbenchFixtures } from "../../../fixtures/workbenchFixtures";
import { renderWithClient } from "../../../test/render";
import { StructuredResultView } from "./StructuredResultView";
import { WorkbenchShell } from "./WorkbenchShell";

describe("WorkbenchShell", () => {
  it("renders the fixture-driven workbench without exposing raw tool ids", () => {
    // shell 测试以 fixture 契约为输入，确保 UI 不直接暴露 provider/tool 内部 id。
    renderWithClient(<WorkbenchShell view={workbenchFixtures.success} />);

    expect(screen.getByTestId("workbench-shell")).toBeInTheDocument();
    expect(screen.getByTestId("workspace-sidebar")).toBeInTheDocument();
    expect(screen.getByTestId("chat-timeline")).toBeInTheDocument();
    expect(screen.getByTestId("inspector")).toBeInTheDocument();
    expect(screen.queryByText(/tool\.fobrain/)).not.toBeInTheDocument();
  });

  it("renders Fobrain visual fixtures as Chinese product UI without debug tokens", () => {
    // 产品视觉验收看可见文本，不允许把调试态英文、tool id 或 schema 细节作为产品内容露出。
    for (const view of Object.values(fobrainVisualFixtures)) {
      cleanup();
      renderWithClient(<WorkbenchShell view={view} />);
      const visibleText = productText(screen.getByTestId("workbench-shell"));

      expect(visibleText).not.toMatch(/[A-Za-z]/);
      expect(visibleText).not.toMatch(/tool\.|schema_version|display_type|StructuredResult|Product Facts|safe projection|fixture|run-/i);
      expect(visibleText).not.toMatch(/[{}]/);
    }
  });

  it("does not stringify nested objects into visible JSON", () => {
    // 嵌套对象值只能显示为产品化摘要，不能把 JSON 结构直接给用户。
    const result: StructuredResult = {
      schema_version: "tool.structured_result.v1",
      status: "resolved",
      data: {
        summary: "对象值展示验证",
        facts: [{ key: "nested", label: "嵌套对象", value: { raw: "value" } }]
      },
      metadata: { safe: true }
    };

    renderWithClient(<StructuredResultView result={result} />);

    expect(screen.getByText("已整理")).toBeInTheDocument();
    expect(productText(screen.getByTestId("structured-result-view"))).not.toMatch(/[{}]|raw|value/);
  });

  it("does not expose unknown tool status codes", () => {
    // 未知状态只能显示产品化兜底，不能把 provider_timeout 等内部枚举透出。
    const view = {
      ...fobrainVisualFixtures.fobrainAssetDetail,
      timeline: fobrainVisualFixtures.fobrainAssetDetail.timeline.map((item) =>
        item.kind === "tool_card" ? { ...item, status: "provider_unknown_state" } : item
      )
    };

    renderWithClient(<WorkbenchShell view={view} />);

    const visibleText = productText(screen.getByTestId("workbench-shell"));
    expect(visibleText).toContain("未知");
    expect(visibleText).not.toContain("provider_unknown_state");
  });

  it("renders terminal pending cards as readonly product records", () => {
    // 终态 pending 只能作为审计/回放记录展示，不能继续提供恢复动作。
    const view = {
      ...workbenchFixtures.approval,
      timeline: workbenchFixtures.approval.timeline.map((item) =>
        item.kind === "approval_card" ? { ...item, status: "cancelled" } : item
      )
    };

    renderWithClient(<WorkbenchShell view={view} />);

    expect(screen.getByText("已取消")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "批准并提交" })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "拒绝" })).not.toBeInTheDocument();
  });
});

function productText(element: HTMLElement) {
  // 测试只检查产品可见文本；Radix 在 jsdom 中注入的 style 文本不是用户可见内容。
  const clone = element.cloneNode(true) as HTMLElement;
  clone.querySelectorAll("style,script").forEach((node) => node.remove());
  return clone.textContent ?? "";
}
