function DashboardStat({
    label = "",
    value = 0 as any,
    color = "text-base-content"
}) {
    return (

        <div className="stat place-items-center">
            <div className={`stat-value ${value ? color : 'text-base-content/50'}`}>
                {value}
            </div>
            <div className={`stat-title ${value ? color : 'text-base-content/50'}`}>
                {label}
            </div>
        </div>
    );
}

export default DashboardStat;