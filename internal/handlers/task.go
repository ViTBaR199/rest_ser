package handlers

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	TaskService *services.TaskServices
}

func NewTaskHandler(taskHandler *services.TaskServices) *TaskHandler {
	return &TaskHandler{TaskService: taskHandler}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TaskService.CreateTask(c, task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully"})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	ToDel := c.Query("to_del")
	if ToDel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no folder id provided"})
		return
	}

	taskID, err := strconv.Atoi(ToDel) // Преобразование строки в int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id format"})
		return
	}
	if err := h.TaskService.DeleteTask(taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfuly"})
}

func (h *TaskHandler) FetchTask(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both start and end parameters are required"})
		return
	}

	startINT, err := strconv.Atoi(start)
	endINT, err1 := strconv.Atoi(end)
	if err != nil || err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id format"})
		return
	}

	rows, err := h.TaskService.FetchTask(startINT, endINT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rows)
}
