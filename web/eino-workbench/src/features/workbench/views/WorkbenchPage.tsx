import { useParams } from "react-router-dom";
import { useWorkbenchView } from "../api/workbenchQueries";
import { useWorkbenchUiStore } from "../state/useWorkbenchUiStore";
import { WorkbenchShell } from "../components/WorkbenchShell";

export function WorkbenchPage() {
  const { workspaceId = "ws-demo" } = useParams();
  const activeFixture = useWorkbenchUiStore((state) => state.activeFixture);
  const view = useWorkbenchView(workspaceId, activeFixture);

  if (view.isError) {
    return <div className="app-failure">Workbench fixture 加载失败。</div>;
  }

  if (!view.data) {
    return <div className="app-loading">正在加载 Workbench...</div>;
  }

  return <WorkbenchShell view={view.data} />;
}
