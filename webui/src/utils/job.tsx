import Const from "../types/const";

export default class JobUtils {
    static jobStateToCode = function (s: string): number {
        switch (s) {
            case "pending":
                return Const.UNKNOWN;
            case "active":
                return Const.INFO;
            case "success":
                return Const.OK;
            case "incomplete":
                return Const.WARNING;
            case "error":
                return Const.CRITICAL;
            default:
                return Const.UNKNOWN;
        }
    };

    // Severity of a job event comes from its type alone. Unknown types are neutral
    static eventTypeToCode = function (t: string): number {
        switch (t) {
            case "client_ping_error":
            case "sync_error":
            case "image_error":
            case "job_error":
                return Const.CRITICAL;
            case "fallback_to_full":
                return Const.WARNING;
            case "client_ping":
            case "sync_end":
            case "image_created":
                return Const.OK;
            default:
                return Const.UNKNOWN;
        }
    };

    // Events after which a job is finished: the job itself can be reloaded to get its final state.
    // A failed client ping ends the job without a job_end event
    static isTerminalEvent = function (t: string): boolean {
        return t === "job_end" || t === "job_error" || t === "client_ping_error";
    };

    // "sync_error" -> "Sync error"
    static eventTypeLabel = function (t: string): string {
        if (!t) {
            return "Unknown";
        }
        const words = t.replace(/_/g, " ");
        return words.charAt(0).toUpperCase() + words.slice(1);
    };
}