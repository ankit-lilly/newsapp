/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{html,templ}"],
  theme: {
    extend: {
      keyframes: {
        shake: {
          "0%, 100%": { transform: "translateX(0)" },
          "25%": { transform: "translateX(-8px)" },
          "75%": { transform: "translateX(8px)" },
        },
      },
      animation: {
        shake: "shake 0.5s ease-in-out",
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["dracula", "lemonade"],
  },
  darkMode: ["class", '[data-theme="dracula"]'],
};
