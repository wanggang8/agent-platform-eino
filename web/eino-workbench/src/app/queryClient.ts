import { QueryClient } from "@tanstack/react-query";

// queryClient 只管理 server state；UI/client state 放在 Zustand store。
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: false
    },
    mutations: {
      retry: false
    }
  }
});
