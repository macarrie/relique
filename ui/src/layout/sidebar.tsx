import React from "react";

import Logo from "../components/logo"
import SidebarMenu from "../components/sidebar_menu"

function Sidebar() {
    return (
        <aside className={`md:w-48 w-12 fixed top-0 left-0 bottom-0 h-screen flex flex-col overflow-y-auto bg-blue-900 text-slate-50`}>
            <div className="flex flex-row items-center text-center py-4 mx-auto mb-4">
                <Logo />
            </div>

            <SidebarMenu />
        </aside>
    );
}

export default Sidebar;

//<aside className={`w-48 fixed top-0 left-0 h-screen flex flex-col overflow-y-auto bg-blue-900 text-slate-50`}>
    //<div className="flex flex-row items-center text-center py-4 mx-auto mb-4">
        //<Logo />
    //</div>

    //<SidebarMenu />
//</aside>
