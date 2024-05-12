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
		return
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
	user_id := c.Query("user_id")
	start := c.Query("start")
	end := c.Query("end")
	folder_id := c.Query("folder_id")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both start and end parameters are required"})
		return
	}

	userINT, err := strconv.Atoi(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	startINT, err := strconv.Atoi(start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start format"})
		return
	}

	endINT, err := strconv.Atoi(end)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end format"})
		return
	}

	var folderINTs []int
	if folder_id != "" {
		folderINT, err := strconv.Atoi(folder_id)
		if err != nil {

		}
		folderINTs = append(folderINTs, folderINT)
	}

	rows, err := h.TaskService.FetchTask(userINT, startINT, endINT, folderINTs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rows == nil {
		c.JSON(http.StatusOK, []models.Task{})
		return
	}

	c.JSON(http.StatusOK, rows)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TaskService.UpdateTask(c, task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (h *TaskHandler) CountTask(c *gin.Context) {
	user_id := c.Query("user_id")
	userINT, err := strconv.Atoi(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	rows, err := h.TaskService.CountTask(userINT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *TaskHandler) CountTaskFavourites(c *gin.Context) {
	user_id := c.Query("user_id")
	userINT, err := strconv.Atoi(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	rows, err := h.TaskService.CountTaskFavourites(userINT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *TaskHandler) FetchTaskFavourites(c *gin.Context) {
	user_id := c.Query("user_id")
	start := c.Query("start")
	end := c.Query("end")
	folder_id := c.Query("folder_id")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both start and end parameters are required"})
		return
	}

	userINT, err := strconv.Atoi(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	startINT, err := strconv.Atoi(start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start format"})
		return
	}

	endINT, err := strconv.Atoi(end)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end format"})
		return
	}

	var folderINTs []int
	if folder_id != "" {
		folderINT, err := strconv.Atoi(folder_id)
		if err != nil {

		}
		folderINTs = append(folderINTs, folderINT)
	}

	rows, err := h.TaskService.FetchTaskFavourites(userINT, startINT, endINT, folderINTs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rows == nil {
		c.JSON(http.StatusNotFound, []models.Task{})
		return
	}

	c.JSON(http.StatusOK, rows)
}
