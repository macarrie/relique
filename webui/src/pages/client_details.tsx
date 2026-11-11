import { useCallback, useEffect, useRef, useState } from 'react';
import { Link, useNavigate, useParams } from "react-router-dom";
import Const from '../types/const';
import Status from '../components/status';
import PageHeader from '../components/page_header';
import Section from '../components/section';
import PropertyList from '../components/property_list';
import ModuleProperties from '../components/module_properties';
import EmptyValue from '../components/empty_value';
import ClientUtils from '../utils/client';
import Client from '../types/client';
import API from '../utils/api';
import Module from '../types/module';
import Image from '../types/image';
import ImageList from '../components/image_list';
import Utils from '../utils/utils';

function ClientDetails() {
    const { client_name } = useParams();
    const [c, setClient] = useState<Client>({} as Client);
    const [loading, setLoading] = useState<boolean>(true);
    const [imgs, setImages] = useState<Image[]>([]);
    const [imgsLoading, setImgsLoading] = useState<boolean>(true);

    const navigate = useNavigate();
    const deleteDialog = useRef<HTMLDialogElement>(null);

    const pingClient = useCallback((name: string) => {
        API.clients.ping(name).then((response) => {
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
        function getClient() {
            setLoading(true);
            if (client_name === undefined) {
                console.log("Client name undefined, cannot get client details");
                return;
            }

            API.clients.get(client_name).then((response) => {
                const c: Client = response.data
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
    }, [client_name, pingClient])

    useEffect(() => {
        function getImageList() {
            setImgsLoading(true);
            API.images.list({ limit: 10000, client: client_name }).then((response) => {
                setImages(response.data.data ?? []);
                setImgsLoading(false);
            }).catch(error => {
                console.log("Cannot get image list", error);
                Utils.notify(Const.CRITICAL, "Cannot get image list", error.toString())
                setImages([]);
            });
        }

        getImageList();
    }, [client_name])

    function reloadConfig() {
        API.config.reload().then(() => {
            Utils.notify(Const.OK, "Configuration reloaded", "Redirecting to client list");
            navigate(`/clients`);
        }).catch(error => {
            console.log("Cannot reload config", error);
            Utils.notify(Const.CRITICAL, "Cannot reload config", error.toString())
        });
    }

    function deleteClient() {
        API.clients.delete(c).then(() => {
            Utils.notify(Const.OK, "Client deleted", "Reloading configuration")
            reloadConfig();
        }).catch(error => {
            console.log("Cannot delete client", error);
            Utils.notify(Const.CRITICAL, "Cannot delete client", error.toString())
        });
    }

    function displayModules(mods: Module[]) {
        if (loading) {
            return <ModuleProperties module={{} as Module} loading />;
        }

        if (!mods || mods.length === 0) {
            return <EmptyValue label="No module configured" />;
        }

        return (
            <div className="grid lg:grid-cols-2 gap-x-10 gap-y-8">
                {mods.map((m: Module) => (
                    <div key={m.name} className="space-y-2">
                        <h3 className="font-display text-base font-semibold">{m.name}</h3>
                        <ModuleProperties module={m} />
                    </div>
                ))}
            </div>
        );
    }

    const sshStatus = (
        <Status
            loading={c.state_is_loading}
            status={c.ssh_alive}
            label={ClientUtils.sshStatusLabel(c.ssh_alive)} />
    );

    return (
        <>
            <PageHeader
                title={<><span className="font-normal text-muted">Client</span> {c.name}</>}
                loading={loading}
                meta={!loading && (
                    <>
                        <span className="code text-muted">{ClientUtils.sshTarget(c)}</span>
                        {sshStatus}
                    </>
                )}
                actions={<>
                    <Link to={`/clients/${client_name}/edit`} className="btn btn-sm">
                        <i className="ri-edit-line"></i>
                        Edit
                    </Link>
                    <button className="btn btn-sm btn-ghost text-error" onClick={() => deleteDialog.current?.showModal()}>
                        <i className="ri-delete-bin-line"></i>
                        Delete
                    </button>
                </>}
            />

            <dialog ref={deleteDialog} className="modal">
                <div className="modal-box space-y-3">
                    <h3 className="font-display text-lg font-semibold">Delete client <span className="code">{c.name}</span>?</h3>
                    <p>The client is removed from the Relique configuration. This cannot be undone.</p>
                    <div className="modal-action">
                        <button type="button" className="btn btn-sm" onClick={() => deleteDialog.current?.close()}>Cancel</button>
                        <button type="button" className="btn btn-sm btn-error" onClick={() => { deleteClient(); deleteDialog.current?.close(); }}>Delete client</button>
                    </div>
                </div>
                <form method="dialog" className="modal-backdrop"><button>close</button></form>
            </dialog>

            <Section title="Connection">
                <PropertyList loading={loading} items={[
                    { label: "Address", value: c.address, mono: true },
                    { label: "SSH user", value: ClientUtils.sshUser(c), mono: true },
                    { label: "SSH port", value: ClientUtils.sshPort(c), mono: true },
                    {
                        label: "SSH status", value: (
                            <div className="space-y-1">
                                {sshStatus}
                                {c.ssh_alive === Const.CRITICAL && c.ssh_alive_message && (
                                    <p className="code text-error whitespace-pre-wrap">{c.ssh_alive_message}</p>
                                )}
                            </div>
                        )
                    },
                ]} />
            </Section>

            <Section title="Modules">
                {displayModules(c.modules)}
            </Section>

            <Section title="Images">
                <ImageList
                    actions={true}
                    data={imgs}
                    paginated={true}
                    sorted={true}
                    loading={imgsLoading}
                    deletable={true}
                    onDeleted={(uuid: string) => setImages(prev => prev.filter(img => img.uuid !== uuid))}
                />
            </Section>
        </>
    );
}

export default ClientDetails;