import React, {useState} from "react";

import Card from './card';

function Tab(props: any) {
    let active: boolean = props.active;
    return <div className={`${props.className} ${!active && "hidden"}`}>{props.children}</div>
}

function Tabs(props: any) {
    const [activeTab, setActiveTab] = useState(props.initialActiveTab || 0);

    function renderTabLine() {
        let activeClass = "border-blue-500 text-blue-500 dark:border-blue-400 dark:text-blue-400"
        let inactiveClass = "border-transparent"
        return (
            <ul className="flex flex-wrap">
                {React.Children.map(props.children, (tab: any) => {
                    let active = activeTab === tab.key
                    return (
                        <li className={`cursor-pointer block flex flex-row items-center px-4 py-3 border-b-2 mr-2 ${active ? activeClass : inactiveClass} ${tab.props.headerClassName}`}
                            key={tab.key} onClick={(e) => {
                            e.preventDefault();
                            setActiveTab(tab.key)
                        }}>
                            {tab.props.title}
                        </li>
                    )
                })}
            </ul>
        )
    }

    function renderTabContent() {
        return (
            <Card className="p-4">
                {React.Children.map(props.children, (tab: any) => {
                    let active = (tab.key === activeTab);

                    return <Tab active={active} key={tab.key} {...tab.props}>{tab.props.children}</Tab>;
                })}
            </Card>
        )
    }

    return <>
        <div className="flex flex-row items-center px-4">
            {props.title && (
                <div className={"flex-grow font-bold text-slate-500 mr-3"}>{props.title}</div>
            )}
            {renderTabLine()}
        </div>
        {renderTabContent()}
    </>
}

export {Tabs, Tab};