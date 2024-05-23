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
		return
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

	rows, err := h.NoteService.FetchNote(userINT, startINT, endINT, folderINTs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rows == nil {
		c.JSON(http.StatusOK, []models.Note{})
		return
	}

	c.JSON(http.StatusOK, rows)
}

func (h *NoteHandlers) UpdateNote(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.NoteService.UpdateNote(c, note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note update successfully"})
}

func (h *NoteHandlers) FetchNoteById(c *gin.Context) {
	note := c.Query("note_id")

	if note == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "both start and end parameters are required"})
		return
	}

	noteINT, err := strconv.Atoi(note)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start format"})
		return
	}

	rows, err := h.NoteService.FetchNoteById(noteINT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rows == (models.Note{}) {
		c.JSON(http.StatusOK, models.Folder{})
		return
	}

	c.JSON(http.StatusOK, rows)
}
