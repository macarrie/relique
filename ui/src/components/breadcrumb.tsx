import React from "react";
import {Link, matchPath, PathMatch, useLocation, useMatch, useParams} from "react-router-dom";
import routes from "../routes";

function Breadcrumb() {
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

    function getCrumbs(currentRoute :any, urlMatch :PathMatch<string> | null) {
        if (urlMatch === null || currentRoute === undefined) {
            return default_breadcrumbs;
        }

        return routes.filter(({path}) => currentRoute.path.includes(path)).map(item => {
            let toReplace = item.path;
            Object.entries(urlMatch.params).forEach(([param, value], index) => {
                if (value !== undefined) {
                    toReplace = toReplace.replace(":" +param, value.toString())
                }
            });
            item.path = toReplace;
            return item;
        });
    }

    let matchedRoute = routes.find(({path}) => {
        if (path === "*") {
            return false;
        }

        let match = matchPath(path, location.pathname);
        if (match !== null) {
            return true;
        }
    });

    let path = matchedRoute === undefined ? "/" : matchedRoute.path;
    let breadcrumbs = getCrumbs(matchedRoute, useMatch(path));

    function renderCrumb(nav :any) {
        if (nav.path === "/") {
            return (
                <li className="inline-flex items-center" key={nav.path}>
                    <Link to={nav.path} className="inline-flex items-center text-sm font-medium text-gray-700 hover:text-blue-700">
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
                    <Link to={nav.path} className="ml-1 text-sm font-medium text-gray-400 md:ml-2">{nav.name}</Link>
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
