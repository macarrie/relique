import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";

import API from "../utils/api";
import Module from "../types/module";
import Card from "./card";
import ModuleCard from "./module_card";

function ModuleDetails() {
    let urlParams = useParams();
    let [m, setModule] = useState<Module | null>(null);
    let [notFound, setNotFound] = useState<boolean>(false);
    let [err, setErr] = useState<boolean>(false);

    useEffect(() => {
        function getModule() {
            let name = urlParams["name"];
            if (name === undefined) {
                console.log("Module name undefined, cannot get details");
                return;
            }

            API.modules.get(name).then((response: any) => {
                setModule(response.data);
            }).catch(error => {
                setErr(true)
                if (error.response.status === 404) {
                    setNotFound(true)
                }
                console.log("Cannot get module details", error);
            });
        }

        getModule();
    }, [urlParams])

    if (m === null) {
        if (err) {
            if (notFound) {
                return <div>Client not found</div>
            }
            return <div>Cannot load client</div>
        }
        return <div>Loading</div>
    }

    return (
        <Card>
            <div className="flex flex-row px-4 py-3 items-center">
                <span className="flex-grow text-xl font-bold">Module details</span>
            </div>

            <ModuleCard className="bg-transparent" module={m} full/>
        </Card>
    );
}

export default ModuleDetails;