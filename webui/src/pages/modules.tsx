import { useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import Module from "../types/module";
import ModuleList from "../components/module_list";

function Modules() {
    let [mods, setModules] = useState<Module[]>([]);

    useEffect(() => {
        function getModuleList() {
            API.modules.list({ limit: 10000 }).then((response: any) => {
                setModules(response.data.data ?? []);
            }).catch(error => {
                console.log("Cannot get job list", error);
                setModules([]);
            });
        }

        getModuleList();
    }, [])

    return (
        <>
            <Card>
                <ModuleList
                    title="All modules"
                    actions={true}
                    data={mods}
                    paginated={true}
                    sorted={true}
                />
            </Card>
        </>
    );
}

export default Modules;