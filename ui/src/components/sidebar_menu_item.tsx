import React from "react";
import { NavLink } from "react-router-dom";


function SidebarMenuItem(props :any) {
    let activeClass = "block bg-slate-50 text-blue-900 font-semibold"
    let inactiveClass = "block text-slate-50 hover:bg-gray-200 hover:text-blue-900"

    return (
        <li>
            <NavLink to={props.link}
                className={({isActive}) => isActive ? activeClass : inactiveClass}>
                    <div className="h-12 flex flex-row items-center">
                        <div className="w-12 block text-center text-xl">
                            <i className={props.icon}></i>
                        </div>
                        <div className="hidden md:block">
                            {props.label}
                        </div>
                    </div>
            </NavLink>
        </li>
    );
}

export default SidebarMenuItem;
