/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{html,templ}"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["dracula", "lemonade"],
  },
  darkMode: ["class", '[data-theme="dracula"]'],
};
