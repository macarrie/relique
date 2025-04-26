import { useEffect, useState } from 'react';
import { Link, useParams } from "react-router-dom";
import Card from '../components/card';
import StatusBadge from '../components/status_badge';
import Job from '../types/job';
import Image from '../types/image';
import API from '../utils/api';
import JobUtils from '../utils/job';
import ModuleCard from '../components/module_card';
import Utils from '../utils/utils';
import ClientCard from '../components/client_card';
import ImageList from '../components/image_list';

function JobDetails() {
    const { job_uuid } = useParams();
    let [j, setJob] = useState<Job>({} as Job);
    let [img, setImage] = useState<Image>({} as Image);

    useEffect(() => {
        function getJob() {
            if (job_uuid === undefined) {
                console.log("Job uuid undefined, cannot get job details");
                return;
            }

            API.jobs.get(job_uuid).then((response: any) => {
                setJob(response.data);
            }).catch(error => {
                console.log("Cannot get job details", error);
                setJob({} as Job);
            });
        }

        getJob();
    }, [job_uuid])

    useEffect(() => {
        function getImage() {
            if (job_uuid === undefined) {
                console.log("Image uuid undefined, cannot get image details");
                return;
            }

            API.images.get(job_uuid).then((response: any) => {
                setImage(response.data);
            }).catch(error => {
                console.log("Cannot get image details", error);
                setImage({} as Image);
            });
        }

        getImage();
    }, [job_uuid])

    return (
        <>
            <Card>
                <div className="px-6 py-4 flex">
                    <h3 className="flex-grow font-bold">
                        General info
                    </h3>
                    <span className="text-l ml-4 code">
                        {j.uuid}
                    </span>
                </div>

                <div className="grid md:grid-cols-2 gap-4 m-4">
                    <Card>
                        <div className="p-4 flex flex-row items-center mb-2">
                            <div className="font-bold">Job settings</div>
                        </div>
                        <table className="table">
                            <tbody>
                                <tr>
                                    <td>Type</td>
                                    <td>{j.job_type}</td>
                                </tr>
                                <tr>
                                    <td>Backup type</td>
                                    <td>{j.backup_type}</td>
                                </tr>
                                <tr>
                                    <td>Running</td>
                                    <td>{j.done ? "Finished" : "In execution"}</td>
                                </tr>
                                <tr>
                                    <td>Status</td>
                                    <td><StatusBadge label={j.status} status={JobUtils.jobStateToCode(j.status)} /></td>
                                </tr>
                                <tr>
                                    <td>Start time</td>
                                    <td>{Utils.formatDate(j.start_time)}</td>
                                </tr>
                                <tr>
                                    <td>End time</td>
                                    <td>{Utils.formatDate(j.end_time)}</td>
                                </tr>
                                <tr>
                                    <td>Storage repository</td>
                                    <td><Link to={`/repo/${j.repository?.name}`}>{j.repository?.name}</Link></td>
                                </tr>
                            </tbody>
                        </table>
                    </Card>
                    <ClientCard client={j.client} />
                    <ModuleCard className="col-span-2" module={j.module} full />
                </div>
            </Card>

            <Card>
                <ImageList
                    title="Generated image"
                    actions={false}
                    data={img?.uuid ? [img] : []}
                    paginated={false}
                    sorted={false}
                />
            </Card>

            <Card>
                <div className="px-6 py-4 flex">
                    <h3 className="flex-grow font-bold">
                        Timeline
                    </h3>
                </div>

                <div className='ml-6 mt-4'>
                    <ul className="timeline timeline-vertical timeline-snap-icon timeline-compact">
                        <li>
                            <div className="timeline-middle w-8 h-8 text-center bg-base-300 rounded-full">
                                <i className="ri-play-fill text-xl align-middle"></i>
                            </div>
                            <div className="timeline-end ml-4 mb-10 pt-2">
                                <h4 className="font-bold text-lg">Job started</h4>
                                <time className="text-base-content/50 italic text-sm">{Utils.formatDate(j.start_time)}</time>
                                <p className="leading-none mt-2">
                                    Backup job starting for client
                                    <Link className="link-primary" to={`/clients/${j.client?.name}`}>
                                        <span className='badge badge-soft badge-neutral badge-sm mx-1'>
                                            {j.client?.name}
                                        </span>
                                    </Link>
                                    and module
                                    <span className="badge badge-soft badge-neutral badge-sm mx-1">{j.module?.name}</span></p>
                            </div>
                            <hr />
                        </li>
                        <li>
                            <hr />
                            <div className="timeline-middle w-8 h-8 text-center bg-base-300 rounded-full">
                                <i className="ri-stop-fill text-xl align-middle"></i>
                            </div>
                            <div className="timeline-end ml-4 mb-10 pt-2">
                                <h4 className="font-bold text-lg">Job ended</h4>
                                <time className="text-base-content/50 italic text-sm">{Utils.formatDate(j.end_time)}</time>
                                <p className="leading-none mt-2">
                                    Backup job ended with status <StatusBadge status={JobUtils.jobStateToCode(j.status)} label={j.status} />
                                </p>
                            </div>
                        </li>
                    </ul>
                </div>
            </Card >
        </>
    );
}

export default JobDetails;