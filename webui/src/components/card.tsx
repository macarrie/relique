function Card(props: any) {
    return (
        <div className={`border border-gray-200 rounded bg-base-100 ${props.className}`}>
            {props.children}
        </div>
    );
}

export default Card;