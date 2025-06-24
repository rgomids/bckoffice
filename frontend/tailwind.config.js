const colors = require('tailwindcss/colors');

module.exports = {
  mode: 'jit',
  darkMode: 'class',
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        background: '#0F0F0F',
        foreground: '#E5E7EB',
        card: '#1A1A1A',
        'border-card': '#262626',
        primary: colors.blue,
        secondary: colors.violet,
      },
    },
  },
  plugins: [],
};
