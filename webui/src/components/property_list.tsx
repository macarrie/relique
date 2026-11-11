import React from "react";
import TextPlaceholder from "./text_placeholder";
import EmptyValue from "./empty_value";

export type Property = {
    label: string,
    value?: React.ReactNode,
    mono?: boolean,
};

function isEmpty(value: React.ReactNode): boolean {
    return value === undefined || value === null || value === "" || (Array.isArray(value) && value.length === 0);
}

function PropertyList({
    items,
    loading = false,
}: {
    items: Property[],
    loading?: boolean,
}) {
    return (
        <dl className="grid grid-cols-[minmax(7rem,11rem)_1fr] gap-x-6 gap-y-2">
            {items.map((item: Property) => (
                <React.Fragment key={item.label}>
                    <dt className="text-muted">{item.label}</dt>
                    <dd className={`min-w-0 break-words ${item.mono ? "code" : ""}`}>
                        {loading ? (
                            <TextPlaceholder />
                        ) : isEmpty(item.value) ? (
                            <EmptyValue />
                        ) : (
                            item.value
                        )}
                    </dd>
                </React.Fragment>
            ))}
        </dl>
    );
}

export default PropertyList;
