import React from "react";

function PaginationButton({
                              children,
                              onClick,
                              disabled,
                              active,
               } :any) {
    return (
        <button className={`w-8 h-8 rounded-full hover:bg-slate-200 hover:dark:bg-slate-700 dark:disabled:text-slate-700 disabled:text-slate-300 ${active ? "!bg-yellow-500/20" : ""}`}
                onClick={onClick}
                disabled={disabled}>
            {children}
        </button>
    );
}
export default PaginationButton;
