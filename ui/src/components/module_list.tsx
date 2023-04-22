import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";

import API from "../utils/api";

import Module from "../types/module";

function ModuleListRowPlaceholder() {
    return (
        <tr className="animate-pulse">
            <td className="py-2 px-3"><div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div></td>
            <td className="py-2 px-3"><div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div></td>
            <td className="py-2 px-3 hidden md:table-cell">
                <div className="flex flex-row space-x-1">
                    <div className="rounded-full h-2 w-12 bg-slate-300 dark:bg-slate-600"></div>
                    <div className="rounded-full h-2 w-12 bg-slate-300 dark:bg-slate-600"></div>
                    <div className="rounded-full h-2 w-12 bg-slate-300 dark:bg-slate-600"></div>
                </div>
            </td>
        </tr>
    );
}

function ModuleListRow(props :any) {
    let mod = props.module;

    return (
        <tr>
            <td className="py-2 px-3"><Link to={`/modules/${mod.name}`}>{mod.name}</Link></td>
            <td className="py-2 px-3 hidden md:table-cell">{mod.module_type}</td>
            <td className="py-2 px-3 space-x-1">{mod.available_variants.map((v: any) => (
                <span className="badge" key={v}>{v}</span>))}</td>
        </tr>
    );
}

function ModuleList(props :any) {
    let [limit, setLimit] = useState(props.limit || 0);
    let [modules, setModuleList] = useState([] as Module[]);
    let [loading, setLoading] = useState(true);

    useEffect(() => {
        setLimit(props.limit);
    }, [props.limit])

    useEffect(() => {
        function getModules() {
            API.modules.list({
                limit: limit,
            }).then((response :any) => {
                setModuleList(response.data);
                setLoading(false)
            }).catch(error => {
                console.log("Cannot get module list", error);
                setLoading(false)
            });
        }

        getModules();
    }, [limit])


    function renderModuleList() {
        if (loading) {
            return (
                <tbody>
                    <ModuleListRowPlaceholder />
                    <ModuleListRowPlaceholder />
                    <ModuleListRowPlaceholder />
                </tbody>
            )
        }

        if (!modules || modules.length === 0) {
            return (
                <tbody>
                    <tr>
                        <td colSpan={3} className={"px-3 py-8 text-center text-3xl italic text-gray-300 dark:text-gray-600"}>
                            No modules
                        </td>
                    </tr>
                </tbody>
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
                <th className="py-2 px-3 hidden md:table-cell">Type</th>
                <th className="py-2 px-3">Variants</th>
            </tr>
            </thead>
            {renderModuleList()}
        </table>
    );
}

export default ModuleList;
