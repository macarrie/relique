import React from "react";
import { Link } from "react-router-dom";

import JobList from "../components/job_list"
import Card from "../components/card"

class Dashboard extends React.Component<any, any> {
    render() {
        return (
        <>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
                <Card className="py-4">
                    <div className="flex-col">
                        <div className="text-center font-display text-slate-700 text-5xl mb-2">82</div>
                        <div className="text-center font-bold text-slate-400">Running jobs</div>
                    </div>
                </Card>
                <Card className="py-4">
                    <div className="flex flex-col">
                        <div className="text-center font-display text-red-700 text-5xl mb-2">13</div>
                        <div className="text-center font-bold text-slate-400">Job errors <span
                            className="italic text-sm text-slate-300">(last 24h)</span></div>
                    </div>
                </Card>
                <Card className="py-4">
                    <div className="flex flex-row">
                        <span className="flex-grow">
                            <div className="text-center font-display text-slate-700 text-5xl mb-2">17</div>
                            <div className="text-center font-bold text-slate-400">Clients</div>
                        </span>
                        <div className="flex flex-col justify-around pr-4">
                            <div className="text-center text-emerald-700 text-l">14 active</div>
                            <div className="text-center text-red-700 text-l">3 inactive</div>
                        </div>
                    </div>
                </Card>
                <Card className="py-4">
                    <div className="flex flex-col">
                        <div className="text-center font-display text-slate-700 text-5xl mb-2">7</div>
                        <div className="text-center font-bold text-slate-400">Installed modules</div>
                    </div>
                </Card>
            </div>

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
