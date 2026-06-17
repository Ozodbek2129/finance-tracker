package models

import (
	"time"
)

// ─── Category Models ───────────────────────────────────────────────────────

type CreateCategoryReq struct {
	Name  string    `json:"name"`
	Icon  string    `json:"icon"`
	Color string    `json:"color"`
}

type CreateCategoryRes struct {
	ID    string `json:"id"`
	Name  string    `json:"name"`
	Icon  string    `json:"icon"`
	Color string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateCategoryReq struct {
	Id	  string    `json:"id"`
	Name  string    `json:"name"`
	Icon  string    `json:"icon"`
	Color string    `json:"color"`
}

type UpdateCategoryRes struct {
	ID    string `json:"id"`
	Name  string    `json:"name"`
	Icon  string    `json:"icon"`
	Color string    `json:"color"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetAllCategoryReq struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"` 
}

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetAllCategoryRes struct {
	Categories []Category `json:"categories"`
}

type DeleteCategoryReq struct {
	Id string `json:"id"`
}

type DeleteCategoryRes struct {
	Message string `json:"message"`
}

// ─── Transaction Models ───────────────────────────────────────────────────────

type Transaction struct {
	ID          string    `json:"id"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	CategoryID  string    `json:"category_id"`
	Category    Category  `json:"category"`   // JOIN bilan keladi
	Type        string    `json:"type"`        // "income" yoki "expense"
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTransactionReq struct {
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	CategoryID  string `json:"category_id"`
	Type        string `json:"type"` // "income" | "expense"
	Date        string `json:"date"` // "2025-01-15" formatida
}

type CreateTransactionRes struct {
	ID          string    `json:"id"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	CategoryID  string    `json:"category_id"`
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateTransactionReq struct {
	ID          string `json:"id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	CategoryID  string `json:"category_id"`
	Type        string `json:"type"`
	Date        string `json:"date"`
}

type UpdateTransactionRes struct {
	ID          string    `json:"id"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	CategoryID  string    `json:"category_id"`
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetAllTransactionReq struct {
	Limit      int64  `json:"limit"`
	Page       int64  `json:"page"`
	Type       string `json:"type"`        // "income" | "expense" | "" (barchasi)
	CategoryID string `json:"category_id"` // bo'sh bo'lsa barchasi
	Month      int    `json:"month"`       // 0 bo'lsa barchasi, aks holda 1-12
	Year       int    `json:"year"`        // 0 bo'lsa barchasi
}

type GetAllTransactionRes struct {
	Transactions []Transaction `json:"transactions"`
	Total        int64         `json:"total"` // pagination uchun
}

type DeleteTransactionReq struct {
	ID string `json:"id"`
}

type DeleteTransactionRes struct {
	Message string `json:"message"`
}

// ─── STATS Models ───────────────────────────────────────────────────────

type MonthlySummary struct {
	TotalIncome  int64 `json:"total_income"`
	TotalExpense int64 `json:"total_expense"`
	Balance      int64 `json:"balance"`
}

type CategorySummary struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	CategoryIcon string `json:"category_icon"`
	TotalAmount  int64  `json:"total_amount"`
	Count        int    `json:"count"`
}


type HighestExpenseDay struct {
	Date   string `json:"date"`
	Amount int64  `json:"amount"`
}

type MostUsedCategory struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Icon             string `json:"icon"`
	TransactionCount int    `json:"transaction_count"`
}

type Insight struct {
	Type    string `json:"type"`    // "warning" | "good" | "info"
	Title   string `json:"title"`
	Message string `json:"message"`
}

type GetStatsRes struct {
	Summary    MonthlySummary    `json:"summary"`
	ByCategory []CategorySummary `json:"by_category"`

	TopExpenseCategory    *CategorySummary   `json:"top_expense_category"`    // Eng ko'p pul ketgan kategoriya
	MostUsedCategory      *MostUsedCategory  `json:"most_used_category"`      // Eng ko'p transaction bo'lgan
	HighestExpenseDay     *HighestExpenseDay `json:"highest_expense_day"`     // Eng ko'p xarajat qilingan kun
	AverageDailyExpense   int64              `json:"average_daily_expense"`   // Kunlik o'rtacha xarajat
	SavingsRate           float64            `json:"savings_rate"`            // Tejash foizi
	ExpenseChangePercent  float64            `json:"expense_change_percent"`  // O'tgan oyga nisbatan xarajat o'zgarishi
	IncomeChangePercent   float64            `json:"income_change_percent"`   // O'tgan oyga nisbatan daromad o'zgarishi
	RiskLevel             string             `json:"risk_level"`              // "low" | "medium" | "high"
	FinancialScore        int                `json:"financial_score"`         // 0-100
	Insights              []Insight          `json:"insights"`                // AI tavsiyalar
}