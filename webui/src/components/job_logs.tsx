import { KeyboardEvent, useState } from "react";
import JobLog from "../types/job_log";
import EmptyValue from "./empty_value";
import TextPlaceholder from "./text_placeholder";

// rsync --progress rewrites the current line with carriage returns: only keep the last state of each line
function cleanOutput(output: string): string {
    return output
        .split("\n")
        .map((line) => {
            const parts = line.split("\r").filter((p) => p !== "");
            return parts.length > 0 ? parts[parts.length - 1] : "";
        })
        .join("\n")
        .trimEnd();
}

// Scroll to the end on display: rsync stats and errors are written last
function scrollToBottom(el: HTMLDivElement | null) {
    if (el) {
        el.scrollTop = el.scrollHeight;
    }
}

function lineNumbers(content: string): string {
    const count = content.split("\n").length;
    return Array.from({ length: count }, (_, i) => i + 1).join("\n");
}

function LogOutput({ label, output, truncated }: { label: string, output: string, truncated: boolean }) {
    const content = cleanOutput(output ?? "");

    return (
        <div className="space-y-1">
            <div className="flex flex-wrap items-baseline gap-x-2">
                <span className="text-label">{label}</span>
                {truncated && <span className="text-sm text-muted">Showing the end of the log only</span>}
            </div>
            {content === "" ? (
                <div><EmptyValue label="Empty" /></div>
            ) : (
                // Numbers and content share the scroll container so they stay aligned. The inner row takes the full
                // content size, so the gutter (border and background) spans the whole log and not only the visible part.
                // The gutter sticks to the left when scrolling long lines horizontally, and is excluded from text
                // selection and screen readers
                <div ref={scrollToBottom} className="code bg-base-100 rounded-box max-h-96 overflow-auto">
                    <div className="flex w-max min-w-full">
                        <pre aria-hidden="true" className="sticky left-0 shrink-0 select-none text-right tabular-nums text-muted bg-base-100 border-r border-base-300 py-3 pl-3 pr-2 whitespace-pre">
                            {lineNumbers(content)}
                        </pre>
                        <pre className="grow py-3 pl-3 pr-3 whitespace-pre">{content}</pre>
                    </div>
                </div>
            )}
        </div>
    );
}

function tabId(idx: number): string {
    return `job-log-tab-${idx}`;
}

function panelId(idx: number): string {
    return `job-log-panel-${idx}`;
}

// rsync output of each synced path of a job, one tab per path
function JobLogs({
    logs,
    loading = false,
}: {
    logs: JobLog[],
    loading?: boolean,
}) {
    // Selection is kept by name so it survives a reload. Falls back to the first tab if the name disappears
    const [selected, setSelected] = useState<string>("");
    const activeIdx = Math.max(0, logs.findIndex((l) => l.name === selected));
    const active = logs[activeIdx];

    // Left/right arrows move between tabs, as expected from a tablist
    function onTabKeyDown(e: KeyboardEvent<HTMLButtonElement>) {
        if (e.key !== "ArrowRight" && e.key !== "ArrowLeft") {
            return;
        }
        e.preventDefault();
        const next = (activeIdx + (e.key === "ArrowRight" ? 1 : -1) + logs.length) % logs.length;
        setSelected(logs[next].name);
        document.getElementById(tabId(next))?.focus();
    }

    if (loading) {
        return (
            <div className="space-y-3">
                <div><TextPlaceholder className="w-48" /></div>
                <div><TextPlaceholder className="w-full" /></div>
                <div><TextPlaceholder className="w-full" /></div>
            </div>
        );
    }

    if (logs.length === 0) {
        return <EmptyValue label="No logs available. Logs are written when each sync ends" />;
    }

    return (
        <div className="space-y-4">
            <div role="tablist" aria-label="Backup paths" className="tabs tabs-border">
                {logs.map((l, idx) => {
                    const isActive = idx === activeIdx;
                    const hasErrors = cleanOutput(l.stderr ?? "") !== "";

                    return (
                        <button key={l.name} id={tabId(idx)} type="button" role="tab"
                            aria-selected={isActive}
                            aria-controls={panelId(idx)}
                            tabIndex={isActive ? 0 : -1}
                            onClick={() => setSelected(l.name)}
                            onKeyDown={onTabKeyDown}
                            className={`tab code gap-1.5 ${isActive ? "tab-active" : ""}`}>
                            {l.backup_path || l.name}
                            {hasErrors && (
                                <i className="ri-error-warning-line text-error" title="Errors logged" aria-label="Errors logged"></i>
                            )}
                        </button>
                    );
                })}
            </div>

            <div key={active.name} id={panelId(activeIdx)} role="tabpanel" aria-labelledby={tabId(activeIdx)} className="space-y-3">
                <LogOutput label="Output" output={active.stdout} truncated={active.stdout_truncated} />
                <LogOutput label="Errors" output={active.stderr} truncated={active.stderr_truncated} />
            </div>
        </div>
    );
}

export default JobLogs;
