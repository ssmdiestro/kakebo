// src/components/ColorSelect.jsx
import { useState, useRef, useEffect } from "react";

export default function ColorSelect({ options, value, onChange, getLabel, getColor, placeholder="Selecciona…" }) {
  const [open, setOpen] = useState(false);
  const ref = useRef(null);
  useEffect(() => {
    const onDoc = (e) => { if (ref.current && !ref.current.contains(e.target)) setOpen(false); };
    document.addEventListener("click", onDoc);
    return () => document.removeEventListener("click", onDoc);
  }, []);
  const selected = options.find(o => getLabel(o) === value);

  const base = {
    width: "100%", padding: "10px 12px", borderRadius: 10,
    color: "#000",
    border: "1px solid #ddd", background: "#fff", cursor: "pointer"
  };

  return (
    <div ref={ref} style={{ position: "relative" }}>
      <button type="button" onClick={() => setOpen(!open)} style={{ ...base, display: "flex", alignItems: "center", justifyContent: "space-between" }}>
        <span style={{ display: "flex", alignItems: "center", gap: 8, color: "#000" }}>
          <span style={{
            width: 12, height: 12, borderRadius: 999, border: "1px solid #999",
            background: selected ? getColor(selected) : "#ccc"
          }} />
          {selected ? getLabel(selected) : <span style={{ color: "#888" }}>{placeholder}</span>}
        </span>
        <span style={{ opacity: .6 }}>▾</span>
      </button>

      {open && (
        <div style={{
          position: "absolute", zIndex: 10, insetInline: 0, marginTop: 6, borderRadius: 10,
          border: "1px solid #e5e7eb", background: "#fff", boxShadow: "0 12px 30px rgba(0,0,0,.12)"
        }}>
          {options.map(o => (
            <div key={getLabel(o)}
              onClick={() => { onChange(getLabel(o)); setOpen(false); }}
              style={{
                padding: "10px 12px", display: "flex", alignItems: "center", gap: 8,
                background: value === getLabel(o) ? "#f5f7ff" : "transparent", cursor: "pointer",
                color: "#000"
              }}>
              <span style={{
                width: 12, height: 12, borderRadius: 999, border: "1px solid #999",
                background: getColor(o)
              }} />
              {getLabel(o)}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
