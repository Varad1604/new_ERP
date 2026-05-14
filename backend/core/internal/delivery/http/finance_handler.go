package http

import (
	"net/http"

	"github.com/enterprise-erp/core/internal/usecase"
	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	financeUseCase *usecase.FinanceUseCase
}

func NewFinanceHandler(financeUseCase *usecase.FinanceUseCase) *FinanceHandler {
	return &FinanceHandler{financeUseCase: financeUseCase}
}

func (h *FinanceHandler) RegisterRoutes(router *gin.RouterGroup) {
	financeGroup := router.Group("/finance")
	{
		financeGroup.POST("/accounts", h.CreateAccount)
		financeGroup.GET("/accounts", h.GetAccounts)
		financeGroup.POST("/journals", h.PostJournalEntry)
	}

	dashboardGroup := router.Group("/dashboard")
	{
		dashboardGroup.GET("/metrics", h.GetDashboardMetrics)
	}
}

func (h *FinanceHandler) GetDashboardMetrics(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not authenticated"})
		return
	}

	metrics, err := h.financeUseCase.GetDashboardMetrics(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve dashboard metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (h *FinanceHandler) CreateAccount(c *gin.Context) {
	var req usecase.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := h.financeUseCase.CreateAccount(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, acc)
}

func (h *FinanceHandler) GetAccounts(c *gin.Context) {
	accounts, err := h.financeUseCase.GetAccounts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve accounts"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *FinanceHandler) PostJournalEntry(c *gin.Context) {
	var req usecase.PostJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	entry, err := h.financeUseCase.PostJournalEntry(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}
