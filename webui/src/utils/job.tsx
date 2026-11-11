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
}