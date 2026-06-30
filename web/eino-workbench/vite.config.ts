import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react()],
  server: {
    host: "127.0.0.1",
    port: 8081,
    strictPort: true,
    fs: {
      allow: ["../.."]
    }
  },
  preview: {
    host: "127.0.0.1",
    port: 8081,
    strictPort: true
  },
  test: {
    environment: "jsdom",
    globals: true,
    include: ["src/**/*.test.ts", "src/**/*.test.tsx"],
    exclude: ["tests/**", "node_modules/**"],
    setupFiles: ["./src/test/setup.ts"]
  }
});
