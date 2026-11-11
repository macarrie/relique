import { Link } from "react-router-dom";
import Repository from "../types/repository";
import PropertyList, { Property } from "./property_list";
import YesNo from "./yes_no";

function RepositoryProperties({
    repo,
    loading = false,
    link = false,
}: {
    repo: Repository,
    loading?: boolean,
    link?: boolean,
}) {
    const items: Property[] = [
        { label: "Name", value: link && repo?.name ? <Link className="text-link" to={`/repositories/${repo.name}`}>{repo.name}</Link> : repo?.name },
        { label: "Type", value: repo?.type },
        { label: "Storage path", value: repo?.path, mono: true },
        { label: "Default", value: <YesNo value={repo?.default} /> },
    ];

    return <PropertyList items={items} loading={loading} />;
}

export default RepositoryProperties;
