import { useState } from "react";
import DashboardStats from "./features/stats/DashboardStats";
import TransactionList from "./features/transactions/TransactionList";
import CategoryManager from "./features/categories/CategoryManager";
import "./App.css";

const TABS = [
  { id: "stats", label: "📊 Statistika" },
  { id: "transactions", label: "💸 Tranzaksiyalar" },
  { id: "categories", label: "🗂 Kategoriyalar" },
];

export default function App() {
  const [tab, setTab] = useState("stats");

  return (
    <div className="app">
      <header className="app-header">
        <div className="app-logo">💰 Finance Tracker</div>
        <nav className="app-nav">
          {TABS.map((t) => (
            <button
              key={t.id}
              className={`nav-btn ${tab === t.id ? "active" : ""}`}
              onClick={() => setTab(t.id)}
            >
              {t.label}
            </button>
          ))}
        </nav>
      </header>

      <main className="app-main">
        {tab === "stats" && <DashboardStats />}
        {tab === "transactions" && <TransactionList />}
        {tab === "categories" && <CategoryManager />}
      </main>
    </div>
  );
}