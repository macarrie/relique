import JobEvent from "../types/job_event";
import Const from "../types/const";
import JobUtils from "../utils/job";
import Utils from "../utils/utils";
import STATUS_COLORS from "../utils/status_colors";
import EmptyValue from "./empty_value";
import TextPlaceholder from "./text_placeholder";

function formatTime(d: string): string {
    const obj = new Date(d);
    if (isNaN(obj.getTime())) {
        return "—";
    }
    return obj.toLocaleTimeString();
}

// Vertical timeline of the events recorded during a job, oldest first.
// Severity comes from the event type; it is always shown with the type label, never by color alone.
function JobEventTimeline({
    events,
    jobStart,
    loading = false,
}: {
    events: JobEvent[],
    jobStart: string,
    loading?: boolean,
}) {
    if (loading) {
        return (
            <div className="space-y-3">
                <div><TextPlaceholder className="w-64" /></div>
                <div><TextPlaceholder className="w-80" /></div>
                <div><TextPlaceholder className="w-56" /></div>
            </div>
        );
    }

    if (events.length === 0) {
        return <EmptyValue label="No events recorded" />;
    }

    return (
        <ol className="ml-1 border-l border-base-300 space-y-3">
            {events.map((e, idx) => {
                const code = JobUtils.eventTypeToCode(e.type);
                const colors = STATUS_COLORS[code] ?? STATUS_COLORS[Const.UNKNOWN];
                // The neutral status dot has the color of the timeline line: use muted so it stays visible
                const dot = code === Const.UNKNOWN ? "bg-muted" : colors.dot;
                const offset = Utils.formatDuration(jobStart, e.time);

                return (
                    <li key={`${e.time}-${idx}`} className="relative pl-5">
                        <span className={`absolute left-[-5.5px] top-[0.45rem] size-2.5 rounded-full ${dot}`} aria-hidden="true"></span>
                        <div className="flex flex-wrap items-baseline gap-x-3">
                            <span className="tabular-nums text-muted">{formatTime(e.time)}</span>
                            {offset && <span className="tabular-nums text-sm text-muted">+{offset}</span>}
                            <span className={`font-medium ${code === Const.UNKNOWN ? "" : colors.text}`}>{JobUtils.eventTypeLabel(e.type)}</span>
                        </div>
                        {e.description && (
                            <p className={`wrap-break-word ${code === Const.CRITICAL ? "code" : ""}`}>{e.description}</p>
                        )}
                    </li>
                );
            })}
        </ol>
    );
}

export default JobEventTimeline;
