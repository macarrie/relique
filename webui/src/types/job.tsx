import Client from "./client";
import Module from "./module";
import Repository from "./repository";

type Job = {
    uuid: string,
    client: Client,
    module: Module,
    repository: Repository,
    status: string,
    done: boolean,
    backup_type: string,
    job_type: string,
    start_time: any,
    end_time: any,
};

export default Job;