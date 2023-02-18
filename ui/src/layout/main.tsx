import React from "react";

import TopBar from "./top_bar"

function Main(props: any) {
    return (
        <div className="w-full flex flex-grow flex-col max-h-screen overflow-y-scroll py-4 px-4">
            <TopBar/>

            <div className="container">
                {props.children}
            </div>
        </div>
    );
}

export default Main;
