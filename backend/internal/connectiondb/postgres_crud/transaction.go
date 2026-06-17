package postgrescrud

import (
	"context"
	"database/sql"
	"errors"
	"finance-tracker/internal/logger"
	models "finance-tracker/internal/model"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type FinanceTracker struct {
	DB  *sql.DB
	Log *slog.Logger
}

func NewFinanceTracker(db *sql.DB) *FinanceTracker {
	return &FinanceTracker{
		DB:  db,
		Log: logger.NewLogger(),
	}
}

// ─── CREATE CATEGORY ──────────────────────────────────────────────────────────────────

func (f *FinanceTracker) CreateCategory(ctx context.Context, req *models.CreateCategoryReq) (*models.CreateCategoryRes, error) {
	id := uuid.NewString()
	newTime := time.Now()

	query := `insert into categories(id, name, icon, color, 
									created_at, updated_at
								) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := f.DB.ExecContext(ctx, query, id, req.Name, req.Icon, req.Color, newTime, newTime)
	if err != nil {
		f.Log.Error("Error inserting category", "err", err)
		return nil, err
	}

	return &models.CreateCategoryRes{
		ID:        id,
		Name:      req.Name,
		Icon:      req.Icon,
		Color:     req.Color,
		CreatedAt: newTime,
	}, nil
}

// ─── UPDATE CATEGORY ──────────────────────────────────────────────────────────────────

func (f *FinanceTracker) UpdateCategory(ctx context.Context, req *models.UpdateCategoryReq) (*models.UpdateCategoryRes, error) {
	query := `
		UPDATE categories
		SET
			name = COALESCE(NULLIF($1, ''), name),
			icon = COALESCE(NULLIF($2, ''), icon),
			color = COALESCE(NULLIF($3, ''), color),
			updated_at = $4
		WHERE id = $5
		AND deleted_at IS NULL
		RETURNING id, name, icon, color, updated_at
	`

	newtime := time.Now()
	var res models.UpdateCategoryRes

	err := f.DB.QueryRowContext(
		ctx,
		query,
		req.Name,
		req.Icon,
		req.Color,
		newtime,
		req.Id,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Icon,
		&res.Color,
		&res.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			f.Log.Error("category not found", "err", err)
			return nil, errors.New("category not found")
		}

		f.Log.Error("Error updating category", "err", err)
		return nil, err
	}

	return &res, nil
}

// ─── GET ALL CATEGORY (filter bilan) ───────────────────────────────────────────────────
func (f *FinanceTracker) GetAllCategory(ctx context.Context, req *models.GetAllCategoryReq) (*models.GetAllCategoryRes, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	query := `select id, name, icon, color, created_at, updated_at
			  from categories
			  where deleted_at IS NULL
			  order by created_at
			  LIMIT $1 OFFSET $2`

	rows, err := f.DB.QueryContext(ctx, query, req.Limit, offset)
	if err != nil {
		f.Log.Error("Error fetching category", "err", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Icon,
			&category.Color,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.GetAllCategoryRes{
		Categories: categories,
	}, nil
}

// ─── DELETE CATEGORY (soft delete) ─────────────────────────────────────────────────────

func (f *FinanceTracker) DeleteCategory(ctx context.Context, req *models.DeleteCategoryReq) (*models.DeleteCategoryRes, error) {
	query := `update categories
			  set deleted_at = $1
			  where id = $2`

	newtime := time.Now()
	res, err := f.DB.ExecContext(ctx, query, newtime, req.Id)
	if err != nil {
		f.Log.Error("Error delete category", "err", err)
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		f.Log.Error("Could not delete category", "err", err)
		return &models.DeleteCategoryRes{
			Message: "Could not delete category",
		}, nil
	}

	return &models.DeleteCategoryRes{
		Message: "Contract Delete Updated successfully",
	}, nil
}

// ─── CREATE TRANSACTION ──────────────────────────────────────────────────────────────────

func (f *FinanceTracker) CreateTransaction(ctx context.Context, req *models.CreateTransactionReq) (*models.CreateTransactionRes, error) {
	id := uuid.NewString()
	now := time.Now()

	// Date string dan time.Time ga o'tkazish
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("date format noto'g'ri, 'YYYY-MM-DD' formatida yuboring")
	}

	query := `
		INSERT INTO transactions (id, amount, description, category_id, type, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = f.DB.ExecContext(ctx, query,
		id, req.Amount, req.Description, req.CategoryID, req.Type, date, now, now,
	)
	if err != nil {
		f.Log.Error("Error creating transaction", "err", err)
		return nil, err
	}

	return &models.CreateTransactionRes{
		ID:          id,
		Amount:      req.Amount,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Type:        req.Type,
		Date:        date,
		CreatedAt:   now,
	}, nil
}

// ─── GET ALL TRANSACTION (filter bilan) ───────────────────────────────────────────────────

func (f *FinanceTracker) GetAllTransactions(ctx context.Context, req *models.GetAllTransactionReq) (*models.GetAllTransactionRes, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit

	where := "WHERE t.deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if req.Type != "" {
		where += fmt.Sprintf(" AND t.type = $%d", argIdx)
		args = append(args, req.Type)
		argIdx++
	}
	if req.CategoryID != "" {
		where += fmt.Sprintf(" AND t.category_id = $%d", argIdx)
		args = append(args, req.CategoryID)
		argIdx++
	}
	if req.Month > 0 {
		where += fmt.Sprintf(" AND EXTRACT(MONTH FROM t.date) = $%d", argIdx)
		args = append(args, req.Month)
		argIdx++
	}
	if req.Year > 0 {
		where += fmt.Sprintf(" AND EXTRACT(YEAR FROM t.date) = $%d", argIdx)
		args = append(args, req.Year)
		argIdx++
	}

	// Jami sonini olish (pagination uchun)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM transactions t %s`, where)
	var total int64
	err := f.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		f.Log.Error("Error counting transactions", "err", err)
		return nil, err
	}

	// LIMIT, OFFSET qo'shamiz
	args = append(args, req.Limit, offset)
	query := fmt.Sprintf(`
		SELECT 
			t.id, t.amount, t.description, t.category_id,
			t.type, t.date, t.created_at, t.updated_at,
			c.id, c.name, c.icon, c.color
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		%s
		ORDER BY t.date DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)

	rows, err := f.DB.QueryContext(ctx, query, args...)
	if err != nil {
		f.Log.Error("Error fetching transactions", "err", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID, &t.Amount, &t.Description, &t.CategoryID,
			&t.Type, &t.Date, &t.CreatedAt, &t.UpdatedAt,
			&t.Category.ID, &t.Category.Name, &t.Category.Icon, &t.Category.Color,
		)
		if err != nil {
			f.Log.Error("Error scanning transaction", "err", err)
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &models.GetAllTransactionRes{
		Transactions: transactions,
		Total:        total,
	}, nil
}

// ─── UPDATE TRANSACTION ──────────────────────────────────────────────────────────────────

func (f *FinanceTracker) UpdateTransaction(ctx context.Context, req *models.UpdateTransactionReq) (*models.UpdateTransactionRes, error) {
	now := time.Now()

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("date format noto'g'ri, 'YYYY-MM-DD' formatida yuboring")
	}

	query := `
		UPDATE transactions
		SET
			amount      = CASE WHEN $1 != 0 THEN $1 ELSE amount END,
			description = COALESCE(NULLIF($2, ''), description),
			category_id = COALESCE(NULLIF($3, '')::uuid, category_id),
			type        = COALESCE(NULLIF($4, ''), type),
			date        = COALESCE($5::date, date),
			updated_at  = $6
		WHERE id = $7
		  AND deleted_at IS NULL
		RETURNING id, amount, description, category_id, type, date, updated_at
	`

	var res models.UpdateTransactionRes
	err = f.DB.QueryRowContext(ctx, query,
		req.Amount, req.Description, req.CategoryID, req.Type, date, now, req.ID,
	).Scan(
		&res.ID, &res.Amount, &res.Description,
		&res.CategoryID, &res.Type, &res.Date, &res.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("transaction topilmadi")
		}
		f.Log.Error("Error updating transaction", "err", err)
		return nil, err
	}

	return &res, nil
}

// ─── DELETE TRANSACTION (soft delete) ─────────────────────────────────────────────────────

func (f *FinanceTracker) DeleteTransaction(ctx context.Context, req *models.DeleteTransactionReq) (*models.DeleteTransactionRes, error) {
	query := `
		UPDATE transactions
		SET deleted_at = $1
		WHERE id = $2
		  AND deleted_at IS NULL
	`

	res, err := f.DB.ExecContext(ctx, query, time.Now().Unix(), req.ID)
	if err != nil {
		f.Log.Error("Error deleting transaction", "err", err)
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return &models.DeleteTransactionRes{
			Message: "Transaction topilmadi yoki allaqachon o'chirilgan",
		}, nil
	}

	return &models.DeleteTransactionRes{
		Message: "Transaction muvaffaqiyatli o'chirildi",
	}, nil
}