package handlers

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	FinanceService *services.FinanceService
}

func NewFinanceHandler(financeService *services.FinanceService) *FinanceHandler {
	return &FinanceHandler{FinanceService: financeService}
}

func (h *FinanceHandler) CreateFinance(c *gin.Context) {
	var finance models.Finance
	if err := c.ShouldBindJSON(&finance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.FinanceService.CreateFinance(c, finance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Finance created successfully"})
}

func (h *FinanceHandler) DeleteFinance(c *gin.Context) {
	var financeDetails struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&financeDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.FinanceService.DeleteFinance(financeDetails.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Finance deleted successfuly"})
}

func (h *FinanceHandler) FetchFinance(c *gin.Context) {
	var financeDetails struct {
		Start int `json:"start"`
		End   int `json:"end"`
		Month int `json:"month"`
	}
	if err := c.ShouldBindJSON(&financeDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := h.FinanceService.FetchFinance(financeDetails.Start, financeDetails.End, financeDetails.Month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": rows})
}
