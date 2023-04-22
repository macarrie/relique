import React, {useEffect, useState} from "react";

import API from "../utils/api";
import ClientUtils from "../utils/client";

import Client from "../types/client";
import Module from "../types/module";

import StatusDot from "./status_dot";
import {Link} from "react-router-dom";

function ClientListRowPlaceholder() {
    return (
        <tr className="animate-pulse">
            <td className="py-2 px-3"><div className="rounded-full h-3 w-3 m-auto bg-slate-300 dark:bg-slate-600"></div></td>
            <td className="py-2 px-3"><div className="rounded-full h-2 w-1/2 bg-slate-300 dark:bg-slate-600"></div></td>
            <td className="py-2 px-3 code"><div className="rounded-full h-2 w-4/5 bg-slate-300 dark:bg-slate-600"></div></td>
            <td className="py-2 px-3 hidden md:table-cell">
                <div className="flex flex-row space-x-1">
                    <div className="rounded-full h-2 w-12 bg-slate-300 dark:bg-slate-600"></div>
                    <div className="rounded-full h-2 w-12 bg-slate-300 dark:bg-slate-600"></div>
                    <div className="rounded-full h-2 w-12 bg-slate-300 dark:bg-slate-600"></div>
                </div>
            </td>
        </tr>
    );
}

function ClientListRow(props :any) {
    let client = props.client;

    function renderModules() {
        if (!client.modules) {
            return <span className="italic text-slate-400 dark:text-slate-600">None</span>;
        }

        let module_names :string[] = client.modules.map((mod :Module) => mod.name)
        return (
            <>
                {module_names.map((mod: any) => (
                    <span className="badge" key={mod}>{mod}</span>
                ))}
            </>
        )
    }

    return (
        <tr>
            <td className="py-2 px-3"><StatusDot status={ClientUtils.alive(client)}/></td>
            <td className="py-2 px-3"><Link to={`/clients/${client.name}`}>{client.name}</Link></td>
            <td className="py-2 px-3 code">{client.address}</td>
            <td className="py-2 px-3 space-x-1 hidden md:table-cell">{renderModules()}</td>
        </tr>
    );
}

function ClientList(props :any) {
    let [limit, setLimit] = useState(props.limit || 0);
    let [clients, setClients] = useState([] as Client[]);
    let [loading, setLoading] = useState(true);

    useEffect(() => {
        setLimit(props.limit);
    }, [props.limit])

    useEffect(() => {
        function getClients() {
            API.clients.list({
                limit: limit,
            }).then((response :any) => {
                setClients(response.data);
                setLoading(false);
            }).catch(error => {
                setLoading(false);
                console.log("Cannot get client list", error);
            });
        }

        getClients();
    }, [limit]);

    function renderClientList() {
        if (loading) {
            return (
                <tbody>
                    <ClientListRowPlaceholder />
                    <ClientListRowPlaceholder />
                    <ClientListRowPlaceholder />
                </tbody>
            )
        }

        if (!clients || clients.length === 0) {
            return (
                <tr>
                    <td colSpan={4} className={"px-3 py-8 text-center text-3xl italic text-gray-300 dark:text-gray-600"}>
                        No clients
                    </td>
                </tr>
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
