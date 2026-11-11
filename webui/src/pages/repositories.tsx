import { useEffect, useState } from "react";

import PageHeader from "../components/page_header";
import Section from "../components/section";
import API from "../utils/api";
import Repository from "../types/repository";
import RepoList from "../components/repo_list";
import Const from "../types/const";
import Utils from "../utils/utils";

function Repositories() {
    const [repos, setRepos] = useState<Repository[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getRepoList() {
            setLoading(true);
            API.repos.list({ limit: 10000 }).then((response) => {
                setRepos(response.data.data ?? []);
                setLoading(false);
            }).catch(error => {
                console.log("Cannot get repo list", error);
                Utils.notify(Const.CRITICAL, "Cannot get repo list", error.toString())
                setRepos([]);
            });
        }

        getRepoList();
    }, [])

    return (
        <>
            <PageHeader title="Repositories" />
            <Section>
                <RepoList
                    actions={true}
                    data={repos}
                    paginated={true}
                    sorted={true}
                    loading={loading}
                />
            </Section>
        </>
    );
}

export default Repositories;