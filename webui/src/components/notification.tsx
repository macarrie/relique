import { ToastContentProps } from "react-toastify/unstyled";

type CustomNotificationProps = ToastContentProps<{
    title: string;
    content: string;
}>;

function Notification({
    data,
}: CustomNotificationProps) {
    return (
        <div className="">
            <h3 className="font-bold">
                {data.title}
            </h3>
            <div className="flex items-center justify-between">
                <p className="text-sm">{data.content}</p>
            </div>
        </div>
    );
}

export default Notification;