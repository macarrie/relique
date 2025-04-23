import Card from "./card";
import StatusBadge from "./status_badge";
import Const from "../types/const";
import Repository from "../types/repository";

function RepositoryCard(props: any) {
    let repo: Repository = props.repo;

    if (repo === null || repo === undefined) {
        return <div>Loading</div>
    }

    return (
        <Card className={props.className}>
            <div className="p-4 flex flex-row items-center mb-2">
                <div className="flex-grow font-bold align-middle">Repository <span className="ml-1 badge badge-neutral badge-soft font-normal">{repo.name}</span></div>
            </div>
            <table className="table">
                <tr>
                    <td>Repository name</td>
                    <td>{repo.name}</td>
                </tr>
                <tr>
                    <td>Type</td>
                    <td>{repo.type}</td>
                </tr>
                <tr>
                    <td>Storage path</td>
                    <td className="code">{repo.path}</td>
                </tr>
                <tr>
                    <td>Use as default</td>
                    <td><StatusBadge label={repo.default ? "true" : "false"} status={repo.default ? Const.OK : Const.CRITICAL} /></td>
                </tr>
            </table>
        </Card>
    );
}

export default RepositoryCard;