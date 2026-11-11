function YesNo({ value = false }) {
    return value ? <span>Yes</span> : <span className="text-muted">No</span>;
}

export default YesNo;
