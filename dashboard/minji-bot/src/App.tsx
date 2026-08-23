import { Routes, Route } from "react-router-dom"
import Home from "./pages/Home"

export function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
    </Routes>
  )
}

export default App
