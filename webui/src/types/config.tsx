import Client from "./client";
import Module from "./module";
import Repository from "./repository";

type WebUIConfig = {
    bind_addr: string,
    port: number,
    ssl_cert: string,
    ssl_key: string,
};

type Config = {
    clients: Client[],
    repositories: Repository[],
    modules: Module[],
    webui: WebUIConfig,
    client_cfg_path: string,
    repo_cfg_path: string,
    module_install_path: string,
    db_path: string,
    catalog_path: string,
    current_file: string,
};

export default Config;
