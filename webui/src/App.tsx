import {
  BrowserRouter,
  Routes,
  Route,
  Navigate
} from "react-router";

import appRoutes from "./routes"

import './App.css'
import Sidebar from './layout/sidebar'
import Main from './layout/main'

function App() {
  return (
    <BrowserRouter basename="/">
      <Sidebar />

      <Main>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace={true} />} />
          {appRoutes.map(({ path, elt }, key) => (
            <Route path={path} element={elt()} key={key} />
          ))}
        </Routes>
      </Main>
    </BrowserRouter >
  );
}

export default App