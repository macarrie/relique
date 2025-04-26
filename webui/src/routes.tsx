import Dashboard from "./pages/dashboard";
import Jobs from "./pages/jobs";
import JobDetails from "./pages/job_details";
import NotFound from "./pages/not_found";
import ClientDetails from "./pages/client_details";
import Clients from "./pages/clients";
import Modules from "./pages/modules";
import ModuleDetails from "./pages/module_details";
import Images from "./pages/images";
import ImageDetails from "./pages/image_details";
import Repositories from "./pages/repositories";
import RepositoryDetails from "./pages/repository_details";

const routes = [
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
        path: "/clients/:client_name",
        name: "Client details",
        elt: () => <ClientDetails />
    },
    {
        path: "/modules",
        name: "All modules",
        elt: () => <Modules />
    },
    {
        path: "/modules/:module_name",
        name: "Module details",
        elt: () => <ModuleDetails />
    },
    {
        path: "/images",
        name: "All images",
        elt: () => <Images />
    },
    {
        path: "/images/:img_uuid",
        name: "Image details",
        elt: () => <ImageDetails />
    },
    {
        path: "/repositories",
        name: "All repositories",
        elt: () => <Repositories />
    },
    {
        path: "/repositories/:repo_name",
        name: "Repository details",
        elt: () => <RepositoryDetails />
    },
    {
        path: "*",
        name: "Not found",
        elt: () => <NotFound />
    },
];

export default routes;