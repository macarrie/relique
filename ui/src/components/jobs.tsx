import React from "react";

import JobList from "../components/job_list"

class Jobs extends React.Component<any, any> {
    render() {
        return (
            <div className="grid grid-cols-4 gap-4">
                <div className="col-span-4 bg-white shadow rounded">
                    <div className="flex flex-row px-4 py-3 items-center">
                        <span className="flex-grow text-xl">Backup/restore jobs</span>
                    </div>
                    <JobList />
                </div>
            </div>
        );
    }
}

export default Jobs;
