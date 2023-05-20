import React, {useCallback, useEffect, useState} from "react";
import {Link} from "react-router-dom";
import Moment from "react-moment";
import {Column} from "react-table";

import API from "../utils/api";

import Job from "../types/job";
import JobUtils from "../utils/job";
import StatusBadge from "./status_badge";
import Table from "./table";
import TableUtils from "../utils/table";

type JobListProps = {
    limit? :number,
    title? :string,
    filtered? :boolean,
    sorted? :boolean,
}

function JobList(props :JobListProps) {
    let [limit, setLimit] = useState(props.limit || 0);
    let [jobs, setJobList] = useState([] as Job[]);
    let [loading, setLoading] = useState(true);

    function uuidDisplay(id :string) {
        return id.split("-")[0]
    }

    const getJobs = useCallback(() => {
        setLoading(true);

        API.jobs.list({
            limit: limit,
        }).then((response :any) => {
            setJobList(response.data || []);
            setLoading(false);
        }).catch(error => {
            console.log("Cannot get job list", error);
            setLoading(false);
        });
    }, [limit])

    useEffect(() => {
        setLimit(props.limit || 0);
    }, [props.limit])

    useEffect(() => {
        getJobs();
    }, [limit, getJobs])

    const columns :Array<Column<Job>> = React.useMemo(() => [
        {
            Header: () => (<div className="py-2 px-3 w-full text-center">ID</div>),
            accessor: (job) => uuidDisplay(job.uuid),
            id: 'id',
            Cell: ({value} :any) => (<div className="py-2 px-3 code"><Link to={`/jobs/${value}`}>{value}</Link></div>),
        },
        {
            Header: () => (<div className="py-2 px-3">Client</div>),
            id: 'client_name',
            accessor: (job) => job.client.name,
            Cell: ({value} :any) => (<div className="py-2 px-3"><Link to={`/clients/${value}`}>{value}</Link></div>),
        },
        {
            Header: () => (<div className="py-2 px-3">Module</div>),
            id: 'module_name',
            accessor: (job) => job.module.name,
            Cell: ({value} :any) => (<div className="py-2 px-3"><span className="badge">{value}</span></div>),
        },
        {
            Header: () => (<div className="py-2 px-3 hidden md:block">Type</div>),
            id: 'job_type',
            accessor: 'job_type',
            Cell: ({value} :any) => (<div className="py-2 px-3 hidden md:block">{value}</div>),
        },
        {
            Header: () => (<div className="py-2 px-3">Status</div>),
            id: 'status',
            accessor: (job) => job.status,
            Cell: ({value} :any) => (<div className="py-2 px-3"><StatusBadge label={value} status={JobUtils.jobStateToCode(value)}/></div>),
        },
        {
            Header: () => (<div className="py-2 px-3 hidden md:block">Start</div>),
            id: 'start_time',
            accessor: 'start_time',
            Cell: ({value} :any) => (<div className="py-2 px-3 hidden md:block"><Moment date={value} format={"DD/MM/YYYY HH:mm:ss"}/></div>),
        },
        {
            Header: () => (<div className="py-2 px-3 hidden md:block">End</div>),
            id: 'end_time',
            accessor: 'end_time',
            Cell: ({value} :any) => (<div className="py-2 px-3 hidden md:block"><Moment date={value} format={"DD/MM/YYYY HH:mm:ss"}/></div>),
        },
    ], []);

    if (loading) {
        return (
            <Table title={props.title}
                   filtered={false}
                   sorted={false}
                   refreshFunc={getJobs}
                   columns={TableUtils.GetPlaceholderColumns(columns)}
                   data={[{}, {}, {}]} />
        );
    }

    return (
        <Table title={props.title}
               filtered={props.filtered}
               sorted={props.sorted}
               refreshFunc={getJobs}
               columns={columns}
               data={jobs} />
    );
}

export default JobList;
