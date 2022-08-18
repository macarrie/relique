import React from "react";
import { Link } from "react-router-dom";

import JobList from "../components/job_list"

class Dashboard extends React.Component<any, any> {
    render() {
        return (
        <>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
                <div className="bg-white shadow px-4 py-6 rounded">
                    <span className="flex-col">
                        <div className="text-center font-serif font-bold text-slate-700 text-5xl mb-4">82</div>
                        <div className="text-center uppercase text-slate-400">Running jobs</div>
                    </span>
                </div>
                <div className="bg-white shadow px-4 py-6 rounded">
                    <div className="flex flex-col">
                        <div className="text-center font-serif font-bold text-red-700 text-5xl mb-4">13</div>
                        <div className="text-center uppercase text-slate-400">Job errors <span className="italic text-sm text-slate-300">(last 24h)</span></div>
                    </div>
                </div>
                <div className="bg-white shadow px-4 py-6 rounded">
                    <div className="flex flex-row">
                        <span className="flex-grow">
                            <div className="text-center font-serif font-bold text-slate-700 text-5xl mb-4">17</div>
                            <div className="text-center uppercase text-slate-400">Clients</div>
                        </span>
                        <div className="flex flex-col justify-around pr-4">
                            <div className="text-center uppercase text-blue-700 text-l">14 active</div>
                            <div className="text-center uppercase text-red-700 text-l">3 inactive</div>
                        </div>
                    </div>
                </div>
                <div className="bg-white shadow px-4 py-6 rounded">
                    <div className="flex flex-col">
                        <div className="text-center font-serif font-bold text-slate-700 text-5xl mb-4">7</div>
                        <div className="text-center uppercase text-slate-400">Installed modules</div>
                    </div>
                </div>
            </div>

            <div className="bg-white shadow rounded">
                <div className="flex flex-row px-4 py-3 items-center">
                    <span className="flex-grow text-xl">Latest jobs</span>
                    <Link to="/jobs" className="bg-blue-500 rounded px-2 py-1 text-slate-50 hover:bg-blue-700 hover:text-slate-50 hover:no-underline uppercase text-xs font-bold">See all</Link>
                </div>
                <JobList limit={5} />
            </div>
        </>
    );
    }
}

export default Dashboard;
