import React from "react";

import ClientList from "../components/client_list"

class Clients extends React.Component<any, any> {
    render() {
        return (
            <div className="grid grid-cols-4 gap-4">
                <div className="col-span-4 bg-white shadow rounded">
                    <div className="flex flex-row px-4 py-3 items-center">
                        <span className="flex-grow text-xl">Clients</span>
                    </div>
                    <ClientList />
                </div>
            </div>
        );
    }
}

export default Clients;
