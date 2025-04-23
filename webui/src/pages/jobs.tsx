import { useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import Job from "../types/job";
import JobList from "../components/job_list";

function Jobs() {
    let [jobs, setJobs] = useState<Job[]>([]);

    useEffect(() => {
        function getJobList() {
            API.jobs.list({ limit: 10000 }).then((response: any) => {
                setJobs(response.data.data ?? []);
            }).catch(error => {
                console.log("Cannot get job list", error);
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
                />
            </Card>
        </>
    );
}

export default Jobs;