import React from "react";

import Logo from "../components/logo"
import SidebarMenu from "../components/sidebar_menu"

function Sidebar(props: any) {
    if (props.mobile) {
        return (
            <aside className={`h-full flex flex-col overflow-y-auto`}>
                <SidebarMenu/>
            </aside>
        );
    }

    return (
        <aside className={`md:w-64 w-12 hidden md:flex flex-col overflow-y-auto`}>
            <div className="flex flex-row items-center text-center py-4 mx-auto mb-4">
                <Logo/>
            </div>

            <SidebarMenu/>
        </aside>
    );
}

export default Sidebar;