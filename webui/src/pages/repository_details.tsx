import { useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import API from '../utils/api';
import Repository from '../types/repository';
import Const from '../types/const';
import Utils from '../utils/utils';
import PageHeader from '../components/page_header';
import Section from '../components/section';
import RepositoryProperties from '../components/repository_properties';

function RepositoryDetails() {
    const { repo_name } = useParams();
    const [repo, setRepo] = useState<Repository>({} as Repository);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getModule() {
            setLoading(true);
            if (repo_name === undefined) {
                console.log("Repository name undefined, cannot get repo details");
                return;
            }

            API.repos.get(repo_name).then((response) => {
                setRepo(response.data);
                setLoading(false);
            }).catch(error => {
                console.log("Cannot get repository details", error);
                Utils.notify(Const.CRITICAL, "Cannot get repo details", error.toString())
                setRepo({} as Repository);
            });
        }

        getModule();
    }, [repo_name])

    return (
        <>
            <PageHeader
                title={<><span className="font-normal text-muted">Repository</span> {repo.name}</>}
                loading={loading}
                meta={!loading && <span>{repo.type} repository{repo.default && ", used by default"}</span>}
            />
            <Section title="Settings">
                <RepositoryProperties loading={loading} repo={repo} />
            </Section>
        </>
    );
}

export default RepositoryDetails;