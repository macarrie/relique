function PaginationButton({
    children,
    onClick,
    disabled,
    active,
}: any) {
    return (
        <button className={`w-8 h-8 rounded-full hover:bg-slate-200 disabled:text-slate-300 disabled:hover:bg-transparent ${active ? "!bg-slate-500/20" : ""}`}
            onClick={onClick}
            disabled={disabled}>
            {children}
        </button>
    );
}
export default PaginationButton;