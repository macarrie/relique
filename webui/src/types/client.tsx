import Module from "./module"

type Client = {
    id: number,
    name: string,
    address: string,
    ssh_user: string,
    ssh_port: number,
    modules: Module[],
    ssh_alive: number,
    ssh_alive_message: string,
    state_is_loading: boolean,
};

export default Client;