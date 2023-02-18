import React from "react";

function Card(props: any) {
    return (
        <div className={`bg-slate-50 dark:bg-slate-800/60 rounded-lg pb-2 ${props.className}`}>
            {props.children}
        </div>
    );
}

export default Card;
