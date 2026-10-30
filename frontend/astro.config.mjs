import { defineConfig } from "astro/config";
import react from "@astrojs/react";

// Static output only — this build is served by the Go binary as plain
// files, with no Node process at request time (design.md "Astro
// builds shells, not content"). The dev-server proxy (task 1.3) sends
// API/auth requests to `breaklab serve` running separately.
// Dev-only: the Go backend (`breaklab serve`) that owns these paths.
const goBackend = "http://127.0.0.1:8787";

export default defineConfig({
  output: "static",
  integrations: [react()],
  vite: {
    server: {
      proxy: {
        "/api": goBackend,
        "/sync": goBackend,
        "/signup": goBackend,
        "/login": goBackend,
        "/logout": goBackend,
        "/account": goBackend,
        "/u": goBackend,
      },
    },
  },
});
