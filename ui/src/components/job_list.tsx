import React from "react";
import {Link} from "react-router-dom";
import Moment from "react-moment";
import {createColumnHelper} from "@tanstack/react-table";

import API from "../utils/api";

import Job from "../types/job";
import JobUtils from "../utils/job";
import StatusBadge from "./status_badge";
import Table from "./table";
import Const from "../types/const";

type JobListProps = {
    limit? :number,
    title? :string,
    filtered? :boolean,
    sorted? :boolean,
    paginated? :boolean,
}

function JobList(props :JobListProps) {
    function uuidDisplay(id :string) {
        if (!id) {
            return "unknown"
        }

        return id.split("-")[0]
    }

    const columnHelper = createColumnHelper<Job>()
    const columns = [
        columnHelper.accessor((job :Job) => job.uuid, {
            header: () => (<div className="py-2 px-3 w-full text-center">ID</div>),
            id: 'id',
            cell: (cell :any) => (<div className="py-2 px-3 code"><Link to={`/jobs/${cell.getValue()}`}>{uuidDisplay(cell.getValue())}</Link></div>),
        }),
        columnHelper.accessor((job :Job) => job.client.name, {
            header: () => (<div className="py-2 px-3">Client</div>),
            id: 'client_name',
            cell: (cell :any) => (<div className="py-2 px-3"><Link to={`/clients/${cell.getValue()}`}>{cell.getValue()}</Link></div>)
        }),
        columnHelper.accessor( (job :Job) => job.module?.name, {
            header: () => (<div className="py-2 px-3">Module</div>),
            id: 'module_name',
            cell: (cell :any) => (<div className="py-2 px-3"><span className="badge">{cell.getValue()}</span></div>),
        }),
        columnHelper.accessor('job_type', {
            header: () => (<div className="py-2 px-3 hidden md:block">Type</div>),
            id: 'job_type',
            cell: (cell :any) => (<div className="py-2 px-3 hidden md:block">{cell.getValue()}</div>),
        }),
        columnHelper.accessor((job :Job) => job.status, {
            header: () => (<div className="py-2 px-3">Status</div>),
            id: 'status',
            cell: (cell :any) => (<div className="py-2 px-3"><StatusBadge label={cell.getValue()} status={JobUtils.jobStateToCode(cell.getValue())}/></div>),
        }),
        columnHelper.accessor('start_time', {
            header: () => (<div className="py-2 px-3 hidden md:block">Start</div>),
            id: 'start_time',
            cell: (cell :any) => (<div className="py-2 px-3 hidden md:block"><Moment date={cell.getValue()} format={"DD/MM/YYYY HH:mm:ss"}/></div>),
        }),
        columnHelper.accessor('end_time', {
            header: () => (<div className="py-2 px-3 hidden md:block">End</div>),
            id: 'end_time',
            cell: (cell :any) => (<div className="py-2 px-3 hidden md:block"><Moment date={cell.getValue()} format={"DD/MM/YYYY HH:mm:ss"}/></div>),
        }),
    ];

    return (
        <Table title={props.title}
               filtered={props.filtered}
               sorted={props.sorted}
               paginated={props.paginated}
               columns={columns}
               defaultPageSize={props.limit || Const.DEFAULT_PAGE_SIZE}
               manualPagination={true}
               fetchDataFunc={API.jobs.list}
        />
    );
}

export default JobList;
