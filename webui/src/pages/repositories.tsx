import { useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import Repository from "../types/repository";
import RepoList from "../components/repo_list";
import Const from "../types/const";
import Utils from "../utils/utils";

function Repositories() {
    let [repos, setRepos] = useState<Repository[]>([]);
    let [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        function getRepoList() {
            setLoading(true);
            API.repos.list({ limit: 10000 }).then((response: any) => {
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
            <Card>
                <RepoList
                    title="All repositories"
                    actions={true}
                    data={repos}
                    paginated={true}
                    sorted={true}
                    loading={loading}
                />
            </Card>
        </>
    );
}

export default Repositories;