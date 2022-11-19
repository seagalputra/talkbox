/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./js/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/forms")],
};
