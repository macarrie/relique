import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";

import API from "../utils/api";

import Module from "../types/module";

function ModuleListRow(props :any) {
    let mod = props.module;

    console.log(mod);
    return (
        <tr>
            <td className="py-2 px-3"><Link to={`/modules/${mod.name}`}>{mod.name}</Link></td>
            <td className="py-2 px-3">{mod.module_type}</td>
            <td className="py-2 px-3 space-x-1">{mod.available_variants.map((v: any) => (
                <span className="badge">{v}</span>))}</td>
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
        <table className="table table-auto w-full">
            <thead>
            <tr>
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
