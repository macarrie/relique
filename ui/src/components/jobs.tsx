import React from "react";

import JobList from "../components/job_list"
import Card from "../components/card"

class Jobs extends React.Component<any, any> {
    render() {
        return (
            <Card>
                <JobList title="Backup/restore jobs" filtered sorted />
            </Card>
        )
        // return (
        //     <Card>
        //         <div className="col-span-4">
        //             <div className="flex flex-row px-4 py-3 items-center">
        //                 <span className="flex-grow text-xl font-bold">Backup/restore jobs</span>
        //             </div>
        //             <JobList/>
        //         </div>
        //     </Card>
        // );
    }
}

export default Jobs;
