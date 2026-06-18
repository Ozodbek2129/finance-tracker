package handler

import (
	corss "finance-tracker/middleware"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(corss.CORSMiddleware())

	api := r.Group("/api/v1")

	categories := api.Group("/categories")
	categories.GET("", h.GetAllCategory)
	categories.POST("", h.CreateCategory)
	categories.PUT("/:id", h.UpdateCategory)
	categories.DELETE("/:id", h.DeleteCategory)

	transactions := api.Group("/transactions")
	transactions.GET("", h.GetAllTransactions)
	transactions.POST("", h.CreateTransaction)
	transactions.PUT("/:id", h.UpdateTransaction)
	transactions.DELETE("/:id", h.DeleteTransaction)

	api.GET("/stats", h.GetStats)

	return r
}