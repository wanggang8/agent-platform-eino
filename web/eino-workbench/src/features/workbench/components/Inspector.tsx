import * as ScrollArea from "@radix-ui/react-scroll-area";
import * as Tabs from "@radix-ui/react-tabs";
import type { InspectorTab, StructuredResult, WorkbenchInspector } from "../../../contracts/generated";
import { useWorkbenchUiStore } from "../state/useWorkbenchUiStore";
import { StructuredResultView } from "./StructuredResultView";

type InspectorProps = {
  readonly inspector: WorkbenchInspector;
  readonly fallbackStructuredResult: StructuredResult | undefined;
};

export function Inspector({ inspector, fallbackStructuredResult }: InspectorProps) {
  const activeInspectorTab = useWorkbenchUiStore((state) => state.activeInspectorTab);
  const setActiveInspectorTab = useWorkbenchUiStore((state) => state.setActiveInspectorTab);
  const fallbackTabs = ["evidence"] as const satisfies readonly InspectorTab[];
  const tabs: readonly InspectorTab[] = inspector.tabs.length > 0 ? inspector.tabs : fallbackTabs;
  const activeTab: InspectorTab = tabs.includes(activeInspectorTab) ? activeInspectorTab : (tabs[0] ?? "evidence");

  return (
    <aside className="inspector" data-testid="inspector" aria-label="Inspector">
      <Tabs.Root value={activeTab} onValueChange={(value) => setActiveInspectorTab(value as InspectorTab)}>
        <Tabs.List className="inspector-tabs">
          {tabs.map((tab) => <Tabs.Trigger value={tab} key={tab}>{tabLabel(tab)}</Tabs.Trigger>)}
        </Tabs.List>
        <ScrollArea.Root className="inspector-scroll">
          <ScrollArea.Viewport>
            <Tabs.Content value="evidence"><EvidencePanel inspector={inspector} fallbackStructuredResult={fallbackStructuredResult} /></Tabs.Content>
            <Tabs.Content value="structured"><StructuredPanel inspector={inspector} fallbackStructuredResult={fallbackStructuredResult} /></Tabs.Content>
            <Tabs.Content value="runtime"><RuntimePanel inspector={inspector} fallbackStructuredResult={fallbackStructuredResult} /></Tabs.Content>
            <Tabs.Content value="audit"><AuditPanel inspector={inspector} fallbackStructuredResult={fallbackStructuredResult} /></Tabs.Content>
          </ScrollArea.Viewport>
          <ScrollArea.Scrollbar orientation="vertical"><ScrollArea.Thumb /></ScrollArea.Scrollbar>
        </ScrollArea.Root>
      </Tabs.Root>
    </aside>
  );
}

function EvidencePanel({ inspector, fallbackStructuredResult }: InspectorProps) {
  const evidence = inspector.evidence ?? [];
  const summary = getStructuredSummary(fallbackStructuredResult);
  if (evidence.length === 0 && !summary) return <EmptyPanel title="暂无证据" />;
  return (
    <div className="inspector-panel">
      {evidence.map((item) => (
        <div className="evidence-row" key={`${item.label}:${item.value}`}>
          <span>{item.label}</span>
          <strong>{item.value}</strong>
        </div>
      ))}
      {summary ? (
        <div className="detail-card">
          <span>安全结果摘要</span>
          <strong>{summary}</strong>
        </div>
      ) : null}
      {fallbackStructuredResult ? (
        <div className="detail-card">
          <span>事实材料</span>
          <strong>StructuredResult</strong>
        </div>
      ) : null}
    </div>
  );
}

function StructuredPanel({ inspector, fallbackStructuredResult }: InspectorProps) {
  const structured = inspector.structured;
  const result = structured?.structured_result ?? fallbackStructuredResult;
  return (
    <div className="inspector-panel">
      <div className="detail-card">
        <span>结构化结果</span>
        <strong>{structured?.result_ref || fallbackStructuredResult ? "已绑定安全结构化结果" : "等待结构化结果"}</strong>
      </div>
      {result ? <StructuredResultView result={result} /> : <EmptyPanel title="暂无结构化结果" />}
    </div>
  );
}

function RuntimePanel({ inspector }: InspectorProps) {
  const runtime = inspector.runtime;
  return (
    <div className="inspector-panel">
      <div className="detail-card">
        <span>运行状态</span>
        <strong>{runtime?.status ?? "unknown"}</strong>
      </div>
      <div className="detail-card">
        <span>模型上下文</span>
        <strong>safe projection only</strong>
      </div>
      <div className="detail-card">
        <span>事实来源</span>
        <strong>Product Facts</strong>
      </div>
      {runtime?.safe_error ? <div className="safe-error">{runtime.safe_error}</div> : null}
    </div>
  );
}

function AuditPanel({ inspector }: InspectorProps) {
  const audit = inspector.audit ?? [];
  if (audit.length === 0) return <EmptyPanel title="暂无审计事件" />;
  return (
    <div className="inspector-panel">
      {audit.map((event) => (
        <div className="audit-row" key={event.audit_id}>
          <span>{event.event_type}</span>
          <strong>{event.safe_summary}</strong>
          <small>{event.actor}</small>
        </div>
      ))}
    </div>
  );
}

function EmptyPanel({ title }: { readonly title: string }) {
  return <div className="empty-panel">{title}</div>;
}

function tabLabel(tab: InspectorTab) {
  const labels: Record<InspectorTab, string> = {
    evidence: "证据链",
    structured: "结构化结果",
    runtime: "运行详情",
    audit: "审计"
  };
  return labels[tab];
}

function getStructuredSummary(result: StructuredResult | undefined) {
  if (!result) return undefined;
  return result.data.summary ?? result.schema_version;
}
