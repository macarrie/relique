import { useCallback, useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from "react-router-dom";
import Status from '../components/status';
import Job from '../types/job';
import Image from '../types/image';
import API from '../utils/api';
import JobUtils from '../utils/job';
import Utils from '../utils/utils';
import Const from '../types/const';
import PageHeader from '../components/page_header';
import Section from '../components/section';
import PropertyList from '../components/property_list';
import ClientProperties from '../components/client_properties';
import ModuleProperties from '../components/module_properties';
import JobRemoveModal from '../components/job_remove_modal';
import JobRetryButton from '../components/job_retry_button';
import JobEventTimeline from '../components/job_event_timeline';
import JobEvent from '../types/job_event';
import JobLogs from '../components/job_logs';
import JobLog from '../types/job_log';

// Events of a running job are reloaded at this interval (ms)
const EVENTS_REFRESH_INTERVAL = 10000;

function JobDetails() {
    const { job_uuid } = useParams();
    const [j, setJob] = useState<Job>({} as Job);
    const [img, setImage] = useState<Image>({} as Image);
    const [loading, setLoading] = useState<boolean>(true);
    const [imgLoading, setImgLoading] = useState<boolean>(true);
    const [removeOpen, setRemoveOpen] = useState<boolean>(false);
    const [events, setEvents] = useState<JobEvent[]>([]);
    const [eventsLoading, setEventsLoading] = useState<boolean>(true);
    const [logs, setLogs] = useState<JobLog[]>([]);
    const [logsLoading, setLogsLoading] = useState<boolean>(true);

    const navigate = useNavigate();

    const getJob = useCallback(() => {
        if (job_uuid === undefined) {
            console.log("Job uuid undefined, cannot get job details");
            return;
        }

        API.jobs.get(job_uuid).then((response) => {
            setJob(response.data);
            setLoading(false);
        }).catch(error => {
            console.log("Cannot get job details", error);
            Utils.notify(Const.CRITICAL, "Cannot get job details", error.toString())
            setJob({} as Job);
        });
    }, [job_uuid]);

    useEffect(() => {
        function initJob() {
            setLoading(true);
            getJob();
        }

        initJob();
    }, [getJob])

    // Resolves with the loaded events, or an empty list when they cannot be loaded
    const loadEvents = useCallback((): Promise<JobEvent[]> => {
        if (job_uuid === undefined) {
            return Promise.resolve([]);
        }

        return API.jobs.events(job_uuid).then((response) => {
            const loaded: JobEvent[] = response.data ?? [];
            setEvents(loaded);
            setEventsLoading(false);
            return loaded;
        }).catch(error => {
            setEventsLoading(false);
            // Unknown job: the job details request already reports it
            if (error.response?.status !== 404) {
                console.log("Cannot get job events", error);
                Utils.notify(Const.CRITICAL, "Job events loading", "Cannot get job events: " + error)
            }
            return [];
        });
    }, [job_uuid]);

    useEffect(() => {
        function initEvents() {
            setEventsLoading(true);
            setEvents([]);
            loadEvents();
        }

        initEvents();
    }, [loadEvents])

    // While the job runs, only its events are reloaded. Once an event marks the end of the job,
    // the job is reloaded once to update its status and actions, which stops the refresh
    useEffect(() => {
        if (loading || j.done) {
            return;
        }

        const timer = setInterval(() => {
            loadEvents().then((loaded) => {
                if (loaded.some((e) => JobUtils.isTerminalEvent(e.type))) {
                    getJob();
                }
            });
        }, EVENTS_REFRESH_INTERVAL);

        return () => clearInterval(timer);
    }, [loading, j.done, loadEvents, getJob])

    // Logs are written when each sync ends: they are loaded again once the job is done
    const loadLogs = useCallback(() => {
        if (job_uuid === undefined) {
            return;
        }

        setLogsLoading(true);
        API.jobs.logs(job_uuid).then((response) => {
            setLogs(response.data ?? []);
            setLogsLoading(false);
        }).catch(error => {
            setLogs([]);
            setLogsLoading(false);
            // Unknown job: the job details request already reports it
            if (error.response?.status !== 404) {
                console.log("Cannot get job logs", error);
                Utils.notify(Const.CRITICAL, "Job logs loading", "Cannot get job logs: " + error)
            }
        });
    }, [job_uuid]);

    useEffect(() => {
        function initLogs() {
            loadLogs();
        }

        initLogs();
    }, [loadLogs, j.done])

    useEffect(() => {
        function getImage() {
            setImgLoading(true);
            if (job_uuid === undefined) {
                console.log("Image uuid undefined, cannot get image details");
                return;
            }

            API.images.get(job_uuid).then((response) => {
                setImage(response.data);
                setImgLoading(false);
            }).catch(error => {
                setImage({} as Image);
                setImgLoading(false);
                // Restore jobs never have an image, and the image of a backup job can have been removed
                if (error.response?.status === 404) {
                    return;
                }
                console.log("Cannot get image details", error);
                Utils.notify(Const.CRITICAL, "Image details loading", "Cannot get image details: " + error)
            });
        }

        getImage();
    }, [job_uuid])

    const jobType = j.job_type ? j.job_type.charAt(0).toUpperCase() + j.job_type.slice(1) : "";
    const duration = Utils.formatDuration(j.start_time, j.end_time);

    return (
        <>
            <PageHeader
                title={<>{jobType} job <span className="code text-2xl font-normal text-muted">{Utils.shortID(j.uuid)}</span></>}
                loading={loading}
                meta={!loading && (
                    <>
                        <Status label={j.status} status={JobUtils.jobStateToCode(j.status)} />
                        <span className="code text-muted">{j.uuid}</span>
                    </>
                )}
                actions={
                    <>
                        <JobRetryButton key={j.uuid} job={j} disabled={loading} />
                        <button type="button" className="btn btn-sm btn-ghost text-error"
                            disabled={loading || !j.done}
                            title={j.done ? "Delete job" : "Job is still running and cannot be deleted"}
                            onClick={() => setRemoveOpen(true)}>
                            <i className="ri-delete-bin-line"></i>
                            Delete
                        </button>
                    </>
                }
            />
            <JobRemoveModal
                job={removeOpen ? j : null}
                onClose={() => setRemoveOpen(false)}
                onRemoved={() => navigate("/jobs")}
            />

            <div className="grid md:grid-cols-2 gap-4">
                <Section title="Execution">
                    <PropertyList loading={loading} items={[
                        { label: "Job type", value: j.job_type },
                        { label: "Backup type", value: j.backup_type },
                        { label: "Started", value: <span className="tabular-nums">{Utils.formatDate(j.start_time)}</span> },
                        { label: "Ended", value: j.done ? <span className="tabular-nums">{Utils.formatDate(j.end_time)}</span> : <span className="text-info font-medium">Running</span> },
                        { label: "Duration", value: j.done && duration ? <span className="tabular-nums">{duration}</span> : undefined },
                        { label: "Repository", value: j.repository?.name && <Link className="text-link" to={`/repositories/${j.repository.name}`}>{j.repository.name}</Link> },
                        { label: "Image", value: !imgLoading && img?.uuid ? <Link className="code text-link" to={`/images/${img.uuid}`}>{Utils.shortID(img.uuid)}</Link> : undefined },
                    ]} />
                </Section>

                <Section title="Client">
                    <ClientProperties client={j.client} loading={loading} link />
                </Section>
            </div>

            <Section title={<>Module {!loading && <span className="font-normal text-muted">{j.module?.name}</span>}</>}>
                <ModuleProperties module={j.module} loading={loading} />
            </Section>

            <Section title="Events" description={!loading && !j.done ? "Refreshed every 10 seconds while the job runs" : undefined}>
                <JobEventTimeline events={events} jobStart={j.start_time} loading={eventsLoading} />
            </Section>

            <Section title="Logs"
                description="rsync output of each synced path"
                actions={
                    <button type="button" className="btn btn-sm" disabled={logsLoading} onClick={loadLogs}>
                        <i className="ri-refresh-line"></i>
                        Reload
                    </button>
                }>
                <JobLogs logs={logs} loading={logsLoading} />
            </Section>
        </>
    );
}

export default JobDetails;