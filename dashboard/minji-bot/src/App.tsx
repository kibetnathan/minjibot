import { Routes, Route } from "react-router-dom"
import Home from "./pages/Home"
import Commands from "./pages/Commands"
import Login from "./pages/Login"
import SignUp from "./pages/SignUp"
import Dashboard from "./pages/Dashboard"

export function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/commands" element={<Commands />} />
      <Route path="/login" element={<Login />} />
      <Route path="/signup" element={<SignUp />} />
      <Route path="/dashboard" element={<Dashboard />} />
    </Routes>
  )
}

export default App
