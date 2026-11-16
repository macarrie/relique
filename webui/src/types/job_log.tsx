type JobLog = {
    name: string,
    backup_path: string,
    stdout: string,
    stdout_truncated: boolean,
    stderr: string,
    stderr_truncated: boolean,
};

export default JobLog;
