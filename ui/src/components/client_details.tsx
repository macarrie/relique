import React, {useEffect, useState} from "react";
import { useParams } from "react-router-dom";

import StatusDot from "./status_dot";
import API from "../utils/api";
import Client from "../types/client";
import Module from "../types/module";
import ModuleCard from "./module_card";
import ClientUtils from "../utils/client";
import Const from "../types/const";
import Card from "./card";

function ClientDetails() {
    const {client_name} = useParams();
    let [client, setClient] = useState<Client | null>(null);
    let [notFound, setNotFound] = useState<boolean>(false);
    let [err, setErr] = useState<boolean>(false);

    useEffect(() => {
        function getClient() {
            if (!client_name) {
                console.log("Cannot get client: no name provided")
                return;
            }

            API.clients.get(client_name).then((response: any) => {
                let cl = response.data;
                if (cl.modules === null) {
                    cl.modules = [];
                }

                setClient(cl);
            }).catch(error => {
                setErr(true)
                if (error.response.status === 404) {
                    setNotFound(true)
                }
                console.log("Cannot get client details", error);
            });
        }

        getClient();
    }, [client_name])

    function displayModules(mods :Module[]) {
        if (mods.length === 0) {
            return <div className={"text-slate-400 italic"}>None</div>;
        }

        return <>
            {mods.map((m :Module) => {
                return <ModuleCard key={m.name} module={m} />
            })}
        </>;

    }

    if (client === null) {
        if (err) {
            if (notFound) {
                return <div>Client not found</div>
            }
            return <div>Cannot load client</div>
        }
        return <div>Loading</div>
    }

    return (
        <Card>
            <div className="flex flex-row px-4 py-3 items-center">
                <span className="flex-grow text-xl font-bold">Client details</span>
            </div>

            <div className="grid md:grid-cols-2 gap-4 m-4">
                <Card className="bg-white bg-opacity-60">
                    <div className="p-4 flex flex-row items-center mb-2">
                        <div className="font-bold text-slate-500 dark:text-slate-200">General info</div>
                    </div>
                    <table className="details-table ml-4">
                        <tr>
                            <td>Name</td>
                            <td>{client.name}</td>
                        </tr>
                        <tr>
                            <td>Address</td>
                            <td className="code text-base">{client.address}</td>
                        </tr>
                        <tr>
                            <td>Port</td>
                            <td className="code text-base">{client.port}</td>
                        </tr>
                    </table>
                </Card>

                <Card className="bg-white bg-opacity-60">
                    <div className="p-4 flex flex-row items-center mb-2">
                        <div className="flex-grow font-bold text-slate-500 dark:text-slate-200">Health</div>
                        <div className={"flex flex-row items-center"}>
                            <div className={"mr-2"}>
                                <StatusDot status={ClientUtils.alive(client)}/>
                            </div>
                            <div className="text-xs">
                                {ClientUtils.GlobalAliveLabel(client)}
                            </div>
                        </div>
                    </div>
                    <table className="details-table ml-4">
                        <tr>
                            <td>API Status</td>
                            <td>
                                <div className={"flex flex-row items-center"}>
                                    <div className={"mr-2"}>
                                        <StatusDot status={client.api_alive}/>
                                    </div>
                                    <div>
                                        {ClientUtils.APIAliveLabel(client)}
                                    </div>
                                </div>
                            </td>
                        </tr>
                        <tr>
                            <td>SSH availability</td>
                            <td>
                                <div className={"flex flex-row items-center"}>
                                    <div className={"mr-2"}>
                                        <StatusDot status={client.ssh_alive}/>
                                    </div>
                                    <div>
                                        <div>
                                            {ClientUtils.SSHAliveLabel(client)}
                                        </div>
                                    </div>
                                </div>
                                {(client.ssh_alive !== Const.OK && client.ssh_alive_message) && (
                                    <div
                                        className="rounded border-l-2 border-red-200 bg-red-100 dark:border-red-900 dark:bg-red-900/50 dark:text-red-200 m-1 ml-5 py-1 px-2 mt-1 text-xs font-mono text-pink-900">
                                        {client.ssh_alive_message}
                                    </div>
                                )}
                            </td>
                        </tr>
                    </table>
                </Card>
            </div>

            <div className="flex flex-col px-4 pt-8 pb-4">
                <div className={"font-bold text-slate-500 dark:text-slate-200 mb-8"}>Modules</div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {displayModules(client.modules)}
                </div>
            </div>
        </Card>
    );
}

export default ClientDetails;
