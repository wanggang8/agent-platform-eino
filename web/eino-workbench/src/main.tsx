import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { router } from "./app/router";
import "./styles/tokens.css";
import "./styles/base.css";
import "./styles/workbench.css";

const root = document.getElementById("root");

if (!root) {
  throw new Error("root element missing");
}

// 应用入口只挂载路由，业务状态由各 feature 自己管理。
createRoot(root).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>
);
