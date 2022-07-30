import React from "react";

import Logo from "../components/logo"
import SidebarMenu from "../components/sidebar_menu"

class Sidebar extends React.Component<any, any> {
    render() {
        return (
            <aside className="fixed top-0 left-0 w-72 h-screen flex flex-col overflow-y-auto bg-blue-900 text-slate-50">
                <div className="flex flex-row items-center text-center py-4 mb-4">
                    <button className="flex-none text-xl mr-2 hover:text-slate-300 hover:bg-blue-800 rounded-full w-10 h-10 ml-2">
                        <i className="ri-menu-line"></i>
                    </button>

                    <Logo />
                </div>

                <SidebarMenu />
            </aside>
        );
    }
}

export default Sidebar;
