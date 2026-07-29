import { useParams } from "react-router-dom";
import { useSubmitWorkbenchMessage, useWorkbenchView } from "../api/workbenchQueries";
import { WorkbenchShell } from "../components/WorkbenchShell";

export function WorkbenchPage() {
  const { workspaceId = "ws-workbench" } = useParams();
  const view = useWorkbenchView(workspaceId);
  const conversationID = view.data?.conversation_id ?? "conversation_initial";
  const submit = useSubmitWorkbenchMessage(workspaceId, conversationID);
  if (view.isError) return <div className="app-failure" role="alert">无法加载工作台。</div>;
  if (!view.data) return <div className="app-loading">正在加载工作台…</div>;
  return <WorkbenchShell view={view.data} isSubmitting={submit.isPending} onSubmit={(content) => submit.mutateAsync(content).then(() => undefined)} />;
}
