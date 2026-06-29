/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
  theme: {
    extend: {
      colors: {
        binding: "#1A1D24",
        paper: "#F7F3EE",
        leather: "#8B5E3C",
        "gold-leaf": "#C9A84C",
        "book-cloth": "#DDBE9F",
        ink: "#2E3332",
        "reading-lamp": "#5A6E6A",
      },
      fontFamily: {
        display: ["Cormorant", "serif"],
        body: ["Inter", "sans-serif"],
        mono: ["JetBrains Mono", "monospace"],
      },
      fontSize: {
        "display": ["4rem", { lineHeight: "1.1", fontWeight: "300" }],
        "page-title": ["1.75rem", { lineHeight: "1.2", fontWeight: "600" }],
        "section": ["1.25rem", { lineHeight: "1.3", fontWeight: "600" }],
        "body": ["0.9375rem", { lineHeight: "1.5", fontWeight: "400" }],
        "data": ["0.8125rem", { lineHeight: "1.4", fontWeight: "500" }],
        "small": ["0.75rem", { lineHeight: "1.4", fontWeight: "400" }],
      },
      spacing: {
        "sidebar": "240px",
      },
      maxWidth: {
        "reading-room": "1280px",
      },
      borderRadius: {
        "card": "0.5rem",
      },
      boxShadow: {
        "index-card": "0 1px 3px rgba(0, 0, 0, 0.08)",
      },
    },
  },
  plugins: [],
};
