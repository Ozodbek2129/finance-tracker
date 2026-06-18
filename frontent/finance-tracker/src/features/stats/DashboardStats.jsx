import { useEffect, useState, useRef } from "react";
import { statsAPI } from "../../api/client";

const now = new Date();

// ── 1. RAQAMLARNI SILLIQ O'STIRUVCHI KOMPONENT (Asosiy qism va Tavsiyalar uchun umumiy) ──
function AnimatedNumber({ value, isPercent = false, trigger = true }) {
  const [current, setCurrent] = useState(0);
  const [isIntersecting, setIsIntersecting] = useState(false);
  const elementRef = useRef(null);

  useEffect(() => {
    if (!trigger) {
      setCurrent(0);
      return;
    }

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsIntersecting(true);
          observer.unobserve(entry.target);
        }
      },
      { threshold: 0.1 }
    );
    if (elementRef.current) observer.observe(elementRef.current);
    return () => observer.disconnect();
  }, [value, trigger]);

  useEffect(() => {
    if (!isIntersecting && trigger) {
      setIsIntersecting(true);
    }
    if (!isIntersecting) return;

    let start = 0;
    const end = Number(value || 0);
    if (end === 0) { setCurrent(0); return; }

    const duration = 1200;
    const startTime = performance.now();

    function updateNumber(currentTime) {
      const elapsedTime = currentTime - startTime;
      if (elapsedTime >= duration) {
        setCurrent(end);
        return;
      }
      const progress = elapsedTime / duration;
      const easeOutProgress = 1 - Math.pow(1 - progress, 3);
      
      const nextValue = isPercent 
        ? parseFloat((easeOutProgress * end).toFixed(1))
        : Math.floor(easeOutProgress * end);
        
      setCurrent(nextValue);
      requestAnimationFrame(updateNumber);
    }
    requestAnimationFrame(updateNumber);
  }, [isIntersecting, value, trigger, isPercent]);

  return <span ref={elementRef}>{isPercent ? current : fmt(current)}</span>;
}

// ── 2. TAVSIYALAR MATNIDAGI RAQAMLAR VA YOZUVLARNI BIR QATORGA JAMLASH ──
function RenderMessageWithAnimation({ message, trigger }) {
  if (!message) return null;
  
  const regex = /(\d+(?:\.\d+)?[%])|(\d+)/g;
  const parts = message.split(regex).filter(Boolean);

  return (
    // display: "block" va ichidagi matn uzilib ketmasligi uchun o'zgarishlar kiritildi
    <p style={{ margin: "6px 0 0 0", fontSize: "14px", opacity: 0.9, lineHeight: "1.6" }}>
      {parts.map((part, index) => {
        if (part.endsWith("%")) {
          const numValue = parseFloat(part.replace("%", ""));
          return (
            <strong key={index} style={{ display: "inline", whiteSpace: "nowrap", padding: "0 4px", fontSize: "15px" }}>
              <AnimatedNumber value={numValue} isPercent={true} trigger={trigger} />%
            </strong>
          );
        }
        if (/^\d+$/.test(part)) {
          const numValue = parseInt(part, 10);
          return (
            <strong key={index} style={{ display: "inline", whiteSpace: "nowrap", padding: "0 4px", fontSize: "15px" }}>
              <AnimatedNumber value={numValue} isPercent={false} trigger={trigger} />
            </strong>
          );
        }
        // Oddiy matn qismlarini inline formatda yonma-yon chiqaradi
        return <span key={index} style={{ display: "inline" }}>{part}</span>;
      })}
    </p>
  );
}

// ── 3. SKROLDA ACHILADIGAN VA RAQAMLARI O'SADIGAN TAVSIYA KARTASI ──
function InsightCard({ ins }) {
  const [visible, setVisible] = useState(false);
  const cardRef = useRef(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setVisible(true);
          observer.unobserve(entry.target);
        }
      },
      {
        rootMargin: "-5% 0px -5% 0px",
        threshold: 0.15
      }
    );

    if (cardRef.current) observer.observe(cardRef.current);
    return () => observer.disconnect();
  }, []);

  return (
    <div
      ref={cardRef}
      className={`insight ${ins.type} ${visible ? "insight-visible" : "insight-hidden"}`}
    >
      <strong style={{ fontSize: "16px", display: "block" }}>{ins.title}</strong>
      <RenderMessageWithAnimation message={ins.message} trigger={visible} />
    </div>
  );
}

// ── 4. ASOSIY STATISTIKA KOMPONENTI ──
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
      <style>{`
        .insight-hidden {
          opacity: 0 !important;
          transform: translateY(35px) scale(0.96);
          transition: opacity 0.4s ease-out, transform 0.4s ease-out;
        }
        .insight-visible {
          opacity: 1 !important;
          transform: translateY(0) scale(1);
          transition: opacity 0.5s ease-out, transform 0.5s ease-out;
        }
        .stats-insights .insight {
          margin-bottom: 12px;
          padding: 14px;
          border-radius: 12px;
        }
      `}</style>

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
          <span className="card-value"><AnimatedNumber value={summary.total_income} /> so'm</span>
        </div>
        <div className="card expense">
          <span className="card-label">Xarajat</span>
          <span className="card-value"><AnimatedNumber value={summary.total_expense} /> so'm</span>
        </div>
        <div className={`card balance ${summary.balance >= 0 ? "positive" : "negative"}`}>
          <span className="card-label">Balans</span>
          <span className="card-value"><AnimatedNumber value={summary.balance} /> so'm</span>
        </div>
      </div>

      {/* Score & Risk */}
      <div className="stats-meta">
        <div className="meta-item">
          <span className="meta-label">Moliyaviy ball</span>
          <div className="score-bar">
            <div className="score-fill" style={{ width: `${financial_score}%` }} />
          </div>
          <span className="meta-value"><AnimatedNumber value={financial_score} isPercent={true} />/100</span>
        </div>
        <div className="meta-item">
          <span className="meta-label">Tejash</span>
          <span className="meta-value"><AnimatedNumber value={savings_rate} isPercent={true} />%</span>
        </div>
        <div className="meta-item">
          <span className="meta-label">Xavf darajasi</span>
          <span className={`risk-badge ${risk_level}`}>
            {risk_level === "low" ? "Past" : risk_level === "medium" ? "O'rta" : "Yuqori"}
          </span>
        </div>
        <div className="meta-item">
          <span className="meta-label">Kunlik o'rtacha</span>
          <span className="meta-value"><AnimatedNumber value={average_daily_expense} /> so'm</span>
        </div>
        {highest_expense_day && (
          <div className="meta-item">
            <span className="meta-label">Eng ko'p xarajat kuni</span>
            <span className="meta-value">
              {highest_expense_day.date} — <AnimatedNumber value={highest_expense_day.amount} /> so'm
            </span>
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
                <span className="cat-pct"><AnimatedNumber value={pct} isPercent={true} />%</span>
                <span className="cat-amount"><AnimatedNumber value={cat.total_amount} /> so'm</span>
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
            <InsightCard key={i} ins={ins} />
          ))}
        </div>
      )}
    </div>
  );
}

function fmt(n) {
  return Number(n || 0).toLocaleString("uz-UZ");
}