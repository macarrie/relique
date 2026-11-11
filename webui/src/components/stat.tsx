import { Link } from "react-router-dom";

const TONES: Record<string, string> = {
    success: "text-success",
    warning: "text-warning",
    error: "text-error",
    info: "text-info",
};

function Stat({
    label,
    value,
    tone,
    to,
    loading = false,
}: {
    label: string,
    value: number | string | undefined,
    tone?: "success" | "warning" | "error" | "info",
    to?: string,
    loading?: boolean,
}) {
    const isZero = value === 0 || value === undefined;
    const color = isZero ? "text-muted" : (tone ? TONES[tone] : "text-base-content");

    const content = (
        <>
            <div className={`font-display text-5xl font-semibold leading-none tabular-nums ${color}`}>
                {loading ? <span className="skeleton inline-block h-8 w-12 align-middle"></span> : (value ?? 0)}
            </div>
            <div className={`mt-2 text-muted ${to ? "group-hover:text-base-content group-hover:underline underline-offset-2" : ""}`}>
                {label}
            </div>
        </>
    );

    if (to) {
        return <Link to={to} className="group block">{content}</Link>;
    }

    return <div>{content}</div>;
}

export default Stat;
