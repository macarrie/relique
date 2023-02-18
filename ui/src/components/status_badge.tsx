import React, {useEffect, useState} from "react";

import Const from "../types/const";

function StatusBadge(props: any) {
    let [st, changeStatus] = useState(props.status);
    let label = props.label;

    function getColor(): string {
        let color: string;
        switch (st) {
            case Const.OK:
                color = "bg-emerald-100 text-emerald-800 hover:bg-emerald-200 dark:bg-emerald-800 dark:text-emerald-100 dark:hover:bg-emerald-700";
                break;
            case Const.WARNING:
                color = "bg-yellow-100 text-yellow-800 hover:bg-yellow-200 dark:bg-yellow-600 dark:text-yellow-100 dark:hover:bg-yellow-700";
                break;
            case Const.CRITICAL:
                color = "bg-red-100 text-red-800 hover:bg-red-200 dark:bg-red-800 dark:text-red-100 dark:hover:bg-red-700";
                break;
            case Const.INFO:
                color = "bg-blue-100 text-blue-800 hover:bg-blue-200 dark:bg-blue-800 dark:text-blue-100 dark:hover:bg-blue-700";
                break;
            default:
                color = "bg-slate-100 text-slate-800 hover:bg-slate-200 dark:bg-slate-700 dark:text-slate-100 dark:hover:bg-slate-600";
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