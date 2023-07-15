import React, {useCallback, useEffect, useMemo, useState} from "react";

import API from "../utils/api";

import Module from "../types/module";
import Table from "./table";
import Const from "../types/const";
import {createColumnHelper} from "@tanstack/react-table";
import {useQuery} from "react-query";
import {Link} from "react-router-dom";

function ModuleList(props :any) {
    let [modules, setModules] = useState([] as Module[]);
    let [loading, setLoading] = useState(true);

    const defaultData = useMemo(() => [], [])
    const fetchDataOptions = {
        limit: 1000,
        offset: 0,
    }
    const dataQuery = useQuery(
        ['modules', fetchDataOptions],
        () => API.modules.list(fetchDataOptions),
        { keepPreviousData: true }
    )

    useEffect(() => {
        setLoading(dataQuery.isLoading || dataQuery.isFetching)
    }, [dataQuery.isLoading, dataQuery.isFetching])

    const getModules = useCallback(function() {
        let moduleList = dataQuery.data?.data.data || defaultData;
        setModules(moduleList);
    }, [dataQuery, defaultData])

    useEffect(() => {
        getModules();
    }, [getModules]);

    function renderVariants(vars :string) {
        return (
            <>
                {vars.split(",").map((v: any) => (
                    <span className="badge" key={v}>{v}</span>
                ))}
            </>
        )
    }

    const columnHelper = createColumnHelper<Module>()
    const columns = [
        columnHelper.accessor('name', {
            header: () => (<div className="py-2 px-3">Name</div>),
            id: 'name',
            cell: (cell :any) => (<div className="py-2 px-3"><Link to={`/modules/${cell.getValue()}`}>{cell.getValue()}</Link></div>),
        }),
        columnHelper.accessor('module_type', {
            header: () => (<div className="py-2 px-3 hidden md:block">Type</div>),
            id: 'type',
            cell: (cell :any) => (<div className="py-2 px-3 hidden md:block">{cell.getValue()}</div>),
        }),
        columnHelper.accessor( (mod) => (mod.available_variants || []).join(", "), {
            header: () => (<div className="py-2 px-3 hidden md:block">Variants</div>),
            id: 'variants',
            cell: (cell :any) => (<div className="py-2 px-3 space-x-1">{renderVariants(cell.getValue())}</div>),
        }),
    ]

    return (
        <Table title={props.title}
               filtered={props.filtered}
               sorted={props.sorted}
               paginated={props.paginated}
               columns={columns}
               defaultPageSize={props.limit || Const.DEFAULT_PAGE_SIZE}
               refreshFunc={getModules}
               data={modules}
               loading={loading}
        />
    );
}

export default ModuleList;
