import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    proxy: {
      "/api": { target: "http://localhost:9001", changeOrigin: true },
    },
    port: 4174,
    strictPort: true,
  },
  preview: {
    port: 4175,
    strictPort: true,
  },
});
