import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import type { M1WorkbenchStreamEvent, M1WorkbenchView as WorkbenchView } from "../../../contracts/generated";

const viewKey = (workspaceId: string) => ["m1-workbench-view", workspaceId] as const;

async function readJSON(response: Response): Promise<WorkbenchView> {
  if (!response.ok) throw new Error("workbench request failed");
  return response.json() as Promise<WorkbenchView>;
}

// useWorkbenchView 只读取 Go API 的安全投影，不导入本地业务 fixture。
export function useWorkbenchView(workspaceId: string) {
  return useQuery({
    queryKey: viewKey(workspaceId),
    queryFn: () => fetch(`/api/workspaces/${encodeURIComponent(workspaceId)}/views/current`).then(readJSON)
  });
}

export function useSubmitWorkbenchMessage(workspaceId: string, conversationId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (content: string) => {
      const accepted = await fetch(`/api/workspaces/${encodeURIComponent(workspaceId)}/messages`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ conversation_id: conversationId, actor_id: "operator", content })
      }).then(readJSON);
      const streamResponse = await fetch(`/api/workspaces/${encodeURIComponent(workspaceId)}/runs/${encodeURIComponent(accepted.run_id)}/stream`);
      if (!streamResponse.ok) throw new Error("workbench stream failed");
      const event = parseLastSSEView(await streamResponse.text());
      return event?.view ?? accepted;
    },
    onSuccess: (view) => queryClient.setQueryData(viewKey(workspaceId), view)
  });
}

function parseLastSSEView(payload: string): M1WorkbenchStreamEvent | undefined {
  const dataLines = payload.split("\n").filter((line) => line.startsWith("data: "));
  const last = dataLines.at(-1);
  return last ? JSON.parse(last.slice(6)) as M1WorkbenchStreamEvent : undefined;
}
