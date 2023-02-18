import React, {useState} from "react";

import Sidebar from "../layout/sidebar";
import Logo from "./logo";

function BurgerMenu() {
    const [open, setOpen] = useState(false);

    if (open) {
        return (
            <div className="bg-blue-50 dark:bg-slate-900 flex flex-col w-full h-full fixed top-0 left-0">
                <div className="flex flex-row px-4 py-4">
                    <button
                        className="shrink text-2xl mr-4 text-slate-600 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100"
                        onClick={() => setOpen(false)}>
                        <i className="ri-menu-2-line"></i>
                    </button>
                    <div className="grow">
                        <Logo/>
                    </div>
                    <div className="shrink">
                        <button
                            className="text-2xl text-slate-600 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100"
                            onClick={() => setOpen(false)}>
                            <i className="ri-close-line"></i>
                        </button>
                    </div>
                </div>
                <div className="grow">
                    {open && (
                        <Sidebar mobile/>
                    )}
                </div>
            </div>
        );
    }

    return (
        <button
            className="md:hidden text-2xl mr-4 text-slate-600 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100"
            onClick={() => setOpen(true)}>
            <i className="ri-menu-2-line"></i>
        </button>
    );
}

export default BurgerMenu;
