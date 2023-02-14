import React from "react";

function Logo() {
    return (
        <span className="text-3xl font-display font-bold text-yellow-500">
            <div className="flex flex-row items-center justify-center">
                <i className="ri-trophy-line md:mr-2 bg-yellow-500 px-1 rounded text-slate-50 hover:bg-gradient-to-br from-yellow-300 to-yellow-500"></i>
                <span className="hidden md:inline">Relique</span>
            </div>
        </span>
    );
}

export default Logo;
