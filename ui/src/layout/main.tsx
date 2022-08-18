import React from "react";

import TopBar from "./top_bar"

class Main extends React.Component<any, any> {
    render() {
        return (
            <div className="relative ml-12 md:ml-48 flex flex-col container bg-slate-50 py-4 px-4">
                <TopBar />
                <hr className="mb-4" />

                {this.props.children}
            </div>
        );
    }
}

export default Main;
