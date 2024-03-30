import React from "react";

import BurgerMenu from "../components/burger_menu"
import Breadcrumb from "../components/breadcrumb"

function TopBar(props :any) {
    return (
        <div className="container flex flex-row mb-4">
        <BurgerMenu sidebarOpen={props.sidebarOpen} setSidebarOpen={props.setSidebarOpen} />
        <Breadcrumb/>
        </div>
    );
}

export default TopBar;
