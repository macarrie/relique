import { Link } from "react-router-dom";

import { useState } from "react";
import { createColumnHelper } from "@tanstack/react-table";
import Status from "./status";
import Job from "../types/job";
import Utils from "../utils/utils";
import JobUtils from "../utils/job";
import Table from "./table";
import JobRemoveModal from "./job_remove_modal";
import JobRetryButton from "./job_retry_button";
import TableUtils, { TableColumns } from "../utils/table";

function JobList({
    actions = true,
    sorted = true,
    paginated = true,
    data = [] as Job[],
    loading = true,
    // Display a delete button on each row
    deletable = false,
    // Called after a successful removal from a row delete button
    onDeleted = undefined as ((uuid: string, imageRemoved: boolean) => void) | undefined,
    // Display a retry button on each backup job row
    retryable = false,
}) {
    const [toRemove, setToRemove] = useState<Job | null>(null);

    const columnHelper = createColumnHelper<Job>()
    const columns: TableColumns<Job> = [
        columnHelper.accessor((job: Job) => job.uuid, {
            header: () => "ID",
            id: 'id',
            cell: (cell) => (<Link to={`/jobs/${cell.getValue()}`} className="code text-link">{Utils.shortID(cell.getValue())}</Link>),
            enableSorting: false,
        }),
        columnHelper.accessor((job: Job) => job.client?.name, {
            header: () => "Client",
            id: 'client_name',
            cell: (cell) => (<Link to={`/clients/${cell.getValue()}`} className="text-link">{cell.getValue()}</Link>)
        }),
        columnHelper.accessor((job: Job) => job.module?.name, {
            header: () => "Module",
            id: 'module_name',
            cell: (cell) => cell.getValue(),
        }),
        columnHelper.accessor('job_type', {
            header: () => "Type",
            id: 'job_type',
            cell: (cell) => cell.getValue(),
        }),
        columnHelper.accessor((job: Job) => job.status, {
            header: () => "Status",
            id: 'status',
            cell: (cell) => (<Status label={cell.getValue()} status={JobUtils.jobStateToCode(cell.getValue())} />),
        }),
        columnHelper.accessor('start_time', {
            header: () => "Start",
            id: 'start_time',
            cell: (cell) => (<span className="tabular-nums">{Utils.formatDate(cell.getValue())}</span>),
            sortingFn: 'datetime',
        }),
        columnHelper.accessor('end_time', {
            header: () => "End",
            id: 'end_time',
            cell: (cell) => (<span className="tabular-nums">{Utils.formatDate(cell.getValue())}</span>),
            sortingFn: 'datetime',
        }),
    ];

    if (deletable || retryable) {
        columns.push(columnHelper.display({
            id: 'actions',
            header: () => (<div></div>),
            cell: (cell) => {
                const job = cell.row.original;
                return (
                    <div className="flex justify-end gap-1">
                        {retryable && <JobRetryButton job={job} compact />}
                        {deletable && (
                            <button
                                type="button"
                                className="btn btn-ghost btn-xs text-error"
                                disabled={!job.done}
                                title={job.done ? "Delete job" : "Job is still running and cannot be deleted"}
                                onClick={() => setToRemove(job)}>
                                <i className="ri-delete-bin-line text-base"></i>
                            </button>
                        )}
                    </div>
                );
            },
        }));
    }

    return (
        <>
        {deletable && (
            <JobRemoveModal
                job={toRemove}
                onClose={() => setToRemove(null)}
                onRemoved={(uuid: string, linkedRemoved: boolean) => onDeleted?.(uuid, linkedRemoved)}
            />
        )}
        <Table
            actions={actions}
            data={loading ? TableUtils.GetPlaceholderRows<Job>() : data}
            columns={loading ? TableUtils.GetPlaceholderColumns(columns) : columns}
            paginated={paginated && !loading}
            sorted={sorted}
        />
        </>
    );
}

export default JobList;
