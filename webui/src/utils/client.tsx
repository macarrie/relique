import Client from "../types/client";
import Const from "../types/const";

export default class ClientUtils {
    static sshUser = function (c: Client): string {
        return c?.ssh_user ? c.ssh_user : Const.DEFAULT_CLIENT_SSH_USER;
    };

    static sshPort = function (c: Client): number {
        return c?.ssh_port ? c.ssh_port : Const.DEFAULT_CLIENT_SSH_PORT;
    };

    static sshTarget = function (c: Client): string {
        return `${ClientUtils.sshUser(c)}@${c?.address}:${ClientUtils.sshPort(c)}`;
    };

    static sshStatusLabel = function (status: number): string {
        switch (status) {
            case Const.OK:
                return "Reachable";
            case Const.CRITICAL:
                return "Unreachable";
            default:
                return "Not checked";
        }
    };
}
