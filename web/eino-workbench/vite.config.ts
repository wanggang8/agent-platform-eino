import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react()],
  server: {
    // dev server 固定本机地址和端口，方便脚本、Playwright 和文档使用同一入口。
    host: "127.0.0.1",
    port: 8081,
    strictPort: true,
    fs: {
      // 允许读取仓库根目录下的契约和 fixture，不扩大到用户目录。
      allow: ["../.."]
    }
  },
  preview: {
    host: "127.0.0.1",
    port: 8081,
    strictPort: true
  },
  test: {
    // 组件测试运行在 jsdom，Playwright 只负责端到端与视觉验收。
    environment: "jsdom",
    globals: true,
    include: ["src/**/*.test.ts", "src/**/*.test.tsx"],
    exclude: ["tests/**", "node_modules/**"],
    setupFiles: ["./src/test/setup.ts"]
  }
});
