import { useMemo } from "react";
import { useParams, useNavigate } from "react-router-dom";
import MonthView from "../components/MonthView";
import NotFound from "../components/NotFound";

function normalizeMonthYear(m, y) {
    let month = Number(m);
    let year = Number(y);

    if (!Number.isInteger(month) || !Number.isInteger(year)) {
        return { valid: false };
    }

    if (month < 1 || month > 12) return { valid: false };
    if (year < 1900 || year > 3000) return { valid: false };

    return { valid: true, month, year };
}

export default function MonthPage() {
    const { month: mParam, year: yParam } = useParams();
    const navigate = useNavigate();

    const { valid, month, year } = useMemo(
        () => normalizeMonthYear(mParam, yParam),
        [mParam, yParam]
    );

    if (!valid) {
        return <NotFound />;
    }

    return (
        <div style={{ padding: 10 }}>
            <MonthView month={month} year={year} />
        </div>
    );
}
