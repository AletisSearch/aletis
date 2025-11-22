import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [tailwindcss()],
  build: {
    // generate .vite/manifest.json in outDir
    manifest: true,
    rollupOptions: {
      // overwrite default .html entry
      input: ["./main.css"],
      output: {
        dir: "./dist",
      },
    },
    modulePreload: {
      polyfill: false,
    },
  },
});
