import { NavLink } from "react-router-dom";

function SidebarMenuItem({
    label = "Menu item",
    target = "/",
    icon = "ri-link",
}) {
    return (
        <li>
            <NavLink
                to={target}
                className={({ isActive }) =>
                    `flex items-center gap-3 py-1.5 pl-[1.125rem] pr-3 border-l-2 group ${isActive
                        ? "border-primary bg-base-300/50 font-semibold text-base-content"
                        : "border-transparent text-base-content/80 hover:bg-base-300/50 hover:text-base-content"
                    }`
                }>
                <i className={`text-lg text-muted group-hover:text-base-content ${icon}`}></i>
                <span>{label}</span>
            </NavLink>
        </li>
    );
}

export default SidebarMenuItem;
