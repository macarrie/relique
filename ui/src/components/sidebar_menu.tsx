import React from "react";

import SidebarMenuItem from "./sidebar_menu_item"

class SidebarMenu extends React.Component<any, any> {
    render() {
        return (
            <ul className="flex-grow space-y-1 pl-2">
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
                <SidebarMenuItem 
                    label="Schedules"
                    icon="ri-calendar-2-line"
                    link="/schedules"
                />
                <SidebarMenuItem 
                    label="Retention"
                    icon="ri-database-2-line"
                    link="/retention"
                />
            </ul>
        );
    }
}

export default SidebarMenu;
