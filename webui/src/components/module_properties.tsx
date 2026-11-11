import { useEffect, useState } from "react";

import Module from "../types/module";
import Card from "./card";
import StatusBadge from "./status_badge";
import Const from "../types/const";
import TextPlaceholder from "./text_placeholder";

function ModuleCard(props: any) {
    let mod: Module = props.module;
    let [showMoreContent, setShowMoreContent] = useState(false);

    useEffect(() => {
        if (props.full) {
            setShowMoreContent(true);
        }
    }, [props.full])

    if (mod?.variant === "") {
        mod.variant = "default";
    }

    function showLess() {
        setShowMoreContent(false);
    }

    function showMore() {
        setShowMoreContent(true);
    }

    function displayBackupPaths(paths: any) {
        let modPaths = paths ?? [];
        if (modPaths.length === 0) {
            return <div className={"text-base-content/50 italic"}>None</div>;
        }

        return (
            <>
                {modPaths.map((path: string) => {
                    return <div key={path} className={"code"}>{path}</div>
                })}
            </>
        )
    }

    return (
        <Card className={props.className}>
            <div className="p-4 flex flex-row items-center mb-2">
                <div className="flex-grow font-bold align-middle">Module
                    {props.loading ? (
                        <span className="ml-2 badge skeleton w-16"></span>
                    ) : (
                        <span className='ml-2 badge badge-neutral badge-soft font-normal'>{mod?.name}</span>
                    )}
                </div>

                {(!props.full && showMoreContent) && (
                    <button className={"text-right button button-small button-text"}
                        onClick={() => showLess()}>Less</button>
                )}
                <button className={`text-right button button-small button-text ${showMoreContent && "hidden"}`}
                    onClick={() => showMore()}>More
                </button>
            </div>
            <table className="table">
                <tr>
                    <td>Module type</td>
                    <td>
                        {props.loading && (
                            <TextPlaceholder />
                        )}
                        {mod?.module_type}
                    </td>
                </tr>
                <tr>
                    <td>Variant</td>
                    <td>
                        {props.loading && (
                            <TextPlaceholder />
                        )}
                        {mod?.variant}
                    </td>
                </tr>
                <tr>
                    <td>Backup type</td>
                    <td>
                        {props.loading && (
                            <TextPlaceholder />
                        )}
                        {mod?.backup_type}
                    </td>
                </tr>
                <tr>
                    <td>Backup paths</td>
                    <td>
                        {props.loading ? (
                            <TextPlaceholder />
                        ) : (
                            <>
                                {displayBackupPaths(mod?.backup_paths)}
                            </>
                        )}
                    </td>
                </tr>
                <tr>
                    <td>Exclusions</td>
                    <td>
                        {props.loading ? (
                            <TextPlaceholder />
                        ) : (
                            <>
                                {displayBackupPaths(mod?.exclude)}
                            </>
                        )}
                    </td>
                </tr>
                <tr>
                    <td>Inclusions</td>
                    <td>
                        {props.loading ? (
                            <TextPlaceholder />
                        ) : (
                            <>
                                {displayBackupPaths(mod?.include)}
                            </>
                        )}
                    </td>
                </tr>
                <tr>
                    <td>Exclude CVS items</td>
                    <td>
                        <StatusBadge loading={props.loading} label={mod?.exclude_cvs ? "true" : "false"} status={mod?.exclude_cvs ? Const.OK : Const.CRITICAL} />
                    </td>
                </tr>
            </table>
        </Card>
    );
}

export default ModuleCard;