package handlers

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"net/http"

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
	var taskDetails struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&taskDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TaskService.DeleteTask(taskDetails.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfuly"})
}

func (h *TaskHandler) FetchTask(c *gin.Context) {
	var taskDetails struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
	if err := c.ShouldBindJSON(&taskDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := h.TaskService.FetchTask(taskDetails.Start, taskDetails.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": rows})
}
