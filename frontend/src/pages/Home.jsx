import { useEffect, useMemo, useState, useCallback } from "react";
import WeekBoard from "../components/WeekBoard";
import { fetchWeekLimitsByDate } from "../api/endpoints";
import Modal from "../components/Modal";
import NewRecordForm from "../components/NewRecordForm";
import MonthTotals from "../components/MonthTotals";

function formatLocalYYYYMMDD(d = new Date()) {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

export default function Home() {
  const [open, setOpen] = useState(false);
  const [limits, setLimits] = useState(null);
  const [error, setError] = useState(null);
  const [refreshKey, setRefreshKey] = useState(0); // 👈 clave de refresco

  const todayStr = useMemo(() => formatLocalYYYYMMDD(), []);
  const currentYear = useMemo(() => new Date().getFullYear(), []);

  // 👉 función para forzar re-fetch en los hijos
  const bumpRefresh = useCallback(() => {
    setRefreshKey(k => k + 1);
  }, []);

  useEffect(() => {
    let cancel = false;
    (async () => {
      try {
        const json = await fetchWeekLimitsByDate(todayStr);
        if (cancel) return;

        const week = json.week ?? json.Week;
        const month = json.month ?? json.Month;
        const startDate = json.startDate ?? json.StartDate;
        const endDate = json.endDate ?? json.EndDate;
        const monthName = json.monthName ?? json.MonthName;

        if (week == null || month == null) {
          throw new Error(`Respuesta sin week/month: ${JSON.stringify(json)}`);
        }

        setLimits({
          week: Number(week),
          month: Number(month),
          startDate,
          endDate,
          monthName,
        });
        setError(null);
      } catch (e) {
        setError(String(e?.message ?? e));
      }
    })();
    return () => { cancel = true; };
  }, [todayStr]);

  if (error) {
    return <div style={{ padding: 16 }}>⚠️ Error cargando semana: {error}</div>;
  }

  if (!limits) {
    return <div style={{ padding: 16 }}>Cargando…</div>;
  }

  return (
    <div style={{ padding: 10 }}>
      <button
        onClick={() => setOpen(true)}
        style={{
          padding: "12px 18px", borderRadius: 12, border: "none",
          background: "linear-gradient(135deg, #10b981, #06b6d4)",
          color: "#fff", cursor: "pointer", boxShadow: "0 8px 20px rgba(16,185,129,.3)",
          position: "fixed",
          bottom: 10,
          right: 10,
          zIndex: 1000
        }}
      >
        ➕
      </button>
      <h1 style={{ textAlign: "center", marginBottom: 10, marginTop: 10 }}> {limits.monthName} - {currentYear}</h1>
      {/* ✅ Pasa refreshKey a todos los WeekBoard para que se re-fetchen */}
      <WeekBoard week={1} month={limits.month} year={currentYear} refreshKey={refreshKey} />
      <WeekBoard week={2} month={limits.month} year={currentYear} refreshKey={refreshKey} />
      <WeekBoard week={3} month={limits.month} year={currentYear} refreshKey={refreshKey} />
      <WeekBoard week={4} month={limits.month} year={currentYear} refreshKey={refreshKey} />
      <WeekBoard week={5} month={limits.month} year={currentYear} refreshKey={refreshKey} />

      {<MonthTotals month={limits.month} year={currentYear} refreshKey={refreshKey} />}
      <Modal open={open} onClose={() => setOpen(false)}>
        <NewRecordForm
          onSuccess={() => {
            setOpen(false);  // cierra modal
            bumpRefresh();   // 👈 fuerza re-fetch de WeekBoard
          }}
          onCancel={() => setOpen(false)}
        />
      </Modal>
    </div>
  );
}
