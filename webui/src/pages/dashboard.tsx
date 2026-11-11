import { useEffect, useState } from 'react';
import Config from "../types/config";
import Section from '../components/section';
import Breadcrumb from '../components/breadcrumb';
import { Link } from 'react-router-dom';
import API from '../utils/api';
import JobList from '../components/job_list';
import Job from '../types/job';
import Const from '../types/const';
import Module from '../types/module';
import Stat from '../components/stat';
import Utils from '../utils/utils';

function Dashboard() {
    const [latestJobs, setLatestJobs] = useState<Job[]>([]);
    const [latestJobsLoading, setLatestJobsLoading] = useState<boolean>(true);

    const [jobStats, setJobStats] = useState<Map<string, number>>(getJobSummaryCounts([]));
    const [jobStatsLoading, setJobStatsLoading] = useState<boolean>(true);

    const [imgStats, setImgStats] = useState<Map<string, number>>(new Map<string, number>);
    const [imgStatsLoading, setImgStatsLoading] = useState<boolean>(true);

    const [serverConfig, setConfig] = useState<Config>({} as Config);
    const [serverConfigLoading, setServerConfigLoading] = useState<boolean>(true);


    const [mods, setModules] = useState<Module[]>([]);
    const [modsLoading, setModsLoading] = useState<boolean>(true);

    function getJobSummaryCounts(jobs: Job[]): Map<string, number> {
        const recap = new Map<string, number>();

        recap.set("total", jobs.length);
        recap.set("running", jobs.filter(elt => elt["status"] === "active").length);
        recap.set("success", jobs.filter(elt => elt["status"] === "success").length);
        recap.set("incomplete", jobs.filter(elt => elt["status"] === "incomplete").length);
        recap.set("error", jobs.filter(elt => elt["status"] === "error").length);

        return recap;
    }

    useEffect(() => {
        function getLatestJobs(nb: number) {
            setLatestJobsLoading(true);
            API.jobs.list({ limit: nb }).then((response) => {
                setLatestJobs(response.data.data ?? []);
                setLatestJobsLoading(false);
            }).catch(error => {
                console.log("Cannot get job list", error);
                Utils.notify(Const.CRITICAL, "Latest jobs", "Cannot get latest job list: " + error)
                setLatestJobs([]);
            });
        }

        getLatestJobs(Const.DASHBOARD_NB_LATEST_JOBS);
    }, [])

    useEffect(() => {
        function getImageStats() {
            setImgStatsLoading(true);
            const imgStats = new Map<string, number>();
            API.images.stats().then((response) => {
                imgStats.set("count", response.data.count);
                imgStats.set("total_size", response.data.total_size);
                setImgStats(imgStats);
                setImgStatsLoading(false);
            }).catch(error => {
                console.log("Cannot get image stats", error);
                Utils.notify(Const.CRITICAL, "Image stats", "Cannot get image stats: " + error)
                setImgStats(imgStats);
            });
        }

        getImageStats();
    }, [])

    useEffect(() => {
        function getTodayJobs() {
            setJobStatsLoading(true);
            const yesterday = new Date(new Date().getTime() - (24 * 60 * 60 * 1000));
            API.jobs.list({ limit: 10000, after: yesterday.toISOString() }).then((response) => {
                setJobStats(getJobSummaryCounts(response.data.data ?? []));
                setJobStatsLoading(false);
            }).catch(error => {
                console.log("Cannot get job list", error);
                Utils.notify(Const.CRITICAL, "Daily jobs info", "Cannot get daily relique jobs: " + error)
                setJobStats(getJobSummaryCounts([]));
            });
        }

        getTodayJobs();
    }, [])

    useEffect(() => {
        function getConfig() {
            setServerConfigLoading(true);
            API.config.get().then((response) => {
                setConfig(response.data ?? {});
                setServerConfigLoading(false);
            }).catch(error => {
                console.log("Cannot get relique config", error);
                Utils.notify(Const.CRITICAL, "Backup policy", "Cannot get relique server config: " + error)
                setConfig({} as Config);
            });
        }

        getConfig();
    }, [])

    useEffect(() => {
        function getModuleList() {
            setModsLoading(true);
            API.modules.list({ limit: 10000 }).then((response) => {
                setModules(response.data.data ?? []);
                setModsLoading(false);
            }).catch(error => {
                console.log("Cannot get module list", error);
                Utils.notify(Const.CRITICAL, "Installed modules", "Cannot get relique modules installed on server: " + error)
                setModules([]);
            });
        }

        getModuleList();
    }, [])

    return (
        <>
            <Breadcrumb />
            <Section
                title="Jobs"
                description="last 24 hours"
                actions={<Link to="/jobs" className="text-link">All jobs</Link>}>
                <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-6">
                    <Stat to="/jobs" loading={jobStatsLoading} value={jobStats.get("total")} label="Total" />
                    <Stat tone="info" loading={jobStatsLoading} value={jobStats.get("running")} label="Running" />
                    <Stat tone="success" loading={jobStatsLoading} value={jobStats.get("success")} label="Success" />
                    <Stat tone="warning" loading={jobStatsLoading} value={jobStats.get("incomplete")} label="Incomplete" />
                    <Stat tone="error" loading={jobStatsLoading} value={jobStats.get("error")} label="Error" />
                </div>
            </Section>

            <div className="grid lg:grid-cols-2 gap-4">
                <Section title="Storage">
                    <div className="grid grid-cols-2 gap-6">
                        <Stat to="/images" loading={imgStatsLoading} value={imgStats.get("count")} label="Images" />
                        <Stat loading={imgStatsLoading} value={Utils.formatSize(imgStats.get("total_size") ?? 0)} label="Size on disk" />
                    </div>
                </Section>
                <Section title="Backup policy">
                    <div className="grid grid-cols-3 gap-6">
                        <Stat to="/clients" loading={serverConfigLoading} value={(serverConfig.clients ?? []).length} label="Clients" />
                        <Stat to="/repositories" loading={serverConfigLoading} value={(serverConfig.repositories ?? []).length} label="Repositories" />
                        <Stat to="/modules" loading={modsLoading} value={(mods ?? []).length} label="Modules" />
                    </div>
                </Section>
            </div>

            <Section
                title="Latest jobs"
                actions={<Link to="/jobs" className="text-link">All jobs</Link>}>
                <JobList
                    data={latestJobs}
                    loading={latestJobsLoading}
                    actions={false}
                    paginated={false}
                    sorted={false} />
            </Section>
        </>
    );
}

export default Dashboard;