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
    ssh_alive :number,
};

export default Client;

