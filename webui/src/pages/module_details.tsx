import { useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import Card from '../components/card';
import API from '../utils/api';
import ModuleCard from '../components/module_card';
import Module from '../types/module';
import Const from '../types/const';
import Utils from '../utils/utils';

function ModuleDetails() {
    const { module_name } = useParams();
    let [mod, setModule] = useState<Module>({} as Module);
    let [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getModule() {
            setLoading(true);
            if (module_name === undefined) {
                console.log("Module name undefined, cannot get module details");
                return;
            }

            API.modules.get(module_name).then((response: any) => {
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
            <Card>
                <div className="px-6 py-4 flex">
                    <h3 className="flex-grow font-bold">
                        Module details
                    </h3>
                </div>
                <div className="m-4">
                    <ModuleCard loading={loading} module={mod} full />
                </div>
            </Card>
        </>
    );
}

export default ModuleDetails;