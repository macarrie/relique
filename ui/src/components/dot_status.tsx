import React from "react";

class DotStatus extends React.Component<any, any> {
    render() {
        return (
            <div className={`w-3 h-3 ${this.props.value ? "bg-green-500" : "bg-red-500"} rounded-full m-auto`}></div>
        );
    }
}
export default DotStatus;