import React, {useEffect, useState} from "react";
import { useParams } from "react-router-dom";

import DotStatus from "./dot_status";
import Const from "../types/const";
import API from "../utils/api";
import Client from "../types/client";
import Module from "../types/module";
import ModuleCard from "./module_card";
import ClientUtils from "../utils/client";

function ClientDetails(props :any) {
    let urlParams = useParams();
    let [client, setClient] = useState<Client | null>(null);

    useEffect(() => {
        function getClient() {
            API.clients.list(urlParams.client_id).then((response :any) => {
                setClient(response.data[0]);
                sshPing();
            }).catch(error => {
                console.log("Cannot get client details", error);
            });
        }

        function sshPing() {
            if (client === null) {
                return;
            }

            API.clients.ssh_ping(client.id).then((response :any) => {
                // @ts-ignore
                client.ssh_alive = Const.OK;
                setClient(client);
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

                // @ts-ignore
                client.ssh_alive = status;
                setClient(client);
            });
        }

        getClient();
    }, [])

    function displayModules(mods :Module[]) {
        if (mods.length === 0) {
            return <div className={"ml-3 text-slate-400 italic"}>None</div>;
        }

        return <>
            {mods.map((m :Module) => {
                return <ModuleCard key={m.name} module={m} />
            })}
        </>;

    }

    if (client === null) {
        return <div>Loading</div>
    }

    return (
        <div className="grid grid-cols-4 gap-4">
            <div className="col-span-4 bg-white shadow rounded">
                <div className="flex flex-row px-4 py-3 items-center">
                    <span className="flex-grow text-xl font-bold">Client details</span>
                </div>

                <div className="flex flex-col px-4 py-3 pb-4 bg-slate-50 space-y-3">
                    <div className={"uppercase font-bold text-slate-500 mb-2"}>General info</div>
                    <div className="flex flex-row content-center">
                        <div className={"w-24 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Name</div>
                        <div className={"flex-grow bg-white py-2 px-3 ml-6 rounded shadow-sm text-slate-900"}>{client.name}</div>
                    </div>
                    <div className="flex flex-row">
                        <div className={"w-24 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Address</div>
                        <div className={"flex-grow bg-white py-2 px-3 ml-6 rounded shadow-sm text-slate-900"}>{client.address}</div>
                    </div>
                    <div className="flex flex-row">
                        <div className={"w-24 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Port</div>
                        <div className={"flex-grow bg-white py-2 px-3 ml-6 rounded shadow-sm text-slate-900"}>{client.port}</div>
                    </div>
                </div>

                <hr />

                <div className="flex flex-col px-4 py-3 bg-slate-50 space-y-3">
                    <div className={"uppercase font-bold text-slate-500 mb-2"}>Health</div>
                    <div className="flex flex-row items-center">
                        <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>Global health</div>
                        <div className={"flex flex-row items-center py-2 px-3 ml-3"}>
                            <div className={"mr-2"}>
                                <DotStatus status={ClientUtils.alive(client)} />
                            </div>
                            <div>
                                {ClientUtils.GlobalAliveLabel(client)}
                            </div>
                        </div>
                    </div>
                    <div className="pl-8 flex flex-row items-center">
                        <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>API status</div>
                        <div className={"flex flex-row items-center py-2 px-3 ml-3"}>
                            <div className={"mr-2"}>
                                <DotStatus status={client.api_alive} />
                            </div>
                            <div>
                                {ClientUtils.APIAliveLabel(client)}
                            </div>
                        </div>
                    </div>
                    <div className="pl-8 flex flex-row items-center">
                        <div className={"w-48 py-2 px-3 font-bold text-sm text-slate-400 uppercase"}>SSH availability</div>
                        <div className={"flex flex-row items-center py-2 px-3 ml-3"}>
                            <div className={"mr-2"}>
                                <DotStatus status={client.ssh_alive} />
                            </div>
                            <div>
                                {ClientUtils.SSHAliveLabel(client)}
                            </div>
                        </div>
                    </div>
                </div>

                <hr />

                <div className="flex flex-col px-4 py-3 pb-4 bg-slate-50 space-y-3">
                    <div className={"uppercase font-bold text-slate-500 mb-2"}>Modules</div>
                    <div className="grid grid-cols-2">
                        {displayModules(client.modules)}
                    </div>
                </div>
            </div>
        </div>
    );
}

export default ClientDetails;
