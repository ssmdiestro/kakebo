import { useState } from "react";
import Modal from "../components/Modal";
import NewRecordForm from "../components/NewRecordForm";

export default function Home() {
  const [open, setOpen] = useState(false);

  return (
    <div style={{ display: "grid", placeItems: "center", minHeight: "70vh" }}>
      <div style={{ textAlign: "center" }}>
        <h1 style={{ marginBottom: 16 }}>Kakebo</h1>
        <button
          onClick={() => setOpen(true)}
          style={{
            padding: "12px 18px", borderRadius: 12, border: "none",
            background: "linear-gradient(135deg, #10b981, #06b6d4)",
            color: "#fff", cursor: "pointer", boxShadow: "0 8px 20px rgba(16,185,129,.3)"
          }}
        >
          ➕ Nuevo registro
        </button>
      </div>

      <Modal open={open} onClose={() => setOpen(false)}>
        <NewRecordForm onSuccess={() => setOpen(false)} />
      </Modal>
    </div>
  );
}
