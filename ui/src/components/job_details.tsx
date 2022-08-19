import React, {useEffect, useState} from "react";
import { useParams } from "react-router-dom";

import API from "../utils/api";
import Job from "../types/job";

function JobDetails(props :any) {
    let urlParams = useParams();
    let [j, setJob] = useState<Job | null>(null);

    useEffect(() => {
        function getJob() {
            let uuid = urlParams["job_uuid"];
            if (uuid === undefined) {
                console.log("Job uuid undefined, cannot get job details");
                return;
            }

            API.jobs.get(uuid).then((response :any) => {
                setJob(response.data);
            }).catch(error => {
                console.log("Cannot get job details", error);
            });
        }

        getJob();
    }, [urlParams])

    if (j === null) {
        return <div>Loading</div>
    }

    return (
        <div>{JSON.stringify(j)}</div>
    );
}

export default JobDetails;