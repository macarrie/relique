import React, {useEffect, useState} from "react";

import Module from "../types/module";
import Card from "./card";

function ModuleCard(props :any) {
    let mod :Module = props.module;
    let [showMoreContent, setShowMoreContent] = useState(false);

    useEffect(() => {
        if (props.full) {
            setShowMoreContent(true);
        }
    }, [props.full])

    if (mod === null) {
        return <div>Loading</div>
    }

    if (mod.variant === "") {
        mod.variant = "default";
    }

    if (mod.schedules === null || mod.schedules.length === 0) {
        mod.schedules = [{name: "none"}];
    }

    function showLess() {
        setShowMoreContent(false);
    }

    function showMore() {
        setShowMoreContent(true);
    }

    function displayScriptName(name :string) {
        if (name === "none" || name === "") {
            return <div className={"text-slate-400 italic"}>None</div>;
        }

        return <div className={"code"}>{name}</div>;
    }

    function displayAdditionalParams(params :any) {
        if (params === null || Object.keys(params).length === 0) {
            return <div className={"text-slate-400 italic"}>None</div>;
        }

        return <div className={"code whitespace-pre"}>{JSON.stringify(params, null, 2)}</div>;
    }

    function displayBackupPaths(paths :any) {
        if (paths === null || paths.length === 0 || (paths.length === 1 && paths[0] === "")) {
            return <div className={"text-slate-400 italic"}>None</div>;
        }

        return (
            <>
                {mod.backup_paths.map((path :string) => {
                    return <div key={path} className={"code"}>{path}</div>
                })}
            </>
        )
    }

    return (
        <Card className={`bg-white bg-opacity-60 ${props.className}`}>
            <div className="p-4 flex flex-row items-center mb-2">
                <div className="flex-grow font-bold text-slate-500 dark:text-slate-200">Module <span
                    className="ml-1 badge">{mod.name}</span></div>

                {(!props.full && showMoreContent) && (
                    <button className={"text-right button button-small button-text"}
                            onClick={() => showLess()}>Less</button>
                )}
                <button className={`text-right button button-small button-text ${showMoreContent && "hidden"}`}
                        onClick={() => showMore()}>More
                </button>
            </div>
            <table className="details-table ml-4 mb-2">
                <tr>
                    <td>Module type</td>
                    <td>{mod.module_type}</td>
                </tr>
                <tr>
                    <td>Variant</td>
                    <td>{mod.variant}</td>
                </tr>
                <tr>
                    <td>Backup type</td>
                    <td>{mod.backup_type}</td>
                </tr>
                <tr>
                    <td>Schedules</td>
                    <td>{mod.schedules.map((s: any) => {
                        return s.name
                    }).join(", ")}</td>
                </tr>
                <tr>
                    <td>Backup paths</td>
                    <td>{displayBackupPaths(mod.backup_paths)}</td>
                </tr>
                <tr className={`${!showMoreContent && "hidden"}`}>
                    <td>Pre backup script</td>
                    <td>{displayScriptName(mod.pre_backup_script)}</td>
                </tr>
                <tr className={`${!showMoreContent && "hidden"}`}>
                    <td>Post backup script</td>
                    <td>{displayScriptName(mod.post_backup_script)}</td>
                </tr>
                <tr className={`${!showMoreContent && "hidden"}`}>
                    <td>Pre restore script</td>
                    <td>{displayScriptName(mod.pre_restore_script)}</td>
                </tr>
                <tr className={`${!showMoreContent && "hidden"}`}>
                    <td>Post restore script</td>
                    <td>{displayScriptName(mod.post_restore_script)}</td>
                </tr>
                <tr className={`${!showMoreContent && "hidden"}`}>
                    <td>Additional params</td>
                    <td>{displayAdditionalParams(mod.params)}</td>
                </tr>
            </table>
        </Card>
    );
}

export default ModuleCard;
