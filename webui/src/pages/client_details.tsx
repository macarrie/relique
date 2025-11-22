import { useCallback, useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from "react-router-dom";
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
import Utils from '../utils/utils';
import TextPlaceholder from '../components/text_placeholder';

function ClientDetails() {
    const { client_name } = useParams();
    let [c, setClient] = useState<Client>({} as Client);
    let [loading, setLoading] = useState<boolean>(true);
    let [imgs, setImages] = useState<Image[]>([]);
    let [imgsLoading, setImgsLoading] = useState<boolean>(true);

    let navigate = useNavigate();

    useEffect(() => {
        function getClient() {
            setLoading(true);
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
                setLoading(false);
                pingClient(client_name);
            }).catch(error => {
                console.log("Cannot get client details", error);
                Utils.notify(Const.CRITICAL, "Cannot get client details", error.toString())
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
            Utils.notify(Const.CRITICAL, "Cannot get ping client", error.toString())
        });
    }, [])

    useEffect(() => {
        function getImageList() {
            setImgsLoading(true);
            API.images.list({ limit: 10000, client: client_name }).then((response: any) => {
                setImages(response.data.data ?? []);
                setImgsLoading(false);
            }).catch(error => {
                console.log("Cannot get image list", error);
                Utils.notify(Const.CRITICAL, "Cannot get image list", error.toString())
                setImages([]);
            });
        }

        getImageList();
    }, [])

    function reloadConfig() {
        API.config.reload().then((_: any) => {
            Utils.notify(Const.OK, "Configuration reloaded", "Redirecting to client list");
            navigate(`/clients`);
        }).catch(error => {
            console.log("Cannot reload config", error);
            Utils.notify(Const.CRITICAL, "Cannot reload config", error.toString())
        });
    }

    function deleteClient() {
        API.clients.delete(c).then((_: any) => {
            Utils.notify(Const.OK, "Client deleted", "Reloading configuration")
            reloadConfig();
        }).catch(error => {
            console.log("Cannot delete client", error);
            Utils.notify(Const.CRITICAL, "Cannot delete client", error.toString())
        });
    }

    function displayModules(mods: Module[]) {
        if (loading) {
            return (
                <TextPlaceholder />
            )
        }

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
                    <div className='space-x-2'>
                        <dialog id="delete_confirm" className="modal">
                            <div className="modal-box">
                                <h3 className="font-bold text-lg">Client deletion</h3>
                                <p className="py-4">Client info will be deleted from Relique configuration. This action is not recoverable.</p>
                                <p className="py-4">Do you want to continue ?</p>
                                <div className="modal-action">
                                    <form method="dialog" className='space-x-2'>
                                        <button type="button" className="btn" onClick={() => (document.getElementById('delete_confirm') as HTMLDialogElement | null)?.close()}>Cancel</button>
                                        <button type="button" className="btn btn-error" onClick={() => { deleteClient(); (document.getElementById('delete_confirm') as HTMLDialogElement | null)?.close(); }}>Delete</button>
                                    </form>
                                </div>
                            </div>
                        </dialog>                       
                        <div className="btn btn-sm btn-error" onClick={() => (document.getElementById('delete_confirm') as HTMLDialogElement | null)?.showModal()}>Delete</div>
                        <Link to={`/clients/${client_name}/edit`} className="cursor-pointer text-slate-500">
                            <div className="btn btn-sm" onClick={() => console.log("Delete client")}>
                                <i className={`text-lg ri-edit-fill`}></i>
                                Edit
                            </div>
                        </Link>
                    </div>
                </div>

                <div className="grid md:grid-cols-2 gap-4 m-4">
                    <ClientCard client={c} link={false} loading={loading} />
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
                                        {loading && <TextPlaceholder />}
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
                    loading={imgsLoading}
                />
            </Card>

        </>
    );
}

export default ClientDetails;