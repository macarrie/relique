import { Link, matchPath, useLocation, useMatch } from "react-router-dom";
import routes from "../routes";
import * as _ from "lodash";

function useCrumbs() {
    let location = useLocation();
    let default_breadcrumbs = [
        {
            name: "Home",
            path: "/"
        },
        {
            name: "Page not found",
            path: "404"
        }
    ]

    let matchedRoute = routes.find((r: any) => {
        if (r.path === "*") {
            return false;
        }

        let match = matchPath(r.path, location.pathname);
        return match !== null;
    }) || { path: "/404" };

    let urlMatch = useMatch(matchedRoute.path);
    if (urlMatch === null) {
        return default_breadcrumbs;
    }

    let currentPath = matchedRoute.path;
    return _.cloneDeep(routes).filter((route: any) => currentPath.includes(route.path))
        .map((item: any) => {
            let toReplace = item.path;
            // @ts-ignore
            Object.entries(urlMatch.params).forEach(([param, value], _) => {
                if (value !== undefined) {
                    toReplace = toReplace.replace(":" + param, value.toString())
                }
            });
            item.path = toReplace;
            return item;
        })
}

function Breadcrumb() {
    let breadcrumbs = useCrumbs();

    function renderCrumb(nav: any) {
        if (nav.path === "/") {
            return (
                <li key={nav.path}>
                    <div className="text-xl mr-2">
                        <i className="ri-home-2-line"></i>
                    </div>
                    <Link to={nav.path} className="">
                        {nav.name}
                    </Link>
                </li>
            );
        }

        return (
            <li aria-current="page" key={nav.path}>
                <Link to={nav.path} className="link-hover text-base-content/70">{nav.name}</Link>
            </li>
        )
    }

    function renderCrumbsLine(crumbs: any) {
        return <>
            {crumbs.map((item: any) => renderCrumb(item))}
        </>
    }

    return (
        <div className="breadcrumbs text-sm">
            <ul>
                {renderCrumbsLine(breadcrumbs)}
            </ul>
        </div>
    );
}

export default Breadcrumb;