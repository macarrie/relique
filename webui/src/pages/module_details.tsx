import { useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import API from '../utils/api';
import Module from '../types/module';
import Const from '../types/const';
import Utils from '../utils/utils';
import PageHeader from '../components/page_header';
import Section from '../components/section';
import ModuleProperties from '../components/module_properties';

function ModuleDetails() {
    const { module_name } = useParams();
    const [mod, setModule] = useState<Module>({} as Module);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getModule() {
            setLoading(true);
            if (module_name === undefined) {
                console.log("Module name undefined, cannot get module details");
                return;
            }

            API.modules.get(module_name).then((response) => {
                setModule(response.data);
                setLoading(false);
            }).catch(error => {
                console.log("Cannot get module details", error);
                Utils.notify(Const.CRITICAL, "Cannot get module details", error.toString())
                setModule({} as Module);
            });
        }

        getModule();
    }, [module_name])

    return (
        <>
            <PageHeader
                title={<><span className="font-normal text-muted">Module</span> {mod.name}</>}
                loading={loading}
                meta={!loading && <span>{mod.module_type} module</span>}
            />
            <Section title="Defaults">
                <ModuleProperties loading={loading} module={mod} installed />
            </Section>
        </>
    );
}

export default ModuleDetails;