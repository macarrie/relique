import { Link } from "react-router-dom";

import { createColumnHelper } from "@tanstack/react-table";
import Image from "../types/image";
import Utils from "../utils/utils";
import Table from "./table";
import TableUtils, { TableColumns } from "../utils/table";
import { useState } from "react";
import ImageRemoveModal from "./image_remove_modal";

function ImageList({
    actions = true,
    sorted = true,
    paginated = true,
    data = [] as Image[],
    loading = true,
    // Display a delete button on each row
    deletable = false,
    // Called after a successful removal from a row delete button
    onDeleted = undefined as ((uuid: string, jobRemoved: boolean) => void) | undefined,
}) {
    const [toRemove, setToRemove] = useState<Image | null>(null);

    const columnHelper = createColumnHelper<Image>()
    const columns: TableColumns<Image> = [
        columnHelper.accessor((img: Image) => img.uuid, {
            header: () => "ID",
            id: 'id',
            cell: (cell) => (<Link to={`/images/${cell.getValue()}`} className="code text-link">{Utils.shortID(cell.getValue())}</Link>),
            enableSorting: false,
        }),
        columnHelper.accessor((img: Image) => img.client?.name, {
            header: () => "Client",
            id: 'client_name',
            cell: (cell) => (<Link to={`/clients/${cell.getValue()}`} className="text-link">{cell.getValue()}</Link>)
        }),
        columnHelper.accessor((img: Image) => img.module?.name, {
            header: () => "Module",
            id: 'module_name',
            cell: (cell) => cell.getValue(),
        }),
        columnHelper.accessor('created_at', {
            header: () => "Date",
            id: 'created_at',
            cell: (cell) => (<span className="tabular-nums">{Utils.formatDate(cell.getValue())}</span>),
            sortingFn: 'datetime',
        }),
        columnHelper.accessor('size_on_disk', {
            header: () => "Size",
            id: 'size_on_disk',
            cell: (cell) => (<span className="tabular-nums">{Utils.formatSize(cell.getValue())}</span>),
        }),
    ];

    if (deletable) {
        columns.push(columnHelper.display({
            id: 'actions',
            header: () => (<div></div>),
            cell: (cell) => (
                <div className="text-right">
                    <button
                        type="button"
                        className="btn btn-ghost btn-xs text-error"
                        title="Delete image"
                        onClick={() => setToRemove(cell.row.original)}>
                        <i className="ri-delete-bin-line text-base"></i>
                    </button>
                </div>
            ),
        }));
    }

    return (
        <>
        {deletable && (
            <ImageRemoveModal
                image={toRemove}
                onClose={() => setToRemove(null)}
                onRemoved={(uuid: string, linkedRemoved: boolean) => onDeleted?.(uuid, linkedRemoved)}
            />
        )}
        <Table
            actions={actions}
            data={loading ? TableUtils.GetPlaceholderRows<Image>() : data}
            columns={loading ? TableUtils.GetPlaceholderColumns(columns) : columns}
            paginated={paginated && !loading}
            sorted={sorted}
        />
        </>
    );
}

export default ImageList;
