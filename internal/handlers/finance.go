package handlers

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"net/http"
	"strconv"

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
	ToDel := c.Query("to_del")
	if ToDel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no folder id provided"})
		return
	}

	financeID, err := strconv.Atoi(ToDel) // Преобразование строки в int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id format"})
		return
	}

	if err := h.FinanceService.DeleteFinance(financeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Finance deleted successfuly"})
}

func (h *FinanceHandler) FetchFinance(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	month := c.Query("month")
	if start == "" || end == "" || month == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both start, end and month parameters are required"})
		return
	}

	startINT, err := strconv.Atoi(start)
	endINT, err1 := strconv.Atoi(end)
	monthINT, err2 := strconv.Atoi(month)
	if err != nil || err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id format"})
		return
	}

	rows, err := h.FinanceService.FetchFinance(startINT, endINT, monthINT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rows)
}
