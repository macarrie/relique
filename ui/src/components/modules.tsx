import React from "react";

import ModuleList from "../components/module_list"
import Card from "../components/card"

class Modules extends React.Component<any, any> {
    render() {
        return (
            <Card>
                <div className="flex flex-row px-4 py-3 items-center">
                    <span className="flex-grow text-xl font-bold">Installed modules</span>
                </div>
                <ModuleList/>
            </Card>
        );
    }
}

export default Modules;
