/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/views/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: ["bumblebee"]
  },
}

