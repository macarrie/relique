import { Link } from "react-router-dom";

import { createColumnHelper } from "@tanstack/react-table";
import Table from "./table";
import Client from "../types/client";
import Module from "../types/module";
import Status from "./status";
import TableUtils from "../utils/table";
import EmptyValue from "./empty_value";
import ClientUtils from "../utils/client";

function ClientList({
    actions = true,
    sorted = true,
    paginated = true,
    data = [] as Client[],
    loading = true,
}) {

    const columnHelper = createColumnHelper<Client>()
    const columns = [
        columnHelper.accessor((client: Client) => client, {
            header: () => (<span className="sr-only">Health</span>),
            id: 'health',
            cell: (cell) => (
                <Status
                    variant="dot"
                    loading={cell.getValue().state_is_loading}
                    status={cell.getValue().ssh_alive}
                    label={ClientUtils.sshStatusLabel(cell.getValue().ssh_alive)} />
            ),
            enableSorting: false,
        }),
        columnHelper.accessor((c: Client) => c.name, {
            header: () => "Name",
            id: 'name',
            cell: (cell) => (<Link className="text-link font-medium" to={`/clients/${cell.getValue()}`}>{cell.getValue()}</Link>),
        }),
        columnHelper.accessor('address', {
            header: () => "Address",
            id: 'address',
            cell: (cell) => (<span className="code">{cell.getValue()}</span>),
        }),
        columnHelper.accessor((client: Client) => (client.modules || []).map((mod: Module) => mod.name).join(", "), {
            header: () => "Modules",
            id: 'modules',
            cell: (cell) => (cell.getValue() ? <span>{cell.getValue()}</span> : <EmptyValue />),
        }),
    ];

    return (
        <Table
            actions={actions}
            data={loading ? TableUtils.GetPlaceholderRows<Client>() : data}
            columns={loading ? TableUtils.GetPlaceholderColumns(columns) : columns}
            paginated={paginated && !loading}
            sorted={sorted}
        />
    );
}

export default ClientList;
