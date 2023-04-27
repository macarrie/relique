import Module from "./module"

type Client = {
    id :number,
    name :string,
    address :string,
    port :number,
    modules :Module[],
    version :string,
    server_address :string,
    server_port :number,
    api_alive :number,
    api_alive_message :string,
    ssh_alive :number,
    ssh_alive_message :string,
};

export default Client;

