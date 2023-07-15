import React from "react";

import JobList from "../components/job_list"
import Card from "../components/card"

class Jobs extends React.Component<any, any> {
    render() {
        return (
            <Card>
                <JobList title="Backup/restore jobs" filtered sorted paginated />
            </Card>
        )
    }
}

export default Jobs;
