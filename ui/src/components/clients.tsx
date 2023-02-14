import React from "react";

import ClientList from "../components/client_list"
import Card from "../components/card"

class Clients extends React.Component<any, any> {
    render() {
        return (
            <Card>
                <div className="flex flex-row px-4 py-3 items-center">
                    <span className="flex-grow text-xl font-bold">Clients</span>
                </div>
                <ClientList/>
            </Card>
        );
    }
}

export default Clients;
