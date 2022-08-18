import React, {useEffect, useState} from "react";

import Module from "../types/module";

function ModuleCard(props :any) {
    let mod :Module = props.module;
    let [showMoreContent, setShowMoreContent] = useState(false);

    if (mod === null) {
        return <div>Loading</div>
    }

    if (mod.variant === "") {
        mod.variant = "default";
    }

    if (mod.schedules.length === 0) {
        mod.schedules = [{name: "none"}];
    }

    function showLess() {
        setShowMoreContent(false);
    }

    function showMore() {
        setShowMoreContent(true);
    }

    function displayScriptName(name :string) {
        if (name === "none") {
            return <div className={"ml-3 text-slate-400 italic"}>None</div>;
        }

        return <div className={"ml-3 font-mono text-pink-500"}>{name}</div>;
    }

    function displayAdditionalParams(params :any) {
        if (Object.keys(params).length === 0) {
            return <div className={"ml-3 text-slate-400 italic"}>None</div>;
        }

        return <div className={"ml-3 font-mono text-pink-500 whitespace-pre"}>{JSON.stringify(params, null, 2)}</div>;
    }

    return (
        <div className={"bg-white m-2 px-4 pt-3 pb-0 rounded shadow flex flex-col divide-y divide-dashed"}>
            <div className={"pb-3"}>
                <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Name</div>
                <div className={"ml-3 font-bold"}>{mod.name}</div>
            </div>

            <div className={"py-3 grid grid-cols-2"}>
                <div>
                    <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Module type</div>
                    <div className={"ml-3"}>{mod.module_type}</div>
                </div>
                <div>
                    <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Variant</div>
                    <div className={"ml-3"}>{mod.variant}</div>
                </div>
            </div>

            <div className={"py-3 grid grid-cols-2"}>
                <div>
                    <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Backup type</div>
                    <div className={"ml-3"}>{mod.backup_type}</div>
                </div>
                <div>
                    <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Schedules</div>
                    <div className={"ml-3"}>{mod.schedules.map((s :any) => { return s.name }).join(", ")}</div>
                </div>
                <div className={`text-right text-blue-600 text-xs col-span-full ${showMoreContent && "hidden"}`} onClick={() => showMore()}>Show more</div>
            </div>

            <div className={`py-3 ${!showMoreContent && "hidden"}`}>
                <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Backup paths</div>
                {mod.backup_paths.map((path :string) => {
                    return <div key={path} className={"ml-3 font-mono text-pink-500"}>{path}</div>
                })}
            </div>

            <div className={`py-3 ${!showMoreContent && "hidden"}`}>
                <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Pre backup script</div>
                {displayScriptName(mod.pre_backup_script)}
            </div>

            <div className={`py-3 ${!showMoreContent && "hidden"}`}>
                <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Post backup script</div>
                {displayScriptName(mod.post_backup_script)}
            </div>
            <div className={`py-3 ${!showMoreContent && "hidden"}`}>
                <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Pre restore script</div>
                {displayScriptName(mod.pre_restore_script)}
            </div>

            <div className={`py-3 ${!showMoreContent && "hidden"}`}>
                <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Post restore script</div>
                {displayScriptName(mod.post_restore_script)}
            </div>

            <div className={`py-3 ${!showMoreContent && "hidden"}`}>
                <div className={"mb-2 font-bold text-xs text-slate-400 uppercase"}>Additional params</div>
                {displayAdditionalParams(mod.params)}
                <div className={"text-right text-blue-600 text-xs"} onClick={() => showLess()}>Show less</div>
            </div>
        </div>
    );
}

export default ModuleCard;
