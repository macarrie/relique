import Card from "./card";
import Const from "../types/const";
import { Link } from "react-router-dom";
import Client from "../types/client";

function ClientCard(props: any) {
    let cl: Client = props.client;

    if (cl === null || cl === undefined) {
        return <div>Loading</div>
    }

    return (
        <Card className={props.className}>
            <div className="p-4 flex flex-row items-center mb-2">
                <div className="flex-grow font-bold leading-none">Client <span className='ml-2 badge badge-neutral badge-soft font-normal'>{cl?.name}</span></div>
                {props.link && (
                    <Link className="link-primary link-hover text-sm" to={`/clients/${cl?.name}`}>
                        Details
                    </Link>
                )}
            </div>
            <table className="table">
                <tbody>
                    <tr>
                        <td>Name</td>
                        <td className=''>{cl?.name}</td>
                    </tr>
                    <tr>
                        <td>Address</td>
                        <td className='code'>{cl?.address}</td>
                    </tr>
                    <tr>
                        <td>SSH user</td>
                        <td className='code'>{cl?.ssh_user == "" ? Const.DEFAULT_CLIENT_SSH_USER : cl?.ssh_user}</td>
                    </tr>
                    <tr>
                        <td>SSH port</td>
                        <td className='code'>{cl?.ssh_port == 0 ? Const.DEFAULT_CLIENT_SSH_PORT : cl?.ssh_port}</td>
                    </tr>
                </tbody>
            </table>
        </Card>
    );
}

export default ClientCard;