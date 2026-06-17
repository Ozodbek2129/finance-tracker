package postgrescrud

import (
	"context"
	"fmt"
	models "finance-tracker/internal/model"
	"math"
)

// ─── GET STATS ────────────────────────────────────────────────────────────────

func (f *FinanceTracker) GetStats(ctx context.Context, month, year int) (*models.GetStatsRes, error) {

	// ── 1. WHERE sharti yasash ────────────────────────────────────────────────
	where := "WHERE t.deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if month > 0 {
		where += fmt.Sprintf(" AND EXTRACT(MONTH FROM t.date) = $%d", argIdx)
		args = append(args, month)
		argIdx++
	}
	if year > 0 {
		where += fmt.Sprintf(" AND EXTRACT(YEAR FROM t.date) = $%d", argIdx)
		args = append(args, year)
		argIdx++
	}

	// ── 2. O'tgan oy WHERE sharti ─────────────────────────────────────────────
	prevMonth, prevYear := prevMonthYear(month, year)
	prevWhere := "WHERE t.deleted_at IS NULL"
	prevArgs := []interface{}{}
	prevIdx := 1

	if prevMonth > 0 {
		prevWhere += fmt.Sprintf(" AND EXTRACT(MONTH FROM t.date) = $%d", prevIdx)
		prevArgs = append(prevArgs, prevMonth)
		prevIdx++
	}
	if prevYear > 0 {
		prevWhere += fmt.Sprintf(" AND EXTRACT(YEAR FROM t.date) = $%d", prevIdx)
		prevArgs = append(prevArgs, prevYear)
	}

	// ── 3. Asosiy summary (income, expense, balance) ──────────────────────────
	summary, err := f.fetchSummary(ctx, where, args)
	if err != nil {
		return nil, err
	}

	// ── 4. O'tgan oy summary (taqqoslash uchun) ───────────────────────────────
	prevSummary, err := f.fetchSummary(ctx, prevWhere, prevArgs)
	if err != nil {
		return nil, err
	}

	// ── 5. Kategoriya bo'yicha xarajatlar ─────────────────────────────────────
	byCategory, err := f.fetchCategoryBreakdown(ctx, where, args)
	if err != nil {
		return nil, err
	}

	// ── 6. Eng ko'p xarajat qilingan kun ─────────────────────────────────────
	highestDay, err := f.fetchHighestExpenseDay(ctx, where, args)
	if err != nil {
		return nil, err
	}

	// ── 7. Kunlik o'rtacha xarajat ────────────────────────────────────────────
	avgDaily, err := f.fetchAverageDailyExpense(ctx, where, args)
	if err != nil {
		return nil, err
	}

	// ── 8. Eng ko'p transaction bo'lgan kategoriya ────────────────────────────
	mostUsed, err := f.fetchMostUsedCategory(ctx, where, args)
	if err != nil {
		return nil, err
	}

	// ── 9. Hisob-kitoblar ─────────────────────────────────────────────────────
	savingsRate := calcSavingsRate(summary.TotalIncome, summary.TotalExpense)
	expenseChange := calcChangePercent(prevSummary.TotalExpense, summary.TotalExpense)
	incomeChange := calcChangePercent(prevSummary.TotalIncome, summary.TotalIncome)
	riskLevel := calcRiskLevel(summary.TotalIncome, summary.TotalExpense)
	score := calcFinancialScore(savingsRate, riskLevel, expenseChange, byCategory, summary)

	// ── 10. Eng ko'p va eng kam xarajat kategoriyalari ───────────────────────
	var topCategory *models.CategorySummary
	if len(byCategory) > 0 {
		topCategory = &byCategory[0]
	}

	// ── 11. AI Insights (tavsiyalar) ─────────────────────────────────────────
	insights := generateInsights(summary, byCategory, savingsRate, expenseChange, riskLevel)

	return &models.GetStatsRes{
		Summary:             *summary,
		ByCategory:          byCategory,
		TopExpenseCategory:  topCategory,
		MostUsedCategory:    mostUsed,
		HighestExpenseDay:   highestDay,
		AverageDailyExpense: avgDaily,
		SavingsRate:         savingsRate,
		ExpenseChangePercent: expenseChange,
		IncomeChangePercent: incomeChange,
		RiskLevel:           riskLevel,
		FinancialScore:      score,
		Insights:            insights,
	}, nil
}

// ─── HELPER: Summary olish ────────────────────────────────────────────────────

func (f *FinanceTracker) fetchSummary(ctx context.Context, where string, args []interface{}) (*models.MonthlySummary, error) {
	query := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)
		FROM transactions t
		%s
	`, where)

	var s models.MonthlySummary
	err := f.DB.QueryRowContext(ctx, query, args...).Scan(&s.TotalIncome, &s.TotalExpense)
	if err != nil {
		f.Log.Error("Error fetching summary", "err", err)
		return nil, err
	}
	s.Balance = s.TotalIncome - s.TotalExpense
	return &s, nil
}

// ─── HELPER: Kategoriya breakdown ─────────────────────────────────────────────

func (f *FinanceTracker) fetchCategoryBreakdown(ctx context.Context, where string, args []interface{}) ([]models.CategorySummary, error) {
	query := fmt.Sprintf(`
		SELECT
			c.id,
			c.name,
			c.icon,
			COALESCE(SUM(t.amount), 0) AS total_amount,
			COUNT(t.id) AS count
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		%s AND t.type = 'expense'
		GROUP BY c.id, c.name, c.icon
		ORDER BY total_amount DESC
	`, where)

	rows, err := f.DB.QueryContext(ctx, query, args...)
	if err != nil {
		f.Log.Error("Error fetching category breakdown", "err", err)
		return nil, err
	}
	defer rows.Close()

	var result []models.CategorySummary
	for rows.Next() {
		var cs models.CategorySummary
		if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.CategoryIcon, &cs.TotalAmount, &cs.Count); err != nil {
			return nil, err
		}
		result = append(result, cs)
	}
	return result, rows.Err()
}

// ─── HELPER: Eng yuqori xarajat kuni ─────────────────────────────────────────

func (f *FinanceTracker) fetchHighestExpenseDay(ctx context.Context, where string, args []interface{}) (*models.HighestExpenseDay, error) {
	query := fmt.Sprintf(`
		SELECT
			TO_CHAR(t.date, 'YYYY-MM-DD'),
			SUM(t.amount) AS total
		FROM transactions t
		%s AND t.type = 'expense'
		GROUP BY t.date
		ORDER BY total DESC
		LIMIT 1
	`, where)

	var h models.HighestExpenseDay
	err := f.DB.QueryRowContext(ctx, query, args...).Scan(&h.Date, &h.Amount)
	if err != nil {
		// Ma'lumot yo'q bo'lsa — nil qaytaramiz, xato emas
		return nil, nil
	}
	return &h, nil
}

// ─── HELPER: Kunlik o'rtacha xarajat ─────────────────────────────────────────

func (f *FinanceTracker) fetchAverageDailyExpense(ctx context.Context, where string, args []interface{}) (int64, error) {
	query := fmt.Sprintf(`
		SELECT COALESCE(AVG(daily_total), 0)
		FROM (
			SELECT SUM(t.amount) AS daily_total
			FROM transactions t
			%s AND t.type = 'expense'
			GROUP BY t.date
		) sub
	`, where)

	var avg float64
	err := f.DB.QueryRowContext(ctx, query, args...).Scan(&avg)
	if err != nil {
		f.Log.Error("Error fetching average daily expense", "err", err)
		return 0, err
	}
	return int64(math.Round(avg)), nil
}

// ─── HELPER: Eng ko'p ishlatilgan kategoriya ──────────────────────────────────

func (f *FinanceTracker) fetchMostUsedCategory(ctx context.Context, where string, args []interface{}) (*models.MostUsedCategory, error) {
	query := fmt.Sprintf(`
		SELECT
			c.id,
			c.name,
			c.icon,
			COUNT(t.id) AS transaction_count
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		%s
		GROUP BY c.id, c.name, c.icon
		ORDER BY transaction_count DESC
		LIMIT 1
	`, where)

	var m models.MostUsedCategory
	err := f.DB.QueryRowContext(ctx, query, args...).Scan(&m.ID, &m.Name, &m.Icon, &m.TransactionCount)
	if err != nil {
		return nil, nil
	}
	return &m, nil
}

// ─── HISOB-KITOB FUNKSIYALARI ─────────────────────────────────────────────────

// O'tgan oy/yilni hisoblash
func prevMonthYear(month, year int) (int, int) {
	if month == 0 || year == 0 {
		return 0, 0
	}
	if month == 1 {
		return 12, year - 1
	}
	return month - 1, year
}

// Tejash foizi: (income - expense) / income * 100
func calcSavingsRate(income, expense int64) float64 {
	if income == 0 {
		return 0
	}
	rate := float64(income-expense) / float64(income) * 100
	return math.Round(rate*100) / 100
}

// O'zgarish foizi: (yangi - eski) / eski * 100
func calcChangePercent(prev, current int64) float64 {
	if prev == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	change := float64(current-prev) / float64(prev) * 100
	return math.Round(change*100) / 100
}

// Risk darajasi
func calcRiskLevel(income, expense int64) string {
	if income == 0 {
		return "high"
	}
	ratio := float64(expense) / float64(income) * 100
	switch {
	case ratio < 50:
		return "low"
	case ratio <= 80:
		return "medium"
	default:
		return "high"
	}
}

// Moliyaviy ball (0-100)
func calcFinancialScore(savingsRate float64, riskLevel string, expenseChange float64, byCategory []models.CategorySummary, summary *models.MonthlySummary) int {
	score := 0

	// Tejash > 30% → +30 ball
	if savingsRate >= 30 {
		score += 30
	} else if savingsRate >= 15 {
		score += 15
	}

	// Xarajat < daromad → +20 ball
	if summary.TotalExpense < summary.TotalIncome {
		score += 20
	}

	// Risk darajasi → +20 ball
	switch riskLevel {
	case "low":
		score += 20
	case "medium":
		score += 10
	}

	// Top kategoriya umumiy xarajatning 40% dan kami → +15 ball
	if len(byCategory) > 0 && summary.TotalExpense > 0 {
		topRatio := float64(byCategory[0].TotalAmount) / float64(summary.TotalExpense) * 100
		if topRatio < 40 {
			score += 15
		}
	}

	// Xarajat o'tgan oyga nisbatan kamaygan → +15 ball
	if expenseChange < 0 {
		score += 15
	}

	if score > 100 {
		score = 100
	}
	return score
}

// ─── AI INSIGHTS (tavsiyalar) ─────────────────────────────────────────────────

func generateInsights(summary *models.MonthlySummary, byCategory []models.CategorySummary, savingsRate float64, expenseChange float64, riskLevel string) []models.Insight {
	var insights []models.Insight

	// 1. Risk darajasi
	switch riskLevel {
	case "high":
		insights = append(insights, models.Insight{
			Type:    "warning",
			Title:   "Yuqori xavf",
			Message: fmt.Sprintf("Daromadingizning 80%% dan ko'prog'i xarajatga ketmoqda. Tejashni boshlash vaqti keldi."),
		})
	case "medium":
		insights = append(insights, models.Insight{
			Type:    "info",
			Title:   "O'rtacha xavf",
			Message: "Xarajatlaringiz daromadingizning 50-80%% orasida. Nazorat qiling.",
		})
	case "low":
		insights = append(insights, models.Insight{
			Type:    "good",
			Title:   "Yaxshi holat",
			Message: "Xarajatlaringiz daromadingizning yarmidan kami. Zo'r!",
		})
	}

	// 2. Tejash foizi
	if savingsRate > 0 {
		insights = append(insights, models.Insight{
			Type:    "good",
			Title:   "Tejash",
			Message: fmt.Sprintf("Bu oy daromadingizning %.1f%% ini tejayapsiz.", savingsRate),
		})
	} else if summary.TotalIncome > 0 {
		insights = append(insights, models.Insight{
			Type:    "warning",
			Title:   "Tejash yo'q",
			Message: "Bu oy tejashingiz yo'q. Xarajatlarni kamaytirishga harakat qiling.",
		})
	}

	// 3. Xarajat o'zgarishi
	if expenseChange > 20 {
		insights = append(insights, models.Insight{
			Type:    "warning",
			Title:   "Xarajat oshdi",
			Message: fmt.Sprintf("Xarajatlaringiz o'tgan oyga nisbatan %.1f%% oshgan.", expenseChange),
		})
	} else if expenseChange < -10 {
		insights = append(insights, models.Insight{
			Type:    "good",
			Title:   "Xarajat kamaydi",
			Message: fmt.Sprintf("Xarajatlaringiz o'tgan oyga nisbatan %.1f%% kamaygan. Davom eting!", -expenseChange),
		})
	}

	// 4. Top kategoriya ogohlantiruvi
	if len(byCategory) > 0 && summary.TotalExpense > 0 {
		top := byCategory[0]
		topRatio := float64(top.TotalAmount) / float64(summary.TotalExpense) * 100
		if topRatio > 50 {
			insights = append(insights, models.Insight{
				Type:    "warning",
				Title:   top.CategoryName,
				Message: fmt.Sprintf("Xarajatlaringizning %.1f%% i '%s' kategoriyasiga ketmoqda.", topRatio, top.CategoryName),
			})
		}
	}

	// 5. Balans manfiy bo'lsa
	if summary.Balance < 0 {
		insights = append(insights, models.Insight{
			Type:    "warning",
			Title:   "Salbiy balans",
			Message: fmt.Sprintf("Bu oy %d so'm qarzga botdingiz. Xarajatlarni ko'rib chiqing.", -summary.Balance),
		})
	}

	return insights
}