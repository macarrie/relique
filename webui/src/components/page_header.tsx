import React from "react";
import Breadcrumb from "./breadcrumb";

function PageHeader({
    title,
    meta,
    actions,
    loading = false,
}: {
    title: React.ReactNode,
    meta?: React.ReactNode,
    actions?: React.ReactNode,
    loading?: boolean,
}) {
    return (
        <header className="space-y-4 pb-6">
            <Breadcrumb />
            <div className="flex flex-wrap items-end gap-x-6 gap-y-3">
                <div className="flex-grow min-w-0">
                    <h1 className="font-display text-3xl font-bold tracking-tight break-words">
                        {loading ? <span className="skeleton inline-block h-8 w-56 align-middle"></span> : title}
                    </h1>
                    {meta && (
                        <div className="mt-1 flex flex-wrap items-center gap-x-4 gap-y-1 text-muted">
                            {meta}
                        </div>
                    )}
                </div>
                {actions && (
                    <div className="flex items-center gap-2">
                        {actions}
                    </div>
                )}
            </div>
        </header>
    );
}

export default PageHeader;
