import { useEffect, useState } from "react";

import Const from "../types/const";

function StatusBadge(props: any) {
    let [st, changeStatus] = useState(props.status);
    let label = props.label;

    function getColor(): string {
        let color: string;
        switch (st) {
            case Const.OK:
                color = "badge-success";
                break;
            case Const.WARNING:
                color = "badge-warning";
                break;
            case Const.CRITICAL:
                color = "badge-error";
                break;
            case Const.INFO:
                color = "badge-info";
                break;
            default:
                color = "badge-neutral";
        }

        return color;
    }

    useEffect(() => {
        changeStatus(props.status);
    }, [props.status]);

    return (
        <span className={`badge badge-sm font-normal ${getColor()}`}>{label}</span>
    );
}

export default StatusBadge;