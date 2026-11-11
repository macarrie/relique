import { Link } from "react-router-dom";
import PageHeader from "../components/page_header";

function NotFound() {
    return (
        <>
            <PageHeader title="Page not found" />
            <p className="text-muted">
                This page does not exist. Go back to the <Link to="/dashboard" className="text-link">dashboard</Link>.
            </p>
        </>
    );
}

export default NotFound;
