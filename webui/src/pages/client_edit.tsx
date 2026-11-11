import React, { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from "react-router-dom";
import PageHeader from '../components/page_header';
import Section from '../components/section';
import Const from '../types/const';
import Client from '../types/client';
import API from '../utils/api';
import Utils from '../utils/utils';
import Module from '../types/module';
import ModuleForm from '../components/module_form';
import lo from 'lodash';


function ClientEdit() {
    const { client_name } = useParams();
    const [client, setClient] = useState<Client>({} as Client);
    const [mods, setModules] = useState<Module[]>([]);

    // true if edition, false if creation
    const editMode = client_name !== undefined;

    const navigate = useNavigate();

    useEffect(() => {
        function getClient() {
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
                setModules(c.modules ?? []);
            }).catch(error => {
                console.log("Cannot get client details", error);
                Utils.notify(Const.CRITICAL, "Cannot get client details", error.toString())
                setClient({} as Client);
            });
        }

        getClient();
    }, [client_name])

    function getFormErrors(): Record<string, string> {
        const errors: Record<string, string> = {};
        if (lo.isEmpty(client.name)) {
            errors.name = "A client name is required";
        }
        if (lo.isEmpty(client.address)) {
            errors.address = "An address is required";
        }
        if (mods.some((m: Module) => lo.isEmpty(m.name))) {
            errors.modules = "Every module needs a name";
        }

        return errors;
    }

    const formErrors = getFormErrors();

    function addModule() {
        setModules([
            ...mods,
            {
                module_type: "generic",
                backup_type: "diff",
            } as Module,
        ]);
    }

    function updateModule(index: number, mod: Module) {
        setModules(mods => mods.map((m: Module, i: number) => {
            if (i === index) {
                return mod
            } else {
                return m
            }
        }))
    }

    function removeModule(index: number) {
        setModules(mods.filter((_m, i) => i !== index));
    }

    function reloadConfig() {
        API.config.reload().then(() => {
            Utils.notify(Const.OK, "Configuration reloaded", "Redirecting to client details");
            navigate(`/clients/${client?.name}`);
        }).catch(error => {
            console.log("Cannot reload config", error);
            Utils.notify(Const.CRITICAL, "Cannot reload config", error.toString())
        });
    }

    function createClient(e: React.MouseEvent<HTMLButtonElement>) {
        e.preventDefault();
        const toSend: Client = { ...client, modules: mods };

        API.clients.create(toSend).then(() => {
            Utils.notify(Const.OK, "Client created", "Reloading configuration")
            reloadConfig()
        }).catch(error => {
            console.log("Cannot create client", error);
            Utils.notify(Const.CRITICAL, "Cannot create client", error.toString())
        });
    }

    function saveClient(e: React.MouseEvent<HTMLButtonElement>) {
        e.preventDefault();
        const toSend: Client = { ...client, modules: mods };

        API.clients.update(toSend).then(() => {
            Utils.notify(Const.OK, "Client saved", "Reloading configuration")
            reloadConfig()
        }).catch(error => {
            console.log("Cannot save modified client", error);
            Utils.notify(Const.CRITICAL, "Cannot save modified client", error.toString())
        });
    }

    const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        let value: string | number
        if (e.target.type == "number") {
            value = parseInt(e.target.value)
        } else {
            value = e.target.value;
        }

        setClient({
            ...client,
            [e.target.name]: value,
        })
    }

    const hasErrors = Object.keys(formErrors).length !== 0;

    return (
        <>
            <PageHeader
                title={editMode ? <><span className="font-normal text-muted">Edit client</span> {client_name}</> : "New client"}
                actions={<>
                    <Link to={editMode ? `/clients/${client_name}` : "/clients"} className="btn btn-sm">
                        Cancel
                    </Link>
                    <button className="btn btn-sm btn-primary" onClick={editMode ? saveClient : createClient} disabled={hasErrors}>
                        {editMode ? "Save changes" : "Create client"}
                    </button>
                </>}
            />

            <Section title="Connection">
                <form className="grid lg:grid-cols-2 gap-x-6 gap-y-1 max-w-5xl" onSubmit={(e) => e.preventDefault()}>
                    <fieldset className="fieldset">
                        <label className="label" htmlFor="client-name">Name</label>
                        <input id="client-name" name="name" type="text" className={`input w-full disabled:bg-base-300/60 disabled:text-base-content/70 ${formErrors.name && client.name !== undefined ? "input-error" : ""}`} required placeholder="web-01" value={client.name ?? ""} onChange={onChange} disabled={editMode ? true : false} />
                        {editMode
                            ? <p className="label">The name cannot be changed after creation</p>
                            : formErrors.name && client.name !== undefined && <p className="text-error">{formErrors.name}</p>}
                    </fieldset>

                    <fieldset className="fieldset">
                        <label className="label" htmlFor="client-address">Address</label>
                        <input id="client-address" name="address" type="text" className={`input w-full font-mono ${formErrors.address && client.address !== undefined ? "input-error" : ""}`} required placeholder="IP or FQDN" value={client.address ?? ""} onChange={onChange} />
                        {formErrors.address && client.address !== undefined && <p className="text-error">{formErrors.address}</p>}
                    </fieldset>

                    <div className="grid grid-cols-2 gap-x-6">
                        <fieldset className="fieldset">
                            <label className="label" htmlFor="client-ssh-user">SSH user</label>
                            <input id="client-ssh-user" name="ssh_user" type="text" className="input w-full disabled:bg-base-300/60 disabled:text-base-content/70 font-mono" placeholder={Const.DEFAULT_CLIENT_SSH_USER} value={client.ssh_user ?? ""} onChange={onChange} />
                        </fieldset>

                        <fieldset className="fieldset">
                            <label className="label" htmlFor="client-ssh-port">SSH port</label>
                            <input id="client-ssh-port" name="ssh_port" type="number" className="input w-full disabled:bg-base-300/60 disabled:text-base-content/70 font-mono" placeholder={String(Const.DEFAULT_CLIENT_SSH_PORT)} value={client.ssh_port || ""} onChange={onChange} />
                        </fieldset>
                    </div>
                </form>
            </Section>

            <Section
                title="Modules"
                description={mods.length === 0 ? "No module configured yet." : undefined}
                actions={
                    <button className="btn btn-sm" onClick={addModule}>
                        <i className="ri-add-line"></i>
                        Add module
                    </button>
                }>
                <div className="grid lg:grid-cols-2 gap-x-10 gap-y-10 max-w-5xl">
                    {mods.map((m: Module, index: number) => (
                        <ModuleForm
                            key={index}
                            id={`module-${index}`}
                            onChange={(mod: Module) => updateModule(index, mod)}
                            onDeleteRequest={() => removeModule(index)}
                            module={m} />
                    ))}
                </div>
            </Section>
        </>
    );
}

export default ClientEdit;