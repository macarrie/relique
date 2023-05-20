import React, {useState} from "react";

function Dropdown(props :any) {
    let [isOpen, setOpen] = useState(false);

    return (
        <div className="relative inline-block text-left">
            <button type="button" className="button button-text button-small" onClick={() => setOpen(!isOpen)}><i className="text-base ri-more-2-line"></i></button>
            {isOpen && (
                <div className="absolute right-0 z-10 space-y-1 py-1 mt-2 rounded min-w-max bg-slate-100 dark:bg-slate-700 shadow">
                    {React.Children.toArray(props.children).map((elt :any, index :number) => (
                        <div key={index} className="mx-1 rounded py-2 px-3 hover:bg-slate-200 dark:hover:bg-slate-600 cursor-pointer" onClick={() => setOpen(false)}>{elt}</div>
                    ))}
                </div>
            )}
        </div>
    );
}

export default Dropdown;
