import { useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import Repository from "../types/repository";
import RepoList from "../components/repo_list";

function Repositories() {
    let [repos, setRepos] = useState<Repository[]>([]);

    useEffect(() => {
        function getRepoList() {
            API.repos.list({ limit: 10000 }).then((response: any) => {
                setRepos(response.data.data ?? []);
            }).catch(error => {
                console.log("Cannot get repo list", error);
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
                />
            </Card>
        </>
    );
}

export default Repositories;