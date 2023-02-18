import React, {useEffect, useState} from "react";

import API from "../utils/api";
import ClientUtils from "../utils/client";

import Client from "../types/client";
import Module from "../types/module";

import StatusDot from "./status_dot";
import {Link} from "react-router-dom";

function ClientListRow(props :any) {
    let client = props.client;
    let module_names :string[] = client.modules.map((mod :Module) => mod.name)

    return (
        <tr>
            <td className="py-2 px-3"><StatusDot status={ClientUtils.alive(client)}/></td>
            <td className="py-2 px-3"><Link to={`/clients/${client.name}`}>{client.name}</Link></td>
            <td className="py-2 px-3 code">{client.address}</td>
            <td className="py-2 px-3 space-x-1 hidden md:table-cell">{module_names.map((mod: any) => (
                <span className="badge">{mod}</span>))}</td>
        </tr>
    );
}

function ClientList(props :any) {
    let [limit, setLimit] = useState(props.limit || 0);
    let [clients, setClients] = useState([] as Client[]);

    useEffect(() => {
        setLimit(props.limit);
    }, [props.limit])

    useEffect(() => {
        function getClients() {
            API.clients.list({
                limit: limit,
            }).then((response :any) => {
                setClients(response.data);
            }).catch(error => {
                console.log("Cannot get client list", error);
            });
        }

        getClients();
    }, [limit]);

    function renderClientList() {
        if (!clients) {
            return (
                <>
                    Loading
                </>
            )
        }

        const clientList = clients.map((client :Client) =>
            <ClientListRow key={client.id} client={client}/>
        );

        return (
            <tbody>
            {clientList}
            </tbody>
        )
    }

    return (
        <table className="table table-auto w-full">
            <thead>
            <tr>
                <th className="py-2 px-3 max-w-min text-center">Health</th>
                <th className="py-2 px-3">Name</th>
                <th className="py-2 px-3">Address</th>
                <th className="py-2 px-3 hidden md:table-cell">Modules</th>
            </tr>
            </thead>
            {renderClientList()}
        </table>
    );
}

export default ClientList;
