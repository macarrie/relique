import React from "react";

import ModuleList from "../components/module_list"
import Card from "../components/card"

class Modules extends React.Component<any, any> {
    render() {
        return (
            <Card>
                <ModuleList title="Installed modules" filtered sorted/>
            </Card>
        );
    }
}

export default Modules;
