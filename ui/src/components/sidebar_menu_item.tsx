import React from "react";
import { NavLink } from "react-router-dom";


function SidebarMenuItem(props :any) {
    let activeClass = "block rounded-lg mx-2 text-slate-700 bg-slate-200 dark:bg-slate-800 dark:text-slate-200 font-semibold hover:text-slate-700 dark:hover:text-slate-200 dark:hover:bg-slate-700"
    let inactiveClass = "block rounded-lg mx-2 text-slate-500 dark:text-slate-400 hover:bg-slate-200 hover:text-slate-700 dark:hover:bg-slate-800 dark:hover:text-slate-200"

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
