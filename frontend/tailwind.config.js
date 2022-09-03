/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx}",
    "./components/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
    fontSize: {
      sm: ['10px', '10px'],
      base: ['16px', '24px'],
      lg: ['20px', '28px'],
      xl: ['30px', '38px'],
    }
  },
  plugins: [require('daisyui')],
}
