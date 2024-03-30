import React, {useCallback, useEffect, useState} from "react";
import {Link} from "react-router-dom";

import API from "../utils/api";
import ClientUtils from "../utils/client";

import Client from "../types/client";
import Module from "../types/module";

import StatusDot from "./status_dot";
import Table from "./table";
import Dropdown from "./dropdown";
import Const from "../types/const";
import Loader from "./loader";
import {createColumnHelper} from "@tanstack/react-table";

function ClientList(props :any) {
    let [clients, setClients] = useState([] as Client[]);
    let [loading, setLoading] = useState(true);

    const getClients = useCallback(() => {
        API.clients.list({
            limit: 1000,
        }).then((response :any) => {
            let clientList = response.data.data || []
            // Update clients status to unknown at load
            clientList.map((c :Client) => {
                c.api_alive = Const.UNKNOWN
                c.ssh_alive = Const.UNKNOWN
                return c
            })
            setClients(clientList);
            setLoading(false);
        }).catch(error => {
            setLoading(false);
        });
    }, [])

    function setClientStateLoading(name :string, state :boolean) {
        setClients(clients => clients.map((c :Client) => {
            if (c.name === name) {
                return {
                    ...c,
                    state_is_loading: state,
                }
            } else {
                return c
            }
        }))
    }

    const pingClient = useCallback((name :string) => {
        setClientStateLoading(name, true)
        API.clients.get(name).then((response :any) => {
            // Wait 500ms before removing loading spinner to avoid blinking
            setTimeout(
                () => setClients(clients => clients.map((c :Client) => {
                        if (c.name === name) {
                            return {
                                ...c,
                                state_is_loading: false,
                                ssh_alive: response.data.ssh_alive,
                                ssh_alive_message: response.data.ssh_alive_message,
                                api_alive: response.data.api_alive,
                                api_alive_message: response.data.api_alive_message,
                            }
                        } else {
                            return c
                        }
                    })),
                500
            )
        }).catch(error => {
            console.log("Cannot get client details", error);
        });
    }, [])

    function pingAllClients() {
        clients.map((c :Client) => pingClient(c.name))
    }

    useEffect(() => {
        getClients();
    }, [getClients]);

    function renderModules(mods :string) {
        if (!mods) {
            return <span className="italic text-slate-400 dark:text-slate-600">None</span>;
        }

        let module_names :string[] = mods.split(",")
        return (
            <div className="flex flex-wrap gap-y-2 gap-x-1">
                {module_names.map((mod: any) => (
                    <div key={mod}>
                        <span className="badge">{mod}</span>
                    </div>
                ))}
            </div>
        )
    }

    function getActions() {
        return (
            <Dropdown>
                <div onClick={() => pingAllClients()}>Ping all clients</div>
            </Dropdown>
        )
    }

    const columnHelper = createColumnHelper<Client>()
    const columns = [
        columnHelper.accessor((client) => {return client}, {
            header: () => (<div className="py-2 px-3 w-full text-center">Health</div>),
            id: 'health',
            cell: (cell :any) => (<div className="py-2 px-3 text-center">{cell.getValue().state_is_loading ? (<Loader label="" />) : (<StatusDot status={ClientUtils.alive(cell.getValue())}/>)}</div>),
        }),
        columnHelper.accessor('name', {
            header: () => (<div className="py-2 px-3">Name</div>),
            id: 'name',
            cell: (cell :any) => (<div className="py-2 px-3"><Link to={`/clients/${cell.getValue()}`}>{cell.getValue()}</Link></div>),
        }),
        columnHelper.accessor('address', {
            header: () => (<div className="py-2 px-3">Address</div>),
            id: 'address',
            cell: (cell :any) => (<div className="py-2 px-3 code">{cell.getValue()}</div>),
        }),
        columnHelper.accessor((client) => (client.modules || []).map((mod :Module) => mod.name).join(", "), {
            header: () => (<div className="py-2 px-3 hidden md:block">Modules</div>),
            id: 'modules',
            cell: (cell :any) => (<div className="py-2 px-3 space-x-1 hidden md:block">{renderModules(cell.getValue())}</div>),
        }),
        columnHelper.accessor((client) => {return client}, {
            header: '',
            id: 'actions',
            cell: (cell :any) => (<Dropdown>
                <div onClick={() => pingClient(cell.getValue().name)}>Ping client</div>
            </Dropdown>),
        }),
    ]

    return (
        <Table title={props.title}
               filtered={props.filtered}
               sorted={props.sorted}
               paginated={props.paginated}
               columns={columns}
               defaultPageSize={props.limit || Const.DEFAULT_PAGE_SIZE}
               refreshFunc={getClients}
               data={clients}
               loading={loading}
               actions={getActions()}
        />
    );
}

export default ClientList;
