import { useEffect, useMemo, useState } from "react";
import { fetchWeekReport } from "../api/endpoints";

const DAY_ORDER = [0, 1, 2, 3, 4, 5, 6];
const DAY_NAME = {
  0: "Lunes",
  1: "Martes",
  2: "Miércoles",
  3: "Jueves",
  4: "Viernes",
  5: "Sábado",
  6: "Domingo",
  9: "Monedero",
};

const COLORS = {
  Supervivencia: "#b9cff7", // azul claro
  "Ocio y Vicio": "#d7ead7", // verde claro
  Compras: "#e8c1d2", // rosa claro
  Null: "#AAAAAA" // gris
};

export default function WeekBoard({ week, month, year, refreshKey }) {
  const [data, setData] = useState(null);
  const [err, setErr] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let alive = true;
    (async () => {
      try {
        setLoading(true);
        setErr(null);
        const json = await fetchWeekReport({ week, month, year }); // el helper ya parsea JSON
        if (alive) setData(json);
      } catch (e) {
        if (alive) setErr(e.message);
      } finally {
        if (alive) setLoading(false);
      }
    })();
    return () => { alive = false; };
  }, [week, month, year, refreshKey]);

  const days = useMemo(() => {
    if (!data?.daySummary) return [];
    // daySummary viene como map[int]DaySummary
    // Lo convertimos a array ordenado por la clave 1..7
    return DAY_ORDER
      .map((k) => [k, data.daySummary[String(k)] || data.daySummary[k]])
      .filter(([, ds]) => !!ds);
  }, [data]);

  // 🧮 Calcular el máximo de filas por categoría (mínimo 10)
  const rowHeights = useMemo(() => {
    const base = { Supervivencia: 10, "Ocio y Vicio": 10, Compras: 10 };
    for (const [, ds] of days) {
      base.Supervivencia = Math.max(base.Supervivencia, flattenRecords(ds?.supervivencia).length);
      base["Ocio y Vicio"] = Math.max(base["Ocio y Vicio"], flattenRecords(ds?.ocioyvicio).length);
      base.Compras = Math.max(base.Compras, flattenRecords(ds?.compras).length);
    }
    return base;
  }, [days]);

  if (loading) return <div style={wrap}>Cargando…</div>;
  if (err) return <div style={wrap}><span style={{ color: "crimson" }}>Error:</span> {err}</div>;
  if (!data) return null;

  return (
    <div style={boardWrap}>
      <div style={gridDays}>
        {days.map(([k, ds]) => (
          <DayColumn key={k} weekdayIndex={k} day={ds} rowHeights={rowHeights} />
        ))}
        <DayColumn key={9} weekdayIndex={9} day={data?.daySummary[1]} rowHeights={rowHeights} />

      </div>

    </div>
  );
}

function DayColumn({ weekdayIndex, day, rowHeights }) {
  const dateStr = day?.date?.realDate;
  const jsDate = dateStr ? new Date(dateStr) : null;
  const dayNum = jsDate ? jsDate.getDate() : "";
  const notValid = day?.date?.year == 2999 && weekdayIndex !== 9;


  const recs = {
    Supervivencia: flattenRecords(day?.supervivencia),
    "Ocio y Vicio": flattenRecords(day?.ocioyvicio),
    Compras: flattenRecords(day?.compras),
  };

  const catList = ["Supervivencia", "Ocio y Vicio", "Compras"];

  const Category = ({ cat }) => (
    <CategoryBox
      title={cat}
      color={notValid ? COLORS["Null"] : COLORS[cat]}
      records={recs[cat]}
      maxRows={rowHeights[cat]}   // 👈 sincroniza altura por categoría
    />
  );

  if (notValid) {
    return (
      <div style={col}>
        <div style={headerTop}><div style={{ fontWeight: 700 }}>{DAY_NAME[weekdayIndex] || ""}</div></div>
        <div style={headerDateNull}>{weekdayIndex === 9 ? "€€" : dayNum}</div>
        {catList.map((c) => <Category key={c} cat={c} />)}
        <div style={totalRow}>
          <div style={{ width: "50%", textAlign: "left", paddingLeft: 6 }}>{fmt(day?.total)}</div>
          <div style={{ width: "50%", textAlign: "right", paddingRight: 6 }}>TOTAL</div>
        </div>
      </div>
    );
  }

  return (
    <div style={col}>
      <div style={headerTop}><div style={{ fontWeight: 700 }}>{DAY_NAME[weekdayIndex] || ""}</div></div>
      <div style={headerDate}>{weekdayIndex === 9 ? "€€" : dayNum}</div>
      {catList.map((c) => <Category key={c} cat={c} />)}
      <div style={totalRow}>
        <div style={{ width: "50%", textAlign: "left", paddingLeft: 6 }}>{fmt(day?.total)}</div>
        <div style={{ width: "50%", textAlign: "right", paddingRight: 6 }}>TOTAL</div>
      </div>
    </div>
  );
}

function CategoryBox({ title, color, records, maxRows }) {
  const count = Math.max(10, Number(maxRows) || 10);
  const rows = records.length < count
    ? [...records, ...Array(count - records.length).fill(null)]
    : records.slice(0, count);

  return (
    <div style={{ ...zone, background: color }}>
      {rows.map((r, idx) => (
        <div key={idx} style={row}>
          <div style={cellLeft}>{r ? fmt(r.amount) : ""}</div>
          <div style={cellRight}>{r ? (r.description || "") : ""}</div>
        </div>
      ))}
    </div>
  );
}

// --- helpers DTO ---
function flattenRecords(categorySummary) {
  if (!categorySummary?.subcategory?.length) return [];
  const all = [];
  for (const sc of categorySummary.subcategory) {
    if (Array.isArray(sc.records)) {
      for (const r of sc.records) {
        all.push(r);
      }
    }
  }
  return all;
}

// --- estilos ---
const boardWrap = { padding: 12, background: "#0b0b0b", color: "#121212", borderRadius: 12, scale: "100%" };
const gridDays = {
  display: "grid",
  gridAutoFlow: "column",
  gridAutoColumns: "minmax(200px, 1fr)",
  gap: 8,
  overflowX: "auto",
};
const col = {
  background: "#e5e7eb",
  border: "1px solid #9ca3af",
  borderRadius: 6,
  display: "flex",
  flexDirection: "column",
  minWidth: 220,
};
const headerTop = { background: "#4b5563", color: "#fff", textAlign: "center", padding: "6px 4px", borderTopLeftRadius: 6, borderTopRightRadius: 6 };
const headerDate = { background: "#9ca3af", color: "#111827", textAlign: "center", padding: "4px 0", fontWeight: 600, borderBottom: "1px solid #6b7280" };
const headerDateNull = { background: "#9ca3af", color: "#AA0000", textAlign: "center", padding: "4px 0", fontWeight: 600, borderBottom: "1px solid #6b7280" };
const zone = { borderBottom: "2px solid #111", padding: 0 };
const row = { display: "flex", borderTop: "1px solid rgba(0,0,0,.15)", height: 24, alignItems: "center" };
const cellLeft = { width: "40%", padding: "0 6px", borderRight: "1px solid rgba(0,0,0,.25)", fontWeight: 600, color: "#0f172a", textAlign: "right" };
const cellRight = { width: "60%", padding: "0 6px", color: "#0f172a", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" };
const totalRow = { display: "flex", background: "#6b7280", color: "#fff", borderTop: "2px solid #111", borderBottomLeftRadius: 6, borderBottomRightRadius: 6, height: 28, alignItems: "center" };
const wrap = { padding: 16, color: "#fff" };

const fmt = (n) => `${(Number(n) || 0).toFixed(2).replace(".", ",")} €`;
