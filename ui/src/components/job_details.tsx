import React, {useEffect, useState} from "react";
import {Link, useParams} from "react-router-dom";

import API from "../utils/api";
import Job from "../types/job";
import ModuleCard from "./module_card";
import Moment from "react-moment";
import JobUtils from "../utils/job";
import JobLogs from "./job_logs";
import {Tabs, Tab} from "./tabs";

function JobDetails() {
    const {job_uuid} = useParams();
    let [j, setJob] = useState<Job | null>(null);

    function getJob() {
        if (job_uuid === undefined) {
            console.log("Job uuid undefined, cannot get job details");
            return;
        }

        API.jobs.get(job_uuid).then((response :any) => {
            setJob(response.data);
        }).catch(error => {
            console.log("Cannot get job details", error);
            setJob(null);
        });
    }

    useEffect(() => {
        getJob();
    }, [])

    if (j === null) {
        return <div>Loading</div>
    }

    return (<>
        <div className={"bg-white shadow rounded"}>
            <div className="flex flex-row px-4 py-3 items-center">
                <span className="flex-grow text-xl font-bold">
                    Job details
                </span>
                <span className="text-l ml-4 font-mono text-pink-500">{j.uuid}</span>
            </div>

            <div className="flex flex-col px-4 py-3 pb-4 bg-slate-50 space-y-3">
                <div className={"uppercase font-bold text-slate-500 mb-2"}>General info</div>
                <div className="flex flex-col md:flex-row content-center">
                    <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Job type</div>
                    <div className={"flex-grow py-2 px-3 text-slate-900"}>{j.job_type}</div>
                </div>
                {j.backup_type === "backup" && (
                    <div className="flex flex-col md:flex-row content-center">
                        <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Backup type</div>
                        <div className={"flex-grow py-2 px-3 text-slate-900"}>diff</div>
                    </div>
                )}
                {j.backup_type === "restore" && (
                    <>
                    <div className="flex flex-col md:flex-row content-center">
                        <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Restore destination</div>
                        <div className={"flex-grow py-2 px-3 font-mono text-pink-500"}>{j.restore_destination}</div>
                    </div>
                    <div className="flex flex-col md:flex-row content-center">
                        <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Restore source job</div>
                        <div className={"flex-grow py-2 px-3 text-slate-900"}><Link to={`/jobs/${j.restore_job_uuid}`}>{j.restore_job_uuid}</Link></div>
                    </div>
                    </>
                )}
                <div className="flex flex-col md:flex-row content-center">
                    <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Running</div>
                    <div className={"flex-grow py-2 px-3 text-slate-900"}>{j.done ? "Finished" : "In execution"}</div>
                </div>
                <div className="flex flex-col md:flex-row content-center">
                    <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Status</div>
                    <div className={"flex-grow py-2 px-3 text-slate-900"}>
                        <span className={`${JobUtils.statusColor(j.status)} font-bold capitalize`}>
                            {j.status}
                        </span>
                    </div>
                </div>
                <div className="flex flex-col md:flex-row content-center">
                    <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Start time</div>
                    <div className={"flex-grow py-2 px-3 text-slate-900"}>
                        <Moment date={j.start_time} format={"DD/MM/YYYY HH:mm:ss"}/>
                    </div>
                </div>
                <div className="flex flex-col md:flex-row content-center">
                    <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>End time</div>
                    <div className={"flex-grow py-2 px-3 text-slate-900"}>
                        <Moment date={j.end_time} format={"DD/MM/YYYY HH:mm:ss"}/>
                    </div>
                </div>
            </div>

            <hr />

            <div className="flex flex-col px-4 py-3 pb-4 bg-slate-50 space-y-3">
                <div className="flex flex-row items-center">
                    <div className={"flex-grow uppercase font-bold text-slate-500 mb-2"}>Client</div>
                    <Link to={`/clients/${j.client.id}`} className="bg-transparent rounded px-2 py-1 text-blue-500 hover:text-blue-900 uppercase text-xs font-bold">
                        Details
                    </Link>
                </div>
                <div className="flex flex-col md:flex-row content-center">
                    <div className={"w-24 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Name</div>
                    <div className={"flex-grow py-2 px-3 md:ml-6 text-slate-900"}>{j.client.name}</div>
                </div>
                <div className="flex flex-col md:flex-row content-center">
                    <div className={"w-24 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Address</div>
                    <div className={"flex-grow py-2 px-3 md:ml-6 text-slate-900"}>{j.client.address}</div>
                </div>
                <div className="flex flex-col md:flex-row content-center">
                    <div className={"w-24 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Port</div>
                    <div className={"flex-grow py-2 px-3 md:ml-6 text-slate-900"}>{j.client.port}</div>
                </div>
            </div>

            <hr />

            <div className="flex flex-col px-4 py-3 pb-4 bg-slate-50 space-y-3">
                <div className={"uppercase font-bold text-slate-500 mb-2"}>Module</div>
                <ModuleCard module={j.module} full />
            </div>
        </div>

        <div className={"bg-white shadow rounded mt-3"}>
            <Tabs title="Logs" initialActiveTab="pre">
                <Tab title="Setup script" key="pre">
                    {j.module.pre_backup_script === "none" ? (
                        <div className="center italic text-slate-400">No pre-backup/restore script configured in module</div>
                    ) : (
                        <div>TODO: Pre script logs contents</div>
                    )}
                </Tab>
                {j.module.backup_paths.map((path :string, index :number) => {
                    return <Tab headerClassName="font-mono text-xs" title={path} key={index}>
                        <JobLogs uuid={job_uuid} path={path} />
                    </Tab>
                })}
                <Tab title="Teardown script" key="post">
                    {j.module.post_backup_script === "none" ? (
                        <div className="center italic text-slate-400">No pre-backup/restore script configured in module</div>
                    ) : (
                        <div>TODO: Post script logs contents</div>
                    )}
                </Tab>
            </Tabs>
        </div>
        </>
    );
}

export default JobDetails;