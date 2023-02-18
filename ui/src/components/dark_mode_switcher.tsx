import React, {useEffect, useState} from "react";

function DarkModeSwitcher() {
    let [theme, setTheme] = useState<String>(getMode());

    useEffect(() => {
        localStorage.theme = theme;
        setMode(theme as string);
    }, [theme])

    // On page load or when changing themes, best to add inline in `head` to avoid FOUC

    // Whenever the user explicitly chooses light mode
    // localStorage.theme = 'light'

    // Whenever the user explicitly chooses dark mode
    // localStorage.theme = 'dark'

    // Whenever the user explicitly chooses to respect the OS preference
    // localStorage.removeItem('theme')
    function getMode() {
        if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
            return "dark";
        } else {
            return "light";
        }
    }

    function setMode(mode: string) {
        if (mode === "dark") {
            document.documentElement.classList.add('dark')
            return;
        }

        document.documentElement.classList.remove('dark')
    }

    return (
        <div className="text-lg">
            {theme === "dark" ? (
                <button onClick={() => setTheme("light")}><i className="ri-sun-line"></i></button>
            ) : (
                <button onClick={() => setTheme("dark")}><i className="ri-moon-line"></i></button>
            )}
        </div>
    );
}

export default DarkModeSwitcher;
