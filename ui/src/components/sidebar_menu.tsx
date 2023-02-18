import React from "react";

import SidebarMenuItem from "./sidebar_menu_item"
import DarkModeSwitcher from "./dark_mode_switcher";

function SidebarMenu() {
    return (
        <>
            <ul className="grow space-y-1">
                <SidebarMenuItem
                    label="Overview"
                    icon="ri-home-2-line"
                    link="/dashboard"
                />
                <SidebarMenuItem
                    label="Jobs"
                    icon="ri-task-line"
                    link="/jobs"
                />
                <SidebarMenuItem
                    label="Clients"
                    icon="ri-device-line"
                    link="/clients"
                />
                <SidebarMenuItem
                    label="Modules"
                    icon="ri-apps-line"
                    link="/modules"
                />
            </ul>

            <ul>
                <li className="px-2 pb-2 flex flex-row items-center">
                    <div
                        className="w-10 h-10 mr-2 bg-slate-400 dark:bg-slate-700 dark:text-slate-300 rounded-full flex items-center justify-center text-slate-100 text-xl">
                        <div className="">A</div>
                    </div>
                    <span className="text-slate-500 dark:text-slate-300 flex-grow">Admin</span>
                    <DarkModeSwitcher/>
                </li>
            </ul>
        </>
    );
}

export default SidebarMenu;
