import { Link } from "react-router-dom";

import { createColumnHelper } from "@tanstack/react-table";
import Table from "./table";
import Module from "../types/module";

function ModuleList({
    title = "",
    actions = true,
    custom_actions = [],
    sorted = true,
    paginated = true,
    data = {} as Module[],
}) {
    function renderVariants(vars: string) {
        return (
            <>
                {vars.split(",").map((v: any) => (
                    <span className="badge badge-neutral badge-sm badge-soft" key={v}>{v}</span>
                ))}
            </>
        )
    }

    const columnHelper = createColumnHelper<Module>()
    const columns = [
        columnHelper.accessor((mod: Module) => mod.name, {
            header: () => (<div>Name</div>),
            id: 'module_name',
            cell: (cell: any) => (<div><Link className="link-primary link-hover" to={`/modules/${cell.getValue()}`}>{cell.getValue()}</Link></div>)
        }),
        columnHelper.accessor('module_type', {
            header: () => (<div>Type</div>),
            id: 'type',
            cell: (cell: any) => (<div>{cell.getValue()}</div>),
        }),
        columnHelper.accessor((mod) => (mod.available_variants ?? []).join(", "), {
            header: () => (<div>Variants</div>),
            id: 'variants',
            cell: (cell: any) => (<div className="space-x-1">{renderVariants(cell.getValue())}</div>),
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

export default ModuleList;