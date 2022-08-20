import Client from "./client"
import Module from "./module"

type Job = {
    uuid :string,
    client :Client,
    module :Module,
    status :string,
    done :boolean,
    backup_type :string,
    job_type :string,
    previous_job_uuid :string,
    restore_job_uuid :string,
    restore_destination :string,
    start_time :any,
    end_time :any,
};

export default Job;

