import React, {useEffect, useState} from "react";
import { useParams } from "react-router-dom";

import API from "../utils/api";
import Module from "../types/module";

function ModuleDetails(props :any) {
    let urlParams = useParams();
    let [m, setModule] = useState<Module | null>(null);

    useEffect(() => {
        function getModule() {
            let name = urlParams["name"];
            if (name === undefined) {
                console.log("Module name undefined, cannot get details");
                return;
            }

            API.modules.get(name).then((response :any) => {
                setModule(response.data);
            }).catch(error => {
                console.log("Cannot get module details", error);
            });
        }

        getModule();
    }, [urlParams])

    if (m === null) {
        return <div>Loading</div>
    }

    return (
        <div>{JSON.stringify(m)}</div>
    );
}

export default ModuleDetails;