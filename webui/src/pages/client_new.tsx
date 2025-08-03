import React, { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";
import Card from '../components/card';
import Const from '../types/const';
import Client from '../types/client';
import API from '../utils/api';
import Utils from '../utils/utils';
import Module from '../types/module';
import ModuleForm from '../components/module_form';
import _ from 'lodash';

function ClientNew() {
    let [client, setClient] = useState<Client>({} as Client);
    let [mods, setModules] = useState<Module[]>([]);
    let [formErrors, setFormErrors] = useState<string[]>([]);
    let navigate = useNavigate();

    function updateFormErrors() {
        let errors: string[] = [];
        if (_.isEmpty(client.name)) {
            errors.push("Empty client name")
        }
        if (_.isEmpty(client.address)) {
            errors.push("Empty client address")
        }

        setFormErrors(errors)
    }

    useEffect(() => {
        updateFormErrors();
    }, [client, mods])

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

    function reloadConfig() {
        API.config.reload().then((_: any) => {
            Utils.notify(Const.OK, "Configuration reloaded", "Redirecting to client details");
            navigate(`/clients/${client?.name}`);
        }).catch(error => {
            console.log("Cannot create client", error);
            Utils.notify(Const.CRITICAL, "Cannot create client", error.toString())
        });
    }

    function createClient(e: React.MouseEvent<HTMLButtonElement>) {
        e.preventDefault();
        let toSend: Client = client;
        toSend.modules = mods

        API.clients.create(toSend).then((_: any) => {
            Utils.notify(Const.OK, "Client created", "Reloading configuration")
            reloadConfig()
        }).catch(error => {
            console.log("Cannot create client", error);
            Utils.notify(Const.CRITICAL, "Cannot create client", error.toString())
        });
    }

    const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        let value: any
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

    return (
        <>
            <Card>
                <div className="px-6 py-4 flex items-center">
                    <h3 className="flex-grow font-bold">
                        General info
                    </h3>

                    <div className=''>
                        <button className='btn btn-success btn-sm' onClick={createClient} disabled={formErrors.length != 0}>Create</button>
                    </div>
                </div>
                <div className='grid grid-cols-2 m-4 gap-4'>
                    <form className='space-y-2'>
                        <fieldset className="fieldset">
                            <label className="label">Name</label>
                            <div className="input validator">
                                <input name="name" type="text" required placeholder="Client name" value={client.name} onChange={onChange} />
                                <span className="badge badge-outline badge-xs">Required</span>
                            </div>
                        </fieldset>

                        <fieldset className="fieldset">
                            <label className="label">Address</label>
                            <div className="input validator">
                                <input name="address" type="text" required placeholder="IP or FQDN" value={client.address} onChange={onChange} />
                                <span className="badge badge-outline badge-xs">Required</span>
                            </div>
                        </fieldset>

                        <fieldset className="fieldset">
                            <label className="label">SSH user</label>
                            <input name="ssh_user" type="text" className="input" placeholder="relique" value={client.ssh_user} onChange={onChange} />
                        </fieldset>

                        <fieldset className="fieldset">
                            <label className="label">SSH port</label>
                            <input name="ssh_port" type="number" className="input" placeholder="22" value={client.ssh_port} onChange={onChange} />
                        </fieldset>
                    </form>
                    <div className="space-y-2">
                        <h3 className='font-bold'>Summary</h3>
                        <table className='table'>
                            <tbody>
                                {formErrors.map((err: string) => (
                                    <tr>
                                        <td className='alert alert-soft alert-error'>
                                            {err}
                                        </td>
                                    </tr>

                                ))}
                                <tr>
                                    <td>
                                        Creating client <span className="badge badge-neutral badge-soft">{client?.name}</span> with address
                                        <span className='code ml-2'>{client?.address || 'unknown'}</span>
                                    </td>
                                </tr>
                                <tr>
                                    <td>
                                        Adding <span className='code'>{mods.length ?? 0}</span> modules to client
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </Card>

            <Card>
                <div className="px-6 py-4 flex items-center">
                    <h3 className="flex-grow font-bold">
                        Modules
                    </h3>

                    <div className=''>
                        <button className='btn btn-sm' onClick={addModule}>Add module</button>
                    </div>
                </div>
                <div className='grid grid-cols-2 gap-4 m-4'>
                    {mods.map((m: Module, index: number) => {
                        return <ModuleForm className="border-dashed p-4" onChange={(mod: Module) => updateModule(index, mod)} module={m} />
                    })}
                </div>
            </Card>
        </>
    );
}

export default ClientNew;