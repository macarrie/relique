import React, {useCallback, useEffect, useState} from "react";

import API from "../utils/api";

import Module from "../types/module";
import {Column} from "react-table";
import Table from "./table";
import TableUtils from "../utils/table";

function ModuleList(props :any) {
    let [limit, setLimit] = useState(props.limit || 0);
    let [modules, setModuleList] = useState([] as Module[]);
    let [loading, setLoading] = useState(true);

    const getModules = useCallback(() => {
        API.modules.list({
            limit: limit,
        }).then((response :any) => {
            setModuleList(response.data || []);
            setLoading(false)
        }).catch(error => {
            console.log("Cannot get module list", error);
            setLoading(false)
        });
    }, [limit])

    useEffect(() => {
        setLimit(props.limit);
    }, [props.limit])

    useEffect(() => {
        getModules();
    }, [getModules])

    function renderVariants(vars :string) {
        return (
            <>
                {vars.split(",").map((v: any) => (
                    <span className="badge" key={v}>{v}</span>
                ))}
            </>
        )
    }

    const columns :Array<Column<Module>> = React.useMemo(() => [
        {
            Header: () => (<div className="py-2 px-3">Name</div>),
            accessor: 'name',
            id: 'name',
            Cell: ({value} :any) => (<div className="py-2 px-3">{value}</div>),
        },
        {
            Header: () => (<div className="py-2 px-3 hidden md:block">Type</div>),
            accessor: 'module_type',
            id: 'type',
            Cell: ({value} :any) => (<div className="py-2 px-3 hidden md:block">{value}</div>),
        },
        {
            Header: () => (<div className="py-2 px-3 hidden md:block">Variants</div>),
            accessor: (mod) => (mod.available_variants || []).join(", "),
            id: 'variants',
            Cell: ({value} :any) => (<div className="py-2 px-3 space-x-1">{renderVariants(value)}</div>),
        },
    ], []);

    if (loading) {
        return (
            <Table title={props.title}
                   filtered={false}
                   sorted={false}
                   columns={TableUtils.GetPlaceholderColumns(columns)}
                   data={[{}, {}, {}]} />
        );
    }

    return (
        <Table title={props.title}
               filtered={props.filtered}
               sorted={props.sorted}
               refreshFunc={getModules}
               columns={columns}
               data={modules} />
    );
}

export default ModuleList;
