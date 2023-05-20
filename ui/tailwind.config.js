/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./src/**/*.{js,jsx,ts,tsx}",
    ],
    theme: {
        extend: {},
        fontFamily: {
            'sans': ['Inter', 'ui-sans-serif', 'system-ui'],
            'serif': ['ui-serif', 'Georgia'],
            'mono': ['Inconsolata', 'ui-monospace', 'SFMono-Regular'],
            'display': ['Inter', 'ui-sans-serif'],
            'body': ['Inter', 'Roboto'],
        }
    },
    plugins: [
        require('@tailwindcss/forms'),
    ],
    darkMode: 'class',
}
