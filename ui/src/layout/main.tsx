import React from "react";

import TopBar from "./top_bar"

class Main extends React.Component<any, any> {
    render() {
        return (
            <div className="relative w-full ml-12 md:ml-48 flex flex-col bg-slate-50 py-4 px-4">
                <TopBar />
                <hr className="mb-4" />

                <div className="container">
                    {this.props.children}
                </div>
            </div>
        );
    }
}

export default Main;
