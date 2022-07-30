import React from "react";

import API from "../utils/api";

import Client from "../types/client";
import Module from "../types/module";

import DotStatus from "./dot_status";

class ClientListRow extends React.Component<any, any> {
    render() {
        let client = this.props.client;
        let module_names :string[] = client.modules.map((mod :Module) => mod.name)

        return (
            <tr className="hover:bg-slate-50">
                <td className="py-2 px-3"><DotStatus value={client.alive} /></td>
                <td className="py-2 px-3">{client.name}</td>
                <td className="py-2 px-3">{client.address}</td>
                <td className="py-2 px-3">{module_names.join(", ")}</td>
            </tr>
        );
    }
}

type State = {
    clients :Client[]
};

class ClientList extends React.Component<any, State> {
    get_clients_interval :number;
    limit :number;

    constructor(props: any) {
        super(props);

        this.get_clients_interval = 0;
        this.limit = this.props.limit ? this.props.limit : 0;
    }

    state :State = {
        clients: [],
    };

    componentDidMount() {
        this.getClients();
    }

    componentWillUnmount() {}

    getClients() {
        API.clients.list({
            limit: this.limit,
        }).then((response :any) => {
            console.log(this.state.clients)
            this.setState({
                clients: response.data,
            });
        }).catch(error => {
            console.log("Cannot get client list", error);
        });
    }

    renderClientList() {
        if (!this.state.clients) {
            return (
                <>
                    Loading
                </>
            )
        }

        const clientList = this.state.clients.map((client :Client) =>
            <ClientListRow key={client.id} client={client}/>
        );

        return (
            <tbody>
            {clientList}
            </tbody>
        )
    }

    render() {
        return (
            <table className="table-auto w-full">
                <thead className="bg-slate-50 uppercase text-slate-500 text-left">
                <tr className="border border-l-0 border-r-0 border-slate-100">
                    <th className="py-2 px-3 max-w-min text-center">Alive</th>
                    <th className="py-2 px-3">Name</th>
                    <th className="py-2 px-3">Address</th>
                    <th className="py-2 px-3">Modules</th>
                </tr>
                </thead>
                {this.renderClientList()}
            </table>
        );
    }
}

export default ClientList;
