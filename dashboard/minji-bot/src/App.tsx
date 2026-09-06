import { Routes, Route } from "react-router-dom"
import Home from "./pages/Home"
import Commands from "./pages/Commands"
import Login from "./pages/Login"
import SignUp from "./pages/SignUp"
import Dashboard from "./pages/Dashboard"
import GuildLogs from "./pages/GuildLogs"
import GuildSettings from "./pages/GuildSettings"
import Profile from "./pages/Profile"
import Diary from "./pages/Diary"

export function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/commands" element={<Commands />} />
      <Route path="/login" element={<Login />} />
      <Route path="/signup" element={<SignUp />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/dashboard/guild/:guildId" element={<GuildLogs />} />
      <Route path="/dashboard/guild/:guildId/settings" element={<GuildSettings />} />
      <Route path="/dashboard/profile" element={<Profile />} />
      <Route path="/dashboard/diary" element={<Diary />} />
    </Routes>
  )
}

export default App
