import React, {useEffect, useState} from "react";

import Const from "../types/const";

function StatusBadge(props: any) {
    let [st, changeStatus] = useState(props.status);
    let label = props.label;

    function getColor(): string {
        let color: string;
        switch (st) {
            case Const.OK:
                color = "bg-emerald-100 text-emerald-800 hover:bg-emerald-200";
                break;
            case Const.WARNING:
                color = "bg-yellow-100 text-yellow-800 hover:bg-yellow-300";
                break;
            case Const.CRITICAL:
                color = "bg-red-100 text-red-800 hover:bg-red-200";
                break;
            case Const.INFO:
                color = "bg-blue-100 text-blue-800 hover:bg-blue-200";
                break;
            default:
                color = "bg-slate-200 text-slate-800 hover:bg-slate-300";
        }

        return color;
    }

    useEffect(() => {
        changeStatus(props.status);
    }, [props.status]);

    return (
        <div className={`badge ${getColor()}`}>{label}</div>
    );
}

export default StatusBadge;