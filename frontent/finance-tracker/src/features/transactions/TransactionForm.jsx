import { useEffect, useState } from "react";
import { transactionsAPI, categoriesAPI } from "../../api/client";

const EMPTY = { amount: "", description: "", categoryId: "", type: "expense", date: today() };

export default function TransactionForm({ initial = null, onSaved, onCancel }) {
  const [form, setForm] = useState(initial ? toForm(initial) : EMPTY);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    categoriesAPI.getAll().then((res) => setCategories(res.data || res || []));
  }, []);

  function set(field, value) {
    setForm((prev) => ({ ...prev, [field]: value }));
  }

  async function handleSubmit() {
    if (!form.amount || !form.categoryId || !form.date) {
      setError("Summa, kategoriya va sanani to'ldiring.");
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const payload = {
        amount: parseInt(form.amount),
        description: form.description,
        category_id: form.categoryId,
        type: form.type,
        date: form.date,
      };
      if (initial) {
        await transactionsAPI.update(initial.id, payload);
      } else {
        await transactionsAPI.create(payload);
      }
      onSaved?.();
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="tf-form">
      <h2>{initial ? "Tranzaksiyani tahrirlash" : "Yangi tranzaksiya"}</h2>

      {error && <div className="tf-error">{error}</div>}

      <label>Tur
        <div className="tf-toggle">
          {["expense", "income"].map((t) => (
            <button
              key={t}
              type="button"
              className={`tf-toggle-btn ${form.type === t ? "active" : ""} ${t}`}
              onClick={() => set("type", t)}
            >
              {t === "expense" ? "Xarajat" : "Daromad"}
            </button>
          ))}
        </div>
      </label>

      <label>Summa (so'm)
        <input
          type="number"
          min="1"
          value={form.amount}
          onChange={(e) => set("amount", e.target.value)}
          placeholder="Masalan: 50000"
        />
      </label>

      <label>Kategoriya
        <select value={form.categoryId} onChange={(e) => set("categoryId", e.target.value)}>
          <option value="">— tanlang —</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>
              {c.icon} {c.name}
            </option>
          ))}
        </select>
      </label>

      <label>Sana
        <input
          type="date"
          value={form.date}
          onChange={(e) => set("date", e.target.value)}
        />
      </label>

      <label>Izoh (ixtiyoriy)
        <input
          type="text"
          value={form.description}
          onChange={(e) => set("description", e.target.value)}
          placeholder="Masalan: Do'kondan xarid"
        />
      </label>

      <div className="tf-actions">
        <button className="btn-save" onClick={handleSubmit} disabled={loading}>
          {loading ? "Saqlanmoqda..." : "Saqlash"}
        </button>
        {onCancel && (
          <button className="btn-cancel" onClick={onCancel} disabled={loading}>
            Bekor qilish
          </button>
        )}
      </div>
    </div>
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}

function toForm(t) {
  return {
    amount: t.amount || "",
    description: t.description || "",
    categoryId: t.categoryId || t.category_id || "",
    type: t.type || "expense",
    date: t.date?.slice(0, 10) || today(),
  };
}