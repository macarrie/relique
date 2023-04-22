import React from "react";
import { Link } from "react-router-dom";

import JobList from "../components/job_list"
import Card from "../components/card"

class Dashboard extends React.Component<any, any> {
    render() {
        return (
        <>
            <Card>
                <div className="flex flex-row px-4 py-3 items-center">
                    <span className="flex-grow text-xl">Latest jobs</span>
                    <Link to="/jobs" className="button button-small button-text">See all</Link>
                </div>
                <JobList limit={10}/>
            </Card>
        </>
    );
    }
}

export default Dashboard;
