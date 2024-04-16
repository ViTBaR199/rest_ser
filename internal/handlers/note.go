package handlers

import (
	"myapp/internal/models"
	"myapp/internal/services"
	"net/http"

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
	var noteDetails struct {
		Id int `json:"id"`
	}
	if err := c.ShouldBindJSON(&noteDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.NoteService.DeleteNote(noteDetails.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfuly"})
}

func (h *NoteHandlers) FetchNote(c *gin.Context) {
	var noteDetails struct {
		Start int `json:"start"`
		End   int `json:"end"`
	}
	if err := c.ShouldBindJSON(&noteDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rows, err := h.NoteService.FetchNote(noteDetails.Start, noteDetails.End)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": rows})
}
