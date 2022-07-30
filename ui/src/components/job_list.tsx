import React from "react";

import API from "../utils/api";

import Job from "../types/job";
import Module from "../types/module";
import Client from "../types/client";
import JobStatus from "../types/job_status";

class JobListRow extends React.Component<any, any> {
    uuidDisplay(id :string) {
        return id.split("-")[0]
    }

    moduleDisplayName(mod :Module) {
        let module_display = <>{mod.name} <span className="italic text-slate-400">({mod.variant})</span></>;
        if (mod.variant === "" || mod.variant === "default") {
            module_display = <>{mod.name}</>;
        }

        return module_display
    }

    clientDisplayName(client :Client) {
        let client_display = <>{client.name}</>;
        if (client.name !== client.address) {
            client_display = <>{client.name}  <span className="italic text-slate-400">({client.address})</span></>;
        }

        return client_display
    }

    statusDisplay(s :JobStatus) {
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

    render() {
        let job = this.props.job;

        return (
            <tr className="hover:bg-slate-50">
                <td className="py-2 px-3">{this.uuidDisplay(job.uuid)}</td>
                <td className="py-2 px-3">{this.clientDisplayName(job.client)}</td>
                <td className="py-2 px-3">{this.moduleDisplayName(job.module)}</td>
                <td className="py-2 px-3">{this.statusDisplay(job.status)}</td>
                <td className="py-2 px-3">{job.start_time}</td>
                <td className="py-2 px-3">{job.end_time}</td>
            </tr>
        );
    }
}

type State = {
    jobs :Job[]
};

class JobList extends React.Component<any, State> {
    get_jobs_interval :number;
    limit :number;

    constructor(props: any) {
        super(props);

        this.get_jobs_interval = 0;
        this.limit = this.props.limit ? this.props.limit : 0;
    }

    state :State = {
        jobs: [],
    };

    componentDidMount() {
        this.getJobs();
    }

    componentWillUnmount() {}

    getJobs() {
        API.jobs.list({
            limit: this.limit,
        }).then((response :any) => {
            this.setState({
                jobs: response.data,
            });
        }).catch(error => {
            console.log("Cannot get job list", error);
        });
    }

    renderJobList() {
        if (!this.state.jobs) {
            return (
                <>
                Loading
                </>
            )
        }

        const jobList = this.state.jobs.map((job :Job) =>
            <JobListRow key={job.uuid} job={job}/>
        );

        return (
            <tbody>
            {jobList}
            </tbody>
        )
    }

    render() {
        return (
            <table className="table-auto w-full">
            <thead className="bg-slate-50 uppercase text-slate-500 text-left">
            <tr className="border border-l-0 border-r-0 border-slate-100">
            <th className="py-2 px-3">ID</th>
            <th className="py-2 px-3">Client</th>
            <th className="py-2 px-3">Module</th>
            <th className="py-2 px-3">Status</th>
            <th className="py-2 px-3">Start</th>
            <th className="py-2 px-3">End</th>
            </tr>
            </thead>
            {this.renderJobList()}
            </table>
        );
    }
}

export default JobList;
