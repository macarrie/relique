import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { AxiosError, AxiosResponse } from "axios";

import Job from "../types/job";
import Const from "../types/const";
import API from "../utils/api";
import Utils from "../utils/utils";

// Display time of the confirmation notification before opening the new job details
const REDIRECT_DELAY_MS = 3000;

// Starts a new job with the same parameters as a finished backup job, then opens the new job details
// once the confirmation notification closes.
// Nothing is rendered for restore jobs: they overwrite data on the client and are never retried automatically.
type Props = {
    job: Job,
    // Compact icon-only button for table rows
    compact?: boolean,
    disabled?: boolean,
};

function JobRetryButton({ job, compact = false, disabled = false }: Props) {
    const [retrying, setRetrying] = useState<boolean>(false);
    const navigate = useNavigate();
    // The redirection happens when the notification closes: skip it if the user already left the page
    const mounted = useRef<boolean>(true);
    useEffect(() => {
        mounted.current = true;
        return () => {
            mounted.current = false;
        };
    }, []);

    if (job.job_type !== "backup") {
        return null;
    }

    function retry() {
        setRetrying(true);
        API.jobs.retry(job.uuid).then((response: AxiosResponse<Job>) => {
            const newUuid = response.data.uuid;
            // The button stays disabled until the redirection to avoid starting the same retry twice
            Utils.notify(Const.OK, "Job retry started", `New job ${Utils.shortID(newUuid)} started. Opening its details.`, {
                autoClose: REDIRECT_DELAY_MS,
                onClose: () => {
                    if (mounted.current) {
                        navigate(`/jobs/${newUuid}`);
                    }
                },
            });
        }).catch((error: AxiosError<{ error?: string }>) => {
            console.log("Cannot retry job", error);
            Utils.notify(Const.CRITICAL, "Cannot retry job", error.response?.data?.error ?? error.toString());
            setRetrying(false);
        });
    }

    // A successful job rarely needs a retry: keep the action available but quieter
    const faded = job.status === "success" ? "opacity-50 hover:opacity-100" : "";
    const title = job.done ? "Retry job with the same parameters" : "Job is still running and cannot be retried";

    return (
        <button
            type="button"
            className={compact ? `btn btn-ghost btn-xs ${faded}` : `btn btn-sm ${faded}`}
            disabled={disabled || retrying || !job.done}
            title={title}
            onClick={retry}>
            <i className={`ri-restart-line ${compact ? "text-base" : ""}`}></i>
            {!compact && "Retry"}
        </button>
    );
}

export default JobRetryButton;
