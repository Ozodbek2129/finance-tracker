import { useEffect, useState, useRef } from "react";
import { transactionsAPI } from "../../api/client";
import TransactionForm from "./TransactionForm";

// ── RAQAMLARNI JOYLASHGANDA SILLIQ O'STIRUVCHI ICHKI KOMPONENT ──
function AnimatedNumber({ value }) {
  const [current, setCurrent] = useState(0);

  useEffect(() => {
    let start = 0;
    const end = Number(value || 0);
    if (end === 0) { setCurrent(0); return; }

    const duration = 1000; // 1 sekundlik tez o'sish
    const startTime = performance.now();

    function updateNumber(currentTime) {
      const elapsedTime = currentTime - startTime;
      if (elapsedTime >= duration) {
        setCurrent(end);
        return;
      }
      const progress = elapsedTime / duration;
      const easeOutProgress = 1 - Math.pow(1 - progress, 3);
      const nextValue = Math.floor(easeOutProgress * end);
      setCurrent(nextValue);
      requestAnimationFrame(updateNumber);
    }
    requestAnimationFrame(updateNumber);
  }, [value]);

  return <span>{fmt(current)}</span>;
}

export default function TransactionList() {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [editing, setEditing] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const limit = 10;

  function load() {
    setLoading(true);
    transactionsAPI
      .getAll({ page, limit })
      .then((res) => setItems(res.data || res || []))
      .finally(() => setLoading(false));
  }

  useEffect(() => { load(); }, [page]);

  async function handleDelete(id) {
    if (!confirm("O'chirishni tasdiqlaysizmi?")) return;
    await transactionsAPI.delete(id);
    load();
  }

  if (editing) {
    return (
      <TransactionForm
        initial={editing}
        onSaved={() => { setEditing(null); load(); }}
        onCancel={() => setEditing(null)}
      />
    );
  }

  if (showForm) {
    return (
      <TransactionForm
        onSaved={() => { setShowForm(false); load(); }}
        onCancel={() => setShowForm(false)}
      />
    );
  }

  return (
    <div className="tl-wrap">
      {/* DINAMIK EFFEKTLAR NAVBATI BILAN PAYDO BO'LISHI UCHUN CSS */}
      <style>{`
        @keyframes listFadeIn {
          from {
            opacity: 0;
            transform: translateY(12px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
        .tl-item {
          animation: listFadeIn 0.4s ease-out forwards;
          opacity: 0; /* Animatsiya boshlanguncha yashirin turadi */
        }
      `}</style>

      <div className="tl-header">
        <h2>Tranzaksiyalar</h2>
        <button className="btn-add" onClick={() => setShowForm(true)}>+ Qo'shish</button>
      </div>

      {loading ? (
        <div className="tl-loading">Yuklanmoqda...</div>
      ) : items.length === 0 ? (
        <div className="tl-empty">Tranzaksiya yo'q. Yangi qo'shing!</div>
      ) : (
        <div className="tl-list">
          {items.map((t, index) => (
            <div 
              key={t.id} 
              className={`tl-item ${t.type}`}
              // Har bir element ketma-ket chiqishi uchun kichik kechikish (delay) qo'shamiz
              style={{ animationDelay: `${index * 0.05}s` }}
            >
              <span className="tl-icon">{t.category?.icon || "💳"}</span>
              <div className="tl-info">
                <span className="tl-cat">{t.category?.name || "—"}</span>
                <span className="tl-desc">{t.description || ""}</span>
              </div>
              <div className="tl-right">
                <span className="tl-amount" style={{ whiteSpace: "nowrap" }}>
                  {t.type === "expense" ? "−" : "+"}
                  {/* Raqamni o'suvchi animatsiya bilan chiqaramiz */}
                  <AnimatedNumber value={t.amount} /> so'm
                </span>
                <span className="tl-date">{t.date?.slice(0, 10)}</span>
              </div>
              <div className="tl-btns">
                <button className="btn-edit" onClick={() => setEditing(t)}>✏️</button>
                <button className="btn-del" onClick={() => handleDelete(t.id)}>🗑</button>
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="tl-pagination">
        <button disabled={page === 1} onClick={() => setPage((p) => p - 1)}>‹ Oldingi</button>
        <span>Sahifa {page}</span>
        <button disabled={items.length < limit} onClick={() => setPage((p) => p + 1)}>Keyingi ›</button>
      </div>
    </div>
  );
}

function fmt(n) {
  return Number(n || 0).toLocaleString("uz-UZ");
}