import { toast } from "react-toastify";
import Notification from "../components/notification";
import Const from "../types/const";

export default class Utils {
    // Returns "—" for missing values and Go zero times (0001-01-01T00:00:00Z)
    static formatDate = function (d: string): string {
        if (!d) {
            return "—";
        }
        const obj = new Date(d);
        if (isNaN(obj.getTime()) || obj.getFullYear() <= 1) {
            return "—";
        }
        return obj.toLocaleString();
    }

    static formatDuration = function (start: string, end: string): string {
        const s = new Date(start), e = new Date(end);
        if (!start || !end || isNaN(s.getTime()) || isNaN(e.getTime()) || s.getFullYear() <= 1 || e.getFullYear() <= 1) {
            return "";
        }
        let secs = Math.max(0, Math.round((e.getTime() - s.getTime()) / 1000));
        const h = Math.floor(secs / 3600);
        const m = Math.floor((secs % 3600) / 60);
        secs = secs % 60;
        if (h > 0) {
            return `${h}h ${m}m`;
        }
        if (m > 0) {
            return `${m}m ${secs}s`;
        }
        return `${secs}s`;
    }

    static shortID = function (id: string): string {
        if (!id) {
            return "unknown";
        }
        return id.split("-")[0];
    }
    static formatSize = function (bytes: number, si = false, dp = 1) {
        const thresh = si ? 1000 : 1024;

        if (Math.abs(bytes) < thresh) {
            return bytes + ' B';
        }

        const units = si
            ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
            : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
        let u = -1;
        const r = 10 ** dp;

        do {
            bytes /= thresh;
            ++u;
        } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


        return bytes.toFixed(dp) + ' ' + units[u];
    }

    // options.autoClose overrides the display time (ms) of the notification.
    // options.onClose is called once the notification is closed, after its display time or when dismissed by the user
    static notify = function (status: number, title: string, content: string, options: { autoClose?: number, onClose?: () => void } = {}) {
        let func
        switch (status) {
            case Const.OK: {
                func = toast.success
                break;
            }
            case Const.WARNING: {
                func = toast.warn
                break;
            }
            case Const.CRITICAL: {
                func = toast.error
                break;
            }
            default: {
                func = toast.info
                break;
            }
        }

        func(Notification, {
            data: {
                title: title,
                content: content,
            },
            ...(options.autoClose !== undefined && { autoClose: options.autoClose }),
            onClose: options.onClose,
        })
    }
}