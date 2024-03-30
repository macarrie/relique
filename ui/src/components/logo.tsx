import React from "react";
import Version from "../version";

function Logo() {
    return (
        <div className="flex flex-col">
            <span className="text-3xl font-display font-bold text-yellow-500">
                <div className="flex flex-row items-center justify-center">
                    <i className="ri-trophy-line mr-2 bg-yellow-500 px-1 rounded text-slate-50 hover:bg-gradient-to-br from-yellow-300 to-yellow-500"></i>
                    <span className="inline">Relique</span>
                </div>
            </span>
            <span className="inline text-right text-xs text-slate-400 dark:text-slate-600 italic">v{Version.Tag}</span>
        </div>
    );
}

export default Logo;
