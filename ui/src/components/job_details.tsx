import React, {useEffect, useState} from "react";
import {Link, useParams} from "react-router-dom";

import API from "../utils/api";
import Job from "../types/job";
import ModuleCard from "./module_card";
import Moment from "react-moment";
import JobUtils from "../utils/job";
import JobLogs from "./job_logs";
import {Tabs, Tab} from "./tabs";
import StatusBadge from "./status_badge";
import Card from "./card";

function JobDetails() {
    const {job_uuid} = useParams();
    let [j, setJob] = useState<Job | null>(null);

    useEffect(() => {
        function getJob() {
            if (job_uuid === undefined) {
                console.log("Job uuid undefined, cannot get job details");
                return;
            }

            API.jobs.get(job_uuid).then((response: any) => {
                setJob(response.data);
            }).catch(error => {
                console.log("Cannot get job details", error);
                setJob(null);
            });
        }

        getJob();
    }, [job_uuid])

    if (j === null) {
        return <div>Loading</div>
    }

    console.log(j)
    return (<>
            <Card>
                <div className="flex flex-row px-4 py-3 items-center">
                <span className="flex-grow text-xl font-bold">
                    Job details
                </span>
                    <span className="text-l ml-4 code">{j.uuid}</span>
                </div>

                <div className="grid md:grid-cols-2 gap-4 m-4">
                    <Card className="bg-white bg-opacity-60">
                        <div className="p-4 flex flex-row items-center mb-2">
                            <div className="font-bold text-slate-500 dark:text-slate-200">General info</div>
                        </div>
                        <table className="details-table ml-4">
                            <tr>
                                <td>Job type</td>
                                <td>{j.job_type}</td>
                            </tr>
                            {j.job_type === "backup" && (
                                <tr>
                                    <td>Backup type</td>
                                    <td>{j.backup_type}</td>
                                </tr>
                            )}
                            {j.backup_type === "restore" && (
                                <>
                                    <tr>
                                        <td>Restore destination</td>
                                        <td>{j.restore_destination}</td>
                                    </tr>
                                    <tr>
                                        <td>Source job</td>
                                        <td><Link to={`/jobs/${j.restore_job_uuid}`}>{j.restore_job_uuid}</Link></td>
                                    </tr>
                                </>
                            )}
                            <tr>
                                <td>Running</td>
                                <td>{j.done ? "Finished" : "In execution"}</td>
                            </tr>
                            <tr>
                                <td>Status</td>
                                <td><StatusBadge label={j.status} status={JobUtils.jobStateToCode(j.status)}/></td>
                            </tr>
                            <tr>
                                <td>Start time</td>
                                <td><Moment date={j.start_time} format={"DD/MM/YYYY HH:mm:ss"}/></td>
                            </tr>
                            <tr>
                                <td>End time</td>
                                <td><Moment date={j.end_time} format={"DD/MM/YYYY HH:mm:ss"}/></td>
                            </tr>
                            <tr>
                                <td>Storage root</td>
                                <td className="code">{j.storage_root}</td>
                            </tr>
                        </table>
                    </Card>
                    <Card className="bg-white bg-opacity-60">
                        <div className="p-4 flex flex-row items-center mb-2">
                            <div className={"flex-grow font-bold text-slate-500 dark:text-slate-200"}>Client</div>
                            <Link className="button button-small button-text" to={`/clients/${j.client.name}`}>
                                Details
                            </Link>
                        </div>
                        <table className="details-table ml-4">
                            <tr>
                                <td>Name</td>
                                <td>{j.client.name}</td>
                            </tr>
                            <tr>
                                <td>Address</td>
                                <td className="code text-base">{j.client.address}</td>
                            </tr>
                            <tr>
                                <td>Port</td>
                                <td className="code text-base">{j.client.port}</td>
                            </tr>
                        </table>
                    </Card>
                </div>

                <div className="flex flex-col px-4 py-3 pb-4 space-y-3">
                    <div className={"font-bold text-slate-500 dark:text-slate-200 mb-2"}>Module</div>
                    <ModuleCard module={j.module} full/>
                </div>
            </Card>

            <Card className="mt-4">
                <Tabs title="Logs" initialActiveTab="pre">
                    <Tab title="Setup script" key="pre">
                        {j.module.pre_backup_script === "none" ? (
                            <div className="center italic text-slate-400">No pre-backup/restore script configured in
                                module</div>
                        ) : (
                            <div>TODO: Pre script logs contents</div>
                        )}
                    </Tab>
                    {j.module.backup_paths.map((path: string, index: number) => {
                        return <Tab headerClassName="font-mono text-xs" title={path} key={index}>
                            <JobLogs uuid={job_uuid} path={path}/>
                        </Tab>
                    })}
                    <Tab title="Teardown script" key="post">
                        {j.module.post_backup_script === "none" ? (
                            <div className="center italic text-slate-400">No pre-backup/restore script configured in
                                module</div>
                        ) : (
                            <div>TODO: Post script logs contents</div>
                        )}
                    </Tab>
                </Tabs>
            </Card>
        </>
    );
}

export default JobDetails;