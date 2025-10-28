import { useMemo, useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import WeekBoard from "./WeekBoard";
import Modal from "./Modal";
import NewRecordForm from "./NewRecordForm";
import MonthTotals from "./MonthTotals";

const MONTH_NAMES_ES = [
    "",
    "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
    "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
];

function addMonths(month, year, delta) {
    let m = month + delta; // puede salirse
    let y = year + Math.floor((m - 1) / 12);
    m = ((m - 1) % 12 + 12) % 12 + 1;
    return { month: m, year: y };
}

export default function MonthView({ month, year }) {
    const [open, setOpen] = useState(false);
    const [refreshKey, setRefreshKey] = useState(0);
    const navigate = useNavigate();

    const monthName = useMemo(() => MONTH_NAMES_ES[Number(month)] ?? "", [month]);
    const bumpRefresh = useCallback(() => setRefreshKey(k => k + 1), []);

    const goPrev = useCallback(() => {
        const { month: m, year: y } = addMonths(month, year, -1);
        navigate(`/${m}/${y}`);
    }, [month, year, navigate]);

    const goNext = useCallback(() => {
        const { month: m, year: y } = addMonths(month, year, +1);
        navigate(`/${m}/${y}`);
    }, [month, year, navigate]);

    const headerWrap = {
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        gap: 12,
        margin: "10px 0",
    };

    const navBtn = {
        width: 36,
        height: 36,
        borderRadius: "999px",
        border: "none",
        cursor: "pointer",
        fontSize: 18,
        lineHeight: "36px",
        boxShadow: "0 4px 12px rgba(0,0,0,.15)",
        background: "linear-gradient(135deg, #e5e7eb, #f3f4f6)",
    };

    return (
        <>
            {/* Botón flotante ➕ */}
            <button
                onClick={() => setOpen(true)}
                style={{
                    padding: "12px 18px", borderRadius: 12, border: "none",
                    background: "linear-gradient(135deg, #10b981, #06b6d4)",
                    color: "#fff", cursor: "pointer", boxShadow: "0 8px 20px rgba(16,185,129,.3)",
                    position: "fixed", bottom: 10, right: 10, zIndex: 1000
                }}
                aria-label="Añadir registro"
                title="Añadir registro"
            >
                ➕
            </button>

            {/* Header con flechas */}
            <div style={headerWrap}>
                <button aria-label="Mes anterior" title="Mes anterior" style={navBtn} onClick={goPrev}>
                    ‹
                </button>
                <h1 style={{ margin: 0, textAlign: "center" }}>
                    {monthName} - {year}
                </h1>
                <button aria-label="Mes siguiente" title="Mes siguiente" style={navBtn} onClick={goNext}>
                    ›
                </button>
            </div>

            {/* Semanas (1..5) */}
            <WeekBoard week={1} month={month} year={year} refreshKey={refreshKey} />
            <WeekBoard week={2} month={month} year={year} refreshKey={refreshKey} />
            <WeekBoard week={3} month={month} year={year} refreshKey={refreshKey} />
            <WeekBoard week={4} month={month} year={year} refreshKey={refreshKey} />
            <WeekBoard week={5} month={month} year={year} refreshKey={refreshKey} />

            <MonthTotals month={month} year={year} refreshKey={refreshKey} />

            <Modal open={open} onClose={() => setOpen(false)}>
                <NewRecordForm
                    onSuccess={() => {
                        setOpen(false);
                        bumpRefresh(); // fuerza re-fetch
                    }}
                    onCancel={() => setOpen(false)}
                />
            </Modal>
        </>
    );
}
