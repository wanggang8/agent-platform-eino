import { Activity, Boxes, CircleCheck, Gauge, Share2, SlidersHorizontal } from "lucide-react";
import type { StructuredResult, WorkbenchView } from "../../../contracts/generated";
import { fixtureLabels, type WorkbenchFixtureKey } from "../../../fixtures/workbenchFixtures";
import { useWorkbenchUiStore } from "../state/useWorkbenchUiStore";
import { Composer } from "./Composer";
import { Inspector } from "./Inspector";
import { Timeline } from "./Timeline";
import { WorkspaceSidebar } from "./WorkspaceSidebar";

const fixtureKeys = Object.keys(fixtureLabels) as WorkbenchFixtureKey[];

type WorkbenchShellProps = {
  readonly view: WorkbenchView;
};

// WorkbenchShell 负责三栏布局和响应式面板切换，业务事实只来自 WorkbenchView。
export function WorkbenchShell({ view }: WorkbenchShellProps) {
  const activeFixture = useWorkbenchUiStore((state) => state.activeFixture);
  const setActiveFixture = useWorkbenchUiStore((state) => state.setActiveFixture);
  const mobilePanel = useWorkbenchUiStore((state) => state.mobilePanel);
  const setMobilePanel = useWorkbenchUiStore((state) => state.setMobilePanel);

  const latestStructuredResult = findLatestStructuredResult(view);

  return (
    <div className="workbench-shell" data-testid="workbench-shell">
      <WorkspaceSidebar activeRunId={view.run_id} />

      <main className="workbench-main" data-mobile-panel={mobilePanel}>
        <header className="workbench-header">
          <div>
            <p className="eyeline">工作区：默认工作区 &gt; 会话</p>
            <h1>查询 10.10.11.69 关联资产与风险</h1>
          </div>
          <div className="header-actions" aria-label="运行环境状态">
            <span><Gauge size={14} /> gpt-5.4-mini</span>
            <span><CircleCheck size={14} /> 连接正常</span>
            <span><Boxes size={14} /> 6/6 可用</span>
            <button type="button" aria-label="分享"><Share2 size={16} /></button>
          </div>
        </header>

        <div className="fixture-switcher" aria-label="Fixture state selector">
          {fixtureKeys.map((key) => (
            <button
              key={key}
              type="button"
              className={activeFixture === key ? "is-active" : ""}
              onClick={() => setActiveFixture(key)}
            >
              {fixtureLabels[key]}
            </button>
          ))}
        </div>

        <div className="mobile-panel-tabs" aria-label="移动端面板切换">
          <button type="button" className={mobilePanel === "chat" ? "is-active" : ""} onClick={() => setMobilePanel("chat")}>
            对话
          </button>
          <button type="button" className={mobilePanel === "inspector" ? "is-active" : ""} onClick={() => setMobilePanel("inspector")}>
            证据
          </button>
        </div>

        <section className="workbench-content">
          <div className="conversation-column">
            <div className="run-status">
              <Activity size={16} />
              <span>运行状态</span>
              <strong>{statusLabel(view.status)}</strong>
            </div>
            <Timeline items={view.timeline} />
            <Composer disabled={view.status === "waiting"} />
          </div>
          <Inspector inspector={view.inspector} fallbackStructuredResult={latestStructuredResult} />
        </section>
      </main>

      <button className="floating-settings" type="button" aria-label="设置与能力">
        <SlidersHorizontal size={18} />
      </button>
    </div>
  );
}

// findLatestStructuredResult 为 Inspector 提供最近的结构化结果兜底展示。
function findLatestStructuredResult(view: WorkbenchView): StructuredResult | undefined {
  for (let index = view.timeline.length - 1; index >= 0; index -= 1) {
    const item = view.timeline[index];
    if (item?.structured_result) return item.structured_result;
  }
  return undefined;
}

// statusLabel 将 run status 映射为中文展示。
function statusLabel(status: WorkbenchView["status"]) {
  const labels: Record<WorkbenchView["status"], string> = {
    created: "已创建",
    running: "运行中",
    waiting: "等待交互",
    succeeded: "已完成",
    failed: "失败",
    cancelled: "已取消",
    stopped: "已停止"
  };
  return labels[status];
}
