/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: '#0f172a',
        secondary: '#1e293b',
        card: '#1e293b',
        border: '#334155',
        accent: '#0ea5e9',
        'accent-dark': '#0284c7',
        success: '#10b981',
        warning: '#f59e0b',
        danger: '#ef4444',
        info: '#0ea5e9',
        purple: '#8b5cf6',
      }
    },
  },
  plugins: [],
}

