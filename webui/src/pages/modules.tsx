import { useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import Module from "../types/module";
import ModuleList from "../components/module_list";
import Const from "../types/const";
import Utils from "../utils/utils";

function Modules() {
    let [mods, setModules] = useState<Module[]>([]);
    let [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getModuleList() {
            setLoading(true);
            API.modules.list({ limit: 10000 }).then((response: any) => {
                setModules(response.data.data ?? []);
                setLoading(false);
            }).catch(error => {
                console.log("Cannot get module list", error);
                Utils.notify(Const.CRITICAL, "Cannot get module list", error.toString())
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
                    loading={loading}
                />
            </Card>
        </>
    );
}

export default Modules;