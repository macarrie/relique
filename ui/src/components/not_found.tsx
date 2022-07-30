import React from 'react';

class NotFound extends React.Component<any, any> {
    render() {
        return (
            <div className="container grow flex justify-center items-center">
                <div className="text-6xl italic text-slate-300">
                    Not found
                </div>
            </div>
        );
    }
}

export default NotFound;
