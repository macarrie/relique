import React, {useEffect, useState} from "react";

import API from "../utils/api";
import ClientUtils from "../utils/client";

import Client from "../types/client";
import Module from "../types/module";
import Const from "../types/const";

import DotStatus from "./dot_status";
import {Link} from "react-router-dom";

function ClientListRow(props :any) {
    let client = props.client;
    let module_names :string[] = client.modules.map((mod :Module) => mod.name)

    return (
        <tr className="hover:bg-slate-50">
            <td className="py-2 px-3"><DotStatus status={ClientUtils.alive(client)} /></td>
            <td className="py-2 px-3"><Link to={`/clients/${client.id}`}>{client.name}</Link></td>
            <td className="py-2 px-3">{client.address}</td>
            <td className="py-2 px-3">{module_names.join(", ")}</td>
        </tr>
    );
}

type State = {
    clients :Client[]
};

function ClientList(props :any) {
    let [limit, setLimit] = useState(props.limit || 0);
    let [clients, setClients] = useState([] as Client[]);

    useEffect(() => {
        function getClients() {
            API.clients.list({
                limit: limit,
            }).then((response :any) => {
                setClients(response.data);
            }).catch(error => {
                console.log("Cannot get client list", error);
            }).finally(() => {
                clients.map((client :Client, index :number) => {
                    sshPing(index);
                });
            });
        }

        function sshPing(index :number) {
            let c = clients[index];
            API.clients.ssh_ping(c.id).then((response :any) => {
                let clientList = clients;
                clientList[index].ssh_alive = Const.OK;
                setClients(clientList);
            }).catch(error => {
                let status :number
                switch (error.response.status) {
                    case 404:
                        status = Const.UNKNOWN;
                        break;
                    case 401:
                        status = Const.CRITICAL;
                        break;
                    default:
                        status = Const.CRITICAL;
                        console.log("Error when getting SSH ping status", error);
                        break;
                }

                let clientList = clients;
                clientList[index].ssh_alive = status;
                setClients(clientList);
            });
        }

        getClients();
    }, []);

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
        <table className="table-auto w-full">
            <thead className="bg-slate-50 uppercase text-slate-500 text-left">
            <tr className="border border-l-0 border-r-0 border-slate-100">
                <th className="py-2 px-3 max-w-min text-center">Health</th>
                <th className="py-2 px-3">Name</th>
                <th className="py-2 px-3">Address</th>
                <th className="py-2 px-3">Modules</th>
            </tr>
            </thead>
            {renderClientList()}
        </table>
    );
}

export default ClientList;
