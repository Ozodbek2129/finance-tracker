import { useEffect, useState } from "react";
import { transactionsAPI } from "../../api/client";
import TransactionForm from "./TransactionForm";

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
          {items.map((t) => (
            <div key={t.id} className={`tl-item ${t.type}`}>
              <span className="tl-icon">{t.category?.icon || "💳"}</span>
              <div className="tl-info">
                <span className="tl-cat">{t.category?.name || "—"}</span>
                <span className="tl-desc">{t.description || ""}</span>
              </div>
              <div className="tl-right">
                <span className="tl-amount">
                  {t.type === "expense" ? "−" : "+"}{fmt(t.amount)} so'm
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