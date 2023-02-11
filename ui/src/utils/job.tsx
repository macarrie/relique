import Const from "../types/const";

export default class JobUtils {
    static statusColor = function (s: string): string {
        let color: string
        switch (s) {
            case "pending":
                color = "text-neutral-700";
                break;
            case "active":
                color = "text-sky-700";
                break;
            case "success":
                color = "text-teal-700";
                break;
            case "incomplete":
                color = "text-orange-700";
                break;
            case "error":
                color = "text-red-700";
                break;
            default:
                color = "text-slate-700";
        }

        return color
    };

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
}
