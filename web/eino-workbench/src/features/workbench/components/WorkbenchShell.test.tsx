import { screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { workbenchFixtures } from "../../../fixtures/workbenchFixtures";
import { renderWithClient } from "../../../test/render";
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
});
