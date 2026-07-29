import { create } from "zustand";

type WorkbenchUiState = {
  readonly activeRightTab: "事实" | "执行记录";
  readonly historyOpen: boolean;
  readonly setActiveRightTab: (tab: "事实" | "执行记录") => void;
  readonly toggleHistory: () => void;
};

// Zustand 只保存 tab/展开态，不复制任何服务端事实。
export const useWorkbenchUiStore = create<WorkbenchUiState>((set) => ({
  activeRightTab: "事实",
  historyOpen: false,
  setActiveRightTab: (activeRightTab) => set({ activeRightTab }),
  toggleHistory: () => set((state) => ({ historyOpen: !state.historyOpen }))
}));
