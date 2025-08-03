import { useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import Job from "../types/job";
import JobList from "../components/job_list";
import Const from "../types/const";
import Utils from "../utils/utils";

function Jobs() {
    let [jobs, setJobs] = useState<Job[]>([]);
    let [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getJobList() {
            setLoading(true);
            API.jobs.list({ limit: 10000 }).then((response: any) => {
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
            <Card>
                <JobList
                    title="All jobs"
                    actions={true}
                    data={jobs}
                    paginated={true}
                    sorted={true}
                    loading={loading}
                />
            </Card>
        </>
    );
}

export default Jobs;