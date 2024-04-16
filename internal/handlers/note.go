package handlers

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoteHandlers struct {
	NoteService *services.NoteService
}

func NewNoteHandlers(noteService *services.NoteService) *NoteHandlers {
	return &NoteHandlers{NoteService: noteService}
}

func (h *NoteHandlers) CreateNote(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.NoteService.CreateNote(c, note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Note created successfully"})
}

func (h *NoteHandlers) DeleteNote(c *gin.Context) {
	ToDel := c.Query("to_del")
	if ToDel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no folder id provided"})
		return
	}

	noteID, err := strconv.Atoi(ToDel) // Преобразование строки в int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder id format"})
		return
	}

	if err := h.NoteService.DeleteNote(noteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfuly"})
}

func (h *NoteHandlers) FetchNote(c *gin.Context) {
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

	rows, err := h.NoteService.FetchNote(startINT, endINT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rows)
}
