package handlers

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FolderHandler struct {
	FolderService *services.FolderService
}

func NewFolderHandler(folderService *services.FolderService) *FolderHandler {
	return &FolderHandler{FolderService: folderService}
}

func HomePage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"title": "Home page"})
}

func (h *FolderHandler) CreateFolder(c *gin.Context) {
	var folder models.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.FolderService.CreateFolder(c, folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Folder created successfully"})
}

func (h *FolderHandler) DeleteFolder(c *gin.Context) {
	ToDel := c.Query("to_del")
	if ToDel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no folder id provided"})
		return
	}

	folderID, err := strconv.Atoi(ToDel) // Преобразование строки в int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id format"})
		return
	}

	if err := h.FolderService.DeleteFolder(folderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfuly"})
}

func (h *FolderHandler) FetchFolder(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	folderType := c.Query("folder_type")
	userId := c.Query("user")
	if folderType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no folder type provided"})
		return
	}

	if start == "" || end == "" || userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both start and end parameters are required"})
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

	userINT, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return
	}

	rows, err := h.FolderService.FetchFolder(startINT, endINT, userINT, folderType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rows == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no information was found about such a user or category"})
		return
	}

	c.JSON(http.StatusOK, rows)
}

func (h *FolderHandler) UpdateFolder(c *gin.Context) {
	var folder models.Folder
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.FolderService.UpdateFolder(c, folder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder updated successfully"})
}
