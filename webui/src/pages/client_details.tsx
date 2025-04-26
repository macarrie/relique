import { useCallback, useEffect, useState } from 'react';
import { useParams } from "react-router-dom";
import Card from '../components/card';
import Const from '../types/const';
import StatusDot from '../components/status_dot';
import Client from '../types/client';
import API from '../utils/api';
import ModuleCard from '../components/module_card';
import Module from '../types/module';
import Image from '../types/image';
import ImageList from '../components/image_list';
import ClientCard from '../components/client_card';

function ClientDetails() {
    const { client_name } = useParams();
    let [c, setClient] = useState<Client>({} as Client);
    let [imgs, setImages] = useState<Image[]>([]);

    useEffect(() => {
        function getClient() {
            if (client_name === undefined) {
                console.log("Client name undefined, cannot get client details");
                return;
            }

            API.clients.get(client_name).then((response: any) => {
                let c: Client = response.data
                c.state_is_loading = true;
                c.ssh_alive = Const.UNKNOWN;
                c.ssh_alive_message = "";
                setClient(c);
                pingClient(client_name);
            }).catch(error => {
                console.log("Cannot get client details", error);
                setClient({} as Client);
            });
        }

        getClient();
    }, [client_name])

    const pingClient = useCallback((name: string) => {
        API.clients.ping(name).then((response: any) => {
            // Wait 500ms before removing loading spinner to avoid blinking
            setTimeout(
                () => setClient(c => {
                    return {
                        ...c,
                        state_is_loading: false,
                        ssh_alive: response.data.ping_error === "" ? Const.OK : Const.CRITICAL,
                        ssh_alive_message: response.data.ping_error ?? "",
                    }
                }),
                500
            )
        }).catch(error => {
            console.log("Cannot ping client", error);
        });
    }, [])

    useEffect(() => {
        function getImageList() {
            API.images.list({ limit: 10000, client: client_name }).then((response: any) => {
                setImages(response.data.data ?? []);
            }).catch(error => {
                console.log("Cannot get image list", error);
                setImages([]);
            });
        }

        getImageList();
    }, [])

    function displayModules(mods: Module[]) {
        if (!mods || mods.length === 0) {
            return <div className={"text-base-content/70 italic"}>None</div>;
        }

        return <>
            {mods.map((m: Module) => {
                return <ModuleCard key={m.name} module={m} full />
            })}
        </>;

    }

    return (
        <>
            <Card>
                <div className="px-6 py-4 flex">
                    <h3 className="flex-grow font-bold">
                        General info
                    </h3>
                </div>

                <div className="grid md:grid-cols-2 gap-4 m-4">
                    <ClientCard client={c} link={false} />
                    <Card>
                        <div className="p-4 flex flex-row items-center mb-2">
                            <div className={"flex-grow font-bold"}>Health</div>
                        </div>
                        <table className="table">
                            <tr>
                                <td>SSH</td>
                                <td>
                                    <div className="flex flex-row items-center mb-2">
                                        <div className="mr-2">
                                            {c.state_is_loading ? (
                                                <span className="text-neutral-300 loading loading-spinner loading-xs"></span>
                                            ) : (
                                                <StatusDot status={c.ssh_alive} />
                                            )}
                                        </div>
                                        <div className="text-sm">
                                            {c.ssh_alive === Const.UNKNOWN && (
                                                <span>SSH connectivity unknown</span>
                                            )}
                                            {c.ssh_alive === Const.OK && (
                                                <span>
                                                    SSH ping successful
                                                </span>
                                            )}
                                            {c.ssh_alive === Const.CRITICAL && (
                                                <span>
                                                    Cannot reach client via SSH
                                                </span>
                                            )}
                                        </div>
                                    </div>
                                    {c.ssh_alive === Const.OK && (
                                        <div className="block alert alert-success alert-soft">
                                            <span>
                                                SSH ping successful from relique server to
                                            </span>
                                            <span className='ml-1 code'>
                                                {c.ssh_user == "" ? Const.DEFAULT_CLIENT_SSH_USER : c.ssh_user}@{c.address}:{c.ssh_port == 0 ? Const.DEFAULT_CLIENT_SSH_PORT : c.ssh_port}
                                            </span>
                                        </div>
                                    )}
                                    {c.ssh_alive === Const.CRITICAL && (
                                        <div className="block alert alert-error alert-soft">{c.ssh_alive_message}</div>
                                    )}
                                </td>
                            </tr>
                        </table>
                    </Card>
                </div>
            </Card>

            <Card>
                <div className="px-6 py-4 flex">
                    <h3 className="flex-grow font-bold">
                        Modules
                    </h3>
                </div>
                <div className="grid md:grid-cols-2 gap-4 m-4">
                    {displayModules(c.modules)}
                </div>
            </Card>

            <Card>
                <ImageList
                    title="Associated images"
                    actions={true}
                    data={imgs}
                    paginated={true}
                    sorted={true}
                />
            </Card>

        </>
    );
}

export default ClientDetails;