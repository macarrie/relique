import React from "react";

class Logo extends React.Component<any, any> {
    render() {
        return (
            <span className="text-3xl font-serif font-bold text-yellow-300 hover:text-yellow-500 hover:no-underline">
                <div className="flex flex-row items-center justify-center">
                    <i className="ri-trophy-line mr-2"></i>
                    <span>Relique</span>
                </div>
            </span>
        );
    }
}

export default Logo;
