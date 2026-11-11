import { useEffect, useState } from 'react';
import { useNavigate, useParams } from "react-router-dom";
import Image from '../types/image';
import API from '../utils/api';
import Utils from '../utils/utils';
import Job from '../types/job';
import JobList from '../components/job_list';
import Const from '../types/const';
import PageHeader from '../components/page_header';
import Section from '../components/section';
import PropertyList from '../components/property_list';
import ClientProperties from '../components/client_properties';
import ModuleProperties from '../components/module_properties';
import RepositoryProperties from '../components/repository_properties';
import ImageRemoveModal from '../components/image_remove_modal';

function ImageDetails() {
    const { img_uuid } = useParams();
    const [img, setImage] = useState<Image>({} as Image);
    const [job, setJob] = useState<Job>({} as Job);
    const [loading, setLoading] = useState<boolean>(true);
    const [jobLoading, setJobLoading] = useState<boolean>(true);
    const [removeOpen, setRemoveOpen] = useState<boolean>(false);

    const navigate = useNavigate();

    useEffect(() => {
        function getImage() {
            setLoading(true);
            if (img_uuid === undefined) {
                console.log("Image uuid undefined, cannot get image details");
                return;
            }

            API.images.get(img_uuid).then((response) => {
                setImage(response.data);
                setLoading(false);
            }).catch(error => {
                console.log("Cannot get image details", error);
                Utils.notify(Const.CRITICAL, "Cannot get image details", error.toString())
                setImage({} as Image);
            });
        }

        getImage();
    }, [img_uuid])

    useEffect(() => {
        function getJob() {
            setJobLoading(true);
            if (img_uuid === undefined) {
                console.log("Job uuid undefined, cannot get job details");
                return;
            }

            API.jobs.get(img_uuid).then((response) => {
                setJob(response.data);
                setJobLoading(false);
            }).catch(error => {
                setJob({} as Job);
                setJobLoading(false);
                // The job that generated the image can have been removed while keeping the image
                if (error.response?.status === 404) {
                    return;
                }
                console.log("Cannot get job details", error);
                Utils.notify(Const.CRITICAL, "Cannot get job details", error.toString())
            });
        }

        getJob();
    }, [img_uuid])

    return (
        <>
            <PageHeader
                title={<>Image <span className="code text-2xl font-normal text-muted">{Utils.shortID(img.uuid)}</span></>}
                loading={loading}
                meta={!loading && (
                    <>
                        <span className="tabular-nums">{Utils.formatDate(img.created_at)}</span>
                        <span className="code text-muted">{img.uuid}</span>
                    </>
                )}
                actions={
                    <button type="button" className="btn btn-sm btn-ghost text-error" disabled={loading} onClick={() => setRemoveOpen(true)}>
                        <i className="ri-delete-bin-line"></i>
                        Delete
                    </button>
                }
            />
            <ImageRemoveModal
                image={removeOpen ? img : null}
                onClose={() => setRemoveOpen(false)}
                onRemoved={() => navigate("/images")}
            />

            <div className="grid md:grid-cols-2 gap-4">
                <Section title="Contents">
                    <PropertyList loading={loading} items={[
                        { label: "Size on disk", value: <span className="tabular-nums">{Utils.formatSize(img.size_on_disk)}</span> },
                        { label: "Elements", value: <span className="tabular-nums">{img.number_of_elements}</span> },
                        { label: "Directories", value: <span className="tabular-nums">{img.number_of_folders}</span> },
                        { label: "Files", value: <span className="tabular-nums">{img.number_of_files}</span> },
                    ]} />
                </Section>
                <Section title="Client">
                    <ClientProperties client={img.client} loading={loading} link />
                </Section>
                <Section title={<>Module {!loading && <span className="font-normal text-muted">{img.module?.name}</span>}</>}>
                    <ModuleProperties module={img.module} loading={loading} />
                </Section>
                <Section title="Repository">
                    <RepositoryProperties repo={img.repository} loading={loading} link />
                </Section>
            </div>

            <Section title="Source job">
                <JobList
                    actions={false}
                    data={job?.uuid ? [job] : []}
                    paginated={false}
                    sorted={false}
                    loading={jobLoading}
                    deletable={true}
                    onDeleted={(_uuid: string, imageRemoved: boolean) => imageRemoved ? navigate("/images") : setJob({} as Job)}
                />
            </Section>
        </>
    );
}

export default ImageDetails;