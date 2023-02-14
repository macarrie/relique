import {Link, matchPath, useLocation, useMatch} from "react-router-dom";
import routes from "../routes";
let _ = require("lodash");

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

    let matchedRoute = routes.find((r :any) => {
        if (r.path === "*") {
            return false;
        }

        let match = matchPath(r.path, location.pathname);
        return match !== null;
    }) || {path: "/404"};

    let urlMatch = useMatch(matchedRoute.path);
    if (urlMatch === null) {
        return default_breadcrumbs;
    }

    let currentPath = matchedRoute.path;
    return _.cloneDeep(routes).filter((route :any) => currentPath.includes(route.path))
        .map((item :any) => {
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

    function renderCrumb(nav :any) {
        if (nav.path === "/") {
            return (
                <li className="inline-flex items-center" key={nav.path}>
                    <Link to={nav.path}
                          className="inline-flex items-center text-sm font-medium text-gray-700 hover:text-slate-900">
                        <div className="text-xl mr-1">
                            <i className="ri-home-2-line"></i>
                        </div>
                        {nav.name}
                    </Link>
                </li>
            );
        }

        return (
            <li aria-current="page" key={nav.path}>
                <div className="flex items-center">
                    <i className="text-2xl text-gray-400 ri-arrow-right-s-line"></i>
                    <Link to={nav.path}
                          className="ml-1 text-sm font-medium text-gray-400 hover:text-slate-700 md:ml-2">{nav.name}</Link>
                </div>
            </li>
        )
    }

    function renderCrumbsLine(crumbs :any) {
        return <>
            {crumbs.map((item :any) => renderCrumb(item))}
        </>
    }

    return (
        <nav className="flex-grow flex" aria-label="breadcrumb">
            <ol className="inline-flex items-center space-x-1">
                {renderCrumbsLine(breadcrumbs)}
            </ol>
        </nav>
    );
}

export default Breadcrumb;
