import { useEffect, useState } from "react";

import Module from "../types/module";
import Const from "../types/const";
import Utils from "../utils/utils";
import API from "../utils/api";

type ModuleFormProps = {
    id: string,
    module: Module,
    className?: string,
    onChange: (mod: Module) => void,
    onDeleteRequest: (mod: Module) => void,
};

function ModuleForm(props: ModuleFormProps) {
    const [mod, setModule] = useState<Module>(props.module);
    const [moduleList, setModules] = useState<Array<Module>>([]);

    useEffect(() => {
        function getModuleList() {
            API.modules.list({ limit: 10000 }).then((response) => {
                setModules(response.data.data ?? []);
            }).catch(error => {
                console.log("Cannot get module list", error);
                Utils.notify(Const.CRITICAL, "Cannot get module list", error.toString())
                setModules([]);
            });
        }

        getModuleList();
    }, [])

    const update = (newMod: Module) => {
        setModule(newMod);
        props.onChange(newMod);
    }

    const handleFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        update({ ...mod, [e.target.name]: e.target.value });
    }

    const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        update({ ...mod, [e.target.name]: e.target.checked });
    }

    const handleArrayChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        update({ ...mod, [e.target.name]: e.target.value.split(',') });
    }

    return (
        <div className={`space-y-3 ${props.className ?? ""}`}>
            <div className="flex items-baseline gap-4">
                <h3 className="flex-grow font-display text-base font-semibold">
                    {mod.name || <span className="text-muted font-normal">Unnamed module</span>}
                </h3>
                <button type="button" className="btn btn-xs btn-ghost text-error" onClick={() => props.onDeleteRequest(mod)}>
                    <i className="ri-close-line"></i>
                    Remove
                </button>
            </div>

            <fieldset className="fieldset">
                <label className="label" htmlFor={`${props.id}-name`}>Name</label>
                <input id={`${props.id}-name`} name="name" type="text" className={`input w-full ${mod.name === "" ? "input-error" : ""}`} required placeholder="Module name" value={mod.name ?? ""} onChange={handleFormChange} />
                {mod.name === "" && <p className="text-error">A module name is required</p>}
            </fieldset>

            <div className="grid grid-cols-2 gap-4">
                <fieldset className="fieldset">
                    <label className="label" htmlFor={`${props.id}-type`}>Module type</label>
                    <select id={`${props.id}-type`} name="module_type" value={mod.module_type} className="select w-full" onChange={handleFormChange}>
                        {mod.module_type && !moduleList.some((m: Module) => m.name === mod.module_type) && (
                            <option key={mod.module_type}>{mod.module_type}</option>
                        )}
                        {moduleList.map((m: Module) => (
                            <option key={m.name}>{m.name}</option>
                        ))}
                    </select>
                </fieldset>

                <fieldset className="fieldset">
                    <label className="label" htmlFor={`${props.id}-backup-type`}>Backup type</label>
                    <select id={`${props.id}-backup-type`} name="backup_type" value={mod.backup_type} className="select w-full" onChange={handleFormChange}>
                        <option>diff</option>
                        <option>full</option>
                    </select>
                </fieldset>
            </div>

            <fieldset className="fieldset">
                <label className="label" htmlFor={`${props.id}-paths`}>Backup paths</label>
                <input id={`${props.id}-paths`} name="backup_paths" type="text" className="input w-full font-mono" placeholder="/etc,/home" value={(mod.backup_paths ?? []).join(',')} onChange={handleArrayChange} />
                <p className="label">Comma-separated</p>
            </fieldset>

            <fieldset className="fieldset">
                <label className="label" htmlFor={`${props.id}-exclude`}>Exclusions</label>
                <input id={`${props.id}-exclude`} name="exclude" type="text" className="input w-full font-mono" value={(mod.exclude ?? []).join(',')} onChange={handleArrayChange} />
            </fieldset>

            <fieldset className="fieldset">
                <label className="label" htmlFor={`${props.id}-include`}>Inclusions</label>
                <input id={`${props.id}-include`} name="include" type="text" className="input w-full font-mono" value={(mod.include ?? []).join(',')} onChange={handleArrayChange} />
            </fieldset>

            <label className="label text-base-content">
                <input name="exclude_cvs" type="checkbox" className="checkbox checkbox-sm" checked={mod.exclude_cvs ?? false} onChange={handleCheckboxChange} />
                Exclude CVS files (.git, .svn…)
            </label>
        </div>
    );
}

export default ModuleForm;
