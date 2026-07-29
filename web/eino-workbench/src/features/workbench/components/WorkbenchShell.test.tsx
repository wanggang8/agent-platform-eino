import { fireEvent, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { M1WorkbenchView as WorkbenchView } from "../../../contracts/generated";
import { renderWithClient } from "../../../test/render";
import { WorkbenchShell } from "./WorkbenchShell";

const resolvedView: WorkbenchView = {
  schema_version: "eino_workbench_view.v2", workspace_id: "ws-1", conversation_id: "conversation-1", run_id: "run-1", status: "resolved",
  messages: [
    { message_id: "m1", role: "user", content: "查询新增的漏洞" },
    { message_id: "m2", role: "assistant", content: "发现 3 条新增漏洞" }
  ],
  result: { query_sequence: 1, summary: "发现 3 条新增漏洞", count: 3, observed_at: "2026-07-27T08:35:00Z" },
  right_panel: { tabs: ["事实", "执行记录"], empty_message: "选择查询结果后查看安全事实。" }
};

describe("M1 WorkbenchShell", () => {
  it("按 nav/main/aside 顺序渲染唯一浅色桌面三栏", () => {
    renderWithClient(<WorkbenchShell view={resolvedView} />);
    const shell = screen.getByTestId("workbench-shell");
    const regions = Array.from(shell.children).filter((node) => ["NAV", "MAIN", "ASIDE"].includes(node.tagName));
    expect(regions.map((node) => node.tagName)).toEqual(["NAV", "MAIN", "ASIDE"]);
    expect(screen.getByTestId("query-result-card")).toHaveTextContent("发现 3 条新增漏洞");
    expect(shell.textContent).not.toMatch(/result_ref|snapshot_item_ref|StructuredResult|审计|审批|暗色|移动端/i);
  });

  it("操作记录只打开本地空态", () => {
    renderWithClient(<WorkbenchShell view={resolvedView} />);
    fireEvent.click(screen.getByRole("button", { name: "操作记录" }));
    expect(screen.getByTestId("history-empty")).toHaveTextContent("操作记录尚未启用");
  });

  it("Composer 提交用户输入且不提供写域控件", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);
    renderWithClient(<WorkbenchShell view={resolvedView} onSubmit={onSubmit} />);
    fireEvent.change(screen.getByLabelText("输入查询"), { target: { value: "查询新增的漏洞" } });
    fireEvent.click(screen.getByRole("button", { name: "发送" }));
    await waitFor(() => expect(onSubmit).toHaveBeenCalledWith("查询新增的漏洞"));
    expect(screen.queryByRole("button", { name: /批准|派发|误报|延时|人员/ })).not.toBeInTheDocument();
  });

  it("失败态不生成伪空结果卡", () => {
    renderWithClient(<WorkbenchShell view={{ ...resolvedView, status: "failed", result: null, safe_error: "无法读取漏洞事实。此次请求不是空结果。" }} />);
    expect(screen.queryByTestId("query-result-card")).not.toBeInTheDocument();
    expect(screen.getByRole("alert")).toHaveTextContent("此次请求不是空结果");
  });
});
