import { useEffect, useState } from "react";

import PageHeader from "../components/page_header";
import Section from "../components/section";
import API from "../utils/api";
import Module from "../types/module";
import ModuleList from "../components/module_list";
import Const from "../types/const";
import Utils from "../utils/utils";

function Modules() {
    const [mods, setModules] = useState<Module[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getModuleList() {
            setLoading(true);
            API.modules.list({ limit: 10000 }).then((response) => {
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
            <PageHeader title="Modules" />
            <Section>
                <ModuleList
                    actions={true}
                    data={mods}
                    paginated={true}
                    sorted={true}
                    loading={loading}
                />
            </Section>
        </>
    );
}

export default Modules;