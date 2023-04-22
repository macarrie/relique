import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";

import API from "../utils/api";
import Module from "../types/module";
import Card from "./card";
import ModuleCard from "./module_card";

function ModuleDetailsPlaceholder() {
    return (
        <Card className="pb-2">
            <div className="flex flex-row px-4 py-3 items-center">
                <span className="flex-grow text-xl font-bold">Module details</span>
            </div>

            <Card className={`bg-white bg-opacity-60 mx-4`}>
                <div className="p-4 flex flex-row items-center mb-2">
                    <div className="flex-grow font-bold text-slate-500 dark:text-slate-200">
                        <div className="rounded-full h-3 w-2/5 bg-slate-300 dark:bg-slate-600"></div>
                    </div>
                </div>
                <table className="details-table ml-4 mb-2">
                    <tr>
                        <td><div className="rounded-full h-2 w-12 bg-slate-300 dark:bg-slate-600"></div></td>
                        <td><div className="rounded-full h-2 w-24 bg-slate-300 dark:bg-slate-600"></div></td>
                    </tr>
                    <tr>
                        <td><div className="rounded-full h-2 w-12 bg-slate-300 dark:bg-slate-600"></div></td>
                        <td><div className="rounded-full h-2 w-24 bg-slate-300 dark:bg-slate-600"></div></td>
                    </tr>
                    <tr>
                        <td><div className="rounded-full h-2 w-12 bg-slate-300 dark:bg-slate-600"></div></td>
                        <td><div className="rounded-full h-2 w-24 bg-slate-300 dark:bg-slate-600"></div></td>
                    </tr>
                </table>
            </Card>
        </Card>
    );
}
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
                return <div className={"text-4xl text-center italic py-8 text-slate-400 dark:text-slate-600"}>Module not found</div>
            }
            return <div className={"text-4xl text-center italic py-8 text-slate-400 dark:text-slate-600"}>Cannot load module</div>
        }
        return <ModuleDetailsPlaceholder />
    }

    return (
        <Card className="pb-2">
            <div className="flex flex-row px-4 py-3 items-center">
                <span className="flex-grow text-xl font-bold">Module details</span>
            </div>

            <ModuleCard className="mx-4" module={m} full/>
        </Card>
    );
}

export default ModuleDetails;