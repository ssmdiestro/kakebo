import { useEffect, useState, useMemo } from "react";
import { fetchMonthReport } from "../api/endpoints";

export default function MonthTotals({ month, year, refreshKey }) {
    const [data, setData] = useState(null);
    const [err, setErr] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        let alive = true;
        (async () => {
            try {
                setLoading(true);
                setErr(null);
                const json = await fetchMonthReport(month, year); // el helper ya parsea JSON
                if (alive) setData(json);
            } catch (e) {
                if (alive) setErr(e.message);
            } finally {
                if (alive) setLoading(false);
            }
        })();
        return () => { alive = false; };
    }, [month, year, refreshKey]);

    const {
        sumasSemanas,
        sumaSupervivenciaSemanas,
        sumaOcioYVicioSemanas,
        sumaComprasSemanas,
        totalMes,
        totalSupervivenciaMes,
        totalOcioYVicioMes,
        totalComprasMes,
    } = useMemo(() => {
        const weeks = data?.weekSums ?? {};
        const res = {
            sumasSemanas: {},
            sumaSupervivenciaSemanas: {},
            sumaOcioYVicioSemanas: {},
            sumaComprasSemanas: {},
            totalMes: 0,
            totalSupervivenciaMes: 0,
            totalOcioYVicioMes: 0,
            totalComprasMes: 0,
        };

        for (const [key, week] of Object.entries(weeks)) {
            const tMes = week?.total ?? 0;
            const tSup = week?.supervivenciaSum?.total ?? 0;
            const tOcio = week?.ocioyvicioSum?.total ?? 0;
            const tCompr = week?.comprasSum?.total ?? 0;

            res.sumasSemanas[key] = tMes;
            res.sumaSupervivenciaSemanas[key] = tSup;
            res.sumaOcioYVicioSemanas[key] = tOcio;
            res.sumaComprasSemanas[key] = tCompr;
        }
        res.totalMes = data?.total ?? 0;
        res.totalSupervivenciaMes = data?.supervivenciaTotal?.total ?? 0;
        res.totalOcioYVicioMes = data?.ocioyvicioTotal?.total ?? 0;
        res.totalComprasMes = data?.comprasTotal?.total ?? 0;
        console.log(res);
        return res;
    }, [data]);

    return <div>totalSupervivencia:</div>;
}