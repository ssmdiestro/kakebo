import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const target = env.VITE_API_TARGET || "http://app:8080";

  console.log(`🟢 [Vite proxy] apuntando a: ${target}`);

  return {
    base: "/",
    plugins: [react()],
    server: {
      host: true,
      port: 5173,
      proxy: {
        "/api": {
          target,
          changeOrigin: true,
          rewrite: (path) => {
            const rewritten = path.replace(/^\/api/, "");
            console.log(`🌀 Proxy reescribe: ${path} → ${rewritten}`);
            return rewritten;
          },
          configure: (proxy, options) => {
            proxy.on("proxyReq", (proxyReq, req, res) => {
              console.log(
                `➡️  Proxy request: ${req.method} ${req.url} → ${options.target}${req.url.replace(/^\/api/, "")}`
              );
            });
            proxy.on("proxyRes", (proxyRes, req, res) => {
              console.log(`⬅️  Proxy response: ${proxyRes.statusCode} para ${req.url}`);
            });
            proxy.on("error", (err, req, res) => {
              console.error(`💥 Proxy error en ${req.url}:`, err.message);
            });
          },
        },
      },
    },
  };
});
