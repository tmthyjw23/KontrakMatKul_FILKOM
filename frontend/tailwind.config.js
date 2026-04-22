/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./lib/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "#050505",
        foreground: "#f4f4f5",
        muted: "#a1a1aa",
        border: "rgba(255, 255, 255, 0.1)",
        panel: "rgba(255, 255, 255, 0.05)",
        "panel-strong": "rgba(255, 255, 255, 0.08)",
      },
      boxShadow: {
        glass: "0 24px 80px rgba(0, 0, 0, 0.45)",
      },
      borderRadius: {
        "4xl": "2rem",
      },
      backgroundImage: {
        "grid-fade":
          "linear-gradient(to right, rgba(255,255,255,0.04) 1px, transparent 1px), linear-gradient(to bottom, rgba(255,255,255,0.04) 1px, transparent 1px)",
      },
    },
  },
  plugins: [],
};
