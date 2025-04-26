import SidebarMenuItem from "../components/sidebar_menu_item";

function Sidebar() {
    return (
        <>
            <aside className="fixed top-0 left-0 z-40 w-64 h-screen bg-base-200">
                <div className="flex items-center m-4 justify-center">
                    <i className="border-2 border-white shadow text-xl ri-trophy-line mr-2 px-2 py-1 rounded-full bg-gradient-to-tr from-secondary to-accent text-white"></i>
                    <span className="text-transparent text-xl font-bold whitespace-nowrap bg-gradient-to-tr from-secondary to-accent bg-clip-text">Relique</span>
                </div>
                <div className="h-full px-3 pt-2 overflow-y-auto">
                    <ul className="space-y-2 font-medium">
                        <SidebarMenuItem target="/dashboard" icon="ri-dashboard-fill" label="Dashboard" />
                        <SidebarMenuItem target="/jobs" icon="ri-list-check-3" label="Jobs" />
                        <SidebarMenuItem target="/clients" icon="ri-device-fill" label="Clients" />
                        <SidebarMenuItem target="/modules" icon="ri-file-code-fill" label="Modules" />
                        <SidebarMenuItem target="/images" icon="ri-stack-fill" label="Images" />
                        <SidebarMenuItem target="/repositories" icon="ri-database-2-fill" label="Repositories" />
                    </ul>
                </div>
            </aside>
        </>
    )
}

export default Sidebar