import defaultTheme from "tailwindcss/defaultTheme";

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{html,templ,tmpl}"],
  theme: {
    extend: {
      screens: {
        print: { raw: "print" },
        mobile: { max: "475px" },
        tablet: { max: "600px" },
        ...defaultTheme.screens,
      },
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
