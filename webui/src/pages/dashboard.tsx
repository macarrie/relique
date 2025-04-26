import { useEffect, useState } from 'react';
import Card from '../components/card';
import { Link } from 'react-router-dom';
import API from '../utils/api';
import JobList from '../components/job_list';
import Job from '../types/job';
import Const from '../types/const';
import Module from '../types/module';
import DashboardStat from '../components/dashboard_stat';
import Utils from '../utils/utils';

function Dashboard() {
    let [latestJobs, setLatestJobs] = useState<Job[]>([]);
    let [jobStats, setJobStats] = useState<Map<string, number>>(getJobSummaryCounts([]));
    let [imgStats, setImgStats] = useState<Map<string, number>>(new Map<string, number>);
    let [serverConfig, setConfig] = useState<any>({});
    let [version, setVersion] = useState<String>("");
    let [mods, setModules] = useState<Module[]>([]);

    function getJobSummaryCounts(jobs: Job[]): Map<string, number> {
        let recap = new Map<string, number>();

        recap.set("total", jobs.length);
        recap.set("running", jobs.filter(elt => elt["status"] === "active").length);
        recap.set("success", jobs.filter(elt => elt["status"] === "success").length);
        recap.set("incomplete", jobs.filter(elt => elt["status"] === "incomplete").length);
        recap.set("error", jobs.filter(elt => elt["status"] === "error").length);

        return recap;
    }

    useEffect(() => {
        function getLatestJobs(nb: number) {
            API.jobs.list({ limit: nb }).then((response: any) => {
                setLatestJobs(response.data.data ?? []);
            }).catch(error => {
                console.log("Cannot get job list", error);
                setLatestJobs([]);
            });
        }

        getLatestJobs(Const.DASHBOARD_NB_LATEST_JOBS);
    }, [])

    useEffect(() => {
        function getImageStats() {
            let imgStats = new Map<string, number>();
            API.images.stats().then((response: any) => {
                imgStats.set("count", response.data.count);
                imgStats.set("total_size", response.data.total_size);
                setImgStats(imgStats);
            }).catch(error => {
                console.log("Cannot get image stats", error);
                setImgStats(imgStats);
            });
        }

        getImageStats();
    }, [])

    useEffect(() => {
        function getVersion() {
            API.config.get_version().then((response: any) => {
                setVersion(response.data.version ?? "unknown");
            }).catch(error => {
                console.log("Cannot get relique version", error);
                setVersion("unknown");
            });
        }

        getVersion();
    }, [])

    useEffect(() => {
        function getTodayJobs() {
            var yesterday = new Date(new Date().getTime() - (24 * 60 * 60 * 1000));
            API.jobs.list({ limit: 10000, after: yesterday.toISOString() }).then((response: any) => {
                setJobStats(getJobSummaryCounts(response.data.data ?? []));
            }).catch(error => {
                console.log("Cannot get job list", error);
                setJobStats(getJobSummaryCounts([]));
            });
        }

        getTodayJobs();
    }, [])

    useEffect(() => {
        function getConfig() {
            API.config.get().then((response: any) => {
                setConfig(response.data ?? {});
            }).catch(error => {
                console.log("Cannot get relique config", error);
                setConfig({});
            });
        }

        getConfig();
    }, [])

    useEffect(() => {
        function getModuleList() {
            API.modules.list({ limit: 10000 }).then((response: any) => {
                setModules(response.data.data ?? []);
            }).catch(error => {
                console.log("Cannot get job list", error);
                setModules([]);
            });
        }

        getModuleList();
    }, [])

    return (
        <>
            <div className="grid grid-cols-4 gap-4">
                <Card className="flex flex-2 col-span-3 items-center p-6">
                    <div className="flex-grow">
                        <h3 className="font-bold">
                            Jobs summary <span className="font-normal italic text-base-content/50">(last 24h)</span>
                        </h3>
                        <div className="stats">
                            <Link to="/jobs">
                                <DashboardStat color="text-base-content" value={(jobStats ?? {}).get("total")} label="Total" />
                            </Link>
                            <DashboardStat color="text-primary" value={jobStats.get("running")} label="Running" />
                            <DashboardStat color="text-success" value={jobStats.get("success")} label="Success" />
                            <DashboardStat color="text-warning" value={jobStats.get("incomplete")} label="Incomplete" />
                            <DashboardStat color="text-error" value={jobStats.get("error")} label="Error" />
                        </div>
                    </div>
                    <div className="flex items-center">
                        <i className="text-6xl text-base-content/20 ri-list-check-3"></i>
                    </div>
                </Card>
                <Card className="flex flex-2 items-center p-6">
                    <div className="flex-grow">
                        <h3 className="font-bold">
                            Relique version
                        </h3>
                        <div className="stats">
                            <DashboardStat color="text-secondary" value={version} label="" />
                        </div>
                    </div>
                </Card>
                <Card className="flex flex-2 col-span-2 items-center p-6">
                    <div className="flex-grow">
                        <h3 className="font-bold">
                            Images
                        </h3>
                        <div className="stats">
                            <Link to="/images">
                                <DashboardStat value={imgStats.get("count")} label="Images" />
                            </Link>
                            <DashboardStat value={Utils.formatSize(imgStats.get("total_size") ?? 0)} label="Size on disk" />
                        </div>
                    </div>
                    <div className="flex items-center">
                        <i className="text-6xl text-base-content/20 ri-stack-fill"></i>
                    </div>
                </Card>
                <Card className="flex flex-2 col-span-2 items-center p-6">
                    <div className="flex-grow">
                        <h3 className="font-bold">
                            Backup policy
                        </h3>
                        <div className="stats">
                            <Link to="/clients">
                                <DashboardStat value={(serverConfig.clients ?? []).length} label="Clients" />
                            </Link>
                            <Link to="/repositories">
                                <DashboardStat value={(serverConfig.repositories ?? []).length} label="Repositories" />
                            </Link>
                            <Link to="/modules">
                                <DashboardStat value={(mods ?? []).length} label="Installed modules" />
                            </Link>
                        </div>
                    </div>
                    <div className="flex items-center">
                        <i className="text-6xl text-base-content/20 ri-shield-check-fill"></i>
                    </div>
                </Card>
            </div>

            <Card>
                <JobList
                    title="Latest jobs"
                    custom_actions={[
                        <Link to="/jobs" className='link-primary link-hover text-sm'>See more</Link>
                    ]}
                    data={latestJobs}
                    actions={false}
                    paginated={false}
                    sorted={false} />
            </Card>
        </>
    );
}

export default Dashboard;