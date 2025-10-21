import { useEffect, useMemo, useState } from "react";
import ColorSelect from "./ColorSelect";


export default function NewRecordForm({ onSuccess }) {
  const [categories, setCategories] = useState([]);
  const [loadingCats, setLoadingCats] = useState(true);
  const [catsError, setCatsError] = useState(null);

  const [sending, setSending] = useState(false);
  const [msg, setMsg] = useState(null);

  const [form, setForm] = useState({
    description: "",
    date: new Date().toISOString().slice(0, 10),
    amount: "",
    notes: "",
    category: "",
    subcategory: "",
  });

  useEffect(() => {
    (async () => {
      try {
        setLoadingCats(true);
        const url = "/api/getCategories"
        const res = await fetch(url);
        if (!res.ok) throw new Error(`GET /getCategories → ${res.status}`);
        const data = await res.json();
        setCategories(Array.isArray(data) ? data : []);
        if (data?.length) {
          setForm(f => ({
            ...f,
            category: data[0].categoryName,
            subcategory: data[0].subCategories?.[0] || ""
          }));
        }
      } catch (e) {
        setCatsError(e.message);
      } finally {
        setLoadingCats(false);
      }
    })();
  }, []);

  const selectedCategory = useMemo(
    () => categories.find(c => c.categoryName === form.category) || null,
    [categories, form.category]
  );
  const subcategories = selectedCategory?.subCategories || [];

  const onChange = (e) => {
    const { name, value } = e.target;
    setForm(f => {
      if (name === "category") {
        const cat = categories.find(c => c.categoryName === value);
        return { ...f, category: value, subcategory: cat?.subCategories?.[0] || "" };
      }
      return { ...f, [name]: value };
    });
  };

  const onSubmit = async (e) => {
    e.preventDefault();
    setMsg(null);
    if (!form.description.trim()) return setMsg("Pon una descripción.");
    if (!/^\d{4}-\d{2}-\d{2}$/.test(form.date)) return setMsg("Fecha en formato YYYY-MM-DD.");
    if (form.amount === "" || Number.isNaN(Number(form.amount))) return setMsg("Cantidad inválida.");
    if (!form.category) return setMsg("Elige una categoría.");
    if (!form.subcategory) return setMsg("Elige una subcategoría.");

    const payload = {
      description: form.description.trim(),
      date: form.date,
      subcategory: { description: form.subcategory, category: form.category },
      amount: Number(form.amount),
      notes: form.notes || "",
    };

    try {
      setSending(true);
      const url = "/api/newRecord"
      const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.error || `Error ${res.status} creando registro`);
      }
      setMsg("Registro creado correctamente ✅");
      // cierra el modal al crear
      setTimeout(() => onSuccess?.(), 400);
    } catch (err) {
      setMsg(err.message);
    } finally {
      setSending(false);
    }
  };

  return (
    <form onSubmit={onSubmit} style={{ display: "grid", gap: 14 }}>
      <h3 style={{ margin: 0, fontSize: 22 }}>Nuevo registro</h3>
      <p style={{ margin: "2px 0 10px", color: "#666" }}>
        Rellena los campos y pulsa <b>Crear registro</b>.
      </p>

      {catsError && <div style={{ color: "crimson" }}>Error categorías: {catsError}</div>}

      <div style={{ display: "grid", gap: 8 }}>
        <label style={{ ...lbl, display: "flex", alignItems: "center", gap: 8 }}>
          Categoría
        </label>
        <ColorSelect
          options={categories}
          value={form.category}
          onChange={(newCat) => {
            const cat = categories.find(c => c.categoryName === newCat);
            setForm(f => ({
              ...f,
              category: newCat,
              subcategory: cat?.subCategories?.[0] || ""
            }));
          }}
          getLabel={(o) => o.categoryName}
          getColor={(o) => o.categoryColor}
        />
      </div>

      <div style={{ display: "grid", gap: 8 }}>
        <label style={lbl}>Subcategoría</label>
        <select 
        style={{
          ...inp,
          WebkitAppearance: "none",
          MozAppearance: "none",
          appearance: "none",
          backgroundImage:
            'url("data:image/svg+xml;utf8,<svg fill=\\"%23999\\" height=\\"12\\" viewBox=\\"0 0 20 20\\" width=\\"12\\" xmlns=\\"http://www.w3.org/2000/svg\\"><path d=\\"M7 8l3 3 3-3z\\"/></svg>")',
          backgroundRepeat: "no-repeat",
          backgroundPosition: "right 12px center",
          backgroundSize: "12px",
        }}
        name="subcategory" value={form.subcategory} onChange={onChange} disabled={!subcategories.length}>
          {subcategories.map(s => <option key={s} value={s}>{s}</option>)}
        </select>
      </div>

      <div style={{ display: "grid", gap: 8 }}>
        <label style={lbl}>Descripción</label>
        <input style={inp} name="description" value={form.description} onChange={onChange} required />
      </div>

      <div style={{ display: "grid", gap: 8 }}>
        <label style={lbl}>Fecha</label>
        <input style={inp} type="date" name="date" value={form.date} onChange={onChange} required />
      </div>

      <div style={{ display: "grid", gap: 8 }}>
        <label style={lbl}>Cantidad</label>
        <input style={inp} type="number" min="0" step="0.50" inputMode="decimal" name="amount" value={form.amount} onChange={onChange} placeholder="0" required />
      </div>

      <div style={{ display: "grid", gap: 8 }}>
        <label style={lbl}>Notas</label>
        <textarea style={{ ...inp, minHeight: 90, resize: "vertical" }} name="notes" value={form.notes} onChange={onChange} placeholder="Notas opcionales…" />
      </div>

      {msg && <div style={{ color: msg.includes("✅") ? "green" : "crimson" }}>{msg}</div>}

      <div style={{ display: "flex", gap: 10, justifyContent: "flex-end", marginTop: 6 }}>
        <button type="button" onClick={() => onSuccess?.()} style={btnSecondary}>Cancelar</button>
        <button type="submit" disabled={sending || loadingCats} style={btnPrimary}>
          {sending ? "Enviando…" : "Crear registro"}
        </button>
      </div>
    </form>
  );
}

const lbl = { fontSize: 14, color: "#444" };
const inp = {
  color: "#000",
  width: "100%", padding: "10px 12px", borderRadius: 10,
  border: "1px solid #ddd", background: "#fff", outline: "none",
  boxShadow: "0 1px 2px rgba(0,0,0,.03)",
  minHeight: 44,      
  fontSize: 16,
  lineHeight: "22px",
  boxSizing: "border-box",
};
const btnPrimary = {
  padding: "10px 14px", borderRadius: 10, border: "none",
  background: "linear-gradient(135deg, #6366f1, #8b5cf6)",
  color: "#fff", cursor: "pointer", boxShadow: "0 6px 16px rgba(99,102,241,.32)"
};
const btnSecondary = {
  padding: "10px 14px", borderRadius: 10, border: "1px solid #ddd",
  background: "#fff", color: "#333", cursor: "pointer"
};
