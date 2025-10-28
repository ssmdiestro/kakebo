import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { fetchWeekLimitsByDate } from "../api/endpoints";

function formatLocalYYYYMMDD(d = new Date()) {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

export default function Home() {
  const navigate = useNavigate();
  const [error, setError] = useState(null);

  const todayStr = useMemo(() => formatLocalYYYYMMDD(), []);

  useEffect(() => {
    let cancel = false;
    (async () => {
      try {
        // Llama a tu endpoint que aplica la lógica del "mes contable"
        const json = await fetchWeekLimitsByDate(todayStr);
        if (cancel) return;

        // Acepta tanto camel como Pascal (como en tu código previo)
        const month = Number(json.month ?? json.Month);
        const year = Number(json.year ?? json.Year ?? new Date().getFullYear());

        if (!Number.isInteger(month) || !Number.isInteger(year)) {
          throw new Error(`Respuesta inválida: ${JSON.stringify(json)}`);
        }

        // Redirige a la ruta normalizada /:month/:year
        navigate(`/${month}/${year}`, { replace: true });
      } catch (e) {
        setError(String(e?.message ?? e));
      }
    })();
    return () => { cancel = true; };
  }, [todayStr, navigate]);

  if (error) {
    // Fallback simple (puedes poner un botón "Reintentar" si quieres)
    return <div style={{ padding: 16 }}>⚠️ Error inicializando: {error}</div>;
  }

  return <div style={{ padding: 16 }}>Cargando…</div>;
}
