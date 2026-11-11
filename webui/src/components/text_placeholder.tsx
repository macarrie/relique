function TextPlaceholder({ className = "w-40" }) {
    return (
        <span className={`skeleton inline-block h-3 align-middle ${className}`}></span>
    )
}

export default TextPlaceholder;
