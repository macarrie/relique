import { useEffect, useState } from "react";

import Card from "../components/card";
import API from "../utils/api";
import Const from "../types/const";
import Utils from "../utils/utils";
import TextPlaceholder from "../components/text_placeholder";

function Config() {
    let [serverConfig, setConfig] = useState<any>({});
    let [loading, setLoading] = useState<boolean>(true);

    function reloadConfig() {
        API.config.reload().then((response: any) => {
            Utils.notify(Const.OK, "Config reloaded", "Configuration has been reloaded on server")
            getConfig();
        }).catch(error => {
            console.log("Cannot create client", error);
            Utils.notify(Const.CRITICAL, "Cannot create client", error.toString())
        });
    }

    function getConfig() {
        setLoading(true);
        API.config.get().then((response: any) => {
            setConfig(response.data ?? {});
            setLoading(false);
        }).catch(error => {
            console.log("Cannot get relique config", error);
            Utils.notify(Const.CRITICAL, "Cannot get config from server", error.toString())
            setConfig({});
        });
    }

    useEffect(() => {
        getConfig();
    }, [])

    return (
        <>
            <Card>
                <div className="p-4 flex flex-row items-center mb-2">
                    <div className="flex-grow font-bold align-middle">Configuration paths</div>
                    <div className="btn btn-sm" onClick={reloadConfig}>Reload config</div>
                </div>
                <div role="alert" className="alert alert-info alert-soft mx-4 mb-2">
                    <span>Non absolute paths are relative to the configuration file used</span>
                </div>
                <table className="table">
                    <tbody>
                        <tr>
                            <td>Current configuration file</td>
                            <td className="code">
                                {loading && (
                                    <TextPlaceholder />
                                )}
                                {serverConfig?.current_file}
                            </td>
                        </tr>
                        <tr>
                            <td>Database</td>
                            <td className="code">
                                {loading && (
                                    <TextPlaceholder />
                                )}
                                {serverConfig?.client_cfg_path}
                            </td>
                        </tr>
                        <tr>
                            <td>Catalog</td>
                            <td className="code">
                                {loading && (
                                    <TextPlaceholder />
                                )}
                                {serverConfig?.catalog_path}
                            </td>
                        </tr>
                        <tr>
                            <td>Clients</td>
                            <td className="code">
                                {serverConfig?.client_cfg_path}
                            </td>
                        </tr>
                        <tr>
                            <td>Repositories</td>
                            <td className="code">
                                {loading && (
                                    <TextPlaceholder />
                                )}
                                {serverConfig?.repo_cfg_path}
                            </td>
                        </tr>
                        <tr>
                            <td>Modules</td>
                            <td className="code">
                                {loading && (
                                    <TextPlaceholder />
                                )}
                                {serverConfig?.module_install_path}
                            </td>
                        </tr>
                    </tbody>
                </table>
            </Card>
            <Card>
                <div className="p-4 flex flex-row items-center mb-2">
                    <div className="flex-grow font-bold align-middle">Web UI settings</div>
                </div>
                <table className="table">
                    <tbody>
                        <tr>
                            <td>Bind address</td>
                            <td className="code">
                                {loading && (
                                    <TextPlaceholder />
                                )}
                                {serverConfig?.webui?.bind_addr}
                            </td>
                        </tr>
                        <tr>
                            <td>Port</td>
                            <td className="code">
                                {loading && (
                                    <TextPlaceholder />
                                )}
                                {serverConfig?.webui?.port}
                            </td>
                        </tr>
                        <tr>
                            <td>SSL cert</td>
                            <td className="code">
                                {loading && (
                                    <TextPlaceholder />
                                )}
                                {serverConfig?.webui?.ssl_cert}
                            </td>
                        </tr>
                        <tr>
                            <td>SSL key</td>
                            <td className="code">
                                {loading && (
                                    <TextPlaceholder />
                                )}
                                {serverConfig?.webui?.ssl_key}
                            </td>
                        </tr>
                    </tbody>
                </table>
            </Card>
        </>
    );
}

export default Config;