import React from "react";

import ModuleList from "../components/module_list"

class Modules extends React.Component<any, any> {
    render() {
        return (
            <div className="grid grid-cols-4 gap-4">
                <div className="col-span-4 bg-white shadow rounded">
                    <div className="flex flex-row px-4 py-3 items-center">
                        <span className="flex-grow text-xl">Installed modules</span>
                    </div>
                    <ModuleList />
                </div>
            </div>
        );
    }
}

export default Modules;
