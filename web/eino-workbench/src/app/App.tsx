import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./queryClient";
import { WorkbenchPage } from "../features/workbench/views/WorkbenchPage";

export function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <WorkbenchPage />
    </QueryClientProvider>
  );
}
