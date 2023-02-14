import React from 'react';
import {
    BrowserRouter,
    Routes,
    Route,
    Navigate
} from "react-router-dom";

import appRoutes from "./routes"

import Sidebar from './layout/sidebar';
import Main from './layout/main';

function App() {
    return (
        <BrowserRouter basename="/ui">
            <div className="flex flex-row min-h-screen bg-blue-50 text-slate-700">
                <Sidebar/>

                <Main>
                    <Routes>
                        <Route path="/" element={<Navigate to="/dashboard" replace={true}/>}/>
                        {appRoutes.map(({path, elt}, key) => (
                            <Route path={path} element={elt()} key={key}/>
                        ))}
                    </Routes>
                </Main>
            </div>
        </BrowserRouter>
    );
}

export default App;
