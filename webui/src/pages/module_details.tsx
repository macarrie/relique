import { useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import Card from '../components/card';
import API from '../utils/api';
import ModuleCard from '../components/module_card';
import Module from '../types/module';

function ModuleDetails() {
    const { module_name } = useParams();
    let [mod, setModule] = useState<Module>({} as Module);

    useEffect(() => {
        function getModule() {
            if (module_name === undefined) {
                console.log("Module name undefined, cannot get module details");
                return;
            }

            API.modules.get(module_name).then((response: any) => {
                setModule(response.data);
            }).catch(error => {
                console.log("Cannot get job details", error);
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
                    <ModuleCard module={mod} full />
                </div>
            </Card>
        </>
    );
}

export default ModuleDetails;