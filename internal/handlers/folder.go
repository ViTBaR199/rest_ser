package handlers

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"net/http"

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
	var folderDetails struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&folderDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.FolderService.DeleteFolder(folderDetails.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfuly"})
}

func (h *FolderHandler) FetchFolder(c *gin.Context) {
	var folderDetails struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
	if err := c.ShouldBindJSON(&folderDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := h.FolderService.FetchFolder(folderDetails.Start, folderDetails.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": rows})
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
