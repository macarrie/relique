import { useCallback, useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import ClientList from "../components/client_list";
import Client from "../types/client";
import Const from "../types/const";

function Clients() {
    let [clients, setClients] = useState<Client[]>([]);

    function getClients() {
        API.clients.list({ limit: 10000 }).then((response: any) => {
            let clientList = response.data.data || []
            // Update clients status to unknown at load
            clientList.map((c: Client) => {
                c.ssh_alive = Const.UNKNOWN
                return c
            })
            setClients(clientList);
        }).catch(error => {
            console.log("Cannot get client list", error);
            setClients([]);
        });
    }

    useEffect(() => {
        getClients();
    }, [])

    function setClientStateLoading(name: string, state: boolean) {
        setClients(clients => clients.map((c: Client) => {
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

    const pingClient = useCallback((name: string) => {
        setClientStateLoading(name, true)
        API.clients.ping(name).then((response: any) => {
            // Wait 500ms before removing loading spinner to avoid blinking
            setTimeout(
                () => setClients(clients => clients.map((c: Client) => {
                    if (c.name === name) {
                        return {
                            ...c,
                            state_is_loading: false,
                            // TODO: Change
                            ssh_alive: response.data.ping_error === "" ? Const.OK : Const.CRITICAL,
                            ssh_alive_message: response.data.ping_error ?? "",
                        }
                    } else {
                        return c
                    }
                })),
                500
            )
        }).catch(error => {
            console.log("Cannot ping client", error);
            setClientStateLoading(name, false)
        });
    }, [])

    function pingAllClients() {
        clients.map((c: Client) => pingClient(c.name))
    }

    return (
        <>
            <Card>
                <ClientList
                    title="All clients"
                    actions={true}
                    custom_actions={[
                        <div className="btn btn-sm" onClick={() => pingAllClients()}>Ping all</div>,
                    ]}
                    data={clients}
                    paginated={true}
                    sorted={true}
                />
            </Card>
        </>
    );
}

export default Clients;