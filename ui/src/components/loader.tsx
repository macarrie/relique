import React, {useEffect, useState} from "react";

function Loader(props :any) {
    let [label, setLabel] = useState(props.label)

    useEffect(() => {
        setLabel(props.label);
    }, [props.label])

    return (
        <div className="flex flex-row justify-center">
            <i className="animate-spin ri-loader-4-line"></i>
            {label !== "" && (
                <span className="ml-2">{label}</span>
            )}
        </div>
    );
}

export default Loader;
