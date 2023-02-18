import React, {useState} from "react";

import API from "../utils/api";

function JobLogs(props :any) {
    let uuid :string = props.uuid;
    let path :string = props.path;

    let [loaded, setLoaded] = useState<boolean>(false);
    let [logs, setLogs] = useState<string>("");
    let [err, setErr] = useState<string>("");


    function getJobLog(path :string) {
        if (uuid === undefined) {
            console.log("Job uuid undefined, cannot get logs");
            return;
        }

        API.jobs.getLogs(uuid, path).then((response :any) => {
            setLogs(response.data);
            setLoaded(true)
        }).catch(error => {
            console.log("Cannot get job logs", error);
            setErr(logRequestErrorMessage(error));
        });
    }

    function logRequestErrorMessage(error :any) :string {
        let baseMsg = "Cannot retrieve job logs (code " +error.response.status +")";
        if (error.response.data) {
            return baseMsg + ": " + error.response.data;
        }

        return baseMsg;
    }

    if (!loaded) {
        return (<div className="flex flex-col content-center">
            <button
                className="button mx-auto my-4"
                onClick={() => getJobLog(path)}>
                Load logs
            </button>
            {err && (
                <div
                    className="mt-4 bg-red-100 px-5 py-3 rounded text-red-900 dark:bg-red-900/50 dark:text-slate-200">{err}</div>
            )}
        </div>);
    }

    let lineCount = logs.split('\n').length;
    let lineCountBlock = ""
    for (let i = 0; i < lineCount; i++) {
        lineCountBlock += (i + 1) + "\n";
    }
    return (<>
        <div className="flex flow-row">
            <div
                className="whitespace-pre font-mono text-pink-500 text-right mr-2 pr-2 border-r-2 border-slate-200 dark:border-slate-600">
                {lineCountBlock}
            </div>

            <code className="block whitespace-pre code overflow-x-scroll">
                {logs}
            </code>
        </div>
    </>);
}

export default JobLogs;
