import React from 'react';
import {
  BrowserRouter,
  Routes,
  Route,
} from "react-router-dom";

import Sidebar from './layout/sidebar';
import Main from './layout/main';
import Dashboard from './components/dashboard';
import Jobs from './components/jobs';
import Clients from './components/clients';
import NotFound from './components/not_found';

function App() {
    return (
        <BrowserRouter>
            <div className="flex flex-row">
                <Sidebar />

                <Main>
                    <Routes>
                        <Route path="/" element={<Dashboard />} />
                        <Route path="/dashboard" element={<Dashboard />} />
                        <Route path="/jobs" element={<Jobs />} />
                        <Route path="/clients" element={<Clients />} />
                        <Route path="*" element={<NotFound />} />
                    </Routes>
                </Main>
            </div>
        </BrowserRouter>
    );
}

export default App;
