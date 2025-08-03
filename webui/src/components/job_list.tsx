import { Link } from "react-router-dom";

import { createColumnHelper } from "@tanstack/react-table";
import { Column } from "react-table";
import StatusBadge from "../components/status_badge";
import Job from "../types/job";
import Utils from "../utils/utils";
import JobUtils from "../utils/job";
import Table from "../components/table";
import React from "react";
import TableUtils from "../utils/table";

function JobList({
    title = "",
    actions = true,
    custom_actions = [] as React.ReactNode,
    sorted = true,
    paginated = true,
    data = {} as Job[],
    loading = true,
}) {

    function uuidDisplay(id: string) {
        if (!id) {
            return "unknown"
        }

        return id.split("-")[0]
    }

    const columnHelper = createColumnHelper<Job>()
    const columns = [
        columnHelper.accessor((job: Job) => job.uuid, {
            header: () => (<div className="w-full">ID</div>),
            id: 'id',
            cell: (cell: any) => (<div className="code"><Link to={`/jobs/${cell.getValue()}`} className="link-hover">{uuidDisplay(cell.getValue())}</Link></div>),
            enableSorting: false,
        }),
        columnHelper.accessor((job: Job) => job.client?.name, {
            header: () => (<div>Client</div>),
            id: 'client_name',
            cell: (cell: any) => (<div><Link to={`/clients/${cell.getValue()}`} className="link-hover link-primary">{cell.getValue()}</Link></div>)
        }),
        columnHelper.accessor((job: any) => job.module?.name, {
            header: () => (<div>Module</div>),
            id: 'module_name',
            cell: (cell: any) => (<div><span className="badge badge-soft badge-sm badge-neutral">{cell.getValue()}</span></div>),
        }),
        columnHelper.accessor('job_type', {
            header: () => (<div className="">Type</div>),
            id: 'job_type',
            cell: (cell: any) => (<div className="">{cell.getValue()}</div>),
        }),
        columnHelper.accessor((job: Job) => job.status, {
            header: () => (<div>Status</div>),
            id: 'status',
            cell: (cell: any) => (<div><StatusBadge label={cell.getValue()} status={JobUtils.jobStateToCode(cell.getValue())} /></div>),
        }),
        columnHelper.accessor('start_time', {
            header: () => (<div className="">Start</div>),
            id: 'start_time',
            cell: (cell: any) => (<div className="">{Utils.formatDate(cell.getValue())}</div>),
            sortingFn: 'datetime',
        }),
        columnHelper.accessor('end_time', {
            header: () => (<div className="">End</div>),
            id: 'end_time',
            cell: (cell: any) => (<div className="">{Utils.formatDate(cell.getValue())}</div>),
            sortingFn: 'datetime',
        }),
    ];

    return (
        <Table
            title={title}
            actions={actions}
            custom_actions={custom_actions}
            data={loading ? [{}, {}, {}] : data}
            columns={loading ? TableUtils.GetPlaceholderColumns(columns as Column[]) : columns}
            paginated={paginated}
            sorted={sorted}
        />
    );
}

export default JobList;