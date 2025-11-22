import { useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import Card from '../components/card';
import Image from '../types/image';
import API from '../utils/api';
import ModuleCard from '../components/module_card';
import Utils from '../utils/utils';
import ClientCard from '../components/client_card';
import RepositoryCard from '../components/repository_card';
import Job from '../types/job';
import JobList from '../components/job_list';
import Const from '../types/const';
import TextPlaceholder from '../components/text_placeholder';

function ImageDetails() {
    const { img_uuid } = useParams();
    let [img, setImage] = useState<Image>({} as Image);
    let [job, setJob] = useState<Job>({} as Job);
    let [loading, setLoading] = useState<boolean>(true);
    let [jobLoading, setJobLoading] = useState<boolean>(true);

    useEffect(() => {
        function getImage() {
            setLoading(true);
            if (img_uuid === undefined) {
                console.log("Image uuid undefined, cannot get image details");
                return;
            }

            API.images.get(img_uuid).then((response: any) => {
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

            API.jobs.get(img_uuid).then((response: any) => {
                setJob(response.data);
                setJobLoading(false);
            }).catch(error => {
                console.log("Cannot get job details", error);
                Utils.notify(Const.CRITICAL, "Cannot get job details", error.toString())
                setJob({} as Job);
            });
        }

        getJob();
    }, [img_uuid])

    return (
        <>
            <Card>
                <div className="px-6 py-4 flex items-center">
                    <h3 className="flex-grow font-bold">
                        General info
                    </h3>
                    {loading && (
                        <div className='flex-grow'>
                            <TextPlaceholder />
                        </div>
                    )}
                    <span className="text-l ml-4 code">
                        {img.uuid}
                    </span>
                </div>

                <div className="grid md:grid-cols-2 gap-4 m-4">
                    <Card>
                        <div className="p-4 flex flex-row items-center mb-2">
                            <div className="font-bold">Image stats</div>
                        </div>
                        <table className="table">
                            <tbody>
                                <tr>
                                    <td>Size on disk</td>
                                    <td>
                                        {loading ? (
                                            <TextPlaceholder />
                                        ) : (
                                            <>
                                                {Utils.formatSize(img.size_on_disk)}
                                            </>
                                        )}
                                    </td>
                                </tr>
                                <tr>
                                    <td>Number of elements (total)</td>
                                    <td>
                                        {loading && (
                                            <TextPlaceholder />
                                        )}
                                        {img.number_of_elements}
                                    </td>
                                </tr>
                                <tr>
                                    <td>Number of directories</td>
                                    <td>
                                        {loading && (
                                            <TextPlaceholder />
                                        )}
                                        {img.number_of_folders}
                                    </td>
                                </tr>
                                <tr>
                                    <td>Number of files</td>
                                    <td>
                                        {loading && (
                                            <TextPlaceholder />
                                        )}
                                        {img.number_of_files}
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </Card>
                    <ClientCard client={img.client} loading={loading} />
                    <ModuleCard module={img.module} loading={loading} full />
                    <RepositoryCard repo={img.repository} loading={loading} />
                </div>
            </Card>


            <Card>
                <JobList
                    title="Generated from job"
                    actions={false}
                    data={[job]}
                    paginated={false}
                    sorted={false}
                    loading={jobLoading}
                />
            </Card>
        </>
    );
}

export default ImageDetails;