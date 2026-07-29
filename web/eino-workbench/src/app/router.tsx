import { createBrowserRouter, Navigate } from "react-router-dom";
import { App } from "./App";

// router 固定 Workbench 为第一屏，不提供营销落地页。
export const router = createBrowserRouter([
  {
    path: "/",
    element: <Navigate to="/workspaces/ws-workbench" replace />
  },
  {
    path: "/workspaces/:workspaceId",
    element: <App />
  }
]);
