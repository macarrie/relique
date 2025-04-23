import { Link } from "react-router-dom";

import { createColumnHelper } from "@tanstack/react-table";
import Table from "./table";
import Repository from "../types/repository";

function RepoList({
    title = "",
    actions = true,
    custom_actions = [],
    sorted = true,
    paginated = true,
    data = {} as Repository[],
}) {
    const columnHelper = createColumnHelper<Repository>()
    const columns = [
        columnHelper.accessor((repo: Repository) => repo.name, {
            header: () => (<div>Name</div>),
            id: 'name',
            cell: (cell: any) => (<div><Link className="link-primary link-hover" to={`/repositories/${cell.getValue()}`}>{cell.getValue()}</Link></div>)
        }),
        columnHelper.accessor((repo: Repository) => repo.type, {
            header: () => (<div>Type</div>),
            id: 'type',
            cell: (cell: any) => (<div>{cell.getValue()}</div>)
        }),
    ];

    return (
        <Table
            title={title}
            actions={actions}
            custom_actions={custom_actions}
            data={data}
            columns={columns}
            paginated={paginated}
            sorted={sorted}
        />
    );
}

export default RepoList;