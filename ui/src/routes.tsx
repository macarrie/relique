import Dashboard from "./components/dashboard";
import Jobs from "./components/jobs";
import JobDetails from "./components/job_details";
import Clients from "./components/clients";
import ClientDetails from "./components/client_details";
import Modules from "./components/modules";
import NotFound from "./components/not_found";
import React from "react";
import ModuleDetails from "./components/module_details";

let routes = [
    {
        path: "/",
        name: "Home",
        elt: () => <Dashboard />
    },
    {
        path: "/dashboard",
        name: "Overview",
        elt: () => <Dashboard />
    },
    {
        path: "/jobs",
        name: "All jobs",
        elt: () => <Jobs />
    },
    {
        path: "/jobs/:job_uuid",
        name: "Job details",
        elt: () => <JobDetails />
    },
    {
        path: "/clients",
        name: "All clients",
        elt: () => <Clients />
    },
    {
        path: "/clients/:client_id",
        name: "Client details",
        elt: () => <ClientDetails />
    },
    {
        path: "/modules",
        name: "All installed modules",
        elt: () => <Modules />
    },
    {
        path: "/modules/:name",
        name: "Module details",
        elt: () => <ModuleDetails />
    },
    {
        path: "*",
        name: "Not found",
        elt: () => <NotFound />
    },
];

export default routes;