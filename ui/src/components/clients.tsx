import React from "react";

import ClientList from "../components/client_list"
import Card from "../components/card"

class Clients extends React.Component<any, any> {
    render() {
        return (
            <Card>
                <ClientList title="Clients" filtered sorted paginated />
            </Card>
        );
    }
}

export default Clients;
