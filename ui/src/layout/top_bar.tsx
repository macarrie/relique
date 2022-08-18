import React from "react";

import Breadcrumb from "../components/breadcrumb"

class TopBar extends React.Component<any, any> {
    render() {
        return (
            <div className="flex flex-row mb-4">
                <Breadcrumb />
            </div>
        );
    }
}

export default TopBar;
