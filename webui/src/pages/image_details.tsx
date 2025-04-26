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

function ImageDetails() {
    const { img_uuid } = useParams();
    let [img, setImage] = useState<Image>({} as Image);
    let [job, setJob] = useState<Job>({} as Job);

    useEffect(() => {
        function getImage() {
            if (img_uuid === undefined) {
                console.log("Image uuid undefined, cannot get image details");
                return;
            }

            API.images.get(img_uuid).then((response: any) => {
                setImage(response.data);
            }).catch(error => {
                console.log("Cannot get image details", error);
                setImage({} as Image);
            });
        }

        getImage();
    }, [img_uuid])

    useEffect(() => {
        function getJob() {
            if (img_uuid === undefined) {
                console.log("Job uuid undefined, cannot get job details");
                return;
            }

            API.jobs.get(img_uuid).then((response: any) => {
                setJob(response.data);
            }).catch(error => {
                console.log("Cannot get job details", error);
                setJob({} as Job);
            });
        }

        getJob();
    }, [img_uuid])

    return (
        <>
            <Card>
                <div className="px-6 py-4 flex">
                    <h3 className="flex-grow font-bold">
                        General info
                    </h3>
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
                                    <td>{Utils.formatSize(img.size_on_disk)}</td>
                                </tr>
                                <tr>
                                    <td>Number of elements (total)</td>
                                    <td>{img.number_of_elements}</td>
                                </tr>
                                <tr>
                                    <td>Number of directories</td>
                                    <td>{img.number_of_folders}</td>
                                </tr>
                                <tr>
                                    <td>Number of files</td>
                                    <td>{img.number_of_files}</td>
                                </tr>
                            </tbody>
                        </table>
                    </Card>
                    <ClientCard client={img.client} />
                    <ModuleCard module={img.module} full />
                    <RepositoryCard repo={img.repository} />
                </div>
            </Card>


            <Card>
                <JobList
                    title="Generated from job"
                    actions={false}
                    data={[job]}
                    paginated={false}
                    sorted={false}
                />
            </Card>
        </>
    );
}

export default ImageDetails;