import Card from "./card";
import StatusBadge from "./status_badge";
import Const from "../types/const";
import Repository from "../types/repository";
import TextPlaceholder from "./text_placeholder";

function RepositoryCard(props: any) {
    let repo: Repository = props.repo;

    return (
        <Card className={props.className}>
            <div className="p-4 flex flex-row items-center mb-2">
                <div className="flex-grow font-bold align-middle">Repository
                    {props.loading ? (
                        <span className="ml-2 badge skeleton w-16"></span>
                    ) : (
                        <span className='ml-2 badge badge-neutral badge-soft font-normal'>{repo?.name}</span>
                    )}
                </div>
            </div>
            <table className="table">
                <tr>
                    <td>Repository name</td>
                    <td>
                        {props.loading && (
                            <TextPlaceholder />
                        )}
                        {repo?.name}
                    </td>
                </tr>
                <tr>
                    <td>Type</td>
                    <td>
                        {props.loading && (
                            <TextPlaceholder />
                        )}
                        {repo?.type}
                    </td>
                </tr>
                <tr>
                    <td>Storage path</td>
                    <td className="code">
                        {props.loading && (
                            <TextPlaceholder />
                        )}
                        {repo?.path}
                    </td>
                </tr>
                <tr>
                    <td>Use as default</td>
                    <td>
                        <StatusBadge loading={props.loading} label={repo?.default ? "true" : "false"} status={repo?.default ? Const.OK : Const.CRITICAL} />
                    </td>
                </tr>
            </table>
        </Card>
    );
}

export default RepositoryCard;