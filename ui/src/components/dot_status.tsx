import React, {useEffect, useState} from "react";

import Const from "../types/const";

function DotStatus(props :any) {
    let [st, changeStatus] = useState(props.status);

    function getDotColor() :string {
        let color :string;
        switch (st) {
            case Const.OK:
                color = "bg-green-500";
                break;
            case Const.WARNING:
                color = "bg-orange-500";
                break;
            case Const.CRITICAL:
                color = "bg-red-500";
                break;
            case Const.INFO:
                color = "bg-blue-500";
                break;
            default:
                color = "bg-slate-500";
        }

        return color;
    }

    useEffect(() => {
        changeStatus(props.status);
    }, [props.status]);

    return (
        <div className={`w-3 h-3 ${getDotColor()} rounded-full m-auto`}></div>
    );
}

export default DotStatus;