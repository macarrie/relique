import { Link } from "react-router-dom";

import { createColumnHelper } from "@tanstack/react-table";
import Table from "./table";
import Module from "../types/module";
import TableUtils from "../utils/table";
import EmptyValue from "./empty_value";

function ModuleList({
    actions = true,
    sorted = true,
    paginated = true,
    data = [] as Module[],
    loading = true,
}) {
    const columnHelper = createColumnHelper<Module>()
    const columns = [
        columnHelper.accessor((mod: Module) => mod.name, {
            header: () => "Name",
            id: 'module_name',
            cell: (cell) => (<Link className="text-link font-medium" to={`/modules/${cell.getValue()}`}>{cell.getValue()}</Link>)
        }),
        columnHelper.accessor('module_type', {
            header: () => "Type",
            id: 'type',
            cell: (cell) => cell.getValue(),
        }),
        columnHelper.accessor((mod: Module) => (mod.available_variants ?? []).join(", "), {
            header: () => "Variants",
            id: 'variants',
            cell: (cell) => (cell.getValue() ? cell.getValue() : <EmptyValue />),
        }),
    ];

    return (
        <Table
            actions={actions}
            data={loading ? TableUtils.GetPlaceholderRows<Module>() : data}
            columns={loading ? TableUtils.GetPlaceholderColumns(columns) : columns}
            paginated={paginated && !loading}
            sorted={sorted}
        />
    );
}

export default ModuleList;
