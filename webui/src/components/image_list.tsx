import { Link } from "react-router-dom";

import { createColumnHelper } from "@tanstack/react-table";
import Image from "../types/image";
import Utils from "../utils/utils";
import Table from "./table";

function ImageList({
    title = "",
    actions = true,
    custom_actions = [],
    sorted = true,
    paginated = true,
    data = {} as Image[],
}) {

    function uuidDisplay(id: string) {
        if (!id) {
            return "unknown"
        }

        return id.split("-")[0]
    }

    const columnHelper = createColumnHelper<Image>()
    const columns = [
        columnHelper.accessor((img: Image) => img.uuid, {
            header: () => (<div className="w-full">ID</div>),
            id: 'id',
            cell: (cell: any) => (<div className="code link-hover"><Link to={`/images/${cell.getValue()}`}>{uuidDisplay(cell.getValue())}</Link></div>),
            enableSorting: false,
        }),
        columnHelper.accessor((img: Image) => img.client?.name, {
            header: () => (<div>Client</div>),
            id: 'client_name',
            cell: (cell: any) => (<div><Link to={`/clients/${cell.getValue()}`} className="link link-primary link-hover">{cell.getValue()}</Link></div>)
        }),
        columnHelper.accessor((img: Image) => img.module?.name, {
            header: () => (<div>Module</div>),
            id: 'module_name',
            cell: (cell: any) => (<div><span className="badge badge-soft badge-neutral badge-sm">{cell.getValue()}</span></div>),
        }),
        columnHelper.accessor('created_at', {
            header: () => (<div>Date</div>),
            id: 'created_at',
            cell: (cell: any) => (<div>{Utils.formatDate(cell.getValue())}</div>),
            sortingFn: 'datetime',
        }),
        columnHelper.accessor('size_on_disk', {
            header: () => (<div>Size</div>),
            id: 'size_on_disk',
            cell: (cell: any) => (<div>{Utils.formatSize(cell.getValue())}</div>),
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

export default ImageList;