export async function fetchWeekReport({ week, month, year }) {
    const url = `/api/getWeekReport?week=${week}&month=${month}&year=${year}`;
    const res = await fetch(url);
    if (!res.ok) {
      const text = await res.text().catch(() => "");
      throw new Error(`HTTP ${res.status} ${text}`);
    }
    return res.json();
  }
  
  export async function fetchWeekLimitsByDate(dateStr) {
    const url = `/api/getWeek/${encodeURIComponent(dateStr)}`;
    const res = await fetch(url);
    if (!res.ok) {
        const text = await res.text().catch(() => "");
        throw new Error(`HTTP ${res.status} ${text}`);
      }
      const data = await res.json();   // ✅ leer una sola vez
      return data;
  }