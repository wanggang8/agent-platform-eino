import { useQuery } from "@tanstack/react-query";
import type { WorkbenchView } from "../../../contracts/generated";
import { workbenchFixtures, type WorkbenchFixtureKey } from "../../../fixtures/workbenchFixtures";

export function useWorkbenchView(workspaceId: string, fixtureKey: WorkbenchFixtureKey) {
  return useQuery<WorkbenchView>({
    queryKey: ["workbench-view", workspaceId, fixtureKey],
    queryFn: async () => workbenchFixtures[fixtureKey]
  });
}
