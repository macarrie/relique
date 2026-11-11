import Const from "../types/const";

const STATUS_COLORS: Record<number, { text: string, dot: string }> = {
    [Const.OK]: { text: "text-success", dot: "bg-success" },
    [Const.WARNING]: { text: "text-warning", dot: "bg-warning" },
    [Const.CRITICAL]: { text: "text-error", dot: "bg-error" },
    [Const.INFO]: { text: "text-info", dot: "bg-info" },
    [Const.UNKNOWN]: { text: "text-muted", dot: "bg-base-300" },
};

// Status is always conveyed by a text label; the dot only reinforces it.
// variant="dot" hides the label visually (kept for screen readers and as tooltip).
function Status({
    status = Const.UNKNOWN,
    label = "",
    variant = "label",
    loading = false,
}: {
    status?: number,
    label?: string,
    variant?: "label" | "dot",
    loading?: boolean,
}) {
    if (loading) {
        return <span className="loading loading-spinner loading-xs text-muted" aria-label="Loading"></span>;
    }

    const colors = STATUS_COLORS[status] ?? STATUS_COLORS[Const.UNKNOWN];

    if (variant === "dot") {
        return (
            <span className="inline-flex items-center" title={label}>
                <span className={`size-2.5 rounded-full ${colors.dot}`}></span>
                <span className="sr-only">{label}</span>
            </span>
        );
    }

    return (
        <span className={`inline-flex items-center gap-1.5 font-medium ${colors.text}`}>
            <span className={`size-2 rounded-full ${colors.dot}`}></span>
            {label}
        </span>
    );
}

export default Status;
