import React from "react";

import JobList from "../components/job_list"
import Card from "../components/card"

class Dashboard extends React.Component<any, any> {
    render() {
        return (
            <Card>
                <JobList title="Latest jobs" limit={10} />
            </Card>
        );
    }
}

export default Dashboard;
