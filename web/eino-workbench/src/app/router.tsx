import { createBrowserRouter, Navigate } from "react-router-dom";
import { App } from "./App";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Navigate to="/workspaces/ws-demo" replace />
  },
  {
    path: "/workspaces/:workspaceId",
    element: <App />
  }
]);
