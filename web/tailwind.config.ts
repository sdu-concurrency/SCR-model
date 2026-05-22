// @ts-expect-error no declared types at this time
import primeuiPlugin from 'tailwindcss-primeui'

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx,json}',
    './formkit.theme.ts' // <-- add your theme file
  ],
  darkMode: 'app-dark',
  plugins: [primeuiPlugin],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter var']
      },
      // primevue colors
      colors: {
        'light-gray': '#ededed',
        'medium-gray': '#706f6f',
        'dark-gray': '#383838',
        'project-red': '#8e1315'
      }
    }
  }
}
