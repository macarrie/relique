import Const from "../types/const";

// Token classes for each status code, shared by every component displaying a status
const STATUS_COLORS: Record<number, { text: string, dot: string }> = {
    [Const.OK]: { text: "text-success", dot: "bg-success" },
    [Const.WARNING]: { text: "text-warning", dot: "bg-warning" },
    [Const.CRITICAL]: { text: "text-error", dot: "bg-error" },
    [Const.INFO]: { text: "text-info", dot: "bg-info" },
    [Const.UNKNOWN]: { text: "text-muted", dot: "bg-base-300" },
};

export default STATUS_COLORS;
