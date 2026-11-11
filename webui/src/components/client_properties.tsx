import { Link } from "react-router-dom";
import Client from "../types/client";
import PropertyList, { Property } from "./property_list";
import ClientUtils from "../utils/client";

function ClientProperties({
    client,
    loading = false,
    link = false,
}: {
    client: Client,
    loading?: boolean,
    link?: boolean,
}) {
    const items: Property[] = [
        { label: "Name", value: link && client?.name ? <Link className="text-link" to={`/clients/${client.name}`}>{client.name}</Link> : client?.name },
        { label: "Address", value: client?.address, mono: true },
        { label: "SSH user", value: ClientUtils.sshUser(client), mono: true },
        { label: "SSH port", value: ClientUtils.sshPort(client), mono: true },
    ];

    return <PropertyList items={items} loading={loading} />;
}

export default ClientProperties;
