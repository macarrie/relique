import { Link, matchPath, useLocation, useMatch } from "react-router-dom";
import routes from "../routes";

type Crumb = {
    name: string,
    path: string,
};

function useCrumbs(): Crumb[] {
    const location = useLocation();
    const default_breadcrumbs: Crumb[] = [
        {
            name: "Home",
            path: "/"
        },
        {
            name: "Page not found",
            path: "404"
        }
    ]

    const matchedRoute = routes.find((r) => {
        if (r.path === "*") {
            return false;
        }

        const match = matchPath(r.path, location.pathname);
        return match !== null;
    }) || { path: "/404" };

    const urlMatch = useMatch(matchedRoute.path);
    if (urlMatch === null) {
        return default_breadcrumbs;
    }

    const currentPath = matchedRoute.path;
    return routes.filter((route) => currentPath.includes(route.path))
        .map((route): Crumb => {
            let path = route.path;
            Object.entries(urlMatch.params).forEach(([param, value]) => {
                if (value !== undefined) {
                    path = path.replace(":" + param, value.toString())
                }
            });
            return { name: route.name, path: path };
        })
}

function Breadcrumb() {
    const breadcrumbs = useCrumbs();

    function renderCrumb(nav: Crumb) {
        return (
            <li key={nav.path}>
                <Link to={nav.path} className="hover:text-base-content hover:underline underline-offset-2">{nav.name}</Link>
            </li>
        )
    }

    function renderCrumbsLine(crumbs: Crumb[]) {
        return <>
            {crumbs.map((item) => renderCrumb(item))}
        </>
    }

    return (
        <div className="breadcrumbs p-0 text-xs text-muted">
            <ul>
                {renderCrumbsLine(breadcrumbs)}
            </ul>
        </div>
    );
}

export default Breadcrumb;