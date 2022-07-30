import React from "react";

class Breadcrumb extends React.Component<any, any> {
    render() {
        return (
            <nav className="flex-grow flex" aria-label="breadcrumb">
                <ol className="inline-flex items-center space-x-1">
                    <li className="inline-flex items-center">
                        <a href="#" className="inline-flex items-center text-sm font-medium text-gray-700 hover:text-blue-700">
                            <div className="text-xl mr-1">
                                <i className="ri-home-2-line"></i>
                            </div>
                            Home
                        </a>
                    </li>
                    <li aria-current="page">
                        <div className="flex items-center">
                            <i className="text-2xl text-gray-400 ri-arrow-right-s-line"></i>
                            <span className="ml-1 text-sm font-medium text-gray-400 md:ml-2">Overview</span>
                        </div>
                    </li>
                </ol>
            </nav>
        );
    }
}

export default Breadcrumb;
