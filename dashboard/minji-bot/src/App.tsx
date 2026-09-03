import { Routes, Route } from "react-router-dom"
import Home from "./pages/Home"
import Commands from "./pages/Commands"

export function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/commands" element={<Commands />} />
    </Routes>
  )
}

export default App
