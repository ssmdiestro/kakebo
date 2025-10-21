import { useEffect } from "react";

export default function Modal({ open, onClose, children }) {
  useEffect(() => {
    const onKey = (e) => e.key === "Escape" && onClose?.();
    if (open) document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div
      onClick={onClose}
      style={{
        position: "fixed", inset: 0, background: "rgba(0,0,0,0.4)",
        display: "grid", placeItems: "center", zIndex: 1000, backdropFilter: "blur(2px)"
      }}
    >
      <div
        onClick={(e) => e.stopPropagation()}
        style={{
          width: "min(720px, 92vw)",
          borderRadius: 16,
          background: "linear-gradient(180deg, #fff, #f7f7fb)",
          boxShadow: "0 12px 40px rgba(0,0,0,.18)",
          padding: 20
        }}
      >
        {children}
      </div>
    </div>
  );
}
