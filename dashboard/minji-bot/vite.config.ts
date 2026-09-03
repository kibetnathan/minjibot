import path from "path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    // Dev proxy: forward /api requests (and the OAuth redirect/callback) to
    // the Echo API on :8080 so the dashboard and API share the :5173 origin.
    // With the proxy, the Discord OAuth flow (login -> callback -> dashboard)
    // stays on the same origin and the session cookie is set correctly.
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
})
