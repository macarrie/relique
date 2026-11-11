import { useEffect, useState } from 'react';
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

function JobDetails() {
    const { job_uuid } = useParams();
    const [j, setJob] = useState<Job>({} as Job);
    const [img, setImage] = useState<Image>({} as Image);
    const [loading, setLoading] = useState<boolean>(true);
    const [imgLoading, setImgLoading] = useState<boolean>(true);
    const [removeOpen, setRemoveOpen] = useState<boolean>(false);

    const navigate = useNavigate();

    useEffect(() => {
        function getJob() {
            setLoading(true);
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
        }

        getJob();
    }, [job_uuid])

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
                    <button type="button" className="btn btn-sm btn-ghost text-error"
                        disabled={loading || !j.done}
                        title={j.done ? "Delete job" : "Job is still running and cannot be deleted"}
                        onClick={() => setRemoveOpen(true)}>
                        <i className="ri-delete-bin-line"></i>
                        Delete
                    </button>
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
        </>
    );
}

export default JobDetails;