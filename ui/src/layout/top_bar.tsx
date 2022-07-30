import React from "react";

import Breadcrumb from "../components/breadcrumb"

class TopBar extends React.Component<any, any> {
    render() {
        return (
            <div className="flex flex-row mb-4">
                <Breadcrumb />

                <div className="space-x-3 pr-3">
                    <button className="text-xl text-slate-400 hover:text-slate-900">
                        <i className="ri-search-line"></i>
                    </button>
                    <button className="text-xl text-slate-400">
                        <i className="ri-mail-line"></i>
                    </button>
                </div>
            </div>
        );
    }
}

export default TopBar;
