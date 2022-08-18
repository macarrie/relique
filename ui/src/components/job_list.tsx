import React, {useEffect, useState} from "react";
import {Link} from "react-router-dom";
import Moment from "react-moment";

import API from "../utils/api";

import Job from "../types/job";
import Module from "../types/module";
import Client from "../types/client";
import JobStatus from "../types/job_status";

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
            client_display = <>{client.name}  <span className="italic text-slate-400">({client.address})</span></>;
        }

        return client_display
    }

    function statusDisplay(s :JobStatus) {
        let status_str :string[] = [
            "Pending",
            "Active",
            "Success",
            "Incomplete",
            "Error",
        ];

        let status_colors :string[] = [
            "text-neutral-700",
            "text-sky-700",
            "text-teal-700",
            "text-orange-700",
            "text-red-700",
        ];

        return <span className={status_colors[s.status]}>{status_str[s.status]}</span>
    }

    let job = props.job;

    return (
        <tr className="hover:bg-slate-50">
            <td className="py-2 px-3 hidden md:table-cell">{uuidDisplay(job.uuid)}</td>
            <td className="py-2 px-3"><Link to={`/clients/${job.client.id}`}>{clientDisplayName(job.client)}</Link></td>
            <td className="py-2 px-3">{moduleDisplayName(job.module)}</td>
            <td className="py-2 px-3">{statusDisplay(job.status)}</td>
            <td className="py-2 px-3 hidden md:table-cell"><Moment date={job.start_time} format={"DD/MM/YYYY HH:mm:ss"}/></td>
            <td className="py-2 px-3 hidden md:table-cell"><Moment date={job.end_time} format={"DD/MM/YYYY HH:mm:ss"}/></td>
        </tr>
    );
}

function JobList(props :any) {
    let [limit, setLimit] = useState(props.limit ? props.limit : 0);
    let [jobs, setJobList] = useState([]);

    useEffect(() => {
        setLimit(props.limit);
    }, [props.limit])

    useEffect(() => {
        function getJobs() {
            API.jobs.list({
                limit: limit,
            }).then((response :any) => {
                setJobList(response.data);
            }).catch(error => {
                console.log("Cannot get job list", error);
            });
        }

        getJobs();
    }, [limit])


    function renderJobList() {
        if (!jobs) {
            return (
                <>
                Loading
                </>
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
        <table className="table-auto w-full">
            <thead className="bg-slate-50 uppercase text-slate-500 text-left">
                <tr className="border border-l-0 border-r-0 border-slate-100">
                    <th className="py-2 px-3 text-center hidden md:table-cell">ID</th>
                    <th className="py-2 px-3">Client</th>
                    <th className="py-2 px-3">Module</th>
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
