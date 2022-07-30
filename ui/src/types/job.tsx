import Client from "./client"
import Module from "./module"
import JobStatus from "./job_status"

type Job = {
    uuid :string,
    client :Client,
    module :Module,
    status :JobStatus,
    done :boolean,
    backup_type :any,
    job_type :any,
    previous_job_uuid :string,
    restore_job_uuid :string,
    restore_destination :string,
    start_time :any,
    end_time :any,
};

export default Job;

