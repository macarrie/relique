import { useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import Card from '../components/card';
import API from '../utils/api';
import Repository from '../types/repository';
import RepositoryCard from '../components/repository_card';

function RepositoryDetails() {
    const { repo_name } = useParams();
    let [repo, setRepo] = useState<Repository>({} as Repository);

    useEffect(() => {
        function getModule() {
            if (repo_name === undefined) {
                console.log("Repository name undefined, cannot get repo details");
                return;
            }

            API.repos.get(repo_name).then((response: any) => {
                setRepo(response.data);
            }).catch(error => {
                console.log("Cannot get job details", error);
                setRepo({} as Repository);
            });
        }

        getModule();
    }, [repo_name])

    return (
        <>
            <Card>
                <div className="px-6 py-4 flex">
                    <h3 className="flex-grow font-bold">
                        Repository details
                    </h3>
                </div>
                <div className="m-4">
                    <RepositoryCard repo={repo} />
                </div>
            </Card>
        </>
    );
}

export default RepositoryDetails;