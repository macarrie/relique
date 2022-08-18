import React from "react";

function Logo() {
    return (
        <span className="text-3xl font-serif font-bold text-yellow-300 hover:text-yellow-500 hover:no-underline">
            <div className="flex flex-row items-center justify-center">
                <i className="ri-trophy-line md:mr-2"></i>
                <span className="hidden md:inline">Relique</span>
            </div>
        </span>
    );
}

export default Logo;
