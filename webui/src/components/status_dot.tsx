import { useEffect, useState } from "react";

import Const from "../types/const";

function StatusDot(props: any) {
    let [st, changeStatus] = useState(props.status);

    function getColor(): string {
        let color: string;
        switch (st) {
            case Const.OK:
                color = "bg-success";
                break;
            case Const.WARNING:
                color = "bg-warning";
                break;
            case Const.CRITICAL:
                color = "bg-error";
                break;
            case Const.INFO:
                color = "bg-info";
                break;
            default:
                color = "bg-neutral/30";
        }

        return color;
    }

    useEffect(() => {
        changeStatus(props.status);
    }, [props.status]);

    return (
        <div className={`w-3 h-3 rounded-full m-auto ${getColor()}`}></div>
    );
}

export default StatusDot;