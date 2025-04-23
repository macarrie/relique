import Breadcrumb from "../components/breadcrumb"

function Main(props: any) {
    return (
        <>
            <div className="p-4 sm:ml-64">
                <div className="mb-4">
                    <Breadcrumb />
                </div>

                <div className="container space-y-4">
                    {props.children}
                </div>
            </div>
        </>
    )
}

export default Main