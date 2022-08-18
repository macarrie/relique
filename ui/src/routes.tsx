import Dashboard from "./components/dashboard";
import Jobs from "./components/jobs";
import Clients from "./components/clients";
import ClientDetails from "./components/client_details";
import NotFound from "./components/not_found";
import React from "react";

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
        path: "*",
        name: "Not found",
        elt: () => <NotFound />
    },
];

export default routes;