import { useEffect, useState } from "react";
import { statsAPI } from "../../api/client";

const now = new Date();

export default function DashboardStats() {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [month, setMonth] = useState(now.getMonth() + 1);
  const [year, setYear] = useState(now.getFullYear());

  useEffect(() => {
    setLoading(true);
    setError(null);
    statsAPI
      .get({ month, year })
      .then(setStats)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [month, year]);

  if (loading) return <div className="stats-loading">Yuklanmoqda...</div>;
  if (error) return <div className="stats-error">Xato: {error}</div>;
  if (!stats) return null;

  const {
    summary,
    savings_rate,
    risk_level,
    financial_score,
    insights,
    by_category,
    highest_expense_day,
    average_daily_expense,
  } = stats;

  return (
    <div className="stats-container">
      {/* Filter */}
      <div className="stats-filter">
        <select value={month} onChange={(e) => setMonth(+e.target.value)}>
          {Array.from({ length: 12 }, (_, i) => (
            <option key={i + 1} value={i + 1}>
              {new Date(2000, i).toLocaleString("uz", { month: "long" })}
            </option>
          ))}
        </select>
        <select value={year} onChange={(e) => setYear(+e.target.value)}>
          {[2023, 2024, 2025, 2026].map((y) => (
            <option key={y} value={y}>{y}</option>
          ))}
        </select>
      </div>

      {/* Summary cards */}
      <div className="stats-cards">
        <div className="card income">
          <span className="card-label">Daromad</span>
          <span className="card-value">{fmt(summary.total_income)} so'm</span>
        </div>
        <div className="card expense">
          <span className="card-label">Xarajat</span>
          <span className="card-value">{fmt(summary.total_expense)} so'm</span>
        </div>
        <div className={`card balance ${summary.balance >= 0 ? "positive" : "negative"}`}>
          <span className="card-label">Balans</span>
          <span className="card-value">{fmt(summary.balance)} so'm</span>
        </div>
      </div>

      {/* Score & Risk */}
      <div className="stats-meta">
        <div className="meta-item">
          <span className="meta-label">Moliyaviy ball</span>
          <div className="score-bar">
            <div className="score-fill" style={{ width: `${financial_score}%` }} />
          </div>
          <span className="meta-value">{financial_score}/100</span>
        </div>
        <div className="meta-item">
          <span className="meta-label">Tejash</span>
          <span className="meta-value">{savings_rate}%</span>
        </div>
        <div className="meta-item">
          <span className="meta-label">Xavf darajasi</span>
          <span className={`risk-badge ${risk_level}`}>
            {risk_level === "low" ? "Past" : risk_level === "medium" ? "O'rta" : "Yuqori"}
          </span>
        </div>
        <div className="meta-item">
          <span className="meta-label">Kunlik o'rtacha</span>
          <span className="meta-value">{fmt(average_daily_expense)} so'm</span>
        </div>
        {highest_expense_day && (
          <div className="meta-item">
            <span className="meta-label">Eng ko'p xarajat kuni</span>
            <span className="meta-value">{highest_expense_day.date} — {fmt(highest_expense_day.amount)} so'm</span>
          </div>
        )}
      </div>

      {/* Category breakdown */}
      {by_category?.length > 0 && (
        <div className="stats-categories">
          <h3>Kategoriyalar bo'yicha xarajat</h3>
          {by_category.map((cat) => {
            const pct = summary.total_expense > 0
              ? Math.round((cat.total_amount / summary.total_expense) * 100)
              : 0;
            return (
              <div key={cat.category_id} className="cat-row">
                <span className="cat-icon">{cat.category_icon}</span>
                <span className="cat-name">{cat.category_name}</span>
                <div className="cat-bar-wrap">
                  <div className="cat-bar" style={{ width: `${pct}%` }} />
                </div>
                <span className="cat-pct">{pct}%</span>
                <span className="cat-amount">{fmt(cat.total_amount)} so'm</span>
              </div>
            );
          })}
        </div>
      )}

      {/* Insights */}
      {insights?.length > 0 && (
        <div className="stats-insights">
          <h3>Tavsiyalar</h3>
          {insights.map((ins, i) => (
            <div key={i} className={`insight ${ins.type}`}>
              <strong>{ins.title}</strong>
              <p>{ins.message}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function fmt(n) {
  return Number(n || 0).toLocaleString("uz-UZ");
}