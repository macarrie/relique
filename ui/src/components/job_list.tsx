import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";
import Moment from "react-moment";

import API from "../utils/api";

import Job from "../types/job";
import Module from "../types/module";
import Client from "../types/client";
import JobUtils from "../utils/job";
import StatusBadge from "./status_badge";

function JobListRowPlaceholder(props :any) {
    return (
        <tr>
            <td className="py-2 px-3 code">
                <div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div>
            </td>
            <td className="py-2 px-3 hidden md:table-cell">
                <div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div>
            </td>
            <td className="py-2 px-3">
                <div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div>
            </td>
            <td className="py-2 px-3">
                <div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div>
            </td>
            <td className="py-2 px-3">
                <div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div>
            </td>
            <td className="py-2 px-3 hidden md:table-cell">
                <div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div>
            </td>
            <td className="py-2 px-3 hidden md:table-cell">
                <div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div>
            </td>
        </tr>
    );
}

function JobListRow(props :any) {
    function uuidDisplay(id :string) {
        return id.split("-")[0]
    }

    function moduleDisplayName(mod :Module) {
        let module_display = <>{mod.name} <span className="italic text-slate-400">({mod.variant})</span></>;
        if (mod.variant === "" || mod.variant === "default") {
            module_display = <>{mod.name}</>;
        }

        return module_display
    }

    function clientDisplayName(client :Client) {
        let client_display = <>{client.name}</>;
        if (client.name !== client.address) {
            client_display = <>{client.name}  <span className="italic hidden md:inline text-slate-400">({client.address})</span></>;
        }

        return client_display
    }

    let job = props.job;

    return (
        <tr>
            <td className="py-2 px-3 code"><Link to={`/jobs/${job.uuid}`}>{uuidDisplay(job.uuid)}</Link></td>
            <td className="py-2 px-3"><Link to={`/clients/${job.client.name}`}>{clientDisplayName(job.client)}</Link>
            </td>
            <td className="py-2 px-3"><span
                className="badge">{moduleDisplayName(job.module)}</span></td>
            <td className="py-2 px-3 hidden md:table-cell">{job.job_type}</td>
            <td className="py-2 px-3"><StatusBadge label={job.status} status={JobUtils.jobStateToCode(job.status)}/>
            </td>
            <td className="py-2 px-3 hidden md:table-cell"><Moment date={job.start_time}
                                                                   format={"DD/MM/YYYY HH:mm:ss"}/></td>
            <td className="py-2 px-3 hidden md:table-cell"><Moment date={job.end_time} format={"DD/MM/YYYY HH:mm:ss"}/>
            </td>
        </tr>
    );
}

function JobList(props :any) {
    let [limit, setLimit] = useState(props.limit || 0);
    let [jobs, setJobList] = useState([] as Job[]);
    let [loading, setLoading] = useState(true);

    useEffect(() => {
        setLimit(props.limit);
    }, [props.limit])

    useEffect(() => {
        function getJobs() {
            API.jobs.list({
                limit: limit,
            }).then((response :any) => {
                setJobList(response.data);
                setLoading(false);
            }).catch(error => {
                console.log("Cannot get job list", error);
                setLoading(false);
            });
        }

        getJobs();
    }, [limit])


    function renderJobList() {
        if (loading) {
            return (
                <tbody>
                    <JobListRowPlaceholder />
                    <JobListRowPlaceholder />
                    <JobListRowPlaceholder />
                </tbody>
            )
        }

        if (!jobs || jobs.length === 0) {
            return (
                <tbody>
                    <tr>
                        <td colSpan={8} className={"px-3 py-8 text-center text-3xl italic text-gray-300 dark:text-gray-600"}>
                            No jobs
                        </td>
                    </tr>
                </tbody>
            )
        }

        const jobList = jobs.map((job :Job) =>
            <JobListRow key={job.uuid} job={job}/>
        );

        return (
            <tbody>
            {jobList}
            </tbody>
        )
    }

    return (
        <table className="table table-auto w-full">
            <thead>
            <tr>
                <th className="py-2 px-3 text-center">ID</th>
                <th className="py-2 px-3">Client</th>
                <th className="py-2 px-3">Module</th>
                <th className="py-2 px-3 hidden md:table-cell">Type</th>
                <th className="py-2 px-3">Status</th>
                <th className="py-2 px-3 hidden md:table-cell">Start</th>
                <th className="py-2 px-3 hidden md:table-cell">End</th>
            </tr>
            </thead>
            {renderJobList()}
        </table>
    );
}

export default JobList;
