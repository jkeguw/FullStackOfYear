/** @type {import('tailwindcss').Config} */
export default {
    content: ['./src/**/*.{vue,ts,tsx}', './index.html'],
    theme: {
        extend: {
            colors: {
                primary: '#1E1E1E',
                secondary: '#333333',
                accent: '#FFFFFF',
                dark: {
                  100: '#4A4A4A',
                  200: '#3A3A3A',
                  300: '#2A2A2A',
                  400: '#1A1A1A',
                  500: '#0A0A0A',
                }
            },
            screens: {
                'xs': '380px',
                'sm': '640px',
                'md': '768px',
                'lg': '1024px',
                'xl': '1280px',
                '2xl': '1536px',
            },
            container: {
                padding: {
                    DEFAULT: '1rem',
                    sm: '1rem',
                    md: '1.5rem',
                    lg: '2rem',
                },
                center: true
            },
        },
    },
    plugins: [],
}