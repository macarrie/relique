import { useEffect, useState } from "react";

import Module from "../types/module";
import Card from "./card";
import Const from "../types/const";
import Utils from "../utils/utils";
import API from "../utils/api";

function ModuleForm(props: any) {
    let [mod, setModule] = useState<Module>(props.module);
    let [moduleList, setModules] = useState<Array<Module>>([]);

    useEffect(() => {
        function getModuleList() {
            API.modules.list({ limit: 10000 }).then((response: any) => {
                setModules(response.data.data ?? []);
            }).catch(error => {
                console.log("Cannot get module list", error);
                Utils.notify(Const.CRITICAL, "Cannot get module list", error.toString())
                setModules([]);
            });
        }

        getModuleList();
    }, [])

    const handleFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        let newMod = {
            ...mod,
            [e.target.name]: e.target.value
        };
        setModule(newMod);
        props.onChange(newMod);
    }

    const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        let newMod = {
            ...mod,
            [e.target.name]: e.target.checked
        };
        setModule(newMod)
        props.onChange(newMod);
    }

    const handleArrayChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        let newMod = {
            ...mod,
            [e.target.name]: e.target.value.split(',')
        };
        setModule(newMod)
        props.onChange(newMod);
    }

    return (
        <Card className={props.className}>
            <button className="float-right text-slate-500 cursor-pointer" onClick={() => {props.onDeleteRequest(mod)}}>
                <i className={`text-lg ri-close-line`}></i>
            </button>

            <form className='space-y-2'>
                <fieldset className="fieldset">
                    <label className="label">Module name</label>
                    <div className="input validator">
                        <input name="name" type="text" required placeholder="Module name" value={mod.name} onChange={handleFormChange} />
                        <span className="badge badge-outline badge-xs">Required</span>
                    </div>
                </fieldset>

                <fieldset className="fieldset">
                    <label className="label">Module type</label>
                    <select name="module_type" defaultValue="generic" value={mod.module_type} className="select" onChange={handleFormChange}>
                        {moduleList.map((m: Module, index: number) => {
                            return (
                                <option key={index}>{m.name}</option>
                            )
                        })}
                    </select>
                </fieldset>

                <fieldset className="fieldset">
                    <label className="label">Backup type</label>
                    <select name="backup_type" defaultValue="diff" value={mod.backup_type} className="select" onChange={handleFormChange}>
                        <option>diff</option>
                        <option>full</option>
                    </select>
                </fieldset>

                <fieldset className="fieldset">
                    <label className="label">Backup paths</label>
                    <div className="input">
                        <input name="backup_paths" type="text" placeholder="Backup paths" value={(mod.backup_paths ?? []).join(',')} onChange={handleArrayChange} />
                    </div>
                </fieldset>

                <fieldset className="fieldset">
                    <label className="label">Exclusions</label>
                    <div className="input">
                        <input name="exclude" type="text" placeholder="Exclusions" value={(mod.exclude ?? []).join(',')} onChange={handleArrayChange} />
                    </div>
                </fieldset>

                <fieldset className="fieldset">
                    <label className="label">Inclusions</label>
                    <div className="input">
                        <input name="include" type="text" placeholder="Inclusions" value={(mod.include ?? []).join(',')} onChange={handleArrayChange} />
                    </div>
                </fieldset>

                <fieldset className="fieldset">
                    <label className="label">
                        <input name="exclude_cvs" type="checkbox" placeholder="Exclude CVS" checked={mod.exclude_cvs} onChange={handleCheckboxChange} />
                        Exclude CVS
                    </label>
                </fieldset>
            </form>
        </Card>
    );
}

export default ModuleForm;