import { Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import MonthPage from "./pages/MonthPage";

export default function App() {
  return (
    <Routes>
      {/* Redirige según tu mes contable (usa fetchWeekLimitsByDate) */}
      <Route path="/" element={<Home />} />

      {/* Vista de un mes concreto, ejemplo: /10/2025 */}
      <Route path="/:month/:year" element={<MonthPage />} />

      {/* (opcional) fallback */}
      <Route path="*" element={<Home />} />
    </Routes>
  );
}
