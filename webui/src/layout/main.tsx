import React from "react";

function Main({ children }: { children?: React.ReactNode }) {
    return (
        <main className="sm:ml-56 px-6 py-6 lg:px-10">
            <div className="max-w-[100rem] mx-auto flex flex-col gap-4">
                {children}
            </div>
        </main>
    )
}

export default Main
