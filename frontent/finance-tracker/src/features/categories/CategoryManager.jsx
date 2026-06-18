import { useEffect, useState, useRef } from "react";
import { categoriesAPI } from "../../api/client";

const EMOJIS = [
  "🍕","🍔","🍜","🍣","🥗","☕","🛒","🏠","🚗","⛽","✈️","🎬","🎮","🎵",
  "👗","👟","💄","📚","🏋️","⚽","🏥","💊","🐶","🌿","💡","📱","💻","🖥️",
  "🔧","🎁","💰","💳","🏦","📦","🧾","🛍️","🍺","🎂","🏖️","🚌","🚇","🚕",
  "🏫","👶","🧸","💒","⛪","🌍","🧹","🪴","🕯️","🎓","🏢","🩺","💉","🦷",
];

function EmojiPicker({ value, onChange }) {
  const [open, setOpen] = useState(false);
  const ref = useRef(null);

  useEffect(() => {
    function handleClick(e) {
      if (ref.current && !ref.current.contains(e.target)) setOpen(false);
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  return (
    <div className="emoji-picker-wrap" ref={ref}>
      <button type="button" className="emoji-trigger" onClick={() => setOpen((o) => !o)}>
        <span className="emoji-current">{value}</span>
        <span className="emoji-arrow">{open ? "▲" : "▼"}</span>
      </button>
      {open && (
        <div className="emoji-grid show">
          {EMOJIS.map((e) => (
            <button
              key={e}
              type="button"
              className={`emoji-btn ${value === e ? "selected" : ""}`}
              onClick={() => { onChange(e); setOpen(false); }}
            >
              {e}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

const EMPTY = { name: "", icon: "📦", color: "#6366f1" };

export default function CategoryManager() {
  const [items, setItems] = useState([]);
  const [form, setForm] = useState(EMPTY);
  const [editId, setEditId] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  function load() {
    categoriesAPI.getAll().then((res) => setItems(res.data || res || []));
  }

  useEffect(() => { load(); }, []);

  function set(field, val) {
    setForm((p) => ({ ...p, [field]: val }));
  }

  function startEdit(cat) {
    setEditId(cat.id);
    setForm({ name: cat.name, icon: cat.icon, color: cat.color || "#6366f1" });
  }

  function cancelEdit() {
    setEditId(null);
    setForm(EMPTY);
    setError(null);
  }

  async function handleSave() {
    if (!form.name.trim()) { setError("Nom kiritilmagan."); return; }
    setLoading(true);
    setError(null);
    try {
      if (editId) {
        await categoriesAPI.update(editId, form);
      } else {
        await categoriesAPI.create(form);
      }
      cancelEdit();
      load();
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(id) {
    if (!confirm("O'chirishni tasdiqlaysizmi?")) return;
    await categoriesAPI.delete(id);
    load();
  }

  return (
    <div className="cm-wrap">
      {/* ── KATEGORIYA ELEMENTLARI UCHUN SILJISH ANIMATSIYALARI ── */}
      <style>{`
        @keyframes catFadeIn {
          from {
            opacity: 0;
            transform: translateY(10px) scale(0.98);
          }
          to {
            opacity: 1;
            transform: translateY(0) scale(1);
          }
        }
        .cm-item {
          animation: catFadeIn 0.35s cubic-bezier(0.25, 1, 0.5, 1) forwards;
          opacity: 0;
          transition: border-left-color 0.2s ease, transform 0.2s ease;
        }
        .cm-item:hover {
          transform: translateX(4px);
        }
        .emoji-grid.show {
          animation: catFadeIn 0.2s ease-out forwards;
        }
        .cm-form, .btn-save, .btn-cancel {
          transition: all 0.2s ease;
        }
      `}</style>

      <h2>Kategoriyalar</h2>

      <div className="cm-form">
        <h3>{editId ? "Kategoriyani tahrirlash" : "Yangi kategoriya"}</h3>
        {error && <div className="cm-error">{error}</div>}
        <div className="cm-row">
          <label>Emoji
            <EmojiPicker value={form.icon} onChange={(e) => set("icon", e)} />
          </label>
          <label>Nom
            <input value={form.name} onChange={(e) => set("name", e.target.value)} placeholder="Masalan: Oziq-ovqat" />
          </label>
          <label>Rang
            <input type="color" value={form.color} onChange={(e) => set("color", e.target.value)} />
          </label>
        </div>
        <div className="cm-form-btns">
          <button className="btn-save" onClick={handleSave} disabled={loading}>
            {loading ? "..." : editId ? "Yangilash" : "Qo'shish"}
          </button>
          {editId && <button className="btn-cancel" onClick={cancelEdit}>Bekor</button>}
        </div>
      </div>

      <div className="cm-list">
        {items.length === 0 ? (
          <div className="cm-empty">Kategoriya yo'q. Birinchi bo'lib qo'shing!</div>
        ) : (
          items.map((c, index) => (
            <div 
              key={c.id} 
              className="cm-item" 
              style={{ 
                borderLeftColor: c.color,
                // Har bir kator ketma-ketlikda navbati bilan chiroyli ko'rinishi uchun
                animationDelay: `${index * 0.04}s` 
              }}
            >
              <span className="cm-icon">{c.icon}</span>
              <span className="cm-name">{c.name}</span>
              <div className="cm-item-btns">
                <button className="btn-edit" onClick={() => startEdit(c)}>✏️</button>
                <button className="btn-del" onClick={() => handleDelete(c.id)}>🗑</button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}