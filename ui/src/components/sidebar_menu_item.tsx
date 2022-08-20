import React from "react";
import { NavLink } from "react-router-dom";


function SidebarMenuItem(props :any) {
    let activeClass = "block rounded-l bg-slate-50 text-blue-900 font-semibold"
    let inactiveClass = "block rounded-l text-slate-50 hover:bg-gray-100 hover:text-blue-900"

    return (
        <li>
            <NavLink to={props.link}
                className={({isActive}) => isActive ? activeClass : inactiveClass}>
                <span className="h-12 px-2 mr-2 md:px-6 flex flex items-center w-full">
                    <div className="flex flex-row text-base">
                        <div className="text-xl">
                            <i className={props.icon}></i>
                        </div>
                        <div className="hidden md:block ml-2 flex items-center">
                            {props.label}
                        </div>
                    </div>
                </span>
            </NavLink>
        </li>
    );
}

export default SidebarMenuItem;
