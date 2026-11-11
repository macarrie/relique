import { ToastContentProps } from "react-toastify/unstyled";

type CustomNotificationProps = ToastContentProps<{
    title: string;
    content: string;
}>;

function Notification({
    data,
}: CustomNotificationProps) {
    return (
        <div className="font-sans text-sm text-base-content">
            <p className="font-semibold">{data.title}</p>
            {data.content && <p className="text-muted break-words">{data.content}</p>}
        </div>
    );
}

export default Notification;
