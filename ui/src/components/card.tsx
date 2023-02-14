import React from "react";

function Card(props: any) {
    return (
        <div className={`bg-slate-50 dark:bg-gray-800/50 rounded-lg ${props.className}`}>
            {props.children}
        </div>
    );
}

export default Card;
