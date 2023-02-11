import Const from "../types/const";
import Client from "../types/client";

export default class ClientUtils {
    static alive = function (c :Client) :number {
        if (c.api_alive === Const.OK && c.ssh_alive === Const.OK) {
            return Const.OK;
        } else if (c.api_alive === Const.OK || c.ssh_alive === Const.OK) {
            return Const.WARNING;
        }

        return Math.max(c.api_alive, c.ssh_alive);
    };

    static APIAliveLabel = function (c :Client) :string {
        if (c.api_alive === Const.OK) {
            return "Client API endpoints reachable"
        } else if (c.api_alive === Const.UNKNOWN) {
            return "Client API status unknown"
        }

        return "Cannot reach client API endpoints"
    };

    static SSHAliveLabel = function (c :Client) :string {
        if (c.ssh_alive === Const.OK) {
            return "Reachable by SSH"
        } else if (c.ssh_alive === Const.UNKNOWN) {
            return "SSH connectivity status unknown"
        }

        return "Cannot reach client via SSH"
    };

    static GlobalAliveLabel = function (c :Client) :string {
        if (c.api_alive === Const.OK && c.ssh_alive === Const.OK) {
            return "Client connectivity OK";
        }
        return "Connectivity issues detected"
    };
}