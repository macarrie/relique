import { useEffect, useState } from "react";

import PageHeader from "../components/page_header";
import Section from "../components/section";
import API from "../utils/api";
import Job from "../types/job";
import JobList from "../components/job_list";
import Const from "../types/const";
import Utils from "../utils/utils";

function Jobs() {
    const [jobs, setJobs] = useState<Job[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getJobList() {
            setLoading(true);
            API.jobs.list({ limit: 10000 }).then((response) => {
                setJobs(response.data.data ?? []);
                setLoading(false);
            }).catch(error => {
                console.log("Cannot get job list", error);
                Utils.notify(Const.CRITICAL, "Cannot get jobs list", error.toString())
                setJobs([]);
            });
        }

        getJobList();
    }, [])

    return (
        <>
            <PageHeader title="Jobs" />
            <Section>
                <JobList
                    actions={true}
                    data={jobs}
                    paginated={true}
                    sorted={true}
                    loading={loading}
                    deletable={true}
                    onDeleted={(uuid: string) => setJobs(prev => prev.filter(j => j.uuid !== uuid))}
                />
            </Section>
        </>
    );
}

export default Jobs;