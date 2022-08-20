import React, {useState} from "react";

function Tab(props :any) {
    return <div>{props.title}: {props.children}</div>
}

function Tabs(props :any) {
    const [activeTab, setActiveTab] = useState(props.children[0].props.title)

    function renderTabTitle(t: string) {
        if (activeTab === t) {
            return (
                <span className="cursor-pointer inline-block px-4 py-3 border-b-2 border-blue-500 text-blue-500">{t}</span>
            )
        }

        return (
            <span className="cursor-pointer inline-block px-4 py-3 border-b-2 border-transparent">{t}</span>
        )
    }

    function renderTabLine() {
        return (
            <ul className="flex flex-wrap border-b">
                {props.children.map((tab :any) => {
                    return (
                        <li className="mr-2" key={tab.props.title} onClick={(e) => {e.preventDefault(); setActiveTab(tab.props.title)}}>
                            {renderTabTitle(tab.props.title)}
                        </li>
                    )
                })}
            </ul>
        )
    }

    function renderTabContent() {
        return (
            <div className="bg-slate-50 p-4">
                {props.children.map((tab :any) => {
                    if (tab.props.title !== activeTab) {
                        return undefined;
                    }
                    return tab.props.children
                })}
            </div>
        )
    }

    return <>
        <div className="flex flex-row items-center px-4">
            {props.title && (
                <div className={"flex-grow uppercase font-bold text-slate-500 mr-3"}>{props.title}</div>
            )}
            {renderTabLine()}
        </div>
        {renderTabContent()}
    </>
}

export {Tabs, Tab};