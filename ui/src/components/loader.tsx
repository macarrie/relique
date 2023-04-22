import React from "react";

function Loader() {
    return (
        <div className="flex flex-row justify-center">
            <i className="animate-spin mr-2 ri-loader-4-line"></i>
            <span>Loading</span>
        </div>
    );
}

export default Loader;
