import { create } from "zustand";
import type { InspectorTab } from "../../../contracts/generated";
import type { WorkbenchFixtureKey } from "../../../fixtures/workbenchFixtures";

type WorkbenchUiState = {
  readonly activeFixture: WorkbenchFixtureKey;
  readonly activeInspectorTab: InspectorTab;
  readonly collapsedToolCards: ReadonlySet<string>;
  readonly mobilePanel: "chat" | "inspector";
  readonly setActiveFixture: (fixture: WorkbenchFixtureKey) => void;
  readonly setActiveInspectorTab: (tab: InspectorTab) => void;
  readonly toggleToolCard: (itemId: string) => void;
  readonly setMobilePanel: (panel: "chat" | "inspector") => void;
};

export const useWorkbenchUiStore = create<WorkbenchUiState>((set) => ({
  activeFixture: "success",
  activeInspectorTab: "evidence",
  collapsedToolCards: new Set(),
  mobilePanel: "chat",
  setActiveFixture: (fixture) => set({ activeFixture: fixture }),
  setActiveInspectorTab: (tab) => set({ activeInspectorTab: tab }),
  toggleToolCard: (itemId) =>
    set((state) => {
      const next = new Set(state.collapsedToolCards);
      if (next.has(itemId)) {
        next.delete(itemId);
      } else {
        next.add(itemId);
      }
      return { collapsedToolCards: next };
    }),
  setMobilePanel: (panel) => set({ mobilePanel: panel })
}));
