import React from "react";

// variant="tile" (default): heading + content on a tinted, borderless tile.
// variant="plain": same, directly on the page background (e.g. nested inside a tile).
function Section({
    title,
    description,
    actions,
    variant = "tile",
    className = "",
    children,
}: {
    title?: React.ReactNode,
    description?: React.ReactNode,
    actions?: React.ReactNode,
    variant?: "plain" | "tile",
    className?: string,
    children?: React.ReactNode,
}) {
    const isTile = variant === "tile";

    return (
        <section className={`space-y-4 ${isTile ? "bg-base-200 rounded-box p-5 lg:p-6" : ""} ${className}`}>
            {(title || actions) && (
                <div className="flex items-baseline gap-4">
                    <div className={`flex-grow min-w-0 ${isTile ? "flex flex-wrap items-baseline gap-x-2" : ""}`}>
                        {title && <h2 className="font-display text-xl font-semibold tracking-tight">{title}</h2>}
                        {description && (
                            <p className="text-muted">
                                {isTile && <span aria-hidden="true">· </span>}
                                {description}
                            </p>
                        )}
                    </div>
                    {actions && (
                        <div className="flex items-center gap-2">
                            {actions}
                        </div>
                    )}
                </div>
            )}
            {children}
        </section>
    );
}

export default Section;
