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

	if err := h.FinanceService.DeleteFinance(c, financeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Finance deleted successfuly"})
}

func (h *FinanceHandler) FetchFinance(c *gin.Context) {
	user_id := c.Query("user_id")
	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both start, end and month parameters are required"})
		return
	}

	userINT, err := strconv.Atoi(user_id)
	startINT, err1 := strconv.Atoi(start)
	endINT, err2 := strconv.Atoi(end)
	if err != nil || err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id format"})
		return
	}

	rows, err := h.FinanceService.FetchFinance(userINT, startINT, endINT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rows == nil {
		c.JSON(http.StatusOK, []models.Finance{})
		return
	}

	c.JSON(http.StatusOK, rows)
}

func (h *FinanceHandler) FetchFinanceIncome(c *gin.Context) {
	user_id := c.Query("user_id")
	start := c.Query("start")
	end := c.Query("end")
	yearMonth := c.Query("yearMonth")
	if start == "" || end == "" || yearMonth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both start, end and month parameters are required"})
		return
	}

	userINT, err3 := strconv.Atoi(user_id)
	startINT, err := strconv.Atoi(start)
	endINT, err1 := strconv.Atoi(end)
	if err != nil || err1 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id format"})
		return
	}

	rows, err := h.FinanceService.FetchFinanceIncome(userINT, startINT, endINT, yearMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rows == nil {
		c.JSON(http.StatusOK, []models.Finance{})
		return
	}

	c.JSON(http.StatusOK, rows)
}

func (h *FinanceHandler) FetchFinanceExpense(c *gin.Context) {
	user_id := c.Query("user_id")
	start := c.Query("start")
	end := c.Query("end")
	yearMonth := c.Query("yearMonth")
	if start == "" || end == "" || yearMonth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both start, end and month parameters are required"})
		return
	}

	userINT, err3 := strconv.Atoi(user_id)
	startINT, err := strconv.Atoi(start)
	endINT, err1 := strconv.Atoi(end)
	if err != nil || err1 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id format"})
		return
	}

	rows, err := h.FinanceService.FetchFinanceExpense(userINT, startINT, endINT, yearMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rows == nil {
		c.JSON(http.StatusOK, []models.Finance{})
		return
	}

	c.JSON(http.StatusOK, rows)
}
