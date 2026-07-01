import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./queryClient";
import { WorkbenchPage } from "../features/workbench/views/WorkbenchPage";

// App 只负责装配全局 QueryClient；具体 Workbench 体验由业务页面承载。
export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <WorkbenchPage />
    </QueryClientProvider>
  );
}
