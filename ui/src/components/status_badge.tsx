import React, {useEffect, useState} from "react";

import Const from "../types/const";

function StatusBadge(props: any) {
    let [st, changeStatus] = useState(props.status);
    let label = props.label;

    function getColor(): string {
        let color: string;
        switch (st) {
            case Const.OK:
                color = "bg-emerald-100 text-emerald-800 hover:bg-emerald-200 dark:bg-emerald-500/70 dark:text-emerald-200 dark:hover:bg-emerald-500/90";
                break;
            case Const.WARNING:
                color = "bg-yellow-100 text-yellow-800 hover:bg-yellow-300 dark:bg-yellow-500/70 dark:text-yellow-100 dark:hover:bg-yellow-500/90";
                break;
            case Const.CRITICAL:
                color = "bg-red-100 text-red-800 hover:bg-red-200 dark:bg-red-500/70 dark:text-red-100 dark:hover:bg-red-500/90";
                break;
            case Const.INFO:
                color = "bg-blue-100 text-blue-800 hover:bg-blue-200 dark:bg-blue-500/70 dark:text-blue-100 dark:hover:bg-blue-500/90";
                break;
            default:
                color = "bg-slate-100 text-slate-800 hover:bg-slate-200 dark:bg-slate-500/70 dark:text-slate-100 dark:hover:bg-slate-500/90";
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