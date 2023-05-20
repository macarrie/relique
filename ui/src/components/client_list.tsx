import React, {useCallback, useEffect, useState} from "react";
import {Link} from "react-router-dom";
import {Column} from "react-table";

import API from "../utils/api";
import ClientUtils from "../utils/client";

import Client from "../types/client";
import Module from "../types/module";

import StatusDot from "./status_dot";
import Table from "./table";
import TableUtils from "../utils/table";
import Dropdown from "./dropdown";
import Const from "../types/const";
import Loader from "./loader";

function ClientList(props :any) {
    let [limit, setLimit] = useState(props.limit || 0);
    let [clients, setClients] = useState([] as Client[]);
    let [loading, setLoading] = useState(true);

    const getClients = useCallback(() => {
        API.clients.list({
            limit: limit,
        }).then((response :any) => {
            let clientList = response.data || []
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
    }, [limit])

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
        setLimit(props.limit);
    }, [props.limit])

    useEffect(() => {
        getClients();
    }, [getClients]);

    function renderModules(mods :string) {
        if (!mods) {
            return <span className="italic text-slate-400 dark:text-slate-600">None</span>;
        }

        let module_names :string[] = mods.split(",")
        return (
            <>
                {module_names.map((mod: any) => (
                    <span className="badge" key={mod}>{mod}</span>
                ))}
            </>
        )
    }

    function getActions() {
        return (
            <Dropdown>
                <div onClick={() => pingAllClients()}>Ping all clients</div>
            </Dropdown>
        )
    }

    const columns :Array<Column<Client>> = React.useMemo(() => [
        {
            Header: () => (<div className="py-2 px-3 w-full text-center">Health</div>),
            accessor: (client) => {return client},
            id: 'health',
            Cell: ({value} :any) => (<div className="py-2 px-3 text-center">{value.state_is_loading ? (<Loader label="" />) : (<StatusDot status={ClientUtils.alive(value)}/>)}</div>),
        },
        {
            Header: () => (<div className="py-2 px-3">Name</div>),
            accessor: 'name',
            id: 'name',
            Cell: ({value} :any) => (<div className="py-2 px-3"><Link to={`/clients/${value}`}>{value}</Link></div>),
        },
        {
            Header: () => (<div className="py-2 px-3">Address</div>),
            accessor: 'address',
            id: 'address',
            Cell: ({value} :any) => (<div className="py-2 px-3 code">{value}</div>),
        },
        {
            Header: () => (<div className="py-2 px-3 hidden md:block">Modules</div>),
            accessor: (client) => (client.modules || []).map((mod :Module) => mod.name).join(", "),
            id: 'modules',
            Cell: ({value} :any) => (<div className="py-2 px-3 space-x-1 hidden md:block">{renderModules(value)}</div>),
        },
        {
            Header: '',
            accessor: (client) => {return client},
            id: 'actions',
            Cell: ({value} :any) => (<Dropdown>
                <div onClick={() => pingClient(value.name)}>Ping client</div>
            </Dropdown>),
        },
    ], [pingClient]);

    if (loading) {
        return (
            <Table title={props.title}
                   filtered={false}
                   sorted={false}
                   refreshFunc={getClients}
                   columns={TableUtils.GetPlaceholderColumns(columns)}
                   data={[{}, {}, {}]} />
        );
    }

    return (
        <Table title={props.title}
               filtered={props.filtered}
               sorted={props.sorted}
               refreshFunc={getClients}
               columns={columns}
               actions={getActions()}
               data={clients} />
    );
}

export default ClientList;
