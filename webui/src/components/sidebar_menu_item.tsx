import { Link } from "react-router-dom";

function SidebarMenuItem({
    label = "Menu item",
    target = "/",
    icon = "ri-link",
}) {
    return (
        <li>
            <Link to={target} className="flex items-center p-2 rounded-lg hover:bg-base-300 group">
                <i className={`text-2xl text-base-content/50 group-hover:text-base-content ${icon}`}></i>
                <span className="ms-3">{label}</span>
            </Link>
        </li>
    );
}

export default SidebarMenuItem;