import { useEffect, useState } from "react";
import Config from "../types/config";

import API from "../utils/api";
import Const from "../types/const";
import Utils from "../utils/utils";
import PageHeader from "../components/page_header";
import Section from "../components/section";
import PropertyList from "../components/property_list";

function ConfigPage() {
    const [serverConfig, setConfig] = useState<Config>({} as Config);
    const [loading, setLoading] = useState<boolean>(true);

    function reloadConfig() {
        API.config.reload().then(() => {
            Utils.notify(Const.OK, "Configuration reloaded", "")
            getConfig();
        }).catch(error => {
            console.log("Cannot reload config", error);
            Utils.notify(Const.CRITICAL, "Cannot reload configuration", error.toString())
        });
    }

    function getConfig() {
        setLoading(true);
        API.config.get().then((response) => {
            setConfig(response.data ?? {});
            setLoading(false);
        }).catch(error => {
            console.log("Cannot get relique config", error);
            Utils.notify(Const.CRITICAL, "Cannot get config from server", error.toString())
            setConfig({} as Config);
        });
    }

    useEffect(() => {
        getConfig();
    }, [])

    const cfg = serverConfig ?? {};

    return (
        <>
            <PageHeader
                title="Configuration"
                actions={
                    <button className="btn btn-sm" onClick={reloadConfig}>
                        <i className="ri-refresh-line"></i>
                        Reload configuration
                    </button>
                }
            />
            <Section title="Paths" description="Relative paths are resolved from the directory of the configuration file.">
                <PropertyList loading={loading} items={[
                    { label: "Configuration file", value: cfg.current_file, mono: true },
                    { label: "Database", value: cfg.db_path, mono: true },
                    { label: "Catalog", value: cfg.catalog_path, mono: true },
                    { label: "Clients", value: cfg.client_cfg_path, mono: true },
                    { label: "Repositories", value: cfg.repo_cfg_path, mono: true },
                    { label: "Modules", value: cfg.module_install_path, mono: true },
                ]} />
            </Section>
            <Section title="Web UI">
                <PropertyList loading={loading} items={[
                    { label: "Bind address", value: cfg.webui?.bind_addr, mono: true },
                    { label: "Port", value: cfg.webui?.port, mono: true },
                    { label: "SSL certificate", value: cfg.webui?.ssl_cert, mono: true },
                    { label: "SSL key", value: cfg.webui?.ssl_key, mono: true },
                ]} />
            </Section>
        </>
    );
}

export default ConfigPage;