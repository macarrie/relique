import { Link } from "react-router-dom";

import { createColumnHelper } from "@tanstack/react-table";
import Table from "./table";
import Repository from "../types/repository";
import TableUtils from "../utils/table";
import EmptyValue from "./empty_value";

function RepoList({
    actions = true,
    sorted = true,
    paginated = true,
    data = [] as Repository[],
    loading = true,
}) {
    const columnHelper = createColumnHelper<Repository>()
    const columns = [
        columnHelper.accessor((repo: Repository) => repo.name, {
            header: () => "Name",
            id: 'name',
            cell: (cell) => (
                <span className="flex items-baseline gap-2">
                    <Link className="text-link font-medium" to={`/repositories/${cell.row.original.name}`}>{cell.getValue()}</Link>
                    {cell.row.original.default && <span className="text-muted">default</span>}
                </span>
            )
        }),
        columnHelper.accessor((repo: Repository) => repo.type, {
            header: () => "Type",
            id: 'type',
            cell: (cell) => cell.getValue(),
        }),
        columnHelper.accessor((repo: Repository) => repo.path, {
            header: () => "Storage path",
            id: 'path',
            cell: (cell) => (cell.getValue() ? <span className="code">{cell.getValue()}</span> : <EmptyValue />),
        }),
    ];

    return (
        <Table
            actions={actions}
            data={loading ? TableUtils.GetPlaceholderRows<Repository>() : data}
            columns={loading ? TableUtils.GetPlaceholderColumns(columns) : columns}
            paginated={paginated && !loading}
            sorted={sorted}
        />
    );
}

export default RepoList;
