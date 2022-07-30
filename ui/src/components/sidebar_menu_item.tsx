import React from "react";
import { NavLink } from "react-router-dom";


class SidebarMenuItem extends React.Component<any, any> {
    render() {
        let activeClass = "block rounded-l bg-slate-50 text-blue-900 font-semibold"
        let inactiveClass = "block rounded-l hover:bg-gray-100 hover:text-blue-900"

        return (
            <li>
                <NavLink to={this.props.link}
                    className={({isActive}) => isActive ? activeClass : inactiveClass}>
                    <span className="h-12 px-6 flex flex items-center w-full">
                        <div className="flex flex-row text-base">
                            <div className="text-xl">
                                <i className={this.props.icon}></i>
                            </div>
                            <div className="ml-2 flex items-center">
                                {this.props.label}
                            </div>
                        </div>
                    </span>
                </NavLink>
            </li>
        );
    }
}

export default SidebarMenuItem;
