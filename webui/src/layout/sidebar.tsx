import { useEffect, useState } from "react";
import SidebarMenuItem from "../components/sidebar_menu_item";
import API from "../utils/api";

function Sidebar() {
    const [version, setVersion] = useState<string>("");

    useEffect(() => {
        API.config.get_version().then((response) => {
            setVersion(response.data.version ?? "unknown");
        }).catch(error => {
            console.log("Cannot get relique version", error);
            setVersion("unknown");
        });
    }, [])

    return (
        <aside className="fixed top-0 left-0 z-40 w-56 h-screen bg-base-200 border-r border-base-300 flex flex-col">
            <div className="flex items-center gap-2.5 px-5 pt-4 pb-3">
                <span className="grid place-items-center size-8 rounded bg-primary text-primary-content">
                    <i className="ri-archive-line text-lg"></i>
                </span>
                <div className="leading-tight">
                    <div className="font-display text-xl font-bold tracking-tight">Relique</div>
                    {version ? (
                        <div className="code text-xs text-muted" title="Relique version">{version}</div>
                    ) : (
                        <span className="skeleton inline-block h-2.5 w-16 align-middle"></span>
                    )}
                </div>
            </div>
            <nav className="flex-grow overflow-y-auto py-2">
                <ul className="space-y-0.5">
                    <SidebarMenuItem target="/dashboard" icon="ri-dashboard-line" label="Dashboard" />
                    <SidebarMenuItem target="/jobs" icon="ri-list-check-3" label="Jobs" />
                    <SidebarMenuItem target="/clients" icon="ri-computer-line" label="Clients" />
                    <SidebarMenuItem target="/modules" icon="ri-file-code-line" label="Modules" />
                    <SidebarMenuItem target="/images" icon="ri-stack-line" label="Images" />
                    <SidebarMenuItem target="/repositories" icon="ri-database-2-line" label="Repositories" />
                </ul>
            </nav>
            <nav className="py-2 border-t border-base-300">
                <ul>
                    <SidebarMenuItem target="/config" icon="ri-settings-3-line" label="Configuration" />
                </ul>
            </nav>
        </aside>
    )
}

export default Sidebar
