import { Link } from "react-router-dom";

import { createColumnHelper } from "@tanstack/react-table";
import Table from "./table";
import Client from "../types/client";
import Module from "../types/module";
import StatusDot from "./status_dot";
import React from "react";

function ClientList({
    title = "",
    actions = true,
    custom_actions = [] as React.ReactNode[],
    sorted = true,
    paginated = true,
    data = {} as Client[],
}) {

    const columnHelper = createColumnHelper<Client>()
    const columns = [
        columnHelper.accessor((client) => { return client }, {
            header: () => (<div className="text-center">Health</div>),
            id: 'health',
            cell: (cell: any) => (<div className="text-center">{cell.getValue().state_is_loading ? <span className="text-neutral-300 loading loading-spinner loading-xs"></span> : <StatusDot status={cell.getValue().ssh_alive} />}</div>),
        }),
        columnHelper.accessor((c: Client) => c.name, {
            header: () => (<div>Name</div>),
            id: 'name',
            cell: (cell: any) => (<div><Link className="link-primary link-hover" to={`/clients/${cell.getValue()}`}>{cell.getValue()}</Link></div>),
        }),
        columnHelper.accessor('address', {
            header: () => (<div>Address</div>),
            id: 'address',
            cell: (cell: any) => (<div className="code">{cell.getValue()}</div>),
        }),
        columnHelper.accessor((client) => (client.modules || []).map((mod: Module) => mod.name).join(", "), {
            header: () => (<div>Modules</div>),
            id: 'modules',
            cell: (cell: any) => (<div className="space-x-1">{renderModules(cell.getValue())}</div>),
        }),
    ];

    function renderModules(mods: string) {
        if (!mods) {
            return <span className="italic text-base-content/20">None</span>;
        }

        let module_names: string[] = mods.split(",")
        return (
            <div className="flex flex-wrap gap-y-2 gap-x-1">
                {module_names.map((mod: any) => (
                    <div key={mod}>
                        <span className="badge badge-sm badge-neutral badge-soft">{mod}</span>
                    </div>
                ))}
            </div>
        )
    }

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

export default ClientList;