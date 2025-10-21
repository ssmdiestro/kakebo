import { useEffect, useMemo, useState } from "react";
import { fetchWeekReport } from "../api/week";

const DAY_ORDER = [0,1,2,3,4,5,6];
const DAY_NAME = {
  0: "Lunes",
  1: "Martes",
  2: "Miércoles",
  3: "Jueves",
  4: "Viernes",
  5: "Sábado",
  6: "Domingo",
};

const COLORS = {
  Supervivencia: "#b9cff7", // azul claro
  "Ocio y Vicio": "#d7ead7", // verde claro
  Compras: "#e8c1d2", // rosa claro
  Null: "#AAAAAA" //gris
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

  if (loading) return <div style={wrap}>Cargando…</div>;
  if (err) return <div style={wrap}><span style={{color:"crimson"}}>Error:</span> {err}</div>;
  if (!data) return null;

  return (
    <div style={boardWrap}>
      <div style={gridDays}>
        {days.map(([k, ds]) => (
          <DayColumn key={k} weekdayIndex={k} day={ds} />
        ))}
      </div>
    </div>
  );
}

function DayColumn({ weekdayIndex, day }) {
  // Extrae la fecha "YYYY-MM-DD" de tu DTO
  const dateStr = day?.date?.realDate;
  const jsDate = dateStr ? new Date(dateStr) : null;
  const dayNum = jsDate ? jsDate.getDate() : "";
  const notValid = day?.date?.year == 2999;
  // Aplana records por categoría (sumando subcategorías)
  const recs = {
    Supervivencia: flattenRecords(day?.supervivencia),
    "Ocio y Vicio": flattenRecords(day?.ocioyvicio),
    Compras: flattenRecords(day?.compras),
  };

  if (notValid){
    return (
        <div style={col}>
          {/* header */}
          <div style={headerTop}>
            <div style={{fontWeight:700}}>{DAY_NAME[weekdayIndex] || ""}</div>
          </div>
          <div style={notValid? headerDateNull: headerDate}>{dayNum}</div>
    
          {/* zonas por categoría */}
          {(["Supervivencia", "Ocio y Vicio", "Compras"]).map((cat) => (
            <CategoryBox key={cat} title={cat} color={COLORS["Null"]} records={recs[cat]} />
          ))}
    
          {/* total del día */}
          <div style={totalRow}>
            <div style={{width:"50%", textAlign:"left", paddingLeft:6}}>
              {fmt(day?.total)}
            </div>
            <div style={{width:"50%", textAlign:"right", paddingRight:6}}>TOTAL</div>
          </div>
        </div>
      );
  } 
  return (
    <div style={col}>
      {/* header */}
      <div style={headerTop}>
        <div style={{fontWeight:700}}>{DAY_NAME[weekdayIndex] || ""}</div>
      </div>
      <div style={headerDate}>{dayNum}</div>

      {/* zonas por categoría */}
      {(["Supervivencia", "Ocio y Vicio", "Compras"]).map((cat) => (
        <CategoryBox key={cat} title={cat} color={COLORS[cat]} records={recs[cat]} />
      ))}

      {/* total del día */}
      <div style={totalRow}>
        <div style={{width:"50%", textAlign:"left", paddingLeft:6}}>
          {fmt(day?.total)}
        </div>
        <div style={{width:"50%", textAlign:"right", paddingRight:6}}>TOTAL</div>
      </div>
    </div>
  );
}

function CategoryBox({ title, color, records }) {
    // Asegura que haya al menos 10 filas en la categoría
  const rows = records.length < 10 
  ? [...records, ...Array(10 - records.length).fill(null)] 
  : records;
  
    return (
      <div style={{ ...zone, background: color }}>
        {/* filas */}
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
  // categorySummary: { subcategory: [{ records: RecordDTO[] }] }
  if (!categorySummary?.subcategory?.length) return [];
  const all = [];
  for (const sc of categorySummary.subcategory) {
    if (Array.isArray(sc.records)) {
      for (const r of sc.records) {
        all.push(r);
      }
    }
  }
  // podrías ordenar por amount/descripcion si quieres
  return all;
}

// --- estilos ---
const boardWrap = { padding: 12, background:"#0b0b0b", color:"#121212", borderRadius: 12, scale:"80%" };
const gridDays = {
    display: "grid",
    gridAutoFlow: "column",
    gridAutoColumns: "minmax(200px, 1fr)", // Cambio: aumento el tamaño mínimo de las columnas
    gap: 8,
    overflowX: "auto",
  };
const col = {
    background: "#e5e7eb",
    border: "1px solid #9ca3af",
    borderRadius: 6,
    display: "flex",
    flexDirection: "column",
    minWidth: 220, // Cambio: Aumento el tamaño mínimo de cada columna (zona de cada día)
  };
const headerTop = { background:"#4b5563", color:"#fff", textAlign:"center", padding:"6px 4px", borderTopLeftRadius:6, borderTopRightRadius:6 };
const headerDate = { background:"#9ca3af", color:"#111827", textAlign:"center", padding:"4px 0", fontWeight:600, borderBottom:"1px solid #6b7280" };
const headerDateNull = { background:"#9ca3af", color:"#AA0000", textAlign:"center", padding:"4px 0", fontWeight:600, borderBottom:"1px solid #6b7280" };
const zone = { borderBottom:"2px solid #111", padding: 0 }; // cada bloque de color
const row = { display:"flex", borderTop:"1px solid rgba(0,0,0,.15)", height: 24, alignItems:"center" };
const cellLeft = { width:"40%", padding:"0 6px", borderRight:"1px solid rgba(0,0,0,.25)", fontWeight:600, color:"#0f172a", textAlign:"right"};
const cellRight = { width:"60%", padding:"0 6px", color:"#0f172a", overflow:"hidden", textOverflow:"ellipsis", whiteSpace:"nowrap" };

const totalRow = { display:"flex", background:"#6b7280", color:"#fff", borderTop:"2px solid #111", borderBottomLeftRadius:6, borderBottomRightRadius:6, height: 28, alignItems:"center" };

const wrap = { padding: 16, color:"#fff" };

const fmt = (n) => `${(Number(n)||0).toFixed(2).replace(".", ",")} €`;
