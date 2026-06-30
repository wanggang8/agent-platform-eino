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

createRoot(root).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>
);
