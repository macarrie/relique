import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";

import API from "../utils/api";

import Module from "../types/module";

function ModuleListRow(props :any) {
    let mod = props.module;

    console.log(mod);
    return (
        <tr className="hover:bg-slate-50">
            <td className="py-2 px-3"><Link to={`/modules/${mod.name}`}>{mod.name}</Link></td>
            <td className="py-2 px-3">{mod.module_type}</td>
            <td className="py-2 px-3">{mod.available_variants.join(", ")}</td>
        </tr>
    );
}

function ModuleList(props :any) {
    let [limit, setLimit] = useState(props.limit || 0);
    let [modules, setModuleList] = useState([]);

    useEffect(() => {
        setLimit(props.limit);
    }, [props.limit])

    useEffect(() => {
        function getModules() {
            API.modules.list({
                limit: limit,
            }).then((response :any) => {
                setModuleList(response.data);
            }).catch(error => {
                console.log("Cannot get module list", error);
            });
        }

        getModules();
    }, [limit])


    function renderModuleList() {
        if (!modules) {
            return (
                <>
                Loading
                </>
            )
        }

        const moduleList = modules.map((mod :Module) =>
            <ModuleListRow key={mod.name} module={mod}/>
        );

        return (
            <tbody>
            {moduleList}
            </tbody>
        )
    }

    return (
        <table className="table-auto w-full">
            <thead className="bg-slate-50 uppercase text-slate-500 text-left">
                <tr className="border border-l-0 border-r-0 border-slate-100">
                    <th className="py-2 px-3">Name</th>
                    <th className="py-2 px-3">Type</th>
                    <th className="py-2 px-3">Variants</th>
                </tr>
            </thead>
            {renderModuleList()}
        </table>
    );
}

export default ModuleList;
