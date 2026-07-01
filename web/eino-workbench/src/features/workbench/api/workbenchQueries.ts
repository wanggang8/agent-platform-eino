import { useQuery } from "@tanstack/react-query";
import type { WorkbenchView } from "../../../contracts/generated";
import { workbenchFixtures, type WorkbenchFixtureKey } from "../../../fixtures/workbenchFixtures";

// useWorkbenchView 当前只读取 fixture 投影；真实 API 接入后仍返回同一 contract 类型。
export function useWorkbenchView(workspaceId: string, fixtureKey: WorkbenchFixtureKey) {
  return useQuery<WorkbenchView>({
    queryKey: ["workbench-view", workspaceId, fixtureKey],
    queryFn: async () => workbenchFixtures[fixtureKey]
  });
}
