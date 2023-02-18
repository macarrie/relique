import React from "react";

import BurgerMenu from "../components/burger_menu"
import Breadcrumb from "../components/breadcrumb"

class TopBar extends React.Component<any, any> {
    render() {
        return (
            <div className="container flex flex-row mb-4">
                <BurgerMenu/>
                <Breadcrumb/>
            </div>
        );
    }
}

export default TopBar;
