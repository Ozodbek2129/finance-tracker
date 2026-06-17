package handler

import (
	"net/http"
	"strconv"

	postgrescrud "finance-tracker/internal/connectiondb/postgres_crud"
	models "finance-tracker/internal/model"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	FT *postgrescrud.FinanceTracker
}

func NewHandler(ft *postgrescrud.FinanceTracker) *Handler {
	return &Handler{FT: ft}
}

// ─── CATEGORY ─────────────────────────────────────────────────────────────────

// POST /categories
func (h *Handler) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.FT.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// PUT /categories/:id
func (h *Handler) UpdateCategory(c *gin.Context) {
	var req models.UpdateCategoryReq
	req.Id = c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.FT.UpdateCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GET /categories?page=1&limit=10
func (h *Handler) GetAllCategory(c *gin.Context) {
	var req models.GetAllCategoryReq

	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}
	req.Page = page

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	req.Limit = limit

	res, err := h.FT.GetAllCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DELETE /categories/:id
func (h *Handler) DeleteCategory(c *gin.Context) {
	req := models.DeleteCategoryReq{Id: c.Param("id")}

	res, err := h.FT.DeleteCategory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ─── TRANSACTION ──────────────────────────────────────────────────────────────

// POST /transactions
func (h *Handler) CreateTransaction(c *gin.Context) {
	var req models.CreateTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type != "income" && req.Type != "expense" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type faqat 'income' yoki 'expense' bo'lishi kerak"})
		return
	}

	res, err := h.FT.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// PUT /transactions/:id
func (h *Handler) UpdateTransaction(c *gin.Context) {
	var req models.UpdateTransactionReq
	req.ID = c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.FT.UpdateTransaction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GET /transactions?page=1&limit=10&type=expense&category_id=...&month=6&year=2026
func (h *Handler) GetAllTransactions(c *gin.Context) {
	var req models.GetAllTransactionReq

	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}
	req.Page = page

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	req.Limit = limit

	req.Type = c.Query("type")
	req.CategoryID = c.Query("category_id")

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 0 {
		month = 0
	}
	req.Month = month

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil || year < 0 {
		year = 0
	}
	req.Year = year

	res, err := h.FT.GetAllTransactions(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DELETE /transactions/:id
func (h *Handler) DeleteTransaction(c *gin.Context) {
	req := models.DeleteTransactionReq{ID: c.Param("id")}

	res, err := h.FT.DeleteTransaction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ─── STATS ────────────────────────────────────────────────────────────────────

// GET /stats?month=6&year=2026
func (h *Handler) GetStats(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 0 {
		month = 0
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil || year < 0 {
		year = 0
	}

	res, err := h.FT.GetStats(c.Request.Context(), month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}