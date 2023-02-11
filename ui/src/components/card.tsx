import React from "react";

function Card(props: any) {
    return (
        <div className={`bg-slate-50 rounded-lg ${props.className}`}>
            {props.children}
        </div>
    );
}

export default Card;
